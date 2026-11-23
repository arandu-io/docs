# Visión general del frontend

- [HTML es la frontera de la aplicación](#html-is-the-application-boundary)
- [Edita, compila y sirve una pantalla](#edit-compile-and-serve-one-screen)
- [Asigna una responsabilidad a cada capa del frontend](#give-each-frontend-layer-one-job)
- [Un fragmento es un contrato de renderizado separado](#a-fragment-is-a-separate-rendering-contract)
- [Prueba el flujo sin `JavaScript` y el flujo mejorado](#test-the-no-javascript-and-enhanced-paths)
- [Interfaces de frontend publicadas](#released-frontend-surfaces)
- [Construye los dos formatos de respuesta](#build-the-two-response-shapes)

<a name="html-is-the-application-boundary"></a>
## HTML es la frontera de la aplicación

Arandu renderiza HTML en Go, usa YPY para solicitar fragmentos renderizados en el servidor e incluye scripts y estilos en el binario de la aplicación. No necesita router en el cliente, capa de datos SPA, runtime Node, dependencia de CDN ni una segunda API JSON para renderizar una página.

<a name="edit-compile-and-serve-one-screen"></a>
## Edita, compila y sirve una pantalla

Edita un `.kyse.go` en resources/views, compila la vista Go tipada y la hoja de estilos de la aplicación y deja que el servidor devuelva una página completa con `ctx.View` o un fragmento nombrado deliberadamente con `ctx.Fragment`.

```sh
aru view:build
aru dev
```

<a name="give-each-frontend-layer-one-job"></a>
## Asigna una responsabilidad a cada capa del frontend

Los enlaces y formularios comunes son la base. YPY mejora ese flujo mientras el servidor sigue siendo responsable de estado, CSRF, autenticación, autorización, validación y persistencia.

- Los controllers componen datos tipados de la página y eligen el formato de respuesta.
- Las vistas Kyse son responsables del marcado, nunca de localizar servicios ni autorizar operaciones.
- Los componentes tipados son responsables de controles reutilizables, no de los registros de la aplicación.
- YPY solicita e intercambia fragmentos del servidor; no se convierte en el estado del dominio.

<a name="a-fragment-is-a-separate-rendering-contract"></a>
## Un fragmento es un contrato de renderizado separado

`Context.Fragment` no elimina el layout de una vista de página completa. Usa una vista separada sin layout como destino del intercambio y conserva el estado 422 de los rechazos: YPY renderiza un 422 dentro del formulario que envió la petición, así que el estado nunca tiene que convertirse en éxito.

> [!WARNING]
> Toda modificación desde el navegador sigue exigiendo el token CSRF vinculado a la sesión y autorización en el servidor. No cargues scripts desde otro origen mientras la política `SecurityHeaders` publicada solo permita scripts del propio origen.

<a name="test-the-no-javascript-and-enhanced-paths"></a>
## Prueba el flujo sin `JavaScript` y el flujo mejorado

Solicita la página completa sin headers de mejora, envía su formulario común y repite con YPY-Request. Comprueba las raíces de documento y fragmento, estados correctos, rechazo de CSRF, etiquetas visibles, foco del teclado y un resultado útil cuando los scripts no se ejecutan.

- `aru view:build` termina correctamente desde un checkout limpio.
- Las páginas completas contienen una sola raíz de documento.
- Los fragmentos no contienen la raíz de un layout.
- Los formularios comunes y mejorados llegan al mismo servicio protegido por Policy.

<a name="released-frontend-surfaces"></a>
## Interfaces de frontend publicadas

`Context.View` y `Context.Fragment` eligen respuestas HTML. `framework/view.NewModule` proporciona el renderizador y las rutas de assets identificados por su contenido; `aru view:build` y `aru dev` compilan y recompilan los fuentes; los componentes Kyse ofrecen controles tipados renderizados en el servidor.

| Superficie | Contrato publicado |
| --- | --- |
| `Context.View` / Fragment | respuestas HTML completas y parciales |
| `framework/view.NewModule` | renderizador y assets embebidos |
| `aru view:build` / `aru dev` | compilación y ciclo de desarrollo |
| `kyse/components` | controles tipados renderizados en el servidor |

<a name="build-the-two-response-shapes"></a>
## Construye los dos formatos de respuesta

Las vistas definen el contrato de renderizado tipado, Kyse define el lenguaje de las plantillas, los assets describen la salida embebida y las respuestas HTTP conservan el estado del protocolo para ambos clientes.

- [Vistas](/docs/{{version}}/frontend/views)
- [Plantillas Kyse](/docs/{{version}}/frontend/kyse-templates)
- [Assets](/docs/{{version}}/frontend/assets)
- [Respuestas HTTP](/docs/{{version}}/http/responses)
