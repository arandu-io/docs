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
| Nosso desenho | 1 diretório, 7 arquivos |

E o custo maior não é o número de arquivos: **o dev Laravel não reconhece nada
ali.** Aggregate, Use Case, Value Object e Command Bus não existem no vocabulário
dele. A frase de onboarding — "se você conhece Laravel, você já sabe 80% do
Arandu" — morre no primeiro `ls`. REGRA 10 e doc 10 seguem valendo: vocabulário
Laravel, layout modular.

**Onde há razão:** Value Object para dinheiro, CPF/CNPJ e e-mail é bom e evita
classe inteira de bug. Cabe como tipo dentro de `<entidade>.entity.go`, sem
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

O texto afirma que se migra de Standard para High Performance "sem alterar a
estrutura da aplicação ou a forma de desenvolvimento, substituindo apenas os
componentes responsáveis pela persistência e comunicação".

Isso não é verdade, e o doc 12 já explica: Scylla é CQL, sem join, sem transação
cruzando partição, o tenant vira chave de partição, e o módulo precisa declarar
`Perfis()`. O `doctor --perfil=performance` reprova módulo que usa JOIN,
transação cruzando agregado ou entidade sem chave de partição.

Prometer troca indolor é exatamente o erro que está custando caro ao Next.js:
vender simples e entregar uma fronteira que exige entender sistema distribuído.

**A formulação honesta:** a API pública não muda — handler, service e policy
continuam idênticos. O modelo de dados muda, e **nem todo módulo é elegível**.

---

## 7. ACEITO — o rótulo dos perfis

"Standard" e "High Performance" comunica melhor que "convencional" e
"performance". O Postgres não é o perfil lento; o Scylla é o perfil distribuído.
Renomear no doc 12 e no `Perfis()`. Mudança de rótulo, custo zero, feita antes
da Fase 2.

---

## 8. Já resolvido, e vai reaparecer: Socket.IO e µWebSockets.js

Ambos são Node. Adotar qualquer um transforma a aplicação em dois processos e
mata a frase que sustenta o produto. O tiering existe dentro do Go:
`coder/websocket` como padrão, `gobwas/ws` para alta carga, mesma interface
(doc 11). Registrado aqui porque reapareceu na consulta externa e vai reaparecer
de novo.

---

## O que muda no plano

- `docs/27-eventos-e-outbox.md` — spec do outbox, com `Grant` no registro e
  medição de atraso no `Collector`.
- Tabela outbox e `Store()` sobem para a **Fase 2**.
- **REGRA 17** no `CLAUDE.md`: o `Grant` vale também na leitura. Não existe
  caminho de consulta sem policy, em nenhum desenho.
- Renomear os perfis para **Standard** e **High Performance**.
- Value Objects entram como tipo, não como camada.
- Event Sourcing fica documentado com gatilho, não construído.

## Nota sobre o método

As três consultas convergiram em quase tudo: não há incumbente, arquitetura
opinativa é o diferencial, IA produz spec e não código, Templ é protagonista.
Convergência entre modelos diferentes é sinal fraco de acerto — todos leem as
mesmas fontes. O que vale mais são as **divergências**, e elas foram três:
Policy opcional na leitura, cinco camadas e Event Sourcing opt-in. As três
seriam decisões ruins, e as três teriam entrado sem ninguém notar.
