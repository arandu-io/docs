# 23 — Estabilidade, versionamento e depreciação

A dor nº 1 do Next.js hoje é mudança constante (doc 22). Um framework que se
vende como garantia e quebra a cada minor perde a tese. Esta política existe
para ser cumprida, não para ser bonita — decidir agora, porque depois da v1 vira
promessa retroativa.

## Versionamento

SemVer, com uma leitura estrita do que é breaking **em Go**:

| Mudança | Classificação |
|---|---|
| Adicionar método a interface pública | **breaking** — quebra quem implementa |
| Remover ou renomear símbolo exportado | breaking |
| Mudar assinatura de função exportada | breaking |
| Mudar comportamento observável sem mudar assinatura | breaking (o pior tipo) |
| Adicionar campo a struct de config | minor |
| Adicionar nova interface opcional (`Bootable`, `Migratable`) | minor |
| Nova regra no `doctor` que reprova código existente | **breaking** |

A última linha é a que costuma escapar. Regra nova de lint que reprova build em
CI é breaking change, mesmo sem tocar em API.

Consequência de desenho, já aplicada: interfaces públicas são **pequenas**, e
capacidade extra entra como interface opcional separada. É o que permite crescer
sem major.

## v0 até a Fase 3

Nada de `v1.0.0` antes de existir cliente real em produção. Em `v0.x` o contrato
é explícito: pode quebrar.

*Corrigido em 07/08/2026:* esta linha terminava em "e quebra vem com nota de
migração". Não veio: as 23 tags publicadas do `framework` não têm nota de
migração no repositório — não há `UPGRADE.md` nem CHANGELOG em lugar nenhum. A
parte que se cumpre é "pode quebrar"; ver "Depreciação" e "O que falta decidir".

O import path segue a regra do Go: `v2` em diante muda o caminho
(`github.com/arandu-io/framework/v2`). Isso é do módulo, não nosso — e é mais um
motivo para major ser raro.

## Cadência e janela de suporte

- **Major:** no máximo uma por ano.
- **Minor:** quando houver o que entregar. Sem calendário artificial.
- **Patch:** sempre que necessário.
- **Suporte:** a major corrente e a anterior por **12 meses** após o lançamento
  da nova. Correção de segurança nas duas.

## Depreciação

A política pretendida, e nada nela roda hoje:

1. Comentário `// Deprecated:` no símbolo, com a alternativa nomeada.
2. Aviso no `aru doctor` — não erro, aviso.
3. Mínimo de **duas minors** com o aviso ativo.
4. Remoção só na próxima major.

Nunca remover e avisar na mesma versão.

*Corrigido em 07/08/2026:* esta seção dizia "nada é removido sem passar por
depreciação", no presente, como se o ciclo estivesse em vigor. O que vale é que
o ciclo **nunca foi aplicado uma única vez**. `git grep "// Deprecated:"` em
todas as 23 tags publicadas do `framework` devolve zero ocorrências, e o `aru
doctor` não tem regra de depreciação — as severidades `Warning` e `Error`
existem em `internal/doctor`, o que falta é uma regra que as use para isso.

E houve o que depreciar. O levantamento de 07/08/2026 rodou `apidiff` par a par
nas tags do `framework` e contou **seis mudanças incompatíveis já publicadas sem
aviso nenhum** — entre elas `auth.New`, `config.Config.DatabaseURL`,
`data.Wrap`, `middleware.Observe` e os cinco verbos do `http.Router`. Uma delas
se confere sem ferramenta alguma: `data.InTransaction(ctx)` na `v0.5.0` virou
`data.InTransaction(ctx, db)` na `v0.13.2`, parâmetro novo, sem depreciação e
sem nota.

Isso não é um defeito de código a consertar: é uma escolha entre adotar o
mecanismo e reescrever a política. Ver "O que falta decidir", abaixo.

## Migração é código, não texto

Migração que não pode ser automática viraria item explícito num `UPGRADE.md`,
com exemplo antes e depois.

*Corrigido em 07/08/2026:* esta linha citava o `UPGRADE.md` no presente, como
arquivo existente. **Não existe `UPGRADE.md` em nenhum dos dez repositórios** —
`find` por `UPGRADE*` na árvore inteira não devolve nada, e também não há
CHANGELOG. As seis quebras da seção anterior saíram sem nota de migração de
espécie alguma.

*Atualizado em 14/08/2026:* o arquivo passou a existir, em um só lugar. O
`framework/UPGRADE.md` foi escrito em 08/08/2026, e é o único entre os catorze
repositórios — `ls arandu-io/*/UPGRADE.md` devolve esse e nada mais. O CHANGELOG
continua não existindo em nenhum: `find` por `CHANGELOG*` na árvore segue sem
devolver nada.

