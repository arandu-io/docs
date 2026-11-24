# Console de debug e tracing

- [Inspecione uma requisição sem tornar o tracing universal](#inspect-one-request-without-making-tracing-universal)
- [Escolha acesso de desenvolvimento ou um segredo explícito de produção](#choose-development-access-or-an-explicit-production-secret)
- [Registre uma requisição e consulte-a pelo ID](#record-a-request-and-read-it-by-request-id)
- [Leia os achados antes da timeline bruta](#read-findings-before-the-raw-timeline)
- [Entenda por que uma requisição não tem registro](#know-why-a-request-has-no-record)
- [Minimize a janela de dados diagnósticos](#minimize-the-diagnostic-data-window)
- [Teste estados aberto e fechado](#test-both-the-open-and-closed-states)
- [Reproduza com escopo restrito e remova o acesso](#reproduce-narrowly-and-remove-access)
- [Conecte traces a logs, saúde e chamadas externas](#connect-traces-to-logs-health-and-outbound-calls)

<a name="inspect-one-request-without-making-tracing-universal"></a>
## Inspecione uma requisição sem tornar o tracing universal

O console de debug é um recorder limitado, dentro do processo, dos Collectors de requisições recentes: distribuição da timeline, SQL com origem e argumentos, dumps, eventos, views, HTTP externo, consultas lentas e suspeitas de N+1. Não é APM hospedado, traces duráveis, histórico de métricas ou consulta entre réplicas.

<a name="choose-development-access-or-an-explicit-production-secret"></a>
## Escolha acesso de desenvolvimento ou um segredo explícito de produção

Desenvolvimento monta o recorder e `/_arandu/debug`. Fora dele, deixe `ARANDU_TRACING_SECRET` vazio salvo necessidade de tracing controlado em produção; nesse caso, use ao menos 16 bytes, TLS e X-Arandu-Trace apenas nas requisições e leituras do console investigadas.

- Mantenha tracing de produção desligado por padrão.
- Trate o segredo de tracing como acesso diagnóstico entre tenants.
- Valide `ARANDU_EDITOR` no boot se links para o código estiverem habilitados.
- Durante um incidente, acesse uma réplica conhecida.

<a name="record-a-request-and-read-it-by-request-id"></a>
## Registre uma requisição e consulte-a pelo ID

Faça uma requisição, guarde X-Request-ID e consulte a lista ou detalhe por `aru trace`. A CLI lê a mesma rota JSON protegida do console do navegador e permite selecionar outro host explicitamente.

```bash
set -euo pipefail
aru dev &
DEV_PID=$!
trap 'kill "$DEV_PID"' EXIT
REQUEST_ID=$(curl --retry 20 --retry-connrefused -si http://127.0.0.1:8080/ | awk -F': ' 'tolower($1) == "x-request-id" { gsub("\\r", "", $2); print $2; exit }')
test -n "$REQUEST_ID" || { echo "response has no X-Request-ID" >&2; exit 1; }
aru trace
aru trace "$REQUEST_ID"
```

<a name="read-findings-before-the-raw-timeline"></a>
## Leia os achados antes da timeline bruta

Leia os achados de lentidão e N+1 antes da timeline bruta. Inspecione origem e duração do SQL, chamadas externas, eventos e dumps temporários sem segredos. Cinco consultas idênticas são uma pista, não prova de repetição indevida.

- Encontre a requisição pelo ID.
- Inspecione os achados antes dos eventos individuais.
- Instrumente dependências pelo Collector.
- Use Gauges compartilhados apenas para valores atuais do processo.

<a name="know-why-a-request-has-no-record"></a>
## Entenda por que uma requisição não tem registro

A ausência de registro pode indicar tracing desligado, header ausente, perda do Collector no middleware, atendimento por outra réplica ou descarte pelo buffer circular limitado. Dependências não instrumentadas ficam sem atribuição; registros protegidos ou ausentes retornam 404 deliberadamente. O recorder usa memória local e desaparece ao reiniciar. Rastrear uma requisição de produção não coleta outras requisições nem substitui um backend durável de telemetria.

<a name="minimize-the-diagnostic-data-window"></a>
## Minimize a janela de dados diagnósticos

> [!WARNING]
> Registros do console podem conter argumentos SQL, dumps, payloads de eventos, caminhos e dados de tenants. Nunca exponha a rota por um proxy de conveniência, registre o segredo, coloque-o na URL ou habilite o header indiscriminadamente. Troque ou remova o acesso após o incidente.

<a name="test-both-the-open-and-closed-states"></a>
## Teste estados aberto e fechado

Teste índice e detalhe em desenvolvimento, 404 com produção desabilitada, 404 com segredo errado, header aceito em tempo constante, coleta por requisição, achados, descarte pelo buffer circular e independência dos probes live e health em relação ao acesso ao console.

- Apenas requisições de produção explicitamente rastreadas são registradas.
- Credenciais erradas não revelam o console.
- O registro mais antigo é removido ao atingir a capacidade.
- Saúde permanece um contrato operacional separado.

<a name="reproduce-narrowly-and-remove-access"></a>
## Reproduza com escopo restrito e remova o acesso

Use logs e saúde primeiro. Se necessário, encaminhe uma requisição controlada a uma réplica conhecida, consulte-a imediatamente, preserve evidência mínima sem dados sensíveis e remova ou troque o segredo. Reiniciar apaga a evidência, mas não corrige a dependência.

<a name="connect-traces-to-logs-health-and-outbound-calls"></a>
## Conecte traces a logs, saúde e chamadas externas

Logs estruturados fornecem correlação durável, clientes instrumentados explicam esperas externas, implantação possui probes e documentação de consultas explica achados SQL.

- [Logs](/docs/{{version}}/operations/logging-observability)
- [Cliente HTTP](/docs/{{version}}/services/http-client)
- [Implantação](/docs/{{version}}/operations/deployment)
- [Consultas](/docs/{{version}}/database/queries)
