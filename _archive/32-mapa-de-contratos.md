# Mapa de contratos — Illuminate\Contracts contra o hesape

**Documento de medição.** Escrito contra o código em 10/08/2026, arquivo por
arquivo, sem estimativa.

Fonte da verdade: `laravel_illuminate/contracts/` — **150 arquivos**, **559
assinaturas** `public function` / `public static function`, **431 nomes de
método distintos** (o que o `plans/cobertura.sh contracts` conta), em 32
diretórios. Cruzado contra `hesape/` na mesma data.

`cobertura.sh contracts` marca `0 de 431 (0%)`, e vai continuar marcando: ele
conta símbolo exportado dentro de `hesape/contracts/`, e esse pacote não tem
nenhum de propósito. Os 431 nomes estão espalhados pelos outros 41 pacotes, e é
disso que este documento é o mapa.

O `hesape/contracts` **não vira pacote**. Em Go a interface mora no pacote que a
consome, e uma árvore de interfaces separada é a segunda maneira de declarar a
mesma coisa (REGRA 9). Mas o `Illuminate\Contracts` não é só uma árvore de
interfaces: é o **inventário da superfície do framework**. É esse inventário que
este documento reconstrói.

## Veredito, em três linhas

| Veredito | Contratos | Fração |
|---|---:|---:|
| **MISSING** — não existe equivalente nenhum | **81** | 54% |
| **PARCIAL** — existe, faltam métodos | **39** | 26% |
| **EXISTS** — o conjunto inteiro respondido | **30** | 20% |

Os 81 MISSING são o produto deste documento. Eles vêm primeiro.

Convenção das tabelas: a coluna **Go** nomeia o tipo e o arquivo; `—` quer dizer
que não há nada. A coluna **Faltando** lista os métodos do contrato PHP que
ninguém responde, com o nome do PHP.

A medição é volátil: o `cookie` era stub inteiro na primeira passagem deste
documento e hoje responde os dois contratos dele; `support.Htmlable`,
`support.Jsonable`, `pipeline.Hub` e `process.ProcessResult` mudaram de veredito
no mesmo dia. Quem reusar estes números confira a data antes.

---

# Parte 1 — Os diretórios com buraco

## Database/ — 13 MISSING, 2 EXISTS, 0 PARCIAL

