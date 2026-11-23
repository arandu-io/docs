# Protección CSRF

- [Por qué una mutación necesita un segundo secreto](#why-a-state-changing-request-needs-a-second-secret)
- [Instale un emisor envolviendo la pila de rutas](#install-one-issuer-around-the-route-stack)
- [Proteja el handler antes de ejecutarlo](#protect-the-handler-before-it-can-run)
- [Cubra formularios y peticiones mejoradas](#cover-forms-and-enhanced-requests)
- [Rechace sin adivinar](#refuse-without-guessing)
- [La rotación cierra la frontera anterior](#rotation-closes-the-old-boundary)
- [Pruebe el rechazo con el mismo cuidado que el éxito](#test-the-rejection-as-carefully-as-the-success)
- [CSRF es una capa de confianza de la petición](#csrf-is-one-layer-of-request-trust)

<a name="why-a-state-changing-request-needs-a-second-secret"></a>
## Por qué una mutación necesita un segundo secreto

El navegador adjunta automáticamente la cookie de sesión a una petición compatible. Eso permite que otro sitio provoque la petición: estar autenticado no demuestra la intención de la página. El token CSRF se vincula al ID de sesión y debe llegar mediante marcado o header que otro sitio no pueda leer.

<a name="install-one-issuer-around-the-route-stack"></a>
## Instale un emisor envolviendo la pila de rutas

`foundation.NewDrivers` construye el único emisor, `drivers.CSRF`, a partir de la clave de la aplicación y `CSRF_TTL`, y lo vincula al store de sesión, `drivers.Sessions`, de modo que el id de sesión solo se acepta después de verificar su cookie firmada. Create instala la verificación en el pipeline global; no hay que cablear nada a mano. El token de un visitante anónimo debe vincularse a un id aleatorio que lleva la cookie firmada de ese mismo visitante, nunca a un id de sesión vacío que compartiría todo visitante: un token que tiene un visitante debe rechazarse en el formulario de otro. Las cookies de sesión de producción siguen exigiendo Secure sobre HTTPS.

<a name="protect-the-handler-before-it-can-run"></a>
## Proteja el handler antes de ejecutarlo

La verificación se ejecuta en el pipeline global después del throttle y antes de todo handler de mutación. Las peticiones seguras GET, HEAD, OPTIONS y TRACE pasan. Las peticiones inseguras leen X-CSRF-Token antes del campo de formulario `_token` y validan contra el id de sesión verificado actual. La única decisión que queda a la aplicación es cuáles de sus propias rutas tienen una guarda propia, nombradas en `WithMiddleware`.

```go example
func WithoutCSRFOnWebhooks(drivers *foundation.Drivers) *foundation.ApplicationBuilder {
	return foundation.Configure(drivers).
		WithMiddleware(func(m *foundation.Middleware) {
			m.ValidateCSRFTokensExcept(
				"/webhooks/",                 // a prefix: every path under it
				"/v1/usage",                  // exactly this path
				"/checkout/{token}/sessions", // one segment of any value
			)
		})
}
```

<a name="cover-forms-and-enhanced-requests"></a>
## Cubra formularios y peticiones mejoradas

Un formulario convencional usa `@csrf`, compilado en campo oculto. El layout también coloca el token en una meta tag ypy-csrf, y la extensión ypy-csrf lo envía como X-CSRF-Token en toda petición del mismo origen que cambia estado; POST, PUT, PATCH y DELETE parciales, sin formulario, no necesitan inventar otra conexión de token. Al sustituir un formulario, renderice un campo nuevo en el fragmento.

Una exención deja una ruta a su propia guarda, como un webhook que verifica la firma del proveedor. Una entrada que termina en / cubre toda ruta debajo de ella, una entrada sin la barra cubre exactamente esa ruta, y un segmento `{name}` coincide con cualquier segmento no vacío, como en un patrón de ruta; una entrada no puede empezar con uno. Un módulo que implementa CSRFExempt declara sus propios receptores: Create los exime junto a los de la aplicación, Boot rechaza una exención que cubra una ruta de escritura de otro dueño, y cada exención recogida se escribe en el log al arrancar.

```kyse
<form method="post" action="/profile">
  @csrf
</form>
```

<a name="refuse-without-guessing"></a>
## Rechace sin adivinar

No cree un token sustituto al tratar la mutación rechazada para continuar la operación. Una mutación con origen cruzado indicado por el navegador devuelve 403. Un token ausente, inválido, expirado o de otra sesión devuelve 419; una petición YPY recibe la instrucción de actualización del framework, YPY-Refresh. La siguiente página proporciona un token actual sin ejecutar la mutación rechazada.

<a name="rotation-closes-the-old-boundary"></a>
## La rotación cierra la frontera anterior

Rote el ID de sesión cuando cambie la autoridad de autenticación. Los tokens del ID anterior dejan de funcionar en la nueva sesión. Nunca ponga tokens en URLs ni logs, ni exima un endpoint de mutación por recibir normalmente peticiones mejoradas: un atacante puede llamarlo directamente.

<a name="test-the-rejection-as-carefully-as-the-success"></a>
## Pruebe el rechazo con el mismo cuidado que el éxito

Ejercite campo válido, header válido, token ausente, alterado, de otra sesión y expirado. En formularios renderizados, cuente formularios de mutación y directivas de token para impedir omisiones silenciosas en nuevos formularios. Verifique también que el rechazo mantiene intacto el registro protegido.

<a name="csrf-is-one-layer-of-request-trust"></a>
## CSRF es una capa de confianza de la petición

El token demuestra que la petición pasó por una página que posee la sesión; no decide quién es el sujeto ni sus permisos. Autenticación de sesión, autorización por policy, validación y headers del mismo origen mantienen responsabilidades separadas y aún se ejecutan con un token válido.
