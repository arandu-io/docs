# 28 — Estrutura do projeto

O doc 13 mapeia o `src/Illuminate/*` — as entranhas do framework. Este mapeia o
**esqueleto da aplicação**: o que `aru new` gera e o que o dev vê ao abrir o
projeto pela primeira vez.

É o documento que mais decide a promessa dos 80%. Ninguém lê o interior do
framework no primeiro dia; todo mundo dá `ls`.

---

## O que `aru new crm` gera

```
meu-app/
├── main.go                          ← public/index.php + artisan
├── migrate.go  ·  background.go     ← os despachos que main.go usa
├── *_test.go                        ← onze testes de esqueleto, ao lado do código
├── go.mod  ·  go.sum
├── arandu.toml
├── .env  ·  .env.example
│
├── app/
│   ├── Http/
│   │   ├── Controllers/
│   │   │   ├── Controller.go
│   │   │   └── HomeController.go
│   │   ├── Middleware/
│   │   └── Requests/
│   ├── Models/
│   │   └── User.go
│   ├── Enums/
│   ├── Policies/
│   │   └── UserPolicy.go
│   ├── Services/                    ← não existe no Laravel
│   ├── Repositories/                ← não existe no Laravel
│   ├── Jobs/  ·  Events/  ·  Listeners/  ·  Mail/
│   └── Providers/
│       └── AppServiceProvider.go
│
├── bootstrap/
│   └── app.go
│
├── config/
│   ├── app.go  ·  auth.go  ·  cache.go  ·  database.go
│   ├── filesystems.go  ·  logging.go  ·  mail.go
│   └── queue.go  ·  services.go  ·  session.go
│
├── database/
│   ├── factories/UserFactory.go
│   ├── migrations/migrations.go
│   └── seeders/{DatabaseSeeder,AdminSeeder,seeders}.go
│
├── resources/
│   ├── css/app.css
│   └── views/
│       ├── layouts/app.kyse.go
│       ├── home.kyse.go
│       ├── layouts_app.go  ·  home.go        ← o gerado, versionado
│       └── page.go                  ← o struct Page, que toda tela embute (doc 14)
│
├── assets/
│   ├── app.css                      ← a saída do Tailwind
│   └── assets.go                    ← o go:embed do app.css, servido no lugar do padrão
│
├── routes/
│   ├── web.go
│   └── console.go
│
├── storage/
│   ├── app/{private,public}/
│   └── framework/{cache,sessions}/
│
├── public/
│   ├── favicon.ico  ·  robots.txt
│   └── public.go                    ← o go:embed dos dois; não é document root
│
├── .github/workflows/ci.yml  ·  .vscode/  ·  .claude/
├── Dockerfile  ·  compose.yml  ·  Taskfile.yml
└── README.md  ·  CONTRIBUTING.md  ·  SECURITY.md  ·  LICENSE.md
```

Os oito diretórios do Laravel — `app`, `bootstrap`, `config`, `database`,
`public`, `resources`, `routes`, `storage` — nos mesmos lugares, com os mesmos
nomes, em CamelCase dentro de `app/` como no PSR-4.

### O que o `aru new` **não** gera, e onde aparece

Quatro coisas costumam ser lembradas como parte do esqueleto e não são. Cada uma
já foi criada por engano ao menos uma vez.

| O quê | Realidade |
|---|---|
| `bootstrap/providers.go` | não existe. `bootstrap/` tem **um** arquivo, `app.go`. No Laravel 11 o `providers.php` é uma lista de classes que o container instancia; sem container (ADR 0001) a lista é a própria fiação, e ela já está no `app.go` |
| `app/Console/Commands/` | não existe. Comando é uma função registrada em `routes/console.go`, no mesmo desenho de `routes/web.go` |
| `app/Rules/` | não existe, e **não há o que colocar nele**: `framework/validation` não expõe contrato de regra — sem interface `Rule`, sem registro, sem tag. Uma regra de domínio é uma função chamada dentro do `Validate() Errors` do request. Ver a seção sobre `app/Rules/` no de-para do `app/` |
| `resources/views/welcome.kyse.go` | não sai do `aru new`. Aparece com `go run github.com/arandu-io/ui@latest auth`, junto com `auth/` e o `HomeController` (ADR 0026) |

