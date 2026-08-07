# ADR 0012-starter-kit-e-comando — O `oka` não é repositório, é `aru make:auth`

**Status:** Aceita. Altera `05-repositorios.md` e `07-pacotes.md`, que previam
`arandu-io/oka` como um dos dez repositórios da Fase 2.

## Contexto

`plans/REFERENCIAS.md` descreve o contrato do `laravel/breeze` literalmente:

> publica views, rotas e controllers no projeto do usuário e sai de cena. Não é
> dependência de runtime. É esse contrato que queremos.

Um repositório Go entregaria o oposto disso. Em Go, a única forma de consumir um
repositório é `import` — e o que se importa não se edita. Editar a tela de login
é a primeira coisa que qualquer pessoa faz.

O Breeze resolve isso sendo um pacote de desenvolvimento com um comando
(`breeze:install`) que copia arquivos e depois é removido do projeto. Em Go o
equivalente do comando já existe: a `aru`, que já é o gerador determinístico com
golden files e blocos `arandu:begin custom` preservados na regeração.

## Decisão

O starter kit é `aru make:auth`. Não existe `arandu-io/oka`.

Os templates vivem em `aru/internal/gen/auth.go` e publicam no projeto o
controller, as rotas e as telas de autenticação, nos diretórios do Laravel
(ADR 0019): `app/Http/Controllers/`, `routes/web.go` e as views em
`resources/views/auth/*.kyse.go`. Nada disso é dependência de runtime, e cada
arquivo é do projeto a partir do momento em que aparece.

O que se publica é o `.kyse.go` — o fonte que o humano lê e edita. O `.go` ao
lado sai do `aru view:build` (ADR 0020).

A portabilidade para a árvore da ADR 0019 **está feita**: o comando emite doze
arquivos, três controllers e as nove views do `laravel/ui`, todos nos caminhos
acima e todos em `.kyse.go`.

```
app/Http/Controllers/Auth/LoginController.go
app/Http/Controllers/Auth/LoginController_handlers.go
app/Http/Controllers/HomeController.go
resources/views/layouts/app.kyse.go
resources/views/home.kyse.go  ·  resources/views/welcome.kyse.go
resources/views/auth/login.kyse.go  ·  register.kyse.go  ·  verify.kyse.go
resources/views/auth/passwords/{confirm,email,reset}.kyse.go
```

Quatro deles — o layout, `home`, `welcome` e o `HomeController` — são
**substituídos**, não pulados, mesmo sem `--force`: em kyse uma página renderiza
com o tipo do layout dela, então o layout e o que o estende são uma unidade só.
O `php artisan ui bootstrap --auth` sobrescreve o `home.blade.php` pelo mesmo
motivo. O `--force` continua governando os outros oito.

**Não existe `modules/authui/` e não existe `views.templ`.** O nome `authui`
sobrevive só como nome de pacote Go dentro de `app/Http/Controllers/Auth/`.
Quem ler isto procurando uma pendência para fechar: não há. Recriar o diretório
vertical seria desfazer a ADR 0019.

## Por quê

**REGRA 9.** Um repositório `oka` seria uma segunda forma de entregar código
gerado ao projeto, ao lado do `make:module` que já existe. Duas formas de fazer
a mesma coisa é exatamente o que a regra proíbe — e a segunda teria que ser
mantida para sempre.

**Não haveria o que importar.** Um pacote `oka` que exporta um `kernel.Module`
com telas embutidas não é o Breeze: é o oposto dele. Quem quisesse mudar uma
palavra da tela teria que forkar.

**O gerador já paga o custo.** Golden files, preservação de bloco custom,
verificação de sintaxe e teste de compilação contra o `framework` já existem para
o `make:module`. O starter kit herda os quatro.

## Consequências

- **Fecha:** não haverá repositório de starter kit, nem um segundo binário. Um
  segundo tema de autenticação, se um dia existir, é uma flag do comando
  (`--style=...`), nunca um repositório novo.
- **Um repositório a menos.** À época isso levava a Fase 2 de dez para nove;
  depois o `database` entrou (ADR 0014) e o `porang` foi dissolvido (ADR 0021).
  A lista atual está em `05-repositorios.md`.
- **O nome `oka` fica livre.** Continua aprovado pelos cinco filtros da REGRA 4 e
  disponível para outro uso, se aparecer um.
- `make:auth` depende da camada de view, e portanto da Fase 2 estar de pé. Isso é
  ordem, não acoplamento: o esqueleto continua autenticando sem ele, com a tela
  mínima do módulo `auth` do core.
