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

Entregas:

- `aru make:module` — gerador determinístico (templates, zero LLM)
- `aru doctor` — as 10 regras de análise estática listadas em `cmd/arandu/doctor.go`
- camada de view: templ + HTMX, com o template base já trazendo o
  `hx-headers` do CSRF (a armadilha nº 1 desse stack)
- `aru make:module --tenant` — escopo de tenant no SQL gerado, vindo do Grant
- migrações coletadas dos módulos, com Atlas
- golden files para toda saída do gerador

Critério de saída:

1. `aru make:module cliente --campos "..."` gera módulo que compila, passa
   `doctor` e tem teste verde sem edição manual
2. `doctor` roda em CI e falha em pelo menos 3 violações plantadas de propósito
3. Um CRUD real de um dos seus produtos foi substituído por módulo gerado
4. Tempo medido: do zero ao CRUD funcionando em menos de 5 minutos

---

## Fase 3 — Observabilidade como produto
**Objetivo:** o cargo-chefe. É aqui que o framework deixa de ser "mais um" e
passa a ter uma razão para alguém trocar.

Entregas:

- console `/_arandu/debug` — ring buffer dos últimos N requests, com queries,
  dumps, eventos, jobs e chamadas externas (equivalente ao Telescope, mas core)
- detector de N+1 com aviso no log e badge no console
- timeline do request: quanto tempo em SQL, em template, em HTTP externo
- `aru trace <request_id>` — reconstrói um request pelo terminal
- link "abrir na IDE" funcionando em VS Code, Cursor, GoLand e Zed
- exportador OpenTelemetry opcional, para quem já tem Grafana/Datadog
- modo produção com header `X-Arandu-Trace` autenticado, para tracing sob demanda

Critério de saída:

1. Um bug real de um dos seus produtos foi diagnosticado exclusivamente pelo
   console, sem `fmt.Println`
2. O overhead do Collector em produção **desligado** é mensurável e zero
   (benchmark comparando com e sem, no mesmo binário)
3. Um dev externo consegue instalar, quebrar de propósito e entender o erro
   sem ler documentação

---

## Fase 4 — Especificação e geração assistida
**Objetivo:** a camada que você desenhou desde o começo, agora com alvo rígido
para compilar contra.

Entregas:

- DSL declarativa (`.arandu.yaml` ou HCL) descrevendo módulo, entidades, campos,
  relações, permissões e telas
- `aru generate` — DSL → código, determinístico, mesma engine da fase 2
- schema JSON da DSL publicado — é isso que faz o LLM gerar spec válida sem
  ter sido treinado no framework
- validação da spec **antes** de gerar: erro de spec nunca vira código quebrado
- `arandu ai` — o LLM produz **apenas a DSL**, nunca Go; a spec é revisada pelo
  humano antes de gerar
- round-trip: código gerado → DSL → código idêntico (prova de que não há drift)

Critério de saída:

1. A partir de uma descrição em português, sai módulo que compila, passa
   `doctor` e tem teste verde, **sem edição manual**
2. Um LLM sem qualquer fine-tuning, recebendo só o schema da DSL, acerta a spec
   em 9 de 10 tentativas
3. Regenerar um módulo existente não perde código customizado (blocos
   protegidos por marcador)

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
