# Reorganização do Arandu para o formato Illuminate

> **Como este documento foi feito.** Quarenta e dois agentes, um por componente
> do Illuminate, cada um lendo o componente inteiro e caçando o equivalente no
> Arandu onde quer que ele estivesse. Depois sete revisores de coerência, que
> leram as especificações **umas contra as outras** — porque foram escritas
> isoladas — e resolveram os choques: dois componentes reivindicando o mesmo
> pacote, coisa sem dono, ciclo de dependência. Nenhum deles tinha permissão de
> editar um arquivo. São 1.070 superfícies especificadas.
>
> Ele existe porque o proprietário parou o trabalho anterior com uma frase que
> estava certa: consertar ponto a ponto não reorganiza estrutura.


**Documento de execução.** Escrito contra o código em 10/08/2026, verificado arquivo por arquivo. Nada aqui foi copiado de `docs/10-de-para-laravel.md` ou `docs/13-estrutura-illuminate.md` sem checar antes — os dois carregam afirmações falsas listadas na seção 7.

---

## 1. O alvo, em uma árvore

Quarenta e dois componentes clonados em `laravel_illuminate/`. Trinta e três viram pacote, oito viram nada, um vira repositório-adapter fora da coleção.

```
github.com/arandu-io/framework/
│
├── arandutest/           Testing                Client, Response, Boot, Harness, assertions, fakes
├── auth/                 Auth\Access            Subject, Action, Policy[T], Grant, Authorize, Allows,
│   │                                            Check, SystemGrant, Guest, ValidTenant, Tenant,
│   │                                            ErrForbidden, WithSubject/SubjectFrom, SignInThrottle
│   ├── events/           Auth\Events            os três nomes de evento do outbox
│   ├── middleware/       Auth\Middleware        Authenticate, RedirectIfAuthenticated, RequireRole,
│   │                                            RequirePassword, EnsureVerified
│   ├── notifications/    Auth\Notifications     construtores de ResetPassword e VerifyEmail
│   ├── passwords/        Auth\Passwords         Broker: SendResetLink, Reset, ConfirmPassword
│   └── users/            EloquentUserProvider   User, Provider, Repo, Policy, Service, migration users
├── broadcasting/         Broadcasting           Hub, Channel, Fanout, Bus, Message, o stream SSE
├── bus/                  Bus                    Batch, Chain, PendingBatch, Store, PushUnique
├── cache/                Cache                  Repository, Store, Locking, ArrayStore, Locks, Lock,
│   │                                            RateLimiter, Limit, Remember[T], Get[T], Pull[T]
│   └── cachetest/                               a suíte de contrato que todo Store passa
├── collections/          Collections            Map, Filter, Reduce, GroupBy, KeyBy, Partition,
│                                                UniqueBy, Sum, First, Last
├── concurrency/          Concurrency            Task[T], Run, RunWith, Recover, PanicError
├── config/               Config                 App, Env, LoadDotenv, String/MustString/Bool/Int/Seconds
├── console/              Console                Command, IO, Registry, ExitError, Exit, GuardEnv
├── cookie/               Cookie                 Jar, Config, Make, Forever, Forget, Read, Write
├── database/             Database               Repository[T,ID], Query, Page[T], DB, Dialect,
│                                                Transaction, Migration, Seeder, Register, Open, ErrNotFound
├── encryption/           Encryption             Signer, Sign, Verify, KeySize, ParseKey, GenerateKey
├── events/               Events                 Event, Recorder, Outbox, Stored, Publisher, Relay,
│                                                On, Listeners, Module
├── exception/            Foundation\Exceptions  Handler, Config, Abort, HTTPError, StatusOf, Recover,
│                                                classify, a página de debug, as páginas 401..503
├── filesystem/           Filesystem             Disk, Adapter, Disks, File, Info, Upload, UploadRules,
│                                                URLSigner, Presigner, Send, Key, CleanKey
├── foundation/           Foundation             Kernel, Module + as interfaces opcionais, Builtins
├── hashing/              Hashing                Make, Check, NeedsRehash, Info, IsHashed,
│                                                MinPasswordLen, MaxPasswordLen
├── httpx/                Http                   Context, State, Reject, Back, LocalPath, Intended,
│   │                                            Bind, WithCookie, WithStatus, Middleware (alias),
│   │                                            Renderer/URLGenerator/ExceptionHandler (interfaces)
│   ├── client/           Http\Client            New, Fake, Stub, AssertSent
│   └── middleware/       Http\Middleware        SecurityHeaders, TrustProxies, CORS, LimitBodySize
├── image/                Image                  Image, Driver, Pipeline, Transformation, Std
├── jsonschema/           JsonSchema             Type, Object/String/Integer/Number/Boolean/Array,
│                                                Union, AnyOf, Document, Validate, Parse
├── log/                  Log                    New, Into, For, With, Field, Fields, Middleware,
│                                                ParseLevel, Capture, os oito níveis PSR-3
├── mail/                 Mail                   Mailable, Envelope, Content, Message, Sent, Mailer,
│   │                                            Pending, Attachment, Headers, Raw, Render, Queue
│   └── transport/        Mail\Transport         SMTP, Log, Array, Resend, SendGrid, Postmark, Failover
├── notifications/        Notifications          Notification, Notifiable, Key, ChannelName, Notifier,
│                                                Channel, Store, Record, Anonymous, Capture
├── number/               Support\Number         FileSize, Currency, Percentage, Ordinal, Abbreviate
├── pagination/           Pagination             LengthAwarePaginator[T], Paginator[T],
│                                                CursorPaginator[T], Cursor, Options, Link, Through
├── pipeline/             Pipeline               Middleware[H], Chain[H], Handler[T], Pipe[T], Identity[T]
├── process/              Process                Runner, System, Command, Result, Invoked, ExitError
├── queue/                Queue                  Job, Handler, Queue, Worker, WorkerOptions, Dispatch,
│                                                Option, HandlerOption, Middleware, Locker, Memory,
│                                                o driver de tabela, o Module
├── routing/              Routing                Router, Route, Routes, Group, Resource, Redirect,
│   │                                            Refuse, SignedRoute, FormatRoutes
│   └── middleware/       Routing\Middleware     Throttle, ValidateSignature
├── scheduler/            Console\Scheduling     Scheduler, Task, Cron, Module
├── session/              Session                Store[T], Handler[T], ArrayHandler[T], Start, From[T],
│                                                Flash, CSRF, ErrTokenMismatch
├── socialite/            Socialite              Provider, Manager, Config, User, Token, Endpoint,
│                                                Mapper, o state assinado, os nove drivers
├── str/                  Support\Str            Slug, Snake, Studly, Camel, Kebab, Headline, Plural,
│                                                Singular, Limit, Squish, Mask, Is, UUID, ULID, Random
├── translation/          Translation            Translator, Lines, Loader, T, Choice, Has, Negotiate,
│                                                Middleware, lang/en/*.json embutido
├── validation/           Validation             Rules, MustCompile, Set, Input, Errors, Failure, refused
└── view/                 View                   Renderer, Register, Page, Layout, Text, Escape, JSON,
    │                                            Class, Yield, Include, CSRF, assets, Module
    ├── compilers/        View\Compilers         o compilador kyse: parser, gerador, o conjunto fechado
    └── middleware/       View\Middleware        ShareErrorsFromSession
```

**Fora da coleção, repositórios de adapter** (REGRA 11, nome funcional em inglês, sem tupi):

| Repositório | Componente Illuminate | O que sobra depois da reorganização |
|---|---|---|
| `arandu-io/database` | *(driver de `Database`)* | só os compartimentos `sqlite`, `pgx`, `mysql` e a suíte de conformidade |
| `arandu-io/kv` | **Redis** | `Client` + `cache.Store` + `session.Handler`; sem API própria |
| `arandu-io/queue` | *(driver de `Queue`)* | só o submódulo `kv` |
| `arandu-io/storage` | *(driver de `Filesystem`)* | `local` e `s3`, retipados contra `filesystem.Adapter` |
| `arandu-io/kyse` | *(parte de `View`)* | a biblioteca de componentes, importada (ADR 0027) |
| `arandu-io/ui` | *(equivalente de `laravel/ui`)* | as nove telas + quatro corpos de e-mail, publicados (ADR 0026) |

**Os oito que não viram pacote nenhum:**

| Componente | Por quê | Onde o trabalho aparece |
|---|---|---|
| **Conditionable** | `if` existe em posição de statement em Go | uma linha no guia de migração |
| **Container** | ADR 0001 | `arandu/bootstrap/app.go` — no esqueleto, não no framework |
| **Contracts** | interface vive no pacote que consome | ~30 interfaces espalhadas, seis satisfeitas de fora do módulo |
| **Html** | apagado do Laravel na 5.1; congelado em 2015 | metade vira `view`, metade vira `kyse/components` |
| **Macroable** | o receptor de um método tem de ser declarado no mesmo pacote | nada |
| **Reflection** | é o mecanismo que a tese rejeita | `aru/internal/typeinfo`, em tempo de build sobre `go/ast` |
| **Remote** | apagado do Laravel na 5.1; sucessor fora de escopo | nada |
| **Workbench** | é um gerador de console | `aru make:package` sobre `aru/internal/gen` — **ainda não existe**, é o destino |

