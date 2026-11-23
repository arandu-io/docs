# Modelos

- [El modelo es la struct de la aplicación](#the-model-is-the-application-s-struct)
- [Declare columnas y construya los metadatos del modelo](#declare-columns-and-construct-model-metadata)
- [Consulte un registro tipado](#query-one-typed-row)
- [Cree, rellene, guarde, elimine y cargue relaciones](#create-fill-save-delete-and-eager-load)
- [Mantenga visibles filas ausentes y metadatos inválidos](#keep-missing-rows-and-invalid-metadata-visible)
- [El ámbito del tenant siempre proviene del Grant](#tenant-scope-always-comes-from-the-grant)
- [Pruebe el model con un tenant vecino hostil](#test-the-model-with-a-hostile-neighboring-tenant)
- [Los modelos no definen el contenido de la respuesta pública](#models-do-not-define-the-public-payload)

<a name="the-model-is-the-application-s-struct"></a>
## El modelo es la struct de la aplicación

Una consulta de model hidrata campos Go declarados. No existe un proxy dinámico de atributos entre código y fila; un campo mal nombrado falla en build. La entidad incorpora `model.Model`, que añade operaciones de instancia y almacenamiento de relaciones, y un `model.NewTable` a nivel de paquete declara la tabla de la que es fila.

<a name="declare-columns-and-construct-model-metadata"></a>
## Declare columnas y construya los metadatos del modelo

Exporte los campos persistidos y describa sus columnas con etiquetas db cuando el nombre predeterminado en `snake_case` no sea suficiente. Configure tabla, clave primaria, columna del tenant, marcas de tiempo, borrado lógico, ámbitos globales y eventos en `model.TableSpec`, donde cada decisión queda visible para quien lee. Ejecute `aru model:build` al añadir o cambiar una entidad: escribe la consulta tipada junto a la struct, y `aru doctor` señala la que esté desactualizada.

<a name="query-one-typed-row"></a>
## Consulte un registro tipado

Inicie la consulta generada de la entidad, componga condiciones y llame a un terminal con contexto y Grant. El puntero devuelto tiene el mismo tipo User declarado, y los métodos del Model incorporado se promueven a él.

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
## Cree, rellene, guarde, elimine y cargue relaciones

Construya una fila nueva con `Users(db)`.`New()` y asigne sus campos uno a uno; una struct literal no tiene conexión, y sus métodos de escritura devuelven `model.ErrUnwired`. Las operaciones terminales de la consulta crean y actualizan conjuntos; las filas devueltas por una consulta o por New pueden ejecutar Save o Delete bajo un Grant. Fill aplica atributos exportados reconocidos e ignora claves desconocidas; `ForceFill` conserva atributos sin procesar desconocidos, pero tampoco alcanza campos no exportados. Ninguno de los dos es una allowlist: construya los cambios a partir de una solicitud tipada, como un mapa literal de las columnas nombradas que el caso de uso puede cambiar, y nunca entregue el mapa de entrada de una solicitud a Fill o Update. With declara relaciones antes de ejecutar, y Load las añade después con otro Grant explícito.

<a name="keep-missing-rows-and-invalid-metadata-visible"></a>
## Mantenga visibles filas ausentes y metadatos inválidos

Separe registro no encontrado de errores de conexión o decodificación cuando el endpoint necesite 404. Tags inválidas, conversión no compatible, configuración de tabla ausente y fallos de escritura devuelven error; no continúe con un valor cero parcialmente hidratado. Los campos no exportados no se persisten por reflexión deliberadamente.

<a name="tenant-scope-always-comes-from-the-grant"></a>
## El ámbito del tenant siempre proviene del Grant

Las operaciones terminales de lectura y escritura aplican el tenant configurado en el modelo usando el tenant del Grant, y la columna de tenant recibe su valor solo del Grant. Al insertar, este tenant confiable sustituye cualquier atributo contradictorio enviado por el llamador. No cuente con la misma protección cuando cambia una fila existente: un mapa de la solicitud puede llevar la columna de tenant, un rol o un hash de contraseña bajo una grafía que el modelo acepta cuando llega a Fill o Update; los cambios reciben solo un mapa literal de columnas nombradas, construido a partir de una solicitud tipada. Una tabla compartida por todos los tenants define Global en su `TableSpec`, la única configuración que quita el filtro de tenant, y debe ser una excepción explícita y revisable.

<a name="test-the-model-with-a-hostile-neighboring-tenant"></a>
## Pruebe el model con un tenant vecino hostil

Cree identificadores o atributos iguales bajo Grants de dos tenants y demuestre el aislamiento en Find, actualización, eliminación, agregaciones y carga anticipada. Compruebe que el atributo de tenant se sustituye al crear, que Save de la instancia conserva el ámbito, que el Grant cero no ejecuta consultas y que una struct literal rechaza toda escritura con `model.ErrUnwired`.

<a name="models-do-not-define-the-public-payload"></a>
## Los modelos no definen el contenido de la respuesta pública

Las relaciones adjuntan filas cargadas al estado incorporado del model; las colecciones añaden operaciones conscientes de filas y la página de paginación lleva la navegación. Los recursos HTTP aún eligen campos que salen del proceso. Una columna nueva no debe convertirse automáticamente en respuesta solo por necesitar exportarse para persistir.
