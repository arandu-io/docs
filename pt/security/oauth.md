# OAuth e login social

- [Um cliente OAuth conclui a troca do código de autorização](#an-oauth-client-completes-the-authorization-code-exchange)
- [Registre callback exato e mantenha o segredo no boot](#register-an-exact-callback-and-keep-the-secret-at-boot)
- [Redirecione com estado e recuse callback incompatível](#redirect-with-state-and-refuse-a-mismatched-callback)
- [Normalize dados do provedor e aplique a policy da aplicação](#normalize-provider-data-then-apply-application-policy)
- [Recuse antes de trocar o token](#keep-refusal-before-the-token-exchange)
- [Estado é obrigatório; vinculação é autorização](#state-is-mandatory-and-linking-is-authorization)
- [Teste redirecionamento e callback como um par](#test-the-redirect-and-callback-as-one-pair)
- [O callback produz entrada de identidade, não uma sessão](#the-callback-produces-identity-input-not-a-session)

<a name="an-oauth-client-completes-the-authorization-code-exchange"></a>
## Um cliente OAuth conclui a troca do código de autorização

O pacote de provedores publicado envia o navegador ao endpoint externo de autorização, verifica estado do callback, troca o código e recupera `oauth.UserData` normalizado. Há construtores para `GitHub`, Google, Facebook, Stripe e Apple, além de `NewProvider` para serviço configurado explicitamente.

Este pacote é cliente OAuth, não servidor de autorização nem login social completo. O esqueleto e o kit UI não fornecem rotas de callback, configuração de credenciais, vinculação de contas, policy de tenant nem model durável de identidade para esse fluxo; isso continua sendo trabalho da aplicação.

<a name="register-an-exact-callback-and-keep-the-secret-at-boot"></a>
## Registre callback exato e mantenha o segredo no boot

- Registre no provedor a URL HTTPS de callback exatamente como a aplicação a enviará.
- Carregue ID e segredo do cliente da configuração de implantação, nunca do fonte nem da resposta ao navegador.
- Decida o tenant da rota de entrada antes de redirecionar.
- Defina policy de vinculação para identidades existentes antes de aceitar callback.

Construa Provider e `CookieStateStore` novos em cada requisição. Provider mantém estado do par de requisições e não é seguro para compartilhamento concorrente.

<a name="redirect-with-state-and-refuse-a-mismatched-callback"></a>
## Redirecione com estado e recuse callback incompatível

```go example
var (
	clientID     = "client-id"
	clientSecret = "client-secret"
	callbackURL  = "https://example.com/auth/github/callback"
)

func githubProvider(w stdhttp.ResponseWriter, r *stdhttp.Request, id, secret, callback string) *providers.Provider {
	state := providers.NewCookieStateStore(w, r)
	return providers.NewGithubProvider(state, id, secret).RedirectURL(callback)
}

func startGithub(w stdhttp.ResponseWriter, r *stdhttp.Request) error {
	return githubProvider(w, r, clientID, clientSecret, callbackURL).Redirect(w, r)
}

func finishGithub(w stdhttp.ResponseWriter, r *stdhttp.Request) error {
	identity, token, err := githubProvider(w, r, clientID, clientSecret, callbackURL).User(r)
	if errors.Is(err, providers.ErrStateMismatch) {
		stdhttp.Error(w, "invalid oauth callback", stdhttp.StatusBadRequest)
		return nil
	}
	if err != nil {
		return err
	}
	_ = identity // Resolve or link through the tenant-aware account service.
	_ = token    // Keep secret; persist only when the provider integration needs it.
	return nil
}
```

A entrada chama `startGithub` e retorna assim que Redirect escreve 302. O callback é outra requisição tratada por `finishGithub`; User verifica e consome o estado antes de trocar o código. Trate o token retornado pelo provedor como segredo mesmo se só precisar de identidade.

<a name="normalize-provider-data-then-apply-application-policy"></a>
## Normalize dados do provedor e aplique a policy da aplicação

`oauth.UserData.Get` e String leem campos normalizados da identidade. Atributos específicos do provedor permanecem disponíveis sem alterar o formato comum. Após o callback, procure o sujeito previamente vinculado; se não existir, siga decisão explícita de vinculação ou registro, sem associar silenciosamente contas pelo mesmo e-mail.

| Construtor | Provedor publicado |
| --- | --- |
| `NewGithubProvider` | `GitHub` |
| `NewGoogleProvider` | Google |
| `NewFacebookProvider` | Facebook |
| `NewStripeProvider` | Stripe |
| `NewAppleProvider` | Apple |
| `NewProvider` | endpoints fornecidos pelo chamador |

`NewAppleProvider` recebe um Services ID, um Team ID, um Key ID e a chave privada no lugar de um client secret, e monta o client secret a cada troca. A Apple responde ao redirect com um POST vindo do próprio site, então a rota de callback aceita POST e fica isenta da proteção CSRF só para esse caminho, com `CSRFExcept`; a verificação de estado ocupa o lugar dela. O site precisa ser servido por HTTPS, e o cookie de estado é enviado com SameSite=None. A documentação do pacote, em `NewAppleProvider`, descreve o fluxo inteiro.

<a name="keep-refusal-before-the-token-exchange"></a>
## Recuse antes de trocar o token

Estado ausente, cookie ausente, estado alterado e replay retornam `ErrStateMismatch`. Callback com `error=access_denied` é recusa, não falha de servidor. Código ausente, falha de troca, resposta malformada e falha de dados do usuário permanecem erros internos distintos; o navegador recebe mensagem limitada sem segredos.

O cookie de estado tem um nome padrão. Iniciar outro OAuth no mesmo navegador antes do primeiro callback pode substituí-lo e fazer o primeiro falhar; trate como fluxo reiniciável, sem enfraquecer a verificação.

<a name="state-is-mandatory-and-linking-is-authorization"></a>
## Estado é obrigatório; vinculação é autorização

> [!WARNING]
> `Provider.Stateless` desliga a defesa CSRF publicada. Use só quando outro mecanismo concreto e testado vincular callback ao cliente iniciador; conveniência não é esse mecanismo.

Nunca vincule automaticamente só porque um serviço externo retornou e-mail. Verifique a garantia do provedor, trate colisões deliberadamente e autorize a transição de tenant. Não coloque segredos do cliente ou tokens em logs, query, payload de sessão ou erros renderizados.

O `CookieStateStore` :hesapeVersion envia seu cookie de estado com Secure por padrão, qualquer que seja o esquema de chegada da requisição, de modo que o cookie mantém o atributo atrás de um proxy que termina TLS. Um servidor de desenvolvimento em http simples chama `Secure(false)` no store, porque um navegador nunca devolve um cookie Secure por http; deixe-o ligado em todos os outros lugares.

<a name="test-the-redirect-and-callback-as-one-pair"></a>
## Teste redirecionamento e callback como um par

- Redirect retorna 302 com ID do cliente, URI exata de retorno, escopos e estado não vazio.
- Callback com estado armazenado alcança os servidores fake de token e usuário.
- Estado ausente, alterado, reutilizado ou de outro navegador falha antes de qualquer troca.
- Recusa e indisponibilidade do provedor têm evidências operacionais distintas sem expor segredos.
- O caminho HTTPS com proxy implantado emite cookie de estado Secure e `HttpOnly`, com SameSite=Lax (SameSite=None para a Apple).
- Vinculação não pode atravessar tenants nem tomar conta por colisão de e-mail.

```sh
go test github.com/arandu-io/hesape/oauth/providers
```

<a name="the-callback-produces-identity-input-not-a-session"></a>
## O callback produz entrada de identidade, não uma sessão

Após a policy aceitar a vinculação, a aplicação pode converter a conta local em `security.Subject` e rotacionar sua sessão. Tokens de API, tokens do provedor e sessões de navegador têm credenciais e ciclos de vida distintos. Mantenha a fronteira externa do provedor fora do repositório de contas.

- [Autenticação](/docs/{{version}}/security/authentication)
- [Sessão HTTP](/docs/{{version}}/http/session)
- [Policies e Grant](/docs/{{version}}/security/authorization)
- [API de provedores OAuth](https://pkg.go.dev/github.com/arandu-io/hesape/oauth/providers)
- [API de dados de usuário OAuth](https://pkg.go.dev/github.com/arandu-io/hesape/oauth)