---

## 2. O endereço

ADR 0043 deixou três formas em aberto. Uma está morta antes da discussão.

| Forma | Import | Custo |
|---|---|---|
| (a) renomear o módulo para `github.com/arandu-io/hesape` | `arandu-io/hesape/httpx` | **Indisponível.** CLAUDE.md REGRA 2: `github.com/arandu-io/framework` é **irreversível**. O proxy do Go serve o caminho antigo para sempre; todo projeto publicado quebra sem caminho de upgrade |
| (b) subdiretório `framework/hesape/` | `arandu-io/framework/hesape/httpx` | uma reescrita mecânica de import em 11 repositórios; o nome aparece uma vez por arquivo e nunca no call site |
| (c) o nome vive nos documentos e na marca | `arandu-io/framework/httpx` | o nome não aparece em import nenhum; a marca tem de ser carregada por docs, README e site |

### A recomendação é (c)

Três razões, em ordem de peso:

1. **Em Go o nome da coleção quase não é digitado.** Em PHP, `Illuminate` aparece em cada `use` porque o namespace é escrito inteiro na importação. Em Go o caminho é escrito uma vez por arquivo e o que se digita em toda linha é o último segmento: `httpx.Context`, `auth.Grant`, `mail.Envelope`. O ganho de (b) sobre (c) é uma palavra a mais numa linha por arquivo. O objetivo declarado — quem vem do Laravel não deve sofrer — é servido pelo **último segmento**, e ele é idêntico nas duas formas.

2. **`framework/hesape/http` põe dois substantivos de coleção em sequência.** `laravel/framework` → `src/Illuminate/Http` funciona porque `Illuminate` é o vendor namespace do PSR-4 e `laravel/framework` é o repositório. Em Go o *módulo* já faz o papel do vendor namespace. `arandu-io/framework/httpx` é o espelho estrutural exato: repositório-framework, componente. Acrescentar `hesape` é duplicar a função.

3. **ADR 0043 não está fechada e ela mesma diz o que fazer.** Filtros 0 e 2 continuam abertos, a própria ADR se revoga se o filtro 2 reprovar, e a frase é literal: *"o nome muda antes de entrar em qualquer import — que é exatamente por que o endereço ficou pendente"*. Escrever `hesape/` em 11 repositórios agora é apostar numa revisão de tupi antigo que ninguém fez.

**O que (c) custa, sem maquiagem:** `hesape` deixa de existir para quem lê código. Ele vive em `docs/`, no README, no site e no material de marca, e um contribuidor que nunca abre a documentação nunca encontra a palavra. Se o proprietário quiser o nome no caminho, **(b) continua disponível depois** — as fronteiras são o caro, o prefixo é `gofmt -r` mais uma reescrita de import, e ela pode ser um commit sozinho depois que a revisão de tupi liberar.

**Todo o resto deste documento usa `hesape/<componente>` como nome lógico e `github.com/arandu-io/framework/<componente>` como endereço real.**

### Duas divergências de nomenclatura contra a REGRA 7, e por quê

| Nome do Illuminate | Nome no Arandu | Razão |
|---|---|---|
| `Http` | **`httpx`** | `package http` colide com `net/http` em **todo** arquivo: `Context.Response` é um `http.ResponseWriter` e `Middleware` é `func(http.Handler) http.Handler`. O `x` é o que Go escreve quando a stdlib já tem a palavra. O identificador digitado (`httpx.Context`) não muda |
| `Testing` | **`arandutest`** | um segmento de import chamado `testing` sombreia o pacote da stdlib que **todo** `_test.go` importa. É o precedente de `net/http/httptest` |

Uma terceira, menor: `Console\Scheduling` fica em `hesape/scheduler` no topo, e não em `hesape/console/scheduling`. Em PHP o sub-namespace é grátis; em Go um subpacote é uma fronteira real, e `scheduler` já é o nome usado em cinco ADRs. Ver seção 7.

---

## 3. A tabela de mudança

Toda linha de código Go que existe hoje aparece aqui. Um pacote sem linha é um pacote que alguém esquece.

### 3.1 `arandu-io/framework`

| Hoje | Destino | Operação |
|---|---|---|
| `framework/arandutest/` | `hesape/arandutest` | move; ganha `Boot`, `Harness`, `FakeDisk`, `QueueRecorder`, `ExceptionFake`, as assertions de mail/queue/notifications/broadcasting |
| `framework/config/config.go` | `hesape/config` | fica o mecanismo: `App`, `Env`, `LoadDotenv`, os readers, a chave |
| `framework/config/database.go` (321 ln) | `hesape/database` | **move.** No Illuminate isso é `Database\Connectors`, não `Config`. Remove o import `config → data`, que hoje aponta ao contrário |
| `framework/config.Config` (11 campos) | **deixa de existir** | `Database`→`database`, `RedisURL`→`cache`+`session`, `SessionTTL`/`CSRFTTL`→`session`, `LogLevel`/`TracingSecret`/`Editor`→`log`. `kernel.New` muda de assinatura |
| `framework/data/` | `hesape/database` | **renomeia.** `Repository`, `Query`, `DB`, `Dialect`, `Transaction`, `Migration`; ganha `Page[T]`, `ErrNotFound`, `Seeder`, `Register`/`Open` |
| `framework/data/repository.go:159` `Tenant` | `hesape/auth` como `auth.Tenant` | **move.** Uma linha (`g.Subject().Tenant`) que hoje mora no pacote de SQL e obriga `cache` e `filesystem` a importarem o banco para saber um tenant |
| `framework/events/` | `hesape/events` | move; ganha `On(name, Publisher)` e `Listeners(...)` |
| `framework/httpx/context.go`, `state.go`, `reject.go` | `hesape/httpx` | move; ganha `Header`, `Path`, `Method`, `IP`, `BearerToken`, `Cookie`, `IsHTMX`, `WantsJSON`, `User`, `Bind`, `WithCookie`, `WithStatus`, `File` |
| `framework/httpx/router.go:19,22` (`Middleware`, `Chain`) | `hesape/pipeline` | **extrai e generifica.** `httpx.Middleware` vira **alias** (`type Middleware = pipeline.Middleware[http.Handler]`) — o `=` é o que deixa `session`, `cookie` e `exception` produzirem middleware sem importar `httpx` |
| `framework/httpx/router.go`, `name.go`, `resource.go` | `hesape/routing` | move; `Routes.URL`→`Routes.Route`; ganha `Group{Prefix,Name,Middleware}`, `Fallback`, `Match`, `Any`, `SignedRoute` |
| `framework/httpx/validate.go` | **apagado** | `ctx.Validate(*validation.Set)` não tem chamador fora do próprio teste. É a segunda forma de validar (REGRA 9) |
| `framework/httpx/middleware/headers.go` | `hesape/httpx/middleware` | move |
| `framework/httpx/middleware/ratelimit.go` | `hesape/routing/middleware` como `Throttle` | move; `Limiter` e `MemoryLimiter` **apagados** em favor de `cache.RateLimiter` |
| `framework/httpx/middleware/auth.go` | `hesape/auth/middleware` | move; renomeia para os nomes do Illuminate |
| `framework/httpx/middleware/csrf.go` | `hesape/foundation` | move (`Foundation\Http\Middleware\VerifyCsrfToken`) |
| `framework/httpx/middleware/observability.go` | `hesape/foundation` | move |
| `framework/httpx/middleware/recover.go` | `hesape/exception` | move, com o comentário de ordenação intacto |
| `framework/httpx/middleware/flash.go` | **dissolvido** | vira metade de `session.Start`, que carrega **e** consome, em vez de só consumir |
| `framework/internal/` (vazio) | **apagado** | |
| `framework/jobs/` | `hesape/queue` | **renomeia.** `Job`, `Handler`, `Queue`, `Worker`; ganha `Dispatch`, `Option`, `HandlerOption`, `Middleware`, `Locker`, `Hooks` |
| `framework/kernel/` | `hesape/foundation` | **renomeia.** `Kernel`, `Module` e as interfaces opcionais mantêm o nome do tipo; ganha `Builtins() []console.Command` |
| `framework/kernel/module.go:98` `Locker` | **apagado** | um só tipo de lock na coleção: `cache.Lock`. `events.Locker` (alias) some junto |
| `framework/kernel/kernel.go:494` `FormatRoutes` | `hesape/routing` | move; passa a usar `console.IO.Table` |
| `framework/mail/mail.go`, `render.go` | `hesape/mail` | move; `Transport.Send` passa a devolver `(Sent, error)`; `Content` ganha `HTML` |
| `framework/mail/transport.go`, `api.go` | `hesape/mail/transport` | **cinde.** `SMTP`, `Log`, `Array`, `Resend`, `SendGrid`; ganha `Postmark`, `Failover`. `ErrRetryable` fica em `hesape/mail` |
| `framework/modules/auth/{user.entity,user.repo,user.policy,user.service,user.request,module}.go` | `hesape/auth/users` | move menos o reset/confirm |
| `framework/modules/auth/user.service.go:539-862` | `hesape/auth/passwords` | move: `ResetPayload`, `FindForReset`, `ResetPassword`, `ConfirmPassword` |
| `framework/modules/auth/handlers.go` | **apagado** | é o perdedor das duas telas de login (ver 6.4) |
| `framework/modules/` | **deixa de existir** | ADR 0019:127-128 aposentou `modules/` como estrutura do framework e o diretório sobreviveu à decisão |
| `framework/observability/logger.go` | `hesape/log` | move; `NewLogger`→`New`, `WithLogger`→`Into`, `Log`→`For`, `RootLogger`→`Middleware`; ganha `With`, `Field`, `Fields`, `Capture` |
| `framework/observability/{collector,recorder,console,console_page,transport,dump,editor}.go` | `hesape/log` | ficam; `Collector` ganha `RecordCache` e `RecordCommand` |
| `framework/observability/errorpage/` | `hesape/exception` | **cinde.** `page.go` racha em `page.go` + `hints.go`; `Render`/`RenderDump` viram métodos não exportados do `Handler` |
| `framework/scheduler/` | `hesape/scheduler` | move; `Locker` vira `*cache.Lock`; `schedule:*` viram valores `console.Command` |
| `framework/security/policy.go` | `hesape/auth` (raiz) | move inteiro. É a metade Access do `Illuminate\Auth` |
| `framework/security/throttle.go` | `hesape/auth` | move; os três métodos de `SignInThrottle` ganham `ctx`; `MemoryThrottle` **apagado** |
| `framework/security/session.go` | `hesape/session` | move; `Load`→`All`, `Rotate`→`Regenerate`, `Destroy`→`Invalidate`, `IDFromRequest`→`ID` |
| `framework/security/flash.go` | `hesape/session` | move, **menos** `LocalPath` (:250) |
| `framework/security/intended.go` | `hesape/httpx` | move. É `Routing\Redirector::intended()`; `session` nunca valida URL, e `LocalPath` é necessário ao `Back` do `reject.go` |
| `framework/security/csrf.go` | `hesape/session` | move. O token é `Session\Store::token()` no Laravel |
| `framework/security/password.go` | `hesape/hashing` | move; `HashPassword`→`Make`, `VerifyPassword`→`Check`; ganha `Info`, `IsHashed`, leitura de bcrypt |
| `framework/security/signature.go` | `hesape/encryption` | move. **Não** vai para `routing`: `routing → httpx → session → cookie → routing` é ciclo |
| `framework/security/` | **deixa de existir sob esse nome** | |
| `framework/storage/` | `hesape/filesystem` | **renomeia.** `Store`→`Adapter` (sem tenant), novo `Disk` (struct, com o Grant), `Path`→`Key`; ganha `Disks`, `Stat`, `Copy`, `Move`, `DeletePrefix`, `Upload`, `URLSigner`, `Send` |
| `framework/validation/` | `hesape/validation` | fica; `Errors` vira `map[string][]Failure{Field,Rule,Key,Args}`; `Humanize` sai para `str.Headline`; `plural(n,noun)` **apagado** em favor de `str.PluralN`; `uuidShape`/`ulidShape` passam a chamar `str.IsUUID`/`IsULID` |
| `framework/view/` | `hesape/view` | fica; ganha `Escape`, `JSON`, `Class`, `OldValues`, `OldChecked`, `Locale`, `T`, `Choice`. **`assets.go` não cinde** |

