# ADR 0027-kyse-e-um-repositorio — a biblioteca de componentes é importada, não copiada

**Status:** Aceita.

## Contexto

A REGRA 3 trava o número de repositórios em onze e exige autorização explícita
para o décimo segundo. O proprietário autorizou este em 08/08/2026, com a frase
*"preciso ser um repositório dedicado a isso, como se fosse uma UI dedicada ao
Arandu, forkada e adaptada dele"*.

Antes disso a mesma coisa foi construída duas vezes, das duas formas erradas:

1. **Dentro do esqueleto.** Todo projeto criado por `aru new` nascia com onze
   componentes que a maioria não usa. Contraria "pouco código".
2. **Como publicador no `arandu-io/ui`.** `ui add button` copiava o
   `button.kyse.go` para dentro do projeto — o modelo do shadcn. Foi recusado:
   quem quer subir uma tela rápido não quer onze arquivos novos na árvore, e o
   `ui` tem zero dependências de propósito, o que impede o componente compilado
   de morar lá (ele importa `framework/view`).

## Decisão

**`github.com/arandu-io/kyse` é um módulo importável.**

```sh
go get github.com/arandu-io/kyse/components
```

Nada é copiado, nada entra em `resources/views/`, e o `aru view:build` do projeto
não tem o que compilar a mais. O `.go` de cada componente é gerado por
`components/build.sh` e **commitado** — módulo cujo gerado não está no git é
módulo que o `go get` não serve.

O `arandu-io/ui` volta a ser só o kit de autenticação, e o comando `add` foi
removido: com componente importável, ele era a segunda forma de desenhar um botão
que a REGRA 9 recusa. Publicar e importar não coexistem.

### O nome

`kyse` passou pelos cinco filtros da REGRA 4: sem gíria ou palavrão em pt-BR;
não é divindade, ofício religioso, etnônimo nem termo de disputa política; em
inglês lê como *the UI kit for kyse*, que é a relação `shadcn/ui`; `kyse` já é o
nome do motor de view em uso desde a ADR 0020; e `arandu-io/kyse` estava livre
no GitHub.

### O que é fork e o que não é

Do [shadcn-htmx](https://github.com/productdevbook/shadcn-htmx) (MIT) vêm o
conjunto de componentes, a API de atributo (`data-variant`, `data-size`) e os
padrões de HTMX — o toast que o endpoint devolve e o `hx-swap="beforeend"`
anexa, por exemplo.

A marcação **não** vem. Lá um botão é uma string de 443 caracteres de utilitária
Tailwind; aqui é `class="btn"`, trinta caracteres, com a regra na folha
vendorizada do Basecoat. Quinze vezes menos, no critério que o proprietário
fixou: o código é escrito por IA, e token é o custo.

## Consequências

- **Um componente é função Go exportada e tipada.** `components.Buton` é
  `undefined` e `Labl` é `unknown field`, os dois em tempo de compilação, na
  linha do `.kyse.go`. Era o furo da versão anterior, resolvida por string em
  runtime: nome errado chegava em produção como espaço em branco.
- **Quem quiser editar um componente envolve-o**, não o edita. É o custo de
  importar em vez de copiar, e é o que o proprietário escolheu sabendo dele.
- **A folha de estilo não vem junto.** Ela mora no esqueleto, em
  `resources/css/basecoat/`, porque o projeto possui o design system. Um projeto
  que apagou aquela pasta recebe componente sem estilo.
- **Onze de oitenta e dois.** O que está adaptado é o que o blog e o painel usam.
  O README diz o número em vez de prometer o catálogo.
- **Doze repositórios.** O `plans/checklist.sh` conta, e passa a contar doze.
