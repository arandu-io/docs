# Deployment

- [Deploy one immutable binary in more than one process](#deploy-one-immutable-binary-in-more-than-one-process)
- [Prepare runtime configuration and shared dependencies](#prepare-runtime-configuration-and-shared-dependencies)
- [Build a versioned image and run its migration job](#build-a-versioned-image-and-run-its-migration-job)
- [Roll web, worker, and scheduled work from the same image](#roll-web-worker-and-scheduled-work-from-the-same-image)
- [Stop the rollout before an unhealthy replica receives traffic](#stop-the-rollout-before-an-unhealthy-replica-receives-traffic)
- [Ship no build context, credentials, or root runtime](#ship-no-build-context-credentials-or-root-runtime)
- [Exercise the container contract before promotion](#exercise-the-container-contract-before-promotion)
- [Roll back the artifact without pretending schema rollback is free](#roll-back-the-artifact-without-pretending-schema-rollback-is-free)
- [Continue with upgrades, migrations, and queue operations](#continue-with-upgrades-migrations-and-queue-operations)

<a name="deploy-one-immutable-binary-in-more-than-one-process"></a>
## Deploy one immutable binary in more than one process

The released starter builds one application binary. The same artifact serves HTTP, drains queues, and runs console work through different arguments, so web and worker cannot drift. Its multi-stage Dockerfile compiles Kyse, embeds assets, disables CGO, trims paths, and runs only `/app` as a non-root distroless process.

<a name="prepare-runtime-configuration-and-shared-dependencies"></a>
## Prepare runtime configuration and shared dependencies

Supply production `APP_KEY`, `APP_URL`, `HTTP_ADDR`, `DATABASE_URL`, and only settings a consumer reads. Production refuses `MAIL_URL=log://` and an unset `MAIL_URL`, so name a real mailer, and refuses `SESSION_SECURE_COOKIE=false`. A missing `APP_ENV` is production; the released image writes ENV `APP_ENV=prod` so the answer is stated rather than inherited. Production migrations need a shared cache lock, so select RESP and supply `REDIS_URL`. Keep `.env`, Git data, storage, and local SQLite outside the image.

- Database and selected RESP services are reachable.
- TLS trust files are present when configured.
- Secrets come from the platform.
- Liveness and readiness probes use separate routes.

<a name="build-a-versioned-image-and-run-its-migration-job"></a>
## Build a versioned image and run its migration job

Build an image tied to a version and full commit, then run migration as a one-off pipeline job before rollout, before any new binary takes jobs. The default serve command listens at `HTTP_ADDR` and writes one application environment line first, with env, debug, and `http_addr`; migrations acquire their shared lock and refuse production in-memory isolation.

```sh
COMMIT_SHA=$(git rev-parse HEAD)
docker build --build-arg VERSION=:aranduVersion --build-arg COMMIT="$COMMIT_SHA" -t registry.example/arandu::aranduVersion .
docker run --rm --env-file ./production.env registry.example/arandu::aranduVersion migrate
```

<a name="roll-web-worker-and-scheduled-work-from-the-same-image"></a>
## Roll web, worker, and scheduled work from the same image

Deploy web, queue workers, and scheduled commands from the same digest. Migrate once before changing replicas, allow SIGTERM to drain work, use `/_arandu/live` for restart decisions, and use `/_arandu/health` to decide traffic readiness. The health route runs at most ten passes over its checks a second, one at a time, and answers a probe past the tenth with the last result, so a load balancer's probes never stack up.

- Run serve for web replicas.
- Run the image with work `--queue=default --workers=4`; work is the internal target of `aru queue:work`.
- Stamp VERSION and COMMIT for incidents.
- Route traffic only from health, never liveness.

<a name="stop-the-rollout-before-an-unhealthy-replica-receives-traffic"></a>
## Stop the rollout before an unhealthy replica receives traffic

View compilation, typed configuration, connector, migration lock, and migration failures stop promotion. A new replica with failed module readiness stays out of service while its structured log names the failing boundary.

> [!WARNING]
> The released image does not migrate in ENTRYPOINT or CMD. Adding migration there races every starting replica. Distroless has no shell, so recovery uses logs, probes, a replacement image, and reproducible commands.

<a name="ship-no-build-context-credentials-or-root-runtime"></a>
## Ship no build context, credentials, or root runtime

> [!WARNING]
> Never copy `.env` or a local database into the image. Preserve nonroot execution, CGO-disabled static builds, the distroless runtime, pinned Aru tooling, and a registry digest tied to the commit. Leave production diagnostics disabled unless an explicit secret is controlled.

<a name="exercise-the-container-contract-before-promotion"></a>
## Exercise the container contract before promotion

From a clean checkout, generate views, test with race, vet, build, create the image, start with production-shaped config, verify both probes, run migrate twice, test previous binary against additive schema, drain a worker, and inspect user and secret absence.

- The image runs as nonroot.
- Health gates traffic and live gates restart.
- Repeated migration is safe.
- SIGTERM drains in-flight work.

<a name="roll-back-the-artifact-without-pretending-schema-rollback-is-free"></a>
## Roll back the artifact without pretending schema rollback is free

Stop traffic to the bad artifact and redeploy the previous immutable image. Additive migrations should leave it working. Run `migrate:rollback` only when the exact batch has a tested Down; destructive column removal takes two releases, not an emergency command. Rollback is a step of your deployment platform: the application ships no automatic rollback, so rehearse redeploying the previous image before you need it.

<a name="continue-with-upgrades-migrations-and-queue-operations"></a>
## Continue with upgrades, migrations, and queue operations

Upgrades choose versions, migrations coordinate schema, queue operations coordinate workers, and observability proves the running artifact.

- [Upgrades](/docs/{{version}}/operations/upgrades)
- [Migrations](/docs/{{version}}/database/migrations)
- [Queue operations](/docs/{{version}}/operations/queues)
- [Observability](/docs/{{version}}/operations/observability)