### 3.2 `arandu-io/aru`

| Hoje | Destino | Operação |
|---|---|---|
| `aru/internal/kyse/` (~1.600 ln) | `hesape/view/compilers` | **move — condicionado** à emenda de CLAUDE.md:283. Hoje é `internal/` da CLI, então nem o teste do framework nem o build do `kyse` conseguem importar, e `OutputPath:274` carrega um fork só para servir um chamador que não pode importá-lo |
| `aru/internal/doctor/doctor.go:404,454,458,480,496` | `aru/internal/typeinfo` | extrai a metade estrutural (`types`, `exprName`, `callName`, `funcBodyContains`, `receiverType`); as regras ficam |
| `aru/internal/gen/factory.go:155,193` | `aru/internal/typeinfo` | `typeExpr` **apagado** em favor de `RenderType` |
| `aru/internal/spec/spec.go:330` `exportedName` | **apagado** | chama `str.Studly`; o próprio doc comment admite a duplicação |
| `aru/internal/spec/schema.go:19-142` | `hesape/jsonschema` | **condicionado** à emenda da ADR 0006. Só a construção move; `Version`, `Types`, `Actions` e o exemplo ficam |
| `aru/internal/manifest/manifest.go:136` `ReadAll` | **apagado** | zero chamadores; varre `root/modules`, diretório que a ADR 0023 aboliu |
| `aru/internal/manifest` | fica | ganha `Write(dir, Module) error` |
| `aru/internal/fonts/fonts.go:222,441` `Slug` | `hesape/str` | vira `str.Slug`; para de derrubar rune acentuado |
| `aru/internal/gen/{spec,mail}.go` (~80 ln de string) | `hesape/str` | `Normalize`→`Snake`, `exported`→`Studly`, `Kebab`, `Humanise`→`Headline`, `upperFirst`→`UpperFirst`, `Module.Table`→`str.Plural` |
| `aru/key.go:13` `appKeyLen` | **apagado** | importa `encryption.KeySize` e `encryption.GenerateKey` |
| `aru/trace.go:63-66` | `config.String` | `APP_PORT` **apagado**: não existe em config nenhum nem no `.env.example` |
| `aru/commands.go:75,89` | | `delegate("work")`/`delegate("routes")` **apagados**; os nomes do Registry são os do Laravel. Ganha `cache:clear`, `cache:forget`, `config:show`, `queue:failed`, `queue:retry`, `queue:forget`, `queue:flush`, `lang:publish`, `view:clear`, `make:notification`, `make:channel`, `make:test`, `make:package` (este último ainda não existe) |
| `aru/internal/gen/{templates,stubs,job,mail,command,listener,event}.go` | fica | **todos os templates mudam**: escrevem o caminho de import novo dentro do código do usuário |

### 3.3 `arandu-io/arandu` (esqueleto) e `arandu-io/examples`

| Hoje | Destino | Operação |
|---|---|---|
| `arandu/config/*.go` (10) | fica | é `config/*.php`. Perde `env/envBool/envInt/envSeconds` (:109-148) para `hesape/config`; perde `parseMailURL` (:111-178) para `hesape/mail`; `Session.{Cookie,Path,Domain,SameSite}` viram `cookie.Config`; **ganha consumidores** para `Cache`, `Filesystems`, `Queue` e `Session.Driver`, que hoje são lidos por nada |
| `arandu/config/session.go:39-40` `Cookie` | **apagado** | contradiz `session.go:16-19`, que fixa o nome de propósito |
| `arandu/bootstrap/app.go` | fica | ganha wiring de cache store, disks, socialite; perde o `panic` de :245 que manda `go get github.com/arandu-io/mail/<provider>` |
| `arandu/bootstrap/console.go:43-132` | `console.Registry` + `foundation.Builtins` | o `switch` de onze casos, hoje copiado em todo projeto, colapsa em duas linhas |
| `arandu/bootstrap/migrate.go` | `hesape/database/commands.go` | move |
| `arandu/bootstrap/background.go:94` | `queue.RunWorker` | move |
| `arandu/routes/console.go:16-82` | `hesape/console` | o **tipo** `Command` sai do esqueleto. Hoje cada projeto tem um tipo nominal diferente e nada pode ser escrito contra dois |
| `arandu/app/Providers/` | fica | forma inalterada; a documentação para de dizer que ele registra fiação |
| `arandu/database/seeders/seeders.go:35,66,98,120` | `hesape/database` | o contrato move; o registro e `Deps` ficam |
| `arandu/database/factories/` | fica | ganha `database.Store[T,ID]` como metade que grava |
| `arandu/tests/TestCase.go` (~150 ln) | `arandutest.Boot` | encolhe para ~6 linhas; a cópia divergente em `examples` desaparece |
| `arandu/storage/app/{public,private}/` | **apagados** | vazios, gitignorados, lidos e escritos por nada |
| `arandu/resources/views/` | fica | ADR 0021 |
| `examples/config/` | **duplicação a resolver** | é cópia byte a byte de `arandu/config/`; não é decisão deste documento |

### 3.4 Os repositórios de adapter

