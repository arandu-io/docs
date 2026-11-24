# Migrations e schema

- [Uma migração é código Go registrado com identidade imutável](#a-migration-is-registered-go-code-with-immutable-identity)
- [Planeje a implantação com versões coexistentes](#design-for-the-mixed-version-rollout-window)
- [Registre uma alteração reversível de esquema](#register-one-reversible-schema-change)
- [Inspecione, simule, aplique e reverta de forma deliberada](#inspect-preview-apply-and-roll-back-deliberately)
- [Entenda resultados de salto, transação e rollback](#understand-skip-transaction-and-rollback-outcomes)
- [Separe a autoridade de publicação da autoridade da requisição](#separate-deployment-authority-from-request-authority)
- [Teste as duas direções do esquema e o binário anterior](#test-both-schema-directions-and-the-old-binary)
- [Restaure o serviço antes de decidir reverter o esquema](#restore-service-before-choosing-a-schema-down-action)
- [Coordene schema, implantação e dados iniciais](#coordinate-schema-deployment-and-seed-data)

<a name="a-migration-is-registered-go-code-with-immutable-identity"></a>
## Uma migração é código Go registrado com identidade imutável

Cada migração implementa `migrations.Migration`, registra-se em init e é ordenada por `GetName`. O repositório registra esse nome exato e o lote somente após o sucesso de Up. O starter usa `arandu_migrations` e importa `database/migrations` no bootstrap para que todos os comandos enxerguem o registro.

<a name="design-for-the-mixed-version-rollout-window"></a>
## Planeje a implantação com versões coexistentes

Documente como o binário anterior se comporta após a chegada do novo esquema. Adicione colunas anuláveis ou com valor padrão e tabelas ignoradas pelo código antigo; distribua renomeações e remoções entre versões. Os comandos de produção exigem um cache RESP compartilhado para o lock de migração.

- Escolha um nome imutável para a migração.
- Planeje a compatibilidade com o binário anterior.
- Teste a gramática do banco escolhido.
- Configure lock entre réplicas em produção.

<a name="register-one-reversible-schema-change"></a>
## Registre uma alteração reversível de esquema

Implemente Up e Down sobre a conexão de migration, registre o valor e afirme `ReversibleMigration` na compilação. Embuta `BaseMigration` para conexão padrão, `ShouldRun` normal e intenção transacional publicada.

```go example
type AddLocaleToUsers struct{ migrations.BaseMigration }

func (AddLocaleToUsers) GetName() string { return "20260830_0004_add_locale_to_users" }
func (AddLocaleToUsers) Up(ctx context.Context, conn migrations.Connection) error {
	return conn.Schema().Table(ctx, "users", func(t *schema.Blueprint) { t.String("locale").Nullable() })
}
func (AddLocaleToUsers) Down(ctx context.Context, conn migrations.Connection) error {
	return conn.Schema().Table(ctx, "users", func(t *schema.Blueprint) { t.DropColumn("locale") })
}
func init() { migrations.Register(AddLocaleToUsers{}) }

var _ migrations.ReversibleMigration = AddLocaleToUsers{}
```

<a name="inspect-preview-apply-and-roll-back-deliberately"></a>
## Inspecione, simule, aplique e reverta de forma deliberada

Gere uma migration, inspecione o status e a saída de pretend, aplique normalmente ou em passos e reverta apenas um último lote revisado ou um passo explícito. Use o schema builder para DDL portátil e Statement ou Select explícitos para um backfill revisado. As colunas de tempo usam precisão de microssegundos por padrão, para que duas linhas escritas no mesmo segundo mantenham a ordem; `schema.DefaultTimePrecision` muda esse padrão, e uma migration `migrations.WidenTimePrecision` eleva as colunas de uma tabela criada com segundos inteiros, que é como as tabelas dos próprios pacotes se atualizam.

- Use `aru make:migration` com `--create` ou `--table`.
- Use `aru migrate:status --pending` e `aru migrate --pretend`.
- Use `aru migrate` ou `aru migrate --step`.
- Use `aru migrate:rollback` com decisão explícita de lote.

<a name="understand-skip-transaction-and-rollback-outcomes"></a>
## Entenda resultados de salto, transação e rollback

ShouldRun=false é ignorado e não é registrado, portanto será reavaliado depois. DDL transacional varia conforme a gramática, e `OutsideTransaction` pode deixar o esquema parcialmente alterado. O rollback valida o lote antes de executar: a ausência de Down sem declaração explícita impede o trabalho; uma migração explicitamente Irreversible permanece aplicada e registrada.

> [!WARNING]
> Uma falha de Down em tempo de execução pode interromper o rollback depois que entradas reversíveis anteriores já foram concluídas.

<a name="separate-deployment-authority-from-request-authority"></a>
## Separe a autoridade de publicação da autoridade da requisição

> [!WARNING]
> Migrations não recebem Grant: são DDL de implantação, não requisições. Dê ao pipeline credencial restrita, vincule valores, revise SQL bruto, não exponha segredos na saída e nunca execute migration em handler HTTP nem ao iniciar réplicas.

<a name="test-both-schema-directions-and-the-old-binary"></a>
## Teste as duas direções do esquema e o binário anterior

Em cada banco descartável suportado, migre duas vezes, inspecione os lotes, execute Down e depois Up, confira o SQL de simulação, injete falhas em Up e Down, teste a validação de irreversibilidade, a disputa e a indisponibilidade do lock, e inicie o binário anterior com o esquema aditivo.

- Repetir migrate é idempotente.
- Down e Up restauram o catálogo esperado.
- A disputa pelo lock permite uma única aplicação.
- O binário anterior continua lendo e escrevendo.

<a name="restore-service-before-choosing-a-schema-down-action"></a>
## Restaure o serviço antes de decidir reverter o esquema

Reimplante o binário anterior quando o esquema aditivo for compatível. Para um lote reversível com erro, inspecione o status e execute o rollback específico já testado. Para DDL parcial fora de transação ou mudanças irreversíveis, inspecione o catálogo real e o backup e escreva uma correção adiante, em vez de editar os registros do repositório de migrações.

`migrate:fresh`, `migrate:reset` e `migrate:refresh` só rodam com `APP_ENV=dev`, e um `APP_ENV` ausente é produção, porque eles esvaziam o schema; o rollback continua sendo uma decisão explícita de lote.

<a name="coordinate-schema-deployment-and-seed-data"></a>
## Coordene schema, implantação e dados iniciais

Implantação possui a janela de rollout, atualizações escolhem versões, seeders inserem registros deliberados e testes de banco exercitam ambas as direções.

- [Publicação](/docs/{{version}}/operations/deployment)
- [Atualizações](/docs/{{version}}/operations/upgrades)
- [Seeders](/docs/{{version}}/database/seeders)
- [Testes de banco de dados](/docs/{{version}}/testing/database)
