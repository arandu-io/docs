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

```toml
nome     = "fulano/crm"
framework = ">= 1.2"
perfis   = ["convencional"]         # ver 12-perfis.md

[permissoes]
rede_externa = true                  # faz chamada HTTP para fora
filesystem   = false
exec         = false
migrations   = true
```

O `doctor` confere a declaração contra os imports reais: módulo que declara
`rede_externa = false` e importa `net/http` como cliente é reprovado. Permissão
declarada e não usada é aviso; usada e não declarada é erro.

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
