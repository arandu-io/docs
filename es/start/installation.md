# Instalación

- [Instala una CLI y compila un programa Go común](#install-one-cli-build-one-ordinary-go-program)
- [Fija la versión de la CLI antes de crear el proyecto](#pin-the-cli-before-creating-the-project)
- [Crea y ejecuta un proyecto local](#create-and-run-a-local-project)
- [Verifica el proyecto generado, no solo la descarga](#prove-the-generated-project-not-just-the-download)
- [Continúa por la raíz de composición](#continue-from-the-composition-root)

> [!WARNING]
> **Limitaciones actuales.** El comando `aru new` publicado exige Git y acceso a la red, realiza un clon superficial del esqueleto inmutable :aranduVersion e intenta compilar las vistas sin interrumpir la creación si esa etapa falla. Verifica las dependencias del proyecto generado antes de usarlo en producción y ejecuta las migraciones como una etapa explícita del pipeline.

<a name="install-one-cli-build-one-ordinary-go-program"></a>
## Instala una CLI y compila un programa Go común

Los proyectos Arandu se compilan en un único binario Go. La CLI aru crea el proyecto, delega en ese binario los comandos que dependen de la aplicación y compila las vistas Kyse; no es una dependencia de runtime ni un segundo servidor.

El esqueleto publicado exige Go 1.26.10 o más nuevo; una toolchain más antigua se niega a compilarlo. También se necesitan Git y acceso a la red mientras la CLI descarga el esqueleto y las primeras herramientas de build con versiones fijas.

<a name="pin-the-cli-before-creating-the-project"></a>
## Fija la versión de la CLI antes de crear el proyecto

```sh
GOWORK=off go install -ldflags='-X main.version=:aruVersion' github.com/arandu-io/aru@:aruVersion
aru version
```

> [!NOTE]
> La tag fija el código fuente de la CLI, y la flag del linker establece la versión que la CLI informa y compara con la mínima exigida por el proyecto. Un `go install` simple en esta tag informa dev porque la herramienta Go no aplica las flags de release del linker. Por separado, aru :aruVersion fija el esqueleto del proyecto en :aranduVersion.

> [!WARNING]
> `aru new` en :aruVersion exige Git y acceso a la red y realiza un clon superficial del esqueleto inmutable :aranduVersion con una sola rama. No ofrece una flag para elegir la versión del esqueleto; verifica el `go.mod` generado antes de continuar cuando necesites una base exacta.

<a name="create-and-run-a-local-project"></a>
## Crea y ejecuta un proyecto local

Elija un nombre de directorio y una ruta de módulo Go. El comando copia `.env.example` a `.env` con una `APP_KEY` nueva, reescribe los imports al módulo elegido e intenta compilar las vistas. Mantiene SQLite como base de datos de desarrollo, y el `APP_ENV=dev` copiado es lo que convierte el proceso en uno de desarrollo: un proceso sin `APP_ENV` se ejecuta como producción.

```sh
aru new ledger --module example.com/acme/ledger
cd ledger
aru migrate
aru dev
```

> [!NOTE]
> El primer inicio del servidor no necesita seed. Para crear un administrador, ejecuta el seeder con nombre y credenciales explícitas, `aru db:seed UserSeeder -e you@example.com -p '<password>' -r` admin: `DatabaseSeeder` no crea ninguna cuenta, y las credenciales pertenecen al comando que ejecutas, no al entorno ni al código.

> [!WARNING]
> Por diseño, `aru new` no ejecuta `go mod tidy` ni inicia el servidor. Una compilación fallida de las vistas produce una advertencia y el proyecto permanece en disco; ejecuta `aru view:build` antes del build, sin dar por hecho que ya existen las vistas generadas.

<a name="prove-the-generated-project-not-just-the-download"></a>
## Verifica el proyecto generado, no solo la descarga

Comprueba las versiones resueltas de los módulos en `go.mod` y ejecuta las verificaciones desde la raíz del proyecto. Instalar correctamente la CLI no garantiza que la configuración, las vistas, el esquema o las pruebas de la aplicación generada sean correctos.

```sh
aru view:build
aru doctor --strict
go test ./...
go build ./...
```

- `aru version` informa :aruVersion.
- `go.mod` declara las versiones de Framework, Hesape, Kyse y conectores que quieres desplegar.
- `aru doctor --strict`, las pruebas y el build terminan correctamente.
- `.env` contiene una clave exclusiva de la aplicación y no se incluye en los commits.

<a name="continue-from-the-composition-root"></a>
## Continúa por la raíz de composición

Lee la guía de directorios antes de mover archivos: los paquetes de vistas generados, las plantillas, la configuración, las rutas y la raíz de composición tienen responsabilidades diferentes. La configuración se carga una vez al iniciar, y el ciclo de la solicitud solo empieza después de construir el grafo explícito de módulos.

- [Estructura de directorios](/docs/{{version}}/start/directory-structure)
- [Configuración](/docs/{{version}}/architecture/configuration)
- [Ciclo de vida de la solicitud](/docs/{{version}}/start/request-lifecycle)
- [Referencia del paquete aru](https://pkg.go.dev/github.com/arandu-io/aru)
