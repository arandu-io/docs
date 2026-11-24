# Colas

- [Persista el trabajo que deba sobrevivir a la petición](#persist-work-that-must-survive-the-request)
- [Elija un almacenamiento cuya frontera transaccional conozca](#choose-a-store-whose-transaction-boundary-you-understand)
- [Crea desde un Grant, envía a la cola y procesa por nombre](#create-from-a-grant-push-and-handle-by-name)
- [Defina la reserva y los reintentos según el comportamiento del handler](#set-lease-and-retry-policy-from-handler-behavior)
- [Sepa si el job se liberó o se apartó](#know-whether-a-job-was-released-or-parked)
- [Nunca reconstruyas autoridad desde el payload](#never-rebuild-authority-from-the-payload)
- [Verifica rollback y entrega duplicada](#prove-rollback-and-duplicate-delivery)
- [Inspeccione, repare, reintente y observe el retraso](#inspect-repair-retry-and-observe-lag)
- [Separa el diseño del productor de las operaciones del worker](#separate-producer-design-from-worker-operations)

<a name="persist-work-that-must-survive-the-request"></a>
## Persista el trabajo que deba sobrevivir a la petición

Una cola separa la aceptación del trabajo de su ejecución. Los jobs de la cola de Hesape llevan tenant, acción, autorizador, payload y UUID estable; el worker reconstruye el Grant antes de invocar el handler registrado. La entrega ocurre al menos una vez, por lo que la idempotencia del handler es un contrato de la aplicación.

<a name="choose-a-store-whose-transaction-boundary-you-understand"></a>
## Elija un almacenamiento cuya frontera transaccional conozca

Create registra el propio módulo de la cola: la tabla de jobs, la salud de retraso y los diagnósticos de la cola que nombra `QUEUE_CONNECTION`, guardada en `drivers.Queue`. Registrar `jobs.NewModule` otra vez en `WithProviders` lo registra dos veces. La cola de base de datos puede grabar su fila en la misma transacción SQL que la escritura que la produce. RESP, seleccionado con `QUEUE_CONNECTION=redis` una vez enlazado `hesape/queue/connectors/redis`, es útil cuando se elige deliberadamente, pero no entra en esa transacción.

- Elija deliberadamente base de datos o RESP.
- Deje el módulo de la cola a Create.
- Mantén servidor web y worker en el mismo artefacto.
- Defina la idempotencia antes de encolar.

<a name="create-from-a-grant-push-and-handle-by-name"></a>
## Crea desde un Grant, envía a la cola y procesa por nombre

Registre el nombre exacto en el `RegisterJobs` del provider, que `bootstrap/app.go` entrega a `WithJobs`, y despache un payload tipado bajo el Grant autorizado con `queue.Dispatch`. `HandlerFor` decodifica el payload en el mismo tipo antes de que se ejecute el handler. Dispatch construye el job con `jobs.New` en la cola predeterminada, así que un tenant o nombre ausente se rechaza, y Push comprueba la autoridad incluso para un job armado a mano para otra cola.

```go example
type SendInvoice struct{ InvoiceID string }

type InvoiceSender interface {
	Send(ctx context.Context, g auth.Grant, jobID, invoiceID string) error
}

type AppServiceProvider struct{ Invoices InvoiceSender }

func (p *AppServiceProvider) RegisterJobs(w *queue.Worker) error {
	w.Handle("invoice.send", queue.HandlerFor(func(ctx context.Context, g auth.Grant, j *jobs.Job, in SendInvoice) error {
		return p.Invoices.Send(ctx, g, j.UUID, in.InvoiceID)
	}))
	return nil
}

func Enqueue(ctx context.Context, drivers *foundation.Drivers, grant auth.Grant, invoiceID string) error {
	if err := queue.Dispatch(ctx, drivers.Queue, grant, "invoice.send", SendInvoice{InvoiceID: invoiceID}); err != nil {
		return err
	}
	return nil
}
```

<a name="set-lease-and-retry-policy-from-handler-behavior"></a>
## Defina la reserva y los reintentos según el comportamiento del handler

`queue:work` construye su worker a partir de `Settings.Queue`: `QUEUE_DEFAULT` nombra la cola, `QUEUE_WORKERS` la concurrencia, `QUEUE_RETRY_AFTER` el lease en segundos y `QUEUE_MAX_ATTEMPTS` los intentos, y `--queue` y `--workers` sobrescriben los dos primeros en una invocación. Una variable sin definir mantiene el valor predeterminado del propio componente de cola: cuatro workers, un lease de cinco minutos, cinco intentos y una espera entre ellos que se duplica desde dos segundos y deja de crecer en una hora. El contexto del handler se cancela pasado un lease, así que un handler que retorna cuando su contexto termina devuelve el job antes de que otro worker pueda tomarlo.

- Limite la concurrencia según las dependencias.
- Defina la reserva por encima de la duración máxima del handler.
- Defina explícitamente la política de intentos y espera progresiva.
- Pase el Recorder de la aplicación a los workers.

<a name="know-whether-a-job-was-released-or-parked"></a>
## Sepa si el job se liberó o se apartó

Un error devuelto consume un intento y libera el job después del backoff; el trabajo agotado se aparca. `ErrManuallyFailed` aparca un fallo permanente de inmediato, y un panic se recupera, se registra con la pila y se aparca en la primera entrega. El worker informa de un job de un batch o de una chain después de ejecutarlo, así que un handler que todavía llama a `bus.Handled` lo cuenta dos veces y empuja dos veces el siguiente eslabón de la chain. Los payloads fallidos siguen siendo datos de cliente protegidos por tenant.

> [!WARNING]
> Un handler que ignora su contexto sigue ejecutándose cuando el lease termina, y otro worker puede tomar el mismo job. Aumentar los intentos no vuelve seguro un efecto externo sin idempotencia.

<a name="never-rebuild-authority-from-the-payload"></a>
## Nunca reconstruyas autoridad desde el payload

> [!WARNING]
> El worker usa el tenant y la acción registrados para su Grant, y un job empujado mientras un operador actuaba en la cuenta lleva a ese operador, de modo que lo que escribe queda registrado como suyo. Un push bajo un Grant de solo lectura se rechaza con `auth.ErrReadOnly`. Nunca acepte `tenant_id`, action, role ni un Grant serializado desde JSON. Pase el Grant del worker a servicios y repositorios, y mantenga la inspección de fallos con alcance de tenant.

<a name="prove-rollback-and-duplicate-delivery"></a>
## Verifica rollback y entrega duplicada

Para la cola de base de datos, encole en una transacción que devuelve un error centinela y demuestre que `PendingSize` permanece en cero. Entregue dos veces el mismo job estable y demuestre un único efecto externo; pruebe también payload malformado, fallo permanente, panic y reserva deliberadamente corta.

- El rollback elimina el registro encolado.
- La entrega duplicada es idempotente.
- Un fallo permanente aparta el trabajo.
- Un tenant no puede inspeccionar trabajo fallido de otro.

<a name="inspect-repair-retry-and-observe-lag"></a>
## Inspeccione, repare, reintente y observe el retraso

Inspeccione fallos y retraso, repare el handler o la dependencia, reintente un job de prueba y sustituya los workers con cierre ordenado. Forget, flush y clear destruyen evidencia o trabajo y exigen el ámbito de tenant revisado más reducido posible.

<a name="separate-producer-design-from-worker-operations"></a>
## Separa el diseño del productor de las operaciones del worker

Operaciones documenta comandos, la programación crea jobs recurrentes, la outbox persiste hechos y la autorización proporciona el límite de Grant. Una cola es un límite interno de ejecución duradera; un webhook usa trabajo en cola para cruzar a un servicio HTTP externo.

- [Patrones de comunicación](/docs/{{version}}/architecture/contracts-adapters)
- [Operaciones de cola](/docs/{{version}}/operations/queues)
- [Programación](/docs/{{version}}/async/scheduling)
- [Eventos y webhooks](/docs/{{version}}/async/events-outbox)
- [Autorización](/docs/{{version}}/security/authorization)
