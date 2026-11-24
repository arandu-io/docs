# Consola de debug y tracing

- [Inspeccione una petición sin hacer universal el tracing](#inspect-one-request-without-making-tracing-universal)
- [Elija acceso de desarrollo o un secreto explícito de producción](#choose-development-access-or-an-explicit-production-secret)
- [Registre una petición y consúltela por su ID](#record-a-request-and-read-it-by-request-id)
- [Lea los hallazgos antes de la timeline sin procesar](#read-findings-before-the-raw-timeline)
- [Entienda por qué una petición no tiene registro](#know-why-a-request-has-no-record)
- [Minimice la ventana de datos diagnósticos](#minimize-the-diagnostic-data-window)
- [Pruebe los estados abierto y cerrado](#test-both-the-open-and-closed-states)
- [Reproduzca con ámbito reducido y elimine el acceso](#reproduce-narrowly-and-remove-access)
- [Conecte traces con logs, salud y llamadas externas](#connect-traces-to-logs-health-and-outbound-calls)

<a name="inspect-one-request-without-making-tracing-universal"></a>
## Inspeccione una petición sin hacer universal el tracing

La consola de debug es un recorder acotado, dentro del proceso, de Collectors de peticiones recientes: distribución de la timeline, SQL con origen y argumentos, volcados, eventos, vistas, HTTP externo, consultas lentas y sospechas de N+1. No es APM alojado, traces duraderos, historial de métricas ni consulta entre réplicas.

<a name="choose-development-access-or-an-explicit-production-secret"></a>
## Elija acceso de desarrollo o un secreto explícito de producción

Desarrollo monta el recorder y `/_arandu/debug`. Fuera de él, deje `ARANDU_TRACING_SECRET` vacío salvo necesidad de tracing controlado en producción; en ese caso, use al menos 16 bytes, TLS y X-Arandu-Trace solo en peticiones y lecturas de consola investigadas.

- Mantenga el tracing de producción desactivado por defecto.
- Trate el secreto de tracing como acceso diagnóstico entre tenants.
- Valide `ARANDU_EDITOR` al arrancar si se habilitan enlaces al código.
- Durante un incidente, acceda a una réplica conocida.

<a name="record-a-request-and-read-it-by-request-id"></a>
## Registre una petición y consúltela por su ID

Haga una petición, guarde X-Request-ID y consulte la lista o el detalle mediante `aru trace`. La CLI lee la misma ruta JSON protegida de la consola del navegador y permite seleccionar otro host explícitamente.

```bash
set -euo pipefail
aru dev &
DEV_PID=$!
trap 'kill "$DEV_PID"' EXIT
REQUEST_ID=$(curl --retry 20 --retry-connrefused -si http://127.0.0.1:8080/ | awk -F': ' 'tolower($1) == "x-request-id" { gsub("\\r", "", $2); print $2; exit }')
test -n "$REQUEST_ID" || { echo "response has no X-Request-ID" >&2; exit 1; }
aru trace
aru trace "$REQUEST_ID"
```

<a name="read-findings-before-the-raw-timeline"></a>
## Lea los hallazgos antes de la timeline sin procesar

Lea los hallazgos de lentitud y N+1 antes de la timeline sin procesar. Inspeccione origen y duración de SQL, llamadas externas, eventos y volcados temporales sin secretos. Cinco consultas idénticas son una pista, no una prueba de repetición indebida.

- Encuentre la petición por ID.
- Inspeccione los hallazgos antes de los eventos individuales.
- Instrumente dependencias mediante el Collector.
- Use Gauges compartidos solo para valores actuales del proceso.

<a name="know-why-a-request-has-no-record"></a>
## Entienda por qué una petición no tiene registro

La ausencia de registro puede indicar tracing desactivado, header ausente, pérdida del Collector en middleware, atención por otra réplica o descarte por el búfer circular limitado. Las dependencias no instrumentadas quedan sin atribución; los registros protegidos o ausentes devuelven 404 deliberadamente. El recorder usa memoria local y desaparece al reiniciar. Rastrear una petición de producción no recoge otras peticiones. Arandu no incluye un exportador de telemetría: las trazas y métricas duraderas pertenecen a un backend que usted mantiene junto a la aplicación, y los logs estructurados llevan `request_id` y tenant para correlacionar ambos.

<a name="minimize-the-diagnostic-data-window"></a>
## Minimice la ventana de datos diagnósticos

> [!WARNING]
> Los registros de consola pueden contener argumentos SQL, volcados, payloads de eventos, rutas y datos de tenants. Nunca exponga la ruta mediante un proxy de conveniencia, registre el secreto, lo incluya en la URL ni habilite el header indiscriminadamente. Rote o elimine el acceso tras el incidente.

<a name="test-both-the-open-and-closed-states"></a>
## Pruebe los estados abierto y cerrado

Pruebe índice y detalle en desarrollo, 404 con producción deshabilitada, 404 con secreto incorrecto, header aceptado en tiempo constante, recopilación por petición, hallazgos, descarte del búfer circular e independencia de probes live y health respecto al acceso a consola.

- Solo se registran peticiones de producción explícitamente rastreadas.
- Las credenciales incorrectas no revelan la consola.
- Se elimina el registro más antiguo al alcanzar la capacidad.
- La salud sigue siendo un contrato operativo independiente.

<a name="reproduce-narrowly-and-remove-access"></a>
## Reproduzca con ámbito reducido y elimine el acceso

Use primero logs y salud. Si hace falta, dirija una petición controlada a una réplica conocida, consúltela de inmediato, preserve evidencia mínima sin datos sensibles y elimine o rote el secreto. Reiniciar borra la evidencia, pero no corrige la dependencia.

<a name="connect-traces-to-logs-health-and-outbound-calls"></a>
## Conecte traces con logs, salud y llamadas externas

Los logs estructurados proporcionan correlación duradera, los clientes instrumentados explican esperas externas, el despliegue posee probes y la documentación de consultas explica hallazgos SQL.

- [Logs](/docs/{{version}}/operations/logging-observability)
- [Cliente HTTP](/docs/{{version}}/services/http-client)
- [Despliegue](/docs/{{version}}/operations/deployment)
- [Consultas](/docs/{{version}}/database/queries)
