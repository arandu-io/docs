# Broadcasting

- [A channel name is a typed publication boundary](#a-channel-name-is-a-typed-publication-boundary)
- [Choose transport and authorization before publishing](#choose-transport-and-authorization-before-publishing)
- [Publish one tenant-scoped event](#publish-one-tenant-scoped-event)
- [Keep channel construction and payloads explicit](#keep-channel-construction-and-payloads-explicit)
- [Do not confuse a typed channel with delivered fan-out](#do-not-confuse-a-typed-channel-with-delivered-fan-out)
- [Authorize subscription and publication separately](#authorize-subscription-and-publication-separately)
- [Test channel isolation and transport refusal](#test-channel-isolation-and-transport-refusal)
- [Use outbox and queues when delivery must survive](#use-outbox-and-queues-when-delivery-must-survive)

> [!WARNING]
> **Current limitations.** Realtime channel authorization is released. General notification fan-out and durable delivery are separate incomplete paths; this guide covers authorized event publication to the configured broadcaster only.

<a name="a-channel-name-is-a-typed-publication-boundary"></a>
## A channel name is a typed publication boundary

The released broadcasting package provides typed channels, events, and driver contracts. It does not prove a fully wired application manager or browser subscription journey, so this page documents the useful component boundary and marks end-to-end integration as limited.

<a name="choose-transport-and-authorization-before-publishing"></a>
## Choose transport and authorization before publishing

Select a concrete driver at bootstrap, define how a trusted Subject may subscribe, and decide whether a dropped connection is acceptable. Use a durable outbox or queue when the event must survive process loss rather than treating broadcast as storage.

- Name the resource in channel construction.
- Derive tenant namespace only from Grant.
- Bound payload size and public fields.
- Choose durable work separately.

<a name="publish-one-tenant-scoped-event"></a>
## Publish one tenant-scoped event

Construct a released channel from the record identifier, build an explicit public payload, and call `Broadcaster.Broadcast` with the authorized Grant. The driver derives the wire tenant namespace through `TenantChannel`; application input never prefixes it.

```go example
func PublishInvoiceUpdate(ctx context.Context, broadcaster broadcasting.Broadcaster, grant security.Grant, invoiceID, reference string) error {
	channel := broadcasting.NewPrivateChannel("invoices." + invoiceID)
	payload := map[string]any{
		"invoice": map[string]any{
			"id":        invoiceID,
			"reference": reference,
		},
	}
	if err := broadcaster.Broadcast(
		ctx, grant, []broadcasting.Channel{channel}, "invoice.updated", payload,
	); err != nil {
		return err
	}
	return nil
}
```

<a name="keep-channel-construction-and-payloads-explicit"></a>
## Keep channel construction and payloads explicit

Keep event names stable, serialize a deliberate resource rather than a model, attach no bearer or session credential, and reconnect clients through their normal authenticated subscription path instead of trusting a channel string supplied by the browser.

- Use explicit public and private channel types.
- Keep event data small and versionable.
- Publish immutable facts.
- Use durable replay only through a stored event system.

<a name="do-not-confuse-a-typed-channel-with-delivered-fan-out"></a>
## Do not confuse a typed channel with delivered fan-out

A driver outage, disconnected client, serialization failure, and denied subscription are different outcomes. The released typed API alone does not guarantee reconnect, replay, ordering, or a manager wired into the starter application.

> [!WARNING]
> The current released contract does not include a complete application broadcasting journey. Before relying on one, prove bootstrap, subscription policy, selected transport, and browser or service delivery together in the tagged application you deploy.

<a name="authorize-subscription-and-publication-separately"></a>
## Authorize subscription and publication separately

> [!WARNING]
> Never accept a tenant prefix in a client channel, place secrets or full models in payloads, or treat knowledge of a private channel name as authorization. `RequestedChannel` rejects a colon, and `TenantChannel` derives the namespace only from the authorized Grant. Authenticate and authorize every subscription.

Broadcast delivery is representation, not authority. A client receiving an id must still pass the normal policy for any subsequent mutation or read.

<a name="test-channel-isolation-and-transport-refusal"></a>
## Test channel isolation and transport refusal

Use a fake driver to assert exact channel and payload, deny a foreign tenant subscription, inject serialization and driver failures, test cancellation, and run the selected transport journey before enabling it in a deployment.

- Channel names contain trusted tenant scope.
- Foreign subscribers are refused.
- Payloads expose only public fields.
- A driver error remains observable.

<a name="use-outbox-and-queues-when-delivery-must-survive"></a>
## Use outbox and queues when delivery must survive

Events and outbox provide durable facts, queues provide durable work, authorization protects subscription, and resources define public payloads.

- [Events and outbox](/docs/{{version}}/async/events-outbox)
- [Queues](/docs/{{version}}/async/queues)
- [Authorization](/docs/{{version}}/security/authorization)
- [HTTP resources](/docs/{{version}}/http/resources)
