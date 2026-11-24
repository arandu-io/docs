# Notificações

- [Uma notificação seleciona canais para uma mensagem pública](#a-notification-selects-channels-around-one-public-message)
- [Assuma a responsabilidade por preferências, endereços e transportes dos destinatários](#own-recipient-preferences-addresses-and-transports)
- [Construa uma mensagem sem armazenar credenciais](#build-a-message-without-storing-credentials)
- [Separe entrega imediata e entrega durável](#separate-immediate-and-durable-delivery)
- [Trate entregas parciais entre canais como estado explícito](#treat-partial-channel-delivery-as-explicit-state)
- [O endereço do destinatário é um dado protegido](#a-recipient-address-is-protected-data)
- [Teste preferências, tenant e novas tentativas](#test-preference-tenant-and-retry-paths)
- [Use e-mail, filas e outbox respeitando suas fronteiras](#use-mail-queues-and-outbox-as-their-real-boundaries)

<a name="a-notification-selects-channels-around-one-public-message"></a>
## Uma notificação seleciona canais para uma mensagem pública

Uma notificação nomeia seus canais em `Via` e monta uma mensagem para cada um. O esqueleto conecta um `Notifier` com o canal de banco de dados em `app/Providers/AppServiceProvider.go`, e o listener `NotifyNoteAuthor` envia a notificação de exemplo `NotePublished` quando o relay da outbox entrega `note.published`. `aru make:notification` gera uma notificação para os canais que você escolher. Acrescente e-mail, preferências do destinatário e entrega em fila de forma explícita, como descreve o restante desta página.

<a name="own-recipient-preferences-addresses-and-transports"></a>
## Assuma a responsabilidade por preferências, endereços e transportes dos destinatários

Resolva um destinatário confiável pelo serviço da aplicação, leia as preferências no escopo do tenant, configure cada transporte e decida quais entregas serão enfileiradas. Não aceite endereços, tenant ou canais habilitados de um payload de job não confiável.

- Carregue o destinatário por meio de um Grant.
- Aplique as preferências de opt-out no momento da entrega.
- Injete transportes concretos.
- Escolha uma chave de idempotência por canal.

<a name="build-a-message-without-storing-credentials"></a>
## Construa uma mensagem sem armazenar credenciais

Dê uma Key estável à notificação, retorne os nomes de seus canais em Via e implemente a interface de representação de cada canal selecionado. Via sozinho não fornece o corpo de e-mail ou banco de dados. Passe o valor completo ao Notifier injetado, que resolve cada canal conectado e continua coletando falhas.

```go example
type InvoicePaid struct {
	InvoiceID string
	URL       string
}

func (InvoicePaid) Key() notifications.Key { return "billing.invoice-paid" }
func (InvoicePaid) Via(notifications.Notifiable) []notifications.ChannelName {
	return []notifications.ChannelName{
		notifications.ChannelMail,
		notifications.ChannelDatabase,
	}
}

func (n InvoicePaid) ToMail(notifications.Notifiable) messages.Mail {
	return messages.NewMail().
		Subject("Your invoice is paid").
		Line("We received your payment.").
		Action("View invoice", n.URL)
}

func (n InvoicePaid) ToDatabase(notifications.Notifiable) messages.Database {
	return messages.NewDatabase(map[string]any{"invoice_id": n.InvoiceID})
}

func NotifyPaid(ctx context.Context, notifier *notifications.Notifier, grant security.Grant, recipient notifications.Notifiable, invoiceID, invoiceURL string) error {
	notice := InvoicePaid{InvoiceID: invoiceID, URL: invoiceURL}
	return notifier.Send(ctx, grant, recipient, notice)
}
```

<a name="separate-immediate-and-durable-delivery"></a>
## Separe entrega imediata e entrega durável

Resolva as preferências imediatamente antes da entrega, renderize resources públicos específicos por canal, enfileire transportes de rede lentos, associe um id estável de entrega e registre o resultado de cada canal sem considerar o sucesso de um como prova da entrega de todos.

- Mantenha a seleção de canais no escopo do tenant.
- Use mail para a representação de e-mail.
- Use queue para tentativas duráveis.
- Registre um resultado por destinatário e canal.

<a name="treat-partial-channel-delivery-as-explicit-state"></a>
## Trate entregas parciais entre canais como estado explícito

Um canal pode ter sucesso enquanto outro falha. Defina o estado de nova tentativa por canal, evite repetir efeitos já concluídos e diferencie endereço inválido, cancelamento explícito, recusa e indisponibilidade do provedor. O pacote de componentes, sozinho, não fornece esse ciclo completo.

<a name="a-recipient-address-is-protected-data"></a>
## O endereço do destinatário é um dado protegido

> [!WARNING]
> Trate e-mail, telefone, tokens de dispositivo e conteúdo da mensagem como dados protegidos do tenant. Nunca registre credenciais ou códigos de recuperação, permita que o payload substitua o destinatário ou ignore uma preferência atual de cancelamento ao tentar novamente.

A autorização decide o acesso ao evento; a preferência de notificação decide se e por onde enviá-lo. Uma não substitui a outra.

<a name="test-preference-tenant-and-retry-paths"></a>
## Teste preferências, tenant e novas tentativas

Verifique se todo nome retornado por Via tem seu método de representação obrigatório. Depois exercite canais habilitados e desabilitados, destinatário de outro tenant, endereço inválido, indisponibilidade do provedor após sucesso em outro canal, entrega duplicada da fila, logs sem dados sensíveis e mudanças de preferência entre o enfileiramento e a execução.

- O opt-out impede o envio pelo canal correspondente.
- Destinatários de outros tenants não podem ser selecionados.
- A entrega duplicada é idempotente.
- A falha parcial preserva o estado dos canais que tiveram sucesso.

<a name="use-mail-queues-and-outbox-as-their-real-boundaries"></a>
## Use e-mail, filas e outbox respeitando suas fronteiras

Mail formata e-mails, filas fornecem entrega durável, eventos e outbox persistem fatos, e autorização protege a consulta dos destinatários.

- [E-mail](/docs/{{version}}/services/mail)
- [Filas](/docs/{{version}}/async/queues)
- [Eventos e outbox](/docs/{{version}}/async/events-outbox)
- [Autorização](/docs/{{version}}/security/authorization)
