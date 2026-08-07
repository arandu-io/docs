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
Route::get(...)          →   r.Get(...)
Controller               →   Controller em app/Http/Controllers
Inertia::render(         →   views.ClientesIndex(
  'Clientes/Index',            views.ClientesIndexData{Clientes: c},
  ['clientes' => $c]         ).Render(ctx, w)
)
props (array)            →   struct tipada da página, nunca map
partial reload           →   hx-get + fragmento
useForm().errors         →   validation.Errors no fragmento do form
shared data              →   valores no context
```

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
| `Route::get('/x', ...)` | `r.Get("/x", ...)` | mesma leitura, em Go |
| `routes/console.php` | `routes/console.go` | igual |
| `config/*.php` (dez arquivos) | `config/*.go` (os mesmos dez) | struct tipada; erro de chave vira erro de compilação |
| `database/migrations/` | `database/migrations/` | igual; `aru migrate` é passo de pipeline, nunca boot |
| `database/seeders/` | `database/seeders/` | igual |
| `database/factories/` | `database/factories/` | igual |
| Event / Listener | `app/Events/` / `app/Listeners/` | — |
| Job / Queue | `app/Jobs/` | a fila padrão é a tabela `jobs` da própria aplicação; RESP no submódulo, NATS só com demanda |
| `artisan` | `aru` | — |
| — | `app/Services/` e `app/Repositories/` | os 20% que o Laravel só recomendava |

### View: as diretivas são as do Blade

O motor de template é o `kyse`, nosso. As diretivas mantêm o nome do Blade de
propósito — não há vocabulário novo a aprender aqui.

| Laravel | Arandu | Muda algo? |
|---|---|---|
| `resources/views/home.blade.php` | `resources/views/home.kyse.go` | o gerado fica ao lado, em `home.go` |
| `@extends('layouts.app')` | `@extends('layouts.app')` | idêntico |
| `@section` / `@endsection` | `@section` / `@endsection` | idêntico |
| `@yield('content')` | `@yield('content')` | idêntico |
| `@include('partials.header')` | `@include('partials.header')` | idêntico |
| `@if` `@elseif` `@else` `@endif` | iguais | idêntico |
| `@foreach` / `@endforeach` | `@foreach` / `@endforeach` | idêntico |
| `@forelse` / `@empty` / `@endforelse` | iguais | idêntico |
| `@for` / `@endfor` | `@for(i := 0; i < n; i++)` | as três cláusulas em Go |
| `@while` / `@endwhile` | iguais | Go tem um laço só, e é `for` |
| `@continue` / `@break` | iguais | idêntico |
| `@csrf` | `@csrf` | idêntico |
| `{{ $var }}` | `{{ var }}` | escapado nos dois |
| `{!! $html !!}` | `{!! html !!}` | cru nos dois |
| `{{-- comentário --}}` | `{{-- comentário --}}` | não chega ao navegador, nos dois |
| `@php` / `@endphp` | `@go` / `@endgo` | só troca o nome da linguagem hospedeira |
| `storage/framework/views/` (cache) | **não existe** | o gerado é o `.go` ao lado, legível e versionável |

## O que o dev Laravel precisa desaprender

| Laravel | Aqui | Por quê |
|---|---|---|
| Service container | wiring explícito, gerado | reflection move erro para runtime |
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
