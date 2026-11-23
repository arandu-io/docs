# Kits iniciais

- [Comece com código próprio e continue responsável por ele](#start-from-owned-source-and-keep-owning-it)
- [Crie um projeto e fixe a versão do publicador](#create-a-project-and-pin-the-publisher)
- [Visualize e publique o kit de autenticação](#preview-and-publish-the-authentication-kit)
- [Republique sem perder código da aplicação silenciosamente](#republish-without-silently-losing-application-code)
- [Rejeite um fluxo parcialmente publicado ou que não compila](#refuse-a-half-published-or-unbuildable-flow)
- [Confira novamente o fluxo de identidade publicado como código da aplicação](#recheck-the-published-identity-flow-as-application-code)
- [Teste o fluxo gerado, não apenas a lista de arquivos](#test-the-generated-journey-not-only-the-file-list)
- [Continue com autenticação, Kyse e atualizações](#continue-with-authentication-kyse-and-upgrades)

<a name="start-from-owned-source-and-keep-owning-it"></a>
## Comece com código próprio e continue responsável por ele

O Arandu tem dois caminhos de entrega de código-fonte. `aru new` cria a aplicação-base a partir do esqueleto; o módulo UI publica um fluxo de autenticação em um projeto existente. Nenhum deles é uma abstração de runtime: os arquivos gerados e publicados se tornam código da aplicação, que deve ser revisado, editado, testado e atualizado deliberadamente.

<a name="create-a-project-and-pin-the-publisher"></a>
## Crie um projeto e fixe a versão do publicador

Instale aru, Git e Go. `aru new` cria um projeto, remove o histórico do esqueleto, escreve um ambiente privado com segredos novos e tenta compilar as views. Faça commit ou backup do projeto antes de executar o publicador UI com tag fixa e confira se `arandu.toml` declara a versão mínima de aru suportada.

- Registre ou verifique o commit do projeto inicial.
- Use uma tag exata do publicador UI.
- Comece com uma árvore de trabalho limpa.
- Execute primeiro o dry-run do publicador.

<a name="preview-and-publish-the-authentication-kit"></a>
## Visualize e publique o kit de autenticação

Visualize e depois publique a UI :uiVersion. Ela escreve arquivos Kyse e Go pertencentes à aplicação e imprime a ligação sem editar `bootstrap/app.go`: o import dos controllers publicados e uma chamada a `authui.New` para listar em `WithProviders`, antes de app, construída a partir das regras de conta do provider e dos drivers, `drivers.Limiter` incluído. Aplique essa ligação, remova de `WithProviders` qualquer módulo de autenticação anterior, construa as views e migre.

```sh
go run github.com/arandu-io/ui@:uiVersion auth --dry-run
go run github.com/arandu-io/ui@:uiVersion auth
aru view:build
aru migrate
```

<a name="republish-without-silently-losing-application-code"></a>
## Republique sem perder código da aplicação silenciosamente

Fixe o publicador na automação e preserve as mudanças da aplicação dentro dos marcadores arandu custom. `--views` republica só as telas e deixa os controllers publicados intactos, então uma correção de segurança num controller exige `--force` depois de revisar o diff; os limites dos formulários de código e as respostas temporizadas da ui :uiVersion moram nos controllers e só chegam com auth `--force`. Cinco arquivos acoplados de layout e página substituem o conteúdo fora dos marcadores mesmo sem force.

- Revise todos os caminhos escritos pelo publicador.
- Siga a ligação impressa para `WithProviders`.
- Mantenha código customizado dentro dos marcadores.
- Compile e teste as rotas após cada publicação.

<a name="refuse-a-half-published-or-unbuildable-flow"></a>
## Rejeite um fluxo parcialmente publicado ou que não compila

O publicador rejeita uma raiz que não seja de projeto, caminho de módulo inválido, argumentos inesperados e versão mínima da aru antiga antes de escrever. Uma falha sequencial do sistema de arquivos ainda pode deixar uma publicação parcial; o controle de versão e a repetição da mesma tag fixa permitem recuperá-la.

> [!WARNING]
> `aru new` :aruVersion seleciona o esqueleto imutável :aranduVersion e não oferece uma flag para substituí-lo. Confira o `go.mod` dessa tag e o projeto gerado antes de considerar a base reproduzível.

<a name="recheck-the-published-identity-flow-as-application-code"></a>
## Confira novamente o fluxo de identidade publicado como código da aplicação

> [!WARNING]
> O kit fornece os fluxos de transporte, mas o projeto é responsável por Policy, resolução de tenant, regras de senha, e-mail, sessão, CSRF e disponibilidade do cadastro. Ao editar, preserve a rotação de sessão, os códigos vinculados à finalidade, a conclusão dos fatores e as verificações de senha recente.

Telas e rotas geradas não autorizam o acesso aos dados da aplicação. Toda operação ainda passa pelo serviço e pela Policy da aplicação.

<a name="test-the-generated-journey-not-only-the-file-list"></a>
## Teste o fluxo gerado, não apenas a lista de arquivos

Partindo de um estado limpo, execute o dry-run, publique a tag exata, confira os arquivos e a composição, compile as views, rode vet, testes com detector de corrida e build. Confira as rotas de autenticação, teste os fluxos de identidade comuns e YPY e republique para comprovar a preservação dos marcadores e o escopo de `--views`.

- O dry-run não escreve arquivos.
- O código publicado compila.
- Os fluxos de autenticação chegam aos serviços publicados.
- Blocos customizados sobrevivem à republicação.

<a name="continue-with-authentication-kyse-and-upgrades"></a>
## Continue com autenticação, Kyse e atualizações

A autenticação define a identidade, o Kyse compila as views publicadas, a extensão do editor entende o projeto e as atualizações preservam os diffs pertencentes à aplicação.

- [Autenticação](/docs/{{version}}/security/authentication)
- [Templates Kyse](/docs/{{version}}/frontend/kyse-templates)
- [Extensão do editor](/docs/{{version}}/ai/assisted-development)
- [Atualizações](/docs/{{version}}/operations/upgrades)
