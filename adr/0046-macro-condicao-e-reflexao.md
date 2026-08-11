# ADR 0046 — Macro, condição e reflexão: o que some, o que fica com outro nome, e o que já entrou sem pedir

**Status:** Aceita. Data: 10/08/2026.
**Supersede parcialmente:** `docs/31-reorganizacao-hesape.md:390`, a linha que
lista `->when()`/`->unless()`/`Conditionable` "em qualquer forma" entre o que a
REGRA 9 rejeita. A rejeição geral cai; entra uma regra estreita, na seção 2.

## Contexto

Três componentes do Illuminate — `Macroable`, `Conditionable` e `Reflection` —
estão marcados em `docs/31` como "não vira pacote nenhum". A marcação está certa
nos três. Ela também responde à pergunta errada.

"Não vira pacote" diz o que **não** existe. Quem está portando uma aplicação
Laravel abre o código com um `Str::macro` na mão e precisa saber o que escrever
no lugar. Uma recusa que não termina em Go concreto não é uma decisão de
arquitetura, é um buraco com justificativa.

Este documento fecha os três com código, e no caminho corrige uma das
marcações — porque a medição contra o clone mostrou que ela já estava
contradita pelo próprio `hesape`.

---

## 1. Macroable — some, e o que se escreve são três formas

`Traits/Macroable.php` tem cinco membros: `macro` registra um closure sob um
nome, `mixin` registra de uma vez todos os métodos públicos e protegidos de um
objeto, `hasMacro` pergunta se o nome está tomado, `flushMacros` esvazia o
registro, e `__call`/`__callStatic` desviam para ele a chamada que a classe
nunca declarou.

Os cinco terminam no mesmo lugar: uma chamada que o tipo não declarou. Em Go o
receptor de um método tem de ser declarado no pacote que declara o tipo, e não
existe gancho para chamada que não resolve — o compilador rejeita o programa.
É a **razão 1** de recusa (interface de linguagem que Go não tem), e é o
componente inteiro.

**A decisão é que existem três substitutos, e qual deles se usa é decidido pelo
que o comportamento novo precisa do valor.** Estão escritos, com código, em
`hesape/macroable/doc.go`. Em resumo:

### (a) Função livre, quando o comportamento só lê o argumento

É a resposta para `Str`, `Arr`, `Number` — todo pacote cuja superfície é função
sobre valor. `Str::macro` é a macro mais comum de uma aplicação Laravel e também
a mais barata de perder, porque já era uma função fingindo ser método.

```go
// Illuminate: Str::macro('initials', ...); Str::initials('Paulo Lima');
// hesape, no pacote que queria a macro:
package report

func Initials(value string) string {
	var out strings.Builder
	for _, word := range strings.Fields(str.Squish(value)) {
		out.WriteString(strings.ToUpper(word[:1]))
	}
	return out.String()
}
```

A chamada vira `report.Initials(name)` em vez de `str.Initials(name)`. A
diferença é o ganho, não o custo: o pacote dono do comportamento aparece na
chamada.

### (b) Tipo próprio embutindo o do framework, quando o comportamento encadeia

É a resposta para `Collection` e para tudo que devolve o receptor. O embutimento
promove os métodos do tipo embutido, então o método novo fica ao lado de
`Filter`, `Map` e `Sort`, e não num segundo vocabulário.

```go
// Illuminate: Collection::macro('published', fn () => $this->filter(...));
type Posts struct {
	collections.Collection[Post]
}

func (p Posts) Published() Posts {
	return Posts{p.Filter(func(post Post, _ int) bool { return post.Published })}
}
```

O custo é uma conversão em cada fronteira — `Posts{c}` na entrada,
`p.Collection` na saída.

### (c) Interface pequena, quando o comportamento cruza fronteira de pacote

É a resposta para `mixin`, o único membro sem equivalente nos dois acima. Lá se
enxerta um objeto inteiro numa classe e uma função escrita contra a classe ganha
uma capacidade em silêncio. Em Go a mesma forma é uma interface, **declarada no
pacote que consome** e satisfeita por quem quiser ser passado. Nada é enxertado;
a exigência fica escrita.

```go
package feed

// Cacheable é o que este pacote exige de um corpo que lhe pedem para escrever.
type Cacheable interface {
	MaxAge() time.Duration
}
```

### O que se perde, e é perda de verdade

