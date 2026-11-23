# Tipos y conversiones de la base de datos

- [Los campos Go definen los tipos de datos](#go-fields-are-the-primary-data-types)
- [Declare intención de almacenamiento en la struct](#declare-storage-intent-in-the-struct)
- [Respete los límites de asignación de la versión publicada](#stay-inside-the-released-assignment-surface)
- [La conversión en tiempo de ejecución puede fallar](#runtime-conversion-can-still-fail)
- [Pruebe la lectura y escritura de valores límite en cada controlador](#round-trip-boundary-values-on-each-driver)
- [Tags y límites actuales de asignación](#tags-and-current-assignment-limits)
- [Los tipos de almacenamiento y respuesta pueden diferir](#storage-types-and-response-types-may-differ)

<a name="go-fields-are-the-primary-data-types"></a>
## Los campos Go definen los tipos de datos

La hidratación del modelo rellena campos declarados en la struct y no dispone de un registro general de casts. La ruta de asignación publicada cubre valores directamente asignables, punteros, `time.Time`, cadenas y slices de bytes, booleanos, conversiones numéricas por reflexión y campos cuya dirección implementa `sql.Scanner`. Esta decodificación es explícita y pertenece al tipo; no promete convertir cualquier representación almacenada en un tipo de dominio.

<a name="declare-storage-intent-in-the-struct"></a>
## Declare intención de almacenamiento en la struct

Use el menor tipo que represente el dominio y admita el rango de la base. Los punteros distinguen SQL NULL del valor cero. Almacene dinero como entero de la menor unidad cuando exija precisión fija; use tags db solo para mapear nombres o excluir campos.

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
## Respete los límites de asignación de la versión publicada

Use punteros para campos anulables solo cuando el conector devuelva un valor compatible. Mantenga fechas, cadenas, bytes, booleanos y números en las representaciones probadas del controlador. En :hesapeVersion, implemente `sql.Scanner` en el puntero del campo para controlar la hidratación y `driver.Valuer` en el valor para controlar los parámetros de escritura. JSON no se decodifica automáticamente; un tipo de dominio puede implementar ambas interfaces para gestionar esa codificación.

<a name="runtime-conversion-can-still-fail"></a>
## La conversión en tiempo de ejecución puede fallar

Un origen y destino incompatibles generan un error de asignación, pero la ruta numérica usa conversión por reflexión; esta guía no afirma que detecte todos los desbordamientos. NULL, fechas y representaciones específicas del controlador deben comprobarse con el conector de producción. Valide los rangos del dominio antes de escribir y no considere un valor cero decodificado como prueba de validez de los datos originales.

<a name="round-trip-boundary-values-on-each-driver"></a>
## Pruebe la lectura y escritura de valores límite en cada controlador

Pruebe los valores numéricos mínimos y máximos aceptados, NULL y no NULL, zonas horarias y precisión, cadenas y slices de bytes en cada dialecto de producción, porque los controladores devuelven representaciones Go diferentes. Pruebe explícitamente los límites de desbordamiento, sin suponer que la reflexión los rechaza. En campos de dominio personalizados, pruebe el ciclo de Scanner y Valuer, incluidos JSON malformado, valores de enum desconocidos y errores de vinculación.

<a name="tags-and-current-assignment-limits"></a>
## Tags y límites actuales de asignación

La tag db nombra columna; un guion excluye campo; la ausencia usa conversión de nombres del framework. La configuración del model elige tabla y comportamiento en Go, sin anotaciones ajenas. En :hesapeVersion, la hidratación respeta `sql.Scanner` y la escritura resuelve `driver.Valuer`; no existen decodificación JSON automática ni registro de casts con nombre.

<a name="storage-types-and-response-types-may-differ"></a>
## Los tipos de almacenamiento y respuesta pueden diferir

Un entero en céntimos en la base de datos puede convertirse en una cadena monetaria formateada en la vista, pero el modelo no debe almacenar presentación. Del mismo modo, las columnas cifradas deben descifrarse en un límite deliberado, no mediante un cast que revele texto sin cifrar cada vez que se cargue el modelo. Los recursos HTTP controlan la representación pública.
