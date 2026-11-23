# Migrations and schema

- [A migration is registered Go code with immutable identity](#a-migration-is-registered-go-code-with-immutable-identity)
- [Design for the mixed-version rollout window](#design-for-the-mixed-version-rollout-window)
- [Register one reversible schema change](#register-one-reversible-schema-change)
- [Inspect, preview, apply, and roll back deliberately](#inspect-preview-apply-and-roll-back-deliberately)
- [Understand skip, transaction, and rollback outcomes](#understand-skip-transaction-and-rollback-outcomes)
- [Separate deployment authority from request authority](#separate-deployment-authority-from-request-authority)
- [Test both schema directions and the old binary](#test-both-schema-directions-and-the-old-binary)
- [Restore service before choosing a schema-down action](#restore-service-before-choosing-a-schema-down-action)
- [Coordinate schema, deployment, and seed data](#coordinate-schema-deployment-and-seed-data)

> [!WARNING]
> **Current limitations.** Migrations and the registry are released as an explicit pipeline path, but generated identity and application ordering still require review. Inspect the exact plan and rollback behavior for the deployed dialect, and never run schema changes during request-process boot.

<a name="a-migration-is-registered-go-code-with-immutable-identity"></a>
## A migration is registered Go code with immutable identity

Each migration implements `migrations.Migration`, registers from init, and is sorted by `GetName`. The repository records that exact name and batch only after Up succeeds. The starter uses `arandu_migrations` and imports `database/migrations` from bootstrap so every command sees the registry.

<a name="design-for-the-mixed-version-rollout-window"></a>
## Design for the mixed-version rollout window

State how the previous binary behaves after the new schema lands. Add columns as nullable or defaulted, add tables old code ignores, and split rename or removal across releases. Production commands require a shared RESP cache for the migration lock.

- Choose an immutable migration name.
- Plan old-binary compatibility.
- Test the selected database grammar.
- Configure a cross-replica lock in production.

<a name="register-one-reversible-schema-change"></a>
## Register one reversible schema change

Implement Up and Down against the migration connection, register the value, and assert `ReversibleMigration` at compile time. Embed `BaseMigration` for the default connection, normal `ShouldRun` behavior, and the released transactional intent.

```go example
type AddLocaleToUsers struct{ migrations.BaseMigration }

func (AddLocaleToUsers) GetName() string { return "20260830_0004_add_locale_to_users" }
func (AddLocaleToUsers) Up(ctx context.Context, conn migrations.Connection) error {
	return conn.Schema().Table(ctx, "users", func(t *schema.Blueprint) { t.String("locale").Nullable() })
}
func (AddLocaleToUsers) Down(ctx context.Context, conn migrations.Connection) error {
	return conn.Schema().Table(ctx, "users", func(t *schema.Blueprint) { t.DropColumn("locale") })
}
func init() { migrations.Register(AddLocaleToUsers{}) }

var _ migrations.ReversibleMigration = AddLocaleToUsers{}
```

<a name="inspect-preview-apply-and-roll-back-deliberately"></a>
## Inspect, preview, apply, and roll back deliberately

Generate a migration, inspect status and pretend output, apply normally or in steps, and roll back only a reviewed last batch or explicit step. Use the schema builder for portable DDL and explicit Statement or Select for a reviewed backfill. Time columns default to microsecond precision, so two rows written within one second keep their order; `schema.DefaultTimePrecision` changes that default, and a `migrations.WidenTimePrecision` migration raises the columns of a table created at whole seconds, which is how the packages' own tables catch up.

- `aru make:migration` with `--create` or `--table`.
- `aru migrate:status --pending` and `aru migrate --pretend`.
- `aru migrate` or `aru migrate --step`.
- `aru migrate:rollback` with an explicit batch decision.

<a name="understand-skip-transaction-and-rollback-outcomes"></a>
## Understand skip, transaction, and rollback outcomes

ShouldRun=false is skipped and not recorded, so it is reconsidered later. Transactional DDL varies by grammar and `OutsideTransaction` may leave partial schema. Rollback preflights the batch; an undeclared missing Down refuses before work, while an explicit Irreversible migration remains applied and recorded.

> [!WARNING]
> Hesape :hesapeVersion does not report an irreversible migration as reverted and does not delete its repository row. A runtime Down failure can still stop after earlier reversible entries already completed.

<a name="separate-deployment-authority-from-request-authority"></a>
## Separate deployment authority from request authority

> [!WARNING]
> Migrations carry no Grant because they are deployment DDL, not requests. Give the pipeline a narrow credential, bind data values, review raw SQL, keep secrets out of progress output, and never run migration from an HTTP handler or replica startup.

<a name="test-both-schema-directions-and-the-old-binary"></a>
## Test both schema directions and the old binary

On each supported disposable database, migrate twice, inspect batch status, run Down then Up, inspect pretend statements, inject Up and Down failures, test irreversible preflight, exercise lock winner and outage, and boot the previous binary against additive schema.

- Repeated migrate is idempotent.
- Down and Up restore the expected catalogue.
- Lock contention applies once.
- The previous binary still reads and writes.

<a name="restore-service-before-choosing-a-schema-down-action"></a>
## Restore service before choosing a schema-down action

Redeploy the previous binary when additive schema is compatible. For a bad reversible batch, inspect status and use the narrow tested rollback. For partial outside-transaction DDL or irreversible change, inspect the real catalogue and backup and write a forward repair instead of editing repository rows.

`migrate:fresh`, `migrate:reset`, and `migrate:refresh` run only with `APP_ENV=dev`, and a missing `APP_ENV` is production, because they empty the schema; rollback remains an explicit batch decision.

<a name="coordinate-schema-deployment-and-seed-data"></a>
## Coordinate schema, deployment, and seed data

Deployment owns the rollout window, upgrades choose releases, seeders populate deliberate records, and database tests exercise both directions.

- [Deployment](/docs/{{version}}/operations/deployment)
- [Upgrades](/docs/{{version}}/operations/upgrades)
- [Seeders](/docs/{{version}}/database/seeders)
- [Database testing](/docs/{{version}}/testing/database)
