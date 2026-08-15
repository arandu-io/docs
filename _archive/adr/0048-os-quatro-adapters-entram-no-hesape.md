# ADR 0048 — Os quatro repositórios de adapter entram no `hesape`

**Status:** Aceita. Data: 10/08/2026.
**Supersede parcialmente:** ADR 0016 (doze repositórios), na parte que separa
`database`, `kv`, `queue` e `storage` do framework.

## Contexto

O `arandu-io/database`, o `arandu-io/kv`, o `arandu-io/queue` e o
`arandu-io/storage` existem como repositórios próprios desde a ADR 0016. A razão
declarada era a REGRA 11: adapter atrás de uma interface não é um modo, e o
driver pesado vive em submódulo porque **em Go não existe dependência
opcional**.

O que mudou é o alvo. O `hesape` replica o `Illuminate` componente por componente
(ADR 0043 e 0044), e no Laravel esses quatro **não são pacotes separados**:

    Illuminate\Database    tem Connectors/ com MySql, Postgres, SQLite dentro
    Illuminate\Redis       tem Connections/ e Connectors/ dentro
    Illuminate\Queue       tem Connectors/ e Jobs/ dentro
    Illuminate\Filesystem  tem o adapter de S3 dentro

Um desenvolvedor que procura `RedisQueue` procura em `Illuminate\Queue`. Aqui ele
teria que saber que mora noutro repositório — o que é exatamente o custo que a
ADR 0044 existe para não cobrar.

O proprietário decidiu em 10/08/2026: os quatro entram, no padrão do Laravel,
mantendo os bancos que a REGRA 11 e a REGRA 12 já nomeiam e nenhum a mais.

## Decisão

**O conteúdo dos quatro repositórios passa a viver dentro do `hesape`, nos
endereços do espelho, e o driver pesado continua em submódulo.**

| origem | destino | módulo |
|---|---|---|
| `database/*.go` | `hesape/database` | o do `hesape` |
| `database/pgx` | `hesape/database/connectors/pgx` | próprio |
| `database/mysql` | `hesape/database/connectors/mysql` | próprio |
| `database/sqlite` | `hesape/database/connectors/sqlite` | próprio |
| `kv/*.go` | `hesape/redis` | **próprio** |
| `queue/*.go` | `hesape/queue` | o do `hesape` |
| `queue/kv` | `hesape/queue/connectors/redis` | próprio |
| `storage/*.go` | `hesape/filesystem` | o do `hesape` |
| `storage/s3` | `hesape/filesystem/s3` | próprio |

O submódulo não é uma concessão: é a única forma de o `go get` do core não
arrastar `pgx`, `go-redis` e o SDK da AWS. O CI do `hesape` reprova qualquer
dependência além de `golang.org/x/crypto`, e essa verificação continua valendo
para o módulo raiz.

### O `redis` inteiro é submódulo, e o Laravel não faz assim

No Laravel, `RedisStore` mora em `Illuminate\Cache`, `RedisQueue` em
`Illuminate\Queue` e o handler de sessão sobre Redis em `Illuminate\Session`.
Cada um no componente que ele serve, e o `predis` é `suggest` no `composer.json`.

Aqui não pode: `hesape/cache` importando o driver faria todo projeto Arandu
baixar `go-redis` para usar o cache em memória. Então **tudo que depende do
driver vive em `hesape/redis`** — a conexão, o store de cache, o lock, o handler
de sessão e o limiter.

Os nomes continuam sendo os do Laravel: `RedisStore`, `RedisLock`,
`RedisQueue`, `CacheBasedSessionHandler`. O que muda é o endereço, e o `doc.go`
de cada componente afetado diz para onde foi.

Esta é a segunda vez que a ausência de dependência opcional em Go move uma peça
de lugar. É a mesma razão da primeira, e a lista de exceções não deve crescer
além dela.

### Só os bancos já decididos

Convencional: **Postgres, MySQL, SQLite** (REGRA 11).
Performance: **Scylla/Cassandra**, que é perfil e não adapter solto (REGRA 12).
RESP: **Dragonfly, Redis, Valkey, KeyDB**, um driver só, e nada pode depender de
RedisJSON, RediSearch ou Lua avançado.

Nenhum banco novo entra por esta ADR. Ela move o que existe.

## O que não é apagado

Os quatro repositórios publicados **continuam existindo** até decisão em
separado. Foram espelhados com `git clone --mirror` em
`~/Developer/arandu-espelhos/` antes de qualquer movimento, e o proxy do Go
continua servindo as versões já publicadas para quem as importa hoje.

Apagar repositório publicado é decisão do proprietário, pedida na hora, e não
consequência automática desta.

## Consequências

- Os imports mudam de `arandu-io/framework/*` para `arandu-io/hesape/*`:
  `security` → `auth`, `data` → `database`, `jobs` → `queue`,
  `storage` → `filesystem`, `httpx` → `http` (ADR 0047), `kernel` →
  `foundation`, `config` → `config`.
- `aru make:module` passa a emitir os imports novos, e o esqueleto publicado
  muda junto.
- A REGRA 3 volta a doze repositórios menos quatro. O `hesape` não é o décimo
  terceiro (ADR 0043): ele é o nome do conjunto, e agora o conjunto é maior.
- A suíte de conformidade do `database` — a que prova que os três drivers
  respondem igual — vem junto e passa a rodar no CI do `hesape`.
