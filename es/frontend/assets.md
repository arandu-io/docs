# Assets sin Node

- [El binario sirve assets identificados por su contenido](#the-binary-serves-content-addressed-assets)
- [Fija la versión mínima de la CLI y la herramienta descargada](#pin-the-cli-floor-and-downloaded-tool)
- [Compila una hoja de estilos dentro de la aplicación](#compile-one-stylesheet-into-the-application)
- [Mantén visibles los fuentes, la salida y la entrega](#keep-source-output-and-serving-visible)
- [Identifica salidas obsoletas y no utilizadas](#recognize-stale-and-unconsumed-output)
- [Mantén los scripts locales y los bytes con origen identificable](#keep-scripts-local-and-bytes-attributable)
- [Verifica bytes, hashes, headers y clases](#prove-bytes-hashes-headers-and-classes)
- [Recompila desde los fuentes y verifica el resultado embebido](#rebuild-from-sources-and-verify-the-embedded-result)
- [Rastrea un asset desde el fuente hasta la respuesta](#trace-an-asset-from-source-to-response)

<a name="the-binary-serves-content-addressed-assets"></a>
## El binario sirve assets identificados por su contenido

`aru view:build` compila el CSS del proyecto con el ejecutable standalone de Tailwind de versión fija y escribe `assets/app.css`. La aplicación incluye esos bytes, y el módulo de vistas del framework sirve CSS, YPY y scripts de interacción mediante una ruta derivada de SHA-256 que cambia con el contenido.

<a name="pin-the-cli-floor-and-downloaded-tool"></a>
## Fija la versión mínima de la CLI y la herramienta descargada

Defina en `arandu.toml` el piso publicado de Aru y la versión de Tailwind. Una versión de Tailwind para la que Aru no trae checksums también necesita tailwindcss-sha256sums, el SHA-256 del archivo de checksums de esa release, en `[tools]`; sin él la descarga se rechaza antes de empezar. El compilador Kyse está dentro de Aru; solo el compilador de hojas de estilo independiente se descarga a la caché de herramientas verificada, y el servidor en ejecución no necesita ni Node ni la herramienta de build.

- Fija `[arandu].aru` en una versión mínima compatible.
- Fija `[tools].tailwindcss` en la versión revisada del binario.
- Mantén `resources/css/app.css` y los fuentes Kyse bajo revisión.
- Importa el paquete de assets de la aplicación en el bootstrap.

<a name="compile-one-stylesheet-into-the-application"></a>
## Compila una hoja de estilos dentro de la aplicación

Ejecuta el build normal de vistas antes de compilar Go. Genera las vistas de la aplicación y escribe la hoja de estilos minificada embebida. Importar el paquete de assets registra esos bytes, y `framework/view.NewModule` monta la ruta con hash.

```sh
aru view:build
go build ./...
```

<a name="keep-source-output-and-serving-visible"></a>
## Mantén visibles los fuentes, la salida y la entrega

El CSS de origen y Kyse siguen perteneciendo a la aplicación. Las vistas generadas normalmente se ignoran en Git, mientras que el consumidor oficial incluye `assets/app.css` compilado en sus commits. Los scripts del framework acompañan al módulo publicado, y `view.URL` emite una ruta local para los bytes realmente registrados.

- Edita `resources/css/app.css` y los fuentes Kyse.
- Regenera `storage/framework/views` en vez de editarlo.
- Incluye en el commit la salida `assets/app.css` del consumidor oficial.
- Mantén habilitada la generación de vistas en el `aru build` normal.

<a name="recognize-stale-and-unconsumed-output"></a>
## Identifica salidas obsoletas y no utilizadas

La ausencia de imports de assets deja silenciosamente la hoja de estilos predeterminada del framework. La falta de paquetes generados rompe los builds limpios, las salidas antiguas de fuentes eliminadas pueden duplicar nombres y un nombre de asset desconocido provoca panic durante el renderizado, en vez de inventar una URL que devolvería 404.

> [!WARNING]
> `aru build --skip-views` puede empaquetar vistas y CSS generados obsoletos. Úsalo solo en un pipeline que ya haya verificado que las salidas corresponden al mismo commit de los fuentes.

<a name="keep-scripts-local-and-bytes-attributable"></a>
## Mantén los scripts locales y los bytes con origen identificable

> [!WARNING]
> Mantén los scripts en el mismo origen, conserva los avisos de licencia, usa `view.URL` en vez de hashes escritos a mano y registra tipos de contenido correctos. Nunca sirvas archivos subidos por usuarios desde el registro confiable de assets de la aplicación ni ejecutes una herramienta del checkout sin versión fija.

<a name="prove-bytes-hashes-headers-and-classes"></a>
## Verifica bytes, hashes, headers y clases

Compila dos veces y compara los bytes. Solicita la hoja de estilos identificada por su contenido y comprueba Content-Type, nosniff, caché inmutable, clases representativas de la aplicación y el cambio de URL tras modificar un byte.

- Los fuentes sin cambios producen CSS idéntico.
- El cuerpo servido coincide con los bytes embebidos.
- Cambiar un byte cambia la URL.
- Las clases exclusivas de la aplicación se conservan tras la compilación.

<a name="rebuild-from-sources-and-verify-the-embedded-result"></a>
## Recompila desde los fuentes y verifica el resultado embebido

Comprueba la versión mínima de aru y la versión fija de Tailwind, regenera los fuentes, recompila y prueba la respuesta del asset. Después inspecciona la URL con hash en la página en ejecución. Vaciar la caché del navegador no repara un binario compilado con los bytes incorrectos.

Si la descarga no está disponible, restaura la caché verificada o compila en una etapa preparatoria con acceso a la red. El servidor de producción no tiene una etapa de publicación de assets.

<a name="trace-an-asset-from-source-to-response"></a>
## Rastrea un asset desde el fuente hasta la respuesta

La visión general explica el contrato de renderizado en el servidor, las vistas son responsables del HTML tipado, Kyse compila las plantillas y los middlewares restringen los orígenes desde los que se cargan scripts.

- [Visión general del frontend](/docs/{{version}}/frontend/overview)
- [Vistas](/docs/{{version}}/frontend/views)
- [Plantillas Kyse](/docs/{{version}}/frontend/kyse-templates)
- [Middleware](/docs/{{version}}/http/middleware)
