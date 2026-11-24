# HTTP session

- [A signed cookie points to a typed Subject](#a-signed-cookie-points-to-a-typed-subject)
- [Choose the backend for the deployment topology](#choose-the-backend-for-the-deployment-topology)
- [Rotate into an authenticated session](#rotate-into-an-authenticated-session)
- [Load, confirm, remember, and end sessions](#load-confirm-remember-and-end-sessions)
- [A missing record stays unauthenticated](#a-missing-record-stays-unauthenticated)
- [Identity changes require id changes](#identity-changes-require-id-changes)
- [Prove the cross-replica path](#prove-the-cross-replica-path)
- [Authentication owns identity; policy owns authority](#authentication-owns-identity-policy-owns-authority)

<a name="a-signed-cookie-points-to-a-typed-subject"></a>
## A signed cookie points to a typed Subject

The framework session store writes a signed session id to the browser and keeps `security.Subject` in a backend. Loading verifies the cookie before reading the record, so request data cannot choose a tenant or account. A session establishes identity; it does not grant permission to application data.

<a name="choose-the-backend-for-the-deployment-topology"></a>
## Choose the backend for the deployment topology

`foundation.NewDrivers` builds the one session store, `drivers.Sessions`, from the application key, `SESSION_LIFETIME`, and the Secure decision, over the backend `SESSION_DRIVER` names. Memory is appropriate for tests and a single local process. Multiple replicas need `SESSION_DRIVER=redis`, the `hesape/redis` connector linked by its blank import in `bootstrap/app.go`, and a reachable `REDIS_URL`. Session storage is independent of `CACHE_STORE`. The cookie is `arandu_session` on path /, for the host that answered, `HttpOnly` and `SameSite` Lax; its name and scope are not settings. `SESSION_TTL`, which `SESSION_LIFETIME` replaced, stops the boot.

| Setting | Released effect |
| --- | --- |
| `SESSION_DRIVER=memory` | one in-process backend |
| `SESSION_DRIVER=redis` | shared RESP-backed session handler |
| `SESSION_LIFETIME` | lifetime in minutes, read by the framework; two hours when unset |
| `SESSION_SECURE_COOKIE` | Secure attribute, read by the framework; on unless `APP_ENV=dev`, and false is refused in production |
| `SESSION_PATH` / `SESSION_DOMAIN` / `SESSION_SAME_SITE` | not settings; any scope other than /, the answering host, and lax stops the boot |

> [!WARNING]
> Never deploy two replicas with the memory backend. Both processes can boot successfully, but a session created by one cannot be loaded by the other.

<a name="rotate-into-an-authenticated-session"></a>
## Rotate into an authenticated session

Use Rotate for sign-in even when the request has no old id. It creates a fresh id, writes the signed cookie, and destroys the pre-authentication record when one exists. Start cannot express replacement and is therefore not the login operation.

```go example
func RotateSession(r *stdhttp.Request, w stdhttp.ResponseWriter, drivers *foundation.Drivers, subject security.Subject) error {
	sessions := drivers.Sessions

	oldID := sessions.IDFromRequest(r)
	_, err := sessions.Rotate(r.Context(), w, oldID, subject)
	if err != nil {
		return fmt.Errorf("rotate session: %w", err)
	}
	return nil
}
```

<a name="load-confirm-remember-and-end-sessions"></a>
## Load, confirm, remember, and end sessions

- Load returns the trusted Subject or an explicit session error.
- Confirm records a recent password confirmation for sensitive routes.
- `Remember(true)` extends the session through the named option.
- Destroy removes the current record and clears both the session and intended-destination cookies.
- `DestroyOthers` removes the same tenant and subject's other records while preserving the selected id.

### Keep sessions in the application database

```go example
func NewSessions(drivers *foundation.Drivers, appKey []byte, ttl time.Duration) (*security.SessionStore, error) {
	handler, err := session.NewDatabaseHandler[security.Subject](drivers.DB, "")
	if err != nil {
		return nil, err
	}
	backend := security.NewSessionBackend(handler)
	sessions := security.NewSessionStore(appKey, ttl, true, backend)
	return sessions, nil
}
```

`SESSION_DRIVER=redis` resolves its connection independently of `CACHE_STORE`, so sessions can be shared while ordinary cache values remain in process memory. Hesape also ships `session.NewDatabaseHandler`, which stores each record in the table `session.RecordMigrations` creates, for more than one replica and no RESP server. It is a component rather than a setting: `SESSION_DRIVER` takes memory or redis, and the global CSRF check and throttle read `drivers.Sessions`, so a store built over the table must replace that one everywhere before it is trusted.

<a name="a-missing-record-stays-unauthenticated"></a>
## A missing record stays unauthenticated

A missing cookie, bad signature, expired record, or absent id does not create a partial Subject. Middleware treats the request as unauthenticated. A shared-backend outage is different: fail the dependent request or health check rather than silently switching to memory and splitting the fleet.

<a name="identity-changes-require-id-changes"></a>
## Identity changes require id changes

Rotate after the final authentication factor and after privilege changes. A pending second factor must not receive an authenticated Subject. Never log session ids, accept an id from form input, derive the tenant from a request parameter, or return a raw backend error containing connection credentials.

`DestroyOthers` is tenant-scoped because the trusted Subject supplies both tenant and subject id. Do not replace it with a backend scan keyed only by an account id.

The store reads a session's roles again on every request through the resolver the provider installs with `drivers.Sessions.ResolveSubjects`, so a demoted account loses its access on its next request. A resolver that answers `hhttp.ErrSubjectGone` ends the session of an account that was deleted, closed, or refused, and a session holding roles with no resolver installed is refused with `hhttp.ErrNoSubjectResolver` rather than trusted as it was written.

<a name="prove-the-cross-replica-path"></a>
## Prove the cross-replica path

Build two application instances over the same RESP server. Rotate a session through the first, replay its cookie to the second, and assert the same tenant and subject load. Then test `SESSION_DRIVER=redis` with `CACHE_STORE=memory`, an unavailable shared store, an unsupported driver, signature tampering, expiry, old-id invalidation, and `DestroyOthers` preserving only the selected id.

```sh
go test ./tests/Feature -run 'Session|Cache'
```

<a name="authentication-owns-identity-policy-owns-authority"></a>
## Authentication owns identity; policy owns authority

`RequireAuth` loads identity from this store. A handler still calls `security.Authorize` and carries the returned Grant into its repository. CSRF binds browser mutations to the current session id, so rotation also invalidates tokens issued for the former boundary.

- [Authentication](/docs/{{version}}/security/authentication)
- [Policies and Grant](/docs/{{version}}/security/authorization)
- [CSRF protection](/docs/{{version}}/http/csrf)
- [Framework session API](https://pkg.go.dev/github.com/arandu-io/framework/security)
- [Hesape RESP session adapter](https://pkg.go.dev/github.com/arandu-io/hesape/redis)
