# Database connections

- [One connection path with instrumentation](#one-connection-path-with-instrumentation)
- [Link one driver and open the configured URL](#link-one-driver-and-open-the-configured-url)
- [Configure pools and transactions at the boundary](#configure-pools-and-transactions-at-the-boundary)
- [Fail boot on an absent driver or invalid connection](#fail-boot-on-an-absent-driver-or-invalid-connection)
- [Run the connector contract against real engines](#run-the-connector-contract-against-real-engines)
- [Configuration, DB, and Connection](#configuration-db-and-connection)
- [Repositories are the application door](#repositories-are-the-application-door)

<a name="one-connection-path-with-instrumentation"></a>
## One connection path with instrumentation

`database.Open` builds the application's DB from typed configuration. The handle knows its dialect, rebinds placeholders, applies the pool policy, and reports statements to the configured collector. Keeping one wrapped handle prevents a second path where queries run without the diagnostics the rest of the application relies on.

<a name="link-one-driver-and-open-the-configured-url"></a>
## Link one driver and open the configured URL

A connector is a separate Go module and registers itself during package initialization. Blank-import exactly the engines this binary is allowed to speak. In an application `foundation.NewDrivers` makes this call with the configuration `DATABASE_URL` and the pool settings produced, answers the handle as `drivers.DB`, and closes the pool at shutdown and at the end of every command; repositories and handlers must not call another environment reader.

```go example
import (
	"github.com/arandu-io/hesape/database"
	_ "github.com/arandu-io/hesape/database/connectors/sqlite"
)

func OpenDatabase(cfg database.Config) (*database.DB, func(), error) {
	db, closeDB, err := database.Open(cfg)
	if err != nil {
		return nil, nil, err
	}
	return db, closeDB, nil
}
```

<a name="configure-pools-and-transactions-at-the-boundary"></a>
## Configure pools and transactions at the boundary

The URL chooses the registered dialect and carries driver connection parameters. Pool limits and lifetime belong to database configuration, not individual repositories. `database.Transaction` keeps all work in the context-bound transaction; pass that context through every call so a nested operation does not escape to the base pool.

<a name="fail-boot-on-an-absent-driver-or-invalid-connection"></a>
## Fail boot on an absent driver or invalid connection

A URL naming a driver the binary did not link cannot be repaired at runtime and should stop startup with the returned error. Treat an unreachable server, invalid credentials, a missing SQLite file, and malformed configuration as distinct boot failures. Do not fall back to another database whose data has different ownership.

<a name="run-the-connector-contract-against-real-engines"></a>
## Run the connector contract against real engines

The connector conformance suite is what proves the shared contract for `PostgreSQL`, `MySQL`, and SQLite. At application level, assert the selected driver is available, a parameterized statement reaches the collector, transaction rollback leaves no row, and shutdown releases the pool. A mock alone cannot prove dialect binding.

<a name="configuration-db-and-connection"></a>
## Configuration, DB, and Connection

The configuration layer and `ParseURL` produce Config; Open creates DB; Transaction scopes work; `AvailableDrivers` shows registered connector names. An application should choose one configuration owner instead of calling `database.Load` again inside its database constructor. DB and Connection are infrastructure below repository enforcement. `ForMigrations` and `ForSchema` expose narrow interfaces for those subsystems instead of turning application repositories into schema tools.

<a name="repositories-are-the-application-door"></a>
## Repositories are the application door

Connection plumbing cannot apply tenant policy to arbitrary SQL, so application modules reach rows through a repository or model terminal that requires a Grant. Migrations deliberately run below that boundary because they operate on schema before any request exists. Connector selection belongs in bootstrap; no handler should import a driver or open another pool.
