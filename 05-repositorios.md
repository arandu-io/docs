# Organização dos repositórios

Espelhado no modelo do Laravel, adaptado para as regras de módulo do Go.
Org: **`github.com/arandu-io`**.

## O padrão que estamos copiando

O Laravel separa três coisas que costumam ser confundidas:

| Papel | Laravel | Arandu |
|---|---|---|
| Código do framework | `laravel/framework` | **`arandu-io/framework`** |
| Esqueleto de projeto | `laravel/laravel` | **`arandu-io/arandu`** |
| Documentação | `laravel/docs` (markdown, branch por versão) | **`arandu-io/docs`** |

Na letra: o repositório homônimo da org é o **esqueleto**, não a biblioteca —
exatamente como `laravel/laravel`. Quem cria projeto usa o nome da marca; quem
importa código usa `framework`.

Consequência em Go, ao contrário do PHP: o import path é o caminho do repo, então
o usuário escreve `github.com/arandu-io/framework/kernel`. No Composer o nome do
pacote é desacoplado do namespace (`Illuminate\...`), aqui não. A marca vive na
org. Decisão travada — mudar depois quebra todo import publicado.

O esqueleto não é biblioteca: ninguém depende dele: você clona uma vez e ele
nunca mais é atualizado inteiro. Isso é o que permite ao framework evoluir sem
brigar com a estrutura de pastas de projetos antigos.

## Mapa completo

São **catorze repositórios**, e são estes: `framework`, `arandu`, `aru`, `ui`,
`kyse`, `hesape`, `mcp`, `database`, `kv`, `queue`, `storage`, `examples`,
`homebrew-tap` e `docs`. O `website` vem depois, na Fase D.

Dez chegaram até a Fase 3. Os quatro últimos foram autorizados **depois** dela,
cada um com ADR: `ui` (ADR 0026), `kyse` (ADR 0027), `mcp` (ADR 0030) e `hesape`
(ADR 0049). Antes disso a conta tinha caído de onze para dez, quando o `porang`
foi dissolvido (ADR 0021).

Treze são públicos; o `docs` é privado. `gh repo list arandu-io` devolve quinze,
contando o `plans`, que não é produto e não conta (REGRA 3).

*Corrigido em 14/08/2026:* este parágrafo dizia dez e listava dez, o que já
estava errado desde 09/08/2026 — o `mcp` existia, com `go.mod` próprio e ADR
própria, e nunca entrou na conta. `docs/31-reorganizacao-hesape.md`, seção 7.2,
item 7, tinha apontado a divergência sem que ela fosse fechada.

*Corrigido em 07/08/2026:* esta frase terminava com "o `.github` guarda
workflow, não código", e a lista do Núcleo trazia `arandu-io/.github` entre os
repositórios. Ele não existe — `gh api repos/arandu-io/.github` responde 404.
Os nove repositórios de código carregam o próprio `.github/workflows`; o `docs`,
que é só markdown, não tem `.github/` nenhum. Criar o repositório de
perfil da org é decisão em aberto — e, pela letra da REGRA 3, é criar o décimo
primeiro repositório da org, o que exige pedido ao proprietário. A regra lista os
dez nominalmente e inclui o `docs`, que é só markdown, então ela não é sobre
"repositório de código". É o mesmo raciocínio que este documento aplica adiante
ao `mail` e ao `ws`.

### Núcleo

Publicados:

```
arandu-io/framework     a biblioteca. Módulo: github.com/arandu-io/framework
                        Deps: stdlib + golang.org/x/crypto. Nada além.
arandu-io/arandu        esqueleto de projeto. O que `aru new` clona.
                        Ninguem depende dele como biblioteca.
arandu-io/aru           CLI. Módulo próprio, binário `aru`.
arandu-io/docs          markdown do guia. Branch por versão (1.x, 2.x).
arandu-io/homebrew-tap  fórmula do `aru` (nome da fórmula: aru, não arandu).
```

Decidido e não criado:

```
arandu-io/website       site que renderiza o docs. Fase D.
                        E uma APLICACAO ARANDU — dogfood. Ver 08-documentacao.md
```

### Adapters (drivers de infraestrutura)

Existem para manter o core sem dependência de terceiros (ADR 0004).

**Esta seção foi substituída por `11-adapters.md`**, que resolveu um problema que
não estava visto aqui: em Go não existe dependência opcional, então um pacote por
fornecedor multiplica repositório sem isolar dependência nenhuma. A divisão é por
**domínio**, com driver pesado em submódulo.

