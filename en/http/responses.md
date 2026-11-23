# HTTP responses

- [Commit one response shape once](#commit-one-response-shape-once)
- [Return a public resource, not a model](#return-a-public-resource-not-a-model)
- [Match the helper to the client contract](#match-the-helper-to-the-client-contract)
- [Do not write past a committed response](#do-not-write-past-a-committed-response)
- [Assert protocol details, not only bodies](#assert-protocol-details-not-only-bodies)
- [Released response helpers](#released-response-helpers)
- [Keep representation at the edge](#keep-representation-at-the-edge)

<a name="commit-one-response-shape-once"></a>
## Commit one response shape once

A response is one explicit status, header set, and body representation. Framework Context helpers preserve this decision for full pages, fragments, declared JSON resources, named redirects, and empty statuses without letting a handler write twice.

<a name="return-a-public-resource-not-a-model"></a>
## Return a public resource, not a model

Return a declared `JsonResource` for JSON rather than a raw model. `ToArray` is the public field allowlist and With supplies top-level metadata; `Context.JSON` selects status and encoding while retaining the resource boundary.

```go example
type Invoice struct {
	ID    string
	Total int64
}

type InvoiceResource struct{ invoice Invoice }

func (r InvoiceResource) ToArray() map[string]any {
	return map[string]any{"id": r.invoice.ID, "total": r.invoice.Total}
}
func (InvoiceResource) With() map[string]any { return nil }

func Show(ctx *fhttp.Context, invoice Invoice) error {
	return ctx.JSON(stdhttp.StatusOK, InvoiceResource{invoice: invoice})
}
```

<a name="match-the-helper-to-the-client-contract"></a>
## Match the helper to the client contract

Use View for a full document, Fragment for a named partial, and `RedirectRoute` for a route-table destination. `RedirectRoute` returns 303 to a normal browser and YPY-Redirect with 204 to an enhanced request so both clients navigate correctly. Use `EventStream` for one-way live server updates; `EventStreamWithOptions` adds heartbeat comments while every event is flushed immediately.

- Choose status before writing the body.
- Set response headers before the first write.
- Return 422 validation fragments as the content of the form they reject, which YPY renders in place.
- Keep JSON resources separate from database models.
- For SSE, create events with id, event, data, retry, or Comment and stop producing as soon as the request context is cancelled.

<a name="do-not-write-past-a-committed-response"></a>
## Do not write past a committed response

Encoding errors, unknown view names, missing route names, and writes after a terminal response are not empty successes. Preserve 404, 409, 422, and 503 meanings instead of reducing every failure to 200 or 500. An SSE stream is not durable replay: disconnect cancellation stops the producer and deliberately omits a closing event. If a long stream must outlive the ordinary server write deadline, clear that deadline for this response with `http.NewResponseController` before streaming.

> [!WARNING]
> YPY renders a 422 inside the form that sent the request and decides every other 4xx and 5xx by its error rules. Changing a 422 to 200 would make metrics and logs call rejected input successful.

<a name="assert-protocol-details-not-only-bodies"></a>
## Assert protocol details, not only bodies

Assert status, content type, headers, and body for normal and enhanced requests. Test `RedirectRoute` in both modes, exact JSON keys, absent sensitive fields, a 422 fragment root, and an encoding or render error. For SSE, assert text/event-stream, no-cache, X-Accel-Buffering: no, multiline data framing, metadata validation, heartbeat comments, a flush per event, and prompt exit on cancellation and application shutdown.

- Normal redirect is 303 with Location.
- Enhanced redirect is 204 with YPY-Redirect.
- `JsonResource` output has only allowlisted keys.
- Validation fragments retain 422 and are the content of the form.
- SSE cancellation stops production without a synthetic final event.

<a name="released-response-helpers"></a>
## Released response helpers

`Context.View` and Fragment render HTML, JSON accepts `hhttp.JsonResource`, `RedirectRoute` reverses a named route, and Status writes a terminal status. `hhttp.Refuse` supplies bounded public errors without exposing internal causes. `ResponseFactory.EventStream` and `EventStreamWithOptions` build SSE responses; `hhttp.NewServerSentEvent` and Comment encode event metadata and data without permitting newline or NUL injection.

<a name="keep-representation-at-the-edge"></a>
## Keep representation at the edge

REST responses finish one request. SSE keeps one HTTP response open for server-to-client updates but does not accept client messages; use `WebSocket` for bidirectional sessions and a durable store when clients must replay missed facts. Resources define JSON fields, views define HTML, controllers choose the representation, and centralized error handling owns failures that escape the action.

- [Communication patterns](/docs/{{version}}/architecture/contracts-adapters)
- [Joaju `WebSocket`](/docs/{{version}}/async/joaju)
- [Resources](/docs/{{version}}/http/resources)
- [Errors](/docs/{{version}}/http/errors)