O diretório com mais MISSING absoluto, e o de leitura mais fácil: **todo o
Eloquent está de fora**, e é decisão registrada (`docs/01-arquitetura.md` rejeita
Active Record). O que existe é o query builder.

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Query/Builder` | EXISTS | `query.Builder` — `database/query/builder.go` | marcador sem método; o builder responde `toSql`→`ToSQL`, `where`, `join`, `orderBy`, `limit`, `union`, `lock` e o resto |
| `Query/Expression` | EXISTS | `query.Expression.GetValue` — `database/query/expression.go` | — |
| `Query/ConditionExpression` | MISSING | — | marcador; nada distingue expressão-condição de expressão |
| `ConcurrencyErrorDetector` | MISSING | — | `causedByConcurrencyError` |
| `LostConnectionDetector` | MISSING | — | `causedByLostConnection` |
| `ModelIdentifier` | MISSING | — | `useCollectionClass`, `getClass`, `useMorphMap` |
| `Events/MigrationEvent` | MISSING | — (`database/events` é stub) | marcador |
| `Eloquent/Builder` | MISSING | — (`database/eloquent` é stub) | marcador |
| `Eloquent/Castable` | MISSING | — | `castUsing` |
| `Eloquent/CastsAttributes` | MISSING | — | `get`, `set` |
| `Eloquent/CastsInboundAttributes` | MISSING | — | `set` |
| `Eloquent/ComparesCastableAttributes` | MISSING | — | `compare` |
| `Eloquent/DeviatesCastableAttributes` | MISSING | — | `increment`, `decrement` |
| `Eloquent/SerializesCastableAttributes` | MISSING | — | `serialize` |
| `Eloquent/SupportsPartialRelations` | MISSING | — | `ofMany`, `isOneOfMany`, `getOneOfManySubQuery` |

O `database.Repository[T, ID]` (`database/repository.go`) é o contrato de dados do
hesape — `Find`, `List`, `Create`, `Update`, `Delete`, todos com `auth.Grant`. Ele
não corresponde a contrato nenhum do Illuminate: é a peça que o Laravel não tem.

## Auth/ — 11 MISSING, 2 PARCIAL, 0 EXISTS

O `auth` do hesape é **autorização**. Autenticação — quem é a pessoa, como ela
prova, como se recupera a senha — não existe.

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Access/Gate` | PARCIAL | `auth.Policy[T].Can`, `auth.Authorize`, `auth.Allows`, `auth.Grant` — `auth/policy.go` | `has`, `define`, `resource`, `policy`, `before`, `after`, `denies`, `any`, `inspect`, `raw`, `getPolicyFor`, `forUser`, `abilities` |
| `Guard` | PARCIAL | `auth.Check`, `auth.SubjectFrom`, `auth.Guest`, `auth.WithSubject`, `Subject.IsGuest` — `auth/context.go` | `id`, `validate`, `hasUser`, `setUser` |
| `Access/Authorizable` | MISSING | — | `can` (existe como `auth.Allows` e como `Policy.Can`, nunca como método do sujeito) |
| `Authenticatable` | MISSING | — (`auth/users` é stub) | `getAuthIdentifierName`, `getAuthIdentifier`, `getAuthPasswordName`, `getAuthPassword`, `getRememberToken`, `setRememberToken`, `getRememberTokenName` |
| `UserProvider` | MISSING | — | `retrieveById`, `retrieveByToken`, `updateRememberToken`, `retrieveByCredentials`, `validateCredentials`, `rehashPasswordIfRequired` |
| `StatefulGuard` | MISSING | — | `attempt`, `once`, `login`, `loginUsingId`, `onceUsingId`, `viaRemember`, `logout` |
| `Factory` | MISSING | — | `guard`, `shouldUse` |
| `CanResetPassword` | MISSING | — (`auth/passwords` é stub) | `getEmailForPasswordReset`, `sendPasswordResetNotification` |
| `PasswordBroker` | MISSING | — | `sendResetLink`, `reset` |
| `PasswordBrokerFactory` | MISSING | — | `broker` |
| `MustVerifyEmail` | MISSING | — | `hasVerifiedEmail`, `markEmailAsVerified`, `markEmailAsUnverified`, `sendEmailVerificationNotification`, `getEmailForVerification` |
| `SupportsBasicAuth` | MISSING | — | `basic`, `onceBasic` |
| `Middleware/AuthenticatesRequests` | MISSING | — (`auth/middleware` é stub) | marcador |

O `auth.SignInThrottle` (`auth/throttle.go`) — `Attempt`, `Refund`, `Clear` — não
corresponde a contrato nenhum: o Laravel resolve isso com o trait
`ThrottlesLogins`, que não é contrato. É a única peça de autenticação que existe,
e é a que protege o formulário que ninguém escreveu ainda.

## Queue/ — 10 MISSING, 3 PARCIAL, 1 EXISTS

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `ShouldBeUnique` | EXISTS | `bus.UniqueJob`, `bus.PushUnique` — `bus/unique.go` | — |
| `Queue` | PARCIAL | `queue.Queue` — `queue/queue.go` | `size`, `delayedSize`, `reservedSize`, `pushOn`, `pushRaw`, `later`, `laterOn`, `bulk`, `getConnectionName`, `setConnectionName` (`pendingSize`→`Pending`, `creationTimeOfOldestPendingJob`→`Oldest`, `pop`→`Reserve`) |
| `Job` | PARCIAL | `queue.Job` (struct) — `queue/queue.go` | `fire`, `release`, `isReleased`, `delete`, `isDeleted`, `isDeletedOrReleased`, `hasFailed`, `markAsFailed`, `fail`, `maxTries`, `maxExceptions`, `timeout`, `retryUntil`, `resolveName`, `resolveQueuedJobClass`, `getConnectionName`, `getRawBody` — no hesape essas operações são do `Queue`, não do `Job` |
| `ShouldBeUniqueUntilProcessing` | PARCIAL | `bus.ReleaseUnique` — `bus/unique.go` | não há marcador que escolha entre soltar na entrega e soltar no fim; quem chama decide |
| `Factory` | MISSING | — | `connection` |
| `Monitor` | MISSING | — | `starting`, `looping`, `failing`, `stopping` (o `console.Observer` observa a CLI, não o worker) |
| `ClearableQueue` | MISSING | — | `clear` |
| `EntityResolver` | MISSING | — | `resolve` |
| `EntityNotFoundException` | MISSING | — | — |
| `QueueableCollection` | MISSING | — | `getQueueableClass`, `getQueueableIds`, `getQueueableRelations`, `getQueueableConnection` |
| `QueueableEntity` | MISSING | — | `getQueueableId`, `getQueueableRelations`, `getQueueableConnection` |
| `ShouldBeEncrypted` | MISSING | — | marcador |
| `ShouldQueue` | MISSING | — | marcador |
| `ShouldQueueAfterCommit` | MISSING | — | marcador |

