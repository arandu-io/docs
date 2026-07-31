# Documentação do Arandu — índice

Ordem de leitura. Cada documento decide algo que o seguinte assume como dado.

| # | Documento | O que decide |
|---|---|---|
| 01 | `01-arquitetura.md` | o caminho obrigatório, o Grant, módulo vertical, o que não portamos do Laravel |
| 02 | `02-observabilidade.md` | o pilar de logs e debug — o diferencial do produto |
| 03 | `03-roadmap-fases.md` | Fases 1 a 4, com critério de saída verificável em cada uma |
| 04 | `04-marca.md` | nome, posicionamento, tom de voz, atribuição da língua |
| 05 | `05-repositorios.md` | como os repositórios se dividem, espelhando o Laravel |
| 06 | `06-escopo.md` | **o que construir e o que não** — 10 repos, não 25 |
| 07 | `07-pacotes.md` | nomes dos pacotes, filtros aplicados e vetos |
| 08 | `08-documentacao.md` | quando documentar, e por que o site é feito no próprio Arandu |
| 09 | `09-uma-forma-so.md` | a regra anti-complexidade, e a evidência no tree do Laravel |
| 10 | `10-de-para-laravel.md` | o de-para de vocabulário; o modelo mental é Laravel + Inertia |
| 11 | `11-adapters.md` | seis pacotes por domínio; adapter não é modo |
| 12 | `12-perfis.md` | perfil convencional e perfil de performance, uma API só |
| 13 | `13-estrutura-illuminate.md` | o que o tree do `Illuminate` ensina sobre organizar o core |
| 14 | `14-view-e-assets.md` | view é binário, nunca Node |
| 15 | `15-multi-tenancy.md` | coluna discriminadora, tenant sempre vindo do `Grant` |
| 16 | `16-agendamento.md` | scheduler e lock por tenant |
| 17 | `17-deploy.md` | migration é passo de pipeline, nunca boot |
| 18 | `18-marketplace.md` | distribuição de módulos |
| 19 | `19-dsl-e-geracao.md` | a DSL é um conjunto fechado |
| 20 | `20-licenca.md` | análise de licença; a decisão está no ADR 0008 |
| 21 | `21-modelo-de-negocio.md` | onde o dinheiro entra |
| — | `packages/` | um documento por pacote do Laravel, com veredito e gatilho |
| — | `adr/` | decisões arquiteturais com justificativa |

Uma observação sobre estes arquivos: eles são os únicos, junto com os ADRs, que
ficam em português. Todo o código e toda a saída do produto são em inglês — ver
`adr/0005-idioma-do-codigo.md`.

Escopo vem antes de nomes de propósito: só se nomeia o que vai existir.

## Decisões já travadas

- Marca **Arandu**, CLI **`aru`** (o binário `arandu` já existe em outro projeto)
- Module path `github.com/arandu-io/framework` — **irreversível**
- Biblioteca em `arandu-io/framework`, esqueleto em `arandu-io/arandu`,
  documentação em `arandu-io/docs` (espelho de `laravel/framework` + `laravel/laravel`)
- Autorização por `security.Grant`: impossível de burlar em tempo de compilação
- Collector e página de erro são **core**, não plugin
- Core sem dependência de terceiros além de `golang.org/x/crypto`
- Sem verticais de nicho (NF-e, SPED, PIX) no namespace de primeira parte
- Sem equivalentes de Octane, Sail, Herd, Valet, Volt e Tinker — resolvem
  problemas do PHP
- APM em produção via OpenTelemetry, não produto próprio
- Documentação de uso é **fase própria, depois da Fase 3**. Comentário de doc em
  Go e ADR continuam obrigatórios durante o desenvolvimento
- `arandu-io/website` é uma **aplicação Arandu** — dogfood, e primeiro cliente real

- Todo o código e toda a saída do produto em **inglês**; português fica na
  conversa e nestes documentos (ADR 0005)
- A CLI é módulo próprio e delega `serve`, `migrate` e `routes` ao binário do
  projeto (ADR 0006)
- Driver de Postgres vive no esqueleto, não no core; migrador mínimo em
  `data.Migrate` até o Atlas da Fase 2 (ADR 0007)
- Licença **MIT**, a mesma do Laravel, nos quatro repositórios (ADR 0008)
- **SQLite por padrão** e Postgres em produção, com SQL portável em vez de query
  builder; migrações com lote, rollback e status, no modelo do Laravel (ADR 0009)
- Estrutura de views no formato do Inertia — `pages`, `layouts`, `components`,
  mais `partials` para os fragmentos do HTMX — na Fase 2 (ADR 0010)
- `SystemGrant` com tenant vazio **não existe**: devolve o `Grant` zero, que não
  passa no `Check`. Aplicação single-tenant usa uma constante, não vazio (REGRA 14)
- Contribuição por **DCO** (`git commit -s`), nunca CLA — o core não precisa ser
  relicenciável, porque o componente pago vive em repositório separado (doc 20)

## Decisões abertas

- `oka` como nome do starter kit — "oca" também é vazio/oco em português
- Nome da camada de view: `porang` confirmado, sem impedimento
- Verificação etimológica das palavras marcadas `[C]` e `[?]` em `07-pacotes.md`
- Tornar os quatro repositórios públicos — decidido que sim, ao fim da Fase 1
- Titular do copyright: `Paulo Lima` pessoa física, no `LICENSE.md`. Trocar para a
  pessoa jurídica se a marca for registrada em nome dela (ADR 0008)
