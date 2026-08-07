# Documentação do Arandu — índice

Ordem de leitura. Cada documento decide algo que o seguinte assume como dado.

| # | Documento | O que decide |
|---|---|---|
| 01 | `01-arquitetura.md` | o caminho obrigatório, o Grant, a estrutura do Laravel, o que não portamos do Laravel |
| 02 | `02-observabilidade.md` | o pilar de logs e debug — o diferencial do produto |
| 03 | `03-roadmap-fases.md` | Fases 1 a 4, com critério de saída verificável |
| 04 | `04-marca.md` | nome, posicionamento, tom de voz, atribuição da língua |
| 05 | `05-repositorios.md` | como os repositórios se dividem |
| 06 | `06-escopo.md` | **o que construir e o que não** |
| 07 | `07-pacotes.md` | nomes dos pacotes, filtros e vetos |
| 08 | `08-documentacao.md` | quando documentar, e por que o site é feito no próprio Arandu |
| 09 | `09-uma-forma-so.md` | **a regra anti-complexidade** |
| 10 | `10-de-para-laravel.md` | a curva de aprendizado: Laravel → Arandu |
| 11 | `11-adapters.md` | os seis pacotes de adapter, e o que conta como modo |
| 12 | `12-perfis.md` | perfil **convencional** e perfil de **performance** |
| 13 | `13-estrutura-illuminate.md` | mapa componente a componente do `Illuminate/*` |
| 14 | `14-view-e-assets.md` | camada de view, toolchain sem Node, limite do Alpine |
| 15 | `15-multi-tenancy.md` | coluna discriminadora, tenant vindo do Grant |
| 16 | `16-agendamento.md` | scheduler no binário, lock distribuído, workers |
| 17 | `17-deploy.md` | binário estático, migration fora do boot |
| 18 | `18-marketplace.md` | `moitara`, selo, `arandu.mod.toml` e permissões |
| 19 | `19-dsl-e-geracao.md` | o que a DSL expressa, o que não, e a saída de emergência |
| 20 | `20-licenca.md` | **MIT e DCO** |
| 21 | `21-modelo-de-negocio.md` | **o que se vende**, e por que muda a Fase 2 |
| 22 | `22-dores-dos-devs.md` | dores em 10 frameworks, e onde estamos expostos |
| 23 | `23-estabilidade.md` | versionamento, depreciação, o que autoriza breaking change |
| 24 | `24-seguranca-processo.md` | canal, prazos, CVE, supply chain |
| 25 | `25-saida-e-nao-objetivos.md` | **quando não usar** e como sair |
| 26 | `26-revisao-externa.md` | o que absorvemos e rejeitamos das consultas externas |
| 27 | `27-eventos-e-outbox.md` | evento de domínio, outbox transacional, idempotência |
| 28 | `28-estrutura-do-projeto.md` | o esqueleto que `aru new` gera — de-para de diretório com o Laravel |
| — | `packages/` | um documento por pacote do Laravel: construir, não construir, gatilho |
| — | `adr/` | decisões arquiteturais com justificativa |

Fora da ordem de leitura, no mesmo diretório:

| Arquivo | O que é |
|---|---|
| `CONTRIBUTING.md` | DCO, o que rodar antes do PR, mensagem de commit, quando abrir ADR |
| `LICENSE.md` | a licença MIT |
| `SECURITY.md` | canal de report, prazos e política de CVE |

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
- View: `kyse` + HTMX + Alpine + Tailwind v4, com o binário standalone. O
  scaffolding de auth espelha o do `laravel/ui`. Assets por `go:embed`, sem CDN
  e sem npm. Ver `14-view-e-assets.md` e ADR 0022
- `arandu-io/website` é uma **aplicação Arandu** — dogfood, e primeiro cliente real
- **Uma forma só de fazer cada coisa.** Um frontend (kyse + HTMX), um runtime,
  um caminho de autorização, um caminho de dados, zero Node no core
- Modelo mental para o usuário: **Laravel + Inertia, sem o build de JavaScript**

- Tenancy por coluna discriminadora; tenant vem sempre do `Grant`. Ver `15`
- Scheduler embutido no binário com lock em `kv`; sem cron do sistema. Ver `16`
- Artefato de deploy é um binário; `aru deploy` não entra agora. Ver `17`
- `arandu.mod.toml` com permissões declaradas existe desde a Fase 2. Ver `18`
- DSL em YAML com JSON Schema publicado, conjunto fechado. Ver `19`
- **MIT + DCO**, a mesma do Laravel; componente pago em repositório separado. Ver `20` e ADR 0008

