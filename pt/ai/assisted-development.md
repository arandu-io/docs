# Desenvolvimento assistido por IA

- [Inteligência do editor, não um runtime de IA na aplicação](#editor-intelligence-not-an-application-ai-runtime)
- [Instale a extensão publicada junto da aru](#install-the-released-extension-beside-aru)
- [Escolha um projeto e abra seu mapa de código](#choose-one-project-and-open-its-source-map)
- [Verifique se editor e CLI usam a mesma raiz](#prove-the-editor-and-cli-agree-on-the-selected-root)
- [Continue pela CLI ou pelas views geradas](#continue-with-the-cli-or-generated-views)

<a name="editor-intelligence-not-an-application-ai-runtime"></a>
## Inteligência do editor, não um runtime de IA na aplicação

A extensão oficial para VS Code cuida do fluxo de edição local de um projeto Arandu. Ela fornece realce de sintaxe Kyse e snippets localmente; a `aru lsp` fornece autocompletar e diagnósticos. Também oferece um Project Map com links para o código, achados do Doctor em Problems e controles explícitos para `aru dev`. A :vscodeVersion não inclui cliente de modelo, ciclo de agente, runtime de prompts ou assistente dentro da aplicação.

<a name="install-the-released-extension-beside-aru"></a>
## Instale a extensão publicada junto da aru

Use VS Code 1.96 ou posterior, instale a CLI aru atual, baixe arandu-:vscodeVersionNoPrefix.vsix da release v:vscodeVersionNoPrefix e instale-o em um workspace local confiável com `arandu.toml`. Workspaces virtuais não são suportados, e workspaces não confiáveis não iniciam a aru.

- Descubra projetos aninhados e com várias raízes a partir de `arandu.toml`.
- Selecione explicitamente um projeto quando houver vários.
- Defina `arandu.aru.path` no escopo do recurso apenas para usar um executável personalizado.
- Não procure aqui um runtime de IA nem uma ponte privada: a extensão não traz nenhum dos dois.

<a name="choose-one-project-and-open-its-source-map"></a>
## Escolha um projeto e abra seu mapa de código

Abra o contêiner de atividades Arandu e selecione uma das raízes encontradas. Essa raiz orienta o LSP, Project Map, Doctor, watchers, terminal de desenvolvimento e tarefa visível de atualização; os nós com localização abrem o arquivo na linha registrada.

```sh
brew install arandu-io/tap/aru
code --install-extension ./arandu-:vscodeVersionNoPrefix.vsix
```

A resolução da aru segue esta ordem: `arandu.aru.path` no escopo do recurso, PATH, `/opt/homebrew/bin/aru` e `/usr/local/bin/aru`.

<a name="prove-the-editor-and-cli-agree-on-the-selected-root"></a>
## Verifique se editor e CLI usam a mesma raiz

Abra um arquivo Kyse e confira o autocompletar ou os diagnósticos da `aru lsp`. Execute o Doctor e verifique se os novos problemas substituem os antigos. Atualize o Project Map após salvar uma alteração relevante, abra um nó com localização e inicie e pare `aru dev` pelos comandos visíveis.

- Todas as views usam a mesma raiz selecionada.
- O Doctor substitui os diagnósticos antigos.
- Os locais de código-fonte do Project Map abrem corretamente.
- O status do terminal de desenvolvimento acompanha os comandos explícitos de início e parada.

<a name="continue-with-the-cli-or-generated-views"></a>
## Continue pela CLI ou pelas views geradas

A extensão nunca executa migrations, seeders ou geradores. Ela verifica a versão estável da aru no máximo uma vez por dia e só atualiza pelo Homebrew após uma escolha explícita. A aru continua sendo a autoridade para LSP, Doctor, grafo e servidor.

Instale e consulte a extensão em seu repositório público e nas releases com tag.

- [Referência da CLI](/docs/{{version}}/cli/reference)
- [Templates Kyse](/docs/{{version}}/frontend/kyse-templates)
- [Extensão para VS Code](https://github.com/arandu-io/vscode-arandu)
- [Releases da extensão](https://github.com/arandu-io/vscode-arandu/releases)
