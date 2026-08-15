# ADR 0047 — Nome de pacote que colide com a biblioteca padrão

**Status:** Aceita e **executada em 14/08/2026**. Data da decisão: 10/08/2026.

`framework/httpx` é `framework/http`, e o renomeio foi atômico nos oito
repositórios que o citavam — `framework`, `arandu`, `examples`, `aru`, `ui`,
`kv`, `queue` e `mcp`. Nenhum identificador exportado mudou: só o caminho.
Zero ocorrências de `httpx` restaram no código, e nenhuma casca ficou para trás,
que é o que esta ADR exige ao dizer que o renomeio é atômico.

Duas armadilhas que só apareceram na execução, registradas para o próximo
renomeio de pacote:

- **No `aru` e no `ui`, a maioria das ocorrências vive dentro de string de
  gerador** — os templates e os `.golden`. O `go build` fica verde enquanto elas
  envelhecem, e quem descobre é o projeto que o gerador produz.
- **O apelido se decide pelo arquivo emitido, não pelo arquivo que emite.** O
  gerador importa `net/http`; o controller que ele escreve pode não importar. Ler
  os imports do gerador põe apelido em arquivo que não tem contra o que apelidar.

`arandutest` continua fora do espelho, pela exceção que esta ADR já declara.

## Contexto

A ADR 0044 decidiu que o nome de função é o do Illuminate, sem exceção, e que a
única alteração permitida é mecânica e tem de estar dita no comentário de doc.
Ela resolveu o identificador e **não** resolveu o nome do pacote.

O espelho do clone manda o diretório se chamar como o componente: `Illuminate\Http`
vira `hesape/http`, `Illuminate\Testing` vira `hesape/testing`. Dois desses
colidem com pacote da biblioteca padrão do Go — `net/http` e `testing` — e a
colisão não é de arquivo, é de **identificador dentro do arquivo de quem importa**:

```go
import (
    "net/http"
    "github.com/arandu-io/hesape/http"   // http redeclarado
)
```

Quem importa os dois é obrigado a apelidar um deles, e um handler que não usa
`net/http` é raro o bastante para não contar.

A medição de 10/08/2026 encontrou o resultado de não ter decidido isso: **quatro
pacotes existem em duplicata**, cada par com o mesmo conteúdo prometido e
implementações diferentes.

| espelho | o que existe hoje | linhas | exportados | importadores |
|---|---|---:|---:|---:|
| `http` | `hesape/http` | 316 | 0 | 0 |
| | `hesape/httpx` | 1.909 | 78 | 13 |
| `testing` | `hesape/testing` | 514 | 0 | 0 |
| | `hesape/arandutest` | 1.102 | 46 | 2 |

Os dois que seguem o espelho são esboço: têm helper privado escrito e **nada
exportado**, porque foram criados pelo espelhamento e nunca preenchidos. Os dois
que fogem dele são os que a aplicação usa.

Isto é a REGRA 9 acontecendo por omissão: duas maneiras de fazer a mesma coisa,
criadas não por decisão mas por ninguém ter decidido.

## Decisão

**O nome do diretório é o do componente Illuminate. Quando ele colide com um
pacote da biblioteca padrão, o custo do apelido é de quem importa, não do
espelho.**

    hesape/http        Illuminate\Http       quem importar os dois apelida
    hesape/testing     Illuminate\Testing    idem

Um apelido no import é uma linha no arquivo de quem chama. Um nome inventado é
uma pergunta em toda busca — quem procura `Illuminate\Http` no nosso código
precisa saber, antes de procurar, que aqui ele se chama outra coisa.

Go apelida import desde sempre, a biblioteca padrão faz isso com
`math/rand/v2`, e o `x` de `httpx` não é uma convenção: é a marca de que o nome
certo estava ocupado.

### O que acontece com os quatro

`httpx` **passa a ser** `http`, e o esboço em `hesape/http` é absorvido — o
`dataGet` que ele tem responde ao `data_get()` que o `Request::input` do Laravel
usa, e é a única coisa aproveitável lá.

`hesape/testing` é absorvido por **`arandutest`**, e este é a exceção declarada
abaixo.

O renomeio é atômico com os treze importadores. Um pacote renomeado com o
importador para trás é o modo de falha que já custou uma noite aqui.

### A exceção, e por que ela é uma só

`arandutest` **fica com o nome fora do espelho**, e a razão não é a colisão: é
que o pacote é importado por arquivo `_test.go`, onde `testing` da stdlib está
presente em praticamente toda linha. Apelidar ali seria apelidar em todo teste
de toda aplicação, o que inverte a conta — o custo deixa de ser de quem importa
os dois e passa a ser de todo mundo.

O `doc.go` dele diz a que componente responde, que é o que mantém a busca
funcionando.

## Consequências

- Some a pergunta "é `http` ou `httpx`?", que hoje se responde olhando qual dos
  dois tem código.
- Quem escreve um handler passa a apelidar: `import stdhttp "net/http"` ou
  `import hhttp "github.com/arandu-io/hesape/http"`. A convenção do projeto é
  apelidar **o nosso**, deixando `http` significando `net/http` como significa
  em todo código Go — o `hesape/http` é o pacote que a aplicação toca por
  último, através do `Context`.
- O `aru` e os esqueletos publicados compilam contra `httpx` e mudam junto.
- `redis` está em disputa parecida com `kv`, e **não** é decidido aqui: o `kv`
  é um repositório publicado, não um diretório do `hesape`, e desmontá-lo é
  decisão de escopo (REGRA 3) e não de nome.

## O que esta ADR não decide

Nome de ferramenta com nome próprio continua sendo guarani e passando pelos
cinco filtros da REGRA 4 — `kyse`, `aru`, `hesape` (ADR 0044). Nada aqui os
toca: um componente do Illuminate não é uma ferramenta com nome próprio, é uma
peça da biblioteca, e o nome dela é o que já está escrito na documentação que a
pessoa leu antes de chegar aqui.
