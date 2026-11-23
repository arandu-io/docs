# Tipos e conversões do banco de dados

- [Os campos Go definem os tipos de dados](#go-fields-are-the-primary-data-types)
- [Declare intenção de armazenamento na struct](#declare-storage-intent-in-the-struct)
- [Respeite os limites de atribuição da versão publicada](#stay-inside-the-released-assignment-surface)
- [A conversão em runtime ainda pode falhar](#runtime-conversion-can-still-fail)
- [Teste ida e volta dos valores-limite em cada driver](#round-trip-boundary-values-on-each-driver)
- [Tags e limites atuais de atribuição](#tags-and-current-assignment-limits)
- [Tipos de armazenamento e resposta podem diferir](#storage-types-and-response-types-may-differ)

<a name="go-fields-are-the-primary-data-types"></a>
## Os campos Go definem os tipos de dados

A hidratação do model preenche campos declarados na struct e não possui um registro geral de casts. O caminho de atribuição publicado cobre valores diretamente atribuíveis, ponteiros, `time.Time`, strings e slices de bytes, booleanos, conversões numéricas por reflexão e campos cujo endereço implementa `sql.Scanner`. Essa decodificação é explícita e pertence ao tipo; não promete converter qualquer representação armazenada em um tipo de domínio.

<a name="declare-storage-intent-in-the-struct"></a>
## Declare intenção de armazenamento na struct

Use o menor tipo que represente o domínio e comporte o intervalo do banco. Ponteiros distinguem SQL NULL do valor zero. Armazene dinheiro como inteiro da menor unidade quando exigir precisão fixa; use tags db só para mapear nomes ou excluir campos.

```go example
type Invoice struct {
	model.Model

	ID        string     `db:"id"`
	Total     int64      `db:"total"`
	DueDate   time.Time  `db:"due_date"`
	DeletedAt *time.Time `db:"deleted_at"`
}
```

<a name="stay-inside-the-released-assignment-surface"></a>
## Respeite os limites de atribuição da versão publicada

Use ponteiros para campos anuláveis apenas quando o conector retornar um valor compatível. Mantenha datas, strings, bytes, booleanos e números nas representações testadas do driver. Na :hesapeVersion, implemente `sql.Scanner` no ponteiro do campo para controlar a hidratação e `driver.Valuer` no valor para controlar os parâmetros de escrita. JSON não é decodificado automaticamente; um tipo de domínio pode implementar ambas as interfaces para cuidar dessa codificação.

<a name="runtime-conversion-can-still-fail"></a>
## A conversão em runtime ainda pode falhar

Origem e destino incompatíveis geram erro de atribuição, mas o caminho numérico usa conversão por reflexão; este guia não afirma que ela detecta todo overflow. NULL, datas e representações específicas do driver precisam ser medidos no conector de produção. Valide os intervalos do domínio antes de escrever e não trate um valor zero decodificado como prova de validade dos dados de origem.

<a name="round-trip-boundary-values-on-each-driver"></a>
## Teste ida e volta dos valores-limite em cada driver

Teste os valores numéricos mínimos e máximos aceitos, NULL e não NULL, fusos e precisão de datas, strings e slices de bytes em cada dialeto de produção, pois os drivers retornam representações Go diferentes. Teste explicitamente os limites de overflow, sem supor que a reflexão os rejeite. Em campos de domínio personalizados, teste o ciclo de Scanner e Valuer, incluindo JSON malformado, valores desconhecidos de enum e erros de vinculação.

<a name="tags-and-current-assignment-limits"></a>
## Tags e limites atuais de atribuição

Tag db nomeia coluna; traço exclui campo; ausência usa conversão de nomes do framework. Configuração do model escolhe tabela e comportamento em Go, sem anotações alheias. Em :hesapeVersion, hidratação respeita `sql.Scanner` e escrita resolve `driver.Valuer`; decodificação JSON automática e registro nomeado de casts não existem.

<a name="storage-types-and-response-types-may-differ"></a>
## Tipos de armazenamento e resposta podem diferir

Um inteiro em centavos no banco pode virar uma string monetária formatada na view, mas o model não deve armazenar apresentação. Da mesma forma, colunas criptografadas precisam ser decifradas em um limite deliberado, não por um cast que revele texto puro sempre que o model for carregado. Os recursos HTTP controlam a representação pública.
