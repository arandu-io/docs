# Operações de fila

- [Opere a fila pelo binário da aplicação](#operate-the-queue-through-the-application-binary)
- [Execute migrations das filas e registre todos os handlers](#migrate-queue-tables-and-register-every-handler)
- [Inicie um worker que conclua trabalho pendente e inspecione o backlog](#start-one-draining-worker-and-inspect-backlog)
- [Pause a entrada, reinicie workers e gerencie jobs falhos](#pause-intake-restart-workers-and-manage-dead-letters)
- [Trate reservas e jobs falhos como estado recuperável](#treat-leases-and-failed-jobs-as-recoverable-state)
- [Mantenha o escopo do tenant em payload, handler e comando](#keep-tenant-scope-in-payload-handler-and-command)
- [Execute fluxos de sucesso, duplicata, retry e encerramento](#run-one-success-duplicate-retry-and-shutdown-journey)
- [Recupere o estado antes de apagá-lo](#recover-state-before-deleting-it)
- [Separe o desenho da fila de sua operação](#separate-queue-design-from-queue-operations)

<a name="operate-the-queue-through-the-application-binary"></a>
## Opere a fila pelo binário da aplicação

Aru encaminha comandos de fila ao binário do projeto porque só ele conhece handlers, banco, provedor de jobs falhos, repositório de lotes, regras de tenant e conexão escolhida. O esqueleto conecta worker e comandos operacionais; não fornece dashboard nem supervisor.

<a name="migrate-queue-tables-and-register-every-handler"></a>
## Execute migrations das filas e registre todos os handlers

Migre jobs, `failed_jobs` e `job_batches`; registre os handlers no `RegisterJobs` do provider, que `bootstrap/app.go` entrega a `WithJobs`; e defina a fila padrão e um lease de nova tentativa maior que o handler mais lento. `QUEUE_CONNECTION=database` guarda os jobs no banco de dados da aplicação, e redis seleciona a fila RESP depois que `hesape/queue/connectors/redis` é vinculado pelo seu import em branco; nomear redis antes disso interrompe o boot com as duas linhas a acrescentar. Os sinais de pausa e reinício entre processos ainda exigem um cache RESP compartilhado.

- `QUEUE_WORKERS` e `QUEUE_MAX_ATTEMPTS` são positivos quando definidos.
- `REDIS_URL` está presente quando o cache compartilhado de sinais é escolhido.
- Payload e handler preservam tenant e proveniência do Grant.
- Worker e web usam a mesma imagem.

<a name="start-one-draining-worker-and-inspect-backlog"></a>
## Inicie um worker que conclua trabalho pendente e inspecione o backlog

Inicie um worker com fila e concorrência explícitas e monitore essa conexão e fila. Padrões de ambiente preenchem valores omitidos; flags substituem apenas seus campos documentados de invocação. SIGTERM para a entrada e conclui o trabalho em andamento.

```sh
aru queue:work --queue=default --workers=4
aru queue:monitor database:default --max=1000
```

<a name="pause-intake-restart-workers-and-manage-dead-letters"></a>
## Pause a entrada, reinicie workers e gerencie jobs falhos

Pause e retome um destino, solicite reinício gracioso dos workers, inspecione falhas por tenant, repita ou descarte IDs selecionados e pode apenas registros falhos ou de lote revisados. `queue:listen` é um fluxo de recarga em desenvolvimento, não o daemon comum de produção.

- `aru queue:pause database:default --for=30m` e `queue:resume`.
- Execute `aru queue:restart` após uma implantação.
- Execute `aru queue:failed --tenant=<id>` antes de retry ou forget.
- `queue:retry`, `queue:forget` e `queue:retry-batch` também leem suas flags depois dos ids.
- Use retry-batch e comandos de poda no escopo do tenant.

<a name="treat-leases-and-failed-jobs-as-recoverable-state"></a>
## Trate reservas e jobs falhos como estado recuperável

Erros do handler geram novas tentativas até `MaxAttempts` e depois estacionam em `failed_jobs`. Um worker morto só libera o job após a reserva; uma reserva menor que a execução permite duplicidade. Corrija handlers, tabelas, stores escolhidos ou cache compartilhado de sinais ausentes antes de confiar nos comandos de recuperação.

<a name="keep-tenant-scope-in-payload-handler-and-command"></a>
## Mantenha o escopo do tenant em payload, handler e comando

> [!WARNING]
> Comandos de jobs falhos e lotes exigem escopo de tenant. Nunca desserialize um Grant ou papel enviado pelo chamador; reautorize nos handlers, minimize payloads, oculte dados sensíveis dos erros e torne efeitos externos idempotentes, pois a entrega ocorre pelo menos uma vez.

<a name="run-one-success-duplicate-retry-and-shutdown-journey"></a>
## Execute fluxos de sucesso, duplicata, retry e encerramento

Exercite sucesso, duplicata, nova tentativa, expiração de lease, recuperação de queda, pausa e retomada entre processos, drenagem por SIGTERM, isolamento de tenant, a fila de banco de dados publicada e a propagação de `QUEUE_DEFAULT`, `QUEUE_WORKERS`, `QUEUE_RETRY_AFTER` e `QUEUE_MAX_ATTEMPTS` até o worker. Verifique que `QUEUE_CONNECTION=redis` sem o conector vinculado interrompe o boot em vez de recair para o banco de dados.

- Um job real é concluído.
- Outro tenant não pode inspecionar nem repetir o job.
- A falha chega ao armazenamento de jobs falhos e permite nova tentativa.
- O encerramento conclui o trabalho sem aceitar novos jobs.

<a name="recover-state-before-deleting-it"></a>
## Recupere o estado antes de apagá-lo

Pause a fila afetada, preserve registros e logs, corrija ou reverta o handler, repita um job de teste e retome. Clear, forget, flush e prune apagam trabalho ou evidência; faça backup e use o filtro mais restrito de tenant e ID.

<a name="separate-queue-design-from-queue-operations"></a>
## Separe o desenho da fila de sua operação

O guia assíncrono define payloads e handlers; logs expõem evidências dos workers; implantação distribui o mesmo artefato; testes de console exercitam a fronteira de comandos delegados.

- [Desenho da fila](/docs/{{version}}/async/queues)
- [Logs](/docs/{{version}}/operations/logging-observability)
- [Implantação](/docs/{{version}}/operations/deployment)
- [Testes de console](/docs/{{version}}/testing/console)
