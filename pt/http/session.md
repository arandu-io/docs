# Sessão HTTP

- [Um cookie assinado aponta para um Subject tipado](#a-signed-cookie-points-to-a-typed-subject)
- [Escolha o backend conforme a topologia de implantação](#choose-the-backend-for-the-deployment-topology)
- [Rotacione para uma sessão autenticada](#rotate-into-an-authenticated-session)
- [Carregue, confirme, lembre e encerre sessões](#load-confirm-remember-and-end-sessions)
- [Registro ausente permanece não autenticado](#a-missing-record-stays-unauthenticated)
- [Mudanças de identidade exigem mudanças de ID](#identity-changes-require-id-changes)
- [Comprove o fluxo entre réplicas](#prove-the-cross-replica-path)
- [Autenticação define identidade; policy define autoridade](#authentication-owns-identity-policy-owns-authority)

> [!WARNING]
> **Limitações atuais.** `foundation.NewDrivers` constrói o store de sessão a partir de `SESSION_DRIVER`, `REDIS_URL`, `SESSION_LIFETIME` e `SESSION_SECURE_COOKIE`, inclusive um store compartilhado sobre RESP para várias réplicas depois que `hesape/redis` é vinculado, e a produção recusa `SESSION_SECURE_COOKIE=false`. O nome e o escopo do cookie são fixos: `arandu_session`, caminho /, o host que respondeu e `SameSite` Lax; o framework recusa `SESSION_TTL` e um `SESSION_PATH`, `SESSION_DOMAIN` ou `SESSION_SAME_SITE` que peça outro escopo.

<a name="a-signed-cookie-points-to-a-typed-subject"></a>
## Um cookie assinado aponta para um Subject tipado

O armazenamento de sessão do framework envia ao navegador um ID assinado e mantém `security.Subject` no backend. O carregamento verifica o cookie antes de ler o registro, impedindo que dados da requisição escolham tenant ou conta. Sessão estabelece identidade; não concede permissão sobre dados.

<a name="choose-the-backend-for-the-deployment-topology"></a>
## Escolha o backend conforme a topologia de implantação

`foundation.NewDrivers` constrói o único store de sessão, `drivers.Sessions`, a partir da chave da aplicação, de `SESSION_LIFETIME` e da decisão sobre Secure, sobre o backend que `SESSION_DRIVER` nomeia. Memory serve para testes e para um único processo local. Várias réplicas precisam de `SESSION_DRIVER=redis`, do conector `hesape/redis` vinculado pelo seu import em branco em `bootstrap/app.go` e de um `REDIS_URL` acessível. O armazenamento de sessão é independente de `CACHE_STORE`.

| Configuração | Efeito publicado |
| --- | --- |
| `SESSION_DRIVER=memory` | um backend dentro do processo |
| `SESSION_DRIVER=redis` | handler de sessão compartilhado com RESP |
| `SESSION_LIFETIME` | duração em minutos, lida pelo framework; duas horas sem valor |
| `SESSION_SECURE_COOKIE` | atributo Secure, lido pelo framework; ligado exceto com `APP_ENV=dev`, e false é recusado em produção |
| `SESSION_PATH` / `SESSION_DOMAIN` / `SESSION_SAME_SITE` | não são configuráveis; qualquer escopo diferente de /, do host que respondeu e de lax interrompe o boot |

> [!WARNING]
> Nunca implante duas réplicas com backend em memória. Ambos os processos podem iniciar, mas uma sessão criada em um não pode ser carregada no outro.

<a name="rotate-into-an-authenticated-session"></a>
## Rotacione para uma sessão autenticada

Use Rotate no login mesmo sem ID anterior. Ele cria ID novo, escreve cookie assinado e destrói o registro pré-autenticação se existir. Start não expressa substituição, portanto não é a operação de login.

```go example
func RotateSession(r *stdhttp.Request, w stdhttp.ResponseWriter, drivers *foundation.Drivers, subject security.Subject) error {
	sessions := drivers.Sessions

	oldID := sessions.IDFromRequest(r)
	_, err := sessions.Rotate(r.Context(), w, oldID, subject)
	if err != nil {
		return fmt.Errorf("rotate session: %w", err)
	}
	return nil
}
```

> [!NOTE]
> O atributo Secure do store segue `SESSION_SECURE_COOKIE` e fica ligado em todo lugar exceto `APP_ENV=dev` quando a variável não está definida. A produção recusa false no boot, porque ali o cookie é a credencial.

<a name="load-confirm-remember-and-end-sessions"></a>
## Carregue, confirme, lembre e encerre sessões

- Load retorna Subject confiável ou erro explícito de sessão.
- Confirm registra confirmação recente de senha para rotas sensíveis.
- `Remember(true)` estende a sessão pela opção nomeada.
- Destroy remove o registro atual e limpa cookies de sessão e destino pretendido.
- `DestroyOthers` remove outros registros do mesmo tenant e sujeito, preservando o ID escolhido.

### Guarde as sessões no banco de dados da aplicação

```go example
func NewSessions(drivers *foundation.Drivers, appKey []byte, ttl time.Duration) (*security.SessionStore, error) {
	handler, err := session.NewDatabaseHandler[security.Subject](drivers.DB, "")
	if err != nil {
		return nil, err
	}
	backend := security.NewSessionBackend(handler)
	sessions := security.NewSessionStore(appKey, ttl, true, backend)
	return sessions, nil
}
```

`SESSION_DRIVER=redis` resolve sua conexão independentemente de `CACHE_STORE`, então as sessões podem ser compartilhadas enquanto os valores comuns de cache ficam na memória do processo. O Hesape também traz `session.NewDatabaseHandler`, que guarda cada registro na tabela que `session.RecordMigrations` cria, para mais de uma réplica sem servidor RESP. É um componente, não uma configuração: `SESSION_DRIVER` aceita memory ou redis, e a verificação global de CSRF e o throttle leem `drivers.Sessions`, então um store construído sobre a tabela precisa substituir aquele em todos os lugares antes de ser confiável.

<a name="a-missing-record-stays-unauthenticated"></a>
## Registro ausente permanece não autenticado

Cookie ausente, assinatura inválida, registro expirado ou ID ausente não criam Subject parcial. Middleware considera a requisição não autenticada. Indisponibilidade do backend compartilhado é diferente: falhe a requisição dependente ou saúde, sem trocar silenciosamente por memória e dividir as instâncias.

> [!WARNING]
> O nome e o escopo do cookie de sessão não são configuráveis. O armazenamento publicado grava `arandu_session` no caminho /, para o host que respondeu, `HttpOnly` e `SameSite` Lax. O framework interrompe o boot com `SESSION_TTL`, que `SESSION_LIFETIME` substituiu, e com `SESSION_PATH`, `SESSION_DOMAIN` ou `SESSION_SAME_SITE` pedindo outro escopo.

<a name="identity-changes-require-id-changes"></a>
## Mudanças de identidade exigem mudanças de ID

Rotacione após o último fator de autenticação e mudanças de privilégio. Um segundo fator pendente não recebe Subject autenticado. Nunca registre IDs de sessão, aceite ID de formulário, derive tenant de parâmetro nem retorne erro bruto do backend com credenciais de conexão.

`DestroyOthers` tem escopo de tenant porque Subject confiável fornece tenant e ID do sujeito. Não substitua por varredura do backend baseada apenas no ID da conta.

O store lê de novo os papéis de uma sessão a cada requisição pelo resolver que o provider instala com `drivers.Sessions.ResolveSubjects`, de modo que uma conta rebaixada perde o acesso na sua próxima requisição. Um resolver que responde `hhttp.ErrSubjectGone` encerra a sessão de uma conta apagada, fechada ou recusada, e uma sessão com papéis sem resolver instalado é recusada com `hhttp.ErrNoSubjectResolver` em vez de ser confiada como foi escrita.

<a name="prove-the-cross-replica-path"></a>
## Comprove o fluxo entre réplicas

Construa duas instâncias sobre o mesmo servidor RESP. Rotacione uma sessão na primeira, reutilize o cookie na segunda e confirme tenant e sujeito iguais. Teste `SESSION_DRIVER=redis` com `CACHE_STORE=memory`, storage indisponível, driver não suportado, assinatura adulterada, expiração, invalidação do ID antigo e `DestroyOthers` preservando só o escolhido.

```sh
go test ./tests/Feature -run 'Session|Cache'
```

<a name="authentication-owns-identity-policy-owns-authority"></a>
## Autenticação define identidade; policy define autoridade

`RequireAuth` carrega identidade desse armazenamento. O handler ainda chama `security.Authorize` e leva o Grant retornado ao repositório. CSRF vincula mutações do navegador ao ID atual da sessão, portanto a rotação também invalida tokens da fronteira anterior.

- [Autenticação](/docs/{{version}}/security/authentication)
- [Policies e Grant](/docs/{{version}}/security/authorization)
- [Proteção CSRF](/docs/{{version}}/http/csrf)
- [API de sessão do framework](https://pkg.go.dev/github.com/arandu-io/framework/security)
- [Adapter de sessão RESP do Hesape](https://pkg.go.dev/github.com/arandu-io/hesape/redis)
