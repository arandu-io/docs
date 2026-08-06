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

- `LISTEN/NOTIFY` do Postgres para latência baixa, com polling de fallback
- lock distribuído em `kv` para não publicar duas vezes com N réplicas
- backoff exponencial por evento, `attempts` e `last_error` gravados
- depois de N falhas, o evento vai para fila morta e **aparece no `doctor`**

## Idempotência é obrigação do consumidor

Entrega é *at-least-once*. O mesmo evento pode chegar duas vezes, e isso não é
bug — é o preço da garantia.

O `id` do outbox é a chave de deduplicação, e o handler gerado pelo `aru` já
nasce com a checagem. Handler sem tratamento de idempotência é reprovado pelo
`doctor`.

## Observabilidade

Aqui o outbox encontra o cargo-chefe. O `Collector` mede:

- eventos gravados neste request, com o evento e o agregado
- **atraso do relay**: quanto tempo o evento mais antigo está pendente
- taxa de falha e tamanho da fila morta

E vira `hint` na página de erro, em inglês como toda mensagem do framework:
*"The oldest unpublished event has been waiting 4m0s. Is the relay running, and
is the publisher accepting?"*. Isso é diagnóstico. Métrica de dashboard ninguém
olha.

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
