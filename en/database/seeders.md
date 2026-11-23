# Seeders

- [Seeders are explicit named application operations](#seeders-are-explicit-named-application-operations)
- [Define narrow dependencies and idempotent intent](#define-narrow-dependencies-and-idempotent-intent)
- [Run the root or one named seeder](#run-the-root-or-one-named-seeder)
- [Create bootstrap accounts through application services](#create-bootstrap-accounts-through-application-services)
- [Stop at the first incomplete operation](#stop-at-the-first-incomplete-operation)
- [Keep secrets out of source, history, and process lists](#keep-secrets-out-of-source-history-and-process-lists)
- [Run every seeder twice on a non-empty database](#run-every-seeder-twice-on-a-non-empty-database)
- [Use factories for data shape and services for policy](#use-factories-for-data-shape-and-services-for-policy)

> [!WARNING]
> **Current limitations.** Seeder registration and execution are released, but factory persistence is not a uniformly verified path. Seed through explicit model or repository calls and make repeated execution safe for the environment where it runs.

<a name="seeders-are-explicit-named-application-operations"></a>
## Seeders are explicit named application operations

A seeder implements `database.Seeder[Deps]` with Name and `Run(context.Context, Deps)`. The project owns Deps, the registry, the fallback root, and execution order, so a seeder reaches only services the provider supplies. The provider's Seed, handed to the builder with `WithSeeders(app.Seed)` in `bootstrap/app.go`, is what `db:seed` and migrate `--seed` run. There is no reflection or filesystem discovery.

<a name="define-narrow-dependencies-and-idempotent-intent"></a>
## Define narrow dependencies and idempotent intent

Declare narrow dependencies containing application services, configured tenant, and optional remaining arguments. Register each seeder and select `DatabaseSeeder` as the no-name fallback. Make a second run converge rather than duplicate rows or silently replace credentials.

- Register every named seeder.
- Choose one fallback root.
- Pass a configured tenant explicitly.
- Make repeated execution observable and idempotent.

<a name="run-the-root-or-one-named-seeder"></a>
## Run the root or one named seeder

`aru db:seed` runs the fallback; a positional name selects one case-insensitively. Unknown names fail with the available registry, and the obsolete `--class` spelling is refused. Remaining arguments belong to that seeder through `database.Flag` and `database.Switch`.

```sh
aru db:seed
aru db:seed UserSeeder -e you@example.com -p '<choose-a-password>'
aru db:seed UserSeeder -e you@example.com -p '<choose-a-password>' -r admin
```

<a name="create-bootstrap-accounts-through-application-services"></a>
## Create bootstrap accounts through application services

The starter root is empty by default: migrating or running `aru db:seed` never creates an account. When an application needs a bootstrap or recovery user, run `UserSeeder` by name with an explicit email and password; add `-r` admin only when that role is intended. Application policy and hashing stay on the normal service path.

- Keep `DatabaseSeeder` empty until the application deliberately adds data.
- Run `UserSeeder` by name to create an account.
- Add `-r` admin only as an explicit authorization decision.
- Leave existing roles unchanged during a password-only recovery.

<a name="stop-at-the-first-incomplete-operation"></a>
## Stop at the first incomplete operation

`database.Seed` wraps a Run failure with the seeder name. Missing registry, dependency, tenant, or explicit `UserSeeder` credential fails before writes, and a root stops rather than skipping over an incomplete prerequisite. Production asks for confirmation unless explicitly forced.

> [!WARNING]
> The framework does not wrap every child seeder in one transaction. Use an application transaction when all-or-nothing behavior is required, or make each step independently idempotent and visible.

<a name="keep-secrets-out-of-source-history-and-process-lists"></a>
## Keep secrets out of source, history, and process lists

> [!WARNING]
> Never create a default account, commit a password, accept a caller-selected tenant, or expose seeding through HTTP. An explicit `UserSeeder` password appears in shell history and process lists; clear that history entry, change the password through the application immediately, and confine bootstrap authority inside an audited application service.

A recovery password supplied on an operator command line must be changed through the application immediately and is not an ordinary deployment seed path.

<a name="run-every-seeder-twice-on-a-non-empty-database"></a>
## Run every seeder twice on a non-empty database

Migrate a non-empty disposable database, run root and each named seeder twice, assert exact tenant and roles without duplicates, exercise unknown names and obsolete flags, fail one child, and prove production confirmation defaults to no.

- Repeated root execution converges.
- Unknown names list valid choices.
- A child failure stops later children.
- Missing secrets leave no partial account.

<a name="use-factories-for-data-shape-and-services-for-policy"></a>
## Use factories for data shape and services for policy

Factories build typed data, migrations own schema, policies carry authority, and deployment decides when a bootstrap operation may run.

- [Factories](/docs/{{version}}/database/factories)
- [Migrations](/docs/{{version}}/database/migrations)
- [Authorization](/docs/{{version}}/security/authorization)
- [Deployment](/docs/{{version}}/operations/deployment)
