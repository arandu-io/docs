# Reestruturação — `Illuminate/*` componente a componente

Mapa direto do `src/Illuminate/` do Laravel (clonado em `reference_laravel/`)
para o Arandu. É a base da frase de onboarding:

> Se você conhece Laravel, já sabe 80% do Arandu. Os outros 20% são o compilador
> cobrando o que o Laravel só recomendava.

Legenda: **core** = dentro de `arandu-io/framework` · **adapter** = pacote de
`11-adapters.md` · **stdlib** = a linguagem já resolve · **não existe** = o
problema é do PHP

## O mapa

| Illuminate | Arandu | Onde |
|---|---|---|
| **Auth** | `security/` — Grant, Policy, Subject | core |
| **Broadcasting** | SSE no core; WebSocket no `ws` | core + adapter |
| **Bus** | não há bus separado: job é job | `queue` |
| **Cache** | `kv` | adapter |
| **Collections** | `slices`, `maps`, genéricos | stdlib |
| **Concurrency** | goroutine + `errgroup` | stdlib |
| **Conditionable** | `if` | não existe |
| **Config** | `config/` — struct tipada, validada no boot | core |
| **Console** | `aru` + bubbletea | CLI |
| **Container** | fiação gerada em `internal/wire` | não existe (ADR 0001) |
| **Contracts** | interface vive no pacote que consome | não existe |
| **Cookie** | `net/http` + cookie assinado no `security/` | stdlib + core |
| **Database** | `database` + sqlc | adapter |
| **Encryption** | `crypto/*` + `security/` | stdlib + core |
| **Events** | `events/` in-process; assíncrono via `queue` | core |
| **Filesystem** | `storage` | adapter |
| **Foundation** | `kernel/` — equivalente direto | core |
| **Hashing** | `security/password.go` (argon2id) | core |
| **Http** | `net/http` + `httpx/` | stdlib + core |
| **Image** | `disintegration/imaging` | terceiro |
| **JsonSchema** | vira o schema da DSL da Fase 4 | Fase 4 |
| **Log** | `observability/` (slog) | core |
| **Macroable** | não há método dinâmico em Go | não existe |
| **Mail** | `mail` | adapter |
| **Notifications** | canal sobre `mail` e `ws` | deferido — gatilho: segundo canal além de e-mail (in-app ou push) pedido por cliente real |
| **Pagination** | `data.Query` com cursor | core |
| **Pipeline** | `httpx.Chain` — é o middleware | core |
| **Process** | `os/exec` | stdlib |
| **Queue** | `queue` | adapter |
| **Redis** | `kv` | adapter |
| **Reflection** | proibido | não existe |
| **Routing** | `httpx/router.go` | core |
| **Session** | `security/session.go` + `kv` | core + adapter |
| **Support** (Str, Arr, Carbon, Number, Uri) | `strings`, `slices`, `time`, `net/url` | stdlib |
| **Testing** | `testing` + `httptest` | stdlib |
| **Translation** | i18n mínimo com `x/text` | deferido — gatilho: primeiro cliente fora do par pt-BR/en. As mensagens do framework nascem em inglês |
| **Validation** | `validation/` — `Validate() Errors` | core |
| **View** | `view/` — motor `kyse`, assets embutidos, HTMX. As views ficam no projeto, em `resources/views/` | core |

## O placar

| Destino | Componentes |
|---|---|
| core | 13 |
| adapter | 7 |
| stdlib resolve | 9 |
| não existe (problema do PHP) | 5 |
| deferido / terceiro / Fase 4 | 4 |

**Nove componentes do Laravel desaparecem porque a stdlib do Go já os tem**, e
cinco desaparecem porque o problema não existe fora do PHP.

A `View` trocou de coluna: era adapter, quando havia um repositório separado
para a camada de view. A ADR 0021 dissolveu esse repositório e o motor `kyse`
(ADR 0020) ficou em `framework/view/`, que é core — um adapter a menos, um core
a mais.

Isso responde por que o core do Arandu tem duas dependências e o Laravel tem
dezenas: metade do `Illuminate` é a camada que, em Go, vem com a linguagem.

## Os cinco que não existem, e por quê

| Componente | Por quê |
|---|---|
| **Container** | reflection em runtime; substituído por fiação gerada e legível |
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
