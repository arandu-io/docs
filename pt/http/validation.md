# Validação

- [Validação restringe entrada; policy concede autoridade](#validation-narrows-input-while-policy-grants-authority)
- [Declare a entrada junto do caso de uso](#declare-the-input-beside-its-use-case)
- [Rejeite uma criação antes de chegar ao armazenamento](#reject-one-create-request-before-storage)
- [Coloque cada regra na fronteira que ela conhece](#put-each-rule-at-the-boundary-it-understands)
- [Preserve erros de campos sem vazar causas](#preserve-field-errors-without-leaking-causes)
- [Defina campos permitidos e limite trabalho caro](#allowlist-fields-and-expensive-work)
- [Teste o menor valor inválido e o limite](#test-the-smallest-failing-value-and-the-boundary)
- [Continue de dados válidos para trabalho autorizado](#continue-from-valid-data-to-authorized-work)

<a name="validation-narrows-input-while-policy-grants-authority"></a>
## Validação restringe entrada; policy concede autoridade

Validação transforma valores não confiáveis do transporte em entrada tipada e erros indexados por campo. Verifica formato e restrições neutras de negócio; autorização e unicidade no banco permanecem operações separadas com sua própria semântica de falha.

<a name="declare-the-input-beside-its-use-case"></a>
## Declare a entrada junto do caso de uso

Mantenha explícitos nomes de transporte, campos de domínio e tamanhos máximos. Não adicione `TenantID` à entrada do navegador: o repositório deriva tenant do Grant.

<a name="reject-one-create-request-before-storage"></a>
## Rejeite uma criação antes de chegar ao armazenamento

Use Bind para ler a requisição num valor tipado cujas tags form nomeiam cada campo aceito e cujo método Rules declara o que cada campo pode ser. `ctx.Bind` roda as regras sobre a entrada bruta antes de converter qualquer coisa e devolve a falha como um `validation.Errors` não vazio, que o adapter de ação transforma em flash e redirecionamento para uma página e em um documento de problema 422 para um cliente JSON. Bind lê apenas o corpo em POST, PUT, PATCH e DELETE, então um valor na query string não satisfaz uma regra nem chega ao registro gravado. Normalize e-mail, identificadores e espaços antes de persistir valores aceitos.

```go example
type RegisterRequest struct {
	Email string `form:"email"`
}

func (RegisterRequest) Rules() validation.Rules {
	return validation.Rules{
		"email": "required|email|max:254",
	}
}

func Register(ctx *fhttp.Context) error {
	var req RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}
	return nil
}
```

<a name="put-each-rule-at-the-boundary-it-understands"></a>
## Coloque cada regra na fronteira que ela conhece

Declare em Rules as regras de obrigatoriedade, tipo, tamanho, pertencimento e comparação, como texto com barras verticais ou como lista por campo. O que cruza campos ou consulta os dados pertence a um método Validate da requisição, que o serviço chama; não repita em Validate uma restrição que Rules já declara. As mensagens saem no idioma da requisição, e os nomes dos campos ficam estáveis para posicionar cada erro junto à sua entrada.

- Valide antes de autorizar quando entrada malformada não deve chegar à policy.
- Autorize antes que uma consulta protegida revele a existência do registro.
- Retorne consistentemente JSON 422 ou flash com redirecionamento de formulário.
- Separe conflitos de persistência de erros de sintaxe.

<a name="preserve-field-errors-without-leaking-causes"></a>
## Preserve erros de campos sem vazar causas

Regras desconhecidas, argumentos inválidos, erros retornados vazios e coerção que esconde entrada malformada são defeitos de implementação. Indisponibilidade de banco nunca é mensagem de validação de campo.

> [!WARNING]
> Não consulte unicidade sem tenant atual e fronteira de autorização. Uma regra aparentemente global pode revelar dados de outro tenant ou rejeitar um valor válido no escopo atual.

<a name="allowlist-fields-and-expensive-work"></a>
## Defina campos permitidos e limite trabalho caro

> [!WARNING]
> Defina campos permitidos e limite bytes do corpo, itens de coleções, strings e uploads antes do trabalho caro. Nunca aceite seletor de tenant, caminho de storage, campo de model ou hash de senha da entrada; mantenha credenciais fora de mensagens, logs e dados antigos de formulário em flash.

<a name="test-the-smallest-failing-value-and-the-boundary"></a>
## Teste o menor valor inválido e o limite

Teste em tabela cada limite de regra com nome real do campo e chave de mensagem; envie depois uma requisição inválida por clientes normais e aprimorados. Confirme que entrada rejeitada não altera registros e que entrada válida chega normalizada ao serviço.

- Limites de obrigatoriedade e máximo são exatos.
- Confirmação entre campos usa nomes de campos aceitos.
- Respostas 422 continuam sendo respostas de rejeição.
- Nenhuma consulta protegida roda fora do escopo do tenant.

<a name="continue-from-valid-data-to-authorized-work"></a>
## Continue de dados válidos para trabalho autorizado

Parsing fornece valores candidatos, controllers escolhem o caso de uso, policies decidem autoridade e respostas renderizam erros dos campos sem mudar seu status.

- [Requisições](/docs/{{version}}/http/requests)
- [Controllers](/docs/{{version}}/http/controllers)
- [Respostas](/docs/{{version}}/http/responses)
