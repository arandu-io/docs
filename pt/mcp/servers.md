# Servidores MCP

- [Exponha recursos da aplicação por meio das policies existentes](#expose-application-capabilities-through-existing-policies)
- [Instale o módulo de servidor e escolha a origem da identidade](#install-the-server-module-and-choose-an-identity-source)
- [Declare uma ferramenta protegida por Policy](#declare-one-policy-backed-tool)
- [Monte Web ou sirva com uma identidade local declarada](#mount-web-or-serve-a-declared-local-subject)
- [Retorne falhas de protocolo que indiquem como agir](#return-actionable-protocol-failures)
- [Mantenha o modelo sujeito ao mesmo limite de Grant](#keep-the-model-behind-the-same-grant-boundary)
- [Teste transporte, schema e recusa entre tenants](#exercise-transport-schema-and-tenant-refusal)
- [Separe a exposição do servidor da automação de clientes](#separate-server-exposure-from-client-automation)

> [!WARNING]
> **Limitações atuais.** MCP :mcpVersion publica os contratos de servidor, ferramenta, recurso, prompt e destino autorizado, mas nenhuma aplicação de referência atual com tag o consome. Verifique a integração com o grafo atual de Framework :frameworkVersion e Hesape :hesapeVersion; um runtime geral de cliente continua adiado.

<a name="expose-application-capabilities-through-existing-policies"></a>
## Exponha recursos da aplicação por meio das policies existentes

O pacote MCP implementa o lado servidor da revisão 2024-11-05 do protocolo. Ele expõe ferramentas, recursos e prompts por HTTP ou stdio com mensagens delimitadas por quebra de linha. Não inclui cliente genérico, sampling, roots, ciclo de agente ou runtime de IA dentro do framework.

<a name="install-the-server-module-and-choose-an-identity-source"></a>
## Instale o módulo de servidor e escolha a origem da identidade

Fixe `github.com/arandu-io/mcp@:mcpVersion` e prepare serviços da aplicação cujos métodos recebam `security.Subject`. Use Web para identidade de sessão ou Local e Start com um Subject configurado explicitamente. Chame `Server.Validate` durante a inicialização.

- Componha o servidor uma vez num provider em app/Providers.
- Use o armazenamento de sessões publicado para obter a identidade HTTP.
- Declare um Subject revisado para o stdio local.
- Use `mcp.Describe` para consultar o inventário.

<a name="declare-one-policy-backed-tool"></a>
## Declare uma ferramenta protegida por Policy

Uma Tool declara Name, Description, Schema e Handle. Os argumentos são validados pelo schema antes de Handle, `Request.Subject` é somente leitura, e o handler chama o mesmo serviço protegido por Policy usado pelo HTTP, em vez de acessar diretamente um repositório. `mcp.JSON` recebe um recurso que lista os seus campos, o mesmo documento que `Context.JSON` escreve, e não um valor solto. `mcp.JSON` recebe um recurso que lista os seus campos, o mesmo documento que `Context.JSON` escreve, e não um valor solto.

```go example
type Post struct{ ID, Title string }

// PostList is the resource a tool answers with: the fields somebody listed.
type PostList []Post

func (l PostList) ToArray() map[string]any {
	items := make([]map[string]any, 0, len(l))
	for _, post := range l {
		items = append(items, map[string]any{"id": post.ID, "title": post.Title})
	}
	return map[string]any{"posts": items}
}

func (PostList) With() map[string]any { return nil }

type PostService struct{}

func (*PostService) List(ctx context.Context, subject security.Subject, status string) ([]Post, error) {
	return nil, nil
}

type ListPosts struct{ service *PostService }

func (ListPosts) Name() string { return "list_posts" }
func (ListPosts) Description() string { return "Lists posts visible to the current subject." }
func (ListPosts) Schema() mcp.Schema {
	return mcp.Object(mcp.String("status", "Status to include").Enum("published", "draft").Required())
}
func (t ListPosts) Handle(ctx context.Context, r mcp.Request) (mcp.Response, error) {
	status, _ := r.String("status")
	posts, err := t.service.List(ctx, r.Subject(), status)
	if err != nil {
		return mcp.Response{}, err
	}
	return mcp.JSON(PostList(posts)), nil
}
```

<a name="mount-web-or-serve-a-declared-local-subject"></a>
## Monte Web ou sirva com uma identidade local declarada

`aru make:mcp-tool`, `make:mcp-resource` e `make:mcp-prompt` escrevem uma tool, um resource ou um prompt sobre um serviço e imprimem a ligação do provider e da rota. Crie o Server com metadados e fatias explícitas de tools, resources e prompts. Monte `mcp.Web` em POST `/mcp` para um Subject vindo da sessão, ou use Start e Local com uma identidade declarada. Use Text, JSON, Error e Describe para resposta e inventário.

- Valide antes de registrar Web.
- Trate `Guest(tenant)` explicitamente em requisições Web não autenticadas.
- Nunca deduza a identidade de Local.
- Reserve stdout exclusivamente para JSON-RPC no modo stdio.

<a name="return-actionable-protocol-failures"></a>
## Retorne falhas de protocolo que indiquem como agir

Falhas de ferramentas retornam conteúdo com `isError`, uma recusa de autorização nunca é um sucesso vazio, ferramentas desconhecidas informam as opções disponíveis, JSON-RPC e argumentos malformados são erros de protocolo, notificações não recebem resposta e mensagens acima de um `MiB` são recusadas. resources/read usa o tipo MIME declarado pelo recurso. `Prompt.Argument.Required` é anunciado, mas não é imposto por prompts/get; por isso, `Prompt.Render` deve validar os argumentos obrigatórios.

> [!WARNING]
> A Aru :aruVersion tem duas superfícies MCP que não são o runtime deste servidor: `aru mcp` serve o próprio servidor da CLI, que responde a um assistente de código sobre o projeto, e os geradores `make:mcp-*` escrevem as tools, os resources e os prompts de uma aplicação. Ela não tem comando `mcp:start` nem `mcp:list` para o servidor da aplicação, que é montado como uma rota. O MCP :mcpVersion não tem consumidor de aplicação publicado atual, então esta página fica limitada até que a integração seja revalidada contra o Framework :frameworkVersion e o Hesape :hesapeVersion.

<a name="keep-the-model-behind-the-same-grant-boundary"></a>
## Mantenha o modelo sujeito ao mesmo limite de Grant

> [!WARNING]
> Toda ferramenta chama um Service ou uma Policy da aplicação com `Request.Subject`. Nunca aceite Subject, tenant, papel ou Grant pelos argumentos. Revise o Subject local explícito como infraestrutura privilegiada e mantenha descrições e instruções livres de segredos.

O schema valida o formato da entrada, não a autoridade de quem chama. A descrição do servidor também é uma superfície de ataque voltada ao modelo e deve ser específica e limitada.

<a name="exercise-transport-schema-and-tenant-refusal"></a>
## Teste transporte, schema e recusa entre tenants

Valide o servidor, inicialize-o e confira o protocolo e suas capacidades. Chame ferramentas com entradas válidas e malformadas, prove que Handle não é executado quando o schema falha, teste dois tenants e Guest, erro de autorização, limite de um `MiB`, ausência de resposta a notificações e saída stdio sem conteúdo extra.

- A revisão do protocolo é 2024-11-05.
- Erros de schema não chamam Handle.
- O comportamento de MIME na leitura de recursos e a limitação de argumentos obrigatórios de prompts estão fixados.
- Entradas de transporte acima do limite são rejeitadas.

<a name="separate-server-exposure-from-client-automation"></a>
## Separe a exposição do servidor da automação de clientes

Autorização e sessões definem a identidade; serviços contêm operações reutilizáveis protegidas por policies; desenvolvimento assistido cobre apenas ferramentas do editor. O runtime de IA e a orquestração de clientes continuam adiados.

- [Autorização](/docs/{{version}}/security/authorization)
- [Sessão](/docs/{{version}}/http/session)
- [Contratos e adapters](/docs/{{version}}/architecture/contracts-adapters)
- [API do servidor MCP](https://pkg.go.dev/github.com/arandu-io/mcp)
