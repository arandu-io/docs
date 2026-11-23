# Filas

- [Persista trabalho que precisa sobreviver à requisição](#persist-work-that-must-survive-the-request)
- [Escolha um armazenamento cuja fronteira transacional você conheça](#choose-a-store-whose-transaction-boundary-you-understand)
- [Crie a partir de um Grant, envie à fila e processe pelo nome](#create-from-a-grant-push-and-handle-by-name)
- [Defina reserva e tentativas conforme o comportamento do handler](#set-lease-and-retry-policy-from-handler-behavior)
- [Saiba se o job foi liberado ou estacionado](#know-whether-a-job-was-released-or-parked)
- [Nunca reconstrua autoridade a partir do payload](#never-rebuild-authority-from-the-payload)
- [Verifique rollback e entrega duplicada](#prove-rollback-and-duplicate-delivery)
- [Inspecione, repare, tente novamente e observe o atraso](#inspect-repair-retry-and-observe-lag)
- [Separe o desenho do produtor das operações do worker](#separate-producer-design-from-worker-operations)

<a name="persist-work-that-must-survive-the-request"></a>
## Persista trabalho que precisa sobreviver à requisição

Uma fila separa a aceitação do trabalho de sua execução. Jobs da fila do Hesape carregam tenant, ação, autorizador, payload e UUID estável; o worker reconstrói o Grant antes de invocar o handler registrado. A entrega ocorre pelo menos uma vez, portanto a idempotência do handler é um contrato da aplicação.

<a name="choose-a-store-whose-transaction-boundary-you-understand"></a>
## Escolha um armazenamento cuja fronteira transacional você conheça

Create registra o próprio módulo da fila: a tabela de jobs, a saúde de atraso e os diagnósticos da fila que `QUEUE_CONNECTION` nomeia, guardada em `drivers.Queue`. Registrar `jobs.NewModule` de novo em `WithProviders` o registra duas vezes. A fila de banco de dados pode gravar sua linha na mesma transação SQL da escrita que a produz. RESP, selecionado por `QUEUE_CONNECTION=redis` depois que `hesape/queue/connectors/redis` é vinculado, é útil quando escolhido deliberadamente, mas não entra nessa transação.

- Escolha deliberadamente banco de dados ou RESP.
- Deixe o módulo da fila para Create.
- Mantenha servidor web e worker no mesmo artefato.
- Defina a idempotência antes de enfileirar.

<a name="create-from-a-grant-push-and-handle-by-name"></a>
## Crie a partir de um Grant, envie à fila e processe pelo nome

Registre o nome exato no `RegisterJobs` do provider, que `bootstrap/app.go` entrega a `WithJobs`, e despache um payload tipado sob o Grant autorizado com `queue.Dispatch`. `HandlerFor` decodifica o payload no mesmo tipo antes de o handler rodar. Dispatch constrói o job com `jobs.New` na fila padrão, então um tenant ou nome ausente é recusado, e Push verifica a autoridade mesmo para um job montado à mão para outra fila.

```go example
type SendInvoice struct{ InvoiceID string }

type InvoiceSender interface {
	Send(ctx context.Context, g auth.Grant, jobID, invoiceID string) error
}

type AppServiceProvider struct{ Invoices InvoiceSender }

func (p *AppServiceProvider) RegisterJobs(w *queue.Worker) error {
	w.Handle("invoice.send", queue.HandlerFor(func(ctx context.Context, g auth.Grant, j *jobs.Job, in SendInvoice) error {
		return p.Invoices.Send(ctx, g, j.UUID, in.InvoiceID)
	}))
	return nil
}

func Enqueue(ctx context.Context, drivers *foundation.Drivers, grant auth.Grant, invoiceID string) error {
	if err := queue.Dispatch(ctx, drivers.Queue, grant, "invoice.send", SendInvoice{InvoiceID: invoiceID}); err != nil {
		return err
	}
	return nil
}
```

<a name="set-lease-and-retry-policy-from-handler-behavior"></a>
## Defina reserva e tentativas conforme o comportamento do handler

`queue:work` constrói seu worker a partir de `Settings.Queue`: `QUEUE_DEFAULT` nomeia a fila, `QUEUE_WORKERS` a concorrência, `QUEUE_RETRY_AFTER` o lease em segundos e `QUEUE_MAX_ATTEMPTS` as tentativas, e `--queue` e `--workers` sobrescrevem os dois primeiros numa invocação. Uma variável não definida mantém o padrão do próprio componente de fila: quatro workers e um lease de cinco minutos, com backoff exponencial limitado entre as tentativas.

- Limite a concorrência conforme as dependências.
- Defina a reserva acima da duração máxima do handler.
- Defina explicitamente a política de tentativas e espera progressiva.
- Passe o Recorder da aplicação aos workers.

<a name="know-whether-a-job-was-released-or-parked"></a>
## Saiba se o job foi liberado ou estacionado

Um erro retornado consome uma tentativa e libera o job depois do backoff; trabalho esgotado é estacionado. `ErrManuallyFailed` estaciona uma falha permanente na hora, e um panic é recuperado, registrado com a pilha e estacionado na primeira entrega. O worker reporta um job de um batch ou de uma chain depois que ele roda, então um handler que ainda chama `bus.Handled` o conta duas vezes e empurra o próximo elo da chain duas vezes. Payloads com falha continuam sendo dados de cliente protegidos por tenant.

> [!WARNING]
> Uma reserva mais curta que a duração do handler permite entregas duplicadas simultâneas. Aumentar as tentativas não torna seguro um efeito externo sem idempotência.

<a name="never-rebuild-authority-from-the-payload"></a>
## Nunca reconstrua autoridade a partir do payload

> [!WARNING]
> O worker usa o tenant e a ação registrados para o seu Grant, e um job empurrado enquanto um operador agia na conta leva esse operador, de modo que o que ele escreve fica registrado como dele. Um push sob um Grant somente leitura é recusado com `auth.ErrReadOnly`. Nunca aceite `tenant_id`, action, role ou um Grant serializado vindos de JSON. Passe o Grant do worker para serviços e repositórios, e mantenha a inspeção de falhas com escopo de tenant.

<a name="prove-rollback-and-duplicate-delivery"></a>
## Verifique rollback e entrega duplicada

Para a fila de banco de dados, enfileire em uma transação que retorna um erro sentinela e comprove que `PendingSize` permanece zero. Entregue duas vezes o mesmo job estável e comprove um único efeito externo; teste também payload malformado, falha permanente, panic e reserva deliberadamente curta.

- O rollback remove o registro enfileirado.
- A entrega duplicada é idempotente.
- Falha permanente estaciona o trabalho.
- Um tenant não pode inspecionar trabalho falho de outro.

<a name="inspect-repair-retry-and-observe-lag"></a>
## Inspecione, repare, tente novamente e observe o atraso

Inspecione falhas e atraso, repare o handler ou a dependência, repita um job de teste e substitua os workers com encerramento gracioso. Forget, flush e clear destroem evidências ou trabalho e exigem o menor escopo de tenant revisado.

<a name="separate-producer-design-from-worker-operations"></a>
## Separe o desenho do produtor das operações do worker

Operações documenta comandos, agendamento cria jobs recorrentes, outbox persiste fatos e autorização fornece a fronteira de Grant. Uma fila é uma fronteira interna de execução durável; um webhook usa trabalho enfileirado para cruzar até um serviço HTTP externo.

- [Padrões de comunicação](/docs/{{version}}/architecture/contracts-adapters)
- [Operações de fila](/docs/{{version}}/operations/queues)
- [Agendamento](/docs/{{version}}/async/scheduling)
- [Eventos e webhooks](/docs/{{version}}/async/events-outbox)
- [Autorização](/docs/{{version}}/security/authorization)
