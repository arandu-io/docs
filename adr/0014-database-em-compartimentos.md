# ADR 0014-database-em-compartimentos — O `arandu-io/database` existe, com um driver por módulo

**Status:** Aceita. **Substitui a ADR 0013**, que decidiu o contrário. Restaura o
que `11-adapters.md` e `05-repositorios.md` já previam.

## Contexto

A ADR 0013 argumentou que `arandu-io/database` seria redundante, porque
`database/sql` já é a interface comum e `framework/data` já é o contrato. O
argumento tinha um furo, e ele é medível:

```
$ cat arandu/go.mod          # antes
require (
    github.com/jackc/pgx/v5 v5.10.0
    modernc.org/sqlite v1.55.0
)
# mais quinze linhas de indireto
```

Todo projeto gerado carregava **os dois drivers**, sempre. Um projeto que só fala
SQLite levava pgx no `go.sum`, no build, no binário e na superfície de
vulnerabilidade — e isso não é hipótese: o `govulncheck` encontrou a
**GO-2026-5004** (SQL injection alcançável a partir de qualquer `List`) num
projeto que poderia não falar Postgres nenhum.

É exatamente o problema que abre o `11-adapters.md`: **em Go não existe
dependência opcional.** A ADR 0013 aplicou a exceção errada — disse que a regra
valia para SES e NATS mas não para SQL, quando o que a regra descreve é o
tamanho do grafo de dependência, e o driver de banco é dos maiores.

## Decisão

`arandu-io/database` existe, com **um módulo por driver**:

```
arandu-io/database/          → registro, política de pool, erro que nomeia o que falta
├── go.mod                     depende só do framework. Zero driver
├── sqlite/go.mod            → modernc.org/sqlite
├── pgx/go.mod               → jackc/pgx
└── mysql/go.mod             → go-sql-driver/mysql
```

## O que se compartilha, e como

Esta é a pergunta que a decisão precisava responder, e a resposta é uma direção
só:

> **O contrato fica no core. O que precisa de driver fica no adapter.**
> `database` importa `framework`; `framework` não importa `database`.

`framework/data` continua dono de `Repository`, `Query`, `DB`, `Dialect`,
`Transaction` e `Migration` — porque é ali que o `Grant` é obrigatório, que é a
razão de um handler não alcançar o banco sem passar por uma Policy. Mover isso
para o adapter moveria a tese do produto para um pacote opcional, e **garantia
opcional não é garantia**.

O adapter fica com tudo que precisa de driver, e só com isso:

- `database.Register(dialect, driverName)` — cada compartimento chama no `init()`
- `database.Open(cfg)` — resolve o driver, abre, ajusta o pool, devolve `*data.DB`
- a política de pool, que sai do `main.go` de cada projeto: pool limitado no
  servidor (o default do `database/sql` é ilimitado, e transforma um pico de
  tráfego em "too many connections" em vez de fila no processo), um escritor só
  no SQLite (mais que isso só converte a espera em "database is locked")
- o `MkdirAll` do diretório do SQLite, que ninguém lembra e cujo erro não nomeia

## O erro que isso compra

Sem driver linkado, o `database/sql` diz `unknown driver "pgx"` — e manda a
pessoa procurar erro de digitação no `.env`. O problema é um import que falta, e
nada naquela mensagem diz isso. Agora:

```
DB_CONNECTION is pgsql and no driver for it is linked into this binary
(linked: sqlite).
Add it:

    go get github.com/arandu-io/database/pgx

and blank-import it in cmd/app/main.go:

    _ "github.com/arandu-io/database/pgx"
```

Nomeia o que está configurado, o que está linkado, o comando e a linha. É a tese
do produto — *debug que diz a causa provável* — aplicada ao primeiro erro que
qualquer pessoa encontra.

## O que **não** muda

- **Trocar de banco continua sendo uma linha no `.env`** (ADR 0009). A lista de
  imports decide quais motores o binário sabe falar; o `.env` decide qual ele
  usa. O esqueleto vem com SQLite e Postgres, que é o que o `.env` oferece
- **Uma `Repository`, um sqlc, um conjunto de queries, um caminho de
  autorização.** REGRA 11 intacta: adapter atrás de interface não é modo
- **O core continua com duas dependências** (ADR 0004)
- `data.Migrate`, `data.Transaction` e o outbox continuam onde estão

## Consequências

- **Fecha:** o contrato de dados não sai do core. Se um dia alguém propuser mover
  `Repository` para cá, a resposta é o parágrafo do meio deste documento
- **A Fase 2 volta a ter dez repositórios**, com `database` no lugar que a doc 11
  já reservava e sem `oka` (ADR 0012)
- **Preço registrado:** as tags viram `pgx/v1.2.0` em vez de `v1.2.0`, e um
  release que toca raiz e compartimento precisa de duas tags em ordem. É o preço
  que o `11-adapters.md` já tinha aceitado
- **CI verifica o que a divisão promete:** falha se um driver alcançar o módulo
  raiz, se o compartimento SQLite carregar Postgres, ou se o de Postgres carregar
  SQLite. Sem isso a divisão para de pagar e ninguém percebe
- **MySQL passa a ser opt-in de verdade**, um `go get` — antes estaria no grafo
  de todo mundo ou em lugar nenhum
