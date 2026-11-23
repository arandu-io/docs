# Notas de release

- [Uma nota de release descreve uma versão imutável](#a-release-note-describes-an-immutable-version)
- [Comece pela tag imutável anterior](#start-from-the-previous-immutable-tag)
- [Descreva a mudança observável](#describe-the-observable-change)
- [Agrupe adições, mudanças, correções e remoções](#group-additions-changes-fixes-and-removals)
- [Pare quando a evidência pertencer a outra revisão](#stop-when-evidence-belongs-to-another-revision)
- [Não publique segredos nem transforme uma correção em instrução de ataque](#do-not-publish-secrets-or-weaponize-a-fix)
- [Calcule a diferença da API antes de criar a tag](#compute-the-api-difference-before-tagging)
- [Corrija com uma nova versão](#correct-with-a-new-version)
- [Guias de atualização descrevem transições mais longas](#upgrade-guides-carry-the-longer-path)

> [!WARNING]
> **Limitações atuais.** Esta página explica evidências imutáveis de release; não autoriza promoção. Gates do repositório e decisão explícita do mantenedor continuam obrigatórios antes de alterar branches, tags ou releases.

<a name="a-release-note-describes-an-immutable-version"></a>
## Uma nota de release descreve uma versão imutável

Cada módulo publicado está abaixo de 1.0 e registra mudanças incompatíveis no changelog. Uma tag é imutável: corrigir um artefato publicado exige uma nova versão, nunca mover a tag antiga. Isso permite recuperar o mesmo código em builds e revisões de incidentes.

<a name="start-from-the-previous-immutable-tag"></a>
## Comece pela tag imutável anterior

Um candidato a release exige revisão limpa, gates aprovados, entrada de changelog e tag anterior conhecida. Compare essa tag com o commit exato pretendido; comparar duas branches móveis produz notas que ninguém consegue reproduzir depois.

<a name="describe-the-observable-change"></a>
## Descreva a mudança observável

Descreva o que o consumidor precisa mudar, o comportamento diferente e a primeira versão do módulo que o contém. Identifique precisamente símbolos removidos ou renomeados e o caminho substituto. Os comandos abaixo ilustram uma inspeção; não são o processo de release nem prova de gates aprovados. Scripts e políticas do repositório continuam sendo a autoridade.

```bash
PREVIOUS_TAG=$(git describe --tags --abbrev=0 HEAD^)
git diff --stat "$PREVIOUS_TAG"...HEAD
GOWORK=off go vet ./...
GOWORK=off go test -race -count=1 -p 1 ./...
```

<a name="group-additions-changes-fixes-and-removals"></a>
## Agrupe adições, mudanças, correções e remoções

Escreva para o consumidor e vincule cada nota de migração ao módulo afetado. Destaque configuração, migrações de dados, segurança e janelas de compatibilidade. Quando módulos precisam mudar juntos, indique o conjunto compatível em vez de deixar o consumidor descobri-lo por erros de build.

<a name="stop-when-evidence-belongs-to-another-revision"></a>
## Pare quando a evidência pertencer a outra revisão

Gate falho, arquivo gerado sujo, comparação de API ausente, changelog ausente ou tag fora da branch prevista bloqueiam o release. Não copie resultados de ontem para as notas de hoje; resultados comprovam um commit e perdem validade quando ele muda.

<a name="do-not-publish-secrets-or-weaponize-a-fix"></a>
## Não publique segredos nem transforme uma correção em instrução de ataque

Notas de release podem descrever comportamento afetado e atualização segura sem credenciais, detalhes privados de incidentes ou instruções de exploração antes da divulgação coordenada. Mudanças de dependências exigem repetir verificações de segurança e licença no grafo do release, não presumir resultados do workspace.

<a name="compute-the-api-difference-before-tagging"></a>
## Calcule a diferença da API antes de criar a tag

Módulos públicos comparam a superfície exportada com a tag anterior e exigem notas de atualização para mudanças relevantes de assinatura. Formatação, vet, race, vulnerabilidades, artefatos gerados e ancestralidade da tag são gates cujos resultados atuais precisam ser medidos. Listá-los aqui não comprova aprovação.

<a name="correct-with-a-new-version"></a>
## Corrija com uma nova versão

Se um artefato ou nota publicados estiverem errados, preserve a tag original, reverta ou corrija na branch de release e publique a próxima versão com correção explícita. Assim consumidores fixam, comparam e revertem sem uma tag resolver para códigos diferentes em duas máquinas.

<a name="upgrade-guides-carry-the-longer-path"></a>
## Guias de atualização descrevem transições mais longas

Notas de release resumem uma versão. Transições de dados ou API em várias etapas pertencem ao guia de atualização, ligado pela nota. Changelogs, política de versionamento semântico, diferenças de API gerada e evidências de rollback são entradas separadas para a mesma decisão; nenhuma prova sozinha que o release está pronto.
