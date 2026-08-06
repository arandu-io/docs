# 19 — DSL e geração assistida

O diferencial do produto estava em um parágrafo do roadmap. Este documento o
especifica.

## O princípio, repetido porque tudo depende dele

**O modelo nunca escreve Go.** O modelo produz uma spec; um gerador determinístico
produz o código. A corretude não depende de o modelo ter sido treinado no Arandu —
depende do nosso compilador.

É isso que responde ao problema de fundo: um framework novo tem zero exemplos no
treino, e o modelo preenche a lacuna com padrões de Laravel e Django que não
existem aqui. Com spec intermediária, isso deixa de importar: o schema cabe no
contexto e é verificável antes de virar código.

## Formato: YAML com JSON Schema publicado

| Candidato | Decisão |
|---|---|
| YAML + JSON Schema | **escolhido** — massa de treino enorme, validação madura, schema publicável |
| HCL | rejeitado — menos treino, e a expressividade extra é exatamente o que não queremos |
| CUE | rejeitado — validação superior, público minúsculo, LLM erra |

O schema é o produto real dessa camada: é ele que faz um modelo sem fine-tuning
nenhum acertar a spec. Ele não é repositório novo — sai do comando `aru schema`,
gerado das mesmas constantes que o validador usa, com exemplo que passa na
própria validação.

## O que a DSL expressa

Conjunto **fechado**, por decisão:

- módulo e a entidade dele
- campos com tipos fechados: `string` `text` `int` `decimal` `money` `bool` `date`
  `timestamp` `uuid` `email`
- por campo, `required` e `unique` — de onde saem as validações geradas
- escopo de tenant, que vem do `Grant` e entra em toda query
- permissões por ação e papel, que abrem a `Policy`

## Considerado e não entrou

Duas coisas estavam na lista original e ficaram de fora, pela REGRA 15 — cada
item que entra na DSL é mantido para sempre:

- **relações `belongs_to` e `has_many`** — só fazem sentido dentro do agregado, e
  o `data.Repository` já é por agregado. Relação na DSL viraria join gerado, e
  join gerado é o começo de um ORM
- **telas `list`, `form` e `show`** — gerariam view, e a forma canônica de uma
  tela ainda não existe. O `make:auth` publica a primeira, e é dela que se
  aprende o que uma tela gerada deve parecer. Gerar antes disso é adivinhar

Ambas voltam à mesa quando houver um módulo real pedindo. Ver `03-roadmap-fases.md`.

## O que a DSL não expressa, e por quê

- lógica de negócio condicional
- SQL customizado
- template livre — toda tela se escreve à mão em `.kyse.go`, com as diretivas do
  `kyse`

Este é o ponto que decide se o produto vive. DSL expressiva demais vira uma
linguagem que você mantém para sempre; restrita demais e o cliente esbarra no
limite no primeiro requisito fora do padrão e abandona.

A saída não é ampliar a DSL, é a **saída de emergência tipada**: o código gerado
tem blocos protegidos por marcador, e regeneração preserva o que está dentro.

```go
// arandu:begin custom
func (s *Service) applyProgressiveDiscount(...) { ... }
// arandu:end custom
```

O que está entre os dois marcadores é preservado na regeração.

Fora do padrão se escreve em Go, com o framework inteiro disponível, e sem perder
a capacidade de regenerar o resto.

## Onde o gerador escreve

Nos diretórios do Laravel, os mesmos que `aru new` já criou. Não existe
`modules/<nome>/`.

| Item da spec | Arquivo emitido |
|---|---|
| entidade | `app/Models/Pedido.go` |
| permissões por ação e papel | `app/Policies/PedidoPolicy.go` |
| acesso a dado | `app/Repositories/PedidoRepository.go` |
| regra de negócio | `app/Services/PedidoService.go`, com os blocos custom |
| validações | `app/Http/Requests/StorePedido.go` |
| ações HTTP | `app/Http/Controllers/PedidoController.go` |
| rotas | linhas em `routes/web.go` |
| campos e tipos | migração em `database/migrations/` e as queries do sqlc |
| o que o módulo declara sobre si | `arandu.mod.toml`, conferido pelo `doctor` |
| teste | o teste do módulo, verde sem edição manual |

**O gerador não escreve view.** Telas ficaram fora da DSL, e o `.kyse.go` se
escreve à mão. O `.go` ao lado dele sai do `aru view:build`, nunca daqui.

O dado que chega a uma view é struct tipada por página, nunca `map` — campo
errado é erro de compilação, e uma regra do `aru doctor` proíbe `map` como dado
de view.

## O pipeline

```
usuário (português)
  → LLM  → spec YAML
  → validação contra o schema        ← erro de spec morre aqui
  → diff apresentado ao humano       ← aprovação explícita
  → gerador determinístico → Go
  → aru doctor + go test
```

Erro de spec nunca vira código quebrado, porque a validação acontece antes da
geração. E o humano aprova a spec, não o código — spec cabe na cabeça, diff
espalhado por `app/`, `database/` e `routes/` não.

## Round-trip

Código gerado → DSL → código idêntico byte a byte, verificado por golden files na
CI. Sem isso não há como garantir que regenerar não destrói trabalho, e sem essa
garantia ninguém regenera — o gerador vira ferramenta de uso único.

## Por que isso não é o Spec Kit

Spec Kit, Kiro e Tessl são agnósticos de linguagem, e é aí que são fracos: spec
genérica gera código genérico, plausível e não verificável. Nossa spec compila
contra um runtime rígido que sabe que `Policy` roda antes de `Repository`, que
toda entidade tem validador e que o tenant vem do `Grant`.

Eles não podem fazer isso porque não são donos de runtime nenhum. Somos.
