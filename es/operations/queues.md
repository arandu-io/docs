# Operaciones de cola

- [Opere la cola mediante el binario de la aplicación](#operate-the-queue-through-the-application-binary)
- [Ejecute migraciones de colas y registre todos los handlers](#migrate-queue-tables-and-register-every-handler)
- [Inicie un worker que complete trabajo pendiente e inspeccione el backlog](#start-one-draining-worker-and-inspect-backlog)
- [Pause la entrada, reinicie workers y gestione jobs fallidos](#pause-intake-restart-workers-and-manage-dead-letters)
- [Trate reservas y jobs fallidos como estado recuperable](#treat-leases-and-failed-jobs-as-recoverable-state)
- [Mantenga el ámbito del tenant en payload, handler y comando](#keep-tenant-scope-in-payload-handler-and-command)
- [Ejecute recorridos de éxito, duplicado, retry y cierre](#run-one-success-duplicate-retry-and-shutdown-journey)
- [Recupere el estado antes de borrarlo](#recover-state-before-deleting-it)
- [Separe el diseño de la cola de su operación](#separate-queue-design-from-queue-operations)

> [!WARNING]
> **Limitaciones actuales.** La CLI y el esqueleto publicados exponen comandos de worker, pausa, retry, jobs fallidos y limpieza, pero la suite etiquetada demuestra principalmente inventario y rechazos; no ejecuta `queue:listen`. Ensaye recuperación y comportamiento del listener con el almacenamiento elegido antes de producción.

<a name="operate-the-queue-through-the-application-binary"></a>
## Opere la cola mediante el binario de la aplicación

Aru reenvía comandos de cola al binario del proyecto porque solo él conoce handlers, base, proveedor de jobs fallidos, repositorio de lotes, reglas de tenant y conexión elegida. El esqueleto conecta un worker y comandos operativos; no proporciona dashboard ni supervisor.

<a name="migrate-queue-tables-and-register-every-handler"></a>
## Ejecute migraciones de colas y registre todos los handlers

Migre jobs, `failed_jobs` y `job_batches`; registre los handlers en el `RegisterJobs` del provider, que `bootstrap/app.go` entrega a `WithJobs`; y defina la cola predeterminada y un lease de reintento mayor que el handler más lento. `QUEUE_CONNECTION=database` guarda los jobs en la base de datos de la aplicación, y redis selecciona la cola RESP una vez enlazado `hesape/queue/connectors/redis` por su import en blanco; nombrar redis antes detiene el arranque con las dos líneas que hay que añadir. Las señales de pausa y reinicio entre procesos siguen exigiendo una caché RESP compartida.

- `QUEUE_WORKERS` y `QUEUE_MAX_ATTEMPTS` son positivos cuando se definen.
- `REDIS_URL` está presente al elegir la caché compartida de señales.
- Payload y handler preservan el tenant y la procedencia del Grant.
- Worker y web usan la misma imagen.

<a name="start-one-draining-worker-and-inspect-backlog"></a>
## Inicie un worker que complete trabajo pendiente e inspeccione el backlog

Inicie un worker con cola y concurrencia explícitas y monitorice esa conexión y cola. Los valores de entorno completan opciones omitidas; las flags sustituyen solo sus campos documentados de invocación. SIGTERM detiene la entrada y completa el trabajo en curso.

```sh
aru queue:work --queue=default --workers=4
aru queue:monitor database:default --max=1000
```

<a name="pause-intake-restart-workers-and-manage-dead-letters"></a>
## Pause la entrada, reinicie workers y gestione jobs fallidos

Pause y reanude un destino, solicite reinicio ordenado de workers, inspeccione fallos por tenant, reintente o descarte IDs seleccionados y limpie solo registros fallidos o de lote revisados. `queue:listen` es un flujo de recarga en desarrollo, no el daemon ordinario de producción.

- `aru queue:pause database:default --for=30m` y `queue:resume`.
- Ejecute `aru queue:restart` tras un despliegue.
- Ejecute `aru queue:failed --tenant=<id>` antes de retry o forget.
- `queue:retry`, `queue:forget` y `queue:retry-batch` también leen sus flags después de los ids.
- Use retry-batch y comandos de limpieza dentro del ámbito del tenant.

<a name="treat-leases-and-failed-jobs-as-recoverable-state"></a>
## Trate reservas y jobs fallidos como estado recuperable

Los errores del handler generan reintentos hasta `MaxAttempts` y luego apartan el trabajo en `failed_jobs`. Un worker muerto solo libera el job tras la reserva; una reserva menor que la ejecución permite duplicidad. Corrija handlers, tablas, stores elegidos o caché compartida de señales ausentes antes de confiar en comandos de recuperación.

> [!WARNING]
> Las pruebas etiquetadas demuestran inventario y rechazos de comandos, pero no ejecutan `queue:listen` ni un recorrido completo de pausa y retry entre procesos. Esta página operativa es limitada hasta que ese recorrido similar a producción se publique y consuma.

<a name="keep-tenant-scope-in-payload-handler-and-command"></a>
## Mantenga el ámbito del tenant en payload, handler y comando

> [!WARNING]
> Los comandos de jobs fallidos y lotes exigen ámbito de tenant. Nunca deserialice un Grant o rol enviado por el llamador; reautorice en los handlers, minimice payloads, oculte datos sensibles de errores y haga idempotentes los efectos externos, pues la entrega ocurre al menos una vez.

<a name="run-one-success-duplicate-retry-and-shutdown-journey"></a>
## Ejecute recorridos de éxito, duplicado, retry y cierre

Ejercite éxito, duplicado, reintento, expiración del lease, recuperación tras caída, pausa y reanudación entre procesos, drenaje por SIGTERM, aislamiento de tenant, la cola de base de datos publicada y la propagación de `QUEUE_DEFAULT`, `QUEUE_WORKERS`, `QUEUE_RETRY_AFTER` y `QUEUE_MAX_ATTEMPTS` hasta el worker. Verifique que `QUEUE_CONNECTION=redis` sin el conector enlazado detiene el arranque en lugar de recurrir a la base de datos.

- Se completa un job real.
- Otro tenant no puede inspeccionar ni reintentar el job.
- El fallo llega al almacenamiento de jobs fallidos y permite reintento.
- El cierre completa el trabajo sin aceptar nuevos jobs.

<a name="recover-state-before-deleting-it"></a>
## Recupere el estado antes de borrarlo

Pause la cola afectada, preserve registros y logs, corrija o revierta el handler, reintente un job de prueba y reanude. Clear, forget, flush y prune borran trabajo o evidencia; haga backup y use el filtro más reducido de tenant e ID.

<a name="separate-queue-design-from-queue-operations"></a>
## Separe el diseño de la cola de su operación

La guía asíncrona define payloads y handlers; los logs muestran evidencia de workers; el despliegue distribuye el mismo artefacto; las pruebas de consola ejercitan la frontera de comandos delegados.

- [Diseño de la cola](/docs/{{version}}/async/queues)
- [Logs](/docs/{{version}}/operations/logging-observability)
- [Despliegue](/docs/{{version}}/operations/deployment)
- [Pruebas de consola](/docs/{{version}}/testing/console)
