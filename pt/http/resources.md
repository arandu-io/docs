# Recursos HTTP

- [O formato da resposta é um tipo declarado](#the-response-shape-is-a-declared-type)
- [Conclua autorização e validação primeiro](#finish-authorization-and-validation-first)
- [Implemente uma vez o contrato de recurso JSON](#implement-the-json-resource-contract-once)
- [Modele envelopes de coleção explicitamente](#model-collection-envelopes-explicitly)
- [Formulário rejeitado continua sendo requisição rejeitada](#a-rejected-form-remains-a-rejected-request)
- [Campos públicos formam uma lista explícita de permissão](#public-fields-are-an-allowlist)
- [Verifique status, formato e fragmento](#assert-status-shape-and-fragment)
- [Recursos ficam entre domínio e transporte](#resources-sit-between-domain-and-transport)

<a name="the-response-shape-is-a-declared-type"></a>
## O formato da resposta é um tipo declarado

Um recurso HTTP é o tipo que o endpoint promete retornar. Tipo e construtor são verificados pelo compilador, mas `JsonResource.ToArray` retorna `map[string]`any: chaves continuam strings de runtime. Testes exatos de formato devem detectar chaves erradas, removidas ou adicionadas por acidente. O recurso explícito mantém campos públicos independentes das colunas do model.

<a name="finish-authorization-and-validation-first"></a>
## Conclua autorização e validação primeiro

Construa recursos apenas com dados já carregados sob Grant e entrada validada em tipo próprio. O serializador não corrige autorização nem interpreta entrada. Defina status e negociação de conteúdo antes de renderizar para que clientes de página e JSON observem o mesmo resultado.

<a name="implement-the-json-resource-contract-once"></a>
## Implemente uma vez o contrato de recurso JSON

`Context.JSON` aceita `JsonResource`: `ToArray` retorna campos que podem sair e With, metadados superiores. Mantenha o model em campo não exportado, construa o recurso a partir de valor já autorizado e responda diretamente pelo handler.

```go example
type Invoice struct {
	ID, Reference string
	Total         int64
}

type InvoiceResource struct{ invoice Invoice }

func NewInvoiceResource(invoice Invoice) InvoiceResource {
	return InvoiceResource{invoice: invoice}
}

func (r InvoiceResource) ToArray() map[string]any {
	return map[string]any{
		"id": r.invoice.ID, "reference": r.invoice.Reference, "total": r.invoice.Total,
	}
}

func (InvoiceResource) With() map[string]any { return nil }

func Show(ctx *fhttp.Context, invoice Invoice) error {
	return ctx.JSON(stdhttp.StatusOK, NewInvoiceResource(invoice))
}
```

<a name="model-collection-envelopes-explicitly"></a>
## Modele envelopes de coleção explicitamente

`Context.JSON` sempre coloca `ToArray` dentro da chave superior data. Em :hesapeVersion, `ResourceCollection.ToArray` já devolve sua própria data; passá-lo a `Context.JSON` gera `data.data`. Defina um recurso de coleção ou formato de resposta explícito e fixe o objeto completo em teste, sem presumir envelope simples. Para página parcial, renderize fragmento separado cuja raiz corresponda ao alvo da troca.

> [!NOTE]
> Trate data como metadado reservado. `Context.JSON` mescla With após criar o payload, então `With()["data"]` sobrescreve o corpo do recurso em vez de ser rejeitado.

<a name="a-rejected-form-remains-a-rejected-request"></a>
## Formulário rejeitado continua sendo requisição rejeitada

Retorne 422 em falha de validação mesmo quando o navegador for inserir o fragmento na página. O YPY renderiza um 422 dentro do formulário que enviou a requisição, mantém o que foi digitado e dá foco ao primeiro campo inválido, então o status não precisa de configuração. Retornar 200 faz logs e métricas registrarem sucesso e trata um formulário rejeitado como sucesso.

<a name="public-fields-are-an-allowlist"></a>
## Campos públicos formam uma lista explícita de permissão

Nunca serialize um model só porque seus campos atuais parecem inofensivos: uma coluna secreta futura entraria no payload sem mudança no endpoint. O mapa `JsonResource` torna adições deliberadas, mas não verifica chaves de string na compilação. Separe metadados With de data e não exponha chaves internas de tenant, digests, caminhos de storage ou detalhes de policy.

<a name="assert-status-shape-and-fragment"></a>
## Verifique status, formato e fragmento

Teste status exato e chaves JSON para a API, depois elemento raiz e mensagens de campo do formulário parcial. Verifique a ausência de campo sensível do model. O teste do recurso deve falhar se o status virar 200 indevidamente, se uma chave mudar ou se o fragmento deixar de apontar para o formulário substituído.

<a name="resources-sit-between-domain-and-transport"></a>
## Recursos ficam entre domínio e transporte

Validadores definem entrada aceita, models definem dados armazenados, views definem HTML e recursos definem saída. A separação é intencional mesmo com campos iniciais parecidos. Recursos de paginação carregam estado de página; tratamento de erros possui a representação estável do problema quando não há recurso de sucesso.
