# Serialización de models

- [El JSON predeterminado es práctico, pero no define la API](#default-json-is-convenient-and-intentionally-not-the-api)
- [Mapee un model a una respuesta explícita](#map-a-model-into-an-explicit-response)
- [Represente colecciones y relaciones opcionales](#represent-collections-and-optional-relations)
- [Deténgase ante valores no compatibles o mapeo incorrecto](#stop-on-unsupported-values-or-broken-mapping)
- [Fije claves públicas y ausencia de campos sensibles](#lock-the-public-keys-and-sensitive-absences)
- [Los DTO y `JsonResource` tienen contratos diferentes](#dtos-and-jsonresource-have-different-contracts)
- [Los logs necesitan una representación segura propia](#logging-needs-its-own-safe-representation)

> [!WARNING]
> **Limitaciones actuales.** Existen representación de model publicada y DTO oficial con campos permitidos, pero ningún recorrido etiquetado demuestra `ToArray` o `ToJSON` del model como contrato público. Prefiera un DTO o `JsonResource` explícito y pruebe sus claves exactas.

<a name="default-json-is-convenient-and-intentionally-not-the-api"></a>
## El JSON predeterminado es práctico, pero no define la API

`encoding/json` serializa campos Go exportados salvo tag o marshaler personalizado. Ese valor predeterminado sirve para fronteras fiables, pero es demasiado amplio para un endpoint público: añadir una columna exportada puede cambiar la salida sin editar el handler ni el contrato de respuesta.

<a name="map-a-model-into-an-explicit-response"></a>
## Mapee un model a una respuesta explícita

Declare los campos prometidos por el endpoint y construya ese tipo con un model ya autorizado. La lista es un DTO Go ordinario serializado mediante `encoding/json`, no el contrato `hhttp.JsonResource` del framework. Reutilice IDs estables y valores de dominio omitiendo tenant, hashes, rutas de storage, timestamps internos y relaciones innecesarias para el cliente.

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
## Represente colecciones y relaciones opcionales

Mapee cada model de la colección y envuelva elementos con datos de paginación cuando el endpoint sea paginado. Incluya un recurso relacionado solo si la relación se cargó deliberadamente y el contrato la exige. Use omitempty solo cuando la ausencia tenga significado definido; no sustituye la elección de campos públicos.

<a name="stop-on-unsupported-values-or-broken-mapping"></a>
## Deténgase ante valores no compatibles o mapeo incorrecto

La codificación JSON puede fallar con valores no admitidos, y las implementaciones de `MarshalJSON` pueden devolver errores. La ausencia de un valor obligatorio del dominio debe interrumpir el mapeo, no generar una respuesta engañosa. Si falla el constructor del recurso, no serialice el modelo como alternativa: volvería a exponer los campos que el recurso debe proteger.

<a name="lock-the-public-keys-and-sensitive-absences"></a>
## Fije claves públicas y ausencia de campos sensibles

Verifique la forma exacta de la respuesta o el conjunto de claves decodificadas y confirme explícitamente la ausencia de campos sensibles del modelo. Añada al modelo o a la entrada del mapeo un campo exportado exclusivo de la fixture y demuestre que no aparece automáticamente. Cubra relaciones opcionales cargadas y no cargadas, además de los errores de codificación.

<a name="dtos-and-jsonresource-have-different-contracts"></a>
## Los DTO y `JsonResource` tienen contratos diferentes

`ToArray` de una fila o de una colección, los controles de visibilidad, los campos ocultos y los atributos añadidos definen la representación en array del paquete model. Un DTO común sigue las etiquetas de `encoding/json` y `MarshalJSON`. En cambio, `Context.JSON` acepta `hhttp.JsonResource`, cuyos métodos `ToArray` y With devuelven mapas en tiempo de ejecución; sus claves exigen pruebas exactas de respuesta. Elija explícitamente el contrato de transporte, sin tratar ambas formas como el mismo recurso.

<a name="logging-needs-its-own-safe-representation"></a>
## Los logs necesitan una representación segura propia

Una lista de campos públicos no elimina automáticamente datos sensibles de logs, eventos ni volcados. Los tipos con credenciales o claves necesitan comportamiento seguro de logging y payloads deliberados. El cifrado protege texto almacenado, pero no vuelve el texto cifrado ni sus metadatos apropiados para toda respuesta pública.
