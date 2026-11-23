# Guia de atualização

- [Atualize o grafo de módulos e a composição pertencente à aplicação](#upgrade-the-module-graph-and-application-owned-wiring)
- [Registre o artefato atual e o estado do schema](#record-the-current-artifact-and-schema-state)
- [Atualize um conjunto compatível de versões em uma branch](#move-one-compatible-release-set-on-a-branch)
- [Revise API, configuração, schema e saída gerada](#review-api-configuration-schema-and-generated-output)
- [Separe falhas de compilação de mudanças silenciosas de comportamento](#separate-compile-failures-from-silent-behavior-changes)
- [Reverifique toda fronteira de identidade e dados](#recheck-every-identity-and-data-boundary)
- [Execute verificações comportamentais no conjunto de versões desejado](#run-behavior-gates-against-the-target-release-set)
- [Restaure o código primeiro e reverta o schema apenas com evidência](#restore-code-first-and-undo-schema-only-with-evidence)
- [Leia as evidências do release antes de mudar produção](#read-release-evidence-before-changing-production)

> [!WARNING]
> **Limitações atuais.** Atualizações de módulos são explícitas e mudanças de schema devem ser compatíveis com rollback. Identidade de migration e relato de rollback irreversível ainda exigem cuidado; inspecione o changelog exato e exercite rollout com versões mistas antes de produção.

<a name="upgrade-the-module-graph-and-application-owned-wiring"></a>
## Atualize o grafo de módulos e a composição pertencente à aplicação

Uma aplicação Arandu reúne módulos Go comuns e código criado do esqueleto. Atualizar `go.mod` muda APIs de bibliotecas; não reescreve bootstrap, configuração, Dockerfile, migrations ou arquivos UI pertencentes ao projeto. Revise juntos os changelogs nas tags e o diff do esqueleto.

<a name="record-the-current-artifact-and-schema-state"></a>
## Registre o artefato atual e o estado do schema

Comece por um commit verde com arquivos de módulo restauráveis, backup adequado ao driver, digest da imagem implantada e lote de migrations registrado em `arandu_migrations`. Leia todas as tags entre a versão instalada e a desejada, nunca a branch de desenvolvimento móvel.

- Inventarie `go list -m` all e `aru version`.
- Registre as versões do esqueleto e da imagem.
- Identifique saída Kyse gerada e adaptações pertencentes à aplicação.
- Planeje a janela de schema com versões diferentes coexistindo.

<a name="move-one-compatible-release-set-on-a-branch"></a>
## Atualize um conjunto compatível de versões em uma branch

Fixe um conjunto revisado de releases, arrume o grafo e regenere as views. Depois siga o `UPGRADE.md` do esqueleto, entrada por entrada, a partir da release de que o projeto foi criado: ele diz o que um projeto muda em `bootstrap/app.go`, app/Providers, configuração, ambiente e testes. O `UPGRADE.md` de cada módulo, em framework, hesape, kyse, ui, joaju e aru, diz o que deixa de compilar ou passa a se comportar diferente e o que escrever no lugar. Porte apenas mudanças intencionais, por revisão, porque `aru new` não é um comando de atualização.

```sh
go get github.com/arandu-io/framework@:frameworkVersion \
  github.com/arandu-io/hesape@:hesapeVersion \
  github.com/arandu-io/kyse@:kyseVersion \
  github.com/arandu-io/hesape/database/connectors/pgx@:pgxVersion \
  github.com/arandu-io/hesape/database/connectors/sqlite@:sqliteVersion \
  github.com/arandu-io/hesape/redis@:redisVersion
go mod tidy
aru view:build
```

<a name="review-api-configuration-schema-and-generated-output"></a>
## Revise API, configuração, schema e saída gerada

Compile após atualizar pequenos grupos de módulos, busque símbolos removidos, compare variáveis de ambiente no parser e no consumidor, regenere views após mudanças Kyse, preserve blocos UI customizados, prepare schemas aditivos e atualize web e worker pelo mesmo artefato.

- Revise símbolos de API e conexão da configuração separadamente.
- Exercite os drivers de cache, sessão e fila escolhidos em produção.
- Gere novamente e compile as views.
- Teste o binário anterior com o schema aditivo desejado.

<a name="separate-compile-failures-from-silent-behavior-changes"></a>
## Separe falhas de compilação de mudanças silenciosas de comportamento

Símbolos removidos falham na compilação; variável de ambiente renomeada, padrão alterado, campo lido mas inativo ou reordenação de middleware podem compilar e quebrar produção. Imports gerados falham sem build de views, e schema destrutivo pode quebrar o binário anterior.

> [!WARNING]
> Não corrija uma atualização adicionando outro leitor local de ambiente sem rastrear o valor até o consumidor em execução. Parsers de componentes do framework e leitores da aplicação podem ler uma variável que nenhum construtor ativo usa.

<a name="recheck-every-identity-and-data-boundary"></a>
## Reverifique toda fronteira de identidade e dados

> [!WARNING]
> Reverifique backends e cookies de sessão, ordem do CSRF, propagação de Grant, hashing de tokens, estado OAuth, exposição de diagnóstico, raízes TLS e armazenamento do rate limit. Atualizar dependências não prova que a composição da aplicação adotou o caminho mais seguro.

<a name="run-behavior-gates-against-the-target-release-set"></a>
## Execute verificações comportamentais no conjunto de versões desejado

Execute geração de views, formatação, vet, testes com race, build, doctor, testes de container semelhantes à produção, sessão entre duas réplicas, migration duas vezes e rollback uma, compatibilidade do schema com o binário anterior, conexão das opções de fila e provas nas tags exatas.

- GOWORK=off `go test -race ./`... passa.
- `aru doctor --strict` passa.
- Views geradas correspondem ao fonte.
- Drivers de produção passam nos seus fluxos de feature.

<a name="restore-code-first-and-undo-schema-only-with-evidence"></a>
## Restaure o código primeiro e reverta o schema apenas com evidência

Reimplante a imagem anterior ou restaure arquivos de módulo e composição da aplicação. Mantenha schema aditivo durante rollback do binário. Reverta uma migration de release só após inspecionar lote e Down testado; nunca use fresh, reset ou refresh como rollback.

<a name="read-release-evidence-before-changing-production"></a>
## Leia as evidências do release antes de mudar produção

Notas de release definem mudanças, migrations definem compatibilidade do schema, implantação define rollout do artefato e configuração expõe riscos silenciosos de composição.

- [Notas de release](/docs/{{version}}/operations/release-notes)
- [Migrations](/docs/{{version}}/database/migrations)
- [Implantação](/docs/{{version}}/operations/deployment)
- [Configuração](/docs/{{version}}/architecture/configuration)
- [Guia de atualização do esqueleto](https://github.com/arandu-io/arandu/blob/:aranduVersion/UPGRADE.md)
