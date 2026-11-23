# Recursos HTTP

- [El formato de respuesta es un tipo declarado](#the-response-shape-is-a-declared-type)
- [Complete primero autorización y validación](#finish-authorization-and-validation-first)
- [Implemente una vez el contrato de recurso JSON](#implement-the-json-resource-contract-once)
- [Modele envoltorios de colección explícitamente](#model-collection-envelopes-explicitly)
- [Un formulario rechazado sigue siendo una petición rechazada](#a-rejected-form-remains-a-rejected-request)
- [Los campos públicos forman una lista explícita de permisos](#public-fields-are-an-allowlist)
- [Verifique estado, formato y fragmento](#assert-status-shape-and-fragment)
- [Los recursos están entre dominio y transporte](#resources-sit-between-domain-and-transport)

<a name="the-response-shape-is-a-declared-type"></a>
## El formato de respuesta es un tipo declarado

Un recurso HTTP es el tipo que el endpoint promete devolver. El compilador verifica tipo y constructor, pero `JsonResource.ToArray` devuelve `map[string]`any: las claves siguen siendo strings de runtime. Las pruebas exactas de formato deben detectar claves incorrectas, eliminadas o añadidas por accidente. El recurso explícito mantiene campos públicos independientes de las columnas del model.

<a name="finish-authorization-and-validation-first"></a>
## Complete primero autorización y validación

Construya recursos solo con datos ya cargados bajo Grant y entrada validada en tipo propio. El serializador no corrige autorización ni interpreta entrada. Defina estado y negociación de contenido antes de renderizar para que clientes de página y JSON observen el mismo resultado.

<a name="implement-the-json-resource-contract-once"></a>
## Implemente una vez el contrato de recurso JSON

`Context.JSON` acepta `JsonResource`: `ToArray` devuelve campos que pueden salir y With, metadatos superiores. Mantenga el model en un campo no exportado, construya el recurso desde un valor ya autorizado y responda directamente desde el handler.

```go example
type Invoice struct {
	ID, Reference string
	Total         int64
}

type InvoiceResource struct{ invoice Invoice }

func NewInvoiceResource(invoice Invoice) InvoiceResource {
	return InvoiceResource{invoice: invoice}
}

func (r InvoiceResource) ToArray() map[string]any {
	return map[string]any{
		"id": r.invoice.ID, "reference": r.invoice.Reference, "total": r.invoice.Total,
	}
}

func (InvoiceResource) With() map[string]any { return nil }

func Show(ctx *fhttp.Context, invoice Invoice) error {
	return ctx.JSON(stdhttp.StatusOK, NewInvoiceResource(invoice))
}
```

<a name="model-collection-envelopes-explicitly"></a>
## Modele envoltorios de colección explícitamente

`Context.JSON` siempre coloca `ToArray` dentro de la clave superior data. En :hesapeVersion, `ResourceCollection.ToArray` ya devuelve su propia data; pasarlo a `Context.JSON` genera `data.data`. Defina un recurso de colección o formato de respuesta explícito y fije el objeto completo en una prueba, sin suponer un envoltorio simple. Para página parcial, renderice un fragmento separado cuya raíz corresponda al destino del intercambio.

> [!NOTE]
> Trate data como metadato reservado. `Context.JSON` mezcla With tras crear el payload, de modo que `With()["data"]` sobrescribe el cuerpo del recurso en lugar de rechazarse.

<a name="a-rejected-form-remains-a-rejected-request"></a>
## Un formulario rechazado sigue siendo una petición rechazada

Devuelva 422 ante fallo de validación aunque el navegador inserte el fragmento en la página. YPY renderiza un 422 dentro del formulario que envió la petición, conserva lo escrito y da foco al primer campo inválido, así que el estado no necesita configuración. Devolver 200 hace que logs y métricas registren éxito y trata un formulario rechazado como éxito.

<a name="public-fields-are-an-allowlist"></a>
## Los campos públicos forman una lista explícita de permisos

Nunca serialice un model solo porque sus campos actuales parecen inofensivos: una columna secreta futura entraría en el payload sin cambiar el endpoint. El mapa `JsonResource` hace deliberadas las adiciones, pero no verifica claves string al compilar. Separe metadatos With de data y no exponga claves internas de tenant, digests, rutas de storage ni detalles de policy.

<a name="assert-status-shape-and-fragment"></a>
## Verifique estado, formato y fragmento

Pruebe estado exacto y claves JSON para la API, luego elemento raíz y mensajes de campo del formulario parcial. Verifique la ausencia de un campo sensible del model. La prueba del recurso debe fallar si el estado se vuelve 200 indebidamente, cambia una clave o el fragmento deja de apuntar al formulario sustituido.

<a name="resources-sit-between-domain-and-transport"></a>
## Los recursos están entre dominio y transporte

Los validadores definen entrada aceptada, los models definen datos almacenados, las vistas definen HTML y los recursos definen salida. La separación es intencional incluso con campos iniciales parecidos. Los recursos de paginación llevan estado de página; el manejo de errores posee la representación estable del problema cuando no hay recurso de éxito.
