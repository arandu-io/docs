# ADR 0004-core-sem-dependencia — Core sem dependência de terceiros

**Status:** Aceita.

## Contexto

Ver docs/01-arquitetura.md.

## Decisão

O core depende apenas da stdlib e de golang.org/x/crypto. Postgres, Redis e NATS entram por módulos separados (arandu-postgres, arandu-redis). Isso é argumento de venda em auditoria de supply chain e evita que uma decisão nossa force uma dependência no cliente.

## Consequência

Revisitar apenas com dado de uso real, não com preferência de estilo.
