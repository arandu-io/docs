# 26 — Revisão externa: o que absorvemos e o que rejeitamos

Consultas paralelas ao ChatGPT e ao Grok sobre a mesma arquitetura. Registro aqui
o que entra, o que fica de fora e o porquê — para não reabrir depois.

---

## 1. Correção de fato: o Goldr existe

O Grok afirmou não ter encontrado o projeto e classificou como alucinação da
pesquisa anterior. Verificado pela API do GitHub nesta mesma sessão:

```
mobiletoly/goldr — 89 stars — Apache-2.0 — criado mai/2026 — push jun/2026
```

É o projeto com **mais tração do nicho** hoje, acima de gothicframework (17) e
gojang (16). Continua valendo como contra-argumento à arquitetura obrigatória
(doc 06), não como concorrente.

Aporte novo e válido do Grok: **GoShip** (`leomorpho/goship`), evolução do
Pagoda, com Go + HTMX + Templ + Ent + Echo + Tailwind + Stripe + Kamal. Não
consegui verificar tração (limite da API); checar antes de citar em qualquer
lugar público.

---

## 2. ABSORVIDO — Transactional Outbox

Lacuna real nossa. O doc 13 coloca Events no core e nunca especificou
consistência. O padrão está certo, e a formulação opinativa também: **não existe
`Publish()` direto**; só `outbox.Store()` dentro da transação; o relay faz parte
do runtime. Isso é a REGRA 9 aplicada a eventos.

Ajustes que fazemos por cima:

- **O registro no outbox carrega o `Grant` que autorizou a escrita.** Auditoria
  completa de graça: quem autorizou, qual ação, qual tenant.
- **O `Collector` mede o atraso do relay.** Evento pendente e latência de
  publicação entram na timeline (doc 02) e nos `hints` da página de erro.
  "Evento pendente há 4 minutos" é diagnóstico, não métrica de dashboard.
- **Sem "modo sync", nem para teste.** O teste usa o mesmo caminho, com o relay
  executado inline pelo helper de teste. Modo é segunda forma (REGRA 9).
- **Fase:** o relay e o NATS ficam na Fase 3, mas a **tabela outbox e o
  `Store()` entram na Fase 2**, com driver sobre `database`. Evento que não foi
  gravado não dá para recuperar depois — igual ao `arandu.mod.toml`.

Vira `docs/27-eventos-e-outbox.md`.

*Corrigido em 07/08/2026:* o registro de fase acima é de quando a decisão foi
absorvida, e envelheceu nos dois itens. O relay existe
(`framework/events/relay.go`, com teste contra banco real em
`arandu/relay_test.go`), e o NATS deixou de ser item a construir: o contrato é
`events.Publisher`, e a única ocorrência da palavra "NATS" no código é um
comentário em `relay.go` dizendo que o framework não escolhe a fila. Pôr esse
driver no core contraria a ADR 0004. O estado se lê no doc 27, não aqui.

---

## 3. ABSORVIDO EM PARTE — CQRS

Separar comando de consulta é barato e bom, e já estamos perto: o
`data.Repository` separa `Find`/`List` de `Create`/`Update`/`Delete`.

**Entra:** nomes de caso de uso separados por intenção, e read model quando a
query pesar de verdade.

**Não entra, e este é o ponto mais perigoso de toda a consulta externa:**

> "Queries podem ter Policy de leitura (opcional)."

Opcional não existe aqui. Se o caminho de leitura pula a Policy, metade da
aplicação fica sem enforcement e a tese do produto morre — e num SaaS
multi-tenant isso é vazamento de dado entre clientes, não sutileza acadêmica.

`data.Repository.List` exige `Grant` hoje e continua exigindo em qualquer
desenho de CQRS. **Read repository é repositório:** exige `Grant`, filtra por
`data.Tenant(g)`, ponto final.

