# Policies y Grant

- [Una policy transforma identidad en autoridad acotada](#a-policy-turns-identity-into-bounded-authority)
- [Lleve un Subject fiable al caso de uso](#carry-a-trusted-subject-to-the-use-case)
- [Autorice la operación y pase su Grant](#authorize-the-operation-and-pass-its-grant)
- [Mantenga la decisión junto a la acción protegida](#keep-the-decision-beside-the-protected-action)
- [Rechace sin revelar otro tenant](#refuse-without-revealing-another-tenant)
- [Nunca reconstruya autoridad desde la entrada de petición](#never-reconstruct-authority-from-request-input)
- [Demuestre rechazo y aislamiento de tenant en almacenamiento](#prove-denial-and-tenant-isolation-at-storage)
- [La identidad inicia la frontera; Grant la completa](#identity-starts-the-boundary-and-grant-finishes-it)

<a name="a-policy-turns-identity-into-bounded-authority"></a>
## Una policy transforma identidad en autoridad acotada

La autenticación produce `security.Subject`: tenant e identidad del actor fiables. La autorización pregunta a una Policy tipada si ese sujeto puede ejecutar una operación sobre un candidato y devuelve `security.Grant` para el repositorio protegido. Una cadena de rol no es autoridad.

<a name="carry-a-trusted-subject-to-the-use-case"></a>
## Lleve un Subject fiable al caso de uso

Cargue Subject mediante middleware de sesión o token y construya policies con los colaboradores necesarios. No cargue un registro de tenant mediante un repositorio sin ámbito solo para decidir si el llamador puede verlo.

- Use el Subject proporcionado por la frontera de autenticación.
- Defina operación y tipo candidato explícitamente.
- Mantenga la consulta por tenant en la policy o repositorio autorizado.
- Defina deliberadamente resultados de visitante y rechazo.

<a name="authorize-the-operation-and-pass-its-grant"></a>
## Autorice la operación y pase su Grant

Llame a `security.Authorize` con policy, sujeto, acción y candidato. Pase solo el Grant devuelto al método de repositorio que lo exige, haciendo imposible omitir autorización en la llamada. Los cambios son un mapa literal de las columnas que esta acción puede escribir, leídas de una solicitud tipada, de modo que ninguna clave enviada llega a la actualización. El adapter aún debe derivar el ámbito mediante `data.Tenant(grant)`.

```go example
const InvoiceUpdate security.Action = "invoice.update"

type Invoice struct{ ID string }

type UpdateInvoiceRequest struct{ Reference, DueOn string }

type Invoices interface {
	Update(ctx context.Context, grant security.Grant, id string, fields map[string]any) error
}

func Update(ctx context.Context, policy security.Policy[Invoice], subject security.Subject, invoices Invoices, invoice Invoice, req UpdateInvoiceRequest) error {
	grant, err := security.Authorize(ctx, policy, subject, InvoiceUpdate, invoice)
	if err != nil {
		return err
	}

	return invoices.Update(ctx, grant, invoice.ID, map[string]any{
		"reference": req.Reference,
		"due_on":    req.DueOn,
	})
}
```

<a name="keep-the-decision-beside-the-protected-action"></a>
## Mantenga la decisión junto a la acción protegida

Use una policy pequeña por recurso de dominio, asigne las operaciones a acciones con nombre estables y mantenga separado el transporte del controller. Reautorice el trabajo en segundo plano y encolado a partir de procedencia duradera en lugar de serializar un Grant de la petición original. Un Grant emitido mientras un operador ve la cuenta de otra persona es de solo lectura: `Grant.ReadOnly` responde true, y toda escritura en repositorio, caché, disco y cola lo rechaza con `auth.ErrReadOnly`, así que deje la acción fuera de una página así en lugar de dejar que falle.

- Autorice antes de lecturas y escrituras protegidas.
- Pase Grant por las fronteras de servicio y repositorio.
- Devuelva la misma ausencia pública para registros ocultos de otro tenant.
- Reconstruya un Subject fiable para trabajo asíncrono.

<a name="refuse-without-revealing-another-tenant"></a>
## Rechace sin revelar otro tenant

Un Subject ausente es no autenticado; un rechazo de policy es prohibido u ocultado como no encontrado; la indisponibilidad del repositorio es operativa. No reduzca eso a un registro vacío ni continúe con Grant cero.

> [!WARNING]
> Un Grant vacío o de otro tenant no puede alcanzar filas del tenant. Un adapter que ignora Grant viola el contrato aunque pasen las pruebas de controller.

<a name="never-reconstruct-authority-from-request-input"></a>
## Nunca reconstruya autoridad desde la entrada de petición

> [!WARNING]
> Nunca acepte tenant, rol, Subject o Grant desde formulario, ruta, query, claim de token sin verificar, payload de job ni argumentos MCP. Nunca almacene autorización en caché más allá del sujeto, acción, recurso y estado de policy que cubrió.

La autenticación y las URLs firmadas demuestran hechos distintos. Una sesión o firma válida aún exige policy para la acción solicitada.

<a name="prove-denial-and-tenant-isolation-at-storage"></a>
## Demuestre rechazo y aislamiento de tenant en almacenamiento

Ejecute un contrato compartido de repositorio con Grants permitidos, denegados, vacíos y ajenos. Ejercite cada acción de policy, visitante, identificadores de otro tenant, cancelación y job que reconstruye identidad por el recorrido fiable de aplicación.

- Un Subject autorizado recibe un Grant acotado.
- Un Subject rechazado no alcanza una operación protegida del adapter.
- Los IDs de otro tenant no revelan existencia.
- Cada repositorio concreto exige Grant.

<a name="identity-starts-the-boundary-and-grant-finishes-it"></a>
## La identidad inicia la frontera; Grant la completa

La autenticación proporciona Subject, las sesiones transportan identidad del navegador, los contratos mantienen Grant en las firmas y los controllers solo traducen transporte.

- [Autenticación](/docs/{{version}}/security/authentication)
- [Sesión](/docs/{{version}}/http/session)
- [Contratos y adapters](/docs/{{version}}/architecture/contracts-adapters)
- [Controllers](/docs/{{version}}/http/controllers)
