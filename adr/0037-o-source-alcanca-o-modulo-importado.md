# ADR 0037-o-source-alcanca-o-modulo-importado — o Tailwind precisa enxergar o `kyse`

**Status:** Aceita.

## Contexto

A ADR 0025 desligou a detecção automática do Tailwind: `@source` passou a ser a
lista inteira. O `resources/css/app.css` declara duas linhas — as views do
projeto e o gerado —, e nenhuma alcança o `kyse`, que pela ADR 0027 é
**importado**, com o fonte no cache de módulos.

Medido: o `assets/app.css` do `examples`, 193.582 bytes, com **zero** ocorrências
de `.text-destructive`, `.w-44`, `.size-3` e `.rounded-full`. O
`components.Field` escreve `text-destructive` na linha de erro de todo
formulário, então a mensagem de validação de toda tela de auth publicada saía na
cor do corpo; os pontos de cor do seletor de tema são `size-3 rounded-full`, ou
seja, spans de tamanho nenhum. Nada falhou — build verde, `aru view:build`
reportando sucesso, página renderizando. Ficou assim por semanas.

## Decisão

**O `aru view:build` resolve o diretório do módulo (`go list -m -f {{.Dir}}
github.com/arandu-io/kyse`) e o entrega ao Tailwind como `@source`.** A entrada
deixa de ser um caminho e passa a ser montada pela CLI, por stdin: o `@import`
absoluto do `app.css` do projeto mais uma linha `@source "<dir>/**/*.go"`.

Duas alternativas foram pesadas e recusadas:

- **o `kyse` publicar a própria camada de utilitários** — é uma segunda folha de
  estilo, com segunda versão de Tailwind e segunda cópia dos tokens, e as duas
  divergem (REGRA 9);
- **listar as classes em `@source inline(...)` ao lado da biblioteca** — é a
  lista que envelhece em silêncio na primeira classe nova, isto é, este mesmo
  defeito com um passo a mais.

## Consequências

- O conserto é da CLI, não do `app.css`. Projeto que já existe recebe atualizando
  o `aru`; editar folha de estilo à mão é conserto que ninguém aplica.
- O `view:build` passa a exigir `go` no PATH. Quando o módulo está no build list
  mas não foi extraído — todo clone novo —, o comando faz `go mod download` dele
  e olha de novo, para `git clone && aru dev` continuar verdadeiro.
- **Diretório vazio tem um significado só: o projeto não exige o `kyse`.** Todo
  o resto é falha e para o build. O discriminador é a **versão**, não o
  diretório nem o código de saída do `go list`: quem fez `go mod vendor` tem
  versão e nenhum diretório (o fonte está em `vendor/`, e é de lá que se lê);
  quem está sem rede no primeiro build tem versão, nenhum diretório e um motivo;
  quem usa `go.work` tem diretório e nenhuma versão. Ler qualquer um desses como
  "não depende" reabre este mesmo defeito por outra porta — medido: mesmo
  projeto, mesma view, só a presença de `vendor/` decidindo se `.size-3` e
  `.rounded-full` existem na folha.
- Projeto que não depende do `kyse` não ganha `@source` nenhum e compila igual.
- +1,3 KB no `examples` (193.582 → 194.871 bytes): é o que as onze componentes
  desenham.
- `TestTheStylesheetCarriesTheClassesTheMarkupRenders`, em `examples/assets/`,
  fixa as classes no arquivo compilado. A 0025 já tinha ensinado que fonte errada
  não falha — ela não contribui, e não avisa que não contribuiu.
