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
| Sanctum | tokens de API | — | decidido core e **não construído**: não existe `security/token.go`. Ver `06-escopo.md` | — |
| Passport | servidor OAuth2 | — | deferido, ver `06-escopo.md` | — |
| Socialite | login por provedor externo | **jopara** | mistura | C |

> Jetstream **não tem equivalente**. Está em manutenção estendida no Laravel; o
> Breeze é o starter kit único recomendado. Seguimos o mesmo.

## Dados e busca

| Laravel | Função | Arandu | Significado | St |
|---|---|---|---|---|
| Scout | busca full-text | *a definir* | `heka` vetado no filtro zero. Deferido, ver `06-escopo.md` | — |
| Ranger | introspecção da aplicação | **epiaka** | ver, enxergar. Nome sem uso previsto: é core, no `aru doctor` (`06-escopo.md`) | C |

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
| Pennant | feature flags | **kanitar** | canitar, cocar — insígnia que se veste. Nome sem uso previsto: o veredito é core, sem repositório (`packages/pennant.md`), e o código não existe | C |

> **Octane não tem equivalente.** Ele existe porque PHP reinicializa a aplicação
> a cada request; um binário Go já boota uma vez e fica em memória. Ver
> `06-escopo.md`.

## Ambiente e deploy

| Laravel | Função | Arandu | Significado | St |
|---|---|---|---|---|
| Forge | provisionar servidor (pago) | — | fora de escopo, e sem nome. Ver a nota abaixo | — |
| Vapor | deploy serverless (pago) | **ybytu** | vento | ? |
| Cloud | plataforma gerenciada (pago) | **amana** | chuva, nuvem | C |
| Envoyer | deploy sem downtime (pago) | **piracema** | subida dos peixes rio acima | C |

*Corrigido em 07/08/2026:* o `taba` estava nomeando duas coisas nesta página —
o Forge aqui e o Nova em "Negócio e administração". Um nome, um pacote: `taba` é
o **painel administrativo** (Nova), Fase 4, e é assim que aparece em
`05-repositorios.md` e em `06-escopo.md`. O Forge fica sem nome porque está em
"NÃO — fora do escopo" do `06`: provisionar servidor é negócio de
infraestrutura, não pacote. Nomear o que não vai existir foi exatamente o que
produziu a colisão. As três linhas restantes desta tabela também estão fora de
escopo e mantêm o nome só como reserva documentada — a `06` reserva na org
apenas os dez repositórios e o `taba`.

## Negócio e administração

| Laravel | Função | Arandu | Significado | St |
|---|---|---|---|---|
| Nova | painel administrativo (pago) | **taba** | aldeia — o povoado que se administra | C |
| Cashier | assinaturas e cobrança | **repy** | preço, valor | ? |

## IA e interface de linha de comando

| Laravel | Função | Arandu | Significado | St |
|---|---|---|---|---|
| Boost | camada de IA / MCP | *a definir* | Fase 4 — nomear no gatilho, não agora | — |
| Prompts | prompts interativos de CLI | — | não existe: a `aru` é só flag e argumento, e o `aru/go.mod` tem um único `require`. Ver `06-escopo.md` | — |
| Folio | roteamento por arquivo | **igarape** | igarapé, caminho de canoa | C |

## View

| Laravel | Função | Arandu | Significado | St |
|---|---|---|---|---|
| Blade | motor de template | **kyse** | faca. O `[C]` é da etimologia, não do nome — ver abaixo | C |
| — | camada de view | — | é core: `framework/view/`, com as views do projeto em `resources/views/` | — |

O motor se chama `kyse` desde a ADR 0020, que tirou o templ da stack: a sintaxe
é a do Blade (`@extends`, `@section`, `@yield`, `{{ }}`), e o arquivo que se
escreve é `.kyse.go`, com o `.go` gerado ao lado. Não há repositório de view: o
`porang` foi dissolvido pela ADR 0021, e o nome fica livre.

*Corrigido em 07/08/2026:* o `kyse` estava marcado `[C]` e listado na tabela de
"Estado conhecido" como nome a confirmar, e continua entre as pendências do
filtro 4 — como se o nome ainda pudesse mudar. Não pode. Ele está no código
publicado, e mudá-lo quebraria arquivo-fonte de projeto de cliente:

- a extensão do fonte é `.kyse.go` — `arandu/resources/views/home.kyse.go`,
  `arandu/resources/views/layouts/app.kyse.go` e os `testdata/` do `aru doctor`
- a tag de build é `//go:build kyse`, e o gerador a emite em
  `aru/internal/gen/authviews.go`
- o motor mora em dois lugares: o compilador é `aru/internal/kyse/` (parser e
  gerador) e o runtime que a view gerada chama é `framework/view/`
- a decisão é a ADR 0020, e o `arandu.toml` registra por escrito que não há
  binário a fixar porque o `kyse` é parte da `aru`

