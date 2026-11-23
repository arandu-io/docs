# Relacionamentos e carregamento antecipado

- [Relações são metadados de model com carregamento explícito](#relations-are-model-metadata-with-explicit-loading)
- [Incorpore Model e defina as duas chaves](#embed-model-and-define-both-keys)
- [Declare uma relação e carregue-a antecipadamente](#declare-and-eager-load-one-relation)
- [Restrinja, conte, associe e sincronize](#constrain-count-attach-and-synchronize)
- [Não confunda relacionamento não carregado com vazio](#do-not-treat-an-unloaded-relation-as-empty)
- [O carregamento antecipado continua sendo um limite de leitura](#the-eager-path-is-still-a-read-boundary)
- [Conte consultas no caso simples de eager loading](#count-statements-for-the-simple-eager-case)
- [Coleções preservam o estado das relações carregadas](#collections-carry-loaded-relation-state-forward)

<a name="relations-are-model-metadata-with-explicit-loading"></a>
## Relações são metadados de model com carregamento explícito

Um relacionamento descreve como duas chaves de models se conectam e como os registros carregados se associam ao pai. Ele nunca é executado apenas porque um campo foi lido. Esse limite explícito mantém o acesso ao banco visível e garante que cada consulta de relacionamento receba o mesmo contexto e Grant da consulta pai.

<a name="embed-model-and-define-both-keys"></a>
## Incorpore Model e defina as duas chaves

As duas entidades embutem `model.Model` e declaram uma tabela. Registre cada relação na tabela do pai com Relate, numa função init ao lado da entidade, antes da primeira consulta: dê a ela um nome estável, a tabela relacionada, as chaves estrangeiras, locais ou proprietárias e os detalhes de pivô quando necessários. Restrições do banco devem impor a relação de chaves independentemente da API de consulta.

<a name="declare-and-eager-load-one-relation"></a>
## Declare uma relação e carregue-a antecipadamente

Escolha um construtor de relação exportado pela versão fixada para as chaves reais, como `HasOne`, `HasMany` ou `BelongsTo`. Variantes through, polimórficas e muitos-para-muitos exigem mais metadados; documente pelas assinaturas exatas da tag, não por contagem memorizada. With nomeia a relação antes do terminal. Related lê os registros carregados como `model.Rows`, e o `PostCollectionOf` gerado os converte para a entidade relacionada.

```go
func init() {
	userTable.Relate("posts", func(u *model.Model) model.Relation {
		return model.HasMany(u, postTable, "user_id", "id")
	})
}

user, err := models.Users(db).With("posts").FirstOrFail(ctx, g)
if err != nil { return err }
rows, loaded := user.Related("posts")
posts := models.PostCollectionOf(rows)
```

<a name="constrain-count-attach-and-synchronize"></a>
## Restrinja, conte, associe e sincronize

Predicados de relacionamento como `WhereHas` filtram os pais, e `WithCount` ou os demais helpers de agregação adicionam campos de resumo. `WithConstraints`, do builder do núcleo, restringe os resultados do carregamento antecipado a partir de um escopo local no bloco custom da entidade. Load e `LoadMissing` operam após a busca. Relacionamentos muitos-para-muitos expõem Attach, Detach e Sync sob um Grant, mantendo as alterações de pivô explícitas e restritas ao tenant.

<a name="do-not-treat-an-unloaded-relation-as-empty"></a>
## Não confunda relacionamento não carregado com vazio

Related informa se o valor foi carregado; verifique esse booleano para não confundir ausência de carregamento com relação vazia. Nomes desconhecidos, tipos de chave incompatíveis, metadados de pivô ausentes e falhas de banco retornam erros. Uma struct literal ou uma cópia por valor de um registro não tem model ao qual anexar a relação, e Load a recusa com `model.ErrUnwired`.

<a name="the-eager-path-is-still-a-read-boundary"></a>
## O carregamento antecipado continua sendo um limite de leitura

Toda consulta de relacionamento deve carregar o Grant e restringir os registros dos dois lados ao tenant. A falta do filtro de tenant no lado relacionado é especialmente perigosa: uma única consulta pode associar registros de outro cliente a vários pais. Nunca aceite nomes de relacionamentos ou colunas de pivô diretamente da requisição.

<a name="count-statements-for-the-simple-eager-case"></a>
## Conte consultas no caso simples de eager loading

Para um carregamento has-many simples e não polimórfico, verifique uma consulta dos pais mais uma consulta dos relacionados, e confirme que pais sem filhos recebem uma coleção carregada vazia. Casos polimórficos e aninhados podem exigir consultas adicionais: teste cada formato documentado separadamente. Inclua chaves coincidentes em outro tenant para provar que seus registros nunca são associados.

<a name="collections-carry-loaded-relation-state-forward"></a>
## Coleções preservam o estado das relações carregadas

Coleções de models podem executar Load ou `LoadMissing` para todos os registros em conjunto; os predicados de relacionamento permanecem no Builder. Os recursos HTTP decidem se os relacionamentos carregados entram na resposta. Chaves estrangeiras protegem a integridade do banco, mas não substituem as políticas nem os filtros de tenant nas leituras e escritas.
