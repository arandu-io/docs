# Fluent query builder

- [Compose before execution](#compose-before-execution)
- [Query open invoices under a Grant](#query-open-invoices-under-a-grant)
- [Filter and order with allowlists](#filter-and-order-with-allowlists)
- [Choose the missing-row contract](#choose-the-missing-row-contract)
- [Verify tenant scope and SQL shape](#verify-tenant-scope-and-sql-shape)
- [Builder and terminal families](#builder-and-terminal-families)
- [Authorization precedes the query](#authorization-precedes-the-query)

<a name="compose-before-execution"></a>
## Compose before execution

A query starts from the constructor `aru model:build` generates beside each entity, such as `models.Invoices(db)`, and composition stays separate from execution. Where, `OrderByDesc`, and Limit describe a statement without sending it. Get is a terminal: it receives the request context and a `security.Grant`, applies the table's tenant column, compiles the dialect-specific SQL, executes it, and returns the generated collection, a slice of the entity's pointers.

Keeping the Grant on the terminal matters. Building a reusable query authorizes nothing, while the statement that reads rows must carry the authority and tenant of the caller that executes it. A zero Grant is rejected with `model.ErrNoTenant` before SQL reaches the connection.

<a name="query-open-invoices-under-a-grant"></a>
## Query open invoices under a Grant

Invoice embeds `model.Model` and `invoiceTable` declares its table; both belong in `app/Models/Invoice.go`. `aru model:build` writes Invoices, `InvoiceQuery`, and `InvoiceCollection` beside them in `InvoiceQuery.go`, a file that is regenerated rather than edited. `OpenInvoices` accepts the instrumented database handle and a Grant already issued for `invoice.list`. The default tenant column is `tenant_id`, so the terminal adds the trusted tenant independently of the status predicate.

The service obtains the Grant through `security.Authorize` with its Policy, authenticated Subject, `invoiceList` action, and the collection resource. Passing the returned Grant into `OpenInvoices` preserves the order used by the released examples application: authorize first, then execute the data operation.

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
## Filter and order with allowlists

Values passed to Where are bindings. Column names, sort directions, selected expressions, and raw fragments are SQL structure and cannot be made safe by a value placeholder. Translate a request choice through a fixed map such as `map[string]string{"newest": "created_at", "total": "total"}`; reject a key the map does not contain before calling `OrderBy` or Select.

Compose additional predicates before the terminal and keep values separate from identifiers. Use a stable secondary ordering when two rows may share the requested column. Limit bounds the returned collection; it does not replace request-layer bounds when the caller is allowed to choose a size.

<a name="choose-the-missing-row-contract"></a>
## Choose the missing-row contract

First and Find return (nil, nil) when no row matches. `FirstOrFail` and `FindOrFail` return an error that matches `model.ErrModelNotFound`. Sole requires exactly one row: it returns `model.ErrModelNotFound` for zero and `model.ErrMultipleRecordsFound` for more than one. Choose the terminal whose cardinality is part of the endpoint contract instead of converting every absence into the same outcome.

Preserve wrapped database and decoding errors so an unavailable connection or invalid stored value is not reported as an ordinary miss. An empty Grant is different again: `model.ErrNoTenant` means execution was refused before a tenant-scoped statement could be built.

<a name="verify-tenant-scope-and-sql-shape"></a>
## Verify tenant scope and SQL shape

Run the listing against a real registered connector with matching open rows in two tenants, a tie at the page boundary, and a closed row in the selected tenant. Assert the first tenant receives only its newest twenty-five open rows in descending `created_at` then id order, the second receives only its own row, and the Collector records one statement containing `tenant_id`, both descending orderings, the limit, and the tenant and status bindings.

Call the same terminal with the zero Grant and a fresh Collector. The result must match `model.ErrNoTenant` and the query count must remain zero. Also execute First and `FirstOrFail` for a row visible only to the other tenant so the optional and failing absence contracts are measured rather than inferred from types.

<a name="builder-and-terminal-families"></a>
## Builder and terminal families

The generated constructor starts a typed query that forwards to `model.Builder`. Where and its predicate variants add bound conditions; Select, joins, grouping, and ordering shape the statement; Limit constrains rows. Get returns the generated collection, such as `InvoiceCollection`. First and Find are optional single-row reads, their `OrFail` forms require a row, and Sole requires one and only one. Pagination and mutation methods are separate terminals and retain the same context-and-Grant boundary.

This guide is verified against hesape :hesapeVersion and framework :frameworkVersion. A later module version may add or change methods, so use the version recorded with the page when comparing a signature.

<a name="authorization-precedes-the-query"></a>
## Authorization precedes the query

A tenant predicate is data isolation, not the entire policy decision. `security.Authorize` asks whether this Subject may perform `invoice.list` and issues the Grant only after the Policy allows it; the model terminal then uses that Grant's tenant. In examples :examplesVersion, `app/Services/CategoryService.go` authorizes before model terminals, `app/Models/Category.go` declares the entity and its table, and `tests/Feature/CategoryModel_test.go` exercises tenant binding and zero-Grant refusal.

Use relationships when related rows are part of the query, pagination when a client needs navigation facts, and an explicit response resource when rows leave the application. None of those layers should accept a column name, tenant, or authority copied directly from request input.
