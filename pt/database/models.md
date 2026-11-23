# Models

- [O model é a struct da aplicação](#the-model-is-the-application-s-struct)
- [Declare colunas e construa os metadados do model](#declare-columns-and-construct-model-metadata)
- [Consulte um registro tipado](#query-one-typed-row)
- [Crie, preencha, salve, exclua e carregue relacionamentos](#create-fill-save-delete-and-eager-load)
- [Mantenha visíveis registros ausentes e metadados inválidos](#keep-missing-rows-and-invalid-metadata-visible)
- [O escopo do tenant sempre vem do Grant](#tenant-scope-always-comes-from-the-grant)
- [Teste o model com um tenant vizinho hostil](#test-the-model-with-a-hostile-neighboring-tenant)
- [Models não definem o conteúdo da resposta pública](#models-do-not-define-the-public-payload)

<a name="the-model-is-the-application-s-struct"></a>
## O model é a struct da aplicação

Uma consulta de model hidrata campos Go declarados. Não existe proxy dinâmico de atributos entre código e registro; um campo com nome errado falha no build. A entidade embute `model.Model`, que adiciona operações de instância e armazenamento de relações, e um `model.NewTable` no nível do pacote declara a tabela da qual ela é registro.

<a name="declare-columns-and-construct-model-metadata"></a>
## Declare colunas e construa os metadados do model

Exporte os campos persistidos e descreva suas colunas com tags db quando o nome padrão em `snake_case` não bastar. Configure tabela, chave primária, coluna do tenant, timestamps, exclusão lógica, escopos globais e eventos em `model.TableSpec`, onde cada escolha fica visível para quem lê. Execute `aru model:build` ao adicionar ou alterar uma entidade: ele escreve a consulta tipada ao lado da struct, e `aru doctor` aponta a que estiver desatualizada.

<a name="query-one-typed-row"></a>
## Consulte um registro tipado

Inicie a consulta gerada da entidade, componha condições e chame terminal com contexto e Grant. O ponteiro retornado tem o mesmo tipo User declarado, e os métodos do Model embutido são promovidos para ele.

```go
type User struct {
	model.Model

	ID    int64  `db:"id"`
	Email string `db:"email"`
}

var userTable = model.NewTable(model.TableSpec{
	Name: "users",
	New:  func() model.Entity { return new(User) },
})

user, err := models.Users(db).Where("email", "=", email).First(ctx, g)
if err != nil { return nil, err }
return user, nil
```

<a name="create-fill-save-delete-and-eager-load"></a>
## Crie, preencha, salve, exclua e carregue relacionamentos

Monte um registro novo com `Users(db)`.`New()` e atribua seus campos um a um; uma struct literal não tem conexão, e seus métodos de escrita retornam `model.ErrUnwired`. As operações terminais da consulta criam e atualizam conjuntos; os registros retornados por uma consulta ou por New podem executar Save ou Delete sob um Grant. Fill aplica atributos exportados reconhecidos e ignora chaves desconhecidas; `ForceFill` preserva atributos brutos desconhecidos, mas também não alcança campos não exportados. Nenhum dos dois é uma allowlist: monte as mudanças a partir de uma requisição tipada, como um mapa literal das colunas nomeadas que o caso de uso pode alterar, e nunca entregue o mapa de entrada de uma requisição a Fill ou Update. With declara relacionamentos antes da execução, e Load os adiciona depois com outro Grant explícito.

<a name="keep-missing-rows-and-invalid-metadata-visible"></a>
## Mantenha visíveis registros ausentes e metadados inválidos

Separe registro não encontrado de erros de conexão ou decodificação quando o endpoint precisar de 404. Tags inválidas, conversão não suportada, configuração de tabela ausente e falhas de escrita retornam erro; não continue com valor zero parcialmente hidratado. Campos não exportados não são persistidos por reflexão deliberadamente.

<a name="tenant-scope-always-comes-from-the-grant"></a>
## O escopo do tenant sempre vem do Grant

As operações terminais de leitura e escrita aplicam o tenant configurado no model usando o tenant do Grant, e a coluna de tenant recebe seu valor apenas do Grant. Na inserção, esse tenant confiável substitui qualquer atributo conflitante enviado pelo chamador. Não conte com a mesma proteção quando uma linha existente muda: um mapa da requisição pode levar a coluna de tenant, um papel ou um hash de senha sob uma grafia que o model aceita quando chega a Fill ou Update; as mudanças recebem apenas um mapa literal de colunas nomeadas, montado a partir de uma requisição tipada. Uma tabela compartilhada por todos os tenants define Global no seu `TableSpec`, a única configuração que remove o filtro de tenant, e deve ser uma exceção explícita e revisável.

<a name="test-the-model-with-a-hostile-neighboring-tenant"></a>
## Teste o model com um tenant vizinho hostil

Crie identificadores ou atributos iguais sob Grants de dois tenants e comprove o isolamento em Find, atualização, exclusão, agregações e carregamento antecipado. Verifique que o atributo de tenant é substituído na criação, que Save da instância mantém o escopo, que o Grant zero não executa consultas e que uma struct literal recusa toda escrita com `model.ErrUnwired`.

<a name="models-do-not-define-the-public-payload"></a>
## Models não definem o conteúdo da resposta pública

Relações anexam registros carregados ao estado embutido do model; coleções adicionam operações conscientes de registros e a página de paginação carrega a navegação. Recursos HTTP ainda escolhem campos que saem do processo. Coluna nova não deve virar resposta automaticamente só porque precisa ser exportada para persistir.
