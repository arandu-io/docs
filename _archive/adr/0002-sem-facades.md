# ADR 0002-sem-facades — Sem facades

**Status:** Aceita.

## Contexto

Ver docs/01-arquitetura.md.

## Decisão

Facade em PHP é proxy estático resolvido pelo container. Em Go exigiria estado global mutável e tornaria teste paralelo inseguro. O acesso é sempre por parâmetro ou por context.

## Consequência

Revisitar apenas com dado de uso real, não com preferência de estilo.
