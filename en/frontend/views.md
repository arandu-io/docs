# Views

- [A view name resolves to compiled Go](#a-view-name-resolves-to-compiled-go)
- [Render typed data end to end](#render-typed-data-end-to-end)
- [Register page, layout, and fragment packages](#register-page-layout-and-fragment-packages)
- [Diagnose names, imports, and stale output](#diagnose-names-imports-and-stale-output)
- [Prove build, escaping, and response shape](#prove-build-escaping-and-response-shape)
- [Released view registry and renderer](#released-view-registry-and-renderer)
- [Write and ship the view source](#write-and-ship-the-view-source)

<a name="a-view-name-resolves-to-compiled-go"></a>
## A view name resolves to compiled Go

A `.kyse.go` source declares its data contract and markup. `aru view:build` writes generated Go under `storage/framework/views`; importing each generated package runs `view.Register`, and `framework/view.NewModule` supplies the one renderer used by `Context.View` and `Context.Fragment`.

<a name="render-typed-data-end-to-end"></a>
## Render typed data end to end

Define page data in the source and pass that exact generated type from the controller. A misspelled field fails compilation at the source line, and the wrong runtime data type returns a named rendering error instead of drawing zero values.

```kyse
//go:build kyse

package views

@go
type GreetingData struct { Name string }
@endgo

<h1>Hello {{ .Name }}</h1>
```

> [!NOTE]
> Page data is a public payload. `Context.View` answers a request that sends Accept: `application/vnd.arandu.view+json` with the view name and every exported field of the data as JSON, including fields the template never prints. Never put secrets, password hashes, tokens, or internal tenant keys in view data; pass only what the page may show.

<a name="register-page-layout-and-fragment-packages"></a>
## Register page, layout, and fragment packages

Import every generated directory in bootstrap because nested directories are distinct Go packages. A nested source such as `invoices/show.kyse.go` registers a dotted name, and a layout-free form view is a separate fragment target.

- Use `ctx.View` with a typed page value for full documents.
- Use `ctx.Fragment` with a layout-free typed view for swaps.
- Keep data composition in controllers and services.
- Regenerate rather than editing `storage/framework/views`.

<a name="diagnose-names-imports-and-stale-output"></a>
## Diagnose names, imports, and stale output

No registered views means generated packages are absent or not imported. An unknown dotted name, wrong data type, duplicate registration, stale generated output, or missing `RendererProvider` must stop the relevant response or boot path visibly.

> [!WARNING]
> Generated application views are normally gitignored. A clean checkout must run `aru view:build` before Go can resolve controllers that import generated types; normal views buffer output and are not unbounded streaming responses.

<a name="prove-build-escaping-and-response-shape"></a>
## Prove build, escaping, and response shape

Compile views, render with the exact data type, and assert HTML escaping, content type, full-page and fragment roots, unknown-name failure, wrong-data failure, and `view.Registered` containing every literal used by a controller.

- `aru view:build` and `go build` succeed.
- User text is escaped.
- A wrong data type is an error.
- Every registered controller view name resolves.

<a name="released-view-registry-and-renderer"></a>
## Released view registry and renderer

`view.Register` and Registered form the generated registry, `view.NewRenderer` buffers a typed render, `view.WrongData` standardizes type mismatch, `Context.View` and Fragment choose HTTP shape, and `framework/view.NewModule` provides the renderer and asset routes.

| Surface | Released contract |
| --- | --- |
| `view.Register` / Registered | generated registration and diagnostics |
| `view.NewRenderer` / `WrongData` | buffered render and typed mismatch |
| `Context.View` / Fragment | HTTP response selection |
| `framework/view.NewModule` | renderer provider and assets |

<a name="write-and-ship-the-view-source"></a>
## Write and ship the view source

Kyse explains source syntax, assets trace compiled bytes, controllers assemble typed data, and response helpers select the final document or fragment.

- [Kyse templates](/docs/{{version}}/frontend/kyse-templates)
- [Assets](/docs/{{version}}/frontend/assets)
- [Controllers](/docs/{{version}}/http/controllers)