E a raiz tem **14 arquivos `.go`**, não um: `main.go`, `migrate.go`,
`background.go` e onze `_test.go`. O `main.go` é o ponto de entrada; os outros
dois são os despachos que ele chama, separados porque `main.go` decide e não
executa. O que fixa o desenho são os três de código; o número total sobe a cada
teste de esqueleto novo, e conferir é `ls arandu/*.go`.

### A migração de `users` não fica no projeto

O Laravel gera `database/migrations/0001_01_01_000000_create_users_table.php`. O
`aru new` **não** gera equivalente: `database/migrations/migrations.go` sai
vazio, com o bloco custom, e a tabela `users` vem do módulo `auth` do framework,
que declara a própria migração em `Migrations()`.

O motivo está escrito no comentário de doc do `migrations.go` gerado, e não em
ADR: repetir aqui o que o kernel já coleta dos módulos aplicaria cada migração
duas vezes. Vale o mesmo para o `outbox` (`events`) e para `jobs` (`queue`).

## Os quatro arquivos que não sobrevivem, e por quê

| Laravel | Arandu | Motivo |
|---|---|---|
| `public/index.php` | `main.go` na raiz | em Go o binário **é** o servidor; não há front controller nem document root |
| `bootstrap/cache/packages.php` `services.php` | **somem** | cache de descoberta existe porque o PHP redescobre a cada request. Compilamos uma vez |
| `tests/Feature/ExampleTest.php` | `*_test.go` ao lado do código | única convenção que não dobra — ver abaixo por quê, e o que exatamente quebra |
| `resources/views/home.blade.php` | `resources/views/home.kyse.go` | o fonte é `.kyse.go`, protegido por `//go:build kyse`; o compilador ignora, o `kyse` lê. O **gerado** vai para a raiz de `resources/views/` — ver abaixo |

Somem também `vendor/`, `.htaccess`, `package.json` e `vite.config.js`.

### Por que não existe `tests/`

*Escrito em 07/08/2026, e a linha da tabela acima dizia menos do que devia — e
uma parte errada: "o Go só executa arquivo terminado em `_test.go`, no pacote do
código testado". A primeira metade é verdade; a segunda não.*

**O Go acha, sim, um teste dentro de `tests/`.** `go test ./...` roda um
`tests/conta_test.go` normalmente — medido. Não é proibição.

O que quebra são duas coisas, e as duas importam:

**1. A cobertura vai para o pacote errado.** Com o teste em `tests/`:

```
probe            coverage: 0.0% of statements
probe/tests   ok coverage: [no statements]
```

O pacote que está sendo testado reporta **zero**, e o pacote de teste reporta
que não tem código. Nenhum dos dois números descreve a realidade, e é a partir
deles que se decide o que falta testar.

Com o mesmo teste ao lado do código: `ok probe — coverage: 100.0% of statements`.

**2. O teste só alcança o que é exportado.** De `tests/`, `probe.naoExportado()`
é `undefined: probe.naoExportado` — erro de compilação. Metade do que se testa
neste projeto é comportamento não exportado: a regra do doctor, o gerador do
kyse, o helper que decide escopo.

### E o Go já tem a distinção que a pasta `tests/` costuma querer dizer

Uma pasta `tests/` normalmente separa "teste de unidade" de "teste que só usa a
API pública". O Go expressa isso **sem pasta**, com dois pacotes no mesmo
diretório:

| arquivo | pacote | enxerga |
|---|---|---|
| `conta_test.go` | `package probe` | tudo, inclusive o não exportado |
| `conta_ext_test.go` | `package probe_test` | só a API pública, como um cliente |

O `probe_test` é o único nome de pacote que o Go deixa conviver com outro no
mesmo diretório, e existe exatamente para isso. Este projeto usa os dois: **50
arquivos** em `package X_test` e **25** em `package X`, e a escolha entre eles é
a pergunta "isto é o contrato ou é a implementação?".

É por isso que a linha da tabela diz *"única convenção que não dobra"*: não é
teimosia com o Laravel, é que a ferramenta atribui cobertura por diretório e a
visibilidade por pacote. Uma pasta `tests/` custa as duas coisas para ganhar uma
árvore mais parecida.

A extensão dupla é o que faz o fonte conviver com o Go sem tooling paralelo:
`.kyse.go` termina em `.go`, então `go build`, `go vet`, `go test` e `go run`
enxergam o arquivo — e a tag de build `kyse`, que nenhum comando define, o
descarta antes do parser. Ver a seção **A extensão `.kyse.go`** adiante.

