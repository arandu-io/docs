# Horizon

**Veredito:** DEFERIDO
**Nome Arandu:** potiron `[?]` — **e o `[?]` é o ponto.** O nome nunca passou os
cinco filtros da REGRA 4, e como o veredito é deferido ninguém volta a ele. Se o
gatilho disparar, o nome passa pelos cinco **antes** de virar código; entrar por
inércia é como um nome reprovado chega a extensão de arquivo.

## O que o Laravel faz

Dashboard e supervisor para filas Redis: throughput, tempo de espera, jobs falhados, retry, balanceamento de workers.

## Decisão

Não construir antes de existir fila em produção.

## Por quê

Dashboard de fila é útil quando há volume. Antes disso é UI bonita sobre uma tabela vazia.

A ordem correta é: existir fila (NATS, Fase 3), existir carga real, e só então descobrir o que você precisa ver. Construir o dashboard antes é adivinhar quais métricas importam.

## O que conteria, se existir

- ring buffer de jobs, com filtro por fila e status
- retry e cancelamento manual, com autorização por Policy
- métricas: throughput, latência, taxa de falha por tipo de job

## Gatilho para construir

Fila rodando em produção há pelo menos um mês, e você abrindo o banco na mão para investigar job falhado.

## Alternativa enquanto não existe

Logs estruturados via slog, mais a UI de monitoramento do próprio NATS.
