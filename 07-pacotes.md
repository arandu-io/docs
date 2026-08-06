# Pacotes de primeira parte

Espelho do ecossistema Laravel (catálogo oficial de 2026), com nomes
tupi-guarani. Sem verticais de nicho.

## A regra de nomeação

Extraída do próprio Laravel: **palavra real cuja acepção literal é metáfora da
função — nunca palavra que descreve a função, nunca sílaba inventada.**

Pulse é "pulso" (monitorar = sentir o batimento). Pennant é "flâmula" (feature
flag = bandeirinha). Sail é "vela" (Docker = navio de contêineres). Nenhum nome
do Laravel é inventado, e nenhum é descritivo.

Consequência: só entram palavras atestadas em dicionário. Sílaba que *soa*
indígena mas não existe cria contexto falso e derruba a história da marca.

## Filtro zero — ler em voz alta em português

**Antes de qualquer análise etimológica**, dizer o nome em voz alta em português
brasileiro e checar contra gíria, palavrão e eufemismo.

Este filtro já reprovou três nomes que passaram por toda a pesquisa etimológica:

- `karai` — eufemismo de "caralho" em PT-BR ("ô carai", "que carai é isso")
- `kuatia` — sílaba inicial dita como "cu"
- `heka` — dito em voz alta é "éca", interjeição de nojo

Nenhum dos três foi pego pela cosmologia nem pelo dicionário guarani. Só pelo
ouvido. Por isso este filtro é o primeiro, não o último.

Passaram na leitura em voz alta: `arandu`, `aru`, `kyse`, `marandu`, `porang`,
`moitara`, `igara`, `taba`, `potiron`.

Na fronteira: `oka` — é a casa indígena, que todo brasileiro reconhece, mas
também é "oco/oca" no sentido de vazio. A decisão continua aberta e deixou de ser
urgente: o starter kit não é mais repositório, é `aru make:auth` (ADR 0012),
então o nome está livre e sem uso previsto.

Ordem correta de checagem:

1. gíria e palavrão em português brasileiro
2. categoria vetada (ver abaixo)
3. leitura em inglês
4. etimologia atestada
5. colisão de registry e domínio

## Regra de admissibilidade

Só entram **substantivos comuns da vida material e cotidiana**: objetos, lugares,
animais, fenômenos naturais, ações, qualidades.

Ficam fora, sem exceção:

| Categoria | Exemplos vetados |
|---|---|
| Divindades | Tupã, Karai Ru Ete, Nhamandu, Jakaira, Sumé |
| Ofícios religiosos | karaí, kunhã-karaí, pajé, xondaro |
| Etnônimos (nomes de povos) | Aruak, Guarani, Kaiowá, Tupi |
| Termos de disputa política ativa | tekoha, aty guasu |
| Textos e conceitos sagrados | ayvu rapyta |
| Entidades espirituais | curupira, boitatá, iara, saci |

Motivo prático além do ético: `karai` traduz no dicionário guarani-português como
"não-índio, senhor" — homem branco. Nomear um produto pago assim, numa marca que
se apresenta honrando a língua, é manchete pronta.

Motivo estrutural: o Laravel usa objetos e fenômenos (pulso, vela, forja,
crepúsculo, flâmula), nunca santos nem cargos religiosos. A regra é a mesma.

## Status

- **[V]** verificado em fonte
- **[C]** confiança alta, sem fonte confirmada
- **[?]** placeholder — precisa de verificação antes de qualquer uso

## Núcleo

| Papel | Laravel | Arandu | Significado | St |
|---|---|---|---|---|
| Biblioteca | `laravel/framework` | **framework** | — | — |
| Esqueleto | `laravel/laravel` | **arandu** | sabedoria, "ouvir o tempo" | V |
| Documentação | `laravel/docs` | **docs** | — | — |
| CLI | `artisan` | **aru** | forma curta de arandu | — |

## Autenticação e acesso

| Laravel | Função | Arandu | Significado | St |
|---|---|---|---|---|
| Breeze | starter kit | — | não é repositório: é `aru make:auth`, ver ADR 0012. O nome `oka` fica livre | — |
| Fortify | auth headless, sem telas | — | é core: pacote `security/` | — |
| Sanctum | tokens de API | — | é core: `security/token.go` | — |
| Passport | servidor OAuth2 | — | deferido, ver `06-escopo.md` | — |
| Socialite | login por provedor externo | **jopara** | mistura | C |

