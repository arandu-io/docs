# Uma forma só

## O problema, com evidência

O Laravel virou três ecossistemas dentro de um. Isso não é opinião: está no tree
clonado em `reference_laravel/`.

**Três camadas de frontend concorrentes**

| Onde | O quê |
|---|---|
| `framework/src/Illuminate/View/Compilers/BladeCompiler.php` | Blade |
| `reverb/src/Pulse/Livewire/`, fixture `pint/.../components/flux.blade.php` | Livewire + Flux |
| `octane/src/Listeners/PrepareInertiaForNextOperation.php` | Inertia + React/Vue |

**Três runtimes**

`octane/src/Commands/StartSwooleCommand.php`, `StartRoadRunnerCommand.php`,
`StartFrankenPhpCommand.php` — cada um com client, ServerStateFile e
ProcessInspector próprios.

**Node dentro do core**

- `framework/src/Illuminate/Support/NodePackageManagers/` — classes `Bun.php`,
  `Npm.php`, `Pnpm.php`, `Yarn.php`. Um framework PHP com suporte a quatro
  gerenciadores de pacote JavaScript.
- `framework/src/Illuminate/Foundation/resources/exceptions/renderer/` tem
  `package.json`, `package-lock.json` e `vite.config.js`. **A página de erro tem
  um build npm próprio dentro do framework.**
- `pint/resources/js/worker.cjs` e `pint/app/PrettierFormatters/` — o formatador
  de PHP precisa de Node para formatar Blade.

**Dashboards em Vue**

`horizon/resources/js/**/*.vue` e `telescope/resources/js/**/*.vue`, com
`webpack.mix.js` e `vite.config.js` cada.

Nada disso foi erro de projeto. Cada peça foi resposta legítima a um pedido real.
O resultado somado é que o desenvolvedor novo precisa escolher três vezes antes
de escrever a primeira linha.

## A regra

**Uma forma só de fazer cada coisa. Sem alternativa oficial.**

| Dimensão | A única forma | Nunca |
|---|---|---|
| View | kyse + HTMX | segundo motor de template; opção React/Vue; Alpine só onde HTMX não alcança |
| CSS | Tailwind v4, binário standalone | segundo framework de CSS; pacote npm |
| Runtime | o binário Go | variantes de servidor, modo "turbo", modo alternativo |
| Autorização | `security.Grant` | atalho, bypass, "modo simples" |
| Dados | SQL parametrizado + `Repository` | ORM, query builder fluente, segunda camada |
| Build | `go build` | npm, Vite, bundler, Node em qualquer lugar do core |
| Config | struct tipada | mapa, config dinâmica, arquivo por ambiente |

O teste antes de aceitar qualquer coisa nova: **isso cria uma segunda maneira de
fazer algo que já tem uma?** Se sim, não entra — nem como opcional, nem como
experimental, nem "só para quem precisa". Opcional vira obrigatório assim que o
primeiro tutorial popular usa.

## O gerado mora em um lugar só

`resources/views/home.kyse.go` é o que se escreve; `storage/framework/views/home.go`
é o que o `aru view:build` emite. Um artefato, um lugar.

Esta seção dizia o contrário — que o gerado ficava **ao lado**, no mesmo
diretório, e que `storage/framework/views/` **não existia**. Nunca foi assim no
código: `aru/internal/kyse/kyse.go` escreve sob `storage/` desde que a árvore do
Laravel entrou, o `.gitignore` de todo projeto ignora aquele diretório, e o
comentário de doc exportado do `OutputPath` afirmava "beside its source" com um
exemplo falso — publicado no `pkg.go.dev`. Corrigido em 10/08/2026.

O motivo de ser sob `storage/` é o mesmo que fez o Laravel escolher aquele
diretório: **o gerado é saída de build, e saída de build não se versiona**. Ao
lado do fonte, cada `aru view:build` sujaria o `git status` com 28 arquivos que
ninguém escreveu e que o revisor teria que ignorar em toda revisão. Sob
`storage/`, uma linha de `.gitignore` resolve, e o diretório inteiro pode ser
apagado sem perder nada.

O que a regra proíbe continua valendo, e é o que o Laravel faz de fato errado
ali: nome derivado de hash. O nosso gerado espelha a árvore do fonte
(`auth/login.kyse.go` → `auth/login.go`), é Go de verdade, é checado pelo
compilador, abre no editor e aparece num diff. A pergunta "onde está o código que
rodou" tem uma resposta só, e ela é legível.

## A página de erro é o caso-teste

Já está implementada em `observability/errorpage/page.go`: HTML e CSS inline num
`const` de `html/template`. Sem npm, sem Vite, sem `package.json`.

Se um dia alguém propuser "melhorar a página de erro com um build de assets", a
resposta é essa seção. Foi exatamente assim que o Node entrou no Laravel.

## Sobre "modos de performance"

Modo é a mesma doença com outro nome. Performance no Arandu vem de **remover
camada**, não de adicionar modo:

- sem ORM: o gerador emite a query e a struct juntas, plano previsível
- sem reflection: wiring explícito, resolvido em compilação
- sem boot por request: o binário sobe uma vez (é o que torna Octane desnecessário)
- sem hidratação: fragmento de HTML, não JSON + render no cliente
- Collector no-op em produção: receiver nil, zero alocação

Quando faltar desempenho, o caminho é ajustar o que existe — índice, pool,
réplica de leitura, cache — nunca ligar um segundo caminho de código.
