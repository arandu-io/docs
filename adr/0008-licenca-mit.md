# ADR 0008-licenca-mit — Licença MIT, a mesma do Laravel

**Status:** Aceita. Fecha a decisão que estava aberta em `00-indice.md`.

## Contexto

A licença precisava estar definida antes da Fase 2 e antes de qualquer
repositório virar público — depois de publicado, mudar para uma licença mais
restritiva não recupera o que já foi distribuído sob a anterior.

## Decisão

MIT, no mesmo formato do Laravel: arquivo `LICENSE.md`, texto padrão MIT,
`Copyright (c) Paulo Lima`, sem ano — o Laravel omite o ano, e ano fixo obriga a
editar o arquivo todo janeiro.

Vale nos quatro repositórios: `framework`, `arandu`, `aru` e `docs`.

Por que MIT e não Apache 2.0, que é o outro padrão do ecossistema Go: MIT é o que
o mercado de framework web espera — Laravel, Rails, Django (BSD), React e Vue são
todos permissivos e curtos. Apache 2.0 traz concessão de patente explícita, o que
é vantagem real para quem tem portfólio de patente e irrelevante aqui, ao custo de
um arquivo que ninguém lê e de um `NOTICE` para manter.

Por que não algo restritivo, tipo BSL ou AGPL: a tese do produto é adoção. Um
framework que a empresa precisa passar pelo jurídico antes de experimentar perde
para o que ela pode simplesmente usar. Os produtos pagos previstos
(`taba`, na Fase 4) são repositórios separados e podem ter licença própria —
MIT no core não obriga MIT no que vem depois.

## Divergência registrada

`20-licenca.md` recomendava Apache-2.0 pela cláusula de patente. O proprietário
decidiu MIT, e a decisão prevalece — o documento foi atualizado para registrar as
duas posições em vez de esconder a que perdeu.

O que se adota do documento 20 mesmo com MIT: **DCO, não CLA**. O contribuidor
assina `git commit -s` afirmando que pode contribuir; ninguém cede direito. Como
o componente pago vive em repositório separado com licença própria, o core nunca
precisa ser relicenciável, e CLA seria custo sem retorno.

## Consequência

O core fica utilizável comercialmente sem contrapartida, inclusive por quem
construa um produto concorrente. É o preço da adoção, e é o mesmo preço que o
Laravel paga.

Revisitar apenas para repositório novo e não publicado — nunca para código que já
saiu sob MIT.
