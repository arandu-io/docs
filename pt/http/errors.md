# Tratamento de erros

- [Um handler classifica, registra e renderiza](#one-handler-classifies-reports-and-renders)
- [Construa o handler pela configuração carregada](#build-the-handler-from-loaded-configuration)
- [Instale recuperação como middleware mais externo](#install-recovery-as-the-outermost-middleware)
- [Retorne uma falha HTTP deliberada](#return-an-intentional-http-failure)
- [Erros desconhecidos continuam sendo erros de servidor](#unknown-errors-stay-server-errors)
- [Detalhe público e causa privada são campos diferentes](#public-detail-and-private-cause-are-different-fields)
- [Comprove representação e registro do erro](#prove-both-representation-and-report)
- [Recupere o serviço, não apenas a requisição](#recover-the-service-not-only-the-request)
- [Acompanhe falhas através das fronteiras](#follow-failures-across-their-boundaries)

<a name="one-handler-classifies-reports-and-renders"></a>
## Um handler classifica, registra e renderiza

Erros atravessam uma fronteira central de recuperação. Falhas esperadas de domínio e transporte preservam significado; panics escapados e erros internos desconhecidos são registrados uma vez e viram resposta de produção limitada, sem stack nem segredos. Mapear um erro de domínio para um status é tarefa da aplicação: retorne um erro que declare o próprio status e nunca deduza o status pelo texto do erro.

<a name="build-the-handler-from-loaded-configuration"></a>
## Construa o handler pela configuração carregada

Create constrói o único handler de exceções a partir da configuração carregada, lê o caminho de módulo do binário para separar os frames de código da aplicação dos do framework e instala Recover no lado mais externo depois das frases de recusa. A página de depuração existe só em desenvolvimento; a validação de produção a mantém desligada. A aplicação registra suas próprias respostas nesse handler em `WithExceptions`.

<a name="install-recovery-as-the-outermost-middleware"></a>
## Instale recuperação como middleware mais externo

Retorne erros das ações e deixe a recuperação escolher a resposta configurada. Em desenvolvimento, a página de erro pode mostrar trechos do fonte; produção envia status genérico e contexto de correlação enquanto logs estruturados preservam a causa encadeada.

```go example
var ErrQuotaExceeded = errors.New("quota exceeded")

func WithErrorHandling(drivers *foundation.Drivers) *foundation.ApplicationBuilder {
	return foundation.Configure(drivers).
		WithExceptions(func(e *exception.Handler) {
			e.DontReport(ErrQuotaExceeded)
			e.ShouldRenderJSONWhen(func(r *stdhttp.Request, _ error) bool {
				return strings.HasPrefix(r.URL.Path, "/api/")
			})
		})
}

func Register(router *fhttp.Router, show func(*fhttp.Context) error) {
	router.Action("GET", "/invoices/{invoice}", show)
}
```

<a name="return-an-intentional-http-failure"></a>
## Retorne uma falha HTTP deliberada

Defina tipos de erro da aplicação para ausência, negação, conflito, estado inválido e indisponibilidade de dependência. O adaptador do router mapeia um erro retornado, puro ou embrulhado com %w: um sentinela de não encontrado responde 404, `security.ErrForbidden` 403, `security.ErrCSRF` 419, `database.ErrUniqueViolation` 409, e um erro com o método `HTTPStatus()` int responde esse status. Dê esse método a um erro de domínio em vez de tratá-lo em cada controller, embrulhe as causas com contexto e evite inferir status comparando strings de erro arbitrárias.

- Retorne erros de validação pelo adapter dedicado.
- Use 404 sem revelar existência em outro tenant.
- Use 409 para conflito real de estado.
- Use 503 para indisponibilidade de dependência obrigatória.

<a name="unknown-errors-stay-server-errors"></a>
## Erros desconhecidos continuam sendo erros de servidor

Escrever parte da resposta antes de retornar erro pode impedir mudar o status final. Panic após enviar headers é registrado, mas não permite substituição limpa; conclua trabalho que pode falhar antes de confirmar a saída.

> [!WARNING]
> Páginas de fonte e dumps de desenvolvimento podem conter segredos. Mantenha o console desligado ou protegido fora do desenvolvimento e remova credenciais, cookies, Grants e URLs de banco dos valores estruturados.

<a name="public-detail-and-private-cause-are-different-fields"></a>
## Detalhe público e causa privada são campos diferentes

> [!WARNING]
> Separe detalhe público da causa privada encadeada. Nunca exponha SQL, URLs de conexão, chaves, tokens, cookies, corpos de requisição, stack traces ou detalhes de policy; mantenha debug desligado em qualquer lugar acessível por pessoas não confiáveis.

<a name="prove-both-representation-and-report"></a>
## Comprove representação e registro do erro

Exercite cada erro esperado pela rota registrada e provoque panic antes e após uma escrita. Confirme que corpos de produção não revelam a causa, o recorder observa um evento e recusa de domínio não vira falha interna.

- Falhas esperadas preservam seu status.
- Erros desconhecidos são registrados uma vez.
- Respostas de produção não contêm stack nem credenciais.
- Ausência por outro tenant é indistinguível de registro inexistente.

<a name="recover-the-service-not-only-the-request"></a>
## Recupere o serviço, não apenas a requisição

Use o ID da resposta para localizar log e trace correlacionados, inspecione diagnósticos do kernel e saúde dos módulos dependentes e corrija ou reverta o release. Repita uma mutação apenas após conhecer seu resultado de commit e garantir idempotência.

<a name="follow-failures-across-their-boundaries"></a>
## Acompanhe falhas através das fronteiras

Observabilidade registra o evento, validação possui erros de campos, controllers retornam falhas e implantação decide se diagnósticos de desenvolvimento podem ser habilitados.

- [Logs e observabilidade](/docs/{{version}}/operations/logging-observability)
- [Controllers](/docs/{{version}}/http/controllers)
- [Validação](/docs/{{version}}/http/validation)
