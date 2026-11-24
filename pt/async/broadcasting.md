# Broadcasting

- [O nome de um canal é uma fronteira tipada de publicação](#a-channel-name-is-a-typed-publication-boundary)
- [Escolha transporte e autorização antes de publicar](#choose-transport-and-authorization-before-publishing)
- [Publique um evento no escopo de um tenant](#publish-one-tenant-scoped-event)
- [Mantenha explícitos a construção dos canais e os payloads](#keep-channel-construction-and-payloads-explicit)
- [Não confunda canal tipado com distribuição concluída](#do-not-confuse-a-typed-channel-with-delivered-fan-out)
- [Autorize assinatura e publicação separadamente](#authorize-subscription-and-publication-separately)
- [Teste isolamento de canais e recusa do transporte](#test-channel-isolation-and-transport-refusal)
- [Use outbox e filas quando a entrega precisar sobreviver](#use-outbox-and-queues-when-delivery-must-survive)

<a name="a-channel-name-is-a-typed-publication-boundary"></a>
## O nome de um canal é uma fronteira tipada de publicação

O pacote de broadcasting oferece canais tipados, eventos e contratos de driver. O esqueleto não conecta um gerenciador de broadcast nem uma assinatura no navegador por você: você escolhe um driver no bootstrap e escreve a policy de assinatura, e esta página cobre a fronteira de canais e eventos entre os dois.

<a name="choose-transport-and-authorization-before-publishing"></a>
## Escolha transporte e autorização antes de publicar

Selecione um driver concreto no bootstrap, defina como um Subject confiável pode se inscrever e decida se a perda de uma conexão é aceitável. Use outbox ou fila durável quando o evento precisar sobreviver à perda do processo, sem tratar broadcast como armazenamento.

- Nomeie o recurso ao construir o canal.
- Derive o namespace do tenant apenas do Grant.
- Limite o tamanho do payload e seus campos públicos.
- Escolha separadamente o trabalho que precisa ser durável.

<a name="publish-one-tenant-scoped-event"></a>
## Publique um evento no escopo de um tenant

Construa um canal com a API publicada a partir do identificador do registro, defina explicitamente o conteúdo público e chame `Broadcaster.Broadcast` com o Grant autorizado. O driver aplica o namespace do tenant por `TenantChannel`; a entrada da aplicação nunca adiciona esse prefixo.

```go example
func PublishInvoiceUpdate(ctx context.Context, broadcaster broadcasting.Broadcaster, grant security.Grant, invoiceID, reference string) error {
	channel := broadcasting.NewPrivateChannel("invoices." + invoiceID)
	payload := map[string]any{
		"invoice": map[string]any{
			"id":        invoiceID,
			"reference": reference,
		},
	}
	if err := broadcaster.Broadcast(
		ctx, grant, []broadcasting.Channel{channel}, "invoice.updated", payload,
	); err != nil {
		return err
	}
	return nil
}
```

<a name="keep-channel-construction-and-payloads-explicit"></a>
## Mantenha explícitos a construção dos canais e os payloads

Mantenha nomes de eventos estáveis, serialize um recurso deliberado em vez de um model, não anexe credenciais bearer ou de sessão e reconecte clientes pelo caminho normal de assinatura autenticada, sem confiar em uma string de canal fornecida pelo navegador.

- Use tipos explícitos para canais públicos e privados.
- Mantenha os dados do evento pequenos e versionáveis.
- Publique fatos imutáveis.
- Faça replay durável apenas por um sistema de eventos persistidos.

<a name="do-not-confuse-a-typed-channel-with-delivered-fan-out"></a>
## Não confunda canal tipado com distribuição concluída

Indisponibilidade do driver, cliente desconectado, falha de serialização e assinatura negada são resultados distintos. A API tipada publicada, sozinha, não garante reconexão, replay, ordenação nem um gerenciador integrado ao esqueleto. Teste juntos o bootstrap, a policy de assinatura, o transporte escolhido e a entrega na aplicação que você implanta.

<a name="authorize-subscription-and-publication-separately"></a>
## Autorize assinatura e publicação separadamente

> [!WARNING]
> Nunca aceite um prefixo de tenant em um canal enviado pelo cliente, coloque segredos ou modelos completos no payload, nem trate o conhecimento do nome de um canal privado como autorização. `RequestedChannel` rejeita dois-pontos, e `TenantChannel` deriva o namespace apenas do Grant autorizado. Autentique e autorize cada inscrição.

A entrega por broadcast é uma representação, não uma autorização. Um cliente que recebe um identificador ainda precisa passar pela policy normal para qualquer leitura ou alteração posterior.

<a name="test-channel-isolation-and-transport-refusal"></a>
## Teste isolamento de canais e recusa do transporte

Use um driver fake para verificar canal e payload exatos, recuse a inscrição de outro tenant, injete falhas de serialização e de driver, teste cancelamento e execute o fluxo do transporte escolhido antes de habilitá-lo na implantação.

- Os nomes de canais contêm o escopo confiável de tenant.
- Assinantes de outro tenant são recusados.
- Payloads expõem apenas campos públicos.
- Um erro do driver permanece observável.

<a name="use-outbox-and-queues-when-delivery-must-survive"></a>
## Use outbox e filas quando a entrega precisar sobreviver

Eventos e outbox fornecem fatos duráveis; filas fornecem trabalho durável; autorização protege a assinatura; recursos definem payloads públicos.

- [Eventos e outbox](/docs/{{version}}/async/events-outbox)
- [Filas](/docs/{{version}}/async/queues)
- [Autorização](/docs/{{version}}/security/authorization)
- [Recursos HTTP](/docs/{{version}}/http/resources)
