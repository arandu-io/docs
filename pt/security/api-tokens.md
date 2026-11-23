# Tokens de API

- [`TokenGuard` autentica uma credencial apresentada](#tokenguard-authenticates-one-presented-credential)
- [Assuma emissão, armazenamento e adapter de requisição](#own-issuance-storage-and-the-request-adapter)
- [Autentique um bearer token previamente provisionado](#authenticate-a-pre-provisioned-bearer-token)
- [Leia credenciais apenas de headers](#read-credentials-from-headers-only)
- [Tokens desconhecidos permanecem respostas 401 normais](#unknown-tokens-remain-ordinary-401-responses)
- [Credencial permanente exige um mecanismo de revogação da aplicação](#a-permanent-credential-needs-an-application-kill-switch)
- [Teste o que o guard faz e o que não faz](#test-what-the-guard-does-and-what-it-does-not-do)
- [Use sessões no navegador e construa o ciclo de vida separadamente](#use-sessions-for-browsers-and-build-lifecycle-separately)

> [!WARNING]
> **Limitações atuais.** `TokenGuard` autentica bearer ou token de entrada, mas não fornece expiração, escopos, rotação ou revogação. Em modo hash, use apenas requisição até User ou Check: Validate envia token bruto. User reduz erros do provedor a nil e `SetRequest` não limpa usuário em cache; construa um guard por requisição e preserve erros operacionais em fronteira própria quando o protocolo exigir distinção.

<a name="tokenguard-authenticates-one-presented-credential"></a>
## `TokenGuard` autentica uma credencial apresentada

`auth.TokenGuard` publicado lê um bearer token, consulta `UserProvider` da aplicação pela conta correspondente e a expõe pelo contrato Guard. Não tem sessão e comprova a credencial novamente em cada requisição.

> [!WARNING]
> Este não é um sistema de emissão ou ciclo de vida. O guard publicado não define expiração, escopos, rotação, revogação, models de token, rotas de gestão nem armazenamento de primeira parte. O token é aceito enquanto o provedor mantiver valor correspondente.

<a name="own-issuance-storage-and-the-request-adapter"></a>
## Assuma emissão, armazenamento e adapter de requisição

Antes de conectar o guard, gere token opaco de alta entropia, mostre-o uma vez, armazene só SHA-256 e implemente `auth.UserProvider` sobre o digest. Adapte também a requisição a `auth.Request`; os tipos HTTP publicados ainda não implementam diretamente todos os métodos.

- Gere credenciais com `crypto/rand`, não nomes, timestamps ou UUIDs escolhidos como identificadores.
- Armazene digest em registro dedicado de credencial, não em coluna reutilizável de senha ou remember token.
- Decida vínculo ao tenant e estado desabilitado na busca do provedor.
- Forneça um caminho operacional separado para substituição e remoção; `TokenGuard` não faz nenhuma delas.

<a name="authenticate-a-pre-provisioned-bearer-token"></a>
## Autentique um bearer token previamente provisionado

Este exemplo cobre só autenticação contra token já existente. Embutir Request do Hesape fornece acessores de query, input, bearer e cookies; os três métodos abaixo completam o contrato `auth.Request` exato.

```go example
type tokenRequest struct{ *hhttp.Request }

func (r tokenRequest) GetUser() string {
	user, _, _ := r.RawRequest().BasicAuth()
	return user
}
func (r tokenRequest) GetPassword() string {
	_, password, _ := r.RawRequest().BasicAuth()
	return password
}
func (r tokenRequest) Context() context.Context {
	return r.RawRequest().Context()
}

func Authenticate(w stdhttp.ResponseWriter, r *stdhttp.Request, users auth.UserProvider) {
	request := tokenRequest{Request: hhttp.NewRequest(r)}
	guard := auth.NewTokenGuard(users, request, "api_token", "api_token_sha256", true)
	user := guard.User()
	if user == nil {
		stdhttp.Error(w, "unauthorized", stdhttp.StatusUnauthorized)
		return
	}
}
```

> [!NOTE]
> Com hash=true, User calcula SHA-256 do token recebido antes de `RetrieveByCredentials`. O campo de armazenamento do provedor precisa conter o digest hexadecimal minúsculo esperado.

<a name="read-credentials-from-headers-only"></a>
## Leia credenciais apenas de headers

`GetTokenForRequest` prefere Authorization: Bearer e usa senha HTTP Basic como fallback. Ignora query e corpo para manter credenciais fora de URLs, histórico, referrer e logs comuns. Crie `TokenGuard` novo por requisição. Em :hesapeVersion, `SetRequest` troca a requisição mas não limpa usuário em cache de User; não é seguro reutilizá-lo após resolução.

> [!WARNING]
> O fallback de senha Basic não pode ser desligado por opção de `TokenGuard`. Se o endpoint aceitar apenas bearer, recuse autorização não Bearer antes de invocar o guard.

<a name="unknown-tokens-remain-ordinary-401-responses"></a>
## Tokens desconhecidos permanecem respostas 401 normais

Tokens ausentes, malformados, desconhecidos e desabilitados geram a mesma resposta não autenticada. Em :hesapeVersion, User também reduz erros de `RetrieveByCredentials` a nil: a rota falha fechada como 401 em indisponibilidade do provedor, mas não a distingue de token desconhecido. Registre a falha no provedor ou use uma fronteira própria que preserve erro quando o protocolo precisar de status operacional. Nunca registre o token apresentado.

> [!WARNING]
> No Hesape :hesapeVersion, `TokenGuard.Validate` não aplica SHA-256 configurado com hash=true, embora User aplique. Não use nem documente Validate como equivalente a User nesse modo. O teste publicado comprova hashing apenas por User e Check.

<a name="a-permanent-credential-needs-an-application-kill-switch"></a>
## Credencial permanente exige um mecanismo de revogação da aplicação

Digest impede que vazamento do banco revele imediatamente o token, mas não fornece expiração nem revogação. O provedor deve recusar registros desabilitados e vincular token ao tenant e conta esperados. Use TLS, não exponha tokens em URLs ou logs, consulte digest por índice e rotacione criando credencial nova antes de remover a antiga.

Não confunda autenticação com autorização. Converta a conta autenticada em Subject e execute a mesma Policy usada por requisições autenticadas por sessão.

<a name="test-what-the-guard-does-and-what-it-does-not-do"></a>
## Teste o que o guard faz e o que não faz

- Um header Bearer resolve a conta e Check tem sucesso.
- Query e corpo nunca autenticam.
- Digest armazenado só resolve com hash=true pelo caminho User.
- Tokens ausentes, alterados e desabilitados retornam 401 sem ecoar a credencial.
- Erro do provedor vira nil e o mesmo 401 por `TokenGuard.User`, enquanto a telemetria do provedor registra a queda.
- Um guard que resolveu um usuário nunca é reutilizado por `SetRequest` em outra requisição.
- Recusa de policy ainda bloqueia token autenticado.

```sh
go test github.com/arandu-io/hesape/auth -run TokenGuard
```

<a name="use-sessions-for-browsers-and-build-lifecycle-separately"></a>
## Use sessões no navegador e construa o ciclo de vida separadamente

Autenticação de navegador usa sessões assinadas e CSRF. `TokenGuard` só cabe quando a aplicação já forneceu o ciclo de credenciais ausente. OAuth obtém identidade de provedor e não substitui um emissor de tokens de API da aplicação.

- [Autenticação](/docs/{{version}}/security/authentication)
- [Policies e Grant](/docs/{{version}}/security/authorization)
- [OAuth e login social](/docs/{{version}}/security/oauth)
- [API `TokenGuard`](https://pkg.go.dev/github.com/arandu-io/hesape/auth#TokenGuard)
