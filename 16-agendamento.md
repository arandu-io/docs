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
        ID:        "billing.close-cycle",
        Spec:      "0 3 * * *",        // cron de 5 campos
        Scope:     kernel.PerTenant,   // ou kernel.Global
        Timeout:   10 * time.Minute,
        Singleton: true,
        Action:    ActionCycleClose,   // vira o SystemGrant
        Run:       m.closeCycle,
    }}
}
```

Em inglês porque é código, e código é inglês (REGRA 6, ADR 0005).

`Scope: PerTenant` faz o scheduler expandir a tarefa para cada tenant ativo,
cada uma com seu `SystemGrant` e seu lock. `Global` roda uma vez — e recebe o
`Grant` zero, porque `SystemGrant` recusa tenant vazio (REGRA 14). Consequência
registrada: **tarefa Global não alcança repositório nenhum.** Isso é restrição,
não esquecimento: trabalho global é limpar arquivo temporário e aquecer cache;
o que lê linha de cliente é `PerTenant`, e ter que dizer isso é o ponto.

A lista de tenants ativos é injetada em `scheduler.Options.Tenants`, porque o
core não sabe onde a aplicação guarda a dela.

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

## Workers de fila — **feito**

Mesmo binário, subcomando: `aru queue:work --queue=default --workers=4`. Não é um
segundo artefato, é a mesma imagem com outro argumento — o que mantém o deploy
com um artefato só (`17-deploy.md`), e impede que o worker rode uma versão
diferente do código que o servidor.

A fila é `arandu-io/queue` (ADR 0016). O driver padrão é a tabela da própria
aplicação, porque um job empurrado dentro de `data.Transaction` é comitado pela
mesma transação da linha que o produziu.

## Comandos — **feitos**

- `aru schedule:list` — o que está registrado, com próxima execução
- `aru schedule:run <id> [--tenant=…]` — execução manual, mesmo lock, mesmo
  `Grant`, mesma instrumentação. Rota diferente seria porta dos fundos, e as
  duas divergiriam
- `aru queue:work [--queue=…] [--workers=…]` — worker de fila

Os três são o mesmo binário com outro argumento, o que mantém o deploy em um
artefato só (`17-deploy.md`).

## O que o scheduler não faz

**Retentativa.** Uma tarefa que falha é registrada no log e no
`kernel.Diagnostic`, e a próxima janela roda de novo. Trabalho que precisa de
orçamento de retry enfileira um job, e a fila cuida — `arandu-io/queue`, com
backoff exponencial e fila morta. Scheduler dispara, fila persiste.

**Cron de segundos.** Cinco campos, não seis. Um framework que oferece cron de
segundo oferece um jeito de escrever busy loop por acidente, e trabalho que
precisa acontecer a cada poucos segundos é worker, não agendamento.

## O que não construímos

Equivalente do Horizon (dashboard de filas) não entra até a Fase 3 — o console
`/_arandu/debug` já mostra job enfileirado e executado. **Gatilho:** cliente real
operando mais de três filas com prioridade distinta.
