# ADR 0026-o-starter-kit-e-um-pacote — `arandu-io/ui`, e a `aru` perde o `make:auth`

**Status:** Aceita.
**Supersede:** a ADR 0012, que decidiu o oposto — "o starter kit é `aru make:auth`.
Não existe `arandu-io/oka`".
**Autoriza:** o décimo primeiro repositório, que a REGRA 3 trava. O pedido é do
proprietário, e é o que essa trava existe para exigir.

## Contexto

O proprietário pediu, desde o começo, que o starter kit fosse **um pacote** —
como o `laravel/ui`: alguém instala, publica o scaffolding no projeto e depois
remove o pacote. E deu a razão, que é a que importa aqui: **no futuro esse
pacote ganha outras dinâmicas.**

A ADR 0012 decidiu o contrário, com um argumento que é verdadeiro e não fecha a
porta que eu achei que fechava:

> Em Go, a única forma de consumir um repositório é `import` — e o que se
> importa não se edita. Editar a tela de login é a primeira coisa que qualquer
> pessoa faz.

A primeira metade continua verdadeira. A segunda não segue dela: **`import` não
é a única forma de consumir um repositório Go.**

```
go run github.com/arandu-io/ui@latest auth
```

roda um módulo publicado e **não acrescenta uma linha ao `go.mod` do projeto** —
verificado, e é uma checagem da `prova-ponta-a-ponta.sh`, não uma afirmação. É o
contrato do `laravel/ui` inteiro, sem a parte de remover:

| Laravel | Aqui |
|---|---|
| `composer require laravel/ui --dev` | — |
| `php artisan ui bootstrap --auth` | `go run github.com/arandu-io/ui@latest auth` |
| `composer remove laravel/ui` | — nada a remover |

O que a ADR 0012 mediu foi o custo de **importar**, e concluiu contra
**publicar**. São coisas diferentes.

## Decisão

**O starter kit é `arandu-io/ui`.** Os templates saem de `aru/internal/gen/` e
passam a viver nele.

**A `aru` perde o `make:auth`.** Não é um atalho que delega, e não é uma segunda
implementação: sai. Duas maneiras de instalar o mesmo scaffolding divergem — é o
que aconteceu com os hooks de git, e cinco repositórios rodaram a cópia errada
por semanas (REGRA 9).

É também o desenho do Laravel, que é o alvo da REGRA 10: **o `artisan` não tem
`make:auth` desde a versão 6.** Quem publica é o `laravel/ui`.

### Não há preset

`php artisan ui` recebe um — `bootstrap`, `vue`, `react` — porque o Laravel
oferece três. O Arandu tem uma pilha: kyse, HTMX, Alpine e Tailwind v4
(ADR 0022). Argumento de preset com um valor legal é uma dimensão que não
existe, então o que varia é o verbo: `auth` hoje, e o que mais ganhar lugar.

## Consequências

- **São onze repositórios.** A REGRA 3 muda de dez para onze, e este documento é
  a autorização que ela pede.
- **O kit evolui sem release da `aru`.** É a razão pela qual o proprietário
  pediu pacote, e é o que a ADR 0012 tirava: hoje uma tela nova exigia publicar
  a CLI inteira.
- **O módulo não tem dependência**, e a CI reprova a primeira. Ele é baixado por
  quem só quer publicar nove views; cada `require` é download que essa pessoa
  paga. Por isso o renderizador de template é uma **cópia** do que está em
  `aru/internal/gen` — importar traria a CLI de volta, que é exatamente o que
  este pacote existe para tirar do caminho. As duas cópias são pequenas o
  bastante para ler lado a lado, e o que divergiria é pego onde importa: as
  golden files comparam a saída publicada byte a byte.
- **O compilador do kyse é `internal` da `aru`**, então este módulo não pode
  importá-lo — e isso está certo: um starter kit que dependesse da CLI voltaria
  a ser parte dela. O teste de "as views compilam" mora na
  `prova-ponta-a-ponta.sh`, onde é mais forte do que seria aqui: publica o kit
  num projeto gerado, roda o `aru view:build` de verdade, o `go build` de
  verdade, e renderiza cada página por HTTP.
- **Quem já tem projeto** com telas publicadas pelo `make:auth` não precisa fazer
  nada: os arquivos são dele desde que chegaram. O comando é que mudou.

## O que não muda

Os treze arquivos, os caminhos, o conteúdo e as cinco substituições sem
`--force` — o layout, `page.go`, `home`, `welcome` e o `HomeController`, que são
uma unidade porque em kyse a página renderiza com o tipo do layout. As golden
files vieram junto e não mudaram um byte.

## Alternativa considerada

**`aru make:auth` ficar como atalho que chama o pacote.** Ganharia descoberta
pelo `aru help`. Rejeitada: continuaria sendo duas maneiras de fazer a mesma
coisa na superfície de quem usa, e o Laravel — que é o modelo — escolheu a
mesma coisa quando tirou o `make:auth` do artisan.
