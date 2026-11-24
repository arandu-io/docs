# Directory structure

- [A directory states who owns the code](#a-directory-states-who-owns-the-code)
- [Run commands from the project root](#run-commands-from-the-project-root)
- [Generate a controller into the application boundary](#generate-a-controller-into-the-application-boundary)
- [Check structure and imports together](#check-structure-and-imports-together)
- [The tree leads to the composition root](#the-tree-leads-to-the-composition-root)

<a name="a-directory-states-who-owns-the-code"></a>
## A directory states who owns the code

The project is an ordinary Go module with explicit application, bootstrap, configuration, database, view, route, generated, and test boundaries. Put a change where its owner can be imported and tested; do not move behavior into main, because package main cannot be imported by the application tests.

<a name="run-commands-from-the-project-root"></a>
## Run commands from the project root

The root is the directory containing `go.mod`, `main.go`, and `arandu.toml`. Aru discovers the project from those three files. `main.go` builds the application with `bootstrap.App` and hands its arguments to `HandleCommand`; no arguments is serve. Source templates live under resources/views; their generated Go packages live under `storage/framework/views` and are rebuilt rather than hand-edited.

| Path | Owner |
| --- | --- |
| app/Http | request types, middleware, controllers |
| app/Services and app/Policies | use cases and authorization |
| app/Providers | providers that build the services, controllers, and handlers over the drivers |
| `app/Console/Commands` | the application's own console commands |
| `bootstrap/app.go` | the one composition file; bootstrap holds nothing else |
| config | the application's own typed sections beside `bootstrap.Settings` |
| `routes/web.go` and `routes/console.go` | HTTP routes; the command list and the schedule |
| database | migrations, factories, seeders |
| resources/views | editable Kyse source |
| `storage/framework/views` | generated Go |
| tests/E2E | real-browser journeys behind the e2e build tag, run by `aru test --e2e` |

<a name="generate-a-controller-into-the-application-boundary"></a>
## Generate a controller into the application boundary

```sh
aru make:controller InvoiceController
aru route:list
```

The generator writes `app/Http/Controllers/InvoiceController.go` and prints the rest. The provider in app/Providers constructs the controller and hands it to `routes/web.go` through the `routes.Deps` value `bootstrap/app.go` gives `WithRouting`, where the action is registered by name; generation does not hide that wiring.

> [!WARNING]
> Never edit a file under `storage/framework/views` as the source of a screen. `aru view:build` replaces generated output from the `.kyse.go` files.

<a name="check-structure-and-imports-together"></a>
## Check structure and imports together

```sh
aru view:build
aru doctor --strict
go test ./...
go build ./...
```

- Application packages do not import main.
- Routes receive constructed controllers rather than building dependencies.
- Generated view packages are current and source templates remain editable.
- New migrations are imported by the composition root and visible to `aru migrate:status`.
- `aru doctor` reports no layout warning: no file beside `bootstrap/app.go`, no program under cmd or scripts, and no command outside `app/Console/Commands`.

<a name="the-tree-leads-to-the-composition-root"></a>
## The tree leads to the composition root

Directory boundaries are useful because the build graph enforces them. Configuration feeds the drivers, providers build services over them, `bootstrap/app.go` composes providers, routes, middleware, and commands, and tests can import the same packages without starting a listener.

- [Modules and composition root](/docs/{{version}}/architecture/modules-and-composition)
- [Request lifecycle](/docs/{{version}}/start/request-lifecycle)
- [Views](/docs/{{version}}/frontend/views)
- [Aru CLI](https://pkg.go.dev/github.com/arandu-io/aru)
