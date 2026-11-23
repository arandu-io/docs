# Estrutura de diretórios

- [O diretório indica a quem o código pertence](#a-directory-states-who-owns-the-code)
- [Execute os comandos na raiz do projeto](#run-commands-from-the-project-root)
- [Gere um controller dentro da aplicação](#generate-a-controller-into-the-application-boundary)
- [Verifique a estrutura e os imports em conjunto](#check-structure-and-imports-together)
- [A árvore leva à raiz de composição](#the-tree-leads-to-the-composition-root)

> [!WARNING]
> **Limitações atuais.** O contrato dos diretórios pertence à aplicação e é verificado na versão marcada do próprio projeto inicial. Não há um consumidor independente com tag, pois ele descreve a árvore desse projeto, não uma integração reutilizável de runtime.

<a name="a-directory-states-who-owns-the-code"></a>
## O diretório indica a quem o código pertence

O projeto é um módulo Go comum, com limites explícitos para aplicação, bootstrap, configuração, banco de dados, views, rotas, arquivos gerados e testes. Coloque cada alteração onde seu responsável possa ser importado e testado; não mova comportamento para main, pois os testes da aplicação não podem importar o pacote main.

<a name="run-commands-from-the-project-root"></a>
## Execute os comandos na raiz do projeto

A raiz é o diretório que contém `go.mod`, `main.go` e `arandu.toml`. A Aru descobre o projeto por esses três arquivos. `main.go` constrói a aplicação com `bootstrap.App` e entrega seus argumentos a `HandleCommand`; sem argumentos é serve. Os templates de origem ficam em resources/views; seus pacotes Go gerados ficam em `storage/framework/views` e são reconstruídos, nunca editados à mão.

| Caminho | Responsável |
| --- | --- |
| app/Http | tipos de requisição, middlewares e controllers |
| app/Services e app/Policies | casos de uso e autorização |
| app/Providers | providers que constroem os serviços, controllers e handlers sobre os drivers |
| `app/Console/Commands` | os comandos de console da própria aplicação |
| `bootstrap/app.go` | o único arquivo de composição; bootstrap não guarda mais nada |
| config | as seções tipadas da própria aplicação ao lado de `bootstrap.Settings` |
| `routes/web.go` e `routes/console.go` | rotas HTTP; a lista de comandos e o agendamento |
| database | migrations, factories e seeders |
| resources/views | código-fonte Kyse editável |
| `storage/framework/views` | Go gerado |
| tests/E2E | jornadas em navegador real atrás da build tag e2e, rodadas por `aru test --e2e` |

<a name="generate-a-controller-into-the-application-boundary"></a>
## Gere um controller dentro da aplicação

```sh
aru make:controller InvoiceController
aru route:list
```

O gerador escreve `app/Http/Controllers/InvoiceController.go` e imprime o restante. O provider em app/Providers constrói o controller e o entrega a `routes/web.go` pelo valor `routes.Deps` que `bootstrap/app.go` passa para `WithRouting`, onde a ação é registrada por nome; a geração não esconde essa ligação.

> [!WARNING]
> Nunca edite um arquivo em `storage/framework/views` como fonte de uma tela. O comando `aru view:build` substitui os arquivos gerados a partir dos arquivos `.kyse.go`.

<a name="check-structure-and-imports-together"></a>
## Verifique a estrutura e os imports em conjunto

```sh
aru view:build
aru doctor --strict
go test ./...
go build ./...
```

- Os pacotes da aplicação não importam main.
- As rotas recebem controllers já construídos, em vez de montar suas dependências.
- Os pacotes de views gerados estão atualizados e os templates de origem continuam editáveis.
- As novas migrations são importadas pela raiz de composição e aparecem em `aru migrate:status`.
- `aru doctor` não relata nenhum aviso de layout: nenhum arquivo ao lado de `bootstrap/app.go`, nenhum programa em cmd ou scripts e nenhum comando fora de `app/Console/Commands`.

<a name="the-tree-leads-to-the-composition-root"></a>
## A árvore leva à raiz de composição

As fronteiras de diretório são úteis porque o grafo de build as impõe. A configuração alimenta os drivers, os providers constroem serviços sobre eles, `bootstrap/app.go` compõe providers, rotas, middleware e comandos, e os testes podem importar os mesmos pacotes sem iniciar um listener.

- [Módulos e raiz de composição](/docs/{{version}}/architecture/modules-and-composition)
- [Ciclo de vida da requisição](/docs/{{version}}/start/request-lifecycle)
- [Views](/docs/{{version}}/frontend/views)
- [CLI aru](https://pkg.go.dev/github.com/arandu-io/aru)
