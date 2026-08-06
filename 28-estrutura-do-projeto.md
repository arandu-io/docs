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
├── go.mod  ·  go.sum
├── arandu.toml
├── .env  ·  .env.example
│
├── app/
│   ├── Http/
│   │   ├── Controllers/
│   │   │   ├── Auth/
│   │   │   │   ├── LoginController.go
│   │   │   │   ├── RegisterController.go
│   │   │   │   └── ResetPasswordController.go
│   │   │   ├── Controller.go
│   │   │   └── HomeController.go
│   │   ├── Middleware/
│   │   └── Requests/
│   ├── Models/
│   │   └── User.go
│   ├── Enums/
│   ├── Policies/
│   ├── Services/                    ← não existe no Laravel
│   ├── Repositories/                ← não existe no Laravel
│   ├── Jobs/  ·  Events/  ·  Listeners/  ·  Mail/
│   └── Providers/
│       └── AppServiceProvider.go
│
├── bootstrap/
│   ├── app.go
│   └── providers.go
│
├── config/
│   ├── app.go  ·  auth.go  ·  cache.go  ·  database.go
│   ├── filesystems.go  ·  logging.go  ·  mail.go
│   └── queue.go  ·  services.go  ·  session.go
│
├── database/
│   ├── factories/UserFactory.go
│   ├── migrations/0001_01_01_000000_create_users_table.go
│   └── seeders/DatabaseSeeder.go
│
├── resources/
│   ├── css/app.css
│   ├── js/app.js
│   └── views/
│       ├── auth/
│       │   ├── passwords/{confirm,email,reset}.kyse.go
│       │   ├── login.kyse.go  ·  register.kyse.go  ·  verify.kyse.go
│       ├── layouts/app.kyse.go
│       ├── home.kyse.go
│       └── welcome.kyse.go
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
│   ├── favicon.ico
│   └── robots.txt
│
├── Dockerfile  ·  compose.yml  ·  Taskfile.yml
└── README.md  ·  SECURITY.md  ·  LICENSE.md
```

Os oito diretórios do Laravel — `app`, `bootstrap`, `config`, `database`,
`public`, `resources`, `routes`, `storage` — nos mesmos lugares, com os mesmos
nomes, em CamelCase dentro de `app/` como no PSR-4.

## Os quatro arquivos que não sobrevivem, e por quê

| Laravel | Arandu | Motivo |
|---|---|---|
| `public/index.php` | `main.go` na raiz | em Go o binário **é** o servidor; não há front controller nem document root |
| `bootstrap/cache/packages.php` `services.php` | **somem** | cache de descoberta existe porque o PHP redescobre a cada request. Compilamos uma vez |
| `tests/Feature/ExampleTest.go` | `*_test.go` ao lado do código | o Go só executa arquivo terminado em `_test.go`, no pacote do código testado. Única convenção que não dobra |
| `resources/views/home.blade.php` | `resources/views/home.kyse.go` | o fonte é `.kyse.go`, protegido por `//go:build kyse`; o compilador ignora, o `kyse` lê. O **gerado** é `home.go`, ao lado |

Somem também `vendor/`, `.htaccess`, `package.json` e `vite.config.js`.

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
| `bootstrap/app.php` | `main.go` | mesma função: o único lugar de configuração |
| `bootstrap/cache/` | **não existe** | cache de rota e config é para linguagem interpretada. Nós compilamos |
| `config/` | `config/` — dez arquivos, um por domínio | **igual ao Laravel**: `app`, `auth`, `cache`, `database`, `filesystems`, `logging`, `mail`, `queue`, `services`, `session`. Cada um exporta uma struct tipada, lida do `.env`; erro de chave é erro de compilação |
| `database/` | `database/` | **igual** — `factories/`, `migrations/`, `seeders/` |
| `public/` | **quase não existe** | asset entra por `go:embed` (doc 14). Sem `storage:link`, sem asset servido de disco |
| `resources/views/` | `resources/views/` | **igual** — layouts, componentes e páginas. É onde as views do dissolvido `porang` foram parar; o runtime de render ficou em `framework/view/` |
| `resources/css/app.css` | `resources/css/app.css` | **igual** — entrada do Bootstrap, vendorizado com a licença e embutido |
| `resources/js/app.js` | `resources/js/app.js` | **igual** — JS à mão, sem npm e sem bundler |
| `public/build/` (saída do Vite) | **não existe** | o CSS do Bootstrap é vendorizado com a licença e embutido por `go:embed`. Sem bundler, sem diretório de saída |
| `routes/web.php` | `routes/web.go` | **igual** |
| `routes/api.php` | o mesmo `routes.go` | o handler decide JSON ou fragmento |
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
- Criar `storage/framework/views/` para o compilado — o gerado fica ao lado do
  fonte, no mesmo pacote. Duas formas para a mesma coisa fere a REGRA 9
