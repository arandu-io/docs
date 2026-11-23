# Middleware

- [One wrapper, one decision](#one-wrapper-one-decision)
- [Install the production envelope](#install-the-production-envelope)
- [Choose the narrowest middleware scope](#choose-the-narrowest-middleware-scope)
- [Middleware must preserve the protocol](#middleware-must-preserve-the-protocol)
- [Prove order and terminal branches](#prove-order-and-terminal-branches)
- [Released middleware API](#released-middleware-api)
- [Put policy at the right boundary](#put-policy-at-the-right-boundary)

<a name="one-wrapper-one-decision"></a>
## One wrapper, one decision

HTTP middleware has the standard `func(http.Handler) http.Handler` shape. The framework composes the global pipeline from the configuration and runs it in one fixed order on the way in and in reverse on the way out; route middleware wraps only the route or group it protects.

<a name="install-the-production-envelope"></a>
## Install the production envelope

The production envelope is already installed. Create composes recovery, the body limits, trusted-proxy normalization, observability, security headers, a global throttle of three hundred requests a minute per session or address, the CSRF check, and method override from the configuration. `bootstrap/app.go` adjusts that pipeline in `WithMiddleware`: it adds middleware at one of three named places and changes the few decisions the pipeline leaves to the application.

```go example
// Stand-ins: each wraps next with the behavior its name describes.
func requestTimer(next stdhttp.Handler) stdhttp.Handler        { return next }
func noStoreOnBearerLinks(next stdhttp.Handler) stdhttp.Handler { return next }
func tenantBanner(next stdhttp.Handler) stdhttp.Handler         { return next }

func WithPipeline(drivers *foundation.Drivers) *foundation.ApplicationBuilder {
	return foundation.Configure(drivers).
		WithMiddleware(func(m *foundation.Middleware) {
			m.Prepend(requestTimer)
			m.BeforeThrottle(noStoreOnBearerLinks)
			m.Append(tenantBanner)
			m.Throttle(cache.PerMinute(600))
			m.ValidateCSRFTokensExcept("/webhooks/", "/checkout/{token}/sessions")
			m.ImageOrigins("https://images.example.com")
		})
}
```

> [!NOTE]
> `TrustProxies` reads X-Forwarded-For only from a peer `TRUSTED_PROXIES` lists, walking the header from the right and skipping trusted hops, so a value a client wrote at the left never becomes its address. Everything that derives identity or a limit from the client address runs after it: `BeforeThrottle`, the global Throttle, and Append.

<a name="choose-the-narrowest-middleware-scope"></a>
## Choose the narrowest middleware scope

The order is fixed: the refusal Sentences, Recover, Prepend, `ValidatePostSize` and `LimitBodySize`, `TrustProxies`, Observe, `SecurityHeaders`, `BeforeThrottle`, Throttle, the CSRF check, `OverrideMethod`, and Append, right before the router. Attach authentication and route-specific controls to a `Router.Group`, `Router.Action`, or verb method. `http.Chain` composes a standalone standard handler with the first middleware outermost.

- Prepend adds middleware right after Recover, before the body limits.
- `BeforeThrottle` holds on a refused request too and reads the real address.
- Append runs after `OverrideMethod`, right before the router.
- Throttle replaces the global budget; Sentences and Observer replace the refusal words and Observe in place.
- Authentication belongs only on protected route groups.

<a name="middleware-must-preserve-the-protocol"></a>
## Middleware must preserve the protocol

A middleware calls `next.ServeHTTP` at most once, returns after a terminal response, sets headers before the first write, and preserves `r.WithContext`. Breaking those rules can duplicate work, discard cancellation and identity, or emit a response with the wrong headers. Recover stays outermost after the sentences: a panic in middleware above it escapes without a page.

> [!WARNING]
> The released Throttle panics for invalid construction and fails open when its backing store errors, logging the failure. Do not present it as a fail-closed sign-in or password-reset protection without an application policy that refuses store outages.

<a name="prove-order-and-terminal-branches"></a>
## Prove order and terminal branches

Dispatch through the composed application and assert where an added middleware runs relative to the throttle and the CSRF check. Prove an untrusted peer cannot forge a forwarded address, a trusted peer can, a terminal refusal never reaches the handler, the 301st request in a minute is refused, and an escaped panic becomes the production status.

- Inbound and outbound order is deterministic.
- Terminal middleware calls no downstream handler.
- Request cancellation survives context enrichment.
- Proxy and throttle failures follow their documented policies.

<a name="released-middleware-api"></a>
## Released middleware API

`foundation.ApplicationBuilder.WithMiddleware` hands over `*foundation.Middleware`, whose Prepend, `BeforeThrottle`, Append, Throttle, `ValidateCSRFTokensExcept`, `ImageOrigins`, Observer, and Sentences adjust the composed pipeline. `http.Chain` composes standard handlers. Framework middleware supplies Observe, `RequireAuth`, and CSRFProtect; Hesape HTTP middleware supplies trusted-proxy normalization and security headers, and the throttle is `hesape/routing/middleware.Throttle` over `drivers.Limiter`.

| Surface | Released contract |
| --- | --- |
| `WithMiddleware` | adjust the pipeline the framework composes |
| Prepend / `BeforeThrottle` / Append | the three places an application adds middleware |
| Throttle / `ImageOrigins` / `ValidateCSRFTokensExcept` | the global budget, image origins, and CSRF exemptions |
| Observer / Sentences | replace Observe and the refusal words in place |
| `RequireAuth` and CSRFProtect | browser identity and mutation checks |

<a name="put-policy-at-the-right-boundary"></a>
## Put policy at the right boundary

Routing chooses scope, authentication loads identity, CSRF protects browser mutations, and rate limiting meters a named key. Middleware composes these boundaries but does not merge their meanings.

- [Routing](/docs/{{version}}/http/routing)
- [Authentication](/docs/{{version}}/security/authentication)
- [CSRF](/docs/{{version}}/http/csrf)
