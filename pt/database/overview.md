# Conexões de banco de dados

- [Um caminho de conexão com instrumentação](#one-connection-path-with-instrumentation)
- [Vincule um driver e abra a URL configurada](#link-one-driver-and-open-the-configured-url)
- [Configure pools e transações na fronteira](#configure-pools-and-transactions-at-the-boundary)
- [Interrompa a inicialização se faltar o driver ou a conexão for inválida](#fail-boot-on-an-absent-driver-or-invalid-connection)
- [Teste o contrato do conector em bancos reais](#run-the-connector-contract-against-real-engines)
- [Configuração, DB e Connection](#configuration-db-and-connection)
- [Repositórios são a porta da aplicação](#repositories-are-the-application-door)

<a name="one-connection-path-with-instrumentation"></a>
## Um caminho de conexão com instrumentação

`database.Open` constrói o DB da aplicação a partir de configuração tipada. O handle conhece o dialeto, adapta placeholders, aplica a política de pool e informa as instruções ao coletor configurado. Manter um único handle encapsulado evita caminhos alternativos que executem consultas sem os diagnósticos usados pelo restante da aplicação.

<a name="link-one-driver-and-open-the-configured-url"></a>
## Vincule um driver e abra a URL configurada

Um conector é um módulo Go separado e se registra durante a inicialização do pacote. Faça import em branco exatamente dos motores que este binário pode falar. Numa aplicação, `foundation.NewDrivers` faz esta chamada com a configuração que `DATABASE_URL` e as configurações do pool produziram, responde o handle como `drivers.DB` e fecha o pool no desligamento e ao fim de cada comando; repositórios e handlers não podem chamar outro leitor de ambiente.

```go example
import (
	"github.com/arandu-io/hesape/database"
	_ "github.com/arandu-io/hesape/database/connectors/sqlite"
)

func OpenDatabase(cfg database.Config) (*database.DB, func(), error) {
	db, closeDB, err := database.Open(cfg)
	if err != nil {
		return nil, nil, err
	}
	return db, closeDB, nil
}
```

<a name="configure-pools-and-transactions-at-the-boundary"></a>
## Configure pools e transações na fronteira

A URL seleciona o dialeto registrado e contém os parâmetros de conexão do driver. Limites e tempo de vida do pool pertencem à configuração do banco, não aos repositórios individuais. `database.Transaction` mantém o trabalho na transação associada ao contexto; passe esse contexto por todas as chamadas para que operações aninhadas não escapem para o pool base.

<a name="fail-boot-on-an-absent-driver-or-invalid-connection"></a>
## Interrompa a inicialização se faltar o driver ou a conexão for inválida

Uma URL que nomeia driver não ligado ao binário não pode ser corrigida em runtime; interrompa o boot com o erro. Diferencie servidor inacessível, credenciais inválidas, arquivo SQLite ausente e configuração malformada. Não use fallback para outro banco com propriedade de dados distinta.

<a name="run-the-connector-contract-against-real-engines"></a>
## Teste o contrato do conector em bancos reais

A suíte de conformidade dos conectores comprova o contrato compartilhado entre `PostgreSQL`, `MySQL` e SQLite. Na aplicação, verifique que o driver selecionado está disponível, que uma consulta parametrizada chega ao coletor, que o rollback não deixa registros e que o encerramento libera o pool. Um mock sozinho não comprova a vinculação no dialeto.

<a name="configuration-db-and-connection"></a>
## Configuração, DB e Connection

Configuração e `ParseURL` produzem Config; Open cria DB; Transaction delimita trabalho; `AvailableDrivers` lista conectores registrados. Escolha um dono da configuração, sem chamar `database.Load` novamente no construtor. DB e Connection são infraestrutura abaixo do repositório. `ForMigrations` e `ForSchema` expõem interfaces restritas sem transformar repositórios em ferramentas de schema.

<a name="repositories-are-the-application-door"></a>
## Repositórios são a porta da aplicação

A conexão não aplica policy de tenant a SQL arbitrário; módulos acessam registros por repositório ou terminal de model que exige Grant. Migrations operam abaixo dessa fronteira por atuar no schema antes de requisições. Escolha conectores no bootstrap; nenhum handler importa driver ou abre outro pool.
