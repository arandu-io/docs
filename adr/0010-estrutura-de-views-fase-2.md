# ADR 0010-estrutura-de-views-fase-2 — Estrutura de views no formato do Inertia, com HTMX

**Status:** Aceita como direção. Implementa na Fase 2, com o `porang`.

## Contexto

Pedido do proprietário: quando a camada de view entrar, a estrutura de pastas
deve ser parecida com a de um projeto Laravel + Inertia + React —
`resources/js/Pages`, `Layouts`, `Components`.

Vale separar duas coisas que costumam vir juntas nessa frase:

- **A organização de arquivos do Inertia é boa e é copiável.** Uma página por
  rota, layouts reaproveitados, componentes isolados. Isso é convenção de pasta,
  não tecnologia.
- **O modelo de execução do Inertia não é.** Inertia existe para ligar um SPA
  React a um backend clássico: o servidor devolve JSON com o nome do componente e
  as props, e o React monta a tela no navegador. Com HTMX o servidor devolve HTML
  pronto, e o navegador troca um pedaço da página. Portar o modelo do Inertia
  para cima do HTMX seria pagar o custo dos dois.

## Decisão

Copiar a **estrutura**, não o mecanismo — a mesma regra que já vale para o
Laravel inteiro (`REGRA 7`).

```
views/
  layouts/     casca da página: <head>, nav, o hx-headers do CSRF
  pages/       uma por rota, o equivalente de resources/js/Pages
  components/  pedaços reaproveitáveis, tipados
  partials/    fragmentos que o HTMX troca, sem layout em volta
```

Dentro de um módulo, a mesma árvore vive em `modules/<nome>/views/` — módulo é um
diretório, e continua sendo publicável inteiro (ADR 0003).

Diferenças que valem registrar, porque são o motivo de não ser uma cópia:

1. **`partials/` não tem equivalente no Inertia.** É a pasta central no HTMX: é o
   que uma resposta de `hx-post` devolve. Ter uma pasta própria para isso deixa
   óbvio o que é página e o que é fragmento.
2. **Props tipadas em vez de JSON.** No Inertia as props atravessam serialização
   e chegam sem tipo no React. Aqui uma página é uma função `templ` com
   parâmetros — se o handler passar o valor errado, não compila.
3. **Sem build de JavaScript.** Não há `vite.config.js`, não há `resources/js`.
   O HTMX é um script, e o Basecoat (MIT) dá a fundação visual.

## Consequência

Fica claro o que a Fase 2 tem que entregar em `arandu-io/porang`: essa árvore, o
template base já com `hx-headers` do CSRF, e o gerador do `aru make:module`
emitindo página, partial e componente nos lugares certos.

Revisitar apenas se o produto passar a precisar de SPA de verdade — o que seria
decisão de stack, e mudaria `01-arquitetura.md` antes deste ADR.
