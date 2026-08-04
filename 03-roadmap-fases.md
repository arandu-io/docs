# Roadmap — fases 1 a 4

Regra que atravessa todas as fases: **cada fase termina com algo em produção**,
não com algo pronto. Se a fase 1 acabar e nenhum sistema real estiver rodando
sobre ela, a fase 2 não começa.

Nome definido: **Arandu** (marca) e **`aru`** (CLI). Import path travado em
`github.com/arandu-io/framework` — a partir daqui é irreversível, porque módulos de
terceiros vão depender dele.

---

## Fase 1 — Núcleo e segurança
**Objetivo:** um binário que sobe, autentica, e não deixa você burlar a
autorização nem por acidente.

Entregas:

- `kernel` — boot, registro de módulos, pipeline, shutdown gracioso, health
- `config` — struct tipada, validação no boot, falha rápido
- `httpx` — router sobre `http.ServeMux`, grupos, middleware, metadados de rota
- `security` — **Grant/Policy** (a peça central), argon2id, CSRF, sessão
- `data` — contrato `Repository` exigindo Grant, `DB` instrumentado
- `validation` — `Validate() Errors`, sem tags nem reflection
- `observability` — slog com escopo de request, request_id, Collector
- `errorpage` — página de erro com stack, snippet, queries, dumps, diagnóstico
- módulo `auth` completo como referência canônica
- `aru key:generate`, `aru serve`, `aru migrate`, `aru routes` — a CLI é módulo
  próprio, e os três últimos delegam ao binário do projeto (ADR 0006)

Critério de saída — **todos** obrigatórios:

1. `go build ./...` e `go test ./...` limpos — **verificado**, com `-race`
2. Existe um teste que prova que chamar repositório sem Grant **não compila** —
   **verificado** em `data/grant_required_test.go`: dois fixtures em
   `data/testdata/`, um chamando repositório sem Grant e outro tentando forjar um
   Grant válido; o teste roda `go vet` e exige a falha com a mensagem esperada
3. Um dos seus produtos atuais roda o login em cima do `auth` — **aberto**,
   depende do proprietário
4. Um panic proposital em dev mostra queries e dumps do request na página —
   **verificado** em `httpx/middleware/pipeline_test.go`
5. Nenhuma dependência de terceiro no core além de `golang.org/x/crypto` —
   **verificado**, e checado em CI. `golang.org/x/sys` entra transitivamente,
   exigido pelo próprio `x/crypto`

Nota sobre o item 2: `//go:build ignore` foi trocado por `testdata/`. Com a tag,
`go vet` falha com "build constraints exclude all Go files" — que é um erro, mas
não o erro que prova a tese. Dentro de `testdata/` o go tool ignora o diretório no
build do módulo e ainda assim compila quando recebe o caminho explícito, o que
permite exigir a mensagem certa.

O que **não** entra na fase 1: templ, HTMX, geração de código, IA, deploy.

---

## Fase 2 — Produtividade e enforcement
**Objetivo:** parar de escrever CRUD à mão, e transformar as regras de
arquitetura em verificação automática.

Entregas — **todas feitas**:

- `aru make:module` — gerador determinístico (templates, zero LLM), com golden
  files e blocos `arandu:begin custom` preservados na regeração
- `aru doctor` — onze regras de análise estática. A décima primeira confere o
  `arandu.mod.toml` contra o AST: permissão usada e não declarada é erro
- `porang` — camada de view em templ + HTMX, com o `hx-headers` do CSRF no
  layout base (a armadilha nº 1 desse stack)
- `aru make:module --tenant` — escopo de tenant no SQL gerado, vindo do Grant
- `aru make:auth` — o starter kit, publicado no projeto e não importado
  (ADR 0012)
- `aru view:build` e `aru dev` — toolchain de view gerenciado pela CLI,
  `templ` e `tailwindcss` standalone, fixados em `arandu.toml`, cacheados em
  `~/.arandu/bin`. Zero Node
