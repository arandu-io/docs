# Rutas

- [Las rutas son contratos ejecutables](#routes-are-executable-contracts)
- [Registre y genere el camino de una ruta](#register-and-reverse-one-route)
- [Grupos, restricciones y raíz exacta](#groups-constraints-and-the-exact-root)
- [Falle durante la construcción, no después de un clic](#fail-while-building-not-after-a-click)
- [Inspeccione y despache la misma tabla](#inspect-and-dispatch-the-same-table)
- [API de enrutamiento publicada](#released-routing-api)
- [Continúe por el recorrido de petición](#continue-the-request-path)

> [!WARNING]
> **Limitaciones actuales.** El router y las URLs con nombre están publicados. Las cachés de ruta y configuración están ausentes intencionalmente; los ejemplos avanzados se limitan a patrones ejercitados por el esqueleto y la suite de feature.

<a name="routes-are-executable-contracts"></a>
## Las rutas son contratos ejecutables

El router asocia un método HTTP y un patrón de ruta Go 1.22 a un handler estándar o acción de controller. Nombrar la ruta al registrarla da a templates, redirecciones, pruebas y CLI una fuente única de dirección. Los grupos añaden prefijo y middleware compartiendo la tabla de rutas padre.

<a name="register-and-reverse-one-route"></a>
## Registre y genere el camino de una ruta

Action adapta `func(*http.Context)` error; Get, Post, Put, Patch y Delete aceptan handlers `net/http`. Nombre la ruta antes del tráfico y use la tabla de rutas para generar el camino con la cantidad exacta de parámetros.

```go example
func Register(router *fhttp.Router, sessions *security.SessionStore, show func(*fhttp.Context) error) (string, error) {
	router.Action("GET", "/invoices/{invoice}", show,
		middleware.RequireAuth(sessions),
	).Name("invoices.show")

	path, err := router.Table().URL("invoices.show", "inv_123")
	if err != nil {
		return "", err
	}
	return path, nil
}
```

<a name="groups-constraints-and-the-exact-root"></a>
## Grupos, restricciones y raíz exacta

Los grupos comparten un prefijo y un límite de política. Las restricciones rechazan identificadores mal formados antes de una action y el patrón de raíz exacta evita que una ruta inicial capture cada subárbol en la gramática `ServeMux` estándar. Un handler `ConnectRPC` puede montarse con Post en su ruta de servicio generada; el mismo listener TLS sirve REST por HTTP/1 y Connect, `gRPC` o gRPC-Web por HTTP/2.

- Use `Group(prefix, middleware...)` para ruta y policy compartidas.
- Use `WhereNumber`, `WhereUuid`, `WhereUlid` o `WhereIn` para el formato de identificadores en la ruta.
- Use `/{name}` para un segmento y `/{path...}` solo cuando la acción posea el resto.
- Use `/{$}` para la raíz exacta de aplicación.
- Para RPC, configure explícitamente límites de bytes, concurrencia, tenant, duración, timeout del servidor, orígenes y cabeceras antes de montar el handler.

<a name="fail-while-building-not-after-a-click"></a>
## Falle durante la construcción, no después de un clic

Los nombres de ruta desconocidos y la cantidad equivocada de parámetros hacen que `Routes.URL` devuelva un error. Las discrepancias de restricción son respuestas 404, no resultados de autorización. Los registros duplicados o en conflicto son defectos de arranque, y el orden de registro debe seguir siendo determinista. El arranque también rechaza una ruta de la aplicación que responde a una ruta que el framework monta bajo `/_arandu/`, como `/_arandu/health`, en un método que la ruta del framework no acepta, nombrando el módulo y el patrón. La causa habitual es `r.Post("/", ...)`, cuyo patrón coincide con toda ruta debajo de la raíz; registre la raíz como "`/{$}`". La autenticación RPC se ejecuta antes de decodificar la petición, pero solo identifica un Subject: los servicios de la aplicación siguen autorizando y emitiendo Grants. La admisión se cierra durante el apagado, los streams observan la cancelación y el framework no abre un listener h2c ni un segundo listener RPC.

> [!WARNING]
> `Route.SetURI` cambia metadatos tras el registro, pero no mueve el patrón instalado en `ServeMux`. Defina la URI final antes del arranque.

<a name="inspect-and-dispatch-the-same-table"></a>
## Inspeccione y despache la misma tabla

Ejecute `aru route:list` y despache la misma tabla mediante httptest. Cubra un identificador válido, una restricción incompatible, un método no admitido, la raíz exacta y la URL con nombre. Para RPC, ejercite llamadas unarias, streaming del cliente, streaming del servidor y bidireccionales mediante Connect, `gRPC` y gRPC-Web; compruebe también rechazo CORS, admisión global y por tenant, límites de cuerpo, cancelación y apagado ordenado en el listener compartido.

- El inventario nombra método, patrón, ruta y módulo.
- Un parámetro válido dentro de la restricción llega a la acción.
- Un parámetro restringido malformado devuelve 404.
- La raíz exacta no captura una ruta desconocida.
- La autorización RPC y las cuatro formas de streaming se ejecutan en el listener de la aplicación.

<a name="released-routing-api"></a>
## API de enrutamiento publicada

`http.NewRouter` y `Router.Action` crean y llenan la tabla. Los métodos de verbo de Router registran handlers estándar; `Router.Group` compone prefijos y middleware; `Route.Name` y `Route.Where*` agregan metadatos estables; `Routes.URL` construye una ruta o informa una incompatibilidad. `rpc.NewHandler` devuelve la ruta de servicio generada y el handler para Factory, Authenticator bearer y Config explícitos. `rpc.Subject` lee la identidad autenticada dentro de un procedimiento; no es un resultado de autorización.

| Superficie | Contrato publicado |
| --- | --- |
| `Router.Action` | registra una acción de controller |
| `Get/Post/Put/Patch/Delete` | registra un handler estándar |
| Group y Name | compone prefijos y direcciones estables |
| `Where*` y `Routes.URL` | restringe parámetros y genera rutas inversas |
| `rpc.NewHandler` | construye un handler Connect, `gRPC` y gRPC-Web con límites explícitos |
| `rpc.Subject` | lee el Subject autenticado antes de que la policy de la aplicación emita un Grant |

<a name="continue-the-request-path"></a>
## Continúe por el recorrido de petición

REST sigue siendo el límite predeterminado de recursos de petición y respuesta. `gRPC` añade llamadas tipadas unarias y por streaming en el mismo servidor HTTP; no es un segundo proceso ni un atajo de autorización. Controllers se encargan de la traducción de transporte, middleware de la política transversal de petición y la firma de URL de la integridad después de elegir la ruta.

- [Patrones de comunicación](/docs/{{version}}/architecture/contracts-adapters)
- [Controllers](/docs/{{version}}/http/controllers)
- [Middleware](/docs/{{version}}/http/middleware)
- [URLs](/docs/{{version}}/http/urls)
