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

Os templates vivem em `aru/internal/gen/auth.go` e produzem quatro arquivos em
`modules/authui/` do projeto: `module.go`, `handlers.go`, `views.templ` e
`arandu.mod.toml`. Nada disso é dependência de runtime, e cada arquivo é do
projeto a partir do momento em que aparece.

## Por quê

**REGRA 9.** Um repositório `oka` seria uma segunda forma de entregar código
gerado ao projeto, ao lado do `make:module` que já existe. Duas formas de fazer
a mesma coisa é exatamente o que a regra proíbe — e a segunda teria que ser
mantida para sempre.

**Não haveria o que importar.** Um pacote `oka` que exporta um `kernel.Module`
com telas embutidas não é o Breeze: é o oposto dele. Quem quisesse mudar uma
palavra da tela teria que forkar.

**O gerador já paga o custo.** Golden files, preservação de bloco custom,
verificação de sintaxe e teste de compilação contra `framework` e `porang` já
existem para o `make:module`. O starter kit herda os quatro.

## Consequências

- **Fecha:** não haverá repositório de starter kit, nem um segundo binário. Um
  segundo tema de autenticação, se um dia existir, é uma flag do comando
  (`--style=...`), nunca um repositório novo.
- **A Fase 2 passa a ter nove repositórios**, não dez: `framework`, `aru`,
  `arandu`, `docs`, `examples`, `porang`, `kv`, `homebrew-tap` e o que a
  ADR 0013 decide sobre `database`.
- **O nome `oka` fica livre.** Continua aprovado pelos cinco filtros da REGRA 4 e
  disponível para outro uso, se aparecer um.
- `make:auth` depende do `porang`, e portanto da Fase 2 estar de pé. Isso é
  ordem, não acoplamento: o esqueleto continua autenticando sem ele, com a tela
  mínima do módulo `auth` do core.
