# Task scheduling

- [Let the scheduler decide when, not how to retry](#let-the-scheduler-decide-when-not-how-to-retry)
- [Declare scope before writing the task body](#declare-scope-before-writing-the-task-body)
- [Schedule one job for every tenant](#schedule-one-job-for-every-tenant)
- [Inspect and invoke the registered schedule](#inspect-and-invoke-the-registered-schedule)
- [Treat a missing lock as a deployment defect](#treat-a-missing-lock-as-a-deployment-defect)
- [Global tasks cannot touch tenant data](#global-tasks-cannot-touch-tenant-data)
- [Run the schedule, queue, and worker end to end](#run-the-schedule-queue-and-worker-end-to-end)
- [Repair the dependency and run the task once](#repair-the-dependency-and-run-the-task-once)
- [Keep durable work in the queue](#keep-durable-work-in-the-queue)

<a name="let-the-scheduler-decide-when-not-how-to-retry"></a>
## Let the scheduler decide when, not how to retry

A scheduled task is a small decision executed on a five-field cron expression. It normally discovers work and enqueues a durable job; the scheduler itself does not retry a failed Run.

<a name="declare-scope-before-writing-the-task-body"></a>
## Declare scope before writing the task body

A `foundation.Task` has a unique ID, valid Spec, Action, and Run. Global receives a zero Grant and cannot access repositories; `PerTenant` receives `SystemGrant(Action, tenant)` for each tenant the function given to `WithTenants` answers. Set Timeout because zero becomes five minutes.

- Name a unique task id.
- Choose Global or `PerTenant` deliberately.
- Answer the tenants with `WithTenants` for `PerTenant`.
- Singleton claims are shared across replicas only when `drivers.Locks` is.

<a name="schedule-one-job-for-every-tenant"></a>
## Schedule one job for every tenant

List the application's tasks in `routes/console.go` and hand them to `WithSchedule`; a module declares its own through Schedule. Create registers the scheduler last, over every module's tasks and the application's, and it runs in the serving process. Keep durable work in a queue handler: the scheduler's Grant becomes the job provenance.

```go example
type CloseCycle struct{ Cycle string }

func Tasks(d *foundation.Drivers) []foundation.Task {
	return []foundation.Task{{
		ID: "billing.close", Spec: "0 3 * * *", Scope: foundation.PerTenant,
		Action: "invoice.send", Timeout: 30 * time.Second, Singleton: true,
		Run: func(ctx context.Context, g auth.Grant) error {
			return queue.Dispatch(ctx, d.Queue, g, "invoice.send", CloseCycle{Cycle: "current"})
		},
	}}
}
```

<a name="inspect-and-invoke-the-registered-schedule"></a>
## Inspect and invoke the registered schedule

`aru schedule:list` reports id, cron, scope, singleton, next, last, and last error. `schedule:run` invokes the same Grant, lock, and instrumentation path; a per-tenant manual run requires `--tenant`, and a manual run of a singleton takes the same claim, so a second one in the same UTC minute does nothing.

- Inspect registration before waiting for cron.
- Run one task manually for a selected tenant.
- Queue slow or retryable work.
- Keep task body safe to invoke twice.

<a name="treat-a-missing-lock-as-a-deployment-defect"></a>
## Treat a missing lock as a deployment defect

Serving processes check near each minute. Without a shared store `drivers.Locks` is nil, and every replica can run a singleton. A singleton claims its task, tenant, and UTC minute (`sched:<tenant>:<id>:<minute>`) and never releases the claim: it expires after an hour, or after the timeout if that is longer, so a run that fails is not repeated by another replica in the same minute. A per-tenant task with no `WithTenants` does not run, and the scheduler reports that rather than skipping in silence. A failed run is not retried, so a task that performs an external side effect directly can lose that work; enqueue a durable idempotent job instead.

<a name="global-tasks-cannot-touch-tenant-data"></a>
## Global tasks cannot touch tenant data

> [!WARNING]
> Never manufacture a tenant inside Global or take one from HTTP, payload, or an environment shortcut. Repository access receives only the per-tenant Grant created from the trusted resolver and declared Action.

<a name="run-the-schedule-queue-and-worker-end-to-end"></a>
## Run the schedule, queue, and worker end to end

Construct the scheduler, call `RunNow`, assert one persisted job, start a worker, and assert the handler receives tenant and payload. Duplicate IDs, bad cron, and a task without Run fail before serving. A missing tenant resolver appears at execution and diagnostics; a missing shared lock must be caught as a deployment configuration test because singleton work otherwise runs on every replica.

- `RunNow` enqueues once.
- The worker receives the expected Grant.
- Invalid schedules fail boot.
- Only the singleton lock winner runs.

<a name="repair-the-dependency-and-run-the-task-once"></a>
## Repair the dependency and run the task once

Use `schedule:list` to separate absent registration from last-run failure. Restore database, queue, resolver, or lock, then invoke one affected tenant. The scheduler does not retry, but recovery can invoke work again; design both the task and any queued handler to be safe to repeat.

<a name="keep-durable-work-in-the-queue"></a>
## Keep durable work in the queue

Queues own durable execution, outbox owns facts after commit, queue operations own recovery, and observability reports schedule failure and lag.

- [Queues](/docs/{{version}}/async/queues)
- [Events and outbox](/docs/{{version}}/async/events-outbox)
- [Queue operations](/docs/{{version}}/operations/queues)
- [Observability](/docs/{{version}}/operations/observability)
