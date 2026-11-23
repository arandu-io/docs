# Controllers

- [Traduzca HTTP y mantenga la policy en el servicio posterior](#translate-http-and-keep-policy-downstream)
- [Inyecte el caso de uso](#inject-the-use-case)
- [Renderice un registro autenticado](#render-one-authenticated-record)
- [Mantenga acciones breves y explícitas](#keep-actions-short-and-explicit)
- [Preserve el significado del fallo](#preserve-failure-meaning)
- [El controller no es autoridad de autorización](#the-controller-is-not-an-authorization-authority)
- [Pruebe transporte y policy por separado](#test-transport-and-policy-separately)
- [Complete la frontera del controller](#complete-the-controller-boundary)

> [!WARNING]
> **Limitaciones actuales.** Los controllers son dependencias ordinarias y explícitas. El framework no resuelve constructores ni autoriza registros automáticamente; cada acción protegida debe obtener y pasar su Grant.

<a name="translate-http-and-keep-policy-downstream"></a>
## Traduzca HTTP y mantenga la policy en el servicio posterior

Un controller lee entrada del transporte, solicita un resultado al servicio de aplicación y elige una respuesta HTTP. No emite Grants, abre otro recorrido de base ni confía en campos de tenant y propiedad del navegador.

<a name="inject-the-use-case"></a>
## Inyecte el caso de uso

Defina la dependencia en el paquete consumidor e inyecte mediante el constructor del controller. Proteja la ruta con middleware de autenticación cuando exija sujeto: `RequireAuth` coloca en el contexto de la solicitud el sujeto que cargó de la sesión, así que `ctx.User` lo informa detrás de la guarda, y `middleware.LoadSubject` hace lo mismo en una ruta pública, dejando pasar a un visitante sin sujeto. El servicio o la policy siguen siendo responsables de autorizar la acción en el registro solicitado.

- Construya los controllers en un provider en app/Providers.
- Registre acciones en rutas con nombre.
- Separe entrada de validación de models persistidos.
- Deje autorización y persistencia al servicio.

<a name="render-one-authenticated-record"></a>
## Renderice un registro autenticado

El plazo del contexto y el Subject autenticado atraviesan la frontera del servicio; el identificador de ruta sigue siendo entrada no fiable, no prueba de acceso. Renderice el resultado ya autorizado mediante vista o recurso declarado.

```go example
type Invoice struct{ ID string }

type InvoiceFinder interface {
	Find(context.Context, security.Subject, string) (Invoice, error)
}

type InvoiceController struct{ invoices InvoiceFinder }
type InvoiceShowData struct{ Invoice Invoice }

func (c InvoiceController) Show(ctx *fhttp.Context) error {
	subject, ok := ctx.User()
	if !ok {
		hhttp.Refuse(ctx.Response, ctx.Request, stdhttp.StatusUnauthorized, "sign in required")
		return nil
	}
	invoice, err := c.invoices.Find(ctx.Ctx(), subject, ctx.Param("invoice"))
	if err != nil {
		return err
	}
	return ctx.View("invoices.show", InvoiceShowData{Invoice: invoice})
}
```

<a name="keep-actions-short-and-explicit"></a>
## Mantenga acciones breves y explícitas

Use Param, Query e Input solo para montar entrada tipada. Devuelva errores operativos a la recuperación central, `validation.Errors` no vacío al recorrido de rechazo y exactamente uno entre View, Fragment, JSON, `RedirectRoute` o Status para responder.

- Pase `ctx.Ctx` al caso de uso.
- Lea identidad mediante la indicación de presencia del contexto.
- Devuelva solo errores de validación no vacíos.
- Elija un helper de respuesta tras el resultado del servicio.

<a name="preserve-failure-meaning"></a>
## Preserve el significado del fallo

Entrada malformada, registros ausentes, acciones denegadas, conflictos y fallos internos son resultados distintos. Convertir todo en 500 borra comportamiento útil; ocultar un error puede devolver 200 con cuerpo vacío por accidente.

> [!WARNING]
> El adapter de acción provoca panic cuando `validation.Errors` devuelto está vacío. Use if `errs.Any()` antes de devolverlo y devuelva nil cuando la validación pase.

<a name="the-controller-is-not-an-authorization-authority"></a>
## El controller no es autoridad de autorización

> [!WARNING]
> Nunca derive tenant de formulario, query ni ruta; nunca cargue datos de tenant sin Subject o Grant; nunca serialice íntegramente un model de persistencia con hashes, claves internas de tenant o campos de auditoría.

Los cuerpos de petición y structs que pueden contener credenciales exigen eliminación explícita en JSON y slog antes de registrarse. La comodidad del controller no autoriza exponer entrada del transporte.

<a name="test-transport-and-policy-separately"></a>
## Pruebe transporte y policy por separado

Use httptest mediante la ruta registrada para verificar estado, cuerpo o redirección, parámetros, rechazo de autenticación y error de servicio. En una prueba unitaria específica, inyecte un finder fake y confirme que recibe el mismo contexto, Subject e ID.

- Un visitante se rechaza antes del caso de uso.
- El Subject y el ID exactos llegan al servicio.
- La validación sigue el contrato de 422 o redirección.
- El comportamiento de la policy sigue probado fuera del controller.

<a name="complete-the-controller-boundary"></a>
## Complete la frontera del controller

Los tipos de petición poseen entrada aceptada, las respuestas poseen semántica de transporte, los recursos poseen campos JSON públicos y la autorización posee Grant. Un controller breve hace visibles esos contratos.

- [Peticiones](/docs/{{version}}/http/requests)
- [Respuestas](/docs/{{version}}/http/responses)
- [Autorización](/docs/{{version}}/security/authorization)
- [Recursos](/docs/{{version}}/http/resources)
