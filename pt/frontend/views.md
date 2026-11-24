# Views

- [O nome de uma view resolve para Go compilado](#a-view-name-resolves-to-compiled-go)
- [Renderize dados tipados de ponta a ponta](#render-typed-data-end-to-end)
- [Registre os pacotes de páginas, layouts e fragmentos](#register-page-layout-and-fragment-packages)
- [Diagnostique nomes, imports e saída desatualizada](#diagnose-names-imports-and-stale-output)
- [Verifique build, escape e formato da resposta](#prove-build-escaping-and-response-shape)
- [Registro de views e renderizador publicados](#released-view-registry-and-renderer)
- [Escreva e distribua o fonte da view](#write-and-ship-the-view-source)

<a name="a-view-name-resolves-to-compiled-go"></a>
## O nome de uma view resolve para Go compilado

Um arquivo `.kyse.go` declara seu contrato de dados e sua marcação. `aru view:build` escreve Go gerado em `storage/framework/views`; o import de cada pacote gerado executa `view.Register`, e `framework/view.NewModule` fornece o renderizador único usado por `Context.View` e `Context.Fragment`.

<a name="render-typed-data-end-to-end"></a>
## Renderize dados tipados de ponta a ponta

Defina os dados da página no fonte e passe exatamente o tipo gerado a partir do controller. Um campo escrito incorretamente falha na compilação na linha do fonte; um tipo de dados incorreto em runtime retorna um erro de renderização identificado em vez de desenhar valores zero.

```kyse
//go:build kyse

package views

@go
type GreetingData struct { Name string }
@endgo

<h1>Hello {{ .Name }}</h1>
```

> [!NOTE]
> Os dados da página são um conteúdo público. `Context.View` responde a uma requisição que envia Accept: `application/vnd.arandu.view+json` com o nome da view e todos os campos exportados dos dados em JSON, inclusive os que o template nunca exibe. Nunca coloque segredos, hashes de senha, tokens ou chaves internas de tenant nos dados da view; passe apenas o que a página pode mostrar.

<a name="register-page-layout-and-fragment-packages"></a>
## Registre os pacotes de páginas, layouts e fragmentos

Importe cada diretório gerado no bootstrap, pois diretórios aninhados são pacotes Go distintos. Uma fonte aninhada como `invoices/show.kyse.go` registra um nome separado por ponto, e uma view de formulário sem layout é um destino de fragmento separado.

- Use `ctx.View` com um valor de página tipado para documentos completos.
- Use `ctx.Fragment` com uma view tipada sem layout para as trocas.
- Mantenha a composição de dados nos controllers e serviços.
- Gere novamente em vez de editar `storage/framework/views`.

<a name="diagnose-names-imports-and-stale-output"></a>
## Diagnostique nomes, imports e saída desatualizada

A ausência de views registradas significa que os pacotes gerados estão ausentes ou não foram importados. Nome pontuado desconhecido, tipo de dados incorreto, registro duplicado, saída gerada antiga ou `RendererProvider` ausente devem interromper de forma visível a resposta ou a inicialização correspondente.

As views geradas da aplicação normalmente são ignoradas pelo Git. Um checkout limpo precisa executar `aru view:build` antes que Go resolva controllers que importam tipos gerados. Views comuns usam buffer de saída e não são respostas de streaming ilimitado.

<a name="prove-build-escaping-and-response-shape"></a>
## Verifique build, escape e formato da resposta

Compile as views, renderize com o tipo exato de dados e confira escape de HTML, tipo de conteúdo, raízes de página completa e fragmento, falhas por nome desconhecido e dados incorretos, e a presença de cada nome literal usado por um controller em `view.Registered`.

- `aru view:build` e `go build` terminam com sucesso.
- O texto do usuário é escapado.
- Um tipo de dados incorreto gera erro.
- Todos os nomes de views usados pelos controllers registrados são resolvidos.

<a name="released-view-registry-and-renderer"></a>
## Registro de views e renderizador publicados

`view.Register` e Registered formam o registro gerado, `view.NewRenderer` renderiza dados tipados em buffer, `view.WrongData` padroniza divergências de tipo, `Context.View` e Fragment escolhem o formato HTTP, e `framework/view.NewModule` fornece o renderizador e as rotas de assets.

| Superfície | Contrato publicado |
| --- | --- |
| `view.Register` / Registered | registro gerado e diagnósticos |
| `view.NewRenderer` / `WrongData` | renderização em buffer e divergência de tipo |
| `Context.View` / Fragment | Seleção da resposta HTTP |
| `framework/view.NewModule` | provedor de renderização e assets |

<a name="write-and-ship-the-view-source"></a>
## Escreva e distribua o fonte da view

O Kyse explica a sintaxe dos templates, os assets permitem rastrear os bytes compilados, os controllers compõem dados tipados e os auxiliares de resposta escolhem o documento ou fragmento final.

- [Templates Kyse](/docs/{{version}}/frontend/kyse-templates)
- [Assets](/docs/{{version}}/frontend/assets)
- [Controllers](/docs/{{version}}/http/controllers)