> Jetstream **não tem equivalente**. Está em manutenção estendida no Laravel; o
> Breeze é o starter kit único recomendado. Seguimos o mesmo.

## Dados e busca

| Laravel | Função | Arandu | Significado | St |
|---|---|---|---|---|
| Scout | busca full-text | *a definir* | `heka` vetado no filtro zero. Deferido, ver `06-escopo.md` | — |
| Ranger | introspecção da aplicação | **epiaka** | ver, enxergar | C |

## Filas e tempo real

| Laravel | Função | Arandu | Significado | St |
|---|---|---|---|---|
| Horizon | dashboard de filas | **potiron** | potirõ, trabalho coletivo | ? |
| Reverb | servidor de websocket | **uirapuru** | pássaro cujo canto se propaga | C |
| Echo | cliente JS de broadcasting | **sabia** | pássaro que responde o canto | C |

## Observabilidade

| Laravel | Função | Arandu | Significado | St |
|---|---|---|---|---|
| Telescope | inspetor de requests (dev) | — | é core: console `/_arandu/debug`, entregue na Fase 3. O nome `marandu` — aviso, informe, advertir [V] — fica livre | — |
| Pulse | performance da aplicação | — | exportar OpenTelemetry, ver `06-escopo.md` | — |
| Nightwatch | APM em produção | — | exportar OpenTelemetry, ver `06-escopo.md` | — |

## Qualidade e testes

| Laravel | Função | Arandu | Significado | St |
|---|---|---|---|---|
| Pint | formatador de código | — | não existe em Go: `gofmt` + `aru doctor`. Ver `packages/pint.md` | — |
| Dusk | testes de navegador | **pytuna** | noite, escuridão | C |
| Pennant | feature flags | **kanitar** | canitar, cocar — insígnia que se veste | C |

> **Octane não tem equivalente.** Ele existe porque PHP reinicializa a aplicação
> a cada request; um binário Go já boota uma vez e fica em memória. Ver
> `06-escopo.md`.

## Ambiente e deploy

| Laravel | Função | Arandu | Significado | St |
|---|---|---|---|---|
| Forge | provisionar servidor (pago) | **taba** | aldeia (o terreno onde se constrói) | C |
| Vapor | deploy serverless (pago) | **ybytu** | vento | ? |
| Cloud | plataforma gerenciada (pago) | **amana** | chuva, nuvem | C |
| Envoyer | deploy sem downtime (pago) | **piracema** | subida dos peixes rio acima | C |

## Negócio e administração

| Laravel | Função | Arandu | Significado | St |
|---|---|---|---|---|
| Nova | painel administrativo (pago) | **taba** | aldeia — o povoado que se administra | C |
| Cashier | assinaturas e cobrança | **repy** | preço, valor | ? |

## IA e interface de linha de comando

| Laravel | Função | Arandu | Significado | St |
|---|---|---|---|---|
| Boost | camada de IA / MCP | *a definir* | Fase 4 — nomear no gatilho, não agora | — |
| Prompts | prompts interativos de CLI | — | usar bubbletea, ver `06-escopo.md` | — |
| Folio | roteamento por arquivo | **igarape** | igarapé, caminho de canoa | C |

## View

| Laravel | Função | Arandu | Significado | St |
|---|---|---|---|---|
| Blade | motor de template | **kyse** | faca | C |
| — | camada de view | — | é core: `framework/view/`, com as views do projeto em `resources/views/` | — |

O motor se chama `kyse` desde a ADR 0020, que tirou o templ da stack: a sintaxe
é a do Blade (`@extends`, `@section`, `@yield`, `{{ }}`), e o arquivo que se
escreve é `.kyse.go`, com o `.go` gerado ao lado. Não há repositório de view: o
`porang` foi dissolvido pela ADR 0021, e o nome fica livre.

## Sem equivalente no Laravel

| Função | Arandu | Significado | St |
|---|---|---|---|
| Registry de módulos (marketplace) | **moitara** | moitará, troca ritual de bens | ? |

## Detalhe por pacote

Cada um dos 14 pacotes principais do Laravel tem documento próprio em
`packages/`, com veredito, justificativa e gatilho. Um só vira repositório novo.

## Triagem de escopo

Este documento define **nomes**. O que de fato será construído, e em que ordem,
está em `06-escopo.md` — dez repositórios até a Fase 3, não vinte e cinco.

## Vetos aplicados

