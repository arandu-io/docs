# OAuth and social login

- [An OAuth client completes the authorization-code exchange](#an-oauth-client-completes-the-authorization-code-exchange)
- [Register an exact callback and keep the secret at boot](#register-an-exact-callback-and-keep-the-secret-at-boot)
- [Redirect with state and refuse a mismatched callback](#redirect-with-state-and-refuse-a-mismatched-callback)
- [Normalize provider data, then apply application policy](#normalize-provider-data-then-apply-application-policy)
- [Keep refusal before the token exchange](#keep-refusal-before-the-token-exchange)
- [State is mandatory, and linking is authorization](#state-is-mandatory-and-linking-is-authorization)
- [Test the redirect and callback as one pair](#test-the-redirect-and-callback-as-one-pair)
- [The callback produces identity input, not a session](#the-callback-produces-identity-input-not-a-session)

<a name="an-oauth-client-completes-the-authorization-code-exchange"></a>
## An OAuth client completes the authorization-code exchange

The released provider package sends a browser to an external authorization endpoint, verifies the callback state, exchanges the code, and retrieves normalized `oauth.UserData`. Constructors exist for `GitHub`, Google, Facebook, Stripe, and Apple, plus `NewProvider` for an explicitly configured service.

This package is an OAuth client, not an authorization server or a complete social-login feature. The skeleton and UI kit provide no callback routes, credential configuration, account linking, tenant policy, or durable identity model for this flow; those stay application work.

<a name="register-an-exact-callback-and-keep-the-secret-at-boot"></a>
## Register an exact callback and keep the secret at boot

- Register the HTTPS callback URL with the provider exactly as the application will send it.
- Load client id and client secret from deployment configuration, never source or a browser response.
- Decide which tenant the entry route belongs to before redirecting.
- Define an account-linking policy for existing identities before accepting a callback.

Construct a new Provider and `CookieStateStore` for each request. A Provider holds request-pair state and is not safe to share between concurrent requests.

<a name="redirect-with-state-and-refuse-a-mismatched-callback"></a>
## Redirect with state and refuse a mismatched callback

```go example
var (
	clientID     = "client-id"
	clientSecret = "client-secret"
	callbackURL  = "https://example.com/auth/github/callback"
)

func githubProvider(w stdhttp.ResponseWriter, r *stdhttp.Request, id, secret, callback string) *providers.Provider {
	state := providers.NewCookieStateStore(w, r)
	return providers.NewGithubProvider(state, id, secret).RedirectURL(callback)
}

func startGithub(w stdhttp.ResponseWriter, r *stdhttp.Request) error {
	return githubProvider(w, r, clientID, clientSecret, callbackURL).Redirect(w, r)
}

func finishGithub(w stdhttp.ResponseWriter, r *stdhttp.Request) error {
	identity, token, err := githubProvider(w, r, clientID, clientSecret, callbackURL).User(r)
	if errors.Is(err, providers.ErrStateMismatch) {
		stdhttp.Error(w, "invalid oauth callback", stdhttp.StatusBadRequest)
		return nil
	}
	if err != nil {
		return err
	}
	_ = identity // Resolve or link through the tenant-aware account service.
	_ = token    // Keep secret; persist only when the provider integration needs it.
	return nil
}
```

The entry route calls `startGithub` and returns as soon as Redirect writes 302. The registered callback is a separate request handled by `finishGithub`; User verifies and consumes its state before exchanging the callback code. Treat the returned provider token as a secret even when the application needs only identity.

<a name="normalize-provider-data-then-apply-application-policy"></a>
## Normalize provider data, then apply application policy

`oauth.UserData.Get` and String read the normalized identity fields. Provider-specific attributes remain available without changing the common shape. After the callback, look up a previously linked provider subject; if none exists, enter an explicit linking or registration decision rather than silently matching any account with the same email address.

| Constructor | Released provider |
| --- | --- |
| `NewGithubProvider` | `GitHub` |
| `NewGoogleProvider` | Google |
| `NewFacebookProvider` | Facebook |
| `NewStripeProvider` | Stripe |
| `NewAppleProvider` | Apple |
| `NewProvider` | caller-supplied endpoints |

`NewAppleProvider` takes a Services ID, a Team ID, a Key ID, and the private key instead of a client secret, and makes the client secret for each exchange. Apple answers the redirect with a POST from its own site, so the callback route accepts POST and is exempt from CSRF protection for that one path with `CSRFExcept`; the state check takes its place. The site must be served over HTTPS, and the state cookie is sent SameSite=None. The package documentation of `NewAppleProvider` describes the whole flow.

<a name="keep-refusal-before-the-token-exchange"></a>
## Keep refusal before the token exchange

Missing state, missing state cookie, changed state, and replay all return `ErrStateMismatch`. A callback with `error=access_denied` is a refusal, not a server fault. Missing codes, token-exchange failures, malformed responses, and user-data failures remain distinct internal errors while the browser receives a bounded, non-secret message.

The state cookie has one default name. Starting another OAuth attempt in the same browser before the first callback can replace it, causing the first callback to fail; report that as a restartable flow rather than weakening verification.

<a name="state-is-mandatory-and-linking-is-authorization"></a>
## State is mandatory, and linking is authorization

> [!WARNING]
> `Provider.Stateless` disables the released CSRF defense. Use it only when another concrete, tested mechanism binds the callback to the initiating client; convenience is not such a mechanism.

Never auto-link solely because an external service returns an email string. Verify the provider's assurance, handle collisions deliberately, and authorize the tenant transition. Do not place client secrets or access tokens in logs, query strings, session payloads, or rendered errors.

`CookieStateStore` :hesapeVersion sends its state cookie Secure by default, whatever scheme the request arrived with, so the cookie keeps the attribute behind a proxy that terminates TLS. A development server on plain http calls `Secure(false)` on the store, because a browser never returns a Secure cookie over http; leave it on everywhere else.

<a name="test-the-redirect-and-callback-as-one-pair"></a>
## Test the redirect and callback as one pair

- Redirect returns 302 with client id, exact redirect URI, scopes, and a non-empty state.
- Callback with the stored state reaches the fake token and user-data servers.
- Missing, altered, reused, and cross-browser state fail before any exchange request.
- Provider refusal and provider outage have different operator evidence without exposing secrets.
- The deployed HTTPS and proxy path emits a Secure, `HttpOnly` state cookie, SameSite=Lax (SameSite=None for Apple).
- Linking cannot cross tenants or seize an account through an email collision.

```sh
go test github.com/arandu-io/hesape/oauth/providers
```

<a name="the-callback-produces-identity-input-not-a-session"></a>
## The callback produces identity input, not a session

After policy accepts the link, the application may turn the local account into `security.Subject` and rotate its own session. API tokens, provider access tokens, and browser sessions are separate credentials with separate lifecycles. Keep the external provider boundary outside the account repository.

- [Authentication](/docs/{{version}}/security/authentication)
- [HTTP session](/docs/{{version}}/http/session)
- [Policies and Grant](/docs/{{version}}/security/authorization)
- [OAuth providers API](https://pkg.go.dev/github.com/arandu-io/hesape/oauth/providers)
- [OAuth user data API](https://pkg.go.dev/github.com/arandu-io/hesape/oauth)
