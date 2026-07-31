# Telescope

**Veredito:** DIVIDIDO — parte core, parte Fase 3
**Nome Arandu:** marandu

## O que o Laravel faz

Inspetor de requests em desenvolvimento: queries, jobs, cache, exceções, mails, com dashboard e histórico.

## Decisão

A **captura** é core. O **histórico e a UI** são o `marandu`, na Fase 3.

## Por quê

Esta é a divisão mais importante do mapa, e é onde o Arandu se diferencia do Laravel.

No Laravel, Telescope é um pacote que você instala. Consequência: a maioria dos projetos não tem, e quando o bug aparece em produção ninguém tem dado nenhum.

No Arandu o `Collector` é core e está ligado por padrão em dev. Custo zero em produção — sem header autorizado, `FromContext` devolve nil e todos os `Record*` são no-op com receiver nil.

O `marandu` fica sendo só o que realmente é opcional: o ring buffer histórico e a interface.

## O que conteria, se existir

No core, já implementado: Collector, Dump, página de erro com diagnóstico.

No `marandu`, Fase 3: histórico dos últimos N requests, timeline por request, `aru trace <request_id>`, detector de N+1 com badge.

## Gatilho para construir

Fase 3, conforme roadmap.

## Alternativa enquanto não existe

A página de erro do core já cobre o caso do dia a dia.
