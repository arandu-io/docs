# 18 — `moitara`: registry de módulos

## O furo que este documento fecha

A tese do produto é que o caminho errado é impossível, não desaconselhado. Módulo
de terceiro é a única porta por onde essa garantia pode vazar — e ela precisa
estar fechada **antes** de existir o primeiro módulo de terceiro, não depois.

Um módulo hostil não consegue fabricar um `security.Grant`: os campos são não
exportados e não há construtor público além de `Authorize`. Mas `SystemGrant`
**é** público, porque trabalho sem request atrás precisa dele — job, tarefa do
scheduler, seeder, comando — e o caminho de login também, onde ainda não existe
sujeito. Esse é o bypass real.

*Corrigido em 07/08/2026:* esta linha dizia "porque jobs e migrations precisam
dele". Migration não precisa: uma migration é DDL, roda pelo `data.Migration` e
não passa por `Repository` nenhum — nenhuma das migrations do esqueleto ou dos
templates do gerador chama `SystemGrant`. O comentário de doc de
`security.SystemGrant` diz para que ele existe: "jobs that run outside a request,
and for the login path, where there is no subject yet".

Fechamento: `aru doctor --strict` reprova `SystemGrant` fora dos diretórios em
que o esqueleto põe trabalho sem request atrás. A regra é
`system-grant-outside-scope`, em `aru/internal/doctor/rules.go`, e o escopo que
ela perdoa é este, por segmento de caminho: `seeders`, `jobs`, `commands`,
`console`, `cmd`, mais os arquivos `main.go` e `routes/console.go`. Fora deles,
só passa a chamada que carrega `//arandu:system-grant <motivo>` na própria linha
ou na de cima — o marcador sem motivo escrito não perdoa nada. Arquivo de teste a
regra nem lê.

*Corrigido em 07/08/2026:* esta linha dizia "fora de job, migration e teste".
Migration não está na lista: `database/migrations` não é escopo de sistema, e um
`SystemGrant` ali é reportado como qualquer outro. A severidade é aviso, não
erro — é o `--strict` que transforma o aviso em reprovação, e por isso a frase só
vale com a flag.

A auditoria da registry roda o doctor em modo strict contra o módulo publicado, e
o selo é o resultado assinado dessa execução.

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
confirmação, roda `doctor` localmente e **imprime** as linhas de fiação, para
serem coladas à mão. Nada é resolvido em runtime — continua tudo em código
commitável.

*Corrigido em 07/08/2026:* esta linha dizia que o `aru add` "regenera o wiring
explícito", e listava isso como trabalho pendente. A decisão está revogada, não
pendente: a ADR 0001, com a nota que a ADR 0023 acrescentou, proíbe o gerador
editar fiação — "não crie `internal/wire/`; fiação gerada ao lado da fiação
escrita à mão seriam duas fiações". O `aru make:module` já é assim: a função
`wiring` em `aru/makemodule.go` imprime as três linhas de `routes/web.go`,
`bootstrap/app.go` e `database/migrations/migrations.go`, e o comando não toca em
nenhum dos três arquivos. Um `aru add` que regenerasse fiação abriria de novo a
porta que a ADR 0001 fechou.

O comando `aru add` ainda não existe — `aru/commands.go` não o registra —, e isso
é coerente com a seção seguinte, que o coloca na Fase 4. O que este documento
descreve aqui é desenho, não entrega.

## Quando construir

Fase 4 ou depois. O que é Fase 2 é apenas o `arandu.mod.toml` e a verificação de
permissões no doctor. Registry sem módulos é vitrine vazia; contrato ausente é
dívida que não se paga depois.