No ar, quatro repositórios, com os submódulos que de fato existem:

```
arandu-io/database    registro, pool e erro; submódulos pgx, sqlite, mysql
arandu-io/kv          go-redis (RESP): cobre Redis, Dragonfly, Valkey, KeyDB
arandu-io/queue       fila sobre a tabela; submódulo kv para volume
arandu-io/storage     interface + disco local; submódulo s3 (R2 como padrão)
```

Desenho, sem repositório e sem data. O `11-adapters.md` os põe em "com demanda":

```
arandu-io/mail        interface + SMTP; submódulos ses, resend, postmark, mailgun
arandu-io/ws          coder/websocket e gobwas/ws
```

Criar qualquer um dos dois é passar de catorze para quinze repositórios, e a
REGRA 3 exige autorização do proprietário antes — como a ADR 0016 registrou para
`queue` e `storage`.

*Corrigido em 07/08/2026:* este bloco listava os seis num só lugar, como se os
seis existissem, e dava `gcs` e `azure` como submódulos do `storage`. A árvore
publicada tem só `s3` (`gh api repos/arandu-io/storage/git/trees/HEAD`), e
`mail` e `ws` não aparecem em `gh repo list arandu-io`. A divisão por domínio
continua sendo a decisão; o que mudou aqui é separar o publicado da planta.

Seis domínios no lugar dos quinze pacotes que a lista por fornecedor produziria.
É o mesmo que o Laravel faz — `Illuminate\Mail` traz SMTP, Mailgun, Postmark,
SES e Resend juntos; `Illuminate\Database` traz MySQL, Postgres e SQLite juntos.

Convenção de nome: **sem prefixo redundante e em inglês funcional**. É
`arandu-io/kv`, não `arandu-io/arandu-redis` — a org já carrega o nome, e
o nome tupi fica nos produtos (`taba`, `moitara`), não na infraestrutura.

### Módulos de primeira parte

O equivalente aos pacotes first-party do Laravel. Cada um implementaria
`kernel.Module` e nada mais. **Nenhum existe hoje**: a triagem de `06-escopo.md`
mandou a maior parte para o core, para biblioteca de terceiro ou para "depois,
com gatilho".

Nomes definidos em `docs/07-pacotes.md`. Os que continuam no mapa:

```
arandu-io/moitara    registry de módulos — ver 18-marketplace.md
arandu-io/uirapuru   servidor websocket (Reverb) — gatilho: SSE não bastar
arandu-io/potiron    dashboard de filas (Horizon) — depois de fila em produção
arandu-io/taba       painel administrativo — pago (Nova), Fase 4
```

O que saiu da lista, e para onde foi:

- `porang` foi dissolvido (ADR 0021): as views ficam em
  `arandu/resources/views/` e o runtime de render em `framework/view/`
- o starter kit **é** repositório: `arandu-io/ui` (ADR 0026, que supersede a 0012)
- o inspetor de requests (`marandu`) é core, no console `/_arandu/debug`,
  desde a Fase 3
- tokens de API estão decididos como core, e **não existem**. Ver a nota abaixo
- `maraca` e `boitata` viraram exportação de OpenTelemetry (ADR 0015)
- Pint e Sail não têm equivalente: `gofmt` e `go run` resolvem

*Corrigido em 07/08/2026:* esta lista dizia "tokens de API são core, em
`security/token.go`". Não há esse arquivo: `ls framework/security/` devolve
`csrf.go`, `password.go`, `policy.go` e `session.go`, mais os testes de cada um,
e nada além. O que vale é o veredito de `docs/packages/sanctum.md` — token de
API é core, sem repositório próprio, e seria uma variação de sessão sobre o
`SessionStore` e o `Grant`, que existem. **Construir ou rebaixar para "depois,
com gatilho" é decisão do proprietário, e está aberta.** Enquanto ela não for
tomada, nenhum documento pode citar `security/token.go` como coisa entregue.

Nota: as primitivas de autenticação (Grant, Policy, sessão, argon2id) ficam no
**core**, não aqui — e essas quatro existem: `security.Grant` em `policy.go`,
argon2id em `password.go`. Segurança é a tese do produto; não pode ser um pacote
opcional que o usuário esquece de instalar. O `arandu-io/ui` publica só as telas.

### Distribuição

```
arandu-io/examples   a aplicacao que demonstra o framework, e a forma canonica
                     que o gerador da Fase 2 vai emitir
```

