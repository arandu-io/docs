# ADR 0045 — O que se escreve no lugar do container

**Status:** Aceita. Data: 10/08/2026.
**Complementa a ADR 0001 e a ADR 0002.** Não reabre nenhuma das duas: a recusa
continua de pé, e o que esta ADR acrescenta é a resposta que faltava a ela.

## Contexto

A ADR 0001 tem doze linhas e recusa o container de injeção de dependência. A ADR
0002 tem cinco e recusa as facades. As duas estão certas e nenhuma das duas diz
o que se escreve no lugar.

Isso passou despercebido enquanto quem escrevia o framework era quem tinha
tomado a decisão. Apareceu na auditoria contra o clone do Illuminate: o
`hesape/container` é um dos cinco diretórios que existem só para dizer que o
componente virou decisão em vez de pacote, e o `doc.go` dele tinha oito linhas,
das quais três apontavam para um documento de reorganização. Alguém vindo do
Laravel que abrisse aquele arquivo procurando `singleton` saía dele sem resposta.

O `Illuminate\Container` tem 64 métodos públicos. A maior parte deles some numa
frase — struct com os campos, função `New` que os recebe, montagem explícita em
`bootstrap/app.go`. Quatro não somem, e são exatamente os quatro pelos quais a
pessoa foi lá: `singleton`, `scoped`, `tagged` e o binding contextual
(`when`/`needs`/`give`). Uma recusa que responde os sessenta fáceis e cala nos
quatro difíceis é uma recusa que manda procurar container no `pkg.go.dev` — e há
container em Go de sobra para achar.

## Decisão

**O `hesape/container` não exporta nada, e o `doc.go` dele é a resposta
completa**, com a tabela dos 64 métodos e uma seção por dificuldade real.

O documento canônico é o comentário de doc, não esta ADR: ele é o que o
`pkg.go.dev` publica e o que aparece no editor de quem está escrevendo o código
naquele momento. Esta ADR registra o porquê e as três perdas.

O código citado nas duas peças foi compilado antes de ser citado — `gofmt`,
`go build`, `go vet` e `go test -race`, num pacote de rascunho fora da árvore.
Documento com trecho de Go que não compila é pior que documento faltando, porque
ninguém confere.

### Por que a recusa se sustenta

Um container resolve nome de tipo em valor, em tempo de execução, por reflexão.
Em PHP isso é quase de graça — tudo já é dinâmico e o sistema de tipos não estava
fazendo o trabalho de qualquer forma. Em Go custa em dois lugares, e o segundo é
o que importa.

Custa em runtime, que é a metade menor: `reflect.New` e busca em mapa por
dependência, a cada requisição, para um grafo que não muda depois do boot.

Custa onde importa no momento em que o grafo está errado. `make(Foo)` sem
binding para `Foo` é panic na primeira requisição que passa por ali, encontrado
por um cliente. O mesmo erro numa chamada de construtor é falha de compilação,
encontrada por quem o cometeu, antes do commit. Go foi escolhido por essa troca;
container gasta ela.

Há um terceiro custo que é fácil não ver, e que não estava em nenhuma das duas
ADRs anteriores. Container é estado global mutável, e teste em Go roda em
paralelo. O Laravel se safa com container compartilhado porque um processo PHP
serve uma requisição; no momento em que deixa de servir — Octane — aparece o
`forgetScopedInstances`, e com ele a classe de bug em que uma requisição enxerga
o estado de outra. Nada aqui tem essa classe de bug, porque nada é compartilhado
sem que uma linha de código compartilhe.

### A resposta curta, e os 64 métodos

Struct com os campos, `New` que os recebe, montagem em `bootstrap/app.go`. A
interface é declarada por quem consome, não por quem fornece, então o construtor
aceitar o tipo concreto é verificado no ponto de chamada, em tempo de
compilação, sem registro em lugar nenhum.

Isso resolve, em uma linha cada:

