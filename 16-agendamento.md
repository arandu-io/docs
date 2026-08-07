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

Nada novo entra no core por causa disso, e o motivo não é o que este documento
dizia. *Corrigido em 07/08/2026:* esta linha afirmava que o `kv` já é dependência
do core. Não é: `framework/go.mod` tem um `require` só, `golang.org/x/crypto`
(ADR 0004). O contrato é a interface `kernel.Locker` — um método `Run(ctx, name,
ttl, fn)` —, quem a implementa é o adaptador, em `kv/lock.go`, com `SetNX`, e a
aplicação injeta a implementação em `scheduler.Options.Locker`. É justamente por
não importar o `kv` que o scheduler reconhece "outra réplica está com o lock"
pela mensagem do erro, na função `isLocked` de `framework/scheduler/scheduler.go`,
e não por um tipo sentinela.

## Falha

Log estruturado. O logger da tarefa carrega `component` e `task`, e no laço por
tenant também `tenant`; a linha de falha acrescenta `duration_ms`, `queries` e
**`error`**. A que traz `sql_ms` é a de sucesso — `task done` —, não a de falha.

O `Collector` instrumenta a execução como se fosse um request — a
timeline de uma tarefa aparece no console de debug igual à de um request. A
instrumentação é condicional: só acontece quando `scheduler.Options.Recorder`
foi passado em `Options`. `Recorder` nil é o padrão de produção — nenhum `Collector` é
construído, e é assim de propósito.

*Corrigido em 07/08/2026:* esta linha dizia "evento emitido" e nomeava as chaves
de log como `task_id` e `tenant_id`. Nenhuma das duas coisas é verdade. O
scheduler não emite evento: `framework/scheduler/` não importa
`framework/events`, e não há chamada a `Recorder.Record` de `events` em lugar
nenhum do pacote — o `Record` que o `execute` chama é o de `observability`, que
alimenta o console de debug, não um barramento. As chaves são `task` e `tenant`,
sem sufixo. O que existe além do log é o `kernel.Diagnostic` do
`scheduler.Module`, que reporta tarefa atrasada e tarefa que falhou na página de
erro, e é o que a seção "O que o scheduler não faz" já descrevia. Emitir evento
de tarefa continua possível dentro do `Run` da própria tarefa, como qualquer
serviço faz; o scheduler não faz por ela.

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
`/_arandu/debug` já lista cada job **executado** como entrada própria, ao lado
dos requests: `framework/jobs/worker.go` grava no `Recorder` um `Recorded` com
`Method: "job"` e o nome do job no `Path`. O **enfileiramento** não tem entrada
própria; ele aparece como o `INSERT INTO jobs` entre as queries do request que o
produziu — **e isso vale só para o driver de tabela.** O `queue/queue.go` grava
por `data.DB.ExecContext`, que passa pelo `RecordQuery`; o `queue/kv` enfileira
com `HSet` + `ZAdd` num pipeline RESP e **não importa `observability`**, então
com ele o enfileiramento não deixa rastro nenhum no console — nem entrada
própria, nem query. É a diferença mais visível entre os dois drivers, e ela não
está no contrato: é consequência de um passar pelo `data` e o outro não.

Os dois dependem de o `Recorder` existir, o que só acontece em dev ou
com `TracingSecret` configurado (`framework/kernel/kernel.go`). **Gatilho:**
cliente real operando mais de três filas com prioridade distinta.
