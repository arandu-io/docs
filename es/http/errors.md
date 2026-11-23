# Manejo de errores

- [Un handler clasifica, registra y renderiza](#one-handler-classifies-reports-and-renders)
- [Construya el handler con la configuración cargada](#build-the-handler-from-loaded-configuration)
- [Instale recuperación como middleware más externo](#install-recovery-as-the-outermost-middleware)
- [Devuelva un fallo HTTP deliberado](#return-an-intentional-http-failure)
- [Los errores desconocidos siguen siendo errores de servidor](#unknown-errors-stay-server-errors)
- [El detalle público y la causa privada son campos distintos](#public-detail-and-private-cause-are-different-fields)
- [Demuestre representación y registro del error](#prove-both-representation-and-report)
- [Recupere el servicio, no solo la petición](#recover-the-service-not-only-the-request)
- [Siga los fallos a través de las fronteras](#follow-failures-across-their-boundaries)

> [!WARNING]
> **Limitaciones actuales.** El middleware de recuperación y la página de error de desarrollo están publicados. La clasificación de errores de dominio sigue siendo contrato de la aplicación; no infiera estados mediante cadenas arbitrarias de error.

<a name="one-handler-classifies-reports-and-renders"></a>
## Un handler clasifica, registra y renderiza

Los errores atraviesan una frontera central de recuperación. Los fallos esperados de dominio y transporte conservan significado; los panics escapados y errores internos desconocidos se registran una vez y se vuelven respuesta de producción acotada, sin stack ni secretos.

<a name="build-the-handler-from-loaded-configuration"></a>
## Construya el handler con la configuración cargada

Create construye el único handler de excepciones a partir de la configuración cargada, lee la ruta de módulo del binario para distinguir los frames de código de la aplicación de los del framework e instala Recover en el lado más externo después de las frases de rechazo. La página de depuración solo existe en desarrollo; la validación de producción la mantiene desactivada. La aplicación registra sus propias respuestas en ese handler en `WithExceptions`.

<a name="install-recovery-as-the-outermost-middleware"></a>
## Instale recuperación como middleware más externo

Devuelva errores desde las acciones y deje que la recuperación elija la respuesta configurada. En desarrollo, la página de error puede mostrar fragmentos de código; producción envía estado genérico y contexto de correlación mientras logs estructurados conservan la causa encadenada.

```go example
var ErrQuotaExceeded = errors.New("quota exceeded")

func WithErrorHandling(drivers *foundation.Drivers) *foundation.ApplicationBuilder {
	return foundation.Configure(drivers).
		WithExceptions(func(e *exception.Handler) {
			e.DontReport(ErrQuotaExceeded)
			e.ShouldRenderJSONWhen(func(r *stdhttp.Request, _ error) bool {
				return strings.HasPrefix(r.URL.Path, "/api/")
			})
		})
}

func Register(router *fhttp.Router, show func(*fhttp.Context) error) {
	router.Action("GET", "/invoices/{invoice}", show)
}
```

<a name="return-an-intentional-http-failure"></a>
## Devuelva un fallo HTTP deliberado

Defina tipos de error de la aplicación para ausencia, denegación, conflicto, estado inválido e indisponibilidad de dependencias. El adaptador del router asigna un error devuelto, solo o envuelto con %w: un centinela de no encontrado responde 404, `security.ErrForbidden` 403, `security.ErrCSRF` 419, `database.ErrUniqueViolation` 409, y un error con el método `HTTPStatus()` int responde ese estado. Dé ese método a un error de dominio en lugar de tratarlo en cada controller, envuelva las causas con contexto y evite inferir estados comparando cadenas de error arbitrarias.

- Devuelva errores de validación mediante el adapter dedicado.
- Use 404 sin revelar existencia en otro tenant.
- Use 409 para un conflicto real de estado.
- Use 503 ante indisponibilidad de una dependencia obligatoria.

<a name="unknown-errors-stay-server-errors"></a>
## Los errores desconocidos siguen siendo errores de servidor

Escribir parte de la respuesta antes de devolver error puede impedir cambiar el estado final. Un panic tras enviar headers se registra, pero no permite una sustitución limpia; complete el trabajo que puede fallar antes de confirmar la salida.

> [!WARNING]
> Las páginas de código y volcados de desarrollo pueden contener secretos. Mantenga la consola desactivada o protegida fuera de desarrollo y elimine credenciales, cookies, Grants y URLs de base de los valores estructurados.

<a name="public-detail-and-private-cause-are-different-fields"></a>
## El detalle público y la causa privada son campos distintos

> [!WARNING]
> Separe el detalle público de la causa privada encadenada. Nunca exponga SQL, URLs de conexión, claves, tokens, cookies, cuerpos de petición, stack traces ni detalles de policy; mantenga debug desactivado donde puedan acceder personas no fiables.

<a name="prove-both-representation-and-report"></a>
## Demuestre representación y registro del error

Ejercite cada error esperado mediante la ruta registrada y provoque panic antes y después de una escritura. Confirme que los cuerpos de producción no revelan la causa, el recorder observa un evento y un rechazo de dominio no se vuelve fallo interno.

- Los fallos esperados preservan su estado.
- Los errores desconocidos se registran una vez.
- Las respuestas de producción no contienen stack ni credenciales.
- La ausencia por otro tenant es indistinguible de un registro inexistente.

<a name="recover-the-service-not-only-the-request"></a>
## Recupere el servicio, no solo la petición

Use el ID de respuesta para localizar log y trace correlacionados, inspeccione diagnósticos del kernel y salud de módulos dependientes y corrija o revierta el release. Repita una mutación solo tras conocer su resultado de commit y garantizar idempotencia.

<a name="follow-failures-across-their-boundaries"></a>
## Siga los fallos a través de las fronteras

La observabilidad registra el evento, la validación posee errores de campos, los controllers devuelven fallos y el despliegue decide si pueden habilitarse diagnósticos de desarrollo.

- [Logs y observabilidad](/docs/{{version}}/operations/logging-observability)
- [Controllers](/docs/{{version}}/http/controllers)
- [Validación](/docs/{{version}}/http/validation)
