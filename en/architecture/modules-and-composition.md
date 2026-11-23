# Modules and composition root

- [The composition root is readable Go](#the-composition-root-is-readable-go)
- [Register one small module](#register-one-small-module)
- [Add lifecycle behavior through optional interfaces](#add-lifecycle-behavior-through-optional-interfaces)
- [Refuse ambiguous composition](#refuse-ambiguous-composition)
- [Boot the same graph a server uses](#boot-the-same-graph-a-server-uses)
- [Required and optional module contracts](#required-and-optional-module-contracts)
- [Interfaces meet at the root](#interfaces-meet-at-the-root)

<a name="the-composition-root-is-readable-go"></a>
## The composition root is readable Go

A module contributes a stable name and routes. Constructors receive every dependency explicitly, `foundation.NewDrivers` builds the concrete adapters from the configuration, and `bootstrap/app.go` lists the modules in `WithProviders`. There is no service container, reflection scan, or hidden resolver between a handler and its collaborators.

<a name="register-one-small-module"></a>
## Register one small module

Define the smallest module that owns the route group, construct its controller from an application service in a provider, and list that value in `WithProviders`. The builder accepts `foundation.Module`, so the compiler proves Name and Routes while bootstrap proves the selected implementation.

```go example
type InvoiceService struct{}

type InvoiceController struct{ service *InvoiceService }

func NewInvoiceController(service *InvoiceService) *InvoiceController {
	return &InvoiceController{service: service}
}

func (c *InvoiceController) Index(ctx *fhttp.Context) error { return ctx.Status(stdhttp.StatusOK) }

type BillingModule struct{ invoices *InvoiceController }

func (m *BillingModule) Name() string { return "billing" }
func (m *BillingModule) Routes(r *fhttp.Router) {
	r.Action("GET", "/invoices", m.invoices.Index).Name("invoices.index")
}

func Compose(drivers *foundation.Drivers, app foundation.Module, service *InvoiceService) (*foundation.Application, error) {
	billing := &BillingModule{invoices: NewInvoiceController(service)}
	return foundation.Configure(drivers).WithProviders(billing, app).Create()
}
```

> [!NOTE]
> `WithProviders` boots the modules in the order it lists them. A module written against `hesape/foundation` whose Routes still takes the handler-only router is listed as `foundation.Component(module)`, which keeps every optional interface the module implements visible to the kernel.

<a name="add-lifecycle-behavior-through-optional-interfaces"></a>
## Add lifecycle behavior through optional interfaces

Lifecycle behavior is expressed through optional interfaces discovered from the registered value. Implement only the contracts the module owns, and never start a goroutine from a constructor or Routes because the application owns process start and shutdown.

- Bootable prepares resources and validates wiring.
- Background starts a loop only when the application runs.
- Closable releases resources during shutdown.
- Migratable, Schedulable, and Health declare schema, tasks, and health checks.
- `RendererProvider` supplies the single renderer used by routes.
- Commandable, CSRFExempt, and `JobHandlers` declare console commands, CSRF exemptions for the module's own receivers, and job handlers; Create collects them from every provider.

<a name="refuse-ambiguous-composition"></a>
## Refuse ambiguous composition

Empty or duplicate names, two renderer providers, a boot failure, and route conflicts are wiring errors. Create and Boot also refuse a command name declared twice, a job name given two handlers, an application route that answers a path the framework mounts under `/_arandu/`, and a module CSRF exemption that covers another owner's write route. They fail before traffic. Registration order is observable for boot, routes, migrations, and shutdown, so keep it intentional instead of relying on map iteration.

<a name="boot-the-same-graph-a-server-uses"></a>
## Boot the same graph a server uses

Boot the same graph the server will use, then inspect its routes and lifecycle rather than testing a controller in isolation. Assert boot order, background start, health result, and reverse-safe cleanup for every resource the module owns.

- Boot returns no error for the complete graph.
- The route inventory names the module and expected handlers.
- Background work starts at Run, not construction.
- Shutdown closes owned resources in a safe order.

<a name="required-and-optional-module-contracts"></a>
## Required and optional module contracts

`foundation.Module` requires `Name()` string and `Routes(*http.Router)`. `foundation.Configure`, `WithProviders`, and Create compose the graph; `Application.HandleCommand`, Boot, Run, Handler, and Shutdown own it. Optional contracts include Bootable, Background, `RendererProvider`, Closable, Diagnostic, Schedulable, Migratable, Health, Ready, Commandable, CSRFExempt, and `JobHandlers`, and `foundation.Component` adapts a `hesape/foundation` module.

<a name="interfaces-meet-at-the-root"></a>
## Interfaces meet at the root

Contracts keep services independent of adapters; configuration selects those adapters; routes expose already-constructed controllers. Directory structure makes each owner visible, while request lifecycle shows when the resulting graph runs.

- [Configuration](/docs/{{version}}/architecture/configuration)
- [Contracts and adapters](/docs/{{version}}/architecture/contracts-adapters)
- [Foundation Module](https://pkg.go.dev/github.com/arandu-io/framework/foundation#Module)
