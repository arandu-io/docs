# Servidores MCP

- [Expón funciones de la aplicación mediante las policies existentes](#expose-application-capabilities-through-existing-policies)
- [Instala el módulo de servidor y elige el origen de la identidad](#install-the-server-module-and-choose-an-identity-source)
- [Declara una herramienta protegida por Policy](#declare-one-policy-backed-tool)
- [Monta Web o sirve con una identidad local declarada](#mount-web-or-serve-a-declared-local-subject)
- [Devuelve fallos de protocolo que indiquen cómo actuar](#return-actionable-protocol-failures)
- [Mantén el modelo sujeto al mismo límite de Grant](#keep-the-model-behind-the-same-grant-boundary)
- [Prueba el transporte, el schema y el rechazo entre tenants](#exercise-transport-schema-and-tenant-refusal)
- [Separa la exposición del servidor de la automatización de clientes](#separate-server-exposure-from-client-automation)

> [!WARNING]
> **Limitaciones actuales.** MCP :mcpVersion publica los contratos de servidor, herramienta, recurso, prompt y destino autorizado, pero ninguna aplicación de referencia actual con etiqueta lo consume. Verifica la integración con el grafo actual de Framework :frameworkVersion y Hesape :hesapeVersion; un runtime general de cliente sigue aplazado.

<a name="expose-application-capabilities-through-existing-policies"></a>
## Expón funciones de la aplicación mediante las policies existentes

El paquete MCP implementa el lado servidor de la revisión 2024-11-05 del protocolo. Expone herramientas, recursos y prompts mediante HTTP o stdio con mensajes delimitados por saltos de línea. No incluye cliente genérico, sampling, roots, bucle de agente ni runtime de IA dentro del framework.

<a name="install-the-server-module-and-choose-an-identity-source"></a>
## Instala el módulo de servidor y elige el origen de la identidad

Fija `github.com/arandu-io/mcp@:mcpVersion` y prepara servicios de la aplicación cuyos métodos reciban `security.Subject`. Usa Web para identidad de sesión o Local y Start con un Subject configurado explícitamente. Llama a `Server.Validate` durante el inicio.

- Componga el servidor una vez en un provider en app/Providers.
- Usa el almacenamiento de sesiones publicado para obtener la identidad HTTP.
- Declara un Subject revisado para stdio local.
- Usa `mcp.Describe` para consultar el inventario.

<a name="declare-one-policy-backed-tool"></a>
## Declara una herramienta protegida por Policy

Una Tool declara Name, Description, Schema y Handle. Los argumentos se validan mediante schema antes de Handle, `Request.Subject` es de solo lectura y el handler llama al mismo servicio protegido por Policy que usa HTTP, en vez de acceder directamente a un repositorio. `mcp.JSON` recibe un recurso que enumera sus campos, el mismo documento que escribe `Context.JSON`, y no un valor suelto. `mcp.JSON` recibe un recurso que enumera sus campos, el mismo documento que escribe `Context.JSON`, y no un valor suelto.

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
## Monta Web o sirve con una identidad local declarada

`aru make:mcp-tool`, `make:mcp-resource` y `make:mcp-prompt` escriben una tool, un resource o un prompt sobre un servicio e imprimen el cableado del provider y de la ruta. Cree el Server con metadatos y slices explícitos de tools, resources y prompts. Monte `mcp.Web` en POST `/mcp` para un Subject derivado de la sesión, o use Start y Local con una identidad declarada. Use Text, JSON, Error y Describe para la respuesta y el inventario.

- Valida antes de registrar Web.
- Trata `Guest(tenant)` explícitamente en solicitudes Web no autenticadas.
- Nunca deduzcas la identidad de Local.
- Reserva stdout exclusivamente para JSON-RPC en modo stdio.

<a name="return-actionable-protocol-failures"></a>
## Devuelve fallos de protocolo que indiquen cómo actuar

Los fallos de herramientas devuelven contenido con `isError`, un rechazo de autorización nunca es un éxito vacío, las herramientas desconocidas indican las opciones disponibles, JSON-RPC y argumentos malformados son errores de protocolo, las notificaciones no reciben respuesta y se rechazan los mensajes superiores a un `MiB`. resources/read usa el tipo MIME declarado por el recurso. `Prompt.Argument.Required` se anuncia, pero prompts/get no lo aplica; por ello, `Prompt.Render` debe validar los argumentos obligatorios.

> [!WARNING]
> Aru :aruVersion tiene dos superficies MCP que no son el runtime de este servidor: `aru mcp` sirve el propio servidor de la CLI, que responde a un asistente de código sobre el proyecto, y los generadores `make:mcp-*` escriben las tools, los resources y los prompts de una aplicación. No tiene comando `mcp:start` ni `mcp:list` para el servidor de la aplicación, que se monta como una ruta. MCP :mcpVersion no tiene un consumidor de aplicación publicado actual, así que esta página queda limitada hasta que la integración se revalide contra Framework :frameworkVersion y Hesape :hesapeVersion.

<a name="keep-the-model-behind-the-same-grant-boundary"></a>
## Mantén el modelo sujeto al mismo límite de Grant

> [!WARNING]
> Toda herramienta llama a un Service o una Policy de la aplicación con `Request.Subject`. Nunca aceptes Subject, tenant, rol ni Grant mediante los argumentos. Revisa el Subject local explícito como infraestructura privilegiada y mantén las descripciones e instrucciones libres de secretos.

El schema valida la forma de la entrada, no la autoridad de quien llama. La descripción del servidor también es una superficie de ataque dirigida al modelo y debe ser específica y limitada.

<a name="exercise-transport-schema-and-tenant-refusal"></a>
## Prueba el transporte, el schema y el rechazo entre tenants

Valida el servidor, inicialízalo y comprueba el protocolo y sus capacidades. Llama a herramientas con entradas válidas y malformadas, demuestra que Handle no se ejecuta cuando falla el schema y prueba dos tenants y Guest, errores de autorización, límite de un `MiB`, ausencia de respuesta a notificaciones y salida stdio sin contenido adicional.

- La revisión del protocolo es 2024-11-05.
- Los errores de schema no llaman a Handle.
- El comportamiento MIME de la lectura de recursos y la limitación de argumentos obligatorios de prompts están fijados.
- Se rechazan las entradas de transporte que superan el límite.

<a name="separate-server-exposure-from-client-automation"></a>
## Separa la exposición del servidor de la automatización de clientes

La autorización y las sesiones definen la identidad; los servicios contienen operaciones reutilizables protegidas por policies; el desarrollo asistido cubre solo herramientas del editor. El runtime de IA y la orquestación de clientes siguen aplazados.

- [Autorización](/docs/{{version}}/security/authorization)
- [Sesión](/docs/{{version}}/http/session)
- [Contratos y adapters](/docs/{{version}}/architecture/contracts-adapters)
- [API del servidor MCP](https://pkg.go.dev/github.com/arandu-io/mcp)