| Hoje | Destino | Operação |
|---|---|---|
| `kv/cache.go:19-119` | **cinde** | chave, codec JSON, checagem de Grant, regra de TTL e `Remember` vão para `hesape/cache`; só as chamadas RESP ficam, como `CacheStore` |
| `kv/lock.go:29-138` | **cinde** | `Locker`/`Lock`/`Run` viram `cache.Locks`; `SETNX` e o `WATCH/MULTI/EXEC` ficam como `AcquireLock`/`ReleaseLock` |
| `kv/limiter.go` | `hesape/cache` | a aritmética de janela fixa move inteira; **o arquivo desaparece** |
| `kv/session.go` | fica | retipado contra `session.Handler[T]`; ganha `Remembered` e `PasswordConfirmedAt`, hoje descartados (:78-83) |
| `kv/kv.go`, `module.go` | fica | `Client` + health |
| `queue/queue.go`, `queue/module.go` | `hesape/queue` | **move para o core.** Importa só `data`, `jobs`, `security` e a stdlib — zero dependência nova. Precedente: `framework/events/module.go:56-65` já tem tabela no core. **Decisão do proprietário** (ver 7) |
| `queue/kv/kv.go:31-60` | **apagado** | é um **segundo** `redis.NewClient` com esquema de prefixo próprio. Passa a receber o `*kv.Client` |
| `database/database.go` (`Register`, `Open`, `DriverName`) | `hesape/database` | **move para o core**, mesma razão. **Decisão do proprietário** |
| `database/{sqlite,pgx,mysql}/` | fica | intocados |
| `storage/local.go`, `storage/s3/s3.go` | fica | retipados contra `filesystem.Adapter`; `s3` ganha `PresignGet` |
| `mcp/schema.go:20-152` | `hesape/jsonschema` | o builder e o `Validate`; `mcp/protocol.go:92` passa a embutir `json.RawMessage` e `mcp/transport.go:136` a chamar `ObjectType.Properties()` |
| `kyse/components/` | fica | ganha `Select`, `Checkbox`, `RadioGroup`, `File`, `Pagination` |
| `ui/views_auth_flow.go` | **cinde** | purposes, TTLs, assinatura e construção de e-mail vão para `hesape/auth/passwords` e `hesape/auth/notifications`. Controllers, telas e corpos de e-mail **ficam publicados** (ADR 0022:95-98) |

---

## 4. O que precisa ser construído

Ordenado por dependência (camada), depois por quanto uma aplicação toca por dia.

### Camada 0 — folhas, não importam nada de `hesape`

| Item | Assinatura | Esforço |
|---|---|---|
| `pipeline` | `type Middleware[H any] func(H) H` / `func Chain[H any](h H, mws ...Middleware[H]) H` | horas |
| `str` — 40 funções | `Slug`, `SlugWith`, `ASCII`, `Snake`, `Studly`, `Camel`, `Kebab`, `UpperFirst`, `Headline`, `Title`, `Plural`, `PluralN`, `Singular`, `Limit`, `LimitWords`, `Words`, `Squish`, `After`, `Before`, `Between`, `Start`, `Finish`, `Mask`, `Is`, `UUID`, `UUID7`, `ULID`, `Random`, `IsUUID`, `IsULID` | dias |
| `collections` — 10 funções | `Map[T,U]`, `Filter[T]`, `Reduce[T,A]`, `First[T]`, `Last[T]`, `GroupBy[T,K]`, `KeyBy[T,K]`, `Partition[T]`, `UniqueBy[T,K]`, `Sum[T,N]` | horas |
| `number` | `FileSize`, `FileSizeBinary`, `Abbreviate`, `Currency`, `Percentage`, `Ordinal`, `Format` | dias |
| `encryption` — o `Signer` movido | `func (s *Signer) Sign(purpose, payload string, ttl time.Duration) string` / `Verify` / `ParseKey` / `GenerateKey` / `KeySize` | horas |
| `pagination` | `LengthAwarePaginator[T]`, `Paginator[T]`, `CursorPaginator[T]`, `Cursor`, `Options`, `Link`, `Through[A,B]`, `CurrentPage`, `OptionsFrom` | semanas |
| `jsonschema` | `Type`, `Object`, `Prop`, `String/Integer/Number/Boolean/Array`, `Union`, `AnyOf`, `Document`, `Validate(schema Type, value any) error` | dias |
| `concurrency` | `func Run[T any](ctx, tasks ...Task[T]) ([]T, error)` / `func Recover(ctx, name string, fn func() error) (err error)` / `PanicError` | dias |
| `process` | `Runner`, `System`, `Command{Name,Args,Dir,Env,Stdin,Timeout,IdleTimeout,MaxOutput,OnOutput}`, `Result`, `Invoked`, `ExitError`, `TimeoutError` | dias |

`concurrency` é o único da lista cuja razão de existir não é ergonomia: `recover()` aparece **uma vez** em todo o framework (`httpx/middleware/recover.go:34`) e cobre só HTTP. Um handler de job, uma tarefa agendada ou uma passada do relay que dê panic derruba o processo.

### Camada 1 — sobre as folhas

| Item | Assinatura | Esforço |
|---|---|---|
| `auth.Allows` | `func Allows[T any](ctx, p Policy[T], s Subject, a Action, resource T) bool` | horas |
| `auth` no contexto | `func WithSubject(ctx, s Subject) context.Context` / `func SubjectFrom(ctx) (Subject, bool)` / `func Check(ctx) bool` | dias |
| `hashing` — leitura de bcrypt | `func Check(plain, encoded string) error` aceitando `$2a$`/`$2b$`/`$2y$`/`$argon2i$` | horas |
| `hashing.Info` / `IsHashed` | `func Info(encoded string) (Params, bool)` / `func IsHashed(v string) bool` | horas |
| `cookie` | `Jar`, `Config{Path,Domain,Secure,SameSite,Partitioned,Unsigned}`, `Make`, `Forever`, `Forget`, `Read`, `Write` | dias |
| `log` — `With` colapsado | `func With(ctx, args ...any) context.Context` / `Field` / `Fields` / `Restore` | dias |
| `log.Capture` | `func Capture() (*slog.Logger, *Records)` | horas |
| `config` — os buracos | `MustString`, `App.Debug`, `App.URL *url.URL`, `App.Timezone *time.Location`, `App.PreviousKeys [][]byte`, `Env.Is/IsProduction` | horas |

### Camada 2 — sobre a camada 1

| Item | Assinatura | Esforço |
|---|---|---|
| `session.Start` + `From[T]` | `func Start[T any](s *Store[T]) pipeline.Middleware[http.Handler]` / `func From[T any](ctx) (T, bool)` — **hoje todo controller chama `sessions.Load` à mão, e o `aru` gera essa chamada** | dias |
| `session.Handler[T]` genérico | `Read`/`Write`/`Destroy`/`DestroyIndex` | dias |
| `cache.Store` + `Repository` | `Get[T]`, `Put`, `Add`, `Forget`, `Has`, `Increment`, `Pull[T]`, `Remember[T]`, `GetMany[T]`, `PutMany`, `Flush`, `Namespace` | semanas |
| `cache.ArrayStore` | `func NewArrayStore() *ArrayStore` — anunciado em `arandu/config/cache.go:18` e implementado em lugar nenhum | dias |
| `cache.Locks` | `func (l *Locks) Lock(name string, ttl time.Duration) *Lock` — **o único tipo de lock da coleção** | dias |
| `cache.RateLimiter` + `Limit` | `Attempt`, `Hit`, `Release`, `Clear`, `AvailableIn`; `PerMinute`/`PerHour`/`PerDay`/`By` | semanas |
| `database.ErrNotFound` | `var ErrNotFound = errors.New("database: no such row")` — muda `Repository[T,ID].Find`, a interface que todo módulo gerado implementa | dias |
| `database.Page[T]` | `List` devolve `Page[T]{Items, Next}` em vez de `[]T` | dias |
| `database` — migrations | `Pretend`, `WithMigrationLock`, `Fresh`, `Wipe`, `HasTable`, `HasColumn`, `Show`, `Table` | dias |
| `database` — erros de driver | `IsUniqueViolation`, `IsForeignKeyViolation`, `IsDeadlock` — hoje é `strings.Contains` copiado em todo repositório gerado | horas |

### Camada 3 — HTTP, view, validação

