# Visão geral do frontend

- [HTML é a fronteira da aplicação](#html-is-the-application-boundary)
- [Edite, compile e sirva uma tela](#edit-compile-and-serve-one-screen)
- [Dê uma responsabilidade a cada camada do frontend](#give-each-frontend-layer-one-job)
- [Um fragmento é um contrato de renderização separado](#a-fragment-is-a-separate-rendering-contract)
- [Teste o fluxo sem `JavaScript` e o fluxo aprimorado](#test-the-no-javascript-and-enhanced-paths)
- [Interfaces de frontend publicadas](#released-frontend-surfaces)
- [Construa os dois formatos de resposta](#build-the-two-response-shapes)

<a name="html-is-the-application-boundary"></a>
## HTML é a fronteira da aplicação

O Arandu renderiza HTML em Go, usa YPY para solicitar fragmentos renderizados no servidor e embute scripts e estilos no binário da aplicação. Não é necessário router no cliente, camada de dados SPA, runtime Node, dependência de CDN ou uma segunda API JSON para renderizar uma página.

<a name="edit-compile-and-serve-one-screen"></a>
## Edite, compile e sirva uma tela

Edite um `.kyse.go` em resources/views, compile a view Go tipada e a folha de estilos da aplicação e deixe o servidor retornar uma página completa com `ctx.View` ou um fragmento nomeado deliberadamente com `ctx.Fragment`.

```sh
aru view:build
aru dev
```

<a name="give-each-frontend-layer-one-job"></a>
## Dê uma responsabilidade a cada camada do frontend

Links e formulários comuns são a base. O YPY aprimora esse fluxo, enquanto o servidor continua responsável por status, CSRF, autenticação, autorização, validação e persistência.

- Os controllers compõem dados tipados da página e escolhem o formato da resposta.
- Views Kyse são responsáveis pela marcação, nunca por localizar serviços ou autorizar operações.
- Componentes tipados são responsáveis por controles reutilizáveis, não pelos registros da aplicação.
- O YPY solicita e troca fragmentos do servidor; ele não se torna o estado do domínio.

<a name="a-fragment-is-a-separate-rendering-contract"></a>
## Um fragmento é um contrato de renderização separado

`Context.Fragment` não remove o layout de uma view de página completa. Use uma view separada e sem layout como destino da troca e preserve o status 422 das recusas: o YPY renderiza um 422 dentro do formulário que enviou a requisição, então o status nunca precisa virar sucesso.

> [!WARNING]
> Toda alteração feita pelo navegador ainda exige o token CSRF vinculado à sessão e autorização no servidor. Não carregue scripts de outra origem enquanto a política `SecurityHeaders` publicada permitir scripts apenas da própria origem.

<a name="test-the-no-javascript-and-enhanced-paths"></a>
## Teste o fluxo sem `JavaScript` e o fluxo aprimorado

Solicite a página completa sem headers de aprimoramento, envie seu formulário comum e repita com YPY-Request. Confira as raízes de documento e fragmento, status corretos, recusa de CSRF, rótulos visíveis, foco do teclado e um resultado útil quando os scripts não executam.

- `aru view:build` termina com sucesso a partir de um checkout limpo.
- Páginas completas contêm uma única raiz de documento.
- Fragmentos não contêm a raiz de um layout.
- Formulários comuns e aprimorados chegam ao mesmo serviço protegido por Policy.

<a name="released-frontend-surfaces"></a>
## Interfaces de frontend publicadas

`Context.View` e `Context.Fragment` escolhem respostas HTML. `framework/view.NewModule` fornece o renderizador e as rotas de assets endereçados pelo conteúdo; `aru view:build` e `aru dev` compilam e recompilam os fontes; os componentes Kyse fornecem controles tipados renderizados no servidor.

| Superfície | Contrato publicado |
| --- | --- |
| `Context.View` / Fragment | respostas HTML completas e parciais |
| `framework/view.NewModule` | renderizador e assets embutidos |
| `aru view:build` / `aru dev` | compilação e ciclo de desenvolvimento |
| `kyse/components` | controles tipados renderizados no servidor |

<a name="build-the-two-response-shapes"></a>
## Construa os dois formatos de resposta

As views definem o contrato de renderização tipado, o Kyse define a linguagem dos templates, os assets descrevem a saída embutida e as respostas HTTP preservam o status do protocolo para ambos os clientes.

- [Views](/docs/{{version}}/frontend/views)
- [Templates Kyse](/docs/{{version}}/frontend/kyse-templates)
- [Assets](/docs/{{version}}/frontend/assets)
- [Respostas HTTP](/docs/{{version}}/http/responses)
