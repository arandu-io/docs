# 16 — Agendamento e workers

## Por que não copiamos o Laravel aqui

O `Illuminate/Console/Scheduling` existe porque PHP não tem processo residente:
um cron do sistema chama `schedule:run` a cada minuto e o PHP decide o que rodar.
São dois artefatos e uma dependência do sistema operacional.

Go tem processo residente. O scheduler é uma goroutine no mesmo binário.

## O contrato

Interface opcional do módulo, no mesmo modelo de `Migrations()`:

```go
func (m *Module) Schedule() []kernel.Task {
    return []kernel.Task{{
        ID:       "billing.fechar-ciclo",
        Spec:     "0 3 * * *",       // cron de 5 campos
        Escopo:   kernel.PorTenant,   // ou kernel.Global
        Timeout:  10 * time.Minute,
        Singleton: true,
        Acao:     ActionCicloFechar,  // vira SystemGrant
    }}
}
```

`Escopo: PorTenant` faz o scheduler expandir a tarefa para cada tenant ativo,
cada uma com seu `SystemGrant` e seu lock. `Global` roda uma vez.

## Réplica: o problema que o cron esconde

Com N réplicas do binário, N schedulers disparam a mesma tarefa. Resolvido com
lock distribuído no `kv` (`SET NX PX`), chave
`sched:<tenant>:<id>:<janela>` e TTL igual ao timeout da tarefa. `Singleton: false`
desliga o lock para tarefas idempotentes.

Isso é possível porque `kv` já é dependência do core. Não entra nada novo.

## Falha

Log estruturado com `task_id`, `tenant_id` e duração, evento emitido, e o
`Collector` instrumenta a execução como se fosse um request — a timeline de uma
tarefa aparece no console de debug igual à de um request.

Retentativa **não** é do scheduler: a tarefa enfileira um job e a fila cuida do
retry. Scheduler dispara, fila persiste.

## Workers de fila

Mesmo binário, subcomando: `aru work --fila=default --workers=4`. Não é um
segundo artefato, é a mesma imagem com outro argumento — o que mantém o deploy
com um artefato só (`17-deploy.md`).

## Comandos

- `aru schedule:list` — o que está registrado, com próxima execução
- `aru schedule:run <id> [--tenant=…]` — execução manual, mesma trilha de auditoria
- `aru work` — worker de fila

## O que não construímos

Equivalente do Horizon (dashboard de filas) não entra até a Fase 3 — o console
`/_arandu/debug` já mostra job enfileirado e executado. **Gatilho:** cliente real
operando mais de três filas com prioridade distinta.
