# De-para: Laravel → Arandu

A curva de aprendizado é o produto. Um desenvolvedor Laravel deve reconhecer o
vocabulário no primeiro contato e só precisar aprender **o que muda de fato**.

## O fluxo mental: Laravel + Inertia, sem o JavaScript

O modelo mental correto não é o page reload inteiro nem a SPA. É **Inertia**.

O Inertia existe para dar sensação de SPA mantendo roteamento e controller no
servidor: sem camada de API, sem router no cliente, sem estado duplicado. É
exatamente o que o HTMX faz — trocando props em JSON por fragmento de HTML.

```
Laravel + Inertia            Arandu
─────────────────────        ─────────────────────
routes/web.php           →   routes/web.go
Route::get(...)          →   r.Action("GET", ...)
Controller               →   Controller em app/Http/Controllers
Inertia::render(         →   ctx.View(
  'Clientes/Index',            "clientes.index",
  ['clientes' => $c]           views.ClientesIndexData{Clientes: c},
)                            )
props (array)            →   struct tipada da página, nunca map
partial reload           →   ctx.Fragment(status, nome, dados) + hx-get
useForm().errors         →   validation.Errors no fragmento do form
shared data              →   valores no context
```

O controller renderiza pelo **nome** da view — `"clientes.index"`, o mesmo ponto
que o Laravel usa —, e o `http.Context` resolve contra o registro. Não existe
função exportada por página que o controller chame direto.

*Corrigido em 07/08/2026:* este bloco dizia
`views.ClientesIndex(...).Render(ctx, w)`. Essa era a forma do `templ`, que a ADR
0020 descartou junto com o motor. O que existe em `framework/http/context.go` é
`func (c *Context) View(name string, data any) error` e o par
`func (c *Context) Fragment(status int, name string, data any) error`; o
esqueleto chama `ctx.View("home", views.HomeData{...})` e o `aru make:module`
gera `ctx.View("purchase-orders.index", ...)`.

A diferença que importa: **não existe build de JavaScript.** Sem Vite, sem npm,
sem hidratação, sem estado no cliente. O HTMX é atributo em HTML — quem sabe HTML
já sabe usar.

## Vocabulário que atravessa igual

### Estrutura e vocabulário do servidor

| Laravel | Arandu | Muda algo? |
|---|---|---|
| `app/Http/Controllers/HomeController.php` | `app/Http/Controllers/HomeController.go` | mesmo nome, mesmo lugar; aqui não pode importar `data` |
| `app/Http/Middleware/` | `app/Http/Middleware/`, `func(http.Handler) http.Handler` | assinatura padrão do Go |
| `app/Http/Requests/` (FormRequest) | `app/Http/Requests/`, struct + `Validate() Errors` | tipado, sem tags nem reflection |
| `app/Models/User.php` | `app/Models/User.go` | entidade é dado; a regra mora em `app/Services/` |
| `app/Policies/` | `app/Policies/`, `Policy[T]` + `Grant` | **aqui é obrigatório, não opcional** |
| `app/Providers/` | `app/Providers/` | sem container: o provider registra a fiação explícita |
| `routes/web.php` | `routes/web.go` | **igual** — central, uma só, não por módulo |
| `Route::get('/x', [C::class,'index'])` | `r.Action("GET", "/x", d.C.Index)` | **`r.Get` existe e não serve**: ele toma `http.HandlerFunc`; action de controller é `func(*Context) error` e registra por `Action` |
| `routes/console.php` | `routes/console.go` | igual |
| `config/*.php` (dez arquivos) | `config/*.go` (os mesmos dez) | struct tipada; erro de chave vira erro de compilação |
| `database/migrations/` | `database/migrations/` | `aru migrate` é passo de pipeline, nunca boot. **E o arquivo não é descoberto**: toda migration entra à mão em `All()`, e esquecer não falha em lugar nenhum — o schema simplesmente não muda |
| `database/seeders/` | `database/seeders/` | igual |
| `database/factories/` | `database/factories/` | mesmo lugar, outro contrato: **não há `->create()`** (a factory não toca o banco), não há faker (os valores são contadores), `->count()` é `MakeMany` e não há `->state()` |
| Event / Listener | `app/Events/` / `app/Listeners/` | **muda tudo**: não existe dispatcher em processo. O evento é linha no outbox dentro da transação, entregue pelo relay cerca de um segundo depois |
| Job / Queue | `app/Jobs/` | a fila padrão é a tabela `jobs` da própria aplicação; RESP no submódulo, NATS só com demanda |
| `artisan` | `aru` | — |
| — | `app/Services/` e `app/Repositories/` | os 20% que o Laravel só recomendava |