| Palavra | Motivo |
|---|---|
| **Anauê** | saudação oficial do integralismo brasileiro, braço estendido. Reivindicada pela Frente Integralista até hoje |
| **Karai** | **"carai" é eufemismo de palavrão em PT-BR** — reprovado no filtro zero. Além disso: divindade (Karai Ru Ete), ofício religioso, e traduz como "não-índio, senhor" |
| **Kuatia** | sílaba inicial dita como "cu" em português |
| **Heka** | dito em voz alta é "éca" — interjeição de nojo em PT-BR |
| **Paje** | ofício religioso, mesma categoria de karaí |
| **Sume** | divindade / herói cultural |
| **Ayvu** | Ayvu Rapyta é texto sagrado mbyá |
| **Curupira / Boitatá** | entidades espirituais |
| **Tupã / Tupan** | divindade; apropriação religiosa |
| **Iara** | figura que afoga homens; e é nome próprio comum no Brasil |
| **Aruak / Aruan** | etnônimo — nome de povo e família linguística, não substantivo comum |
| **Itaipu** | marca consolidada (usina); colisão de reconhecimento |
| **Jaguar** | marca automotiva; colisão de trademark |
| **Tekoha** | termo central da luta territorial guarani-kaiowá; uso comercial é apropriação |
| **Yby, Ybyra, Yuxi, Yuna, Yari, Aymu** | `y` é vogal em guarani mas consoante em inglês; `Yby` não tem vogal legível para anglófono |
| **Nheenga** | `nh` impronunciável em inglês |
| **Akar, Kuar, Namu, Anira, Kambi, Tain, Arani, Arita, Moara, Guaru, Aruani, Kairu, Poti** | não reconhecidas como palavra atestada — risco de contexto falso |

## Checagem de nome — o que vale para Go

**Go não tem registry central.** O namespace de um pacote é o próprio caminho do
repositório (`github.com/arandu-io/kv`), resolvido pelo `proxy.golang.org` e
indexado no `pkg.go.dev`. Não existe nada para alguém ocupar preventivamente.

Consequência: **npm, PyPI e crates.io são irrelevantes** para nomear pacote Go.
npm é JavaScript, PyPI é Python, crates é Rust. Nenhum bloqueia nada aqui.

Onde checar, e só aqui:

| Onde | Por que | Bloqueia? |
|---|---|---|
| Org + nome de repo no GitHub | **é o module path** — o namespace real | sim, é definitivo |
| `pkg.go.dev` | outro módulo com nome parecido gera confusão e SEO ruim | não bloqueia, mas pesa |
| Fórmula do Homebrew | nome de binário duplicado não coexiste no PATH | sim, para a CLI |
| Domínio `.dev` / `.com` | marca | não bloqueia pacote |
| INPI classes 9 e 42 | marca | juridicamente sim |

Não há exceção: npm não entra no grafo, em pacote nenhum. O CSS é **Bootstrap**,
vendorizado no repositório do `arandu` com a licença e embutido por `go:embed`
(ADR 0022). O Tailwind e o Basecoat caíram junto — o Basecoat era plugin de
Tailwind distribuído por npm, e era o único motivo pelo qual a exceção existia.
Ver `14-view-e-assets.md` e a REGRA 13 — esta linha revoga a exceção que este
documento abria antes.

### Estado conhecido

| Nome | Situação |
|---|---|
| `arandu` | `.io`, `.app` e `.studio` ocupados; CLI `arandu` já existe (por isso a CLI é `aru`); org `@Arandu` e `Arandusoft` tomadas — daí `arandu-io` |
| `kyse`, `oka`, `taba`, `marandu`, `igara`, `moitara`, `potiron` | sem impedimento conhecido. Confirmar org/repo no GitHub e pkg.go.dev |

Correção registrada: uma versão anterior deste documento descartou `oka` e `taba`
por estarem ocupados em npm e PyPI. Isso estava errado — a checagem não se aplica
a Go. Os dois seguem disponíveis.

## A fazer antes de registrar

1. Checar cada `[?]` e `[C]` no *Dicionário de tupi antigo* (Navarro, 2013)
2. Confirmar org e nome de repo no GitHub, e buscar no `pkg.go.dev`
3. Reservar na org `arandu-io` só os nomes que existirão até a Fase 3

Regra a levar para a verificação: *não quero palavra que descreva a função,
quero palavra atestada cuja acepção literal seja metáfora dela.*
