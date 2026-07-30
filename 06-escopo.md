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
| Prompts | bubbletea, survey | a CLI usa uma delas |
| Dusk | chromedp, rod | documentar, não empacotar |
| Socialite | goth | adapter fino no core, se alguém pedir |
| Pulse / Nightwatch | OpenTelemetry | **exportar OTel.** Construir APM próprio é competir com Grafana e Datadog |

O caso do Pulse/Nightwatch é o mais importante: nosso Collector já captura tudo
em dev. Em produção, exportar OTel e deixar o cliente usar o Grafana que ele já
paga. `maraca` e `boitata` saem do mapa.

## CORE — dentro do framework

| Laravel | Onde vive no Arandu |
|---|---|
| Telescope (parte dev) | `observability/` — Collector e página de erro, já implementados |
| Fortify | `security/` — Grant, Policy, sessão, argon2id. Segurança é a tese; não pode ser pacote opcional |
| Sanctum | tokens de API são um adapter de sessão. `security/token.go` |
| Pennant | feature flags são um mapa e um middleware. ~200 linhas |
| Ranger | `aru doctor` já faz introspecção |

Sem nome indígena para essa camada: é o pacote `security` do core. O nome `sume` foi vetado (divindade) — ver a regra de admissibilidade em `07-pacotes.md`.

## CONSTRUIR — Fase 1

Quatro repositórios. Nada mais.

```
arandu-io/framework    a biblioteca
arandu-io/arandu       esqueleto de projeto
arandu-io/aru          CLI
arandu-io/docs         markdown do guia
```

Critério: sem esses quatro não existe produto.

## CONSTRUIR — Fase 2

```
arandu-io/postgres     adapter — sem ele o framework não roda
arandu-io/redis        adapter — sessão e rate limit distribuídos
arandu-io/porang       camada de view (templ + HTMX + Basecoat)
arandu-io/oka          starter kit: telas de auth
```

`porang` é o único com valor competitivo real da lista — é a camada que liga
componente a Grant, validação e CSRF, e nenhuma biblioteca genérica pode fazer.

## CONSTRUIR — Fase 3

```
arandu-io/marandu      console histórico de requests
arandu-io/moitara      registry de módulos
```

`moitara` só faz sentido depois de existir módulo para distribuir. Antes disso é
infraestrutura sem conteúdo.

## DEPOIS — só com pedido concreto

| Nome | Função | Gatilho para construir |
|---|---|---|
| *a definir* | busca full-text | quando full-text do Postgres não bastar para um cliente real |
| `uirapuru` / `sabia` | websocket | quando SSE do HTMX não bastar. SSE cobre a maioria dos casos |
| `potiron` | dashboard de filas | depois de existir fila em produção |
| `taba` | painel administrativo | Fase 4. É o único dos "pagos" que faz sentido |
| *a definir* | camada de IA | Fase 4, depende de tudo acima. Nomear no gatilho |
| `igarape` | roteamento por arquivo | provavelmente nunca — conflita com módulos verticais |
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
| Construir na Fase 1 | **4** |
| Construir na Fase 2 | **4** |
| Construir na Fase 3 | **2** |
| Core (sem repo) | 5 |
| Não existe em Go | 5 |
| Biblioteca de terceiro | 4 |
| Depois, com gatilho | 7 |
| Fora de escopo | 4 |

**Dez repositórios até a Fase 3**, contra 25 se o espelho fosse cego. E dos dez,
quatro são obrigatórios para existir produto.

## Reserva de nomes

Reservar na org **só** os dez que vão existir até a Fase 3, mais `taba`,
pelo risco de colisão futura. Os demais ficam documentados aqui
e reservados quando o gatilho acontecer — nome reservado sem repo vira lixo na
org e confunde quem chega.
