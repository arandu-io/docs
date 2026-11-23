# Models

- [The model is the application's struct](#the-model-is-the-application-s-struct)
- [Declare columns and construct model metadata](#declare-columns-and-construct-model-metadata)
- [Query one typed row](#query-one-typed-row)
- [Create, fill, save, delete, and eager load](#create-fill-save-delete-and-eager-load)
- [Keep missing rows and invalid metadata visible](#keep-missing-rows-and-invalid-metadata-visible)
- [Tenant scope always comes from the Grant](#tenant-scope-always-comes-from-the-grant)
- [Test the model with a hostile neighboring tenant](#test-the-model-with-a-hostile-neighboring-tenant)
- [Models do not define the public payload](#models-do-not-define-the-public-payload)

<a name="the-model-is-the-application-s-struct"></a>
## The model is the application's struct

A model query hydrates declared Go fields. There is no dynamic attribute proxy between application code and the row, so a misspelled field fails at build time. The entity embeds `model.Model`, which adds instance operations and relation storage, and a package-level `model.NewTable` declares the table it is a row of.

<a name="declare-columns-and-construct-model-metadata"></a>
## Declare columns and construct model metadata

Export fields that are persisted and describe their column names with db tags when the default snake-case name is not enough. Configure the table, primary key, tenant column, timestamps, soft deletes, global scopes, and events in `model.TableSpec`, where every choice is visible to a reader. Run `aru model:build` after adding or changing an entity: it writes the typed query beside the struct, and `aru doctor` reports one that is stale.

<a name="query-one-typed-row"></a>
## Query one typed row

Start the generated query for the entity, compose conditions, and call a terminal with context and Grant. The returned pointer is the same User type the application declared, and the methods of the embedded Model are promoted onto it.

```go
type User struct {
	model.Model

	ID    int64  `db:"id"`
	Email string `db:"email"`
}

var userTable = model.NewTable(model.TableSpec{
	Name: "users",
	New:  func() model.Entity { return new(User) },
})

user, err := models.Users(db).Where("email", "=", email).First(ctx, g)
if err != nil { return nil, err }
return user, nil
```

<a name="create-fill-save-delete-and-eager-load"></a>
## Create, fill, save, delete, and eager load

Build a new row with `Users(db)`.`New()` and assign its fields one by one; a struct literal has no connection, and its write methods return `model.ErrUnwired`. Query terminals create and update sets, while rows returned by a query or by New can Save or Delete themselves under a Grant. Fill applies recognized exported attributes and ignores unknown keys; `ForceFill` preserves unknown raw attributes but still cannot reach unexported fields. Neither is an allowlist: build the changes from a typed request as a literal map of the named columns the use case may change, and never hand a request's input map to Fill or Update. With declares relations before execution, and Load adds them later with another explicit Grant.

<a name="keep-missing-rows-and-invalid-metadata-visible"></a>
## Keep missing rows and invalid metadata visible

Handle record-not-found separately from connection or decoding errors when the endpoint needs a 404. Invalid tags, unsupported field conversion, a missing table configuration, and write failures return errors; do not continue with a partially hydrated zero value. An unexported field is intentionally not persisted by reflection.

<a name="tenant-scope-always-comes-from-the-grant"></a>
## Tenant scope always comes from the Grant

Read and write terminals apply the tenant configured for the model using the Grant's tenant, and the tenant column takes its value from the Grant alone. On insert, that trusted tenant replaces any conflicting attribute supplied by the caller. Do not rely on the same protection when an existing row changes: a request map can carry the tenant column, a role, or a password hash under a spelling the model accepts once it reaches Fill or Update; the changes take only a literal map of named columns built from a typed request. A table shared by every tenant sets Global in its `TableSpec`, the only setting that drops the tenant filter, and should be an explicit, reviewable exception.

<a name="test-the-model-with-a-hostile-neighboring-tenant"></a>
## Test the model with a hostile neighboring tenant

Create identical identifiers or attributes under two tenant Grants and prove Find, update, delete, aggregates, and eager loads stay isolated. Verify the tenant attribute is overwritten on creation, instance Save keeps scope, a zero Grant cannot execute, and a struct literal refuses every write with `model.ErrUnwired`.

<a name="models-do-not-define-the-public-payload"></a>
## Models do not define the public payload

Relationships attach loaded rows to embedded model state, collections add row-aware set operations, and a pagination page carries navigation facts. HTTP resources still select the fields that leave the process. A new database column should not automatically become a response field merely because the model must export it for persistence.
