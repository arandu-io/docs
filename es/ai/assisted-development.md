# Desarrollo asistido por IA

- [Inteligencia del editor, no un runtime de IA en la aplicación](#editor-intelligence-not-an-application-ai-runtime)
- [Instala la extensión publicada junto con aru](#install-the-released-extension-beside-aru)
- [Elige un proyecto y abre su mapa de código](#choose-one-project-and-open-its-source-map)
- [Verifica que el editor y la CLI usen la misma raíz](#prove-the-editor-and-cli-agree-on-the-selected-root)
- [Continúa con la CLI o las vistas generadas](#continue-with-the-cli-or-generated-views)

<a name="editor-intelligence-not-an-application-ai-runtime"></a>
## Inteligencia del editor, no un runtime de IA en la aplicación

La extensión oficial para VS Code gestiona el flujo de edición local de un proyecto Arandu. Proporciona resaltado de sintaxis Kyse y snippets localmente; `aru lsp` proporciona autocompletado y diagnósticos. También ofrece un Project Map enlazado al código, resultados de Doctor en Problems y controles explícitos para `aru dev`. La :vscodeVersion no incluye cliente de modelo, bucle de agente, runtime de prompts ni asistente dentro de la aplicación.

<a name="install-the-released-extension-beside-aru"></a>
## Instala la extensión publicada junto con aru

Usa VS Code 1.96 o posterior, instala la CLI aru actual, descarga arandu-:vscodeVersionNoPrefix.vsix de la release v:vscodeVersionNoPrefix e instálalo en un workspace local de confianza que contenga `arandu.toml`. Los workspaces virtuales no están soportados y los no confiables no inician aru.

- Descubre proyectos anidados y con varias raíces a partir de `arandu.toml`.
- Selecciona explícitamente un proyecto cuando haya varios.
- Configura `arandu.aru.path` en el ámbito del recurso solo para usar un ejecutable personalizado.
- No busques aquí un runtime de IA ni un puente privado: la extensión no incluye ninguno de los dos.

<a name="choose-one-project-and-open-its-source-map"></a>
## Elige un proyecto y abre su mapa de código

Abre el contenedor de actividades Arandu y selecciona una de las raíces encontradas. Esa raíz determina LSP, Project Map, Doctor, watchers, terminal de desarrollo y tarea visible de actualización; los nodos con ubicación abren el archivo en la línea registrada.

```sh
brew install arandu-io/tap/aru
code --install-extension ./arandu-:vscodeVersionNoPrefix.vsix
```

La resolución de aru sigue este orden: `arandu.aru.path` en el ámbito del recurso, PATH, `/opt/homebrew/bin/aru` y `/usr/local/bin/aru`.

<a name="prove-the-editor-and-cli-agree-on-the-selected-root"></a>
## Verifica que el editor y la CLI usen la misma raíz

Abre un archivo Kyse y comprueba el autocompletado o los diagnósticos de `aru lsp`. Ejecuta Doctor y verifica que los nuevos problemas sustituyan a los anteriores. Actualiza Project Map después de guardar un cambio relevante, abre un nodo con ubicación e inicia y detén `aru dev` mediante los comandos visibles.

- Todas las vistas usan la misma raíz seleccionada.
- Doctor sustituye los diagnósticos anteriores.
- Las ubicaciones de código fuente del Project Map se abren correctamente.
- El estado del terminal de desarrollo sigue los comandos explícitos de inicio y parada.

<a name="continue-with-the-cli-or-generated-views"></a>
## Continúa con la CLI o las vistas generadas

La extensión nunca ejecuta migraciones, seeders ni generadores. Comprueba la versión estable de aru como máximo una vez al día y solo actualiza mediante Homebrew tras una elección explícita. Aru sigue siendo la autoridad para LSP, Doctor, grafo y servidor.

Instala y consulta la extensión en su repositorio público y en sus releases etiquetadas.

- [Referencia de la CLI](/docs/{{version}}/cli/reference)
- [Plantillas Kyse](/docs/{{version}}/frontend/kyse-templates)
- [Extensión para VS Code](https://github.com/arandu-io/vscode-arandu)
- [Releases de la extensión](https://github.com/arandu-io/vscode-arandu/releases)