| Illuminate | Go |
| --- | --- |
| `make` / `resolve` / `get` | o campo já está lá: `h.users` |
| `makeWith` | argumentos de função |
| `build` | a função `New` |
| `factory` | um valor de func: `type ReportFactory func(ctx) (*Report, error)` |
| `call` / `wrap` | chame o método; o que não varia é campo |
| `bind` | passe o valor ao construtor |
| `instance` | você já tem o valor; entregue |
| `alias` | quem consome declara a interface de que precisa |
| `extend` | um decorador que implementa a mesma interface |
| `bindIf` / `singletonIf` / `scopedIf` | ou a linha está no `bootstrap/app.go` ou não está; default vai no construtor |
| `bound` / `has` | o compilador; na borda opcional, `!= nil` num campo com nome |
| `resolved` / `isShared` / `isAlias` | não há slot preguiçoso nem lifetime para interrogar |
| `getBindings` | `bootstrap/app.go`, lido de cima a baixo |
| `currentlyResolving` | o stack trace, que já diz melhor |
| `forgetInstance` / `forgetInstances` / `flush` / `forgetScopedInstances` | construa outra `App` |
| `beforeResolving` / `resolving` / `afterResolving` | faça uma vez, no `Build`, à vista |
| `rebinding` / `refresh` | não se troca binding de aplicação rodando; se a necessidade for real, `atomic.Pointer` |
| `bindMethod` / `hasMethodBinding` / `callMethodBinding` | campo de func, ou método no tipo |
| `getInstance` / `setInstance` | é o container global; é o que a ADR 0001 recusou |

O `call` merece nota, porque é onde mora um terço da reflexão do componente
(`BoundMethod.php`): ele injeta as dependências de um método a partir da
assinatura. Em Go, o que o método precisa e não recebe é campo, posto ali pelo
construtor, e o que varia por chamada é argumento.

O ganho maior da lista é a linha do `flush`. Aqueles quatro métodos existem para
desfazer container compartilhado entre testes. Aqui um teste que precisa de
aplicação limpa chama o mesmo `Build` que o binário chama, recebe os próprios
valores, e roda em paralelo com todos os outros, porque não há nada
compartilhado para esquecer.

### `singleton`

`$this->app->singleton(...)` quer dizer: construa uma vez e entregue a mesma a
todos. Em Go isso é o que uma variável é.

```go
func Build(cfg Config, db *sql.DB) App {
	users := NewSQLUsers(db)

	return App{
		Profile: NewProfileController(users, log),
		Billing: NewBillingController(users),
		Cleanup: NewCleanupTask(users),
	}
}
```

Os três seguram o mesmo `*SQLUsers`. Nada declara isso; é o que passar um
ponteiro significa. Não existe flag de lifetime porque não existe um segundo
lifetime: valor construído uma vez e passado é compartilhado, valor construído
duas vezes são dois valores, e as duas coisas estão visíveis nas quatro linhas.

A versão do Laravel carrega um aviso que esta não carrega. Como quem decide é o
container, `bind` e `singleton` produzem código idêntico no ponto de uso, e
errar — serviço com estado registrado com `bind`, e cada injeção recebe um novo;
valor específico da requisição registrado com `singleton`, e ele vaza entre
requisições — é invisível até produção. Aqui o erro tem de ser digitado: duas
chamadas a `NewSQLUsers`, em duas linhas, uma do lado da outra.

O que o singleton traz em Go, e o PHP não precisa pensar, é que valor
compartilhado é tocado por várias goroutines ao mesmo tempo. Compartilhado com
estado mutável quer dizer mutex ou atômico, e `go test -race` é quem cobra. É
obrigação real, e é a razão de a coisa compartilhada quase sempre ser um handle
para algo que já é seguro concorrentemente — `*sql.DB`, `*slog.Logger`,
`*http.Client` — e não um saco de campos.

### `scoped`

É o que não tem resposta de uma linha, e o que decide se a pessoa aceita a
ausência do container.

`scoped` é uma instância por requisição: construída na primeira vez em que é
pedida dentro daquela requisição, compartilhada pelo resto dela, descartada no
fim. Em Go há **duas formas**, e escolher entre elas é a única decisão.

**Forma A — valor no `context.Context` da requisição.**

```go
type grantKey struct{}

func WithGrant(ctx context.Context, g Grant) context.Context {
	return context.WithValue(ctx, grantKey{}, g)
}

func GrantFrom(ctx context.Context) (Grant, error) {
	g, ok := ctx.Value(grantKey{}).(Grant)
	if !ok {
		return Grant{}, ErrNoGrant
	}
	return g, nil
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		g := grantFor(r)
		next.ServeHTTP(w, r.WithContext(WithGrant(r.Context(), g)))
	})
}
```

A chave é tipo não exportado e é struct, não string. Não é estilo: é o que
impede outro pacote de escrever naquele slot e o que torna colisão impossível —
a propriedade que faz isso ser seguro onde um mapa global não seria. O acessor
devolver erro em vez de valor zero é a mesma disciplina: `Grant` ausente tem de
ser tratado, não defaultado (REGRA 14).

**Forma B — campo de um valor que o handler constrói e é dono.**

