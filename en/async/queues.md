# Queues

- [Persist work that must survive the request](#persist-work-that-must-survive-the-request)
- [Choose a store whose transaction boundary you understand](#choose-a-store-whose-transaction-boundary-you-understand)
- [Create from a Grant, push, and handle by name](#create-from-a-grant-push-and-handle-by-name)
- [Set lease and retry policy from handler behavior](#set-lease-and-retry-policy-from-handler-behavior)
- [Know whether a job was released or parked](#know-whether-a-job-was-released-or-parked)
- [Never rebuild authority from the payload](#never-rebuild-authority-from-the-payload)
- [Prove rollback and duplicate delivery](#prove-rollback-and-duplicate-delivery)
- [Inspect, repair, retry, and observe lag](#inspect-repair-retry-and-observe-lag)
- [Separate producer design from worker operations](#separate-producer-design-from-worker-operations)

<a name="persist-work-that-must-survive-the-request"></a>
## Persist work that must survive the request

A queue separates accepting work from performing it. Hesape queue jobs carry tenant, action, authorizer, payload, and stable UUID, and a worker rebuilds the Grant before invoking the named handler. Delivery is at least once, so handler idempotence is an application contract.

<a name="choose-a-store-whose-transaction-boundary-you-understand"></a>
## Choose a store whose transaction boundary you understand

Create registers the queue's module itself: the jobs table, lag health, and diagnostics of the queue `QUEUE_CONNECTION` names, which `drivers.Queue` holds. Registering `jobs.NewModule` again in `WithProviders` registers it twice. The database queue can commit its row in the same SQL transaction as the producing write. RESP, selected by `QUEUE_CONNECTION=redis` once `hesape/queue/connectors/redis` is linked, is useful when chosen deliberately but cannot join that transaction.

- Select database or RESP deliberately.
- Leave the queue module to Create.
- Keep web and worker on one artifact.
- Choose idempotency before enqueue.

<a name="create-from-a-grant-push-and-handle-by-name"></a>
## Create from a Grant, push, and handle by name

Register the exact name in the provider's `RegisterJobs`, which `bootstrap/app.go` hands to `WithJobs`, and dispatch a typed payload under the authorized Grant with `queue.Dispatch`. `HandlerFor` decodes the payload into the same type before the handler runs. Dispatch builds the job with `jobs.New` on the default queue, so a missing tenant or name is refused, and Push checks authority even for a job built by hand for another queue.

```go example
type SendInvoice struct{ InvoiceID string }

type InvoiceSender interface {
	Send(ctx context.Context, g auth.Grant, jobID, invoiceID string) error
}

type AppServiceProvider struct{ Invoices InvoiceSender }

func (p *AppServiceProvider) RegisterJobs(w *queue.Worker) error {
	w.Handle("invoice.send", queue.HandlerFor(func(ctx context.Context, g auth.Grant, j *jobs.Job, in SendInvoice) error {
		return p.Invoices.Send(ctx, g, j.UUID, in.InvoiceID)
	}))
	return nil
}

func Enqueue(ctx context.Context, drivers *foundation.Drivers, grant auth.Grant, invoiceID string) error {
	if err := queue.Dispatch(ctx, drivers.Queue, grant, "invoice.send", SendInvoice{InvoiceID: invoiceID}); err != nil {
		return err
	}
	return nil
}
```

<a name="set-lease-and-retry-policy-from-handler-behavior"></a>
## Set lease and retry policy from handler behavior

`queue:work` builds its worker from `Settings.Queue`: `QUEUE_DEFAULT` names the queue, `QUEUE_WORKERS` the concurrency, `QUEUE_RETRY_AFTER` the lease in seconds, and `QUEUE_MAX_ATTEMPTS` the tries, and `--queue` and `--workers` override the first two for one invocation. An unset variable keeps the queue component's own default: four workers, a five-minute lease, five tries, and a wait between tries that doubles from two seconds and stops growing at one hour. A handler's context is cancelled after one lease, so a handler that returns when its context ends hands the job back before another worker can take it.

- Bound concurrency by dependencies.
- Set lease above maximum handler duration.
- Set attempt and backoff policy explicitly.
- Carry the application Recorder into workers.

<a name="know-whether-a-job-was-released-or-parked"></a>
## Know whether a job was released or parked

A returned error consumes an attempt and releases after backoff; exhausted work parks. `ErrManuallyFailed` parks a permanent failure immediately, and a panic is recovered, recorded with stack, and parked on first delivery. The worker reports a job of a batch or a chain after it runs, so a handler that still calls `bus.Handled` counts it twice and pushes the chain's next link twice. Failed payloads remain tenant-protected customer data.

> [!WARNING]
> A handler that ignores its context keeps running after the lease ends, and another worker can take the same job. Increasing attempts does not make a non-idempotent external side effect safe.

<a name="never-rebuild-authority-from-the-payload"></a>
## Never rebuild authority from the payload

> [!WARNING]
> The worker uses recorded tenant and action for its Grant, and a job pushed while an operator acted in the account carries that operator, so what it writes is recorded as theirs. A push under a read-only Grant is refused with `auth.ErrReadOnly`. Never accept `tenant_id`, action, role, or a serialized Grant from JSON. Pass the worker Grant to services and repositories, and keep failure inspection tenant-scoped.

<a name="prove-rollback-and-duplicate-delivery"></a>
## Prove rollback and duplicate delivery

For database queue, push inside a transaction that returns a sentinel and prove `PendingSize` stays zero. Deliver the same stable job twice and prove one external effect; also test malformed payload, permanent failure, panic, and a deliberately short lease.

- Rollback removes the enqueued row.
- Duplicate delivery is idempotent.
- Permanent failure parks.
- A foreign tenant cannot inspect failed work.

<a name="inspect-repair-retry-and-observe-lag"></a>
## Inspect, repair, retry, and observe lag

Inspect failures and lag, repair the handler or dependency, retry a canary, and replace workers gracefully. Forget, flush, and clear destroy evidence or work and require the narrowest reviewed tenant scope.

<a name="separate-producer-design-from-worker-operations"></a>
## Separate producer design from worker operations

Operations documents commands, scheduling creates recurring jobs, outbox persists facts, and authorization supplies the Grant boundary. A queue is an internal durable execution boundary; a webhook uses queued work to cross into an external HTTP service.

- [Communication patterns](/docs/{{version}}/architecture/contracts-adapters)
- [Queue operations](/docs/{{version}}/operations/queues)
- [Scheduling](/docs/{{version}}/async/scheduling)
- [Events and webhooks](/docs/{{version}}/async/events-outbox)
- [Authorization](/docs/{{version}}/security/authorization)