### Comando: o que o dev Laravel digita de memória

O alvo do produto é quem **migra** uma aplicação Laravel existente. Numa migração
não se gera o módulo inteiro: porta-se um controller, depois uma migration,
depois um middleware. `php artisan make:controller` é memória muscular, e a
tabela abaixo é o que ela encontra do outro lado.

Enquanto uma linha desta tabela não tiver comando, ela é uma promessa em aberto.
O `checklist.sh` verifica a tabela contra `aru help`.

| `artisan` | `aru` | nota |
|---|---|---|
| `make:controller` | `make:controller` | `--resource` gera as sete ações, por `http.Resource` |
| `make:model` | `make:model` | **o Model é struct de dados, não Active Record** — quem toca o banco é o Repository, e ele exige `security.Grant` (doc 01, REGRA 17). É a diferença que mais surpreende quem migra |
| `make:migration` | `make:migration` | `--create=` e `--table=` como no artisan, mas **`--fields` é obrigatório** — `aru make:migration add_status --table=posts` sozinho é erro |
| `make:middleware` | `make:middleware` | — |
| `make:request` | `make:request` | validação tipada; erro de campo é erro de compilação |
| `make:seeder` | `make:seeder` | — |
| `make:factory` | `make:factory` | — |
| `make:enum` | `make:enum` | `--values` é obrigatório: em Go o enum é exatamente o boilerplate, e um tipo sozinho não emite linha útil |
| `make:policy` | `make:policy` | já existe. Nasce negando tudo, e é isso que faz o padrão ser seguro |
| `make:job` | `make:job` | fila é a tabela `jobs` da própria aplicação |
| `make:event` | `make:event` | o evento vai no outbox, na mesma transação da escrita (doc 27) |
| `make:listener` | `make:listener` | — |
| `make:command` | `make:command` | — |
| `make:mail` | `make:mail` | existe desde que o `framework/mail` entrou no core. Escreve o mailable e as duas partes da mensagem. Esta linha dizia que não existia |
| `composer require laravel/ui --dev` | — | `go run <módulo>@latest` roda um módulo publicado sem tocar no `go.mod` |
| `ui bootstrap --auth` | `go run github.com/arandu-io/ui@latest auth` | as mesmas nove telas, HTMX + Alpine + Tailwind v4 no lugar do Bootstrap (ADR 0022). Sem preset: há uma pilha só |
| `composer remove laravel/ui` | — | não há o que remover |
| — | `make:module` | não tem equivalente: é o `make:model --all` com Policy, Service, Repository e as quatro telas de uma vez |
| `route:list` | `route:list` | — |
| `migrate` | `migrate` | — |
| `migrate:rollback` | `migrate:rollback` | — |
| `migrate:fresh` | `migrate:fresh` | — |
| `migrate:status` | `migrate:status` | — |
| `db:seed` | `db:seed` | o nome é **posicional**: `--class=` é recusado com erro, e é justamente o que a mão digita sozinha |
| `queue:work` | `queue:work` | mesmo nome, dois flags contra dezesseis: só `--queue` (uma fila, não lista) e `--workers` |
| `schedule:run` | `schedule:run` | **outro comando com o mesmo nome**: exige um id e roda **uma** tarefa. O laço de agendamento vive dentro do `aru serve` |
| `key:generate` | `key:generate` | a chave faz outro trabalho: aqui ela é **só** chave de HMAC (sessão, CSRF, links assinados). Não cifra nada — não existe `Crypt` |
| `serve` | `serve` | e `aru dev`, que recompila a view ao salvar |
| `tinker` | **não entra** | REPL para linguagem dinâmica; `docs/06-escopo.md` veta |
| `optimize`, `config:cache` | **não entram** | resolvem custo de boot do PHP. Aqui a config é struct |
| `view:clear` | **falta** | existe sim um diretório de view compilada que envelhece, em `storage/framework/views/`, no mesmo lugar do Laravel. Apagá-lo é o remédio de folclore que o `aru view:build` já deveria tornar desnecessário |
| `storage:link` | **não entra** | o disco público é servido pelo binário, não por link simbólico |
| — | `doctor` | não tem equivalente: é o que faz a arquitetura ser garantida pelo compilador |

### View: as diretivas são as do Blade

O motor de template é o `kyse`, nosso. As diretivas mantêm o nome do Blade de
propósito — não há vocabulário novo a aprender aqui.

