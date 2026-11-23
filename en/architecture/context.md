# Context

- [Context carries cancellation, not dependencies](#context-carries-cancellation-not-dependencies)
- [Pass the request lifetime into the service](#pass-the-request-lifetime-into-the-service)
- [Derive a bounded child for one dependency](#derive-a-bounded-child-for-one-dependency)
- [Cancellation is an outcome, not a log storm](#cancellation-is-an-outcome-not-a-log-storm)
- [Cancel the parent and observe the leaf](#cancel-the-parent-and-observe-the-leaf)
- [Standard and request context APIs](#standard-and-request-context-apis)
- [Cancellation crosses every explicit adapter](#cancellation-crosses-every-explicit-adapter)

<a name="context-carries-cancellation-not-dependencies"></a>
## Context carries cancellation, not dependencies

A request context says when work is no longer useful and carries a small set of request-scoped facts across package boundaries. Pass it first to every operation that can block. Constructors still carry durable dependencies; configuration and services do not belong in context values.

<a name="pass-the-request-lifetime-into-the-service"></a>
## Pass the request lifetime into the service

`http.Context.Ctx` returns the underlying request context. Pass it through the controller into database, cache, external HTTP, lock, and concurrent operations so disconnects and deadlines reach the leaf that is doing work. The authenticated Subject travels beside it into the service, because a route parameter is an untrusted identifier and never evidence of access; a request without a Subject is refused before the service runs.

```go example
type Invoice struct{ ID, Reference string }

func (i Invoice) ToArray() map[string]any {
	return map[string]any{"id": i.ID, "reference": i.Reference}
}

func (Invoice) With() map[string]any { return nil }

type InvoiceFinder interface {
	Find(ctx context.Context, subject security.Subject, id string) (Invoice, error)
}

type InvoiceController struct{ invoices InvoiceFinder }

func (c *InvoiceController) Show(ctx *fhttp.Context) error {
	subject, ok := ctx.User()
	if !ok {
		return ctx.Status(stdhttp.StatusUnauthorized)
	}
	invoice, err := c.invoices.Find(ctx.Ctx(), subject, ctx.Param("invoice"))
	if err != nil {
		return err
	}
	return ctx.JSON(stdhttp.StatusOK, invoice)
}
```

<a name="derive-a-bounded-child-for-one-dependency"></a>
## Derive a bounded child for one dependency

Derive a shorter child deadline when one dependency needs a tighter bound. Authentication middleware can attach a Subject with `auth.WithSubject`; readers use `auth.SubjectFrom` and handle the boolean instead of assuming identity exists.

- Use `context.WithTimeout` for a bounded dependency call.
- Call cancel with defer immediately after derivation.
- Read Subject with the reported presence boolean.
- Transfer detached work only to an explicit owner.

<a name="cancellation-is-an-outcome-not-a-log-storm"></a>
## Cancellation is an outcome, not a log storm

Propagate `context.Canceled` and `context.DeadlineExceeded` so the boundary decides whether a response remains useful. Do not replace a request context with `context.Background` to make a failing call continue, and never retain `*http.Context` after the handler returns.

> [!WARNING]
> `context.WithoutCancel` removes the deadline as well as cancellation. Detached work must establish a new deadline, report its failures, and make repeated execution safe.

<a name="cancel-the-parent-and-observe-the-leaf"></a>
## Cancel the parent and observe the leaf

Cancel the parent and assert a blocked repository or client call returns. Expire a child deadline while a sibling on the parent remains usable, and run the race detector against request-scoped values.

- Cancellation reaches the leaf call.
- A child timeout does not cancel its parent.
- A missing Subject is unauthenticated; only `auth.Guest(tenant)` enters policy as a guest.
- No request object is used after the handler returns.

<a name="standard-and-request-context-apis"></a>
## Standard and request context APIs

The standard library owns Context, `WithCancel`, `WithTimeout`, `WithDeadline`, Cause, and `WithoutCancel`. Framework `http.Context` exposes Ctx for the request lifetime; authentication packages supply narrow Subject context helpers.

<a name="cancellation-crosses-every-explicit-adapter"></a>
## Cancellation crosses every explicit adapter

Concurrency derives sibling contexts from this value. Controllers begin the chain, services and repositories preserve it, and outbound clients build requests with it. Background modules receive a process-owned context instead of borrowing a request.

- [Standard context](https://pkg.go.dev/context)
- [Controllers](/docs/{{version}}/http/controllers)
- [Concurrency](/docs/{{version}}/architecture/concurrency)