- Passar `map[string]any` para a view — o dado chega como struct tipada por
  página, e o `doctor` reprova o `map`
- Criar `modules/<nome>/` vertical — a estrutura é a do Laravel, diretório a
  diretório
- Criar `tests/` — Go usa `_test.go` ao lado
- Criar `public/js` e `public/css` — asset é embutido (doc 14, REGRA 13)
- Escrever arquivo `.gen.go` de `internal/wire/` à mão — é gerado; edite e a
  próxima geração apaga

---

## Anexo — três ajustes sobre a proposta original

### 1. `resources/assets/` é Laravel 5, não Laravel 12

O `resources/assets/sass` e `resources/assets/js` existiam na época do Elixir. O
Laravel moderno tirou o nível `assets/`: hoje é `resources/css/app.css` e
`resources/js/app.js`, direto.

Copiamos o atual. Um nível a menos, e o dev que abrir vai reconhecer.

### 2. `public/` não fica dentro de `resources/`

Nem no Laravel. Lá, `resources/` é **fonte** e `public/` é **saída** — o
document root para onde o servidor aponta. Misturar os dois é o erro.

E aqui `public/` praticamente não existe: asset é embutido (REGRA 13), não
servido de disco. O que o Laravel joga em `public/build/` — a saída do Vite — em
nós não tem equivalente: o CSS do Bootstrap é vendorizado com a licença dentro
de `resources/css/` e entra no binário por `go:embed`.

Se sobrar alguma coisa que precise ser servida de disco em dev, ela vai em
`public/` na **raiz**, não dentro de `resources/`.

Restrição técnica do `go:embed` que decide o layout: a diretiva só alcança
arquivo no mesmo diretório ou abaixo. Então o arquivo com o `//go:embed` fica na
raiz do projeto ou dentro de `resources/`, nunca em `internal/`.

### 3. `resources/js/app.js` é permitido — e isso não fura a REGRA 13

Ponto que vale esclarecer, porque é fácil ler a regra como "nada de JavaScript".

**Não é isso.** A regra proíbe **Node no build**: npm, `package.json`, bundler,
Vite. Ela não proíbe arquivo `.js`.

| Permitido | Proibido |
|---|---|
| `app.js` escrito à mão, embutido e servido | `package.json` |
| definição de componente Alpine em JS puro | `npm install` |
| `<script>` apontando para asset embutido | Vite, webpack, esbuild |
| — | `import` que precise de transpilação |

Ou seja: `resources/js/alpine/components.js` pode existir. Só precisa ser JS que
o navegador executa direto, sem passar por build.

**Uma ressalva de conteúdo, não de forma.** Um diretório `alpine/components/`
crescendo é sinal de que estado está migrando para o cliente — e o limite do
doc 14 continua valendo: estado do cliente, efêmero e invisível ao servidor. As
três ao mesmo tempo.

Se o arquivo passar a ter `fetch`, store global ou regra de negócio, o `doctor`
reprova. Não pelo arquivo existir, mas pelo que ele faz.

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
| `app/Jobs/` `Events/` `Listeners/` `Mail/` `Rules/` | iguais |
| `app/Console/Commands/` | `app/Console/Commands/` |
| — | **`app/Services/`** — regra de negócio |
| — | **`app/Repositories/`** — acesso a dado, exige `Grant` |

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

O par de arquivos existe nos dois. No Laravel o gerado fica escondido em
`storage/`; aqui ele fica visível ao lado, no mesmo diretório e no mesmo pacote
— e é código Go de verdade, checado pelo compilador.

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
@csrf
{{ variavel }}       escapado
{!! html !!}         cru
@go · @endgo         o @php
```

Mesmos nomes, mesma sintaxe. A única adição é `@go` … `@endgo`, que é o `@php`
com o nome da linguagem hospedeira trocado — mesmo raciocínio do Blade, e a
mesma forma de diretiva `@`. Não existe `{{ go }}`: chave dupla é interpolação,
bloco de código é `@`.

O conjunto é fechado. Caso fora do padrão se escreve em Go, dentro do bloco
`// arandu:begin custom` (REGRA 15).
