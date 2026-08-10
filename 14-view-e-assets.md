# 14 — View e assets (`kyse`)

Este é o documento da camada de view. Ele fecha três coisas: qual é o motor de
template, de onde vem o CSS, e onde a view mora agora que o `porang` deixou de
existir.

O templ caiu (ADR 0020). O `porang` foi dissolvido
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

**A exceção cai. A regra fica.** O Tailwind entra pelo binário standalone, que
não é npm — e é o que fecha a porta que a exceção tinha aberto.

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

### O par fonte/gerado fica no mesmo pacote — lado a lado só na raiz

```
resources/views/home.kyse.go         ← escreve-se
resources/views/home.go              ← o aru view:build gera, ao lado

resources/views/auth/login.kyse.go   ← escreve-se
resources/views/auth_login.go        ← o gerado, achatado um diretório acima
```

*Corrigido em 07/08/2026:* esta seção dizia "lado a lado", sem ressalva. Vale
para a view na raiz de `resources/views/` e não vale para a aninhada. **Go tem
um pacote por diretório**: com o gerado dentro de `resources/views/auth/`, esse
diretório seria outro pacote, e a struct que o layout declara não estaria
visível para a página que o estende — que é exatamente o que faz uma página
herdar o dado do layout. Então o fonte mantém a árvore do Laravel e o gerado vai
para a raiz de `resources/views/`, com o caminho no nome e `/` virando `_`. Quem
faz isso é `OutputPath`, em `aru/internal/kyse/kyse.go`. Ver ADR 0020, seção
"o gerado é flat".

Achatar o fonte também seria a outra saída, e custaria o que a árvore existe
para dar: o desenvolvedor Laravel abre `resources/views/auth/` e encontra
`login`, `register` e `passwords/` onde espera. O nome de render não muda com o
achatamento — `auth/login.kyse.go` continua sendo `ctx.View("auth.login", …)`,
que é como o Laravel a chama.

O par existe também no Laravel. Lá o gerado fica escondido em
`storage/framework/views/`; aqui ele fica visível, no mesmo pacote — e é código
Go de verdade, checado pelo compilador.

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
@forelse · @empty · @endforelse
@for · @endfor
@while · @endwhile
@continue · @break
@csrf
{{ .Campo }}         escapado
{!! .Html !!}        cru
{{-- nota --}}       não chega ao navegador
@go … @endgo         o @php, com Go de verdade dentro
```

São vinte e três nomes, e são exatamente estes: o conjunto vive em dois mapas de
`aru/internal/kyse/kyse.go` — `blockDirectives`, com os sete que abrem região
(`section`, `if`, `foreach`, `forelse`, `for`, `while`, `go`) mais os sete
`@end…` correspondentes, e `inlineDirectives`, com os nove que emitem no lugar
(`extends`, `yield`, `include`, `csrf`, `elseif`, `else`, `continue`, `break`,
`empty`). A função `Directives()` devolve os vinte e três, e
`TestEveryDirectiveEmitsSomething` percorre a lista inteira: onze ele exercita
com uma view mínima e confere o traço no Go gerado (`if`, `foreach`, `for`,
`while`, `yield`, `include`, `csrf`, `forelse`, `continue`, `break`, `empty`), e
doze ele pula por nome — `section`, `extends`, `else`, `elseif`, `go` e os sete
`@end…` —, que têm teste próprio. Uma diretiva que entre no conjunto sem estar
em nenhuma das duas listas reprova o teste. Três vezes uma diretiva foi
declarada sem caso no gerador e o nó sumiu em silêncio, com build verde.

**A interpolação começa com ponto.** `{{ .Campo }}` vira `d.Campo`, onde `d` é o
struct daquela página; expressão sem ponto é passada adiante como Go, então
`{{ variavel }}` não é o campo `Variavel` — é um identificador solto, e o build
quebra. Chamada de método vale igual: `{{ .PageTitle() }}`.

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

### `Layout` e `Page`: como uma página tipada cabe num layout só

Se cada página tem o próprio struct, o layout não pode ser tipado pelo struct de
nenhuma delas. É por isso que existe um par de nomes que este documento não
mencionava e sem o qual a camada não se explica:

- **`Layout`** é a interface, declarada no `@go` de
  `arandu/resources/views/layouts/app.kyse.go`. Ela pede exatamente o que o
  layout desenha, e nada além: `PageTitle`, `BrandName`, `CSRFToken`,
  `SignedIn`, `SignedInName`, `HomeLink`, `LoginLink`, `LogoutLink` e
  `RegisterLink`.
- **`Page`** é o struct embutível, em `arandu/resources/views/page.go` — Go
  comum, sem tag de build, porque não renderiza nada. Ele carrega o título, a
  marca, o token de CSRF e a navegação, e implementa `Layout`.

A página embute `Page` e ganha o contrato de graça:

```go
type InvoicesIndexData struct {
	Page                     // a moldura
	Invoices   []InvoiceRow  // o que só esta página tem
	NextCursor string
}

