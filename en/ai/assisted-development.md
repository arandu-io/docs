# AI-assisted development

- [Editor intelligence, not an application AI runtime](#editor-intelligence-not-an-application-ai-runtime)
- [Install the released extension beside Aru](#install-the-released-extension-beside-aru)
- [Choose one project and open its source map](#choose-one-project-and-open-its-source-map)
- [Prove the editor and CLI agree on the selected root](#prove-the-editor-and-cli-agree-on-the-selected-root)
- [Continue with the CLI or generated views](#continue-with-the-cli-or-generated-views)

> [!WARNING]
> **Current limitations.** The VS Code extension is preview tooling for project selection, commands, Kyse syntax, snippets, and Aru integration. It does not provide an in-application AI runtime or a private assistant bridge, and command availability still depends on the selected project and local Aru installation.

<a name="editor-intelligence-not-an-application-ai-runtime"></a>
## Editor intelligence, not an application AI runtime

The first-party VS Code extension owns the local editing loop around an Arandu project. The extension supplies Kyse syntax highlighting and snippets locally; `aru lsp` supplies completion and diagnostics. It also provides a source-linked Project Map, Doctor findings in Problems, and explicit controls for `aru dev`. :vscodeVersion ships no model client, agent loop, prompt runtime, or in-application assistant.

<a name="install-the-released-extension-beside-aru"></a>
## Install the released extension beside Aru

Use VS Code 1.96 or newer, install the current Aru CLI, download arandu-:vscodeVersionNoPrefix.vsix from the :vscodeVersion release, and install it in a trusted local workspace containing `arandu.toml`. Virtual workspaces are unsupported and untrusted workspaces do not start Aru.

- Discover nested and multi-root projects from `arandu.toml`.
- Select a project explicitly when several exist.
- Set resource-scoped `arandu.aru.path` only for a custom executable.
- Keep runtime AI and the private bridge outside this preview tooling page.

<a name="choose-one-project-and-open-its-source-map"></a>
## Choose one project and open its source map

Open the Arandu activity container and select one discovered root. That root drives LSP, Project Map, Doctor, watchers, development terminal, and visible update task; source-located nodes open the file at the recorded line.

```sh
brew install arandu-io/tap/aru
code --install-extension ./arandu-:vscodeVersionNoPrefix.vsix
```

> [!NOTE]
> Aru resolution is resource-scoped `arandu.aru.path`, PATH, `/opt/homebrew/bin/aru`, then `/usr/local/bin/aru`.

<a name="prove-the-editor-and-cli-agree-on-the-selected-root"></a>
## Prove the editor and CLI agree on the selected root

Open a Kyse file and confirm `aru lsp` completion or diagnostics, run Doctor and see fresh Problems replace stale ones, refresh Project Map after a relevant save, open a located node, and start then stop `aru dev` through visible commands.

- All views use the same selected root.
- Doctor replaces stale diagnostics.
- Project Map source locations open correctly.
- Development terminal status follows explicit start and stop.

<a name="continue-with-the-cli-or-generated-views"></a>
## Continue with the CLI or generated views

The extension never runs migrations, seeders, or generators; it checks stable Aru at most once per day and upgrades Homebrew only after an explicit choice. Aru remains authority for LSP, Doctor, graph, and serving.

Install and inspect the extension at its public repository and tagged releases.

- [CLI reference](/docs/{{version}}/cli/reference)
- [Kyse templates](/docs/{{version}}/frontend/kyse-templates)
- [VS Code extension](https://github.com/arandu-io/vscode-arandu)
- [Extension releases](https://github.com/arandu-io/vscode-arandu/releases)
