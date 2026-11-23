# Logs e observabilidade

- [Escreva eventos estruturados pelo contexto da requisição](#write-structured-events-through-the-request-context)
- [Defina o nível e preserve a ordem dos middlewares](#set-the-level-and-preserve-middleware-order)
- [Emita um evento de domínio pesquisável](#emit-one-searchable-domain-event)
- [Instrumente esperas externas e trabalho em segundo plano](#instrument-outbound-waits-and-background-work)
- [Torne visível a ausência de telemetria](#make-missing-telemetry-visible)
- [Registre identificadores, nunca credenciais ou autoridade](#log-identifiers-never-credentials-or-authority)
- [Verifique os campos na fronteira do handler](#assert-fields-at-the-handler-boundary)
- [Diagnostique a partir do último sinal durável](#diagnose-from-the-last-durable-signal)
- [Abra a timeline da requisição apenas quando necessário](#open-the-request-timeline-only-when-needed)

> [!WARNING]
> **Limitações atuais.** Os logs estruturados e a coleta de requisições estão publicados, e o logger raiz segue `LOG_LEVEL` e `LOG_FORMAT`. `LOG_SAMPLING` é lido e não tem efeito, exportação e retenção são escolhas de implantação, e nenhum segredo ou detalhe interno de Grant pertence a uma linha de log.

<a name="write-structured-events-through-the-request-context"></a>
## Escreva eventos estruturados pelo contexto da requisição

Create constrói um único logger slog raiz a partir das configurações de logging: a saída padrão em `LOG_LEVEL`, no `LOG_FORMAT`, que é texto em desenvolvimento e JSON nos demais quando vazio. Cada linha leva o ambiente como channel e nenhuma posição de código. Observe deriva um logger de requisição com `request_id`, method e path e, ao concluir, registra status, duração, bytes, número de consultas e tempo de SQL.

<a name="set-the-level-and-preserve-middleware-order"></a>
## Defina o nível e preserve a ordem dos middlewares

Defina um `LOG_LEVEL` publicado e, quando o padrão do ambiente não servir ao coletor, `LOG_FORMAT`. O pipeline composto já executa Observe depois de Recover, dos limites de corpo e de `TrustProxies`. Passe `context.Context` por controllers, serviços e trabalho externo para que `observability.Log` recupere o logger e o Collector do escopo.

- Use o logger raiz configurado.
- Preserve o contexto da requisição.
- Mantenha recuperação fora da observação.
- Envie stdout ao coletor da plataforma.

<a name="emit-one-searchable-domain-event"></a>
## Emita um evento de domínio pesquisável

Emita texto estável de eventos com identificadores estruturados por `observability.Log(ctx)`. O logger está sempre disponível; construir outro logger global perde campos da requisição e o nível configurado.

```go example
type InvoiceSender interface {
	Send(ctx context.Context, id string) error
}

type InvoiceService struct{ sender InvoiceSender }

func (s *InvoiceService) Send(ctx context.Context, id string) error {
	observability.Log(ctx).Info("invoice send requested", "invoice_id", id)
	return s.sender.Send(ctx, id)
}
```

<a name="instrument-outbound-waits-and-background-work"></a>
## Instrumente esperas externas e trabalho em segundo plano

Use `observability.Client` para chamadas externas limitadas, `WithLogger` para trabalho em segundo plano controlado e o Recorder da aplicação nas opções de fila ou agendador que o aceitam. A plataforma cuida de coleta de stdout, retenção, indexação, rotação e alertas.

- Correlacione eventos por `request_id`.
- Defina timeout para clientes externos.
- Mantenha explícito o escopo dos loggers de segundo plano.
- Trate retenção externa como configuração de implantação.

<a name="make-missing-telemetry-visible"></a>
## Torne visível a ausência de telemetria

Perder o contexto usa o logger padrão e perde a correlação. Um cliente não instrumentado omite tempo externo e pode esperar indefinidamente sem timeout. Métodos do Collector desabilitado são no-ops intencionais, não prova de instrumentação executada.

> [!WARNING]
> No framework :frameworkVersion, `LOG_SAMPLING` é lido e validado em `Settings.Logging` e não chega a nenhum logger; não o opere como controle. Uma consulta no agregador de logs que lê source não encontra nada, porque as linhas levam channel no lugar.

<a name="log-identifiers-never-credentials-or-authority"></a>
## Registre identificadores, nunca credenciais ou autoridade

> [!WARNING]
> Nunca registre senhas, tokens, cookies, payloads de sessão, `APP_KEY`, segredos de tracing, corpos de requisição, registros sem limite ou detalhes completos do Grant. Argumentos SQL e dumps são dados diagnósticos sensíveis entre tenants; produção recusa nível debug.

<a name="assert-fields-at-the-handler-boundary"></a>
## Verifique os campos na fronteira do handler

Capture um `slog.Handler`, faça requisições com ids válidos e malformados, verifique os campos de conclusão, exercite uma chamada externa instrumentada, confira JSON em produção e texto em desenvolvimento, prove que `LOG_FORMAT` sobrescreve o padrão do ambiente e que o filtro de nível vale, e verifique que `LOG_SAMPLING` não tem efeito.

- IDs válidos de requisição são devolvidos e registrados.
- IDs malformados são substituídos.
- O tempo externo chega ao Collector.
- `LOG_SAMPLING` não é anunciado como efetivo.

<a name="diagnose-from-the-last-durable-signal"></a>
## Diagnostique a partir do último sinal durável

Correlacione `request_id` nos logs estruturados duráveis, reproduza com tracing autorizado apenas quando seguro e restaure um `LOG_LEVEL` suportado se a configuração impedir o boot. Corrija um coletor externo na plataforma; a aplicação não possui arquivo local de log a recuperar.

<a name="open-the-request-timeline-only-when-needed"></a>
## Abra a timeline da requisição apenas quando necessário

O console diagnóstico detalha uma requisição, clientes HTTP contribuem com tempos externos, implantação possui a coleta de logs e configuração define os controles efetivos.

- [Observabilidade diagnóstica](/docs/{{version}}/operations/observability)
- [Cliente HTTP](/docs/{{version}}/services/http-client)
- [Implantação](/docs/{{version}}/operations/deployment)
- [Configuração](/docs/{{version}}/architecture/configuration)
