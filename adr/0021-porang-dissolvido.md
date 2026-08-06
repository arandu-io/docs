# ADR 0021-porang-dissolvido — O `porang` é dissolvido: view mora no esqueleto, render mora no framework

**Status:** Aceita. Fecha a porta que a ADR 0010 abriu ao prometer a árvore de
views "em `arandu-io/porang`". De onze repositórios para dez.

## Contexto

O `porang` era um repositório separado com layout, componentes, assets embutidos
e o runtime de render. A justificativa era simetria: se o Arandu tem um pacote
para cada coisa, a view também tem o dela.

O Laravel não faz isso. Em Laravel, `resources/views/` mora em `laravel/laravel`
— o esqueleto, o que o desenvolvedor clona e edita. O que fica em
`laravel/framework` é o motor: `Illuminate/View`, o compilador, o localizador de
arquivo. Não existe repositório intermediário com as telas.

E o desenho não correspondia a onde o trabalho acontece. Quem abre um projeto
Arandu vai editar layout e componente todo dia. Deixar isso em outro repositório
significa que a primeira mudança de tela é um fork, ou um vendor, ou um pedido de
pull request num projeto que não é o dele. A REGRA 9 diz que não pode existir uma
segunda maneira de fazer algo que já tem uma — e "editar a view" já tinha a sua:
abrir o arquivo no projeto.

O `porang` também é a âncora do `docs/21-modelo-de-negocio.md`, que elege
`porang-pro` como primeira receita. Isso pesa a favor de manter o repositório, e
não é suficiente: repositório de código não se desenha para acomodar plano de
venda.

## Decisão

**O `porang` é dissolvido.** As peças vão para onde o Laravel as põe:

| O que | De | Para |
|---|---|---|
| layout e componentes | `porang/layout/`, `porang/components/` | `arandu/resources/views/` |
| runtime de render | `porang/render.go` | `framework/view/` |
| assets embutidos e rota | `porang/assets.go`, `porang/module.go` | `framework/view/` |

Os dez repositórios passam a ser: `framework`, `arandu`, `aru`, `database`,
`kv`, `queue`, `storage`, `examples`, `homebrew-tap`, `docs`.

Duas restrições viajam junto com o código:

1. **`framework/view` continua sendo `kernel.Module`, não função solta.** Essa
   parte não é preferência de estilo, é um bug conhecido. O `porang.Mount` era
   função exportada com **zero call sites** em todos os repositórios: o kernel
   não chamava, o `main.go` gerado não chamava, e a tela de login que o
   `aru make:auth` escreve pedia três assets — folha de estilo, HTMX e Alpine —
   e recebia 404 nos três. Foi corrigido no `porang v0.5.0` virando módulo.
   Função precisa ser lembrada; módulo aparece no `Register`, ao lado de eventos,
   jobs e agendador, que é onde alguém lendo o `main.go` já procura para saber do
   que a aplicação é feita. Reintroduzir a função reintroduz o 404.

2. **Nada de perda de trabalho.** Das 1.805 linhas do `porang`, 1.029 são geradas
   pelo templ. Com o templ fora, elas seriam regeradas pelo kyse de qualquer
   forma. O que migra de verdade é o restante.

## O que isso fecha

A pergunta "onde fica a view" tem uma resposta só, e ela é a mesma que o
desenvolvedor Laravel já sabe: no projeto, em `resources/views/`. A ADR 0010 fica
valendo na estrutura de pastas e cai na localização — onde ela dizia
`arandu-io/porang`, leia-se `arandu/resources/views/`.

Fecha também a porta de um repositório de UI de primeira parte. Componente de
primeira parte é arquivo do esqueleto, não dependência.

## O que isso não fecha, e fica em aberto

**O modelo de negócio perde a âncora.** O `docs/21-modelo-de-negocio.md` elege
`porang-pro` como primeira receita, pela única razão de haver precedente
comprovado no nicho. Sem o `porang`, não há repositório-irmão fechado para
pendurar o catálogo pago, e o argumento de que "componente premium vende para
quem usa templ puro" já tinha caído junto com o templ (ADR do `kyse`).

Isso não se resolve aqui. Vender componente continua possível — o formato é
outro, e o documento 21 precisa ser reescrito antes de a Fase 2 fechar. Registrar
o custo agora é o ponto: a dissolução é boa decisão de arquitetura e deixa uma
decisão comercial em aberto, e as duas coisas são verdade ao mesmo tempo.

## Referências

- ADR 0010 — estrutura de views no formato do Inertia
- `docs/05-repositorios.md` — a lista, agora com dez
- `docs/14-view-e-assets.md`
- `docs/21-modelo-de-negocio.md` — o que perde a âncora
- `porang v0.5.0` — a correção do `Mount` sem call site
