# Database types and casts

- [Go fields are the primary data types](#go-fields-are-the-primary-data-types)
- [Declare storage intent in the struct](#declare-storage-intent-in-the-struct)
- [Stay inside the released assignment surface](#stay-inside-the-released-assignment-surface)
- [Runtime conversion can still fail](#runtime-conversion-can-still-fail)
- [Round-trip boundary values on each driver](#round-trip-boundary-values-on-each-driver)
- [Tags and current assignment limits](#tags-and-current-assignment-limits)
- [Storage types and response types may differ](#storage-types-and-response-types-may-differ)

<a name="go-fields-are-the-primary-data-types"></a>
## Go fields are the primary data types

Model hydration targets declared struct fields and has no general cast registry. The released assignment path covers directly assignable values, pointers, `time.Time`, strings and byte slices, booleans, numeric reflection conversions, and fields whose address implements `sql.Scanner`. That is an explicit, type-owned decoding surface, not a promise that arbitrary storage representations become domain types.

<a name="declare-storage-intent-in-the-struct"></a>
## Declare storage intent in the struct

Use the smallest type that honestly represents the domain while matching the database range. Pointers distinguish SQL NULL from a type's zero value. Store money as an integer of the smallest unit when fixed precision is required, and use db tags only to map names or exclude a field.

```go example
type Invoice struct {
	model.Model

	ID        string     `db:"id"`
	Total     int64      `db:"total"`
	DueDate   time.Time  `db:"due_date"`
	DeletedAt *time.Time `db:"deleted_at"`
}
```

<a name="stay-inside-the-released-assignment-surface"></a>
## Stay inside the released assignment surface

Use pointers for nullable fields only where the selected connector returns a compatible value, and keep time, string, bytes, booleans, and numeric fields inside tested driver representations. In :hesapeVersion, implement `sql.Scanner` on a field's pointer to own hydration and `driver.Valuer` on its value to own write bindings. JSON is not decoded automatically; a domain type may implement both interfaces when it owns that encoding.

<a name="runtime-conversion-can-still-fail"></a>
## Runtime conversion can still fail

An incompatible source and destination returns an assignment error, but the numeric path uses reflection conversion and this guide does not claim it detects every overflow. NULL, time, and driver-specific value shapes must be measured on the production connector. Validate domain ranges before writes and do not treat a decoded zero value as proof that source data was valid.

<a name="round-trip-boundary-values-on-each-driver"></a>
## Round-trip boundary values on each driver

Test minimum and maximum accepted numeric values, NULL and non-NULL cases, time zones and precision, strings, and byte slices against each production dialect because drivers return different Go shapes. Explicitly test overflow boundaries instead of assuming reflection rejects them. For custom domain fields, round-trip Scanner and Valuer behavior, including malformed JSON, unknown enum values, and binding errors.

<a name="tags-and-current-assignment-limits"></a>
## Tags and current assignment limits

A db tag names the column; a dash excludes the field; an absent tag uses the framework's name conversion. Model configuration chooses table and behavior in Go rather than hiding those decisions in unrelated annotations. In :hesapeVersion, hydration honors `sql.Scanner` and write bindings resolve `driver.Valuer`; automatic JSON decoding and a named cast registry remain absent.

<a name="storage-types-and-response-types-may-differ"></a>
## Storage types and response types may differ

A database integer in cents can become a formatted currency string in a view, but the model should not store presentation. Likewise, encrypted columns need decryption at a deliberate boundary, not a cast that silently reveals plaintext whenever the model is loaded. HTTP resources own the public representation.
