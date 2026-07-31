# 15 — Multi-tenancy

Tenancy é do framework, não do módulo. Se cada módulo resolver isolamento do seu
jeito, o primeiro que errar vaza dado entre clientes — e num SaaS isso é o fim.

## O modelo escolhido: coluna discriminadora

| Modelo | Decisão | Motivo |
|---|---|---|
| Banco por tenant | rejeitado | multiplica migration, pool e backup por cliente; inviável com centenas |
| Schema por tenant (Postgres) | rejeitado | o sqlc gera contra **um** schema; N schemas quebram a geração |
| **Coluna `tenant_id`** | **escolhido** | uma migration, um pool, escopo vindo do `Grant` |

No perfil de performance a mesma coluna vira **chave de partição** — ver
`12-perfis.md`. É a única diferença entre os perfis nesse assunto.

## Onde o tenant entra, e de onde ele nunca vem

Resolução no middleware `Tenant`, **antes** da sessão, nesta ordem:

1. domínio ou subdomínio (`cliente.app.com`) — o padrão
2. header `X-Tenant` — só quando `Env == dev`
3. path (`/t/<slug>/…`) — opcional, para quem não tem domínio por cliente

O tenant resolvido entra em `Subject.Tenant`, gravado na sessão. A partir daí a
**única** fonte para SQL é `data.Tenant(g)`, que lê do `Grant`.

**Nunca** de parâmetro de rota, corpo do request, query string ou header em
produção. Um tenant vindo do cliente é o cliente escolhendo os dados que quer ler.

## O que o framework prefixa automaticamente

Isolar a tabela e esquecer o resto é o erro clássico. Tudo abaixo carrega tenant:

- chave de sessão e de cache (`kv`)
- chave de rate limit — senão um tenant esgota a cota dos outros
- chave de lock do scheduler (ver `16-agendamento.md`)
- `tenant_id` no logger de escopo e no `Collector`
- nome de arquivo e prefixo em `storage`

## Jobs, migrations e o `SystemGrant`

Job não tem sessão, então o tenant viaja no payload e o worker monta
`security.SystemGrant(acao, tenant)`. A assinatura já exige tenant explícito — não
existe `SystemGrant` sem tenant, de propósito.

Migration é **global**: schema é compartilhado, roda uma vez.

## Verificações do `aru doctor`

- módulo declarado com `--tenant` cujo SQL não filtra por tenant → **erro**
- `tenant` lido de `r.PathValue`, `r.FormValue` ou header fora de dev → **erro**
- chave de cache, lock ou rate limit sem prefixo de tenant → **erro**
- `SystemGrant` com tenant vazio → **erro**

## Single-tenant não é um segundo caminho

Quem não precisa de tenancy roda com `tenant_id` fixo em uma constante. Mesmo
código, mesmo SQL, mesmo `Grant`. Não existe "modo single-tenant" — seria uma
segunda forma de fazer a mesma coisa (REGRA 9).

## Gatilho para revisar

Exigência regulatória de isolamento **físico** por cliente. Nesse caso a saída é
um adapter separado com banco por tenant, não um modo dentro do adapter atual.
