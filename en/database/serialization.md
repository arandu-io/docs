# Model serialization

- [Default JSON is convenient and intentionally not the API](#default-json-is-convenient-and-intentionally-not-the-api)
- [Map a model into an explicit response](#map-a-model-into-an-explicit-response)
- [Represent collections and optional relations](#represent-collections-and-optional-relations)
- [Stop on unsupported values or broken mapping](#stop-on-unsupported-values-or-broken-mapping)
- [Lock the public keys and sensitive absences](#lock-the-public-keys-and-sensitive-absences)
- [DTOs and `JsonResource` have different contracts](#dtos-and-jsonresource-have-different-contracts)
- [Logging needs its own safe representation](#logging-needs-its-own-safe-representation)

> [!WARNING]
> **Current limitations.** The released model representation and an official allowlisted DTO exist, but no tagged application journey proves model `ToArray` or `ToJSON` as a public response contract. Prefer an explicit DTO or `JsonResource` and test its exact keys.

<a name="default-json-is-convenient-and-intentionally-not-the-api"></a>
## Default JSON is convenient and intentionally not the API

An exported Go struct field is serialized by `encoding/json` unless a tag or custom marshaler says otherwise. That default is useful inside trusted boundaries and too broad for a public endpoint: adding an exported database column can change output without the handler or response contract being edited.

<a name="map-a-model-into-an-explicit-response"></a>
## Map a model into an explicit response

Declare the fields the endpoint promises and construct that type from an already authorized model. The listing is an ordinary Go DTO serialized by `encoding/json`; it is not the framework's `hhttp.JsonResource` contract. The response can reuse stable identifiers and domain values while omitting tenant columns, password hashes, storage paths, internal timestamps, and relation state that the client does not need.

```go example
type User struct{ ID, Name string }

type UserResource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewUserResource(user User) UserResource {
	return UserResource{ID: user.ID, Name: user.Name}
}
```

<a name="represent-collections-and-optional-relations"></a>
## Represent collections and optional relations

Map each model in a collection and wrap items with pagination facts when the endpoint is paginated. Include a related resource only when the relation was deliberately loaded and the contract calls for it. Use omitempty only when absence has defined meaning; it is not a substitute for choosing public fields.

<a name="stop-on-unsupported-values-or-broken-mapping"></a>
## Stop on unsupported values or broken mapping

JSON encoding can fail for unsupported values and custom `MarshalJSON` implementations can return errors. A missing required domain value should fail mapping rather than become misleading output. Do not fall back to serializing the model when a resource constructor fails, because that reopens the field leak the resource exists to close.

<a name="lock-the-public-keys-and-sensitive-absences"></a>
## Lock the public keys and sensitive absences

Assert the exact response shape or decoded key set and explicitly assert sensitive model fields are absent. Add a fixture-only exported field to the model or mapping input and prove it does not appear automatically. Cover loaded and unloaded optional relations and encoding error paths.

<a name="dtos-and-jsonresource-have-different-contracts"></a>
## DTOs and `JsonResource` have different contracts

`ToArray` on a row or a collection, visibility controls, hidden fields, and appended attributes describe the model package's array representation. An ordinary DTO follows `encoding/json` tags and `MarshalJSON`. `Context.JSON` instead accepts `hhttp.JsonResource`, whose `ToArray` and With methods return runtime maps; those string keys require exact response tests. Choose one transport contract explicitly rather than calling both shapes the same resource.

<a name="logging-needs-its-own-safe-representation"></a>
## Logging needs its own safe representation

A response allowlist does not automatically redact structured logs, events, or debug dumps. Types carrying credentials or key material should expose safe logging behavior and deliberate event payloads. Encryption protects stored plaintext but does not make ciphertext or its metadata appropriate for every public response.
