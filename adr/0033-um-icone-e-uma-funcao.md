# ADR 0033-um-icone-e-uma-funcao — os ícones são Phosphor, vendorizados, um por função

**Status:** Aceita.

## Contexto

O `kyse` desenha botão, alerta e diálogo, e nenhum deles tem ícone. A folha do
Basecoat já dimensiona `svg` descendente de `.btn`, `.alert` e `.badge` — o buraco
é a arte, não o estilo.

O pacote escolhido é o [Phosphor](https://github.com/phosphor-icons/core), MIT,
1512 ícones em seis pesos. A REGRA 13 fecha npm e CDN, e a CSP é
`script-src 'self'`, então sobra vendorizar. A pergunta que restou é **como**:
1512 ícones são 600 KB de geometria, e uma aplicação usa trinta.

## Decisão

**`github.com/arandu-io/kyse/icons`, uma função exportada por ícone, peso
regular, gerada por `internal/icongen` a partir de um commit fixado.**

```go
{!! icons.Trash(icons.Props{Label: "Delete this post"}) !!}
```

O linker do Go prova que uma função não é alcançável e a descarta com os dados
dela. Não prova isso de uma busca em mapa, e `go:embed` não é código — é dado
vivo. Medido, com `html/template` já no binário nos quatro casos:

| binário | bytes | sobre a base |
|---|---|---|
| base, sem ícone | 3.249.106 | — |
| um ícone | 3.249.282 | **+176 B** |
| cinquenta ícones | 3.265.794 | +16.688 B |
| os 1512 | 4.431.202 | +1.182.096 B |
| **um** ícone por mapa com nome | 4.488.626 | +1.239.520 B |

A última linha é a alternativa recusada: quem desenha uma lixeira por nome paga
7043 vezes mais do que quem chama a função, porque o mapa mantém as 1512 vivas.
`go:embed` do diretório tem o mesmo perfil, com a diferença de que o erro de
nome chega em produção como quadrado em branco em vez de `undefined:` na
compilação — o mesmo argumento da ADR 0027.

**Um peso, o regular.** Seis pesos são seis vezes a fonte e, pior, `Trash` e
`TrashBold`: dois nomes para a mesma ideia e uma pergunta em toda chamada. O
conjunto de ícones é decisão do design system, como a folha de estilo, e a
REGRA 9 diz que ela se toma uma vez.

**Sem `Size` e sem `Class`.** O `svg` sai `currentColor` a `1em`, e atributo de
apresentação perde para qualquer regra CSS — o tamanho continua sendo decidido
onde já era, na folha que o projeto possui.

## Consequências

- **Só geometria atravessa a fronteira.** O gerador recusa qualquer arquivo que
  não seja o invólucro conhecido com um único `<path d="…">`, e valida o `d`
  contra o alfabeto de comandos SVG. Nenhum ícone pode carregar `<script>`,
  `onload=` ou um segundo elemento, porque nenhum ícone carrega marcação.
- **885 KB de Go gerado no repositório**, em vinte e seis arquivos por letra
  inicial. É o custo, e ele é de compilação e de diff, não de binário.
- **O pin é o commit**, não a tag: o Phosphor versiona por `package.json` (2.1.1)
  e a tag mais nova do repositório é a v2.0.8. Subir de versão é editar uma
  constante e revisar o diff.
- **A CI regenera e compara.** Ícone editado à mão compila, e a próxima geração
  o desfaz em silêncio.
- **A licença viaja.** `icons/LICENSE.md` é o MIT do Phosphor, copiado do
  checkout, sem alteração.
