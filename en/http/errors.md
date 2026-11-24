# Error handling

- [One handler classifies, reports, and renders](#one-handler-classifies-reports-and-renders)
- [Build the handler from loaded configuration](#build-the-handler-from-loaded-configuration)
- [Install recovery as the outermost middleware](#install-recovery-as-the-outermost-middleware)
- [Return an intentional HTTP failure](#return-an-intentional-http-failure)
- [Unknown errors stay server errors](#unknown-errors-stay-server-errors)
- [Public detail and private cause are different fields](#public-detail-and-private-cause-are-different-fields)
- [Prove both representation and report](#prove-both-representation-and-report)
- [Recover the service, not only the request](#recover-the-service-not-only-the-request)
- [Follow failures across their boundaries](#follow-failures-across-their-boundaries)

<a name="one-handler-classifies-reports-and-renders"></a>
## One handler classifies, reports, and renders

Errors cross one centralized recovery boundary. Expected domain and transport failures keep their meaning, while escaped panics and unknown internal errors are recorded once and become a bounded production response without stack frames or secrets. Mapping a domain error to a status is the application's job: return an error that states its own status, and never infer one from the error text.

<a name="build-the-handler-from-loaded-configuration"></a>
## Build the handler from loaded configuration

Create builds the one exception handler from the loaded configuration, reads the binary's module path to tell the application's source frames from the framework's, and installs Recover outermost after the refusal sentences. The debug page exists only in development; production validation keeps it off. The application registers its own answers on that handler in `WithExceptions`.

<a name="install-recovery-as-the-outermost-middleware"></a>
## Install recovery as the outermost middleware

Return errors from controller actions and let recovery select the configured response. In development the error page can show source frames; production sends a generic status and correlation context while structured logs retain the wrapped cause.

```go example
var ErrQuotaExceeded = errors.New("quota exceeded")

func WithErrorHandling(drivers *foundation.Drivers) *foundation.ApplicationBuilder {
	return foundation.Configure(drivers).
		WithExceptions(func(e *exception.Handler) {
			e.DontReport(ErrQuotaExceeded)
			e.ShouldRenderJSONWhen(func(r *stdhttp.Request, _ error) bool {
				return strings.HasPrefix(r.URL.Path, "/api/")
			})
		})
}

func Register(router *fhttp.Router, show func(*fhttp.Context) error) {
	router.Action("GET", "/invoices/{invoice}", show)
}
```

<a name="return-an-intentional-http-failure"></a>
## Return an intentional HTTP failure

Define application error types for absence, denial, conflict, invalid state, and dependency unavailability. The router adapter maps a returned error, bare or wrapped with %w: a not-found sentinel answers 404, `security.ErrForbidden` 403, `security.ErrCSRF` 419, `database.ErrUniqueViolation` 409, and an error with an `HTTPStatus()` int method answers that status. Give a domain error that method rather than switching on it in each controller, wrap causes with context, and avoid inferring statuses by matching arbitrary error strings.

- Return validation errors through their dedicated adapter.
- Use 404 without revealing cross-tenant existence.
- Use 409 for a real state conflict.
- Use 503 for a required dependency outage.

<a name="unknown-errors-stay-server-errors"></a>
## Unknown errors stay server errors

Writing part of a response before returning an error can make the final status impossible to change. A panic after headers were sent is recorded but cannot produce a clean replacement; handlers should finish fallible work before committing output.

> [!WARNING]
> Development source pages and request dumps can contain secrets. Keep the development console disabled or protected outside development, and redact credentials, cookies, Grants, and database URLs from structured values.

<a name="public-detail-and-private-cause-are-different-fields"></a>
## Public detail and private cause are different fields

> [!WARNING]
> Keep public detail separate from the private wrapped cause. Never expose SQL, connection URLs, keys, tokens, cookies, request bodies, stack traces, or policy internals; keep debug false anywhere an untrusted person can reach.

<a name="prove-both-representation-and-report"></a>
## Prove both representation and report

Exercise each expected error through the registered route, then panic before and after a write. Assert production bodies reveal no cause, the recorder observes one event, and a domain refusal does not become an internal fault.

- Expected failures retain their status.
- Unknown errors are recorded once.
- Production responses contain no stack or credential.
- A cross-tenant miss is indistinguishable from absence.

<a name="recover-the-service-not-only-the-request"></a>
## Recover the service, not only the request

Use the response request id to find the correlated log and trace, inspect kernel diagnostics and dependent module health, then fix or roll back the release. Replay a mutation only after its commit outcome is known and the operation is idempotent.

<a name="follow-failures-across-their-boundaries"></a>
## Follow failures across their boundaries

Observability records the event, validation owns field errors, controllers return failures, and deployment decides whether development diagnostics can be enabled.

- [Logging and observability](/docs/{{version}}/operations/logging-observability)
- [Controllers](/docs/{{version}}/http/controllers)
- [Validation](/docs/{{version}}/http/validation)
