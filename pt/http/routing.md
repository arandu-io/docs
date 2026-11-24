# Rotas

- [Rotas são contratos executáveis](#routes-are-executable-contracts)
- [Registre e gere o caminho de uma rota](#register-and-reverse-one-route)
- [Grupos, restrições e raiz exata](#groups-constraints-and-the-exact-root)
- [Falhe durante a construção, não após um clique](#fail-while-building-not-after-a-click)
- [Inspecione e despache a mesma tabela](#inspect-and-dispatch-the-same-table)
- [API de roteamento publicada](#released-routing-api)
- [Continue pelo caminho da requisição](#continue-the-request-path)

<a name="routes-are-executable-contracts"></a>
## Rotas são contratos executáveis

O router associa método HTTP e padrão de caminho Go 1.22 a um handler padrão ou ação de controller. Nomear a rota no registro dá a templates, redirecionamentos, testes e CLI uma fonte única de endereço. Grupos adicionam prefixo e middleware compartilhando a tabela de rotas pai.

<a name="register-and-reverse-one-route"></a>
## Registre e gere o caminho de uma rota

Action adapta `func(*http.Context)` error; Get, Post, Put, Patch e Delete aceitam handlers `net/http`. Nomeie a rota antes do tráfego e use a tabela de rotas para gerar o caminho com a contagem exata de parâmetros.

```go example
func Register(router *fhttp.Router, sessions *security.SessionStore, show func(*fhttp.Context) error) (string, error) {
	router.Action("GET", "/invoices/{invoice}", show,
		middleware.RequireAuth(sessions),
	).Name("invoices.show")

	path, err := router.Table().URL("invoices.show", "inv_123")
	if err != nil {
		return "", err
	}
	return path, nil
}
```

<a name="groups-constraints-and-the-exact-root"></a>
## Grupos, restrições e raiz exata

Grupos compartilham um prefixo e uma fronteira de política. Restrições recusam identificadores malformados antes da action e o padrão de raiz exata impede que uma rota inicial capture toda subárvore na gramática do `ServeMux` padrão. Um handler `ConnectRPC` pode ser montado com Post em seu caminho de serviço gerado; o mesmo listener TLS serve REST por HTTP/1 e Connect, `gRPC` ou gRPC-Web por HTTP/2.

- Use `Group(prefix, middleware...)` para caminho e policy compartilhados.
- Use `WhereNumber`, `WhereUuid`, `WhereUlid` ou `WhereIn` para formato de identificadores na rota.
- Use `/{name}` para um segmento e `/{path...}` só quando a ação possuir o restante.
- Use `/{$}` para a raiz exata da aplicação.
- Para RPC, configure explicitamente limites de bytes, concorrência, tenant, duração, timeout do servidor, origens e cabeçalhos antes de montar o handler.

<a name="fail-while-building-not-after-a-click"></a>
## Falhe durante a construção, não após um clique

Nomes de rota desconhecidos e a contagem errada de parâmetros fazem `Routes.URL` retornar erro. Divergências de restrição são respostas 404, não resultados de autorização. Registros duplicados ou conflitantes são defeitos de boot, e a ordem de registro precisa permanecer determinística. O boot também recusa uma rota da aplicação que responde a um caminho que o framework monta sob `/_arandu/`, como `/_arandu/health`, num método que a rota do framework não aceita, nomeando o módulo e o padrão. A causa habitual é `r.Post("/", ...)`, cujo padrão corresponde a todo caminho abaixo da raiz; registre a raiz como "`/{$}`". A autenticação RPC roda antes da decodificação da requisição, mas apenas identifica um Subject: os serviços da aplicação ainda autorizam e emitem Grants. A admissão fecha durante o desligamento, os streams observam o cancelamento, e o framework não abre um listener h2c nem um segundo listener RPC.

> [!NOTE]
> `Route.SetURI` altera metadados após o registro, mas não move o padrão instalado no `ServeMux`. Defina a URI final antes do boot.

<a name="inspect-and-dispatch-the-same-table"></a>
## Inspecione e despache a mesma tabela

Execute `aru route:list` e despache a mesma tabela pelo httptest. Cubra um identificador válido, incompatibilidade de restrição, método não suportado, a raiz exata e a URL nomeada. Para RPC, exercite chamadas unárias, streaming do cliente, streaming do servidor e bidirecionais por Connect, `gRPC` e gRPC-Web; comprove também recusa de CORS, admissão global e por tenant, limites de corpo, cancelamento e desligamento gracioso no listener compartilhado.

- O inventário nomeia método, padrão, rota e módulo.
- Parâmetro válido dentro da restrição chega à ação.
- Parâmetro restrito malformado retorna 404.
- A raiz exata não captura um caminho desconhecido.
- A autorização RPC e as quatro formas de streaming executam no listener da aplicação.

<a name="released-routing-api"></a>
## API de roteamento publicada

`http.NewRouter` e `Router.Action` criam e preenchem a tabela. Métodos de verbo de Router registram handlers padrão; `Router.Group` compõe prefixos e middleware; `Route.Name` e `Route.Where*` adicionam metadados estáveis; `Routes.URL` constrói um caminho ou relata incompatibilidade. `rpc.NewHandler` retorna o caminho de serviço gerado e o handler para Factory, Authenticator bearer e Config explícitos. `rpc.Subject` lê a identidade autenticada em um procedimento; não é um resultado de autorização.

| Superfície | Contrato publicado |
| --- | --- |
| `Router.Action` | registra uma ação de controller |
| `Get/Post/Put/Patch/Delete` | registra um handler padrão |
| Group e Name | compõe prefixos e endereços estáveis |
| `Where*` e `Routes.URL` | restringe parâmetros e gera caminhos reversos |
| `rpc.NewHandler` | constrói um handler Connect, `gRPC` e gRPC-Web com limites explícitos |
| `rpc.Subject` | lê o Subject autenticado antes que a policy da aplicação emita um Grant |

<a name="continue-the-request-path"></a>
## Continue pelo caminho da requisição

REST permanece a fronteira padrão de recursos com requisição e resposta. `gRPC` acrescenta chamadas tipadas unárias e por streaming no mesmo servidor HTTP; não é um segundo processo nem um atalho de autorização. Controllers cuidam da tradução de transporte, middleware cuida da política transversal da requisição e assinatura de URL cuida da integridade depois da seleção da rota.

- [Padrões de comunicação](/docs/{{version}}/architecture/contracts-adapters)
- [Controllers](/docs/{{version}}/http/controllers)
- [Middleware](/docs/{{version}}/http/middleware)
- [URLs](/docs/{{version}}/http/urls)
