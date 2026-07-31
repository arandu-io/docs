# ADR 0015-otel-tem-gatilho — O exportador OpenTelemetry espera um cliente que o use

**Status:** Aceita. Detalha `03-roadmap-fases.md` (entregas da Fase 3) e
`06-escopo.md`.

## Contexto

A Fase 3 lista, entre as entregas: *"exportador OpenTelemetry opcional, para
quem já tem Grafana/Datadog"*. E o `06-escopo.md` fecha a questão do APM:

> **exportar OTel.** Construir APM próprio é competir com Grafana e Datadog.

A decisão de **não** construir APM está tomada e continua valendo. O que este
documento decide é *quando* construir o exportador.

Três coisas ficaram claras ao chegar aqui:

1. **O critério de saída da Fase 3 não cita OTel.** Os três critérios são: um bug
   real diagnosticado só pelo console, overhead zero medido, e um dev externo
   entendendo um erro sem ler documentação. Nenhum passa por exportador.
2. **O SDK do OTel é grande.** `go.opentelemetry.io/otel` mais o exporter OTLP e
   as dependências gRPC somam mais que todo o resto do grafo do framework junto.
   Em Go não existe dependência opcional (doc 11), então isso é módulo separado —
   o décimo repositório.
3. **O desenho correto depende de onde a instrumentação entra**, e há duas
   respostas plausíveis: envolver o `http.Handler` e o driver `database/sql` com
   `otelhttp` e `otelsql`, ou dar ao core um ponto de extensão próprio. A
   primeira é barata e usa o que já existe; a segunda evita depender de dois
   pacotes de terceiros no caminho quente. **Escolher sem um cliente real é
   adivinhar**, e a escolha errada custa uma major.

## Decisão

Não construir agora. O exportador entra como `arandu-io/otel` quando o gatilho
disparar.

**Gatilho:** um usuário — interno ou externo — com Grafana, Datadog ou Honeycomb
já em produção, que peça para ver o request do Arandu junto do resto do trace
dele. Nesse dia, o desenho se decide contra o que ele já usa, não contra o que
imaginamos que ele usaria.

## Por quê

**Complexidade especulativa é o que o `05-repositorios.md` já proíbe** em outro
contexto: *"fatiar antes de alguém pedir é complexidade especulativa"*. Vale
igual aqui. Um exportador sem consumidor é código que ninguém executa, com
dependências que todo mundo audita.

**O que a Fase 3 entrega já resolve o problema dela.** O console, a timeline, o
`aru trace`, o diagnóstico na página de erro e o `X-Arandu-Trace` cobrem o
momento em que alguém está debugando. OTel cobre outro momento — a correlação
com o resto da infraestrutura — e esse momento chega com o primeiro cliente que
tem infraestrutura.

**O que já está pronto para o dia do gatilho:**

- o `Collector` já registra query, dump, evento e chamada externa com duração e
  origem: é exatamente o material de um span
- `Timeline` já separa SQL, render e externo
- o `X-Request-ID` já entra e sai do request, e `sanitizeRequestID` já rejeita o
  que não é hexadecimal — o mesmo campo que carrega `traceparent` depois
- `kernel.Diagnostic` já é o ponto onde um módulo reporta estado

Nada disso precisa mudar para o exportador existir. É o que torna adiar barato.

## Consequências

- **Fecha:** APM próprio segue fora de escopo, para sempre. Este documento adia
  o exportador, não reabre a construção de produto de observabilidade.
- A Fase 3 fecha com os três critérios de saída, e a entrega de OTel fica
  registrada como adiada em vez de silenciosamente omitida.
- `arandu-io/otel` continua reservado como décimo repositório. A REGRA 3 —
  *dez repositórios até a Fase 3, não crie o décimo primeiro sem pedir* — segue
  intacta.
- Se o gatilho não disparar até a Fase 4 fechar, vale reavaliar se ele existe: um
  exportador que ninguém pediu em duas fases provavelmente não é o que falta.
