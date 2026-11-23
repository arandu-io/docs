# Testing overview

- [Tests boot the same application graph without serving a port](#tests-boot-the-same-application-graph-without-serving-a-port)
- [Generate views and isolate external state](#generate-views-and-isolate-external-state)
- [Run the released project gates](#run-the-released-project-gates)
- [Prove behavior at the narrowest complete boundary](#prove-behavior-at-the-narrowest-complete-boundary)
- [Choose HTTP, console, database, or contract detail next](#choose-http-console-database-or-contract-detail-next)

> [!WARNING]
> **Current limitations.** Testing helpers are released across several modules rather than one unified facade. Compose ordinary Go tests from the helper that owns the boundary, and keep package-specific setup visible in the test.

<a name="tests-boot-the-same-application-graph-without-serving-a-port"></a>
## Tests boot the same application graph without serving a port

Arandu uses ordinary Go tests plus application helpers. A test may construct one service, boot the composed application and dispatch its Handler, or run a command adapter. The valuable boundary is the smallest one that still proves the behavior, not a separate testing framework.

<a name="generate-views-and-isolate-external-state"></a>
## Generate views and isolate external state

Run `aru view:build` on a clean checkout because generated view packages may be gitignored. Use disposable database and shared-store fixtures, explicit fake transports, controlled clocks where offered, and environment restoration through `t.Setenv`.

- Keep GOWORK=off release checks reproducible.
- Isolate database and tenant per test.
- Prevent stray network requests.
- Use `t.Cleanup` for every process-owned resource.

<a name="run-the-released-project-gates"></a>
## Run the released project gates

Run format, view generation, vet, race tests, build, and strict Doctor from the project root. A green parser or unit test is not enough when the behavior depends on bootstrap, a selected driver, or generated output. `aru test` runs `go test`, and `aru test --e2e` runs the journeys under tests/E2E, which sit behind the e2e build tag, drive a real browser against the application on a port of their own, and skip on a developer's machine where no browser is found.

```sh
aru view:build
unformatted=$(find . -name '*.go' -not -path '*/testdata/*' -not -name '*.kyse.go' -exec gofmt -l {} +)
test -z "$unformatted" || { echo "$unformatted"; exit 1; }
GOWORK=off go vet ./...
GOWORK=off go test -race -count=1 ./...
GOWORK=off go build ./...
aru doctor --strict
```

<a name="prove-behavior-at-the-narrowest-complete-boundary"></a>
## Prove behavior at the narrowest complete boundary

Use unit tests for pure policy and transformation, contract tests for every adapter, feature tests through the real application graph, and production-shaped journeys for cross-replica state. Assert outcomes and refusal paths rather than only that values parse or constructors return.

- A changed configuration reaches observable behavior.
- Every network adapter has a no-stray-request test.
- Tenant isolation is exercised at storage.
- Race tests cover shared state and shutdown.

<a name="choose-http-console-database-or-contract-detail-next"></a>
## Choose HTTP, console, database, or contract detail next

HTTP tests dispatch routes, console tests exercise delegated commands, database tests prove transactions and Grants, and test doubles capture component boundaries without global state.

- [HTTP testing](/docs/{{version}}/testing/http)
- [Console testing](/docs/{{version}}/testing/console)
- [Database testing](/docs/{{version}}/testing/database)
- [Test doubles](/docs/{{version}}/testing/test-doubles)
