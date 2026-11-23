# HTTP client

- [Every outbound request begins with context and a deadline](#every-outbound-request-begins-with-context-and-a-deadline)
- [Send one bounded JSON request](#send-one-bounded-json-request)
- [Add headers, retries, and decoding deliberately](#add-headers-retries-and-decoding-deliberately)
- [Distinguish refusal, timeout, and malformed response](#distinguish-refusal-timeout-and-malformed-response)
- [Test against a local protocol server](#test-against-a-local-protocol-server)
- [Released HTTP client surface](#released-http-client-surface)
- [Instrument the call and keep application policy outside it](#instrument-the-call-and-keep-application-policy-outside-it)

> [!WARNING]
> **Current limitations.** The client propagates context and limits retries to safe requests, but no first-party application journey currently exercises every verb. `Batch.Send` also lacks complete behavioral coverage; prefer the single-request path until that edge is verified.

<a name="every-outbound-request-begins-with-context-and-a-deadline"></a>
## Every outbound request begins with context and a deadline

The Hesape HTTP client wraps `net/http` with context-aware verbs, request options, response decoding, and middleware. The component is released but has no current tagged first-party application consumer, so this guide is limited to the verified package contract rather than claiming default starter wiring.

<a name="send-one-bounded-json-request"></a>
## Send one bounded JSON request

Construct a Factory at boot, set its response-byte ceiling, derive a child deadline, and send typed JSON through a `PendingRequest`. Check transport and then call Throw for a 4xx or 5xx before decoding the eager response.

```go example
func PostPayload(ctx context.Context, endpoint string, payload map[string]any) error {
	factory := client.NewFactory(nil).MaxResponseBytes(2 << 20)
	request := factory.CreatePendingRequest().Timeout(5 * time.Second).AcceptJSON()
	callCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	response, err := request.Post(callCtx, endpoint, payload)
	if err != nil {
		return err
	}
	if err := response.Throw(nil); err != nil {
		return err
	}
	return nil
}
```

<a name="add-headers-retries-and-decoding-deliberately"></a>
## Add headers, retries, and decoding deliberately

Configure base URL and immutable headers at construction, add request-specific authorization without logging it, propagate context, bound response bodies, and retry only idempotent or explicitly idempotency-keyed operations with capped backoff.

- Use context on every verb.
- Set client and operation deadlines.
- Limit response bytes before decoding.
- Instrument through the application Collector.

<a name="distinguish-refusal-timeout-and-malformed-response"></a>
## Distinguish refusal, timeout, and malformed response

DNS, connection, TLS, cancellation, deadline, non-2xx status, and invalid response bodies have different evidence. Never turn an unavailable dependency into an empty success or retry an unsafe mutation merely because the transport closed.

> [!WARNING]
> A timeout of zero or raw default client can hold resources indefinitely. Redact URLs containing credentials and cap error-body capture so a remote service cannot fill logs or memory.

<a name="test-against-a-local-protocol-server"></a>
## Test against a local protocol server

Use `httptest.Server` for method, context cancellation, headers, JSON body, a bounded large response, redirects, slow response, retryable status, unsafe mutation, malformed JSON, and TLS trust. Assert no credential reaches logs.

- Cancellation reaches the server.
- Retries follow idempotency policy.
- Response limits fire before unbounded allocation.
- Observability records URL, status, and duration without secrets.

<a name="released-http-client-surface"></a>
## Released HTTP client surface

`NewFactory`, `Factory.MaxResponseBytes`, `Factory.CreatePendingRequest`, `PendingRequest` verbs, `Response.Throw`, and response decoding form the released surface. Response eagerly reads and closes the body; callers do not close an exported Body field.

| Surface | Released contract |
| --- | --- |
| `NewFactory` / `MaxResponseBytes` | construct the guarded factory |
| `CreatePendingRequest` / Timeout | describe one request |
| `Get/Post/Put/Patch/Delete` | context-aware HTTP verbs |
| `Response.Throw` / Object | status and typed decoding |

<a name="instrument-the-call-and-keep-application-policy-outside-it"></a>
## Instrument the call and keep application policy outside it

Context supplies cancellation, logging records outbound waits, contracts hide provider details from domain services, and testing owns the protocol fixture.

- [Context](/docs/{{version}}/architecture/context)
- [Logging](/docs/{{version}}/operations/logging-observability)
- [Contracts](/docs/{{version}}/architecture/contracts-adapters)
- [Test doubles](/docs/{{version}}/testing/test-doubles)
