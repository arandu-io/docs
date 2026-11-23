# Contributing

- [Repository gates define a reviewable change](#repository-gates-define-a-reviewable-change)
- [Read the repository contract before editing](#read-the-repository-contract-before-editing)
- [Make one focused change and its regression](#make-one-focused-change-and-its-regression)
- [Prove every new guard can fail](#prove-every-new-guard-can-fail)
- [Keep one change attributable to one reason](#keep-one-change-attributable-to-one-reason)

> [!WARNING]
> **Current limitations.** The contribution contract is verified by the tagged repository's own CI and layout guard. It has no independent tagged consumer because it governs changes to that repository rather than a runtime package.

<a name="repository-gates-define-a-reviewable-change"></a>
## Repository gates define a reviewable change

Formatting, vet, and race-enabled tests are the common floor, with repository-specific generation, dependency, vulnerability, and conformance checks added beside them. The commands run locally and in CI; a contribution reports their actual results rather than treating the command names as approval.

```bash
unformatted=$(find . -name '*.go' -not -path '*/testdata/*' -not -name '*.kyse.go' -exec gofmt -l {} +)
test -z "$unformatted" || { echo "$unformatted"; exit 1; }
GOWORK=off go vet ./...
GOWORK=off go test -race -count=1 -p 1 ./...
```

<a name="read-the-repository-contract-before-editing"></a>
## Read the repository contract before editing

Use the Go version and module workspace policy the repository declares, inspect generated-file ownership, and identify the focused tests for the behavior being changed. Security reports follow the private reporting policy; a public issue is not the first place to disclose a vulnerability.

<a name="make-one-focused-change-and-its-regression"></a>
## Make one focused change and its regression

Reproduce the failure, add a test that observes the broken behavior, implement the smallest coherent correction, and run that test again. Update public documentation when the observable API changes. Keep unrelated formatting or generated churn out of the patch so the behavioral evidence remains reviewable.

<a name="prove-every-new-guard-can-fail"></a>
## Prove every new guard can fail

Plant the violation, run the focused guard, observe its failure, restore the source, and observe success. A scanning test also asserts it inspected a plausible number of targets so a stale pattern cannot pass by matching nothing. Keep generated view sources and intentionally invalid fixtures out of generic formatting checks for their documented reasons.

<a name="keep-one-change-attributable-to-one-reason"></a>
## Keep one change attributable to one reason

Separate mechanical rewrites from behavioral changes when that keeps review evidence legible. A commit message explains the reason the diff cannot reconstruct. Repository-local instructions, generated-file ownership, release rules, and security reporting procedures still apply in addition to the common commands.
