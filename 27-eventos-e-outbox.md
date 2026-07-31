# 27 — Eventos de domínio e Outbox

O doc 13 coloca Events no core e nunca especificou consistência. Esta é a spec.
A decisão veio da revisão externa (doc 26) e é uma das poucas coisas que
absorvemos de lá inteiras.

## O problema

O fluxo ingênuo perde dado nos dois sentidos:

```go
repo.Save(ctx, cliente)        // 1. grava
bus.Publish(cliente.Eventos()) // 2. publica
```

Se o processo morre entre 1 e 2, o evento nunca sai. Se o passo 2 acontece e a
transação faz rollback, o resto do sistema reage a algo que não existe.

## A regra

**Não existe `Publish()` direto.** Só `outbox.Store()`, dentro da mesma
transação da escrita. Um relay publica depois. Isso é a REGRA 9 aplicada a
eventos: uma forma só.

```go
return s.tx.Dentro(ctx, func(ctx context.Context) error {
    if err := s.repo.Update(ctx, g, cliente); err != nil {
        return err
    }
    return s.outbox.Store(ctx, g, cliente.PuxarEventos())
})
```

Repare no `g`: o **`Grant` entra no registro do outbox**. Quem autorizou, qual
ação, qual tenant. Auditoria completa sem tabela extra.

## A tabela

```sql
CREATE TABLE outbox (
    id           uuid PRIMARY KEY,
    tenant_id    uuid NOT NULL,
    evento       text NOT NULL,          -- "cliente.criado"
    agregado     text NOT NULL,
    agregado_id  text NOT NULL,
    payload      jsonb NOT NULL,
    autorizado_por text NOT NULL,        -- subject do Grant
    acao         text NOT NULL,          -- action do Grant
    ocorrido_em  timestamptz NOT NULL,
    publicado_em timestamptz,
    tentativas   int NOT NULL DEFAULT 0,
    ultimo_erro  text
);

CREATE INDEX ON outbox (ocorrido_em) WHERE publicado_em IS NULL;
```

## O relay

Faz parte do runtime, no mesmo binário (mesma lógica do scheduler, doc 16):

- `LISTEN/NOTIFY` do Postgres para latência baixa, com polling de fallback
- lock distribuído em `kv` para não publicar duas vezes com N réplicas
- backoff exponencial por evento, `tentativas` e `ultimo_erro` gravados
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

E vira `hint` na página de erro: *"evento `fatura.paga` pendente há 4 minutos —
o relay está rodando?"*. Isso é diagnóstico. Métrica de dashboard ninguém olha.

## Sem modo síncrono

Nem para teste. O teste usa o mesmo caminho, com o relay executado inline pelo
helper (`arandutest.DrenarOutbox(t)`). "Modo sync só em teste" é segunda forma
de fazer a mesma coisa, e a segunda forma sempre vaza para produção.

## Fases

| O quê | Fase |
|---|---|
| Tabela `outbox` e `outbox.Store()`, driver sobre `database` | **2** |
| Relay, fila morta, `kv` para lock, NATS | 3 |
| Painel de eventos no console `/_arandu/debug` | 3 |

O `Store()` sobe para a Fase 2 pelo mesmo motivo do `arandu.mod.toml`: evento
que não foi gravado não se recupera depois. O relay pode esperar; a gravação não.
