# De-para: Laravel → Arandu

A curva de aprendizado é o produto. Um desenvolvedor Laravel deve reconhecer o
vocabulário no primeiro contato e só precisar aprender **o que muda de fato**.

## O fluxo mental: Laravel + Inertia, sem o JavaScript

O modelo mental correto não é Blade nem SPA. É **Inertia**.

O Inertia existe para dar sensação de SPA mantendo roteamento e controller no
servidor: sem camada de API, sem router no cliente, sem estado duplicado. É
exatamente o que o HTMX faz — trocando props em JSON por fragmento de HTML.

```
Laravel + Inertia            Arandu
─────────────────────        ─────────────────────
Route::get(...)          →   r.Get(...)
Controller               →   handler
Inertia::render(         →   views.ClienteIndex(
  'Clientes/Index',            clientes,
  ['clientes' => $c]         ).Render(ctx, w)
)
props (array)            →   parâmetros tipados do componente
partial reload           →   hx-get + fragmento
useForm().errors         →   validation.Errors no fragmento do form
shared data              →   valores no context
```

A diferença que importa: **não existe build de JavaScript.** Sem Vite, sem npm,
sem hidratação, sem estado no cliente. O HTMX é atributo em HTML — quem sabe HTML
já sabe usar.

## Vocabulário que atravessa igual

| Laravel | Arandu | Muda algo? |
|---|---|---|
| Controller | handler | fino nos dois; aqui não pode importar `data` |
| Middleware | `func(http.Handler) http.Handler` | assinatura padrão do Go |
| FormRequest | struct + `Validate() Errors` | tipado, sem tags nem reflection |
| Policy | `Policy[T]` + `Grant` | **aqui é obrigatório, não opcional** |
| Migration | migração no módulo | coletada pelo kernel |
| `routes/web.php` | `Routes(r *httpx.Router)` no módulo | por módulo, não global |
| `config/*.php` | struct tipada | erro de chave vira erro de compilação |
| Event / Listener | evento / handler | — |
| Job / Queue | job | NATS |
| `artisan` | `aru` | — |

## O que o dev Laravel precisa desaprender

| Laravel | Aqui | Por quê |
|---|---|---|
| Service container | wiring explícito, gerado | reflection move erro para runtime |
| Facades | passagem por parâmetro ou context | exigiria estado global mutável |
| Eloquent / Active Record | entity + repository + sqlc | SQL previsível, injection impossível |
| `$fillable` | struct de request explícita | a classe de bug não existe |
| Blade | templ | tipagem em tempo de compilação |
| Model com regra dentro | service | entidade é dado, não comportamento |

## A frase de onboarding

> Se você conhece Laravel, você já sabe 80% do Arandu. Os outros 20% são o
> compilador cobrando o que o Laravel só recomendava.

## O que isso exige da documentação (Fase D)

Um guia "Vindo do Laravel" com esta tabela na primeira página. É o documento de
maior conversão possível para o público-alvo, e ele já está escrito aqui —
falta só virar prosa quando a API congelar.
