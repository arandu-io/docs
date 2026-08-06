# 14 — View e assets (`kyse`)

Este é o documento da camada de view. Ele fecha três coisas: qual é o motor de
template, de onde vem o CSS, e onde a view mora agora que o `porang` deixou de
existir.

O templ caiu (ADR 0020). O Tailwind caiu (ADR 0022). O `porang` foi dissolvido
(ADR 0021). O que não mudou é a restrição que decide tudo — e é por ela que se
começa.

## A restrição que decide tudo

Um projeto Arandu tem que continuar sendo:

```
git clone && aru dev
```

Sem `node_modules`, sem `package.json`, sem lockfile de JavaScript, sem Node
instalado na máquina. Se o desenvolvedor precisar de Node para rodar um projeto,
a frase de venda morre e a REGRA 9 vira decoração.

Isso não proíbe ter etapa de build. Proíbe que a etapa de build seja Node.

## A exceção que existia, e caiu

`09-uma-forma-so.md` proíbe: *Build → `go build`; nunca npm, Vite, bundler, Node
em qualquer lugar do core.*

`07-pacotes.md` abria: *única exceção onde npm importa: `porang`, se distribuir
plugin de Tailwind (o Basecoat é distribuído por npm).*

As duas não podiam valer ao mesmo tempo. E a exceção era precisamente o caminho
pelo qual o Node entrou no Laravel — a página de erro com `package.json` e
`vite.config.js` dentro do framework não nasceu de uma decisão de arquitetura,
nasceu de uma exceção pontual como essa.

**A exceção cai. A regra fica.** Com o Tailwind fora e o Bootstrap vendorizado,
ela deixa de ter até motivo de existir.

## O motor: `kyse`

`kyse` é o motor de template do Arandu — o lugar que o Blade ocupa no Laravel.
Vem do guarani *kysê*, faca. O acento cai no nome de arquivo e no import path,
por restrição do Go.

Ele substitui o templ. O motivo não é o templ ser ruim: é que a REGRA 10 diz que
o modelo mental é Laravel, e o dev Laravel reconhece `@extends` e `@section`
imediatamente. Motor próprio também remove uma dependência de terceiro do caminho
crítico de build.

### A extensão termina em `.go`, e isso foi verificado

O `.blade.php` termina em `.php` porque o PHP é a linguagem hospedeira. O
`.kyse.go` termina em `.go` pelo mesmo motivo. **A extensão do fonte é
`.kyse.go`, não `.kyse`.** O arquivo abre com a tag de build que o exclui do
compilador:

```go
//go:build kyse

package views

@extends('layouts.app')
```

Testado com Go 1.26.5 e cache limpo: `go build ./...`, `go vet ./...`,
`go test ./...` e `go run` passam, com o arquivo gerado no mesmo diretório e o
pacote importável.

Dois detalhes que a verificação expôs, e que valem estar escritos:

- o `gofmt` é o único que ignora build tag. O CI o chama filtrando:
  `gofmt -l $(find . -name '*.go' -not -name '*.kyse.go')`
- a cláusula `package` é obrigatória quando existe outro `.go` no diretório. Sem
  ela funciona apenas enquanto o diretório estiver sozinho, e quebra depois

### O par fonte/gerado fica lado a lado

```
resources/views/home.kyse.go   ← escreve-se
resources/views/home.go        ← o aru view:build gera, mesmo pacote
```

O par existe também no Laravel. Lá o gerado fica escondido em
`storage/framework/views/`; aqui ele fica visível ao lado, no mesmo diretório e
no mesmo pacote — e é código Go de verdade, checado pelo compilador.

**Não existe `storage/framework/views/`.** Cache de template compilado é
necessidade de linguagem interpretada. Nós compilamos uma vez, e o resultado é
binário.

### Diretivas — os nomes do Blade

```
@extends('layouts.app')
@section('content') … @endsection
@yield('content')
@include('partials.header')
@if · @elseif · @else · @endif
@foreach · @endforeach
@csrf
{{ variavel }}       escapado
{!! html !!}         cru
@go … @endgo         o @php, com Go de verdade dentro
```

A sintaxe do bloco de código é **`@go` … `@endgo`**, nunca `{{ go }}`. Mesmos
nomes do Blade no resto, mesma sintaxe. A única adição é o `@go`, que é o `@php`
com o nome da linguagem hospedeira trocado — mesmo raciocínio do Blade.

O conjunto é fechado, pela REGRA 15. Caso fora do padrão se escreve em Go, dentro
de `@go`, não numa diretiva nova.

## O dado é struct tipada por página, nunca `map`

Esta é a linha que sustenta a tese do produto na camada de view.

Cada página tem um struct de dados próprio. O handler monta o struct, o `kyse`
gera uma função que recebe esse struct, e campo errado no template é **erro de
compilação** — não página em branco em produção.

`map[string]any` é o oposto disso: erro de digitação vira `nil` silencioso, e o
compilador não tem o que checar. É exatamente a classe de bug que o Go permitiria
evitar e o Blade não evita.

Uma regra de `aru doctor` reprova `map` como dado de view.

> Não é o desenvolvedor que garante a arquitetura. É o compilador. Numa view
> `map`, o compilador não garante nada.

## Toolchain: nenhum binário baixado

O `kyse` é compilado pela própria `aru`, via `aru view:build`. Não há binário de
terceiro para baixar, nem versão para fixar em `arandu.toml`, nem cache em
`~/.arandu/bin`.

Isso é uma simplificação em relação ao desenho anterior, que baixava dois
binários — `templ` e o `tailwindcss` standalone. Ambos saíram:

