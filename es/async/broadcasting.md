# Broadcasting

- [El nombre de un canal es una frontera tipada de publicación](#a-channel-name-is-a-typed-publication-boundary)
- [Elija transporte y autorización antes de publicar](#choose-transport-and-authorization-before-publishing)
- [Publique un evento dentro del ámbito de un tenant](#publish-one-tenant-scoped-event)
- [Mantenga explícitos la construcción de canales y los payloads](#keep-channel-construction-and-payloads-explicit)
- [No confundas un canal tipado con distribución completada](#do-not-confuse-a-typed-channel-with-delivered-fan-out)
- [Autorice la suscripción y la publicación por separado](#authorize-subscription-and-publication-separately)
- [Pruebe el aislamiento de canales y el rechazo del transporte](#test-channel-isolation-and-transport-refusal)
- [Use outbox y colas cuando la entrega deba sobrevivir](#use-outbox-and-queues-when-delivery-must-survive)

> [!WARNING]
> **Limitaciones actuales.** La autorización de canales en tiempo real está publicada. La distribución general de notificaciones y la entrega duradera son flujos separados aún incompletos; esta guía solo cubre la publicación autorizada de eventos en el broadcaster configurado.

<a name="a-channel-name-is-a-typed-publication-boundary"></a>
## El nombre de un canal es una frontera tipada de publicación

El paquete de broadcasting publicado ofrece canales tipados, eventos y contratos de driver. No demuestra la composición completa de un gestor en la aplicación ni el recorrido de suscripción en el navegador. Esta guía documenta la frontera útil del componente y señala los límites de la integración de extremo a extremo.

<a name="choose-transport-and-authorization-before-publishing"></a>
## Elija transporte y autorización antes de publicar

Selecciona un driver concreto en el bootstrap, define cómo puede suscribirse un Subject confiable y decide si es aceptable perder una conexión. Usa una outbox o cola duradera cuando el evento deba sobrevivir a la pérdida del proceso, sin tratar el broadcast como almacenamiento.

- Nombre el recurso al construir el canal.
- Deriva el namespace del tenant solo del Grant.
- Limite el tamaño del payload y sus campos públicos.
- Elija por separado el trabajo que deba ser duradero.

<a name="publish-one-tenant-scoped-event"></a>
## Publique un evento dentro del ámbito de un tenant

Construya un canal con la API publicada a partir del identificador del registro, defina explícitamente el contenido público y llame a `Broadcaster.Broadcast` con el Grant autorizado. El driver aplica el espacio de nombres del tenant mediante `TenantChannel`; la entrada de la aplicación nunca añade ese prefijo.

```go example
func PublishInvoiceUpdate(ctx context.Context, broadcaster broadcasting.Broadcaster, grant security.Grant, invoiceID, reference string) error {
	channel := broadcasting.NewPrivateChannel("invoices." + invoiceID)
	payload := map[string]any{
		"invoice": map[string]any{
			"id":        invoiceID,
			"reference": reference,
		},
	}
	if err := broadcaster.Broadcast(
		ctx, grant, []broadcasting.Channel{channel}, "invoice.updated", payload,
	); err != nil {
		return err
	}
	return nil
}
```

<a name="keep-channel-construction-and-payloads-explicit"></a>
## Mantenga explícitos la construcción de canales y los payloads

Mantenga estables los nombres de eventos, serialice un recurso deliberado en lugar de un model, no adjunte credenciales bearer ni de sesión y reconecte clientes por el recorrido normal de suscripción autenticada, sin confiar en una cadena de canal proporcionada por el navegador.

- Use tipos explícitos para canales públicos y privados.
- Mantén los datos del evento pequeños y versionables.
- Publique hechos inmutables.
- Realice replay duradero solo mediante un sistema de eventos persistidos.

<a name="do-not-confuse-a-typed-channel-with-delivered-fan-out"></a>
## No confundas un canal tipado con distribución completada

La indisponibilidad del driver, un cliente desconectado, un fallo de serialización y una suscripción denegada son resultados distintos. La API tipada publicada, por sí sola, no garantiza reconexión, replay, orden ni un gestor integrado en el esqueleto.

> [!WARNING]
> El contrato publicado actual no incluye un recorrido completo de broadcasting de la aplicación. Antes de depender de él, demuestre juntos el arranque, la policy de suscripción, el transporte elegido y la entrega en navegador o servicio en la aplicación con la etiqueta que despliega.

<a name="authorize-subscription-and-publication-separately"></a>
## Autorice la suscripción y la publicación por separado

> [!WARNING]
> Nunca aceptes un prefijo de tenant en un canal enviado por el cliente, incluyas secretos o modelos completos en el payload ni consideres conocer el nombre de un canal privado como autorización. `RequestedChannel` rechaza los dos puntos y `TenantChannel` deriva el namespace solo del Grant autorizado. Autentica y autoriza cada suscripción.

La entrega por broadcast es una representación, no una autorización. Un cliente que recibe un identificador aún debe pasar por la policy normal para cualquier lectura o modificación posterior.

<a name="test-channel-isolation-and-transport-refusal"></a>
## Pruebe el aislamiento de canales y el rechazo del transporte

Usa un driver fake para verificar canal y payload exactos, rechaza la suscripción de otro tenant, inyecta fallos de serialización y de driver, prueba la cancelación y ejecuta el flujo del transporte elegido antes de habilitarlo en el despliegue.

- Los nombres de canales contienen el ámbito confiable de tenant.
- Los suscriptores de otro tenant se rechazan.
- Los payloads exponen solo campos públicos.
- Un error del driver permanece observable.

<a name="use-outbox-and-queues-when-delivery-must-survive"></a>
## Use outbox y colas cuando la entrega deba sobrevivir

Eventos y outbox proporcionan hechos duraderos; las colas, trabajo duradero; la autorización protege la suscripción; los recursos definen payloads públicos.

- [Eventos y outbox](/docs/{{version}}/async/events-outbox)
- [Colas](/docs/{{version}}/async/queues)
- [Autorización](/docs/{{version}}/security/authorization)
- [Recursos HTTP](/docs/{{version}}/http/resources)