A camada de view não aparece aqui porque não é repositório, e ela vive em dois
lugares: o desenvolvedor edita `resources/views/` no próprio projeto; o
**runtime de render** do `kyse` e os assets embutidos são core, em
`framework/view/` (`render.go`, `runtime.go` e o `go:embed` de `assets.go`); e o
**compilador** do `kyse` é parte da `aru`, em `aru/internal/kyse/`, chamado pelo
`aru view:build`.

## O que é diferente do Laravel, e importa

**1. Sem subtree split.** O Laravel desenvolve em monorepo e publica pacotes
somente-leitura por automação. Em Go, um módulo por repositório, cada um com
sua própria tag semver. Mais simples, e não precisa de robô de release.

**2. Não fatiar o core agora.** O Laravel fatiou `illuminate/*` depois de anos,
quando gente de fora começou a querer só as Collections. Fatiar antes de alguém
pedir é complexidade especulativa.

**3. A regra de versão major do Go é brutal — planeje para não usar.**
Em Go, major version vai no import path: `v2` significa
`github.com/arandu-io/framework/v2`, e **todo** módulo do ecossistema precisa
atualizar o import. O Laravel não tem esse problema. Consequência prática:
estabilizar a API pública o mais rápido possível e ficar em `v1` por anos. Cada
`v2` é uma migração forçada em todo projeto de todo cliente.

**4. Matriz de compatibilidade.** Cada adapter e módulo declara no seu `go.mod`
a versão de core que suporta. O MVS do Go escolhe a maior versão exigida, o que
normalmente funciona — mas dois módulos que exijam minors incompatíveis
fragmentam o ecossistema em silêncio. O `aru doctor` deve conferir isso.

**5. `aru` em repositório separado do core.** Se a CLI viver dentro do módulo do
framework, todo projeto que importa o framework arrasta as dependências da CLI.
Separado, o core continua com **uma** dependência direta, `golang.org/x/crypto`,
mais o `golang.org/x/sys` indireto que vem com ela — é o `framework/go.mod`
inteiro.

*Corrigido em 07/08/2026:* esta linha dizia "duas dependências", e o número não
era engano solto: vem da **ADR 0017**, que escolheu contar a indireta. Os dois
jeitos de contar descrevem o mesmo `go.mod`; o que não pode é cada documento usar
um sem dizer. Aqui e no `CLAUDE.md` a contagem é a de `require` direto — **uma**
—, e a ADR 0017 continua válida no que decide, que é a CLI poder ter dependência
e o core não.

## Ordem de criação

Fase 1: `framework`, `arandu`, `aru`, `docs` — **feito**, os quatro existem,
compilam e estão publicados. `examples` entrou junto: escrever à mão o módulo que
o gerador vai emitir é o passo que transforma "template que pareceu bom" em
"template já usado"
Fase 2: `kv`, `database`, `homebrew-tap` — **feito**, os três existem e estão
publicados. *Corrigido em 07/08/2026 (ADR 0026):* o starter kit **é** um
repositório — `arandu-io/ui`, o espelho do `laravel/ui` —, e a `aru` perdeu o
`make:auth`. O que a ADR 0012 mediu foi o custo de **importar** um kit, e
concluiu contra **publicar** um: `go run <módulo>@latest` roda um módulo
publicado sem acrescentar uma linha ao `go.mod` de quem chama. O nome `oka` fica
livre; o texto abaixo é o registro do desenho anterior.

O desenho anterior dizia: `oka` saiu da lista, o starter kit é `aru make:auth`, não um
repositório (ADR 0012). O `porang` chegou a existir e foi dissolvido pela
ADR 0021 — view no esqueleto, runtime em `framework/view/`. O `database` é um
repositório com quatro módulos — o contrato fica no core e cada driver tem o
seu (ADR 0014)
Fase 3: `queue` e `storage` — **feito**. São o nono e o décimo repositórios,
autorizados pela ADR 0016; a REGRA 3 passa a valer com dez
Depois da Fase 3: `ui` (ADR 0026), `kyse` (ADR 0027), `mcp` (ADR 0030) e
`hesape` (ADR 0049) — quatro autorizações fora da ordem de fases, e a REGRA 3
passa a valer com catorze
Fase D: `website` (aplicacao Arandu, sobre o `kyse` e o `framework/view`)
Fase 4+: `taba` (painel administrativo) e a camada de IA, ainda sem nome.
Ver a triagem em `06-escopo.md`

Registrar os nomes na org desde já é gratuito e evita colisão interna. Criar os
repositórios vazios, não.
