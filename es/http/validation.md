# Validación

- [La validación restringe entrada; la policy concede autoridad](#validation-narrows-input-while-policy-grants-authority)
- [Declare la entrada junto al caso de uso](#declare-the-input-beside-its-use-case)
- [Rechace una creación antes de llegar al almacenamiento](#reject-one-create-request-before-storage)
- [Coloque cada regla en la frontera que conoce](#put-each-rule-at-the-boundary-it-understands)
- [Preserve errores de campos sin filtrar causas](#preserve-field-errors-without-leaking-causes)
- [Defina campos permitidos y limite trabajo costoso](#allowlist-fields-and-expensive-work)
- [Pruebe el menor valor inválido y el límite](#test-the-smallest-failing-value-and-the-boundary)
- [Continúe de datos válidos a trabajo autorizado](#continue-from-valid-data-to-authorized-work)

<a name="validation-narrows-input-while-policy-grants-authority"></a>
## La validación restringe entrada; la policy concede autoridad

La validación transforma valores no fiables del transporte en entrada tipada y errores indexados por campo. Verifica formato y restricciones neutrales de negocio; autorización y unicidad en base siguen siendo operaciones separadas con su propia semántica de fallo.

<a name="declare-the-input-beside-its-use-case"></a>
## Declare la entrada junto al caso de uso

Mantenga explícitos nombres de transporte, campos de dominio y tamaños máximos. No añada `TenantID` a la entrada del navegador: el repositorio deriva el tenant del Grant.

<a name="reject-one-create-request-before-storage"></a>
## Rechace una creación antes de llegar al almacenamiento

Use Bind para leer la solicitud en un valor tipado cuyas etiquetas form nombran cada campo aceptado y cuyo método Rules declara lo que cada campo puede ser. `ctx.Bind` ejecuta las reglas sobre la entrada sin procesar antes de convertir nada y devuelve el fallo como un `validation.Errors` no vacío, que el adapter de acción convierte en flash y redirección para una página y en un documento de problema 422 para un cliente JSON. Bind lee solo el cuerpo en POST, PUT, PATCH y DELETE, así que un valor en la query string no satisface una regla ni llega al registro guardado. Normalice correo, identificadores y espacios antes de persistir valores aceptados.

```go example
type RegisterRequest struct {
	Email string `form:"email"`
}

func (RegisterRequest) Rules() validation.Rules {
	return validation.Rules{
		"email": "required|email|max:254",
	}
}

func Register(ctx *fhttp.Context) error {
	var req RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}
	return nil
}
```

<a name="put-each-rule-at-the-boundary-it-understands"></a>
## Coloque cada regla en la frontera que conoce

Declare en Rules las reglas de obligatoriedad, tipo, tamaño, pertenencia y comparación, como texto con barras verticales o como lista por campo. Lo que abarca varios campos o consulta los datos pertenece a un método Validate de la solicitud, que el servicio llama; no repita en Validate una restricción que Rules ya declara. Los mensajes salen en el idioma de la solicitud, y los nombres de los campos se mantienen estables para colocar cada error junto a su entrada.

- Valide antes de autorizar cuando la entrada malformada no deba llegar a la policy.
- Autorice antes de que una consulta protegida revele la existencia del registro.
- Devuelva consistentemente JSON 422 o flash con redirección de formulario.
- Separe conflictos de persistencia de errores de sintaxis.

<a name="preserve-field-errors-without-leaking-causes"></a>
## Preserve errores de campos sin filtrar causas

Reglas desconocidas, argumentos inválidos, errores devueltos vacíos y coerción que oculta entrada malformada son defectos de implementación. La indisponibilidad de base nunca es un mensaje de validación de campo.

> [!WARNING]
> No consulte unicidad sin el tenant actual y la frontera de autorización. Una regla aparentemente global puede revelar datos de otro tenant o rechazar un valor válido en el ámbito actual.

<a name="allowlist-fields-and-expensive-work"></a>
## Defina campos permitidos y limite trabajo costoso

> [!WARNING]
> Defina campos permitidos y limite bytes de cuerpo, elementos de colecciones, strings y uploads antes del trabajo costoso. Nunca acepte selector de tenant, ruta de storage, campo de model ni hash de contraseña desde la entrada; mantenga credenciales fuera de mensajes, logs y datos antiguos de formulario en flash.

<a name="test-the-smallest-failing-value-and-the-boundary"></a>
## Pruebe el menor valor inválido y el límite

Pruebe en tabla cada límite de regla con nombre real de campo y clave de mensaje; envíe después una petición inválida mediante clientes normales y mejorados. Confirme que la entrada rechazada no cambia filas y que la válida llega normalizada al servicio.

- Los límites de obligatoriedad y máximo son exactos.
- La confirmación entre campos usa nombres de campos aceptados.
- Las respuestas 422 siguen siendo respuestas de rechazo.
- Ninguna consulta protegida se ejecuta fuera del ámbito del tenant.

<a name="continue-from-valid-data-to-authorized-work"></a>
## Continúe de datos válidos a trabajo autorizado

El parsing proporciona valores candidatos, los controllers eligen el caso de uso, las policies deciden autoridad y las respuestas renderizan errores de campos sin cambiar su estado.

- [Peticiones](/docs/{{version}}/http/requests)
- [Controllers](/docs/{{version}}/http/controllers)
- [Respuestas](/docs/{{version}}/http/responses)
