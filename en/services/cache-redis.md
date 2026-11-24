# Cache and Redis

- [The application uses a tenant-scoped repository](#the-application-uses-a-tenant-scoped-repository)
- [Begin with the in-process store](#begin-with-the-in-process-store)
- [Store, remember, count, and invalidate](#store-remember-count-and-invalidate)
- [A cache miss is expected; a store error is not](#a-cache-miss-is-expected-a-store-error-is-not)
- [Test isolation and expiry through the repository](#test-isolation-and-expiry-through-the-repository)
- [Repository and store boundaries](#repository-and-store-boundaries)
- [A protocol-compatible server is a deployment choice](#a-protocol-compatible-server-is-a-deployment-choice)

<a name="the-application-uses-a-tenant-scoped-repository"></a>
## The application uses a tenant-scoped repository

The cache Repository adds JSON encoding, expiry, events, and tenant-prefixed keys over a Store. Application code holds the repository, not the raw backend. Every ordinary cache operation takes a Grant, so the tenant comes from authorization rather than a request field or hand-built prefix.

<a name="begin-with-the-in-process-store"></a>
## Begin with the in-process store

The array store gives development and unit tests a backend with the same Repository surface. Construct it once, then Put and retrieve under the same Grant. A different tenant produces a different physical key even when the application key text is identical.

```go example
type Dashboard struct{ Total int }

func CacheDashboard(ctx context.Context, g security.Grant, value Dashboard) (Dashboard, error) {
	repository := cache.New(cache.NewArrayStore())
	if err := repository.Put(ctx, g, "dashboard", value, 5*time.Minute); err != nil {
		return Dashboard{}, err
	}
	return cache.Get[Dashboard](ctx, repository, g, "dashboard")
}
```

<a name="store-remember-count-and-invalidate"></a>
## Store, remember, count, and invalidate

Put replaces a value for a bounded lifetime; Forever is the explicit unbounded choice; Has, Missing, Forget, and Touch manage presence; Increment and Decrement update counters. The generic Get and Remember helpers decode into the type the caller requests, keeping assertion logic out of handlers.

In an application the repository is already built: `foundation.NewDrivers` answers `drivers.Store`, the `cache.Repository` over the store `CACHE_STORE` names, and `drivers.Cache` holds every store by name. memory is always there; redis exists once the `hesape/redis` connector is linked in `bootstrap/app.go` and `REDIS_URL` names the endpoint. A provider hands `drivers.Store` to the services that cache.

<a name="a-cache-miss-is-expected-a-store-error-is-not"></a>
## A cache miss is expected; a store error is not

Branch on the package's missing-value behavior and compute the source value once. Propagate connection, encoding, and decoding errors instead of treating every failure as a miss: silently rebuilding during a backend outage can multiply database work until the primary service fails too. Never cache a failed or partially authorized result.

A read-only Grant, one issued while an operator only views somebody else's account, reads and never writes: every method that stores, removes, flushes, or locks an entry refuses it with `auth.ErrReadOnly`, and Remember and Flexible answer what they compute without storing it. Ask `Grant.ReadOnly` to leave a write out of such a page rather than letting it fail.

<a name="test-isolation-and-expiry-through-the-repository"></a>
## Test isolation and expiry through the repository

Write the same logical key with two tenant Grants and assert each reads only its own value. Advance or bound time to prove expiry, verify Remember calls its loader only on a miss, and run the shared store contract against the network driver before selecting it in production configuration.

<a name="repository-and-store-boundaries"></a>
## Repository and store boundaries

Repository methods cover values, batches, counters, namespaces, tags when supported, and overlapping-work locks. Store implementations cover physical persistence. `WithoutOverlapping` receives a Grant and tenant-prefixes its lock key like other application operations. `FlushLocks` is the fleet-wide exception and takes no Grant because it clears the store's lock namespace. Those boundaries do not make raw stores an application API.

<a name="a-protocol-compatible-server-is-a-deployment-choice"></a>
## A protocol-compatible server is a deployment choice

The network driver speaks the protocol shared by Redis-compatible servers without making product-specific commands part of the cache contract. Select the connection and store in configuration, keep application calls unchanged, and size TTLs from the freshness of the source data. Durable records and session state are not cache entries merely because the backend can retain bytes.
