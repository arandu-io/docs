# 30 — o que falta em relação ao padrão do Laravel

Levantado em 08/08/2026, contra o código que existe: quatorze pacotes em
`framework/`, um módulo (`auth`), trinta e quatro comandos na `aru`, quatro
repositórios de adapter e o `kyse` com onze componentes.

Não é uma lista de desejos. Cada linha diz **o que o Laravel tem**, **o que
existe aqui hoje** e **o que falta de fato** — e onde a resposta é "não vai
existir", diz por quê, porque metade da superfície do Laravel existe para
remediar perigo de linguagem dinâmica que Go não tem (REGRA 7).

Os pacotes de primeira parte do Laravel — Sanctum, Telescope, Horizon, Scout,
Cashier e os outros — já têm veredito em `docs/packages/`. Este documento é
sobre o **framework**, não sobre eles.

---

## 1. O que já espelha o Laravel e está fechado

Para a lista abaixo ser lida como o que falta, e não como o que existe:

| Laravel | Arandu | Onde |
|---|---|---|
| Estrutura de diretórios | idêntica | `app/`, `bootstrap/`, `config/`, `database/`, `resources/`, `routes/`, `storage/`, `tests/` |
| `artisan` | `aru`, 34 comandos | `make:*`, `migrate*`, `db:seed`, `route:list`, `schedule:*`, `queue:work`, `doctor` |
| Blade | kyse | `.kyse.go` → Go em tempo de build, erro na linha do fonte |
| `laravel/ui` | `arandu-io/ui` | nove telas, publicadas |
| shadcn/ui | `arandu-io/kyse` | onze componentes, importados |
| Migrations | `database/migrations` | `aru migrate`, rollback, status, fresh |
| Seeders e Factories | `database/seeders`, `database/factories` | |
| Form Requests | `app/Http/Requests` | `validation` |
| Policies | `app/Policies` | e são **obrigatórias**: sem `Grant` não há repositório |
| Middleware | `app/Http/Middleware` | pipeline no `bootstrap` |
| Eventos e Listeners | `events`, `app/Events`, `app/Listeners` | com outbox transacional |
| Jobs e filas | `jobs`, `arandu-io/queue` | `aru queue:work` |
| Scheduler | `scheduler` | `aru schedule:run` |
| Storage | `storage`, `arandu-io/storage` | |
| Telescope | `observability` + console em `/_arandu/debug` | core, ADR |
| Tests | `tests/Feature`, `tests/Unit`, `tests/TestCase.go` | |

---

## 2. Falta, e é bloqueante para um SaaS

### 2.1 Mail — **FEITO**

**Laravel:** `Mail::to()->send(new OrderShipped($order))`, Mailables, markdown
mail, `MAIL_MAILER`, filas por padrão.

**Aqui:** `framework/mail` tem `Mailable`, `Envelope`, `Content`, `Mailer` e
cinco transportes — smtp, log, array, **resend** e **sendgrid**. Os dois
provedores ficam no core e não em submódulo, porque nenhum dos dois arrasta SDK:
são um POST com corpo JSON sobre `net/http` (ADR 0031). Amazon SES continua fora,
e é o único que precisaria de assinatura de requisição.

`Content` tem `View` e `TextView`, ambas por nome de view. Só tinha `Text string`,
literal Go, e por causa disso a view de texto que o projeto gerava **não tinha
como ser enviada**: toda mensagem saía só em HTML.

`ErrRetryable` separa 429 e 5xx de endereço recusado — sem isso um rate limit
engole o link de verificação em silêncio.

**Ainda falta:** `aru make:mail`, e envio pela fila como caminho de uma linha
(hoje é um job que chama `Send`, o que é o desenho certo, mas não tem gerador).

### 2.2 Cache

**Laravel:** `Cache::remember()`, tags, `CACHE_STORE`.

**Aqui:** `arandu-io/kv` cobre sessão, rate limit e lock sobre RESP. Não há
`Cache` como fachada de aplicação, nem `remember`.

