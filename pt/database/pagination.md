# Paginação

- [Um paginador fornece navegação, não apenas registros](#a-paginator-answers-navigation-not-only-rows)
- [Escolha paginação com total, simples ou por cursor](#choose-length-aware-simple-or-cursor-pagination)
- [Preserve os filtros e uma ordenação estável](#preserve-filters-and-stable-ordering)
- [Recuse entrada malformada e posições impossíveis](#reject-malformed-input-and-impossible-positions)
- [Teste limites e alterações concorrentes](#test-boundaries-and-concurrent-change)
- [Três formatos de operações terminais](#three-terminal-shapes)
- [Paginação faz parte do contrato de resposta](#pagination-is-part-of-the-response-contract)

<a name="a-paginator-answers-navigation-not-only-rows"></a>
## Um paginador fornece navegação, não apenas registros

Um slice sozinho não distingue última página cheia de página cheia com mais dados. Terminais de paginação retornam itens tipados e estado para navegação estável, executando sob o mesmo Grant e tenant de Get.

<a name="choose-length-aware-simple-or-cursor-pagination"></a>
## Escolha paginação com total, simples ou por cursor

Paginate executa a consulta de dados e uma contagem para calcular o total de páginas. `SimplePaginate` dispensa a contagem total e apenas verifica se existe outra página. `CursorPaginate` usa valores ordenados codificados em um cursor opaco, evitando grandes deslocamentos em resultados extensos ou que mudam frequentemente.

```go
invoices, page, err := models.Invoices(db).
	OrderBy("created_at").
	Paginate(ctx, g, 25, requestedPage, pagination.Options{})
if err != nil { return nil, nil, err }
return invoices, page, nil
```

<a name="preserve-filters-and-stable-ordering"></a>
## Preserve os filtros e uma ordenação estável

Componha busca e filtros independentes de tenant antes do terminal, valide tamanho de página na requisição e copie parâmetros aprovados aos links. Cursor exige ordem determinística com desempate único para timestamps iguais não duplicarem nem ocultarem registros entre requisições.

<a name="reject-malformed-input-and-impossible-positions"></a>
## Recuse entrada malformada e posições impossíveis

Os terminais de paginação publicados normalizam página abaixo de um para 1. A validação HTTP ainda pode rejeitar esse valor conforme seu contrato, mas o terminal não retorna erro de página negativa. Tamanho ilimitado, cursores malformados ou indecodificáveis são responsabilidade da requisição. Página além do fim segue o contrato documentado de vazio ou não encontrado.

<a name="test-boundaries-and-concurrent-change"></a>
## Teste limites e alterações concorrentes

Teste páginas abaixo de um e sua normalização para 1; depois resultados vazios, uma página exata, última página cheia, um registro extra, tamanho máximo permitido e outro tenant com mesma ordenação. Para cursores, insira registro entre requisições e confirme a continuação documentada, sem confiar no texto interno do cursor.

<a name="three-terminal-shapes"></a>
## Três formatos de operações terminais

Paginate retorna os registros e uma `*pagination.LengthAwarePage` com totais; `SimplePaginate` retorna os registros e uma `*pagination.Page` focada em anterior/próxima; `CursorPaginate` recebe `pagination.Cursor` interpretado e retorna os registros e uma `*pagination.CursorPage`. Os três recebem contexto, Grant, tamanho de página e opções após compor a consulta.

<a name="pagination-is-part-of-the-response-contract"></a>
## Paginação faz parte do contrato de resposta

Os recursos HTTP decidem como serializar itens, totais e links; a página apenas fornece esses dados. A validação da consulta define filtros e colunas de ordenação permitidos. Use Chunk ou Cursor para processamento em lotes no servidor sem navegação para o usuário: um objeto de página é desnecessário quando cada registro será processado uma única vez.