**Um pacote de terceiro não consegue mais acrescentar método ao `Collection` de
todo mundo.** No Laravel, um `spatie/laravel-collection-macros` registra seus
métodos uma vez e toda `Collection` do processo passa a tê-los, inclusive as
construídas por código que nunca ouviu falar do pacote. Nenhuma das três formas
faz isso: `report.Initials` tem de ser importada pelo nome, `Posts` tem de ser
construída, e uma função tipada contra `collections.Collection[Post]` não
enxerga `Published` faça o que fizer qualquer outro pacote.

Isso elimina uma categoria de pacote de ecossistema que funcionava bem em PHP, e
põe uma conversão explícita em toda fronteira onde um tipo embrulhado encontra
uma assinatura do framework. O que se compra: toda chamada de método nomeia o
pacote que a declarou, nenhum import muda o sentido de código onde ele não
aparece, e `flushMacros` — que só existe porque macro é estado global mutável
que vaza entre testes — não tem o que esvaziar.

---

## 2. Conditionable — `when` e `unless` entram, e só em tipo fluente

Esta é a marcação que estava errada, e a evidência estava dentro do repositório.

`docs/31:390` põe `->when()`/`->unless()` entre o que a REGRA 9 rejeita, com o
argumento de que "Go tem `if` em posição de statement". Enquanto isso,
`hesape/collections/enumerates.go:544-590` já carrega **seis** desses métodos —
`When`, `Unless`, `WhenEmpty`, `WhenNotEmpty`, `UnlessEmpty`, `UnlessNotEmpty` —
escritos, testados e com o comentário de doc nomeando a trait. Registro que
afirma o contrário do código é registro que ninguém vai consultar duas vezes.

### O argumento, contra o `database/query/builder.go` que acabou de ser escrito

O builder é fluente de verdade: `Where`, `OrderBy`, `Limit`, `Join` e mais
cinquenta métodos mutam e devolvem o receptor. Com filtro condicional, o `if`
custa isto:

```go
q := db.Table("posts")
if published {
	q = q.Where("published", true)
}
q = q.OrderBy("id")
```

contra isto:

```go
db.Table("posts").
	When(published, func(q *query.Builder) { q.Where("published", true) }, nil).
	OrderBy("id")
```

A primeira forma não é errada e ninguém fica impedido de escrevê-la. Ela é mais
longa por uma variável que existe só para ser reatribuída — e a reatribuição é
onde o builder é usado velho: `q.Where(...)` sem o `q =` na frente **compila** e
descarta o filtro. A falha tem uma forma só, e é um filtro silenciosamente não
aplicado — a mesma classe do `WHERE IN` vazio compilado para nada.

### Por que isso não é a segunda forma que a REGRA 9 proíbe

`if` decide se um **statement** roda. `When` é um elo numa cadeia de chamadas,
que é uma **expressão**. Dentro da cadeia o `if` não é utilizável sem cortá-la e
inventar um nome intermediário; fora dela o `When` não tem em que se pendurar.
Não há duas maneiras de fazer a mesma coisa: há uma maneira para cada uma das
duas coisas.

### A regra estreita

**`When` e `Unless` entram num tipo quando todos os métodos dele devolvem o
receptor, para que as chamadas encadeiem. Só então.** Em tipo cujo método
devolve resultado, `when` é `if` com parênteses e não entra.

Onde estão hoje:

| Tipo | Métodos |
|---|---|
| `collections.Collection[T]` | `When`, `Unless`, `WhenEmpty`, `WhenNotEmpty`, `UnlessEmpty`, `UnlessNotEmpty` (já existiam) |
| `query.Builder` | `When`, `Unless` (`hesape/database/query/conditionable.go`, escrito aqui) |

Candidatos que a regra alcança e que ficam por escrever, cada um no seu pacote:
`mail.Pending`, `httpx` client `PendingRequest`, e o registrador de rota — os
três são fluentes, e os três usam `Conditionable` no Laravel. Ficam registrados
como pendência, não como decisão nova.

### As três mudanças mecânicas, e o que é pulado

A condição é `bool` e não o `mixed` truthy do PHP: não há truthiness em que se
apoiar, e quem escreve a expressão é quem sabe o que a torna verdadeira. O
callback muta o builder e não devolve nada, em vez de devolver um valor que o
PHP depois coalesce com `$this` — todo método do `Builder` já muta e devolve o
receptor, então callback que devolvesse algo seria uma segunda convenção.
Callback `nil` com condição verdadeira deixa o valor intocado em vez de estourar,
que é o que `Collection.When` já faz.

