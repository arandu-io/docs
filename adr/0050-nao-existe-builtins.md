# ADR 0050 — Não existe `Builtins`: o despacho de comando é o `switch` do esqueleto

**Status:** Aceita. Data: 14/08/2026.

## A porta que se fecha

O `framework` deixa de poder fornecer uma lista de comandos prontos. Não haverá
`foundation.Builtins()`, nem registro de comando no `Application`, nem gancho
por onde um módulo acrescente um comando ao binário do projeto.

## Contexto

`docs/31-reorganizacao-hesape.md:190` atribuiu ao `foundation` um
`Builtins() []console.Command` — *"o que o binário do projeto responde sem o
`aru`"*. Quando o `foundation` foi escrito, em 14/08/2026, ele foi o único item
da especificação que não virou código, e a razão não foi falta de tempo: a forma
não sai do que um `Application` tem.

O que o esqueleto faz hoje é um `switch` de onze casos em
`arandu/bootstrap/console.go:60-105`: `serve`, `migrate`, `migrate:rollback`,
`migrate:status`, `migrate:fresh`, `routes`, `db:seed`, `schedule:list`,
`schedule:run`, `work` e `Version`.

Desses onze, **três** são respondíveis a partir de um `Application` — `serve`,
`routes` e `Version`. Os outros oito não:

- os quatro de `migrate` e o `db:seed` precisam de um `*data.DB`, que o
  `Application` não abre e não segura;
- `db:seed` chama código do projeto, que o framework não conhece;
- `schedule:list` e `schedule:run` precisam do módulo de scheduler que o projeto
  registrou por último;
- `work` precisa de um store de fila vindo de `github.com/arandu-io/queue`, um
  módulo separado **que importa o framework**. O framework não pode importá-lo
  de volta.

## Decisão

**`Builtins` não existe, em forma nenhuma.** O despacho de comando é o `switch`
do esqueleto, escrito à mão, no repositório do projeto.

## As três razões, em ordem de peso

**1. O nome não é do Laravel.** A [ADR 0044](0044-nome-de-funcao-e-do-laravel-nome-de-ferramenta-e-guarani.md)
é literal: nome de função é o do Laravel, sem exceção. `Builtins` não aparece em
nenhum lugar do Illuminate — foi inventado pela especificação. Um símbolo
exportado com nome que o desenvolvedor Laravel nunca viu falha o filtro antes de
qualquer discussão sobre utilidade.

**2. O que o Laravel tem no lugar é exatamente o que rejeitamos.**
`Illuminate\Foundation\Providers\ArtisanServiceProvider` registra **110**
comandos, e o faz sendo um *service provider*: cada um entra no container por
nome e é resolvido de lá quando o `artisan` decide qual rodar. A
[ADR 0001](0001-sem-container-di.md) rejeitou o container e a
[ADR 0002](0002-sem-facades.md) rejeitou a fachada que o alcança. Sem os dois,
não existe o lugar onde os 110 são registrados, e portar a lista sem portar o
mecanismo é portar a metade que não funciona.

**3. Um gancho de registro é o container voltando pela porta dos fundos.** A
única forma de o `work` chegar ao binário a partir do framework seria o
`Application` aceitar comandos de quem se registrasse — um mapa de nome para
função, preenchido em tempo de execução por quem for importado. Isso é um
container com outro nome, e a ADR 0001 não o rejeitou pelo nome.

## O que se escreve no lugar, e por que ele é melhor

O `switch` do esqueleto. Ele é feio de propósito e essa é a propriedade:

```go
switch command {
case "migrate":
    return data.Migrate(ctx, db, migrations.All())
case "work":
    return worker.Run(ctx)
}
```

Está no repositório do projeto, o desenvolvedor lê os onze casos de uma vez, e
acrescentar o décimo segundo é escrever um `case`. Nada é resolvido por nome em
tempo de execução, nada aparece no binário porque um pacote foi importado, e a
pergunta *"de onde vem este comando?"* se responde lendo um arquivo.

É a mesma decisão do `bootstrap/app.go`, pela mesma razão que a ADR 0001 dá: o
*composition root* é visível, e o que o compilador verifica é o que roda.

## O que esta ADR não decide

**O `aru` continua sendo a CLI.** Nada aqui muda o que
`github.com/arandu-io/aru` oferece — `aru new`, `aru make:*`, `aru doctor`,
`aru view:build`. A [ADR 0006](0006-cli-em-modulo-separado.md) já separou as
duas coisas: a `aru` é um binário próprio, e o que ela precisa do projeto ela
obtém delegando ao binário do projeto. Esta ADR é sobre o binário do projeto, e
não sobre a ferramenta.

**O gerador pode gerar o `case`.** Se um dia `aru make:command` existir, o que
ele faz é escrever um `case` no `switch` e o arquivo do comando ao lado — a
mesma coisa que `aru make:module` já faz com as três linhas de fiação que ele
imprime para colar (REGRA 12). Gerar não é registrar em tempo de execução.

## Consequência

- `docs/31-reorganizacao-hesape.md:190` fica superado no ponto do `Builtins`. O
  resto daquela linha — `Kernel` e `Module` mantendo o nome do tipo — foi
  cumprido.
- `framework/foundation/doc.go` deixa de descrever o `Builtins` como decisão
  pendente e passa a citar esta ADR.
- O `switch` do esqueleto passa a ser superfície com decisão registrada. Quem
  propuser substituí-lo por registro automático está reabrindo a ADR 0001.