- `data.Transaction` e `events.Outbox` — evento gravado na mesma transação da
  escrita, nunca `Publish()` direto (doc 27)
- `arandu-io/kv` — cache, sessão, rate limit e lock sobre RESP (doc 11)
- `arandu-io/homebrew-tap` — `brew install arandu-io/tap/aru`
- migrações coletadas dos módulos. **Sem Atlas**: `data.Migrate` cobre o que a
  Fase 2 precisa, e um segundo migrador seria uma segunda forma de fazer a mesma
  coisa

Critério de saída:

1. `aru make:module cliente --campos "..."` gera módulo que compila, passa
   `doctor` e tem teste verde sem edição manual — **verificado**, e o checklist
   gera num projeto temporário e compila o resultado a cada execução
2. `doctor` roda em CI e falha em pelo menos 3 violações plantadas de propósito —
   **verificado**: oito violações plantadas em `aru/internal/doctor/testdata`,
   e o esqueleto passa no próprio `doctor --strict`
3. Um CRUD real de um dos seus produtos foi substituído por módulo gerado —
   **aberto**, depende do proprietário
4. Tempo medido: do zero ao CRUD funcionando em menos de 5 minutos —
   **verificado**, medido em 2 segundos

O que **não** entra na Fase 2: relay do outbox, fila morta e painel de eventos
(Fase 3, doc 27); registry de módulos (Fase 4, doc 18); documentação de uso
(fase própria, doc 08).

---

## Fase 3 — Observabilidade como produto
**Objetivo:** o cargo-chefe. É aqui que o framework deixa de ser "mais um" e
passa a ter uma razão para alguém trocar.

Entregas:

- console `/_arandu/debug` — ring buffer dos últimos 200 requests, com queries,
  dumps, eventos e chamadas externas (equivalente ao Telescope, mas core) —
  **feito**
- detector de N+1 com aviso no log, nomeando o statement e a contagem, e
  diagnóstico no topo da página do console — **feito**
- timeline do request: quanto tempo em SQL, em template, em HTTP externo, e
  quanto sobra — **feito**
- `aru trace <request_id>` — reconstrói um request pelo terminal, lendo o
  console da aplicação que está de pé — **feito**
- link "abrir na IDE" em VS Code, Cursor, GoLand e Zed — **feito**, uma
  implementação só, compartilhada pela página de erro e pelo console
- modo produção com header `X-Arandu-Trace` autenticado — **feito**
- relay do outbox, fila morta e atraso no health check (doc 27) — **feito**
- `kernel.Diagnostic`: o módulo diz o que está errado, e aparece na página de
  erro — **feito**
- agendamento: `kernel.Task`, scheduler in-process com lock, escopo por tenant,
  e a tarefa instrumentada como request no console (doc 16) — **feito**
- fila: `arandu-io/queue` sobre a tabela e sobre RESP, com backoff, fila morta e
  `aru work` — **feito**, ADR 0016
- storage: `arandu-io/storage` em disco e sobre o protocolo S3, com R2 como
  padrão — **feito**, ADR 0016
- deploy: `aru build`, `aru build --docker` e o Dockerfile distroless do
  esqueleto (doc 17) — **feito**
- exportador OpenTelemetry opcional — **adiado com gatilho**, ver ADR 0015

Critério de saída:

1. Um bug real de um dos seus produtos foi diagnosticado exclusivamente pelo
   console, sem `fmt.Println` — **aberto**, depende do proprietário
2. O overhead do Collector em produção **desligado** é mensurável e zero —
   **verificado**, e como teste em vez de benchmark: código instrumentado e
   código sem instrumentação alocam exatamente igual, e 22 chamadas de `Record`
   com o Collector ausente custam 0 alocações e 0 bytes
