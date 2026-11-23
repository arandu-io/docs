# Concurrency

- [Concurrent work still has one owner](#concurrent-work-still-has-one-owner)
- [Run independent reads with a ceiling](#run-independent-reads-with-a-ceiling)
- [Defer work deliberately and recover owned boundaries](#defer-work-deliberately-and-recover-owned-boundaries)
- [A failed set returns no partial slice](#a-failed-set-returns-no-partial-slice)
- [Prove the ceiling and cancellation](#prove-the-ceiling-and-cancellation)
- [Finite task operations](#finite-task-operations)
- [Background loops belong to modules](#background-loops-belong-to-modules)

<a name="concurrent-work-still-has-one-owner"></a>
## Concurrent work still has one owner

The concurrency helpers run a finite set of typed tasks under one context. They preserve argument order, cancel siblings after the first failure, wait for every started goroutine, and convert panics into errors rather than letting one task terminate the process.

<a name="run-independent-reads-with-a-ceiling"></a>
## Run independent reads with a ceiling

Use `RunWith` when independent reads can proceed together but share a finite connection or upstream quota. Results remain in task argument order, not completion order, and a positive ceiling prevents the fan-out from consuming every available resource.

```go example
type Invoice struct{ ID string }

type InvoiceFinder interface {
	Find(ctx context.Context, id string) (Invoice, error)
}

func LoadPair(ctx context.Context, invoices InvoiceFinder, firstID, secondID string) ([]Invoice, error) {
	results, err := concurrency.RunWith(ctx, 4,
		func(ctx context.Context) (Invoice, error) { return invoices.Find(ctx, firstID) },
		func(ctx context.Context) (Invoice, error) { return invoices.Find(ctx, secondID) },
	)
	if err != nil {
		return nil, err
	}
	return results, nil
}
```

<a name="defer-work-deliberately-and-recover-owned-boundaries"></a>
## Defer work deliberately and recover owned boundaries

Run has no ceiling and `RunWith` adds one; zero or a negative value means no ceiling. Recover guards one named call site and returns `PanicError` with the value and stack. Defer returns a callback and runs nothing until invoked.

- Use a ceiling below the database pool or upstream quota.
- Give deferred work its own deadline.
- Call the returned deferred callback exactly once.
- Keep long-lived loops in module lifecycle contracts.

<a name="a-failed-set-returns-no-partial-slice"></a>
## A failed set returns no partial slice

The first task failure cancels siblings and returns a nil results slice. A task that ignores its derived context can keep the entire call waiting after its result was discarded. An already-cancelled parent starts nothing; cancellation during execution returns the parent error.

> [!WARNING]
> Defer detaches request cancellation and deadlines with `context.WithoutCancel`. The new owner must set a deadline, report failures, and make duplicate execution safe; invoking the callback twice runs the tasks twice.

<a name="prove-the-ceiling-and-cancellation"></a>
## Prove the ceiling and cancellation

Count simultaneous task entries to prove the ceiling, fail one task to prove a sibling observes cancellation, reverse completion order to prove result order, and panic inside a task to match `*concurrency.PanicError`.

- Maximum simultaneous entries equals the configured ceiling.
- A blocked sibling exits after the first error.
- Results retain argument order.
- Deferred work does not run before its callback.

<a name="finite-task-operations"></a>
## Finite task operations

`Task[T]`, Run, `RunWith`, Defer, Recover, and `PanicError` form the direct API. `Manager[T]` resolves configured drivers, but the built-in process, fork, and sync names currently reach the same goroutine behavior; prefer direct functions unless configuration is itself required.

<a name="background-loops-belong-to-modules"></a>
## Background loops belong to modules

A finite fan-out inside a request differs from a scheduler, queue worker, or relay loop. Long-lived work implements the application lifecycle, propagates a process-owned shutdown context, and exposes health. Cross-replica exclusivity requires a shared lock, not goroutines.

- [Concurrency API](https://pkg.go.dev/github.com/arandu-io/hesape/concurrency)
- [Context](/docs/{{version}}/architecture/context)
- [Modules](/docs/{{version}}/architecture/modules-and-composition)
