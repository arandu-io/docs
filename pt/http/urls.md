# URLs e URLs assinadas

- [Nomes de rota são a fonte dos endereços](#route-names-are-the-address-source)
- [Gere e redirecione por nome](#build-and-redirect-by-name)
- [Separe caminhos, hosts e assinaturas](#keep-paths-hosts-and-signatures-separate)
- [Não esconda a lacuna da integração do framework](#do-not-hide-the-framework-bridge-gap)
- [Altere cada byte assinado](#change-every-signed-byte)
- [Superfícies de URL publicadas e sua fronteira](#released-url-surfaces-and-boundary)
- [Use links conforme seu nível real de confiança](#use-links-at-their-real-trust-level)

<a name="route-names-are-the-address-source"></a>
## Nomes de rota são a fonte dos endereços

Rotas nomeadas são a fonte dos endereços da aplicação. Gere pela tabela em vez de concatenar segmentos; assine apenas caminho canônico e valores limitados de query quando o link precisar detectar adulteração.

<a name="build-and-redirect-by-name"></a>
## Gere e redirecione por nome

Gere primeiro o caminho nomeado e falhe em nome desconhecido ou parâmetros incompatíveis. Um assinador pode adicionar expiração e assinatura a um link de propósito único; verifique antes de o handler usar qualquer valor assinado da query.

```go example
func SignInvoiceLink(router *fhttp.Router, signer *encryption.Signer, invoiceID string) (string, error) {
	path, err := router.Table().URL("invoices.show", invoiceID)
	if err != nil {
		return "", err
	}

	token := signer.Sign("invoice.show", path, 15*time.Minute)
	verifiedPath, err := signer.Verify("invoice.show", token)
	if err != nil {
		return "", err
	}
	return verifiedPath, nil
}
```

<a name="keep-paths-hosts-and-signatures-separate"></a>
## Separe caminhos, hosts e assinaturas

Mantenha decisões de host e protocolo na configuração confiável da implantação. Gere URLs absolutas só após normalizar o proxy confiável, preserve a ordem canônica dos parâmetros na assinatura e use expiração curta em ações por e-mail.

- Nomeie rotas durante o registro.
- Trate erros de geração reversa de URL no código da aplicação.
- Assine o destino canônico, não texto fornecido pelo chamador.
- Verifique expiração e assinatura antes de qualquer mutação.

<a name="do-not-hide-the-framework-bridge-gap"></a>
## Não esconda a lacuna da integração do framework

Nome de rota desconhecido, contagem incorreta de parâmetros, query alterada, assinatura expirada e host não confiável são recusas distintas. Nunca use destino sem assinatura como fallback quando assinar falhar. `hesape/routing.SignedRoute` recebe a tabela de rotas de `hesape/routing`, não a de `framework/http.Router`; com o router do framework, obtenha o caminho por `Table().URL` e assine-o com `security.Signer`, como no exemplo acima.

> [!WARNING]
> Assinatura válida comprova integridade e expiração, não identidade ou permissão. O handler ainda carrega Subject atual, autoriza a ação e garante uso único no armazenamento durável.

<a name="change-every-signed-byte"></a>
## Altere cada byte assinado

Gere uma URL nomeada, assine, verifique e altere um segmento do caminho, valor de query e expiração. Teste URLs absolutas na topologia real de proxies confiáveis e impeça que Host controlado pelo atacante vire origem canônica.

- Parâmetros de rota são escapados e contados.
- Um valor alterado falha na verificação.
- Um link expirado falha na verificação.
- A autorização ainda executa após verificar a assinatura.

<a name="released-url-surfaces-and-boundary"></a>
## Superfícies de URL publicadas e sua fronteira

`Route.Name` registra identidade estável e `Routes.URL` gera seu caminho. O assinador publicado de security cria e verifica destinos assinados limitados. Essas APIs não fazem busca de conta, autorização nem consumo de uso único.

<a name="use-links-at-their-real-trust-level"></a>
## Use links conforme seu nível real de confiança

Roteamento possui o modelo do caminho, proxies confiáveis estabelecem protocolo e host externos, autorização possui a ação e códigos descartáveis possuem o uso único durável.

- [Rotas](/docs/{{version}}/http/routing)
- [Autorização](/docs/{{version}}/security/authorization)
- [Verificação de e-mail](/docs/{{version}}/security/email-verification)
