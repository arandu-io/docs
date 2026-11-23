# Routing

- [Routes are executable contracts](#routes-are-executable-contracts)
- [Register and reverse one route](#register-and-reverse-one-route)
- [Groups, constraints, and the exact root](#groups-constraints-and-the-exact-root)
- [Fail while building, not after a click](#fail-while-building-not-after-a-click)
- [Inspect and dispatch the same table](#inspect-and-dispatch-the-same-table)
- [Released routing API](#released-routing-api)
- [Continue the request path](#continue-the-request-path)

> [!WARNING]
> **Current limitations.** The router and named URLs are released. Route and configuration caches are intentionally absent, and advanced routing examples are limited to patterns exercised by the skeleton and feature suite.

<a name="routes-are-executable-contracts"></a>
## Routes are executable contracts

The router binds an HTTP method and a Go 1.22 path pattern to either a standard handler or a controller action. Naming the route at registration gives templates, redirects, tests, and the CLI one source of truth for its address. A group adds a prefix and middleware while sharing the parent route table.

<a name="register-and-reverse-one-route"></a>
## Register and reverse one route

Action adapts `func(*http.Context)` error, while Get, Post, Put, Patch, and Delete accept standard `net/http` handlers. Name the route before serving traffic and use the route table to reverse its path with the exact parameter count.

```go example
func Register(router *fhttp.Router, sessions *security.SessionStore, show func(*fhttp.Context) error) (string, error) {
	router.Action("GET", "/invoices/{invoice}", show,
		middleware.RequireAuth(sessions),
	).Name("invoices.show")

	path, err := router.Table().URL("invoices.show", "inv_123")
	if err != nil {
		return "", err
	}
	return path, nil
}
```

<a name="groups-constraints-and-the-exact-root"></a>
## Groups, constraints, and the exact root

Groups share a prefix and policy boundary. Constraints refuse malformed identifiers before an action runs, and the exact-root pattern prevents a home route from catching every subtree in the standard `ServeMux` grammar. A `ConnectRPC` handler can be mounted with Post on its generated service path; the same TLS listener then serves REST over HTTP/1 and Connect, `gRPC`, or gRPC-Web over HTTP/2.

- Use `Group(prefix, middleware...)` for shared path and policy.
- Use `WhereNumber`, `WhereUuid`, `WhereUlid`, or `WhereIn` for route-level identifier shape.
- Use `/{name}` for one segment and `/{path...}` only when the action owns the remainder.
- Use `/{$}` for the exact application root.
- For RPC, configure byte, concurrency, tenant, duration, server-timeout, origin, and header limits explicitly before mounting the handler.

<a name="fail-while-building-not-after-a-click"></a>
## Fail while building, not after a click

Unknown route names and the wrong parameter count make `Routes.URL` return an error. Constraint mismatches are 404 responses, not authorization results. Duplicate or conflicting registrations are boot defects and registration order must stay deterministic. Boot also refuses an application route that answers a path the framework mounts under `/_arandu/`, such as `/_arandu/health`, in a method the framework route does not take, naming the module and the pattern. The usual cause is `r.Post("/", ...)`, whose pattern matches every path below the root; register the root as "`/{$}`". RPC authentication runs before request decoding, but it identifies a Subject only: application services still authorize and issue Grants. Admission closes during shutdown, streams observe cancellation, and the framework does not open an h2c or second RPC listener.

> [!WARNING]
> `Route.SetURI` changes metadata after registration but cannot move the pattern already installed in the underlying `ServeMux`. Define the final URI before the application boots.

<a name="inspect-and-dispatch-the-same-table"></a>
## Inspect and dispatch the same table

Run `aru route:list` and dispatch the same table through httptest. Cover a valid identifier, a constraint mismatch, an unsupported method, the exact root, and the named URL. For RPC, exercise unary, client-streaming, server-streaming, and bidirectional calls through Connect, `gRPC`, and gRPC-Web; also prove CORS refusal, global and per-tenant admission, body limits, cancellation, and graceful shutdown on the shared listener.

- The inventory names the method, pattern, route, and module.
- A valid constrained parameter reaches the action.
- A malformed constrained parameter returns 404.
- The exact root does not swallow an unknown path.
- RPC authorization and all four streaming shapes run on the application listener.

<a name="released-routing-api"></a>
## Released routing API

`http.NewRouter` and `Router.Action` create and populate the table. Router verb methods register standard handlers; `Router.Group` composes prefixes and middleware; `Route.Name` and `Route.Where*` add stable metadata; `Routes.URL` builds a path or reports a mismatch. `rpc.NewHandler` returns the generated service path and handler for an explicit Factory, bearer Authenticator, and Config. `rpc.Subject` reads the authenticated identity inside a procedure; it is not an authorization result.

| Surface | Released contract |
| --- | --- |
| `Router.Action` | register a controller action |
| `Get/Post/Put/Patch/Delete` | register a standard handler |
| Group and Name | compose prefixes and stable addresses |
| `Where*` and `Routes.URL` | constrain parameters and reverse paths |
| `rpc.NewHandler` | build one Connect, `gRPC`, and gRPC-Web handler with explicit limits |
| `rpc.Subject` | read the authenticated Subject before application policy issues a Grant |

<a name="continue-the-request-path"></a>
## Continue the request path

REST remains the default request-response resource boundary. `gRPC` adds typed unary and streaming calls on the same HTTP server; it is not a second process or an authorization shortcut. Controllers own transport translation, middleware owns cross-cutting request policy, and URL signing owns integrity after a route is selected.

- [Communication patterns](/docs/{{version}}/architecture/contracts-adapters)
- [Controllers](/docs/{{version}}/http/controllers)
- [Middleware](/docs/{{version}}/http/middleware)
- [URLs](/docs/{{version}}/http/urls)
