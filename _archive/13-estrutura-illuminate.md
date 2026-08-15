# Reestruturação — `Illuminate/*` componente a componente

Mapa direto do `src/Illuminate/` do Laravel (clonado em `reference_laravel/`)
para o Arandu. É a base da frase de onboarding:

> Se você conhece Laravel, já sabe 80% do Arandu. Os outros 20% são o compilador
> cobrando o que o Laravel só recomendava.

Legenda: **core** = dentro de `arandu-io/framework` · **adapter** = pacote de
`11-adapters.md` · **CLI** = dentro de `arandu-io/aru` · **stdlib** = a
linguagem já resolve · **não existe** = o problema é do PHP

A coluna **Onde** diz o destino arquitetural, não o estado de entrega. Quatro
adapters estão publicados (`database`, `kv`, `queue`, `storage`); `mail` e `ws`
seguem decididos e sem gatilho, o que deixa as linhas **Mail** e a metade
WebSocket de **Broadcasting** como desenho.

## O mapa

| Illuminate | Arandu | Onde |
|---|---|---|
| **Auth** | `security/` — Grant, Policy, Subject | core |
| **Broadcasting** | SSE no core e WebSocket no `ws`. **Nenhum dos dois existe.** O `ws` está decidido e **sem gatilho** — ADR 0016, "com demanda"; o gatilho que os docs 05 e 06 registram é do `joaju`, e é "quando o SSE do HTMX não bastar" | core + adapter |
| **Bus** | não há bus separado: job é job; o pacote é o `queue` | adapter |
| **Cache** | `kv` | adapter |
| **Collections** | `slices`, `maps`, genéricos | stdlib |
| **Concurrency** | goroutine, `sync` e `context`; o `errgroup` é `x/sync` e nada no core o importa | stdlib |
| **Conditionable** | `if` | não existe |
| **Config** | `config/` — struct tipada, validada no boot | core |
| **Console** | `aru` sobre o `flag` da stdlib | CLI |
| **Container** | fiação explícita em `bootstrap/app.go`, escrita à mão | não existe (ADR 0001) |
| **Contracts** | interface vive no pacote que consome | não existe |
| **Cookie** | `net/http` + cookie assinado no `security/` | stdlib + core |
| **Database** | `database` + `database/sql` | adapter |
| **Encryption** | `crypto/*` + `security/` | stdlib + core |
| **Events** | `events/` in-process; assíncrono via `queue` | core |
| **Filesystem** | `storage` | adapter |
| **Foundation** | `kernel/` — equivalente direto | core |
| **Hashing** | `security/password.go` (argon2id) | core |
| **Http** | `net/http` + `http/` | stdlib + core |
| **Image** | `disintegration/imaging` | terceiro |
| **JsonSchema** | vira o schema da DSL da Fase 4 | Fase 4 |
| **Log** | `observability/` (slog) | core |
| **Macroable** | não há método dinâmico em Go | não existe |
| **Mail** | `mail` | adapter |
| **Notifications** | canal sobre `mail` e `ws` | deferido — gatilho: segundo canal além de e-mail (in-app ou push) pedido por cliente real |
| **Pagination** | `data.Query` com cursor | core |
| **Pipeline** | `http.Chain` — é o middleware | core |
| **Process** | `os/exec` | stdlib |
| **Queue** | `queue` | adapter |
| **Redis** | `kv` | adapter |
| **Reflection** | proibido | não existe |
| **Routing** | `http/router.go` | core |
| **Session** | `security/session.go` + `kv` | core + adapter |
| **Support** (Str, Arr, Carbon, Number, Uri) | `strings`, `slices`, `time`, `net/url` | stdlib |
| **Testing** | `testing` + `httptest` | stdlib |
| **Translation** | i18n mínimo com `x/text` | deferido — gatilho: primeiro cliente fora do par pt-BR/en. As mensagens do framework nascem em inglês |
| **Validation** | `validation/` — `Validate() Errors` | core |
| **View** | `view/` — motor `kyse`, assets embutidos, HTMX. As views ficam no projeto, em `resources/views/` | core |

*Corrigido em 07/08/2026:* a linha do **Console** dizia "`aru` + bubbletea". A
`aru` não usa bubbletea nem survey: o `aru/go.mod` tem **um** require,
`gopkg.in/yaml.v3`, e todo subcomando lê argumento com o `flag` da stdlib. O
único require é compromisso declarado na própria `go.mod` — a CLI pode se dar ao
luxo de uma dependência, e quem importa o framework não paga por ela (ADR 0006).
Bubbletea e survey não entram.

