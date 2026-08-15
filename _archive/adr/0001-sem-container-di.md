# ADR 0001-sem-container-di — Sem container de injeção de dependência

**Status:** Aceita. **Superada em um ponto pela ADR 0023** — o lugar do wiring e
quem o escreve. Ver a nota na Decisão.

## Contexto

Ver docs/01-arquitetura.md.

## Decisão

O wiring é explícito e gerado pela CLI em internal/wire/modules.gen.go. Reflection em Go custa em performance e, pior, move erro de compilação para runtime — o oposto do motivo de escolher Go. O custo é verbosidade no arquivo de wiring, que ninguém edita à mão.

> **Nota — o que a ADR 0023 mudou.** O `internal/wire/modules.gen.go` nunca
> chegou a existir. O que se decidiu aqui — wiring explícito, sem reflection e
> sem container — continua valendo inteiro; mudaram o lugar e o autor. A fiação
> é `bootstrap/app.go`, o equivalente do `bootstrap/app.php` do Laravel, e o
> `aru make:module` **imprime** as três linhas em vez de editar o arquivo. A
> frase "que ninguém edita à mão" está invertida hoje: **só** se edita à mão, e
> é isso que faz o arquivo dizer o que a aplicação é.
>
> Não crie `internal/wire/`. Fiação gerada ao lado da fiação escrita à mão
> seriam duas fiações (REGRA 9).

## Consequência

Revisitar apenas com dado de uso real, não com preferência de estilo.