**Também não entra:** `CommandHandler[C,R]` e `QueryHandler[Q,R]` genéricos como
camada de despacho. É indireção sem ganho — em Go a chamada direta do service é
tipada, rastreável e aparece no stack. Bus de comando é a facade do mundo
tipado: você troca "quem chama quem" por "descubra em runtime".

---

## 4. REJEITADO — a arquitetura de cinco camadas

O Grok propõe Domain / Application / Infrastructure / Transport / Presentation,
com Value Objects, Aggregates, Commands, Queries, DTOs e Use Cases.

**A proposta contradiz o próprio levantamento dele.** O mesmo texto aponta
"muito boilerplate", "DI excessivo" e "over-engineering para apps médios" como a
dor principal do NestJS e do Spring — e então propõe uma estrutura com **mais**
camadas que os dois.

A conta:

| | Módulo `customer` |
|---|---|
| Proposta externa | 5 subdiretórios, ~15 arquivos, interface de repositório duplicada (write e read) |
| Nosso desenho | 12 arquivos, cada um no diretório do Laravel que já existe no projeto |

*Corrigido em 07/08/2026:* a linha dizia 9 arquivos. São **12**, conferidos
rodando o gerador:
`aru make:module purchase_order --fields "title:string!,amount:money" --dry-run`
lista sete arquivos em `app/`, uma migração em `database/migrations/` e quatro
views `.kyse.go` em `resources/views/`. As telas entraram depois (ADR 0023) e a
conta nunca foi refeita. A diferença encolheu, e o argumento não dependia dela: o
que se compara é onde cada arquivo mora, não quantos são.

E o custo maior não é o número de arquivos: **o dev Laravel não reconhece nada
ali.** Aggregate, Use Case, Value Object e Command Bus não existem no vocabulário
dele. A frase de onboarding — "se você conhece Laravel, você já sabe 80% do
Arandu" — morre no primeiro `ls`. REGRA 10 e doc 10 seguem valendo: vocabulário
Laravel, e a estrutura de diretórios do Laravel (ADR 0019).

**Onde há razão:** Value Object para dinheiro, CPF/CNPJ e e-mail é bom e evita
classe inteira de bug. Cabe como tipo em `app/Models/`, junto da entidade, sem
virar camada nem diretório.

---

## 5. REJEITADO — Event Sourcing opt-in por módulo

"Opt-in" é modo com outro nome, e aqui é pior que o Octane: já temos dois perfis
(doc 12). Somar sourcing por módulo dá três modelos de dados e a combinação
perfil × sourcing — cada módulo declarando duas dimensões, e o gerador emitindo
quatro variantes de cada arquivo.

**Gatilho escrito**, para não reabrir por entusiasmo: exigência regulatória de
histórico imutável **como fonte da verdade**. Auditoria não é gatilho — o audit
log da camada de service já resolve auditoria.

---

## 6. CORREÇÃO ao documento do ChatGPT — a promessa de migração

O texto afirma que se migra de um perfil para o outro — ele chama os dois de
"Standard" e "High Performance" — "sem alterar a estrutura da aplicação ou a
forma de desenvolvimento, substituindo apenas os componentes responsáveis pela
persistência e comunicação".

Isso não é verdade, e o doc 12 já explica: Scylla é CQL, sem join, sem transação
cruzando partição, e o tenant vira chave de partição. No desenho, o módulo
declara `Profiles()` e o `doctor` reprova, sob a flag de perfil, módulo que usa
JOIN, transação cruzando agregado ou entidade sem chave de partição.

*Corrigido em 07/08/2026:* este parágrafo descrevia a checagem no presente, e a
flag em português. Nada disso existe: a interface `Tiered` com `Profiles()` não
está declarada em nenhum `.go` do `framework`, o `aru doctor` registra só
`-strict`, e nenhuma regra dele menciona `JOIN`. Os rótulos `conventional` e
`performance` existem — mas no campo `profiles` do `arandu.mod.toml`, que é
onde a seção 7 abaixo os trava. E a flag, quando entrar, é `--profile`: flag em
português viola a REGRA 6. Ver `docs/12-perfis.md`, que já lista as flags reais
de cada comando. O veredito da seção não depende da ferramenta e continua de pé:
o que torna a troca de perfil não-indolor é o modelo de dados, não a falta de um
lint.