## Validation/ — 7 MISSING, 4 PARCIAL, 0 EXISTS

O hesape valida com regra compilada em `Set` (`validation/compile.go`). O Laravel
valida com objeto de regra. As duas coisas fazem o mesmo trabalho e nenhum
contrato de regra do Illuminate tem equivalente.

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Validator` | PARCIAL | `validation.Set.Validate` — `validation/set.go` | `failed`, `sometimes`, `after` (`validate`/`validated`→ retorno `(Input, Errors)`, `fails`→`Errors.Any`, `errors`→`Errors`) |
| `Factory` | PARCIAL | `validation.Compile`, `validation.MustCompile`, `validation.WithMessages` | `extend`, `extendImplicit`, `replacer` — não há registro de regra definida pelo usuário |
| `Rule` | PARCIAL | funções `validation.Required`, `Email`, `MinLen`, `MaxLen`, `Confirmed`, `NotZero`, `Unique` | `passes`/`message` não são métodos: a função escreve direto no `Errors` |
| `ValidationRule` | PARCIAL | `validation.Validatable.Validate() Errors` — `validation/validator.go` | forma diferente: valida o objeto inteiro, não um atributo com callback `$fail` |
| `CompilableRules` | MISSING | — | `compile` |
| `DataAwareRule` | MISSING | — | `setData` |
| `ValidatorAwareRule` | MISSING | — | `setValidator` |
| `ImplicitRule` | MISSING | — | marcador |
| `InvokableRule` | MISSING | — | `__invoke` — **razão 1**: método mágico do PHP |
| `UncompromisedVerifier` | MISSING | — | `verify` (regra `uncompromised`, contra Have I Been Pwned) |
| `ValidatesWhenResolved` | MISSING | — | `validateResolved` — **razão 2**: depende do container resolver o form request |

## Container/ — 6 MISSING, decisão registrada

Não é buraco: é ADR 0001. A fiação é explícita e escrita à mão em
`bootstrap/app.go`. Todos os seis entram pela **razão 2**.

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Container` | MISSING | — | os 24 métodos: `bound`, `alias`, `tag`, `tagged`, `bind`, `bindMethod`, `bindIf`, `singleton`, `singletonIf`, `scoped`, `scopedIf`, `extend`, `instance`, `addContextualBinding`, `when`, `factory`, `flush`, `make`, `call`, `resolved`, `beforeResolving`, `resolving`, `afterResolving`, `get` |
| `ContextualBindingBuilder` | MISSING | — | `needs`, `give`, `giveTagged`, `giveConfig` |
| `ContextualAttribute` | MISSING | — | marcador |
| `SelfBuilding` | MISSING | — | marcador |
| `BindingResolutionException` | MISSING | — | — |
| `CircularDependencyException` | MISSING | — | — |

## Broadcasting/ — 5 MISSING, 2 PARCIAL

`hesape/broadcasting` é stub. O único pedaço vivo está no canal de notificação.

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Broadcaster` | PARCIAL | `channels.Broadcaster.Push` — `notifications/channels/broadcast.go` | `auth`, `validAuthenticationResponse` — não há autorização de canal privado |
| `ShouldBroadcast` | PARCIAL | `channels.BroadcastNotification.ToBroadcast` — `notifications/channels/broadcast.go` | `broadcastOn` — o canal vem do `Notifiable`, não do evento |
| `Factory` | MISSING | — | `connection` |
| `HasBroadcastChannel` | MISSING | — | `broadcastChannelRoute`, `broadcastChannel` |
| `ShouldBroadcastNow` | MISSING | — | marcador |
| `ShouldBeUnique` | MISSING | — | marcador |
| `ShouldRescue` | MISSING | — | marcador |

## Support/ — 6 MISSING, 0 PARCIAL, 6 EXISTS

O diretório que mais mudou: `Htmlable` e `Jsonable` passaram a existir como
interface, não como método solto num tipo ou outro.

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Arrayable` | EXISTS | `arr.Arrayer` — `support/arr/arr.go` | — |
| `Htmlable` | EXISTS | `support.Htmlable`, `support.HtmlString` — `support/htmlstring.go` | — |
| `Jsonable` | EXISTS | `support.Jsonable` — `support/htmlstring.go` (`toJson`→`ToJson() (string, error)`, porque `encoding/json` falha com erro e não tem flags) | — |
| `MessageBag` | EXISTS | `support.MessageBag` — `support/messagebag.go` | — |
| `MessageProvider` | EXISTS | `support.MessageProvider` — `support/messagebag.go` | — |
| `ValidatedData` | EXISTS | `validation.Input` — `validation/set.go` | `ArrayAccess`/`IteratorAggregate` pulados por serem interface de linguagem do PHP (**razão 1**) |
| `Renderable` | MISSING | — | `render` (existe `view.Renderer.Render`, que é o motor, não o valor) |
| `Responsable` | MISSING | — | `toResponse` |
| `CanBeEscapedWhenCastToString` | MISSING | — | `escapeWhenCastingToString` |
| `DeferringDisplayableValue` | MISSING | — | `resolveDisplayableValue` |
| `HasOnceHash` | MISSING | — | `onceHash` |
| `DeferrableProvider` | MISSING | — | `provides` — **razão 2**: só serve service provider (ADR 0002) |

