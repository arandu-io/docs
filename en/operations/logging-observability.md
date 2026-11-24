# Logging and observability

- [Write structured events through the request context](#write-structured-events-through-the-request-context)
- [Set the level and preserve middleware order](#set-the-level-and-preserve-middleware-order)
- [Emit one searchable domain event](#emit-one-searchable-domain-event)
- [Instrument outbound waits and background work](#instrument-outbound-waits-and-background-work)
- [Make missing telemetry visible](#make-missing-telemetry-visible)
- [Log identifiers, never credentials or authority](#log-identifiers-never-credentials-or-authority)
- [Assert fields at the handler boundary](#assert-fields-at-the-handler-boundary)
- [Diagnose from the last durable signal](#diagnose-from-the-last-durable-signal)
- [Open the request timeline only when needed](#open-the-request-timeline-only-when-needed)

<a name="write-structured-events-through-the-request-context"></a>
## Write structured events through the request context

Create builds one root slog logger from the logging settings: the standard output at `LOG_LEVEL`, in `LOG_FORMAT`, which is text in development and JSON elsewhere when empty. Each line carries the environment as its channel and no source position. Observe derives a request logger with `request_id`, method, and path, then records status, duration, bytes, query count, and SQL time on completion.

<a name="set-the-level-and-preserve-middleware-order"></a>
## Set the level and preserve middleware order

Set a released `LOG_LEVEL` and, when the environment's default is wrong for the collector, `LOG_FORMAT`. The composed pipeline already runs Observe after Recover, the body limits, and `TrustProxies`. Pass `context.Context` through controllers, services, and outbound work so `observability.Log` retrieves the scoped logger and Collector. `LOG_SAMPLING` is read into `Settings.Logging` but changes no logger, so it does not reduce log volume. Export and retention belong to the platform that collects stdout.

- Use the configured root logger.
- Preserve request context.
- Keep recovery outside observation.
- Send stdout to the platform collector.

<a name="emit-one-searchable-domain-event"></a>
## Emit one searchable domain event

Emit stable event text with structured identifiers through `observability.Log(ctx)`. The logger is always available; constructing a second package-global logger loses request fields and the configured level.

```go example
type InvoiceSender interface {
	Send(ctx context.Context, id string) error
}

type InvoiceService struct{ sender InvoiceSender }

func (s *InvoiceService) Send(ctx context.Context, id string) error {
	observability.Log(ctx).Info("invoice send requested", "invoice_id", id)
	return s.sender.Send(ctx, id)
}
```

<a name="instrument-outbound-waits-and-background-work"></a>
## Instrument outbound waits and background work

Use `observability.Client` for bounded outbound calls, `WithLogger` for controlled background ownership, and the application Recorder where queue or scheduler options accept it. The platform owns stdout collection, retention, indexing, rotation, and alerts.

- Correlate events with `request_id`.
- Give outbound clients a timeout.
- Keep background loggers explicitly scoped.
- Treat external retention as deployment configuration.

<a name="make-missing-telemetry-visible"></a>
## Make missing telemetry visible

Dropping context falls back to the default logger and loses correlation. A raw client omits outbound time and may wait indefinitely without a timeout. Disabled Collector methods are intentional no-ops, not evidence that instrumentation ran. A query on the log aggregator that reads `source` finds nothing, because lines carry `channel` instead.

<a name="log-identifiers-never-credentials-or-authority"></a>
## Log identifiers, never credentials or authority

> [!WARNING]
> Never log passwords, tokens, cookies, session payloads, `APP_KEY`, tracing secrets, request bodies, unbounded records, or full Grant internals. SQL arguments and debug dumps are sensitive cross-tenant diagnostic data, and production refuses debug level.

<a name="assert-fields-at-the-handler-boundary"></a>
## Assert fields at the handler boundary

Capture a `slog.Handler`, make requests with valid and malformed request ids, assert completion fields, exercise an instrumented outbound call, verify production JSON and development text, prove that `LOG_FORMAT` overrides the environment's default and that level filtering holds, and assert `LOG_SAMPLING` has no effect.

- Valid request ids are echoed and logged.
- Malformed ids are replaced.
- Outbound time reaches the Collector.
- `LOG_SAMPLING` is not advertised as effective.

<a name="diagnose-from-the-last-durable-signal"></a>
## Diagnose from the last durable signal

Correlate `request_id` across durable structured logs, reproduce under authorized tracing only when safe, and restore a supported `LOG_LEVEL` if configuration blocks boot. Repair an external collector in the platform; the application has no local log file to recover.

<a name="open-the-request-timeline-only-when-needed"></a>
## Open the request timeline only when needed

The diagnostic console expands one request, HTTP clients contribute outbound timing, deployment owns log collection, and configuration defines the effective controls.

- [Diagnostic observability](/docs/{{version}}/operations/observability)
- [HTTP client](/docs/{{version}}/services/http-client)
- [Deployment](/docs/{{version}}/operations/deployment)
- [Configuration](/docs/{{version}}/architecture/configuration)
