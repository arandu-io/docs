# 22 — Dores dos desenvolvedores nos frameworks

Levantamento do que os desenvolvedores reclamam nos principais frameworks, e o
que cada dor significa para o Arandu. Cada item termina em **validado** (decisão
nossa já responde), **exposto** (podemos cair no mesmo buraco) ou **lacuna**
(nenhum documento cobre). Pesquisa de julho/2026; fontes ao final.

---

## 1. AdonisJS — o espelho

Comece por aqui. O AdonisJS é "o Laravel do Node": opinativo, batteries
included, ORM próprio, CLI própria, TypeScript-first. É o análogo mais próximo
do que o Arandu quer ser em Go — e a lista de reclamações dele é, literalmente,
a nossa lista de riscos.

| Dor | Detalhe |
|---|---|
| **Arquitetura híbrida** | Quando o requisito não cabe na camada padronizada, o time faz gambiarra. Parte do código segue a convenção, parte evolui fora dela. Resultado: perda de homogeneidade e aumento do custo total de propriedade no longo prazo |
| **Lock-in por brick** | ORM, migrações, middleware e CLI são próprios. Trocar qualquer peça por componente externo exige reescrita significativa — e o próximo update do framework arrisca quebrar a adaptação |
| **Ecossistema pequeno** | Menos plugins que Express ou NestJS. Integração específica (monitoramento, biometria, evento complexo) exige conector caseiro |
| **Comunidade menor** | Menos gente para responder quando algo trava |
| **Mudanças rápidas** | Framework relativamente novo, updates frequentes, esforço extra só para acompanhar |
| **Curva de aprendizado** | Conjunto rico de features e convenções próprias |
| **Overhead de desempenho** | Comparado a microframework, o pacote completo custa |
| Documentação | Aparece como ponto forte e, em outra fonte, como "pobre" — sinal de cobertura desigual |

**Leitura para nós.** Cinco dessas sete são estruturais de "framework opinativo
novo e pequeno", não defeitos de execução do Adonis. Vamos ter todas. A única
que dá para atacar por desenho é a primeira — e é a mais cara.

**E o posicionamento deles é bom.** A página oficial de resposta à comunidade
diz que os problemas do ecossistema JS não são de ferramenta, são de **coesão**:
existem ORMs excelentes, routers rápidos, validadores tipados e libs de auth
ótimas, mas nunca foram desenhados para trabalhar juntos — "móvel de IKEA de
cinco coleções diferentes: cada peça é bem feita, nada encaixa". O argumento
deles: quando um único time desenha router, validador, ORM, mailer e cliente de
teste, os tipos fluem por tudo.

Esse é o melhor argumento de venda de framework integrado que apareceu na
pesquisa, e vale mais para nós do que para eles: em Go, tipo que flui de ponta a
ponta é verificado pelo compilador.

---

## 2. HTMX — o risco que herdamos inteiro

Apostamos nele. Então importa mais saber onde ele falha do que onde ele brilha.

### O que dizem contra

| Dor | Detalhe |
|---|---|
| **Sobrefragmentação e "pop-in"** | A falha mais comum: a página faz quinze chamadas separadas e o usuário vê os pedaços aparecendo um a um. O site parece quebrado ou lento |
| **Spaghetti de endpoints** | Vários endpoints atualizando pedaços diferentes do DOM de forma independente geram estado imprevisível e inconsistência visual |
| **Estado no cliente é o erro nº 1** | "Se você está escrevendo JavaScript para controlar qual aba está ativa, quais filtros estão aplicados ou quais itens estão selecionados, você está brigando com a arquitetura do htmx" |
| **Offline não funciona** | Depende de hipermídia vinda do servidor |
| **Interatividade de alta frequência** | Ferramenta de desenho, jogo rápido, canvas: latência de rede inviabiliza |
| **Colaboração em tempo real** | Times escolheram HTMX para dashboard colaborativo e passaram meses costurando WebSocket e estado no cliente que o React resolveria nativo |
| Latência é o teto | Toda interação é ida e volta ao servidor |

### O que joga a favor

O htmx tem 14 KB minificado, 5,5 KB gzipado; com extensões, uns 20 KB na rede.
O app React mediano em 2026 entrega cerca de 487 KB de JavaScript comprimido
**antes** de qualquer código seu rodar, e o time-to-interactive num Android
mediano em 4G estrangulado fica entre 4 e 7 segundos no p75.

Sem virtual DOM, sem hidratação, sem passo de reconciliação.

### Consequências diretas para o Arandu

