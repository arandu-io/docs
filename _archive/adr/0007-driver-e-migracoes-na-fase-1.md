# ADR 0007-driver-e-migracoes-na-fase-1 — Driver no esqueleto, migrador mínimo no core

**Status:** Aceita. Fecha uma porta em duas frentes; a segunda é temporária por decisão.

## Contexto

A Fase 1 tem que entregar "um binário que sobe, autentica" e `aru migrate`
(`03-roadmap-fases.md`). Duas coisas faltavam para isso, e nenhuma podia sair do
core sem quebrar o ADR 0004:

- um driver de Postgres — o core não pode depender de terceiro;
- um aplicador de migração — o roadmap coloca Atlas na Fase 2, não na 1.

## Decisão

**1. O driver vive no esqueleto.** `arandu-io/arandu` declara `jackc/pgx/v5` e
importa `pgx/v5/stdlib` em `cmd/app/main.go`. Isso não viola o ADR 0004: o
esqueleto **não é biblioteca** — ninguém o importa, você o clona uma vez. O core
segue com stdlib e `golang.org/x/crypto`.

O adapter `arandu-io/postgres` da Fase 2 continua fazendo sentido: pool com
tunning, health, `LISTEN/NOTIFY` e os tipos que o gerador emitir.

**2. O migrador mínimo é core, em `data.Migrate`.** Usa apenas `database/sql`:
tabela `arandu_migrations`, cada migração numa transação junto com o registro, e
parada no primeiro erro. Sai da Fase 2 quando o Atlas entrar.

Consequências que valem registrar por serem irreversíveis na prática:

- **`roles` é `jsonb`, não `text[]`.** Ler um array de Postgres exige tipo
  específico de driver; `jsonb` lê como `[]byte` com qualquer driver, o que
  mantém o core sem dependência.
- **`APP_KEY` aceita o prefixo `base64:`.** 32 bytes aleatórios não são
  imprimíveis, então `aru key:generate` precisa emitir valor codificado. A forma
  crua continua aceita.
- **O SQL do migrador é PostgreSQL.** `timestamptz`, `now()`, placeholder `$1`.
  Postgres é o banco suportado na Fase 1; outro banco é decisão futura, com
  adapter próprio.

## Consequência

O esqueleto depende do framework por versão (`v0.1.0`), sem `replace`. Enquanto o
repositório for privado, quem clonar precisa de `GOPRIVATE=github.com/arandu-io`
e de acesso de leitura — o CI faz isso com um token. Quando os repositórios
virarem públicos, as duas linhas do workflow deixam de ser necessárias e nada no
`go.mod` muda.

Durante o desenvolvimento simultâneo dos dois repositórios, use um `go.work`
local apontando para os dois diretórios; ele não é versionado, e é o que evita
ter que tagear o framework a cada alteração.
