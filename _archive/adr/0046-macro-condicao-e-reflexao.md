# ADR 0046 — Macro, condição e reflexão: o que some, o que fica com outro nome, e o que já entrou sem pedir

**Status:** Aceita. Data: 10/08/2026.
**Supersede parcialmente:** `docs/31-reorganizacao-hesape.md:390`, a linha que
lista `->when()`/`->unless()`/`Conditionable` "em qualquer forma" entre o que a
REGRA 9 rejeita, e `docs/31:116`, a linha da tabela que dá "`if` existe em
posição de statement em Go" como motivo da recusa. A rejeição geral cai nas
duas; entra uma regra estreita, na seção 2.

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
`hesape/collections/enumerates.go:635-682` já carrega **seis** desses métodos —
`When`, `Unless`, `WhenEmpty`, `WhenNotEmpty`, `UnlessEmpty`, `UnlessNotEmpty` —
escritos, testados e com o comentário de doc nomeando a trait. Registro que
afirma o contrário do código é registro que ninguém vai consultar duas vezes.

### O argumento, contra o `database/query/builder.go` que acabou de ser escrito

O builder é fluente de verdade, e dá para medir: em `builder.go`, **67 dos 84
métodos exportados devolvem `*Builder`** — `Where`, `OrderBy`, `Limit`, `Join` e
o resto mutam e devolvem o receptor. Encadear não é um estilo possível ali, é o
modo normal de usar o tipo. Com filtro condicional, o `if` custa isto:

```go
q = q.From("posts")
if published {
	q = q.Where("published", true)
}
q = q.OrderBy("id")
```

contra isto:

```go
q.From("posts").
	When(published, func(q *query.Builder) { q.Where("published", true) }, nil).
	OrderBy("id")
```

As duas formas foram compiladas contra o builder como ele está, num pacote de
rascunho fora da árvore, antes de serem citadas. A grafia Laravel da primeira
linha é `db.table('posts')`, que é `Connection::table` e ainda não está escrita;
`From` é `Builder::from`, e é o ponto de entrada que existe hoje.

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

Candidatos que a regra alcança e que ficam por escrever, cada um **no seu
pacote, por quem o escreve** — esta ADR não os toca:

| Candidato | Estado hoje | Usa `Conditionable` no PHP |
|---|---|---|
| `mail.Pending` (`mail/mail.go:222`) | escrito e fluente (`CC`, `BCC` devolvem `*Pending`) | `PendingMail` sim |
| `routing.Route` (`routing/name.go:13`) | escrito e fluente (`Name` devolve `*Route`) | `Route` sim |
| `http/client.PendingRequest` | **não escrito**: `http/client/` tem só `doc.go` e subpastas | `PendingRequest` sim |

A ADR anterior a esta revisão dizia "`httpx` client `PendingRequest`". Errado:
`httpx` é o contexto de requisição, e o cliente HTTP é `http/client`. Corrigido
aqui para que a pendência aponte para o pacote certo.

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
assim mesmo.** Então foi contada — e a contagem se mexeu três vezes durante a
redação, o que é informação e não ruído. Leia a seção sabendo disso: **a tabela
é a auditoria linha a linha de oito arquivos; o total do repositório está no fim
da seção e é maior.**

Auditados linha a linha, **oito arquivos não-teste, vinte e um pontos de
chamada**:

| Arquivo:linha | Símbolo | Para quê | Dá para trocar por genérico ou interface? |
|---|---|---|---|
| `support/arr/arr.go:57` | `Accessible` | `Arr::accessible`: o valor pode ser indexado — mapa, slice ou array | **Não**, enquanto a assinatura for `any`. Um type switch sobre `map[string]any` e `[]any` responderia `false` para `map[string]int`, que é resposta errada, não resposta mais estrita |
| `support/arr/arr.go:169` | `typeName` | o `gettype()` de dentro das mensagens de erro de tipo | **Cosmético.** `fmt.Sprintf("%T", v)` é o mesmo `reflect` com outra grafia; troca o import, não o mecanismo |
| `support/arr/arr.go:375-384` | `From` | desreferenciar ponteiro e escolher o ramo antes do round-trip por JSON | **Não.** O `reflect` aqui só escolhe o ramo; quem faz o trabalho é a reflexão do `encoding/json`, e ela sai junto com a assinatura `any`, não antes |
| `support/arr/arr.go:500` | `IsAssoc` | `Arr::isAssoc` | **Não**, mesma razão da primeira linha |
| `support/arr/arr.go:509` | `IsList` | `Arr::isList` | **Não**, mesma razão |
| `support/arr/list.go:297` | `Wrap` | `Arr::wrap`: percorrer slice de qualquer tipo para produzir `[]any` | **Parcial.** `Wrap[T any](s []T) []any` cobre o caso tipado e deixa o caso `any` ainda em `reflect` — duas grafias de `Arr::wrap`, que a REGRA 9 recusa |
| `support/optional.go:33` | `Optional.Get` | `Optional::__get`: campo de struct por nome, chave de mapa, atravessando ponteiro | **Não.** Uma interface `interface{ Get(string) any }` só responderia por tipo que optou por ela, e `optional()` existe exatamente para o valor que já se tem na mão |
| `support/interactswithdata.go:311` | `matchesEnumCase` | `BackedEnum::tryFrom`: comparar o valor cru com o backing do case, sem que `type Status string` escape do type switch | **Sim, por interface, e com preço.** `interface{ Value() any }` é o `BackedEnum` do PHP escrito em Go, mas exige o método em todo enum da aplicação; hoje um `type Status string` sem método nenhum funciona. É o único da tabela em que a troca é possível e a decisão é de custo, não de mecanismo |
| `support/testing/fakes/fakes.go:46` | `class` | `get_class()`: o tipo sob o qual o fake arquiva o valor gravado | **Não.** O fake guarda valores heterogêneos numa lista; o tipo foi apagado na inserção, e genérico não responde depois disso |
| `support/testing/fakes/fakes.go:69` | `classToken` | `SomeJob::class` normalizado | **Não**, mesma razão |
| `support/testing/fakes/fakes.go:91` | `tokenName` | o nome que a mensagem de asserção imprime | **Não**, mesma razão |
| `support/testing/fakes/fakes.go:122` | `instanceOf` | `$valor instanceof $tipo`, com token de interface satisfeito também pelo ponteiro | **Não.** É a pergunta de subtipo em runtime, que é o serviço |
| `support/testing/fakes/fakes.go:153` | `restore` | `unserialize(serialize())`: o round-trip que o job faz pela fila | **Não.** `restore[T any](v T) T` resolveria, mas o registro é `[]any` — a erosão do tipo é na lista, não aqui |
| `testing/t.go:200` | `assertIsArray` | `PHPUnit::assertIsArray` sobre valor decodificado | **Não**, mesma razão da primeira linha |
| `testing/t.go:239` | `identical` | `DeepEqual` sobre a forma canônica, para o `===` do PHP | **Não.** A comparação estrutural é o serviço, não um detalhe da implementação |
| `testing/t.go:302` | `canonical` | normalizar slice, mapa e ponteiro do que `json.Unmarshal` devolveu | **Não**, mesma razão de `From` |
| `testing/t.go:367` | `isNil` | o `null` do PHP, inclusive interface não-nula carregando valor nil | **Não.** É a única forma de fazer essa pergunta em Go |
| `testing/t.go:408` | `count` | o `count()` do PHP | **Não**, mesma razão da primeira linha |
| `view/runtime.go:74` | `funcName` | detectar que `{{ .Name }}` interpolou o método em vez do resultado, para dizer isso no panic | **Sim, e é a única que sai de vez.** O kyse compila o template e conhece o tipo: a mesma checagem cabe em `aru view:build`, sobre `go/ast` |
| `log/context/repository.go:736` | `stackContains` | `in_array` com `$strict`: `Context::stackContains` | **Não.** O `mixed` está na assinatura porque a pilha de contexto guarda o que lhe derem |
| `log/context/repository.go:871` | `looseEqual` | o `==` do PHP: `in_array` sem `$strict` | **Não**, mesma razão de `identical` |

Leitura da tabela:

- **Dez dos vinte e um estão em pacote que só existe para servir teste** —
  `testing/t.go` e `support/testing/fakes/fakes.go`. Não são arquivos `_test.go`,
  então a contagem os inclui de propósito; mas nada que atenda uma requisição
  chega neles.
- **Nove dos onze restantes são a mesma pergunta**, feita a um valor cujo tipo
  estático é `any`: é mapa, é slice, qual o tamanho, é nil, são iguais. Estão em
  `support/arr`, `support/optional` e `log/context`, os três espelhando
  superfície PHP tipada `mixed`. Genérico não responde pergunta sobre valor cujo
  tipo já foi apagado na chamada — o `any` está na assinatura porque `Arr::get`,
  `Optional::__get` e `Context::stackContains` recebem `mixed`. Sair de `reflect`
  ali é reprojetar a assinatura, não trocar o mecanismo.
- **Um é cosmético** (`typeName`) e **um é parcial** (`Wrap`), e o parcial custa
  uma segunda grafia.
- **Um pode ir para tempo de build** (`view/runtime.go:74`), e fica como gatilho:
  quando o `kyse` tiver a checagem de tipo do template, a mensagem de panic vira
  erro de compilação e o `reflect` sai do `view`.
- **Um tem resposta de interface e a resposta foi recusada por preço**
  (`matchesEnumCase`): `BackedEnum` em Go é um método a mais em todo enum de toda
  aplicação, para tirar um `reflect.Kind` de uma função de dez linhas. Fica
  registrado como troca conhecida, não como pendência.
- **Nenhum dos vinte e um decide autorização, roteamento, despacho ou
  persistência.** Medido, não afirmado: `auth`, `routing`, `validation`, `queue`,
  `events`, `database`, `session`, `encryption`, `hashing`, `http` e `cache`
  importam `reflect` em **zero** arquivos não-teste.

