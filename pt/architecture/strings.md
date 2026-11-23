# Strings

- [Biblioteca padrão primeiro](#standard-library-first)
- [Convenções de nomes compartilhadas](#shared-naming-conventions)
- [Limites e slugs vazios](#boundaries-and-empty-slugs)

<a name="standard-library-first"></a>
## Biblioteca padrão primeiro

Use a biblioteca padrão para operações comuns com strings. strings fornece busca, divisão, substituição, remoção de espaços, comparação e conversão de maiúsculas/minúsculas sem acrescentar convenções da aplicação. strconv interpreta e formata valores primitivos; unicode e utf8 cuidam de caracteres e codificação. Manter essas chamadas visíveis facilita reconhecer o comportamento familiar do Go.

Use str quando vários pacotes precisarem concordar com uma convenção de nomes que a biblioteca padrão não define: converter um campo em rótulo, derivar um identificador no estilo Go, flexionar um substantivo inglês, converter texto latino suportado em ASCII ou produzir um slug de endereço. Uma transformação usada por um único recurso fica com esse recurso, em vez de se tornar um auxiliar global.

<a name="shared-naming-conventions"></a>
## Convenções de nomes compartilhadas

Studly e Camel constroem identificadores a partir de palavras separadas por hífen, sublinhado ou espaço. Snake e Kebab tornam explícita a conversão oposta. Headline transforma um valor semelhante a um campo em rótulo; o framework :frameworkVersion usa exatamente essa função quando a validação torna o nome de um campo legível. Assim, Password Confirmation é a saída compartilhada, não um texto escolhido separadamente por cada formulário.

Slug recebe seu separador explicitamente porque Go não tem argumentos padrão. Funções de pacote são a forma mais clara para uma operação. `str.Of` está disponível quando uma sequência de transformações descreve um valor, e `ToString` retorna a string final; os dois caminhos executam o mesmo comportamento de pacote.

```go example
func NamingExamples() (string, string, string) {
	return str.Studly("invoice_item"), // InvoiceItem
		str.Headline("password_confirmation"), // Password Confirmation
		str.Slug("Quarterly Report 2026", "-") // quarterly-report-2026
}
```

> [!NOTE]
> Mantenha o valor original como fonte de verdade. Um identificador derivado pode ser produzido novamente; armazenar várias variantes de capitalização cria oportunidades independentes de divergência.

<a name="boundaries-and-empty-slugs"></a>
## Limites e slugs vazios

ASCII usa uma tabela de conversão documentada, não um transliterador geral. Ela cobre Latin-1 Supplement e Latin Extended-A; sistemas de escrita não suportados são descartados, não adivinhados. Slug começa por essa conversão, portanto um título escrito apenas com caracteres não suportados pode gerar um slug vazio. Confira o resultado e use uma chave estável emitida separadamente quando houver possibilidade de saída vazia.

Um slug é apresentação, não identidade, autorização ou lista de permissão do banco. Identificadores de tenant, chaves de armazenamento, parâmetros de rota e escolhas de coluna ainda precisam ser validados pelo subsistema que os consome. Frases em linguagem humana pertencem aos catálogos de tradução: mudar limites de palavras ou capitalização não traduz significado.