| Item | Assinatura | Esforço |
|---|---|---|
| `exception.Abort` + `HTTPError` + `classify` | `func Abort(status int, message string) error` / `func StatusOf(err error) (int, bool)` — **hoje todo erro que sai de um handler vira 500** (`router.go:195` dá panic) | horas |
| `exception.Handler` | `NewHandler(Config)`, `Report`, `Render`, `RenderForConsole`, `Renderable`, `Reportable`, `DontReport`, `Level`, `ShouldRenderJSONWhen` | dias |
| páginas de erro de produção | `func RenderStatus(w, r, status int, message string)` + `html/template` para 401, 403, 404, 419, 429, 500, 503 — **hoje não existe nenhuma; produção responde `text/plain`** | dias |
| `httpx.Context` — os acessores | `Header`, `Path`, `FullURL`, `Method`, `IP`, `BearerToken`, `Cookie`, `IsHTMX`, `WantsJSON`, `NotFound`, `RedirectRoute` | horas |
| `httpx.Context.User` | `func (c *Context) User() (auth.Subject, bool)` | dias |
| `httpx.Context.Bind` | `func (c *Context) Bind(dst any) error` — hoje cada controller preenche a struct à mão | dias |
| `httpx.Context.WithStatus` | `func (c *Context) WithStatus(message string) *Context` + `Page.Status()` — a mensagem de sucesso que sobrevive a um POST-redirect-GET; **não existe** | dias |
| `httpx.File` + `filesystem.Upload` | `func File(c *Context, field string) (filesystem.Upload, error)` — nada em `framework/` faz parse de multipart | dias |
| `middleware.TrustProxies` | `func TrustProxies(trusted []netip.Prefix) httpx.Middleware` — **urgente**: `ratelimit.go:58` lê `RemoteAddr` e nunca `X-Forwarded-For`, então atrás de proxy o mundo inteiro divide uma chave de rate limit |
| `middleware.LimitBodySize` | `func LimitBodySize(max int64) httpx.Middleware` — `kernel.go:386` monta o servidor sem limite nenhum de corpo | horas |
| `middleware.CORS` | `func CORS(o CORSOptions) httpx.Middleware` | dias |
| `filesystem.Disk` + `Adapter` | `Put`, `Get`, `Stat`, `Exists`, `Delete`, `Copy`, `Move`, `List`, `DeletePrefix` | dias |
| `filesystem.Send` | `func Send(w, r, g auth.Grant, d *Disk, key string, opt SendOptions) error` — **nada no repositório serve um arquivo guardado por HTTP hoje** | dias |
| `filesystem.URLSigner` | `TemporaryURL`, `Redeem`, `Presigner` | dias |
| `filesystem.UploadRules` | `func (u Upload) Check(r UploadRules) error` — o "contrato próprio" que `validation/refused.go:46` promete e ninguém implementou | dias |
| `view.Escape` / `JSON` / `Class` | `func Escape(v any) string` / `func JSON(v any) template.JS` / `func Class(pairs ...any) string` | horas |
| `view.Page` — os buracos | `OldValues`, `OldChecked`, `Locale`, `T`, `Choice`, `Status` | horas |
| `validation.Errors` → `[]Failure` | `type Failure struct{ Field, Rule, Key string; Args []any }` | dias |
| `translation` | `Translator`, `Lines`, `T`, `Choice`, `Has`, `Loader`, `NewFS`, `NewMap`, `Negotiate`, `Middleware`, `WithPlural`, `WithNamespace`, `WithMissingKey` + `lang/en/{validation,auth,passwords,pagination}.json` | semanas |

### Camada 4 — console, fila, e-mail, eventos

| Item | Assinatura | Esforço |
|---|---|---|
| `console.Command` no framework | `type Command struct{ Name, Description string; Hidden bool; Isolated string; Run func(ctx, *IO) error }` — **hoje o tipo é declarado no esqueleto** | horas |
| `console.IO` | `Info`, `Line`, `Comment`, `Warn`, `Error`, `NewLine`, `Args`, `Flags`, `Confirm`, `Ask`, `Secret`, `Choice`, `Table`, `TwoColumnDetail`, `Task`, `Progress` — **não existe nada; o stub gerado imprime com `fmt.Println`** | dias |
| `console.Registry` | `New`, `Add`, `Handle`, `Names`, `Help`, `Call`, `CallSilent`, `Observe`, `WithLocks` | dias |
| `console.ExitError` / `GuardEnv` | `func Exit(code int, format string, a ...any) error` / `func GuardEnv(env config.Env, want config.Env, action string) error` | horas |
| `queue` — opções por job | `On`, `After`, `At`, `Unique`; `MaxAttempts`, `Backoff`, `Timeout` | dias |
| `queue` — controle no handler | `func Release(after time.Duration) error` / `func Fail(cause error) error` / `var ErrDiscard` / `type Failing interface{ Failed(...) }` | horas |
| `queue` — middleware | `WithoutOverlapping`, `RateLimited`, `ThrottlesExceptions` | dias |
| `queue` — dead letter | `Forget`, `Flush` + `queue:failed`/`queue:retry`/`queue:forget`/`queue:flush` | dias |
| `queue.Memory` | driver `sync` e fake de teste **num tipo só** | dias |
| `mail` — anexos | `Attachment{Name,ContentType,Inline,Path,Bytes,Key}` + `Resolve(ctx, g, s)` | dias |
| `mail.Sent` | `Transport.Send` devolve `(Sent, error)` — hoje o message-id do provedor é descartado | horas |
| `mail.Pending.Queue` | `func (p *Pending) Queue(ctx, g, q queue.Queue, m Mailable, opts ...queue.Option) error` | dias |
| `mail.Render` / `Headers` / `Raw` / `Postmark` / `Failover` | | dias |
| `events` — fan-out | `func On(name string, h Publisher) Publisher` / `func Listeners(list ...Publisher) Publisher` — `aru make:listener` já supõe que existe | horas |
| `events` — assertions | `Collected.AssertPublished` / `AssertNotPublished` / `AssertPublishedCount` | horas |

### Camada 5 — o que fica em cima de tudo

| Item | Assinatura | Esforço |
|---|---|---|
| `auth/passwords.Broker` | `SendResetLink`, `Reset`, `ConfirmPassword` — o mecanismo sai do `ui` e vira código do framework; controller e telas ficam publicados | semanas |
| `auth/notifications` | `SendVerification`, `SendReset` (construtores; os corpos ficam no `ui`) | semanas |
| `auth.Provider` | a costura que deixa uma aplicação autenticar contra a própria tabela | semanas |
| `bus` | `Batch`, `PendingBatch{Add,Then,Catch,Finally,AllowFailures,OnQueue,Dispatch}`, `Chain`, `Store`, `Find`, `Cancel`, `Batching`, `PushUnique` | semanas |
| `notifications` | `Notification`, `Notifiable`, `Key`, `Notifier{Send,SendMany,Suppress}`, `Channel`, `Store`, `Record`, `Policy`, `Anonymous`, `Capture`, a tabela | semanas |
| `broadcasting` | `Hub`, `Channel`, `Private`, `Subscribe`, `Fanout`, `Bus`, `Message`, `ConnID`, o stream SSE, o `Module` | semanas |
| `socialite` | `Provider`, `Manager`, `Config`, `User`, `Token`, `Endpoint`, `Mapper`, o state assinado, nove drivers, `Fake` | semanas |
| `image` | `Image`, `Driver`, `Pipeline`, `Transformation`, `Std`, `Cover`, `Orient`, `Store`, `Respond`, limites de decodificação | semanas |
| `arandutest` | `Boot`, `Harness`, `Application`, `AssertDatabaseHas/Missing/Count/Empty`, `Mailbox`, `QueueRecorder`, `MemoryStore`, `View`, `Artisan`, `Client.ActingAs`, `AssertAuthenticated`, `Response.Invalid` | semanas |
| `aru make:package` (ainda não existe) | `func GeneratePackage(p Package) ([]File, error)` + `manifest.Write` | dias |

---

## 5. O que fica de fora

**Rejeitado por decisão de arquitetura já registrada** — ADR 0001 (container), ADR 0002 (facades), `docs/01-arquitetura.md:136-139`:

> `app()`, `resolve()`, `Container::getInstance`, `bind`, `singleton`, `scoped`, `extend`, `tag`, `alias`, `resolving`, `flush`, contextual binding, atributos `#[Singleton]`/`#[Bind]`/`#[Config]`/`#[CurrentUser]`, `ServiceProvider::register`, method injection, route model binding, todas as 44 facades, `Cache::`, `Mail::`, `Storage::`, `Artisan::`, `Socialite::`, `Notification::`, `Broadcast::`, `Lang::`, `Config::get/set`, `cache()`, `logger()`, `env()` fora de `config/`, `Eloquent`, `$fillable`, `hasMany`/`belongsTo`/`with()`, `loadMorph`, model events, observers, casts, scopes.

**Rejeitado pela REGRA 9 / `docs/09-uma-forma-so.md:56`** (segunda maneira de fazer o que já tem uma):

> `Str::of()` e todo o `Stringable`; `->when()`/`->unless()`/`Conditionable` em qualquer forma; `abort_if`/`throw_if`/`rescue`; `Cache::store()`/`driver()`/`extend()`/`tags()`; `Mail::mailer()`; `dispatch()->onConnection()`; `Broadcast::driver()`/`extend()`; `Session::get/put` por chave; `Cookie::queue()` e as seis funções da fila; `$this->components->info()` (o segundo conjunto de saída do artisan); `PromptsForMissingInput`; `->only()`/`->except()` em `Route::resource`; `Kernel::$middlewarePriority`; `setScopes` ao lado de `scopes`; `:Attribute`/`:ATTRIBUTE`; `lang/en.json` (o segundo formato de catálogo); `->stateless()`; `->redirectUrl()`; `session.lottery`; `queue:listen`; `sendmail`.

**Impossível em Go, não indesejado:**

> `Macroable` (o receptor tem de ser declarado no mesmo pacote); `HigherOrderWhenProxy` (`__get`/`__call`); `Form::model` (`object_get`); `$c->map->name`; `ArrayAccess`/`Iterator`/`Countable` no paginador; `Blade::render` de string em runtime; `SerializesModels`; closures em `->then()`/`->catch()` de batch; `$loop` sem mudar o conjunto fechado de diretivas.

**Rejeitado pela REGRA 13 (zero Node) / `docs/09-uma-forma-so.md:48`:**

> `laravel-echo` e `pusher-js`; `install:broadcasting`; markdown mailables e `Markdown::theme`; os nove Blade de paginação (`useBootstrapFive`, `useTailwind`); `@vite` como manifesto; qualquer `package.json`.

