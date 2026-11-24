# Assets sem Node

- [O binário serve assets endereçados pelo conteúdo](#the-binary-serves-content-addressed-assets)
- [Fixe a versão mínima da CLI e a ferramenta baixada](#pin-the-cli-floor-and-downloaded-tool)
- [Compile uma folha de estilos para dentro da aplicação](#compile-one-stylesheet-into-the-application)
- [Mantenha visíveis os fontes, a saída e a entrega](#keep-source-output-and-serving-visible)
- [Identifique saídas antigas e não utilizadas](#recognize-stale-and-unconsumed-output)
- [Mantenha scripts locais e bytes com origem identificável](#keep-scripts-local-and-bytes-attributable)
- [Verifique bytes, hashes, headers e classes](#prove-bytes-hashes-headers-and-classes)
- [Recompile a partir dos fontes e confira o resultado embutido](#rebuild-from-sources-and-verify-the-embedded-result)
- [Rastreie um asset do fonte à resposta](#trace-an-asset-from-source-to-response)

<a name="the-binary-serves-content-addressed-assets"></a>
## O binário serve assets endereçados pelo conteúdo

`aru view:build` compila o CSS do projeto com o executável standalone do Tailwind de versão fixa e escreve `assets/app.css`. A aplicação embute esses bytes, e o módulo de views do framework serve CSS, YPY e scripts de interação por um caminho derivado de SHA-256 que muda junto com o conteúdo.

<a name="pin-the-cli-floor-and-downloaded-tool"></a>
## Fixe a versão mínima da CLI e a ferramenta baixada

Defina em `arandu.toml` o piso publicado da Aru e a versão do Tailwind. Uma versão do Tailwind para a qual a Aru não traz checksums também precisa de tailwindcss-sha256sums, o SHA-256 do arquivo de checksums dessa release, em `[tools]`; sem ele o download é recusado antes de começar. O compilador Kyse fica dentro da Aru; só o compilador de folhas de estilo avulso é baixado para o cache de ferramentas verificado, e o servidor em execução não precisa nem de Node nem da ferramenta de build.

- Fixe `[arandu].aru` em uma versão mínima compatível.
- Fixe `[tools].tailwindcss` na versão revisada do binário.
- Mantenha `resources/css/app.css` e os fontes Kyse sob revisão.
- Importe o pacote de assets da aplicação no bootstrap.

<a name="compile-one-stylesheet-into-the-application"></a>
## Compile uma folha de estilos para dentro da aplicação

Execute o build normal de views antes da compilação Go. Ele gera as views da aplicação e escreve a folha de estilos minificada embutida. O import do pacote de assets registra esses bytes, e `framework/view.NewModule` monta a rota com hash.

```sh
aru view:build
go build ./...
```

<a name="keep-source-output-and-serving-visible"></a>
## Mantenha visíveis os fontes, a saída e a entrega

O CSS de origem e o Kyse continuam pertencendo à aplicação. Views geradas normalmente são ignoradas pelo Git, enquanto o consumidor oficial commita `assets/app.css` compilado. Os scripts do framework acompanham o módulo publicado, e `view.URL` emite um caminho local para os bytes realmente registrados.

- Edite `resources/css/app.css` e os fontes Kyse.
- Gere `storage/framework/views` novamente em vez de editá-lo.
- Commite a saída `assets/app.css` do consumidor oficial.
- Mantenha a geração de views habilitada no `aru build` normal.

<a name="recognize-stale-and-unconsumed-output"></a>
## Identifique saídas antigas e não utilizadas

Imports de assets ausentes deixam silenciosamente a folha de estilos padrão do framework. Pacotes gerados ausentes quebram builds limpos, saídas antigas de fontes removidas podem duplicar nomes e um nome de asset desconhecido causa panic durante a renderização, em vez de inventar uma URL que retornaria 404.

> [!WARNING]
> `aru build --skip-views` pode empacotar views e CSS gerados desatualizados. Use-o apenas em um pipeline que já tenha verificado que as saídas correspondem ao mesmo commit dos fontes.

<a name="keep-scripts-local-and-bytes-attributable"></a>
## Mantenha scripts locais e bytes com origem identificável

> [!WARNING]
> Mantenha os scripts na mesma origem, preserve os avisos de licença, use `view.URL` em vez de hashes escritos à mão e registre tipos de conteúdo corretos. Nunca sirva uploads de usuários pelo registro confiável de assets da aplicação nem execute uma ferramenta fornecida pelo checkout sem versão fixa.

<a name="prove-bytes-hashes-headers-and-classes"></a>
## Verifique bytes, hashes, headers e classes

Compile duas vezes e compare os bytes. Solicite a folha de estilos endereçada pelo conteúdo e confira Content-Type, nosniff, cache imutável, classes representativas da aplicação e a mudança de URL após alterar um byte.

- Fontes inalterados produzem CSS idêntico.
- O corpo servido corresponde aos bytes embutidos.
- Alterar um byte altera a URL.
- As classes exclusivas da aplicação permanecem após a compilação.

<a name="rebuild-from-sources-and-verify-the-embedded-result"></a>
## Recompile a partir dos fontes e confira o resultado embutido

Confira a versão mínima da aru e a versão fixa do Tailwind, gere novamente os fontes, recompile e teste a resposta do asset. Depois inspecione a URL com hash na página em execução. Limpar o cache do navegador não corrige um binário compilado com os bytes errados.

Se o download estiver indisponível, restaure o cache verificado ou compile em uma etapa preparatória com acesso à rede. O servidor de produção não tem etapa de publicação de assets.

<a name="trace-an-asset-from-source-to-response"></a>
## Rastreie um asset do fonte à resposta

A visão geral explica o contrato de renderização no servidor, as views são responsáveis pelo HTML tipado, o Kyse compila os templates e os middlewares restringem de onde os scripts podem ser carregados.

- [Visão geral do frontend](/docs/{{version}}/frontend/overview)
- [Views](/docs/{{version}}/frontend/views)
- [Templates Kyse](/docs/{{version}}/frontend/kyse-templates)
- [Middleware](/docs/{{version}}/http/middleware)