Pulado, pela **razão 1**: `HigherOrderWhenProxy`, e com ele as formas de zero e
de um argumento de `when`/`unless`. O proxy existe para capturar a próxima
leitura de propriedade ou chamada de método por `__get`/`__call` e aplicar a
condição a ela. A forma de três argumentos — condição, callback, default — é
tudo o que o `hesape` carrega, e é a que aplicação Laravel escreve na esmagadora
maioria das vezes.

Não entra também a condição como closure, `when(fn ($q) => ...)`: ali o
argumento é um `bool` que o chamador calcula antes da chamada, que é a mesma
expressão uma linha acima. Carregar as duas seria duas grafias de uma coisa.

---

## 3. Reflection — não vira pacote, e aqui está o que já usa `reflect`

`Reflector` lê uma classe que lhe entregaram: `isCallable` pergunta se o par
`[$objeto, 'metodo']` resolve, `getParameterClassName` e o plural leem os type
hints de um parâmetro, `getClassAttributes` sobe a cadeia de herança colhendo
atributos PHP 8, e `isParameterBackedEnumWithStringBackingType` faz uma pergunta
que nenhum tipo Go coloca. `ReflectsClosures` faz o mesmo com um closure, para
que `Event::listen(fn (OrderShipped $e) => ...)` descubra a que evento se
inscreve pelo tipo do parâmetro. O `helpers.php` acrescenta `lazy` e `proxy`,
sobre `newLazyGhost` e `newLazyProxy`.

Todos servem a um chamador que existe no Laravel e não existe aqui: o container
resolvendo um construtor (ADR 0001), a facade desviando uma chamada estática
(ADR 0002), o dispatcher inferindo inscrição a partir de assinatura. Os
chamadores foram embora primeiro; o leitor ficou sem para quem ler. O listener
que o PHP infere de um closure aqui é um parâmetro de tipo, conhecido do
compilador. Onde a informação é mesmo necessária, a `aru` a lê de `go/ast` em
tempo de build, em `aru/internal/typeinfo`, e emite Go.

**Recusar construir o pacote vale menos do que saber onde a reflexão entrou
assim mesmo.** Então foi contada. Em 487 arquivos `.go`, o `hesape` importa
`reflect` em **quatro arquivos não-teste**, em **doze pontos de chamada**:

| Arquivo:linha | Símbolo | Para quê | Dá para trocar por genérico ou interface? |
|---|---|---|---|
| `support/arr/arr.go:57` | `Accessible` | `Arr::accessible`: o valor pode ser indexado — mapa, slice ou array | **Não**, enquanto a assinatura for `any`. Um type switch sobre `map[string]any` e `[]any` responderia `false` para `map[string]int`, que é resposta errada, não resposta mais estrita |
| `support/arr/arr.go:169` | `typeName` | o `gettype()` de dentro das mensagens de erro de tipo | **Cosmético.** `fmt.Sprintf("%T", v)` é o mesmo `reflect` com outra grafia; troca o import, não o mecanismo |
| `support/arr/arr.go:375-384` | `From` | desreferenciar ponteiro e escolher o ramo antes do round-trip por JSON | **Não.** O `reflect` aqui só escolhe o ramo; quem faz o trabalho é a reflexão do `encoding/json`, e ela sai junto com a assinatura `any`, não antes |
| `support/arr/arr.go:500` | `IsAssoc` | `Arr::isAssoc` | **Não**, mesma razão da primeira linha |
| `support/arr/arr.go:509` | `IsList` | `Arr::isList` | **Não**, mesma razão |
| `support/arr/list.go:297` | `Wrap` | `Arr::wrap`: percorrer slice de qualquer tipo para produzir `[]any` | **Parcial.** `Wrap[T any](s []T) []any` cobre o caso tipado e deixa o caso `any` ainda em `reflect` — duas grafias de `Arr::wrap`, que a REGRA 9 recusa |
| `testing/t.go:200` | `assertIsArray` | `PHPUnit::assertIsArray` sobre valor decodificado | **Não**, mesma razão da primeira linha |
| `testing/t.go:239` | `identical` | `DeepEqual` sobre a forma canônica, para o `===` do PHP | **Não.** A comparação estrutural é o serviço, não um detalhe da implementação |
| `testing/t.go:302` | `canonical` | normalizar slice, mapa e ponteiro do que `json.Unmarshal` devolveu | **Não**, mesma razão de `From` |
| `testing/t.go:367` | `isNil` | o `null` do PHP, inclusive interface não-nula carregando valor nil | **Não.** É a única forma de fazer essa pergunta em Go |
| `testing/t.go:408` | `count` | o `count()` do PHP | **Não**, mesma razão da primeira linha |
| `view/runtime.go:74` | `funcName` | detectar que `{{ .Name }}` interpolou o método em vez do resultado, para dizer isso no panic | **Sim, e é a única.** O kyse compila o template e conhece o tipo: a mesma checagem cabe em `aru view:build`, sobre `go/ast` |