## Os três diretórios novos — e eles são os 20%

`app/Services/`, `app/Repositories/` e `app/Policies/` obrigatório.

No Laravel, Service e Repository são convenção de time organizado, não do
framework — quem quer, cria; quem não quer, engorda o Controller. Aqui são
parte do esqueleto e o `doctor` cobra.

É literalmente a frase de onboarding, virada em diretório: *se você conhece
Laravel, você já sabe 80% do Arandu; os outros 20% são o compilador cobrando o
que o Laravel só recomendava.*

---

## De-para do diretório raiz

| Laravel | Arandu | Nota |
|---|---|---|
| `app/` | `app/` | **igual**, com `Services/` e `Repositories/` a mais |
| `bootstrap/app.php` | `bootstrap/app.go` | **igual** — o único lugar de fiação. Escrito à mão, e é o que o `aru make:module` manda você editar (ADR 0023) |
| `public/index.php` + `artisan` | `main.go` | o binário é o servidor **e** a linha de comando. Ele despacha; quem compõe é o `bootstrap/app.go` |
| `bootstrap/cache/` | **não existe** | cache de rota e config é para linguagem interpretada. Nós compilamos |
| `config/` | `config/` — dez arquivos, um por domínio | **igual ao Laravel**: `app`, `auth`, `cache`, `database`, `filesystems`, `logging`, `mail`, `queue`, `services`, `session`. Cada um exporta uma struct tipada, lida do `.env`; erro de chave é erro de compilação |
| `database/` | `database/` | **igual** — `factories/`, `migrations/`, `seeders/` |
| `public/` | **quase não existe** | asset entra por `go:embed` (doc 14). Sem `storage:link`, sem asset servido de disco. O que sobra é `favicon.ico` e `robots.txt` — endereços que o cliente escolhe, e por isso não podem ter hash no caminho — embutidos por `public/public.go`, com uma rota por nome |
| `resources/views/` | `resources/views/` | **igual** — layouts, componentes e páginas. É onde as views do dissolvido `porang` foram parar; o runtime de render ficou em `framework/view/` |
| `resources/css/app.css` | `resources/css/app.css` | **igual** — a entrada do Tailwind; `aru view:build` compila para `assets/app.css`, embutido por `go:embed` |
| `resources/js/app.js` | **não existe** | *corrigido em 07/08/2026.* Esta linha dizia "**igual** — JS à mão, sem npm e sem bundler". O `aru new` não gera `resources/js/`: o esqueleto tem `resources/css/` e `resources/views/`, e não há um `.js` no repositório. HTMX e Alpine vêm embutidos em `framework/view/assets/` |
| `public/build/` (saída do Vite) | **não existe** | o CSS compilado entra no binário por `go:embed`. Sem bundler, sem diretório de saída servido do disco |
| `routes/web.php` | `routes/web.go` | **igual** |
| `routes/api.php` | o mesmo `routes/web.go` | o handler decide JSON ou fragmento |
| `routes/console.php` | `routes/console.go` | **igual** |
| `storage/logs/` | **stdout** | log estruturado vai para o agregador, não para arquivo |
| `tests/` | `*_test.go` ao lado do código | convenção do Go, e é melhor: o teste anda junto |
| `vendor/` | **não existe** | `go.sum` e cache de módulo |

Três dessas linhas são ganho direto, e valem ser ditas ao dev que chega:
**sem `bootstrap/cache`, sem `vendor/`, sem `storage/logs`**. Três fontes
clássicas de "funciona na minha máquina" que somem.

## Onde o `kyse` entra

Duas camadas, como no Laravel:

| Laravel | Arandu |
|---|---|
| `resources/views/layouts/app.blade.php` | `resources/views/layouts/app.kyse.go` |
| `resources/views/components/button.blade.php` | `resources/views/components/button.kyse.go` |
| `resources/views/dashboard.blade.php` | `resources/views/dashboard.kyse.go` |
| `resources/views/clientes/index.blade.php` | `resources/views/clientes/index.kyse.go` |

## Erros que este documento previne

- Salvar a view como `home.kyse` — o fonte é `home.kyse.go`, com
  `//go:build kyse` na primeira linha e cláusula `package`