3. Um dev externo consegue instalar, quebrar de propósito e entender o erro
   sem ler documentação — **aberto**, depende de pessoa

---

## Fase 4 — Especificação e geração assistida
**Objetivo:** a camada que você desenhou desde o começo, agora com alvo rígido
para compilar contra.

Entregas:

- DSL declarativa em YAML (`module.arandu.yaml`) descrevendo módulo, campos e
  permissões — **feito**. Sem relações e sem telas: ver "o que ficou de fora"
- `aru generate` — DSL → código, determinístico, **a mesma engine** do
  `make:module`, com os mesmos golden files — **feito**
- `aru schema` — o JSON Schema da DSL, gerado das mesmas constantes que o
  validador usa, com exemplo que passa na própria validação — **feito**
- validação da spec **antes** de gerar, reportando todos os problemas de uma
  vez, e sem escrever arquivo nenhum quando reprova — **feito**
- round-trip: código gerado → DSL → código idêntico byte a byte, verificado por
  teste — **feito**
- `arandu ai` — **não construído**. O que ele faria é pedir a DSL a um modelo, e
  isso é o que qualquer assistente já faz melhor com o schema na mão. O produto
  desta camada é o schema, não um wrapper de chamada

Critério de saída:

1. A partir de uma descrição, sai módulo que compila, passa `doctor` e tem teste
   verde, **sem edição manual** — **verificado**, e o checklist gera num projeto
   temporário e roda os três a cada execução
2. Um LLM sem fine-tuning, recebendo só o schema, acerta a spec em 9 de 10
   tentativas — **aberto**, depende de medição
3. Regenerar não perde o código customizado — **verificado**

### O que ficou de fora da DSL, e por quê

`relações` e `telas` estavam na lista e não entraram. O motivo é a REGRA 15: a
DSL é um conjunto fechado, e cada coisa que entra nela é mantida para sempre.

- **relações** — `belongs_to` e `has_many` só fazem sentido dentro do agregado,
  e o `data.Repository` já é por agregado. Uma relação na DSL viraria join
  gerado, e join gerado é o começo de um ORM
- **telas** — `list`, `form` e `show` gerariam templ, e a forma canônica de uma
  tela ainda não existe: `make:auth` publica a primeira, e é dela que se aprende
  o que uma tela gerada deve parecer. Gerar antes disso é adivinhar

Ambas voltam à mesa quando houver um módulo real pedindo, que é o mesmo gatilho
de todo o resto.

---

---

## Fase D — Documentação
**Entra depois da Fase 3.** Não em paralelo, não durante.

Política completa em `08-documentacao.md`. Em resumo: nada de guia, tutorial ou
site enquanto a API se move. Comentário de doc em Go e ADR seguem obrigatórios
durante todo o desenvolvimento — não são documentação, são código e decisão.

Entregas:

- API pública do core congelada; nenhuma mudança de assinatura durante a fase
- guias escritos a partir de `docs/01` a `08`
- `arandu-io/docs` — markdown puro, um branch por versão
- `arandu-io/website` — **aplicação Arandu**, não Next.js nem gerador estático.
  Depende do `porang`, e é por isso que esta fase vem depois da 3
- toda mensagem de erro do framework revisada contra o tom de voz de `04-marca.md`
- página de origem do nome revisada por quem estude tupi antigo

Critério de saída, único e verificável:

> Alguém de fora instala, cria um CRUD e sobe em produção seguindo só o site, sem
> perguntar nada a você.

---

## Depois da fase 4

Não planejar agora. As três coisas que provavelmente aparecem — marketplace de
módulos, módulos verticais brasileiros (NF-e, SPED, PIX, boleto), e deploy
integrado — dependem de coisas que só a fase 3 vai revelar: se tem gente usando,
quem está usando, e o que essa gente pede.

Se ninguém estiver usando ao fim da fase 3, o produto certo é a plataforma
interna, e as fases 1 a 3 já pagaram por si.
