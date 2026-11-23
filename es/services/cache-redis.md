# Caché y Redis

- [La aplicación usa un repositorio con ámbito de tenant](#the-application-uses-a-tenant-scoped-repository)
- [Comience con almacenamiento dentro del proceso](#begin-with-the-in-process-store)
- [Almacene, memorice, cuente e invalide](#store-remember-count-and-invalidate)
- [Una ausencia en caché es esperable; un error del almacenamiento, no](#a-cache-miss-is-expected-a-store-error-is-not)
- [Pruebe aislamiento y expiración mediante el repositorio](#test-isolation-and-expiry-through-the-repository)
- [Fronteras del repositorio y almacenamiento](#repository-and-store-boundaries)
- [Un servidor compatible con el protocolo es una elección de despliegue](#a-protocol-compatible-server-is-a-deployment-choice)

> [!WARNING]
> **Limitaciones actuales.** `foundation.NewDrivers` construye `drivers.Store`, el `cache.Repository` con alcance de tenant sobre el store que nombra `CACHE_STORE`, y el esqueleto publicado entrega el store subyacente a los códigos de un solo uso y al segundo factor en lugar de usar la caché a través del repositorio. Pruebe el espacio de nombres basado en Grant en la aplicación que usa la caché a través de él.

<a name="the-application-uses-a-tenant-scoped-repository"></a>
## La aplicación usa un repositorio con ámbito de tenant

Repository de caché añade JSON, expiración, eventos y claves prefijadas por tenant sobre Store. La aplicación mantiene el repositorio, no el backend sin procesar. Cada operación ordinaria recibe Grant, por lo que el tenant proviene de autorización, no de un campo de petición ni de un prefijo manual.

<a name="begin-with-the-in-process-store"></a>
## Comience con almacenamiento dentro del proceso

El almacenamiento array ofrece a desarrollo y pruebas unitarias un backend con la misma superficie Repository. Constrúyalo una vez; use Put y recupere con el mismo Grant. Otro tenant genera una clave física distinta aunque el texto de la clave de aplicación sea idéntico.

```go example
type Dashboard struct{ Total int }

func CacheDashboard(ctx context.Context, g security.Grant, value Dashboard) (Dashboard, error) {
	repository := cache.New(cache.NewArrayStore())
	if err := repository.Put(ctx, g, "dashboard", value, 5*time.Minute); err != nil {
		return Dashboard{}, err
	}
	return cache.Get[Dashboard](ctx, repository, g, "dashboard")
}
```

<a name="store-remember-count-and-invalidate"></a>
## Almacene, memorice, cuente e invalide

Put sustituye un valor con duración acotada; Forever es la elección explícita sin plazo; Has, Missing, Forget y Touch gestionan presencia; Increment y Decrement actualizan contadores. Get y Remember genéricos decodifican al tipo solicitado, dejando aserciones de tipo fuera de handlers.

En una aplicación el repositorio ya está construido: `foundation.NewDrivers` responde `drivers.Store`, el `cache.Repository` sobre el store que nombra `CACHE_STORE`, y `drivers.Cache` guarda todos los stores por nombre. memory siempre existe; redis existe una vez enlazado el conector `hesape/redis` en `bootstrap/app.go` y cuando `REDIS_URL` nombra el endpoint. Un provider entrega `drivers.Store` a los servicios que usan caché.

<a name="a-cache-miss-is-expected-a-store-error-is-not"></a>
## Una ausencia en caché es esperable; un error del almacenamiento, no

Trate la ausencia según el paquete y calcule el valor de origen una vez. Propague errores de conexión, codificación y decodificación sin tratarlos como cache miss: reconstruir silenciosamente durante una caída puede multiplicar consultas hasta tumbar el servicio principal. Nunca almacene resultados fallidos ni parcialmente autorizados.

Un Grant de solo lectura, emitido mientras un operador solo ve la cuenta de otra persona, lee y nunca escribe: todo método que guarda, elimina, vacía o bloquea una entrada lo rechaza con `auth.ErrReadOnly`, y Remember y Flexible responden lo que calculan sin guardarlo. Pregunte a `Grant.ReadOnly` para dejar una escritura fuera de una página así en lugar de dejar que falle.

<a name="test-isolation-and-expiry-through-the-repository"></a>
## Pruebe aislamiento y expiración mediante el repositorio

Escriba la misma clave lógica con Grants de dos tenants y confirme que cada uno lee solo su valor. Avance o limite el tiempo para demostrar expiración, verifique que Remember llama al loader solo ante ausencia y ejecute el contrato compartido contra el driver de red antes de elegirlo en producción.

<a name="repository-and-store-boundaries"></a>
## Fronteras del repositorio y almacenamiento

Los métodos de Repository cubren valores, lotes, contadores, namespaces, tags cuando se admiten y bloqueos contra trabajo superpuesto. Store cubre persistencia física. `WithoutOverlapping` recibe Grant y prefija el bloqueo por tenant como otras operaciones. `FlushLocks` es la excepción global: no recibe Grant porque limpia el namespace de bloqueos del almacenamiento. Eso no convierte stores sin procesar en API de aplicación.

<a name="a-protocol-compatible-server-is-a-deployment-choice"></a>
## Un servidor compatible con el protocolo es una elección de despliegue

El driver de red habla el protocolo común de servidores compatibles con Redis sin incluir comandos específicos de producto en el contrato de caché. Elija conexión y store en configuración, mantenga llamadas de aplicación y dimensione TTL según la actualidad necesaria de los datos de origen. Los registros duraderos y sesiones no se vuelven caché solo porque el backend retenga bytes.
