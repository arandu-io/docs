# Perfis — convencional e performance

O Arandu prevê **dois conjuntos de recursos**: perfil convencional e perfil de
performance. Não são dois frameworks e não são dois caminhos de código escritos
à mão. São duas fiações da mesma arquitetura.

Este documento é **desenho**. O perfil convencional está construído; o de
performance **não tem uma linha de código**. As seções que descrevem o que muda
entre os dois marcam onde a promessa para; as duas que descrevem o que é
*idêntico* nos dois não marcam — e vale ler com isso em mente, porque não há
segundo perfil para trocar. *Corrigido em 07/08/2026:* a abertura dizia que o Arandu "entrega" os dois
conjuntos, no presente, o que fazia a página inteira ler como descrição do que
existe.

## Por que é viável aqui

Porque existe gerador determinístico. `aru make:module` emite o Repository com o
**SQL parametrizado escrito à mão nos templates do gerador**, a partir da mesma
especificação que emitiria o CQL. A duplicação fica **dentro do gerador**, não
espalhada no código do cliente.

Um framework onde o dev escreve a query à mão não consegue fazer isso sem dobrar
o trabalho de todo mundo. Aqui a conta fica com a engine.

*Corrigido em 07/08/2026:* esta seção dizia que o `aru make:module` emite
`queries.sql` e `queries.cql`. Nenhum `.sql` é emitido, e nunca foi: o gerador
produz doze arquivos — oito `.go` (controller, model, policy, repository,
service, request, teste e migration) e quatro `.kyse.go` (index, show, create,
edit) —, com o SQL em literal Go dentro do Repository. O `sqlc`, que era o
mecanismo por trás da promessa de `.sql`, nunca entrou (ADR 0024) e virou
gatilho: ele passa a valer no dia em que o perfil de performance exigir a
segunda emissão de queries, que é exatamente o assunto deste documento.

## O que é idêntico nos dois perfis

Isto não muda, e é o que preserva a REGRA 9:

- `security.Grant` e Policy — mesmo caminho obrigatório de autorização
- contrato de módulo `Routes(r *http.Router)`
- validação, CSRF, sessão, rate limit, headers
- Collector, página de erro, `aru doctor`
- `kyse` + HTMX + Alpine na view, com o CSS Tailwind compilado e embutido
- toda a CLI

**A API pública é a mesma.** Trocar de perfil não muda uma linha do handler,
do service ou da policy.

## O que muda

| Camada | Convencional | Performance |
|---|---|---|
| Banco | PostgreSQL, MySQL ou SQLite + `database/sql` | ScyllaDB + gocql |
| KV / cache / sessão | Dragonfly (RESP) | Dragonfly (RESP) |
| Filas | tabela `jobs` da aplicação, ou `queue/kv` (RESP) | NATS |
| Tempo real | SSE (`net/http`) | WebSocket (`gobwas/ws`) |
| Escopo de tenant | coluna `tenant_id` | chave de partição |

Dragonfly é o mesmo nos dois: protocolo RESP, mesmo driver, custo zero.

**A coluna da direita é inteiramente desenho.** Não há uma linha de código do
perfil de performance em lugar nenhum do projeto: `grep -rni
"gocql\|scylla\|cassandra\|nats\|gobwas"` sobre `framework`, `aru`, `arandu`,
`database`, `kv`, `queue`, `storage` e `examples` devolve duas ocorrências, e
nenhuma é código: um comentário em `framework/events/relay.go`, dizendo que o
framework **não** escolhe o transporte, e uma linha de prosa do
`framework/README.md` afirmando que Postgres, Redis e NATS chegam por módulos
adapter separados — o adapter NATS não existe, e essa linha é uma promessa em
aberto no README daquele repositório, fora deste diretório. Os dialetos que
`framework/data` aceita são `sqlite`, `pgsql` e `mysql`, e não há um quarto. A
coluna da esquerda também não está completa: SSE não existe — ver
`13-estrutura-illuminate.md`.

## A restrição que o perfil de performance impõe

CQL não é SQL: sem join, sem transação entre partições, sem foreign key, schema
modelado por query.

Consequência, e ela precisa estar clara para quem escreve módulo:

> **O perfil de performance só suporta módulos com acesso agregado.**
> Uma entidade por vez, sem join entre entidades, sem transação cruzando
> agregados.

Isso não é limitação nossa — é a natureza de banco wide-column. A boa notícia é
que o contrato `data.Repository` **já é agregado por construção**: `Find`,
`List`, `Create`, `Update`, `Delete` sobre uma entidade, com `Grant`. Nunca
expusemos join nem query ad-hoc. O desenho já estava certo.

## Módulo declara o perfil que suporta

O que existe hoje é **declarativo, não uma interface Go**: o campo `profiles` do
`arandu.mod.toml`, lido por `aru/internal/manifest`. Os quatro adapters
publicados declaram `profiles = ["conventional", "performance"]`; os módulos de
exemplo declaram `["conventional"]`. Os dois rótulos são esses, e não há um
terceiro.

*Corrigido em 07/08/2026:* aqui havia um bloco Go declarando uma interface
`Tiered` com `Profiles() []Profile`, apresentada como se fosse código do
framework. Ela nunca foi escrita: `grep -rn "Tiered" framework aru arandu` não
devolve nada, e sai com status 1. Nem `Tiered`, nem `Profile`, nem
`ProfileConventional`, nem `ProfilePerformance` existem. Manter o bloco fazia a
página parecer referência de API, e o que ela descreve é desenho.

O critério de recusa do `aru doctor` para o perfil de performance seria: módulo
sem suporte declarado, query com `JOIN`, subquery correlacionada ou
`BEGIN`/`COMMIT` cruzando agregados, entidade sem chave de partição, e tenant que
não vem do `Grant` para compor a partição. **Três dos quatro não existem**, e o
subcomando não aceita a flag — ver a seção seguinte.

O quarto já existe na forma geral, e é das regras mais fortes do `doctor` hoje:
`tenant-from-request` e `tenant-from-header` reprovam tenant vindo de path,
corpo, query ou header em vez do `Grant` — é a REGRA 14 sendo cobrada. O que
falta é a metade da chave de partição, porque não existe partição.

Regra do ecossistema: **módulo de primeira parte suporta os dois perfis.** Módulo
de comunidade declara o que suporta, e a registry mostra isso na listagem — a
registry é a `moitara` do `18-marketplace.md`, Fase 4, e também não existe ainda.

## Como se escolhe

```bash
aru new crm                          # convencional, o padrão
aru new crm --profile=performance
aru make:module client               # SQL parametrizado no Repository
aru make:module client --profile=both    # o mesmo, mais o CQL do agregado
```

**Nenhuma flag de perfil existe ainda** — nem aqui, nem no `aru doctor` da seção
anterior. O perfil de performance é trabalho posterior ao convencional, e a `aru`
de hoje tem exatamente isto:

| Comando | Flags que existem |
|---|---|
| `aru new` | `-module` |
| `aru make:module` | `-fields`, `-tenant`, `-force`, `-dry-run` |
| `aru doctor` | `-strict` |

Os blocos acima são o **desenho**, não a CLI. E quando as flags entrarem, elas
entram em inglês como todas as outras (REGRA 6): `--profile`, não `--perfil`.

O perfil vive em `bootstrap/app.go` — a fiação da aplicação. O diretório
`internal/wire/` **não existe**, e nada nele é gerado: o `aru make:module`
*imprime* as três linhas de fiação, e quem as cola é você (ADR 0023). Trocar de
perfil é reescrever essas linhas e regerar as queries, não reescrever a
aplicação.

Um diretório de fiação gerada convivendo com o `bootstrap/app.go` escrito à mão
seriam duas fiações, que é exatamente o que a REGRA 9 proíbe.

## Ordem de construção

O perfil convencional é a Fase 2. O perfil de performance entra depois do
convencional estar em produção, porque ele **compila contra o mesmo contrato** —
sem o contrato estável, o segundo perfil vira adivinhação.

Não é "se houver demanda": é sequência técnica. O alvo precisa existir antes de
ter dois atiradores.

## O que continua fora

Socket.IO e µWebSockets.js são Node. Nenhum perfil traz um segundo runtime para
dentro da aplicação — isso não é performance, é o problema que o Octane existe
para remediar no PHP. Ver `09-uma-forma-so.md`.