## Foundation/ — 4 MISSING, 1 PARCIAL

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `ExceptionRenderer` | PARCIAL | `exception.Handler.Render`, `exception.Views` — `exception/handler.go` | assinatura diferente: escreve na resposta em vez de devolver string |
| `Application` | MISSING | — (`foundation` tem só as interfaces de módulo) | os 33 métodos: `version`, `basePath`, `bootstrapPath`, `configPath`, `databasePath`, `langPath`, `publicPath`, `resourcePath`, `storagePath`, `environment`, `runningInConsole`, `runningUnitTests`, `hasDebugModeEnabled`, `maintenanceMode`, `isDownForMaintenance`, `registerConfiguredProviders`, `register`, `registerDeferredProvider`, `resolveProvider`, `boot`, `booting`, `booted`, `bootstrapWith`, `getLocale`, `getNamespace`, `getProviders`, `hasBeenBootstrapped`, `loadDeferredProviders`, `setLocale`, `shouldSkipMiddleware`, `terminating`, `terminate` — os de provider entram pela **razão 2**, os de caminho não |
| `MaintenanceMode` | MISSING | — | `activate`, `deactivate`, `active`, `data` |
| `CachesConfiguration` | MISSING | — | `configurationIsCached`, `getCachedConfigPath`, `getCachedServicesPath` |
| `CachesRoutes` | MISSING | — | `routesAreCached`, `getCachedRoutesPath` |

`foundation.Bootable`, `Background`, `Closable`, `Diagnostic`, `Schedulable`,
`Migratable` e `Health` (`foundation/module.go`) são o contrato de módulo do
hesape. Não correspondem a contrato do Illuminate: substituem o service provider.

## Redis/ — 4 MISSING

`hesape/redis` é stub inteiro. O `kv` é repositório à parte (ADR 0016) e não
entra nesta contagem.

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Connection` | MISSING | — | `subscribe`, `psubscribe`, `command` |
| `Connector` | MISSING | — | `connect`, `connectToCluster` |
| `Factory` | MISSING | — | `connection` |
| `LimiterTimeoutException` | MISSING | — | — (`cache.ErrLockTimeout` é o do cache, não o do limiter Redis) |

## Mail/ — 3 MISSING, 2 PARCIAL

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Mailer` | PARCIAL | `mail.Mailer`, `mail.Pending` — `mail/mail.go` | `raw`, `sendNow` (`to`→`To`/`ToAddress`, `cc`→`Pending.CC`, `bcc`→`Pending.BCC`, `send`→`Pending.Send`) |
| `Mailable` | PARCIAL | `mail.Mailable` (`Envelope`+`Content`) — `mail/mail.go` | `send`, `queue`, `later`, `cc`, `bcc`, `to`, `locale`, `mailer` — no hesape o destinatário é do `Mailer`, não do `Mailable` |
| `Factory` | MISSING | — | `mailer` |
| `MailQueue` | MISSING | — | `queue`, `later` |
| `Attachable` | MISSING | — | `toMailAttachment` |

Curiosidade que vale registrar: `fakes.MailFake`
(`support/testing/fakes/mailfake.go`) **tem** `Raw`, `SendNow`, `Queue` e `Later`. O
dublê de teste responde quatro métodos que o `mail.Mailer` de verdade não tem, o
que significa que um teste passa exercitando uma API que produção não oferece.

