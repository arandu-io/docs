# Email verification

- [A purpose-bound code confirms one email transition](#a-purpose-bound-code-confirms-one-email-transition)
- [Share the one-time store and configure mail delivery](#share-the-one-time-store-and-configure-mail-delivery)
- [Issue, deliver, and consume through the service](#issue-deliver-and-consume-through-the-service)
- [Rotate codes when the address or request changes](#rotate-codes-when-the-address-or-request-changes)
- [Invalid and expired codes reveal the same bounded result](#invalid-and-expired-codes-reveal-the-same-bounded-result)
- [Verification is not account linking](#verification-is-not-account-linking)
- [Exercise issuance, replay, expiry, and tenant isolation](#exercise-issuance-replay-expiry-and-tenant-isolation)
- [Finish verification before granting the final session](#finish-verification-before-granting-the-final-session)

<a name="a-purpose-bound-code-confirms-one-email-transition"></a>
## A purpose-bound code confirms one email transition

Hesape onetime codes implement the released primitive and UI :uiVersion supplies the application flow. The service issues a purpose-bound six-digit code, mail delivers its only clear copy, and Consume atomically accepts it once before the account's verified state is written.

<a name="share-the-one-time-store-and-configure-mail-delivery"></a>
## Share the one-time store and configure mail delivery

Configure the shared `CodeStore`, cryptographic key, mail transport, and tenant-aware user service. Keep verification purpose separate from password reset and two-factor codes even when they share the same storage implementation.

- Use an exact purpose and tenant.
- Render the code in the verification mailable.
- Keep raw codes out of durable logs.
- Send through the application mail bridge.

<a name="issue-deliver-and-consume-through-the-service"></a>
## Issue, deliver, and consume through the service

Construct the subject from trusted tenant, user id, and normalized captured address. Issue and mail its code, then Consume that same purpose and subject before `MarkVerified` performs its conditional account write. A repeated submission cannot reuse the spent code.

```go
const verifyPurpose = "verify-email"
subject := user.AccountTenant() + "\x00" + user.AccountID() + "\x00" + account.NormalizeEmail(user.AccountEmail())

code, err := codes.Issue(ctx, verifyPurpose, subject)
if err != nil { return err }
if err := mailer.Send(ctx, VerifyEmail{User: user, Code: code}); err != nil { return err }

if err := codes.Consume(ctx, verifyPurpose, subject, submitted); err != nil { return err }
_, _, err = users.MarkVerified(ctx, user.AccountTenant(), user.AccountID(), user.AccountEmail())
return err
```

<a name="rotate-codes-when-the-address-or-request-changes"></a>
## Rotate codes when the address or request changes

Reissue after an address change, supersede former codes according to the store contract, keep repeated form submissions bounded, and require a new sign-in after verification rather than turning code consumption into a session.

- Bind purpose, tenant, account, and expiry.
- Consume before recording success.
- Send a replacement without disclosing account existence.
- Keep mail delivery failure observable and retryable.

<a name="invalid-and-expired-codes-reveal-the-same-bounded-result"></a>
## Invalid and expired codes reveal the same bounded result

Unknown, altered, expired, already consumed, wrong-purpose, and wrong-tenant codes all fail without exposing which account exists. One purpose and subject may present twenty wrong codes in twenty-four hours across every code issued; after that Consume refuses every code with `ErrTooManyFailures` until the window closes, and issuing a new code does not reset the count. The screens of ui :uiVersion answer that refusal with "too many attempts, try again later" rather than calling the code invalid. A mail outage is separate from an invalid POST; a database failure after Consume needs an idempotent repair path.

Do not mark verified from an account id alone. The released POST flow relies on the one-time store, purpose, and captured-address boundary.

<a name="verification-is-not-account-linking"></a>
## Verification is not account linking

> [!WARNING]
> Never auto-link accounts solely by matching an email string. Do not log raw codes, put them in analytics, URLs, or redirect targets. Accept the code only in a CSRF-protected POST. Verification confirms control of one delivered address; policy still decides tenant membership and account changes.

Use TLS and a short expiry. A browser session may carry flow state, but the durable one-time record is what prevents replay across processes.

<a name="exercise-issuance-replay-expiry-and-tenant-isolation"></a>
## Exercise issuance, replay, expiry, and tenant isolation

Run the UI form against Hesape Codes with valid, altered, expired, consumed, wrong-purpose, wrong-tenant, replacement, mail failure, and cross-replica submissions. Assert the account changes once and no response reveals another account.

- Issue and Consume succeed once.
- Replay and expiry fail.
- Another purpose or tenant cannot consume.
- The UI consumer writes verified state only after Consume.

<a name="finish-verification-before-granting-the-final-session"></a>
## Finish verification before granting the final session

Authentication determines when verification is required, password reset uses a separate purpose, sessions rotate after identity transitions, and mail transports the code.

- [Authentication](/docs/{{version}}/security/authentication)
- [Password reset](/docs/{{version}}/security/password-reset)
- [Session](/docs/{{version}}/http/session)
- [Mail](/docs/{{version}}/services/mail)
