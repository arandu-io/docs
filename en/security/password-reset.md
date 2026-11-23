# Password reset

- [A reset consumes one purpose-bound credential](#a-reset-consumes-one-purpose-bound-credential)
- [Prepare code storage, mail, hashing, and sessions](#prepare-code-storage-mail-hashing-and-sessions)
- [Validate, consume, write, then destroy other sessions](#validate-consume-write-then-destroy-other-sessions)
- [Keep request and completion responses non-enumerating](#keep-request-and-completion-responses-non-enumerating)
- [Do not consume a code before the candidate is valid](#do-not-consume-a-code-before-the-candidate-is-valid)
- [A reset changes every active identity boundary](#a-reset-changes-every-active-identity-boundary)
- [Run the complete released UI journey](#run-the-complete-released-ui-journey)
- [Reset uses one-time codes but is not email verification](#reset-uses-one-time-codes-but-is-not-email-verification)

<a name="a-reset-consumes-one-purpose-bound-credential"></a>
## A reset consumes one purpose-bound credential

UI :uiVersion implements the current password reset over Hesape onetime codes. The flow accepts a non-enumerating request, delivers a purpose-bound code, validates the new password, consumes the code, writes the hash, and destroys the account's other sessions.

<a name="prepare-code-storage-mail-hashing-and-sessions"></a>
## Prepare code storage, mail, hashing, and sessions

Wire Codes, the tenant-aware account rules (`*account.Users`), hasher, mailer, and the same `SessionStore` used for login. Build the reset origin from configuration and keep reset purpose distinct from verification and two-factor recovery.

- Use short-lived one-time storage.
- Send through the configured mail transport.
- Validate password and confirmation before Consume.
- Make session invalidation shared across replicas.

<a name="validate-consume-write-then-destroy-other-sessions"></a>
## Validate, consume, write, then destroy other sessions

After validation, look up the account in the trusted tenant, Consume a code bound to its current address and password fingerprint, then call the application's compare-and-swap reset. The released consumer attempts `DestroyOthers` after the password write and logs a failure rather than undoing the new password.

```go
if err := ctx.Bind(&request); err != nil { return err }
user, err := users.Lookup(ctx, tenant, request.Email)
if err != nil { return err }
fingerprint := account.Fingerprint(user.PasswordHash())
subject := user.AccountTenant() + "\x00" + user.AccountID() + "\x00" +
	account.NormalizeEmail(user.AccountEmail()) + "\x00" + fingerprint
if err := codes.Consume(ctx, "reset-password", subject, request.Code); err != nil { return err }
updated, err := users.ResetPassword(ctx, user.AccountTenant(), user.AccountID(), user.AccountEmail(), fingerprint, request.Password)
if err != nil { return err }
if err := sessions.DestroyOthers(ctx, updated.(*models.User).Subject(), ""); err != nil {
	observability.Log(ctx).Error("signing the account's other sessions out", "error", err)
}
return nil
```

<a name="keep-request-and-completion-responses-non-enumerating"></a>
## Keep request and completion responses non-enumerating

Answer reset requests the same whether an account exists, meter requests with an explicit security policy, allow a fresh request to supersede stale mail, and preserve the current newly authenticated session only when the released flow intentionally selected it.

- Use one public response for known and unknown email.
- Rate-limit request and consume separately.
- Keep code and password out of old input and logs.
- Notify the account through an independent message after success.

<a name="do-not-consume-a-code-before-the-candidate-is-valid"></a>
## Do not consume a code before the candidate is valid

Invalid password must not consume a usable code. Unknown, expired, replayed, wrong-purpose, and wrong-tenant codes fail uniformly, and after twenty wrong codes for one purpose and subject within twenty-four hours Consume refuses every code with `ErrTooManyFailures` until the window closes, however many new codes are sent; the screens of ui :uiVersion say "too many attempts, try again later" then. A write failure after Consume requires an explicit recovery path; never mark success because the code was valid.

> [!WARNING]
> The released flow is validation, Consume, password compare-and-swap, then best-effort `DestroyOthers`. A session-store outage is logged but does not undo the password change, so the current consumer cannot promise that revocation succeeded.

<a name="a-reset-changes-every-active-identity-boundary"></a>
## A reset changes every active identity boundary

> [!WARNING]
> Never accept account or tenant authority from the reset form, disclose account existence, log the code or plaintext, preserve unrelated sessions by default, or use a permanent API token as a reset credential. Use TLS and configured origins.

Reset proves control of the delivered one-time credential; policy still controls the tenant account being changed.

<a name="run-the-complete-released-ui-journey"></a>
## Run the complete released UI journey

Run known and unknown requests, valid completion, invalid confirmation, altered, expired, replayed, wrong-purpose, and cross-tenant codes, concurrent consumption, mail outage, password write failure, and two replicas sharing session state.

- Public request responses do not enumerate.
- Invalid input leaves the code usable.
- Concurrent Consume succeeds once.
- Other sessions fail after success.

<a name="reset-uses-one-time-codes-but-is-not-email-verification"></a>
## Reset uses one-time codes but is not email verification

Email verification uses another purpose, hashing protects the new password, sessions enforce invalidation, and rate limiting bounds online attempts.

- [Email verification](/docs/{{version}}/security/email-verification)
- [Hashing](/docs/{{version}}/security/hashing)
- [Session](/docs/{{version}}/http/session)
- [Rate limiting](/docs/{{version}}/services/rate-limiting)
