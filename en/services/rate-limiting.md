# Rate limiting

- [A limiter meters one explicit key in one store](#a-limiter-meters-one-explicit-key-in-one-store)
- [Construct the limiter and throttle a route](#construct-the-limiter-and-throttle-a-route)
- [Choose keys and refusal policy by threat](#choose-keys-and-refusal-policy-by-threat)
- [The released throttle fails open on store errors](#the-released-throttle-fails-open-on-store-errors)
- [Prove windows, concurrency, and outage behavior](#prove-windows-concurrency-and-outage-behavior)
- [Released cache limiter and middleware](#released-cache-limiter-and-middleware)
- [Place the throttle before expensive authentication work](#place-the-throttle-before-expensive-authentication-work)

<a name="a-limiter-meters-one-explicit-key-in-one-store"></a>
## A limiter meters one explicit key in one store

The current released path is `cache.NewRateLimiter` plus routing `middleware.Throttle`. The framework already builds one limiter, `drivers.Limiter`, counting in the store `CACHE_STORE` names, and the global pipeline throttles every request at three hundred a minute per session or address; `WithMiddleware` changes that budget with `Middleware.Throttle`. framework `v0.41.0` removed the older framework limiter surface; copy no examples that import it.

<a name="construct-the-limiter-and-throttle-a-route"></a>
## Construct the limiter and throttle a route

Pass `drivers.Limiter` to Throttle with a positive limit, decay, key function, and refusal callback for a route that needs a tighter budget than the global one. Middleware panics at construction for missing collaborators or non-positive policy. The published sign-in screens count their code forms in the same limiter.

```go example
import middleware "github.com/arandu-io/hesape/routing/middleware"

func LoginThrottle(drivers *foundation.Drivers, clientKey func(*stdhttp.Request) string) func(stdhttp.Handler) stdhttp.Handler {
	limiter := drivers.Limiter
	throttle := middleware.Throttle(limiter, cache.PerMinute(5),
		func(r *stdhttp.Request) string { return "login:" + clientKey(r) },
		func(w stdhttp.ResponseWriter, r *stdhttp.Request, status int, message string) {
			stdhttp.Error(w, message, status)
		},
	)
	return throttle
}
```

<a name="choose-keys-and-refusal-policy-by-threat"></a>
## Choose keys and refusal policy by threat

Choose account plus trusted network keys for credential boundaries, normalize keys, keep tenant scope where known, send Retry-After, use a shared store for fleet-wide budgets, and separate general abuse limits from strict authentication refusal policy.

- Construct once at boot.
- Normalize bounded keys.
- Use a shared production store.
- Return an explicit 429 refusal.

<a name="the-released-throttle-fails-open-on-store-errors"></a>
## The released throttle fails open on store errors

The released Throttle logs a backing-store error and calls the handler: it fails open. That preserves availability for ordinary endpoints but cannot be presented as fail-closed protection for sign-in, password reset, or another security boundary without an application policy that refuses on store outage.

> [!WARNING]
> Memory stores split counters across replicas and reset on restart. A shared store outage and key-cardinality attack need explicit health, alerting, and refusal decisions.

<a name="prove-windows-concurrency-and-outage-behavior"></a>
## Prove windows, concurrency, and outage behavior

Test first request through limit, next refusal, Retry-After, window expiry, concurrent increments, two app instances on one store, untrusted proxy headers, key normalization, store outage, and boot panic for invalid configuration.

- The shared budget spans replicas.
- Concurrency cannot exceed the intended count.
- Store outage follows the declared fail policy.
- Credential hashing is never reached after refusal.

<a name="released-cache-limiter-and-middleware"></a>
## Released cache limiter and middleware

`cache.NewRateLimiter` supplies the current store-backed counter; `routing/middleware.Throttle` adapts it to HTTP. The application chooses keys, thresholds, shared store, refusal representation, and fail-open or fail-closed security policy.

| Surface | Released contract |
| --- | --- |
| `cache.NewRateLimiter` | store-backed fixed-window decisions |
| `middleware.Throttle` | HTTP key and refusal adapter |
| Retry-After | client retry guidance |
| shared cache store | fleet-wide counter state |

<a name="place-the-throttle-before-expensive-authentication-work"></a>
## Place the throttle before expensive authentication work

Middleware controls ordering, authentication and reset define sensitive threats, cache selects shared state, and hashing explains the expensive work a credential limiter protects.

- [Middleware](/docs/{{version}}/http/middleware)
- [Authentication](/docs/{{version}}/security/authentication)
- [Cache and Redis](/docs/{{version}}/services/cache-redis)
- [Hashing](/docs/{{version}}/security/hashing)
