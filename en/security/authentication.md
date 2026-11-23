# Authentication

- [Authentication turns credentials into a Subject](#authentication-turns-credentials-into-a-subject)
- [Wire the user store, password hashing, throttle, and sessions](#wire-the-user-store-password-hashing-throttle-and-sessions)
- [Verify credentials and rotate only after the final factor](#verify-credentials-and-rotate-only-after-the-final-factor)
- [Load, remember, confirm, and sign out](#load-remember-confirm-and-sign-out)
- [Return one credential failure](#return-one-credential-failure)
- [Keep enumeration and fixation closed](#keep-enumeration-and-fixation-closed)
- [Test the complete request path](#test-the-complete-request-path)
- [Authorization begins after identity](#authorization-begins-after-identity)

<a name="authentication-turns-credentials-into-a-subject"></a>
## Authentication turns credentials into a Subject

Authentication answers who is making the request. A validated request supplies normalized credentials to the authentication service; the provider retrieves the account and checks the password hash; the successful account becomes a `security.Subject`; the session store signs and persists who the person is, and each later request reads the roles again through the session resolver. Permission is a separate policy decision.

<a name="wire-the-user-store-password-hashing-throttle-and-sessions"></a>
## Wire the user store, password hashing, throttle, and sessions

The application needs a repository that scopes account lookup by the resolved tenant, an application key, factor service, session backend, and sign-in throttle. Production replicas require a shared session backend. The released application throttle is in-process memory, so each replica owns a separate failure budget; add a shared enforced control before claiming fleet-wide lockout.

<a name="verify-credentials-and-rotate-only-after-the-final-factor"></a>
## Verify credentials and rotate only after the final factor

Validate the form, then call the application's credential seam with the trusted tenant, normalized address, and plain password. If the account requires a second factor, write the short-lived signed pending attempt and create no authenticated Subject yet. Only a password-only account or a completed authenticator/recovery challenge rotates into the final session.

```go
type CredentialVerifier interface {
	Verify(context.Context, string, string, string) (security.Subject, error)
}

subject, err := credentials.Verify(r.Context(), tenant, input.Email, input.Password)
if err != nil { return err }

required, err := factors.Required(r.Context(), subject.Tenant, subject.ID)
if err != nil { return err }
if required {
	return pending.Write(w, subject, input.Remember)
}

oldID := sessions.IDFromRequest(r)
_, err = sessions.Rotate(r.Context(), w, oldID, subject, security.Remember(input.Remember))
return err
```

<a name="load-remember-confirm-and-sign-out"></a>
## Load, remember, confirm, and sign out

`RequireAuth`, `RequireRole`, and `RequireConfirmedPassword` load the signed session, read its roles again through the session resolver, and put its Subject on the request context, where handlers read it with `ctx.User`; `middleware.LoadSubject` does the same on a public route and lets a guest through with no Subject. Remember extends the session using the named option. Password confirmation records a recent proof for sensitive routes. Sign-out invalidates the current session; signing out other devices destroys the same subject's other records within the same tenant.

<a name="return-one-credential-failure"></a>
## Return one credential failure

An unknown address and a wrong password return the same public error. Lookup outages remain operational errors and do not spend a credential-failure budget. A locked attempt returns a bounded retry delay before account lookup. Invalid or expired session cookies become unauthenticated requests rather than partially populated subjects.

<a name="keep-enumeration-and-fixation-closed"></a>
## Keep enumeration and fixation closed

Throttle by both normalized identity and trusted client before the expensive password check, and hold failed comparisons to the service's timing policy. The official memory throttle splits budgets across replicas. Never log the password, session id, pending-factor payload, recovery code, or application key; rotate only after the final enabled factor and on privilege change.

The published screens count every submission of a form that sends or checks a one-time code in `drivers.Limiter`, ten an hour per address typed and thirty per fifteen minutes per client, before anything is looked up. The reset request and the verification resend answer after the same second whether the address has an account or not, and the two-factor setup screen and the recovery codes go out with Cache-Control: no-store, private.

<a name="test-the-complete-request-path"></a>
## Test the complete request path

Exercise successful sign-in, wrong password, unknown account, lockout, backend outage, remember behavior, session rotation, sign-out, and another tenant with the same normalized address. Assert public failures are indistinguishable while logs and status codes still distinguish an outage from rejected credentials for operators.

<a name="authorization-begins-after-identity"></a>
## Authorization begins after identity

A Subject is not a Grant. Once authentication identifies the actor, a handler calls `security.Authorize` for the intended action and resource, then passes the returned `security.Grant` to the data boundary. Repositories derive tenancy with `data.Tenant(g)`. CSRF independently proves that a browser mutation came through the bound session.