## A virada de agosto/2026

Quatro decisões substituíram desenho anterior. Se algum documento contradisser o
que está abaixo, o documento está velho.

1. **Estrutura idêntica à do Laravel**, diretório a diretório — `app/Http/
   Controllers`, `app/Models`, `app/Policies`, `bootstrap/`, `config/`,
   `database/`, `resources/views`, `routes/`, `storage/`, `public/`. O
   `modules/` vertical **caiu**. Motivo: reconhecimento do dev Laravel vale mais
   que elegância de organização (doc 28)
2. **Motor de template próprio: `kyse`**, do guarani *kysê*, faca — o Blade do
   Arandu. O **templ caiu**. `resources/views/home.kyse.go` é o que se escreve,
   `resources/views/home.go` é o que o `aru view:build` gera, e
   `@go`/`@endgo` é o `@php` com Go de verdade dentro (doc 14, ADR 0020)
3. **O scaffolding de auth espelha o do `laravel/ui`**: as mesmas nove views,
   nos mesmos caminhos. A stack não muda — HTMX, Alpine e Tailwind v4 (ADR 0022)
4. **O `porang` foi dissolvido.** No Laravel a view mora no esqueleto: as views
   foram para `arandu/resources/views/`, o runtime para `framework/view/`.
   De onze repositórios para dez (ADR 0021)

### A extensão termina em `.go`, e isso foi verificado

O `.blade.php` termina em `.php` porque o PHP é a linguagem hospedeira. O
`.kyse.go` termina em `.go` pelo mesmo motivo. O arquivo abre com a tag de build
que o exclui do compilador:

```go
//go:build kyse

package views

@extends('layouts.app')
```

Testado com Go 1.26.5 e cache limpo: `go build ./...`, `go vet ./...`,
`go test ./...` e `go run` passam, com o arquivo gerado no mesmo diretório e o
pacote importável. O `gofmt` é o único que ignora build tag, e o CI o chama
filtrando `*.kyse.go`.

A cláusula `package` é obrigatória — sem ela funciona apenas enquanto o
diretório não tiver nenhum outro `.go`.

### Propagação

Tabela de propagação não fica aqui. Uma lista de "quem ainda não foi
atualizado" envelhece na primeira correção e passa a mentir sobre o próprio
repositório. O estado se verifica por comando, com `./plans/checklist.sh`.

## Decisões abertas

Aguardando decisão sua:

- **A primeira receita, agora sem o `porang`** (doc 21) — a dissolução tirou a
  âncora do modelo de negócio: não há mais repositório-irmão para pendurar o
  catálogo pago, e o argumento de vender componente para quem usa templ puro caiu
  junto com o templ. Vender componente continua possível; o formato é outro, e o
  doc 21 precisa ser reescrito (ADR 0021)
- **`Collector` medir orçamento de latência** (doc 22) — recomendado; é o modo de morrer do Livewire
- Verificação etimológica das palavras marcadas `[C]` e `[?]` em `07-pacotes.md`

- Tração do GoShip não verificada (limite da API) — checar antes de citar em público

## Estado do código

Fases 1 a 4 entregues, verificadas por `./plans/checklist.sh` — que verifica por
comando, não por afirmação.

**Dez repositórios** na organização. O trio central espelha o Laravel:

| Repositório | Espelha | Papel |
|---|---|---|
| `framework` | `laravel/framework` | a biblioteca: kernel, security, data, httpx, view |
| `arandu` | `laravel/laravel` | o esqueleto de projeto que o `aru new` clona |
| `aru` | `artisan` | a CLI: `new`, `make:*`, `migrate`, `doctor`, `view:build` |
| `database` | — | `Open` e os compartimentos `pgx`, `mysql`, `sqlite` |
| `kv` | — | cache, sessão, lock e rate limit sobre RESP |
| `queue` | — | fila sobre a tabela `jobs`, e sobre RESP no submódulo |
| `storage` | — | arquivo em disco, e o protocolo S3 no submódulo |
| `examples` | — | aplicação de exemplo, compilada pelo CI |
| `homebrew-tap` | — | fórmula do `aru` |
| `docs` | `laravel/docs` | estas decisões |

O `porang` era o décimo primeiro e foi **dissolvido**: no Laravel a view mora no
esqueleto, não no framework. As views foram para `arandu/resources/views/` e o
runtime de render para `framework/view/` (ADR 0021).

A REGRA 3 vale com dez. `queue` e `storage` foram autorizados pelo proprietário
na Fase 3 (ADR 0016).
