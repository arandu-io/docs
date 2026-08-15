# Escopo real — o que construir e o que não

O Laravel tem ~25 pacotes de primeira parte porque tem 12 anos e um mercado de
milhões de desenvolvedores. Cada um nasceu porque alguém pediu.

Criar 25 repositórios no dia 1 é dívida de manutenção antes de existir usuário.
Este documento classifica cada item do mapa em `07-pacotes.md`.

## Classificação

- **CONSTRUIR** — necessidade real, repositório próprio
- **CORE** — existe, mas dentro do framework, sem repo separado
- **NÃO EXISTE** — resolve um problema de PHP/Laravel que Go não tem
- **TERCEIRO** — já existe biblioteca Go madura; usar, não reescrever
- **DEPOIS** — só quando houver pedido concreto
- **NÃO** — fora do escopo de um framework

## NÃO EXISTE — o problema é do PHP, não nosso

| Laravel | Por que não existe em Go |
|---|---|
| **Octane** | Go boota uma vez e fica em memória. É o comportamento padrão. Octane existe porque PHP reinicializa a app a cada request. |
| **Sail** | Docker de dev existe porque a stack PHP local é dolorosa (php-fpm, nginx, extensões). Em Go, `go run` resolve. Um `compose.yml` no esqueleto cobre Postgres e Redis. Não é repo, é arquivo. |
| **Herd / Valet** | Mesmo motivo, e é macOS-only. |
| **Volt** | Específico do Livewire. Não temos Livewire. |
| **Tinker** | REPL não é viável em linguagem compilada. |

Cinco itens eliminados sem custo. O nome `tata` (Octane) e `igara` (Sail) voltam
para o pool.

## TERCEIRO — já existe e é melhor que o nosso seria

| Laravel | Biblioteca Go | Decisão |
|---|---|---|
| Prompts | bubbletea, survey | **nenhuma das duas hoje.** A `aru` não tem prompt interativo. Ver a nota abaixo |
| Dusk | chromedp, rod | documentar, não empacotar |
| Socialite | goth | adapter fino no core, se alguém pedir |
| Pulse / Nightwatch | OpenTelemetry | **exportar OTel.** Construir APM próprio é competir com Grafana e Datadog |

*Corrigido em 07/08/2026:* a linha do Prompts dizia "a CLI usa uma delas". Não
usa. O `aru/go.mod` tem um único `require`, `gopkg.in/yaml.v3`, e ele está lá
pelo parser da DSL, não por prompt; `grep -rn "bubbletea\|survey" aru/` não
devolve nada, e os únicos `os.Stdin` da CLI (`delegate.go`, `dev.go`) são
repasse para subprocesso, não leitura de resposta. Ter só essa dependência é
compromisso declarado, e a ADR 0017 fixou o preço de entrar na lista: um ADR por
dependência, não um `go get`. As duas bibliotecas continuam sendo a recomendação
**se** o prompt interativo for construído. Hoje ele não existe.

O caso do Pulse/Nightwatch é o mais importante: nosso Collector já captura tudo
em dev. Em produção, exportar OTel e deixar o cliente usar o Grafana que ele já
paga. `maraca` e `boitata` saem do mapa.

## CORE — dentro do framework

| Laravel | Onde vive no Arandu |
|---|---|
| Telescope (parte dev) | `observability/` — Collector, página de erro e o console `/_arandu/debug`, já implementados |
| Blade | `view/` — o **runtime de render** do `kyse` e os assets embutidos; o **compilador** é parte da `aru`, em `aru/internal/kyse/`. As views ficam no projeto, em `resources/views/`. Ver ADR 0020 e 0021 |
| Fortify | `security/` — Grant, Policy, sessão, argon2id. Segurança é a tese; não pode ser pacote opcional |
| Sanctum | tokens de API são um adapter de sessão. **Decidido, não construído**: não existe `security/token.go` |
| Pennant | feature flags são um mapa e um middleware, ~200 linhas. **Decidido, não construído**: não existe `features/` no framework |
| Ranger | `aru doctor` já faz introspecção |

*Corrigido em 07/08/2026:* duas destas seis linhas descreviam código que não
existe, e a seção "Classificação" define CORE como "existe, mas dentro do
framework" — Sanctum e Pennant não passam nesse teste:

- `ls framework/security/` devolve `csrf.go`, `password.go`, `policy.go` e
  `session.go`. Não há `token.go`
- `ls framework/` não tem `features/`, e `grep -ril "feature.flag\|pennant"` em
  `framework/`, `aru/` e `arandu/` não devolve nada

O veredito de cada um continua sendo o de `docs/packages/`, e é ele que vale:
`pennant.md` diz "entra no core quando alguém precisar da primeira flag" — o
gatilho é a primeira flag, e ela não apareceu. `sanctum.md` diz core e traz o campo
"Gatilho para construir" com "Não se aplica." — sem completar a frase, ao
contrário do `pennant.md`, que diz de que depende. **Construir os dois, ou reclassificá-los para "depois, com
gatilho", é decisão do proprietário, e está aberta.**

As outras quatro linhas foram conferidas e estão certas: `framework/observability/`
tem `collector.go`, `console.go` e `errorpage/`, e o console responde em
`/_arandu/debug`; `framework/view/` tem o runtime e `assets/`; `security.Grant`
está em `policy.go` e o argon2id em `password.go`; `aru/internal/doctor/` existe.

Sem nome indígena para essa camada: é o pacote `security` do core. O nome `sume` foi vetado (divindade) — ver a regra de admissibilidade em `07-pacotes.md`.

