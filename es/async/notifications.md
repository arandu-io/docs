# Notificaciones

- [Una notificación selecciona canales para un mensaje público](#a-notification-selects-channels-around-one-public-message)
- [Asume la responsabilidad de preferencias, direcciones y transportes de los destinatarios](#own-recipient-preferences-addresses-and-transports)
- [Construye un mensaje sin almacenar credenciales](#build-a-message-without-storing-credentials)
- [Separa la entrega inmediata de la duradera](#separate-immediate-and-durable-delivery)
- [Trate las entregas parciales entre canales como estado explícito](#treat-partial-channel-delivery-as-explicit-state)
- [La dirección del destinatario es un dato protegido](#a-recipient-address-is-protected-data)
- [Pruebe preferencias, tenant y reintentos](#test-preference-tenant-and-retry-paths)
- [Use correo, colas y outbox respetando sus fronteras](#use-mail-queues-and-outbox-as-their-real-boundaries)

<a name="a-notification-selects-channels-around-one-public-message"></a>
## Una notificación selecciona canales para un mensaje público

Una notificación nombra sus canales en `Via` y construye un mensaje para cada uno. El esqueleto conecta un `Notifier` con el canal de base de datos en `app/Providers/AppServiceProvider.go`, y su listener `NotifyNoteAuthor` envía la notificación de ejemplo `NotePublished` cuando el relay de la outbox entrega `note.published`. `aru make:notification` genera una notificación para los canales que usted elija. Agregue correo, preferencias del destinatario y entrega en cola de forma explícita, como describe el resto de esta página.

<a name="own-recipient-preferences-addresses-and-transports"></a>
## Asume la responsabilidad de preferencias, direcciones y transportes de los destinatarios

Resuelva un destinatario fiable mediante el servicio de la aplicación, lea preferencias dentro del ámbito del tenant, configure cada transporte y decida qué entregas se encolarán. No acepte direcciones, tenant ni canales habilitados desde un payload de job no fiable.

- Cargue el destinatario mediante un Grant.
- Aplica las preferencias de opt-out en el momento de entrega.
- Inyecte transportes concretos.
- Elija una clave de idempotencia por canal.

<a name="build-a-message-without-storing-credentials"></a>
## Construye un mensaje sin almacenar credenciales

Asigna una Key estable a la notificación, devuelve los nombres de sus canales en Via e implementa la interfaz de representación de cada canal seleccionado. Via por sí solo no proporciona un cuerpo de correo ni de base de datos. Pasa el valor completo al Notifier inyectado, que resuelve cada canal conectado y sigue recopilando fallos.

```go example
type InvoicePaid struct {
	InvoiceID string
	URL       string
}

func (InvoicePaid) Key() notifications.Key { return "billing.invoice-paid" }
func (InvoicePaid) Via(notifications.Notifiable) []notifications.ChannelName {
	return []notifications.ChannelName{
		notifications.ChannelMail,
		notifications.ChannelDatabase,
	}
}

func (n InvoicePaid) ToMail(notifications.Notifiable) messages.Mail {
	return messages.NewMail().
		Subject("Your invoice is paid").
		Line("We received your payment.").
		Action("View invoice", n.URL)
}

func (n InvoicePaid) ToDatabase(notifications.Notifiable) messages.Database {
	return messages.NewDatabase(map[string]any{"invoice_id": n.InvoiceID})
}

func NotifyPaid(ctx context.Context, notifier *notifications.Notifier, grant security.Grant, recipient notifications.Notifiable, invoiceID, invoiceURL string) error {
	notice := InvoicePaid{InvoiceID: invoiceID, URL: invoiceURL}
	return notifier.Send(ctx, grant, recipient, notice)
}
```

<a name="separate-immediate-and-durable-delivery"></a>
## Separa la entrega inmediata de la duradera

Resuelve las preferencias justo antes de la entrega, renderiza resources públicos específicos por canal, encola transportes de red lentos, asocia un id estable de entrega y registra el resultado de cada canal sin considerar el éxito de uno como prueba de la entrega de todos.

- Mantenga la selección de canales dentro del ámbito del tenant.
- Usa mail para la representación de correo.
- Usa queue para reintentos duraderos.
- Registre un resultado por destinatario y canal.

<a name="treat-partial-channel-delivery-as-explicit-state"></a>
## Trate las entregas parciales entre canales como estado explícito

Un canal puede tener éxito mientras otro falla. Defina el estado de reintento por canal, evite repetir efectos ya completados y distinga dirección inválida, cancelación explícita, rechazo e indisponibilidad del proveedor. El paquete de componentes, por sí solo, no proporciona este ciclo completo.

<a name="a-recipient-address-is-protected-data"></a>
## La dirección del destinatario es un dato protegido

> [!WARNING]
> Trate el correo, el teléfono, los tokens de dispositivo y el contenido del mensaje como datos protegidos del tenant. Nunca registre credenciales ni códigos de recuperación, permita que el payload sustituya al destinatario ni ignore una preferencia actual de cancelación al reintentar.

La autorización decide el acceso al evento; la preferencia de notificación decide si enviarlo y por dónde. Una no sustituye a la otra.

<a name="test-preference-tenant-and-retry-paths"></a>
## Pruebe preferencias, tenant y reintentos

Verifique que todo nombre devuelto por Via tenga su método de representación obligatorio. Luego ejercite canales habilitados y deshabilitados, destinatario de otro tenant, dirección inválida, indisponibilidad del proveedor tras el éxito de otro canal, entrega duplicada de la cola, logs sin datos sensibles y cambios de preferencia entre encolado y ejecución.

- El opt-out impide el envío por el canal correspondiente.
- No se pueden seleccionar destinatarios de otros tenants.
- La entrega duplicada es idempotente.
- El fallo parcial preserva el estado de los canales que tuvieron éxito.

<a name="use-mail-queues-and-outbox-as-their-real-boundaries"></a>
## Use correo, colas y outbox respetando sus fronteras

Mail formatea correos, las colas proporcionan entrega duradera, los eventos y la outbox persisten hechos y la autorización protege la consulta de destinatarios.

- [Correo](/docs/{{version}}/services/mail)
- [Colas](/docs/{{version}}/async/queues)
- [Eventos y outbox](/docs/{{version}}/async/events-outbox)
- [Autorización](/docs/{{version}}/security/authorization)
