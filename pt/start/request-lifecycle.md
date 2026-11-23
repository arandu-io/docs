# Ciclo de vida da requisição

- [Cada requisição percorre um pipeline explícito](#one-request-crosses-one-explicit-pipeline)
- [Identifique a rota e seus middlewares](#know-the-route-and-the-middleware-that-own-it)
- [Observe uma resposta completa](#observe-a-complete-response)
- [Teste pelo Handler quando a porta não importar](#test-through-handler-when-the-port-is-irrelevant)
- [Cada etapa tem uma responsabilidade](#each-stage-has-one-question)

<a name="one-request-crosses-one-explicit-pipeline"></a>
## Cada requisição percorre um pipeline explícito

O servidor Go aceita a conexão, aplica os middlewares globais, encontra a rota, aplica seus middlewares, cria um contexto HTTP, chama o handler registrado e escreve a resposta. Os testes podem inicializar a mesma aplicação composta usada pelo servidor sem abrir uma porta.

<a name="know-the-route-and-the-middleware-that-own-it"></a>
## Identifique a rota e seus middlewares

Use o inventário de rotas em vez de adivinhar a partir de uma URL. O framework compõe o pipeline global a partir da configuração, e `bootstrap/app.go` o ajusta em `WithMiddleware`. O middleware de rota é anexado onde a rota é declarada. Um controller recebe seus colaboradores pelo construtor; ele não os resolve durante a requisição.

```sh
aru route:list
aru dev
```

<a name="observe-a-complete-response"></a>
## Observe uma resposta completa

```sh
curl -i http://127.0.0.1:8080/
```

A resposta atravessou o pipeline global nesta ordem: as frases de recusa, Recover, o que Prepend acrescentou, os dois limites de corpo, `TrustProxies`, Observe, `SecurityHeaders`, o que `BeforeThrottle` acrescentou, o Throttle global, a verificação de CSRF, `OverrideMethod` e o que Append acrescentou. Depois vêm a correspondência de rota, o middleware da rota, o controller e o renderer. Numa rota que altera estado, autenticação, autorização, CSRF e validação continuam sendo verificações separadas, e não uma guarda genérica.

> [!NOTE]
> Uma requisição sob `/_arandu/` pula o pipeline da aplicação somente quando uma rota do framework a responde: uma sonda, o console de depuração, o reload ou um asset. Qualquer outro caminho sob esse prefixo é tráfego da aplicação e atravessa o pipeline inteiro. Create compõe a aplicação mas não abre nenhum listener; `HandleCommand` inicializa os módulos e executa o servidor para serve, o comando sem argumentos, enquanto um comando que só precisa da configuração e do banco de dados não inicializa nada.

<a name="test-through-handler-when-the-port-is-irrelevant"></a>
## Teste pelo Handler quando a porta não importar

```go example
func assertHomeIsOK(t *testing.T, app *foundation.Application) {
	request := httptest.NewRequest(stdhttp.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	app.Handler().ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}
}
```

- Um panic se torna a resposta de erro configurada e é registrado uma única vez.
- O cancelamento de uma requisição chega às chamadas de banco de dados e de clientes por meio do contexto.
- Os testes distinguem falhas de autenticação de recusas da Policy.
- A resposta é escrita uma única vez, com o status e o tipo de conteúdo esperados.

<a name="each-stage-has-one-question"></a>
## Cada etapa tem uma responsabilidade

O roteamento escolhe o handler. Os middlewares cuidam das políticas transversais da requisição. Os controllers traduzem o transporte em uma chamada de serviço. Os serviços autorizam os casos de uso, e os repositórios acessam os dados no escopo do tenant. Views e resources definem o formato da saída. Separar essas responsabilidades permite rastrear uma requisição que falhou.

- [Rotas](/docs/{{version}}/http/routing)
- [Middleware](/docs/{{version}}/http/middleware)
- [Controllers](/docs/{{version}}/http/controllers)
- [Contexto](/docs/{{version}}/architecture/context)
- [Aplicação do Framework](https://pkg.go.dev/github.com/arandu-io/framework/foundation)