- **Sobrefragmentação vira problema do gerador, não do usuário.** Se
  `aru make:module` emitir tela com quinze fragmentos, todo app gerado nasce com
  o defeito. Precisa de diretriz escrita e, de preferência, verificação no
  `doctor`.
- **O limite do Alpine que escrevemos no doc 14 está certo.** "Estado do
  cliente, efêmero e invisível ao servidor" é exatamente o que a comunidade
  aponta como o erro nº 1 de quem usa HTMX mal. Validado por fonte externa.
- **Precisamos de um "quando não usar o Arandu" escrito.** Offline, alta
  frequência e colaboração em tempo real não são nossos. Dizer isso primeiro é
  mais forte do que deixar o cliente descobrir depois — e alimenta a lacuna da
  história de saída.

---

## 3. Livewire — a lição mais cara da pesquisa

Relevante porque é o concorrente conceitual direto: UI reativa dirigida pelo
servidor.

O time do Lychee passou **um ano e meio** reimplementando o front-end em
Livewire. Em desenvolvimento local tudo parecia rápido, as requisições
instantâneas, a reatividade palpável. Ao subir para servidor externo, bateu na
realidade da latência. O profiling mostrou que o custo maior estava na
serialização: transformar linha de SQL em objeto de modelo, e o objeto Carbon
sendo inicializado toda vez, gastando tempo só para descobrir o fuso do
servidor. A conclusão deles foi que Livewire e renderização no servidor não eram
boa opção para o produto — depois de seguirem o hype.

Em produção com Filament, o relato é o mesmo: interface dolorosamente lenta, e a
lista de mitigação é `wire:model.defer`, lazy loading, eager loading para evitar
N+1, minificação e CDN. Ou seja: o framework devolve o problema para o dev.

**A lição, e ela é desconfortável:** arquitetura dirigida pelo servidor morre de
latência e de custo de serialização por request, e **isso é invisível em
desenvolvimento local**.

Nossa posição é melhor por construção — HTMX manda fragmento de HTML, não faz
round-trip de estado serializado e assinado como o Livewire, e não há boot por
request. Mas o modo de morrer é o mesmo, e nosso `Collector` hoje conta query e
não mede **orçamento de latência**. Isso é item novo.

---

## 4. Laravel

| Dor | Detalhe |
|---|---|
| **"Magia demais"** | A reclamação nº 1. Coisas acontecem e não dá para rastrear. A defesa oficial — "é só metaprogramação de PHP bem documentada" — confirma a dor |
| Facades como estado global | Debate permanente |
| **Dor de upgrade** | Relato concreto: duas semanas depurando a subida de Laravel 10 para 11 porque um pacote dependia de API interna removida. "A magia que deixa o desenvolvimento rápido deixa a depuração lenta" |
| Teto de desempenho | Chega uma hora em que o overhead aparece |
| Encoraja arquitetura bagunçada | Não impõe camada |
| Complexo demais para app pequeno | Principal ressalva em reviews agregadas |
| Eloquent e N+1 | Ver bloco 7 |

## 5. Symfony

| Dor | Detalhe |
|---|---|
| **Curva íngreme** | Assusta até quem vem do Laravel |
| **Verbosidade** | Mais código para o mesmo resultado |
| **Fadiga de decisão** | "Formas demais de fazer a mesma coisa" |
| **Fragmentação da comunidade** | Múltiplas soluções para o mesmo problema comum |
| Desenvolvimento mais lento | Exige pensar mais antes |
| Setup inicial pesado | Muita configuração e dependência |
| Excesso de documentação | Volume atrapalha achar a resposta |

Repare que "fadiga de decisão" e "fragmentação" são o custo exato de **não** ter
a REGRA 9. O Symfony é o contra-exemplo que justifica a regra.

## 6. Next.js — o caso mais instrutivo

| Dor | Detalhe |
|---|---|
| **Mudança constante** | "Next.js me perdeu com as mudanças constantes e a complexidade. O benefício não justifica a sobrecarga cognitiva" |
| **Duas formas de tudo** | Pages Router e App Router coexistindo; Server Components mudando o modelo mental |
| Marketing mentiu | Vendeu-se como "só React com extras" quando exige entender fronteira cliente/servidor. A crítica não é ter ficado complexo — é ter prometido simples |
| CVE recorrente | Cadeia de vulnerabilidades no App Router ao longo de 2025–2026, com upgrade urgente repetido |
| Lock-in de fornecedor | Acoplamento percebido com a Vercel |
| Build lento | Time relatou corte de 60% ao migrar para Astro |
| "Kubernetes do frontend" | A frase que circula |

Para onde vão: TanStack Start, Astro, SvelteKit. E o elogio é sempre o mesmo —
explícito, previsível, tipado, sem esconder o que acontece.

