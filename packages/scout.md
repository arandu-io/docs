# Scout

**Veredito:** DEFERIDO
**Nome Arandu:** a definir (heka vetado no filtro zero)

## O que o Laravel faz

Busca full-text com drivers para Algolia, Meilisearch e banco relacional. Sincroniza índice automaticamente quando o modelo muda.

## Decisão

Não construir antes de um cliente real esbarrar no limite do Postgres.

## Por quê

O full-text do Postgres com `tsvector` e índice GIN resolve busca para a esmagadora maioria dos SaaS. Ele é transacional, não precisa de serviço extra, e não tem problema de sincronização de índice.

O Scout existe porque no Laravel a sincronização com serviço externo é chata. Só vale construir quando alguém provar que precisa de serviço externo.

## O que conteria, se existir

- interface Indexable e sincronização por evento
- adapter Postgres (tsvector) e adapter Meilisearch
- respeito ao escopo de tenant vindo do Grant — busca que ignora tenant vaza dado entre clientes

## Gatilho para construir

Cliente com volume ou requisito (typo tolerance, facetas) que o Postgres não atende.

## Alternativa enquanto não existe

`tsvector` + GIN, com as queries escritas à mão nos templates do gerador — o
sqlc não entrou (ADR 0024).
