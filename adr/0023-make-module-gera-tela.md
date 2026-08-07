# ADR 0023-make-module-gera-tela — `aru make:module` gera as quatro telas, e devolve HTML

**Status:** Aceita.

## Contexto

O gerador emitia nove arquivos em `modules/<nome>/` e o handler respondia JSON.
Depois da ADR 0019 a árvore virou a do Laravel, e o `porang` foi dissolvido no
`framework/view` com o `kyse` como motor de template (ADR 0020, 0021).

Um gerador que entrega CRUD em JSON, num framework cuja tese é "Laravel +
Inertia sem JavaScript", entrega metade: o desenvolvedor gera o módulo e ainda
tem a tela inteira pela frente — que é justamente a parte em que ele mais erra
sozinho, porque envolve CSRF, 422, HTMX e o contrato do layout.

O doc 19 dizia "o gerador não escreve view", com a justificativa de que a forma
canônica de uma tela ainda não existia. Ela existe: o `make:auth` publicou as
nove telas do `laravel/ui` e o `layouts.app` do esqueleto fixou o contrato — uma
interface `Layout` com `PageTitle`.

## Decisão

`aru make:module` emite doze arquivos, todos na árvore do Laravel:

```
app/Models/Pedido.go                  app/Http/Requests/PedidoRequest.go
app/Policies/PedidoPolicy.go          app/Http/Controllers/PedidoController.go
app/Repositories/PedidoRepository.go  app/Http/Controllers/PedidoController_test.go
app/Services/PedidoService.go         database/migrations/<data>_create_pedidos_table.go
resources/views/pedidos/{index,show,create,edit}.kyse.go
```

O controller implementa as sete interfaces de `httpx` e responde `ctx.View` com
struct tipada por página. **Nenhum handler gerado responde JSON.**

Três consequências que fecham porta:

- **A DSL não cresceu.** Não existe `screens:` na spec. As telas saem dos campos
  que ela já declara (REGRA 15)
- **`arandu.mod.toml` sai do `make:module`.** Ele descrevia um módulo publicável
  em `modules/<nome>/`, e esse diretório não existe mais
- **O gerador não escreve fiação.** O campo em `routes.Deps`, o construtor em
  `bootstrap/app.go` e a migração em `All()` são três linhas impressas pelo
  comando e coladas à mão. Fiação explícita que o gerador edita sozinho deixa de
  ser explícita

Um detalhe técnico decide o resto do desenho: **os pacotes agora são
compartilhados**. `app/Repositories` guarda o repositório de todos os módulos, e
Go tem um pacote por diretório — então nenhum nome de pacote pode ser genérico.
`columns` virou `pedidoColumns`, `ErrNotFound` virou `models.ErrPedidoNotFound`,
`scan` e `isConflict` viraram métodos do repositório, e `NewID` foi substituído
por `data.NewID`. Sem isso, o segundo módulo gerado não compila.

## Custo conhecido

As regras do `aru doctor` continuam escritas contra `modules/` (custo já
registrado na ADR 0019) e agora também contra `arandu.mod.toml`, que o
`make:module` não emite mais.

## Referências

- ADR 0019 — a árvore é a do Laravel
- ADR 0020 — kyse
- ADR 0021 — porang dissolvido
- `docs/19-dsl-e-geracao.md`, `docs/28-estrutura-do-projeto.md`
