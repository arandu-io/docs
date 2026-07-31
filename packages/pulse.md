# Pulse

**Veredito:** NÃO CONSTRUIR — OpenTelemetry
**Nome Arandu:** —

## O que o Laravel faz

Métricas de performance da aplicação em produção: requests lentos, queries lentas, uso de fila, usuários ativos.

## Decisão

Exportar OpenTelemetry. Não construir produto próprio.

## Por quê

Esta é a decisão que mais economiza tempo em todo o escopo.

APM é um mercado com Datadog, Grafana e New Relic, todos com times de dezenas de pessoas. Construir um concorrente sozinho, como pacote de um framework, produz algo que ninguém vai preferir ao que já paga.

E o cliente **já tem** um. A pergunta dele não é "onde vejo minha latência", é "por que o Arandu não aparece no meu Grafana".

O Collector já tem todos os dados. Falta só o exportador.

## O que conteria, se existir

Um exportador OTel no core, opcional, desligado por padrão.

## Gatilho para construir

Nenhum para produto próprio. O exportador entra na Fase 3.

## Alternativa enquanto não existe

OpenTelemetry, para o backend que o cliente já usa.
