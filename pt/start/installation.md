# Instalação

- [Instale uma CLI e compile um programa Go comum](#install-one-cli-build-one-ordinary-go-program)
- [Fixe a versão da CLI antes de criar o projeto](#pin-the-cli-before-creating-the-project)
- [Crie e execute um projeto local](#create-and-run-a-local-project)
- [Verifique o projeto gerado, não apenas o download](#prove-the-generated-project-not-just-the-download)
- [Continue pela raiz de composição](#continue-from-the-composition-root)

<a name="install-one-cli-build-one-ordinary-go-program"></a>
## Instale uma CLI e compile um programa Go comum

Projetos Arandu são compilados em um único binário Go. A CLI aru cria o projeto, delega a esse binário os comandos que dependem da aplicação e compila as views Kyse; ela não é uma dependência de runtime nem um segundo servidor.

O esqueleto publicado exige Go 1.26.10 ou mais novo; uma toolchain mais antiga se recusa a compilá-lo. Git e acesso à rede também são necessários enquanto a CLI baixa o esqueleto e as primeiras ferramentas de build com versões fixas.

<a name="pin-the-cli-before-creating-the-project"></a>
## Fixe a versão da CLI antes de criar o projeto

```sh
GOWORK=off go install -ldflags='-X main.version=:aruVersion' github.com/arandu-io/aru@:aruVersion
aru version
```

> [!NOTE]
> A tag fixa o código-fonte da CLI, e a flag do linker define a versão que a CLI informa e compara com a versão mínima exigida pelo projeto. Um `go install` simples nessa tag informa dev, pois a ferramenta Go não aplica as flags de release do linker. O `aru new` cria o projeto a partir da versão mais recente do esqueleto, a :aranduVersion.

<a name="create-and-run-a-local-project"></a>
## Crie e execute um projeto local

Escolha um nome de diretório e um caminho de módulo Go. O comando copia `.env.example` para `.env` com uma `APP_KEY` nova, reescreve os imports para o módulo escolhido e tenta compilar as views. Ele mantém o SQLite como banco de desenvolvimento, e o `APP_ENV=dev` copiado é o que torna o processo um processo de desenvolvimento: um processo sem `APP_ENV` roda como produção.

```sh
aru new ledger --module example.com/acme/ledger
cd ledger
aru migrate
aru dev
```

> [!NOTE]
> A primeira inicialização do servidor não precisa de seed. Para criar um administrador, rode o seeder nomeado com credenciais explícitas, `aru db:seed UserSeeder -e you@example.com -p '<password>' -r` admin: `DatabaseSeeder` não cria nenhuma conta, e as credenciais pertencem ao comando que você roda, não ao ambiente nem ao código.

> [!NOTE]
> O `aru new` não executa `go mod tidy` nem inicia o servidor. Se a compilação das views falhar, o projeto é criado mesmo assim, e o comando indica que você rode `aru view:build`.

<a name="prove-the-generated-project-not-just-the-download"></a>
## Verifique o projeto gerado, não apenas o download

Confira as versões resolvidas dos módulos no `go.mod` e execute as verificações a partir da raiz do projeto. Instalar a CLI com sucesso não garante que a configuração, as views, o schema ou os testes da aplicação gerada estejam corretos.

```sh
aru view:build
aru doctor --strict
go test ./...
go build ./...
```

- `aru version` informa :aruVersion.
- O `go.mod` declara as versões de Framework, Hesape, Kyse e conectores que você pretende publicar.
- `aru doctor --strict`, testes e build terminam com sucesso.
- O `.env` contém uma chave exclusiva da aplicação e não é commitado.

<a name="continue-from-the-composition-root"></a>
## Continue pela raiz de composição

Leia o guia de diretórios antes de mover arquivos: pacotes de views gerados, templates, configuração, rotas e raiz de composição têm responsabilidades diferentes. A configuração é carregada uma vez na inicialização, e o ciclo da requisição só começa depois que o grafo explícito de módulos foi construído.

- [Estrutura de diretórios](/docs/{{version}}/start/directory-structure)
- [Configuração](/docs/{{version}}/architecture/configuration)
- [Ciclo de vida da requisição](/docs/{{version}}/start/request-lifecycle)
- [Referência do pacote aru](https://pkg.go.dev/github.com/arandu-io/aru)
