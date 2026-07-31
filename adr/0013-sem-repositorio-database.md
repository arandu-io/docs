# ADR 0013-sem-repositorio-database — O `database` já existe, e é `framework/data`

**Status:** Aceita. Altera `11-adapters.md` e `05-repositorios.md`, que previam
`arandu-io/database` como repositório da Fase 2.

## Contexto

O `11-adapters.md` lista seis pacotes por domínio, e o primeiro é:

| Pacote | Núcleo | Submódulos |
|---|---|---|
| `database` | interface `Repository` + pgx | `mysql`, `sqlite` |

E fecha com: *"Fase 2: `database` (pgx) e `kv`. Sem os dois o framework não roda."*

O framework roda. A Fase 1 fechou com login funcionando em **SQLite sem instalar
nada** e em **PostgreSQL 18**, mudando só o `.env`. O caminho que tornou isso
possível foi outro:

- `framework/data` define `Repository`, `Query`, `DB`, `Dialect`, `Transaction` e
  as migrations. Fala `database/sql`, e por isso **não importa driver nenhum** —
  o core continua com duas dependências (ADR 0004)
- o **projeto** escolhe o driver com um blank import em `cmd/app/main.go`:
  `_ "github.com/jackc/pgx/v5/stdlib"`, `_ "modernc.org/sqlite"`. O esqueleto não
  é biblioteca, ninguém o importa, e por isso pode depender de driver (ADR 0006)
- `Dialect.Rebind` converte `?` em `$N`, então uma query serve os três bancos

Criar `arandu-io/database` agora significaria uma segunda `Repository`, uma
segunda `DB`, um segundo lugar onde a migration mora — para os mesmos três
bancos que já funcionam.

## Decisão

Não existe `arandu-io/database`. O domínio de dados é `framework/data`, e o
driver é escolhido pelo projeto.

## Por quê

**REGRA 9 e REGRA 11.** *Adapter atrás de uma interface não é um modo* — e aqui
a interface é a `database/sql` da stdlib, que já é o ponto de troca. O que a
doc 11 queria evitar (um pacote por fornecedor) já está evitado, e por um
mecanismo mais barato: nenhum pacote por fornecedor, porque a stdlib é o
contrato.

**O motivo que criou a regra não se aplica.** A doc 11 existe porque em Go não há
dependência opcional: se `mail` importar o SDK do SES, todo mundo carrega o SDK.
Isso é verdade para SES, S3 e NATS — SDKs pesados sem interface comum na stdlib.
Não é verdade para SQL: `database/sql` **é** a interface comum, e o driver entra
por efeito colateral de import, no binário final, escolhido por quem faz o
binário.

**O que o `database` traria já está feito.** `Repository` com `Grant`
obrigatório, escopo de tenant vindo do `Grant`, instrumentação no `Collector`,
migrations com batch e rollback, portabilidade por `Rebind`. Mover para outro
repositório seria custo de release sem ganho.

## Consequências

- **Fecha:** o contrato de dados fica no core, e não será fatiado depois sem
  quebrar todo import de todo projeto. Isso é coerente com `05-repositorios.md`:
  *"não fatiar o core agora"*
- **Fica de pé o resto da doc 11.** `kv`, `mail`, `storage`, `queue` e `ws`
  continuam pacotes por domínio em repositório próprio, porque cada um traz
  dependência pesada que o core não pode carregar. `kv` foi criado na Fase 2
- **O driver é responsabilidade do projeto**, e o esqueleto já traz os dois que
  importam. Adicionar MySQL é um blank import
- **Se um dia houver um adapter que não cabe em `database/sql`** — Scylla é o
  candidato, e a doc 12 já o trata como perfil e não como driver — a decisão
  volta à mesa. O gatilho está escrito em `11-adapters.md`: carga real medida
  batendo no limite do Postgres depois de índice, particionamento, réplica de
  leitura e pool ajustados