**Rejeitado pela REGRA 2 / REGRA 17 (nada alcança dado sem Grant, inclusive leitura):**

> `Rule::unique` e `Rule::exists`; `Storage::url()` e visibilidade pública; `storage:link`; `->create()` em factory; `$request->file()->store()` como one-liner; `Broadcast::channel()` como segundo caminho de autorização; middleware `can:` (ADR 0039:41-44).

**Rejeitado por ADR 0004 (core sem dependência):**

> `x/text` para plurais CLDR e formatação de moeda; `doctrine/inflector`; `voku/portable-ascii`; `league/commonmark`; `nesbot/carbon`; `opis/json-schema`; `goth`; `markbates/goth`; OAuth 1.0a; pty (`Process::tty`); SES.

**Apagado do próprio Laravel** — `Html` (5.1), `Remote` (5.1), `Workbench` (5.1). Especificar qualquer um seria divergir *do* Laravel usando um nome que diz o contrário.

**Deferido com gatilho escrito, não excluído:**

| Item | Gatilho | Fonte |
|---|---|---|
| WebSocket (`joaju`) | quando o SSE do HTMX não bastar | `docs/06-escopo.md:150` |
| Reverb / protocolo Pusher | idem | `docs/packages/reverb.md:11` |
| `jsonschema.Parse` | o primeiro consumidor que precise ler um schema que não escreveu | — |
| tradução (o pacote existe; os idiomas não) | primeiro cliente fora do par pt-BR/en | `docs/13:57` |

---

## 6. A ordem de execução

Nove fases. Cada uma termina com `gofmt -l $(find . -name '*.go' -not -name '*.kyse.go') && go vet ./... && go test -race ./...` verde nos **doze** repositórios que têm código Go. Os outros dois dos catorze não entram no portão: `homebrew-tap` não tem arquivo `.go` e `docs` é markdown.

*Corrigido em 14/08/2026:* esta linha dizia onze, que era o número em 10/08/2026, quando o documento foi escrito. O `hesape` virou repositório próprio pela ADR 0049 e entra no portão como os outros.

### Fase 0 — Consertos independentes. Nada se move.

| Conserto | Onde | Por quê agora |
|---|---|---|
| `ui` não compila | `ui/views.go` — 11 chamadas passam `Error: .X` para `components.FieldProps`, que não tem campo `Error` (`kyse/components/field.kyse.go:24-50`) | `go run github.com/arandu-io/ui@latest auth` — o **primeiro** comando que um dev Laravel roda — publica nove telas que não buildam. `ui/views_test.go:64-70` adia o compile check, então os testes do módulo estão verdes |
| `view/page.go` (~:272) manda usar `Error:` | doc comment publicado no pkg.go.dev | é a instrução que produziu as 11 chamadas erradas |
| keyset com `created_at` fixo | `framework/modules/auth/user.repo.go:280-286` calcula o limite em `created_at` e ordena por `column` da allowlist, que aceita `"email"` | páginas que pulam e páginas que repetem, em silêncio, no módulo de auth. `aru/internal/gen/templates.go:261-266` documenta esse mesmo bug como **já corrigido** no gerador — a correção chegou ao template e nunca ao módulo |
| `panic` apontando repositório inexistente | `arandu/bootstrap/app.go:245` manda `go get github.com/arandu-io/mail/<provider>` | `examples` foi corrigido; o esqueleto que o usuário clona não |
| `MAIL_URL=smtps://` disca em claro | `arandu/config/mail.go:65-69` faz parse de `Encryption` e `mail.SMTP` não tem o campo | URL documentada que não funciona |
| segundo cliente RESP | `queue/kv/kv.go:31-60` | duas formas de conectar em RESP, com layouts de chave diferentes |
| texto da REGRA 14 | CLAUDE.md:325-326 | diz que sessão, rate limit e lock de scheduler são prefixados por tenant. `kv/session.go:31-34`, `kv/limiter.go:74-79` e `kv/lock.go:129-133` dizem o contrário, com razão, e `docs/15:100-105` documenta as três exceções. **Essa frase mandou dois agentes desta auditoria pelo caminho errado** |

### Fase 1 — Decisões e ADRs. Nada se move.

| ADR | Fecha |
|---|---|
| **0044 — o que é recusado e nunca foi escrito** | `Cache::forever` (recusado em `kv/cache.go:121` sem ADR), `Storage::append/prepend`, `file`/`database` cache store, `sendmail`, `queue:listen`, `queue:pause/resume`, `ShouldBeEncrypted`, `EncryptedPrivateChannel`, `@method`, `@push`/`@stack`, `@switch`, `@includeIf`, `Html::obfuscate`, `Form::selectRange`, `session.block`, `Gate::define`/`Gate::policy` (o registro; a facade já está em 0002) |
| **0045 — revoga o `goth`** | `docs/packages/socialite.md:12` manda um adapter sobre `markbates/goth`; goth traz `gorilla/sessions`, que é um segundo store de sessão. ADR 0031 já resolveu o caso idêntico (Resend/SendGrid) a favor do core sobre `net/http` |
| **0046 — emenda a ADR 0006** | a `aru` pode importar `github.com/arandu-io/framework`. Necessário para consolidar `jsonschema`; o texto atual diz *"módulo próprio, sem dependência"*, e a razão declarada é unidirecional (projetos não devem pagar pela CLI), o que não é violado |
| **0047 — emenda CLAUDE.md:283** | "kyse é parte da `aru`" deixa de ser verdade quando o compilador vira `hesape/view/compilers`. Todas as propriedades da REGRA 13 sobrevivem (sem download, sem npm, sem Node, código compilado no binário); a frase não |
| **0048 — a exceção da URL pré-assinada** | `filesystem.Presigner` devolve um 302 para uma capability bearer resgatada sem Policy. É exceção real à REGRA 17. Se recusada, o argumento de egresso do R2 em `docs/11:90-93` tem de sair, porque está justificando uma economia que a arquitetura declina |
| **0049 — a tabela `users` é do framework ou do projeto** | `arandu/app/Models/User.go:29` é `type User = auth.User`, alias para struct de oito campos fixos. No Laravel a migration está em `laravel/laravel` e `App\Models\User` é classe real. Um dev Laravel adiciona uma coluna no dia um e aqui não consegue. **Maior divergência não resolvida do projeto** |
| **ADR 0043 — os filtros** | revisão de tupi antigo (filtro 2), leitura em voz alta (filtro 0), leitura em inglês (filtro 3). Trava a forma (b) do endereço, não trava nada mais |
| **Correções de documentação** | `docs/13` (linhas 22, 24, 25, 26, 29, 30, 33, 35, 37, 38, 40, 42, 43, 45, 46, 47, 48, 50, 53, 54, 55, 57, 59, 72-73, 82, 90-105, 124, 126); `docs/10` (57, 58, 59, 60, 61, 94, 95, 97, 113, 123, 134, e as ausências: auth, session, cookie, broadcasting, notifications, pagination, Str, translation, Concurrency, Bus); `docs/30` (70-76, 161, 166, 192-200, 210); `docs/05` (94-96, 140-141); `docs/07` (37, 90, 94, 112); `docs/11` (20-29, 42, 45, 55, 167); `docs/12` (55); `docs/28` (122, 131, 157-160, 204, 365-379); `docs/adr/0001` (:12); `docs/adr/0009` (:35, :41-42); `docs/adr/0022` (:16-28); `docs/02` (9-11); `docs/06` (57) |

### Fase 2 — As folhas. Aditivo. **Não quebra API publicada.**

`pipeline`, `str`, `number`, `collections`, `jsonschema`, `pagination`, `concurrency`, `process`.

A única extração é `Chain` de `httpx/router.go:22` — 21 chamadas dentro de `framework/`. `httpx.Middleware` vira alias, então nenhum arquivo de aplicação muda de forma. Atenção: um `[]func(http.Handler) http.Handler` **não** é atribuível a `[]httpx.Middleware`, e já não era; o teste do stub (`aru/internal/gen/stubs.go:307`) tem de cobrir o nível de slice, não só o de elemento.

`jsonschema` fica bloqueado até a ADR 0046. Sem ela, consolida-se só o `mcp` e a duplicação com `aru` fica de pé.

### Fase 3 — A cisão de `security`. **Quebra API publicada em dez repositórios.**

A maior e a mais atômica. `framework/security` racha em cinco; `security.Grant` vira `auth.Grant`.

| Custo | Número |
|---|---|
| arquivos que importam `framework/security` | **129** |
| referências a `security.Grant` | **215** |
| repositórios afetados | framework, aru, arandu, examples, ui, database, kv, queue, storage, mcp |
| repositórios publicados que quebram | framework, arandu, examples, ui, kv, queue, storage, database, mcp |

Junto, nesta fase, porque não separa:

