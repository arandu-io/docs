# ADR 0009-multi-banco-sqlite-default — SQLite por padrão, Postgres em produção, com SQL portável

**Status:** Aceita. Substitui a parte do ADR 0007 que dizia "Postgres é o banco suportado na Fase 1".

## Contexto

Decisão do proprietário: seguir o modelo do Laravel — SQLite no `.env` por
padrão, e o cliente troca para Postgres ou MySQL quando quiser. Motivo prático:
projeto novo tem que rodar sem instalar nada.

O problema é que **o Laravel só consegue isso porque tem Query Builder**. Ele
escreve `$table->string('email')->unique()` e gera o DDL de cada dialeto; escreve
`User::where(...)` e gera o SELECT de cada dialeto. Nós rejeitamos essa camada em
`01-arquitetura.md` — a tese é SQL previsível, que se lê no repositório.

Havia três caminhos:

1. Portar um query builder — resolve, e mata a tese do produto.
2. SQL por dialeto em cada repositório — dobra o código e o gerador da Fase 2.
3. **Escrever SQL no subconjunto que todos os bancos compartilham.**

## Decisão

Caminho 3. Portabilidade por escrita, não por abstração.

A camada inteira é `data.Dialect`, e ela faz **uma** coisa: reescrever o
placeholder. Toda query se escreve com `?`, e o Postgres recebe `$1, $2`.
Literais entre aspas são preservados. Não há tradução de tipo, nem de função,
nem de DDL.

O que isso obriga, e que virou regra para todo repositório:

| Em vez de | Use | Porque |
|---|---|---|
| `uuid` + `gen_random_uuid()` | `TEXT` + id gerado na aplicação | são três grafias diferentes da mesma ideia |
| `jsonb` / `text[]` | `TEXT` com JSON | array de Postgres exige tipo de driver; `jsonb` não existe no SQLite |
| `timestamptz DEFAULT now()` | `TIMESTAMP`, valor vindo do Go | default de banco não é portável |
| índice funcional / collation | valor normalizado na escrita | `UNIQUE` simples passa a bastar |
| `RETURNING` | `INSERT` + valor já conhecido em Go | MySQL não tem `RETURNING` |

Configuração no padrão do Laravel: `DB_CONNECTION`, `DB_DATABASE`, `DB_HOST`,
`DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`, com `DATABASE_URL` ganhando de todos
quando existe. Padrão: `sqlite` em `database/database.sqlite`.

Migrações também no modelo do Laravel: número de lote por execução,
`migrate:rollback` desfazendo o último lote, `migrate:status`, e `migrate:fresh`
que **se recusa a rodar fora de `APP_ENV=dev`** — prompt de confirmação não é
proteção às três da manhã.

Driver SQLite: `modernc.org/sqlite`, que é Go puro e não exige cgo. O core
continua com stdlib e `x/crypto`.

*Corrigido em 07/08/2026:* esta linha dizia que o driver "vive no esqueleto,
junto do `pgx`". A **ADR 0014** o moveu para compartimento — `arandu-io/database`
com um `go.mod` por driver. O esqueleto só faz o import em branco
(`arandu/main.go:37-38`), que é o que liga o driver no binário sem que ele
apareça no grafo de quem não o usa.

## Consequência

**MySQL é suportado, e não é a recomendação.** — *corrigido em 07/08/2026.*

Esta consequência dizia o contrário: *"MySQL está declarado, não suportado. A
conexão funciona; os repositórios não, porque usam `RETURNING`."* Ela contradizia
a própria tabela quatro linhas acima, que lista `RETURNING` entre o que **não** se
usa justamente porque o MySQL não tem. Nenhum repositório usa `RETURNING` — a
busca devolve zero ocorrências fora de comentário.

O erro se propagou para dois lugares publicados: o comentário de `DialectMySQL`
em `framework/data/dialect.go`, que mandava o leitor a esta ADR para confirmar, e
o `database/mysql`, que dizia o oposto e estava certo. A suíte de conformidade da
ADR 0018 roda as queries geradas contra servidor de verdade nos três engines, o
MySQL entre eles, e passa.

O que continua verdade é a recomendação: Postgres é onde a história de migração,
o DDL transacional e o relay do outbox surpreendem menos.

**MongoDB não entra.** Não é SQL, não tem transação com a mesma semântica, e nem
o Laravel o suporta sem pacote de terceiro. Um adapter próprio seria outro
produto.

O preço a pagar é conhecido: alguns recursos bons de Postgres — `jsonb` com
índice GIN, tipos de array, full-text nativo — ficam fora do que um módulo
first-party pode usar. Módulo de terceiro que queira exigir Postgres pode; o
`aru doctor` da Fase 2 deve avisar quando isso acontecer.

Revisitar se aparecer módulo first-party que precise de recurso exclusivo de um
banco. A resposta então é adapter, não query builder.
