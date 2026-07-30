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

## Decisões abertas

- `oka` como nome do starter kit — "oca" também é vazio/oco em português
- Nome da camada de view: `porang` confirmado, sem impedimento
- Verificação etimológica das palavras marcadas `[C]` e `[?]` em `07-pacotes.md`
- Tornar os quatro repositórios públicos — decidido que sim, ao fim da Fase 1
- Titular do copyright: `Paulo Lima` pessoa física, no `LICENSE.md`. Trocar para a
  pessoa jurídica se a marca for registrada em nome dela (ADR 0008)

## Estado do código

Fase 1 compila e passa nos testes. Quatro repositórios criados na org, todos
privados por enquanto:

| Repositório | Estado |
|---|---|
| `framework` | 10 pacotes, `go build`, `go vet` e `go test -race` limpos |
| `aru` | CLI, módulo próprio sem dependência; `key:generate`, `serve`, `migrate`, `routes` |
| `arandu` | esqueleto; sobe, serve `/auth/login`, health e shutdown gracioso |
| `docs` | estes documentos e os ADRs |

Critérios de saída da Fase 1 (`03-roadmap-fases.md`): 1, 2, 4 e 5 verificados por
teste automatizado. O critério 3 — um produto real rodando o login em cima do
`auth` — depende do proprietário e continua aberto.

Duas pendências conhecidas: sem Docker no ambiente onde o código foi escrito, o
caminho de login não foi exercitado contra um Postgres de verdade; e o `replace`
para `../framework` no `go.mod` do esqueleto sai quando o framework virar público
e ganhar tag.
