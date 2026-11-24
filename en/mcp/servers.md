# MCP servers

- [Expose application capabilities through existing policies](#expose-application-capabilities-through-existing-policies)
- [Install the server module and choose an identity source](#install-the-server-module-and-choose-an-identity-source)
- [Declare one policy-backed tool](#declare-one-policy-backed-tool)
- [Mount Web or serve a declared local subject](#mount-web-or-serve-a-declared-local-subject)
- [Return actionable protocol failures](#return-actionable-protocol-failures)
- [Keep the model behind the same Grant boundary](#keep-the-model-behind-the-same-grant-boundary)
- [Exercise transport, schema, and tenant refusal](#exercise-transport-schema-and-tenant-refusal)
- [Separate server exposure from client automation](#separate-server-exposure-from-client-automation)

<a name="expose-application-capabilities-through-existing-policies"></a>
## Expose application capabilities through existing policies

The MCP package is the server side of protocol revision 2024-11-05. It exposes tools, resources, and prompts over HTTP or newline-framed stdio. It does not ship a general client, sampling, roots, an agent loop, or an in-framework AI runtime.

<a name="install-the-server-module-and-choose-an-identity-source"></a>
## Install the server module and choose an identity source

Pin `github.com/arandu-io/mcp@:mcpVersion`, prepare application services whose methods accept `security.Subject`, use Web for session identity or Local and Start with one explicit configured Subject, and call `Server.Validate` during boot.

- Compose the server once in a provider in app/Providers.
- Use the released session store for HTTP identity.
- Declare a reviewed Subject for local stdio.
- Use `mcp.Describe` for inventory.

<a name="declare-one-policy-backed-tool"></a>
## Declare one policy-backed tool

A Tool declares Name, Description, Schema, and Handle. Arguments are schema-validated before Handle, `Request.Subject` is read-only, and the handler calls the same policy-backed service used by HTTP rather than reaching a repository directly. `mcp.JSON` takes a resource that lists its fields, the same document `Context.JSON` writes, and not a bare value.

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
## Mount Web or serve a declared local subject

`aru make:mcp-tool`, `make:mcp-resource`, and `make:mcp-prompt` write one tool, resource, or prompt over a service and print the provider and route wiring. Create Server with metadata and explicit tool, resource, and prompt slices. Mount `mcp.Web` on POST `/mcp` for session-derived Subject, or use Start and Local with declared identity. Use Text, JSON, Error, and Describe for response and inventory.

- Validate before registering Web.
- Treat `Guest(tenant)` deliberately on unauthenticated Web requests.
- Never infer Local identity.
- Keep stdout exclusively JSON-RPC in stdio mode.

<a name="return-actionable-protocol-failures"></a>
## Return actionable protocol failures

Tool failures return content with `isError`, authorization refusal is never an empty success, unknown tools name available choices, malformed JSON-RPC and arguments are protocol errors, notifications have no response, and messages above one `MiB` are refused. resources/read uses the MIME type declared by the resource. `prompts/get` refuses a call that omits a `Prompt.Argument.Required` argument or names one the prompt did not declare, before `Prompt.Render` runs.

Aru :aruVersion has two MCP surfaces that are not this server's runtime: `aru mcp` serves the CLI's own server, which answers a coding assistant about the project, and the `make:mcp-*` generators write an application's tools, resources, and prompts. It has no `mcp:start` or `mcp:list` command for the application's server, which is mounted as a route.

<a name="keep-the-model-behind-the-same-grant-boundary"></a>
## Keep the model behind the same Grant boundary

> [!WARNING]
> Every tool calls an application Service or Policy with `Request.Subject`. Never accept Subject, tenant, role, or Grant from arguments. Review the explicit local Subject as privileged infrastructure and keep descriptions and instructions free of secrets.

Schema validates input shape, not authority. The server-facing description itself is model-facing attack surface and must remain specific and bounded.

<a name="exercise-transport-schema-and-tenant-refusal"></a>
## Exercise transport, schema, and tenant refusal

Validate the server, initialize and assert protocol and capabilities, call valid and malformed tools, prove Handle is skipped on schema error, exercise two tenants and Guest, test authorization error, one-MiB ceiling, notification silence, and clean stdio output.

- Protocol revision is 2024-11-05.
- Schema errors do not call Handle.
- Resource read MIME behavior and required prompt arguments are pinned.
- Oversized transport input is bounded.

<a name="separate-server-exposure-from-client-automation"></a>
## Separate server exposure from client automation

Authorization and sessions define identity, services contain reusable policy-backed operations, and assisted development covers editor tooling only. Runtime AI and client orchestration remain deferred.

- [Authorization](/docs/{{version}}/security/authorization)
- [Session](/docs/{{version}}/http/session)
- [Contracts and adapters](/docs/{{version}}/architecture/contracts-adapters)
- [MCP server API](https://pkg.go.dev/github.com/arandu-io/mcp)