- Editar `resources/views/home.go` — é o gerado; a próxima geração apaga
- Criar `storage/framework/views/` para o compilado — o gerado fica na raiz de
  `resources/views/`, no mesmo pacote do fonte. Duas formas para a mesma coisa
  fere a REGRA 9
- Procurar o gerado de `auth/login.kyse.go` dentro de `resources/views/auth/` —
  ele é `resources/views/auth_login.go`, um diretório acima
- Passar `map[string]any` para a view — o dado chega como struct tipada por
  página, e o `doctor` reprova o `map`
- Criar `modules/<nome>/` vertical — a estrutura é a do Laravel, diretório a
  diretório
- Criar `tests/` — Go usa `_test.go` ao lado
- Criar `public/js` e `public/css` — asset é embutido (doc 14, REGRA 13)
- Criar `internal/wire/` — o diretório não existe e nunca existiu. A fiação é
  `bootstrap/app.go`, escrita à mão: o `aru make:module` **imprime** as três
  linhas e não edita arquivo nenhum (ADR 0023). Fiação gerada ao lado da fiação
  escrita à mão seriam duas fiações (REGRA 9)

---

## Anexo — três ajustes sobre a proposta original

### 1. `resources/assets/` é Laravel 5, não Laravel 12

O `resources/assets/sass` e `resources/assets/js` existiam na época do Elixir. O
Laravel moderno tirou o nível `assets/`: hoje é `resources/css/app.css` e
`resources/js/app.js`, direto.

Copiamos o nível a menos, não os dois diretórios: aqui existe
`resources/css/app.css` e não existe `resources/js/` — ver o anexo 3. O dev que
abrir reconhece o caminho do CSS, que é o que a entrada do Tailwind precisa ser.

### 2. `public/` não fica dentro de `resources/`

Nem no Laravel. Lá, `resources/` é **fonte** e `public/` é **saída** — o
document root para onde o servidor aponta. Misturar os dois é o erro.

E aqui `public/` praticamente não existe: asset é embutido (REGRA 13), não
servido de disco. O que o Laravel joga em `public/build/` — a saída do Vite — em
nós vai para **`assets/app.css`**, na raiz do projeto, e entra no binário por
`go:embed`.

Fonte e saída não moram juntos: `resources/css/app.css` é a **entrada** do
Tailwind e `assets/app.css` é o que o `aru view:build` escreve. Só o segundo é
embutido, e editar o segundo é trabalho perdido — a próxima compilação o
sobrescreve.

Se sobrar alguma coisa que precise ser servida de disco em dev, ela vai em
`public/` na **raiz**, não dentro de `resources/`.

Restrição técnica do `go:embed` que decide o layout: a diretiva só alcança
arquivo no mesmo diretório ou abaixo. Então o arquivo com o `//go:embed` fica na
raiz do projeto ou dentro de `resources/`, nunca em `internal/`.

### 3. Arquivo `.js` não fura a REGRA 13 — e o esqueleto não tem nenhum

*Corrigido em 07/08/2026:* este anexo se chamava "`resources/js/app.js` é
permitido" e legislava sobre esse arquivo, terminando em "se o arquivo passar a
ter `fetch`, store global ou regra de negócio, o `doctor` reprova". Duas coisas
eram falsas.

**O arquivo existiu, e foi apagado em 07/08/2026** — commit `650f9b3` do
`arandu`, que registra: tinha 37 linhas anexando o `X-CSRF-Token` no
`htmx:configRequest`, e nada o embutia, servia ou referenciava. Hoje
`arandu/resources/` tem `css/` e `views/`, e não há
um `.js` no esqueleto. O `aru` não menciona `resources/js` em lugar nenhum, e
o `framework/view` não teria como servi-lo: ele embute HTMX, Alpine e o CSS, e a
única troca que a aplicação pode fazer é `RegisterStylesheet` — não existe
`RegisterScript`. Servir seria possível por outro caminho, o mesmo que o
esqueleto já usa para o favicon e o `robots.txt` em `public/public.go`:
`go:embed`, uma tabela de content type e uma rota por nome. O que não existe é
lugar para o arquivo, não mecanismo.

