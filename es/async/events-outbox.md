# Eventos y outbox transaccional

- [Confirma juntos el estado de negocio y sus hechos](#commit-business-state-and-its-facts-together)
- [Conecta tabla, relay, publicador y lock entre réplicas](#wire-the-table-relay-publisher-and-replica-lock)
- [Almacene el evento en la transacción de la escritura](#store-the-event-inside-the-write-transaction)
- [Procese por el recorrido de publicación de producción](#drain-through-the-production-publisher-path)
- [Prepare cada consumidor para entregas repetidas](#design-every-consumer-for-repeated-delivery)
- [Tenant y autoridad proceden del Grant](#tenant-and-authority-come-from-the-grant)
- [Demuestre rollback, replay y comportamiento entre réplicas](#prove-rollback-replay-and-replica-behavior)
- [Continúe con workers de cola y programación](#continue-with-queue-workers-and-scheduling)

<a name="commit-business-state-and-its-facts-together"></a>
## Confirma juntos el estado de negocio y sus hechos

Usa la outbox duradera para un hecho que debe salir del proceso después de una escritura en la base de datos. `events.Outbox.Store` recibe el mismo contexto y Grant, registra tenant, action y autorizador y rechaza llamadas fuera de `data.Transaction` con `ErrNoTransaction`. El despacho dentro del proceso no cierra la ventana de fallo entre commit y publicación.

<a name="wire-the-table-relay-publisher-and-replica-lock"></a>
## Conecta tabla, relay, publicador y lock entre réplicas

Construya el outbox con `events.NewOutbox` sobre `drivers.DB` en un provider, entregue a `events.NewRelay` una lista `events.Listeners` cuyas entradas construye `events.On`, una por nombre de evento, y liste `events.WithRelay(relay)` en `WithProviders` en lugar de solo el `NewModule` dueño de la tabla. Varias réplicas exigen un lock distribuido: pase `drivers.Locker`, que solo es nil cuando todo store se queda dentro del proceso, y eso solo sirve para un único proceso servidor.

- Migre la tabla de la outbox.
- Inyecte un único Publisher de producción.
- Registra el ciclo de vida del relay.
- Usa un lock compartido entre las réplicas que atienden solicitudes.

<a name="store-the-event-inside-the-write-transaction"></a>
## Almacene el evento en la transacción de la escritura

Escribe el cambio de negocio y almacena su hecho inmutable con Store en la misma transacción SQL. Si una operación falla, ninguna se confirma. Los consumidores pueden recargar por id, pero no deben suponer que el estado mutable aún coincida con el evento ocurrido.

```go example
type Invoices interface {
	MarkPaid(ctx context.Context, grant security.Grant, id string) error
}

func MarkPaid(ctx context.Context, db *data.DB, invoices Invoices, outbox *events.Outbox, grant security.Grant, invoiceID string) error {
	return data.Transaction(ctx, db, func(ctx context.Context) error {
		if err := invoices.MarkPaid(ctx, grant, invoiceID); err != nil {
			return err
		}
		return outbox.Store(ctx, grant, []events.Event{{
			Name: "invoice.paid", Aggregate: "invoice", AggregateID: invoiceID,
			Payload: map[string]any{"invoice_id": invoiceID}, OccurredAt: time.Now().UTC(),
		}})
	})
}
```

<a name="drain-through-the-production-publisher-path"></a>
## Procese por el recorrido de publicación de producción

`Relay.Drain` procesa las filas pendientes a través del mismo publisher que el bucle que sirve. Inspeccione las operaciones de pendientes, aparcados, reintento, retraso y marcado; por defecto el intervalo es de un segundo, el lote de 100, diez intentos y un lock de relay de treinta segundos, salvo configuración. Para un consumidor HTTP externo, componga `webhook.Manager` con la base de datos, el resolutor de secretos y opciones acotadas, y adapte el Publisher del outbox existente en lugar de crear un segundo outbox. Registre el módulo de webhook de Hesape con `foundation.Component` en `WithProviders`: Create recoge sus handlers de jobs y las exenciones de CSRF de sus receptores, así que la aplicación no escribe ni una llamada a `RegisterJobHandlers` ni una línea `ValidateCSRFTokensExcept` para ellos, y llamar también a `RegisterJobHandlers` registra cada job dos veces y detiene el arranque. Su prefijo de gestión y su prefijo de recepción no pueden superponerse, porque el prefijo de recepción está exento de la verificación de CSRF.

- Procese los registros mediante el Publisher de producción.
- Monitoriza eventos pendientes y retraso.
- Aparte el trabajo al agotar los intentos.
- Reintente solo después de corregir el consumidor o el transporte.
- Registre cada endpoint de webhook con una URL HTTP o HTTPS absoluta y un secreto de al menos 32 bytes.

<a name="design-every-consumer-for-repeated-delivery"></a>
## Prepare cada consumidor para entregas repetidas

La entrega ocurre al menos una vez. Una publicación puede tener éxito antes de que falle la confirmación, por lo que el mismo `Stored.ID` puede llegar de nuevo. Los consumidores deduplican por ese ID estable y hacen idempotentes los efectos; las filas estacionadas permanecen como evidencia. El despacho de webhook crea atómicamente un snapshot de entrega y un job de cola para cada evento y endpoint. Los fallos de transporte, 408, 425, 429 y 5xx se reintentan; las respuestas 4xx comunes son permanentes y se respeta Retry-After.

> [!WARNING]
> Una publicación directa después del commit puede perderse permanentemente si el proceso termina entre ambas operaciones. Use la outbox siempre que ese hecho deba salir del proceso.

<a name="tenant-and-authority-come-from-the-grant"></a>
## Tenant y autoridad proceden del Grant

> [!WARNING]
> No coloque en un evento un tenant elegido por quien llama para confiar después en él. Store persiste tenant, action y `authorized_by` desde el Grant; la entrega de webhook hereda tenant y action de ese Grant. Las peticiones POST firmadas llevan X-Arandu-Signature, X-Arandu-Timestamp y X-Arandu-Delivery-ID; HMAC-SHA256 cubre timestamp, ID de entrega y cuerpo. Mantenga la resolución del endpoint detrás del cliente protegido contra SSRF e imponga el límite de payload de un mebibyte, límites de respuesta y redirección, timeout y rotación de secreto.

Los payloads contienen identificadores y hechos públicos inmutables, no credenciales, Grants sin procesar ni registros completos con secretos. Una firma de webhook válida prueba la integridad del mensaje, no permiso para leer datos de otro tenant.

<a name="prove-rollback-replay-and-replica-behavior"></a>
## Demuestre rollback, replay y comportamiento entre réplicas

Devuelva un error centinela desde `data.Transaction` y compruebe que no queda ninguna fila de negocio ni de outbox. Luego haga que la publicación tenga éxito y falle la confirmación, drene dos veces y compruebe la deduplicación del consumidor y el bloqueo entre réplicas. Para webhooks, verifique la firma byte a byte, la rotación del secreto anterior, un job por evento y endpoint, idempotencia del despacho duplicado, rechazo de SSRF, límites de respuesta, clasificación de reintentos, claims cercados, estado de auditoría y limpieza por retención.

- El rollback no deja un registro en la outbox.
- El replay mantiene un único efecto externo.
- Solo quien obtiene el lock procesa los eventos.
- El retraso de eventos retenidos aparece en los diagnósticos.
- Un ID de entrega de webhook duplicado nunca crea un segundo efecto externo.

<a name="continue-with-queue-workers-and-scheduling"></a>
## Continúe con workers de cola y programación

Las colas ejecutan intentos duraderos de webhook, la programación descubre trabajo, la observabilidad expone el retraso del relay y la autorización define el límite del Grant registrado. Los webhooks salientes duraderos están separados de los callbacks opcionales del ciclo de vida Pusher en Joaju.

- [Patrones de comunicación](/docs/{{version}}/architecture/contracts-adapters)
- [Colas](/docs/{{version}}/async/queues)
- [`WebSocket` de Joaju](/docs/{{version}}/async/joaju)
- [Observabilidad](/docs/{{version}}/operations/observability)
- [Autorización](/docs/{{version}}/security/authorization)
