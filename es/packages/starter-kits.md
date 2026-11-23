# Kits iniciales

- [Empieza con código propio y sigue siendo responsable de él](#start-from-owned-source-and-keep-owning-it)
- [Crea un proyecto y fija la versión del publicador](#create-a-project-and-pin-the-publisher)
- [Previsualiza y publica el kit de autenticación](#preview-and-publish-the-authentication-kit)
- [Vuelve a publicar sin perder silenciosamente código de la aplicación](#republish-without-silently-losing-application-code)
- [Rechaza un flujo parcialmente publicado o que no compile](#refuse-a-half-published-or-unbuildable-flow)
- [Vuelve a comprobar el flujo de identidad publicado como código de la aplicación](#recheck-the-published-identity-flow-as-application-code)
- [Prueba el flujo generado, no solo la lista de archivos](#test-the-generated-journey-not-only-the-file-list)
- [Continúa con autenticación, Kyse y actualizaciones](#continue-with-authentication-kyse-and-upgrades)

<a name="start-from-owned-source-and-keep-owning-it"></a>
## Empieza con código propio y sigue siendo responsable de él

Arandu tiene dos vías de entrega de código fuente. `aru new` crea la aplicación base desde el esqueleto; el módulo UI publica un flujo de autenticación en un proyecto existente. Ninguno es una abstracción de runtime: los archivos generados y publicados se convierten en código de la aplicación que debe revisarse, editarse, probarse y actualizarse de forma deliberada.

<a name="create-a-project-and-pin-the-publisher"></a>
## Crea un proyecto y fija la versión del publicador

Instala aru, Git y Go. `aru new` crea un proyecto, elimina el historial del esqueleto, escribe un entorno privado con secretos nuevos e intenta compilar las vistas. Haz un commit o una copia de seguridad del proyecto antes de ejecutar el publicador UI con tag fija y comprueba que `arandu.toml` declare la versión mínima de aru soportada.

- Registra o verifica el commit del proyecto inicial.
- Usa una tag exacta del publicador UI.
- Empieza con un árbol de trabajo limpio.
- Ejecuta primero el dry-run del publicador.

<a name="preview-and-publish-the-authentication-kit"></a>
## Previsualiza y publica el kit de autenticación

Previsualice y después publique UI :uiVersion. Escribe archivos Kyse y Go que pertenecen a la aplicación e imprime el cableado sin editar `bootstrap/app.go`: el import de los controllers publicados y una llamada a `authui.New` para listar en `WithProviders`, antes de app, construida a partir de las reglas de cuenta del provider y de los drivers, `drivers.Limiter` incluido. Aplique ese cableado, quite de `WithProviders` cualquier módulo de autenticación anterior, construya las vistas y migre.

```sh
go run github.com/arandu-io/ui@:uiVersion auth --dry-run
go run github.com/arandu-io/ui@:uiVersion auth
aru view:build
aru migrate
```

<a name="republish-without-silently-losing-application-code"></a>
## Vuelve a publicar sin perder silenciosamente código de la aplicación

Fije el publicador en la automatización y conserve los cambios de la aplicación dentro de los marcadores arandu custom. `--views` republica solo las pantallas y deja intactos los controllers publicados, así que una corrección de seguridad en un controller exige `--force` tras revisar el diff; los límites de los formularios de código y las respuestas temporizadas de ui :uiVersion viven en los controllers y solo llegan con auth `--force`. Cinco archivos acoplados de layout y página reemplazan el contenido fuera de los marcadores incluso sin force.

- Revisa todas las rutas que escribe el publicador.
- Siga el cableado impreso para `WithProviders`.
- Mantén el código personalizado dentro de los marcadores.
- Compila y prueba las rutas después de cada publicación.

<a name="refuse-a-half-published-or-unbuildable-flow"></a>
## Rechaza un flujo parcialmente publicado o que no compile

El publicador rechaza una raíz que no sea de proyecto, una ruta de módulo inválida, argumentos inesperados y una versión mínima antigua de aru antes de escribir. Un fallo secuencial del sistema de archivos aún puede dejar una publicación parcial; el control de versiones y la repetición de la misma tag fija permiten recuperarla.

> [!WARNING]
> `aru new` :aruVersion selecciona el esqueleto inmutable :aranduVersion y no ofrece una flag para sustituirlo. Comprueba el `go.mod` de esa tag y el proyecto generado antes de considerar reproducible la base.

<a name="recheck-the-published-identity-flow-as-application-code"></a>
## Vuelve a comprobar el flujo de identidad publicado como código de la aplicación

> [!WARNING]
> El kit proporciona los flujos de transporte, pero el proyecto es responsable de Policy, resolución de tenant, reglas de contraseña, correo, sesión, CSRF y disponibilidad del registro. Al editar, conserva la rotación de sesión, los códigos vinculados a su finalidad, la finalización de factores y las comprobaciones de contraseña reciente.

Las pantallas y rutas generadas no autorizan el acceso a los datos de la aplicación. Toda operación sigue pasando por el servicio y la Policy de la aplicación.

<a name="test-the-generated-journey-not-only-the-file-list"></a>
## Prueba el flujo generado, no solo la lista de archivos

Desde un estado limpio, ejecuta el dry-run, publica la tag exacta, inspecciona los archivos y la composición, compila las vistas y ejecuta vet, pruebas con detector de carreras y build. Comprueba las rutas de autenticación, prueba los flujos de identidad normales y YPY y vuelve a publicar para demostrar la conservación de marcadores y el alcance de `--views`.

- El dry-run no escribe archivos.
- El código publicado compila.
- Los flujos de autenticación llegan a los servicios publicados.
- Los bloques personalizados sobreviven a la republicación.

<a name="continue-with-authentication-kyse-and-upgrades"></a>
## Continúa con autenticación, Kyse y actualizaciones

La autenticación define la identidad, Kyse compila las vistas publicadas, la extensión del editor entiende el proyecto y las actualizaciones conservan los diffs de la aplicación.

- [Autenticación](/docs/{{version}}/security/authentication)
- [Plantillas Kyse](/docs/{{version}}/frontend/kyse-templates)
- [Extensión del editor](/docs/{{version}}/ai/assisted-development)
- [Actualizaciones](/docs/{{version}}/operations/upgrades)