**Falta:** uma API `cache` no core sobre o `kv`, com `Remember`, `Forget` e TTL,
prefixada por tenant (REGRA 14). **Sem tags** — elas exigem estrutura que a
REGRA 11 proíbe depender (`RedisJSON`, Lua avançado).

### 2.3 Notifications

**Laravel:** `$user->notify()`, canais (mail, database, broadcast, SMS).

**Aqui:** não existe.

**Falta:** depende do Mail. Provavelmente **não deve** ser um sistema de canais
genérico: o Laravel precisa dele porque a notificação é polimórfica em runtime;
aqui um `Notifier` tipado por canal resolve sem inventar despacho dinâmico.

### 2.4 Verificação de e-mail e registro — **FEITO no `examples`**

**Laravel:** `MustVerifyEmail`, `Auth::routes()`.

**Aqui:** o `framework/auth` ganhou `Register`, `MarkVerified`, `EnsureUser`,
`Lookup` e `Names`; `User` ganhou `Name` e `VerifiedAt`; `Subject` ganhou
`Verified`. O `examples` escreveu os handlers como referência
(`RegisterController.go`) e as duas telas publicadas passaram a ter rota.

O link é **assinado e não armazenado** (`security.Signer`), com propósito dentro
da assinatura, e não abre sessão. Confirmar o endereço compra uma coisa só, e é
o ponto: `CommentPolicy` recusa `CommentCreate` sem ela. ADR 0032.

Auto-registro passa pela policy: convidado declarado (ADR 0029) sobre
`ActionUserCreate`, e a policy decide olhando o **candidato** — sem papel, ou não.

**Ainda falta:** portar isso para o `arandu-io/ui`. O starter kit continua
publicando `register.kyse.go` e `verify.kyse.go` **sem** handler, e a instrução
que ele imprime ainda mostra a assinatura de quatro argumentos do
`authui.New` — que é a dele e está coerente, mas é a versão sem registro.

---

## 3. Falta, e muda o desenho

### 3.1 Relações

**Laravel:** `hasMany`, `belongsTo`, `with()` para evitar N+1.

**Aqui:** não existe, e é deliberado — Active Record foi rejeitado
(`docs/01-arquitetura.md`). O `examples` lê o thread de comentários com um
método de repositório escrito à mão, e é isso que a tese pede.

**Falta decidir:** o `aru make:module` sabe gerar `post_id:string!` e **não sabe
gerar o método que o segue**. A lacuna real não é "relação" — é o gerador não
emitir `ForPost` quando um campo é chave estrangeira. Isso é gerador, não ORM.

### 3.2 Paginação

**Laravel:** `paginate()` e um componente de links.

**Aqui:** cursor keyset no `data.Query`, e o `examples` desenha "Older posts" à
mão.

**Falta:** um componente `Pagination` no `kyse` e um tipo `data.Page` carregando
cursor e total. **Sem offset**: `LIMIT ? OFFSET ?` numa tabela grande é o
problema que o keyset existe para não ter.

### 3.3 Upload de arquivo

**Laravel:** `$request->file('avatar')->store()`.

**Aqui:** `storage` e `arandu-io/storage` existem. Não há caminho de `Request` a
`Storage`, nem validação de tipo, nem limite de tamanho.

**Falta:** `http.Context.File(name)`, validação no Form Request, e um
componente de upload. É pequeno e é visível — todo CRUD real tem uma imagem.

### 3.4 Autorização na view

**Laravel:** `@can('update', $post)`.

**Aqui:** a REGRA 15 fecha a DSL, então `@can` não entra. O `examples` resolve
passando URL vazia: um botão sem endereço não é desenhado.

**Falta:** nada a construir — falta **documentar** que esse é o padrão, porque é
a primeira coisa que um dev Laravel procura.

---

## 4. Falta, e é a `aru`

