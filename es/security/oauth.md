# OAuth e inicio de sesión social

- [Un cliente OAuth completa el intercambio del código de autorización](#an-oauth-client-completes-the-authorization-code-exchange)
- [Registre un callback exacto y mantenga el secreto en el arranque](#register-an-exact-callback-and-keep-the-secret-at-boot)
- [Redirija con estado y rechace un callback incompatible](#redirect-with-state-and-refuse-a-mismatched-callback)
- [Normalice datos del proveedor y aplique la policy de aplicación](#normalize-provider-data-then-apply-application-policy)
- [Rechace antes de intercambiar el token](#keep-refusal-before-the-token-exchange)
- [El estado es obligatorio; la vinculación es autorización](#state-is-mandatory-and-linking-is-authorization)
- [Pruebe redirección y callback como un par](#test-the-redirect-and-callback-as-one-pair)
- [El callback produce entrada de identidad, no una sesión](#the-callback-produces-identity-input-not-a-session)

<a name="an-oauth-client-completes-the-authorization-code-exchange"></a>
## Un cliente OAuth completa el intercambio del código de autorización

El paquete de proveedores publicado envía el navegador al endpoint externo de autorización, verifica el estado del callback, intercambia el código y recupera `oauth.UserData` normalizado. Hay constructores para `GitHub`, Google, Facebook, Stripe y Apple, además de `NewProvider` para un servicio configurado explícitamente.

Este paquete es un cliente OAuth, no un servidor de autorización ni un login social completo. El esqueleto y el kit UI no proporcionan rutas de callback, configuración de credenciales, vinculación de cuentas, policy de tenant ni model duradero de identidad para ese flujo; eso sigue siendo trabajo de la aplicación.

<a name="register-an-exact-callback-and-keep-the-secret-at-boot"></a>
## Registre un callback exacto y mantenga el secreto en el arranque

- Registre en el proveedor la URL HTTPS de callback exactamente como la enviará la aplicación.
- Cargue ID y secreto del cliente desde configuración de despliegue, nunca desde el código ni desde una respuesta al navegador.
- Decida el tenant de la ruta de entrada antes de redirigir.
- Defina una policy de vinculación para identidades existentes antes de aceptar un callback.

Construya Provider y `CookieStateStore` nuevos en cada petición. Provider mantiene estado del par de peticiones y no es seguro compartirlo concurrentemente.

<a name="redirect-with-state-and-refuse-a-mismatched-callback"></a>
## Redirija con estado y rechace un callback incompatible

```go example
var (
	clientID     = "client-id"
	clientSecret = "client-secret"
	callbackURL  = "https://example.com/auth/github/callback"
)

func githubProvider(w stdhttp.ResponseWriter, r *stdhttp.Request, id, secret, callback string) *providers.Provider {
	state := providers.NewCookieStateStore(w, r)
	return providers.NewGithubProvider(state, id, secret).RedirectURL(callback)
}

func startGithub(w stdhttp.ResponseWriter, r *stdhttp.Request) error {
	return githubProvider(w, r, clientID, clientSecret, callbackURL).Redirect(w, r)
}

func finishGithub(w stdhttp.ResponseWriter, r *stdhttp.Request) error {
	identity, token, err := githubProvider(w, r, clientID, clientSecret, callbackURL).User(r)
	if errors.Is(err, providers.ErrStateMismatch) {
		stdhttp.Error(w, "invalid oauth callback", stdhttp.StatusBadRequest)
		return nil
	}
	if err != nil {
		return err
	}
	_ = identity // Resolve or link through the tenant-aware account service.
	_ = token    // Keep secret; persist only when the provider integration needs it.
	return nil
}
```

La entrada llama a `startGithub` y retorna en cuanto Redirect escribe 302. El callback es otra petición tratada por `finishGithub`; User verifica y consume el estado antes de intercambiar el código. Trate el token devuelto por el proveedor como secreto aunque solo necesite identidad.

<a name="normalize-provider-data-then-apply-application-policy"></a>
## Normalice datos del proveedor y aplique la policy de aplicación

`oauth.UserData.Get` y String leen campos normalizados de identidad. Los atributos específicos del proveedor siguen disponibles sin cambiar el formato común. Tras el callback, busque el sujeto previamente vinculado; si no existe, siga una decisión explícita de vinculación o registro, sin asociar silenciosamente cuentas por el mismo correo.

| Constructor | Proveedor publicado |
| --- | --- |
| `NewGithubProvider` | `GitHub` |
| `NewGoogleProvider` | Google |
| `NewFacebookProvider` | Facebook |
| `NewStripeProvider` | Stripe |
| `NewAppleProvider` | Apple |
| `NewProvider` | endpoints proporcionados por el llamador |

`NewAppleProvider` recibe un Services ID, un Team ID, un Key ID y la clave privada en lugar de un client secret, y construye el client secret en cada intercambio. Apple responde a la redirección con un POST desde su propio sitio, así que la ruta de callback acepta POST y queda exenta de la protección CSRF solo para esa ruta, con `CSRFExcept`; la verificación de estado ocupa su lugar. El sitio debe servirse por HTTPS, y la cookie de estado se envía con SameSite=None. La documentación del paquete, en `NewAppleProvider`, describe el flujo completo.

<a name="keep-refusal-before-the-token-exchange"></a>
## Rechace antes de intercambiar el token

Estado ausente, cookie ausente, estado alterado y replay devuelven `ErrStateMismatch`. Un callback con `error=access_denied` es rechazo, no fallo de servidor. Código ausente, fallo de intercambio, respuesta malformada y fallo de datos de usuario siguen siendo errores internos distintos; el navegador recibe un mensaje acotado sin secretos.

La cookie de estado tiene un nombre predeterminado. Iniciar otro OAuth en el mismo navegador antes del primer callback puede sustituirla y hacer fallar el primero; trátelo como un flujo reiniciable, sin debilitar la verificación.

<a name="state-is-mandatory-and-linking-is-authorization"></a>
## El estado es obligatorio; la vinculación es autorización

> [!WARNING]
> `Provider.Stateless` desactiva la defensa CSRF publicada. Úselo solo cuando otro mecanismo concreto y probado vincule el callback al cliente iniciador; la comodidad no es ese mecanismo.

Nunca vincule automáticamente solo porque un servicio externo devolvió correo. Verifique la garantía del proveedor, trate colisiones deliberadamente y autorice la transición de tenant. No ponga secretos del cliente ni tokens en logs, query, payload de sesión ni errores renderizados.

`CookieStateStore` :hesapeVersion envía su cookie de estado con Secure por defecto, sea cual sea el esquema con el que llegó la petición, así que la cookie conserva el atributo detrás de un proxy que termina TLS. Un servidor de desarrollo en http simple llama a `Secure(false)` en el store, porque un navegador nunca devuelve una cookie Secure por http; déjelo activado en todos los demás casos.

<a name="test-the-redirect-and-callback-as-one-pair"></a>
## Pruebe redirección y callback como un par

- Redirect devuelve 302 con ID de cliente, URI exacta de retorno, ámbitos y estado no vacío.
- Un callback con estado almacenado alcanza los servidores fake de token y usuario.
- Un estado ausente, alterado, reutilizado o de otro navegador falla antes de cualquier intercambio.
- El rechazo y la indisponibilidad del proveedor tienen evidencias operativas distintas sin exponer secretos.
- El recorrido HTTPS con proxy desplegado emite una cookie de estado Secure y `HttpOnly`, con SameSite=Lax (SameSite=None para Apple).
- La vinculación no puede cruzar tenants ni apropiarse de una cuenta por colisión de correo.

```sh
go test github.com/arandu-io/hesape/oauth/providers
```

<a name="the-callback-produces-identity-input-not-a-session"></a>
## El callback produce entrada de identidad, no una sesión

Tras aceptar la vinculación la policy, la aplicación puede convertir la cuenta local en `security.Subject` y rotar su sesión. Los tokens de API, tokens del proveedor y sesiones de navegador tienen credenciales y ciclos de vida distintos. Mantenga la frontera externa del proveedor fuera del repositorio de cuentas.

- [Autenticación](/docs/{{version}}/security/authentication)
- [Sesión HTTP](/docs/{{version}}/http/session)
- [Policies y Grant](/docs/{{version}}/security/authorization)
- [API de proveedores OAuth](https://pkg.go.dev/github.com/arandu-io/hesape/oauth/providers)
- [API de datos de usuario OAuth](https://pkg.go.dev/github.com/arandu-io/hesape/oauth)
