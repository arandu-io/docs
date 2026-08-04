# ADR 0017-dependencia-da-cli — A CLI pode ter dependência; o core continua sem

**Status:** Aceita. Detalha a ADR 0004 (core sem dependência) e a ADR 0006 (CLI
em módulo separado), sem alterar nenhuma das duas.

## Contexto

Até a Fase 4 o `aru` tinha **zero** dependências, e o `plans/checklist.sh` e o CI
verificavam isso. Não era uma decisão registrada — era uma propriedade que
aconteceu de ser verdade, e que dois guards passaram a defender.

A DSL quebrou isso. O formato escolhido pelo `19-dsl-e-geracao.md` é YAML, com
justificativa: massa de treino enorme, validação madura, schema publicável. Go
não tem parser de YAML na biblioteca padrão.

As alternativas foram consideradas e são piores:

- **Escrever um parser** — YAML é grande. Um subconjunto seria exatamente o
  formato que um modelo eventualmente escreve fora, e o erro apareceria como
  "spec inválida" para um documento que qualquer outro parser aceita
- **Trocar YAML por JSON** — desfaz a decisão do doc 19 pelo motivo errado, e
  JSON não tem comentário, que é onde a spec explica a si mesma
- **TOML** — a stdlib também não tem, e a massa de treino é menor

## Decisão

**O core continua com duas dependências: `golang.org/x/crypto` e o `x/sys` que
ela arrasta.** A ADR 0004 não muda, e é a que importa: todo projeto que importa
o framework paga o grafo dele.

**A CLI pode ter dependência, uma a uma, com justificativa.** Hoje é uma:

| Dependência | Por quê | Grafo próprio |
|---|---|---|
| `gopkg.in/yaml.v3` | o formato da DSL, escolhido pelo doc 19 | nenhum |

Essa separação é o motivo de a ADR 0006 ter posto a CLI em módulo próprio. O
texto dela já dizia: *"se a CLI viver dentro do módulo do framework, todo projeto
que importa o framework arrasta as dependências da CLI"*. A CLI existe separada
justamente para poder ter dependências.

## A regra, e o que a torna real

Entrar na tabela acima custa **um ADR**, não um `go get`. O critério:

1. resolve um problema que a stdlib não resolve;
2. tem grafo próprio pequeno ou nulo;
3. é mantida, e a API é estável há anos.

`yaml.v3` passa nos três: zero dependências próprias, e o formato não muda.

O `plans/checklist.sh` e o CI do `aru` passam a verificar contra a **lista**, e
não contra o zero. Um `go get` de qualquer outra coisa reprova, com o nome do
pacote na mensagem. É o mesmo atrito de antes, aplicado onde ele ainda faz
sentido.

## Consequências

- **Fecha:** o core nunca ganha dependência por conveniência de ferramenta. Se
  algo do `aru` precisar de biblioteca, ele fica no `aru`
- O guard antigo — "a CLI não tem dependência nenhuma" — some, e some **por
  decisão registrada**, não porque alguém o afrouxou para o CI passar
- A biblioteca de prompt interativo que o `go.mod` do `aru` mencionava desde a
  Fase 1 continua sem entrar: não há prompt interativo, e não haverá enquanto
  flag e arquivo resolverem
