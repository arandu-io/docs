# Queue operations

- [Operate the queue through the application binary](#operate-the-queue-through-the-application-binary)
- [Migrate queue tables and register every handler](#migrate-queue-tables-and-register-every-handler)
- [Start one draining worker and inspect backlog](#start-one-draining-worker-and-inspect-backlog)
- [Pause intake, restart workers, and manage dead letters](#pause-intake-restart-workers-and-manage-dead-letters)
- [Treat leases and failed jobs as recoverable state](#treat-leases-and-failed-jobs-as-recoverable-state)
- [Keep tenant scope in payload, handler, and command](#keep-tenant-scope-in-payload-handler-and-command)
- [Run one success, duplicate, retry, and shutdown journey](#run-one-success-duplicate-retry-and-shutdown-journey)
- [Recover state before deleting it](#recover-state-before-deleting-it)
- [Separate queue design from queue operations](#separate-queue-design-from-queue-operations)

<a name="operate-the-queue-through-the-application-binary"></a>
## Operate the queue through the application binary

Aru forwards queue commands to the project binary because only that binary knows handlers, database, failed-job provider, batch repository, tenant rules, and selected connection. The starter wires a worker and operational commands; it does not ship a dashboard or supervisor.

<a name="migrate-queue-tables-and-register-every-handler"></a>
## Migrate queue tables and register every handler

Migrate jobs, `failed_jobs`, and `job_batches`; register handlers in the provider's `RegisterJobs`, which `bootstrap/app.go` hands to `WithJobs`; and set the default queue and a retry lease longer than the slowest handler. `QUEUE_CONNECTION=database` keeps jobs in the application's database, and redis selects the RESP queue once `hesape/queue/connectors/redis` is linked by its blank import; naming redis before that stops the boot with the two lines to add. Cross-process pause and restart signals still require a shared RESP cache.

- `QUEUE_WORKERS` and `QUEUE_MAX_ATTEMPTS` are positive when set.
- `REDIS_URL` is present when shared signal cache is selected.
- Payload and handler preserve tenant and Grant provenance.
- Worker and web use the same image.

<a name="start-one-draining-worker-and-inspect-backlog"></a>
## Start one draining worker and inspect backlog

Start a worker with explicit queue and concurrency, then monitor that connection and queue. Environment defaults fill omitted values; flags override only their documented invocation fields. SIGTERM stops intake and drains in-flight work.

```sh
aru queue:work --queue=default --workers=4
aru queue:monitor database:default --max=1000
```

<a name="pause-intake-restart-workers-and-manage-dead-letters"></a>
## Pause intake, restart workers, and manage dead letters

Pause and resume a target, request graceful worker restart, inspect tenant-scoped failures, retry or forget selected ids, and prune only reviewed failed or batch rows. `queue:listen` is a development reload workflow, not the ordinary production daemon.

- `aru queue:pause database:default --for=30m` and `queue:resume`.
- `aru queue:restart` after a deployment.
- `aru queue:failed --tenant=<id>` before retry or forget.
- `queue:retry`, `queue:forget`, and `queue:retry-batch` read their flags after the ids too.
- Use tenant-scoped retry-batch and prune commands.

<a name="treat-leases-and-failed-jobs-as-recoverable-state"></a>
## Treat leases and failed jobs as recoverable state

Handler errors retry until `MaxAttempts` and then park in `failed_jobs`. A dead worker releases only after its lease; a lease shorter than runtime can execute twice. Missing handlers, tables, selected stores, or shared signal cache must be repaired before recovery commands are trusted.

<a name="keep-tenant-scope-in-payload-handler-and-command"></a>
## Keep tenant scope in payload, handler, and command

> [!WARNING]
> Failed-job and batch commands require tenant scope. Never deserialize a caller-supplied Grant or role; re-authorize in handlers, minimize payloads, redact failure text, and make every external side effect idempotent because delivery is at least once.

<a name="run-one-success-duplicate-retry-and-shutdown-journey"></a>
## Run one success, duplicate, retry, and shutdown journey

Exercise success, duplicate, retry, lease expiry, crash recovery, cross-process pause and resume, SIGTERM drain, tenant isolation, the released database queue, and propagation of `QUEUE_DEFAULT`, `QUEUE_WORKERS`, `QUEUE_RETRY_AFTER`, and `QUEUE_MAX_ATTEMPTS` into the worker. Assert `QUEUE_CONNECTION=redis` without the linked connector stops the boot rather than falling back to the database.

- A real job drains.
- A wrong tenant cannot inspect or retry it.
- Failure reaches dead-letter storage and can be retried.
- Shutdown drains without taking new work.

<a name="recover-state-before-deleting-it"></a>
## Recover state before deleting it

Pause the affected queue, preserve rows and logs, fix or roll back the handler, retry a canary, then resume. Clear, forget, flush, and prune delete waiting work or evidence; back up and use the narrowest tenant and id filter.

<a name="separate-queue-design-from-queue-operations"></a>
## Separate queue design from queue operations

The async guide defines payload and handler design, logging exposes worker evidence, deployment ships the same artifact, and console tests exercise the delegated command boundary.

- [Queue design](/docs/{{version}}/async/queues)
- [Logging](/docs/{{version}}/operations/logging-observability)
- [Deployment](/docs/{{version}}/operations/deployment)
- [Console testing](/docs/{{version}}/testing/console)