- `data.Tenant` → `auth.Tenant`. Sem isso, `cache` e `filesystem` importam o pacote de SQL para saber um tenant.
- `kernel.Locker` e `events.Locker` apagados (o lock passa a ser `cache.Lock`, mas o *tipo* já pode ser declarado aqui e implementado na fase 5).
- `SignInThrottle` ganha `ctx` nos três métodos, **antes** de `MemoryThrottle` ser apagado.
- `session.go:16-19` vence `arandu/config/session.go:68`: `SESSION_COOKIE` deixa de ser lido e o campo é apagado.
- `aru/internal/doctor/rules.go:1110` casa com o identificador `Rotate`. Renomear para `Regenerate` sem editar essa regra **desliga em silêncio a checagem de session fixation** — é a única regra do doctor cuja falha é invisível.

Os geradores (`aru/internal/gen/templates.go`, `stubs.go`) escrevem `security.` dentro do código do usuário. Todo projeto já gerado precisa de reescrita mecânica documentada em `UPGRADE.md`. **Não** deixar um pacote `framework/security` de fachada reexportando: seria um segundo caminho de import (REGRA 9).

### Fase 4 — A cisão de `httpx` e de `observability`. **Quebra.**

`httpx` racha em quatro (httpx, routing, pipeline-adoção, e `validate.go` apagado); `httpx/middleware` racha em cinco; `observability` racha em dois; `kernel` vira `foundation`.

| Custo | Número |
|---|---|
| arquivos que importam `framework/httpx` | **71** |
| repositórios | framework(31), examples(18), aru(8), arandu(8), ui(2), mcp(1) |

Ordem interna: `pipeline` (fase 2) → alias `httpx.Middleware` → `ExceptionHandler` declarado em `httpx` e implementado em `exception` → `router.go:195` deixa de dar panic. **A partir do momento em que o classificador entra, `auth.ErrForbidden` responde 403 e `database.ErrNotFound` responde 404. Todo teste que hoje afirma 500 muda.** Não dá para faseá-lo: um classificador pela metade é pior que nenhum.

Efeito colateral obrigatório: `docs/00:55`, `02`, `03:26`, `03:214`, `06:57`, `09:88`, `12:41`, `14:393`, `25:54` e `27:137` nomeiam `observability/errorpage`. Todos ficam falsos no instante do move.

### Fase 5 — `config`, `cache`, `kv` e `filesystem`. **Quebra.**

`framework/config` racha em três; `kv` racha na linha do contrato; `storage` vira `filesystem`.

- `framework/config.Config` (11 campos) deixa de existir → `kernel.New` muda de assinatura, e é o construtor mais chamado do esqueleto e das duas suítes de teste.
- `kv` perde `Cache`, `Locker` e `Limiter`: quatro tipos exportados, doze métodos.
- `middleware.RateLimit` muda de assinatura; `MemoryThrottle` sai.
- `storage.Store` → `filesystem.Adapter`; `Path` → `Key`; novo `Disk`. Dois módulos publicados (`storage` e `storage/s3`), duas tags coordenadas.
- `arandu/config/{cache,filesystems,queue}.go` ganham consumidor pela primeira vez.

**Três coisas escrevem em `arandu/bootstrap/app.go` nesta fase.** Sequenciar, ou aceitar três reescritas conflitantes de uma função de 150 linhas.

**Risco de segurança específico:** `security/throttle.go:58-95` documenta uma corrida que a primeira implementação errou — o orçamento tem de ser tomado **antes** da senha ser checada, sob um lock só. Reexprimir isso como `Hit`/`Release`/`Clear` sobre um `Store` reabre a corrida através de um round trip de rede, onde o lock deixou de ser um mutex. Testes antes de apagar o código antigo, não depois.

### Fase 6 — `queue`, `database` e os adapters entram no core. **Quebra caminho de módulo publicado.**

`framework/jobs` → `hesape/queue`. `framework/data` → `hesape/database`. `queue/queue.go`+`module.go` e `database/database.go` sobem para o core.

- `jobs.Job` ganha `BatchID` e `Chain` — muda a struct, a tabela (`queue/module.go:69-93`) e o `Push`/`Reserve` de **ambos** os drivers. REGRA 16: colunas nulláveis com default, e um driver que descarte os campos em silêncio perde a cadeia de jobs permanentemente. **Os dois drivers na mesma release.**
- `arandu-io/queue` fica com um módulo só, num subdiretório. `arandu-io/database` idem. É mudança quebrante de caminho publicado.
- O comentário `framework/jobs/jobs.go:6-9` (*"a driver is a separate module... a core that carried a Redis client would put it in every project's go.sum"*) tem de ser reescrito no mesmo commit, ou vira afirmação falsa publicada no pkg.go.dev. A razão vale para o driver **kv**; não vale para o de tabela, que importa só stdlib e o core.
- **Colisão de nome a decidir aqui:** `hesape/queue` e `arandu-io/queue/kv`; `hesape/database` e `arandu-io/database/pgx`. Todo arquivo de wiring que importe os dois precisa de alias.

### Fase 7 — `console`, o compilador da view, `translation`, `validation.Errors`. **Quebra.**

- O tipo `Command` sai do esqueleto: `routes.Command` deixa de ser atribuível a um campo `console.Command`. `Run` muda de `func(ctx, []string) error` para `func(ctx, *console.IO) error`. **Todo comando já gerado para de compilar.**
- `aru help` dentro de um projeto passa a delegar para listar os comandos da aplicação — o que faz um comando instantâneo passar a depender de `go run`.
- `aru/internal/kyse` → `hesape/view/compilers` (condicionado à ADR 0047). Todo `.go` gerado por `aru view:build` e os dez arquivos gerados de `kyse/components` mudam de header de import; a árvore compilada inteira é regerada e `kyse` re-tagueado.
- `validation.Errors` vira `map[string][]Failure`: quebra `httpx/validate.go`, `view/page.go`, `security/flash.go`, todo controller e view gerados, e qualquer aplicação que escreva `errs["email"][0]`. **Entra antes da 1.0 ou nunca entra.**
- `Page.FieldError`, `FieldErrors`, `HasErrors`, `ErrorSummary` e `OldOr` **mantêm a assinatura exata**: `kyse/components/page.go:31` declara `Page interface{ FieldError(name string) string; OldOr(...) string }`, e é essa interface que faz `arandu-io/kyse` não importar nada do framework.

### Fase 8 — Construir o que falta. Aditivo.

`arandutest` primeiro (tudo o mais precisa de teste), depois `bus`, `notifications`, `socialite`, `image`, `broadcasting`.

`broadcasting` é o último e tem dois bloqueios anteriores: ~~`kv` precisa ganhar `Publish`/`Subscribe` (não tem — verificado, zero
ocorrências)~~ **Falso, corrigido em 14/08/2026:** `hesape/redis/connections`
tem os três — `Publish` (`commands.go:506`), `Subscribe` e `PSubscribe`
(`subscriptions.go:27,37`), e a decisão HTTP/1.1 vs HTTP/2 tem de ser tomada **antes da primeira linha**. `framework/kernel/reload.go:15-50` é um post-mortem escrito neste repositório desse mesmo mecanismo quebrando: um `EventSource` aberto gasta uma das seis conexões que o HTTP/1.1 dá por origem, e `aru dev` serve HTTP/1.1 em texto claro.

---

## 7. O que este documento não resolve

### 7.1 Onde as especificações discordaram e o revisor escolheu

| Questão | Vencedor | Perdedor | Razão |
|---|---|---|---|
| **`Bus` é pacote próprio?** | `hesape/bus`, separado de `hesape/queue` | a revisão de área "work", que atribuiu batch/chain a `hesape/queue` de passagem | o Laravel publica `illuminate/bus` separado de `illuminate/queue`; a fronteira é *exigível* (`bus` compõe `jobs.Job` e nunca toca a tabela; `queue` chama um `BatchRecorder` declarado no próprio pacote), e o dev digita `Bus::batch`. **A revisão de "work" não tinha a spec de `bus` à frente** |
| **Onde fica o `Signer`?** | `hesape/encryption` | a revisão de "security", que mandou para `hesape/routing` | `routing → httpx → session → cookie → routing` é ciclo. O `Signer` é primitiva de folha. `hesape/encryption` guarda-o e diz no doc comment que não há `Crypt`; a metade de *construção* de URL assinada (`SignedRoute`, `ValidateSignature`) vai mesmo para `routing`, como no Illuminate |
| **`intended.go`** | `hesape/httpx` | a revisão de "security", que mandou para `hesape/session` | `Routing\Redirector::intended()` é o endereço no Laravel; `session` nunca valida URL; `LocalPath` é preciso ao `Back` do `reject.go`, que está em `httpx` |
| **`State`, `Reject`, `Back`** | `hesape/httpx` | a spec de `session` | `illuminate/http` **requer** `illuminate/session`, e `illuminate/session` não requer nenhum dos dois. `withInput`/`withErrors` são `Http\RedirectResponse`. O inverso é ciclo |
| **Renomear `arandu-io/kv` para `redis`** | `kv` fica | a spec de `redis` | o adapter está **fora** de `hesape` e não é governado pela REGRA 7; é governado pela REGRA 11 (nome funcional) e por `docs/11:60-66`, que escolheu `kv` porque é um protocolo e quatro produtos. `Illuminate\Redis` mapeia para `arandu-io/kv`, não para um pacote do core |
| **Contratos de paginação** (`Paginated`, `LengthAware`) | **OUT** | a spec de `pagination`, que os pôs no pacote produtor | `docs/13:32` e `:123`: a interface vive no pacote que **consome**. Zero consumidores existem, e o componente kyse recebe `[]Link`, não interface |
| **Cindir `view/assets`** | não cinde | a spec de `view` | o Illuminate não tem `View\Assets` (o equivalente de Vite é `Foundation\Vite`), e `view.URL("app.css")` está dentro de `{{ }}` em todo layout. Contagem de linhas não é fronteira |
| **Endereço do `hesape`** | (c) — o nome nos documentos | (b) — `framework/hesape/` | ver seção 2. (b) continua disponível depois, ao custo de um commit mecânico |
| **`SESSION_COOKIE` configurável** | fixo (o framework vence) | `arandu/config/session.go:68` | `session.go:16-19` documenta que nome configurável quebra a ligação do CSRF. Reverter reabre `session.go:449-457` e é outra decisão |
| **Prefixo de tenant no lock do `--isolated`** | **sem** prefixo | a spec de `console` | `kv/lock.go:129-133` e `docs/15:100-105`: um lock por tenant deixaria N réplicas rodarem a tarefa cada uma para um tenant diferente. A spec foi induzida ao erro pelo texto da REGRA 14 |
| **`hesape/container`** | não existe, em forma nenhuma | — | o resíduo já está em `foundation`; o *composition root* está em `arandu/bootstrap/app.go:101`, no **esqueleto**. `hesape` fica um componente a menos que o Illuminate de propósito, e o guia tem de dizer isso na primeira página |
| **`hesape/scheduler` no topo** | topo | `hesape/console/scheduling` | tecnicamente as duas funcionam (subpacote em Go não cria dependência do pai). `scheduler` já está em cinco ADRs e é o que se procura. Julgamento, não regra |

