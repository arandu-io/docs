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

O schema publicado em `arandu-io/schema` é o produto real dessa camada: é ele que
faz um modelo sem fine-tuning nenhum acertar a spec.

## O que a DSL expressa

Conjunto **fechado**, por decisão:

- módulo e entidades
- campos com tipos fechados: `string` `text` `int` `decimal` `money` `bool` `date`
  `timestamp` `uuid` `enum` `ref`
- relações `belongs_to` e `has_many` — sempre dentro do agregado
- permissões por ação e papel, que geram a `Policy`
- validações declarativas (obrigatório, faixa, tamanho, formato)
- telas `list`, `form`, `show` com colunas e campos
- jobs e tarefas agendadas

## O que a DSL não expressa, e por quê

- lógica de negócio condicional
- SQL customizado
- template livre

Este é o ponto que decide se o produto vive. DSL expressiva demais vira uma
linguagem que você mantém para sempre; restrita demais e o cliente esbarra no
limite no primeiro requisito fora do padrão e abandona.

A saída não é ampliar a DSL, é a **saída de emergência tipada**: o código gerado
tem blocos protegidos por marcador, e regeneração preserva o que está dentro.

```go
// arandu:inicio custom — preservado na regeração
func (s *Service) aplicarDescontoProgressivo(...) { ... }
// arandu:fim custom
```

Fora do padrão se escreve em Go, com o framework inteiro disponível, e sem perder
a capacidade de regenerar o resto.

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
geração. E o humano aprova a spec, não o código — spec cabe na cabeça, diff de
sete arquivos não.

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
