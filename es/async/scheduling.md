# Programación de tareas

- [Deje que el programador decida cuándo ejecutar, no cómo reintentar](#let-the-scheduler-decide-when-not-how-to-retry)
- [Declare el ámbito antes de escribir el cuerpo de la tarea](#declare-scope-before-writing-the-task-body)
- [Programa un job para cada tenant](#schedule-one-job-for-every-tenant)
- [Inspecciona y ejecuta la programación registrada](#inspect-and-invoke-the-registered-schedule)
- [Trate la ausencia de bloqueo como un defecto de despliegue](#treat-a-missing-lock-as-a-deployment-defect)
- [Las tareas Global no pueden acceder a datos de tenant](#global-tasks-cannot-touch-tenant-data)
- [Pruebe programación, cola y worker de extremo a extremo](#run-the-schedule-queue-and-worker-end-to-end)
- [Repara la dependencia y ejecuta la tarea una vez](#repair-the-dependency-and-run-the-task-once)
- [Mantenga el trabajo duradero en la cola](#keep-durable-work-in-the-queue)

<a name="let-the-scheduler-decide-when-not-how-to-retry"></a>
## Deje que el programador decida cuándo ejecutar, no cómo reintentar

Una tarea programada es una pequeña decisión ejecutada por una expresión cron de cinco campos. Normalmente descubre trabajo y encola un job duradero; el programador no repite una ejecución Run fallida.

<a name="declare-scope-before-writing-the-task-body"></a>
## Declare el ámbito antes de escribir el cuerpo de la tarea

Una `foundation.Task` tiene un ID único, una Spec válida, Action y Run. Global recibe un Grant cero y no llega a los repositorios; `PerTenant` recibe `SystemGrant(Action, tenant)` para cada tenant que responde la función entregada a `WithTenants`. Defina Timeout, porque cero se convierte en cinco minutos.

- Defina un identificador único para la tarea.
- Elija Global o `PerTenant` explícitamente.
- Responda los tenants con `WithTenants` para `PerTenant`.
- Las reclamaciones singleton solo se comparten entre réplicas cuando `drivers.Locks` lo está.

<a name="schedule-one-job-for-every-tenant"></a>
## Programa un job para cada tenant

Liste las tareas de la aplicación en `routes/console.go` y entréguelas a `WithSchedule`; un módulo declara las suyas mediante Schedule. Create registra el programador al final, sobre las tareas de todos los módulos y las de la aplicación, y se ejecuta en el proceso que sirve. Mantenga el trabajo duradero en un handler de cola: el Grant del programador se convierte en la procedencia del job.

```go example
type CloseCycle struct{ Cycle string }

func Tasks(d *foundation.Drivers) []foundation.Task {
	return []foundation.Task{{
		ID: "billing.close", Spec: "0 3 * * *", Scope: foundation.PerTenant,
		Action: "invoice.send", Timeout: 30 * time.Second, Singleton: true,
		Run: func(ctx context.Context, g auth.Grant) error {
			return queue.Dispatch(ctx, d.Queue, g, "invoice.send", CloseCycle{Cycle: "current"})
		},
	}}
}
```

<a name="inspect-and-invoke-the-registered-schedule"></a>
## Inspecciona y ejecuta la programación registrada

`aru schedule:list` informa id, cron, ámbito, singleton, próxima ejecución, última ejecución y último error. `schedule:run` utiliza el mismo recorrido de Grant, bloqueo e instrumentación; una ejecución manual por tenant exige `--tenant`, y la ejecución manual de un singleton toma la misma reclamación, así que una segunda en el mismo minuto UTC no hace nada.

- Inspecciona el registro antes de esperar al cron.
- Ejecute una tarea manualmente para un tenant seleccionado.
- Encole trabajo lento o que pueda requerir reintentos.
- Prepare el cuerpo de la tarea para ejecutarse dos veces con seguridad.

<a name="treat-a-missing-lock-as-a-deployment-defect"></a>
## Trate la ausencia de bloqueo como un defecto de despliegue

Los procesos que sirven comprueban cerca de cada minuto. Sin un store compartido, `drivers.Locks` es nil y cada réplica puede ejecutar un singleton. Un singleton reclama su tarea, el tenant y el minuto UTC (`sched:<tenant>:<id>:<minuto>`) y nunca libera la reclamación: expira tras una hora, o tras el timeout si es mayor, así que una ejecución que falla no la repite otra réplica en el mismo minuto. Una tarea por tenant sin `WithTenants` no se ejecuta, y el programador lo informa en lugar de saltarla en silencio. Una ejecución que falla no se reintenta, así que una tarea que ejecuta directamente un efecto externo puede perder ese trabajo; encole un job duradero e idempotente.

<a name="global-tasks-cannot-touch-tenant-data"></a>
## Las tareas Global no pueden acceder a datos de tenant

> [!WARNING]
> Nunca inventes un tenant dentro de Global ni lo obtengas de HTTP, payload o atajos de entorno. El acceso a repositorios recibe solo el Grant por tenant creado por el resolvedor confiable y la Action declarada.

<a name="run-the-schedule-queue-and-worker-end-to-end"></a>
## Pruebe programación, cola y worker de extremo a extremo

Construya el programador, llame a `RunNow`, confirme un job persistido, inicie un worker y confirme que el handler recibe tenant y payload. IDs duplicados, cron inválido y tareas sin Run fallan antes de servir. La ausencia del resolvedor de tenants aparece en ejecución y diagnósticos; la falta de bloqueo compartido debe detectarse mediante una prueba de configuración del despliegue, pues el singleton se ejecutaría en cada réplica.

- `RunNow` encola una vez.
- El worker recibe el Grant esperado.
- Las programaciones inválidas impiden el arranque.
- Solo quien obtiene el lock singleton ejecuta.

<a name="repair-the-dependency-and-run-the-task-once"></a>
## Repara la dependencia y ejecuta la tarea una vez

Usa `schedule:list` para distinguir la ausencia de registro de un fallo en la última ejecución. Restaura base de datos, cola, resolvedor o lock e invoca un tenant afectado. El scheduler no reintenta automáticamente, pero la recuperación puede ejecutar el trabajo de nuevo; haz seguros de repetir tanto la tarea como cualquier handler encolado.

<a name="keep-durable-work-in-the-queue"></a>
## Mantenga el trabajo duradero en la cola

Las colas gestionan la ejecución duradera, la outbox los hechos tras el commit, las operaciones de colas la recuperación y la observabilidad informa de fallos y retrasos de programación.

- [Colas](/docs/{{version}}/async/queues)
- [Eventos y outbox](/docs/{{version}}/async/events-outbox)
- [Operaciones de cola](/docs/{{version}}/operations/queues)
- [Observabilidad](/docs/{{version}}/operations/observability)
