# ADR 0019-estrutura-do-laravel — A árvore é a do Laravel, diretório a diretório

**Status:** Aceita. Supersede a ADR 0003.

## Contexto

A ADR 0003 decidiu módulo vertical: uma feature em um diretório, tudo do domínio
junto. A justificativa era boa e continua sendo — o custo de entender uma feature
cai de seis arquivos espalhados para um `ls`.

O problema é que o alvo do produto não é o desenvolvedor que gosta de módulo
vertical. É o desenvolvedor Laravel, e a frase de onboarding é "se você conhece
Laravel, você já sabe 80% do Arandu".

Um `modules/billing/` quebra essa frase no primeiro `ls`. Não importa que a
organização seja melhor: quem chega não está avaliando organização, está
procurando `app/Http/Controllers`. Se não acha, o que sobra é um framework Go
novo com vocabulário do Laravel colado por fora — que é exatamente o que a
REGRA 10 existe para evitar.

Reconhecimento vale mais que elegância de organização. Essa é a decisão, e ela
custa a elegância.

## Decisão

**A árvore de diretórios de um projeto Arandu é a do Laravel, diretório a
diretório.** O `modules/<nome>/` vertical cai.

```
app/Http/Controllers    app/Models       app/Jobs        bootstrap/
app/Http/Middleware     app/Enums        app/Events      config/
app/Http/Requests       app/Policies     app/Listeners   database/
                        app/Services     app/Mail        resources/
                        app/Repositories app/Providers   routes/
                                                         storage/
                                                         public/
```

`config/` tem dez arquivos, como o Laravel. `database/` tem `factories`,
`migrations` e `seeders`. `resources/` tem `views`, `css` e `js`. `routes/` tem
`web.go` e `console.go`.

### O que muda de verdade — os outros 20%

`app/Services/`, `app/Repositories/` e `app/Policies/` são **obrigatórios**, e o
`aru doctor` cobra.

No Laravel os três são convenção de time organizado: existem quando alguém decide
que existem, e o framework não tem opinião. Aqui são esqueleto. Controller não
alcança repositório, repositório não é alcançável sem policy, e regra de negócio
não mora no controller — nada disso é sugestão, é o que a REGRA 17 e a tese do
produto exigem.

É por isso que a árvore ser idêntica não torna o framework idêntico. O
desenvolvedor Laravel reconhece cada diretório e aprende só o que mudou: três
deles deixaram de ser opcionais.

### Decisões que vêm junto

- **Motor de template próprio: `kyse`** (guarani *kyse*, faca). O templ cai. A
  extensão do fonte é `.kyse.go` — verificado com Go 1.26.5 e cache limpo:
  `go build`, `go vet`, `go test` e `go run` passam com `//go:build kyse` no
  topo. O `gofmt` é o único que ignora build tag, e o CI o chama com
  `gofmt -l $(find . -name '*.go' -not -name '*.kyse.go')`. A cláusula `package`
  é obrigatória quando há outro `.go` no diretório.

  O gerado fica ao lado, no mesmo diretório e no mesmo pacote:
  `resources/views/home.kyse.go` → `resources/views/home.go`. **Não existe
  `storage/framework/views/`.** Nada de view sai da árvore.

  As diretivas têm o nome do Blade: `@extends`, `@section`/`@endsection`,
  `@yield`, `@include`, `@if`/`@elseif`/`@else`/`@endif`,
  `@foreach`/`@endforeach`, `@csrf`, `{{ }}` escapado, `{!! !!}` cru, e
  `@go`/`@endgo` no lugar do `@php`. Bloco de código é `@go` ... `@endgo`,
  nunca `{{ go }}`.

- **O scaffolding de auth espelha o do `laravel/ui`** — as mesmas nove views. A
  stack de CSS não muda: Tailwind v4 com o binário standalone, vendorizado com a
  licença e
  embutido por `go:embed`.

- **O `porang` foi dissolvido.** De onze repositórios para dez. Views e
  componentes foram para `arandu/resources/views/`; o runtime de render e os
  assets embutidos, para `framework/view/`.

- **O dado chega à view como struct tipada por página, nunca `map`.** É o que
  sustenta a tese: campo errado é erro de compilação, não página em branco em
  produção. Uma regra de `aru doctor` proíbe `map` como dado de view.

Os dez repositórios: `framework` (= laravel/framework), `arandu`
(= laravel/laravel), `aru` (= artisan), `database`, `kv`, `queue`, `storage`,
`examples`, `homebrew-tap`, `docs`. Licença MIT, ADR 0008 — qualquer menção a
Apache-2.0 é regressão.

## Custo técnico conhecido — pago

O `aru doctor` foi escrito contra `modules/`. O campo `file.module` só era
preenchido para caminhos que começavam com `modules/`, e ficava vazio para todo o
resto. Com a árvore nova ele ficaria vazio sempre — e regra que testa
`f.module == ""` e sai não reporta nada. Ela não falha: ela passa.

**Buraco de doctor é buraco na tese.** Um doctor que passa porque não olhou é
pior que doctor nenhum: ele dá um verde.

A detecção foi reescrita: `classify` deriva **categoria** e **entidade** do
caminho — `app/Policies/InvoicePolicy.go` é `Policies`/`Invoice` — e as regras
passaram a raciocinar por entidade, não por diretório de módulo.
`project.blind` pergunta pela entidade, `hasRedaction` pelo diretório do tipo
(que é o pacote do receiver) e `usedPermissions` lê o `arandu.mod.toml` da raiz,
porque a unidade de distribuição é o módulo Go (doc 18), não um subdiretório.

- `repositoryNeedsPolicy`, `policyMustBeOpened`, `sensitiveFieldNeedsRedaction`,
  `declaredPermissionsMatchTheCode` — reescritas por entidade
- `handlerMustNotReachData` virou `controllerMustNotReachData`
- `moduleMustNotImportModule` **cai** com a estrutura que policiava, e é
  substituída por `controllerMustNotReachTheRepository` — a mesma fronteira, na
  direção que a árvore nova tem
- duas regras novas que o doc 14 exige: `view-data-is-a-map` e
  `view-does-not-exist`

`TestTheTreeIsLaravels` falha se o caminho mudar de novo, que é o que esta seção
pedia.

## O que isso não fecha

Agrupar por domínio dentro de `app/` continua possível para quem quiser, do mesmo
jeito que é possível no Laravel. O que sai é o `modules/` como estrutura do
framework — com contrato de módulo, gerador e regra de doctor apontando para ele.

## Referências

- ADR 0003 — módulos verticais (superseded)
- ADR 0008 — licença MIT
- ADR 0010 — estrutura de views
- `docs/10-de-para-laravel.md` — o modelo mental é Laravel + Inertia
- `docs/09-uma-forma-so.md`
- `aru/internal/doctor/doctor.go` — `classify`, categoria e entidade

Supersede a ADR 0003.
