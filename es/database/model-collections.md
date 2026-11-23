# Colecciones de models

- [Una colección de models es un slice de structs de aplicación](#a-model-collection-is-a-slice-of-application-structs)
- [Consulte e inspeccione el conjunto tipado](#query-and-inspect-the-typed-set)
- [Cargue, seleccione, compare y convierta](#load-select-compare-and-convert)
- [Separe vacío, ausente y no consultable](#keep-empty-missing-and-unqueryable-distinct)
- [Pruebe por separado las operaciones en memoria y en almacenamiento](#prove-set-behavior-and-storage-behavior-separately)
- [Grupos de métodos de colección](#collection-method-groups)
- [Use paginación cuando el cliente necesite una ventana de resultados](#use-pagination-when-the-client-needs-a-window)

<a name="a-model-collection-is-a-slice-of-application-structs"></a>
## Una colección de models es un slice de structs de aplicación

La colección generada, como `UserCollection`, es un slice de punteros a la misma entidad que devuelve un terminal único. Sus métodos reenvían a `model.Rows`, que conoce claves, atributos cargados, relaciones y persistencia. Al ser un slice, se recorre, indexa y ordena con el lenguaje y la biblioteca estándar, sin que la capa de model duplique una biblioteca de colecciones.

<a name="query-and-inspect-the-typed-set"></a>
## Consulte e inspeccione el conjunto tipado

Get devuelve la colección generada directamente. len describe el conjunto, `ModelKeys` extrae las claves primarias y Find busca en la colección ya cargada, sin otra consulta a la base de datos. Estas operaciones en memoria no necesitan Grant porque no acceden al almacenamiento.

```go
users, err := models.Users(db).OrderBy("name").Get(ctx, g)
if err != nil { return err }

ids := users.ModelKeys()
user := users.Find(id)
```

<a name="load-select-compare-and-convert"></a>
## Cargue, seleccione, compare y convierta

Load y `LoadMissing` recuperan relaciones para el conjunto, y `LoadCount` añade conteos de relaciones. Only, Except y Unique operan por las claves de los modelos y devuelven el mismo tipo de colección, y Pluck lee un atributo cargado. Rows entrega el conjunto a `model.Rows` para Diff, Intersect, Merge, `LoadAggregate` y `ToQuery`, y el `UserCollectionOf` generado convierte las filas de vuelta.

<a name="keep-empty-missing-and-unqueryable-distinct"></a>
## Separe vacío, ausente y no consultable

Find devuelve nil cuando ninguna fila del conjunto tiene la clave; `FindOrFail` convierte la ausencia en error. `ToQuery` en un conjunto vacío devuelve `model.ErrEmptyCollection`; las relaciones y Fresh pueden fallar en la base. Una struct literal añadida al slice no tiene model, y Load rechaza el conjunto con `model.ErrUnwired`.

<a name="prove-set-behavior-and-storage-behavior-separately"></a>
## Pruebe por separado las operaciones en memoria y en almacenamiento

Pruebe comparaciones de clave, unicidad y vacío en filas construidas con New o con el Make de una factory, que no hacen I/O. Integre Load, Fresh, Push y `ToQuery` con dos tenants y recuento de consultas. En `LoadMissing`, precargue solo parte del conjunto y confirme la búsqueda de lo ausente sin descartar relaciones ya adjuntas.

<a name="collection-method-groups"></a>
## Grupos de métodos de colección

La colección generada trae `ModelKeys`, Find, `FindOrFail`, Contains, Pluck, Only, Except, Unique, Load, `LoadMissing`, `LoadCount`, Fresh, `ToArray` y Rows. `model.Rows` añade Count, `IsEmpty`, First, Diff, Intersect, Merge, Push, `LoadAggregate`, `ToQuery`, `MakeHidden` y `MakeVisible`. Los métodos de visibilidad ajustan la representación del modelo en array, no el tipo declarado de la respuesta HTTP.

<a name="use-pagination-when-the-client-needs-a-window"></a>
## Use paginación cuando el cliente necesite una ventana de resultados

Una colección contiene filas ya devueltas; no conoce total de páginas ni cursores siguientes. Los terminales de paginación devuelven la colección junto a una página con el estado de navegación. Para procesamiento sin límite en servidor, Chunk, Each o Cursor controlan memoria sin montar primero toda la tabla en una colección.
