# Os 14 pacotes principais do Laravel — um documento por projeto

Cada arquivo é auto-contido: dá para entregar um deles a alguém sem contexto e a
pessoa saber o que fazer.

## Placar

| Veredito | Quantos |
|---|---|
| Não existe em Go — o problema é do PHP | **3** |
| Não construir — biblioteca de terceiro é melhor | **2** |
| Core, sem repositório próprio | **2** |
| Não é projeto — vive dentro de outro | **1** |
| Deferido, com gatilho escrito | **5** |
| Construir na Fase 3 | ~~1~~ **0** |

**Nenhum repositório novo saiu desta lista de catorze.**

*Atualizado em 07/08/2026:* o único candidato era o `marandu`, a UI do Telescope,
e ele foi entregue **como core** — console, timeline, `aru trace` e detector de
N+1 vivem em `framework/observability/`. O repositório foi dissolvido antes de
existir, e o nome está livre. Nenhum gatilho dos outros treze disparou.

## A lista

| Pacote | Veredito | Nome | Quando |
|---|---|---|---|
| [Telescope](telescope.md) | dividido: captura é core, UI é pacote | `marandu` | Fase 3 |
| [Pulse](pulse.md) | exportar OpenTelemetry | — | Fase 3, no core |
| [Octane](octane.md) | não existe em Go | — | nunca |
| [Pint](pint.md) | não existe em Go | — | nunca |
| [Sail](sail.md) | não existe em Go | — | nunca |
| [Dusk](dusk.md) | usar chromedp | — | nunca |
| [Socialite](socialite.md) | usar goth, adapter fino | — | com demanda |
| [Sanctum](sanctum.md) | core | — | Fase 1–2 |
| [Pennant](pennant.md) | core | — | quando houver a 1ª flag |
| [Echo](echo.md) | vive na camada de view | — | com o `uirapuru` |
| [Scout](scout.md) | deferido | a definir | Postgres não bastar |
| [Horizon](horizon.md) | deferido | `potiron` | fila em produção |
| [Reverb](reverb.md) | deferido | `uirapuru` | SSE não bastar |
| [Cashier](cashier.md) | deferido | a definir | Fase 4 |

## Correção em relação ao mapa anterior

**Pint saiu.** Estava nomeado `katu` e classificado como pacote a construir. Errado:
o Pint existe porque o PHP não tem formatador canônico. Go tem `gofmt` na
toolchain desde o dia um, sem configuração e sem debate de estilo. O que sobrava
de útil — checar regras de arquitetura, não estilo — já é o `aru doctor`.

Mesmo raciocínio de Octane e Sail: são soluções para problemas da linguagem, não
do domínio.

## Como ler

Três perguntas resolvem quase todo caso:

1. **O problema existe em Go?** Octane, Sail e Pint morrem aqui.
2. **Já existe biblioteca Go madura?** Dusk, Socialite e Pulse morrem aqui.
3. **É pequeno o bastante para ser core?** Pennant e Sanctum morrem aqui.

O que sobrevive às três é candidato a repositório — e mesmo assim precisa de
gatilho escrito antes de existir.
