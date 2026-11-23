# Relationships and eager loading

- [Relations are model metadata with explicit loading](#relations-are-model-metadata-with-explicit-loading)
- [Embed Model and define both keys](#embed-model-and-define-both-keys)
- [Declare and eager load one relation](#declare-and-eager-load-one-relation)
- [Constrain, count, attach, and synchronize](#constrain-count-attach-and-synchronize)
- [Do not treat an unloaded relation as empty](#do-not-treat-an-unloaded-relation-as-empty)
- [The eager path is still a read boundary](#the-eager-path-is-still-a-read-boundary)
- [Count statements for the simple eager case](#count-statements-for-the-simple-eager-case)
- [Collections carry loaded relation state forward](#collections-carry-loaded-relation-state-forward)

<a name="relations-are-model-metadata-with-explicit-loading"></a>
## Relations are model metadata with explicit loading

A relation describes how two model keys connect and how loaded rows attach to the parent. It never runs merely because a field is read. That explicit boundary keeps database work visible and ensures every relationship query receives the same context and Grant as its parent query.

<a name="embed-model-and-define-both-keys"></a>
## Embed Model and define both keys

Both entities embed `model.Model` and declare a table. Register each relation on the parent's table with Relate, in an init function beside the entity, before the first query: give it a stable name, the related table, the foreign and local or owner keys, and pivot details where needed. Database constraints should enforce the key relationship independently of the query API.

<a name="declare-and-eager-load-one-relation"></a>
## Declare and eager load one relation

Choose a relation constructor that the pinned model release exports for the keys you actually have, such as `HasOne`, `HasMany`, or `BelongsTo`. Through, polymorphic, and many-to-many variants have additional metadata and should be documented only from their exact tagged signatures rather than from a memorized family count. With names the relation before the terminal. Related reads the loaded rows back as `model.Rows`, and the generated `PostCollectionOf` converts them to the related entity.

```go
func init() {
	userTable.Relate("posts", func(u *model.Model) model.Relation {
		return model.HasMany(u, postTable, "user_id", "id")
	})
}

user, err := models.Users(db).With("posts").FirstOrFail(ctx, g)
if err != nil { return err }
rows, loaded := user.Related("posts")
posts := models.PostCollectionOf(rows)
```

<a name="constrain-count-attach-and-synchronize"></a>
## Constrain, count, attach, and synchronize

Relation predicates such as `WhereHas` filter parents, and `WithCount` or the other aggregate helpers add summary fields. The core builder's `WithConstraints` narrows eager results from a local scope in the entity's custom block. Load and `LoadMissing` operate after retrieval. Many-to-many relations expose Attach, Detach, and Sync under a Grant so pivot changes remain tenant-scoped and explicit.

<a name="do-not-treat-an-unloaded-relation-as-empty"></a>
## Do not treat an unloaded relation as empty

Related reports whether the value was loaded; branch on that boolean rather than confusing not loaded with no related rows. Unknown relation names, incompatible key types, missing pivot metadata, and database failures return errors. A struct literal or a value copy of a row has no model to attach a relation to, and Load refuses it with `model.ErrUnwired`.

<a name="the-eager-path-is-still-a-read-boundary"></a>
## The eager path is still a read boundary

Every relation query must carry the Grant and scope both parent and related rows. A missing related-side tenant predicate is especially dangerous because one query can attach another customer's rows to many parents at once. Never accept a relation name or pivot column directly from request input.

<a name="count-statements-for-the-simple-eager-case"></a>
## Count statements for the simple eager case

For a simple non-polymorphic has-many load, assert one parent query plus one related query and verify parents with no children receive an empty loaded collection. Polymorphic and nested cases may require additional statements, so test their documented shape separately. Include matching keys in another tenant to prove they never attach.

<a name="collections-carry-loaded-relation-state-forward"></a>
## Collections carry loaded relation state forward

Model collections can Load or `LoadMissing` for every row together, while relationship predicates stay on Builder. HTTP resources decide whether loaded relations become output. Database foreign keys protect integrity, but they do not replace policy or tenant filtering on reads and writes.