### 7.2 Decisões que só o proprietário pode tomar

1. **A tabela `users` do framework** (ADR 0049 acima). `arandu/app/Models/User.go:29` é um alias. `auth.Provider` torna a substituição possível sem tornar a tabela padrão extensível. Quem precisar de uma coluna a mais escreve tabela, repositório e policy próprios — e aí reset, verificação e throttle, que todos tomam `users.User`, não valem para ele. **Este documento não fecha isso.**
2. **Driver de tabela e registro subindo para o core** (`queue`, `database`). Contradiz `framework/jobs/jobs.go:6-9` e `docs/11:41`, e muda o que dois repositórios publicados contêm. É REGRA 3 na margem.
3. **`aru` importar `framework`** (ADR 0046). Sem isso, a consolidação de `jsonschema` não executa e a duplicação `aru`/`mcp` fica. Custo lateral: `aru` ganha `x/crypto` transitivo, e `aru schema` passa a depender de qual versão do framework a CLI foi compilada — `mcp` fixa `v0.25.3`; se `aru` fixar outra, os dois consumidores emitem schemas diferentes do mesmo pacote.
4. **O compilador kyse sair da `aru`** (ADR 0047).
5. **A exceção da URL pré-assinada** (ADR 0048), ou a retirada do argumento de egresso do R2 de `docs/11:90-93`.
6. **`jsonschema.Validate` vs `framework/validation`.** Não são a mesma coisa — uma checa struct/request contra strings de regra e produz frase para humano; a outra checa JSON decodificado de um modelo de linguagem contra um documento de schema, antes de qualquer handler. O Laravel também as separa (`Illuminate\JsonSchema` não tem validador; os testes dele usam `opis/json-schema`, que a CLAUDE.md:56 proíbe). Dois pacotes, nenhum implementado em termos do outro. Os resultados piores — `validation` ganhar `FromSchema`, ou `jsonschema` ganhar strings de regra — são a segunda forma que a REGRA 9 existe para impedir.
7. **A contagem de repositórios.** ~~CLAUDE.md:66 diz doze. `arandu-io/mcp` existe, tem `go.mod` próprio e foi autorizado pela ADR 0030 (09/08/2026), o que faria treze. A REGRA 3 conta repositórios; o número na regra e o número na árvore não batem.~~ **Fechado em 14/08/2026:** são catorze, com o `hesape` (ADR 0049). A REGRA 3 do `CLAUDE.md` e o `05-repositorios.md` passaram a listar os catorze nominalmente, e as demais páginas que ainda diziam dez foram corrigidas junto.
8. **A duplicação `examples/config/` vs `arandu/config/`.** Byte a byte. Toda edição no esqueleto é duas edições enquanto durar.
9. **`session.block` / a janela de lost-update.** Um `Put` de registro inteiro tipado tem essa janela. Não está decidido nem para um lado nem para o outro.
10. **HTTP/1.1 vs HTTP/2 para o stream de broadcasting.** Ou o `Hub` recusa abrir stream sob HTTP/1.1 (e a coisa não funciona em desenvolvimento), ou abre e reintroduz o defeito que `reload.go:19-31` documenta. Antes da primeira linha.

### 7.3 O que a reorganização não entrega, e nenhum nome esconde

- **Autenticação de máquina.** Não há bearer token, API key nem `TokenGuard` em lugar nenhum, e **não há decisão registrada em `docs/` nem em `docs/adr/`**. `mcp/transport.go:50` já cai em `security.Guest(tenant)` para um chamador não autenticado. Um framework de SaaS sem resposta para "como uma máquina se autentica" tem um buraco que esta reorganização não preenche.
- **`JsonResource`.** Nenhuma regra o exclui e nenhum código o fornece. Provavelmente deve ser OUT (REGRA 10: o modelo é fragmento HTMX), mas um OUT sem citação é um MISSING.
- **Stack trace para erro devolvido.** `frames.go:32 Capture` roda no ponto do recover. Um erro devolvido normalmente não carrega stack nenhum. No caminho *novo e melhor* — devolve o erro, ganha 403 — a página de debug mostra onde o handler devolveu, não onde o erro nasceu. Corrigir custa capturar stack em todo `errors.New` de `security`, `data`, `storage`, `jobs`, `events` e todo módulo, a uma alocação cada. É o maior gap honesto contra o Ignition e é propriedade da linguagem, não do layout.
- **Plurais corretos** para ru, ar, pl, cs, lt, sl, mt, lv, cy, ro, ga. ADR 0004 barra `x/text` e a REGRA 7 barra transcrever o switch do Zend. `WithPlural` é a costura, e o doc comment tem de dizer isso no primeiro parágrafo.
- **Formatação de número, moeda e data com locale.** Um dev Laravel puxa `Number::currency()` e `diffForHumans()` no mesmo reflexo de `__()`. Os dois precisam de `x/text`. Primeira reclamação depois do lançamento.
- **Componentes com filhos.** `{{ $slot }}` não tem equivalente e não tem desenho; todo `Props` em `kyse/components` é escalar plano. O primeiro componente de um dev Laravel é `<x-card>…markup arbitrário…</x-card>` e não é exprimível. **Maior lacuna de capacidade da camada de view**, maior que qualquer diretiva.
- **`route()` numa view.** Não existe e não vai existir. Toda URL é um campo que o controller preenche. Invisível até alguém escrever `{{ route('login') }}` e receber erro de compilação nomeando um campo que não existe — pelo menos é alto.
- **`$loop`.** Separador, zebra, first/last e count são escritos à mão em `@go`.
- **`@method` e `@push`/`@stack`.** Duas das vinte diretivas mais usadas do Blade, ausentes **sem decisão em lugar nenhum**. A recomendação deste documento é recusar `_method` (o HTMX carrega o verbo real) e registrar a consequência: um formulário que precise funcionar sem JavaScript não tem caminho para PUT nem DELETE. Mas a recusa não tem citação hoje e por isso continua MISSING.
- **Escape só de corpo HTML.** `template.HTMLEscapeString` cobre `<>&'"` e está correto dentro de atributo com aspas, mas não há escape de contexto URL nem de contexto JS, porque kyse gera `io.WriteString` puro. `href="{{ .UserSuppliedURL }}"` emite `javascript:` alegremente. O Laravel tem o mesmo buraco; igualar o Laravel aqui é herdar um foot-gun conhecido, e essa é uma decisão a tomar com conhecimento de causa, não por omissão.
- **Substituição de colaborador em teste.** `$this->mock()`, `Mail::fake()`, `Storage::fake()`, `Http::fake()` repousam em `instance()`. Qualquer colaborador que ainda não seja parâmetro de construtor é intestável isoladamente até alguém torná-lo um.
- **Pacotes de terceiro não se auto-registram.** `composer require` + package discovery dá bindings funcionando sem edição; aqui todo módulo é uma linha que alguém digita em `bootstrap/app.go`. Custo de ecossistema real: só o autor da aplicação sai ganhando; o autor do pacote sai perdendo.
- **Nada disto faz o Eloquent aparecer.** `Post::find(1)->update()` não tem contrapartida e nunca terá: persistência é método de repositório que exige prova de uma decisão de autorização, inclusive na leitura. Essa é a fração da migração que não é vocabulário, e nenhuma renomeação a reduz.