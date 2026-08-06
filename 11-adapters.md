# Adapters — pacotes por domínio

Seis pacotes no lugar de quinze. Um pacote por **domínio**, não por fornecedor.

O Laravel faz igual e é a referência aqui: `Illuminate\Mail` traz SMTP, Mailgun,
Postmark, SES e Resend no mesmo pacote; `Illuminate\Database` traz MySQL,
Postgres, SQLite e SQL Server juntos. Nunca um pacote por fornecedor.

## A diferença de Go que decide o formato

**Em Go não existe dependência opcional.** No Composer o driver do SES é
`suggest` — só baixa se você pedir. Em Go, se `mail` importar o SDK da AWS,
**todo mundo** que importar `mail` carrega o SDK no `go.sum`, no build e na
superfície de vulnerabilidade. Build tag não resolve: a dependência continua no
`go.mod`.

Solução: **pacote por domínio, driver pesado em submódulo.**

```
arandu-io/mail/                 um repositório, um conceito
├── go.mod       →  .../mail          interface + SMTP (stdlib, zero dep)
├── ses/go.mod   →  .../mail/ses      SDK da AWS só aqui
├── resend/go.mod
├── postmark/go.mod
└── mailgun/go.mod
```

Quem quer o básico: `go get github.com/arandu-io/mail`. Quem precisa de SES:
`go get github.com/arandu-io/mail/ses`. Um repositório, uma documentação, um
issue tracker, isolamento de dependência preservado.

Preço real: as tags viram `ses/v1.2.0` em vez de `v1.2.0`.

Isso é exceção à regra de um módulo por repositório (`05-repositorios.md`).
Matriz de drivers é o caso em que multi-módulo paga o próprio custo.

## Os seis domínios

| Pacote | Núcleo (sem dep pesada) | Submódulos |
|---|---|---|
| `database` | registro, pool e erro; zero driver | `pgx`, `sqlite`, `mysql` |
| `kv` | go-redis (RESP) | nenhum |
| `mail` | interface + SMTP (`net/smtp`) | `ses`, `resend`, `postmark`, `mailgun` |
| `storage` | interface + disco local | `s3`, `gcs`, `azure` |
| `queue` | interface + driver sobre `kv` e sobre `database` | `nats`, `sqs`, `rabbitmq` |
| `ws` | `coder/websocket` e `gobwas/ws` | nenhum |

Dois padrões que aparecem sozinhos e mantêm a conta baixa:

1. **Onde há protocolo comum, um driver serve muitos produtos.** RESP cobre
   Dragonfly, Redis, Valkey e KeyDB. A API do S3 cobre S3, R2, Spaces e MinIO.
2. **Onde a stdlib resolve, não há dependência.** SMTP e disco local.

Por isso o pacote se chama **`kv`**, não `redis`: o nome amarra a identidade a um
produto que nem é o recomendado.

## Nomes em inglês funcional

Nada de nome tupi nos adapters. É o que o usuário digita no import, e clareza
vence marca ali — `Illuminate\Mail`, não nome de fantasia. O tupi fica nos
produtos (`taba`, `moitara`), não na infraestrutura.

## Adapter não é modo

`09-uma-forma-so.md` proíbe segundo caminho de código. A linha divisória:

> **Adapter atrás de uma interface não é um modo. Segundo modelo de dados é.**

Postgres, MySQL e SQLite falam SQL: um `Repository`, um sqlc, um conjunto de
queries, um caminho de autorização. Trocar o driver não duplica nada.

## Alta performance — o que entra e o que não

### Cloudflare R2 — o padrão de object storage

Mesma lógica do Dragonfly abaixo: um produto recomendado atrás de um protocolo
que já é padrão. A API do S3 cobre R2, S3, Spaces, B2 e MinIO, então adotar R2
custa uma linha de configuração.

**O motivo é a conta.** R2 não cobra egresso. Num SaaS que devolve arquivo para
quem o enviou — nota fiscal, relatório, anexo — o egresso é a maior parte da
fatura de storage em S3, e cresce com o número de clientes. É o tipo de custo
que só aparece quando o produto começa a dar certo.

