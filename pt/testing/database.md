# Testes de banco de dados

- [Testes de banco comprovam transações, dialetos e escopo do tenant](#database-tests-prove-transactions-dialects-and-tenant-scope)
- [Crie um schema descartável e um tenant explícito](#create-disposable-schema-and-one-explicit-tenant)
- [Execute o contrato do repositório com um Grant real](#run-the-repository-contract-with-a-real-grant)
- [Exercite as fronteiras de transação, migration e consulta](#exercise-transaction-migration-and-query-boundaries)
- [Não deixe estado compartilhado tornar a suíte dependente da ordem](#do-not-let-shared-state-make-a-green-suite-order-dependent)
- [Todo fixture e asserção permanecem no escopo do tenant](#every-fixture-and-assertion-stays-tenant-scoped)
- [Execute o contrato em todo conector suportado](#run-the-contract-against-every-supported-connector)
- [Use factories para a estrutura e migrations para o estado do catálogo](#use-factories-for-shape-and-migrations-for-catalogue-state)

<a name="database-tests-prove-transactions-dialects-and-tenant-scope"></a>
## Testes de banco comprovam transações, dialetos e escopo do tenant

Um teste de banco usa conector, gramática de consulta, model, transação e migration reais escolhidos pela aplicação. Fakes em memória são contratos úteis, mas não comprovam placeholders SQL, comportamento DDL, locks, cancelamento ou predicados de tenant.

<a name="create-disposable-schema-and-one-explicit-tenant"></a>
## Crie um schema descartável e um tenant explícito

Dê a cada teste um banco só dele e execute as migrations pelo conector publicado: `test.Serve` sobe a aplicação sobre uma cópia de um modelo Postgres já migrado ou sobre um arquivo SQLite temporário, e as opções Migrate, Seed, SQLite e Database dizem como o schema é montado e em que servidor o banco é criado. O isolamento vem do armazenamento exclusivo do teste; nunca reutilize credenciais de produção nem um schema compartilhado de desenvolvimento.

- Isole o schema por teste ou pacote.
- Execute migrations reais.
- Use sementes determinísticas nas factories.
- Limpe com rollback transacional ou descarte dos recursos próprios.

<a name="run-the-repository-contract-with-a-real-grant"></a>
## Execute o contrato do repositório com um Grant real

Exercite o repositório protegido com um Grant real e inspecione o banco independentemente pelas asserções de `framework/test`. Inclua `tenant_id` em toda comparação de registros do cliente para que outro tenant não faça a asserção passar por acidente.

```go example
func assertInvoiceStored(ctx context.Context, t *testing.T, db *database.DB, invoiceID string) {
	test.AssertDatabaseHas(t, ctx, db, "invoices", test.Match{
		"id":        invoiceID,
		"tenant_id": "acme",
	})
}
```

<a name="exercise-transaction-migration-and-query-boundaries"></a>
## Exercite as fronteiras de transação, migration e consulta

Teste commit e rollback, aninhamento onde suportado, cancelamento, conflitos, ausência, ordenação, paginação, locks, migrations Up e Down, escritas parciais de factories e inserção atômica de fila ou outbox quando compartilham SQL.

- Use fixtures explícitos de dialeto SQL.
- Verifique registros e predicados de tenant.
- Teste cancelamento do contexto durante I/O.
- Execute rollback transacional como um caminho real de falha.

<a name="do-not-let-shared-state-make-a-green-suite-order-dependent"></a>
## Não deixe estado compartilhado tornar a suíte dependente da ordem

Schema, relógio, IDs ou registros globais compartilhados podem fazer testes passarem apenas em determinada ordem. Detecte registros vazados, nomes fixos de migration, conexões retidas, colisões de portas em paralelo e suposições de que SQLite se comporta como todo dialeto de produção. Uma expectativa de mock de que o SQL foi chamado não prova que a consulta devolveu apenas linhas autorizadas; mantenha ao menos um contrato com conector real para cada caminho de implantação suportado.

<a name="every-fixture-and-assertion-stays-tenant-scoped"></a>
## Todo fixture e asserção permanecem no escopo do tenant

> [!WARNING]
> Use dados sintéticos de tenant, nunca despeje registros de outros tenants ao falhar, vincule valores em vez de interpolar SQL e impeça que Grants de teste vazem para caminhos de requisição da aplicação.

`security.SystemGrant` recebe ação e tenant explícitos. Registre a justificativa do teste no nome ou comentário e mantenha o Doctor estrito recusando uso não autorizado em produção.

<a name="run-the-contract-against-every-supported-connector"></a>
## Execute o contrato em todo conector suportado

Execute sucesso, ausência, duplicata, Grant vazio, tenant alheio, cancelamento, rollback, repetição de migration, Down e Up e casos paralelos sem race contra cada conector suportado na CI.

- Todo conector satisfaz o mesmo contrato de repositório.
- Dados de outro tenant ficam invisíveis.
- O rollback não deixa registros.
- Migrations são repetíveis e reversíveis onde declarado.

<a name="use-factories-for-shape-and-migrations-for-catalogue-state"></a>
## Use factories para a estrutura e migrations para o estado do catálogo

Models definem o mapeamento, factories geram registros tipados, migrations possuem o schema e contratos permitem executar a mesma suíte comportamental em cada adapter.

- [Models](/docs/{{version}}/database/models)
- [Factories](/docs/{{version}}/database/factories)
- [Migrations](/docs/{{version}}/database/migrations)
- [Contratos](/docs/{{version}}/architecture/contracts-adapters)
