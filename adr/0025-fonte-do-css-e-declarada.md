# ADR 0025-fonte-do-css-e-declarada — o Tailwind não descobre fonte sozinho

**Status:** Aceita.

## Contexto

O Tailwind v4 detecta fontes por conta própria, **além** do que o `@source`
declara: ele caminha a árvore a partir do arquivo de CSS e lê nomes de classe de
tudo que encontra, respeitando `.gitignore`.

Essa varredura alcançava `assets/app.css` — a **saída** do próprio comando. O
Tailwind lia os nomes de classe do build anterior e os realimentava, então o CSS
deixava de ser função das views:

- `aru view:build` duas vezes seguidas escrevia dois arquivos diferentes;
- uma classe apagada da última view que a usava sobrevivia mais um build;
- o `app.css` do framework, artefato versionado, mudava a cada recompilação.

O `@source` do framework agravava: apontava para `../layout/**/*.templ` e
`../components/**/*.templ`, duas pastas que a ADR 0021 moveu e um motor que a
ADR 0020 substituiu. As duas globs casavam **nada**. Glob que não casa não
falha — ela não contribui, e esconde que não contribuiu. Todos os 13.660 bytes
do arquivo vinham da varredura automática que aquelas linhas pareciam controlar.

Ninguém notou porque o checklist compara arquivos com o que está publicado, e os
dois estavam consistentes entre si — errados juntos.

## Decisão

**`@import "tailwindcss" source(none)` nos dois arquivos de CSS fonte.** A
detecção automática fica desligada e as linhas `@source` passam a ser a lista
inteira.

- `arandu/resources/css/app.css`: as views, fonte e gerado
- `framework/view/assets/app.src.css`: as duas páginas do próprio framework —
  a de erro e a do console

## Consequências

- O build vira função das entradas. Primeira execução já é o ponto fixo.
- O `app.css` do esqueleto cai de 9.939 para 7.819 bytes e o do framework de
  13.660 para 4.849. A diferença era utilitário raspado da saída anterior, que
  nenhuma página usa.
- **Fonte nova precisa de linha nova.** É o custo, e é o ponto: uma pasta de
  views que ninguém declarou compila para nada, visivelmente, em vez de quase
  funcionar porque a varredura automática a pegou.
- O fallback do framework passa a ser só a camada base. Ele é servido a um
  projeto que apagou a própria folha de estilo, e essas páginas ficam sem
  utilitário. É coerente — o framework não tem view desde a ADR 0021 — mas
  contradiz o comentário do `viewbuild.go` que promete "as views ainda renderizam
  estilizadas". **Fica aberto:** ou o comentário muda, ou o fallback deixa de
  existir. Não foi decidido aqui.

## Como isso não volta

Três verificações, porque a decisão tem três formas de se perder:

1. `TestTheStylesheetDoesNotReadItsOwnOutput`, nos dois repositórios — o
   `@import` tem que carregar `source(none)`;
2. no framework, **toda** `@source` tem que casar algum arquivo;
3. na `prova-ponta-a-ponta.sh`, `aru view:build` roda duas vezes e o md5 tem que
   ser o mesmo. É a única que não depende de reconhecer a causa.