var _ Layout = InvoicesIndexData{}
```

É o que faz **toda** tela que o `aru make:module` gera — index, show, create e
edit — renderizar pelo layout que o `aru make:auth` instala, sem que os dois
comandos saibam um do outro: os templates em `aru/internal/gen/viewtemplates.go`
embutem `Page` e afirmam `var _ Layout` em tempo de compilação, e o
`aru/internal/gen/authviews.go` reescreve `page.go` e redeclara a mesma
interface junto do layout novo — porque o `make:auth` roda em projeto gerado por
qualquer versão do `aru new`, inclusive uma anterior a esse contrato.

O `kyse` fecha o circuito na hora de escolher o tipo de render: `RenderType`, em
`aru/internal/kyse/kyse.go`, devolve **a interface** para um layout — o arquivo
que contém `@yield` — e **o struct** para uma página. Ler o primeiro tipo do
`@go` para os dois casos foi um defeito real: o layout do `make:auth` ficou
tipado pela struct da tela de login, e toda página do `make:module` passou a
responder

```
view "layouts.app" takes AuthPage and got views.InvoicesIndexData
```

com `go build` verde, porque a discordância é uma asserção de tipo em tempo de
execução.

## Toolchain: um binário baixado, não dois

O `kyse` é compilado pela própria `aru`, via `aru view:build`: não há binário de
terceiro para baixar, nem versão para fixar em `arandu.toml`, nem cache para ele
em `~/.arandu/bin`.

O `tailwindcss` standalone **continua sendo baixado** — é o que a seção seguinte
descreve, e o que a REGRA 13 autoriza: binário único não é npm. O desenho
anterior baixava dois; hoje é um:

| Saiu | Entrou |
|---|---|
| `templ` (binário baixado) | `aru view:build`, dentro da CLI |
| `tailwindcss` standalone | **continua** — é binário único, não npm |

`aru dev` roda o watch do `kyse` mais o servidor. Um comando, um terminal.

## CSS: Tailwind v4, pelo binário standalone

O CSS é **Tailwind v4**. O `aru view:build` baixa o binário standalone,
verifica o SHA-256 contra o `sha256sums.txt` da release e cacheia em
`~/.arandu/bin`. Binário único não é npm: não há `package.json`, não há
`node_modules`, não há lockfile de JavaScript.

`resources/css/app.css` é a folha do projeto, e `assets/app.css` a saída
compilada, que entra no binário por `go:embed`. **A entrada entregue ao Tailwind
não é o arquivo: é montada pelo `aru` e passada por stdin** — o `@import`
absoluto do `app.css` do projeto mais um `@source` apontando para o diretório do
`kyse`, que é módulo importado (ADR 0027) e por isso mora fora da árvore que o
projeto sabe declarar. Sem essa linha, nenhuma classe que só os componentes
desenham chega à página, e nada falha. Ver ADR 0037.

O diretório vem do `go list -m`, com o `vendor/` do projeto quando ele existe.
Quando o módulo é exigido e o fonte não está na máquina, o comando **para e
diz** — folha de estilo sem as classes dos componentes é defeito que só aparece
na tela.

O que o `laravel/ui` empresta é a **estrutura**, não o CSS: as nove views de
`resources/views/`, com a mesma função e nos mesmos caminhos, escritas em
utilitário Tailwind no lugar das classes do Bootstrap. Ver ADR 0022.

Se manter o vendor de um componente ficar caro, a alternativa é escrever o CSS
— não é reintroduzir npm.

## Assets: `go:embed`, sem CDN

`htmx.min.js`, `alpine.min.js` e o `app.css` compilado entram no binário via
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

Verificação no `aru doctor`, regra `alpine-reaches-the-server`: ela varre
`x-data`, `x-init`, `x-effect`, `x-on:*` e `@*` — aspas simples ou duplas — e
reprova o que contenha `fetch(`, `axios`, `XMLHttpRequest`, `$store`,
`navigator.sendBeacon`, `EventSource(` ou `new WebSocket`. A varredura é só
dentro dos `*.kyse.go` de `resources/views/`; o `doctor` não abre arquivo `.js`
(doc 28, anexo 3).

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

Componente é uma view em `resources/views/components/`, incluída por `@include`:

```go
//go:build kyse

package views

@include('components.alert')
```

O diretório é convenção, e nada no código o trata: o `aru new` não o cria, e a
palavra `components` não aparece em `.go` nenhum da `aru`. No `framework/view`
ela aparece — na camada `components` do Tailwind, dentro do CSS —, o que não é
tratamento de diretório. Ele passa a existir quando alguém escreve o primeiro
componente.

*Corrigido em 07/08/2026:* esta seção dizia que o `kyse` resolve o `@include`
**em tempo de build**, numa chamada da função Go gerada com o struct daquela
view, e que o componente dá "erro em tempo de compilação". Não é o que o gerador
faz. `@include('components.alert')` emite
`view.Include(w, "components.alert", data)` (`aru/internal/kyse/generate.go`), e
`view.Include` procura o nome no registro de **páginas**, que cada página gerada
preenche no `init()` — `view.Register("home", renderHome)` —, repassando o dado
da página que incluiu (`framework/view/runtime.go`). São dois registros, e essa
distinção tem consequência: um layout gerado chama `view.RegisterLayout`, que
grava noutro mapa, lido só por `RenderInto`. Logo,
`@include('layouts.app')` devolve `unknownView` — layout não se inclui, se
estende. **Nome errado é erro de render**, não de
compilação: a função devolve `unknownView(name)`. O layout é resolvido do mesmo
jeito, por `RenderInto`, e o tipo é afirmado dentro da função gerada
(`data.(HomeData)`, com `view.WrongData` na discordância) — é o mesmo mecanismo
de runtime que a seção do `Layout` já descreve.

O que o compilador garante é o **campo**: a view renderiza com o struct que ela
declara, então `d.Title` inexistente não compila, e é essa metade que sustenta
a tese aqui.

Tag customizada continua fora, e o motivo é a REGRA 9 mais o dado: seria um
segundo sistema de componentes ao lado do `@include`, com registro próprio e sem
struct por view. O `kyse` já dá composição e dado tipado, com uma forma só de
escrever componente.

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

*Corrigido em 07/08/2026:* esta seção anunciava `aru new --api`, gerando o
esqueleto sem `resources/views/` e sem os assets embutidos. **Esse flag não
existe.** `aru new` declara um único flag, `--module` (`aru/new.go`), e o
comando clona o esqueleto inteiro, sempre igual. Não há gatilho nem fase
declarada para o modo API: o assunto está aberto, e entra por ADR se entrar.

O que existe hoje, e não depende de flag, é a metade da promessa que já vale: o
handler decide se devolve fragmento ou JSON, e o resto do contrato é idêntico —
não há `routes/api.go` separado (doc 28). Um projeto que só sirva JSON tem
`resources/views/` e os assets embutidos no binário sem usá-los.

## Licença

MIT (ADR 0008), no framework e no vendor. Qualquer menção a Apache-2.0 neste
documento ou em documento vizinho é regressão.
