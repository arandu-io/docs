# Estructura de directorios

- [El directorio indica a quién pertenece el código](#a-directory-states-who-owns-the-code)
- [Ejecuta los comandos desde la raíz del proyecto](#run-commands-from-the-project-root)
- [Genera un controller dentro de la aplicación](#generate-a-controller-into-the-application-boundary)
- [Verifica la estructura y los imports en conjunto](#check-structure-and-imports-together)
- [El árbol conduce a la raíz de composición](#the-tree-leads-to-the-composition-root)

<a name="a-directory-states-who-owns-the-code"></a>
## El directorio indica a quién pertenece el código

El proyecto es un módulo Go común, con límites explícitos para aplicación, bootstrap, configuración, base de datos, vistas, rutas, archivos generados y pruebas. Coloca cada cambio donde su responsable pueda importarse y probarse; no muevas comportamiento a main, porque las pruebas de la aplicación no pueden importar el paquete main.

<a name="run-commands-from-the-project-root"></a>
## Ejecuta los comandos desde la raíz del proyecto

La raíz es el directorio que contiene `go.mod`, `main.go` y `arandu.toml`. Aru descubre el proyecto por esos tres archivos. `main.go` construye la aplicación con `bootstrap.App` y entrega sus argumentos a `HandleCommand`; sin argumentos es serve. Las plantillas fuente están en resources/views; sus paquetes Go generados están en `storage/framework/views` y se reconstruyen, nunca se editan a mano.

| Ruta | Responsable |
| --- | --- |
| app/Http | tipos de solicitud, middlewares y controllers |
| app/Services y app/Policies | casos de uso y autorización |
| app/Providers | providers que construyen los servicios, controllers y handlers sobre los drivers |
| `app/Console/Commands` | los comandos de consola propios de la aplicación |
| `bootstrap/app.go` | el único archivo de composición; bootstrap no guarda nada más |
| config | las secciones tipadas propias de la aplicación junto a `bootstrap.Settings` |
| `routes/web.go` y `routes/console.go` | rutas HTTP; la lista de comandos y la programación |
| database | migraciones, factories y seeders |
| resources/views | código fuente Kyse editable |
| `storage/framework/views` | Go generado |
| tests/E2E | recorridos en navegador real detrás de la build tag e2e, ejecutados con `aru test --e2e` |

<a name="generate-a-controller-into-the-application-boundary"></a>
## Genera un controller dentro de la aplicación

```sh
aru make:controller InvoiceController
aru route:list
```

El generador escribe `app/Http/Controllers/InvoiceController.go` e imprime el resto. El provider en app/Providers construye el controller y lo entrega a `routes/web.go` a través del valor `routes.Deps` que `bootstrap/app.go` pasa a `WithRouting`, donde la acción se registra por nombre; la generación no oculta ese cableado.

> [!WARNING]
> Nunca edites un archivo de `storage/framework/views` como fuente de una pantalla. `aru view:build` sustituye los archivos generados a partir de los archivos `.kyse.go`.

<a name="check-structure-and-imports-together"></a>
## Verifica la estructura y los imports en conjunto

```sh
aru view:build
aru doctor --strict
go test ./...
go build ./...
```

- Los paquetes de la aplicación no importan main.
- Las rutas reciben controllers ya construidos en vez de montar sus dependencias.
- Los paquetes de vistas generados están actualizados y las plantillas de origen siguen siendo editables.
- Las nuevas migraciones se importan desde la raíz de composición y aparecen en `aru migrate:status`.
- `aru doctor` no informa ningún aviso de layout: ningún archivo junto a `bootstrap/app.go`, ningún programa en cmd o scripts y ningún comando fuera de `app/Console/Commands`.

<a name="the-tree-leads-to-the-composition-root"></a>
## El árbol conduce a la raíz de composición

Los límites de directorio son útiles porque el grafo de build los impone. La configuración alimenta los drivers, los providers construyen servicios sobre ellos, `bootstrap/app.go` compone providers, rutas, middleware y comandos, y los tests pueden importar los mismos paquetes sin iniciar un listener.

- [Módulos y raíz de composición](/docs/{{version}}/architecture/modules-and-composition)
- [Ciclo de vida de la solicitud](/docs/{{version}}/start/request-lifecycle)
- [Vistas](/docs/{{version}}/frontend/views)
- [CLI aru](https://pkg.go.dev/github.com/arandu-io/aru)