## View/ — 3 MISSING, 1 PARCIAL

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Factory` | PARCIAL | `view.Register`, `view.Registered`, `view.Renderer` — `view/render.go` | `file`, `share`, `composer`, `creator`, `addNamespace`, `replaceNamespace` (`exists`→`Registered`, `make`→`Render`) |
| `Engine` | MISSING | — (`view/engines` é stub) | `get` — o kyse compila para função Go; não há motor plugável, e é REGRA 9 |
| `View` | MISSING | — | `name`, `with`, `getData`, `render` — não existe objeto de view; `view.Page` é a struct de dados |
| `ViewCompilationException` | MISSING | — | — (`view.WrongData` é outro erro: dado do tipo errado, não template quebrado) |

## Routing/ — 2 MISSING, 3 PARCIAL

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Registrar` | PARCIAL | `routing.Router`, `routing.Resource[C]` — `routing/router.go` | `options`, `substituteBindings`, `substituteImplicitBindings` |
| `ResponseFactory` | PARCIAL | `http.Context` — `http/context.go`; `filesystem.Send`; `http.Intended` | `make`, `noContent`, `jsonp`, `eventStream`, `stream`, `streamJson`, `streamDownload`, `redirectToAction`, `redirectGuest` |
| `UrlGenerator` | PARCIAL | `support.UrlGenerator` (`To`, `Route`, `SignedRoute`, `Action`), `support.TemporarySignedRoute`, `routing.Routes.Route`/`Must`, `http.Context.FullURL`, `http.Back`, `view.Assets` | `secure`, `query`, `getRootControllerNamespace`, `setRootControllerNamespace` |
| `UrlRoutable` | MISSING | — | `getRouteKey`, `getRouteKeyName`, `resolveRouteBinding`, `resolveChildRouteBinding` |
| `BindingRegistrar` | MISSING | — | `bind`, `getBindingCallback` |

## Cache/ — 1 MISSING, 2 PARCIAL, 4 EXISTS

O diretório mais bem coberto que ainda tem furo.

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Repository` | EXISTS | `cache.Repository` + `Remember[T]`, `Get[T]`, `Pull[T]`, `Sear[T]`, `RememberForever[T]` — `cache/repository.go` | — |
| `Lock` | EXISTS | `cache.Lock` — `cache/lock.go` (`get`→`Get`, `block`→`Block`, `release`, `owner`, `forceRelease`) | — |
| `LockProvider` | EXISTS | `cache.Locks.Lock`, `cache.Locks.RestoreLock` — `cache/lock.go` | — |
| `LockTimeoutException` | EXISTS | `cache.ErrLockTimeout` — `cache/lock.go` | — |
| `Store` | PARCIAL | `cache.Store` — `cache/store.go` (declara `Get`, `Put`, `Add`, `Forget`, `Increment`, `Flush`) | `many`, `putMany`, `decrement`, `forever`, `touch`, `getPrefix` **não estão na interface** — o `ArrayStore` implementa todos, então um Store novo compila sem eles |
| `CanFlushLocks` | PARCIAL | `cache.CanFlushLocks` — `cache/store.go` | `hasSeparateLockStore` está no `ArrayStore`, não na interface |
| `Factory` | MISSING | — | `store` |

## Console/ — 1 MISSING, 2 PARCIAL, 1 EXISTS

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Isolatable` | EXISTS | `console.Command.Isolated`, `console.Registry.WithLocks` — `console/command.go` | — |
| `Kernel` | PARCIAL | `console.Registry` — `console/registry.go` | `bootstrap`, `queue`, `output`, `terminate` (`handle`→`Handle`, `call`→`Call`, `all`→`Names`) |
| `Application` | PARCIAL | `console.Registry.Call`, `console.Registry.CallSilent` | `output` |
| `PromptsForMissingInput` | MISSING | — | marcador — `console.IO.Ask` existe, mas nada liga argumento ausente a pergunta |

## Events/ — 1 MISSING, 1 PARCIAL, 1 EXISTS

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `ShouldDispatchAfterCommit` | EXISTS | `events.Outbox.Store` + `events.ErrNoTransaction` — `events/events.go` | não é marcador: **todo** evento é after-commit por construção, e um fora de transação é recusado |
| `Dispatcher` | PARCIAL | `events.Publisher` (`On`, `Publish`, `Listeners`), `events.Outbox` — `events/listeners.go` | `hasListeners`, `subscribe`, `until`, `push`, `flush`, `forget`, `forgetPushed` |
| `ShouldHandleEventsAfterCommit` | MISSING | — | marcador (é sobre o *listener*, não sobre o evento) |