## CONSTRUIR — Fase 1

Cinco repositórios. Nada mais.

```
arandu-io/framework    a biblioteca
arandu-io/arandu       esqueleto de projeto
arandu-io/aru          CLI
arandu-io/docs         markdown do guia
arandu-io/examples     a aplicação que demonstra o framework
```

Critério: sem os quatro primeiros não existe produto. O `examples` entrou junto
porque escrever à mão o módulo que o gerador vai emitir é o que transforma
"template que pareceu bom" em "template já usado".

## CONSTRUIR — Fase 2

```
arandu-io/database      adapter de SQL — sem ele o framework não roda
arandu-io/kv            adapter de chave-valor — sessão e rate limit distribuídos
arandu-io/homebrew-tap  fórmula do `aru`
```

Os nomes seguem `11-adapters.md`: um pacote por **domínio**, não por fornecedor.
`database` cobre Postgres, MySQL e SQLite; `kv` cobre Redis, Dragonfly, Valkey e
KeyDB, porque falam o mesmo protocolo.

A camada de view estava nesta lista, como `arandu-io/porang`. Não é mais
repositório: a ADR 0021 dissolveu o `porang` — as views ficam em
`arandu/resources/views/`, no projeto que o desenvolvedor edita, o runtime de
render do `kyse` mais os assets embutidos ficam em `framework/view/`, que é
core, e o compilador do `kyse` é parte da `aru` (`aru/internal/kyse/`). O valor
competitivo continua sendo o mesmo — é a camada que liga componente a Grant,
validação e CSRF, e nenhuma biblioteca genérica faz isso — e é justamente por
isso que ela é core, não dependência opcional.

O starter kit **voltou** para a lista em 07/08/2026, pela ADR 0026: é
`arandu-io/ui`, o espelho do `laravel/ui`, e a `aru` perdeu o `make:auth`. A
linha abaixo é o registro do desenho anterior.

O desenho anterior dizia: `oka` não é repositório, é o comando `aru make:auth`
(ADR 0012).

## CONSTRUIR — Fase 3

```
arandu-io/queue        fila sobre a tabela do projeto; driver RESP para volume
arandu-io/storage      disco local e protocolo S3, com R2 como padrão
```

Os dois foram autorizados pela ADR 0016.

O console histórico de requests não virou repositório: é core, em
`/_arandu/debug`, entregue na Fase 3. O nome `marandu` fica livre. O `moitara`
segue sem data — registry só faz sentido depois de existir módulo para
distribuir; antes disso é infraestrutura sem conteúdo.

## DEPOIS — só com pedido concreto

| Nome | Função | Gatilho para construir |
|---|---|---|
| *a definir* | busca full-text | quando full-text do Postgres não bastar para um cliente real |
| `joaju` | websocket | quando SSE do HTMX não bastar. SSE cobre a maioria dos casos |
| `potiron` | dashboard de filas | depois de existir fila em produção |
| `taba` | painel administrativo | Fase 4. É o único dos "pagos" que faz sentido |
| *a definir* | camada de IA | Fase 4, depende de tudo acima. Nomear no gatilho |
| `igarape` | roteamento por arquivo | provavelmente nunca — `routes/web.go` já é o caminho único de rota, e roteamento por arquivo seria o segundo (REGRA 9). Rota descoberta por convenção também some do arquivo onde se lê o que a aplicação expõe, e com ela a Policy que a protege |
| — | servidor OAuth2 | nicho pequeno, esforço enorme. Recomendar biblioteca existente |

## NÃO — fora do escopo

| Laravel | Por que não |
|---|---|
| Forge | provisionar servidor é negócio de infraestrutura: SRE, suporte 24/7, billing, SLA |
| Vapor | idem, com serverless |
| Cloud | plataforma gerenciada é o produto mais caro de operar que existe |
| Envoyer | Dokploy já resolve. Não competir com o que já funciona |

Esses quatro são empresas separadas, não pacotes. O Laravel só os tem porque a
Laravel LLC virou empresa de infraestrutura depois de o framework dominar.

## O placar

| Categoria | Itens |
|---|---|
| Construir na Fase 1 | **5** — framework, arandu, aru, docs, examples |
| Construir na Fase 2 | **3** — database, kv, homebrew-tap |
| Construir na Fase 3 | **2** — queue, storage |
| Autorizados depois da Fase 3 | **4** — ui (ADR 0026), kyse (ADR 0027), mcp (ADR 0030), hesape (ADR 0049) |
| Core (sem repo) | 6 — quatro construídos, dois decididos e não construídos |
| Não existe em Go | 5 |
| Biblioteca de terceiro | 4 |
| Depois, com gatilho | 7 |
| Fora de escopo | 4 |

**Dez repositórios até a Fase 3**, contra 25 se o espelho fosse cego. E dos dez,
quatro são obrigatórios para existir produto. Hoje a org tem **catorze**: os
quatro da quarta linha entraram depois, cada um com a autorização que a REGRA 3
exige.

## Reserva de nomes

Reservar na org **só** os catorze que existem, mais `taba`, pelo risco de
colisão futura. Os demais ficam documentados aqui
e reservados quando o gatilho acontecer — nome reservado sem repo vira lixo na
org e confunde quem chega.

*Corrigido em 14/08/2026:* o placar parava na Fase 3 e a reserva dizia dez, e as
duas linhas eram lidas como o total da org. O `ui` e o `kyse` já faltavam antes
do `mcp`, que existe desde 09/08/2026; o `hesape` entrou em 14/08/2026. A lista
nominal está em `05-repositorios.md`.