Prometer troca indolor é exatamente o erro que está custando caro ao Next.js:
vender simples e entregar uma fronteira que exige entender sistema distribuído.

**A formulação honesta:** a API pública não muda — handler, service e policy
continuam idênticos. O modelo de dados muda, e **nem todo módulo é elegível**.

---

## 7. RECUSADO — renomear os perfis

A proposta era trocar "convencional" e "performance" por "Standard" e "High
Performance", com o argumento de que o Postgres não é o perfil lento e o Scylla é
o perfil distribuído. Foi aceita aqui e depois recusada, por dois motivos.

O rótulo não é custo zero: `conventional` e `performance` estão no campo
`profiles` do `arandu.mod.toml`, que já foi publicado. Manifesto de módulo é
contrato — renomear quebra todo módulo que declarou o valor antigo, e um período
de aceitar os dois nomes seria duas formas de dizer a mesma coisa (REGRA 9).

E o preço de ter tentado ficou visível: por um tempo o corpus carregou três
rótulos para dois perfis. **Valem `conventional` e `performance`**, no manifesto,
na CLI e no doc 12.

---

## 8. Já resolvido, e vai reaparecer: Socket.IO e µWebSockets.js

Ambos são Node. Adotar qualquer um transforma a aplicação em dois processos e
mata a frase que sustenta o produto. Essa parte não depende de nada ser
construído: é veto, e vale hoje.

O que muda é a alternativa. O doc 11 desenha o tiering dentro do Go —
`coder/websocket` como padrão, `gobwas/ws` para alta carga, mesma interface — e
esse desenho segue de pé, **sem uma linha escrita**.

*Corrigido em 07/08/2026:* esta seção dizia "o tiering existe dentro do Go". Não
existe. Não há pacote `ws` no projeto, e `coder/websocket` e `gobwas/ws` não
aparecem em nenhum `.go` nem em nenhum `go.mod` — nem sequer como dependência.
WebSocket é trabalho com gatilho escrito (doc 06: "quando SSE do HTMX não
bastar"), e o `aru make:event` não emite `broadcastOn` exatamente por isso.

A única coisa que existe hoje com a palavra WebSocket é uma regra do `doctor`
que **proíbe** `new WebSocket` dentro de `x-data` do Alpine — o oposto de uma
implementação. Registrado aqui porque o assunto reapareceu na consulta externa e
vai reaparecer de novo: o veto ao Node é resposta pronta, o tiering em Go ainda
é promessa.

---

## O que muda no plano

- `docs/27-eventos-e-outbox.md` — spec do outbox, com `Grant` no registro e
  medição de atraso no `Collector`.
- Tabela outbox e `Store()` sobem para a **Fase 2**.
- **REGRA 17** no `CLAUDE.md`: o `Grant` vale também na leitura. Não existe
  caminho de consulta sem policy, em nenhum desenho.
- Os rótulos dos perfis **não mudam**: `conventional` e `performance`.
- Value Objects entram como tipo, não como camada.
- Event Sourcing fica documentado com gatilho, não construído.

## Nota sobre o método

As três consultas convergiram em quase tudo: não há incumbente, arquitetura
opinativa é o diferencial, IA produz spec e não código, o templ é protagonista.
Convergência entre modelos diferentes é sinal fraco de acerto — todos leem as
mesmas fontes. O último item envelheceu: o templ caiu em agosto de 2026 e o motor
de view passou a ser o `kyse` (ADR 0020), o que é um exemplo do que a
convergência vale. O que vale mais são as **divergências**, e elas foram três:
Policy opcional na leitura, cinco camadas e Event Sourcing opt-in. As três
seriam decisões ruins, e as três teriam entrado sem ninguém notar.