**E o `doctor` não lê `.js`.** `parseViews`, em
`aru/internal/doctor/doctor.go`, percorre `resources/views/` e coleta
exclusivamente os `*.kyse.go` — é essa lista que a regra
`alpine-reaches-the-server` varre. Ela reprova `x-data`, `x-init`, `x-effect`,
`x-on:*` e `@*` que contenham `fetch(`, `axios`, `XMLHttpRequest`, `$store`,
`navigator.sendBeacon`, `EventSource(` ou `new WebSocket`, **dentro da view**.
Um `.js` ao lado passaria intocado, porque o `doctor` não o abre.

O que continua valendo, e é por isso que o anexo não sai: é fácil ler a REGRA 13
como "nada de JavaScript", e não é isso. Ela proíbe **Node no build** — npm,
`package.json`, bundler, Vite — e não proíbe arquivo `.js`.

| Permitido | Proibido |
|---|---|
| `.js` escrito à mão, embutido e servido | `package.json` |
| definição de componente Alpine em JS puro | `npm install` |
| `<script>` apontando para asset embutido | Vite, webpack, esbuild |
| — | `import` que precise de transpilação |

Hoje isso é permissão sem uso: não há diretório no esqueleto nem registro no
`framework/view` para um script da aplicação. Criar um sem ADR seria abrir uma
segunda origem de JavaScript ao lado do `go:embed` do framework (REGRA 9). Se
algum dia entrar, o limite do doc 14 continua valendo: estado do cliente,
efêmero e invisível ao servidor, as três ao mesmo tempo.

E há verificação — só não é do `doctor`, cuja checagem para na view. Quem
impede é o próprio esqueleto: `TestResourcesHoldNoJavaScript` percorre
`resources/` e reprova qualquer `.js` — *"is never built, embedded or served"* —
e `TestTheOnlyScriptsServedAreTheEmbeddedOnes` fixa que todo `<script>` emitido
aponta para os assets embutidos.

O que significa que a permissão da tabela acima é teórica: um `.js` posto em
`resources/` hoje reprova o test suite do esqueleto, e exercê-la custa apagar um
teste — que é o atrito certo para uma segunda origem de JavaScript.

---

## De-para do `app/`

| Laravel | Arandu |
|---|---|
| `app/Http/Controllers/HomeController.php` | `app/Http/Controllers/HomeController.go` |
| `app/Http/Middleware/` | `app/Http/Middleware/` |
| `app/Http/Requests/StorePedido.php` | `app/Http/Requests/StorePedido.go` |
| `app/Models/User.php` | `app/Models/User.go` |
| `app/Enums/Status.php` | `app/Enums/Status.go` |
| `app/Policies/PedidoPolicy.php` | `app/Policies/PedidoPolicy.go` |
| `app/Providers/AppServiceProvider.php` | `app/Providers/AppServiceProvider.go` |
| `app/Jobs/` `Events/` `Listeners/` `Mail/` | iguais |
| `app/Rules/` | **não existe** — ver abaixo |
| `app/Console/Commands/` | **não existe** — comando é uma função em `routes/console.go` |
| — | **`app/Services/`** — regra de negócio |
| — | **`app/Repositories/`** — acesso a dado, exige `Grant` |

### `app/Rules/` não existe, e não é omissão

No Laravel, `app/Rules/` guarda classes que implementam `ValidationRule` e são
referenciadas por nome numa lista de regras — o mecanismo de string do
`Validator`. **`framework/validation` não tem esse mecanismo**: não há interface
`Rule`, não há registro, não há tag de struct, e a validação não usa reflection.
O contrato é `Validatable`, com um único método `Validate() Errors`, e as regras
são funções comuns — `Required`, `MinLen`, `MaxLen`, `Email`, `NotZero` — que
recebem o `Errors` e escrevem nele.

Uma regra de domínio, então, é uma função Go chamada de dentro do `Validate()`
do request, no mesmo pacote em que ela é usada. Sem diretório, sem registro e
sem indireção.

Criar `app/Rules/` só faz sentido junto com um contrato de regra, e um contrato
de regra é uma segunda forma de validar ao lado do `Validate() Errors` que já
existe (REGRA 9). Se algum dia entrar, entra por ADR — não por alguém notando
que o Laravel tem o diretório.

Nome de arquivo em CamelCase, como no PSR-4. Não é convenção de Go — mas é o
que faz o dev Laravel reconhecer, e é o objetivo. O nome de **pacote** dentro do
arquivo segue Go (minúscula): `package controllers`.

