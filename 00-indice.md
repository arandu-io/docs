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
| 11 | `11-adapters.md` | os pacotes de adapter por domínio — quatro existem —, e o que conta como modo |
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
- Adapter por **domínio**, nunca por fornecedor, com driver pesado em submódulo.
  Existem **quatro**: `database`, `kv`, `queue` e `storage`. O `11-adapters.md`
  desenha seis — `mail` e `ws` continuam sem gatilho ("com demanda") e sem
  repositório
- **Dois perfis**: convencional (Postgres) e performance (Scylla). Mesma API,
  mesma Policy, mesmo contrato de módulo — muda a fiação gerada. Ver `12-perfis.md`
- Dragonfly como padrão de KV nos dois perfis
- Dos 14 pacotes principais do Laravel, **nenhum** virou repositório novo. Ver
  `packages/`

  *Corrigido em 07/08/2026:* esta linha dizia que **um** viraria. O candidato
  único era o `marandu`, a UI do Telescope, e ele não existe: o console
  histórico de requests ficou core, em `framework/observability/console.go`,
  servido em `/_arandu/debug` desde a Fase 3 — ver `06-escopo.md`. Os outros
  treze vereditos se dividem em cinco categorias, não três — "não existe em Go"
  (3), "não construir" (2), "core" (2), "não é projeto" (1) e "deferido com
  gatilho" (5) —, e nenhum gatilho disparou. O placar em `packages/00-indice.md` ainda conta um
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

O que o `./plans/checklist.sh` respondeu em 07/08/2026, por critério de saída do
`03-roadmap-fases.md`:

| Fase | Estado |
|---|---|
| 1 | espera pessoa — 30 de 31; o critério que falta não é de código |
| 2 | fechada — 13 de 13 |
| 3 | espera pessoa — 3 de 5; os dois que faltam não são de código |
| 4 | fechada — 5 de 5 |

Os três critérios que faltam dependem de pessoa, não de commit: um produto real
rodando o login sobre o auth (Fase 1), um bug real diagnosticado só pelo console
e um dev externo que instala, quebra e entende o erro sem documentação (Fase 3).

Nenhuma verificação falhou na rodada, e duas do instrumento foram consertadas
no mesmo dia para que isso fosse verdade.

A primeira: a seção *O que os documentos afirmam, e o que existe* procura
`aru <nome>` nos documentos e cobra que o nome esteja no `help`, com uma lista
de exceções que era **escrita à mão**. Ela não distinguia documento que
**promete** um comando de documento que o cita **para dizer que não existe** —
então cada correção honesta feita aqui aumentava a falha. A lista saiu: hoje o
comando ou existe no `help`, ou o documento diz na mesma linha que ele não
existe. O atrito continua, só que no documento, que é onde alguém lê.

A segunda: a guarda que cobra que todo nome de regra citado exista no `doctor`
pegava qualquer palavra com dois hifens entre crases — e pegou dois nomes de
servidor DNS no doc 04, o que deixou a Fase 2 reportando **aberta**. Agora ela
só olha linhas que falam de regra ou do doctor.

As duas eram do mesmo feitio: instrumento que reprova o documento por dizer a
verdade empurra o documento a calar sobre a lacuna.

Duas verificações ficam de fora enquanto não houver banco de pé — conformidade
em Postgres e em MySQL, por `./plans/testar-com-bancos-reais.sh`.

*Corrigido em 07/08/2026:* esta linha dizia "Fases 1 a 4 entregues, verificadas
por `./plans/checklist.sh`". O checklist não diz isso: ele fecha as Fases 2 e 4
e marca as Fases 1 e 3 como "espera pessoa". Citar o instrumento e afirmar o
contrário do que ele responde é pior do que não citá-lo. O estado se lê rodando
o comando; a tabela acima é uma fotografia, e envelhece.

**Catorze repositórios** na organização. O trio central espelha o Laravel:

| Repositório | Espelha | Papel |
|---|---|---|
| `framework` | `laravel/framework` | a biblioteca: kernel, security, data, httpx, view |
| `arandu` | `laravel/laravel` | o esqueleto de projeto que o `aru new` clona |
| `aru` | `artisan` | a CLI: `new`, `make:*`, `migrate`, `doctor`, `view:build` |
| `ui` | `laravel/ui` | o starter kit de auth, publicado no projeto (ADR 0026) |
| `kyse` | — | a biblioteca de componentes, importada e não copiada (ADR 0027) |
| `hesape` | `illuminate/*` | a coleção de componentes que o `framework` importa (ADR 0049) |
| `mcp` | `laravel/mcp` | servidor Model Context Protocol, com `Grant` em toda ferramenta (ADR 0030) |
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

A REGRA 3 vale com catorze, e cada passo acima de dez tem autorização registrada:
`ui` (ADR 0026), `kyse` (ADR 0027), `mcp` (ADR 0030) e `hesape` (ADR 0049).
`queue` e `storage` foram autorizados pelo proprietário na Fase 3 (ADR 0016).

*Corrigido em 14/08/2026:* esta seção dizia dez e listava dez. Já estava errada
antes do `hesape`: o `ui` e o `kyse` nunca entraram na tabela, e o `mcp` existe
desde 09/08/2026 com `go.mod` próprio e ADR própria. A lista que vale é a de
`05-repositorios.md`.
