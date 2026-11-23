# Ciclo de vida de la solicitud

- [Cada solicitud recorre un pipeline explícito](#one-request-crosses-one-explicit-pipeline)
- [Identifica la ruta y sus middlewares](#know-the-route-and-the-middleware-that-own-it)
- [Observa una respuesta completa](#observe-a-complete-response)
- [Prueba mediante Handler cuando el puerto no importe](#test-through-handler-when-the-port-is-irrelevant)
- [Cada etapa tiene una responsabilidad](#each-stage-has-one-question)

<a name="one-request-crosses-one-explicit-pipeline"></a>
## Cada solicitud recorre un pipeline explícito

El servidor Go acepta la conexión, aplica los middlewares globales, encuentra una ruta, aplica sus middlewares, crea un contexto HTTP, llama al handler registrado y escribe la respuesta. Las pruebas pueden inicializar la misma aplicación compuesta que usa el servidor sin abrir un puerto.

<a name="know-the-route-and-the-middleware-that-own-it"></a>
## Identifica la ruta y sus middlewares

Use el inventario de rutas en lugar de adivinar a partir de una URL. El framework compone el pipeline global a partir de la configuración, y `bootstrap/app.go` lo ajusta en `WithMiddleware`. El middleware de ruta se adjunta donde se declara la ruta. Un controller recibe sus colaboradores por el constructor; no los resuelve durante la petición.

```sh
aru route:list
aru dev
```

<a name="observe-a-complete-response"></a>
## Observa una respuesta completa

```sh
curl -i http://127.0.0.1:8080/
```

La respuesta atravesó el pipeline global en este orden: las frases de rechazo, Recover, lo que añadió Prepend, los dos límites de cuerpo, `TrustProxies`, Observe, `SecurityHeaders`, lo que añadió `BeforeThrottle`, el Throttle global, la verificación de CSRF, `OverrideMethod` y lo que añadió Append. Después vienen la coincidencia de ruta, el middleware de la ruta, el controller y el renderer. En una ruta que cambia estado, autenticación, autorización, CSRF y validación siguen siendo comprobaciones separadas y no una guarda genérica.

> [!NOTE]
> Una petición bajo `/_arandu/` omite el pipeline de la aplicación solo cuando una ruta del framework la responde: una sonda, la consola de depuración, el reload o un asset. Cualquier otra ruta bajo ese prefijo es tráfico de la aplicación y atraviesa el pipeline completo. Create compone la aplicación pero no abre ningún listener; `HandleCommand` inicializa los módulos y ejecuta el servidor para serve, el comando sin argumentos, mientras que un comando que solo necesita la configuración y la base de datos no inicializa nada.

<a name="test-through-handler-when-the-port-is-irrelevant"></a>
## Prueba mediante Handler cuando el puerto no importe

```go example
func assertHomeIsOK(t *testing.T, app *foundation.Application) {
	request := httptest.NewRequest(stdhttp.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	app.Handler().ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}
}
```

- Un panic se convierte en la respuesta de error configurada y se registra una sola vez.
- La cancelación de una solicitud llega a las llamadas posteriores a bases de datos y clientes mediante el contexto.
- Las pruebas distinguen los fallos de autenticación de los rechazos de la Policy.
- La respuesta se escribe una sola vez, con el estado y el tipo de contenido previstos.

<a name="each-stage-has-one-question"></a>
## Cada etapa tiene una responsabilidad

El enrutamiento elige el handler. Los middlewares gestionan las políticas transversales de la solicitud. Los controllers traducen el transporte en una llamada de servicio. Los servicios autorizan los casos de uso y los repositorios acceden a los datos dentro del ámbito del tenant. Las vistas y los resources definen el formato de salida. Separar estas responsabilidades permite rastrear una solicitud fallida.

- [Rutas](/docs/{{version}}/http/routing)
- [Middleware](/docs/{{version}}/http/middleware)
- [Controllers](/docs/{{version}}/http/controllers)
- [Contexto](/docs/{{version}}/architecture/context)
- [Aplicación de Framework](https://pkg.go.dev/github.com/arandu-io/framework/foundation)
