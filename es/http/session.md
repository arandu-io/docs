# Sesión HTTP

- [Una cookie firmada apunta a un Subject tipado](#a-signed-cookie-points-to-a-typed-subject)
- [Elija el backend según la topología de despliegue](#choose-the-backend-for-the-deployment-topology)
- [Rote hacia una sesión autenticada](#rotate-into-an-authenticated-session)
- [Cargue, confirme, recuerde y termine sesiones](#load-confirm-remember-and-end-sessions)
- [Un registro ausente permanece sin autenticar](#a-missing-record-stays-unauthenticated)
- [Los cambios de identidad exigen cambios de ID](#identity-changes-require-id-changes)
- [Demuestre el recorrido entre réplicas](#prove-the-cross-replica-path)
- [La autenticación define identidad; la policy define autoridad](#authentication-owns-identity-policy-owns-authority)

> [!WARNING]
> **Limitaciones actuales.** `foundation.NewDrivers` construye el store de sesión a partir de `SESSION_DRIVER`, `REDIS_URL`, `SESSION_LIFETIME` y `SESSION_SECURE_COOKIE`, incluido un store compartido sobre RESP para varias réplicas una vez enlazado `hesape/redis`, y producción rechaza `SESSION_SECURE_COOKIE=false`. El nombre y el alcance de la cookie son fijos: `arandu_session`, ruta /, el host que respondió y `SameSite` Lax; el framework rechaza `SESSION_TTL` y un `SESSION_PATH`, `SESSION_DOMAIN` o `SESSION_SAME_SITE` que pida otro alcance.

<a name="a-signed-cookie-points-to-a-typed-subject"></a>
## Una cookie firmada apunta a un Subject tipado

El almacenamiento de sesión del framework envía al navegador un ID firmado y mantiene `security.Subject` en el backend. La carga verifica la cookie antes de leer el registro, impidiendo que los datos de petición elijan tenant o cuenta. La sesión establece identidad; no concede permiso sobre datos.

<a name="choose-the-backend-for-the-deployment-topology"></a>
## Elija el backend según la topología de despliegue

`foundation.NewDrivers` construye el único store de sesión, `drivers.Sessions`, a partir de la clave de la aplicación, `SESSION_LIFETIME` y la decisión sobre Secure, sobre el backend que nombra `SESSION_DRIVER`. Memory sirve para tests y para un único proceso local. Varias réplicas necesitan `SESSION_DRIVER=redis`, el conector `hesape/redis` enlazado por su import en blanco en `bootstrap/app.go` y un `REDIS_URL` alcanzable. El almacenamiento de sesión es independiente de `CACHE_STORE`.

| Configuración | Efecto publicado |
| --- | --- |
| `SESSION_DRIVER=memory` | un backend dentro del proceso |
| `SESSION_DRIVER=redis` | handler de sesión compartido con RESP |
| `SESSION_LIFETIME` | duración en minutos, leída por el framework; dos horas si no se define |
| `SESSION_SECURE_COOKIE` | atributo Secure, leído por el framework; activo salvo con `APP_ENV=dev`, y false se rechaza en producción |
| `SESSION_PATH` / `SESSION_DOMAIN` / `SESSION_SAME_SITE` | no son configurables; cualquier alcance distinto de /, el host que respondió y lax detiene el arranque |

> [!WARNING]
> Nunca despliegue dos réplicas con backend en memoria. Ambos procesos pueden arrancar, pero una sesión creada en uno no puede cargarse en el otro.

<a name="rotate-into-an-authenticated-session"></a>
## Rote hacia una sesión autenticada

Use Rotate al iniciar sesión aunque no haya ID anterior. Crea un ID nuevo, escribe la cookie firmada y destruye el registro previo a autenticación si existe. Start no expresa sustitución, por lo que no es la operación de login.

```go example
func RotateSession(r *stdhttp.Request, w stdhttp.ResponseWriter, drivers *foundation.Drivers, subject security.Subject) error {
	sessions := drivers.Sessions

	oldID := sessions.IDFromRequest(r)
	_, err := sessions.Rotate(r.Context(), w, oldID, subject)
	if err != nil {
		return fmt.Errorf("rotate session: %w", err)
	}
	return nil
}
```

> [!NOTE]
> El atributo Secure del store sigue a `SESSION_SECURE_COOKIE` y está activo en todas partes salvo `APP_ENV=dev` cuando la variable no está definida. Producción rechaza false al arrancar, porque allí la cookie es la credencial.

<a name="load-confirm-remember-and-end-sessions"></a>
## Cargue, confirme, recuerde y termine sesiones

- Load devuelve un Subject fiable o un error explícito de sesión.
- Confirm registra confirmación reciente de contraseña para rutas sensibles.
- `Remember(true)` amplía la sesión mediante la opción con nombre.
- Destroy elimina el registro actual y limpia cookies de sesión y destino previsto.
- `DestroyOthers` elimina otros registros del mismo tenant y sujeto, preservando el ID elegido.

### Guarde las sesiones en la base de datos de la aplicación

```go example
func NewSessions(drivers *foundation.Drivers, appKey []byte, ttl time.Duration) (*security.SessionStore, error) {
	handler, err := session.NewDatabaseHandler[security.Subject](drivers.DB, "")
	if err != nil {
		return nil, err
	}
	backend := security.NewSessionBackend(handler)
	sessions := security.NewSessionStore(appKey, ttl, true, backend)
	return sessions, nil
}
```

`SESSION_DRIVER=redis` resuelve su conexión independientemente de `CACHE_STORE`, así que las sesiones pueden compartirse mientras los valores comunes de caché quedan en la memoria del proceso. Hesape también incluye `session.NewDatabaseHandler`, que guarda cada registro en la tabla que crea `session.RecordMigrations`, para más de una réplica sin servidor RESP. Es un componente, no una configuración: `SESSION_DRIVER` acepta memory o redis, y la verificación global de CSRF y el throttle leen `drivers.Sessions`, así que un store construido sobre la tabla debe reemplazar a ese en todas partes antes de ser confiable.

<a name="a-missing-record-stays-unauthenticated"></a>
## Un registro ausente permanece sin autenticar

Una cookie ausente, firma inválida, registro expirado o ID ausente no crean un Subject parcial. El middleware considera la petición no autenticada. La indisponibilidad del backend compartido es distinta: falle la petición dependiente o la salud, sin cambiar silenciosamente a memoria y dividir las instancias.

> [!WARNING]
> El nombre y el alcance de la cookie de sesión no son configurables. El almacenamiento publicado escribe `arandu_session` en la ruta /, para el host que respondió, `HttpOnly` y `SameSite` Lax. El framework detiene el arranque con `SESSION_TTL`, que `SESSION_LIFETIME` reemplazó, y con `SESSION_PATH`, `SESSION_DOMAIN` o `SESSION_SAME_SITE` pidiendo otro alcance.

<a name="identity-changes-require-id-changes"></a>
## Los cambios de identidad exigen cambios de ID

Rote tras el último factor de autenticación y cambios de privilegio. Un segundo factor pendiente no recibe Subject autenticado. Nunca registre IDs de sesión, acepte ID de formulario, derive tenant de un parámetro ni devuelva un error bruto del backend con credenciales de conexión.

`DestroyOthers` tiene ámbito de tenant porque el Subject fiable proporciona tenant e ID del sujeto. No lo sustituya por un barrido del backend basado solo en el ID de cuenta.

El store vuelve a leer los roles de una sesión en cada petición mediante el resolver que el provider instala con `drivers.Sessions.ResolveSubjects`, así que una cuenta degradada pierde su acceso en su siguiente petición. Un resolver que responde `hhttp.ErrSubjectGone` termina la sesión de una cuenta eliminada, cerrada o rechazada, y una sesión con roles sin resolver instalado se rechaza con `hhttp.ErrNoSubjectResolver` en lugar de confiar en ella tal como se escribió.

<a name="prove-the-cross-replica-path"></a>
## Demuestre el recorrido entre réplicas

Construya dos instancias sobre el mismo servidor RESP. Rote una sesión en la primera, reutilice la cookie en la segunda y confirme el mismo tenant y sujeto. Pruebe `SESSION_DRIVER=redis` con `CACHE_STORE=memory`, storage indisponible, driver no compatible, firma alterada, expiración, invalidación del ID antiguo y `DestroyOthers` preservando solo el elegido.

```sh
go test ./tests/Feature -run 'Session|Cache'
```

<a name="authentication-owns-identity-policy-owns-authority"></a>
## La autenticación define identidad; la policy define autoridad

`RequireAuth` carga identidad desde ese almacenamiento. El handler aún llama a `security.Authorize` y lleva el Grant devuelto al repositorio. CSRF vincula mutaciones del navegador al ID actual de sesión, por lo que la rotación también invalida tokens de la frontera anterior.

- [Autenticación](/docs/{{version}}/security/authentication)
- [Policies y Grant](/docs/{{version}}/security/authorization)
- [Protección CSRF](/docs/{{version}}/http/csrf)
- [API de sesión del framework](https://pkg.go.dev/github.com/arandu-io/framework/security)
- [Adapter de sesión RESP de Hesape](https://pkg.go.dev/github.com/arandu-io/hesape/redis)
