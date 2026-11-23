# Recuperação de senha

- [Recuperação consome uma credencial vinculada ao propósito](#a-reset-consumes-one-purpose-bound-credential)
- [Prepare armazenamento de códigos, e-mail, hashing e sessões](#prepare-code-storage-mail-hashing-and-sessions)
- [Valide, consuma, grave e destrua outras sessões](#validate-consume-write-then-destroy-other-sessions)
- [Evite enumeração nas respostas de pedido e conclusão](#keep-request-and-completion-responses-non-enumerating)
- [Não consuma um código antes de validar o candidato](#do-not-consume-a-code-before-the-candidate-is-valid)
- [Recuperar senha altera todas as fronteiras ativas de identidade](#a-reset-changes-every-active-identity-boundary)
- [Execute o fluxo completo publicado da UI](#run-the-complete-released-ui-journey)
- [Recuperação usa códigos únicos, mas não é verificação de e-mail](#reset-uses-one-time-codes-but-is-not-email-verification)

<a name="a-reset-consumes-one-purpose-bound-credential"></a>
## Recuperação consome uma credencial vinculada ao propósito

UI :uiVersion implementa recuperação atual por códigos onetime do Hesape. O fluxo aceita pedido sem enumeração, entrega código de propósito definido, valida nova senha, consome código, grava hash e destrói outras sessões da conta.

<a name="prepare-code-storage-mail-hashing-and-sessions"></a>
## Prepare armazenamento de códigos, e-mail, hashing e sessões

Ligue Codes, as regras de conta com escopo de tenant (`*account.Users`), o hasher, o mailer e o mesmo `SessionStore` usado no login. Monte a origem da redefinição a partir da configuração e mantenha o propósito de redefinição separado do de verificação e da recuperação do segundo fator.

- Use armazenamento de uso único com duração curta.
- Envie pelo transporte de e-mail configurado.
- Valide senha e confirmação antes de Consume.
- Compartilhe invalidação de sessão entre réplicas.

<a name="validate-consume-write-then-destroy-other-sessions"></a>
## Valide, consuma, grave e destrua outras sessões

Após validar, busque a conta no tenant confiável, consuma código vinculado ao endereço e à impressão da senha atuais e execute o reset compare-and-swap da aplicação. O consumidor publicado tenta `DestroyOthers` após gravar a senha e registra a falha sem desfazer a nova senha.

```go
if err := ctx.Bind(&request); err != nil { return err }
user, err := users.Lookup(ctx, tenant, request.Email)
if err != nil { return err }
fingerprint := account.Fingerprint(user.PasswordHash())
subject := user.AccountTenant() + "\x00" + user.AccountID() + "\x00" +
	account.NormalizeEmail(user.AccountEmail()) + "\x00" + fingerprint
if err := codes.Consume(ctx, "reset-password", subject, request.Code); err != nil { return err }
updated, err := users.ResetPassword(ctx, user.AccountTenant(), user.AccountID(), user.AccountEmail(), fingerprint, request.Password)
if err != nil { return err }
if err := sessions.DestroyOthers(ctx, updated.(*models.User).Subject(), ""); err != nil {
	observability.Log(ctx).Error("signing the account's other sessions out", "error", err)
}
return nil
```

<a name="keep-request-and-completion-responses-non-enumerating"></a>
## Evite enumeração nas respostas de pedido e conclusão

Responda a pedidos de recuperação igualmente com ou sem conta existente; limite requisições por política explícita de segurança, permita que um pedido novo substitua e-mail antigo e preserve a sessão recém-autenticada apenas quando o fluxo publicado a tiver escolhido deliberadamente.

- Use uma resposta pública para e-mails conhecidos e desconhecidos.
- Limite pedido e consumo separadamente.
- Mantenha código e senha fora de dados antigos do formulário e logs.
- Notifique a conta por uma mensagem independente após sucesso.

<a name="do-not-consume-a-code-before-the-candidate-is-valid"></a>
## Não consuma um código antes de validar o candidato

Uma senha inválida não pode consumir um código utilizável. Códigos desconhecidos, expirados, repetidos, de outro propósito e de outro tenant falham da mesma forma, e depois de vinte códigos errados para um propósito e um subject em vinte e quatro horas Consume recusa todo código com `ErrTooManyFailures` até a janela fechar, quantos códigos novos forem enviados; as telas da ui :uiVersion dizem então "too many attempts, try again later". Uma falha de escrita depois de Consume exige um caminho de recuperação explícito; nunca marque sucesso porque o código era válido.

> [!WARNING]
> O fluxo publicado valida, executa Consume, compare-and-swap da senha e então `DestroyOthers` como melhor esforço. Falha no armazenamento de sessões é registrada sem desfazer a senha; o consumidor atual não pode garantir revogação bem-sucedida.

<a name="a-reset-changes-every-active-identity-boundary"></a>
## Recuperar senha altera todas as fronteiras ativas de identidade

> [!WARNING]
> Nunca aceite autoridade de conta ou tenant do formulário de recuperação, revele existência de conta, registre código ou senha, preserve outras sessões por padrão nem use token permanente de API como credencial de recuperação. Use TLS e origens configuradas.

Recuperação comprova controle da credencial de uso único entregue; a policy ainda controla a conta do tenant alterada.

<a name="run-the-complete-released-ui-journey"></a>
## Execute o fluxo completo publicado da UI

Teste pedidos conhecidos e desconhecidos, conclusão válida, confirmação inválida, códigos alterados, expirados, repetidos, de propósito errado e outro tenant, consumo concorrente, falha de e-mail, falha de gravação da senha e duas réplicas com sessões compartilhadas.

- Respostas públicas de pedido não enumeram contas.
- Entrada inválida mantém o código utilizável.
- Consume concorrente tem sucesso uma vez.
- Outras sessões falham após sucesso.

<a name="reset-uses-one-time-codes-but-is-not-email-verification"></a>
## Recuperação usa códigos únicos, mas não é verificação de e-mail

Verificação de e-mail usa outro propósito, hashing protege a nova senha, sessões aplicam invalidação e rate limit restringe tentativas online.

- [Verificação de e-mail](/docs/{{version}}/security/email-verification)
- [Hashing](/docs/{{version}}/security/hashing)
- [Sessão](/docs/{{version}}/http/session)
- [Limitação de requisições](/docs/{{version}}/services/rate-limiting)
