# Eventos e outbox transacional

- [Commite o estado de negócio e seus fatos juntos](#commit-business-state-and-its-facts-together)
- [Conecte tabela, relay, publicador e lock entre réplicas](#wire-the-table-relay-publisher-and-replica-lock)
- [Armazene o evento na transação da escrita](#store-the-event-inside-the-write-transaction)
- [Processe pelo caminho de publicação de produção](#drain-through-the-production-publisher-path)
- [Prepare todo consumidor para entregas repetidas](#design-every-consumer-for-repeated-delivery)
- [Tenant e autoridade vêm do Grant](#tenant-and-authority-come-from-the-grant)
- [Comprove rollback, replay e comportamento entre réplicas](#prove-rollback-replay-and-replica-behavior)
- [Continue com workers de fila e agendamento](#continue-with-queue-workers-and-scheduling)

<a name="commit-business-state-and-its-facts-together"></a>
## Commite o estado de negócio e seus fatos juntos

Use a outbox durável para um fato que precisa sair do processo após uma escrita no banco. `events.Outbox.Store` recebe o mesmo contexto e Grant, registra tenant, action e autorizador e rejeita chamadas fora de `data.Transaction` com `ErrNoTransaction`. O despacho dentro do processo não fecha a janela de falha entre commit e publicação.

<a name="wire-the-table-relay-publisher-and-replica-lock"></a>
## Conecte tabela, relay, publicador e lock entre réplicas

Construa o outbox com `events.NewOutbox` sobre `drivers.DB` num provider, entregue a `events.NewRelay` uma lista `events.Listeners` cujas entradas `events.On` constrói, uma por nome de evento, e liste `events.WithRelay(relay)` em `WithProviders` em vez de apenas o `NewModule` dono da tabela. Várias réplicas exigem um lock distribuído: passe `drivers.Locker`, que só é nil quando todo store fica dentro do processo, e isso só serve para um único processo servidor.

- Migre a tabela da outbox.
- Injete um único Publisher de produção.
- Registre o ciclo de vida do relay.
- Use um lock compartilhado entre as réplicas que atendem requisições.

<a name="store-the-event-inside-the-write-transaction"></a>
## Armazene o evento na transação da escrita

Escreva a alteração de negócio e armazene seu fato imutável com Store na mesma transação SQL. Se uma das operações falhar, nenhuma é commitada. Consumidores podem recarregar pelo id, mas não devem presumir que o estado mutável ainda corresponda ao evento ocorrido.

```go example
type Invoices interface {
	MarkPaid(ctx context.Context, grant security.Grant, id string) error
}

func MarkPaid(ctx context.Context, db *data.DB, invoices Invoices, outbox *events.Outbox, grant security.Grant, invoiceID string) error {
	return data.Transaction(ctx, db, func(ctx context.Context) error {
		if err := invoices.MarkPaid(ctx, grant, invoiceID); err != nil {
			return err
		}
		return outbox.Store(ctx, grant, []events.Event{{
			Name: "invoice.paid", Aggregate: "invoice", AggregateID: invoiceID,
			Payload: map[string]any{"invoice_id": invoiceID}, OccurredAt: time.Now().UTC(),
		}})
	})
}
```

<a name="drain-through-the-production-publisher-path"></a>
## Processe pelo caminho de publicação de produção

`Relay.Drain` processa as linhas pendentes pelo mesmo publisher do loop que serve. Inspecione as operações de pendentes, estacionados, nova tentativa, atraso e marcação; o padrão é intervalo de um segundo, lote de 100, dez tentativas e um lock de relay de trinta segundos, salvo configuração. Para um consumidor HTTP externo, componha `webhook.Manager` com o banco de dados, o resolvedor de segredos e opções limitadas, e adapte o Publisher do outbox existente em vez de criar um segundo outbox. Registre o módulo de webhook do Hesape com `foundation.Component` em `WithProviders`: Create coleta seus handlers de jobs e as isenções de CSRF dos seus receptores, então a aplicação não escreve nem uma chamada a `RegisterJobHandlers` nem uma linha `ValidateCSRFTokensExcept` para eles, e chamar `RegisterJobHandlers` também registra cada job duas vezes e interrompe o boot. Seu prefixo de gestão e seu prefixo de recepção não podem se sobrepor, porque o prefixo de recepção é isento da verificação de CSRF.

- Processe os registros pelo Publisher de produção.
- Monitore eventos pendentes e atraso.
- Estacione após esgotar as tentativas.
- Tente novamente somente após corrigir o consumidor ou o transporte.
- Registre cada endpoint de webhook com uma URL HTTP ou HTTPS absoluta e um segredo de pelo menos 32 bytes.

<a name="design-every-consumer-for-repeated-delivery"></a>
## Prepare todo consumidor para entregas repetidas

A entrega ocorre pelo menos uma vez. Uma publicação pode ter sucesso antes de a confirmação falhar, então o mesmo `Stored.ID` pode chegar novamente. Consumidores deduplicam por esse ID estável e tornam os efeitos idempotentes; linhas estacionadas permanecem como evidência. O despacho de webhook cria atomicamente um snapshot de entrega e um job de fila para cada evento e endpoint. Falhas de transporte, 408, 425, 429 e 5xx tentam novamente; respostas 4xx comuns são permanentes, e Retry-After é respeitado.

> [!WARNING]
> Uma publicação direta após o commit pode se perder permanentemente se o processo terminar entre as duas operações. Use a outbox sempre que esse fato precisar sair do processo.

<a name="tenant-and-authority-come-from-the-grant"></a>
## Tenant e autoridade vêm do Grant

> [!WARNING]
> Não coloque em um evento um tenant escolhido pelo chamador para depois confiar nele. Store persiste tenant, action e `authorized_by` a partir do Grant; a entrega de webhook herda tenant e action desse Grant. Requisições POST assinadas carregam X-Arandu-Signature, X-Arandu-Timestamp e X-Arandu-Delivery-ID; HMAC-SHA256 cobre timestamp, ID da entrega e corpo. Mantenha a resolução do endpoint atrás do cliente protegido contra SSRF e imponha o limite de payload de um mebibyte, limites de resposta e redirecionamento, timeout e rotação de segredo.

Payloads contêm identificadores e fatos públicos imutáveis, não credenciais, Grants brutos nem registros completos com segredos. Uma assinatura válida de webhook comprova a integridade da mensagem, não permissão para ler dados de outro tenant.

<a name="prove-rollback-replay-and-replica-behavior"></a>
## Comprove rollback, replay e comportamento entre réplicas

Retorne um erro sentinela de `data.Transaction` e comprove que nenhuma linha de negócio nem de outbox permanece. Depois faça a publicação ter sucesso e a confirmação falhar, drene duas vezes e comprove deduplicação do consumidor e bloqueio entre réplicas. Para webhooks, verifique a assinatura byte a byte, rotação do segredo anterior, um job por evento e endpoint, idempotência de despacho duplicado, recusa de SSRF, limites de resposta, classificação de tentativas, claims cercados, status de auditoria e limpeza por retenção.

- O rollback não deixa registro na outbox.
- O replay mantém um único efeito externo.
- Apenas quem obtém o lock processa os eventos.
- O atraso de eventos retidos aparece nos diagnósticos.
- Um ID de entrega de webhook duplicado nunca cria um segundo efeito externo.

<a name="continue-with-queue-workers-and-scheduling"></a>
## Continue com workers de fila e agendamento

Filas executam tentativas duráveis de webhook, agendamento descobre trabalho, observabilidade expõe atraso do relay e autorização define a fronteira do Grant registrado. Webhooks de saída duráveis são separados dos callbacks opcionais do ciclo de vida Pusher no Joaju.

- [Padrões de comunicação](/docs/{{version}}/architecture/contracts-adapters)
- [Filas](/docs/{{version}}/async/queues)
- [`WebSocket` do Joaju](/docs/{{version}}/async/joaju)
- [Observabilidade](/docs/{{version}}/operations/observability)
- [Autorização](/docs/{{version}}/security/authorization)
