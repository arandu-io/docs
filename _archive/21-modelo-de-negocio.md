# 21 — Modelo de negócio

Os documentos de arquitetura não têm nenhuma linha sobre o que se vende. Isso não
é lacuna comercial: é lacuna de arquitetura, porque a resposta muda a prioridade
da Fase 2.

**Estado deste documento: a âncora caiu.** A primeira receita eleita aqui era o
`porang-pro`, e o `porang` foi dissolvido pela ADR 0021. O raciocínio de
licença e de separação continua de pé; o produto concreto que ele vendia, não.
Ver "A âncora que caiu com o `porang`", abaixo.

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
| **Componentes premium** (catálogo pago, em repositório separado) | design de qualidade, não engenharia | caminho preferido; o formato está em aberto |
| **Admin gerado** (equivalente Nova) | muita UI, licença por projeto | Fase 3+ |
| **Cloud e deploy** | operar infra 24/7, time de plantão | não, sozinho |
| **Suporte e licença empresarial** | seu tempo, um cliente por vez | receita cedo, não escala |

## Recomendação

Componente pago primeiro. Três razões, e nenhuma delas dependia do nome do
repositório:

1. É o único modelo com precedente de receita comprovado no nicho.
2. Casa com a licença de `20-licenca.md`: core MIT, pago em repositório separado,
   com licença comercial própria.
3. Não compete com o GitHub nem com a Vercel. Compete com o tempo do dev montando
   tela de CRUD, e esse tempo tem preço conhecido.

## O que a decisão do `kyse` muda aqui

Motor de template próprio é ativo e é passivo ao mesmo tempo.

**Ativo:** o catálogo pago passa a ser insubstituível. Componente escrito em
`.kyse.go` não roda em projeto templ, então ele só serve a quem adotou o Arandu.
Isso é bom para receita recorrente e ruim para alcance — o Goilerplate vende
blocos para qualquer projeto templ, nós não.

**Passivo:** o `kyse` precisa de LSP, realce de sintaxe para VS Code, GoLand e
Neovim, e mensagem de erro com arquivo e linha. Isso é manutenção permanente que
o templ dava de graça, e sai do orçamento da Fase 2.

Consequência prática: **o catálogo pago deixa de vender para o mercado templ.**
O teste barato muda — em vez de publicar blocos e medir demanda geral, você mede
adoção do Arandu primeiro. É teste mais lento e mais caro.

## A âncora que caiu com o `porang`

A ADR 0021 dissolveu o `porang`: layout e componentes foram para
`arandu/resources/views/`, o runtime de render para `framework/view/`, e os onze
repositórios viraram dez. A decisão é boa de arquitetura e deixa um custo aqui.

**Não existe mais repositório-irmão de UI para pendurar o catálogo pago.** O
`porang-pro` era a primeira receita porque havia um `porang` aberto embaixo dele,
no mesmo desenho do templUI e do Goilerplate. Sem o `porang`, componente de
primeira parte é arquivo do esqueleto, não dependência — e a ADR fecha
explicitamente a porta de um repositório de UI de primeira parte.

O que continua valendo:

- o pago fica **fora do core**, em repositório separado, com licença comercial
  própria (`20-licenca.md`)
- o core continua **MIT**, e não precisa ser relicenciável — por isso DCO e não CLA
- o que se vende é tempo de montar tela, não o framework

**O que está em aberto, e não se resolve neste documento:** qual é o produto
concreto. Não há candidato decidido, e inventar um aqui seria pior que registrar
a lacuna. O documento precisa ser reescrito com uma resposta antes de a Fase 2
fechar — é o prazo que a própria ADR 0021 fixou.

Enquanto isso, o teste barato que este documento propunha não tem como rodar: ele
dependia de publicar blocos durante a Fase 2 e medir demanda. Sem produto
definido, a validação de receita fica para depois do framework servir HTML de
verdade.

## O que fica de fora, e continua de fora

Verticais brasileiros (NF-e, SPED, PIX) — decisão sua, registrada em
`06-escopo.md`. Consequência assumida: abrimos mão do fosso regulatório que
nenhum concorrente global replicaria, em troca de um produto que vende em inglês
para o mundo. Se algum dia entrarem, entram como módulo de comunidade na registry,
fora do namespace de primeira parte.
