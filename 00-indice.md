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
| 22 | `22-dores-dos-devs.md` | pesquisa: o que os desenvolvedores reclamam nos frameworks |
| 23 | `23-estabilidade.md` | versionamento, depreciação e o que autoriza um breaking change |
| 24 | `24-seguranca-processo.md` | como uma vulnerabilidade é reportada, tratada e divulgada |
| 25 | `25-saida-e-nao-objetivos.md` | quando **não** usar o Arandu, e como sair dele |
| 26 | `26-revisao-externa.md` | a revisão externa: o que absorvemos e o que rejeitamos |
| 27 | `27-eventos-e-outbox.md` | eventos de domínio e outbox; sem `Publish()` direto |
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
- O starter kit é `aru make:auth`, **não um repositório**: em Go o que se importa
  não se edita, e editar a tela de login é a primeira coisa que se faz (ADR 0012)
- O contrato de dados fica no core e o driver no adapter, **um módulo por
  driver**: um projeto que só fala SQLite não carrega Postgres no go.sum nem na
  superfície de vulnerabilidade (ADR 0014, substitui a 0013)
- Eventos só por **outbox**, dentro da transação da escrita. Não existe
  `Publish()` direto, e `Store` fora de `data.Transaction` é erro (doc 27)
- Todo módulo declara `arandu.mod.toml`, e o `doctor` confere a declaração de
  permissões contra o AST: usada e não declarada é erro (doc 18)
- Toolchain de view gerenciado pela `aru` em `~/.arandu/bin`, fixado em
  `arandu.toml`. `templ` e `tailwindcss` standalone, zero Node (doc 14)
- O relay do outbox faz **polling**, não `LISTEN/NOTIFY`: o caminho específico
  do Postgres poria driver no core e daria a SQLite e MySQL um segundo caminho
  de código (doc 27)
- Exportador OpenTelemetry **adiado com gatilho** — um cliente real com Grafana,
  Datadog ou Honeycomb já em produção. APM próprio segue fora de escopo para
  sempre (ADR 0015)
- Titular do copyright é a **pessoa jurídica**: HYZIS - Serviços Digitais LTDA,
  CNPJ 57.283.262/0001-28. Vale em todos os repositórios, e é em nome dela que a
  marca vai ao INPI (ADR 0008)

## Decisões abertas

- Nome da camada de view: `porang` confirmado, sem impedimento
- Verificação etimológica das palavras marcadas `[C]` e `[?]` em `07-pacotes.md`
- Redesenhar o `docs` como aplicação Arandu, e só então torná-lo público. Depende
  do `porang`, ou seja, da Fase 2
- O esqueleto pesa 23 MB de binário por causa do SQLite puro-Go (`modernc.org`,
  que carrega um compilador C). Aceito pelo ganho de rodar sem instalar nada;
  revisitar se o tamanho incomodar em produção — quem usa só Postgres remove o
  import e o driver sai do binário

## Estado do código

Fase 1 compila e passa nos testes. Cinco repositórios na org, quatro deles
públicos desde 31/07/2026 — `docs` segue privado até ser redesenhado como
aplicação Arandu:

| Repositório | Estado |
|---|---|
| `framework` | 10 pacotes, `go build`, `go vet` e `go test -race` limpos |
| `aru` | CLI, módulo próprio sem dependência; `key:generate`, `serve`, `migrate`, `routes` |
| `arandu` | esqueleto; sobe, serve `/auth/login`, health e shutdown gracioso. `database/seeders` no padrão Laravel |
| `docs` | estes documentos e os ADRs |
| `examples` | a aplicação de demonstração: dois módulos escritos à mão e um tour guiado |

Os 27 documentos estão publicados, sem buraco na sequência. Para o estado
verificado do projeto a qualquer momento: `./plans/checklist.sh`.

Critérios de saída da Fase 1 (`03-roadmap-fases.md`): 1, 2, 4 e 5 verificados por
teste automatizado. O critério 3 — um produto real rodando o login em cima do
`auth` — depende do proprietário e continua aberto.

O caminho completo — `aru migrate`, `aru db:seed`, formulário de login, senha
errada recusada, senha certa emitindo sessão — foi executado nos **dois** bancos:
SQLite, sem nada instalado, e PostgreSQL 18, trocando apenas o `.env`.
