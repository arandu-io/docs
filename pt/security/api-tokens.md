# Tokens de API

- [`TokenGuard` autentica uma credencial apresentada](#tokenguard-authenticates-one-presented-credential)
- [Resolva bearer tokens com `RequireToken`](#resolve-bearer-tokens-with-requiretoken)
- [Assuma emissão, armazenamento e adapter de requisição](#own-issuance-storage-and-the-request-adapter)
- [Autentique um bearer token previamente provisionado](#authenticate-a-pre-provisioned-bearer-token)
- [Leia credenciais apenas de headers](#read-credentials-from-headers-only)
- [Tokens desconhecidos permanecem respostas 401 normais](#unknown-tokens-remain-ordinary-401-responses)
- [Credencial permanente exige um mecanismo de revogação da aplicação](#a-permanent-credential-needs-an-application-kill-switch)
- [Teste o que o guard faz e o que não faz](#test-what-the-guard-does-and-what-it-does-not-do)
- [Use sessões no navegador e construa o ciclo de vida separadamente](#use-sessions-for-browsers-and-build-lifecycle-separately)

<a name="tokenguard-authenticates-one-presented-credential"></a>
## `TokenGuard` autentica uma credencial apresentada

`auth.TokenGuard` publicado lê um bearer token, consulta `UserProvider` da aplicação pela conta correspondente e a expõe pelo contrato Guard. Não tem sessão e comprova a credencial novamente em cada requisição.

`TokenGuard` é um guard, não um sistema de tokens. Ele não define expiração, escopos, rotação nem revogação, e um token aceito continua aceito enquanto o provedor mantiver um valor correspondente. A aplicação emite, armazena e revoga os próprios tokens.

<a name="resolve-bearer-tokens-with-requiretoken"></a>
## Resolva bearer tokens com `RequireToken`

Para um grupo de rotas de API, `middleware.RequireToken`, de `github.com/arandu-io/hesape/http/middleware`, lê o bearer token do header Authorization e entrega ao `TokenResolver` da aplicação apenas o seu digest SHA-256, como `middleware.DigestToken` o calcula. Token ausente, token desconhecido e token que o resolvedor informa como `ErrUnknownToken` (nunca emitido, revogado, expirado) recebem o mesmo 401. Qualquer outro erro do resolvedor é uma falha em responder e termina em 500, para que uma queda do armazenamento não diga a um cliente com token válido que ele foi retirado. A rota nunca lê a query string, um cookie nem a sessão.

Uma requisição aprovada leva o `auth.Subject` resolvido no contexto, exatamente como atrás de `RequireAuth`, e a Policy continua decidindo cada registro. O `PersonalAccessTokenService` do esqueleto é o lado da emissão: mostra o token uma única vez e guarda só o digest, aceita uma expiração opcional, revoga um token ou todos os tokens de uma conta (todos, quando a conta é encerrada ou a senha é redefinida) e resolve um digest para o Subject atual da conta.

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

O fallback de senha Basic não pode ser desligado por opção de `TokenGuard`. Se o endpoint aceitar apenas bearer, recuse autorização não Bearer antes de invocar o guard, ou use `RequireToken`, que lê só o esquema Bearer.

<a name="unknown-tokens-remain-ordinary-401-responses"></a>
## Tokens desconhecidos permanecem respostas 401 normais

Tokens ausentes, malformados, desconhecidos e desabilitados geram a mesma resposta não autenticada. Em :hesapeVersion, User também reduz erros de `RetrieveByCredentials` a nil: a rota falha fechada como 401 em indisponibilidade do provedor, mas não a distingue de token desconhecido. Registre a falha no provedor ou use uma fronteira própria que preserve erro quando o protocolo precisar de status operacional. Nunca registre o token apresentado.

No Hesape :hesapeVersion, `TokenGuard.Validate` não aplica o SHA-256 configurado com hash=true, embora User aplique. Não use Validate como equivalente a User nesse modo. O teste publicado comprova hashing apenas por User e Check.

<a name="a-permanent-credential-needs-an-application-kill-switch"></a>
## Credencial permanente exige um mecanismo de revogação da aplicação

Digest impede que vazamento do banco revele imediatamente o token, mas sozinho não fornece expiração nem revogação; o serviço de tokens do esqueleto acrescenta uma expiração opcional e apaga a linha ao revogar. O provedor deve recusar registros desabilitados e vincular token ao tenant e conta esperados. Use TLS, não exponha tokens em URLs ou logs, consulte digest por índice e rotacione criando credencial nova antes de remover a antiga.

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
