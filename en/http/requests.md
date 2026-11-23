# HTTP requests

- [Read transport data without inventing identity](#read-transport-data-without-inventing-identity)
- [Assemble one typed input contract](#assemble-one-typed-input-contract)
- [Choose the reader that preserves meaning](#choose-the-reader-that-preserves-meaning)
- [Bound the body before decoding it](#bound-the-body-before-decoding-it)
- [Exercise every input channel and limit](#exercise-every-input-channel-and-limit)
- [Released request surfaces](#released-request-surfaces)
- [Move from transport to a trusted operation](#move-from-transport-to-a-trusted-operation)

<a name="read-transport-data-without-inventing-identity"></a>
## Read transport data without inventing identity

The preferred request boundary is the slim framework Context plus standard `net/http` primitives. Read route parameters, query values, headers, and a bounded body only where the action needs them; the broad Hesape Request wrapper remains a compatibility surface rather than a required controller dependency.

<a name="assemble-one-typed-input-contract"></a>
## Assemble one typed input contract

Bind the few accepted fields into a dedicated input type whose form tags name every field that may arrive, normalize canonical values, and pass only validated data to the service. A type with a Rules method has those rules run over the raw input first. Bind writes only tagged fields, trims each value, and reports a value that does not convert as a field error instead of a silent zero. Preserve `ctx.Ctx` so cancellation and authenticated identity stay attached to every downstream call.

```go example
type CreateInvoiceInput struct {
	Reference string `form:"reference"`
	Amount    int64  `form:"amount"`
}

func Create(ctx *fhttp.Context) error {
	var input CreateInvoiceInput
	if err := ctx.Bind(&input); err != nil {
		return err
	}
	return nil
}
```

<a name="choose-the-reader-that-preserves-meaning"></a>
## Choose the reader that preserves meaning

Use trusted-proxy middleware before reading a peer-derived client address. Limit request bodies before decoding, close uploaded streams, and keep route identifiers separate from writable fields so a form cannot replace ownership or tenancy. Input and Bind read one map: the query string for GET and HEAD, and the body alone for POST, PUT, PATCH, and DELETE, so a value in the action URL is read explicitly with Query.

- Use Param for matched route segments.
- Use Query for optional URL controls.
- Use Bind for declared form fields; unsafe methods read only the body.
- Use the standard Request for headers, body, and context.

<a name="bound-the-body-before-decoding-it"></a>
## Bound the body before decoding it

Malformed numbers, unsupported media types, oversized bodies, missing required values, and cancelled reads must remain explicit. Do not silently coerce a malformed value to zero or continue after a decoder stopped at a partial payload.

> [!WARNING]
> Forwarded headers are caller input until `TrustProxies` accepts the immediate peer. Never key authentication throttles, audit identity, or signed links from an untrusted forwarding value.

<a name="exercise-every-input-channel-and-limit"></a>
## Exercise every input channel and limit

Send valid, missing, malformed, oversized, and cancelled requests through the real route. Assert the typed input given to the service, and prove tenant or ownership fields supplied by the caller are ignored or refused.

- Only declared input reaches the service.
- Body and upload limits fire before allocation grows unbounded.
- Cancellation stops parsing and downstream work.
- Untrusted forwarded values do not replace the peer address.

<a name="released-request-surfaces"></a>
## Released request surfaces

Framework `http.Context` exposes Ctx, Param, Query, Input, Bind, User, response helpers, and the underlying request and writer. Standard `net/http` owns headers, bodies, cookies, and request context. The larger Hesape Request remains available when its compatibility helpers are specifically needed.

<a name="move-from-transport-to-a-trusted-operation"></a>
## Move from transport to a trusted operation

Validation converts this transport input into typed application input that satisfies the declared rules. Controllers choose the use case, trusted proxies establish the network boundary, and responses encode the result.

- [Validation](/docs/{{version}}/http/validation)
- [Controllers](/docs/{{version}}/http/controllers)
- [Responses](/docs/{{version}}/http/responses)
