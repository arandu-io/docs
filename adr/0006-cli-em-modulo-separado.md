# ADR 0006-cli-em-modulo-separado — A CLI é módulo próprio, e delega ao binário do projeto

**Status:** Aceita. **Superada em um ponto pela ADR 0019** — o caminho do ponto
de entrada. Ver a nota na Decisão.

## Contexto

O esqueleto trazia a CLI em `cmd/aru/` dentro do módulo do framework. Isso
contradiz `05-repositorios.md`, item 5: se a CLI vive no módulo do framework,
todo projeto que importa o framework arrasta as dependências dela — e a CLI vai
depender de bubbletea na Fase 2.

Há um segundo problema, que só aparece quando se tenta implementar `aru routes`:
o wiring é explícito e não existe container nem carregamento de plugin (ADR 0001).
Logo, **uma CLI compilada separadamente não tem como saber quais módulos uma
aplicação registrou.** Não é limitação de implementação; é consequência direta de
não ter reflection resolvendo dependência.

## Decisão

Duas partes.

1. A CLI vive em `arandu-io/aru`, módulo próprio, sem dependência. `cmd/aru/`
   sai do framework.

2. Os comandos se dividem por aquilo de que precisam:

   | Precisa | Comandos | Onde roda |
   |---|---|---|
   | só arquivo | `key:generate`, `new`, `make:*`, `doctor` | na própria CLI |
   | os módulos registrados | `serve`, `migrate`, `routes` | delegado ao binário do projeto |

   A delegação é `go run ./cmd/app <comando>` a partir da raiz do projeto. Por
   isso o esqueleto tem esse despacho em `cmd/app/main.go`: é lá que os módulos
   existem.

> **Nota — o que a ADR 0019 mudou.** O `cmd/app/` deixou de existir. A árvore
> passou a espelhar o Laravel, e no Laravel o ponto de entrada fica na raiz: o
> `main.go` do projeto substitui `public/index.php` **e** o `artisan`. A
> delegação hoje é `go run . <comando>` — a constante `appPackage` em
> `aru/delegate.go` vale `"."`, e o `projectRoot` procura `go.mod`, `main.go` e
> `arandu.toml` juntos.
>
> O **motivo** da decisão não mudou: uma CLI compilada à parte não conhece os
> módulos registrados, então ela executa o projeto. Só o caminho mudou.
> `cmd/app/main.go` não existe em projeto nenhum; não o recrie.

## Consequência

`aru serve` exige a toolchain Go na máquina — o que é verdade em
desenvolvimento por definição. Em produção roda o binário do projeto direto, sem
a CLI.

Revisitar apenas se aparecer uma forma de listar módulos sem executar a
aplicação, o que exigiria abandonar o wiring explícito.
