# Installation

- [Install one CLI, build one ordinary Go program](#install-one-cli-build-one-ordinary-go-program)
- [Pin the CLI before creating the project](#pin-the-cli-before-creating-the-project)
- [Create and run a local project](#create-and-run-a-local-project)
- [Prove the generated project, not just the download](#prove-the-generated-project-not-just-the-download)
- [Continue from the composition root](#continue-from-the-composition-root)

> [!WARNING]
> **Current limitations.** The released `aru new` command requires Git and network access, shallow-clones the immutable skeleton :aranduVersion, and treats the generated-view build as best effort. Verify the generated module graph before production use, and run migrations as an explicit pipeline step.

<a name="install-one-cli-build-one-ordinary-go-program"></a>
## Install one CLI, build one ordinary Go program

Arandu projects compile into a single Go binary. The aru CLI creates the project, delegates application-aware commands to that binary, and builds Kyse views; it is not a runtime dependency or a second server.

The released skeleton requires Go 1.26.10 or newer; an older toolchain refuses to build it. Git and network access are also needed while the CLI fetches the skeleton and the first pinned build tools.

<a name="pin-the-cli-before-creating-the-project"></a>
## Pin the CLI before creating the project

```sh
GOWORK=off go install -ldflags='-X main.version=:aruVersion' github.com/arandu-io/aru@:aruVersion
aru version
```

> [!NOTE]
> The tag pins the CLI source, and the linker flag stamps the version that the CLI reports and compares with a project's required floor. A plain `go install` at this tag reports dev because the Go tool does not apply release linker flags. Aru :aruVersion independently pins the project skeleton to :aranduVersion.

> [!WARNING]
> `aru new` in :aruVersion requires Git and network access and shallow-clones the immutable skeleton :aranduVersion with a single branch. It exposes no skeleton-version flag; verify the generated `go.mod` before continuing when an exact project baseline matters.

<a name="create-and-run-a-local-project"></a>
## Create and run a local project

Choose a directory name and a Go module path. The command copies `.env.example` to `.env` with a fresh `APP_KEY`, rewrites imports to the chosen module, and attempts to compile the views. It leaves SQLite as the development database, and the copied `APP_ENV=dev` is what makes the process a development one: a process with no `APP_ENV` runs as production.

```sh
aru new ledger --module example.com/acme/ledger
cd ledger
aru migrate
aru dev
```

> [!NOTE]
> The first server boot needs no seed. When you want an admin, run the named seeder with explicit credentials, `aru db:seed UserSeeder -e you@example.com -p '<password>' -r` admin: `DatabaseSeeder` creates no account, and the credentials belong to the command you run, not to the environment or the source.

> [!WARNING]
> `aru new` deliberately does not run `go mod tidy` or start the server. A failed view build is reported as a warning and the project remains on disk; run `aru view:build` before building rather than assuming generated views exist.

<a name="prove-the-generated-project-not-just-the-download"></a>
## Prove the generated project, not just the download

Confirm the resolved module versions in `go.mod`, then run the project checks from its root. A successful CLI installation alone says nothing about the generated application's configuration, views, schema, or tests.

```sh
aru view:build
aru doctor --strict
go test ./...
go build ./...
```

- `aru version` reports :aruVersion.
- `go.mod` names the framework, Hesape, Kyse, and connector versions you intend to deploy.
- `aru doctor --strict`, tests, and build all exit successfully.
- `.env` contains a unique application key and is not committed.

<a name="continue-from-the-composition-root"></a>
## Continue from the composition root

Read the directory guide before moving files: generated view packages, source templates, configuration, routes, and the composition root have different ownership. Configuration is loaded once at boot, and the request lifecycle begins only after the explicit module graph has been built.

- [Directory structure](/docs/{{version}}/start/directory-structure)
- [Configuration](/docs/{{version}}/architecture/configuration)
- [Request lifecycle](/docs/{{version}}/start/request-lifecycle)
- [Aru package reference](https://pkg.go.dev/github.com/arandu-io/aru)