| Saiu | Entrou |
|---|---|
| `templ` (binário baixado) | `aru view:build`, dentro da CLI |
| `tailwindcss` standalone (binário baixado) | Bootstrap vendorizado, sem etapa de compilação de CSS |

`aru dev` roda o watch do `kyse` mais o servidor. Um comando, um terminal.

## CSS: Bootstrap, vendorizado

O CSS é **Bootstrap**, como no `laravel/ui`. O Tailwind caiu, e com ele o binário
standalone e a necessidade de compilar CSS no build.

O Bootstrap entra vendorizado no repositório do `arandu`, com o arquivo de
licença e a atribuição, e é embutido por `go:embed`. npm não entra no grafo de
dependências em nenhuma hipótese.

Com o Tailwind fora, o Basecoat sai junto — ele é um plugin de Tailwind, e era o
único motivo pelo qual a exceção de npm tinha sido aberta.

Se manter o vendor ficar caro, a alternativa é escrever o CSS — não é
reintroduzir npm.

## Assets: `go:embed`, sem CDN

`htmx.min.js`, `alpine.min.js` e o CSS do Bootstrap entram no binário via
`go:embed` e são servidos por `framework/view` em `/_arandu/assets/<hash>/…` com
cache imutável.

Sem CDN, e não por preferência: o middleware `SecurityHeaders` define
`script-src 'self'`. Servir HTMX pelo unpkg exigiria afrouxar a CSP — pagar
segurança por conveniência. É uma diferença observável em relação ao Gojang, que
carrega htmx por unpkg.

Consequência que vale explicitar: **o deploy continua sendo um binário.** Sem
passo de publicação de assets, sem CDN para invalidar, sem `storage:link`.

## Alpine: o limite, escrito

A REGRA 9 diz *"Alpine só onde o HTMX não alcança"*. Sem um teste, isso é opinião,
e opinião não segura ninguém em revisão de PR.

**Alpine é permitido quando o estado é, ao mesmo tempo:**

1. puramente do cliente,
2. efêmero — morre no reload sem perda,
3. invisível para o servidor.

Permitido: dropdown aberto/fechado, aba selecionada, foco, máscara de input,
confirmação antes de submeter.

Proibido: qualquer coisa que vá ao servidor, qualquer coisa persistida, qualquer
coisa derivada de dado que o servidor já tem.

O teste em uma frase: **se você está escrevendo regra de negócio em Alpine, o
componente devia ser um fragmento HTMX.**

Verificação no `aru doctor`: reprova `x-data` que contenha `fetch`, `axios`,
`$store` global, ou `$watch` com efeito de rede.

## Bug conhecido: o Alpine embutido não roda sob a nossa CSP

O `alpine.min.js` embutido hoje é a **build padrão**. Ela avalia expressão de
diretiva construindo função em runtime, com
`Object.getPrototypeOf(async function(){}).constructor` — o equivalente a `eval`.
Isso exige `unsafe-eval` na `script-src`.

A nossa CSP é `script-src 'self'`, sem `unsafe-eval`, e não vai deixar de ser:
afrouxar a CSP para acomodar uma build de JavaScript é exatamente a troca que
este documento recusa em todos os outros pontos.

**Consequência hoje: toda diretiva `x-data` está quebrada no browser.** O
navegador bloqueia a avaliação e o componente não inicializa.

A solução é a build `@alpinejs/csp`, publicada pelo próprio Alpine para este
caso. Ela não avalia string: as expressões deixam de ser escritas inline na
diretiva e passam a ser métodos e propriedades declarados em objeto de
componente. Ou seja, **muda como se escreve Alpine** — não é troca de arquivo e
pronto, e o código de exemplo existente precisa ser reescrito.

Duas coisas ficam registradas: o embed passa a ser a build CSP, e a CSP não muda.
A build CSP é vendorizada com licença, como o resto — não vem por npm.

## Componente é view `kyse` tipada

Não existe `<nx-modal />`, `<x-card>` nem web component registrado em runtime.

Componente é uma view em `resources/views/components/`, incluída por `@include`,
que o `kyse` resolve em tempo de build para uma chamada da função Go gerada, com
o struct tipado daquela view:

```go
//go:build kyse

package views

@include('components.alert')
```

Tag customizada seria um segundo sistema de componentes ao lado do `kyse`, com
resolução em runtime e erro só na tela. O `kyse` já dá composição, tipagem e erro
em tempo de compilação.

## Onde a view mora — o `porang` foi dissolvido

O `porang` era o décimo primeiro repositório. Ele não existe mais, e a razão é a
mesma que decidiu a estrutura de diretórios: **no Laravel, a view mora no
esqueleto, não no framework.**

| O que era do `porang` | Onde está |
|---|---|
| views e componentes | `arandu/resources/views/` |
| runtime de render e assets embutidos | `framework/view/` |

Dez repositórios, não onze. A REGRA 3 continua valendo.

## O que a camada de view não toca

A página de erro. Ela continua `const` de `html/template` em
`observability/errorpage`, sem depender de `framework/view`, sem CSS compilado e
sem asset embutido.

O motivo ficou mais forte com o `kyse`: a página de erro precisa renderizar
quando o resto está quebrado — inclusive **quando o build do `kyse` falhou**. Uma
página de erro escrita em `kyse` não renderizaria exatamente no momento em que
mais importa.

## Modo API

`aru new --api` gera o esqueleto sem `resources/views/` e sem os assets
embutidos. Nada mais muda: o handler decide se devolve fragmento ou JSON, e o
resto do contrato é idêntico.

## Licença

MIT (ADR 0008), no framework e no vendor. Qualquer menção a Apache-2.0 neste
documento ou em documento vizinho é regressão.
