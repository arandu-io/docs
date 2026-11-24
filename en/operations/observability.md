# Debug console and tracing

- [Inspect one request without making tracing universal](#inspect-one-request-without-making-tracing-universal)
- [Choose development access or an explicit production secret](#choose-development-access-or-an-explicit-production-secret)
- [Record a request and read it by request id](#record-a-request-and-read-it-by-request-id)
- [Read findings before the raw timeline](#read-findings-before-the-raw-timeline)
- [Know why a request has no record](#know-why-a-request-has-no-record)
- [Minimize the diagnostic data window](#minimize-the-diagnostic-data-window)
- [Test both the open and closed states](#test-both-the-open-and-closed-states)
- [Reproduce narrowly and remove access](#reproduce-narrowly-and-remove-access)
- [Connect traces to logs, health, and outbound calls](#connect-traces-to-logs-health-and-outbound-calls)

<a name="inspect-one-request-without-making-tracing-universal"></a>
## Inspect one request without making tracing universal

The debug console is a bounded in-process recorder for recent request Collectors: timeline shares, SQL with origin and arguments, dumps, events, views, outbound HTTP, slow statements, and suspected N+1 queries. It is not hosted APM, durable traces, metric history, or a cross-replica query service.

<a name="choose-development-access-or-an-explicit-production-secret"></a>
## Choose development access or an explicit production secret

Development mounts the recorder and `/_arandu/debug`. Outside development, leave `ARANDU_TRACING_SECRET` empty unless controlled production tracing is required; then use at least 16 bytes, TLS, and X-Arandu-Trace only on the request and console reads being investigated.

- Keep production tracing off by default.
- Treat the tracing secret as cross-tenant diagnostic access.
- Validate `ARANDU_EDITOR` at boot if source links are enabled.
- Target a known replica during an incident.

<a name="record-a-request-and-read-it-by-request-id"></a>
## Record a request and read it by request id

Make a request, retain its X-Request-ID, and inspect the list or detail through `aru trace`. The CLI reads the same protected debug JSON route as the browser console and can target another host explicitly.

```bash
set -euo pipefail
aru dev &
DEV_PID=$!
trap 'kill "$DEV_PID"' EXIT
REQUEST_ID=$(curl --retry 20 --retry-connrefused -si http://127.0.0.1:8080/ | awk -F': ' 'tolower($1) == "x-request-id" { gsub("\\r", "", $2); print $2; exit }')
test -n "$REQUEST_ID" || { echo "response has no X-Request-ID" >&2; exit 1; }
aru trace
aru trace "$REQUEST_ID"
```

<a name="read-findings-before-the-raw-timeline"></a>
## Read findings before the raw timeline

Read slow and N+1 findings before the raw timeline, inspect SQL origin and duration, outbound calls, events, and temporary non-secret dumps. Five identical statements are a diagnosis lead, not proof that repetition is wrong.

- Find the request by id.
- Inspect findings before every event.
- Instrument dependencies through the Collector.
- Use shared Gauges only for current-value process measurements.

<a name="know-why-a-request-has-no-record"></a>
## Know why a request has no record

No record can mean tracing was disabled, the header was absent, middleware lost the Collector, another replica served the request, or the bounded ring evicted it. An uninstrumented dependency remains unattributed, and protected or absent records deliberately answer 404. The recorder is local memory and disappears on restart. Tracing one production request does not collect other requests. Arandu ships no telemetry exporter: durable traces and metrics belong to a backend you run beside the application, and the structured logs carry `request_id` and tenant so the two can be correlated.

<a name="minimize-the-diagnostic-data-window"></a>
## Minimize the diagnostic data window

> [!WARNING]
> Console rows can include SQL arguments, dumps, event payloads, paths, and tenant data. Never expose the route through a convenience proxy, log the secret, put it in a URL, or enable the header broadly. Rotate or remove access after the incident.

<a name="test-both-the-open-and-closed-states"></a>
## Test both the open and closed states

Test development index and detail, production disabled 404, wrong-secret 404, constant-time accepted header, request-scoped collection, findings, ring eviction, and the independence of live and health probes from console access.

- Only explicitly traced production requests are recorded.
- Wrong credentials reveal no console.
- The oldest record is evicted at capacity.
- Health remains a separate operational contract.

<a name="reproduce-narrowly-and-remove-access"></a>
## Reproduce narrowly and remove access

Use logs and health first. If necessary, route one controlled request to a known replica, fetch it immediately, preserve minimal redacted evidence, then remove or rotate the secret. Restarting clears evidence but does not repair the dependency.

<a name="connect-traces-to-logs-health-and-outbound-calls"></a>
## Connect traces to logs, health, and outbound calls

Structured logs provide durable correlation, instrumented clients explain outbound waits, deployment owns probes, and query documentation explains SQL findings.

- [Logging](/docs/{{version}}/operations/logging-observability)
- [HTTP client](/docs/{{version}}/services/http-client)
- [Deployment](/docs/{{version}}/operations/deployment)
- [Queries](/docs/{{version}}/database/queries)
