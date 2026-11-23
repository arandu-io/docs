# Notifications

- [A notification selects channels around one public message](#a-notification-selects-channels-around-one-public-message)
- [Own recipient preferences, addresses, and transports](#own-recipient-preferences-addresses-and-transports)
- [Build a message without storing credentials](#build-a-message-without-storing-credentials)
- [Separate immediate and durable delivery](#separate-immediate-and-durable-delivery)
- [Treat partial channel delivery as explicit state](#treat-partial-channel-delivery-as-explicit-state)
- [A recipient address is protected data](#a-recipient-address-is-protected-data)
- [Test preference, tenant, and retry paths](#test-preference-tenant-and-retry-paths)
- [Use mail, queues, and outbox as their real boundaries](#use-mail-queues-and-outbox-as-their-real-boundaries)

> [!WARNING]
> **Current limitations.** Notification primitives are released, but there is no complete first-party fan-out journey with durable authorization and retry semantics. Keep channel delivery explicit and use the transactional outbox for work that must survive process failure.

<a name="a-notification-selects-channels-around-one-public-message"></a>
## A notification selects channels around one public message

The released notification primitives describe a notification and its channels, but no tagged first-party application proves a complete preference-to-transport journey. Use them as a limited composition surface and keep actual mail, queue, and recipient policy explicit.

<a name="own-recipient-preferences-addresses-and-transports"></a>
## Own recipient preferences, addresses, and transports

Resolve a trusted recipient through the application service, read channel preferences inside tenant scope, configure each transport, and decide which deliveries must be queued. Do not accept addresses, tenant, or enabled channels from an untrusted job payload.

- Load recipient through a Grant.
- Apply opt-out preferences at delivery time.
- Inject concrete transports.
- Choose an idempotency key per channel.

<a name="build-a-message-without-storing-credentials"></a>
## Build a message without storing credentials

Give the notification a stable Key, return its channel names from Via, and implement the representation interface for every selected channel. Via alone does not supply a mail or database body. Pass the complete value to the injected Notifier, which resolves each wired channel and continues collecting failures.

```go example
type InvoicePaid struct {
	InvoiceID string
	URL       string
}

func (InvoicePaid) Key() notifications.Key { return "billing.invoice-paid" }
func (InvoicePaid) Via(notifications.Notifiable) []notifications.ChannelName {
	return []notifications.ChannelName{
		notifications.ChannelMail,
		notifications.ChannelDatabase,
	}
}

func (n InvoicePaid) ToMail(notifications.Notifiable) messages.Mail {
	return messages.NewMail().
		Subject("Your invoice is paid").
		Line("We received your payment.").
		Action("View invoice", n.URL)
}

func (n InvoicePaid) ToDatabase(notifications.Notifiable) messages.Database {
	return messages.NewDatabase(map[string]any{"invoice_id": n.InvoiceID})
}

func NotifyPaid(ctx context.Context, notifier *notifications.Notifier, grant security.Grant, recipient notifications.Notifiable, invoiceID, invoiceURL string) error {
	notice := InvoicePaid{InvoiceID: invoiceID, URL: invoiceURL}
	return notifier.Send(ctx, grant, recipient, notice)
}
```

<a name="separate-immediate-and-durable-delivery"></a>
## Separate immediate and durable delivery

Resolve preferences just before delivery, render channel-specific public resources, queue slow network transports, attach a stable delivery id, and record each channel outcome without turning another channel's success into proof of all delivery.

- Keep channel selection tenant-scoped.
- Use mail for email representation.
- Use queue for durable retry.
- Record one outcome per recipient and channel.

<a name="treat-partial-channel-delivery-as-explicit-state"></a>
## Treat partial channel delivery as explicit state

One channel can succeed while another fails. Decide retry state per channel, avoid resending successful side effects, and distinguish invalid address, explicit opt-out, provider refusal, and provider outage. The component package alone does not supply this end-to-end lifecycle.

> [!WARNING]
> The current released contract does not include a complete application notification subsystem. Before relying on one, prove recipient preference, queue, transport, idempotency, and operator recovery together in the tagged application you deploy.

<a name="a-recipient-address-is-protected-data"></a>
## A recipient address is protected data

> [!WARNING]
> Treat email, phone, device tokens, and message content as protected tenant data. Never log credentials or reset codes, allow a payload to override the recipient, or bypass a user's current opt-out at retry time.

Authorization decides access to the event; notification preference decides whether and where to send it. Neither substitutes for the other.

<a name="test-preference-tenant-and-retry-paths"></a>
## Test preference, tenant, and retry paths

Assert that every name returned by Via has its required representation method, then exercise enabled and disabled channels, a foreign-tenant recipient, invalid address, provider outage after another channel succeeds, duplicate queue delivery, redacted logs, and preference changes between enqueue and handle.

- Opt-out prevents its channel.
- Foreign recipients cannot be selected.
- Duplicate delivery is idempotent.
- Partial failure preserves successful channel state.

<a name="use-mail-queues-and-outbox-as-their-real-boundaries"></a>
## Use mail, queues, and outbox as their real boundaries

Mail formats email, queues provide durable delivery, events and outbox persist facts, and authorization protects recipient lookup.

- [Mail](/docs/{{version}}/services/mail)
- [Queues](/docs/{{version}}/async/queues)
- [Events and outbox](/docs/{{version}}/async/events-outbox)
- [Authorization](/docs/{{version}}/security/authorization)