## Filesystem/ — 1 MISSING, 2 PARCIAL, 2 EXISTS

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Factory` | EXISTS | `filesystem.Disks.Disk` — `filesystem/disk.go` | — |
| `FileNotFoundException` | EXISTS | `filesystem.ErrNotFound` — `filesystem/disk.go` | — |
| `Filesystem` | PARCIAL | `filesystem.Disk`, `filesystem.Adapter` — `filesystem/disk.go` | `path`, `readStream`, `putFile`, `putFileAs`, `writeStream`, `getVisibility`, `setVisibility`, `prepend`, `append`, `files`, `allFiles`, `directories`, `allDirectories`, `makeDirectory`, `deleteDirectory` (`size`/`lastModified` estão em `Info`) |
| `Cloud` | PARCIAL | `filesystem.URLSigner.TemporaryURL`, `filesystem.Presigner` — `filesystem/url.go` | `url` — só existe URL assinada e temporária; não há URL pública permanente |
| `LockTimeoutException` | MISSING | — | — não há lock de arquivo |

## Http/ — 1 MISSING

`hesape/http` é stub (só `data.go`, com helpers não exportados). O tráfego HTTP
vive em `http` e `routing`.

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Kernel` | MISSING | — | `bootstrap`, `handle`, `terminate`, `getApplication` — `routing.Router.ServeHTTP` atende requisição, mas não há kernel com ciclo de vida |

## Log/ — 1 MISSING

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `ContextLogProcessor` | MISSING | — | marcador sobre a `ProcessorInterface` do Monolog; `log.With`, `log.Fields`, `log.Into` e `log.ShareContext` (`log/logger.go`) fazem o trabalho, sem interface |

## Session/ — 1 MISSING, 1 PARCIAL

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Session` | PARCIAL | `session.Store[T]`, `session.CSRF`, `session.Flash`, `http.Intended` — `session/session.go` | `getName`, `setName`, `setId`, `save`, `exists`, `has`, `get`, `pull`, `put`, `remove`, `forget`, `flush`, `migrate`, `isStarted`, `getHandler`, `handlerNeedsRequest`, `setRequestOnHandler` — o `Record[T]` é tipado, então não há bag chave-valor |
| `Middleware/AuthenticatesSessions` | MISSING | — (`session/middleware` é stub) | marcador |

---

# Parte 2 — Os diretórios sem MISSING

## Pagination/ — 3 PARCIAL

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `LengthAwarePaginator` | PARCIAL | `pagination.LengthAwarePaginator[T]` — `pagination/lengthaware.go` (`getUrlRange`→`GetURLRange`, `total`, `lastPage`) | `render` |
| `Paginator` | PARCIAL | `pagination.Paginator[T]` — `pagination/simple.go` | `render` |
| `CursorPaginator` | PARCIAL | `pagination.CursorPaginator[T]` — `pagination/cursorpaginator.go` | `appends`, `fragment`, `withQueryString`, `path`, `isEmpty`, `isNotEmpty`, `render` |

`render` falta nos três pelo mesmo motivo: os links viram `Links()` /
`LinkCollection()`, e quem desenha é a view. É decisão, mas não está escrita em
ADR nenhuma.

## Bus/ — 2 PARCIAL

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Dispatcher` | PARCIAL | `bus.Chain.Dispatch`, `bus.Queue.Push`, `queue.Worker.Handle` — `bus/chain.go` | `dispatchSync`, `dispatchNow`, `dispatchAfterResponse`, `hasCommandHandler`, `getCommandHandler`, `pipeThrough` (`chain`→`NewChain`, `map`→`Worker.Handle`) |
| `QueueingDispatcher` | PARCIAL | `bus.Store.Find`, `bus.NewBatch`, `bus.Queue.Push` — `bus/batch.go` | os três métodos são respondidos, mas por três tipos: não há um despachante |

## Debug/ — 2 PARCIAL

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `ExceptionHandler` | PARCIAL | `exception.Handler` — `exception/handler.go` (`report`, `render`, `renderForConsole`) | `shouldReport` não é exportado (é `dontReport`, não exportado) |
| `ShouldntReport` | PARCIAL | `exception.Config.DontReport []error` — `exception/handler.go` | é lista de sentinelas comparada com `errors.Is`, não marcador no tipo do erro |

