# Logs y observabilidad

- [Escriba eventos estructurados mediante el contexto de petición](#write-structured-events-through-the-request-context)
- [Defina el nivel y preserve el orden de los middlewares](#set-the-level-and-preserve-middleware-order)
- [Emita un evento de dominio que pueda buscarse](#emit-one-searchable-domain-event)
- [Instrumente esperas externas y trabajo de fondo](#instrument-outbound-waits-and-background-work)
- [Haga visible la ausencia de telemetría](#make-missing-telemetry-visible)
- [Registre identificadores, nunca credenciales ni autoridad](#log-identifiers-never-credentials-or-authority)
- [Verifique los campos en la frontera del handler](#assert-fields-at-the-handler-boundary)
- [Diagnostique desde la última señal duradera](#diagnose-from-the-last-durable-signal)
- [Abra la timeline de petición solo cuando haga falta](#open-the-request-timeline-only-when-needed)

<a name="write-structured-events-through-the-request-context"></a>
## Escriba eventos estructurados mediante el contexto de petición

Create construye un único logger slog raíz a partir de la configuración de logging: la salida estándar en `LOG_LEVEL`, en `LOG_FORMAT`, que es texto en desarrollo y JSON en los demás casos cuando está vacío. Cada línea lleva el entorno como channel y ninguna posición de código. Observe deriva un logger de petición con `request_id`, method y path y, al terminar, registra estado, duración, bytes, número de consultas y tiempo de SQL.

<a name="set-the-level-and-preserve-middleware-order"></a>
## Defina el nivel y preserve el orden de los middlewares

Defina un `LOG_LEVEL` publicado y, cuando el valor predeterminado del entorno no sirva al colector, `LOG_FORMAT`. El pipeline compuesto ya ejecuta Observe después de Recover, de los límites de cuerpo y de `TrustProxies`. Pase `context.Context` por controllers, servicios y trabajo saliente para que `observability.Log` recupere el logger y el Collector del alcance. `LOG_SAMPLING` se lee en `Settings.Logging`, pero no cambia ningún logger, así que no reduce el volumen de logs. La exportación y la retención corresponden a la plataforma que recoge el stdout.

- Use el logger raíz configurado.
- Preserve el contexto de petición.
- Mantenga la recuperación fuera de la observación.
- Envíe stdout al colector de la plataforma.

<a name="emit-one-searchable-domain-event"></a>
## Emita un evento de dominio que pueda buscarse

Emita texto estable de eventos con identificadores estructurados mediante `observability.Log(ctx)`. El logger siempre está disponible; construir otro logger global pierde campos de petición y el nivel configurado.

```go example
type InvoiceSender interface {
	Send(ctx context.Context, id string) error
}

type InvoiceService struct{ sender InvoiceSender }

func (s *InvoiceService) Send(ctx context.Context, id string) error {
	observability.Log(ctx).Info("invoice send requested", "invoice_id", id)
	return s.sender.Send(ctx, id)
}
```

<a name="instrument-outbound-waits-and-background-work"></a>
## Instrumente esperas externas y trabajo de fondo

Use `observability.Client` para llamadas externas acotadas, `WithLogger` para trabajo de fondo controlado y el Recorder de la aplicación en opciones de cola o programador que lo acepten. La plataforma gestiona recopilación de stdout, retención, indexación, rotación y alertas.

- Correlacione eventos mediante `request_id`.
- Defina timeout para clientes externos.
- Mantenga explícito el ámbito de los loggers de fondo.
- Trate la retención externa como configuración de despliegue.

<a name="make-missing-telemetry-visible"></a>
## Haga visible la ausencia de telemetría

Perder el contexto utiliza el logger predeterminado y pierde correlación. Un cliente sin instrumentar omite tiempo externo y puede esperar indefinidamente sin timeout. Los métodos del Collector deshabilitado son no-ops intencionales, no prueba de instrumentación ejecutada. Una consulta en el agregador de logs que lee `source` no encuentra nada, porque las líneas llevan `channel` en su lugar.

<a name="log-identifiers-never-credentials-or-authority"></a>
## Registre identificadores, nunca credenciales ni autoridad

> [!WARNING]
> Nunca registre contraseñas, tokens, cookies, payloads de sesión, `APP_KEY`, secretos de tracing, cuerpos de petición, registros sin límite ni detalles completos de Grant. Los argumentos SQL y volcados son datos diagnósticos sensibles entre tenants; producción rechaza el nivel debug.

<a name="assert-fields-at-the-handler-boundary"></a>
## Verifique los campos en la frontera del handler

Capture un `slog.Handler`, haga peticiones con ids válidos y malformados, verifique los campos de finalización, ejercite una llamada saliente instrumentada, compruebe JSON en producción y texto en desarrollo, pruebe que `LOG_FORMAT` sobrescribe el valor predeterminado del entorno y que el filtro de nivel se cumple, y verifique que `LOG_SAMPLING` no tiene efecto.

- Los IDs válidos de petición se devuelven y registran.
- Los IDs malformados se sustituyen.
- El tiempo externo llega al Collector.
- `LOG_SAMPLING` no se anuncia como efectivo.

<a name="diagnose-from-the-last-durable-signal"></a>
## Diagnostique desde la última señal duradera

Correlacione `request_id` en los logs estructurados duraderos, reproduzca con tracing autorizado solo cuando sea seguro y restaure un `LOG_LEVEL` compatible si la configuración impide el arranque. Repare un colector externo en la plataforma; la aplicación no tiene un archivo local de log que recuperar.

<a name="open-the-request-timeline-only-when-needed"></a>
## Abra la timeline de petición solo cuando haga falta

La consola diagnóstica detalla una petición, los clientes HTTP aportan tiempos externos, el despliegue posee la recopilación de logs y la configuración define los controles efectivos.

- [Observabilidad diagnóstica](/docs/{{version}}/operations/observability)
- [Cliente HTTP](/docs/{{version}}/services/http-client)
- [Despliegue](/docs/{{version}}/operations/deployment)
- [Configuración](/docs/{{version}}/architecture/configuration)
