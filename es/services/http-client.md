# Cliente HTTP

- [Cada petición externa empieza con contexto y plazo](#every-outbound-request-begins-with-context-and-a-deadline)
- [Envíe una petición JSON acotada](#send-one-bounded-json-request)
- [Añada headers, reintentos y decodificación deliberadamente](#add-headers-retries-and-decoding-deliberately)
- [Distinga rechazo, timeout y respuesta malformada](#distinguish-refusal-timeout-and-malformed-response)
- [Pruebe contra un servidor local del protocolo](#test-against-a-local-protocol-server)
- [Superficie publicada del cliente HTTP](#released-http-client-surface)
- [Instrumente la llamada y mantenga fuera la política de aplicación](#instrument-the-call-and-keep-application-policy-outside-it)

> [!WARNING]
> **Limitaciones actuales.** El cliente propaga contexto y restringe reintentos a peticiones seguras, pero ningún recorrido actual de primera parte ejercita todos los verbos. `Batch.Send` tampoco tiene cobertura de comportamiento completa; prefiera petición única hasta verificar esa frontera.

<a name="every-outbound-request-begins-with-context-and-a-deadline"></a>
## Cada petición externa empieza con contexto y plazo

El cliente HTTP de Hesape envuelve `net/http` con verbos conscientes del contexto, opciones de petición, decodificación y middleware. El componente está publicado, pero no hay consumidor actual de primera parte etiquetado; esta guía se limita al contrato verificado del paquete, sin prometer composición predeterminada en el esqueleto.

<a name="send-one-bounded-json-request"></a>
## Envíe una petición JSON acotada

Construya Factory al arrancar, defina el límite de bytes de respuesta, derive un plazo hijo y envíe JSON tipado mediante `PendingRequest`. Verifique el transporte y llame a Throw para 4xx o 5xx antes de decodificar la respuesta ya cargada.

```go example
func PostPayload(ctx context.Context, endpoint string, payload map[string]any) error {
	factory := client.NewFactory(nil).MaxResponseBytes(2 << 20)
	request := factory.CreatePendingRequest().Timeout(5 * time.Second).AcceptJSON()
	callCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	response, err := request.Post(callCtx, endpoint, payload)
	if err != nil {
		return err
	}
	if err := response.Throw(nil); err != nil {
		return err
	}
	return nil
}
```

<a name="add-headers-retries-and-decoding-deliberately"></a>
## Añada headers, reintentos y decodificación deliberadamente

Configure URL base y headers inmutables en la construcción, añada autorización específica sin registrarla, propague contexto, limite cuerpos de respuesta y reintente solo operaciones idempotentes o con clave explícita de idempotencia, usando backoff acotado.

- Use contexto en cada verbo.
- Defina plazos para cliente y operación.
- Limite los bytes de respuesta antes de decodificar.
- Instrumente mediante el Collector de la aplicación.

<a name="distinguish-refusal-timeout-and-malformed-response"></a>
## Distinga rechazo, timeout y respuesta malformada

DNS, conexión, TLS, cancelación, plazo, estado no 2xx y cuerpo inválido tienen evidencias distintas. Nunca convierta una dependencia indisponible en éxito vacío ni reintente una mutación insegura solo porque el transporte se cerró.

> [!WARNING]
> Un timeout cero o cliente predeterminado sin configuración pueden retener recursos indefinidamente. Oculte URLs con credenciales y limite la captura del cuerpo de error para impedir que un servicio externo llene logs o memoria.

<a name="test-against-a-local-protocol-server"></a>
## Pruebe contra un servidor local del protocolo

Use `httptest.Server` para método, cancelación, headers, cuerpo JSON, respuesta grande acotada, redirecciones, respuesta lenta, estado reintentable, mutación insegura, JSON malformado y confianza TLS. Verifique que las credenciales no lleguen a los logs.

- La cancelación llega al servidor.
- Los reintentos siguen la política de idempotencia.
- Los límites de respuesta actúan antes de una asignación sin techo.
- La observabilidad registra URL, estado y duración sin secretos.

<a name="released-http-client-surface"></a>
## Superficie publicada del cliente HTTP

`NewFactory`, `Factory.MaxResponseBytes`, `Factory.CreatePendingRequest`, verbos de `PendingRequest`, `Response.Throw` y decodificación forman la API publicada. Response lee y cierra el cuerpo de inmediato; el llamador no cierra un campo Body exportado.

| Superficie | Contrato publicado |
| --- | --- |
| `NewFactory` / `MaxResponseBytes` | construye la factory protegida |
| `CreatePendingRequest` / Timeout | describe una petición |
| `Get/Post/Put/Patch/Delete` | verbos HTTP conscientes del contexto |
| `Response.Throw` / Object | estado y decodificación tipada |

<a name="instrument-the-call-and-keep-application-policy-outside-it"></a>
## Instrumente la llamada y mantenga fuera la política de aplicación

El contexto proporciona cancelación, los logs registran esperas externas, los contratos ocultan detalles de proveedores a servicios de dominio y las pruebas poseen el fixture del protocolo.

- [Contexto](/docs/{{version}}/architecture/context)
- [Logs](/docs/{{version}}/operations/logging-observability)
- [Contratos](/docs/{{version}}/architecture/contracts-adapters)
- [Dobles de prueba](/docs/{{version}}/testing/test-doubles)
