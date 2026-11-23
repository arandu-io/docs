# HTTP resources

- [The response shape is a declared type](#the-response-shape-is-a-declared-type)
- [Finish authorization and validation first](#finish-authorization-and-validation-first)
- [Implement the JSON resource contract once](#implement-the-json-resource-contract-once)
- [Model collection envelopes explicitly](#model-collection-envelopes-explicitly)
- [A rejected form remains a rejected request](#a-rejected-form-remains-a-rejected-request)
- [Public fields are an allowlist](#public-fields-are-an-allowlist)
- [Assert status, shape, and fragment](#assert-status-shape-and-fragment)
- [Resources sit between domain and transport](#resources-sit-between-domain-and-transport)

<a name="the-response-shape-is-a-declared-type"></a>
## The response shape is a declared type

An HTTP resource is the type an endpoint promises to return. The resource type and constructor are checked by the compiler, but `JsonResource.ToArray` returns `map[string]`any, so response keys remain runtime strings. Exact response-shape tests must catch a misspelled, removed, or accidentally added key. The explicit resource still keeps public fields independent from whatever columns happen to exist on a model.

<a name="finish-authorization-and-validation-first"></a>
## Finish authorization and validation first

Build a resource only from data already loaded under a Grant and from request input already validated into its own type. The serializer is not the place to repair authorization or parse input. Decide the HTTP status and content negotiation before rendering so page and JSON callers observe the same outcome.

<a name="implement-the-json-resource-contract-once"></a>
## Implement the JSON resource contract once

`Context.JSON` accepts a `JsonResource`: `ToArray` returns the fields that may leave and With returns top-level metadata about the response. Keep the model in an unexported field, construct the resource from an already authorized value, and answer it directly from the handler.

```go example
type Invoice struct {
	ID, Reference string
	Total         int64
}

type InvoiceResource struct{ invoice Invoice }

func NewInvoiceResource(invoice Invoice) InvoiceResource {
	return InvoiceResource{invoice: invoice}
}

func (r InvoiceResource) ToArray() map[string]any {
	return map[string]any{
		"id": r.invoice.ID, "reference": r.invoice.Reference, "total": r.invoice.Total,
	}
}

func (InvoiceResource) With() map[string]any { return nil }

func Show(ctx *fhttp.Context, invoice Invoice) error {
	return ctx.JSON(stdhttp.StatusOK, NewInvoiceResource(invoice))
}
```

<a name="model-collection-envelopes-explicitly"></a>
## Model collection envelopes explicitly

`Context.JSON` always puts `ToArray` under a top-level data key. In :hesapeVersion `ResourceCollection.ToArray` already returns its own data key, so passing that helper through `Context.JSON` produces `data.data`. Define an explicit collection resource or response shape and pin the complete object in a test instead of assuming a flat envelope. For partial-page requests, render a separate fragment whose root matches the swap target.

> [!NOTE]
> Treat data as reserved metadata. `Context.JSON` merges With after creating the payload, so `With()["data"]` overwrites the resource body rather than being rejected.

<a name="a-rejected-form-remains-a-rejected-request"></a>
## A rejected form remains a rejected request

Return 422 for validation failure even when the browser will swap the returned fragment into the page. YPY renders a 422 inside the form that sent the request, keeps what was typed and focuses the first invalid field, so the status needs no configuration. Returning 200 makes logs and metrics call the attempt successful and treats a rejected form as a success.

<a name="public-fields-are-an-allowlist"></a>
## Public fields are an allowlist

Never encode a model merely because its current fields look harmless. A later secret column would join the payload without the endpoint changing. A `JsonResource` map makes additions deliberate but does not make string keys compile-time checked. Keep With metadata keys distinct from data, and avoid exposing internal tenant keys, digests, storage paths, or policy details.

<a name="assert-status-shape-and-fragment"></a>
## Assert status, shape, and fragment

Test the exact status and JSON keys for an API request, then test the root element and field messages for the partial-page form. Include an assertion that a sensitive model field is absent. A resource test should fail when the response accidentally becomes 200, when a key changes, or when a fragment no longer targets the form being replaced.

<a name="resources-sit-between-domain-and-transport"></a>
## Resources sit between domain and transport

Request validators define accepted input, models define stored data, views define HTML, and resources define output. Keeping those types separate is intentional even when their first version has similar fields. Pagination resources carry page state, while error handling owns the stable problem representation used when no success resource exists.
