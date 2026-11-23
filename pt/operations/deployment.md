# Implantação

- [Implante um binário imutável em mais de um processo](#deploy-one-immutable-binary-in-more-than-one-process)
- [Prepare configuração de runtime e dependências compartilhadas](#prepare-runtime-configuration-and-shared-dependencies)
- [Construa uma imagem versionada e execute seu job de migration](#build-a-versioned-image-and-run-its-migration-job)
- [Atualize web, worker e tarefas agendadas pela mesma imagem](#roll-web-worker-and-scheduled-work-from-the-same-image)
- [Interrompa o rollout antes que uma réplica sem saúde receba tráfego](#stop-the-rollout-before-an-unhealthy-replica-receives-traffic)
- [Não inclua contexto de build, credenciais nem execução como root](#ship-no-build-context-credentials-or-root-runtime)
- [Exercite o contrato do container antes da promoção](#exercise-the-container-contract-before-promotion)
- [Reverta o artefato sem tratar rollback de schema como gratuito](#roll-back-the-artifact-without-pretending-schema-rollback-is-free)
- [Continue com atualizações, migrations e operações de fila](#continue-with-upgrades-migrations-and-queue-operations)

> [!WARNING]
> **Limitações atuais.** Binário, container e etapa explícita de migration estão publicados. Um contrato completo de saúde e rollback automático ainda não é universal: defina probes pelas dependências realmente conectadas e ensaie rollback fora do boot de requisições.

<a name="deploy-one-immutable-binary-in-more-than-one-process"></a>
## Implante um binário imutável em mais de um processo

O esqueleto publicado constrói um binário da aplicação. O mesmo artefato serve HTTP, processa filas e executa console com argumentos diferentes, evitando divergência entre web e worker. Seu Dockerfile multi-stage compila Kyse, embute assets, desativa CGO, remove caminhos e executa apenas `/app` como processo distroless não root.

<a name="prepare-runtime-configuration-and-shared-dependencies"></a>
## Prepare configuração de runtime e dependências compartilhadas

Forneça `APP_KEY`, `APP_URL`, `HTTP_ADDR` e `DATABASE_URL` de produção, e apenas as configurações que algum consumidor lê. A produção recusa `MAIL_URL=log://` e um `MAIL_URL` não definido, então nomeie um mailer real, e recusa `SESSION_SECURE_COOKIE=false`. Um `APP_ENV` ausente é produção; a imagem publicada escreve ENV `APP_ENV=prod` para que a resposta fique declarada em vez de herdada. As migrations de produção precisam de um lock de cache compartilhado, então selecione RESP e forneça `REDIS_URL`. Mantenha `.env`, dados do Git, storage e o SQLite local fora da imagem.

- Banco e serviços RESP escolhidos estão acessíveis.
- Arquivos de confiança TLS estão presentes quando configurados.
- Segredos vêm da plataforma.
- Probes de liveness e readiness usam rotas separadas.

<a name="build-a-versioned-image-and-run-its-migration-job"></a>
## Construa uma imagem versionada e execute seu job de migration

Construa uma imagem amarrada a uma versão e a um commit completo e depois rode a migration como um job avulso do pipeline antes do rollout, antes que qualquer binário novo pegue jobs. O comando serve padrão escuta em `HTTP_ADDR` e escreve primeiro uma linha application environment, com env, debug e `http_addr`; as migrations adquirem seu lock compartilhado e recusam o isolamento em memória em produção.

```sh
COMMIT_SHA=$(git rev-parse HEAD)
docker build --build-arg VERSION=:aranduVersion --build-arg COMMIT="$COMMIT_SHA" -t registry.example/arandu::aranduVersion .
docker run --rm --env-file ./production.env registry.example/arandu::aranduVersion migrate
```

<a name="roll-web-worker-and-scheduled-work-from-the-same-image"></a>
## Atualize web, worker e tarefas agendadas pela mesma imagem

Implante web, workers de fila e comandos agendados a partir do mesmo digest. Migre uma vez antes de trocar as réplicas, permita que SIGTERM drene o trabalho, use `/_arandu/live` para decisões de reinício e `/_arandu/health` para decidir a prontidão para tráfego. A rota de saúde executa no máximo dez passadas por segundo sobre suas verificações, uma de cada vez, e responde a uma sonda além da décima com o último resultado, de modo que as sondas de um balanceador nunca se acumulam.

- Execute serve nas réplicas web.
- Execute a imagem com work `--queue=default --workers=4`; work é o destino interno de `aru queue:work`.
- Grave VERSION e COMMIT para investigar incidentes.
- Direcione tráfego pela saúde, nunca pela simples liveness.

<a name="stop-the-rollout-before-an-unhealthy-replica-receives-traffic"></a>
## Interrompa o rollout antes que uma réplica sem saúde receba tráfego

Falhas de compilação de views, configuração tipada, conector, lock ou execução de migration impedem a promoção. Uma réplica nova com falha de prontidão de módulo fica fora de serviço enquanto o log estruturado identifica a fronteira que falhou.

> [!WARNING]
> A imagem publicada não executa migrations em ENTRYPOINT ou CMD. Adicioná-las ali cria concorrência entre réplicas que iniciam. Distroless não tem shell: a recuperação usa logs, probes, imagem substituta e comandos reproduzíveis.

<a name="ship-no-build-context-credentials-or-root-runtime"></a>
## Não inclua contexto de build, credenciais nem execução como root

> [!WARNING]
> Nunca copie `.env` ou banco local para a imagem. Preserve execução nonroot, builds estáticos com CGO desativado, runtime distroless, Aru fixada e digest do registro vinculado ao commit. Deixe diagnósticos de produção desligados salvo controle de um segredo explícito.

<a name="exercise-the-container-contract-before-promotion"></a>
## Exercite o contrato do container antes da promoção

Em um checkout limpo, gere views, teste com race, execute vet e build, crie a imagem, inicie com configuração semelhante à produção, verifique ambos os probes, execute migrate duas vezes, teste o binário anterior com schema aditivo, conclua trabalho de um worker e confira usuário de execução e ausência de segredos.

- A imagem executa como nonroot.
- Health controla tráfego e live controla reinício.
- Repetir uma migration é seguro.
- SIGTERM conclui o trabalho em andamento.

<a name="roll-back-the-artifact-without-pretending-schema-rollback-is-free"></a>
## Reverta o artefato sem tratar rollback de schema como gratuito

Retire o tráfego do artefato defeituoso e reimplante a imagem imutável anterior. Migrations aditivas devem mantê-la funcional. Execute `migrate:rollback` apenas se o lote exato tiver Down testado; remover colunas de forma destrutiva exige dois releases, não um comando de emergência.

<a name="continue-with-upgrades-migrations-and-queue-operations"></a>
## Continue com atualizações, migrations e operações de fila

Atualizações escolhem versões, migrations coordenam schema, operações de fila coordenam workers e observabilidade comprova o artefato em execução.

- [Atualizações](/docs/{{version}}/operations/upgrades)
- [Migrations](/docs/{{version}}/database/migrations)
- [Operações de fila](/docs/{{version}}/operations/queues)
- [Observabilidade](/docs/{{version}}/operations/observability)