## O `kyse`

`kyse` é o motor de template do Arandu — o lugar que o Blade ocupa no Laravel.
Vem do guarani **kysê**, faca, atestado em dicionário. O acento cai no nome de
arquivo e no import path, por restrição do Go.

| Laravel | Arandu |
|---|---|
| `resources/views/home.blade.php` | `resources/views/home.kyse.go` |
| `resources/views/layouts/app.blade.php` | `resources/views/layouts/app.kyse.go` |
| `storage/framework/views/a3f9…php` (cache) | `resources/views/home.go` (gerado) |
| `storage/framework/views/…` de uma view aninhada | `resources/views/auth_login.go` (gerado, achatado) |

O par de arquivos existe nos dois. No Laravel o gerado fica escondido em
`storage/`; aqui ele fica visível, no mesmo pacote — e é código Go de verdade,
checado pelo compilador. Na raiz de `resources/views/` ele fica ao lado do
fonte; para uma view aninhada, não — ver "Onde o gerado fica" adiante.

Uma forma só: `storage/framework/views/` **não existe**. O compilado não sai da
árvore do fonte.

## A extensão `.kyse.go`

O fonte da view é `.kyse.go`, não `.kyse`. O arquivo abre assim:

```go
//go:build kyse

package views

@extends('layouts.app')
```

A tag de build `kyse` não é definida por nenhum comando, então o toolchain
descarta o arquivo antes de tentar entendê-lo. Verificado com Go 1.26.5 e cache
limpo: `go build ./...`, `go vet ./...`, `go test ./...` e `go run` passam.

Dois detalhes que a verificação expôs:

- A cláusula `package` é obrigatória quando existe outro `.go` no diretório.
- O `gofmt` é o único que ignora a tag de build — ele formata o que a tag
  exclui. Por isso o CI o chama filtrando:
  `gofmt -l $(find . -name '*.go' -not -name '*.kyse.go')`.

O ganho de terminar em `.go`: nenhuma ferramenta de Go precisa ser ensinada
sobre o `kyse`. Editor, `go list`, `go vet` e CI tratam o diretório como
qualquer outro pacote.

## Diretivas — idênticas ao Blade

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
@go · @endgo         o @php
```

São vinte e três nomes, e o conjunto inteiro está nos dois mapas de
`aru/internal/kyse/kyse.go`. A lista de referência é a do doc 14.

**A interpolação começa com ponto:** `{{ .Campo }}` vira `d.Campo`, onde `d` é o
struct daquela página. Sem o ponto a expressão é passada adiante como Go, então
`{{ variavel }}` não lê campo nenhum — quebra o build.

Mesmos nomes, mesma sintaxe. A única adição é `@go` … `@endgo`, que é o `@php`
com o nome da linguagem hospedeira trocado — mesmo raciocínio do Blade, e a
mesma forma de diretiva `@`. Não existe `{{ go }}`: chave dupla é interpolação,
bloco de código é `@`.

O conjunto é fechado. Caso fora do padrão se escreve em Go, dentro de `@go` …
`@endgo` (REGRA 15). Numa view **gerada** pelo `aru make:module`, o `@go` já vem
com um par `// arandu:begin custom` … `// arandu:end custom` dentro dele: é ali
que o que a regeneração precisa preservar tem de ficar.

### Onde o gerado fica, e por que nem sempre ao lado

Para uma view na raiz, o gerado fica ao lado: `resources/views/home.kyse.go`
gera `resources/views/home.go`.

Para uma view aninhada, não. `resources/views/auth/login.kyse.go` gera
`resources/views/auth_login.go` — na raiz de `resources/views/`, com o caminho
no nome.

**É restrição da linguagem, não preferência.** Go tem um pacote por diretório.
Com o gerado ao lado de um fonte aninhado, `resources/views/auth/` seria um
pacote diferente de `resources/views/`, e a struct que o layout declara não
seria visível para a página que o estende — que é exatamente o que faz uma
página herdar o dado do layout.

A alternativa seria achatar o fonte também, e aí se perde o que a estrutura
existe para dar: o desenvolvedor Laravel abrindo `resources/views/auth/` e
encontrando `login`, `register` e `passwords/` onde espera.

Descoberto ao implementar, não ao decidir. Estava escrito "ao lado" para todos.
