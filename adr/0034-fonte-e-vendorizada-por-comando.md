# ADR 0034-fonte-e-vendorizada-por-comando — `aru font:add`

**Status:** Aceita.

## Contexto

A CSP é `default-src 'self'` e a REGRA 13 fecha CDN e npm. Isso deixava a
tipografia com uma saída só: um empilhamento de fontes locais.

O `examples` fez exatamente isso, e o resultado é o defeito que motivou esta ADR:

```css
--font-display: "Iowan Old Style", "Palatino Linotype", Palatino,
  "Book Antiqua", "Hoefler Text", Georgia, "Times New Roman", serif;
```

Isso é **Iowan Old Style no Mac, Georgia no Windows, Liberation Serif no Linux,
Times no Android**. Quatro blogs diferentes, e a identidade decidida pelo sistema
operacional de quem lê.

A alternativa manual — baixar um `.woff2`, largar num diretório e escrever o
`@font-face` à mão — erra duas coisas que ninguém acerta de cabeça: o
`unicode-range`, que precisa bater com o que o arquivo cobre, e as métricas do
fallback, que são três números que não se adivinham.

## Decisão

**A fonte entra por comando, uma vez, e vira arquivo do projeto.**

```sh
aru font:add "Young Serif" --as display
aru font:add "Public Sans" --as body --weight 400..700 --subset latin,latin-ext
aru font:list
aru font:remove display
```

Mesma forma do binário do Tailwind que a `aru` já fixa: baixa, verifica,
registra, e a partir daí é um arquivo. **Nada é buscado em build, e nada em
runtime.**

O comando escreve cinco coisas:

| | |
|---|---|
| `assets/fonts/*.woff2` | os bytes, commitados |
| `assets/fonts/OFL-*.txt` | a licença, ao lado do que ela cobre |
| `assets/fonts.go` | `go:embed` + `fonts.Register`, gerado |
| `resources/css/fonts.css` | `@font-face`, fallback métrico, `--font-display` |
| `arandu.toml` `[fonts.*]` | família, peso, subset, licença, métricas |

### Cinco decisões dentro dela

**Só `latin` por padrão.** `U+0000-00FF` cobre todo acento que português e inglês
usam. `latin-ext` é polonês, tcheco e turco — e dobra os bytes por um idioma que
ninguém pediu. Um teste fixa esse fato.

**Dois papéis, `display` e `body`, e o conjunto é fechado.** Um terceiro seria um
terceiro arquivo em todo binário; a resposta para "e legenda?" é um peso, não uma
família.

**Fallback com métrica casada.** `font-display: swap` garante reflow: a página
pinta no sistema e cada linha se move quando a fonte chega. O comando lê
`head`/`hhea`/`OS/2` do TrueType — sessenta linhas de stdlib, sem parser de
terceiro — e emite `ascent-override`, `descent-override`, `line-gap-override` e
`size-adjust`. A troca muda os desenhos sem mover nada. Sem métrica legível,
**nenhum override é escrito**: um número inventado desloca numa direção que
ninguém mediu.

**O `src:` carrega o hash da própria fonte.** Uma `url()` relativa resolve contra
`/_arandu/assets/<hashDoCSS>/` e **alcança** o arquivo — mas o handler compara o
hash do caminho com o do conteúdo e serve divergência com `no-cache`, de
propósito, para que referência velha degrade em página lenta e não em página
quebrada. Relativa, portanto, é a fonte rebaixada **a cada visita**, em silêncio.
`view.AssetHash` virou exportada por isso: é contrato atravessando fronteira de
repositório, a CLI não pode importar o framework, e os dois lados têm teste
fixando a mesma resposta.

**A licença viaja junto.** A OFL exige, e binário com fonte embutida é
redistribuição. Um comando que vendoriza os bytes e deixa a licença para trás
coloca em falta todo projeto que o usou.

### O que continua fora

Sem hospedagem de fonte, sem subsetter próprio e sem descoberta automática de
peso. Subsetter exigiria compilador de fonte no toolchain; os subsets publicados
já são por escrito, que é o corte que importa.

## Consequências

- O corpo continua no stack do sistema: **zero bytes**, renderiza no primeiro
  paint. Só o display é vendorizado — é onde a identidade mora, e são poucas
  palavras por página.
- O `examples` custa **18 KB** e desenha igual em todo sistema operacional.
- **A camada de fonte mora no `kyse`, não no framework.** Quais formatos
  entram, quanto um `.woff2` custa contra um `.ttf` e o que um
  `<link rel=preload>` precisa dizer são decisões tipográficas — são da pilha
  que desenha a página. O framework ficou com `view.RegisterAsset`, que é nome,
  content type e bytes, sem opinião sobre o que carrega.
- `kyse/fonts.Preloads()` emite o `<link rel=preload>` de cada face registrada,
  com `crossorigin` — sem ele o preload é uma **segunda** requisição em vez de
  cache quente, que é a forma mais comum de um preload deixar a página mais
  lenta.
- `aru font:search` e `aru font:info` tornam o catálogo navegável: 1942 famílias,
  e um comando que só aceita o nome exato é um comando que exige um navegador
  aberto do lado.
- `aru font:add --file` aceita face própria — `.woff2`, `.ttf` ou `.otf` — com a
  mesma fiação. `--family` é obrigatório, porque nome de arquivo não é nome de
  família.
- `font-src 'self'` foi escrito na CSP. O `default-src` já cobria; explicitar faz
  a linha dizer o que vale, para quem for perguntar se uma URL do Google
  funcionaria.
