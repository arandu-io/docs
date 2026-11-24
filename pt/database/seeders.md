# Seeders

- [Seeders são operações explícitas da aplicação, identificadas por nome](#seeders-are-explicit-named-application-operations)
- [Defina dependências restritas e comportamento idempotente](#define-narrow-dependencies-and-idempotent-intent)
- [Execute o seeder raiz ou um seeder pelo nome](#run-the-root-or-one-named-seeder)
- [Crie contas iniciais pelos serviços da aplicação](#create-bootstrap-accounts-through-application-services)
- [Pare na primeira operação incompleta](#stop-at-the-first-incomplete-operation)
- [Mantenha segredos fora do fonte, histórico e lista de processos](#keep-secrets-out-of-source-history-and-process-lists)
- [Execute cada seeder duas vezes num banco não vazio](#run-every-seeder-twice-on-a-non-empty-database)
- [Use factories para os dados e serviços para as políticas](#use-factories-for-data-shape-and-services-for-policy)

<a name="seeders-are-explicit-named-application-operations"></a>
## Seeders são operações explícitas da aplicação, identificadas por nome

Um seeder implementa `database.Seeder[Deps]` com Name e `Run(context.Context, Deps)`. O projeto é dono de Deps, do registro, da raiz de fallback e da ordem de execução, então um seeder alcança apenas os serviços que o provider fornece. O Seed do provider, entregue ao builder com `WithSeeders(app.Seed)` em `bootstrap/app.go`, é o que `db:seed` e migrate `--seed` executam. Não há reflexão nem descoberta pelo sistema de arquivos.

<a name="define-narrow-dependencies-and-idempotent-intent"></a>
## Defina dependências restritas e comportamento idempotente

Declare dependências restritas aos serviços da aplicação, ao tenant configurado e aos argumentos restantes opcionais. Registre cada seeder e selecione `DatabaseSeeder` como padrão quando não houver nome. A segunda execução deve convergir, sem duplicar registros nem substituir credenciais silenciosamente.

- Registre cada seeder nomeado.
- Escolha uma raiz de fallback.
- Passe explicitamente um tenant configurado.
- Torne execuções repetidas observáveis e idempotentes.

<a name="run-the-root-or-one-named-seeder"></a>
## Execute o seeder raiz ou um seeder pelo nome

`aru db:seed` executa o seeder padrão; um nome posicional seleciona outro sem diferenciar maiúsculas de minúsculas. Nomes desconhecidos falham com a lista dos disponíveis, e a opção obsoleta `--class` é recusada. Os argumentos restantes pertencem ao seeder, acessados por `database.Flag` e `database.Switch`.

```sh
aru db:seed
aru db:seed UserSeeder -e you@example.com -p '<choose-a-password>'
aru db:seed UserSeeder -e you@example.com -p '<choose-a-password>' -r admin
```

<a name="create-bootstrap-accounts-through-application-services"></a>
## Crie contas iniciais pelos serviços da aplicação

A raiz do starter não cria conta alguma: migrar ou executar `aru db:seed` nunca cria, e em desenvolvimento ela popula as notas de exemplo, cujos autores não têm senha utilizável. Quando uma aplicação precisar de um usuário de bootstrap ou recuperação, execute `UserSeeder` pelo nome com e-mail e senha explícitos; adicione `-r` admin somente quando essa função for intencional. A policy da aplicação e o hashing permanecem no caminho normal do serviço.

- Mantenha contas fora de `DatabaseSeeder` até que a aplicação as adicione deliberadamente.
- Execute `UserSeeder` pelo nome para criar uma conta.
- Adicione `-r` admin somente como uma decisão explícita de autorização.
- Não altere os papéis existentes em uma recuperação apenas de senha.

<a name="stop-at-the-first-incomplete-operation"></a>
## Pare na primeira operação incompleta

`database.Seed` envolve uma falha de Run com o nome do seeder. Registro, dependência, tenant ou credencial explícita do `UserSeeder` ausentes falham antes da escrita, e uma raiz para em vez de ignorar um pré-requisito incompleto. Produção pede confirmação, salvo quando forçada explicitamente.

> [!WARNING]
> O framework não envolve todos os seeders filhos numa única transação. Use transação da aplicação quando exigir tudo ou nada, ou torne cada etapa idempotente e observável de forma independente.

<a name="keep-secrets-out-of-source-history-and-process-lists"></a>
## Mantenha segredos fora do fonte, histórico e lista de processos

> [!WARNING]
> Nunca crie uma conta padrão, commite uma senha, aceite um tenant escolhido pelo chamador ou exponha seeding por HTTP. Uma senha explícita do `UserSeeder` aparece no histórico do shell e na lista de processos; limpe essa entrada do histórico, altere a senha imediatamente pela aplicação e restrinja a autoridade de bootstrap a um serviço auditado da aplicação.

Senha de recuperação passada na linha de comando deve ser trocada imediatamente pela aplicação; não é o caminho comum de seed de implantação.

<a name="run-every-seeder-twice-on-a-non-empty-database"></a>
## Execute cada seeder duas vezes num banco não vazio

Migre um banco descartável com dados, execute duas vezes o seeder raiz e cada seeder nomeado, verifique tenant e papéis exatos sem duplicatas, teste nomes desconhecidos e flags obsoletas, provoque a falha de um filho e comprove que a confirmação de produção tem não como padrão.

- Execuções repetidas do seeder raiz convergem.
- Nomes desconhecidos listam opções válidas.
- Falha de um filho impede os seguintes.
- Segredos ausentes não deixam conta parcial.

<a name="use-factories-for-data-shape-and-services-for-policy"></a>
## Use factories para os dados e serviços para as políticas

Factories constroem dados tipados, migrations possuem schema, policies carregam autoridade e implantação decide quando uma operação inicial pode executar.

- [Factories](/docs/{{version}}/database/factories)
- [Migrações](/docs/{{version}}/database/migrations)
- [Autorização](/docs/{{version}}/security/authorization)
- [Publicação](/docs/{{version}}/operations/deployment)
