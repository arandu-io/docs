# 18 — `moitara`: registry de módulos

## O furo que este documento fecha

A tese do produto é que o caminho errado é impossível, não desaconselhado. Módulo
de terceiro é a única porta por onde essa garantia pode vazar — e ela precisa
estar fechada **antes** de existir o primeiro módulo de terceiro, não depois.

Um módulo hostil não consegue fabricar um `security.Grant`: os campos são não
exportados e não há construtor público além de `Authorize`. Mas `SystemGrant` **é**
público, porque jobs e migrations precisam dele. Esse é o bypass real.

Fechamento: `aru doctor --strict` reprova `SystemGrant` fora de job, migration e
teste. A auditoria da registry roda o doctor em modo strict contra o módulo
publicado, e o selo é o resultado assinado dessa execução.

## `moitara` não hospeda pacote

Go não tem registry central: o namespace é o próprio repositório
(`github.com/fulano/modulo`), resolvido pelo `proxy.golang.org`. Tentar ser um npm
seria construir infraestrutura que a linguagem já tem.

`moitara` é **índice curado mais selo de auditoria**. Três níveis:

| Nível | Significado |
|---|---|
| `oficial` | publicado sob `arandu-io` |
| `auditado` | passou `doctor --strict` e revisão humana, com selo assinado e versão fixada |
| `comunidade` | apenas listado. Sem selo, sem promessa |

## `arandu.mod.toml` — a decisão que precisa sair agora

Todo módulo declara metadado, e isso tem que existir desde a **Fase 2**: módulo
gerado sem esse arquivo não tem como ganhar o metadado depois sem quebrar quem já
instalou.

**As chaves são em inglês, como todo identificador (REGRA 6).** O arquivo é lido
por `aru/internal/manifest`, que recusa chave desconhecida em vez de ignorá-la —
permissão escrita errado é permissão que parece declarada e não está. Este bloco
é o que o parser aceita, e nada além dele:

```toml
# profiles: conventional | performance, ver 12-perfis.md
name = "fulano/crm"
framework = ">= 1.2"
profiles = ["conventional"]

[permissions]
# network: o módulo faz chamada HTTP para fora
network = true
filesystem = false
exec = false
migrations = true
```

O leitor é um parser de TOML **mínimo**, escrito à mão porque o core não toma
dependência de terceiro (REGRA 2): ele entende seção, `chave = valor`, string
entre aspas, lista e booleano. Comentário só em **linha inteira**, começando com
`#`. Comentário no fim de uma linha de valor não é reconhecido e vira parte do
valor — `network = true  # faz chamada` para com
`network must be true or false`. Por isso o bloco acima comenta acima, nunca ao
lado.

Quatro permissões, e só essas quatro: `network`, `filesystem`, `exec` e
`migrations`. Qualquer outra é erro de leitura, com arquivo e linha. A seção se
chama `[permissions]`; `[permissoes]` é seção desconhecida e para o parser na
primeira linha, com `unknown section [permissoes]`.

O `doctor` confere a declaração contra os imports reais: módulo que declara
`network = false` e importa `net/http` como cliente é reprovado. Permissão
declarada e não usada é aviso; usada e não declarada é erro. É a regra
`permission-not-declared` do `aru doctor`.

**Quem escreve o arquivo hoje.** O `aru make:module` deixou de emiti-lo na
ADR 0023, junto com o `modules/<nome>/` que ele descrevia: o gerador escreve na
árvore do Laravel, e um projeto Arandu não é um módulo publicável. Enquanto a
registry não existe, o `arandu.mod.toml` é escrito à mão por quem for publicar —
e o leitor ainda procura os módulos em `modules/`, que é dívida registrada na
ADR 0023, não decisão.

## Instalação

```
aru add github.com/fulano/crm
```

Resolve pelo proxy do Go, checa o índice, exibe nível e permissões para
confirmação, roda `doctor` localmente e regenera o wiring explícito. Nada é
resolvido em runtime — continua tudo em código commitável.

## Quando construir

Fase 4 ou depois. O que é Fase 2 é apenas o `arandu.mod.toml` e a verificação de
permissões no doctor. Registry sem módulos é vitrine vazia; contrato ausente é
dívida que não se paga depois.
