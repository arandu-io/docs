# Controllers

- [Translate HTTP and keep policy downstream](#translate-http-and-keep-policy-downstream)
- [Inject the use case](#inject-the-use-case)
- [Render one authenticated record](#render-one-authenticated-record)
- [Keep actions short and explicit](#keep-actions-short-and-explicit)
- [Preserve failure meaning](#preserve-failure-meaning)
- [The controller is not an authorization authority](#the-controller-is-not-an-authorization-authority)
- [Test transport and policy separately](#test-transport-and-policy-separately)
- [Complete the controller boundary](#complete-the-controller-boundary)

<a name="translate-http-and-keep-policy-downstream"></a>
## Translate HTTP and keep policy downstream

A controller reads transport input, asks an application service for an outcome, and chooses one HTTP response. It does not mint authorization grants, open a second database path, or trust tenant and ownership fields supplied by the browser.

<a name="inject-the-use-case"></a>
## Inject the use case

Define the dependency in the consumer package and inject it through the controller constructor. Protect the route with authentication middleware when the action requires a subject: `RequireAuth` puts the subject it loaded from the session on the request context, so `ctx.User` reports it behind the guard, and `middleware.LoadSubject` does the same on a public route while letting a guest through with no subject. The service or policy remains responsible for authorizing the requested action on the requested record.

- Construct controllers in a provider in app/Providers.
- Register actions on named routes.
- Keep validation input separate from stored models.
- Let the service own authorization and persistence.

<a name="render-one-authenticated-record"></a>
## Render one authenticated record

The context deadline and authenticated Subject cross the service boundary, while the path identifier remains untrusted input rather than evidence of access. Render the already authorized result through a view or declared resource.

```go example
type Invoice struct{ ID string }

type InvoiceFinder interface {
	Find(context.Context, security.Subject, string) (Invoice, error)
}

type InvoiceController struct{ invoices InvoiceFinder }
type InvoiceShowData struct{ Invoice Invoice }

func (c InvoiceController) Show(ctx *fhttp.Context) error {
	subject, ok := ctx.User()
	if !ok {
		hhttp.Refuse(ctx.Response, ctx.Request, stdhttp.StatusUnauthorized, "sign in required")
		return nil
	}
	invoice, err := c.invoices.Find(ctx.Ctx(), subject, ctx.Param("invoice"))
	if err != nil {
		return err
	}
	return ctx.View("invoices.show", InvoiceShowData{Invoice: invoice})
}
```

<a name="keep-actions-short-and-explicit"></a>
## Keep actions short and explicit

Use Param, Query, and Input only to assemble typed input. Return operational errors to centralized recovery, non-empty `validation.Errors` to the reject path, and exactly one of View, Fragment, JSON, `RedirectRoute`, or Status for the response.

- Pass `ctx.Ctx` into the use case.
- Read identity through the context's reported presence.
- Return only non-empty validation errors.
- Choose one response helper after the service outcome.

<a name="preserve-failure-meaning"></a>
## Preserve failure meaning

Malformed input, absent records, denied actions, conflicts, and internal faults are distinct service outcomes. Turning every error into 500 erases useful behavior; swallowing one can accidentally answer 200 with an empty body.

> [!WARNING]
> The action adapter panics when returned `validation.Errors` is empty. Use if `errs.Any()` before returning it, and return nil when validation passed.

<a name="the-controller-is-not-an-authorization-authority"></a>
## The controller is not an authorization authority

> [!WARNING]
> Never derive the tenant from a form, query, or path; never load tenant-owned data without Subject or Grant; never serialize a persistence model wholesale when it contains hashes, internal tenant keys, or audit fields.

Request bodies and structs that can hold credentials need explicit JSON and slog redaction before they are logged. A controller's convenience is not permission to expose transport input.

<a name="test-transport-and-policy-separately"></a>
## Test transport and policy separately

Use httptest through the registered route to assert status, body or redirect, parameters, authentication refusal, and a service error. In a focused unit test, inject a fake finder and assert the same context, Subject, and id reach it.

- A guest is refused before the use case.
- The exact Subject and id reach the service.
- Validation follows the 422 or redirect contract.
- Policy behavior remains tested outside the controller.

<a name="complete-the-controller-boundary"></a>
## Complete the controller boundary

Request types own accepted input, responses own transport semantics, resources own public JSON fields, and authorization owns the Grant. A short controller makes each of those contracts visible.

- [Requests](/docs/{{version}}/http/requests)
- [Responses](/docs/{{version}}/http/responses)
- [Authorization](/docs/{{version}}/security/authorization)
- [Resources](/docs/{{version}}/http/resources)
