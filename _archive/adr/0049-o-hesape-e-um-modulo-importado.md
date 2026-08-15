# ADR 0049 — O `hesape` é um módulo importado, e o `framework` é o Foundation

**Status:** Aceita. Data: 14/08/2026. Fecha o endereço que a
[ADR 0043](0043-a-colecao-de-componentes-chama-hesape.md) deixou pendente, e
**reverte duas frases dela** — ver "O que esta ADR desfaz".

## A porta que se fecha

O `hesape` deixa de poder ser um diretório dentro do `framework`. As formas (b)
e (c) de `docs/31-reorganizacao-hesape.md`, seção 2, morrem aqui: o módulo
`github.com/arandu-io/hesape` está publicado, com a `v0.1.0` no proxy do Go, e o
proxy serve o caminho para sempre.

## Contexto

A ADR 0043 nomeou o conjunto de componentes e deixou o endereço aberto,
deliberadamente. `docs/31`, seção 2, pesou as três formas e **recomendou a (c)** —
o nome viveria nos documentos e na marca, e o código seguiria em
`arandu-io/framework/httpx`. As 576 linhas daquele documento, com as 1.070
superfícies especificadas, foram escritas sob essa premissa.

Duas coisas aconteceram depois.

A primeira: a tabela de mudança de `docs/31`, seção 3.1, foi executada — mas
dentro de um módulo próprio, não dentro do `framework`. São 173 pacotes e 1.202
arquivos, publicados como `github.com/arandu-io/hesape v0.1.0`.

A segunda: o proprietário decidiu, e a frase é a decisão — *o `hesape` é
repositório fora do `framework` e importado nele*. O fluxo é
`arandu → framework → hesape`, em repositórios separados.

O estado que essas duas produziram é o pior possível e não podia continuar:
**duas árvores vivas** com o mesmo conteúdo, e uma correção feita na errada se
perde. Cada um dos 16 pacotes de produto do `framework` tem correspondente no
`hesape`, sob o nome do Illuminate. A duplicação é inteira, não parcial.

## Decisão

**`github.com/arandu-io/hesape` é um módulo próprio, e o `framework` o declara
em `require`.**

O `framework` fica com exatamente o que o Laravel **não** publica como pacote
separado:

| `laravel/framework` | `arandu-io/framework` |
|---|---|
| `src/Illuminate/*`, com subtree split em `illuminate/*` | módulo `github.com/arandu-io/hesape` |
| `Illuminate\Foundation` — **sem** split, só existe aqui | `framework/foundation` |
| `config/`, os 17 arquivos | `framework/config` |
| `config-stubs/` | `framework/configstubs` |
| `bin/split.sh` e o `splitsh-lite` de 7 MB | nada |
| `types/`, as asserções do PHPStan | nada |

As duas últimas linhas são o argumento inteiro. O `laravel/framework` carrega
`src/Illuminate` porque o Composer **exige** monorepo mais subtree split para
publicar `illuminate/http` como repositório instalável. É artefato de
empacotamento do PHP. Em Go, `github.com/arandu-io/hesape/routing` já é
importável sozinho, sem splitter, sem binário de 7 MB versionado e sem um
segundo lugar onde o mesmo código existe.

`Foundation` não está no `replace` do `composer.json` do Laravel — não existe
`illuminate/foundation`. Ele só é entregue dentro do `laravel/framework`. É
esse o papel que sobra para o nosso `framework`, e é por isso que ele continua
existindo em vez de virar uma casca.

### O `foundation` fica em dois lugares, e isso não é uma segunda forma

`hesape/foundation` guarda o **vocabulário**: `Module` e as sete interfaces
opcionais. `framework/foundation` guarda o **boot**: `Application`,
`RendererProvider`, `Builtins` e os bootstrappers.

O corte é forçado, não estético. Os quatro adapters que a
[ADR 0048](0048-os-quatro-adapters-entram-no-hesape.md) dobrou para dentro do
`hesape` implementam `Module`. Se `Module` subisse para o `framework`, o
`hesape` passaria a requerer o `framework` — dependência circular de módulo, e
uma segunda linha num `go.mod` que tem uma só (ADR 0004).

