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
é explícito: pode quebrar, e quebra vem com nota de migração.

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

Nada é removido sem passar por depreciação:

1. Comentário `// Deprecated:` no símbolo, com a alternativa nomeada.
2. Aviso no `aru doctor` — não erro, aviso.
3. Mínimo de **duas minors** com o aviso ativo.
4. Remoção só na próxima major.

Nunca remover e avisar na mesma versão.

## Migração é código, não texto

Migração que não pode ser automática vira item explícito no `UPGRADE.md`, com
exemplo antes e depois.

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

Na CI, contra cada release:

- `go build` da aplicação de exemplo da versão anterior contra a nova
- golden files do gerador — saída idêntica byte a byte, salvo mudança declarada
- `apidiff` entre tags, com o diff aparecendo na revisão do PR

Sem essas três, a política é declaração de intenção.

## O que autoriza um breaking change

Só três coisas:

1. Falha de segurança que não dá para corrigir de forma compatível.
2. Erro de contrato descoberto antes da v1.
3. Mudança na stdlib ou no toolchain do Go que force a mão.

"Ficou mais elegante" não é motivo. "Descobrimos um jeito melhor" não é motivo.
