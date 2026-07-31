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

### Núcleo

```
arandu-io/framework     a biblioteca. Módulo: github.com/arandu-io/framework
                        Deps: stdlib + golang.org/x/crypto. Nada além.
arandu-io/arandu        esqueleto de projeto. O que `aru new` clona.
                        Ninguem depende dele como biblioteca.
arandu-io/aru           CLI. Módulo próprio, binário `aru`.
arandu-io/docs          markdown do guia. Branch por versão (1.x, 2.x).
arandu-io/website       site que renderiza o docs.
                        E uma APLICACAO ARANDU — dogfood. Ver 08-documentacao.md
arandu-io/.github       workflows e templates compartilhados da org.
arandu-io/homebrew-tap  fórmula do `aru` (nome da fórmula: aru, não arandu).
```

### Adapters (drivers de infraestrutura)

Existem para manter o core sem dependência de terceiros (ADR 0004).

**Esta seção foi substituída por `11-adapters.md`**, que resolveu um problema que
não estava visto aqui: em Go não existe dependência opcional, então um pacote por
fornecedor multiplica repositório sem isolar dependência nenhuma. A divisão é por
**domínio**, com driver pesado em submódulo:

```
arandu-io/database    Repository + pgx; submódulos mysql, sqlite
arandu-io/kv          go-redis (RESP): cobre Redis, Dragonfly, Valkey, KeyDB
arandu-io/mail        interface + SMTP; submódulos ses, resend, postmark, mailgun
arandu-io/storage     interface + disco local; submódulos s3, gcs, azure
arandu-io/queue       interface + driver sobre kv e sobre database; submódulos nats, sqs
arandu-io/ws          coder/websocket e gobwas/ws
```

Seis pacotes no lugar dos quinze que a lista por fornecedor produziria. É o mesmo
que o Laravel faz — `Illuminate\Mail` traz SMTP, Mailgun, Postmark, SES e Resend
juntos; `Illuminate\Database` traz MySQL, Postgres e SQLite juntos.

Convenção de nome: **sem prefixo redundante e em inglês funcional**. É
`arandu-io/database`, não `arandu-io/arandu-postgres` — a org já carrega o nome, e
o nome tupi fica nos produtos (`porang`, `marandu`), não na infraestrutura.

### Módulos de primeira parte

O equivalente aos pacotes first-party do Laravel. Cada um implementa
`kernel.Module` e nada mais.

```
Nomes definidos em `docs/07-pacotes.md`. Os principais:

arandu-io/oka        starter kit (Breeze)
arandu-io/curupira   tokens de API (Sanctum)
arandu-io/oka        starter kit (Breeze)
arandu-io/marandu    inspetor de requests em dev (Telescope)
arandu-io/maraca     performance da aplicação (Pulse)
arandu-io/boitata    APM em produção (Nightwatch)
arandu-io/katu       formatador de código (Pint)
arandu-io/igara      ambiente Docker de dev (Sail)
arandu-io/uirapuru   servidor websocket (Reverb)
arandu-io/porang     camada de view em templ
arandu-io/moitara    registry de módulos
arandu-io/taba       painel administrativo — pago (Nova)
```

Nota: as primitivas de autenticação (Grant, Policy, sessão, argon2id) ficam no
**core**, não aqui. Segurança é a tese do produto; não pode ser um pacote
opcional que o usuário esquece de instalar. `starter-auth` traz só as telas.

### UI e distribuição

```
arandu-io/porang     camada de view em templ, ligada aos contratos do core
arandu-io/moitara    registry de módulos — ver 18-marketplace.md
arandu-io/examples   apps de exemplo
```

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
Separado, o core continua com duas dependências.

## Ordem de criação

Fase 1: `framework`, `arandu`, `aru`, `docs` — **feito**, os quatro existem,
compilam e estão publicados
Fase 2: `database`, `kv`, `porang`, `oka`, `homebrew-tap`
Fase D: `website` (aplicacao Arandu, depende do porang)
Fase 4+: `taba` (painel administrativo) e a camada de IA, ainda sem nome.
Ver a triagem em `06-escopo.md`

Registrar os nomes na org desde já é gratuito e evita colisão interna. Criar os
repositórios vazios, não.
