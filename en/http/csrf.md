# CSRF protection

- [Why a state-changing request needs a second secret](#why-a-state-changing-request-needs-a-second-secret)
- [Install one issuer around the route stack](#install-one-issuer-around-the-route-stack)
- [Protect the handler before it can run](#protect-the-handler-before-it-can-run)
- [Cover forms and enhanced requests](#cover-forms-and-enhanced-requests)
- [Refuse without guessing](#refuse-without-guessing)
- [Rotation closes the old boundary](#rotation-closes-the-old-boundary)
- [Test the rejection as carefully as the success](#test-the-rejection-as-carefully-as-the-success)
- [CSRF is one layer of request trust](#csrf-is-one-layer-of-request-trust)

<a name="why-a-state-changing-request-needs-a-second-secret"></a>
## Why a state-changing request needs a second secret

A browser attaches a session cookie to a matching request automatically. That convenience also lets another site cause the request, so a signed-in cookie is not proof that the page meant to submit it. The CSRF token is tied to the session id and must arrive through markup or a header the other site cannot read.

<a name="install-one-issuer-around-the-route-stack"></a>
## Install one issuer around the route stack

`foundation.NewDrivers` builds the one issuer, `drivers.CSRF`, from the application key and `CSRF_TTL`, and binds it to the session store, `drivers.Sessions`, so the session id is accepted only after its signed cookie verifies. Create installs the check in the global pipeline; nothing has to be wired by hand. A guest's token must be bound to a random id carried in that guest's own signed cookie, never to an empty session id that every visitor would share: a token one visitor holds must be refused on another visitor's form. Production session cookies still require Secure over HTTPS.

<a name="protect-the-handler-before-it-can-run"></a>
## Protect the handler before it can run

The check runs in the global pipeline after the throttle and before every mutation handler. Safe GET, HEAD, OPTIONS, and TRACE requests pass. Unsafe requests read X-CSRF-Token before the `_token` form field and validate against the current verified session id. The one decision left to the application is which of its own paths carry a guard of their own, named in `WithMiddleware`.

```go example
func WithoutCSRFOnWebhooks(drivers *foundation.Drivers) *foundation.ApplicationBuilder {
	return foundation.Configure(drivers).
		WithMiddleware(func(m *foundation.Middleware) {
			m.ValidateCSRFTokensExcept(
				"/webhooks/",                 // a prefix: every path under it
				"/v1/usage",                  // exactly this path
				"/checkout/{token}/sessions", // one segment of any value
			)
		})
}
```

<a name="cover-forms-and-enhanced-requests"></a>
## Cover forms and enhanced requests

A conventional form carries `@csrf`, which compiles to the hidden field. The application layout also places the token in a ypy-csrf meta tag, and the ypy-csrf extension sends it as X-CSRF-Token on every same-origin request that changes state, so partial-page POST, PUT, PATCH, and DELETE requests without a form do not invent token plumbing. When a response replaces a form, render a fresh field into the replacement.

An exemption leaves a path to its own guard, such as a webhook that verifies the provider's signature. An entry ending in / covers every path under it, an entry without one covers exactly that path, and a `{name}` segment matches any one non-empty segment, as a route pattern does; an entry may not begin with one. A module that implements CSRFExempt declares its receivers itself: Create exempts them beside the application's own, Boot refuses one that covers a write route of another owner, and every collected exemption is written to the log at boot.

```kyse
<form method="post" action="/profile">
  @csrf
</form>
```

<a name="refuse-without-guessing"></a>
## Refuse without guessing

Do not mint a replacement while handling the rejected mutation and continue. A browser-reported cross-origin mutation returns 403. A missing, invalid, expired, or other-session token returns 419; a YPY request receives the framework refresh instruction, YPY-Refresh. The next rendered page supplies a current token while the rejected mutation remains untouched.

<a name="rotation-closes-the-old-boundary"></a>
## Rotation closes the old boundary

Rotate the session id when authentication authority changes. A token bound to the former id then stops working with the new session. Never put tokens in URLs or logs, and do not exempt a state-changing endpoint merely because it normally receives an enhanced request; an attacker can call the endpoint directly.

<a name="test-the-rejection-as-carefully-as-the-success"></a>
## Test the rejection as carefully as the success

Exercise a valid field, a valid header, a missing token, a modified token, a token from another session, and an expired token. For rendered forms, count state-changing forms and their token directives so a newly added form cannot silently omit protection. Also assert that rejection leaves the protected record unchanged.

<a name="csrf-is-one-layer-of-request-trust"></a>
## CSRF is one layer of request trust

The token proves that a request came through a page holding the session; it does not decide who the subject is or what that subject may do. Session authentication, policy authorization, input validation, and same-origin response headers keep their separate responsibilities and still run on a request whose token validates.