```go
type RequestScope struct {
	Tx    *sql.Tx
	Start time.Time
	Log   *slog.Logger
}

func (c *OrderController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tx, err := c.db.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, "database unavailable", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	scope := &RequestScope{Tx: tx, Start: time.Now(), Log: c.log.With("path", r.URL.Path)}

	if err := c.place(r.Context(), scope); err != nil {
		http.Error(w, "could not place order", http.StatusInternalServerError)
		return
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, "could not place order", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
```

**Quando cada uma:**

- **Forma B quando o valor tem consumidor nomeado numa assinatura que é sua.** A
  transação, o cache de entidades da requisição, os erros de validação
  acumulados: tudo que os toca é código seu, e passar como argumento é o que faz
  o compilador conferir que chegaram. É o padrão. Custa um parâmetro.

- **Forma A quando o valor tem de atravessar código que não tem razão nenhuma
  para saber dele.** O `Grant` é posto por middleware, atravessa o router, o
  kernel e o pipeline — nada disso é seu e nada disso deve ganhar um parâmetro
  por causa dele — e é lido lá embaixo por um repositório. Prazo e cancelamento
  viajam do mesmo jeito, que é por isso que o `context.Context` existe.

**O custo da Forma A, dito com clareza:** a checagem de tipo se foi. O
`GrantFrom` faz asserção em runtime, e um handler montado sem o middleware
compila e falha na primeira requisição. É exatamente o modo de falha que este
framework recusa em todo o resto — então se gasta de propósito, nos poucos
valores que de fato atravessam o framework, e se gasta atrás de um acessor que
devolve erro em vez de valor zero.

O que **não** se faz é a terceira forma: mapa no nível do pacote indexado por ID
de requisição. Isso é um container com ergonomia pior, e vaza toda vez que um
handler retorna cedo.

Uma propriedade vem de graça e vale ser dita. No Laravel a instância `scoped`
mora num container que sobrevive à requisição, então sob Octane ela precisa ser
esquecida entre requisições, e esquecer de esquecer é vazamento de dado entre
usuários. Aqui um valor de escopo de requisição é variável local na goroutine do
handler: quando o handler retorna, virou lixo. Não há o que esquecer, e o
`forgetScopedInstances` não tem contraparte porque não tem problema a resolver.

### `tagged`

`tag([A::class, B::class], 'reports')` e depois `tagged('reports')` junta um
conjunto de implementações e entrega em grupo. A resposta em Go é um slice, e a
tag é o nome do parâmetro:

```go
type Notifier struct{ channels []Channel }

func NewNotifier(channels ...Channel) *Notifier {
	return &Notifier{channels: channels}
}

// bootstrap/app.go
notifier := NewNotifier(MailChannel{}, SMSChannel{})
```

Três coisas da versão do Laravel não sobrevivem, e duas delas é bom perder.

**A string.** `tagged('reports')` com erro de digitação devolve iterável vazio e
a funcionalidade silenciosamente não faz nada — não há erro, porque tag vazia é
legal. Slice não se digita errado, e elemento de tipo errado é falha de build.

**A ordem.** Passa a ser a ordem da lista de argumentos, que é onde quem procura
vai procurar, em vez de ser a ordem em que os registros aconteceram.

**A preguiça.** Essa é perda real. O `RewindableGenerator.php` existe para que
`tagged` devolva um gerador: os membros só são construídos se iterados, e a
interface `Countable` permite contar sem construir. Isso importa quando uma tag
tem vinte serviços e a requisição usa dois. Se um membro for caro o bastante
para importar, o tipo do elemento vira fábrica — `[]func() Channel` — e a
preguiça fica visível no tipo em vez de escondida no container. Na prática os
membros são handles, e construir todos no boot custa nada e falha no boot em vez
de falhar na vigésima requisição.

`giveTagged` e `giveConfig` são as formas contextuais da mesma coisa: um
argumento de slice e um campo de config.

### `when` / `needs` / `give`

O binding contextual é onde o container é mais convincente em PHP e menos
necessário em Go:

```php
$this->app->when(PhotoController::class)
          ->needs(Filesystem::class)
          ->give(fn () => Storage::disk('local'));
```

Ele existe porque o container resolve **por tipo**, então dois consumidores que
precisam da mesma interface receberiam a mesma implementação, e isso é a saída
de emergência. Go não resolve por tipo. Cada um recebe o que você passa:

```go
photos   := NewPhotoController(LocalDisk{Root: "storage/app"})
invoices := NewInvoiceController(S3Disk{Bucket: "invoices"})
```

É a funcionalidade inteira, em duas linhas, sem builder, sem cadeia fluente e
sem um terceiro conceito para aprender. `when` é o argumento, `needs` é o tipo
do parâmetro, `give` é o valor. Binding contextual aninhado — dar a uma subárvore
inteira outra implementação — é argumento de construtor passado adiante, o que dá
mais o que digitar que o PHP e é também a única forma de descobrir depois qual
subárvore recebeu o quê.

