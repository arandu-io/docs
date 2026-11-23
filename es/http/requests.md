# Peticiones HTTP

- [Lea datos del transporte sin inventar identidad](#read-transport-data-without-inventing-identity)
- [Monte un contrato tipado de entrada](#assemble-one-typed-input-contract)
- [Elija el lector que preserva el significado](#choose-the-reader-that-preserves-meaning)
- [Limite el cuerpo antes de decodificar](#bound-the-body-before-decoding-it)
- [Ejercite cada canal de entrada y límite](#exercise-every-input-channel-and-limit)
- [Superficies de petición publicadas](#released-request-surfaces)
- [Pase del transporte a una operación fiable](#move-from-transport-to-a-trusted-operation)

<a name="read-transport-data-without-inventing-identity"></a>
## Lea datos del transporte sin inventar identidad

La frontera preferida es Context reducido del framework con primitivas `net/http`. Lea parámetros de ruta, query, headers y cuerpo acotado solo donde la acción lo necesite; el Request amplio de Hesape es compatibilidad, no dependencia obligatoria del controller.

<a name="assemble-one-typed-input-contract"></a>
## Monte un contrato tipado de entrada

Use Bind para leer solo los campos aceptados en un tipo de entrada dedicado, cuyas etiquetas form nombran cada campo que puede llegar, normalice valores canónicos y pase solo datos validados al servicio. Un tipo con método Rules tiene esas reglas ejecutadas primero sobre la entrada sin procesar. Bind escribe solo campos con etiqueta, recorta cada valor e informa un valor que no se convierte como error del campo, en lugar de un cero silencioso. Preserve `ctx.Ctx` para mantener cancelación e identidad autenticada en llamadas posteriores.

```go example
type CreateInvoiceInput struct {
	Reference string `form:"reference"`
	Amount    int64  `form:"amount"`
}

func Create(ctx *fhttp.Context) error {
	var input CreateInvoiceInput
	if err := ctx.Bind(&input); err != nil {
		return err
	}
	return nil
}
```

<a name="choose-the-reader-that-preserves-meaning"></a>
## Elija el lector que preserva el significado

Use middleware de proxy fiable antes de leer una dirección de cliente derivada del par. Limite cuerpos antes de decodificar, cierre streams de upload y separe identificadores de ruta de campos escribibles, impidiendo que formularios sustituyan propiedad o tenant. Input y Bind leen un único mapa: la query string en GET y HEAD, y solo el cuerpo en POST, PUT, PATCH y DELETE, así que un valor en la URL de la acción se lee explícitamente con Query.

- Use Param para segmentos coincidentes de ruta.
- Use Query para controles opcionales de URL.
- Use Bind para campos de formulario declarados; los métodos inseguros leen solo el cuerpo.
- Use Request estándar para headers, cuerpo y contexto.

<a name="bound-the-body-before-decoding-it"></a>
## Limite el cuerpo antes de decodificar

Números malformados, tipos de medio no compatibles, cuerpos excesivos, valores obligatorios ausentes y lecturas canceladas deben seguir explícitos. No convierta silenciosamente un valor inválido en cero ni continúe tras detenerse el decoder en un payload parcial.

> [!WARNING]
> Los headers reenviados son entrada del llamador hasta que `TrustProxies` acepta el par inmediato. Nunca derive throttles de autenticación, identidad de auditoría ni enlaces firmados de reenvíos no fiables.

<a name="exercise-every-input-channel-and-limit"></a>
## Ejercite cada canal de entrada y límite

Envíe peticiones válidas, ausentes, malformadas, excesivas y canceladas por la ruta real. Verifique la entrada tipada enviada al servicio y demuestre que los campos de tenant o propiedad del llamador se ignoran o rechazan.

- Solo la entrada declarada llega al servicio.
- Los límites de cuerpo y upload actúan antes de una asignación sin techo.
- La cancelación detiene el parsing y el trabajo posterior.
- Los valores de reenvío no fiables no sustituyen la dirección del par.

<a name="released-request-surfaces"></a>
## Superficies de petición publicadas

`http.Context` del framework expone Ctx, Param, Query, Input, Bind, User, helpers de respuesta y request y writer subyacentes. `net/http` estándar posee headers, cuerpos, cookies y contexto. El Request más amplio de Hesape sigue disponible cuando se necesiten sus helpers de compatibilidad.

<a name="move-from-transport-to-a-trusted-operation"></a>
## Pase del transporte a una operación fiable

La validación transforma entrada del transporte en entrada tipada que satisface reglas declaradas. Los controllers eligen el caso de uso, los proxies fiables establecen la frontera de red y las respuestas codifican el resultado.

- [Validación](/docs/{{version}}/http/validation)
- [Controllers](/docs/{{version}}/http/controllers)
- [Respuestas](/docs/{{version}}/http/responses)
