# Mail

- [A mailable owns representation and the transport owns delivery](#a-mailable-owns-representation-and-the-transport-owns-delivery)
- [Configure one mailer and public origin at boot](#configure-one-mailer-and-public-origin-at-boot)
- [Build and send one typed message](#build-and-send-one-typed-message)
- [Queue slow delivery with tenant provenance](#queue-slow-delivery-with-tenant-provenance)
- [Keep render, attachment, and transport failures distinct](#keep-render-attachment-and-transport-failures-distinct)
- [Mail is an external disclosure boundary](#mail-is-an-external-disclosure-boundary)
- [Capture the exact MIME contract](#capture-the-exact-mime-contract)
- [Use mail inside purpose-bound application flows](#use-mail-inside-purpose-bound-application-flows)

> [!WARNING]
> **Current limitations.** The Hesape mail component and the framework bridge render and deliver the documented messages. Attachment storage and asynchronous authorization are not complete end to end; attach only data already authorized and queue mail through persisted trusted provenance.

<a name="a-mailable-owns-representation-and-the-transport-owns-delivery"></a>
## A mailable owns representation and the transport owns delivery

Hesape :hesapeVersion owns mail messages and transports, while `github.com/arandu-io/framework/mail` :frameworkVersion supplies the compatibility bridge used by the application. The example below uses that bridge's `Address{Email: ...}` and `ToAddress` surface; Hesape's lower-level Mailer uses To instead.

<a name="configure-one-mailer-and-public-origin-at-boot"></a>
## Configure one mailer and public origin at boot

`foundation.NewDrivers` builds `drivers.Mail` from `MAIL_URL`, `MAIL_FROM_ADDRESS`, and `MAIL_FROM_NAME`. `smtp://` upgrades with STARTTLS when the server offers it, `smtps://` speaks implicit TLS on port 465, and `resend://` and `sendgrid://` reach those providers. `log://` writes the whole message into the log and is refused in production, and so is an unset `MAIL_URL`, which means the same. A provider injects `drivers.Mail` into services or queued handlers rather than reading environment or opening SMTP from a controller.

- Use configured From identity.
- Build links from trusted `APP_URL`.
- Keep transport credentials secret.
- Decide immediate versus durable delivery.

<a name="build-and-send-one-typed-message"></a>
## Build and send one typed message

Build a typed message from public application data and send it through the injected transport with context. The controller or domain service supplies recipient and purpose; the transport never discovers accounts or tenant policy.

```go example
type VerifyEmail struct{ Code string }

func (VerifyEmail) Envelope() mail.Envelope {
	return mail.Envelope{Subject: "Verify your address"}
}
func (m VerifyEmail) Content() mail.Content {
	return mail.Content{Text: "Use this code: " + m.Code}
}

type Recipient struct{ Email, Name string }

func SendVerification(ctx context.Context, mailer *mail.Mailer, recipient Recipient, code string) error {
	return mailer.ToAddress(mail.Address{
		Email: recipient.Email, Name: recipient.Name,
	}).Send(ctx, VerifyEmail{Code: code})
}
```

<a name="queue-slow-delivery-with-tenant-provenance"></a>
## Queue slow delivery with tenant provenance

Queue network delivery when the response need not wait, preserve tenant and action in the job Grant, render text and HTML consistently, use bounded application-owned attachments, and keep `log://` only for development evidence. To put an activity log or a filter in front of every message, call `drivers.WrapMail` before `drivers.Mail` is handed to anything: a service that took the mailer earlier keeps the one it took.

- Pass context into delivery.
- Use one public origin.
- Queue with Grant provenance.
- Attach only authorized bounded content.

<a name="keep-render-attachment-and-transport-failures-distinct"></a>
## Keep render, attachment, and transport failures distinct

Rendering, address validation, attachment open or read, TLS, authentication, provider refusal, and timeout each need operator evidence. Do not mark a verification or reset message sent when only construction succeeded.

> [!WARNING]
> The released bridge does not make every attachment or Grant path automatically safe. Open attachments inside the authorized delivery operation, bound bytes, and ensure queued handlers reconstruct current policy rather than trusting payload paths.

<a name="mail-is-an-external-disclosure-boundary"></a>
## Mail is an external disclosure boundary

> [!WARNING]
> Never log SMTP credentials, one-time codes in production, full recipient lists, attachment contents, or provider response bodies containing personal data. Prevent header injection, use TLS, and derive recipient and tenant through application services.

A delivered link or code is still purpose-bound application input. Mail delivery itself does not authenticate the callback or authorize an account change.

<a name="capture-the-exact-mime-contract"></a>
## Capture the exact MIME contract

Use a capture transport for From, To, subject, text and HTML, headers, attachments, cancellation, provider refusal, and redaction. Run the application bridge and queued path, not only message-builder unit tests.

- Captured MIME has expected public fields.
- A transport failure remains retryable evidence.
- Cancellation stops delivery.
- Secrets are absent from logs.

<a name="use-mail-inside-purpose-bound-application-flows"></a>
## Use mail inside purpose-bound application flows

Email verification and password reset define purposes, queues provide durable delivery, notifications compose channels, and configuration selects the transport.

- [Email verification](/docs/{{version}}/security/email-verification)
- [Password reset](/docs/{{version}}/security/password-reset)
- [Queues](/docs/{{version}}/async/queues)
- [Configuration](/docs/{{version}}/architecture/configuration)
