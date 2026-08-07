# Perfis — convencional e performance

O Arandu entrega **dois conjuntos de recursos**: perfil convencional e perfil de
performance. Não são dois frameworks e não são dois caminhos de código escritos
à mão. São duas fiações da mesma arquitetura.

## Por que é viável aqui

Porque existe gerador determinístico. `aru make:module` emite `queries.sql` e
`queries.cql` a partir da mesma especificação. A duplicação fica **dentro do
gerador**, não espalhada no código do cliente.

Um framework onde o dev escreve a query à mão não consegue fazer isso sem dobrar
o trabalho de todo mundo. Aqui a conta fica com a engine.

## O que é idêntico nos dois perfis

Isto não muda, e é o que preserva a REGRA 9:

- `security.Grant` e Policy — mesmo caminho obrigatório de autorização
- contrato de módulo `Routes(r *httpx.Router)`
- validação, CSRF, sessão, rate limit, headers
- Collector, página de erro, `aru doctor`
- `kyse` + HTMX + Alpine na view, com o CSS Tailwind compilado e embutido
- toda a CLI

**A API pública é a mesma.** Trocar de perfil não muda uma linha do handler,
do service ou da policy.

## O que muda

| Camada | Convencional | Performance |
|---|---|---|
| Banco | PostgreSQL + `database/sql` | ScyllaDB + gocql |
| KV / cache / sessão | Dragonfly (RESP) | Dragonfly (RESP) |
| Filas | `kv` ou Postgres | NATS |
| Tempo real | SSE (`net/http`) | WebSocket (`gobwas/ws`) |
| Escopo de tenant | coluna `tenant_id` | chave de partição |

Dragonfly é o mesmo nos dois: protocolo RESP, mesmo driver, custo zero.

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

```go
// Tiered is optional. Without it, the module supports the conventional
// profile only.
type Tiered interface {
    Profiles() []Profile // ProfileConventional, ProfilePerformance
}
```

Os dois rótulos são `conventional` e `performance`, e são os mesmos que o
`arandu.mod.toml` já publica no campo `profiles`.

`aru doctor --perfil=performance` falha quando:

1. um módulo do projeto não declara suporte ao perfil
2. um `queries.sql` usa `JOIN`, subquery correlacionada ou `BEGIN`/`COMMIT`
   cruzando agregados
3. uma entidade não tem chave de partição definida
4. o tenant não vem do `Grant` para compor a partição

Regra do ecossistema: **módulo de primeira parte suporta os dois perfis.** Módulo
de comunidade declara o que suporta, e a registry mostra isso na listagem.

## Como se escolhe

```bash
aru new crm                          # convencional, o padrão
aru new crm --perfil=performance
aru make:module cliente              # emite queries.sql
aru make:module cliente --perfil=ambos   # emite queries.sql e queries.cql
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
