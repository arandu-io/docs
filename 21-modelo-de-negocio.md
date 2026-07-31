# 21 — Modelo de negócio

Catorze documentos de arquitetura sem nenhuma linha sobre o que se vende. Isso não
é lacuna comercial: é lacuna de arquitetura, porque a resposta muda a prioridade
da Fase 2.

## O que o mercado já respondeu

- **Framework não monetiza.** Concorrência gratuita infinita e público que resolve
  por conta própria.
- **O Laravel não vende o Laravel.** Vende Forge, Vapor, Nova e Cloud — serviços e
  painel.
- **O Goilerplate vende 222 blocos premium**, não o boilerplate. É o único
  precedente de receita nesse nicho exato (Go + templ + HTMX), e o autor mantém o
  templUI aberto embaixo.

## Os quatro caminhos

| Caminho | Exige | Veredito |
|---|---|---|
| **Componentes premium** (`porang-pro`) | design de qualidade, não engenharia | **começar por aqui** |
| **Admin gerado** (equivalente Nova) | muita UI, licença por projeto | Fase 3+ |
| **Cloud e deploy** | operar infra 24/7, time de plantão | não, sozinho |
| **Suporte e licença empresarial** | seu tempo, um cliente por vez | receita cedo, não escala |

## Recomendação

`porang-pro` primeiro. Três razões:

1. É o único modelo com precedente de receita comprovado no nicho.
2. Casa com a licença de `20-licenca.md`: core Apache, pago em repositório separado.
3. Não compete com o GitHub nem com a Vercel. Compete com o tempo do dev montando
   tela de CRUD, e esse tempo tem preço conhecido.

## A consequência de arquitetura

Se a receita está no `porang`, **o `porang` sobe de prioridade**. Hoje ele está
depois do core na Fase 2. Deveria entrar junto, porque é o artefato que valida
demanda antes de o framework estar completo — e componente premium pode ser
vendido para quem usa templ puro, sem adotar o Arandu.

Isso inverte o teste barato que estava implícito no plano: em vez de terminar o
framework e depois descobrir se alguém paga, publicar um pacote de blocos `porang`
durante a Fase 2 e medir. Se ninguém compra bloco, ninguém compra ecossistema.

## O que fica de fora, e continua de fora

Verticais brasileiros (NF-e, SPED, PIX) — decisão sua, registrada em
`06-escopo.md`. Consequência assumida: abrimos mão do fosso regulatório que
nenhum concorrente global replicaria, em troca de um produto que vende em inglês
para o mundo. Se algum dia entrarem, entram como módulo de comunidade na registry,
fora do namespace de primeira parte.
