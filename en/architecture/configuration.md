# Configuration

- [Configuration is a typed boot input](#configuration-is-a-typed-boot-input)
- [Load once and fail before serving traffic](#load-once-and-fail-before-serving-traffic)
- [Use one owner for each setting](#use-one-owner-for-each-setting)
- [Reject invalid and unwired values visibly](#reject-invalid-and-unwired-values-visibly)
- [Trace a value into the constructor that uses it](#trace-a-value-into-the-constructor-that-uses-it)
- [Loaders and typed configuration](#loaders-and-typed-configuration)
- [Configuration chooses adapters; code keeps the contract](#configuration-chooses-adapters-code-keeps-the-contract)

> [!WARNING]
> **Current limitations.** The framework reads every driver setting into `bootstrap.Settings` and builds the drivers from it, and the released skeleton adds only its own sections. `LOG_SAMPLING` is still read and validated without reaching the logger, so test every additional variable at the consumer that should read it.

<a name="configuration-is-a-typed-boot-input"></a>
## Configuration is a typed boot input

The application reads configuration once, validates it before opening the server, and passes typed values into constructors. Handlers do not query environment variables and components do not look up string keys at the point of use.

> [!WARNING]
> The framework owns every setting that builds a driver. `bootstrap.LoadSettings` reads the application, HTTP, database, cache, session, queue, mail, filesystem, and logging variables into `bootstrap.Settings`, and the skeleton's `config.Config` embeds those settings beside its own Auth, Geo, and Services sections. Do not read a framework variable a second time under `config/`: two readers of one name drift apart.

<a name="load-once-and-fail-before-serving-traffic"></a>
## Load once and fail before serving traffic

`bootstrap/app.go` loads the configuration, has `foundation.NewDrivers` build every driver the settings name, builds the application's provider over those drivers, and hands the same drivers to `foundation.Configure`. Create returns the composed application or the first contradiction it finds.

```go
cfg, err := appconfig.Load()
if err != nil {
	return nil, err
}

drivers, err := foundation.NewDrivers(cfg.Settings)
if err != nil {
	return nil, err
}

app, err := providers.NewAppServiceProvider(cfg, drivers)
if err != nil {
	_ = drivers.Close()
	return nil, err
}

return foundation.Configure(drivers).WithProviders(app).Create()
```

> [!NOTE]
> A driver changes by changing the environment, never in code. `CACHE_STORE`, `SESSION_DRIVER`, `QUEUE_CONNECTION`, `MAIL_URL`, and `FILESYSTEM_DISK` name what `NewDrivers` builds, and a value naming a driver the binary does not link stops the boot with the two lines that add it.

<a name="use-one-owner-for-each-setting"></a>
## Use one owner for each setting

| Setting family | Owner in the skeleton | Unit or shape |
| --- | --- | --- |
| `APP_*`, `HTTP_*`, `TRUSTED_PROXIES`, and observability | framework `bootstrap.Settings` | typed values |
| `DATABASE_URL` and `DB_*` | framework `bootstrap.Settings` | one URL; pool counts and seconds |
| `SESSION_DRIVER`, `SESSION_LIFETIME`, `SESSION_SECURE_COOKIE`, and `CSRF_TTL` | framework `bootstrap.Settings` | memory or redis; minutes; boolean; positive seconds |
| `CACHE_STORE`, `CACHE_PREFIX`, `CACHE_TTL`, and `REDIS_*` | framework `bootstrap.Settings` | memory or redis; prefix; seconds; one URL and TLS files |
| `QUEUE_CONNECTION`, `QUEUE_DEFAULT`, `QUEUE_WORKERS`, `QUEUE_RETRY_AFTER`, and `QUEUE_MAX_ATTEMPTS` | framework `bootstrap.Settings` | database or redis; queue name; positive counts; retry is seconds |
| `MAIL_URL` and `MAIL_FROM_*` | framework `bootstrap.Settings` | one transport URL; sender |
| `FILESYSTEM_*` | framework `bootstrap.Settings` | local, r2, or s3; root, bucket, endpoint, region, and keys |
| `LOG_LEVEL`, `LOG_FORMAT`, and `LOG_SAMPLING` | framework `bootstrap.Settings` | level; json or text; one in N |
| `ARANDU_TENANT_ID`, `AUTH_*`, `GEO_*`, and service credentials | the application's own sections under `config/` | typed fields |

The `.env` file fills values that the process environment did not already provide. Deployment environment values therefore win. Keep secrets outside source control, and use the example file only as a catalogue of candidate inputs—not as proof that every candidate reaches a consumer.

### `APP_ENV` decides what production refuses

`APP_ENV` is dev, staging, or prod. The `.env` copied from `.env.example` says dev; a process with no `APP_ENV` is prod, which turns off the debug console, the stack on the error page, and the commands that empty a database. A development machine names `APP_ENV=dev`, and a container image sets `APP_ENV=prod` rather than leaving the answer to an unset variable.

### Session names are not aliases

The session lifetime has one name. The framework reads `SESSION_LIFETIME`, in minutes, and `foundation.NewDrivers` builds the session store from it. `SESSION_TTL`, the count of seconds the skeleton used to read, stops the boot with the `SESSION_LIFETIME` value to write instead.

<a name="reject-invalid-and-unwired-values-visibly"></a>
## Reject invalid and unwired values visibly

Malformed booleans, non-positive counts, invalid URLs, missing keys, and unsupported store names should stop boot with the variable named in the error. Do not replace an invalid explicit value with a default; that makes the process report a deployment other than the one its operator requested.

> [!WARNING]
> Production refuses two values at boot. `SESSION_SECURE_COOKIE=false` would send the session, the credential itself, in the clear; a deployment served over http is a staging one and sets `APP_ENV=staging`. `MAIL_URL=log://`, which is also the value of an unset `MAIL_URL`, would write every message, password-reset codes included, into the log; production sets `smtp://`, `smtps://`, `resend://`, or `sendgrid://`.

> [!WARNING]
> In framework :frameworkVersion, the session cookie's name and scope are not settings: it is `arandu_session`, on path /, for the host that answered. `SESSION_PATH`, `SESSION_DOMAIN`, or `SESSION_SAME_SITE` asking for another scope stops the boot, `SESSION_SECURE` stops it naming `SESSION_SECURE_COOKIE`, and `SESSION_COOKIE` is read by nothing.

<a name="trace-a-value-into-the-constructor-that-uses-it"></a>
## Trace a value into the constructor that uses it

A parser test proves syntax only. For every operational setting, add or find a test that changes the environment, builds the application, and observes the downstream component. The session replica test and queue worker settings tests are the model: they fail if wiring silently restores a constructor default.

```sh
aru doctor --strict
go test ./tests/Unit ./tests/Feature
```

- Invalid explicit values make boot fail and name the setting.
- An absent value produces the documented default.
- A changed value reaches the component's observable behavior.
- Secrets and connection credentials are redacted from errors.

<a name="loaders-and-typed-configuration"></a>
## Loaders and typed configuration

`bootstrap.LoadSettings` reads the environment once and answers `bootstrap.Settings`, one typed section per driver beside the framework's own configuration. In the skeleton, `config.Load` adds the application's sections and `config.From` does the same over settings a test supplies. `foundation.NewDrivers` builds the database, cache stores, locks, session store, CSRF issuer, rate limiter, queue, mailer, and disks from the settings, and `foundation.Configure` composes the application over those drivers.

- [Framework bootstrap configuration](https://pkg.go.dev/github.com/arandu-io/framework/foundation/bootstrap)
- [Framework drivers and builder](https://pkg.go.dev/github.com/arandu-io/framework/foundation)
- [Hesape configuration utilities](https://pkg.go.dev/github.com/arandu-io/hesape/config)

<a name="configuration-chooses-adapters-code-keeps-the-contract"></a>
## Configuration chooses adapters; code keeps the contract

A store, driver, or transport choice belongs in the configuration, and the composition root builds it once through `foundation.NewDrivers`. Application services continue to depend on the stable interface they receive. This keeps configuration from becoming a service locator and makes a test able to construct the same graph with explicit substitutes.

- [Modules and composition root](/docs/{{version}}/architecture/modules-and-composition)
- [HTTP session](/docs/{{version}}/http/session)
- [Contracts and adapters](/docs/{{version}}/architecture/contracts-adapters)
