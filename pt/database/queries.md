# Construtor fluente de consultas

- [Componha antes de executar](#compose-before-execution)
- [Consulte faturas abertas sob um Grant](#query-open-invoices-under-a-grant)
- [Filtre e ordene com listas de valores permitidos](#filter-and-order-with-allowlists)
- [Escolha o contrato de registro ausente](#choose-the-missing-row-contract)
- [Verifique escopo de tenant e formato do SQL](#verify-tenant-scope-and-sql-shape)
- [Famílias de builders e operações terminais](#builder-and-terminal-families)
- [A autorização vem antes da consulta](#authorization-precedes-the-query)

<a name="compose-before-execution"></a>
## Componha antes de executar

Uma consulta começa pelo construtor que `aru model:build` gera ao lado de cada entidade, como `models.Invoices(db)`, e a composição fica separada da execução. Where, `OrderByDesc` e Limit descrevem a instrução sem enviá-la. Get é uma operação terminal: recebe o contexto da requisição e um `security.Grant`, aplica a coluna do tenant da tabela, compila o SQL do dialeto, executa e retorna a coleção gerada, um slice de ponteiros da entidade.

Manter o Grant na operação terminal é essencial. Construir uma consulta reutilizável não autoriza nada; a instrução que lê os registros deve carregar a autoridade e o tenant de quem a executa. Um Grant zero é rejeitado com `model.ErrNoTenant` antes que o SQL chegue à conexão.

<a name="query-open-invoices-under-a-grant"></a>
## Consulte faturas abertas sob um Grant

Invoice embute `model.Model` e `invoiceTable` declara sua tabela; os dois ficam em `app/Models/Invoice.go`. `aru model:build` escreve Invoices, `InvoiceQuery` e `InvoiceCollection` ao lado deles em `InvoiceQuery.go`, arquivo que é regenerado e não editado. `OpenInvoices` recebe banco instrumentado e Grant já emitido para `invoice.list`. A coluna de tenant padrão é `tenant_id`: o terminal adiciona tenant confiável independentemente do filtro de status.

O serviço obtém o Grant por `security.Authorize`, com sua Policy, o Subject autenticado, a ação `invoiceList` e o recurso da coleção. Passar esse Grant para `OpenInvoices` preserva a ordem usada na aplicação de exemplos publicada: primeiro autorizar, depois executar a operação de dados.

```go
type Invoice struct {
	model.Model

	ID        int64     `db:"id"`
	TenantID  string    `db:"tenant_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}

var invoiceTable = model.NewTable(model.TableSpec{
	Name: "invoices",
	New:  func() model.Entity { return new(Invoice) },
})

const invoiceList security.Action = "invoice.list"

func OpenInvoices(ctx context.Context, db *database.DB, grant security.Grant) (models.InvoiceCollection, error) {
	return models.Invoices(db).
		Where("status", "=", "open").
		OrderByDesc("created_at").
		OrderByDesc("id").
		Limit(25).
		Get(ctx, grant)
}
```

<a name="filter-and-order-with-allowlists"></a>
## Filtre e ordene com listas de valores permitidos

Valores passados a Where são parâmetros vinculados. Nomes de colunas, direções de ordenação, expressões selecionadas e fragmentos brutos fazem parte da estrutura SQL e não ficam seguros com um placeholder de valor. Traduza a escolha da requisição por um mapa fixo, como `map[string]``string{"newest": "created_at", "total": "total"}`; rejeite chaves ausentes antes de chamar `OrderBy` ou Select.

Componha os predicados adicionais antes da operação terminal e mantenha valores separados de identificadores. Use uma ordenação secundária estável quando dois registros puderem compartilhar o valor da coluna solicitada. Limit restringe a coleção retornada, mas não substitui os limites da camada de requisição quando o chamador escolhe o tamanho.

<a name="choose-the-missing-row-contract"></a>
## Escolha o contrato de registro ausente

First e Find retornam (nil, nil) sem correspondência. `FirstOrFail` e `FindOrFail` retornam erro compatível com `model.ErrModelNotFound`. Sole exige exatamente um: `model.ErrModelNotFound` para zero e `model.ErrMultipleRecordsFound` para mais de um. Escolha o terminal cuja cardinalidade faça parte do contrato, sem tratar toda ausência igualmente.

Preserve os erros encapsulados de banco e decodificação para que uma conexão indisponível ou um valor armazenado inválido não sejam tratados como simples ausência. Um Grant vazio é outro caso: `model.ErrNoTenant` indica que a execução foi recusada antes de construir uma instrução limitada ao tenant.

<a name="verify-tenant-scope-and-sql-shape"></a>
## Verifique escopo de tenant e formato do SQL

Execute a listagem num conector real registrado com registros abertos de dois tenants, empate na borda da página e registro fechado no tenant escolhido. Confirme que o primeiro recebe apenas seus 25 abertos mais recentes em `created_at` e id decrescentes, o segundo apenas seu registro, e Collector registra uma consulta com `tenant_id`, ambas as ordenações, limite e bindings de tenant e status.

Chame a mesma operação terminal com o Grant zero e um Collector novo. O resultado deve corresponder a `model.ErrNoTenant`, e a contagem de consultas deve permanecer zero. Execute também First e `FirstOrFail` para um registro visível apenas ao outro tenant, medindo os contratos de ausência opcional e de erro, sem inferi-los pelos tipos.

<a name="builder-and-terminal-families"></a>
## Famílias de builders e operações terminais

O construtor gerado inicia uma consulta tipada que repassa a `model.Builder`. Where e variantes adicionam condições vinculadas; Select, joins, agrupamento e ordenação moldam SQL; Limit restringe registros. Get retorna a coleção gerada, como `InvoiceCollection`. First e Find são leituras opcionais; `OrFail` exige registro e Sole exige exatamente um. Paginação e mutações são terminais separados com a mesma fronteira de contexto e Grant.

Este guia foi verificado contra hesape :hesapeVersion e framework :frameworkVersion. Versões posteriores podem adicionar ou mudar métodos; use a versão registrada na página ao comparar assinaturas.

<a name="authorization-precedes-the-query"></a>
## A autorização vem antes da consulta

Predicado de tenant isola dados, mas não é toda a policy. `security.Authorize` pergunta se Subject pode executar `invoice.list` e só emite Grant após permissão; o terminal usa o tenant do Grant. Em examples :examplesVersion, `app/Services/CategoryService.go` autoriza antes dos terminais, `app/Models/Category.go` declara a entidade e sua tabela e `tests/Feature/CategoryModel_test.go` testa tenant e recusa de Grant vazio.

Use relações quando registros relacionados fazem parte da consulta, paginação quando o cliente precisa de navegação e recurso explícito quando registros saem da aplicação. Nenhuma camada aceita coluna, tenant ou autoridade copiados diretamente da requisição.