**`aru upgrade` — ainda não existe, e tem gatilho.** Ele rodaria as migrações
mecânicas (renomear símbolo, mover pacote, ajustar assinatura) como um codemod.
Não foi construído porque o framework está em `v0.x` e ainda não houve uma
mudança mecânica para automatizar — escrever o codemod antes da primeira seria
adivinhar a forma dela.

**Gatilho:** a primeira mudança de minor que exija a mesma edição em mais de um
lugar de um projeto real. Nesse dia o codemod se escreve contra aquela mudança,
que é o único jeito de saber se a ferramenta serve.

O alvo, escrito porque é o que se cobra da gente: **subir de minor não deve
tomar mais que dez minutos** num projeto médio.

## Como isso é verificado

Três verificações foram desenhadas para rodar na CI contra cada release. Uma
existe:

| Verificação | Estado |
|---|---|
| golden files do gerador — saída idêntica byte a byte, salvo mudança declarada | **existe** — `aru/internal/gen/golden_test.go`, com os fixtures em `internal/gen/testdata` e `-update` para regravar |
| `go build` da aplicação de exemplo da versão anterior contra a nova | **não existe** — a CI do `examples` constrói a versão corrente contra a corrente, nada mais |
| `apidiff` entre tags, com o diff aparecendo na revisão do PR | **existe desde 08/08/2026** — passo `api diff against the last release` na CI do `framework`, pacote a pacote contra a última tag |

*Atualizado em 08/08/2026:* a segunda das três passou a existir, e a política
deixou de ser declaração de intenção.

O passo `api diff against the last release` compara a árvore de trabalho com a
última tag, pacote a pacote, e imprime toda mudança incompatível no log do
build. Ele **não proíbe quebrar** — em `v0.x` quebrar é permitido, e congelar
uma forma sobre a qual ninguém construiu é pior. Ele proíbe quebrar **em
silêncio**: mudança incompatível sem entrada no `UPGRADE.md` reprova, num diff
que o revisor está olhando, em vez de no build de alguém três semanas depois.

O `UPGRADE.md` do `framework` existe, e foi escrito **por medição**: `apidiff`
entre a `v0.1.0` e a `v0.13.3`, sobre os quinze pacotes. O levantamento anterior
falava em seis quebras publicadas; são mais, e estão lá, agrupadas por pacote —
`data.Wrap`, `data.AppliedMigrations`, `Query.Filter`, os cinco verbos do
`Router` e o `Route.Name`, o `middleware.Observe`, quatro campos do `Collector`
e o `HandleDebugConsole`, `kernel.Routes`, `config.DatabaseURL` e o `auth.New`.

Falta a terceira: o build da aplicação de exemplo da versão anterior contra a
nova.

*Corrigido em 07/08/2026:* esta seção afirmava as três como fatos da CI, e só a
primeira era verdade. Dos nove `.github/workflows/ci.yml` do projeto, oito são de
repositório Go e rodam `gofmt`, `go vet`, `go test -race` e `govulncheck`. O
`go build` está só em três deles (`aru`, `arandu` e `examples`), e cada um
carrega ainda a sua própria guarda: orçamento de dependência no `framework` e no
`aru`, isolamento de driver no `database` e no `queue`, ausência de Node no `kv`,
no `queue` e no `storage`, `doctor` e build da imagem no `arandu`. O nono é o do
`homebrew-tap`, que não roda Go nenhum: `brew style`, `brew audit` e
`brew install`. Nenhum dos nove roda `apidiff`, e nenhum constrói a versão
anterior contra a nova.

A frase que fechava a seção estava certa, e agora se lê com o resultado na mão:

> Sem essas três, a política é declaração de intenção.

## O que falta decidir

Esta é a única parte deste documento que não é decisão tomada, e ela é do
proprietário. A política acima está escrita por inteiro; o mecanismo que a
cumpre não foi construído, e seis quebras já saíram publicadas por baixo dela.

São dois caminhos, e não dá para ficar entre os dois:

1. **Adotar o mecanismo.** `apidiff` na CI do `framework`, `UPGRADE.md` por
   versão e `// Deprecated:` em uso a partir da próxima remoção. É o que este
   documento já decidiu e ninguém implementou.
2. **Reescrever a política.** Declarar o pré-v1 explicitamente sem garantia de
   compatibilidade e reduzir esta página ao que se cumpre hoje, deixando o ciclo
   de depreciação para quando a v1 se aproximar.

O que não se sustenta é a terceira opção em vigor: a política escrita como se
valesse, sem nada que a verifique.

## O que autoriza um breaking change

Só três coisas:

1. Falha de segurança que não dá para corrigir de forma compatível.
2. Erro de contrato descoberto antes da v1.
3. Mudança na stdlib ou no toolchain do Go que force a mão.

"Ficou mais elegante" não é motivo. "Descobrimos um jeito melhor" não é motivo.
