# ADR 0016-queue-e-storage — `queue` e `storage` como décimo e décimo primeiro repositórios

**Status:** Aceita. Consome a autorização que a REGRA 3 exige e fecha a lista de
`11-adapters.md` para a Fase 3.

## Contexto

A REGRA 3 do `CLAUDE.md` diz:

> Dez repositórios até a Fase 3. Não crie o décimo primeiro sem pedir.

Antes deste trabalho existiam **nove**: `framework`, `aru`, `arandu`, `docs`,
`examples`, `porang`, `kv`, `database` e `homebrew-tap`. A ADR 0015 tinha
reservado a décima vaga para `otel`, que ficou adiado com gatilho.

E o `11-adapters.md:127` diz, sem ambiguidade:

> Fase 3: `queue` (sobre `kv`, depois `nats`), `storage` (local, depois `s3`).

Os dois faltavam. Pior: o `16-agendamento.md` descreve `aru work` — um worker de
fila — como parte do agendamento, e o `03-roadmap-fases.md` lista **jobs** no
console da Fase 3. Sem fila, três documentos descreviam algo que não existia.

## Decisão

Criar `arandu-io/queue` e `arandu-io/storage`, chegando a **onze repositórios**.

O proprietário autorizou explicitamente, respondendo *"incluir tudo e seguir
para a próxima fase"* à pergunta sobre criar os dois. É essa autorização que a
REGRA 3 pede, e ela fica registrada aqui em vez de na memória de uma conversa.

## O formato, herdado da ADR 0014

Cada um repete a divisão que o `database` já provou:

```
arandu-io/queue/            fila sobre a tabela da aplicação — zero dependência
└── kv/                     driver RESP, go.mod próprio

arandu-io/storage/          disco local — zero dependência
└── s3/                     protocolo S3, go.mod próprio
```

O contrato fica no core — `framework/jobs` e `framework/storage` — porque é ali
que o `Grant` é obrigatório. Mover isso para o adapter tornaria a garantia
opcional, e garantia opcional não é garantia.

## Duas decisões dentro dos dois

**A fila padrão é a tabela, não o Redis.** Um job empurrado dentro de
`data.Transaction` é comitado pela mesma transação da linha que o produziu — é a
garantia do outbox aplicada a trabalho em vez de a evento. Redis não oferece
isso, e por isso é o segundo driver, não o primeiro.

**O storage padrão é R2.** Ver `docs/11-adapters.md`: R2 não cobra egresso, e
num SaaS que devolve arquivos para quem os enviou o egresso é a maior parte da
conta em todo o resto — a linha que cresce junto com o sucesso. O pacote se
chama `s3` porque é o nome do **protocolo**, que cobre R2, S3, Spaces, B2 e
MinIO; o produto recomendado aparece em `s3.R2(...)` e no README, exatamente
como o Dragonfly aparece no `kv`.

**Sem SDK em nenhum dos dois.** SigV4 são duzentas linhas contra um SDK da AWS
que traz cem módulos, cadeia de credencial própria, política de retry própria e
regras de contexto próprias — e o algoritmo não muda desde 2012, enquanto a
superfície do SDK muda a cada trimestre. O CI reprova se qualquer SDK de nuvem
entrar no grafo.

## Consequências

- **Fecha:** onze repositórios até a Fase 4. O décimo segundo volta a exigir
  pedido, e a REGRA 3 continua valendo com o número atualizado
- `otel` perde a vaga reservada pela ADR 0015 e passa a ser o décimo segundo,
  quando o gatilho dela disparar. O adiamento não muda
- `mail` e `ws` (doc 11) continuam sem gatilho: "com demanda"
- O `aru work` do doc 16 passa a existir de verdade, e o console da Fase 3 passa
  a mostrar job com timeline, que era o que faltava para a linha do roadmap ser
  verdade