## 7. Django e o problema do ORM

| Dor | Detalhe |
|---|---|
| **N+1 é a dor dominante** | "O ORM torna perigosamente fácil escrever query que funciona com 100 linhas em dev e colapsa com 100.000 em produção" |
| Lazy loading implícito | Acessar `book.author` dispara query e ninguém vê |
| O ORM não otimiza sozinho | Auditoria manual, sempre |
| Só aparece em produção | Passa em review, passa em staging, quebra com tráfego real |

Vale para Django, Rails, Eloquent e Hibernate. É a dor mais transversal da
pesquisa inteira.

## 8. Rails

| Dor | Detalhe |
|---|---|
| Magia que não dá para grepar | Método gerado dinamicamente não aparece em busca no código |
| Convenção mágica exige decorar | "Você tem que saber muita coisa" só para nomear certo |
| **Mass assignment** | Caso histórico: dev não declara `attr_accessible`, atacante injeta campo no POST e assume a conta. A mitigação depende da disciplina do dev |
| Consumo de memória | Reclamação persistente |

## 9. Spring Boot e NestJS

| Dor | Detalhe |
|---|---|
| **Fadiga de boilerplate** | Mesmo com auto-configuração, a verbosidade atrasa entrega |
| DI onde não precisava | "500 arquivos de boilerplate para um hello world com dependência injetada e ORM mapeado" |
| **Memória** | 256–512 MB de heap como baseline, antes de servir um request |
| **Startup** | 8 a 15 segundos — inviabiliza serverless, atrasa o ciclo de dev |
| Magia por anotação | Difícil depurar sem conhecer o mecanismo |
| NestJS: cerimônia herdada do Angular | Decorators, módulos, providers, container |

## 10. Express e frameworks Go

**Express:** flexibilidade excelente para protótipo, "catastrófica para time
grande construindo sistema multi-tenant complexo" — falta de estrutura.

**Go:** aqui a dor não é com um framework, é com a ideia de framework.

| Dor | Detalhe |
|---|---|
| **Lock-in** | "Difícil de remover depois de integrado fundo" |
| **Abandono** | O Gorilla parou de ser mantido pelo time original em nov/2023, e 17% dos devs Go ainda usavam em 2025 |
| **Fiber quebra o ecossistema** | Usa fasthttp, então middleware genérico de `net/http` não funciona — apontado como o maior problema real do Fiber |
| Migração custa meses | Lock-in, deprecação e migração cara: os três riscos nomeados |
| `net/http` do 1.22 esvaziou o argumento | Roteamento nativo melhorou; puxar dependência pesada virou questão |

---

## Os padrões que atravessam tudo

Por frequência, somando todos os ecossistemas:

1. **Magia que não se rastreia** — Laravel, Rails, Spring, Next.js.
2. **N+1 e o ORM que facilita escrever código lento** — Django, Rails, Laravel,
   Hibernate. A mais cara em produção e a mais invisível em dev.
3. **Segunda forma de fazer a mesma coisa** — Next.js no extremo; Symfony chama
   de fadiga de decisão.
4. **Esteira de upgrade** — mudança constante, migração cara, CVE urgente.
5. **Boilerplate e cerimônia** — Spring, NestJS, Symfony.
6. **Lock-in e risco de abandono** — Go e AdonisJS.
7. **Peso: memória e startup** — Spring no extremo.
8. **Ecossistema e comunidade pequenos** — AdonisJS. Custo de ser novo.
9. **Marketing que promete simplicidade e entrega complexidade** — a crítica
   mais dura ao Next.js.

---

## Auditoria: onde o Arandu está

### Validado

| Dor | Nossa resposta |
|---|---|
| Magia que não se rastreia | Sem container, sem facades, sem reflection; wiring explícito e commitado (ADR 0001, 0002) |
| Mass assignment | Request struct explícita; a classe de bug não existe (doc 13) |
| N+1 | sqlc elimina lazy loading implícito e o `Collector` detecta e nomeia repetição, com arquivo:linha (doc 02) |
| Duas formas de tudo / fadiga de decisão | REGRA 9, escrita antes desta pesquisa |
| Memória e startup | Binário único, sem boot por request |
| Fasthttp quebrando ecossistema | Recusado; ficamos em `net/http` |
| Estado do cliente em HTMX | O limite do Alpine no doc 14 é exatamente o erro nº 1 apontado pela comunidade |
| Regra que ninguém verifica | `doctor` transforma regra em checagem de CI |
| Marketing que promete demais | Tom sóbrio (doc 04) |