O precedente do Laravel é exato: `Illuminate\Contracts\Foundation\Application`
**é** splitado; `Illuminate\Foundation\Application` **não é**. O contrato desce,
a implementação fica.

### A ponte tem data de morte escrita

Treze módulos `go.mod` em oito repositórios pinam `framework v0.25.4`, e
`arandu` mais `examples` sozinhos são 125 arquivos. Trocar tudo num corte só
deixa todos vermelhos ao mesmo tempo, e `aru` e `ui` quebram **em silêncio** —
as referências deles ao `framework` vivem dentro de string de gerador, onde o
`go build` não olha.

Então cada pacote antigo do `framework` vira casca fina de reexport para o
`hesape`, e os consumidores migram um a um. Dos 140 símbolos que os treze de
fato consomem, 116 sobrevivem com o mesmo nome e viram `type X = hesape.X`.

Isso cria, durante a transição, dois caminhos de import para o mesmo tipo — o
que a [REGRA 9](../09-uma-forma-so.md) proíbe. A tensão é real e é aceita **com
prazo escrito**: o `doc.go` de toda casca carrega a linha que diz em que versão
ela some, e o CI reprova casca sem prazo. Ponte sem data de morte não é ponte,
é garfo.

## O que esta ADR desfaz

Duas frases da ADR 0043 deixam de valer, e ficam registradas aqui em vez de
serem apagadas de lá:

1. *"`hesape` **não é o décimo terceiro repositório** (REGRA 3)."* Agora é. A
   autorização é esta ADR, que é o que a REGRA 3 existe para exigir.
2. *"o nome muda antes de entrar em qualquer import — que é exatamente por que o
   endereço ficou pendente"*. O nome **entrou** no import.

## O custo que esta decisão aceita, sem maquiagem

**Os filtros 0 e 2 da REGRA 4 continuam abertos.** A ADR 0043 é explícita: o
filtro 0 (voz alta em português) nunca foi avaliado por quem tem o ouvido, e o
filtro 2 (categoria vetada) tem pendência registrada no `checklist.sh` — luz e
resplendor estão no centro do Ayvu Rapyta, onde Ñamandu se cria pelo resplendor
do próprio coração, e a REGRA 4 veta divindade e ofício religioso. O filtro 3
também segue aberto: `hesape` lido por falante de inglês soa *"he's ape"*, e o
produto é global e em inglês.

A ADR 0043 amarrou o endereço a esses filtros justamente para que uma reprovação
custasse barato. Com a `v0.1.0` publicada, não custa mais: o proxy do Go serve
`github.com/arandu-io/hesape` para sempre, e uma troca de nome depois disto é a
mesma quebra que a REGRA 2 chama de irreversível no caminho do `framework`.

Isso não reabre a decisão — ela está tomada. Fica escrito porque o registro
existe para carregar o preço junto com a escolha, e porque a checagem do filtro
2 continua no `checklist.sh` como item de pessoa: se ela reprovar, o que se
decide é conviver com o nome, não trocá-lo.

## Consequência

- `framework/go.mod` ganha `require github.com/arandu-io/hesape`. A guarda de
  orçamento de dependência do CI passa a permitir `github.com/arandu-io/*` e
  continua reprovando terceiros: `arandu-io` é código nosso, e o que a guarda
  existe para barrar é dependência que o core não escreveu e não pode consertar.
- `framework/modules/` deixa de existir. A
  [ADR 0019](0019-estrutura-do-laravel.md) já o havia aposentado como estrutura
  do framework e o diretório sobreviveu à decisão.
- `docs/31-reorganizacao-hesape.md`, seção 2, fica superada quanto ao endereço.
  A seção 3.1 continua valendo inteira: ela diz o que foi para onde, e isso não
  mudou.
- `docs/05-repositorios.md` e a REGRA 3 do `CLAUDE.md` passam a contar o
  `hesape` entre os repositórios.
