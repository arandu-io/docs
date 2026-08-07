# 27 — Eventos de domínio e Outbox

O doc 13 coloca Events no core e nunca especificou consistência. Esta é a spec.
A decisão veio da revisão externa (doc 26) e é uma das poucas coisas que
absorvemos de lá inteiras.

## O problema

O fluxo ingênuo perde dado nos dois sentidos:

```go
repo.Save(ctx, customer)       // 1. write
bus.Publish(customer.Events()) // 2. publish
```

Se o processo morre entre 1 e 2, o evento nunca sai. Se o passo 2 acontece e a
transação faz rollback, o resto do sistema reage a algo que não existe.

## A regra

**Não existe `Publish()` direto.** Só `outbox.Store()`, dentro da mesma
transação da escrita. Um relay publica depois. Isso é a REGRA 9 aplicada a
eventos: uma forma só.

```go
return data.Transaction(ctx, s.db, func(ctx context.Context) error {
    if err := s.repo.Update(ctx, g, customer); err != nil {
        return err
    }
    return s.outbox.Store(ctx, g, customer.PullEvents())
})
```

Repare no `g`: o **`Grant` entra no registro do outbox**. Quem autorizou, qual
ação, qual tenant. Auditoria completa sem tabela extra.

## A tabela

Só o subconjunto portável de tipos: cada um se escreve igual em SQLite,
PostgreSQL e MySQL, que é o que deixa um schema só servir aos três (ADR 0009).
`jsonb` seria a grafia de um motor, e o payload é escrito e lido como JSON em
texto de qualquer jeito.

```sql
CREATE TABLE outbox (
    id            VARCHAR(255) PRIMARY KEY,
    tenant_id     VARCHAR(255) NOT NULL,
    event         TEXT NOT NULL,          -- "customer.created"
    aggregate     TEXT NOT NULL,
    aggregate_id  TEXT NOT NULL,
    payload       TEXT NOT NULL,
    authorized_by TEXT NOT NULL,          -- the Grant's subject
    action        TEXT NOT NULL,          -- the Grant's action
    occurred_at   TIMESTAMP NOT NULL,
    published_at  TIMESTAMP,
    attempts      INTEGER NOT NULL DEFAULT 0,
    last_error    TEXT
);

-- The relay reads unpublished events oldest first. A partial index would be
-- tighter, and MySQL does not have one; the two leading columns give the same
-- scan on every engine.
CREATE INDEX idx_outbox_pending ON outbox (published_at, occurred_at);
CREATE INDEX idx_outbox_tenant  ON outbox (tenant_id, occurred_at);
```

A fila morta entrou depois, em migração separada — a primeira já tinha rodado em
algum lugar. A coluna é nullable, então o binário anterior continua funcionando
durante o rollout: ele simplesmente nunca a escreve (REGRA 16).

```sql
ALTER TABLE outbox ADD COLUMN failed_at TIMESTAMP;

-- "Pending" now means neither published nor parked.
CREATE INDEX idx_outbox_unfinished ON outbox (failed_at, published_at, occurred_at);
```

## O relay

Faz parte do runtime, no mesmo binário (mesma lógica do scheduler, doc 16):

- **polling**, com intervalo de 1s por padrão (`RelayOptions.Interval`). O
  `LISTEN/NOTIFY` foi considerado e **recusado**: é específico do Postgres,
  colocaria um driver no core e daria a SQLite e MySQL um segundo caminho de
  código (REGRA 9 e REGRA 11). Um segundo de latência num publish de fundo não é
  o problema que este framework existe para resolver
- lock distribuído em `kv` para não publicar duas vezes com N réplicas
  (`RelayOptions.Locker`, `LockTTL` de 30s)
- backoff por evento, `attempts` e `last_error` gravados
- depois de `MaxAttempts` falhas (10 por padrão) o evento é **estacionado**:
  `failed_at` preenchido, e ele sai do conjunto pendente até alguém chamar
  `Outbox.Retry`

## Idempotência é obrigação do consumidor

Entrega é *at-least-once*. O mesmo evento pode chegar duas vezes, e isso não é
bug — é o preço da garantia.

**A contrapartida que o framework oferece é uma só: o `id` do registro do outbox
viaja com o evento** (`Stored.ID`), e é uma chave de deduplicação estável — o
mesmo evento reentregue chega com o mesmo `id`. O que o consumidor faz com ela é
com ele.

E é preciso ser explícito sobre o que **não** existe, porque já esteve escrito
aqui que existia:

- **não há handler gerado.** O `aru make:module` não emite listener nenhum;
  `app/Events/` e `app/Listeners/` saem vazios do `aru new`. Não há checagem de
  idempotência "que já nasce" em lugar algum
- **não há regra de `doctor` que reprove handler sem idempotência.** O contrato
  do lado de fora é `events.Publisher` — uma interface de **saída**. Quem
  consome está noutro processo, muitas vezes noutra linguagem, e o `doctor` lê o
  AST deste repositório: ele não tem o consumidor para analisar

Se um dia houver um contrato de consumidor de primeira parte, aí cabe regra. Até
lá, esta seção é um aviso, não uma garantia.

## Observabilidade

Aqui o outbox encontra o cargo-chefe, e em **três** lugares distintos — vale
separar, porque cada um roda em um momento diferente:

| Onde | O quê | Como |
|---|---|---|
| `Collector` | eventos gravados **neste request**, com nome e payload | `Collector.RecordEvent` / `Events()`, visíveis no console `/_arandu/debug` |
| Página de erro | atraso do relay e contagem da fila morta | `events.Module.Diagnose`, um `kernel.Diagnostic` — consultado quando a página renderiza |
| Health check | o relay parou | `events.Module.Health` falha quando o evento mais antigo passa de `maxLag`, 1 minuto |

O `Diagnose` é o que produz o `hint`, em inglês como toda mensagem do framework:

> *"The oldest unpublished event has been waiting 4m0s. Is the relay running, and
> is the publisher accepting?"*

> *"3 event(s) gave up after repeated failures, the most recent being
> invoice.paid: connection refused. They stay in the outbox until retried."*

Isso é diagnóstico. Métrica de dashboard ninguém olha.

**A fila morta não aparece no `aru doctor`, e não pode aparecer.** O `doctor` faz
análise estática: lê o AST do projeto e nunca executa código nem abre conexão de
banco. Contar eventos estacionados exige consultar a tabela `outbox` de um
ambiente — que é exatamente o que o `Diagnose` e o `Health` fazem, dentro do
processo que já tem o banco aberto. Quem quiser a lista em código chama
`Outbox.Parked`.

## Sem modo síncrono

Nem para teste. O teste usa o mesmo caminho, com o relay executado inline pelo
helper (`arandutest.DrainOutbox`). "Modo sync só em teste" é segunda forma
de fazer a mesma coisa, e a segunda forma sempre vaza para produção.

## Fases

| O quê | Fase |
|---|---|
| Tabela `outbox` e `outbox.Store()`, driver sobre `database` | **2** |
| Relay, fila morta, `kv` para lock, NATS | 3 |
| Painel de eventos no console `/_arandu/debug` | 3 |

O `Store()` sobe para a Fase 2 pelo mesmo motivo do `arandu.mod.toml`: evento
que não foi gravado não se recupera depois. O relay pode esperar; a gravação não.
