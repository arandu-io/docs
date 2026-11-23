# Limitação de requisições

- [Um limiter mede uma chave explícita em um armazenamento](#a-limiter-meters-one-explicit-key-in-one-store)
- [Construa o limiter e limite uma rota](#construct-the-limiter-and-throttle-a-route)
- [Escolha chaves e política de recusa conforme a ameaça](#choose-keys-and-refusal-policy-by-threat)
- [O throttle publicado permite acesso se o armazenamento falhar](#the-released-throttle-fails-open-on-store-errors)
- [Comprove janelas, concorrência e comportamento em falhas](#prove-windows-concurrency-and-outage-behavior)
- [Limiter de cache e middleware publicados](#released-cache-limiter-and-middleware)
- [Posicione o throttle antes do trabalho caro de autenticação](#place-the-throttle-before-expensive-authentication-work)

<a name="a-limiter-meters-one-explicit-key-in-one-store"></a>
## Um limiter mede uma chave explícita em um armazenamento

O caminho publicado atual é `cache.NewRateLimiter` mais o middleware de roteamento Throttle. O framework já constrói um limiter, `drivers.Limiter`, contando no store que `CACHE_STORE` nomeia, e o pipeline global limita toda requisição a trezentas por minuto por sessão ou endereço; `WithMiddleware` muda esse orçamento com `Middleware.Throttle`. O framework `v0.41.0` removeu a antiga superfície de limiter do framework; não copie exemplos que a importam.

<a name="construct-the-limiter-and-throttle-a-route"></a>
## Construa o limiter e limite uma rota

Passe `drivers.Limiter` para Throttle com limite positivo, janela, função de chave e callback de recusa numa rota que precise de um orçamento mais apertado que o global. O middleware entra em panic na construção quando faltam colaboradores ou a política não é positiva. As telas de login publicadas contam seus formulários de código no mesmo limiter.

```go example
import middleware "github.com/arandu-io/hesape/routing/middleware"

func LoginThrottle(drivers *foundation.Drivers, clientKey func(*stdhttp.Request) string) func(stdhttp.Handler) stdhttp.Handler {
	limiter := drivers.Limiter
	throttle := middleware.Throttle(limiter, cache.PerMinute(5),
		func(r *stdhttp.Request) string { return "login:" + clientKey(r) },
		func(w stdhttp.ResponseWriter, r *stdhttp.Request, status int, message string) {
			stdhttp.Error(w, message, status)
		},
	)
	return throttle
}
```

<a name="choose-keys-and-refusal-policy-by-threat"></a>
## Escolha chaves e política de recusa conforme a ameaça

Em fronteiras de credenciais, combine conta e rede confiável nas chaves; normalize-as, preserve tenant quando conhecido, envie Retry-After e use armazenamento compartilhado para limites entre instâncias. Separe limites gerais de abuso da política estrita de recusa de autenticação.

- Construa uma vez no boot.
- Normalize chaves limitadas.
- Use armazenamento compartilhado em produção.
- Retorne uma recusa 429 explícita.

<a name="the-released-throttle-fails-open-on-store-errors"></a>
## O throttle publicado permite acesso se o armazenamento falhar

O Throttle publicado registra erros do armazenamento e chama o handler: falha aberto. Isso mantém disponibilidade em endpoints comuns, mas não representa proteção fail-closed de login, recuperação de senha ou outra fronteira de segurança sem uma policy da aplicação que recuse acesso quando o armazenamento cair.

> [!WARNING]
> Armazenamentos em memória dividem contadores entre réplicas e zeram no reinício. Falha de armazenamento compartilhado e ataques por cardinalidade de chaves exigem decisões explícitas de saúde, alertas e recusa.

<a name="prove-windows-concurrency-and-outage-behavior"></a>
## Comprove janelas, concorrência e comportamento em falhas

Teste requisições até o limite, recusa seguinte, Retry-After, expiração da janela, incrementos concorrentes, duas instâncias no mesmo armazenamento, headers de proxy não confiáveis, normalização de chaves, falha do armazenamento e panic no boot com configuração inválida.

- O limite compartilhado abrange réplicas.
- A concorrência não ultrapassa a contagem prevista.
- Falha do armazenamento segue a política declarada.
- Hashing de credenciais nunca é alcançado após recusa.

<a name="released-cache-limiter-and-middleware"></a>
## Limiter de cache e middleware publicados

`cache.NewRateLimiter` fornece o contador atual baseado em armazenamento; `routing/middleware.Throttle` adapta-o a HTTP. A aplicação escolhe chaves, limites, armazenamento compartilhado, representação de recusa e política de segurança fail-open ou fail-closed.

| Superfície | Contrato publicado |
| --- | --- |
| `cache.NewRateLimiter` | decisões de janela fixa baseadas no armazenamento |
| `middleware.Throttle` | adapter de chave HTTP e recusa |
| Retry-After | orientação de novas tentativas para o cliente |
| armazenamento compartilhado de cache | estado dos contadores entre instâncias |

<a name="place-the-throttle-before-expensive-authentication-work"></a>
## Posicione o throttle antes do trabalho caro de autenticação

Middleware controla a ordem, autenticação e recuperação definem ameaças sensíveis, cache escolhe estado compartilhado e hashing explica o trabalho caro protegido pelo limiter de credenciais.

- [Middleware](/docs/{{version}}/http/middleware)
- [Autenticação](/docs/{{version}}/security/authentication)
- [Cache e Redis](/docs/{{version}}/services/cache-redis)
- [Hashing](/docs/{{version}}/security/hashing)
