# 14 — View e assets (`porang`)

`porang` era a única camada nomeada sem especificação. Este documento fecha isso
e resolve uma contradição que já estava escrita nos documentos anteriores.

## A contradição

`09-uma-forma-so.md` proíbe: *Build → `go build`; nunca npm, Vite, bundler, Node
em qualquer lugar do core.*

`07-pacotes.md` abre: *única exceção onde npm importa: `porang`, se distribuir
plugin de Tailwind (o Basecoat é distribuído por npm).*

As duas não podem valer ao mesmo tempo. E a exceção é precisamente o caminho pelo
qual o Node entrou no Laravel — a página de erro com `package.json` e
`vite.config.js` dentro do framework não nasceu de uma decisão de arquitetura,
nasceu de uma exceção pontual como essa.

**A exceção cai. A regra fica.**

## A restrição que decide tudo

Um projeto Arandu tem que continuar sendo:

```
git clone && aru dev
```

Sem `node_modules`, sem `package.json`, sem lockfile de JavaScript, sem Node
instalado na máquina. Se o desenvolvedor precisar de Node para rodar um projeto,
a frase de venda morre e a REGRA 9 vira decoração.

Isso não proíbe ter etapa de build. Proíbe que a etapa de build seja Node.

## Toolchain: só binário

Duas ferramentas, ambas binário único, ambas gerenciadas pela `aru` — baixadas na
versão fixada em `arandu.toml` e cacheadas em `~/.arandu/bin`, como o Go faz com
o toolchain:

| Ferramenta | Faz | Por que não viola |
|---|---|---|
| `templ` | `.templ` → `.go` | binário Go |
| `tailwindcss` (standalone CLI) | `.css` → `.css` | binário único, sem runtime Node |

O Tailwind publica uma CLI standalone exatamente para quem não quer Node no
projeto. É essa que usamos. **Nunca** `npx tailwindcss`, nunca `@tailwindcss/cli`
via npm, nunca `package.json` no esqueleto gerado.

`aru dev` roda os dois em watch mais o servidor. Um comando, um terminal.

## Assets: `go:embed`, sem CDN

`htmx.min.js`, `alpine.min.js` e o CSS compilado entram no binário via `go:embed`
e são servidos por `porang` em `/_arandu/assets/<hash>/…` com cache imutável.

Sem CDN, e não por preferência: o middleware `SecurityHeaders` define
`script-src 'self'`. Servir HTMX pelo unpkg exigiria afrouxar a CSP — pagar
segurança por conveniência. É uma diferença observável em relação ao Gojang, que
carrega htmx por unpkg.

Consequência que vale explicitar: **o deploy continua sendo um binário.** Sem
passo de publicação de assets, sem CDN para invalidar, sem `storage:link`.

## Basecoat: vendorizar, não depender

O CSS entra vendorizado no repositório do `porang`, com o arquivo de licença e a
atribuição. npm não entra no grafo de dependências em nenhuma hipótese.

Se manter o vendor ficar caro, a alternativa é escrever os componentes — não é
reintroduzir npm.

## Alpine: o limite, escrito

A REGRA 9 diz *"Alpine só onde o HTMX não alcança"*. Sem um teste, isso é opinião,
e opinião não segura ninguém em revisão de PR.

**Alpine é permitido quando o estado é, ao mesmo tempo:**

1. puramente do cliente,
2. efêmero — morre no reload sem perda,
3. invisível para o servidor.

Permitido: dropdown aberto/fechado, aba selecionada, foco, máscara de input,
confirmação antes de submeter.

Proibido: qualquer coisa que vá ao servidor, qualquer coisa persistida, qualquer
coisa derivada de dado que o servidor já tem.

O teste em uma frase: **se você está escrevendo regra de negócio em Alpine, o
componente devia ser um fragmento HTMX.**

Verificação no `aru doctor`: reprova `x-data` que contenha `fetch`, `axios`,
`$store` global, ou `$watch` com efeito de rede.

## Componente é função templ

Não existe `<nx-modal />`, `<x-card>` nem web component registrado em runtime.
Componente é função templ tipada:

```go
@porang.Modal(porang.ModalProps{Titulo: "Cliente", Tamanho: porang.LG})
```

Tag customizada seria um segundo sistema de componentes ao lado do templ, com
resolução em runtime e erro só na tela. O templ já dá composição, tipagem e erro
em tempo de compilação.

## O que `porang` não toca

A página de erro. Ela continua `const` de `html/template` em
`observability/errorpage`, sem depender de `porang`, sem CSS compilado e sem
asset embutido. Motivo: ela precisa renderizar quando o resto está quebrado —
inclusive quando o build de assets falhou.

## Perfil e modo API

`aru new --api` não instala `porang`. Nada mais muda: o contrato do módulo é
`Routes(r *httpx.Router)`, e o handler decide se devolve fragmento ou JSON.
