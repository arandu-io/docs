# 15 — Multi-tenancy

Tenancy é do framework, não do módulo. Se cada módulo resolver isolamento do seu
jeito, o primeiro que errar vaza dado entre clientes — e num SaaS isso é o fim.

## O modelo escolhido: coluna discriminadora

| Modelo | Decisão | Motivo |
|---|---|---|
| Banco por tenant | rejeitado | multiplica migration, pool e backup por cliente; inviável com centenas |
| Schema por tenant (Postgres) | rejeitado | o gerador emite contra **um** schema; N schemas multiplicam a emissão |
| **Coluna `tenant_id`** | **escolhido** | uma migration, um pool, escopo vindo do `Grant` |

No perfil de performance a mesma coluna vira **chave de partição** — ver
`12-perfis.md`. É a única diferença entre os perfis nesse assunto.

## Onde o tenant entra, e de onde ele nunca vem

**Não existe middleware `Tenant`.** O tenant é resolvido em **um** ponto do
sistema: o login. O contrato é `auth.TenantResolver`, uma função
`func(*http.Request) string` recebida por `auth.New` na fiação explícita do
`bootstrap/app.go`.

```go
// arandu/bootstrap/app.go, verbatim
auth.New(authService, auth.FixedTenant(cfg.Auth.Tenant))
```

O resolver é consultado **só** no login, que é o único momento em que não há
sessão a quem perguntar. Dali em diante o tenant vive em `Subject.Tenant`,
gravado na sessão, e a **única** fonte para SQL é `data.Tenant(g)`, que lê do
`Grant`.

Um resolver que leia o host é trabalho previsto e ainda **não escrito**; hoje o
único implementado é `FixedTenant`, que devolve sempre a mesma constante. Não é
"modo single-tenant": é o mesmo caminho de código com uma constante no lugar do
resolver, e trocar por outro resolver não muda uma query sequer.

**Nunca** de parâmetro de rota, corpo do request, query string ou header depois
do login. Um tenant vindo do cliente a cada request é o cliente escolhendo os
dados que quer ler.

Um desenho com header `X-Tenant` em dev ou com prefixo de path `/t/<slug>/` já
esteve escrito aqui e **foi descartado**: cada fonte extra é mais um lugar de
onde o tenant pode chegar errado, e resolver só no login é estritamente mais
seguro que resolver a cada request (REGRA 14). Quem for implementar host-based:
é mais um `TenantResolver`, chamado no mesmo lugar, e nada além disso.

## O que o framework escopa por tenant — e o que deliberadamente não escopa

Isolar a tabela e esquecer o resto é o erro clássico. Mas prefixar tudo também
é erro, e o motivo é caso a caso.

**Escopado, e cobrado pela assinatura** — todo método pede `Grant`, e sem tenant
devolve erro em vez de chave global:

| O quê | Como |
|---|---|
| `kv.Cache` | chave `<prefixo>:<tenant>:<namespace>:<id>`; sem tenant, `kv.ErrNoTenant` |
| `storage` | o tenant é segmento do caminho; sem tenant, `storage.ErrNoTenant` |
| `jobs` | `TenantID` vem do `Grant` no `Push`, e o worker recusa job cujo tenant não bate (`jobs.ErrForged`) |
| `events` / outbox | o registro carrega `tenant_id`, vindo do `Grant` |
| SQL de módulo `--tenant` | `WHERE tenant_id = ?` com `data.Tenant(g)` |

**Não escopado, de propósito**, cada um com a razão escrita no código:

| O quê | Por quê |
|---|---|
| chave de sessão (`kv/session.go`) | o id da sessão é **quem diz** qual é o tenant; prefixar exigiria saber a resposta antes da pergunta. São 256 bits de entropia e o cookie leva um HMAC do id |
| chave de rate limit (`kv/limiter.go`) | o rate limit roda **antes** da autenticação justamente nas rotas onde mais importa, e ali não há `Grant`. A chave vem de `middleware.KeyByIP` ou `KeyBySession` |
| lock do scheduler (`kv/lock.go`) | um lock por tenant deixaria N réplicas rodarem a mesma tarefa, cada uma para um tenant. Quem expande por tenant é a tarefa `kernel.PerTenant`, com `scheduler.Options.Tenants` (ver `16-agendamento.md`) |

O logger de escopo de request leva `request_id`, `method` e `path`, **não**
`tenant`. O do worker leva `tenant`, porque ali o job carrega um. Se o
`tenant_id` na linha de log de request passar a ser necessário, é mudança no
`middleware.Observe` — não está lá hoje.

## Jobs, migrations e o `SystemGrant`

Job não tem sessão, então o tenant viaja no payload e o worker monta
`security.SystemGrant(acao, tenant)`. A assinatura já exige tenant explícito — não
existe `SystemGrant` sem tenant, de propósito.

Migration é **global**: schema é compartilhado, roda uma vez.

## Verificações do `aru doctor`

As regras que existem, pelo nome com que o `doctor` as reporta:

| Regra | O que reprova |
|---|---|
| `sql-without-tenant-scope` | módulo declarado com `--tenant` cujo SQL não filtra por tenant |
| `tenant-from-request` | `tenant` lido de `r.PathValue`, `r.FormValue` ou da query |
| `tenant-from-header` | `tenant` lido de header |
| `system-grant-without-tenant` | `SystemGrant` com tenant vazio |
| `system-grant-outside-scope` | `SystemGrant` fora de job, migration e teste (`--strict`) |

**Não existe — e não deve existir — regra de "chave de cache, lock ou rate limit
sem prefixo de tenant".** Ela contradiz as três decisões da seção anterior:
sessão, rate limit e lock são sem prefixo por desenho, então a regra reprovaria
o próprio `kv`. E para o único caso em que o prefixo é obrigatório — o cache —
não há o que verificar estaticamente: `kv.Cache` só tem métodos que pedem
`Grant`, e um cache sem tenant devolve `ErrNoTenant` em runtime. A assinatura já
fecha a porta, que é sempre preferível a uma regra de lint (REGRA 2).

## Single-tenant não é um segundo caminho

Quem não precisa de tenancy roda com `tenant_id` fixo em uma constante. Mesmo
código, mesmo SQL, mesmo `Grant`. Não existe "modo single-tenant" — seria uma
segunda forma de fazer a mesma coisa (REGRA 9).

## Gatilho para revisar

Exigência regulatória de isolamento **físico** por cliente. Nesse caso a saída é
um adapter separado com banco por tenant, não um modo dentro do adapter atual.
