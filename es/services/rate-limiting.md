# Limitación de peticiones

- [Un limiter mide una clave explícita en un almacenamiento](#a-limiter-meters-one-explicit-key-in-one-store)
- [Construya el limiter y limite una ruta](#construct-the-limiter-and-throttle-a-route)
- [Elija claves y política de rechazo según la amenaza](#choose-keys-and-refusal-policy-by-threat)
- [El throttle publicado permite acceso si falla el almacenamiento](#the-released-throttle-fails-open-on-store-errors)
- [Demuestre ventanas, concurrencia y comportamiento ante fallos](#prove-windows-concurrency-and-outage-behavior)
- [Limiter de caché y middleware publicados](#released-cache-limiter-and-middleware)
- [Coloque el throttle antes del trabajo costoso de autenticación](#place-the-throttle-before-expensive-authentication-work)

<a name="a-limiter-meters-one-explicit-key-in-one-store"></a>
## Un limiter mide una clave explícita en un almacenamiento

El camino publicado actual es `cache.NewRateLimiter` más el middleware de enrutamiento Throttle. El framework ya construye un limiter, `drivers.Limiter`, que cuenta en el store que nombra `CACHE_STORE`, y el pipeline global limita cada petición a trescientas por minuto por sesión o dirección; `WithMiddleware` cambia ese presupuesto con `Middleware.Throttle`. El framework `v0.41.0` eliminó la antigua superficie de limiter del framework; no copie ejemplos que la importen.

<a name="construct-the-limiter-and-throttle-a-route"></a>
## Construya el limiter y limite una ruta

Pase `drivers.Limiter` a Throttle con un límite positivo, una ventana, una función de clave y un callback de rechazo en una ruta que necesite un presupuesto más estricto que el global. El middleware entra en panic al construirse cuando faltan colaboradores o la política no es positiva. Las pantallas de inicio de sesión publicadas cuentan sus formularios de código en el mismo limiter.

```go example
import middleware "github.com/arandu-io/hesape/routing/middleware"

func LoginThrottle(drivers *foundation.Drivers, clientKey func(*stdhttp.Request) string) func(stdhttp.Handler) stdhttp.Handler {
	limiter := drivers.Limiter
	throttle := middleware.Throttle(limiter, cache.PerMinute(5),
		func(r *stdhttp.Request) string { return "login:" + clientKey(r) },
		func(w stdhttp.ResponseWriter, r *stdhttp.Request, status int, message string) {
			stdhttp.Error(w, message, status)
		},
	)
	return throttle
}
```

<a name="choose-keys-and-refusal-policy-by-threat"></a>
## Elija claves y política de rechazo según la amenaza

En fronteras de credenciales, combine cuenta y red fiable en las claves; normalícelas, preserve tenant cuando se conozca, envíe Retry-After y use almacenamiento compartido para límites entre instancias. Separe límites generales de abuso de la política estricta de rechazo de autenticación.

- Construya una vez al arrancar.
- Normalice claves acotadas.
- Use almacenamiento compartido en producción.
- Devuelva un rechazo 429 explícito.

<a name="the-released-throttle-fails-open-on-store-errors"></a>
## El throttle publicado permite acceso si falla el almacenamiento

El Throttle publicado registra errores del almacenamiento y llama al handler: falla abierto. Eso mantiene disponibilidad en endpoints comunes, pero no representa protección fail-closed de login, recuperación de contraseña ni otra frontera de seguridad sin una policy de aplicación que rechace acceso si cae el almacenamiento.

> [!WARNING]
> Los almacenamientos en memoria dividen contadores entre réplicas y se reinician al arrancar. Un fallo del almacenamiento compartido y ataques por cardinalidad de claves exigen decisiones explícitas de salud, alertas y rechazo.

<a name="prove-windows-concurrency-and-outage-behavior"></a>
## Demuestre ventanas, concurrencia y comportamiento ante fallos

Pruebe peticiones hasta el límite, rechazo siguiente, Retry-After, vencimiento de ventana, incrementos concurrentes, dos instancias sobre el mismo almacenamiento, headers de proxy no fiables, normalización de claves, fallo de almacenamiento y panic al arrancar con configuración inválida.

- El límite compartido abarca réplicas.
- La concurrencia no supera la cantidad prevista.
- Un fallo de almacenamiento sigue la política declarada.
- El hashing de credenciales nunca se alcanza tras un rechazo.

<a name="released-cache-limiter-and-middleware"></a>
## Limiter de caché y middleware publicados

`cache.NewRateLimiter` proporciona el contador actual basado en almacenamiento; `routing/middleware.Throttle` lo adapta a HTTP. La aplicación elige claves, límites, almacenamiento compartido, representación de rechazo y política de seguridad fail-open o fail-closed.

| Superficie | Contrato publicado |
| --- | --- |
| `cache.NewRateLimiter` | decisiones de ventana fija basadas en almacenamiento |
| `middleware.Throttle` | adapter de clave HTTP y rechazo |
| Retry-After | orientación de reintentos para el cliente |
| almacenamiento compartido de caché | estado de contadores entre instancias |

<a name="place-the-throttle-before-expensive-authentication-work"></a>
## Coloque el throttle antes del trabajo costoso de autenticación

El middleware controla el orden, autenticación y recuperación definen amenazas sensibles, la caché elige estado compartido y el hashing explica el trabajo costoso protegido por el limiter de credenciales.

- [Middleware](/docs/{{version}}/http/middleware)
- [Autenticación](/docs/{{version}}/security/authentication)
- [Caché y Redis](/docs/{{version}}/services/cache-redis)
- [Hashing](/docs/{{version}}/security/hashing)