## Notifications/ — 2 PARCIAL

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Dispatcher` | PARCIAL | `notifications.Notifier.Send`, `SendMany` — `notifications/notifier.go` | `sendNow` — não há distinção entre enfileirar e entregar agora |
| `Factory` | PARCIAL | `notifications.Notifier.Channels` — `notifications/notifier.go` | `channel($name)` — lista nomes, mas não devolve o canal |

## Translation/ — 2 PARCIAL, 1 EXISTS

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `HasLocalePreference` | EXISTS | `notifications.Localized.PreferredLocale` — `notifications/notification.go` | — |
| `Translator` | PARCIAL | `translation.Translator` — `translation/translation.go` (`get`→`Get`, `choice`→`Choice`, `getLocale`→`Locale`) | `setLocale` — o idioma vem do contexto (`translation.WithLocale`), não de estado no tradutor |
| `Loader` | PARCIAL | `translation.Loader.Load` — `translation/loader.go` | `addNamespace`, `addJsonPath`, `namespaces` — `WithNamespace` é opção do `Translator`, não método do loader |

## Process/ — 1 PARCIAL, 1 EXISTS

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `ProcessResult` | EXISTS | `process.ProcessResult` — `process/result.go` (`command`, `successful`, `failed`, `exitCode`, `output`, `seeInOutput`, `errorOutput`, `seeInErrorOutput`, `throw`, `throwIf`) | — |
| `InvokedProcess` | PARCIAL | `process.InvokedProcess` — `process/invoked.go` (`id`→`ID`, `signal`, `running`, `output`, `errorOutput`, `latestOutput`, `latestErrorOutput`, `wait`) | `command`, `waitUntil` |

## Pipeline/ — 2 EXISTS

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Pipeline` | EXISTS | `pipeline.Pipeline[T]` — `pipeline/pipeline.go` (`send`→`Send`, `through`→`Through`, `then`→`Then`) | `via` pulado pela **razão 1**: escolhe por nome qual método chamar em cada pipe, que é despacho dinâmico do PHP; um `Pipe[T]` tipado tem um método só |
| `Hub` | EXISTS | `pipeline.Hub[T].Pipe` — `pipeline/hub.go` | — |

## Cookie/ — 2 EXISTS

