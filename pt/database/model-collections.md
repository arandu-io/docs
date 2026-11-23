# Coleções de models

- [Uma coleção de models é um slice de structs da aplicação](#a-model-collection-is-a-slice-of-application-structs)
- [Consulte e inspecione o conjunto tipado](#query-and-inspect-the-typed-set)
- [Carregue, selecione, compare e converta](#load-select-compare-and-convert)
- [Separe vazio, ausente e não consultável](#keep-empty-missing-and-unqueryable-distinct)
- [Teste separadamente as operações em memória e no armazenamento](#prove-set-behavior-and-storage-behavior-separately)
- [Grupos de métodos de coleção](#collection-method-groups)
- [Use paginação quando o cliente precisar de uma janela de resultados](#use-pagination-when-the-client-needs-a-window)

<a name="a-model-collection-is-a-slice-of-application-structs"></a>
## Uma coleção de models é um slice de structs da aplicação

A coleção gerada, como `UserCollection`, é um slice de ponteiros para a mesma entidade que um terminal único retorna. Seus métodos repassam a `model.Rows`, que conhece chaves, atributos carregados, relações e persistência. Por ser um slice, ela é percorrida, indexada e ordenada com a linguagem e a biblioteca padrão, sem que a camada de model duplique uma biblioteca de coleções.

<a name="query-and-inspect-the-typed-set"></a>
## Consulte e inspecione o conjunto tipado

Get retorna a coleção gerada diretamente. len descreve o conjunto, `ModelKeys` extrai as chaves primárias e Find busca na coleção já carregada, sem nova consulta ao banco. Essas operações em memória não precisam de Grant porque não acessam o armazenamento.

```go
users, err := models.Users(db).OrderBy("name").Get(ctx, g)
if err != nil { return err }

ids := users.ModelKeys()
user := users.Find(id)
```

<a name="load-select-compare-and-convert"></a>
## Carregue, selecione, compare e converta

Load e `LoadMissing` buscam relacionamentos para o conjunto, e `LoadCount` adiciona contagens de relações. Only, Except e Unique operam pelas chaves dos models e retornam o mesmo tipo de coleção, e Pluck lê um atributo carregado. Rows entrega o conjunto a `model.Rows` para Diff, Intersect, Merge, `LoadAggregate` e `ToQuery`, e o `UserCollectionOf` gerado converte os registros de volta.

<a name="keep-empty-missing-and-unqueryable-distinct"></a>
## Separe vazio, ausente e não consultável

Find retorna nil quando nenhum registro do conjunto tem a chave; `FindOrFail` torna ausência um erro. `ToQuery` em conjunto vazio retorna `model.ErrEmptyCollection`; relações e Fresh podem falhar no banco. Uma struct literal acrescentada ao slice não tem model, e Load recusa o conjunto com `model.ErrUnwired`.

<a name="prove-set-behavior-and-storage-behavior-separately"></a>
## Teste separadamente as operações em memória e no armazenamento

Teste comparações de chave, unicidade e vazio em registros montados com New ou com o Make de uma factory, que não fazem I/O. Integre Load, Fresh, Push e `ToQuery` com dois tenants e contagem de consultas. Em `LoadMissing`, pré-carregue só parte do conjunto e confirme busca apenas do ausente sem descartar relações já anexadas.

<a name="collection-method-groups"></a>
## Grupos de métodos de coleção

A coleção gerada traz `ModelKeys`, Find, `FindOrFail`, Contains, Pluck, Only, Except, Unique, Load, `LoadMissing`, `LoadCount`, Fresh, `ToArray` e Rows. `model.Rows` acrescenta Count, `IsEmpty`, First, Diff, Intersect, Merge, Push, `LoadAggregate`, `ToQuery`, `MakeHidden` e `MakeVisible`. Métodos de visibilidade ajustam a representação do model em array, não o tipo declarado da resposta HTTP.

<a name="use-pagination-when-the-client-needs-a-window"></a>
## Use paginação quando o cliente precisar de uma janela de resultados

Uma coleção contém registros já retornados; não conhece total de páginas nem cursores seguintes. Os terminais de paginação retornam a coleção ao lado de uma página com o estado de navegação. Para processamento sem limite no servidor, Chunk, Each ou Cursor controlam memória sem montar primeiro a tabela inteira numa coleção.
