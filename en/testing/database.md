# Database testing

- [Database tests prove transactions, dialects, and tenant scope](#database-tests-prove-transactions-dialects-and-tenant-scope)
- [Create disposable schema and one explicit tenant](#create-disposable-schema-and-one-explicit-tenant)
- [Run the repository contract with a real Grant](#run-the-repository-contract-with-a-real-grant)
- [Exercise transaction, migration, and query boundaries](#exercise-transaction-migration-and-query-boundaries)
- [Do not let shared state make a green suite order-dependent](#do-not-let-shared-state-make-a-green-suite-order-dependent)
- [Every fixture and assertion stays tenant-scoped](#every-fixture-and-assertion-stays-tenant-scoped)
- [Run the contract against every supported connector](#run-the-contract-against-every-supported-connector)
- [Use factories for shape and migrations for catalogue state](#use-factories-for-shape-and-migrations-for-catalogue-state)

> [!WARNING]
> **Current limitations.** Database helpers, models, and tenant assertions are released. Factory persistence and migration identity still have limits, so tests should create only the rows they need through the same verified data path as the application.

<a name="database-tests-prove-transactions-dialects-and-tenant-scope"></a>
## Database tests prove transactions, dialects, and tenant scope

A database test uses the actual connector, query grammar, model, transaction, and migration code selected by the application. Memory fakes are useful contracts, but they cannot prove SQL placeholders, DDL behavior, lock semantics, cancellation, or tenant predicates.

<a name="create-disposable-schema-and-one-explicit-tenant"></a>
## Create disposable schema and one explicit tenant

Give each test a database of its own and migrate it through the released connector: `test.Serve` boots the application over a copy of a migrated Postgres template or a temporary SQLite file, and its Migrate, Seed, SQLite, and Database options say how the schema is built and which server the database is made on. Isolation comes from storage owned by that test; never reuse production credentials or a shared developer schema.

- Isolate schema per test or package.
- Run real migrations.
- Use deterministic factory seeds.
- Clean with transaction rollback or owned disposal.

<a name="run-the-repository-contract-with-a-real-grant"></a>
## Run the repository contract with a real Grant

Exercise the protected repository with a real Grant, then inspect the database independently with the assertions of `framework/test`. Include `tenant_id` in every customer-owned row match so another tenant cannot make the assertion pass accidentally.

```go example
func assertInvoiceStored(ctx context.Context, t *testing.T, db *database.DB, invoiceID string) {
	test.AssertDatabaseHas(t, ctx, db, "invoices", test.Match{
		"id":        invoiceID,
		"tenant_id": "acme",
	})
}
```

<a name="exercise-transaction-migration-and-query-boundaries"></a>
## Exercise transaction, migration, and query boundaries

Test commit and rollback, nested behavior where supported, cancellation, conflicts, absence, ordering, pagination, locking, migrations Up and Down, factory partial writes, and queue or outbox atomic insertion when those share SQL.

- Use explicit SQL dialect fixtures.
- Assert rows and tenant predicates.
- Test context cancellation during I/O.
- Run transaction rollback as a real failure path.

<a name="do-not-let-shared-state-make-a-green-suite-order-dependent"></a>
## Do not let shared state make a green suite order-dependent

Shared schema, clock, ids, or global registrations can make tests pass only in one order. Detect leaked rows, fixed migration names, held connections, parallel port collisions, and assumptions that SQLite behavior equals every production dialect.

> [!WARNING]
> A mock expectation that SQL was called does not prove the query returned only authorized rows. Keep at least one concrete connector contract for each supported deployment path.

<a name="every-fixture-and-assertion-stays-tenant-scoped"></a>
## Every fixture and assertion stays tenant-scoped

> [!WARNING]
> Use synthetic tenant data, never dump cross-tenant rows on failure, bind values rather than interpolating SQL, and ensure test helper Grants cannot leak into application request paths.

`security.SystemGrant` takes an explicit action and tenant. Keep the test's reason in its name or comment, and let strict Doctor continue to reject unauthorized production use.

<a name="run-the-contract-against-every-supported-connector"></a>
## Run the contract against every supported connector

Run success, absence, duplicate, zero Grant, foreign tenant, cancellation, rollback, migration repeat, Down and Up, and race-safe parallel cases against each supported connector in CI.

- Every connector satisfies the same repository contract.
- Cross-tenant data is invisible.
- Rollback leaves no row.
- Migrations are repeatable and reversible where declared.

<a name="use-factories-for-shape-and-migrations-for-catalogue-state"></a>
## Use factories for shape and migrations for catalogue state

Models define mapping, factories generate typed rows, migrations own schema, and contracts let each adapter run the same behavioral suite.

- [Models](/docs/{{version}}/database/models)
- [Factories](/docs/{{version}}/database/factories)
- [Migrations](/docs/{{version}}/database/migrations)
- [Contracts](/docs/{{version}}/architecture/contracts-adapters)
