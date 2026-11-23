# Contratos y adapters

- [El consumidor es responsable de la interfaz útil más pequeña](#the-consumer-owns-the-smallest-useful-interface)
- [Incluye autorización en el contrato del repositorio](#put-authorization-in-the-repository-seam)
- [Selecciona adapters al iniciar](#select-adapters-at-boot)
- [Evita interfaces que oculten autoridad o responsabilidad](#avoid-interfaces-that-hide-authority-or-ownership)
- [Ejecuta un contrato contra cada adapter](#run-one-contract-against-each-adapter)
- [Los contratos de Framework se mantienen acotados](#framework-contracts-stay-narrow)
- [La composición hace visible la sustitución](#composition-makes-substitution-visible)

<a name="the-consumer-owns-the-smallest-useful-interface"></a>
## El consumidor es responsable de la interfaz útil más pequeña

Un contrato nombra el comportamiento que necesita un consumidor. Un adapter lo satisface estructuralmente y se selecciona en la raíz de composición. El código de dominio permanece independiente de SQL, RESP, SMTP o un servicio externo sin añadir un registro ni un paquete universal de abstracciones.

<a name="put-authorization-in-the-repository-seam"></a>
## Incluye autorización en el contrato del repositorio

Incluye autorización en el contrato del repositorio para que sustituir el almacenamiento no elimine el ámbito del tenant. El servicio declara solo la operación que llama y la aserción de compilación mantiene el adapter SQL concreto alineado con la interfaz de ese consumidor.

```go example
type Invoice struct{ ID string }

type InvoiceRepository interface {
	Find(context.Context, security.Grant, string) (Invoice, error)
}

type InvoiceService struct{ invoices InvoiceRepository }

func NewInvoiceService(invoices InvoiceRepository) *InvoiceService {
	return &InvoiceService{invoices: invoices}
}

type SQLInvoiceRepository struct{}

func (*SQLInvoiceRepository) Find(context.Context, security.Grant, string) (Invoice, error) {
	return Invoice{}, nil
}

var _ InvoiceRepository = (*SQLInvoiceRepository)(nil)
```

<a name="select-adapters-at-boot"></a>
## Selecciona adapters al iniciar

`foundation.NewDrivers` construye el conector de base de datos, el store de caché, el transporte de correo y el disco del sistema de archivos que nombra la configuración, y un provider pasa cada uno, o un fake de test, una vez. Los adaptadores traducen los errores de ausencia, conflicto e indisponibilidad del driver cuando quienes llaman necesitan bifurcar, conservando la causa para los operadores.

- Valida las dependencias obligatorias en los constructores.
- Mantén los comportamientos opcionales en interfaces pequeñas separadas.
- Encapsula los errores del driver sin exponer credenciales ni SQL.
- Coloca aserciones de compilación junto a los adapters.

<a name="avoid-interfaces-that-hide-authority-or-ownership"></a>
## Evita interfaces que oculten autoridad o responsabilidad

Un repositorio que solo recibe un id para datos de un tenant no tiene dónde demostrar autoridad. Las interfaces grandes obligan a los fakes a implementar métodos sin uso y acoplan consumidores no relacionados; las interfaces pequeñas duplicadas son más seguras cuando las semánticas realmente difieren.

> [!WARNING]
> No compenses la ausencia de autoridad con un tenant global, un singleton local de solicitud ni una lectura de entorno dentro del adapter. Pasa `security.Grant` en la operación que accede al almacenamiento protegido.

<a name="run-one-contract-against-each-adapter"></a>
## Ejecuta un contrato contra cada adapter

Ejecuta un contrato compartido contra cada adapter. Prueba éxito, ausencia, aislamiento de tenant, cancelación y mapeo de errores mediante la interfaz del consumidor, sin acceder a detalles de implementación de un driver concreto.

- Los Grants cero y de otros tenants no devuelven datos protegidos.
- Los adapters de memoria y de red cumplen el mismo comportamiento.
- La cancelación retorna con prontitud.
- Las aserciones de compilación cubren cada adapter concreto.

<a name="framework-contracts-stay-narrow"></a>
## Los contratos de Framework se mantienen acotados

Elija el contrato de comunicación según la interacción, no según un ranking de popularidad. REST es petición y respuesta para recursos públicos; `gRPC` es un límite de servicio tipado que también admite streaming; una cola persiste trabajo asíncrono; `WebSocket` es una sesión bidireccional en tiempo real; los eventos enviados por el servidor (SSE) transmiten actualizaciones del servidor en una respuesta HTTP; y un webhook entrega un evento firmado a un endpoint HTTP externo.

| Superficie | Contrato publicado |
| --- | --- |
| REST | recursos de petición y respuesta; use para APIs públicas o dirigidas al navegador |
| `gRPC` | llamadas tipadas unarias o por streaming; use Connect, `gRPC` o gRPC-Web en el servidor HTTP de la aplicación |
| Cola | trabajo asíncrono duradero con entrega al menos una vez |
| `WebSocket` | sesiones bidireccionales persistentes para interacción de baja latencia |
| SSE | actualizaciones unidireccionales del servidor mediante una respuesta HTTP cancelable |
| Webhook | entrega duradera y firmada de eventos a un servicio externo |

<a name="composition-makes-substitution-visible"></a>
## La composición hace visible la sustitución

REST, `gRPC`, colas, `WebSocket`, SSE y webhooks resuelven problemas diferentes de propiedad y entrega. La configuración elige el adaptador concreto, la raíz de composición lo construye y el servicio conserva solo su contrato. Context lleva la vida útil de la petición y Grant lleva el ámbito de tenant autorizado, haciendo visible la sustitución en un único diff de composición.

No sustituya el almacenamiento duradero por un transporte en tiempo real: `WebSocket` y SSE no reproducen por sí solos el trabajo perdido. Una cola ejecuta trabajo interno, mientras que un webhook cruza un límite HTTP externo. El streaming `gRPC` comparte el proceso servidor y el ciclo de apagado en lugar de crear un segundo listener.

- [Enrutamiento HTTP y `gRPC`](/docs/{{version}}/http/routing)
- [Respuestas HTTP y SSE](/docs/{{version}}/http/responses)
- [Colas](/docs/{{version}}/async/queues)
- [`WebSocket` de Joaju](/docs/{{version}}/async/joaju)
- [Eventos y webhooks](/docs/{{version}}/async/events-outbox)
