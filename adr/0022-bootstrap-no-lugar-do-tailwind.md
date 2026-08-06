# ADR 0022-bootstrap-no-lugar-do-tailwind — O CSS é Bootstrap, como no `laravel/ui`

**Status:** Aceita. Decisão do proprietário. Altera `14-view-e-assets.md` e
`20-licenca.md`, e remove a exceção de npm que o `07-pacotes.md` abria.

## Contexto

A REGRA 10 diz que o modelo mental é o do desenvolvedor Laravel: ele deve
reconhecer o vocabulário imediatamente e só aprender o que muda de fato. A
ADR 0019 aplicou isso à estrutura de diretórios, e a ADR 0020 ao motor de
template. Falta o CSS.

O projeto básico que esse desenvolvedor conhece é o do `php artisan ui bootstrap
--auth`. Ele sabe o que é `layouts/app`, sabe onde fica a tela de login, e sabe
o nome das classes. É a mesma razão que decidiu `@extends` em vez de sintaxe
nova.

O desenho anterior era Tailwind, compilado pelo binário standalone que a `aru`
baixava. Esse caminho existe no código: `aru/internal/toolchain/`, com
verificação de SHA-256, escrita atômica, cache em `~/.arandu/bin/` e versão
fixada em `arandu.toml`. Ele funciona — e é infraestrutura inteira mantida para
produzir um arquivo `.css`.

## Decisão

**O CSS é Bootstrap, vendorizado no repositório do `arandu`, com a licença, e
embutido por `go:embed`.**

As nove views do starter kit são as do `laravel/ui`, com os mesmos nomes:

```
layouts/app
home
welcome
auth/login
auth/register
auth/verify
auth/passwords/confirm
auth/passwords/email
auth/passwords/reset
```

O Tailwind sai. Com ele saem o binário standalone, o download, a fixação de
versão em `arandu.toml` e a compilação de CSS no build.

O Basecoat sai junto: ele é plugin de Tailwind, distribuído por npm, e era o
único motivo pelo qual a exceção de npm em `07-pacotes.md` tinha sido aberta.

## Por quê

**O alvo é o dev Laravel, e é esse o projeto básico que ele conhece.** O
argumento é o mesmo da estrutura de diretórios. Não faria sentido copiar
`app/Http/Controllers` e `resources/views` diretório a diretório e depois
entregar uma tela de login que ele não reconhece.

**Some uma etapa de build inteira.** Bootstrap vendorizado não compila: o
arquivo já é o arquivo servido. A `aru` deixa de baixar um binário de terceiro,
e o `aru dev` deixa de ter um processo a mais para orquestrar e para falhar.

**A REGRA 13 fica mais fácil de cumprir, não mais difícil.** A regra permite ter
etapa de build e proíbe que ela seja Node. O caminho mais seguro para essa
promessa é o que tem menos etapas de build — e agora nenhuma delas produz CSS.

**A exceção de npm perde até o motivo de existir.** Não é mais uma exceção que
resistimos a usar: é uma exceção sem caso de uso.

## O contra, registrado

O `laravel/ui` com Bootstrap é scaffolding da era do Laravel 6. Os starter kits
do Laravel 12 usam Tailwind. Quem chega do Laravel recente vai estranhar — vai
abrir `resources/views/auth/login` esperando classes de utilitário e encontrar
`form-control` e `btn btn-primary`.

Isso não é ignorância da alternativa, é escolha. A decisão foi do proprietário, e
o motivo é que o projeto básico com Bootstrap é o que ele quer como base do
framework.

Se um dia isso mudar, muda por ADR, e muda para todo mundo de uma vez. Não vira
flag. Duas folhas de estilo suportadas seriam duas formas de escrever a mesma
tela, e a REGRA 9 proíbe.

## O custo

Os componentes que existiam em classes Tailwind são reescritos em Bootstrap. É
trabalho de reescrita de markup, não de arquitetura — as views já são tipadas e o
struct de dados de cada página não muda.

O vendor passa a ser mantido à mão: atualizar Bootstrap é trocar o arquivo e a
atribuição no `THIRD_PARTY.md`. Se isso ficar caro, a alternativa é escrever o
CSS. Não é reintroduzir npm.

## O que isso não fecha

Não decide nada sobre HTMX nem sobre Alpine. Eles continuam entrando por
`go:embed`, servidos com `script-src 'self'`, pelos motivos de `docs/14`.

Não toca na página de erro. Ela continua sem CSS compilado e sem asset embutido,
porque precisa renderizar quando o build quebrou.

## Referências

- ADR 0008 — licença MIT, que vale também para o vendor
- ADR 0019 — estrutura do Laravel, diretório a diretório
- ADR 0020 — `kyse` no lugar do templ
- ADR 0021 — dissolução do `porang`
- `docs/14-view-e-assets.md` — o documento da camada de view
- `docs/09-uma-forma-so.md` e REGRA 13
