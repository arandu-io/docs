# Model factories

- [Factories build typed rows for tests and seeds](#factories-build-typed-rows-for-tests-and-seeds)
- [Configure the model before defining its factory](#configure-the-model-before-defining-its-factory)
- [Make deterministic data or create authorized rows](#make-deterministic-data-or-create-authorized-rows)
- [Compose count, state, sequence, and relationships](#compose-count-state-sequence-and-relationships)
- [Account for partial persistence and callback errors](#account-for-partial-persistence-and-callback-errors)
- [Creation still carries the caller's Grant](#creation-still-carries-the-caller-s-grant)
- [Prove repeatability, isolation, and tenant ownership](#prove-repeatability-isolation-and-tenant-ownership)
- [Use factories from focused tests and idempotent seeders](#use-factories-from-focused-tests-and-idempotent-seeders)

<a name="factories-build-typed-rows-for-tests-and-seeds"></a>
## Factories build typed rows for tests and seeds

`aru make:factory` writes a typed factory for an entity, such as `factories.Users(db)` in `database/factories/UserFactory.go`, and `aru model:build` keeps everything outside its custom block current. It holds one valid default and composes deterministic state around it over the non-generic `factories.Factory`. The definition returns User, states receive `*User`, and results are the same rows query terminals return. Misspelled fields fail compilation instead of disappearing in a map.

<a name="configure-the-model-before-defining-its-factory"></a>
## Configure the model before defining its factory

Declare the entity and its table in app/Models and run `aru model:build`. The typed factory is built on the entity's generated query, so the table supplies the key settings and tenant convention and the db argument supplies the connection; there is no model value to construct by hand.

- Declare the entity and its table.
- Keep one valid default definition.
- Use deterministic Faker values.
- Supply Grant only to persistence operations.

<a name="make-deterministic-data-or-create-authorized-rows"></a>
## Make deterministic data or create authorized rows

`MakeOne` performs no I/O and needs no context or Grant; `CreateOne` writes through the model with both. Default seed is fixed at 1, and `Seed(n)` makes a failure reproducible.

```go
func defineUser(f faker.Faker) models.User {
	return models.User{Name: f.Name(), Email: f.Unique().Email(), Active: true}
}

made, err := factories.Users(db).MakeOne()
if err != nil { return err }
stored, err := factories.Users(db).Count(3).Create(ctx, grant)
if err != nil { return err }
```

<a name="compose-count-state-sequence-and-relationships"></a>
## Compose count, state, sequence, and relationships

Count and Seed choose volume and repeatability; State runs in order; Sequence cycles; `AfterMaking` and `AfterCreating` add callbacks. The typed factory forwards those; Has, which creates children after parents, and `ForParent`, which creates a parent first, are on the core factory and are reached from a method in the custom block. Every modifier returns a new factory.

- `Count(n)` and `Seed(n)` control a reproducible batch.
- State and Sequence compose typed changes.
- `AfterCreating` receives context and Grant.
- Has and `ForParent` require an explicit link callback.

<a name="account-for-partial-persistence-and-callback-errors"></a>
## Account for partial persistence and callback errors

Make can fail during model assignment. Create saves sequentially, so a later Save or callback error can leave earlier rows because Factory does not wrap the whole batch. Count below zero becomes zero, and a valid but wrong relationship callback creates misleading data.

> [!WARNING]
> Wrap factory creation in a caller-owned transaction when atomicity matters, and assert relationship keys rather than assuming a compiling callback connected the intended records.

<a name="creation-still-carries-the-caller-s-grant"></a>
## Creation still carries the caller's Grant

> [!WARNING]
> A factory is not an authorization bypass. Create and relationship callbacks use the supplied Grant; never manufacture a tenant or use system authority in a request path. Production seeds contain no fixed passwords, tokens, or real personal data.

<a name="prove-repeatability-isolation-and-tenant-ownership"></a>
## Prove repeatability, isolation, and tenant ownership

Prove the same seed repeats and another diverges, base factories do not mutate, Make performs no DB calls, Create writes the expected tenant and count, denied Grants refuse, relation keys match, and partial callback failure is handled deliberately.

- Seeded values repeat.
- Make has no I/O.
- Create preserves tenant authority.
- Modifiers do not leak between tests.

<a name="use-factories-from-focused-tests-and-idempotent-seeders"></a>
## Use factories from focused tests and idempotent seeders

Models define persistence, relationships define links, seeders own named bootstrap operations, and database tests isolate each created dataset.

- [Models](/docs/{{version}}/database/models)
- [Relationships](/docs/{{version}}/database/relationships)
- [Seeders](/docs/{{version}}/database/seeders)
- [Database testing](/docs/{{version}}/testing/database)
