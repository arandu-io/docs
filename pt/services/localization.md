# Localização

- [Um tradutor, uma decisão de idioma por requisição](#one-translator-one-locale-decision-per-request)
- [Registre catálogos e resolva uma chave](#register-catalogues-and-resolve-a-key)
- [Substitua valores, escolha plurais e adicione namespaces](#replace-values-choose-plurals-and-add-namespaces)
- [Fallback sem resolução retorna a chave](#an-unresolved-fallback-returns-the-key)
- [Teste a origem do idioma e o catálogo](#test-the-locale-source-and-the-catalogue)
- [Operações do tradutor](#translator-operations)
- [Rotas e texto possuem contratos de idioma distintos](#routes-and-prose-have-separate-language-contracts)

> [!WARNING]
> **Limitações atuais.** O carregamento de traduções, o fallback, `NewFromDirectory` e o `WithTranslator` do builder estão publicados e testados, mas o esqueleto publicado não traz catálogo e não chama `WithTranslator`. Registre o catálogo e o locale da requisição explicitamente e depois verifique o caminho completo da aplicação.

<a name="one-translator-one-locale-decision-per-request"></a>
## Um tradutor, uma decisão de idioma por requisição

Translator carrega frases por chave de caminhos registrados e mantém grupos em cache com segurança para requisições concorrentes. O idioma pertence ao contexto da requisição, escolhido uma vez por middleware ou escopo explícito. Código seguinte lê essa decisão sem reinterpretar URL ou header de idioma.

<a name="register-catalogues-and-resolve-a-key"></a>
## Registre catálogos e resolva uma chave

Construa o tradutor uma vez a partir do diretório de catálogos com `translation.NewFromDirectory`, nomeando o locale de resposta e o de fallback, e entregue-o ao builder com `WithTranslator`. Toda requisição passa então a levar o catálogo ao lado do nome da aplicação, de modo que `view.Page` responde Translator e Locale e um componente desenha as próprias palavras no idioma da requisição. Get verifica a cadeia de fallback configurada; quando nem o catálogo pedido nem o de fallback resolvem o item, o handler de chave ausente roda e a própria chave não resolvida é devolvida.

```go example
func WithTranslations(drivers *foundation.Drivers) (*foundation.Application, error) {
	translator, err := translation.NewFromDirectory(os.DirFS("."), "lang", "en", "en")
	if err != nil {
		return nil, err
	}

	return foundation.Configure(drivers).WithTranslator(translator).Create()
}
```

<a name="replace-values-choose-plurals-and-add-namespaces"></a>
## Substitua valores, escolha plurais e adicione namespaces

Get substitui placeholders nomeados, Choice seleciona o plural pelo número, Has verifica cobertura e `AddNamespace` associa traduções de pacote ao próprio hint. `InLocale` serve a trabalho fora da requisição; quando middleware já fornece idioma, use essa fonte única em vez de aninhar outra escolha.

<a name="an-unresolved-fallback-returns-the-key"></a>
## Fallback sem resolução retorna a chave

O fallback mantém uma página utilizável quando o idioma está incompleto, mas não inventa texto. Se o item continuar ausente e o handler de chaves não fornecer substituição, Get retorna a chave, expondo o defeito do catálogo. Testes de paridade devem comparar conjuntos de chaves e recusar valores copiados do idioma fonte, salvo exceções curtas e justificadas. Um idioma não suportado nunca deve virar caminho de leitura arbitrária de arquivos.

<a name="test-the-locale-source-and-the-catalogue"></a>
## Teste a origem do idioma e o catálogo

Envie requisições por todos os prefixos de idioma suportados e confira interface traduzida e rotas estáveis. Exercite substituições e limites de plural, depois procure chaves ausentes, extras e valores suspeitosamente idênticos. O teste também deve verificar quantidade plausível de chaves para um loader quebrado não passar sem verificar nada.

<a name="translator-operations"></a>
## Operações do tradutor

`NewFromDirectory` constrói o tradutor da aplicação e Locales lista os locales que seu diretório traz, que é a lista contra a qual o middleware de locale negocia. New, Get, Choice, Has e `HasForLocale` formam o caminho normal de leitura. `AddPath` e `AddJSONPath` registram catálogos da aplicação; `AddNamespace` registra catálogos de pacotes; `DetermineLocalesUsing` personaliza a ordem de fallback; `HandleMissingKeysUsing` observa chaves não resolvidas. `SetLocale` muda o padrão e não substitui o contexto de locale por requisição.

<a name="routes-and-prose-have-separate-language-contracts"></a>
## Rotas e texto possuem contratos de idioma distintos

Interface, datas, números e frases do produto pertencem à tradução. Identificadores de API, exemplos Go e corpo do manual podem permanecer deliberadamente no idioma fonte quando esse for o contrato publicado. Helpers de URL devem receber o idioma escolhido pelo contexto de roteamento, evitando divergência entre links, middleware e catálogos.
