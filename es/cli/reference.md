# Consola y generación de código de aru

- [Aru es el adapter del binario que pertenece al proyecto](#aru-is-the-adapter-around-a-project-owned-binary)
- [Familias de comandos publicadas en :aruVersion](#released-command-families-in-aruversion)
- [Consulta la guía de la tarea detrás de cada comando](#use-the-task-guide-behind-each-command)

> [!WARNING]
> **Limitaciones actuales.** El inventario de comandos, los generadores de código, la compilación de vistas y doctor están publicados. Las migraciones, la persistencia de factories y la recuperación de colas conservan las limitaciones descritas en sus propias guías; la CLI no hace estos flujos más seguros de lo que permiten sus contratos.

<a name="aru-is-the-adapter-around-a-project-owned-binary"></a>
## Aru es el adapter del binario que pertenece al proyecto

Aru crea e inspecciona proyectos, compila Kyse, genera Go común y delega los comandos que dependen de la aplicación en el binario del proyecto. Ejecuta `aru help` para consultar el inventario, `aru about` para el contrato del proyecto y `aru version` para la CLI instalada. Esta página describe :aruVersion; no incluye comandos que solo existen en una rama de desarrollo.

<a name="released-command-families-in-aruversion"></a>
## Familias de comandos publicadas en :aruVersion

Identifique con about, version, `action:list` y `imports:catalog`. Cree, ejecute y verifique con new, `key:generate`, dev, serve, build, test, `route:list`, trace, doctor, lsp, mcp, `view:build`, `vendor:publish`, `skills:sync` y `openapi:check`; la aplicación nativa tiene `native:run`, `native:dev`, `native:build` y `native:icons`. Los recorridos en navegador real se escriben con `make:e2e`, se graban con `e2e:record`, se explican tras un fallo con `e2e:explain` y se preparan con `e2e:install-browsers`. Genere con schema, generate, `model:build`, `make:module`, `make:model`, `make:service`, `make:resource`, `make:notification`, `make:client`, `make:mcp-tool`, `make:mcp-resource`, `make:mcp-prompt`, `make:migration`, `make:controller`, `make:middleware`, `make:request`, `make:factory`, `make:seeder`, `make:job`, `make:mail`, `make:command`, `make:listener`, `make:event`, `make:enum`, `make:policy`, `make:test`, `make:contract` y `make:e2e`. `model:build` reescribe el archivo de query junto a cada entidad, que nunca se edita a mano; todo otro archivo generado pertenece a la aplicación y debe compararse y compilarse, y un generador imprime el cableado para app/Providers, las rutas y el builder en lugar de escribirlo.

Los comandos de base de datos son migrate, `migrate:rollback`, `migrate:status`, `migrate:fresh`, que solo se ejecuta con `APP_ENV=dev`, `db:seed`, `schedule:list` y `schedule:run`. Los comandos de cola son `queue:work`, listen, restart, pause, resume, clear, monitor, failed, retry, forget, flush, prune-failed, retry-batch y prune-batches. Los comandos de fuentes son `font:search`, info, add, list y remove; los archivos, las licencias y el CSS vendorizados se convierten en assets del proyecto.

| Superficie | Contrato publicado |
| --- | --- |
| `aru new <name> --module <path>` | clonar y adaptar un proyecto inicial, crear un entorno privado y la configuración del editor |
| `aru dev` / serve / build / test | compilar vistas y ejecutar, empaquetar o probar el proyecto |
| `aru route:list` / trace / doctor / lsp | inspeccionar contratos de runtime y del editor |
| `aru mcp` / `skills:sync` / `vendor:publish` / `openapi:check` | responder a un asistente, mantener al día las Skills y los archivos publicados y verificar contratos de API |
| `aru native:run` / `native:dev` / `native:build` | construir y ejecutar la aplicación nativa |
| `aru schema` / generate / `make:*` | generación de código de la aplicación basada en schema o mediante comandos específicos |
| `aru migrate*` / `db:seed` / `schedule:*` | operaciones delegadas de base de datos y programación |
| `aru queue:*` | operaciones delegadas de workers, señales, backlog, jobs fallidos y lotes |
| `aru font:*` | descubrimiento de fuentes incorporadas y gestión de assets del proyecto |

```sh
brew install arandu-io/tap/aru
aru version
aru help
```

<a name="use-the-task-guide-behind-each-command"></a>
## Consulta la guía de la tarea detrás de cada comando

`aru new` :aruVersion realiza un clon superficial del proyecto inicial inmutable :aranduVersion, exige Git y red, crea secretos nuevos, intenta compilar las vistas y no ejecuta `go mod tidy` ni inicia el servidor. Comprueba las dependencias resultantes; los comandos delegados fallan si el binario del proyecto no los implementa.

`aru doctor` también avisa cuando un proyecto se aparta del layout del builder: bootstrap-file-beside-app, program-under-cmd, command-outside-commands y program-under-scripts. El catálogo dice lo que existe. Las guías de instalación, migration, cola, despliegue y editor son dueñas de las decisiones operativas, la seguridad y la recuperación.

- [Instalación](/docs/{{version}}/start/installation)
- [Migraciones](/docs/{{version}}/database/migrations)
- [Operaciones de colas](/docs/{{version}}/operations/queues)
- [Despliegue](/docs/{{version}}/operations/deployment)
- [Extensión del editor](/docs/{{version}}/ai/assisted-development)