### Os `Attributes/`

Os dezesseis arquivos de `Attributes/` são binding contextual escrito no
parâmetro do construtor. Cada um é campo de struct preenchido no
`bootstrap/app.go`, ou leitura da requisição:

| Atributo | Go |
| --- | --- |
| `#[Auth('web')]` | campo `auth *auth.Service`; o guard é escolhido na fiação |
| `#[Authenticated]` | `GrantFrom(ctx)` — o `Grant` carrega quem chamou |
| `#[CurrentUser]` | `GrantFrom(ctx)`; ele herda de `Authenticated` |
| `#[Cache('redis')]` | campo `cache *cache.Repository` |
| `#[Config('app.name')]` | campo da struct de config tipada: `cfg.App.Name` |
| `#[Context('trace')]` | `context.WithValue` com chave de tipo não exportado |
| `#[Database('reporting')]` | campo `db *data.DB` |
| `#[DB]` | idem; herda de `Database` |
| `#[Give(Foo::class)]` | o argumento que você passa |
| `#[Log('billing')]` | campo `log *slog.Logger`, ou `log.With("component", ...)` |
| `#[RouteParameter('id')]` | `r.PathValue("id")` no handler |
| `#[RequestAttribute('x')]` | valor no `r.Context()`, igual ao `#[Context]` |
| `#[Storage('s3')]` | campo `disk *filesystem.Disk` |
| `#[Tag('reports')]` | campo de slice; ver `tagged` |
| `#[Singleton]` / `#[Scoped]` | ver as seções acima |
| `#[Bind]` / `#[BindWhen]` | uma linha no `bootstrap/app.go`; a condição é um `if` |

Dois merecem separação.

`#[Authenticated]` e `#[CurrentUser]` resolvem o usuário autenticado a partir do
container, o que quer dizer que qualquer objeto, em qualquer lugar, pode
perguntar quem está logado. Aqui quem chamou chega como `security.Grant` no
contexto da requisição e alcança a camada de dados como argumento, porque a
REGRA 14 faz do `Grant` a única fonte de tenant para SQL e a REGRA 17 o exige
também na leitura. Tornar isso alcançável sem passar é precisamente o que este
framework não faz.

`#[Config]` é aquele em que a substituição é melhor, não apenas equivalente:
`#[Config('app.nmae')]` injeta `null` e a falha aparece em outro lugar;
`cfg.App.Nmae` não compila.

## Consequências

- O `hesape/container` fica sem símbolo exportado, permanentemente. Se aparecer
  um, algo deu errado. O `cobertura.sh` já conta o componente na coluna
  "decisão": zero ali não é dívida, é veredito — e agora o veredito está escrito.

- Fica registrado o que se perde, em três itens, porque fingir o contrário
  tornaria o resto menos confiável:

  **Verbosidade.** O `Build` do `bootstrap/app.go` é longo e cresce com a
  aplicação. É o preço de o arquivo ser a resposta a "de onde vem isto", e o
  `aru make:module` imprime as linhas para colar em vez de editá-lo, de modo que
  o arquivo nunca é algo que um gerador escreveu (ADR 0001, ADR 0023).

  **Preguiça.** O container constrói na primeira vez em que pedem; o `Build`
  constrói tudo no boot. Para handle — pool, client, logger — isso é de graça e
  falha no boot em vez de falhar na primeira requisição, o que é melhor. Para
  algo caro e raro, a resposta é campo de func, e o adiamento passa a ser coisa
  do tipo em vez de propriedade do framework.

  **Substituição em runtime.** Um pacote não alcança uma aplicação rodando para
  trocar um binding. Teste troca construindo a aplicação com outro argumento,
  que é a mesma substituição feita no único momento em que ela é segura.

- A ADR 0001 continua com doze linhas e continua sendo a decisão. Esta é a
  resposta que faltava a ela, e o lugar de procurar quando alguém perguntar
  "então como eu faço `scoped`?".

## Alternativa considerada

**Escrever um container pequeno, sem reflexão, com registro tipado por
genéricos.** Rejeitada pela REGRA 9. Ele funcionaria, e seria a segunda maneira
de montar a aplicação, ao lado do `bootstrap/app.go` que já existe. Duas
fiações, uma legível de cima a baixo e outra que se lê saltando entre registros,
é exatamente o custo que a regra existe para recusar — e a versão pequena não
resolve nenhum dos três itens da lista de perdas acima, porque nenhum deles vem
da falta de container.