*Corrigido em 07/08/2026:* a linha do **Broadcasting** prometia "SSE no core;
WebSocket no `ws`" sem dizer que nenhum dos dois foi escrito, e sem gatilho.
Não há SSE: `text/event-stream` não aparece em nenhum `.go` do `framework`, e o
repositório `ws` não existe (ADR 0016 deixou `mail` e `ws` "com demanda"). O
detalhe que fecha a porta pelo lado do navegador: `framework/view/assets/`
embute exatamente três arquivos — `htmx.min.js`, `alpine.min.js` e `app.css` —
e a extensão SSE do HTMX não está entre eles. Como a CSP é `script-src 'self'`,
não há como puxá-la de CDN. O documento passa a registrar a ausência em vez
de contar com o mecanismo — sem inventar gatilho, porque o registro já tem os
seus: o `ws` está "com demanda" (ADR 0016) e o `joaju` espera o SSE do HTMX
não bastar (docs 05 e 06).

## O placar

São **38 linhas** no mapa, e cada uma conta uma vez. Linha com dois destinos
conta no **primeiro** que aparece: `stdlib + core` conta em stdlib,
`core + adapter` conta em core. É a regra que torna o placar recontável.

| Destino | Componentes |
|---|---|
| core | 13 |
| adapter | 7 |
| stdlib resolve | 8 |
| não existe (problema do PHP) | 5 |
| deferido / terceiro / Fase 4 | 4 |
| CLI | 1 |
| **total** | **38** |

*Corrigido em 07/08/2026:* o placar dizia 9 em stdlib e não tinha a linha da CLI,
somando 38 por compensação entre um erro para mais e uma omissão. A conta real,
pela regra acima: 11 linhas `core` + 2 `core + adapter` = 13; 7 `adapter`;
5 `stdlib` + 3 `stdlib + core` = 8; 4 `não existe` + 1 `não existe (ADR 0001)`
= 5; 1 `terceiro` + 1 `Fase 4` + 2 `deferido` = 4; 1 `CLI`. Somam 38, que é o
número de linhas.

**Oito componentes do Laravel caem na stdlib do Go.** Cinco desaparecem por
inteiro — Collections, Concurrency, Process, Support e Testing — e três (Cookie,
Encryption e Http) ficam na stdlib com uma camada fina de core em cima. Outros
cinco desaparecem porque o problema não existe fora do PHP.

A `View` trocou de coluna: era adapter, quando havia um repositório separado
para a camada de view. A ADR 0021 dissolveu esse repositório e o motor `kyse`
(ADR 0020) ficou em `framework/view/`, que é core — um adapter a menos, um core
a mais.

Isso responde por que o core do Arandu tem **uma dependência direta** —
`golang.org/x/crypto`, mais o `golang.org/x/sys` que ela arrasta como indireta —
e o Laravel tem dezenas: metade do `Illuminate` é a camada que, em Go, vem com a
linguagem.

## Os cinco que não existem, e por quê

| Componente | Por quê |
|---|---|
| **Container** | reflection em runtime; substituído por fiação escrita à mão em `bootstrap/app.go` — ver a linha do Container acima, na mesma tabela, que sempre esteve certa |
| **Contracts** | Go põe interface no consumidor, não num pacote central de contratos |
| **Macroable** | adicionar método a uma classe em runtime. Impossível e indesejável |
| **Conditionable** | `->when()` existe porque encadeamento fluente não aceita `if` |
| **Reflection** | é o mecanismo que a tese do produto rejeita |

## Correção sobre o Symfony

O Laravel **não substituiu o Symfony — roda sobre ele.** No tree clonado:
`Illuminate/Session/SymfonySessionDecorator.php`, `Foundation/Exceptions/Whoops/`,
e todo o `Illuminate/Console` sobre o Symfony Console.

Consequência arquitetural: **a stdlib do Go é o Symfony do Arandu.** `net/http`,
`html/template`, `database/sql`, `log/slog`, `crypto`. Você não está construindo
só o Laravel — está construindo Laravel e Symfony ao mesmo tempo, e a metade
Symfony já veio com a linguagem.

## O que NÃO copiar do tree

Evidências de complexidade que o Arandu não herda:

- **Três frontends**: Blade, Livewire (`reverb/src/Pulse/Livewire/`) e Inertia
  (`octane/src/Listeners/PrepareInertiaForNextOperation.php`)
- **Três runtimes**: Swoole, RoadRunner e FrankenPHP, cada um com client,
  `ServerStateFile` e `ProcessInspector` próprios
- **Node no core**: `Support/NodePackageManagers/` tem `Bun.php`, `Npm.php`,
  `Pnpm.php`, `Yarn.php`
- **Build npm na página de erro**: `Foundation/resources/exceptions/renderer/`
  tem `package.json`, `vite.config.js` e `package-lock.json`

Ver `09-uma-forma-so.md`.