### Exposto

| Dor | Nosso risco |
|---|---|
| **Arquitetura híbrida** | A dor central do AdonisJS. Quando o requisito não couber no padrão, metade do código do cliente sai da convenção. Nosso antídoto é o bloco `// arandu:begin custom`, que hoje está enterrado no doc 19 como detalhe da Fase 4 |
| **Boilerplate e cerimônia** | Policy por entidade, Grant em toda assinatura, wiring explícito, `arandu.mod.toml`. Sem gerador cobrindo tudo, viramos o Spring |
| **Sobrefragmentação de HTMX** | Se o gerador emitir tela com muitos fragmentos, todo app do cliente nasce com "pop-in" |
| **Latência invisível em dev** | O modo de morrer do Livewire. O Collector conta query, não mede orçamento de latência |
| **Ecossistema e comunidade pequenos** | Estrutural. Não tem solução de desenho, só de tempo |
| **Lock-in** | Somos deliberadamente opinativos, e os três projetos vivos do nicho se posicionam contra isso |

### Lacuna — nenhum dos 21 documentos cobre

**1. Política de estabilidade e versionamento.** A dor nº 1 do Next.js é mudança
constante, e o relato do Laravel 10→11 mostra o custo concreto: duas semanas de
depuração por API interna removida. Não temos nada escrito sobre compatibilidade,
cadência, janela de suporte ou o que autoriza breaking change. Decidir **antes**
da v1.

**2. Política de divulgação de vulnerabilidade.** Sem `SECURITY.md`, sem canal,
sem embargo, sem processo de CVE. Para um framework cuja tese é "o caminho
errado é impossível", improvisar isso sob pressão é contradição performática.

**3. História de saída e não-objetivos.** "Como sair do Arandu" e "quando não
usar o Arandu" não existem. O risco de abandono é máximo num projeto de uma
pessoa — o Gorilla parou e 17% do ecossistema ficou pendurado. E os
não-objetivos do HTMX (offline, alta frequência, colaboração em tempo real) são
matéria-prima pronta para esse documento.

---

## O que isso muda no plano

- As três lacunas viram `23-estabilidade.md`, `24-seguranca-processo.md` e
  `25-saida-e-nao-objetivos.md`. As duas primeiras são pré-requisito de
  repositório público, junto com a licença (doc 20).
- **Bloco de escape sobe para a Fase 2.** Deixou de ser detalhe da DSL: é o
  antídoto contra a dor central do AdonisJS, e módulo gerado sem marcador de
  bloco protegido não dá para retrofitar.
- **O `Collector` ganha orçamento de latência**, não só contagem de query. O
  Livewire mostra que o app morre no deploy, não no localhost.
- **Diretriz de fragmentação** no gerador de views, com verificação no `doctor`.
- **Coesão vira o pitch.** O argumento do AdonisJS — peças excelentes que nunca
  foram desenhadas para encaixar — é mais forte em Go, onde o encaixe é
  verificado pelo compilador.
- **O detector de N+1 merece a primeira tela do site.** É a dor mais citada em
  quatro ecossistemas e a demonstração leva dez segundos.

---

## Fontes

- Edana, *Advantages and Drawbacks of AdonisJS* — jan/2026
- adonisjs.com/we-hear-you (posicionamento oficial sobre coesão)
- SaaSHub, fichas de prós e contras do AdonisJS e do Symfony
- Vibe Forem, *HTMX in 2026: Why Hypermedia is Dominating the Modern Web*
- DEV, *htmx in 2026: When You Don't Need React (And When You Absolutely Do)*
- Medium/Hex Shift, *Why Most Developers Are Using HTMX Badly*
- BirJob, *HTMX in Production 2026* (números de bundle)
- Brent Haskins, *React vs HTMX in 2026*
- LycheeOrg, *Livewire performance problems* (o caso de 1,5 ano)
- Laracasts, *Painfully Slow Interface with Filament/Livewire in Production*
- fsck.sh, *Laravel vs Symfony 2025* (upgrade 10→11)
- Laracasts, *The Many Misconceptions of Laravel*
- LogRocket, *Is Next.js still developer-friendly?*
- Appwrite e buildmvpfast sobre a fuga do Next.js
- AppSignal e techbuddies sobre N+1 no Django
- coreitpro, *Does Ruby on Rails Have Too Much Magic?*
- Out Plane, *Spring Boot vs NestJS* — fev/2026
- JetBrains, *Popular Golang Web Frameworks* — abr/2026
- Encore, *Gin vs Echo vs Fiber 2026*
- Three Dots Labs, *When you shouldn't use frameworks in Go*
