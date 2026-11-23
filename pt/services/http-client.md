# Cliente HTTP

- [Toda requisição externa começa com contexto e prazo](#every-outbound-request-begins-with-context-and-a-deadline)
- [Envie uma requisição JSON limitada](#send-one-bounded-json-request)
- [Adicione headers, tentativas e decodificação deliberadamente](#add-headers-retries-and-decoding-deliberately)
- [Distinga recusa, timeout e resposta malformada](#distinguish-refusal-timeout-and-malformed-response)
- [Teste contra um servidor local do protocolo](#test-against-a-local-protocol-server)
- [Superfície publicada do cliente HTTP](#released-http-client-surface)
- [Instrumente a chamada e mantenha a política da aplicação fora dela](#instrument-the-call-and-keep-application-policy-outside-it)

> [!WARNING]
> **Limitações atuais.** O cliente propaga contexto e restringe tentativas a requisições seguras, mas nenhum fluxo atual de primeira parte exercita todos os verbos. `Batch.Send` também não tem cobertura comportamental completa; prefira requisição única até verificar essa fronteira.

<a name="every-outbound-request-begins-with-context-and-a-deadline"></a>
## Toda requisição externa começa com contexto e prazo

O cliente HTTP do Hesape envolve `net/http` com verbos conscientes de contexto, opções de requisição, decodificação e middleware. O componente está publicado, mas não há consumidor atual de primeira parte em tag; este guia se limita ao contrato verificado do pacote, sem prometer composição padrão no esqueleto.

<a name="send-one-bounded-json-request"></a>
## Envie uma requisição JSON limitada

Construa Factory no boot, defina teto de bytes da resposta, derive um prazo filho e envie JSON tipado por `PendingRequest`. Verifique o transporte e chame Throw para 4xx ou 5xx antes de decodificar a resposta já carregada.

```go example
func PostPayload(ctx context.Context, endpoint string, payload map[string]any) error {
	factory := client.NewFactory(nil).MaxResponseBytes(2 << 20)
	request := factory.CreatePendingRequest().Timeout(5 * time.Second).AcceptJSON()
	callCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	response, err := request.Post(callCtx, endpoint, payload)
	if err != nil {
		return err
	}
	if err := response.Throw(nil); err != nil {
		return err
	}
	return nil
}
```

<a name="add-headers-retries-and-decoding-deliberately"></a>
## Adicione headers, tentativas e decodificação deliberadamente

Configure URL base e headers imutáveis na construção, adicione autorização específica sem registrá-la, propague contexto, limite corpos de resposta e repita apenas operações idempotentes ou com chave explícita de idempotência, usando backoff limitado.

- Use contexto em todo verbo.
- Defina prazos para cliente e operação.
- Limite os bytes da resposta antes de decodificar.
- Instrumente pelo Collector da aplicação.

<a name="distinguish-refusal-timeout-and-malformed-response"></a>
## Distinga recusa, timeout e resposta malformada

DNS, conexão, TLS, cancelamento, prazo, status não 2xx e corpo inválido têm evidências distintas. Nunca transforme dependência indisponível em sucesso vazio nem repita mutação insegura apenas porque o transporte fechou.

> [!WARNING]
> Timeout zero ou cliente padrão sem configuração podem reter recursos indefinidamente. Oculte URLs com credenciais e limite a captura do corpo de erro para impedir que um serviço externo encha logs ou memória.

<a name="test-against-a-local-protocol-server"></a>
## Teste contra um servidor local do protocolo

Use `httptest.Server` para método, cancelamento, headers, corpo JSON, resposta grande limitada, redirecionamentos, resposta lenta, status repetível, mutação insegura, JSON malformado e confiança TLS. Verifique que credenciais não chegam aos logs.

- O cancelamento chega ao servidor.
- Novas tentativas seguem a política de idempotência.
- Limites da resposta atuam antes de alocação sem teto.
- Observabilidade registra URL, status e duração sem segredos.

<a name="released-http-client-surface"></a>
## Superfície publicada do cliente HTTP

`NewFactory`, `Factory.MaxResponseBytes`, `Factory.CreatePendingRequest`, verbos de `PendingRequest`, `Response.Throw` e decodificação compõem a API publicada. Response lê e fecha o corpo imediatamente; o chamador não fecha um campo Body exportado.

| Superfície | Contrato publicado |
| --- | --- |
| `NewFactory` / `MaxResponseBytes` | constrói a factory protegida |
| `CreatePendingRequest` / Timeout | descreve uma requisição |
| `Get/Post/Put/Patch/Delete` | verbos HTTP conscientes do contexto |
| `Response.Throw` / Object | status e decodificação tipada |

<a name="instrument-the-call-and-keep-application-policy-outside-it"></a>
## Instrumente a chamada e mantenha a política da aplicação fora dela

Contexto fornece cancelamento, logs registram esperas externas, contratos ocultam detalhes de provedores dos serviços de domínio e testes possuem o fixture do protocolo.

- [Contexto](/docs/{{version}}/architecture/context)
- [Logs](/docs/{{version}}/operations/logging-observability)
- [Contratos](/docs/{{version}}/architecture/contracts-adapters)
- [Dublês de teste](/docs/{{version}}/testing/test-doubles)
