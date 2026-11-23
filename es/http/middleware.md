# Middleware

- [Un wrapper, una decisión](#one-wrapper-one-decision)
- [Instale el perímetro de producción](#install-the-production-envelope)
- [Elija el menor ámbito de middleware](#choose-the-narrowest-middleware-scope)
- [El middleware debe preservar el protocolo](#middleware-must-preserve-the-protocol)
- [Demuestre orden y ramas terminales](#prove-order-and-terminal-branches)
- [API de middleware publicada](#released-middleware-api)
- [Coloque la policy en la frontera correcta](#put-policy-at-the-right-boundary)

<a name="one-wrapper-one-decision"></a>
## Un wrapper, una decisión

El middleware HTTP tiene la forma estándar `func(http.Handler) http.Handler`. El framework compone el pipeline global a partir de la configuración y lo ejecuta en un único orden fijo a la ida y en orden inverso a la vuelta; el middleware de ruta envuelve solo la ruta o el grupo que protege.

<a name="install-the-production-envelope"></a>
## Instale el perímetro de producción

El envoltorio de producción ya está instalado. Create compone, a partir de la configuración, la recuperación, los límites de cuerpo, la normalización de proxies confiables, la observabilidad, las cabeceras de seguridad, un throttle global de trescientas peticiones por minuto por sesión o dirección, la verificación de CSRF y la sobrescritura de método. `bootstrap/app.go` ajusta ese pipeline en `WithMiddleware`: añade middleware en uno de tres lugares con nombre y cambia las pocas decisiones que el pipeline deja a la aplicación.

```go example
// Stand-ins: each wraps next with the behavior its name describes.
func requestTimer(next stdhttp.Handler) stdhttp.Handler        { return next }
func noStoreOnBearerLinks(next stdhttp.Handler) stdhttp.Handler { return next }
func tenantBanner(next stdhttp.Handler) stdhttp.Handler         { return next }

func WithPipeline(drivers *foundation.Drivers) *foundation.ApplicationBuilder {
	return foundation.Configure(drivers).
		WithMiddleware(func(m *foundation.Middleware) {
			m.Prepend(requestTimer)
			m.BeforeThrottle(noStoreOnBearerLinks)
			m.Append(tenantBanner)
			m.Throttle(cache.PerMinute(600))
			m.ValidateCSRFTokensExcept("/webhooks/", "/checkout/{token}/sessions")
			m.ImageOrigins("https://images.example.com")
		})
}
```

> [!NOTE]
> `TrustProxies` lee X-Forwarded-For solo de un par que `TRUSTED_PROXIES` lista, recorriendo la cabecera desde la derecha y saltando los saltos confiables, de modo que un valor que el cliente escribió a la izquierda nunca se convierte en su dirección. Todo lo que deriva identidad o un límite de la dirección del cliente se ejecuta después: `BeforeThrottle`, el Throttle global y Append.

<a name="choose-the-narrowest-middleware-scope"></a>
## Elija el menor ámbito de middleware

El orden es fijo: las Sentences de rechazo, Recover, Prepend, `ValidatePostSize` y `LimitBodySize`, `TrustProxies`, Observe, `SecurityHeaders`, `BeforeThrottle`, Throttle, la verificación de CSRF, `OverrideMethod` y Append, justo antes del router. Adjunte la autenticación y los controles específicos de ruta a un `Router.Group`, `Router.Action` o método de verbo. `http.Chain` compone un handler estándar independiente con el primer middleware en el lado más externo.

- Prepend añade middleware justo después de Recover, antes de los límites de cuerpo.
- `BeforeThrottle` también actúa en una petición rechazada y lee la dirección real.
- Append se ejecuta después de `OverrideMethod`, justo antes del router.
- Throttle reemplaza el presupuesto global; Sentences y Observer reemplazan en el mismo lugar las palabras de rechazo y Observe.
- La autenticación pertenece solo a grupos de rutas protegidos.

<a name="middleware-must-preserve-the-protocol"></a>
## El middleware debe preservar el protocolo

Un middleware llama a `next.ServeHTTP` como máximo una vez, retorna después de una respuesta terminal, define las cabeceras antes de la primera escritura y preserva `r.WithContext`. Romper esas reglas puede duplicar trabajo, descartar cancelación e identidad o emitir una respuesta con las cabeceras equivocadas. Recover queda en el lado más externo después de las frases: un panic en un middleware por encima de él escapa sin página.

> [!WARNING]
> El Throttle publicado provoca panic por construcción inválida y falla abierto si falla el almacenamiento, registrando el error. No lo presente como protección fail-closed de login o recuperación sin una policy de aplicación que rechace la indisponibilidad.

<a name="prove-order-and-terminal-branches"></a>
## Demuestre orden y ramas terminales

Despache a través de la aplicación compuesta y verifique dónde se ejecuta un middleware añadido respecto al throttle y a la verificación de CSRF. Pruebe que un par no confiable no puede falsificar una dirección reenviada, que un par confiable sí puede, que un rechazo terminal nunca llega al handler, que la petición 301 en un minuto se rechaza y que un panic escapado se convierte en el estado de producción.

- El orden de entrada y salida es determinista.
- El middleware terminal no llama al handler posterior.
- La cancelación de petición sobrevive al enriquecimiento del contexto.
- Los fallos de proxy y throttle siguen sus políticas documentadas.

<a name="released-middleware-api"></a>
## API de middleware publicada

`foundation.ApplicationBuilder.WithMiddleware` entrega `*foundation.Middleware`, cuyos Prepend, `BeforeThrottle`, Append, Throttle, `ValidateCSRFTokensExcept`, `ImageOrigins`, Observer y Sentences ajustan el pipeline compuesto. `http.Chain` compone handlers estándar. El middleware del framework proporciona Observe, `RequireAuth` y CSRFProtect; el middleware HTTP de Hesape proporciona la normalización de proxies confiables y las cabeceras de seguridad, y el throttle es `hesape/routing/middleware.Throttle` sobre `drivers.Limiter`.

| Superficie | Contrato publicado |
| --- | --- |
| `WithMiddleware` | ajusta el pipeline que compone el framework |
| Prepend / `BeforeThrottle` / Append | los tres lugares donde una aplicación añade middleware |
| Throttle / `ImageOrigins` / `ValidateCSRFTokensExcept` | el presupuesto global, los orígenes de imágenes y las exenciones de CSRF |
| Observer / Sentences | reemplazan en el mismo lugar Observe y las palabras de rechazo |
| `RequireAuth` y CSRFProtect | verificaciones de identidad y mutación del navegador |

<a name="put-policy-at-the-right-boundary"></a>
## Coloque la policy en la frontera correcta

El enrutamiento elige ámbito, la autenticación carga identidad, CSRF protege mutaciones del navegador y rate limit mide una clave con nombre. El middleware compone esas fronteras sin mezclar sus significados.

- [Rutas](/docs/{{version}}/http/routing)
- [Autenticación](/docs/{{version}}/security/authentication)
- [CSRF](/docs/{{version}}/http/csrf)
