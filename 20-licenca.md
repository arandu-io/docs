# 20 — Licença

**Decisão tomada: MIT.** Ver `adr/0008-licenca-mit.md`. O restante deste
documento é a análise que levou à decisão, mantida porque registra o que foi
pesado — e porque a recomendação escrita aqui era outra.

É irreversível na prática: depois do primeiro contribuidor externo, relicenciar
exige o consentimento dele.

## As opções, com a consequência de cada uma

| Licença | Ganha | Perde |
|---|---|---|
| **MIT** | adoção máxima, zero atrito, padrão do ecossistema Go | qualquer um fecha um fork e vende; nenhuma concessão de patente |
| **Apache-2.0** | tudo do MIT **mais** concessão explícita de patente e exigência de `NOTICE` | arquivo a mais no repositório |
| **BSL 1.1** (Dragonfly, HashiCorp) | impede concorrente de oferecer como serviço gerenciado; converte para open em N anos | adoção despenca. Ninguém adota framework cuja licença precisa passar no jurídico |
| **Open core + comercial** | monetização direta | exige CLA desde o primeiro PR externo, e comunidade desconfia |

## Recomendação deste documento, e a decisão do proprietário

Este documento recomendava **Apache-2.0** no framework, nos adapters e na CLI:
a cláusula de patente não custa nada, e a empresa que passa código por jurídico
encontra menos atrito. É o que o Goldr escolheu.

**O proprietário decidiu MIT**, a mesma licença do Laravel, e a decisão vale.
O raciocínio está no `adr/0008`: o alvo declarado do produto é o desenvolvedor
que vem de Laravel, e a licença que ele já conhece é uma barreira a menos. A
concessão de patente do Apache protege quem tem portfólio de patente para
oferecer — o que não é o caso aqui.

O que a escolha de MIT deixa de fora, e que vale saber:

- sem concessão explícita de patente. O risco é teórico para um framework web,
  e real para quem publica implementação de algoritmo patenteado. Não é o caso.
- sem `NOTICE`. Um arquivo a menos para manter.

Não BSL: o produto é um framework, e framework sem adoção não tem o que monetizar.
BSL faz sentido para banco de dados, onde o hyperscaler é a ameaça real. Aqui a
ameaça é a irrelevância, não o fork.

O componente pago, se existir, tem **licença comercial própria** e vive em
repositório separado. Nunca sob a mesma licença do core — misturar os dois é o
erro que faz a comunidade desconfiar do projeto inteiro.

Qual repositório é esse está em aberto: o `porang-pro` era a resposta, e o
`porang` foi dissolvido (ADR 0021). A regra de licença não depende dessa
resposta — o que ela exige é que o pago fique fora do core, sob licença própria.

## A decisão que anda junto: DCO, não CLA

- **DCO** (`Signed-off-by`): contribuidor afirma que pode contribuir. Leve, aceito
  por padrão no mundo Go.
- **CLA**: contribuidor cede direitos, o que permite relicenciar depois.

Recomendo **DCO** — a menos que exista chance real de dual-license no futuro.
Se existir, tem que ser CLA desde o primeiro PR externo; adotar CLA depois é
recontatar todo mundo, um a um.

Como o componente pago fica em repositório separado com licença própria, o core
não precisa ser relicenciável. Por isso: DCO.

## O que fazer, na ordem

1. `LICENSE.md` (MIT) na raiz de todo repositório `arandu-io` — **feito**
2. `CONTRIBUTING.md` exigindo `git commit -s` — **feito**
3. Cabeçalho de licença **não** em cada arquivo — o Go não usa essa prática
4. `THIRD_PARTY.md` no `framework`, com a atribuição de HTMX, Alpine e Tailwind
   embutidos em `framework/view/assets/` — **feito**. E não fica a cargo da
   memória: `TestEveryEmbeddedAssetIsCredited` e
   `TestTheCreditedVersionsAreTheEmbeddedOnes`, em
   `framework/view/third_party_test.go`, reprovam arquivo que aparece em
   `view/assets/` sem entrada no aviso e versão registrada que deixou de bater
   com a que está dentro do arquivo.

   *Corrigido em 07/08/2026:* este item mandava o `THIRD_PARTY.md` para o
   `arandu`. O arquivo está no `framework`, que é onde os assets embutidos moram
   desde a ADR 0021 — procurando o arquivo nos dez repositórios, existe um só, o
   `framework/THIRD_PARTY.md`.

   *Atualizado em 14/08/2026:* os repositórios são catorze, e o arquivo já não é
   um. Existem dois, byte a byte idênticos: `framework/THIRD_PARTY.md` e
   `hesape/view/THIRD_PARTY.md`, porque o `view/` inteiro — aviso, teste de
   crédito e os seis arquivos de `view/assets/` — foi para o `hesape` (ADR 0049)
   sem sair do `framework`. É a duplicação que a própria ADR 0049 chama de duas
   árvores vivas: a correção feita em uma não chega à outra. Qual das duas fica
   segue os assets embutidos, e isso é a transição que a ADR 0049 descreve.

5. `THIRD_PARTY.md` no `arandu` — **não feito, e é furo aberto.** O
   `arandu/assets/app.css` é saída compilada do Tailwind: a primeira linha do
   arquivo é `/*! tailwindcss v4.3.3 | MIT License | https://tailwindcss.com */`,
   e ele entra no binário por `go:embed` em `arandu/assets/assets.go`. O
   esqueleto não tem `THIRD_PARTY.md`, e o teste de crédito do `framework` não
   alcança isso — ele lê o diretório `framework/view/assets/` contra o
   `framework/THIRD_PARTY.md`, e nada olha para `arandu/assets/`. O único teste
   que existe ao lado do arquivo, `TestTheBrowserGetsThisProjectsStylesheet`,
   compara md5 de bytes servidos; não é sobre crédito.

   A consequência é que o furo se propaga: o `aru new` clona o esqueleto inteiro
   e só apaga o `.git`, então todo projeto gerado nasce redistribuindo Tailwind
   sem o aviso de copyright que a licença MIT dele exige — o mesmo raciocínio que
   o `framework/THIRD_PARTY.md` escreve para justificar a própria existência
   ("`go:embed` puts these bytes in the executable, so every user of Arandu
   becomes a redistributor"). Registrado aqui porque é o lugar onde a decisão de
   licença mora, e não resolvido: fechar exige mexer em código, que este
   documento não faz.

Isso deveria ter acontecido antes de qualquer repositório ir a público, e não
aconteceu: `gh repo view` mostra `arandu` e `framework` como `PUBLIC`, com o item
5 em aberto. O prazo não vale mais como prazo — vale como dívida datada.