| Comando do `artisan` | Estado |
|---|---|
| `make:mail` | falta, depende de 2.1 |
| `make:notification` | falta, depende de 2.3 |
| `make:test` | **falta**, e agora que existe `tests/` é uma linha |
| `make:resource` | não vai existir: é serialização de API, e `ctx.JSON` já resolve |
| `make:observer` | não vai existir: são os eventos |
| `make:rule` | falta — regra de validação reutilizável |
| `make:cast`, `make:scope` | não existem: são Eloquent |
| `tinker` | não vai existir (REGRA 3) |
| `optimize`, `config:cache` | não existem em Go: não há bootstrap para cachear |
| `storage:link` | não existe: não há document root |
| `about` | falta — um resumo do que está fiado, e é barato |
| `queue:failed`, `queue:retry` | **faltam** — a dead letter existe e não há como listar nem retentar |
| `db:wipe` | falta (o `migrate:fresh` cobre parcialmente) |
| `vendor:publish` | não existe: publicar é `go run <module>@latest` (ADR 0026) |

---

## 5. Falta, e é infraestrutura de desenvolvimento

### 5.1 Hot reload de verdade

`aru dev` reconstrói. Não há injeção de fragmento nem preservação de estado; um
`hx-boost` recarrega a página inteira. O Laravel também não tem, mas o Vite tem,
e é a comparação que um dev faz.

### 5.2 Realce de sintaxe do kyse

Sem ele, um `.kyse.go` é texto cru — o editor tenta lê-lo como Go e marca a
marcação como erro. Uma gramática TextMate resolve, e está declarada como custo
permanente desde a ADR 0020.

### 5.3 `arandutest`

O pacote existe. Não há `tests/TestCase.go` **no framework** oferecendo o que
toda aplicação repete: subir o kernel, autenticar como um sujeito, afirmar
sobre uma resposta HTML. Cada projeto reescreve.

**Falta:** `arandutest.App(t)`, `arandutest.ActingAs(subject)`,
`arandutest.See(t, html, "texto")`. É o equivalente do `RefreshDatabase` e do
`actingAs`, e é o que faz teste de feature valer a pena escrever.

---

## 6. Falta, e é sobre operar

| | |
|---|---|
| **Health check detalhado** | existe `/_arandu/health`; falta um por dependência |
| **Graceful shutdown** | o kernel para; falta drenar a fila antes |
| **Rate limit por rota** | existe global; falta `.Middleware(throttle("60,1"))`. A entrada tem o seu, por identidade, sob o `Authenticate` — ADR 0038 |
| **Sessão em RESP** | o `kv` oferece; o esqueleto usa memória e diz que é errado com duas réplicas |
| **`queue:failed`** | ver 4 |
| **Métricas** | `observability` exporta OpenTelemetry; falta o exemplo de fiação |

---

## 7. O que decidido **não** existe

Repetido aqui para a lista acima não ser lida como incompleta:

- **Eloquent / Active Record** — `Repository` + SQL parametrizado (ADR 0001)
- **Container de DI** — a fiação é `bootstrap/app.go`, lida de cima a baixo
- **Facades** — ADR 0002
- **`$fillable` / mass assignment** — não existe atribuição em massa
- **Service Providers** — `bootstrap` compõe
- **Octane, Sail, Herd, Valet, Volt, Tinker** — resolvem problemas do PHP
- **Blade components como tag** — componente é função tipada (REGRA 13)
- **`@can`, `@auth`, `@error`** — a DSL é fechada (REGRA 15)
- **Offset pagination** — ver 3.2
- **Cache tags** — ver 2.2

---

## Ordem sugerida

Por dependência, não por tamanho:

1. ~~**Mail** (2.1)~~ — feito. Sobra o `make:mail`
2. **`arandutest`** (5.3) — todo o resto fica mais barato de provar
3. ~~**Registro e verificação** (2.4)~~ — feito no `examples`; falta o `ui`
4. **Upload** (3.3) e **paginação** (3.2) — todo CRUD real precisa dos dois
5. **Cache** (2.2) e **`queue:failed`** (4) — operação
6. **Realce do kyse** (5.2) — a primeira impressão de quem abre o editor

Notifications (2.3) vem depois de Mail, e provavelmente menor do que o Laravel
faz.
