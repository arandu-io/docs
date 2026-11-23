# Verificação de e-mail

- [Um código vinculado a propósito confirma uma transição de e-mail](#a-purpose-bound-code-confirms-one-email-transition)
- [Compartilhe armazenamento de uso único e configure e-mail](#share-the-one-time-store-and-configure-mail-delivery)
- [Emita, entregue e consuma pelo serviço](#issue-deliver-and-consume-through-the-service)
- [Renove códigos quando o endereço ou pedido mudar](#rotate-codes-when-the-address-or-request-changes)
- [Códigos inválidos e expirados revelam o mesmo resultado limitado](#invalid-and-expired-codes-reveal-the-same-bounded-result)
- [Verificação não é vinculação de conta](#verification-is-not-account-linking)
- [Exercite emissão, replay, expiração e isolamento de tenant](#exercise-issuance-replay-expiry-and-tenant-isolation)
- [Conclua verificação antes de conceder a sessão final](#finish-verification-before-granting-the-final-session)

<a name="a-purpose-bound-code-confirms-one-email-transition"></a>
## Um código vinculado a propósito confirma uma transição de e-mail

Códigos onetime do Hesape implementam a primitiva publicada e UI :uiVersion fornece o fluxo. O serviço emite código de seis dígitos vinculado ao propósito, e-mail entrega sua única cópia legível e Consume aceita atomicamente uma vez antes da escrita do estado verificado.

<a name="share-the-one-time-store-and-configure-mail-delivery"></a>
## Compartilhe armazenamento de uso único e configure e-mail

Configure `CodeStore` compartilhado, chave criptográfica, transporte de e-mail e serviço de usuário consciente de tenant. Separe propósito de verificação dos de recuperação de senha e dois fatores, mesmo compartilhando armazenamento.

- Use propósito e tenant exatos.
- Renderize o código no mailable de verificação.
- Mantenha códigos brutos fora de logs duráveis.
- Envie pela ponte de e-mail da aplicação.

<a name="issue-deliver-and-consume-through-the-service"></a>
## Emita, entregue e consuma pelo serviço

Construa o sujeito com tenant confiável, ID de usuário e endereço capturado normalizado. Emita e envie o código; depois execute Consume para o mesmo propósito e sujeito antes de `MarkVerified` gravar condicionalmente a conta. Uma repetição não reutiliza o código consumido.

```go
const verifyPurpose = "verify-email"
subject := user.AccountTenant() + "\x00" + user.AccountID() + "\x00" + account.NormalizeEmail(user.AccountEmail())

code, err := codes.Issue(ctx, verifyPurpose, subject)
if err != nil { return err }
if err := mailer.Send(ctx, VerifyEmail{User: user, Code: code}); err != nil { return err }

if err := codes.Consume(ctx, verifyPurpose, subject, submitted); err != nil { return err }
_, _, err = users.MarkVerified(ctx, user.AccountTenant(), user.AccountID(), user.AccountEmail())
return err
```

<a name="rotate-codes-when-the-address-or-request-changes"></a>
## Renove códigos quando o endereço ou pedido mudar

Reemita após mudança de endereço, substitua códigos anteriores conforme o contrato do store, limite envios repetidos e exija novo login após verificar, sem transformar consumo de código em sessão.

- Vincule propósito, tenant, conta e expiração.
- Consuma antes de registrar sucesso.
- Envie substituição sem revelar existência da conta.
- Mantenha falhas de entrega de e-mail observáveis e passíveis de nova tentativa.

<a name="invalid-and-expired-codes-reveal-the-same-bounded-result"></a>
## Códigos inválidos e expirados revelam o mesmo resultado limitado

Códigos desconhecidos, alterados, expirados, já consumidos, de outro propósito e de outro tenant falham sem expor qual conta existe. Um propósito e um subject podem apresentar vinte códigos errados em vinte e quatro horas somando todos os códigos emitidos; depois disso Consume recusa todo código com `ErrTooManyFailures` até a janela fechar, e emitir um código novo não zera a contagem. As telas da ui :uiVersion respondem a essa recusa com "too many attempts, try again later" em vez de chamar o código de inválido. Uma queda de e-mail é diferente de um POST inválido; uma falha de banco de dados depois de Consume precisa de um caminho de reparo idempotente.

> [!WARNING]
> Não marque como verificado apenas por ID de conta. O POST publicado depende de armazenamento de uso único, propósito e endereço capturado.

<a name="verification-is-not-account-linking"></a>
## Verificação não é vinculação de conta

> [!WARNING]
> Nunca vincule contas só pelo mesmo e-mail. Não registre códigos brutos nem os envie a analytics, URLs ou redirecionamentos. Aceite código apenas em POST protegido por CSRF. Verificação confirma controle de um endereço; a policy ainda decide associação ao tenant e mudanças de conta.

Use TLS e expiração curta. Sessão de navegador pode carregar estado do fluxo, mas o registro durável de uso único impede replay entre processos.

<a name="exercise-issuance-replay-expiry-and-tenant-isolation"></a>
## Exercite emissão, replay, expiração e isolamento de tenant

Teste o formulário UI contra Codes do Hesape com códigos válidos, alterados, expirados, consumidos, de propósito/tenant errado, substituição, falha de e-mail e envio entre réplicas. Confirme uma única mudança na conta e nenhuma exposição de outra conta.

- Issue e Consume têm sucesso uma vez.
- Replay e expiração falham.
- Outro propósito ou tenant não pode consumir.
- O consumidor UI grava estado verificado só após Consume.

<a name="finish-verification-before-granting-the-final-session"></a>
## Conclua verificação antes de conceder a sessão final

Autenticação determina quando exigir verificação, recuperação de senha usa outro propósito, sessões rotacionam após transições de identidade e e-mail transporta o código.

- [Autenticação](/docs/{{version}}/security/authentication)
- [Recuperação de senha](/docs/{{version}}/security/password-reset)
- [Sessão](/docs/{{version}}/http/session)
- [E-mail](/docs/{{version}}/services/mail)
