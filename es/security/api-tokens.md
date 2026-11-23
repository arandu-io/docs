# Tokens de API

- [`TokenGuard` autentica una credencial presentada](#tokenguard-authenticates-one-presented-credential)
- [Asuma emisión, almacenamiento y adapter de petición](#own-issuance-storage-and-the-request-adapter)
- [Autentique un bearer token previamente provisionado](#authenticate-a-pre-provisioned-bearer-token)
- [Lea credenciales solo desde headers](#read-credentials-from-headers-only)
- [Los tokens desconocidos siguen siendo respuestas 401 normales](#unknown-tokens-remain-ordinary-401-responses)
- [Una credencial permanente exige un mecanismo de revocación de aplicación](#a-permanent-credential-needs-an-application-kill-switch)
- [Pruebe lo que hace y lo que no hace el guard](#test-what-the-guard-does-and-what-it-does-not-do)
- [Use sesiones en el navegador y construya el ciclo de vida por separado](#use-sessions-for-browsers-and-build-lifecycle-separately)

> [!WARNING]
> **Limitaciones actuales.** `TokenGuard` autentica bearer o token de entrada, pero no proporciona expiración, ámbitos, rotación ni revocación. En modo hash, use solo petición hacia User o Check: Validate envía el token sin procesar. User reduce errores del proveedor a nil y `SetRequest` no limpia el usuario en caché; construya un guard por petición y preserve errores operativos en una frontera propia cuando el protocolo exija distinguirlos.

<a name="tokenguard-authenticates-one-presented-credential"></a>
## `TokenGuard` autentica una credencial presentada

El `auth.TokenGuard` publicado lee un bearer token, consulta `UserProvider` de aplicación por la cuenta correspondiente y la expone mediante el contrato Guard. No tiene sesión y demuestra la credencial de nuevo en cada petición.

> [!WARNING]
> Este no es un sistema de emisión ni ciclo de vida. El guard publicado no define expiración, ámbitos, rotación, revocación, models de token, rutas de gestión ni almacenamiento de primera parte. El token se acepta mientras el proveedor mantenga un valor correspondiente.

<a name="own-issuance-storage-and-the-request-adapter"></a>
## Asuma emisión, almacenamiento y adapter de petición

Antes de conectar el guard, genere un token opaco de alta entropía, muéstrelo una vez, almacene solo SHA-256 e implemente `auth.UserProvider` sobre el digest. Adapte también la petición a `auth.Request`; los tipos HTTP publicados aún no implementan directamente todos los métodos.

- Genere credenciales con `crypto/rand`, no nombres, timestamps ni UUIDs elegidos como identificadores.
- Almacene el digest en un registro dedicado de credencial, no en una columna reutilizable de contraseña o remember token.
- Decida vínculo al tenant y estado deshabilitado en la búsqueda del proveedor.
- Proporcione un recorrido operativo separado para sustitución y eliminación; `TokenGuard` no realiza ninguna.

<a name="authenticate-a-pre-provisioned-bearer-token"></a>
## Autentique un bearer token previamente provisionado

Este ejemplo cubre solo autenticación contra un token ya existente. Incorporar Request de Hesape proporciona accesores de query, input, bearer y cookies; los tres métodos siguientes completan el contrato `auth.Request` exacto.

```go example
type tokenRequest struct{ *hhttp.Request }

func (r tokenRequest) GetUser() string {
	user, _, _ := r.RawRequest().BasicAuth()
	return user
}
func (r tokenRequest) GetPassword() string {
	_, password, _ := r.RawRequest().BasicAuth()
	return password
}
func (r tokenRequest) Context() context.Context {
	return r.RawRequest().Context()
}

func Authenticate(w stdhttp.ResponseWriter, r *stdhttp.Request, users auth.UserProvider) {
	request := tokenRequest{Request: hhttp.NewRequest(r)}
	guard := auth.NewTokenGuard(users, request, "api_token", "api_token_sha256", true)
	user := guard.User()
	if user == nil {
		stdhttp.Error(w, "unauthorized", stdhttp.StatusUnauthorized)
		return
	}
}
```

> [!NOTE]
> Con hash=true, User calcula SHA-256 del token recibido antes de `RetrieveByCredentials`. El campo de almacenamiento del proveedor debe contener el digest hexadecimal en minúsculas esperado.

<a name="read-credentials-from-headers-only"></a>
## Lea credenciales solo desde headers

`GetTokenForRequest` prefiere Authorization: Bearer y usa contraseña HTTP Basic como fallback. Ignora query y cuerpo para mantener credenciales fuera de URLs, historial, referrer y logs comunes. Cree un `TokenGuard` nuevo por petición. En :hesapeVersion, `SetRequest` cambia la petición pero no limpia el usuario en caché de User; no es seguro reutilizarlo tras resolución.

> [!WARNING]
> El fallback de contraseña Basic no puede desactivarse mediante una opción de `TokenGuard`. Si el endpoint acepta solo bearer, rechace autorización no Bearer antes de invocar el guard.

<a name="unknown-tokens-remain-ordinary-401-responses"></a>
## Los tokens desconocidos siguen siendo respuestas 401 normales

Los tokens ausentes, malformados, desconocidos y deshabilitados generan la misma respuesta no autenticada. En :hesapeVersion, User también reduce errores de `RetrieveByCredentials` a nil: la ruta falla cerrada como 401 ante indisponibilidad del proveedor, pero no la distingue de token desconocido. Registre el fallo en el proveedor o use una frontera propia que preserve el error cuando el protocolo necesite estado operativo. Nunca registre el token presentado.

> [!WARNING]
> En Hesape :hesapeVersion, `TokenGuard.Validate` no aplica SHA-256 configurado con hash=true, aunque User sí lo hace. No use ni documente Validate como equivalente a User en ese modo. La prueba publicada demuestra hashing solo mediante User y Check.

<a name="a-permanent-credential-needs-an-application-kill-switch"></a>
## Una credencial permanente exige un mecanismo de revocación de aplicación

El digest impide que una filtración de base revele inmediatamente el token, pero no proporciona expiración ni revocación. El proveedor debe rechazar registros deshabilitados y vincular el token al tenant y cuenta esperados. Use TLS, no exponga tokens en URLs ni logs, consulte el digest por índice y rote creando una credencial nueva antes de eliminar la antigua.

No confunda autenticación con autorización. Convierta la cuenta autenticada en Subject y ejecute la misma Policy usada por peticiones autenticadas mediante sesión.

<a name="test-what-the-guard-does-and-what-it-does-not-do"></a>
## Pruebe lo que hace y lo que no hace el guard

- Un header Bearer resuelve la cuenta y Check tiene éxito.
- La query y el cuerpo nunca autentican.
- El digest almacenado solo se resuelve con hash=true mediante User.
- Los tokens ausentes, alterados y deshabilitados devuelven 401 sin repetir la credencial.
- Un error del proveedor se vuelve nil y el mismo 401 mediante `TokenGuard.User`, mientras la telemetría del proveedor registra la caída.
- Un guard que resolvió un usuario nunca se reutiliza mediante `SetRequest` en otra petición.
- El rechazo de policy aún bloquea un token autenticado.

```sh
go test github.com/arandu-io/hesape/auth -run TokenGuard
```

<a name="use-sessions-for-browsers-and-build-lifecycle-separately"></a>
## Use sesiones en el navegador y construya el ciclo de vida por separado

La autenticación de navegador usa sesiones firmadas y CSRF. `TokenGuard` solo corresponde cuando la aplicación ya proporcionó el ciclo de credenciales ausente. OAuth obtiene identidad de proveedor y no sustituye un emisor de tokens de API de aplicación.

- [Autenticación](/docs/{{version}}/security/authentication)
- [Policies y Grant](/docs/{{version}}/security/authorization)
- [OAuth e inicio de sesión social](/docs/{{version}}/security/oauth)
- [API `TokenGuard`](https://pkg.go.dev/github.com/arandu-io/hesape/auth#TokenGuard)