- O esqueleto pesa 23 MB de binário por causa do SQLite puro-Go (`modernc.org`,
  que carrega um compilador C). Aceito pelo ganho de rodar sem instalar nada;
  revisitar se o tamanho incomodar em produção — quem usa só Postgres remove o
  import e o driver sai do binário

## Estado do código

Fase 1 compila e passa nos testes. Quatro repositórios criados na org, todos
privados por enquanto:

| Repositório | Estado |
|---|---|
| `framework` | 10 pacotes, `go build`, `go vet` e `go test -race` limpos |
| `aru` | CLI, módulo próprio sem dependência; `key:generate`, `serve`, `migrate`, `routes` |
| `arandu` | esqueleto; sobe, serve `/auth/login`, health e shutdown gracioso. `database/seeders` no padrão Laravel |
| `docs` | estes documentos e os ADRs |
| `examples` | a aplicação de demonstração: dois módulos escritos à mão e um tour guiado |

Critérios de saída da Fase 1 (`03-roadmap-fases.md`): 1, 2, 4 e 5 verificados por
teste automatizado. O critério 3 — um produto real rodando o login em cima do
`auth` — depende do proprietário e continua aberto.

O caminho completo — `aru migrate`, `aru db:seed`, formulário de login, senha
errada recusada, senha certa emitindo sessão — foi executado nos **dois** bancos:
SQLite, sem nada instalado, e PostgreSQL 18, trocando apenas o `.env`.


# Documentação do Arandu — índice

Ordem de leitura. Cada documento decide algo que o seguinte assume como dado.

| # | Documento | O que decide |
|---|---|---|
| 01 | `01-arquitetura.md` | o caminho obrigatório, o Grant, módulo vertical, o que não portamos do Laravel |
| 02 | `02-observabilidade.md` | o pilar de logs e debug — o diferencial do produto |
| 03 | `03-roadmap-fases.md` | Fases 1 a 4, com critério de saída verificável em cada uma |
| 04 | `04-marca.md` | nome, posicionamento, tom de voz, atribuição da língua |
| 05 | `05-repositorios.md` | como os repositórios se dividem, espelhando o Laravel |
| 06 | `06-escopo.md` | **o que construir e o que não** — 10 repos, não 25 |
| 07 | `07-pacotes.md` | nomes dos pacotes, filtros aplicados e vetos |
| 08 | `08-documentacao.md` | quando documentar, e por que o site é feito no próprio Arandu |
| 09 | `09-uma-forma-so.md` | **a regra anti-complexidade** — uma forma só de fazer cada coisa |
| 10 | `10-de-para-laravel.md` | a curva de aprendizado: Laravel → Arandu, e o fluxo Inertia sem JS |
| 11 | `11-adapters.md` | os seis pacotes de adapter, e o que conta como modo |
| 12 | `12-perfis.md` | perfil convencional e perfil de performance |
| 13 | `13-estrutura-illuminate.md` | mapa componente a componente do `Illuminate/*` |
| 14 | `14-view-e-assets.md` | camada de view, toolchain sem Node, limite do Alpine |
| — | `packages/` | um documento por pacote do Laravel: construir, não construir, e o gatilho |
| — | `adr/` | decisões arquiteturais com justificativa |

Escopo vem antes de nomes de propósito: só se nomeia o que vai existir.

## Decisões já travadas

- Marca **Arandu**, CLI **`aru`** (o binário `arandu` já existe em outro projeto)
- Module path `github.com/arandu-io/framework` — **irreversível**
- Biblioteca em `arandu-io/framework`, esqueleto em `arandu-io/arandu`,
  documentação em `arandu-io/docs` (espelho de `laravel/framework` + `laravel/laravel`)
- Autorização por `security.Grant`: impossível de burlar em tempo de compilação
- Collector e página de erro são **core**, não plugin
- Core sem dependência de terceiros além de `golang.org/x/crypto`
- Sem verticais de nicho (NF-e, SPED, PIX) no namespace de primeira parte
- Sem equivalentes de Octane, Sail, Herd, Valet, Volt e Tinker — resolvem
  problemas do PHP
- APM em produção via OpenTelemetry, não produto próprio
- Seis pacotes de adapter por domínio, driver pesado em submódulo
- **Dois perfis**: convencional (Postgres) e performance (Scylla). Mesma API,
  mesma Policy, mesmo contrato de módulo — muda a fiação gerada. Ver `12-perfis.md`
- Dragonfly como padrão de KV nos dois perfis
- Dos 14 pacotes principais do Laravel, **um** vira repositório novo. Ver `packages/`
- Documentação de uso é **fase própria, depois da Fase 3**. Comentário de doc em
  Go e ADR continuam obrigatórios durante o desenvolvimento
- View: `templ` + `tailwindcss` standalone, binário único gerenciado pela `aru`.
  Assets por `go:embed`, sem CDN e sem npm. Ver `14-view-e-assets.md`
- `arandu-io/website` é uma **aplicação Arandu** — dogfood, e primeiro cliente real
- **Uma forma só de fazer cada coisa.** Um frontend (templ + HTMX), um runtime,
  um caminho de autorização, um caminho de dados, zero Node no core
- Modelo mental para o usuário: **Laravel + Inertia, sem o build de JavaScript**

## Decisões abertas

- `oka` como nome do starter kit — "oca" também é vazio/oco em português
- Nome da camada de view: `porang` confirmado, sem impedimento
- Verificação etimológica das palavras marcadas `[C]` e `[?]` em `07-pacotes.md`
- Licença do core, antes da Fase 2

## Estado do código

Esqueleto da Fase 1 escrito e **não compilado** — não havia toolchain Go no
ambiente onde foi gerado. Rodar `gofmt -w . && go build ./...` antes de qualquer
outra coisa.
