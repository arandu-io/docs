# Package development

- [A package is a Go module with an explicit entry point](#a-package-is-a-go-module-with-an-explicit-entry-point)
- [Implement the required module surface](#implement-the-required-module-surface)
- [Add optional lifecycle contracts deliberately](#add-optional-lifecycle-contracts-deliberately)
- [Make missing wiring fail at construction](#make-missing-wiring-fail-at-construction)
- [Test the package through its exported contract](#test-the-package-through-its-exported-contract)
- [Required and optional interfaces](#required-and-optional-interfaces)
- [Put primitives in components and state with its owner](#put-primitives-in-components-and-state-with-its-owner)

> [!WARNING]
> **Current limitations.** The package skeleton is released and self-tested, but no independent tagged application imports it. Treat its module, policy, resource, migration, and route surfaces as a verified starting point and add a consuming application journey before claiming integration.

<a name="a-package-is-a-go-module-with-an-explicit-entry-point"></a>
## A package is a Go module with an explicit entry point

An installable package has its own module path and exposes behavior through normal Go APIs. Framework integration uses `foundation.Module`, whose required contract is only a stable Name and route registration. Dependencies stay explicit in the constructor and no container or reflection resolves them behind the composition root.

<a name="implement-the-required-module-surface"></a>
## Implement the required module surface

Name returns a lowercase stable identifier and Routes receives the application's router. Optional capabilities are structural interfaces discovered only when implemented; migrations and scheduled work are not methods every module must pretend to have. The application constructs the module and lists it in `WithProviders` in `bootstrap/app.go`.

```go example
type Handler struct{}

func (*Handler) Index(ctx *fhttp.Context) error { return ctx.Status(stdhttp.StatusOK) }

type Module struct {
	handler *Handler
}

func (*Module) Name() string { return "invoices" }

func (m *Module) Routes(r *fhttp.Router) {
	r.Action("GET", "/invoices", m.handler.Index).Name("invoices.index")
}
```

<a name="add-optional-lifecycle-contracts-deliberately"></a>
## Add optional lifecycle contracts deliberately

Implement `foundation.Migratable` only when the module owns schema; its exact method is `Migrations() []foundation.Migration`. Implement `foundation.Schedulable` only when it declares kernel work; its exact method is `Schedule() []foundation.Task`. Keep views, translations, policies, handlers, and tests within the module that owns them. A package declares work; it does not start an unmanaged goroutine during import.

<a name="make-missing-wiring-fail-at-construction"></a>
## Make missing wiring fail at construction

Reject nil required dependencies in the module constructor and let duplicate route names, invalid configuration, or migration conflicts surface during bootstrap. Do not silently substitute in-memory adapters when production wiring requested a shared backend. An optional interface should be absent, not implemented with a no-op that implies support.

<a name="test-the-package-through-its-exported-contract"></a>
## Test the package through its exported contract

Compile a small consumer that constructs the module, registers routes, and exercises one request. Run migration or scheduling contract tests only for interfaces the module implements. The package's own tests verify its repositories and policies; the consuming application still tests its chosen wiring and configuration.

<a name="required-and-optional-interfaces"></a>
## Required and optional interfaces

`foundation.Module` requires `Name()` string and `Routes(*http.Router)`. `foundation.Migratable` requires `Migrations() []foundation.Migration`, and `foundation.Schedulable` requires `Schedule() []foundation.Task`; both remain optional and do not enlarge Module. Constructors remain ordinary functions, and route handlers remain ordinary net/http-compatible handlers through the framework router.

<a name="put-primitives-in-components-and-state-with-its-owner"></a>
## Put primitives in components and state with its owner

A reusable calculation or adapter can live in the component library without becoming an installable module. A capability that owns schema, routes, and policy belongs in a package. The application composition root is where package adapters meet concrete configuration, which keeps the package reusable without hiding how it runs.
