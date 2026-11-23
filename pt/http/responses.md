# Respostas HTTP

- [Confirme um formato de resposta uma única vez](#commit-one-response-shape-once)
- [Retorne um recurso público, não um model](#return-a-public-resource-not-a-model)
- [Escolha o helper conforme o contrato do cliente](#match-the-helper-to-the-client-contract)
- [Não escreva depois de confirmar a resposta](#do-not-write-past-a-committed-response)
- [Verifique detalhes do protocolo, não só corpos](#assert-protocol-details-not-only-bodies)
- [Helpers de resposta publicados](#released-response-helpers)
- [Mantenha representação na borda](#keep-representation-at-the-edge)

<a name="commit-one-response-shape-once"></a>
## Confirme um formato de resposta uma única vez

Uma resposta é um status explícito, conjunto de headers e representação do corpo. Helpers de Context preservam essa decisão para páginas completas, fragmentos, recursos JSON declarados, redirecionamentos nomeados e status vazios, impedindo escrita duplicada pelo handler.

<a name="return-a-public-resource-not-a-model"></a>
## Retorne um recurso público, não um model

Retorne `JsonResource` declarado para JSON, não model bruto. `ToArray` define os campos públicos permitidos e With fornece metadados superiores; `Context.JSON` escolhe status e codificação preservando a fronteira do recurso.

```go example
type Invoice struct {
	ID    string
	Total int64
}

type InvoiceResource struct{ invoice Invoice }

func (r InvoiceResource) ToArray() map[string]any {
	return map[string]any{"id": r.invoice.ID, "total": r.invoice.Total}
}
func (InvoiceResource) With() map[string]any { return nil }

func Show(ctx *fhttp.Context, invoice Invoice) error {
	return ctx.JSON(stdhttp.StatusOK, InvoiceResource{invoice: invoice})
}
```

<a name="match-the-helper-to-the-client-contract"></a>
## Escolha o helper conforme o contrato do cliente

Use View para um documento completo, Fragment para uma parcial nomeada e `RedirectRoute` para um destino da tabela de rotas. `RedirectRoute` retorna 303 a um navegador comum e YPY-Redirect com 204 a uma requisição aprimorada, para que ambos naveguem corretamente. Use `EventStream` para atualizações unidirecionais ao vivo do servidor; `EventStreamWithOptions` adiciona comentários de heartbeat enquanto cada evento é descarregado imediatamente.

- Escolha o status antes de escrever o corpo.
- Defina headers da resposta antes da primeira escrita.
- Retorne fragmentos 422 como o conteúdo do formulário que rejeitam, que o YPY renderiza no lugar.
- Separe recursos JSON de models de banco.
- Para SSE, crie eventos com id, event, data, retry ou Comment e pare a produção assim que o contexto da requisição for cancelado.

<a name="do-not-write-past-a-committed-response"></a>
## Não escreva depois de confirmar a resposta

Erros de codificação, nomes de view desconhecidos, nomes de rota ausentes e escritas após uma resposta terminal não são sucessos vazios. Preserve os significados de 404, 409, 422 e 503 em vez de reduzir toda falha a 200 ou 500. Um stream SSE não é reprodução durável: o cancelamento da desconexão para o produtor e omite deliberadamente um evento de encerramento. Se um stream longo precisar ultrapassar o prazo comum de escrita do servidor, limpe esse prazo para esta resposta com `http.NewResponseController` antes do streaming.

> [!WARNING]
> O YPY renderiza um 422 dentro do formulário que enviou a requisição e decide todo outro 4xx e 5xx pelas suas regras de erro. Mudar um 422 para 200 faria métricas e logs tratarem entrada rejeitada como sucesso.

<a name="assert-protocol-details-not-only-bodies"></a>
## Verifique detalhes do protocolo, não só corpos

Verifique status, tipo de conteúdo, cabeçalhos e corpo para requisições comuns e aprimoradas. Teste `RedirectRoute` nos dois modos, chaves JSON exatas, ausência de campos sensíveis, raiz de fragmento 422 e erro de codificação ou renderização. Para SSE, verifique text/event-stream, no-cache, X-Accel-Buffering: no, enquadramento de dados multilinha, validação de metadados, comentários de heartbeat, um flush por evento e saída rápida no cancelamento e desligamento da aplicação.

- Redirecionamento normal é 303 com Location.
- Redirecionamento aprimorado é 204 com YPY-Redirect.
- A saída `JsonResource` contém apenas chaves permitidas.
- Fragmentos de validação preservam 422 e são o conteúdo do formulário.
- O cancelamento de SSE para a produção sem criar um evento final sintético.

<a name="released-response-helpers"></a>
## Helpers de resposta publicados

`Context.View` e Fragment renderizam HTML, JSON aceita `hhttp.JsonResource`, `RedirectRoute` reverte uma rota nomeada e Status escreve um status terminal. `hhttp.Refuse` fornece erros públicos limitados sem expor causas internas. `ResponseFactory.EventStream` e `EventStreamWithOptions` constroem respostas SSE; `hhttp.NewServerSentEvent` e Comment codificam metadados e dados do evento sem permitir injeção de nova linha ou NUL.

<a name="keep-representation-at-the-edge"></a>
## Mantenha representação na borda

Respostas REST encerram uma requisição. SSE mantém uma resposta HTTP aberta para atualizações do servidor ao cliente, mas não aceita mensagens do cliente; use `WebSocket` para sessões bidirecionais e armazenamento durável quando clientes precisarem reproduzir fatos perdidos. Resources definem campos JSON, views definem HTML, controllers escolhem a representação e o tratamento central de erros cuida das falhas que escapam da action.

- [Padrões de comunicação](/docs/{{version}}/architecture/contracts-adapters)
- [`WebSocket` do Joaju](/docs/{{version}}/async/joaju)
- [Recursos](/docs/{{version}}/http/resources)
- [Erros](/docs/{{version}}/http/errors)
