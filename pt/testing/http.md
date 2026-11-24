# Testes HTTP

- [Despache o pipeline registrado com httptest](#dispatch-the-registered-pipeline-with-httptest)
- [Inicialize a composição real e controle os adapters externos](#boot-the-real-graph-and-keep-external-adapters-controlled)
- [Envie uma requisição pelo Handler](#send-one-request-through-handler)
- [Verifique protocolo, sessão, YPY e caminhos de policy](#assert-protocol-session-ypy-and-policy-paths)
- [Não deixe a recuperação de produção esconder um defeito do teste](#do-not-let-production-recovery-hide-a-test-defect)
- [Mantenha credenciais e dados de tenants fora de fixtures e saídas](#keep-credentials-and-tenant-data-out-of-fixtures-and-output)
- [Cubra a matriz completa de rotas](#cover-the-complete-route-matrix)
- [Leve detalhes de serviços e adapters a contratos focados](#move-service-and-adapter-details-to-focused-contracts)

<a name="dispatch-the-registered-pipeline-with-httptest"></a>
## Despache o pipeline registrado com httptest

Um teste de feature HTTP despacha `foundation.Application.Handler` por `net/http/httptest`. Ele atravessa middleware global, correspondência de rotas, middleware de rota, controller, serviço, renderizador e resposta de erro sem abrir um listener.

<a name="boot-the-real-graph-and-keep-external-adapters-controlled"></a>
## Inicialize a composição real e controle os adapters externos

Construa os mesmos módulos e configurações usados pelo servidor, substitua apenas adapters externos explícitos, compile views e use banco e backend de sessão dedicados. Preserve a ordem dos middlewares de produção.

- Registre rotas e renderizador reais.
- Use um banco de teste com fixtures de tenants.
- Capture e-mail e HTTP de saída.
- Escolha sessão em memória ou compartilhada conforme o comportamento testado.

<a name="send-one-request-through-handler"></a>
## Envie uma requisição pelo Handler

Crie uma requisição padrão, despache-a ao Handler da aplicação e verifique status, cabeçalhos e corpo. Quando a jornada atravessa login, CSRF, rotação e logout, `test.NewClient(t, app.Handler())`, de `framework/test`, é o mesmo despacho com um pote de cookies, o token CSRF lido da última página carregada e sem seguir redirecionamentos. Seu `*test.Response` verifica com `AssertStatus`, `AssertOk`, `AssertSee`, `AssertDontSee`, `AssertRedirect` e os métodos Assert para cabeçalhos, cookies e JSON. `test.Serve` sobe a aplicação numa porta livre sobre um banco só do teste, e `test.As` entra como um sujeito pela sessão da própria aplicação.

```go example
func assertHomeIsOK(t *testing.T, app *foundation.Application) {
	request := httptest.NewRequest(stdhttp.MethodGet, "/invoices/inv_123", nil)
	response := httptest.NewRecorder()
	app.Handler().ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}
}
```

<a name="assert-protocol-session-ypy-and-policy-paths"></a>
## Verifique protocolo, sessão, YPY e caminhos de policy

Teste representação normal e YPY-Request, redirecionamentos nomeados, recursos JSON, validação de formulários, CSRF, proxy confiável, autenticação, recusa de policy e cancelamento. Verifique o resultado do serviço ou da persistência além do texto da resposta.

- Verifique Content-Type e o status exato.
- Mantenha cookies ao longo do fluxo real.
- Teste página completa e fragmento separadamente.
- Inspecione os efeitos após a resposta.

<a name="do-not-let-production-recovery-hide-a-test-defect"></a>
## Não deixe a recuperação de produção esconder um defeito do teste

Um panic deve chegar à fronteira do teste quando configurado para expor defeitos; não aprove um teste apenas porque a recuperação de produção converteu o panic em 500. Escritas parciais, erros engolidos, rotas desconhecidas e views geradas desatualizadas exigem casos explícitos. Um fake de controller em teste unitário não prova middleware, nomes de rota, cookies de sessão, CSRF, registro do renderer nem mapeamento de erros; mantenha ao menos uma jornada completa de rota para cada fluxo público.

<a name="keep-credentials-and-tenant-data-out-of-fixtures-and-output"></a>
## Mantenha credenciais e dados de tenants fora de fixtures e saídas

> [!WARNING]
> Use credenciais e tenants sintéticos, oculte dados sensíveis na saída de falha, nunca chame provedores reais e verifique que um tenant alheio recebe a mesma ausência limitada. Teste Secure, `HttpOnly`, `SameSite`, assinatura e atributos CSRF como detalhes do protocolo.

Um helper que pula middleware pode facilitar testes e deixar a produção errada; autentique pela sessão publicada ou construa deliberadamente um Subject apenas em testes de serviço.

<a name="cover-the-complete-route-matrix"></a>
## Cubra a matriz completa de rotas

Cubra método e caminho incorretos, requisição válida, entrada malformada, visitante, Subject recusado, outro tenant, YPY, preferência JSON, dependência indisponível, panic, cancelamento da requisição e rotação de sessão. Execute a suíte com detecção de race.

- Toda rota possui caminhos de sucesso e recusa.
- Registros de outro tenant permanecem invisíveis.
- Requisições aprimoradas mantêm o status correto.
- O cancelamento chega às chamadas seguintes.

<a name="move-service-and-adapter-details-to-focused-contracts"></a>
## Leve detalhes de serviços e adapters a contratos focados

Controllers traduzem o transporte, o ciclo da requisição explica o pipeline, dublês capturam adapters externos e testes de banco comprovam isolamento no armazenamento.

- [Controllers](/docs/{{version}}/http/controllers)
- [Ciclo de vida da requisição](/docs/{{version}}/start/request-lifecycle)
- [Dublês de teste](/docs/{{version}}/testing/test-doubles)
- [Testes de banco de dados](/docs/{{version}}/testing/database)
