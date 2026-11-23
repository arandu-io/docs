# Correo

- [El mailable posee la representación; el transporte posee la entrega](#a-mailable-owns-representation-and-the-transport-owns-delivery)
- [Configure un mailer y el origen público al arrancar](#configure-one-mailer-and-public-origin-at-boot)
- [Construya y envíe un mensaje tipado](#build-and-send-one-typed-message)
- [Encole entregas lentas con procedencia de tenant](#queue-slow-delivery-with-tenant-provenance)
- [Separe fallos de renderizado, adjuntos y transporte](#keep-render-attachment-and-transport-failures-distinct)
- [El correo es una frontera de divulgación externa](#mail-is-an-external-disclosure-boundary)
- [Capture el contrato MIME exacto](#capture-the-exact-mime-contract)
- [Use correo dentro de flujos de aplicación vinculados a un propósito](#use-mail-inside-purpose-bound-application-flows)

> [!WARNING]
> **Limitaciones actuales.** El componente de correo de Hesape y el puente del framework renderizan y entregan los mensajes documentados. Almacenamiento de adjuntos y autorización asíncrona no están completos de extremo a extremo; adjunte solo datos ya autorizados y encole correo con procedencia fiable persistida.

<a name="a-mailable-owns-representation-and-the-transport-owns-delivery"></a>
## El mailable posee la representación; el transporte posee la entrega

Hesape :hesapeVersion posee mensajes y transportes de correo; `github.com/arandu-io/framework/mail` :frameworkVersion proporciona el puente de compatibilidad usado por la aplicación. El ejemplo siguiente usa `Address{Email: ...}` y `ToAddress` de ese puente; Mailer de nivel inferior de Hesape usa To.

<a name="configure-one-mailer-and-public-origin-at-boot"></a>
## Configure un mailer y el origen público al arrancar

`foundation.NewDrivers` construye `drivers.Mail` a partir de `MAIL_URL`, `MAIL_FROM_ADDRESS` y `MAIL_FROM_NAME`. `smtp://` hace el upgrade con STARTTLS cuando el servidor lo ofrece, `smtps://` habla TLS implícito en el puerto 465, y `resend://` y `sendgrid://` llegan a esos proveedores. `log://` escribe el mensaje completo en el log y se rechaza en producción, igual que un `MAIL_URL` sin definir, que significa lo mismo. Un provider inyecta `drivers.Mail` en servicios o handlers encolados en lugar de leer el entorno o abrir SMTP desde un controller.

- Use la identidad From configurada.
- Construya enlaces mediante el `APP_URL` fiable.
- Mantenga secretas las credenciales de transporte.
- Elija entre entrega inmediata y duradera.

<a name="build-and-send-one-typed-message"></a>
## Construya y envíe un mensaje tipado

Construya un mensaje tipado con datos públicos de aplicación y envíelo mediante el transporte inyectado con contexto. El controller o servicio de dominio proporcionan destinatario y propósito; el transporte nunca descubre cuentas ni policy de tenant.

```go example
type VerifyEmail struct{ Code string }

func (VerifyEmail) Envelope() mail.Envelope {
	return mail.Envelope{Subject: "Verify your address"}
}
func (m VerifyEmail) Content() mail.Content {
	return mail.Content{Text: "Use this code: " + m.Code}
}

type Recipient struct{ Email, Name string }

func SendVerification(ctx context.Context, mailer *mail.Mailer, recipient Recipient, code string) error {
	return mailer.ToAddress(mail.Address{
		Email: recipient.Email, Name: recipient.Name,
	}).Send(ctx, VerifyEmail{Code: code})
}
```

<a name="queue-slow-delivery-with-tenant-provenance"></a>
## Encole entregas lentas con procedencia de tenant

Encole la entrega por red cuando la respuesta no necesita esperar, conserve tenant y acción en el Grant del job, renderice texto y HTML de forma coherente, use adjuntos acotados y propios de la aplicación y mantenga `log://` solo como evidencia de desarrollo. Para poner un log de actividad o un filtro delante de cada mensaje, llame a `drivers.WrapMail` antes de que `drivers.Mail` se entregue a nada: un servicio que tomó el mailer antes se queda con el que tomó.

- Pase contexto a la entrega.
- Use un único origen público.
- Encole con procedencia de Grant.
- Adjunte solo contenido autorizado y acotado.

<a name="keep-render-attachment-and-transport-failures-distinct"></a>
## Separe fallos de renderizado, adjuntos y transporte

Renderizado, validación de dirección, apertura o lectura de adjuntos, TLS, autenticación, rechazo del proveedor y timeout necesitan evidencia operativa propia. No marque mensajes de verificación o recuperación como enviados si solo su construcción tuvo éxito.

> [!WARNING]
> La integración publicada no vuelve automáticamente seguros todos los recorridos de adjuntos o Grant. Abra adjuntos dentro de la entrega autorizada, limite bytes y haga que los handlers de cola reconstruyan la policy actual sin confiar en rutas del payload.

<a name="mail-is-an-external-disclosure-boundary"></a>
## El correo es una frontera de divulgación externa

> [!WARNING]
> Nunca registre credenciales SMTP, códigos de un solo uso en producción, listas completas de destinatarios, adjuntos ni respuestas del proveedor con datos personales. Impida inyección de headers, use TLS y derive destinatario y tenant mediante servicios de aplicación.

Un enlace o código entregado sigue siendo entrada de aplicación vinculada a un propósito. La entrega de correo no autentica el callback ni autoriza cambios de cuenta.

<a name="capture-the-exact-mime-contract"></a>
## Capture el contrato MIME exacto

Use transporte de captura para From, To, asunto, texto y HTML, headers, adjuntos, cancelación, rechazo del proveedor y eliminación de datos sensibles. Ejecute la integración de aplicación y el recorrido encolado, no solo pruebas unitarias del constructor de mensajes.

- El MIME capturado tiene los campos públicos esperados.
- Un fallo de transporte permanece como evidencia para un reintento.
- La cancelación detiene la entrega.
- Los secretos no aparecen en los logs.

<a name="use-mail-inside-purpose-bound-application-flows"></a>
## Use correo dentro de flujos de aplicación vinculados a un propósito

Verificación de correo y recuperación de contraseña definen propósitos, las colas proporcionan entrega duradera, las notificaciones componen canales y la configuración elige el transporte.

- [Verificación de correo](/docs/{{version}}/security/email-verification)
- [Recuperación de contraseña](/docs/{{version}}/security/password-reset)
- [Colas](/docs/{{version}}/async/queues)
- [Configuración](/docs/{{version}}/architecture/configuration)
