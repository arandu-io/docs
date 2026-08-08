# ADR 0028-o-chrome-mora-no-framework — a view gerada fica ao lado do fonte

**Status:** Aceita.

## Contexto

O gerado do kyse caía achatado: `auth/login.kyse.go` compilava para
`auth_login.go`, um diretório acima. Quinze arquivos no topo de
`resources/views/`, ao lado das pastas que os originaram.

O motivo registrado no código era:

> *"não é ao lado quando o fonte é aninhado, e isso é uma restrição da
> linguagem, não uma preferência. **Go tem um pacote por diretório.** Com o
> gerado ao lado de um fonte aninhado, `resources/views/layouts/` seria um
> pacote diferente de `resources/views/`, e a struct que o layout declara não
> seria visível para a página que o estende."*

A primeira frase é verdade e a conclusão não decorre dela: pacote aninhado
importa o que precisa, como qualquer outro. Foi medido com um módulo que compila
e roda.

O que fazia o achatamento **parecer** necessário era o chrome — `Page` e
`Layout` — morar no pacote de views do projeto.

## Decisão

**`view.Page` e `view.Layout` mudam para `framework/view`**, e o gerado passa a
ficar ao lado do fonte.

- Todo projeto Arandu declarava as mesmas noventa linhas de título, marca, token
  e quatro links. Ninguém as escreveu e ninguém as lia.
- Todo arquivo gerado já importa `framework/view`, então a página aninhada
  enxerga o contrato sem importar mais nada.
- `view.Page` ganha `PageDescription` e `CanonicalURL`, para uma tela dizer o que
  é sem o layout adivinhar.

O layout do projeto passa a **não declarar tipo nenhum**. `RenderType` devolve
`view.Layout` quando não há interface — e é assim que o defeito que custou um
ciclo não volta: sem `@go` no layout, não existe tipo ali para o compilador de
view escolher errado.

Um layout que desenha mais declara a própria interface. A área de admin do
`examples` faz isso, e é o caso que o desenho deixa aberto.

## Consequências

- **Um blank import por diretório de views** no `bootstrap`, ao lado do que já
  existia. É a mesma registração em que todo o desenho se apoia.
- `resources/views/` mostra a árvore que a pessoa escreveu. O
  `.vscode/settings.json` esconde o gerado da listagem e da busca, com cláusula
  `when` — arquivo Go escrito à mão no mesmo diretório continua visível.
- **O gerado continua no git.** Módulo cujo gerado falta é módulo que não compila
  num clone novo, e `pkg.go.dev` fica vazio.
- `make:module` emite um pacote por recurso, e o controller importa esse pacote.
- Mudança compatível no framework: `apidiff` reporta só `Page: added` e
  `Layout: added`.
