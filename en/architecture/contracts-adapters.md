# Contracts and adapters

- [The consumer owns the smallest useful interface](#the-consumer-owns-the-smallest-useful-interface)
- [Put authorization in the repository seam](#put-authorization-in-the-repository-seam)
- [Select adapters at boot](#select-adapters-at-boot)
- [Avoid interfaces that hide authority or ownership](#avoid-interfaces-that-hide-authority-or-ownership)
- [Run one contract against each adapter](#run-one-contract-against-each-adapter)
- [Framework contracts stay narrow](#framework-contracts-stay-narrow)
- [Composition makes substitution visible](#composition-makes-substitution-visible)

<a name="the-consumer-owns-the-smallest-useful-interface"></a>
## The consumer owns the smallest useful interface

A contract names the behavior one consumer needs. An adapter satisfies it structurally and is selected in the composition root. Domain code stays independent of SQL, RESP, SMTP, or an external service without adding a registry or universal abstraction package.

<a name="put-authorization-in-the-repository-seam"></a>
## Put authorization in the repository seam

Put authorization in the repository seam so changing storage cannot remove tenant scope. The service declares only the operation it calls, and the compile-time assertion keeps the concrete SQL adapter aligned with that consumer-owned shape.

```go example
type Invoice struct{ ID string }

type InvoiceRepository interface {
	Find(context.Context, security.Grant, string) (Invoice, error)
}

type InvoiceService struct{ invoices InvoiceRepository }

func NewInvoiceService(invoices InvoiceRepository) *InvoiceService {
	return &InvoiceService{invoices: invoices}
}

type SQLInvoiceRepository struct{}

func (*SQLInvoiceRepository) Find(context.Context, security.Grant, string) (Invoice, error) {
	return Invoice{}, nil
}

var _ InvoiceRepository = (*SQLInvoiceRepository)(nil)
```

<a name="select-adapters-at-boot"></a>
## Select adapters at boot

`foundation.NewDrivers` builds the database connector, cache store, mail transport, and filesystem disk the configuration names, and a provider passes each one, or a test fake, once. Adapters translate driver absence, conflict, and unavailable errors when callers need to branch while preserving the cause for operators.

- Validate required collaborators in constructors.
- Keep optional behavior in a separate small interface.
- Wrap driver errors without leaking credentials or SQL.
- Place compile-time assertions beside adapters.

<a name="avoid-interfaces-that-hide-authority-or-ownership"></a>
## Avoid interfaces that hide authority or ownership

A repository accepting only an id for tenant-owned data has nowhere to prove authority. Large interfaces force fakes to implement unused methods and couple unrelated consumers; duplicated small interfaces are safer when semantics genuinely differ.

> [!WARNING]
> Do not recover missing authority with a global tenant, request-local singleton, or adapter-side environment lookup. Pass `security.Grant` at the operation that reaches protected storage.

<a name="run-one-contract-against-each-adapter"></a>
## Run one contract against each adapter

Run one shared contract against every adapter. Exercise success, absence, tenant isolation, cancellation, and error mapping through the consumer interface rather than reaching into a concrete driver's implementation.

- Zero and foreign Grants return no protected data.
- Memory and network adapters meet the same behavior.
- Cancellation returns promptly.
- Compile-time assertions cover every concrete adapter.

<a name="framework-contracts-stay-narrow"></a>
## Framework contracts stay narrow

Choose the communication contract from the interaction, not from a popularity ranking. REST is request-response for public resources; `gRPC` is a typed service boundary that also supports streaming; a queue persists asynchronous work; `WebSocket` is a bidirectional live session; server-sent events (SSE) stream server updates over one HTTP response; and a webhook delivers a signed event to an external HTTP endpoint.

| Surface | Released contract |
| --- | --- |
| REST | request-response resources; use for public or browser-facing APIs |
| `gRPC` | typed unary or streaming service calls; use Connect, `gRPC`, or gRPC-Web on the application HTTP server |
| Queue | durable asynchronous work with at-least-once delivery |
| `WebSocket` | persistent bidirectional sessions for low-latency interaction |
| SSE | one-way server updates over a cancellable HTTP response |
| Webhook | durable signed event delivery to an external service |

<a name="composition-makes-substitution-visible"></a>
## Composition makes substitution visible

REST, `gRPC`, queues, `WebSocket`, SSE, and webhooks solve different ownership and delivery problems. Configuration chooses the concrete adapter, the composition root constructs it, and the service keeps only its contract. Context carries request lifetime and Grant carries authorized tenant scope, making substitution visible in one wiring diff.

Do not substitute a live transport for durable storage: `WebSocket` and SSE do not replay missed work by themselves. A queue executes internal work, while a webhook crosses an external HTTP boundary. `gRPC` streaming shares the serving process and shutdown lifecycle rather than creating a second listener.

- [HTTP routing and `gRPC`](/docs/{{version}}/http/routing)
- [HTTP responses and SSE](/docs/{{version}}/http/responses)
- [Queues](/docs/{{version}}/async/queues)
- [Joaju `WebSocket`](/docs/{{version}}/async/joaju)
- [Events and webhooks](/docs/{{version}}/async/events-outbox)
