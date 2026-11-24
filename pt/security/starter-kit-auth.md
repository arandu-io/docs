# Autenticação do starter kit

- [O publicador grava um fluxo de identidade pertencente à aplicação](#the-publisher-writes-an-application-owned-identity-flow)
- [Fixe a versão UI e proteja a árvore de trabalho](#pin-the-ui-release-and-protect-the-working-tree)
- [Visualize, publique, conecte e compile](#preview-publish-wire-and-compile)
- [Preserve blocos customizados numa republicação revisada](#preserve-custom-blocks-during-a-reviewed-republish)
- [Recuse uma superfície de autenticação parcialmente publicada](#refuse-a-half-published-authentication-surface)
- [Telas geradas não se tornam autorização](#generated-screens-do-not-become-authorization)
- [Teste cada transição de identidade publicada](#test-every-published-identity-transition)
- [Assuma o código-fonte e suas futuras atualizações](#own-the-source-and-its-future-upgrades)

<a name="the-publisher-writes-an-application-owned-identity-flow"></a>
## O publicador grava um fluxo de identidade pertencente à aplicação

O publicador auth de UI :uiVersion grava views, controllers, helpers e mailables num projeto existente. Viram fonte da aplicação, não pacote de runtime atualizável sem revisar bootstrap, policy, sessão, e-mail e decisões de tenant.

<a name="pin-the-ui-release-and-protect-the-working-tree"></a>
## Fixe a versão UI e proteja a árvore de trabalho

Comece com projeto commitado ou com backup, fixe a tag exata da UI e confirme a versão mínima de Aru em `arandu.toml`. Faça dry-run e inspecione imports e registro de bootstrap impressos: o publicador não edita `bootstrap/app.go`.

- Use `github.com/arandu-io/ui@:uiVersion`.
- Comece com diff limpo.
- Execute `--dry-run` antes de escrever.
- Revise as mudanças de composição impressas.

<a name="preview-publish-wire-and-compile"></a>
## Visualize, publique, conecte e compile

Visualize cada caminho, publique a tag exata, acrescente a `WithProviders`, antes de app, a chamada impressa a `authui.New` com `drivers.Limiter`, remova qualquer módulo de autenticação anterior, construa o Kyse e migre. O publicador não acrescenta dependência de módulo e deixa todos os arquivos escritos como propriedade da aplicação.

```sh
go run github.com/arandu-io/ui@:uiVersion auth --dry-run
go run github.com/arandu-io/ui@:uiVersion auth
aru view:build
aru migrate
```

<a name="preserve-custom-blocks-during-a-reviewed-republish"></a>
## Preserve blocos customizados numa republicação revisada

Use `--views` apenas para republicar as telas: ele deixa intactos os controllers publicados, então uma correção de segurança num controller só chega por `--force` após revisar o diff. Mantenha alterações da aplicação em marcadores `arandu:begin` custom. Cinco arquivos fortemente acoplados são sobrescritos mesmo sem `--force`, preservando blocos marcados.

- Fixe versões exatas do publicador na automação.
- Revise cada caminho substituído.
- Mantenha lógica customizada dentro dos marcadores.
- Ajuste políticas de registro, senha, tenant e conta no código próprio.

<a name="refuse-a-half-published-authentication-surface"></a>
## Recuse uma superfície de autenticação parcialmente publicada

O publicador recusa raiz sem projeto, caminho de módulo inválido, argumentos posicionais não suportados e versão mínima antiga de Aru antes de escrever. Arquivos existentes não acoplados permanecem salvo force; uma escrita bem-sucedida ainda exige composição impressa, geração de views, migration e testes.

> [!WARNING]
> Blocos customizados são correspondidos por posição e templates esperam um por arquivo. Código fora dos marcadores pode ser substituído; dry-run e diff versionado são obrigatórios antes de republicar.

<a name="generated-screens-do-not-become-authorization"></a>
## Telas geradas não se tornam autorização

O kit inclui telas de login, registro, verificação, recuperação, confirmação e dois fatores. O projeto decide abertura de registro, policy, resolução de tenant, regras de senha, e-mail, backend de sessão e ordem CSRF. Crie a sessão final só após os fatores exigidos.

Rotas e formulários publicados não são permissão. Preserve Rotate, propósitos de uso único, middleware de senha recente e verificações de Grant no servidor ao editar.

<a name="test-every-published-identity-transition"></a>
## Teste cada transição de identidade publicada

Publique de uma árvore limpa, inspecione composição, gere views, formate, execute vet, race e build; inventarie rotas de auth; exercite rejeição de login normal e YPY, registro, verificação, recuperação, recuperação de dois fatores, logout, transporte de e-mail, preservação de blocos customizados e escopo `--views`.

- Dry-run não escreve nada.
- Rotas publicadas compilam e renderizam.
- Todo fluxo de identidade alcança sua fronteira de segurança publicada.
- Republicar preserva blocos customizados.

<a name="own-the-source-and-its-future-upgrades"></a>
## Assuma o código-fonte e suas futuras atualizações

Autenticação define identidade, dois fatores completam o fluxo, Kyse compila telas e atualizações explicam os diffs do publicador pertencentes à aplicação.

- [Autenticação](/docs/{{version}}/security/authentication)
- [Sessão HTTP](/docs/{{version}}/http/session)
- [Templates Kyse](/docs/{{version}}/frontend/kyse-templates)
- [Atualizações](/docs/{{version}}/operations/upgrades)
