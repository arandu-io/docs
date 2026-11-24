# HTTP tests

- [Dispatch the registered pipeline with httptest](#dispatch-the-registered-pipeline-with-httptest)
- [Boot the real graph and keep external adapters controlled](#boot-the-real-graph-and-keep-external-adapters-controlled)
- [Send one request through Handler](#send-one-request-through-handler)
- [Assert protocol, session, YPY, and policy paths](#assert-protocol-session-ypy-and-policy-paths)
- [Do not let production recovery hide a test defect](#do-not-let-production-recovery-hide-a-test-defect)
- [Keep credentials and tenant data out of fixtures and output](#keep-credentials-and-tenant-data-out-of-fixtures-and-output)
- [Cover the complete route matrix](#cover-the-complete-route-matrix)
- [Move service and adapter details to focused contracts](#move-service-and-adapter-details-to-focused-contracts)

<a name="dispatch-the-registered-pipeline-with-httptest"></a>
## Dispatch the registered pipeline with httptest

An HTTP feature test dispatches `foundation.Application.Handler` through `net/http/httptest`. It crosses global middleware, route matching, route middleware, controller, service, renderer, and error response without opening a listener.

<a name="boot-the-real-graph-and-keep-external-adapters-controlled"></a>
## Boot the real graph and keep external adapters controlled

Build the same modules and configuration as serving code, substitute only explicit external adapters, compile views, and use a dedicated database and session backend. Keep middleware order identical to production.

- Register real routes and renderer.
- Use a test database with tenant fixtures.
- Capture mail and outbound HTTP.
- Choose memory or shared session according to the behavior under test.

<a name="send-one-request-through-handler"></a>
## Send one request through Handler

Create a standard request, dispatch it to the application Handler, and assert status, headers, and body. When the journey spans login, CSRF, rotation, and logout, `test.NewClient(t, app.Handler())` from `framework/test` is the same dispatch with a cookie jar, the CSRF token read off the last page it loaded, and no redirect following. Its `*test.Response` asserts with `AssertStatus`, `AssertOk`, `AssertSee`, `AssertDontSee`, `AssertRedirect`, and the Assert methods for headers, cookies, and JSON. `test.Serve` boots the application on a free port over a database of the test's own, and `test.As` signs a subject in through the application's own session.

```go example
func assertHomeIsOK(t *testing.T, app *foundation.Application) {
	request := httptest.NewRequest(stdhttp.MethodGet, "/invoices/inv_123", nil)
	response := httptest.NewRecorder()
	app.Handler().ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}
}
```

<a name="assert-protocol-session-ypy-and-policy-paths"></a>
## Assert protocol, session, YPY, and policy paths

Test normal and YPY-Request representation, named redirects, JSON resources, form validation, CSRF, trusted proxy, authentication, policy denial, and cancellation. Assert the service or persistent outcome as well as response text.

- Assert Content-Type and exact status.
- Carry cookies across the real flow.
- Test full page and fragment separately.
- Inspect side effects after the response.

<a name="do-not-let-production-recovery-hide-a-test-defect"></a>
## Do not let production recovery hide a test defect

A panic should reach the testing boundary when configured to expose defects; do not write a green test merely because production recovery converted it to 500. Partial writes, swallowed errors, unknown routes, and stale generated views need explicit cases. A controller unit fake does not prove middleware, route names, session cookies, CSRF, renderer registration, or error mapping, so keep at least one complete route journey for every public flow.

<a name="keep-credentials-and-tenant-data-out-of-fixtures-and-output"></a>
## Keep credentials and tenant data out of fixtures and output

> [!WARNING]
> Use synthetic credentials and tenants, redact failure output, never call live providers, and assert a foreign tenant receives the same bounded absence. Test Secure, `HttpOnly`, `SameSite`, signature, and CSRF attributes as protocol details.

A bypassed middleware helper can make tests easy and production wrong; authenticate through the released session or deliberately construct Subject only in a service-level test.

<a name="cover-the-complete-route-matrix"></a>
## Cover the complete route matrix

Cover method and path mismatch, valid request, malformed input, guest, denied Subject, another tenant, YPY, JSON preference, dependency outage, panic, request cancellation, and session rotation. Run the suite with race detection.

- Every route has a happy and refusal path.
- Cross-tenant records remain invisible.
- Enhanced requests retain truthful status.
- Cancellation reaches downstream calls.

<a name="move-service-and-adapter-details-to-focused-contracts"></a>
## Move service and adapter details to focused contracts

Controllers translate transport, the request lifecycle explains the pipeline, test doubles capture external adapters, and database tests prove storage isolation.

- [Controllers](/docs/{{version}}/http/controllers)
- [Request lifecycle](/docs/{{version}}/start/request-lifecycle)
- [Test doubles](/docs/{{version}}/testing/test-doubles)
- [Database testing](/docs/{{version}}/testing/database)
