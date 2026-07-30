# ADR 0005-idioma-do-codigo — Todo o código e toda a saída em inglês

**Status:** Aceita. Substitui a REGRA 6 do `CLAUDE.md`, que mandava comentar em português.

## Contexto

O produto é para o mercado global (`04-marca.md`). O esqueleto da Fase 1 foi
escrito com identificadores em inglês, mas comentários, mensagens de erro, saída
da CLI, textos da página de erro e mensagens de log em português.

Três consequências práticas disso, e a primeira é a que fecha a porta:

1. O `pkg.go.dev` gera a referência da API a partir dos comentários de doc. Doc
   em português é a primeira coisa que um desenvolvedor de fora vê, e derruba a
   adoção antes de qualquer avaliação técnica.
2. Mensagem de erro é a documentação mais lida do framework (`08-documentacao.md`).
   Erro em português num binário distribuído no mundo é um bug de produto.
3. Corrigir depois significa reabrir cada arquivo — e, se houver versão publicada,
   mudar texto de erro que gente já tem em runbook e em busca do Google.

## Decisão

Inglês em **tudo** que existe no repositório de código: identificadores,
comentários de doc, comentários internos, mensagens de erro, mensagens de log,
saída da CLI, textos da página de erro, nomes de teste e mensagens de falha de
teste.

Português segue em dois lugares, e só neles: a conversa entre o proprietário e
quem trabalha no código, e este conjunto de documentos de decisão (`docs/01` a
`08`, `adr/`), que é registro interno e não é lido por usuário do framework.

Consequência para a REGRA 6 do `CLAUDE.md`: a parte que dizia "comentários de
código em português, sem acento" está revogada. A regra de não usar acento em
arquivo `.go` perde o motivo junto, porque texto em inglês não tem acento.

## Consequência

Revisitar apenas se o produto deixar de ter mercado global como alvo — o que
mudaria `04-marca.md` primeiro, não este ADR.
