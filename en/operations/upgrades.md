# Upgrade guide

- [Upgrade the module graph and application-owned wiring](#upgrade-the-module-graph-and-application-owned-wiring)
- [Record the current artifact and schema state](#record-the-current-artifact-and-schema-state)
- [Move one compatible release set on a branch](#move-one-compatible-release-set-on-a-branch)
- [Review API, configuration, schema, and generated output](#review-api-configuration-schema-and-generated-output)
- [Separate compile failures from silent behavior changes](#separate-compile-failures-from-silent-behavior-changes)
- [Recheck every identity and data boundary](#recheck-every-identity-and-data-boundary)
- [Run behavior gates against the target release set](#run-behavior-gates-against-the-target-release-set)
- [Restore code first and undo schema only with evidence](#restore-code-first-and-undo-schema-only-with-evidence)
- [Read release evidence before changing production](#read-release-evidence-before-changing-production)

> [!WARNING]
> **Current limitations.** Module upgrades are explicit and rollback-compatible schema changes are required. Migration identity and irreversible rollback reporting still need care; inspect the exact changelog and exercise mixed-version rollout before production.

<a name="upgrade-the-module-graph-and-application-owned-wiring"></a>
## Upgrade the module graph and application-owned wiring

An Arandu application is ordinary Go modules plus application code created from a starter. Updating `go.mod` changes library APIs; it does not rewrite bootstrap, configuration, Dockerfile, migrations, or published UI files owned by the project. Review tagged changelogs and the starter diff together.

<a name="record-the-current-artifact-and-schema-state"></a>
## Record the current artifact and schema state

Start from a green commit with restorable module files, a database backup appropriate to the driver, the deployed image digest, and the migration batch recorded by `arandu_migrations`. Read every tag between installed and target versions, never the moving development branch.

- Inventory `go list -m` all and `aru version`.
- Record the starter and image versions.
- Identify generated Kyse output and application-owned adaptations.
- Plan the mixed-version schema window.

<a name="move-one-compatible-release-set-on-a-branch"></a>
## Move one compatible release set on a branch

Pin one reviewed release set, tidy the graph, and regenerate views. Then follow the skeleton's `UPGRADE.md` entry by entry from the release the project was created from: it names what a project changes in `bootstrap/app.go`, app/Providers, configuration, the environment, and tests. Each module's own `UPGRADE.md`, in framework, hesape, kyse, ui, joaju, and aru, names what stops compiling or behaves differently and what to write instead. Port only intentional changes by review, because `aru new` is not an upgrade command.

```sh
go get github.com/arandu-io/framework@:frameworkVersion \
  github.com/arandu-io/hesape@:hesapeVersion \
  github.com/arandu-io/kyse@:kyseVersion \
  github.com/arandu-io/hesape/database/connectors/pgx@:pgxVersion \
  github.com/arandu-io/hesape/database/connectors/sqlite@:sqliteVersion \
  github.com/arandu-io/hesape/redis@:redisVersion
go mod tidy
aru view:build
```

<a name="review-api-configuration-schema-and-generated-output"></a>
## Review API, configuration, schema, and generated output

Compile after small module groups, search removed symbols, compare every environment name at parser and consumer, regenerate views after Kyse changes, preserve custom UI blocks, stage additive schema changes, and roll web and worker from the same artifact.

- Review API symbols and configuration wiring separately.
- Exercise the cache, session, and queue drivers selected in production.
- Regenerate and compile view output.
- Test the previous binary against the additive target schema.

<a name="separate-compile-failures-from-silent-behavior-changes"></a>
## Separate compile failures from silent behavior changes

Removed symbols fail at compile time, while a renamed environment variable, changed default, dead parsed field, or middleware reordering can compile and still break production. Generated imports fail when view build is skipped, and destructive schema can break the old binary.

> [!WARNING]
> Do not fix an upgrade by adding another local environment reader without tracing the value into the runtime consumer. Framework-owned component parsers and application-owned readers can both parse a variable that neither live constructor uses.

<a name="recheck-every-identity-and-data-boundary"></a>
## Recheck every identity and data boundary

> [!WARNING]
> Recheck session backends and cookies, CSRF order, Grant propagation, token hashing, OAuth state, diagnostic exposure, TLS roots, and rate-limit store selection. A dependency update is not proof that application wiring adopted the safer path.

<a name="run-behavior-gates-against-the-target-release-set"></a>
## Run behavior gates against the target release set

Run view generation, formatting, vet, race tests, build, doctor, production-shaped container tests, session across two replicas, migration twice and rollback once, previous-binary schema compatibility, queue option wiring, and exact tagged proof checks.

- GOWORK=off `go test -race ./`... passes.
- `aru doctor --strict` passes.
- Generated views match source.
- Production drivers pass their feature journeys.

<a name="restore-code-first-and-undo-schema-only-with-evidence"></a>
## Restore code first and undo schema only with evidence

Redeploy the previous image or restore module files and application wiring. Leave an additive schema in place during binary rollback. Undo a release migration only after inspecting its batch and tested Down; never use fresh, reset, or refresh as rollback.

<a name="read-release-evidence-before-changing-production"></a>
## Read release evidence before changing production

Release notes define changes, migrations define schema compatibility, deployment defines artifact rollout, and configuration exposes silent wiring risks.

- [Release notes](/docs/{{version}}/operations/release-notes)
- [Migrations](/docs/{{version}}/database/migrations)
- [Deployment](/docs/{{version}}/operations/deployment)
- [Configuration](/docs/{{version}}/architecture/configuration)
- [Skeleton upgrade guide](https://github.com/arandu-io/arandu/blob/:aranduVersion/UPGRADE.md)
