# Middleware

- [Um wrapper, uma decisão](#one-wrapper-one-decision)
- [Instale o perímetro de produção](#install-the-production-envelope)
- [Escolha o menor escopo de middleware](#choose-the-narrowest-middleware-scope)
- [Middleware precisa preservar o protocolo](#middleware-must-preserve-the-protocol)
- [Comprove ordem e ramos terminais](#prove-order-and-terminal-branches)
- [API de middleware publicada](#released-middleware-api)
- [Posicione a policy na fronteira correta](#put-policy-at-the-right-boundary)

<a name="one-wrapper-one-decision"></a>
## Um wrapper, uma decisão

O middleware HTTP tem a forma padrão `func(http.Handler) http.Handler`. O framework compõe o pipeline global a partir da configuração e o executa numa única ordem fixa na ida e na ordem inversa na volta; o middleware de rota envolve apenas a rota ou o grupo que protege.

<a name="install-the-production-envelope"></a>
## Instale o perímetro de produção

O envelope de produção já está instalado. Create compõe, a partir da configuração, a recuperação, os limites de corpo, a normalização de proxies confiáveis, a observabilidade, os cabeçalhos de segurança, um throttle global de trezentas requisições por minuto por sessão ou endereço, a verificação de CSRF e a sobrescrita de método. `bootstrap/app.go` ajusta esse pipeline em `WithMiddleware`: acrescenta middleware em um de três lugares nomeados e muda as poucas decisões que o pipeline deixa para a aplicação.

```go example
// Stand-ins: each wraps next with the behavior its name describes.
func requestTimer(next stdhttp.Handler) stdhttp.Handler        { return next }
func noStoreOnBearerLinks(next stdhttp.Handler) stdhttp.Handler { return next }
func tenantBanner(next stdhttp.Handler) stdhttp.Handler         { return next }

func WithPipeline(drivers *foundation.Drivers) *foundation.ApplicationBuilder {
	return foundation.Configure(drivers).
		WithMiddleware(func(m *foundation.Middleware) {
			m.Prepend(requestTimer)
			m.BeforeThrottle(noStoreOnBearerLinks)
			m.Append(tenantBanner)
			m.Throttle(cache.PerMinute(600))
			m.ValidateCSRFTokensExcept("/webhooks/", "/checkout/{token}/sessions")
			m.ImageOrigins("https://images.example.com")
		})
}
```

> [!NOTE]
> `TrustProxies` lê X-Forwarded-For só de um par que `TRUSTED_PROXIES` lista, percorrendo o cabeçalho da direita e pulando saltos confiáveis, de modo que um valor que o cliente escreveu à esquerda nunca vira o seu endereço. Tudo o que deriva identidade ou limite do endereço do cliente roda depois dele: `BeforeThrottle`, o Throttle global e Append.

<a name="choose-the-narrowest-middleware-scope"></a>
## Escolha o menor escopo de middleware

A ordem é fixa: as Sentences de recusa, Recover, Prepend, `ValidatePostSize` e `LimitBodySize`, `TrustProxies`, Observe, `SecurityHeaders`, `BeforeThrottle`, Throttle, a verificação de CSRF, `OverrideMethod` e Append, logo antes do router. Anexe autenticação e controles específicos de rota a um `Router.Group`, `Router.Action` ou método de verbo. `http.Chain` compõe um handler padrão avulso com o primeiro middleware no lado mais externo.

- Prepend acrescenta middleware logo depois de Recover, antes dos limites de corpo.
- `BeforeThrottle` vale também numa requisição recusada e lê o endereço real.
- Append roda depois de `OverrideMethod`, logo antes do router.
- Throttle substitui o orçamento global; Sentences e Observer substituem no mesmo lugar as palavras de recusa e o Observe.
- Autenticação pertence apenas a grupos de rotas protegidos.

<a name="middleware-must-preserve-the-protocol"></a>
## Middleware precisa preservar o protocolo

Um middleware chama `next.ServeHTTP` no máximo uma vez, retorna depois de uma resposta terminal, define cabeçalhos antes da primeira escrita e preserva `r.WithContext`. Quebrar essas regras pode duplicar trabalho, descartar cancelamento e identidade ou emitir uma resposta com os cabeçalhos errados. Recover fica no lado mais externo depois das frases: um panic num middleware acima dele escapa sem página.

> [!WARNING]
> Throttle publicado causa panic por construção inválida e falha aberto se o armazenamento falhar, registrando o erro. Não o apresente como proteção fail-closed de login ou recuperação sem policy da aplicação que recuse a indisponibilidade.

<a name="prove-order-and-terminal-branches"></a>
## Comprove ordem e ramos terminais

Despache pela aplicação composta e verifique onde um middleware acrescentado roda em relação ao throttle e à verificação de CSRF. Prove que um par não confiável não consegue forjar um endereço encaminhado, que um par confiável consegue, que uma recusa terminal nunca alcança o handler, que a 301ª requisição em um minuto é recusada e que um panic escapado vira o status de produção.

- A ordem de entrada e saída é determinística.
- Middleware terminal não chama o handler seguinte.
- Cancelamento da requisição sobrevive ao enriquecimento do contexto.
- Falhas de proxy e throttle seguem suas políticas documentadas.

<a name="released-middleware-api"></a>
## API de middleware publicada

`foundation.ApplicationBuilder.WithMiddleware` entrega `*foundation.Middleware`, cujos Prepend, `BeforeThrottle`, Append, Throttle, `ValidateCSRFTokensExcept`, `ImageOrigins`, Observer e Sentences ajustam o pipeline composto. `http.Chain` compõe handlers padrão. O middleware do framework fornece Observe, `RequireAuth` e CSRFProtect; o middleware HTTP do Hesape fornece a normalização de proxies confiáveis e os cabeçalhos de segurança, e o throttle é `hesape/routing/middleware.Throttle` sobre `drivers.Limiter`.

| Superfície | Contrato publicado |
| --- | --- |
| `WithMiddleware` | ajusta o pipeline que o framework compõe |
| Prepend / `BeforeThrottle` / Append | os três lugares onde uma aplicação acrescenta middleware |
| Throttle / `ImageOrigins` / `ValidateCSRFTokensExcept` | o orçamento global, as origens de imagens e as isenções de CSRF |
| Observer / Sentences | substituem no mesmo lugar o Observe e as palavras de recusa |
| `RequireAuth` e CSRFProtect | verificações de identidade e mutação do navegador |

<a name="put-policy-at-the-right-boundary"></a>
## Posicione a policy na fronteira correta

Roteamento escolhe escopo, autenticação carrega identidade, CSRF protege mutações no navegador e rate limit mede uma chave nomeada. Middleware compõe essas fronteiras sem misturar seus significados.

- [Rotas](/docs/{{version}}/http/routing)
- [Autenticação](/docs/{{version}}/security/authentication)
- [CSRF](/docs/{{version}}/http/csrf)
