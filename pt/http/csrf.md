# Proteção CSRF

- [Por que uma mutação precisa de um segundo segredo](#why-a-state-changing-request-needs-a-second-secret)
- [Instale um emissor envolvendo a pilha de rotas](#install-one-issuer-around-the-route-stack)
- [Proteja o handler antes de sua execução](#protect-the-handler-before-it-can-run)
- [Cubra formulários e requisições aprimoradas](#cover-forms-and-enhanced-requests)
- [Recuse sem adivinhar](#refuse-without-guessing)
- [A rotação encerra a fronteira anterior](#rotation-closes-the-old-boundary)
- [Teste a rejeição com o mesmo cuidado do sucesso](#test-the-rejection-as-carefully-as-the-success)
- [CSRF é uma camada de confiança da requisição](#csrf-is-one-layer-of-request-trust)

<a name="why-a-state-changing-request-needs-a-second-secret"></a>
## Por que uma mutação precisa de um segundo segredo

O navegador anexa automaticamente o cookie de sessão a uma requisição compatível. Isso permite que outro site provoque a requisição: estar autenticado não prova a intenção da página. O token CSRF se vincula ao ID da sessão e deve chegar por marcação ou header que outro site não possa ler.

<a name="install-one-issuer-around-the-route-stack"></a>
## Instale um emissor envolvendo a pilha de rotas

`foundation.NewDrivers` constrói o único emissor, `drivers.CSRF`, a partir da chave da aplicação e de `CSRF_TTL`, e o vincula ao store de sessão, `drivers.Sessions`, de modo que o id de sessão só é aceito depois que seu cookie assinado é verificado. Create instala a verificação no pipeline global; nada precisa ser ligado à mão. O token de um visitante anônimo precisa estar vinculado a um id aleatório levado no cookie assinado desse próprio visitante, nunca a um id de sessão vazio que todo visitante compartilharia: um token que um visitante tem precisa ser recusado no formulário de outro. Os cookies de sessão de produção continuam exigindo Secure sobre HTTPS.

<a name="protect-the-handler-before-it-can-run"></a>
## Proteja o handler antes de sua execução

A verificação roda no pipeline global depois do throttle e antes de todo handler de mutação. Requisições seguras GET, HEAD, OPTIONS e TRACE passam. Requisições inseguras leem X-CSRF-Token antes do campo de formulário `_token` e validam contra o id de sessão verificado atual. A única decisão deixada para a aplicação é quais dos seus próprios caminhos têm uma guarda própria, nomeados em `WithMiddleware`.

```go example
func WithoutCSRFOnWebhooks(drivers *foundation.Drivers) *foundation.ApplicationBuilder {
	return foundation.Configure(drivers).
		WithMiddleware(func(m *foundation.Middleware) {
			m.ValidateCSRFTokensExcept(
				"/webhooks/",                 // a prefix: every path under it
				"/v1/usage",                  // exactly this path
				"/checkout/{token}/sessions", // one segment of any value
			)
		})
}
```

<a name="cover-forms-and-enhanced-requests"></a>
## Cubra formulários e requisições aprimoradas

Um formulário convencional usa `@csrf`, compilado em campo oculto. O layout também coloca o token em uma meta tag ypy-csrf, e a extensão ypy-csrf o envia como X-CSRF-Token em toda requisição de mesma origem que altera estado; POST, PUT, PATCH e DELETE parciais, sem formulário, não precisam inventar outra conexão de token. Ao substituir formulário, renderize um campo novo no fragmento.

Uma isenção deixa um caminho para a sua própria guarda, como um webhook que verifica a assinatura do provedor. Uma entrada terminada em / cobre todo caminho abaixo dela, uma entrada sem a barra cobre exatamente aquele caminho, e um segmento `{name}` corresponde a qualquer segmento não vazio, como num padrão de rota; uma entrada não pode começar com um. Um módulo que implementa CSRFExempt declara os próprios receptores: Create os isenta ao lado dos da aplicação, Boot recusa uma isenção que cubra uma rota de escrita de outro dono, e toda isenção coletada é escrita no log no boot.

```kyse
<form method="post" action="/profile">
  @csrf
</form>
```

<a name="refuse-without-guessing"></a>
## Recuse sem adivinhar

Não crie um token substituto ao tratar a mutação rejeitada para continuar a operação. Mutação com origem cruzada informada pelo navegador retorna 403. Token ausente, inválido, expirado ou de outra sessão retorna 419; uma requisição YPY recebe a instrução de atualização do framework, YPY-Refresh. A próxima página fornece token atual sem executar a mutação rejeitada.

<a name="rotation-closes-the-old-boundary"></a>
## A rotação encerra a fronteira anterior

Rotacione o ID da sessão quando a autoridade de autenticação mudar. Tokens do ID anterior deixam de funcionar na nova sessão. Nunca coloque tokens em URLs ou logs, nem isente um endpoint de mutação por receber normalmente requisições aprimoradas: um atacante pode chamá-lo diretamente.

<a name="test-the-rejection-as-carefully-as-the-success"></a>
## Teste a rejeição com o mesmo cuidado do sucesso

Exercite campo válido, header válido, token ausente, alterado, de outra sessão e expirado. Nos formulários renderizados, conte formulários de mutação e diretivas de token para impedir omissão silenciosa em novos formulários. Verifique também que a rejeição mantém o registro protegido intacto.

<a name="csrf-is-one-layer-of-request-trust"></a>
## CSRF é uma camada de confiança da requisição

O token comprova que a requisição passou por uma página que possui a sessão; não decide quem é o sujeito nem suas permissões. Autenticação de sessão, autorização por policy, validação e headers de mesma origem mantêm responsabilidades separadas e ainda executam quando o token é válido.
