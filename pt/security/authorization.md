# Policies e Grant

- [Uma policy transforma identidade em autoridade limitada](#a-policy-turns-identity-into-bounded-authority)
- [Leve Subject confiável ao caso de uso](#carry-a-trusted-subject-to-the-use-case)
- [Autorize a operação e passe seu Grant](#authorize-the-operation-and-pass-its-grant)
- [Mantenha a decisão junto da ação protegida](#keep-the-decision-beside-the-protected-action)
- [Recuse sem revelar outro tenant](#refuse-without-revealing-another-tenant)
- [Nunca reconstrua autoridade pela entrada da requisição](#never-reconstruct-authority-from-request-input)
- [Comprove recusa e isolamento de tenant no armazenamento](#prove-denial-and-tenant-isolation-at-storage)
- [Identidade inicia a fronteira; Grant a completa](#identity-starts-the-boundary-and-grant-finishes-it)

<a name="a-policy-turns-identity-into-bounded-authority"></a>
## Uma policy transforma identidade em autoridade limitada

Autenticação produz `security.Subject`: tenant e identidade do ator confiáveis. Autorização pergunta a uma Policy tipada se esse sujeito pode executar uma operação em um candidato e retorna `security.Grant` para o repositório protegido. Uma string de papel não é autoridade.

<a name="carry-a-trusted-subject-to-the-use-case"></a>
## Leve Subject confiável ao caso de uso

Carregue Subject por middleware de sessão ou token e construa policies com os colaboradores necessários. Não carregue um registro de tenant por repositório sem escopo apenas para decidir se o chamador pode vê-lo.

- Use Subject fornecido pela fronteira de autenticação.
- Defina operação e tipo candidato explicitamente.
- Mantenha a consulta por tenant na policy ou repositório autorizado.
- Defina deliberadamente resultados de visitante e recusa.

<a name="authorize-the-operation-and-pass-its-grant"></a>
## Autorize a operação e passe seu Grant

Chame `security.Authorize` com policy, sujeito, ação e candidato. Passe apenas o Grant retornado ao método de repositório que o exige, tornando impossível omitir autorização na chamada. As alterações são um mapa literal das colunas que esta ação pode escrever, lidas de uma requisição tipada, de modo que nenhuma chave enviada chega à atualização. O adapter ainda deve derivar o escopo por `data.Tenant(grant)`.

```go example
const InvoiceUpdate security.Action = "invoice.update"

type Invoice struct{ ID string }

type UpdateInvoiceRequest struct{ Reference, DueOn string }

type Invoices interface {
	Update(ctx context.Context, grant security.Grant, id string, fields map[string]any) error
}

func Update(ctx context.Context, policy security.Policy[Invoice], subject security.Subject, invoices Invoices, invoice Invoice, req UpdateInvoiceRequest) error {
	grant, err := security.Authorize(ctx, policy, subject, InvoiceUpdate, invoice)
	if err != nil {
		return err
	}

	return invoices.Update(ctx, grant, invoice.ID, map[string]any{
		"reference": req.Reference,
		"due_on":    req.DueOn,
	})
}
```

<a name="keep-the-decision-beside-the-protected-action"></a>
## Mantenha a decisão junto da ação protegida

Use uma policy pequena por recurso de domínio, mapeie operações para ações nomeadas estáveis e mantenha o transporte do controller separado. Nunca serialize o Grant da requisição para trabalho em segundo plano. Um job de fila, uma tarefa agendada e uma entrega de webhook rodam sob um Grant que o framework reconstrói a partir do tenant e da ação registrados, e do operador quando alguém agia na conta. Um listener do outbox recebe o evento armazenado, com `TenantID`, `Action` e `AuthorizedBy`, e monta o próprio Grant com `auth.SystemGrant(<action>, stored.TenantID)` numa linha marcada com `//arandu:system-grant <motivo>`, que o `aru doctor` verifica. Um Grant emitido enquanto um operador visualiza a conta de outra pessoa é somente leitura: `Grant.ReadOnly` responde true, e toda escrita em repositório, cache, disco e fila o recusa com `auth.ErrReadOnly`, então deixe a ação fora de uma página assim em vez de deixá-la falhar.

- Autorize antes de leituras e escritas protegidas.
- Passe Grant pelas fronteiras de serviço e repositório.
- Retorne a mesma ausência pública para registros ocultos de outro tenant.
- Reconstrua Subject confiável para trabalho assíncrono.

<a name="refuse-without-revealing-another-tenant"></a>
## Recuse sem revelar outro tenant

Subject ausente é não autenticado; recusa de policy é proibida ou ocultada como não encontrado; indisponibilidade de repositório é operacional. Não reduza isso a registro vazio nem continue com Grant zero.

> [!WARNING]
> Grant vazio ou de outro tenant não pode alcançar registros do tenant. Um adapter que ignora Grant viola o contrato mesmo se testes de controller passarem.

<a name="never-reconstruct-authority-from-request-input"></a>
## Nunca reconstrua autoridade pela entrada da requisição

> [!WARNING]
> Nunca aceite tenant, papel, Subject ou Grant de formulário, caminho, query, claim de token sem verificação, payload de job ou argumentos MCP. Nunca armazene autorização em cache além do sujeito, ação, recurso e estado de policy que ela cobriu.

Autenticação e URLs assinadas comprovam fatos diferentes. Sessão ou assinatura válida ainda exigem policy para a ação solicitada.

<a name="prove-denial-and-tenant-isolation-at-storage"></a>
## Comprove recusa e isolamento de tenant no armazenamento

Execute contrato compartilhado de repositório com Grants permitidos, negados, vazios e estrangeiros. Exercite cada ação da policy, visitante, identificadores de outro tenant, cancelamento e job que reconstrói identidade pelo caminho confiável da aplicação.

- Subject autorizado recebe Grant limitado.
- Subject recusado não alcança operação protegida do adapter.
- IDs de outro tenant não revelam existência.
- Todo repositório concreto impõe Grant.

<a name="identity-starts-the-boundary-and-grant-finishes-it"></a>
## Identidade inicia a fronteira; Grant a completa

Autenticação fornece Subject, sessões transportam identidade do navegador, contratos mantêm Grant nas assinaturas e controllers só traduzem transporte.

- [Autenticação](/docs/{{version}}/security/authentication)
- [Sessão](/docs/{{version}}/http/session)
- [Contratos e adapters](/docs/{{version}}/architecture/contracts-adapters)
- [Controllers](/docs/{{version}}/http/controllers)
