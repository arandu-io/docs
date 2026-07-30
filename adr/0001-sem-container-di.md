# ADR 0001-sem-container-di — Sem container de injeção de dependência

**Status:** Aceita.

## Contexto

Ver docs/01-arquitetura.md.

## Decisão

O wiring é explícito e gerado pela CLI em internal/wire/modules.gen.go. Reflection em Go custa em performance e, pior, move erro de compilação para runtime — o oposto do motivo de escolher Go. O custo é verbosidade no arquivo de wiring, que ninguém edita à mão.

## Consequência

Revisitar apenas com dado de uso real, não com preferência de estilo.
