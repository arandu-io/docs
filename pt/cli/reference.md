# Console e geração de código da aru

- [A aru é o adapter do binário pertencente ao projeto](#aru-is-the-adapter-around-a-project-owned-binary)
- [Famílias de comandos publicadas na :aruVersion](#released-command-families-in-aruversion)
- [Consulte o guia da tarefa por trás de cada comando](#use-the-task-guide-behind-each-command)

> [!WARNING]
> **Limitações atuais.** O inventário de comandos, os geradores de código, a compilação de views e o doctor estão publicados. Migrations, persistência de factories e recuperação de filas mantêm as limitações descritas em seus próprios guias; a CLI não torna esses fluxos mais seguros do que seus contratos permitem.

<a name="aru-is-the-adapter-around-a-project-owned-binary"></a>
## A aru é o adapter do binário pertencente ao projeto

A aru cria e inspeciona projetos, compila Kyse, gera Go comum e delega os comandos que dependem da aplicação ao binário do projeto. Execute `aru help` para consultar o inventário, `aru about` para o contrato do projeto e `aru version` para a CLI instalada. Esta página descreve a :aruVersion; comandos existentes apenas na branch de desenvolvimento não entram aqui.

<a name="released-command-families-in-aruversion"></a>
## Famílias de comandos publicadas na :aruVersion

Identifique com about, version, `action:list` e `imports:catalog`. Crie, rode e verifique com new, `key:generate`, dev, serve, build, test, `route:list`, trace, doctor, lsp, mcp, `view:build`, `vendor:publish`, `skills:sync` e `openapi:check`; a aplicação nativa tem `native:run`, `native:dev`, `native:build` e `native:icons`. As jornadas em navegador real são escritas com `make:e2e`, gravadas com `e2e:record`, explicadas após uma falha com `e2e:explain` e preparadas com `e2e:install-browsers`. Gere com schema, generate, `model:build`, `make:module`, `make:model`, `make:service`, `make:resource`, `make:notification`, `make:client`, `make:mcp-tool`, `make:mcp-resource`, `make:mcp-prompt`, `make:migration`, `make:controller`, `make:middleware`, `make:request`, `make:factory`, `make:seeder`, `make:job`, `make:mail`, `make:command`, `make:listener`, `make:event`, `make:enum`, `make:policy`, `make:test`, `make:contract` e `make:e2e`. `model:build` reescreve o arquivo de query ao lado de cada entidade, que nunca é editado à mão; todo outro arquivo gerado pertence à aplicação e precisa ser comparado e compilado, e um gerador imprime a ligação para app/Providers, as rotas e o builder em vez de escrevê-la.

Os comandos de banco de dados são migrate, `migrate:rollback`, `migrate:status`, `migrate:fresh`, que só roda com `APP_ENV=dev`, `db:seed`, `schedule:list` e `schedule:run`. Os comandos de fila são `queue:work`, listen, restart, pause, resume, clear, monitor, failed, retry, forget, flush, prune-failed, retry-batch e prune-batches. Os comandos de fontes são `font:search`, info, add, list e remove; arquivos, licenças e CSS vendorizados viram assets do projeto.

| Superfície | Contrato publicado |
| --- | --- |
| `aru new <name> --module <path>` | clonar e adaptar um projeto inicial, criar ambiente privado e configurações do editor |
| `aru dev` / serve / build / test | compilar views e rodar, empacotar ou testar o projeto |
| `aru route:list` / trace / doctor / lsp | inspecionar contratos de runtime e do editor |
| `aru mcp` / `skills:sync` / `vendor:publish` / `openapi:check` | responder a um assistente, manter Skills e arquivos publicados em dia e verificar contratos de API |
| `aru native:run` / `native:dev` / `native:build` | construir e rodar a aplicação nativa |
| `aru schema` / generate / `make:*` | geração de código da aplicação orientada por schema ou por comandos específicos |
| `aru migrate*` / `db:seed` / `schedule:*` | operações delegadas de banco de dados e agendamento |
| `aru queue:*` | operações delegadas de workers, sinais, backlog, jobs com falha e lotes |
| `aru font:*` | descoberta de fontes incorporadas e gerenciamento de assets do projeto |

```sh
brew install arandu-io/tap/aru
aru version
aru help
```

<a name="use-the-task-guide-behind-each-command"></a>
## Consulte o guia da tarefa por trás de cada comando

O `aru new` :aruVersion faz um clone superficial do projeto inicial imutável :aranduVersion, exige Git e rede, cria segredos novos, tenta compilar as views e não executa `go mod tidy` nem inicia o servidor. Confira as dependências resultantes; comandos delegados falham se o binário do projeto não os implementar.

`aru doctor` também avisa quando um projeto se afasta do layout do builder: bootstrap-file-beside-app, program-under-cmd, command-outside-commands e program-under-scripts. O catálogo diz o que existe. Os guias de instalação, migration, fila, implantação e editor são donos das escolhas operacionais, da segurança e da recuperação.

- [Instalação](/docs/{{version}}/start/installation)
- [Migrations](/docs/{{version}}/database/migrations)
- [Operações de filas](/docs/{{version}}/operations/queues)
- [Publicação](/docs/{{version}}/operations/deployment)
- [Extensão do editor](/docs/{{version}}/ai/assisted-development)
