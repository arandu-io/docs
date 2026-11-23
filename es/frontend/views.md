# Vistas

- [El nombre de una vista resuelve a Go compilado](#a-view-name-resolves-to-compiled-go)
- [Renderiza datos tipados de extremo a extremo](#render-typed-data-end-to-end)
- [Registra los paquetes de páginas, layouts y fragmentos](#register-page-layout-and-fragment-packages)
- [Diagnostica nombres, imports y salidas obsoletas](#diagnose-names-imports-and-stale-output)
- [Verifica build, escape y formato de respuesta](#prove-build-escaping-and-response-shape)
- [Registro de vistas y renderizador publicados](#released-view-registry-and-renderer)
- [Escribe y distribuye el fuente de la vista](#write-and-ship-the-view-source)

<a name="a-view-name-resolves-to-compiled-go"></a>
## El nombre de una vista resuelve a Go compilado

Un archivo `.kyse.go` declara su contrato de datos y su marcado. `aru view:build` escribe Go generado en `storage/framework/views`; importar cada paquete generado ejecuta `view.Register`, y `framework/view.NewModule` proporciona el único renderizador usado por `Context.View` y `Context.Fragment`.

<a name="render-typed-data-end-to-end"></a>
## Renderiza datos tipados de extremo a extremo

Define los datos de la página en el fuente y pasa exactamente el tipo generado desde el controller. Un campo mal escrito falla al compilarse en la línea del fuente; un tipo de datos incorrecto en runtime devuelve un error de renderizado identificado en vez de mostrar valores cero.

```kyse
//go:build kyse

package views

@go
type GreetingData struct { Name string }
@endgo

<h1>Hello {{ .Name }}</h1>
```

> [!NOTE]
> Los datos de la página son un contenido público. `Context.View` responde a una solicitud que envía Accept: `application/vnd.arandu.view+json` con el nombre de la vista y todos los campos exportados de los datos en JSON, incluidos los que la plantilla nunca muestra. Nunca ponga secretos, hashes de contraseña, tokens ni claves internas de tenant en los datos de la vista; pase solo lo que la página puede mostrar.

<a name="register-page-layout-and-fragment-packages"></a>
## Registra los paquetes de páginas, layouts y fragmentos

Importa cada directorio generado en el bootstrap, porque los directorios anidados son paquetes Go distintos. Una fuente anidada como `invoices/show.kyse.go` registra un nombre separado por punto, y una vista de formulario sin layout es un destino de fragmento separado.

- Usa `ctx.View` con un valor de página tipado para documentos completos.
- Usa `ctx.Fragment` con una vista tipada sin layout para los intercambios.
- Mantén la composición de datos en controllers y servicios.
- Vuelve a generar en vez de editar `storage/framework/views`.

<a name="diagnose-names-imports-and-stale-output"></a>
## Diagnostica nombres, imports y salidas obsoletas

La ausencia de vistas registradas significa que faltan los paquetes generados o no se han importado. Un nombre con puntos desconocido, un tipo de datos incorrecto, un registro duplicado, una salida generada antigua o la falta de `RendererProvider` deben interrumpir visiblemente la respuesta o el inicio correspondiente.

> [!WARNING]
> Las vistas generadas de la aplicación normalmente se ignoran en Git. Un checkout limpio debe ejecutar `aru view:build` antes de que Go resuelva controllers que importan tipos generados. Las vistas comunes usan buffer de salida y no son respuestas de streaming ilimitado.

<a name="prove-build-escaping-and-response-shape"></a>
## Verifica build, escape y formato de respuesta

Compila las vistas, renderiza con el tipo exacto de datos y comprueba escape de HTML, tipo de contenido, raíces de página completa y fragmento, fallos por nombre desconocido y datos incorrectos, y la presencia en `view.Registered` de cada nombre literal usado por un controller.

- `aru view:build` y `go build` terminan correctamente.
- El texto del usuario se escapa.
- Un tipo de datos incorrecto produce un error.
- Todos los nombres de vistas usados por controllers registrados se resuelven.

<a name="released-view-registry-and-renderer"></a>
## Registro de vistas y renderizador publicados

`view.Register` y Registered forman el registro generado, `view.NewRenderer` renderiza datos tipados en un buffer, `view.WrongData` estandariza las discrepancias de tipo, `Context.View` y Fragment eligen el formato HTTP y `framework/view.NewModule` proporciona el renderizador y las rutas de assets.

| Superficie | Contrato publicado |
| --- | --- |
| `view.Register` / Registered | registro generado y diagnósticos |
| `view.NewRenderer` / `WrongData` | renderizado en buffer y discrepancia de tipo |
| `Context.View` / Fragment | Selección de respuesta HTTP |
| `framework/view.NewModule` | proveedor de renderizado y assets |

<a name="write-and-ship-the-view-source"></a>
## Escribe y distribuye el fuente de la vista

Kyse explica la sintaxis de las plantillas, los assets permiten rastrear los bytes compilados, los controllers componen datos tipados y los auxiliares de respuesta eligen el documento o fragmento final.

- [Plantillas Kyse](/docs/{{version}}/frontend/kyse-templates)
- [Assets](/docs/{{version}}/frontend/assets)
- [Controllers](/docs/{{version}}/http/controllers)