Ressalva registrada: R2 não implementa a API inteira do S3. O que este adapter
usa — PUT, GET, HEAD, DELETE e ListObjectsV2 com paginação — é suportado; ACL
por objeto, versionamento e classes de armazenamento não são, e nenhum deles
está no contrato `storage.Store`.

### Dragonfly — entra como padrão recomendado

Compatível com o protocolo Redis (RESP2 e RESP3); o `go-redis` conecta sem
mudança de código. Multi-threaded, substitui um cluster Redis por um nó em
muitos casos.

**Custo de adoção: zero.** Não é adapter novo, é a mesma string de conexão.

Ressalvas que o `kv` não pode ignorar:

- BSL 1.1 — livre para auto-hospedar, proibido oferecer como serviço gerenciado.
  Converte para Apache 2.0 em 2029
- Sentinel não suportado
- eviction limitada a `allkeys-lru` e `volatile-lru`
- lacunas em módulos (Search, JSON, TimeSeries) e em Lua avançado

**Consequência para o core: nada em `kv` pode depender de RedisJSON, RediSearch
ou Lua avançado.** Se depender, o Dragonfly deixa de ser drop-in.

### WebSocket — dois drivers, uma interface

`coder/websocket` como padrão (idiomático, próximo da stdlib) e `gobwas/ws` para
alta concorrência (zero-copy, epoll). Mesma API, in-process, sem Node.

Descartados: **Socket.IO e µWebSockets.js são Node.js.** Adotar qualquer um
transforma uma aplicação Arandu em dois processos e dois runtimes, e mata a
frase que sustenta o produto. Socket.IO é pior ainda: é um protocolo que exige
cliente JavaScript no navegador, e o front aqui é HTMX sem bundler.

### Scylla / Cassandra — não construir

CQL não é SQL: sem join, sem transação entre partições, sem foreign key, schema
modelado por query. Não é trocar driver — é **uma segunda camada de dados**, com
um segundo conjunto de queries para cada módulo, disciplina de modelagem
diferente, e o escopo de tenant virando chave de partição.

Dobra o custo de tudo que vier depois, para resolver um problema de escala que
ninguém no roadmap tem.

**E o ganho não está aí.** Num SaaS o gargalo quase nunca é throughput do engine:
é N+1, índice faltando, serialização e custo de boot. Go sem ORM, sem reflection,
sem boot por request e com sqlc já removeu os três maiores. Postgres em hardware
decente faz dezenas de milhares de transações por segundo.

**Gatilho:** cliente real com carga **medida**, não estimada, batendo no limite do
Postgres depois de índice, particionamento, réplica de leitura e pool ajustados.
Nesse dia, avaliar — e a avaliação começa pelo custo de manter duas camadas.

## Ordem de construção

Fase 2: `database` e `kv` — **feitos**. O contrato de dados fica no core
(`framework/data`, onde o `Grant` é obrigatório) e o `database` fica com tudo que
precisa de driver, um módulo por driver. Ver ADR 0014.
Fase 3: `queue` e `storage` — **feitos**, e são o nono e o décimo repositórios,
autorizados pela ADR 0016. Eram o décimo e o décimo primeiro até a ADR 0021
dissolver o `porang`.

A fila padrão é a **tabela da própria aplicação**, não o Redis: um job empurrado
dentro de `data.Transaction` é comitado pela mesma transação da linha que o
produziu, que é a garantia do outbox aplicada a trabalho. O driver RESP é o
segundo, para volume que a tabela não aguenta.

O object storage padrão é o **Cloudflare R2**. Não cobra egresso, e num SaaS que
devolve arquivo para quem o enviou o egresso é a maior parte da conta em todo o
resto — é a linha que cresce junto com o sucesso. O pacote se chama `s3` porque
é o nome do protocolo, que cobre R2, S3, Spaces, B2 e MinIO; o produto
recomendado aparece em `s3.R2(...)`, como o Dragonfly aparece no `kv`.
Com demanda: `mail`, `ws`, e os submódulos de cada domínio.
