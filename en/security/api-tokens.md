# API tokens

- [`TokenGuard` authenticates one presented credential](#tokenguard-authenticates-one-presented-credential)
- [Resolve bearer tokens with `RequireToken`](#resolve-bearer-tokens-with-requiretoken)
- [Own issuance, storage, and the request adapter](#own-issuance-storage-and-the-request-adapter)
- [Authenticate a pre-provisioned bearer token](#authenticate-a-pre-provisioned-bearer-token)
- [Read credentials from headers only](#read-credentials-from-headers-only)
- [Unknown tokens remain ordinary 401 responses](#unknown-tokens-remain-ordinary-401-responses)
- [A permanent credential needs an application kill switch](#a-permanent-credential-needs-an-application-kill-switch)
- [Test what the guard does and what it does not do](#test-what-the-guard-does-and-what-it-does-not-do)
- [Use sessions for browsers and build lifecycle separately](#use-sessions-for-browsers-and-build-lifecycle-separately)

<a name="tokenguard-authenticates-one-presented-credential"></a>
## `TokenGuard` authenticates one presented credential

The released `auth.TokenGuard` reads a bearer token, asks an application-owned `UserProvider` for the matching account, and exposes that account through the Guard contract. It has no session and proves the credential again on every request.

`TokenGuard` is a guard, not a token system. It does not define expiration, scopes, rotation, or revocation, and a token it accepts stays acceptable for as long as the provider keeps a matching value. The application issues, stores, and revokes its tokens.

<a name="resolve-bearer-tokens-with-requiretoken"></a>
## Resolve bearer tokens with `RequireToken`

For an API route group, `middleware.RequireToken` from `github.com/arandu-io/hesape/http/middleware` reads the bearer token from the Authorization header and hands the application's `TokenResolver` only its SHA-256 digest, as `middleware.DigestToken` computes it. A missing token, an unknown one, and one the resolver reports as `ErrUnknownToken` (never issued, revoked, expired) all get the same 401. Any other resolver error is a failure to answer and ends as a 500, so a store outage does not tell a client with a good token that it was taken away. The route never reads the query string, a cookie, or the session.

A request that passes carries the resolved `auth.Subject` on its context, exactly as behind `RequireAuth`, and the Policy still decides every record. The skeleton's `PersonalAccessTokenService` is the issuing side: it shows a token once and stores only its digest, accepts an optional expiry, revokes one token or all of an account's tokens (all of them when the account is closed or its password is reset), and resolves a digest to the account's current Subject.

<a name="own-issuance-storage-and-the-request-adapter"></a>
## Own issuance, storage, and the request adapter

Before wiring the guard, the application must generate a high-entropy opaque token, show it once, store only its SHA-256 digest, and implement `auth.UserProvider` over that digest. It must also adapt the incoming request to `auth.Request`; the released HTTP request types do not yet satisfy every method of that interface directly.

- Generate credentials with `crypto/rand`, not names, timestamps, or UUIDs chosen for identifiers.
- Store a digest in a dedicated credential record, not a reusable password or remember-token column.
- Decide tenant binding and disabled state in the provider lookup.
- Provide a separate operator path for replacement and removal; `TokenGuard` does not do either.

<a name="authenticate-a-pre-provisioned-bearer-token"></a>
## Authenticate a pre-provisioned bearer token

This example covers only request authentication against a token that already exists. Embedding the Hesape request supplies query, input, bearer, and cookie accessors; the three methods below complete the exact `auth.Request` contract.

```go example
type tokenRequest struct{ *hhttp.Request }

func (r tokenRequest) GetUser() string {
	user, _, _ := r.RawRequest().BasicAuth()
	return user
}
func (r tokenRequest) GetPassword() string {
	_, password, _ := r.RawRequest().BasicAuth()
	return password
}
func (r tokenRequest) Context() context.Context {
	return r.RawRequest().Context()
}

func Authenticate(w stdhttp.ResponseWriter, r *stdhttp.Request, users auth.UserProvider) {
	request := tokenRequest{Request: hhttp.NewRequest(r)}
	guard := auth.NewTokenGuard(users, request, "api_token", "api_token_sha256", true)
	user := guard.User()
	if user == nil {
		stdhttp.Error(w, "unauthorized", stdhttp.StatusUnauthorized)
		return
	}
}
```

> [!NOTE]
> With hash=true, the User path hashes the wire token with SHA-256 before calling `RetrieveByCredentials`. The provider's storage field must contain the lowercase hexadecimal digest expected by that path.

<a name="read-credentials-from-headers-only"></a>
## Read credentials from headers only

`GetTokenForRequest` prefers Authorization: Bearer and falls back to the HTTP Basic password. It deliberately ignores query strings and request bodies, keeping credentials out of URLs, browser history, referrer headers, and ordinary access logs. Construct a new `TokenGuard` for every request. In :hesapeVersion, `SetRequest` replaces the request but does not clear the user already cached by User, so it is not a safe rebind after resolution.

The Basic-password fallback cannot be disabled by a `TokenGuard` option. If the endpoint must accept bearer credentials only, reject non-Bearer authorization before invoking the guard, or use `RequireToken`, which reads the Bearer scheme only.

<a name="unknown-tokens-remain-ordinary-401-responses"></a>
## Unknown tokens remain ordinary 401 responses

Missing, malformed, unknown, and disabled tokens produce the same unauthenticated response. In :hesapeVersion, User also collapses every `RetrieveByCredentials` error to nil, so the shown route fails closed as 401 during a provider outage but cannot distinguish that outage from an unknown token. Record the lookup failure inside the provider, or use an application-owned authentication seam that preserves the error when the protocol must return an operational status. Never log the presented token while diagnosing a lookup.

In Hesape :hesapeVersion, `TokenGuard.Validate` does not apply the configured SHA-256 transformation when hash=true, although User does. Do not use Validate as equivalent to User in hashed mode. The released test proves hashing only through User and Check.

<a name="a-permanent-credential-needs-an-application-kill-switch"></a>
## A permanent credential needs an application kill switch

A digest protects a database disclosure from immediately revealing the wire token, but it does not provide expiry or revocation semantics by itself; the skeleton's token service adds an optional expiry and deletes the row on revoke. The provider should refuse disabled records and bind every token to the expected tenant and account. Use TLS, keep tokens out of URLs and logs, compare the stored digest through indexed lookup, and rotate by creating a new credential before removing the old one.

Do not interpret authentication as authorization. Convert the authenticated account into the application's Subject and run the same Policy that a session-authenticated request uses.

<a name="test-what-the-guard-does-and-what-it-does-not-do"></a>
## Test what the guard does and what it does not do

- A Bearer header resolves the account and Check succeeds.
- The query string and body never authenticate.
- A stored digest resolves only when hash=true through User.
- Missing, altered, and disabled tokens return 401 without echoing the credential.
- A provider error becomes nil and the same 401 through `TokenGuard.User`, while provider telemetry records the outage.
- A guard that resolved one user is never reused through `SetRequest` for another request.
- Policy denial still blocks an authenticated token.

```sh
go test github.com/arandu-io/hesape/auth -run TokenGuard
```

<a name="use-sessions-for-browsers-and-build-lifecycle-separately"></a>
## Use sessions for browsers and build lifecycle separately

Browser authentication uses signed sessions and CSRF protection. `TokenGuard` is appropriate only where an application has already supplied the missing credential lifecycle. OAuth is a client flow for obtaining provider identity and is not a replacement for an application's API-token issuer.

- [Authentication](/docs/{{version}}/security/authentication)
- [Policies and Grant](/docs/{{version}}/security/authorization)
- [OAuth and social login](/docs/{{version}}/security/oauth)
- [`TokenGuard` API](https://pkg.go.dev/github.com/arandu-io/hesape/auth#TokenGuard)