O `[C]` que sobra mede a **etimologia**, que é o que a coluna sempre mediu:
*kysê*, faca, ainda não foi conferida no Navarro. **Encerrar ou não o filtro 4
para este nome é decisão do proprietário, e está aberta** — a diferença é que
ela não decide mais o nome, só o que este documento pode afirmar sobre a origem
dele.

## Sem equivalente no Laravel

| Função | Arandu | Significado | St |
|---|---|---|---|
| Registry de módulos (marketplace) | **moitara** | moitará, troca ritual de bens | ? |

## Detalhe por pacote

Cada um dos 14 pacotes principais do Laravel tem documento próprio em
`packages/`, com veredito, justificativa e gatilho. **Nenhum deles vira
repositório novo hoje.**

*Corrigido em 07/08/2026:* esta linha terminava com "Um só vira repositório
novo". O "um" era a UI do Telescope, o `marandu`, e ela não virou repositório: o
console histórico é core, em `framework/observability/` — `ConsolePath =
"/_arandu/debug"` em `console.go`, e o ring buffer de 200 requests em
`recorder.go` (`DefaultRecorderSize`) — e foi entregue na Fase 3. É o que a
tabela de Observabilidade desta mesma página já diz, e o `06-escopo.md` também. Os
deferidos de `06-escopo.md` só viram repositório quando o gatilho acontecer, e
cada um passa pela REGRA 3 antes. Os arquivos `packages/00-indice.md` e
`packages/telescope.md` ainda carregam a contagem antiga e o `marandu` como
pacote da Fase 3.

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

Não há exceção: npm não entra no grafo, em pacote nenhum. O que caiu foi o
**npm**, não o Tailwind. O CSS é **Tailwind v4** e continua sendo (ADR 0022); o
que saiu foi o Basecoat, que era plugin de Tailwind distribuído por npm e o
único motivo pelo qual a exceção existia.

O Tailwind não é vendorizado. É o **binário standalone**, que a `aru` baixa,
verifica contra o `sha256sums.txt` da release e cacheia em `~/.arandu/bin`, com
a versão presa em `arandu.toml` (`tailwindcss = "v4.3.3"`). O código está em
`aru/internal/toolchain/`. Binário único não é npm: não há `package.json`, não
há `node_modules`, não há lockfile de JavaScript.

O que é vendorizado, com licença e `THIRD_PARTY.md`, são os assets que entram no
binário por `go:embed`: `htmx.min.js` (0BSD) e `alpine.min.js` (MIT). O
`app.css` **não é vendorizado** — é saída do binário do Tailwind sobre fonte
nossa, e o `THIRD_PARTY.md` diz isso com todas as letras; o crédito existe
porque a saída carrega o preflight do Tailwind, que é outra coisa. O CSS-fonte fica em dois lugares — `framework/view/assets/app.src.css`
para a página de erro e o console, e `arandu/resources/css/app.css` para as
views do projeto. Os dois abrem com `@import "tailwindcss" source(none)`
(ADR 0025).

Ver `14-view-e-assets.md` e a REGRA 13 — este trecho revoga a exceção que este
documento abria antes.

*Corrigido em 07/08/2026:* o parágrafo dizia que "o Tailwind e o Basecoat caíram
junto" e que o Tailwind era "vendorizado no repositório do `arandu` com a
licença". As duas afirmações estão erradas e se contradiziam entre si — CSS
vendorizado não teria por que ter binário fixado no `arandu.toml`. O Tailwind
ficou, e chega por download verificado por SHA-256.

### Estado conhecido

| Nome | Situação |
|---|---|
| `arandu` | `.io`, `.app` e `.studio` ocupados; CLI `arandu` já existe (por isso a CLI é `aru`); org `@Arandu` e `Arandusoft` tomadas — daí `arandu-io` |
| `kyse` | não é repositório: é o compilador em `aru/internal/kyse/` mais o runtime em `framework/view/`. O nome já está no código publicado — extensão `.kyse.go` e tag `//go:build kyse` |
| `oka`, `taba`, `marandu`, `igara`, `moitara`, `potiron` | sem impedimento conhecido. Confirmar org/repo no GitHub e pkg.go.dev |

Correção registrada: uma versão anterior deste documento descartou `oka` e `taba`
por estarem ocupados em npm e PyPI. Isso estava errado — a checagem não se aplica
a Go. Os dois seguem disponíveis.

## A fazer antes de registrar

1. Checar cada `[?]` e `[C]` no *Dicionário de tupi antigo* (Navarro, 2013). O
   `kyse` está nessa lista pela **etimologia**, não pelo nome: o nome é
   já em código publicado — ver a seção View
2. Confirmar org e nome de repo no GitHub, e buscar no `pkg.go.dev`
3. Reservar na org `arandu-io` só os nomes que existirão até a Fase 3

Regra a levar para a verificação: *não quero palavra que descreva a função,
quero palavra atestada cuja acepção literal seja metáfora dela.*
