# Pruebas de base de datos

- [Las pruebas de base demuestran transacciones, dialectos y ámbito del tenant](#database-tests-prove-transactions-dialects-and-tenant-scope)
- [Cree un esquema desechable y un tenant explícito](#create-disposable-schema-and-one-explicit-tenant)
- [Ejecute el contrato del repositorio con un Grant real](#run-the-repository-contract-with-a-real-grant)
- [Ejercite las fronteras de transacción, migración y consulta](#exercise-transaction-migration-and-query-boundaries)
- [No deje que el estado compartido haga que la suite dependa del orden](#do-not-let-shared-state-make-a-green-suite-order-dependent)
- [Cada fixture y aserción permanece dentro del ámbito del tenant](#every-fixture-and-assertion-stays-tenant-scoped)
- [Ejecute el contrato en cada conector compatible](#run-the-contract-against-every-supported-connector)
- [Use factories para la estructura y migraciones para el estado del catálogo](#use-factories-for-shape-and-migrations-for-catalogue-state)

> [!WARNING]
> **Limitaciones actuales.** Los helpers de base, models y aserciones de tenant están publicados. La persistencia de factories y la identidad de migraciones aún tienen límites; cree en las pruebas solo las filas necesarias, por el mismo recorrido de datos verificado de la aplicación.

<a name="database-tests-prove-transactions-dialects-and-tenant-scope"></a>
## Las pruebas de base demuestran transacciones, dialectos y ámbito del tenant

Una prueba de base usa el conector, la gramática de consulta, el model, la transacción y la migración reales elegidos por la aplicación. Los fakes en memoria son contratos útiles, pero no demuestran placeholders SQL, comportamiento DDL, bloqueos, cancelación ni predicados de tenant.

<a name="create-disposable-schema-and-one-explicit-tenant"></a>
## Cree un esquema desechable y un tenant explícito

Dé a cada prueba una base propia y ejecute las migraciones mediante el conector publicado: `test.Serve` levanta la aplicación sobre una copia de una plantilla Postgres ya migrada o sobre un archivo SQLite temporal, y las opciones Migrate, Seed, SQLite y Database indican cómo se arma el esquema y en qué servidor se crea la base. El aislamiento proviene del almacenamiento exclusivo de la prueba; nunca reutilice credenciales de producción ni un esquema de desarrollo compartido.

- Aísle el esquema por prueba o paquete.
- Ejecute migraciones reales.
- Use semillas deterministas en las factories.
- Limpie con rollback transaccional o descarte de los recursos propios.

<a name="run-the-repository-contract-with-a-real-grant"></a>
## Ejecute el contrato del repositorio con un Grant real

Ejercite el repositorio protegido con un Grant real e inspeccione la base de forma independiente con las aserciones de `framework/test`. Incluya `tenant_id` en cada comparación de filas del cliente para que otro tenant no haga pasar la aserción por accidente.

```go example
func assertInvoiceStored(ctx context.Context, t *testing.T, db *database.DB, invoiceID string) {
	test.AssertDatabaseHas(t, ctx, db, "invoices", test.Match{
		"id":        invoiceID,
		"tenant_id": "acme",
	})
}
```

<a name="exercise-transaction-migration-and-query-boundaries"></a>
## Ejercite las fronteras de transacción, migración y consulta

Pruebe commit y rollback, anidamiento donde se admita, cancelación, conflictos, ausencia, orden, paginación, bloqueos, migraciones Up y Down, escrituras parciales de factories e inserción atómica de cola u outbox cuando compartan SQL.

- Use fixtures explícitos de dialecto SQL.
- Verifique filas y predicados de tenant.
- Pruebe la cancelación del contexto durante I/O.
- Ejecute el rollback transaccional como un recorrido real de fallo.

<a name="do-not-let-shared-state-make-a-green-suite-order-dependent"></a>
## No deje que el estado compartido haga que la suite dependa del orden

Un esquema, reloj, IDs o registros globales compartidos pueden hacer que las pruebas pasen solo en cierto orden. Detecte filas filtradas, nombres fijos de migración, conexiones retenidas, colisiones de puertos en paralelo y supuestos de que SQLite se comporta como cada dialecto de producción.

> [!WARNING]
> Una expectativa de mock sobre una llamada SQL no demuestra que la consulta devuelva solo filas autorizadas. Mantenga al menos un contrato con conector concreto por cada despliegue compatible.

<a name="every-fixture-and-assertion-stays-tenant-scoped"></a>
## Cada fixture y aserción permanece dentro del ámbito del tenant

> [!WARNING]
> Use datos sintéticos de tenant, nunca vuelque filas de otros tenants al fallar, vincule valores en vez de interpolar SQL e impida que Grants de prueba lleguen a recorridos de petición de la aplicación.

`security.SystemGrant` recibe acción y tenant explícitos. Registre el motivo de la prueba en su nombre o comentario y mantenga el Doctor estricto rechazando usos no autorizados en producción.

<a name="run-the-contract-against-every-supported-connector"></a>
## Ejecute el contrato en cada conector compatible

Ejecute éxito, ausencia, duplicado, Grant vacío, otro tenant, cancelación, rollback, repetición de migración, Down y Up y casos paralelos sin race contra cada conector compatible en CI.

- Cada conector satisface el mismo contrato de repositorio.
- Los datos de otro tenant son invisibles.
- El rollback no deja filas.
- Las migraciones son repetibles y reversibles donde se declara.

<a name="use-factories-for-shape-and-migrations-for-catalogue-state"></a>
## Use factories para la estructura y migraciones para el estado del catálogo

Los models definen el mapeo, las factories generan filas tipadas, las migraciones poseen el esquema y los contratos permiten ejecutar la misma suite de comportamiento en cada adapter.

- [Models](/docs/{{version}}/database/models)
- [Factories](/docs/{{version}}/database/factories)
- [Migraciones](/docs/{{version}}/database/migrations)
- [Contratos](/docs/{{version}}/architecture/contracts-adapters)
