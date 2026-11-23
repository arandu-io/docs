# Request lifecycle

- [One request crosses one explicit pipeline](#one-request-crosses-one-explicit-pipeline)
- [Know the route and the middleware that own it](#know-the-route-and-the-middleware-that-own-it)
- [Observe a complete response](#observe-a-complete-response)
- [Test through Handler when the port is irrelevant](#test-through-handler-when-the-port-is-irrelevant)
- [Each stage has one question](#each-stage-has-one-question)

<a name="one-request-crosses-one-explicit-pipeline"></a>
## One request crosses one explicit pipeline

The Go server accepts the connection, applies global middleware, matches one route, applies route middleware, constructs an HTTP context, calls the registered handler, and writes the response. The same composed application used by the server can be booted by tests without opening a port.

<a name="know-the-route-and-the-middleware-that-own-it"></a>
## Know the route and the middleware that own it

Use the route inventory rather than guessing from a URL. The framework composes the global pipeline from the configuration, and `bootstrap/app.go` adjusts it in `WithMiddleware`. Route middleware is attached where the route is declared. A controller receives collaborators through its constructor; it does not resolve them during the request.

```sh
aru route:list
aru dev
```

<a name="observe-a-complete-response"></a>
## Observe a complete response

```sh
curl -i http://127.0.0.1:8080/
```

The response has crossed the global pipeline in this order: the refusal sentences, Recover, what Prepend added, the two body limits, `TrustProxies`, Observe, `SecurityHeaders`, what `BeforeThrottle` added, the global Throttle, the CSRF check, `OverrideMethod`, and what Append added. Then come route matching, route middleware, the controller, and the renderer. For a state-changing route, authentication, authorization, CSRF, and validation remain separate checks rather than one generic guard.

> [!NOTE]
> A request under `/_arandu/` skips the application's pipeline only when a framework route answers it: a probe, the debug console, the reload, or an asset. Any other path under that prefix is application traffic and crosses the whole pipeline. Create composes the application but opens no listener; `HandleCommand` boots the modules and runs the server for serve, the command with no arguments, while a command that needs only the configuration and the database boots nothing.

<a name="test-through-handler-when-the-port-is-irrelevant"></a>
## Test through Handler when the port is irrelevant

```go example
func assertHomeIsOK(t *testing.T, app *foundation.Application) {
	request := httptest.NewRequest(stdhttp.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	app.Handler().ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}
}
```

- A panic becomes the configured error response and is recorded once.
- A cancelled request reaches downstream database and client calls through context.
- Authentication and policy denial are distinguishable in tests.
- A response is written once, with the intended status and content type.

<a name="each-stage-has-one-question"></a>
## Each stage has one question

Routing chooses the handler. Middleware handles cross-cutting request policy. Controllers translate transport into a service call. Services authorize use cases and repositories perform tenant-scoped data access. Views and resources own the output shape. Keeping those questions separate makes a failed request traceable.

- [Routing](/docs/{{version}}/http/routing)
- [Middleware](/docs/{{version}}/http/middleware)
- [Controllers](/docs/{{version}}/http/controllers)
- [Context](/docs/{{version}}/architecture/context)
- [Framework application](https://pkg.go.dev/github.com/arandu-io/framework/foundation)
