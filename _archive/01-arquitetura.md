# Arquitetura

## O caminho obrigatório

```
Request
  ↓  Recover        panic → página de debug (dev) / 500 mudo (prod)
  ↓  Observe        request_id, logger de escopo, Collector
  ↓  Headers        CSP, HSTS, nosniff
  ↓  RateLimit      por IP e por subject
  ↓  CSRF           header (HTMX) ou campo de form
  ↓  Handler        extrai entrada, delega. Não fala com banco.
  ↓  Request        Validate() → Errors por campo
  ↓  Service        regra de negócio
  ↓  Policy         Authorize() → emite Grant
  ↓  Repository     exige Grant. Sem Grant não compila.
  ↓  DB             instrumentado no Collector
```

A diferença em relação ao Laravel está em uma linha: **o Grant**.

No Laravel, `Gate::authorize()` é uma chamada que você pode esquecer de fazer, e
nada te avisa. Aqui, `security.Grant` tem apenas campos não exportados. Como toda
assinatura de repositório exige um `Grant`, não existe caminho do handler até o
banco sem um `Grant`. Isso é o sistema de tipos, e não é convenção.

O que o compilador **não** decide é *qual* `Grant`. `SystemGrant` é exportado e
compila em qualquer pacote, inclusive num handler, e o `Grant` que ele devolve
passa em `Check`. Quem segura essa porta é um lint — o `aru doctor` — e a
validação de tenant em runtime. A seção abaixo é sobre exatamente isso, e a
distinção importa: prometer garantia de compilação onde há lint é a diferença
entre a promessa e o mecanismo.

### A segunda porta: `security.SystemGrant`

`Authorize` não é o único construtor público. Existem outros dois, e os dois são
deliberados:

```go
func SystemGrant(a Action, tenant string) Grant   // security
func GrantFor(j Job) security.Grant               // jobs, e delega ao de cima
```

Ele existe porque nem todo trabalho tem sessão atrás: job, seeder, tarefa do
scheduler e o próprio caminho de login rodam sem subject. Sem ele, esses
caminhos precisariam de uma policy falsa — e uma policy falsa é pior que uma
saída explícita, porque não aparece na auditoria.

O que o segura, em três camadas:

1. **A assinatura.** O tenant é obrigatório e validado. Tenant vazio ou fora de
   `^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$` devolve um `Grant` **com razão**, que
   reprova em `Check` dizendo qual foi. Não é o `Grant` zero — o zero é
   justamente o que não carrega razão nenhuma, e `Check` ramifica em
   `g.reason != ""` para distinguir os dois enganos. Um `SystemGrant` que lê todos os clientes não é
   expressável (REGRA 14, e `15-multi-tenancy.md`)
2. **O `aru doctor`**, pelo nome com que ele reporta.
   `system-grant-without-tenant` (erro) pega o tenant literal vazio antes de
   alguém esperar a falha em runtime.
   `system-grant-outside-scope` (aviso, e `--strict` transforma aviso em falha)
   reprova a chamada fora dos lugares em que o esqueleto põe trabalho sem
   request atrás: `main.go`, `routes/console.go` e qualquer arquivo sob um
   segmento `seeders`, `jobs`, `commands`, `console` ou `cmd`. É segmento de
   caminho, não nome de função — nome de função é convenção que um rename
   derrota. A exceção deliberada se escreve na linha, com motivo:
   `//arandu:system-grant <razão>`. Marcador sem razão não desculpa nada
3. **A regra `tenant-from-request`**, que segue o valor em vez de ler o nome: um
   dado que a requisição controla, atribuído a uma variável e usado como tenant
   de um construtor de `Grant`, é achado — mesmo com o nome mais inocente

   O alcance dela é declarado, e é menor do que a frase acima sugere sozinha:
   **uma função por vez**. Um valor que atravessa chamada não é seguido. Está
   escrito no comentário da própria regra, em `aru/internal/doctor/rules.go`:
   *"a doctor that claims more coverage than it has is worse than one that
   claims less"*.

   A forma embutida — `SystemGrant(ActionView, ctx.Query("org"))`, sem passar
   por variável — **passava batido até 07/08/2026**, e é a mais curta de
   escrever. A primeira passagem da regra procurava atribuição, e sem nenhuma
   ela desistia antes da segunda. Hoje é achado.

*Corrigido em 07/08/2026:* este parágrafo dizia que não há "nenhum construtor
público além de `security.Authorize`". Há, e são **dois**.

O primeiro é `security.SystemGrant`, em `framework/security/policy.go`. Quem lia
concluía que a saída de emergência não existe, e a primeira vez que a
encontrasse no código seria sem saber o que a guarda.

O segundo apareceu na revisão adversarial do mesmo dia, e não estava guardado:
**`jobs.GrantFor`**, em `framework/jobs/jobs.go`. Ele existe para o worker
reemitir exatamente o que o push autorizou, lendo a ação e o tenant de um `Job`
— e `Job` é struct exportada, que qualquer pacote preenche:

```go
jobs.GrantFor(jobs.Job{Action: "customer.delete", TenantID: ctx.Query("org")})
```

Isso chega ao banco com permissão que ninguém concedeu, sob um tenant escolhido
por quem chamou. As três regras casavam o nome literal `security.SystemGrant`,
então essa linha **não produzia achado nenhum** — nem fora de escopo, nem tenant
da requisição, nada. Reproduzido em projeto gerado: `aru doctor` respondia
`no findings`.

As regras passaram a casar um **conjunto** de construtores, e
`TestEveryGrantConstructorIsGuarded` lê o que o `framework` exporta e reprova
função nova que devolva `Grant` e não esteja na lista. A tese continua de pé, e
agora ela é verificada: as saídas são duas, são nomeadas, e cada chamada de
qualquer uma é a auditoria.

## Estrutura: a do Laravel, diretório a diretório

O esqueleto é o do Laravel, com os mesmos nomes e nos mesmos lugares —
`app/Http/Controllers`, `app/Models`, `app/Policies`, `app/Providers`,
`bootstrap/`, `config/`, `database/`, `resources/views`, `routes/`, `storage/`,
`public/`. Ver `28-estrutura-do-projeto.md`.

O objetivo é único e vale mais que qualquer elegância de organização: **o dev
que vem do Laravel abre o projeto e reconhece tudo.**

Somam-se três diretórios que o Laravel não tem, e eles são a diferença do
produto:

```
app/Services/       ← regra de negócio, obrigatório
app/Repositories/   ← acesso a dado, exige Grant
app/Policies/       ← existe no Laravel, aqui é obrigatório
```

No Laravel, Service e Repository são convenção de time organizado. Aqui são
esqueleto, e o `doctor` cobra. É a frase de onboarding virada em diretório.

## O que foi deliberadamente deixado de fora

| Laravel | Decisão aqui | Por quê |
|---|---|---|
| Service container | wiring explícito, escrito à mão em `bootstrap/app.go` | reflection em Go custa e esconde |
| Facades | não existe | exigiria estado global mutável |
| Eloquent | SQL parametrizado + repository | SQL previsível, injection impossível |
| `$fillable` | não existe | request struct explícita elimina a classe de bug |
| Blade | kyse | tipagem em tempo de compilação |

As três últimas linhas são o ponto: metade da superfície de segurança do Laravel
existe para remediar perigos de linguagem dinâmica. Portar a mitigação de um bug
que Go não permite é complexidade sem retorno.
