# ADR 0030-mcp-e-um-repositorio — a aplicação exposta a um assistente passa pela policy

**Status:** Aceita.

## Contexto

O proprietário autorizou em 09/08/2026, pedindo o equivalente do `laravel/mcp`.
A REGRA 3 exige decisão explícita para cada repositório, e este é o décimo
terceiro.

O Model Context Protocol é como um assistente alcança um programa: o programa
declara ferramentas que ele pode chamar, recursos que pode ler e prompts que
pode usar. O `laravel/mcp` foi rastreado na documentação — `Server`, `Tool`,
`Resource`, `Prompt`, `make:mcp-*`, `routes/ai.php`, `Mcp::web` e `Mcp::local`.

## Decisão

**`github.com/arandu-io/mcp`, importável, com a forma do `laravel/mcp` e uma
diferença que é o motivo dele existir.**

**Ferramenta alcança dado, e todo caminho até dado neste framework carrega um
`security.Grant`.** O `Subject` está no `Request` e a ferramenta não tem outro
jeito de chamar um serviço:

```go
func (t ListPosts) Handle(ctx context.Context, r mcp.Request) (mcp.Response, error) {
	found, err := t.svc.List(ctx, r.Subject(), data.Query{Limit: 20})
	…
}
```

Isso não é asseio. Um servidor MCP entrega a um modelo de linguagem as chaves da
aplicação, e a versão em que uma ferramenta consulta o banco direto seria o maior
buraco que este projeto poderia publicar — e publicaria **em silêncio**, porque as
respostas pareceriam certas. A policy que recusa uma ferramenta recusa pelo mesmo
motivo que recusa um controller.

### De onde vem o sujeito, e por que os dois transportes diferem

| | |
|---|---|
| `mcp.Web` | da sessão, exatamente como qualquer requisição HTTP. Sem sessão é convidado, e o que um convidado pode é resposta da policy (ADR 0029) |
| `mcp.Local` | de configuração, sobre stdio. **Não há sessão num pipe**, então a identidade é declarada onde o servidor é registrado e fica visível em `routes/ai.go` |

O local **recebe** um `Subject` em vez de assumir um. Uma aplicação que deixa um
assistente agir como administrador escreveu isso num lugar que um revisor lê.

### Três decisões menores, cada uma com teste

- **Recusa é erro, nunca resultado vazio.** Modelo que recebe lista vazia conclui
  que não há nada e conta isso a alguém. Modelo que ouve que não pode, para.
- **Argumento não declarado é recusado.** Modelo que inventa um e não é avisado
  continua inventando, e a ferramenta age sobre uma chamada que entendeu pela
  metade.
- **Ferramenta sem descrição não sobe.** A descrição é o que o modelo lê para
  decidir se chama — é a string de maior alavancagem do pacote, e sem ela a
  ferramenta é chamada a esmo. É erro de declaração, então é no boot.

### O que não existe de propósito

Sem atributos, porque Go não tem: nome e descrição são métodos. Sem facade e sem
registro dinâmico — um servidor declara suas ferramentas numa fatia, então
ferramenta que existe e não é alcançável aparece num arquivo. Sem sampling e sem
roots: capacidade declarada e não servida é capacidade que o cliente reporta como
servidor quebrado.

## Consequências

- **Treze repositórios.** O `plans/checklist.sh` passa a contar treze.
- Uma dependência, e é o framework — a CI recusa a segunda.
- Nada é escrito em stdout além de uma resposta: linha de log ali é erro de
  parse no cliente, e é a forma mais comum de um servidor stdio parecer quebrado
  funcionando.
- Falta o `aru make:mcp-server`, `make:mcp-tool` e `mcp:start`, e o exemplo no
  `examples`. O pacote existe antes do gerador de propósito: gerador escrito
  contra API que ainda se move é trabalho feito duas vezes.
