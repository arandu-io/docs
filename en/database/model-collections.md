# Model collections

- [A model collection is a slice of application structs](#a-model-collection-is-a-slice-of-application-structs)
- [Query and inspect the typed set](#query-and-inspect-the-typed-set)
- [Load, select, compare, and convert](#load-select-compare-and-convert)
- [Keep empty, missing, and unqueryable distinct](#keep-empty-missing-and-unqueryable-distinct)
- [Prove set behavior and storage behavior separately](#prove-set-behavior-and-storage-behavior-separately)
- [Collection method groups](#collection-method-groups)
- [Use pagination when the client needs a window](#use-pagination-when-the-client-needs-a-window)

<a name="a-model-collection-is-a-slice-of-application-structs"></a>
## A model collection is a slice of application structs

The generated collection, such as `UserCollection`, is a slice of pointers to the same entity a single terminal returns. Its methods forward to `model.Rows`, which understands model keys, loaded attributes, relations, and persistence. Being a slice, it ranges, indexes, and sorts with the language and the standard library, so the model layer does not duplicate a collection library.

<a name="query-and-inspect-the-typed-set"></a>
## Query and inspect the typed set

Get returns the generated collection directly. len describes the set, `ModelKeys` extracts primary keys, and Find searches the already loaded collection without another database call. None of those in-memory operations needs a Grant because they do not reach storage.

```go
users, err := models.Users(db).OrderBy("name").Get(ctx, g)
if err != nil { return err }

ids := users.ModelKeys()
user := users.Find(id)
```

<a name="load-select-compare-and-convert"></a>
## Load, select, compare, and convert

Load and `LoadMissing` fetch relations for the set, and `LoadCount` adds relation counts. Only, Except, and Unique operate by model keys and return the same collection type, and Pluck reads one loaded attribute. Rows hands the set to `model.Rows` for Diff, Intersect, Merge, `LoadAggregate`, and `ToQuery`, and the generated `UserCollectionOf` converts rows back.

<a name="keep-empty-missing-and-unqueryable-distinct"></a>
## Keep empty, missing, and unqueryable distinct

Find returns nil when no row in the set has the key, while `FindOrFail` makes absence an error. `ToQuery` on an empty set returns `model.ErrEmptyCollection`; relation and Fresh operations can fail at the database. A struct literal appended to the slice has no model, and Load refuses the set with `model.ErrUnwired`.

<a name="prove-set-behavior-and-storage-behavior-separately"></a>
## Prove set behavior and storage behavior separately

Unit-test key comparisons, uniqueness, and empty behavior on rows built with New or a factory's Make, which perform no I/O. Integration-test Load, Fresh, Push, and `ToQuery` with two tenants and statement counts. For `LoadMissing`, pre-load only part of the set and assert it fetches what is absent without discarding already attached rows.

<a name="collection-method-groups"></a>
## Collection method groups

The generated collection carries `ModelKeys`, Find, `FindOrFail`, Contains, Pluck, Only, Except, Unique, Load, `LoadMissing`, `LoadCount`, Fresh, `ToArray`, and Rows. `model.Rows` adds Count, `IsEmpty`, First, Diff, Intersect, Merge, Push, `LoadAggregate`, `ToQuery`, `MakeHidden`, and `MakeVisible`. Visibility methods adjust the model-to-array representation, not the declared HTTP response type.

<a name="use-pagination-when-the-client-needs-a-window"></a>
## Use pagination when the client needs a window

A collection is the rows already returned; it does not know total pages or continuation cursors. The pagination terminals return the collection beside a page that holds the navigation state. For unbounded server-side processing, Chunk, Each, or Cursor keeps memory controlled instead of first building one collection containing the entire table.
