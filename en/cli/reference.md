# Aru console and code generation

- [Aru is the adapter around a project-owned binary](#aru-is-the-adapter-around-a-project-owned-binary)
- [Released command families in :aruVersion](#released-command-families-in-aruversion)
- [Use the task guide behind each command](#use-the-task-guide-behind-each-command)

> [!WARNING]
> **Current limitations.** The command inventory, code generators, view build, and doctor are released. Migration, factory persistence, and queue recovery retain the limitations stated in their dedicated guides; the CLI does not make those paths safer than their underlying contracts.

<a name="aru-is-the-adapter-around-a-project-owned-binary"></a>
## Aru is the adapter around a project-owned binary

Aru creates and inspects a project, compiles Kyse, generates ordinary Go, and delegates application-aware commands to that project's binary. Run `aru help` for inventory, `aru about` for the project contract, and `aru version` for the installed CLI. This page describes :aruVersion; commands present only on a development branch are absent.

<a name="released-command-families-in-aruversion"></a>
## Released command families in :aruVersion

Identify with about, version, `action:list`, and `imports:catalog`. Create, run, and check with new, `key:generate`, dev, serve, build, test, `route:list`, trace, doctor, lsp, mcp, `view:build`, `vendor:publish`, `skills:sync`, and `openapi:check`; the native application has `native:run`, `native:dev`, `native:build`, and `native:icons`. Real-browser journeys are written with `make:e2e`, recorded with `e2e:record`, explained after a failure with `e2e:explain`, and prepared with `e2e:install-browsers`. Generate with schema, generate, `model:build`, `make:module`, `make:model`, `make:service`, `make:resource`, `make:notification`, `make:client`, `make:mcp-tool`, `make:mcp-resource`, `make:mcp-prompt`, `make:migration`, `make:controller`, `make:middleware`, `make:request`, `make:factory`, `make:seeder`, `make:job`, `make:mail`, `make:command`, `make:listener`, `make:event`, `make:enum`, `make:policy`, `make:test`, `make:contract`, and `make:e2e`. `model:build` rewrites the query file beside each entity, which is never edited by hand; every other generated file is application-owned and must be diffed and compiled, and a generator prints the wiring for app/Providers, routes, and the builder instead of writing it.

Database commands are migrate, `migrate:rollback`, `migrate:status`, `migrate:fresh`, which runs only with `APP_ENV=dev`, `db:seed`, `schedule:list`, and `schedule:run`. Queue commands are `queue:work`, listen, restart, pause, resume, clear, monitor, failed, retry, forget, flush, prune-failed, retry-batch, and prune-batches. Font commands are `font:search`, info, add, list, and remove; vendored files, licenses, and CSS become project assets.

| Surface | Released contract |
| --- | --- |
| `aru new <name> --module <path>` | clone and rewrite a starter, create private environment and editor settings |
| `aru dev` / serve / build / test | compile views and run, package, or test the project |
| `aru route:list` / trace / doctor / lsp | inspect runtime and editor contracts |
| `aru mcp` / `skills:sync` / `vendor:publish` / `openapi:check` | answer an assistant, keep Skills and published files in step, and check API contracts |
| `aru native:run` / `native:dev` / `native:build` | build and run the native application |
| `aru schema` / generate / `make:*` | schema-driven and granular application-owned source generation |
| `aru migrate*` / `db:seed` / `schedule:*` | delegated database and schedule operations |
| `aru queue:*` | delegated worker, signal, backlog, failed-job, and batch operations |
| `aru font:*` | vendored font discovery and project asset management |

```sh
brew install arandu-io/tap/aru
aru version
aru help
```

<a name="use-the-task-guide-behind-each-command"></a>
## Use the task guide behind each command

`aru new` :aruVersion shallow-clones the immutable starter :aranduVersion, requires Git and network, writes fresh secrets, attempts view build, and neither runs `go mod tidy` nor starts the server. Verify the resulting module graph; delegated commands fail if the project binary does not implement them.

`aru doctor` also warns when a project drifts from the builder layout: bootstrap-file-beside-app, program-under-cmd, command-outside-commands, and program-under-scripts. The catalogue says what exists. Installation, migration, queue, deployment, and editor guides own operational choices, safety, and recovery.

- [Installation](/docs/{{version}}/start/installation)
- [Migrations](/docs/{{version}}/database/migrations)
- [Queue operations](/docs/{{version}}/operations/queues)
- [Deployment](/docs/{{version}}/operations/deployment)
- [Editor extension](/docs/{{version}}/ai/assisted-development)
