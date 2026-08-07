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
- **telas na DSL** — não existe `screens:` na spec, e não vai existir. O que
  mudou é outra coisa: o gerador **passou a emitir as quatro telas do CRUD**
  (`index`, `show`, `create`, `edit`) a partir dos campos que a spec já declara,
  sem nenhuma chave nova. A forma canônica de uma tela deixou de ser hipótese
  quando o `make:auth` publicou as dela e o `layouts.app` do esqueleto fixou o
  contrato — uma interface `Layout` com `PageTitle`

A relação volta à mesa quando houver um módulo real pedindo. Ver
`03-roadmap-fases.md`.

## O que a DSL não expressa, e por quê

- lógica de negócio condicional
- SQL customizado
- template livre — a spec não descreve markup. As quatro telas geradas saem dos
  campos, e a partir daí são suas: editadas à mão em `.kyse.go`, com as
  diretivas do `kyse` e nenhuma outra

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
| validações | `app/Http/Requests/PedidoRequest.go`, com `StorePedido` e `UpdatePedido` |
| ações HTTP | `app/Http/Controllers/PedidoController.go`, as sete do `Resource` |
| telas | `resources/views/pedidos/{index,show,create,edit}.kyse.go` |
| campos e tipos | migração em `database/migrations/`, listada em `All()` |
| teste | `app/Http/Controllers/PedidoController_test.go`, verde sem edição manual |
| a spec | `database/specs/pedido.yaml`, commitada com o código |

**O gerador escreve a view fonte, nunca a gerada.** O `.kyse.go` sai daqui; o
`.go` ao lado dele sai do `aru view:build`, e editá-lo é trabalho perdido na
próxima compilação.

**O gerador não escreve fiação.** Três linhas ficam por conta de quem gerou — o
campo do `Deps` e o `Resource` em `routes/web.go`, o construtor dentro do literal
`routes.Deps` de `bootstrap/app.go`, e a migração no `All()` de
`database/migrations/migrations.go`. O comando imprime as três. Um
gerador que editasse esses arquivos sozinho tiraria deles a única propriedade
que justificam ter: dizer, por escrito, do que a aplicação é feita.

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

## Round-trip: pela spec, não pelo código

Spec → Go, e de novo Spec → Go, byte a byte igual. O ciclo fecha pela spec, não
pelo código: o `aru generate` grava a spec normalizada em
`database/specs/<módulo>.yaml` junto com o código que ela produziu, e regenerar
lê esse arquivo de volta. É o que `TestTheRoundTripIsByteForByte`, em
`aru/generate_test.go`, verifica — gera, regenera a partir da spec que a primeira
execução salvou, e compara os bytes de `app/`. Os golden files de
`aru/internal/gen/` prendem a outra ponta, e a CI roda os dois no
`go test -race ./...`.

Sem isso não há como garantir que regenerar não destrói trabalho, e sem essa
garantia ninguém regenera — o gerador vira ferramenta de uso único. O que protege
o trabalho escrito à mão é outra coisa: os blocos `arandu:begin custom`.

*Corrigido em 07/08/2026:* esta seção dizia "código gerado → DSL → código
idêntico byte a byte", que promete um caminho de volta, de Go para spec. Ele não
existe. Em `aru/internal/spec/` há `Read`, `Parse`, `Write` e `Marshal` — YAML
para struct e struct para YAML, e nada mais; nenhuma função lê Go gerado e
devolve spec. O `go/parser` aparece na `aru` em `makerequest.go`,
`makemigration.go`, `internal/gen/factory.go` e no `doctor`, sempre para ler um
modelo que já existe ou auditar código, nunca para reconstruir a spec. A ida é
determinística e verificada; a volta nunca foi construída, e não está prometida
aqui. O comentário de doc de `TestTheRoundTripIsByteForByte` ainda repete a
formulação antiga; o corpo do teste faz o que esta seção descreve.

## Por que isso não é o Spec Kit

Spec Kit, Kiro e Tessl são agnósticos de linguagem, e é aí que são fracos: spec
genérica gera código genérico, plausível e não verificável. Nossa spec compila
contra um runtime rígido que sabe que `Policy` roda antes de `Repository`, que
toda entidade tem validador e que o tenant vem do `Grant`.

Eles não podem fazer isso porque não são donos de runtime nenhum. Somos.
