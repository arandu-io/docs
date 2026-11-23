# Requisições HTTP

- [Leia dados do transporte sem inventar identidade](#read-transport-data-without-inventing-identity)
- [Monte um contrato tipado de entrada](#assemble-one-typed-input-contract)
- [Escolha o leitor que preserva o significado](#choose-the-reader-that-preserves-meaning)
- [Limite o corpo antes de decodificar](#bound-the-body-before-decoding-it)
- [Exercite cada canal de entrada e limite](#exercise-every-input-channel-and-limit)
- [Superfícies de requisição publicadas](#released-request-surfaces)
- [Passe do transporte a uma operação confiável](#move-from-transport-to-a-trusted-operation)

<a name="read-transport-data-without-inventing-identity"></a>
## Leia dados do transporte sem inventar identidade

A fronteira preferida é Context enxuto do framework com primitivas `net/http`. Leia parâmetros de rota, query, headers e corpo limitado apenas onde a ação precisa; Request amplo do Hesape é compatibilidade, não dependência obrigatória de controller.

<a name="assemble-one-typed-input-contract"></a>
## Monte um contrato tipado de entrada

Use Bind para ler apenas os campos aceitos num tipo de entrada dedicado, cujas tags form nomeiam cada campo que pode chegar, normalize valores canônicos e passe só dados validados ao serviço. Um tipo com método Rules tem essas regras executadas primeiro sobre a entrada bruta. Bind escreve apenas campos com tag, apara cada valor e informa um valor que não converte como erro do campo, em vez de um zero silencioso. Preserve `ctx.Ctx` para manter cancelamento e identidade autenticada nas chamadas seguintes.

```go example
type CreateInvoiceInput struct {
	Reference string `form:"reference"`
	Amount    int64  `form:"amount"`
}

func Create(ctx *fhttp.Context) error {
	var input CreateInvoiceInput
	if err := ctx.Bind(&input); err != nil {
		return err
	}
	return nil
}
```

<a name="choose-the-reader-that-preserves-meaning"></a>
## Escolha o leitor que preserva o significado

Use middleware de proxy confiável antes de ler endereço de cliente derivado do par. Limite corpos antes de decodificar, feche streams de upload e separe identificadores de rota dos campos graváveis, impedindo que formulários substituam propriedade ou tenant. Input e Bind leem um único mapa: a query string em GET e HEAD, e apenas o corpo em POST, PUT, PATCH e DELETE, então um valor na URL da ação é lido explicitamente com Query.

- Use Param para segmentos de rota correspondentes.
- Use Query para controles opcionais da URL.
- Use Bind para campos de formulário declarados; métodos inseguros leem apenas o corpo.
- Use Request padrão para headers, corpo e contexto.

<a name="bound-the-body-before-decoding-it"></a>
## Limite o corpo antes de decodificar

Números malformados, tipos de mídia não suportados, corpos excessivos, valores obrigatórios ausentes e leituras canceladas devem permanecer explícitos. Não converta silenciosamente um valor inválido em zero nem continue após um decoder parar em payload parcial.

> [!WARNING]
> Headers encaminhados são entrada do chamador até `TrustProxies` aceitar o par imediato. Nunca derive throttles de autenticação, identidade de auditoria ou links assinados de encaminhamento não confiável.

<a name="exercise-every-input-channel-and-limit"></a>
## Exercite cada canal de entrada e limite

Envie requisições válidas, ausentes, malformadas, excessivas e canceladas pela rota real. Verifique a entrada tipada enviada ao serviço e comprove que campos de tenant ou propriedade do chamador são ignorados ou recusados.

- Apenas entrada declarada chega ao serviço.
- Limites de corpo e upload atuam antes de alocação sem teto.
- O cancelamento interrompe parsing e trabalho seguinte.
- Valores de encaminhamento não confiáveis não substituem o endereço do par.

<a name="released-request-surfaces"></a>
## Superfícies de requisição publicadas

`http.Context` do framework expõe Ctx, Param, Query, Input, Bind, User, helpers de resposta e request e writer subjacentes. `net/http` padrão possui headers, corpos, cookies e contexto. Request mais amplo do Hesape continua disponível quando seus helpers de compatibilidade forem necessários.

<a name="move-from-transport-to-a-trusted-operation"></a>
## Passe do transporte a uma operação confiável

Validação transforma entrada do transporte em entrada tipada que satisfaz regras declaradas. Controllers escolhem o caso de uso, proxies confiáveis estabelecem a fronteira de rede e respostas codificam o resultado.

- [Validação](/docs/{{version}}/http/validation)
- [Controllers](/docs/{{version}}/http/controllers)
- [Respostas](/docs/{{version}}/http/responses)
