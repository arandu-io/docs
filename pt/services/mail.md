# E-mail

- [O mailable possui a representação; o transporte possui a entrega](#a-mailable-owns-representation-and-the-transport-owns-delivery)
- [Configure um mailer e a origem pública no boot](#configure-one-mailer-and-public-origin-at-boot)
- [Construa e envie uma mensagem tipada](#build-and-send-one-typed-message)
- [Enfileire entregas lentas com proveniência de tenant](#queue-slow-delivery-with-tenant-provenance)
- [Separe falhas de renderização, anexos e transporte](#keep-render-attachment-and-transport-failures-distinct)
- [E-mail é uma fronteira de divulgação externa](#mail-is-an-external-disclosure-boundary)
- [Capture o contrato MIME exato](#capture-the-exact-mime-contract)
- [Use e-mail dentro de fluxos da aplicação vinculados a um propósito](#use-mail-inside-purpose-bound-application-flows)

> [!WARNING]
> **Limitações atuais.** O componente de e-mail do Hesape e a ponte do framework renderizam e entregam as mensagens documentadas. Armazenamento de anexos e autorização assíncrona não estão completos de ponta a ponta; anexe apenas dados já autorizados e enfileire e-mails com proveniência confiável persistida.

<a name="a-mailable-owns-representation-and-the-transport-owns-delivery"></a>
## O mailable possui a representação; o transporte possui a entrega

Hesape :hesapeVersion possui mensagens e transportes de e-mail; `github.com/arandu-io/framework/mail` :frameworkVersion fornece a ponte de compatibilidade usada pela aplicação. O exemplo abaixo usa `Address{Email: ...}` e `ToAddress` dessa ponte; Mailer de nível inferior do Hesape usa To.

<a name="configure-one-mailer-and-public-origin-at-boot"></a>
## Configure um mailer e a origem pública no boot

`foundation.NewDrivers` constrói `drivers.Mail` a partir de `MAIL_URL`, `MAIL_FROM_ADDRESS` e `MAIL_FROM_NAME`. `smtp://` faz o upgrade com STARTTLS quando o servidor oferece, `smtps://` fala TLS implícito na porta 465, e `resend://` e `sendgrid://` alcançam esses provedores. `log://` escreve a mensagem inteira no log e é recusado em produção, assim como um `MAIL_URL` não definido, que significa o mesmo. Um provider injeta `drivers.Mail` em serviços ou handlers enfileirados em vez de ler o ambiente ou abrir SMTP num controller.

- Use a identidade From configurada.
- Construa links pelo `APP_URL` confiável.
- Mantenha secretas as credenciais de transporte.
- Escolha entre entrega imediata e durável.

<a name="build-and-send-one-typed-message"></a>
## Construa e envie uma mensagem tipada

Construa uma mensagem tipada com dados públicos da aplicação e envie pelo transporte injetado com contexto. Controller ou serviço de domínio fornecem destinatário e propósito; o transporte nunca descobre contas nem policy de tenant.

```go example
type VerifyEmail struct{ Code string }

func (VerifyEmail) Envelope() mail.Envelope {
	return mail.Envelope{Subject: "Verify your address"}
}
func (m VerifyEmail) Content() mail.Content {
	return mail.Content{Text: "Use this code: " + m.Code}
}

type Recipient struct{ Email, Name string }

func SendVerification(ctx context.Context, mailer *mail.Mailer, recipient Recipient, code string) error {
	return mailer.ToAddress(mail.Address{
		Email: recipient.Email, Name: recipient.Name,
	}).Send(ctx, VerifyEmail{Code: code})
}
```

<a name="queue-slow-delivery-with-tenant-provenance"></a>
## Enfileire entregas lentas com proveniência de tenant

Enfileire a entrega pela rede quando a resposta não precisa esperar, preserve tenant e ação no Grant do job, renderize texto e HTML de forma coerente, use anexos limitados e pertencentes à aplicação e mantenha `log://` apenas como evidência de desenvolvimento. Para pôr um log de atividades ou um filtro na frente de toda mensagem, chame `drivers.WrapMail` antes de `drivers.Mail` ser entregue a qualquer coisa: um serviço que pegou o mailer antes continua com o que pegou.

- Passe contexto à entrega.
- Use uma única origem pública.
- Enfileire com proveniência de Grant.
- Anexe apenas conteúdo autorizado e limitado.

<a name="keep-render-attachment-and-transport-failures-distinct"></a>
## Separe falhas de renderização, anexos e transporte

Renderização, validação de endereço, abertura ou leitura de anexos, TLS, autenticação, recusa do provedor e timeout precisam de evidência operacional própria. Não marque mensagens de verificação ou recuperação como enviadas se apenas sua construção teve sucesso.

> [!WARNING]
> A integração publicada não torna automaticamente seguros todos os caminhos de anexo ou Grant. Abra anexos dentro da entrega autorizada, limite bytes e faça handlers de fila reconstruírem a policy atual sem confiar nos caminhos do payload.

<a name="mail-is-an-external-disclosure-boundary"></a>
## E-mail é uma fronteira de divulgação externa

> [!WARNING]
> Nunca registre credenciais SMTP, códigos de uso único em produção, listas completas de destinatários, anexos ou respostas de provedor com dados pessoais. Impeça injeção de headers, use TLS e derive destinatário e tenant pelos serviços da aplicação.

Um link ou código entregue continua sendo entrada da aplicação vinculada a um propósito. A entrega de e-mail não autentica o callback nem autoriza mudanças de conta.

<a name="capture-the-exact-mime-contract"></a>
## Capture o contrato MIME exato

Use transporte de captura para From, To, assunto, texto e HTML, headers, anexos, cancelamento, recusa do provedor e remoção de dados sensíveis. Execute a integração da aplicação e o caminho enfileirado, não apenas testes unitários do construtor de mensagem.

- O MIME capturado tem os campos públicos esperados.
- Uma falha de transporte permanece como evidência para nova tentativa.
- O cancelamento interrompe a entrega.
- Segredos não aparecem nos logs.

<a name="use-mail-inside-purpose-bound-application-flows"></a>
## Use e-mail dentro de fluxos da aplicação vinculados a um propósito

Verificação de e-mail e recuperação de senha definem propósitos, filas fornecem entrega durável, notificações compõem canais e configuração escolhe o transporte.

- [Verificação de e-mail](/docs/{{version}}/security/email-verification)
- [Recuperação de senha](/docs/{{version}}/security/password-reset)
- [Filas](/docs/{{version}}/async/queues)
- [Configuração](/docs/{{version}}/architecture/configuration)
