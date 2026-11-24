# Controllers

- [Traduza HTTP e mantenha a policy no serviço seguinte](#translate-http-and-keep-policy-downstream)
- [Injete o caso de uso](#inject-the-use-case)
- [Renderize um registro autenticado](#render-one-authenticated-record)
- [Mantenha ações curtas e explícitas](#keep-actions-short-and-explicit)
- [Preserve o significado da falha](#preserve-failure-meaning)
- [O controller não é autoridade de autorização](#the-controller-is-not-an-authorization-authority)
- [Teste transporte e policy separadamente](#test-transport-and-policy-separately)
- [Complete a fronteira do controller](#complete-the-controller-boundary)

<a name="translate-http-and-keep-policy-downstream"></a>
## Traduza HTTP e mantenha a policy no serviço seguinte

Um controller lê entrada do transporte, solicita resultado ao serviço da aplicação e escolhe uma resposta HTTP. Não emite Grants, abre outro caminho de banco nem confia em campos de tenant e propriedade do navegador.

<a name="inject-the-use-case"></a>
## Injete o caso de uso

Defina a dependência no pacote consumidor e injete pelo construtor do controller. Proteja a rota com middleware de autenticação quando exigir sujeito: `RequireAuth` coloca no contexto da requisição o sujeito que carregou da sessão, então `ctx.User` o informa atrás da guarda, e `middleware.LoadSubject` faz o mesmo numa rota pública, deixando passar um visitante sem sujeito. Serviço ou policy continuam responsáveis por autorizar a ação no registro solicitado.

- Construa os controllers num provider em app/Providers.
- Registre ações em rotas nomeadas.
- Separe entrada de validação de models persistidos.
- Deixe autorização e persistência com o serviço.

<a name="render-one-authenticated-record"></a>
## Renderize um registro autenticado

O prazo do contexto e o Subject autenticado atravessam a fronteira do serviço; o identificador do caminho continua sendo entrada não confiável, não prova de acesso. Renderize o resultado já autorizado por view ou recurso declarado.

```go example
type Invoice struct{ ID string }

type InvoiceFinder interface {
	Find(context.Context, security.Subject, string) (Invoice, error)
}

type InvoiceController struct{ invoices InvoiceFinder }
type InvoiceShowData struct{ Invoice Invoice }

func (c InvoiceController) Show(ctx *fhttp.Context) error {
	subject, ok := ctx.User()
	if !ok {
		hhttp.Refuse(ctx.Response, ctx.Request, stdhttp.StatusUnauthorized, "sign in required")
		return nil
	}
	invoice, err := c.invoices.Find(ctx.Ctx(), subject, ctx.Param("invoice"))
	if err != nil {
		return err
	}
	return ctx.View("invoices.show", InvoiceShowData{Invoice: invoice})
}
```

<a name="keep-actions-short-and-explicit"></a>
## Mantenha ações curtas e explícitas

Use Param, Query e Input apenas para montar entrada tipada. Retorne erros operacionais à recuperação central, `validation.Errors` não vazio ao caminho de rejeição e exatamente um entre View, Fragment, JSON, `RedirectRoute` ou Status para responder.

- Passe `ctx.Ctx` ao caso de uso.
- Leia identidade pela indicação de presença do contexto.
- Retorne apenas erros de validação não vazios.
- Escolha um helper de resposta após o resultado do serviço.

<a name="preserve-failure-meaning"></a>
## Preserve o significado da falha

Entrada malformada, registros ausentes, ações negadas, conflitos e falhas internas são resultados distintos. Converter tudo em 500 apaga comportamento útil; engolir um erro pode retornar 200 com corpo vazio por acidente.

> [!WARNING]
> O adapter de ação causa panic quando `validation.Errors` retornado está vazio. Use if `errs.Any()` antes de retorná-lo e retorne nil quando a validação passar.

<a name="the-controller-is-not-an-authorization-authority"></a>
## O controller não é autoridade de autorização

> [!WARNING]
> Nunca derive tenant de formulário, query ou caminho; nunca carregue dados de tenant sem Subject ou Grant; nunca serialize integralmente um model de persistência com hashes, chaves internas de tenant ou campos de auditoria.

Corpos de requisição e structs que podem conter credenciais exigem remoção explícita em JSON e slog antes de serem registrados. Conveniência do controller não autoriza expor entrada do transporte.

<a name="test-transport-and-policy-separately"></a>
## Teste transporte e policy separadamente

Use httptest pela rota registrada para verificar status, corpo ou redirecionamento, parâmetros, recusa de autenticação e erro de serviço. Num teste unitário focado, injete um finder fake e confirme que recebe o mesmo contexto, Subject e ID.

- Um visitante é recusado antes do caso de uso.
- Subject e ID exatos chegam ao serviço.
- Validação segue o contrato de 422 ou redirecionamento.
- O comportamento da policy continua testado fora do controller.

<a name="complete-the-controller-boundary"></a>
## Complete a fronteira do controller

Tipos de requisição possuem entrada aceita, respostas possuem semântica de transporte, recursos possuem campos JSON públicos e autorização possui Grant. Um controller curto deixa esses contratos visíveis.

- [Requisições](/docs/{{version}}/http/requests)
- [Respostas](/docs/{{version}}/http/responses)
- [Autorização](/docs/{{version}}/security/authorization)
- [Recursos](/docs/{{version}}/http/resources)
