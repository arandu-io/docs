# Agendamento de tarefas

- [Deixe o agendador decidir quando executar, não como repetir](#let-the-scheduler-decide-when-not-how-to-retry)
- [Declare o escopo antes de escrever o corpo da tarefa](#declare-scope-before-writing-the-task-body)
- [Agende um job para cada tenant](#schedule-one-job-for-every-tenant)
- [Inspecione e execute a agenda registrada](#inspect-and-invoke-the-registered-schedule)
- [Trate a ausência de lock como defeito de implantação](#treat-a-missing-lock-as-a-deployment-defect)
- [Tarefas Global não podem acessar dados de tenant](#global-tasks-cannot-touch-tenant-data)
- [Teste agendamento, fila e worker de ponta a ponta](#run-the-schedule-queue-and-worker-end-to-end)
- [Repare a dependência e execute a tarefa uma vez](#repair-the-dependency-and-run-the-task-once)
- [Mantenha o trabalho durável na fila](#keep-durable-work-in-the-queue)

<a name="let-the-scheduler-decide-when-not-how-to-retry"></a>
## Deixe o agendador decidir quando executar, não como repetir

Uma tarefa agendada é uma pequena decisão executada por uma expressão cron de cinco campos. Normalmente ela descobre trabalho e enfileira um job durável; o agendador não repete uma execução Run que falhou.

<a name="declare-scope-before-writing-the-task-body"></a>
## Declare o escopo antes de escrever o corpo da tarefa

Uma `foundation.Task` tem ID único, Spec válida, Action e Run. Global recebe um Grant zero e não alcança repositórios; `PerTenant` recebe `SystemGrant(Action, tenant)` para cada tenant que a função entregue a `WithTenants` responde. Defina Timeout, porque zero vira cinco minutos.

- Defina um identificador único para a tarefa.
- Escolha Global ou `PerTenant` explicitamente.
- Responda os tenants com `WithTenants` para `PerTenant`.
- As reivindicações singleton só são compartilhadas entre réplicas quando `drivers.Locks` é.

<a name="schedule-one-job-for-every-tenant"></a>
## Agende um job para cada tenant

Liste as tarefas da aplicação em `routes/console.go` e entregue-as a `WithSchedule`; um módulo declara as próprias por Schedule. Create registra o agendador por último, sobre as tarefas de todos os módulos e as da aplicação, e ele roda no processo que serve. Mantenha o trabalho durável num handler de fila: o Grant do agendador vira a procedência do job.

```go example
type CloseCycle struct{ Cycle string }

func Tasks(d *foundation.Drivers) []foundation.Task {
	return []foundation.Task{{
		ID: "billing.close", Spec: "0 3 * * *", Scope: foundation.PerTenant,
		Action: "invoice.send", Timeout: 30 * time.Second, Singleton: true,
		Run: func(ctx context.Context, g auth.Grant) error {
			return queue.Dispatch(ctx, d.Queue, g, "invoice.send", CloseCycle{Cycle: "current"})
		},
	}}
}
```

<a name="inspect-and-invoke-the-registered-schedule"></a>
## Inspecione e execute a agenda registrada

`aru schedule:list` informa id, cron, escopo, singleton, próxima execução, última execução e último erro. `schedule:run` usa o mesmo caminho de Grant, lock e instrumentação; uma execução manual por tenant exige `--tenant`.

- Inspecione o registro antes de esperar pelo cron.
- Execute uma tarefa manualmente para um tenant selecionado.
- Enfileire trabalho lento ou sujeito a novas tentativas.
- Prepare o corpo da tarefa para executar duas vezes com segurança.

<a name="treat-a-missing-lock-as-a-deployment-defect"></a>
## Trate a ausência de lock como defeito de implantação

Os processos que servem verificam perto de cada minuto. Sem um store compartilhado, `drivers.Locks` é nil e toda réplica pode executar um singleton. As chaves de lock cobrem tarefa, tenant e minuto durante o timeout; uma tarefa por tenant sem `WithTenants` não roda, e o agendador relata isso em vez de pular em silêncio.

> [!WARNING]
> O agendador não tem orçamento de tentativas. Uma tarefa que executa um efeito externo diretamente pode perder ou repetir trabalho; enfileire um job durável e idempotente.

<a name="global-tasks-cannot-touch-tenant-data"></a>
## Tarefas Global não podem acessar dados de tenant

> [!WARNING]
> Nunca invente um tenant dentro de Global nem o obtenha de HTTP, payload ou atalho de ambiente. O acesso a repositórios recebe apenas o Grant por tenant criado pelo resolvedor confiável e pela Action declarada.

<a name="run-the-schedule-queue-and-worker-end-to-end"></a>
## Teste agendamento, fila e worker de ponta a ponta

Construa o agendador, chame `RunNow`, confirme um job persistido, inicie um worker e confirme que o handler recebe tenant e payload. IDs duplicados, cron inválido e tarefa sem Run falham antes de servir. A ausência do resolvedor de tenants aparece na execução e nos diagnósticos; a falta de lock compartilhado deve ser detectada por um teste de configuração de implantação, pois o singleton seria executado em toda réplica.

- `RunNow` enfileira uma vez.
- O worker recebe o Grant esperado.
- Agendamentos inválidos impedem a inicialização.
- Apenas quem obtém o lock singleton executa.

<a name="repair-the-dependency-and-run-the-task-once"></a>
## Repare a dependência e execute a tarefa uma vez

Use `schedule:list` para distinguir ausência de registro de falha na última execução. Restaure banco de dados, fila, resolvedor ou lock e invoque um tenant afetado. O scheduler não repete tentativas automaticamente, mas a recuperação pode executar o trabalho de novo; torne seguros de repetir tanto a tarefa quanto qualquer handler enfileirado.

<a name="keep-durable-work-in-the-queue"></a>
## Mantenha o trabalho durável na fila

Filas cuidam da execução durável, a outbox cuida dos fatos após o commit, as operações de filas cuidam da recuperação e a observabilidade informa falhas e atrasos do agendamento.

- [Filas](/docs/{{version}}/async/queues)
- [Eventos e outbox](/docs/{{version}}/async/events-outbox)
- [Operações de fila](/docs/{{version}}/operations/queues)
- [Observabilidade](/docs/{{version}}/operations/observability)
