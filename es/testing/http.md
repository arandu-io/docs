# Pruebas HTTP

- [Despache el pipeline registrado con httptest](#dispatch-the-registered-pipeline-with-httptest)
- [Inicie la composición real y controle los adapters externos](#boot-the-real-graph-and-keep-external-adapters-controlled)
- [Envíe una petición por el Handler](#send-one-request-through-handler)
- [Verifique protocolo, sesión, YPY y recorridos de policy](#assert-protocol-session-ypy-and-policy-paths)
- [No deje que la recuperación de producción oculte un defecto de la prueba](#do-not-let-production-recovery-hide-a-test-defect)
- [Mantenga credenciales y datos de tenants fuera de fixtures y salidas](#keep-credentials-and-tenant-data-out-of-fixtures-and-output)
- [Cubra la matriz completa de rutas](#cover-the-complete-route-matrix)
- [Lleve los detalles de servicios y adapters a contratos específicos](#move-service-and-adapter-details-to-focused-contracts)

> [!WARNING]
> **Limitaciones actuales.** El harness de peticiones ejercita el router y el middleware reales. `test.NewClient` guarda cookies y el token CSRF, y su Response verifica estado, texto y redirecciones; verifique el estado protegido directamente junto a él.

<a name="dispatch-the-registered-pipeline-with-httptest"></a>
## Despache el pipeline registrado con httptest

Una prueba de feature HTTP despacha `foundation.Application.Handler` mediante `net/http/httptest`. Atraviesa middleware global, coincidencia de rutas, middleware de ruta, controller, servicio, renderizador y respuesta de error sin abrir un listener.

<a name="boot-the-real-graph-and-keep-external-adapters-controlled"></a>
## Inicie la composición real y controle los adapters externos

Construya los mismos módulos y configuración del servidor, sustituya solo adapters externos explícitos, compile vistas y use una base y un backend de sesión dedicados. Preserve el orden de los middlewares de producción.

- Registre rutas y renderizador reales.
- Use una base de prueba con fixtures de tenants.
- Capture correo y HTTP saliente.
- Elija sesión en memoria o compartida según el comportamiento probado.

<a name="send-one-request-through-handler"></a>
## Envíe una petición por el Handler

Cree una petición estándar, despáchela al Handler de la aplicación y verifique estado, cabeceras y cuerpo. Cuando el recorrido abarca inicio de sesión, CSRF, rotación y cierre de sesión, `test.NewClient(t, app.Handler())`, de `framework/test`, es el mismo despacho con un tarro de cookies, el token CSRF leído de la última página cargada y sin seguir redirecciones. Su `*test.Response` verifica con `AssertStatus`, `AssertOk`, `AssertSee`, `AssertDontSee`, `AssertRedirect` y los métodos Assert para cabeceras, cookies y JSON. `test.Serve` levanta la aplicación en un puerto libre sobre una base propia de la prueba, y `test.As` inicia sesión como un sujeto mediante la sesión de la propia aplicación.

```go example
func assertHomeIsOK(t *testing.T, app *foundation.Application) {
	request := httptest.NewRequest(stdhttp.MethodGet, "/invoices/inv_123", nil)
	response := httptest.NewRecorder()
	app.Handler().ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}
}
```

<a name="assert-protocol-session-ypy-and-policy-paths"></a>
## Verifique protocolo, sesión, YPY y recorridos de policy

Pruebe representación normal y YPY-Request, redirecciones con nombre, recursos JSON, validación de formularios, CSRF, proxy fiable, autenticación, rechazo de policy y cancelación. Verifique el resultado del servicio o la persistencia además del texto de respuesta.

- Verifique Content-Type y el estado exacto.
- Mantenga cookies a lo largo del recorrido real.
- Pruebe la página completa y el fragmento por separado.
- Inspeccione los efectos tras la respuesta.

<a name="do-not-let-production-recovery-hide-a-test-defect"></a>
## No deje que la recuperación de producción oculte un defecto de la prueba

Un panic debe llegar a la frontera de la prueba cuando se configure para exponer defectos; no apruebe una prueba solo porque la recuperación de producción convirtió el panic en 500. Escrituras parciales, errores ocultos, rutas desconocidas y vistas generadas desactualizadas necesitan casos explícitos.

> [!WARNING]
> Un fake unitario del controller no demuestra middleware, nombres de rutas, cookies de sesión, CSRF, registro del renderizador ni mapeo de errores. Mantenga al menos un recorrido completo de ruta por cada flujo público.

<a name="keep-credentials-and-tenant-data-out-of-fixtures-and-output"></a>
## Mantenga credenciales y datos de tenants fuera de fixtures y salidas

> [!WARNING]
> Use credenciales y tenants sintéticos, oculte datos sensibles en las salidas de fallo, nunca llame a proveedores reales y verifique que otro tenant reciba la misma ausencia acotada. Pruebe Secure, `HttpOnly`, `SameSite`, firma y atributos CSRF como detalles del protocolo.

Un helper que omite middleware puede facilitar las pruebas y dejar la producción incorrecta; autentique mediante la sesión publicada o construya deliberadamente un Subject solo en pruebas de servicio.

<a name="cover-the-complete-route-matrix"></a>
## Cubra la matriz completa de rutas

Cubra método y ruta incorrectos, petición válida, entrada malformada, visitante, Subject rechazado, otro tenant, YPY, preferencia JSON, dependencia indisponible, panic, cancelación de petición y rotación de sesión. Ejecute la suite con detección de race.

- Cada ruta tiene recorridos de éxito y rechazo.
- Los registros de otro tenant permanecen invisibles.
- Las peticiones mejoradas mantienen el estado correcto.
- La cancelación llega a las llamadas posteriores.

<a name="move-service-and-adapter-details-to-focused-contracts"></a>
## Lleve los detalles de servicios y adapters a contratos específicos

Los controllers traducen el transporte, el ciclo de petición explica el pipeline, los dobles capturan adapters externos y las pruebas de base demuestran aislamiento en almacenamiento.

- [Controllers](/docs/{{version}}/http/controllers)
- [Ciclo de vida de la petición](/docs/{{version}}/start/request-lifecycle)
- [Dobles de prueba](/docs/{{version}}/testing/test-doubles)
- [Pruebas de base de datos](/docs/{{version}}/testing/database)