Era stub inteiro na primeira medição deste documento, no mesmo dia. Hoje o
`cookie.CookieJar` responde os dois contratos, e `cookie/middleware` traz
`EncryptCookies` e `AddQueuedCookiesToResponse`.

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Factory` | EXISTS | `cookie.CookieJar.Make`, `.Forever`, `.Forget` — `cookie/jar.go` | — |
| `QueueingFactory` | EXISTS | `cookie.CookieJar.Queue`, `.Unqueue`, `.GetQueuedCookies` — `cookie/jar.go` | — |

## Concurrency/ — 1 PARCIAL

| Contrato | Veredito | Go | Faltando |
|---|---|---|---|
| `Driver` | PARCIAL | `concurrency.Run[T]`, `concurrency.RunWith[T]` — `concurrency/run.go` | `defer` — não há `DeferredCallback` |

## Config/, Encryption/, Hashing/, JsonSchema/ — completos

Quatro diretórios, sete contratos, nenhum furo.

| Contrato | Veredito | Go |
|---|---|---|
| `Config/Repository` | EXISTS | `config.Repository` — `config/repository.go` (`has`, `get`, `all`, `set`, `prepend`, `push`) |
| `Encryption/Encrypter` | EXISTS | `encryption.Encrypter` + `Encrypt[T]`/`Decrypt[T]` — `encryption/encrypter.go` (`getKey`, `getAllKeys`, `getPreviousKeys`) |
| `Encryption/StringEncrypter` | EXISTS | `encryption.Encrypter.EncryptString`/`DecryptString` — `encryption/encrypter.go` |
| `Encryption/DecryptException` | EXISTS | `encryption.ErrDecrypt` |
| `Encryption/EncryptException` | EXISTS | `encryption.ErrEncrypt` |
| `Hashing/Hasher` | EXISTS | `hashing.Hasher` — `hashing/hashers.go` (`info`, `make`, `check`, `needsRehash`), com `BcryptHasher`, `ArgonHasher`, `Argon2IdHasher` |
| `JsonSchema/JsonSchema` | EXISTS | `jsonschema.Object`, `Array`, `String`, `Integer`, `Number`, `Boolean` — `jsonschema/type.go` |

---

# O que a lista mostra

**Autenticação não existe. São 11 contratos dos 13 de `Auth/`.**
O `auth` do hesape é autorização — `Grant`, `Policy[T]`, `Authorize` — e faz isso
melhor que o Laravel, porque o compilador cobra. Mas não há `Authenticatable`,
`UserProvider`, `StatefulGuard`, `PasswordBroker` nem `MustVerifyEmail`. Quem vem
do Laravel abre o pacote esperando `Auth::attempt($credentials)` e encontra a
metade do problema que vem *depois* de saber quem é a pessoa. Um projeto Arandu
hoje escreve o login à mão: buscar o usuário, comparar o hash, iniciar a sessão,
lembrar o dispositivo, mandar o e-mail de verificação, expirar o token de reset.
São seis peças, todas com armadilha de segurança, e nenhuma tem dono. O
`auth.SignInThrottle` já está pronto — é o cadeado do formulário que ninguém
escreveu. É o buraco que impede a primeira aplicação real de existir sem código
de framework escrito no projeto, que é exatamente o que o framework existe para
evitar.

**A camada de fila entrega, mas não se administra. 10 MISSING em 14.**
`queue.Queue` e `queue.Worker` fazem o essencial: empurrar, reservar com lease,
confirmar, falhar com backoff, estacionar, retomar. O que falta é tudo em volta.
`Monitor` (`starting`, `looping`, `failing`, `stopping`) é o que instrumenta um
worker em produção, e sem ele não há como plugar métrica nem alerta sem mexer no
loop. `Factory`/`connection` é o que permite mais de uma conexão de fila, e sua
ausência amarra a aplicação a uma. `ShouldQueue` e `ShouldBeEncrypted` não são
cerimônia: são a diferença entre um job dizer o que ele é e cada chamador
lembrar. Somado ao `Job` PARCIAL — 17 métodos ausentes, porque no hesape as
operações são do `Queue` e não do `Job` — o resultado é uma fila que roda e uma
fila que ninguém opera.

**Dois componentes ainda são stub inteiro: `broadcasting` e `redis`.**
São 9 contratos MISSING sem uma linha de Go. Broadcasting derruba o canal
`ShouldBroadcast` de notificação real: o `channels.Broadcaster` só tem `Push`,
sem `auth` de canal privado, então **não existe canal privado** — qualquer um que
abra o stream vê o que passa por ele. Redis ausente não é grave por si (o `kv` é
repositório à parte, ADR 0016), mas `Connection.subscribe` é o transporte que o
broadcasting precisaria, então os dois buracos são um só.

O terceiro stub desta lista saiu dela hoje. O `cookie` era MISSING nos dois
contratos e agora responde os dois, com `CookieJar` e o middleware que criptografa
e anexa a fila à resposta — o que tira três lugares que decidiam `Secure`,
`SameSite` e `HttpOnly` por conta própria e põe um.

**Todo o Eloquent está fora, e isso é decisão, não buraco.**
13 dos 15 contratos de `Database/` são MISSING e nove deles são casts e relações
de Active Record, rejeitados em `docs/01-arquitetura.md`. O `database.Repository[T,
ID]` é o substituto e é melhor: exige `auth.Grant` na assinatura, então não há
consulta sem policy (REGRA 17). O que vale registrar é o que **não** é Eloquent e
mesmo assim falta: `LostConnectionDetector` e `ConcurrencyErrorDetector`. Os dois
são o que faz um pool sobreviver a um failover e uma transação sobreviver a um
deadlock, e nenhum depende de ORM. Hoje um erro de conexão perdida e um erro de
SQL malformado são o mesmo `error` para quem chama.

**O que a contagem não diz.**
Os 30 EXISTS são reais e completos. Mas dois dos 39 PARCIAL são de um tipo
específico e pior: o contrato existe como interface Go e **está incompleto na
própria interface**. `cache.Store` declara seis métodos e não declara `many`,
`putMany`, `decrement`, `forever`, `touch` nem `getPrefix`, embora o `ArrayStore`
implemente todos os seis. Um Store novo — Redis, Postgres — compila sem eles e
passa nos testes. `cache.CanFlushLocks` tem o mesmo problema com
`hasSeparateLockStore`. Interface incompleta é pior que método faltando, porque o
compilador para de cobrar.

E há o inverso, no `mail`: o `fakes.MailFake` responde `Raw`, `SendNow`, `Queue`
e `Later`, que o `mail.Mailer` de produção não tem. O dublê é maior que o
original, então existe teste verde exercitando API que a aplicação não pode
chamar.
