# Relaciones y carga anticipada

- [Las relaciones son metadatos del model con carga explícita](#relations-are-model-metadata-with-explicit-loading)
- [Incorpore Model y defina ambas claves](#embed-model-and-define-both-keys)
- [Declare una relación y cárguela anticipadamente](#declare-and-eager-load-one-relation)
- [Restrinja, cuente, asocie y sincronice](#constrain-count-attach-and-synchronize)
- [No confunda una relación no cargada con una vacía](#do-not-treat-an-unloaded-relation-as-empty)
- [La carga anticipada sigue siendo un límite de lectura](#the-eager-path-is-still-a-read-boundary)
- [Cuente consultas en el caso simple de eager loading](#count-statements-for-the-simple-eager-case)
- [Las colecciones preservan el estado de relaciones cargadas](#collections-carry-loaded-relation-state-forward)

<a name="relations-are-model-metadata-with-explicit-loading"></a>
## Las relaciones son metadatos del model con carga explícita

Una relación describe cómo se conectan dos claves de modelos y cómo se asocian los registros cargados al padre. Nunca se ejecuta solo porque se lea un campo. Este límite explícito mantiene visible el acceso a la base de datos y garantiza que cada consulta de relación reciba el mismo contexto y Grant que la consulta principal.

<a name="embed-model-and-define-both-keys"></a>
## Incorpore Model y defina ambas claves

Ambas entidades incorporan `model.Model` y declaran una tabla. Registre cada relación en la tabla del padre con Relate, en una función init junto a la entidad, antes de la primera consulta: dele un nombre estable, la tabla relacionada, las claves foráneas, locales o propietarias y los detalles de pivote cuando hagan falta. Las restricciones de base deben imponer la relación de claves independientemente de la API de consulta.

<a name="declare-and-eager-load-one-relation"></a>
## Declare una relación y cárguela anticipadamente

Elija un constructor de relación exportado por la versión fijada para las claves reales, como `HasOne`, `HasMany` o `BelongsTo`. Las variantes through, polimórficas y muchos-a-muchos exigen más metadatos; documente con las firmas exactas de la etiqueta, no con recuentos memorizados. With nombra la relación antes del terminal. Related lee las filas cargadas como `model.Rows`, y el `PostCollectionOf` generado las convierte a la entidad relacionada.

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
## Restrinja, cuente, asocie y sincronice

Los predicados de relación como `WhereHas` filtran los padres, y `WithCount` o los demás helpers de agregación añaden campos de resumen. `WithConstraints`, del builder del núcleo, restringe los resultados de la carga anticipada desde un ámbito local en el bloque custom de la entidad. Load y `LoadMissing` operan después de recuperar los datos. Las relaciones muchos-a-muchos exponen Attach, Detach y Sync bajo un Grant, manteniendo los cambios de pivote explícitos y limitados al tenant.

<a name="do-not-treat-an-unloaded-relation-as-empty"></a>
## No confunda una relación no cargada con una vacía

Related informa si se cargó el valor; compruebe ese booleano para no confundir falta de carga con relación vacía. Nombres desconocidos, tipos de clave incompatibles, metadatos de pivote ausentes y fallos de base devuelven errores. Una struct literal o una copia por valor de una fila no tiene model al que adjuntar la relación, y Load la rechaza con `model.ErrUnwired`.

<a name="the-eager-path-is-still-a-read-boundary"></a>
## La carga anticipada sigue siendo un límite de lectura

Toda consulta de relación debe llevar el Grant y restringir los registros de ambos lados al tenant. La falta del filtro de tenant en el lado relacionado es especialmente peligrosa: una sola consulta puede asociar registros de otro cliente a varios padres. Nunca acepte nombres de relaciones o columnas de pivote directamente de la solicitud.

<a name="count-statements-for-the-simple-eager-case"></a>
## Cuente consultas en el caso simple de eager loading

Para una carga has-many simple y no polimórfica, verifique una consulta de padres más una de relacionados, y confirme que los padres sin hijos reciben una colección cargada vacía. Los casos polimórficos y anidados pueden necesitar consultas adicionales: pruebe por separado cada forma documentada. Incluya claves coincidentes en otro tenant para demostrar que sus registros nunca se asocian.

<a name="collections-carry-loaded-relation-state-forward"></a>
## Las colecciones preservan el estado de relaciones cargadas

Las colecciones de modelos pueden ejecutar Load o `LoadMissing` para todos los registros en conjunto; los predicados de relación permanecen en Builder. Los recursos HTTP deciden si las relaciones cargadas forman parte de la respuesta. Las claves foráneas protegen la integridad de la base de datos, pero no sustituyen las políticas ni los filtros de tenant al leer y escribir.
