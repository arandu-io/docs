# Despliegue

- [Despliegue un binario inmutable en más de un proceso](#deploy-one-immutable-binary-in-more-than-one-process)
- [Prepare configuración de runtime y dependencias compartidas](#prepare-runtime-configuration-and-shared-dependencies)
- [Construya una imagen versionada y ejecute su job de migración](#build-a-versioned-image-and-run-its-migration-job)
- [Actualice web, worker y tareas programadas con la misma imagen](#roll-web-worker-and-scheduled-work-from-the-same-image)
- [Detenga el rollout antes de que una réplica no saludable reciba tráfico](#stop-the-rollout-before-an-unhealthy-replica-receives-traffic)
- [No incluya contexto de build, credenciales ni ejecución como root](#ship-no-build-context-credentials-or-root-runtime)
- [Ejercite el contrato del contenedor antes de la promoción](#exercise-the-container-contract-before-promotion)
- [Revierta el artefacto sin tratar el rollback del esquema como gratuito](#roll-back-the-artifact-without-pretending-schema-rollback-is-free)
- [Continúe con actualizaciones, migraciones y operaciones de cola](#continue-with-upgrades-migrations-and-queue-operations)

> [!WARNING]
> **Limitaciones actuales.** El binario, el contenedor y el paso explícito de migración están publicados. Un contrato completo de salud y rollback automático aún no es universal: defina probes según las dependencias realmente conectadas y ensaye el rollback fuera del arranque de peticiones.

<a name="deploy-one-immutable-binary-in-more-than-one-process"></a>
## Despliegue un binario inmutable en más de un proceso

El esqueleto publicado construye un binario de aplicación. El mismo artefacto sirve HTTP, procesa colas y ejecuta consola con argumentos distintos, evitando divergencias entre web y worker. Su Dockerfile multi-stage compila Kyse, incorpora assets, desactiva CGO, elimina rutas y ejecuta solo `/app` como proceso distroless no root.

<a name="prepare-runtime-configuration-and-shared-dependencies"></a>
## Prepare configuración de runtime y dependencias compartidas

Proporcione `APP_KEY`, `APP_URL`, `HTTP_ADDR` y `DATABASE_URL` de producción, y solo la configuración que algún consumidor lee. Producción rechaza `MAIL_URL=log://` y un `MAIL_URL` sin definir, así que nombre un mailer real, y rechaza `SESSION_SECURE_COOKIE=false`. Un `APP_ENV` ausente es producción; la imagen publicada escribe ENV `APP_ENV=prod` para que la respuesta quede declarada y no heredada. Las migrations de producción necesitan un lock de caché compartido, así que seleccione RESP y proporcione `REDIS_URL`. Mantenga `.env`, los datos de Git, storage y el SQLite local fuera de la imagen.

- La base y los servicios RESP elegidos están accesibles.
- Los archivos de confianza TLS están presentes cuando se configuran.
- Los secretos provienen de la plataforma.
- Las probes de liveness y readiness usan rutas separadas.

<a name="build-a-versioned-image-and-run-its-migration-job"></a>
## Construya una imagen versionada y ejecute su job de migración

Construya una imagen atada a una versión y a un commit completo y después ejecute la migration como un job puntual del pipeline antes del rollout, antes de que cualquier binario nuevo tome jobs. El comando serve predeterminado escucha en `HTTP_ADDR` y escribe primero una línea application environment, con env, debug y `http_addr`; las migrations adquieren su lock compartido y rechazan el aislamiento en memoria en producción.

```sh
COMMIT_SHA=$(git rev-parse HEAD)
docker build --build-arg VERSION=:aranduVersion --build-arg COMMIT="$COMMIT_SHA" -t registry.example/arandu::aranduVersion .
docker run --rm --env-file ./production.env registry.example/arandu::aranduVersion migrate
```

<a name="roll-web-worker-and-scheduled-work-from-the-same-image"></a>
## Actualice web, worker y tareas programadas con la misma imagen

Despliegue web, workers de cola y comandos programados desde el mismo digest. Migre una vez antes de cambiar las réplicas, permita que SIGTERM drene el trabajo, use `/_arandu/live` para las decisiones de reinicio y `/_arandu/health` para decidir si recibe tráfico. La ruta de salud ejecuta como máximo diez pasadas por segundo sobre sus comprobaciones, una a la vez, y responde a una sonda posterior a la décima con el último resultado, así que las sondas de un balanceador nunca se acumulan.

- Ejecute serve en las réplicas web.
- Ejecute la imagen con work `--queue=default --workers=4`; work es el destino interno de `aru queue:work`.
- Registre VERSION y COMMIT para investigar incidentes.
- Dirija tráfico según la salud, nunca según la mera liveness.

<a name="stop-the-rollout-before-an-unhealthy-replica-receives-traffic"></a>
## Detenga el rollout antes de que una réplica no saludable reciba tráfico

Los fallos de compilación de vistas, configuración tipada, conector, bloqueo o ejecución de migración impiden la promoción. Una réplica nueva con fallo de disponibilidad de módulo queda fuera de servicio mientras el log estructurado identifica la frontera fallida.

> [!WARNING]
> La imagen publicada no ejecuta migraciones en ENTRYPOINT ni CMD. Añadirlas allí genera carreras entre réplicas que arrancan. Distroless no tiene shell: la recuperación usa logs, probes, imagen sustituta y comandos reproducibles.

<a name="ship-no-build-context-credentials-or-root-runtime"></a>
## No incluya contexto de build, credenciales ni ejecución como root

> [!WARNING]
> Nunca copie `.env` ni una base local a la imagen. Preserve ejecución nonroot, builds estáticos con CGO desactivado, runtime distroless, Aru fijada y digest del registro vinculado al commit. Deje los diagnósticos de producción desactivados salvo control de un secreto explícito.

<a name="exercise-the-container-contract-before-promotion"></a>
## Ejercite el contrato del contenedor antes de la promoción

En un checkout limpio, genere vistas, pruebe con race, ejecute vet y build, cree la imagen, arranque con configuración similar a producción, verifique ambas probes, ejecute migrate dos veces, pruebe el binario anterior con esquema aditivo, complete trabajo de un worker y compruebe usuario de ejecución y ausencia de secretos.

- La imagen se ejecuta como nonroot.
- Health controla el tráfico y live controla el reinicio.
- Repetir una migración es seguro.
- SIGTERM completa el trabajo en curso.

<a name="roll-back-the-artifact-without-pretending-schema-rollback-is-free"></a>
## Revierta el artefacto sin tratar el rollback del esquema como gratuito

Retire el tráfico del artefacto defectuoso y redespliegue la imagen inmutable anterior. Las migraciones aditivas deben mantenerla funcional. Ejecute `migrate:rollback` solo si el lote exacto tiene Down probado; eliminar columnas de forma destructiva exige dos releases, no un comando de emergencia.

<a name="continue-with-upgrades-migrations-and-queue-operations"></a>
## Continúe con actualizaciones, migraciones y operaciones de cola

Las actualizaciones eligen versiones, las migraciones coordinan el esquema, las operaciones de cola coordinan workers y la observabilidad demuestra el artefacto en ejecución.

- [Actualizaciones](/docs/{{version}}/operations/upgrades)
- [Migraciones](/docs/{{version}}/database/migrations)
- [Operaciones de cola](/docs/{{version}}/operations/queues)
- [Observabilidad](/docs/{{version}}/operations/observability)