Leitura da tabela:

- **Nove dos doze são a mesma pergunta**, feita a um valor cujo tipo estático é
  `any`: é mapa, é slice, qual o tamanho, é nil. Estão em `support/arr` e em
  `testing`, os dois espelhando superfície PHP tipada `mixed`. Genérico não
  responde pergunta sobre valor cujo tipo já foi apagado na chamada — o `any`
  está na assinatura porque `Arr::get` e `assertEquals` recebem `mixed`. Sair de
  `reflect` ali é reprojetar a assinatura, não trocar o mecanismo.
- **Um é cosmético** (`typeName`) e **um é parcial** (`Wrap`), e o parcial custa
  uma segunda grafia.
- **Um pode ir para tempo de build** (`view/runtime.go:74`), e fica como gatilho:
  quando o `kyse` tiver a checagem de tipo do template, a mensagem de panic vira
  erro de compilação e o `reflect` sai do `view`.
- **Nenhum dos doze está em caminho de requisição, e nenhum decide autorização,
  roteamento, despacho ou persistência.** `routing`, `validation`, `security`,
  `queue`, `events` e `database` não importam `reflect` em arquivo não-teste.
  É a tese se sustentando, medida em vez de afirmada.

Em arquivo de teste o `reflect` aparece em cinco lugares, quatro deles
`DeepEqual`. O quinto, `filesystem/disk_test.go:41`, percorre os métodos da
interface `Adapter` para provar que **todos** recebem um `Grant` — reflexão
usada para verificar invariante de arquitetura, que é o uso que este projeto
quer, e num teste, que é onde ele cabe.

---

## Consequências

- Os três diretórios continuam existindo, com `doc.go` e sem símbolo exportado.
  Quem procura acha a razão, e agora acha também o que escrever no lugar.
- `docs/31:390` perde `->when()`/`->unless()`/`Conditionable` da lista de
  rejeitados pela REGRA 9 e ganha a regra estreita da seção 2. `docs/13:125`, que
  já dizia que `->when()` existe porque encadeamento fluente não aceita `if`,
  estava certo e passa a ser o texto de referência.
- `query.Builder` ganha `When` e `Unless`. Três outros tipos fluentes ficam como
  pendência declarada, não como porta aberta: a regra diz quando entram.
- A extensão de terceiro sobre `Collection` acaba, e isso é perda registrada. Um
  pacote como `spatie/laravel-collection-macros` não tem tradução para Go.
- A contagem de `reflect` vira um número verificável, e a linha de comando para
  refazê-la está no `hesape/reflection/doc.go`. Se ela subir, sobe visível.

## Alternativa considerada

**Rejeitar `Conditionable` inteiro, como `docs/31` dizia.** Rejeitada pela
evidência: o `collections` já carregava seis métodos da trait antes desta
discussão, e a regra geral teria de ser aplicada removendo-os. Removê-los
tornaria `Collection().Filter().When().Map()` impossível de escrever sem cortar
a cadeia, no tipo do framework em que encadear é o modo normal de usar. A
proibição geral era mais cara do que a regra estreita e já não estava sendo
obedecida.

**Registrar macro num mapa global de `name -> func(any) any`, com despacho por
string.** Rejeitada: reconstrói `__call` à mão, troca erro de compilação por
panic em runtime, e cai na razão da seção 3 — é exatamente o mecanismo que a
tese do produto recusa.

**Construir `reflection` como pacote de utilitários sobre `reflect`, "para quando
alguém precisar".** Rejeitada: um pacote publicado é um convite. Os doze pontos
medidos acima existem porque alguém precisou de um deles, não porque havia um
pacote; publicar um faria a curva subir, e a única alavanca contra isso é não
existir um lugar óbvio para pôr o décimo terceiro.
