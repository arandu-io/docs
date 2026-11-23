# Constructor fluido de consultas

- [Componga antes de ejecutar](#compose-before-execution)
- [Consulte facturas abiertas bajo un Grant](#query-open-invoices-under-a-grant)
- [Filtre y ordene con listas de valores permitidos](#filter-and-order-with-allowlists)
- [Elija el contrato de registro ausente](#choose-the-missing-row-contract)
- [Verifique ámbito de tenant y formato de SQL](#verify-tenant-scope-and-sql-shape)
- [Familias de builders y operaciones terminales](#builder-and-terminal-families)
- [La autorización precede a la consulta](#authorization-precedes-the-query)

<a name="compose-before-execution"></a>
## Componga antes de ejecutar

Una consulta empieza por el constructor que `aru model:build` genera junto a cada entidad, como `models.Invoices(db)`, y la composición queda separada de la ejecución. Where, `OrderByDesc` y Limit describen la sentencia sin enviarla. Get es una operación terminal: recibe el contexto de la solicitud y un `security.Grant`, aplica la columna del tenant de la tabla, compila el SQL del dialecto, ejecuta y devuelve la colección generada, un slice de punteros de la entidad.

Mantener el Grant en la operación terminal es esencial. Construir una consulta reutilizable no autoriza nada; la sentencia que lee los registros debe llevar la autoridad y el tenant de quien la ejecuta. Un Grant cero se rechaza con `model.ErrNoTenant` antes de que el SQL llegue a la conexión.

<a name="query-open-invoices-under-a-grant"></a>
## Consulte facturas abiertas bajo un Grant

Invoice incorpora `model.Model` e `invoiceTable` declara su tabla; ambos van en `app/Models/Invoice.go`. `aru model:build` escribe Invoices, `InvoiceQuery` e `InvoiceCollection` junto a ellos en `InvoiceQuery.go`, un archivo que se regenera y no se edita. `OpenInvoices` recibe base instrumentada y Grant ya emitido para `invoice.list`. La columna de tenant predeterminada es `tenant_id`: el terminal añade el tenant fiable independientemente del filtro de estado.

El servicio obtiene el Grant mediante `security.Authorize`, con su Policy, el Subject autenticado, la acción `invoiceList` y el recurso de la colección. Pasar ese Grant a `OpenInvoices` conserva el orden de la aplicación de ejemplos publicada: primero autorizar, después ejecutar la operación de datos.

```go
type Invoice struct {
	model.Model

	ID        int64     `db:"id"`
	TenantID  string    `db:"tenant_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}

var invoiceTable = model.NewTable(model.TableSpec{
	Name: "invoices",
	New:  func() model.Entity { return new(Invoice) },
})

const invoiceList security.Action = "invoice.list"

func OpenInvoices(ctx context.Context, db *database.DB, grant security.Grant) (models.InvoiceCollection, error) {
	return models.Invoices(db).
		Where("status", "=", "open").
		OrderByDesc("created_at").
		OrderByDesc("id").
		Limit(25).
		Get(ctx, grant)
}
```

<a name="filter-and-order-with-allowlists"></a>
## Filtre y ordene con listas de valores permitidos

Los valores pasados a Where son parámetros vinculados. Los nombres de columnas, las direcciones de ordenación, las expresiones seleccionadas y los fragmentos sin procesar forman parte de la estructura SQL y no se protegen con un marcador de valor. Traduzca la elección de la solicitud mediante un mapa fijo, como `map[string]``string{"newest": "created_at", "total": "total"}`; rechace claves ausentes antes de llamar a `OrderBy` o Select.

Componga los predicados adicionales antes de la operación terminal y mantenga los valores separados de los identificadores. Use un orden secundario estable cuando dos registros puedan compartir el valor de la columna solicitada. Limit restringe la colección devuelta, pero no sustituye los límites de la capa de solicitud cuando el llamador elige el tamaño.

<a name="choose-the-missing-row-contract"></a>
## Elija el contrato de registro ausente

First y Find devuelven (nil, nil) sin coincidencia. `FirstOrFail` y `FindOrFail` devuelven error compatible con `model.ErrModelNotFound`. Sole exige exactamente uno: `model.ErrModelNotFound` para cero y `model.ErrMultipleRecordsFound` para más de uno. Elija el terminal cuya cardinalidad forme parte del contrato, sin tratar toda ausencia igual.

Conserve los errores encapsulados de base de datos y decodificación para que una conexión no disponible o un valor almacenado inválido no se traten como una simple ausencia. Un Grant vacío es otro caso: `model.ErrNoTenant` indica que la ejecución fue rechazada antes de construir una sentencia limitada al tenant.

<a name="verify-tenant-scope-and-sql-shape"></a>
## Verifique ámbito de tenant y formato de SQL

Ejecute la lista en un conector real registrado con filas abiertas de dos tenants, empate en el borde de página y una fila cerrada en el tenant elegido. Confirme que el primero recibe solo sus 25 abiertos más recientes por `created_at` e id descendentes, el segundo solo su fila y Collector registra una consulta con `tenant_id`, ambos órdenes, límite y bindings de tenant y estado.

Llame a la misma operación terminal con el Grant cero y un Collector nuevo. El resultado debe corresponder a `model.ErrNoTenant`, y el número de consultas debe permanecer en cero. Ejecute también First y `FirstOrFail` para un registro visible solo para el otro tenant, midiendo los contratos de ausencia opcional y de error, sin inferirlos de los tipos.

<a name="builder-and-terminal-families"></a>
## Familias de builders y operaciones terminales

El constructor generado inicia una consulta tipada que reenvía a `model.Builder`. Where y variantes añaden condiciones vinculadas; Select, joins, agrupación y orden moldean SQL; Limit restringe filas. Get devuelve la colección generada, como `InvoiceCollection`. First y Find son lecturas opcionales; `OrFail` exige una fila y Sole exactamente una. Paginación y mutaciones son terminales separados con la misma frontera de contexto y Grant.

Esta guía se verificó contra hesape :hesapeVersion y framework :frameworkVersion. Las versiones posteriores pueden añadir o cambiar métodos; use la versión registrada en la página al comparar firmas.

<a name="authorization-precedes-the-query"></a>
## La autorización precede a la consulta

El predicado de tenant aísla datos, pero no es toda la policy. `security.Authorize` pregunta si Subject puede ejecutar `invoice.list` y solo emite Grant tras permiso; el terminal usa el tenant del Grant. En examples :examplesVersion, `app/Services/CategoryService.go` autoriza antes de terminales, `app/Models/Category.go` declara la entidad y su tabla y `tests/Feature/CategoryModel_test.go` prueba tenant y rechazo de Grant vacío.

Use relaciones cuando filas relacionadas formen parte de la consulta, paginación cuando el cliente necesite navegación y un recurso explícito cuando las filas salgan de aplicación. Ninguna capa acepta columna, tenant ni autoridad copiados directamente de la petición.
