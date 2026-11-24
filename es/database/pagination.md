# Paginación

- [Un paginador proporciona navegación, no solo registros](#a-paginator-answers-navigation-not-only-rows)
- [Elija paginación con total, simple o por cursor](#choose-length-aware-simple-or-cursor-pagination)
- [Conserve los filtros y un orden estable](#preserve-filters-and-stable-ordering)
- [Rechace entrada malformada y posiciones imposibles](#reject-malformed-input-and-impossible-positions)
- [Pruebe límites y cambios concurrentes](#test-boundaries-and-concurrent-change)
- [Tres formas de operaciones terminales](#three-terminal-shapes)
- [La paginación forma parte del contrato de respuesta](#pagination-is-part-of-the-response-contract)

<a name="a-paginator-answers-navigation-not-only-rows"></a>
## Un paginador proporciona navegación, no solo registros

Un slice solo no distingue la última página llena de una página llena con más datos. Los terminales de paginación devuelven elementos tipados y estado para navegación estable, ejecutando bajo el mismo Grant y tenant de Get.

<a name="choose-length-aware-simple-or-cursor-pagination"></a>
## Elija paginación con total, simple o por cursor

Paginate ejecuta la consulta de datos y un conteo para calcular el total de páginas. `SimplePaginate` evita el conteo total y solo comprueba si existe otra página. `CursorPaginate` usa valores ordenados codificados en un cursor opaco, evitando grandes desplazamientos en resultados extensos o que cambian con frecuencia.

```go
invoices, page, err := models.Invoices(db).
	OrderBy("created_at").
	Paginate(ctx, g, 25, requestedPage, pagination.Options{})
if err != nil { return nil, nil, err }
return invoices, page, nil
```

<a name="preserve-filters-and-stable-ordering"></a>
## Conserve los filtros y un orden estable

Componga búsqueda y filtros independientes del tenant antes del terminal, valide tamaño de página en la petición y copie parámetros aprobados a los enlaces. El cursor exige orden determinista con desempate único para que timestamps iguales no dupliquen ni oculten filas entre peticiones.

<a name="reject-malformed-input-and-impossible-positions"></a>
## Rechace entrada malformada y posiciones imposibles

Los terminales de paginación publicados normalizan páginas inferiores a uno a 1. La validación HTTP aún puede rechazar ese valor según su contrato, pero el terminal no devuelve error de página negativa. Tamaño ilimitado, cursores malformados o no decodificables son responsabilidad de petición. Una página posterior al final sigue el contrato documentado de vacío o no encontrado.

<a name="test-boundaries-and-concurrent-change"></a>
## Pruebe límites y cambios concurrentes

Pruebe páginas inferiores a uno y su normalización a 1; luego resultados vacíos, una página exacta, última página llena, una fila extra, tamaño máximo permitido y otro tenant con igual orden. Para cursores, inserte una fila entre peticiones y confirme la continuación documentada, sin confiar en el texto interno del cursor.

<a name="three-terminal-shapes"></a>
## Tres formas de operaciones terminales

Paginate devuelve las filas y una `*pagination.LengthAwarePage` con totales; `SimplePaginate` devuelve las filas y una `*pagination.Page` centrada en anterior/siguiente; `CursorPaginate` recibe un `pagination.Cursor` interpretado y devuelve las filas y una `*pagination.CursorPage`. Los tres reciben contexto, Grant, tamaño de página y opciones tras componer la consulta.

<a name="pagination-is-part-of-the-response-contract"></a>
## La paginación forma parte del contrato de respuesta

Los recursos HTTP deciden cómo serializar elementos, totales y enlaces; la página solo aporta esos datos. La validación de la consulta define los filtros y las columnas de ordenación permitidos. Use Chunk o Cursor para procesar lotes en el servidor sin navegación para el usuario: un objeto de página es innecesario cuando cada registro se procesa una sola vez.
