# Events and transactional outbox

- [Commit business state and its facts together](#commit-business-state-and-its-facts-together)
- [Wire the table, relay, publisher, and replica lock](#wire-the-table-relay-publisher-and-replica-lock)
- [Store the event inside the write transaction](#store-the-event-inside-the-write-transaction)
- [Drain through the production publisher path](#drain-through-the-production-publisher-path)
- [Design every consumer for repeated delivery](#design-every-consumer-for-repeated-delivery)
- [Tenant and authority come from the Grant](#tenant-and-authority-come-from-the-grant)
- [Prove rollback, replay, and replica behavior](#prove-rollback-replay-and-replica-behavior)
- [Continue with queue workers and scheduling](#continue-with-queue-workers-and-scheduling)

<a name="commit-business-state-and-its-facts-together"></a>
## Commit business state and its facts together

Use the durable outbox for a fact that must leave the process after a database write. `events.Outbox.Store` takes the same context and Grant, records tenant, action, and authorizer, and refuses calls outside `data.Transaction` with `ErrNoTransaction`. In-process dispatch cannot close the commit-to-publish crash window.

<a name="wire-the-table-relay-publisher-and-replica-lock"></a>
## Wire the table, relay, publisher, and replica lock

Build the outbox with `events.NewOutbox` over `drivers.DB` in a provider, hand `events.NewRelay` an `events.Listeners` list whose entries `events.On` builds, one per event name, and list `events.WithRelay(relay)` in `WithProviders` rather than only the table-owning `NewModule`. Multiple replicas require a distributed lock: pass `drivers.Locker`, which is nil only when every store is kept in-process, and that is right for a single serving process alone.

- Migrate the outbox table.
- Inject one production Publisher.
- Register the relay lifecycle.
- Use a shared lock across serving replicas.

<a name="store-the-event-inside-the-write-transaction"></a>
## Store the event inside the write transaction

Write the business change and Store its immutable fact in the same SQL transaction. If either operation fails, neither commits; consumers can reload by id but must not assume mutable state still matches the occurred event.

```go example
type Invoices interface {
	MarkPaid(ctx context.Context, grant security.Grant, id string) error
}

func MarkPaid(ctx context.Context, db *data.DB, invoices Invoices, outbox *events.Outbox, grant security.Grant, invoiceID string) error {
	return data.Transaction(ctx, db, func(ctx context.Context) error {
		if err := invoices.MarkPaid(ctx, grant, invoiceID); err != nil {
			return err
		}
		return outbox.Store(ctx, grant, []events.Event{{
			Name: "invoice.paid", Aggregate: "invoice", AggregateID: invoiceID,
			Payload: map[string]any{"invoice_id": invoiceID}, OccurredAt: time.Now().UTC(),
		}})
	})
}
```

<a name="drain-through-the-production-publisher-path"></a>
## Drain through the production publisher path

`Relay.Drain` processes pending rows through the same publisher as the serving loop. Inspect pending, parked, retry, lag, and marking operations; default interval is one second, batch 100, ten attempts, and a thirty-second relay lock unless configured. For an external HTTP consumer, compose `webhook.Manager` with the database, secret resolver, and bounded options, and adapt the existing outbox Publisher instead of creating a second outbox. Register hesape's webhook module with `foundation.Component` in `WithProviders`: Create collects its job handlers and the CSRF exemptions of its receivers, so the application writes neither a `RegisterJobHandlers` call nor a `ValidateCSRFTokensExcept` line for them, and calling `RegisterJobHandlers` as well registers each job twice and stops the boot. Its management prefix and its receive prefix may not overlap, because the receive prefix is exempt from the CSRF check.

- Drain through the production Publisher.
- Monitor pending and lag.
- Park after the attempt budget.
- Retry only after fixing the consumer or transport.
- Register each webhook endpoint with an absolute HTTP or HTTPS URL and a secret of at least 32 bytes.

<a name="design-every-consumer-for-repeated-delivery"></a>
## Design every consumer for repeated delivery

Delivery is at least once. A publish can succeed before acknowledgement fails, so the same `Stored.ID` may arrive again. Consumers deduplicate by that stable id and make effects idempotent; parked rows remain evidence instead of being discarded. Webhook dispatch snapshots one delivery row and one queue job atomically for each event and endpoint. Transport failures, 408, 425, 429, and 5xx retry; ordinary 4xx responses are permanent, and Retry-After is honored.

> [!WARNING]
> A direct publish after commit can be permanently lost when the process exits between the two operations. Use the outbox whenever that fact is required outside the process.

<a name="tenant-and-authority-come-from-the-grant"></a>
## Tenant and authority come from the Grant

> [!WARNING]
> Do not put a caller-selected tenant in an event and trust it later. Store persists tenant, action, and `authorized_by` from the Grant; webhook delivery inherits tenant and action from that Grant. Signed POST requests carry X-Arandu-Signature, X-Arandu-Timestamp, and X-Arandu-Delivery-ID; HMAC-SHA256 covers timestamp, delivery id, and body. Keep endpoint resolution behind the SSRF-guarded client, enforce the one-mebibyte payload ceiling, response and redirect bounds, timeout, and secret rotation.

Payloads contain identifiers and immutable public facts, not credentials, raw Grants, or full secret-bearing records. A valid webhook signature proves message integrity, not permission to read other tenant data.

<a name="prove-rollback-replay-and-replica-behavior"></a>
## Prove rollback, replay, and replica behavior

Return a sentinel from `data.Transaction` and prove neither business nor outbox row remains. Then make publish succeed while acknowledgement fails, drain twice, and prove consumer deduplication and multi-replica lock behavior. For webhooks, verify the signature byte-for-byte, previous-secret rotation, one job per event and endpoint, duplicate dispatch idempotency, SSRF refusal, response bounds, retry classification, fenced claims, audit status, and retention cleanup.

- Rollback leaves no outbox row.
- Replay keeps one external effect.
- Only the lock winner drains.
- Parked lag appears in diagnostics.
- A duplicate webhook delivery id never creates a second external effect.

<a name="continue-with-queue-workers-and-scheduling"></a>
## Continue with queue workers and scheduling

Queues execute durable webhook attempts, scheduling discovers work, observability exposes relay lag, and authorization defines the recorded Grant boundary. Durable outgoing webhooks are separate from Joaju's optional Pusher lifecycle callbacks.

- [Communication patterns](/docs/{{version}}/architecture/contracts-adapters)
- [Queues](/docs/{{version}}/async/queues)
- [Joaju `WebSocket`](/docs/{{version}}/async/joaju)
- [Observability](/docs/{{version}}/operations/observability)
- [Authorization](/docs/{{version}}/security/authorization)