| Laravel | Arandu | Muda algo? |
|---|---|---|
| `resources/views/home.blade.php` | `resources/views/home.kyse.go` | o gerado é `storage/framework/views/home.go` — **sob `storage/`, não ao lado**, e a árvore é espelhada: `auth/login.kyse.go` gera `auth/login.go`, sem achatar |
| `@extends('layouts.app')` | `@extends('layouts.app')` | idêntico |
| `@section` / `@endsection` | `@section` / `@endsection` | idêntico |
| `@yield('content')` | `@yield('content')` | idêntico |
| `@include('partials.header')` | `@include('partials.header')` | idêntico |
| `@if` `@elseif` `@else` `@endif` | iguais | idêntico |
| `@foreach` / `@endforeach` | `@foreach` / `@endforeach` | idêntico |
| `@forelse` / `@empty` / `@endforelse` | iguais | idêntico |
| `@for` / `@endfor` | `@for(i := 0; i < n; i++)` | as três cláusulas em Go |
| `@while` / `@endwhile` | iguais | Go tem um laço só, e é `for` |
| `@continue` / `@break` | iguais | **sem condição**: o Blade aceita `@break($i > 3)`, aqui só a forma nua. A diferença é silenciosa |
| `@csrf` | `@csrf` | idêntico |
| `{{ $var }}` | `{{ .Var }}` | escapado nos dois; o **ponto** é o que faz virar campo do dado |
| `{!! $html !!}` | `{!! .Html !!}` | cru nos dois; mesma regra do ponto |
| `{{-- comentário --}}` | `{{-- comentário --}}` | não chega ao navegador, nos dois |
| `@php` / `@endphp` | `@go` / `@endgo` | só troca o nome da linguagem hospedeira |
| `storage/framework/views/` | **é exatamente onde o gerado fica** | mesmo endereço do Laravel e mesmo motivo: é saída de build, gitignorada. O que difere é que o nosso é Go legível com a árvore espelhada, não um hash |

A tabela acima é o conjunto **inteiro** de diretivas: sete de bloco (`section`,
`if`, `foreach`, `forelse`, `for`, `while`, `go`, cada uma com seu `end`) e nove
inline (`extends`, `yield`, `include`, `csrf`, `elseif`, `else`, `continue`,
`break`, `empty`). Não há décima sétima, e é isso que a REGRA 15 quer dizer: o
que não couber se escreve em Go, dentro de `@go`. As duas listas estão em
`aru/internal/kyse/kyse.go`, e `TestEveryDirectiveEmitsSomething` percorre
`Directives()` para garantir que cada membro está exercitado **ou**
explicitamente pulado. Onze são exercitados ali; doze — as diretivas de
fechamento, mais `@go`, `@section`, `@extends`, `@else` e `@elseif` — têm testes
próprios e entram no mapa `skip`.

**A regra do ponto.** Só expressão iniciada por ponto vira campo do dado da
página: `{{ .Name }}` gera `view.Text(d.Name)`, onde `d` é o dado já convertido
para o tipo da página. O resto passa direto para o Go — é o que faz
`{{ len(d.Items) }}`, `{{ view.URL("app.css") }}` e `{{ item }}` (a variável
ligada por um `@foreach`) funcionarem. A consequência é que `{{ name }}`, sem
ponto, gera `view.Text(name)` e o build para em `undefined: name`. É a diferença
com o `$` do Blade. Não é a única do capítulo — a tabela acima lista as
outras, e a coluna "Muda algo?" é onde elas estão.

## O que o dev Laravel precisa desaprender

| Laravel | Aqui | Por quê |
|---|---|---|
| Service container | fiação explícita em `bootstrap/app.go`, escrita à mão | reflection move erro para runtime |
| Facades | passagem por parâmetro ou context | exigiria estado global mutável |
| Eloquent / Active Record | entity + repository + SQL parametrizado | SQL previsível, injection impossível |
| `$fillable` | struct de request explícita | a classe de bug não existe |
| `compact()` / `with([...])` na view | struct tipada por página | campo errado vira erro de compilação; `map` é reprovado pelo `doctor` |
| Model com regra dentro | service | entidade é dado, não comportamento |

## A frase de onboarding

> Se você conhece Laravel, você já sabe 80% do Arandu. Os outros 20% são o
> compilador cobrando o que o Laravel só recomendava.

## O que isso exige da documentação (Fase D)

Um guia "Vindo do Laravel" com esta tabela na primeira página. É o documento de
maior conversão possível para o público-alvo, e ele já está escrito aqui —
falta só virar prosa quando a API congelar.