A versão anterior desta seção dizia que nenhum dos pontos estava "em caminho de
requisição". Era forte demais e sai: `support/arr`, `support/optional` e
`log/context` são alcançáveis de dentro de uma requisição. O que se sustenta é a
frase estreita acima — nenhum deles **decide** nada.

Em arquivo de teste o `reflect` aparece em cinco arquivos e dezenove pontos,
dezoito deles `DeepEqual`. O décimo nono, `filesystem/disk_test.go:41`, percorre
os métodos da interface `Adapter` para provar que **nenhum** recebe um `Grant` —
reflexão usada para verificar invariante de arquitetura, que é o uso que este
projeto quer, e num teste, que é onde ele cabe.

### Segunda medição, no fim do mesmo dia — e o que fica para auditar

Entre a auditoria acima e o fechamento desta ADR, cinco arquivos novos passaram
a importar `reflect`. Foram **contados, não auditados linha a linha**, e a
distinção está aqui de propósito: número medido por comando é uma coisa,
veredito por ponto de chamada é outra, e misturar os dois é como uma tabela
envelhece sem avisar.

| Arquivo | Ocorrências de `reflect.` | O que é, em uma linha |
|---|---|---|
| `collections/arr/access.go` | 41 | `Arr` completo sobre mapa e slice de qualquer tipo — chave, `Get`/`Set`/`Forget`, `IsAssoc`, `IsList` |
| `collections/arr/typed.go` | 4 | as variantes tipadas do mesmo `Arr` |
| `support/testing/fakes/queuefake.go` | 4 | `reflect.TypeFor` como token de classe, e `DeepEqual` sobre o round-trip |
| `support/testing/fakes/busfake.go` | 2 | o mesmo, no `BusFake` |
| `jsonschema/deserializer.go` | 1 | `DeepEqual` comparando ramos irmãos do schema |

**Total do repositório na última medição: 13 arquivos não-teste, 98 ocorrências
de `reflect.` fora de comentário.** Contra 8 e 21 algumas horas antes, e 4 e 12
de manhã.

Dois recados, e nenhum deles é alarme:

1. **A subida é quase toda um pacote só.** `collections/arr` responde por 45 das
   98 ocorrências, e é a reescrita do `Arr` — exatamente a superfície que a
   tabela já classificou como irredutível, pelo mesmo motivo (`Arr::get` recebe
   `mixed`). Não é reflexão nova em lugar novo; é mais do mesmo lugar.
2. **`collections/arr` é o que a próxima auditoria olha primeiro**, e a pergunta
   a fazer lá é uma só: `collections/arr` e `support/arr` são dois `Arr`? Se
   forem, é a REGRA 9 e não esta ADR que decide, e a decisão é de quem escreve os
   dois pacotes — não desta ADR, que não os toca.

O que **não** mudou com a subida: `auth`, `routing`, `validation`, `queue`,
`events`, `database`, `session`, `encryption`, `hashing`, `http` e `cache`
continuam em **zero** arquivos não-teste com `reflect`. Foi remedido depois da
última contagem, não herdado da primeira. É a frase que importa, e é a que se
mantém.

---

## Consequências

- Os três diretórios continuam existindo, com `doc.go` e sem símbolo exportado.
  Quem procura acha a razão, e agora acha também o que escrever no lugar.
- `docs/31:390` perde `->when()`/`->unless()`/`Conditionable` da lista de
  rejeitados pela REGRA 9, e `docs/31:116` perde a linha de tabela que dava "`if`
  existe em posição de statement" como motivo; as duas ganham a regra estreita da
  seção 2. `docs/13:125`, que já dizia que `->when()` existe porque encadeamento
  fluente não aceita `if`, estava certo e passa a ser o texto de referência.
  `docs/13:28`, que diz `Conditionable` → "`if`" → "não existe", cai junto.
- `query.Builder` ganha `When` e `Unless`. Três outros tipos fluentes ficam como
  pendência declarada, não como porta aberta: a regra diz quando entram.
- A extensão de terceiro sobre `Collection` acaba, e isso é perda registrada. Um
  pacote como `spatie/laravel-collection-macros` não tem tradução para Go.
- A contagem de `reflect` vira um número verificável, e a linha de comando para
  refazê-la está no `hesape/reflection/doc.go`. Se ela subir, sobe visível — e
  **subiu três vezes antes desta ADR fechar**: 4 arquivos e 12 pontos de manhã,
  8 e 21 na auditoria linha a linha, 13 arquivos e 98 ocorrências na medição
  final, com `collections/arr` respondendo por quase metade. O número na tabela
  envelhece em horas; o comando não. É por isso que ele está no `doc.go`, e é por
  isso que quem citar qualquer número desta ADR deve rodá-lo antes.
- A frase que a subida **não** moveu é a única que decide alguma coisa: os onze
  pacotes que autorizam, roteiam, despacham e persistem seguem em zero. Ela foi
  remedida depois da última contagem, não herdada da primeira.

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
