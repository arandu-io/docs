# Respuestas HTTP

- [Confirme un formato de respuesta una sola vez](#commit-one-response-shape-once)
- [Devuelva un recurso público, no un model](#return-a-public-resource-not-a-model)
- [Elija el helper según el contrato del cliente](#match-the-helper-to-the-client-contract)
- [No escriba después de confirmar la respuesta](#do-not-write-past-a-committed-response)
- [Verifique detalles del protocolo, no solo cuerpos](#assert-protocol-details-not-only-bodies)
- [Helpers de respuesta publicados](#released-response-helpers)
- [Mantenga la representación en el borde](#keep-representation-at-the-edge)

<a name="commit-one-response-shape-once"></a>
## Confirme un formato de respuesta una sola vez

Una respuesta es un estado explícito, un conjunto de headers y una representación del cuerpo. Los helpers de Context preservan esa decisión para páginas completas, fragmentos, recursos JSON declarados, redirecciones con nombre y estados vacíos, impidiendo escritura duplicada del handler.

<a name="return-a-public-resource-not-a-model"></a>
## Devuelva un recurso público, no un model

Devuelva un `JsonResource` declarado para JSON, no un model sin procesar. `ToArray` define los campos públicos permitidos y With proporciona metadatos superiores; `Context.JSON` elige estado y codificación preservando la frontera del recurso.

```go example
type Invoice struct {
	ID    string
	Total int64
}

type InvoiceResource struct{ invoice Invoice }

func (r InvoiceResource) ToArray() map[string]any {
	return map[string]any{"id": r.invoice.ID, "total": r.invoice.Total}
}
func (InvoiceResource) With() map[string]any { return nil }

func Show(ctx *fhttp.Context, invoice Invoice) error {
	return ctx.JSON(stdhttp.StatusOK, InvoiceResource{invoice: invoice})
}
```

<a name="match-the-helper-to-the-client-contract"></a>
## Elija el helper según el contrato del cliente

Use View para un documento completo, Fragment para un parcial con nombre y `RedirectRoute` para un destino de la tabla de rutas. `RedirectRoute` devuelve 303 a un navegador normal y YPY-Redirect con 204 a una petición mejorada para que ambos naveguen correctamente. Use `EventStream` para actualizaciones unidireccionales en vivo del servidor; `EventStreamWithOptions` agrega comentarios de heartbeat mientras cada evento se descarga inmediatamente.

- Elija el estado antes de escribir el cuerpo.
- Defina headers de respuesta antes de la primera escritura.
- Devuelva fragmentos 422 como el contenido del formulario que rechazan, que YPY renderiza en su lugar.
- Separe recursos JSON de models de base.
- Para SSE, cree eventos con id, event, data, retry o Comment y detenga la producción tan pronto se cancele el contexto de la petición.

<a name="do-not-write-past-a-committed-response"></a>
## No escriba después de confirmar la respuesta

Los errores de codificación, nombres de view desconocidos, nombres de ruta ausentes y escrituras después de una respuesta terminal no son éxitos vacíos. Conserve los significados de 404, 409, 422 y 503 en vez de reducir cada fallo a 200 o 500. Un stream SSE no es reproducción duradera: la cancelación por desconexión detiene al productor y omite deliberadamente un evento de cierre. Si un stream largo debe superar el plazo ordinario de escritura del servidor, borre ese plazo para esta respuesta con `http.NewResponseController` antes del streaming.

> [!WARNING]
> YPY renderiza un 422 dentro del formulario que envió la petición y decide cualquier otro 4xx y 5xx con sus reglas de error. Cambiar un 422 a 200 haría que métricas y logs trataran una entrada rechazada como éxito.

<a name="assert-protocol-details-not-only-bodies"></a>
## Verifique detalles del protocolo, no solo cuerpos

Compruebe estado, tipo de contenido, cabeceras y cuerpo para peticiones normales y mejoradas. Pruebe `RedirectRoute` en ambos modos, claves JSON exactas, ausencia de campos sensibles, una raíz de fragmento 422 y un error de codificación o renderizado. Para SSE, compruebe text/event-stream, no-cache, X-Accel-Buffering: no, encuadre de datos multilínea, validación de metadatos, comentarios de heartbeat, un flush por evento y salida rápida en la cancelación y el apagado de la aplicación.

- La redirección normal es 303 con Location.
- La redirección mejorada es 204 con YPY-Redirect.
- La salida `JsonResource` contiene solo claves permitidas.
- Los fragmentos de validación preservan 422 y son el contenido del formulario.
- La cancelación de SSE detiene la producción sin crear un evento final sintético.

<a name="released-response-helpers"></a>
## Helpers de respuesta publicados

`Context.View` y Fragment renderizan HTML, JSON acepta `hhttp.JsonResource`, `RedirectRoute` invierte una ruta con nombre y Status escribe un estado terminal. `hhttp.Refuse` proporciona errores públicos acotados sin exponer causas internas. `ResponseFactory.EventStream` y `EventStreamWithOptions` construyen respuestas SSE; `hhttp.NewServerSentEvent` y Comment codifican metadatos y datos del evento sin permitir inyección de nueva línea o NUL.

<a name="keep-representation-at-the-edge"></a>
## Mantenga la representación en el borde

Las respuestas REST terminan una petición. SSE mantiene una respuesta HTTP abierta para actualizaciones del servidor al cliente, pero no acepta mensajes del cliente; use `WebSocket` para sesiones bidireccionales y almacenamiento duradero cuando los clientes deban reproducir hechos perdidos. Resources define campos JSON, views define HTML, controllers eligen la representación y el manejo centralizado de errores atiende fallos que escapan de la action.

- [Patrones de comunicación](/docs/{{version}}/architecture/contracts-adapters)
- [`WebSocket` de Joaju](/docs/{{version}}/async/joaju)
- [Recursos](/docs/{{version}}/http/resources)
- [Errores](/docs/{{version}}/http/errors)
