# Cache e Redis

- [A aplicação usa um repositório com escopo de tenant](#the-application-uses-a-tenant-scoped-repository)
- [Comece com armazenamento dentro do processo](#begin-with-the-in-process-store)
- [Armazene, memorize, conte e invalide](#store-remember-count-and-invalidate)
- [Ausência no cache é esperada; erro do armazenamento, não](#a-cache-miss-is-expected-a-store-error-is-not)
- [Teste isolamento e expiração pelo repositório](#test-isolation-and-expiry-through-the-repository)
- [Fronteiras do repositório e armazenamento](#repository-and-store-boundaries)
- [Servidor compatível com o protocolo é escolha de implantação](#a-protocol-compatible-server-is-a-deployment-choice)

> [!WARNING]
> **Limitações atuais.** `foundation.NewDrivers` constrói `drivers.Store`, o `cache.Repository` com escopo de tenant sobre o store que `CACHE_STORE` nomeia, e o esqueleto publicado entrega o store subjacente aos códigos de uso único e ao segundo fator em vez de fazer cache pelo repositório. Teste o namespace baseado em Grant na aplicação que faz cache por ele.

<a name="the-application-uses-a-tenant-scoped-repository"></a>
## A aplicação usa um repositório com escopo de tenant

Repository de cache adiciona JSON, expiração, eventos e chaves prefixadas por tenant sobre Store. A aplicação mantém o repositório, não o backend bruto. Toda operação comum recebe Grant, portanto o tenant vem da autorização, não de campo da requisição nem prefixo manual.

<a name="begin-with-the-in-process-store"></a>
## Comece com armazenamento dentro do processo

O armazenamento array fornece a desenvolvimento e testes unitários um backend com a mesma superfície Repository. Construa-o uma vez; use Put e recupere com o mesmo Grant. Outro tenant gera uma chave física diferente mesmo quando o texto da chave da aplicação é igual.

```go example
type Dashboard struct{ Total int }

func CacheDashboard(ctx context.Context, g security.Grant, value Dashboard) (Dashboard, error) {
	repository := cache.New(cache.NewArrayStore())
	if err := repository.Put(ctx, g, "dashboard", value, 5*time.Minute); err != nil {
		return Dashboard{}, err
	}
	return cache.Get[Dashboard](ctx, repository, g, "dashboard")
}
```

<a name="store-remember-count-and-invalidate"></a>
## Armazene, memorize, conte e invalide

Put substitui um valor com duração limitada; Forever é a escolha explícita sem prazo; Has, Missing, Forget e Touch gerenciam presença; Increment e Decrement atualizam contadores. Get e Remember genéricos decodificam no tipo solicitado, deixando asserções de tipo fora dos handlers.

Numa aplicação o repositório já está construído: `foundation.NewDrivers` responde `drivers.Store`, o `cache.Repository` sobre o store que `CACHE_STORE` nomeia, e `drivers.Cache` guarda todos os stores por nome. memory sempre existe; redis existe depois que o conector `hesape/redis` é vinculado em `bootstrap/app.go` e `REDIS_URL` nomeia o endpoint. Um provider entrega `drivers.Store` aos serviços que fazem cache.

<a name="a-cache-miss-is-expected-a-store-error-is-not"></a>
## Ausência no cache é esperada; erro do armazenamento, não

Trate a ausência conforme o pacote e calcule o valor de origem uma vez. Propague erros de conexão, codificação e decodificação sem tratá-los como cache miss: reconstruir silenciosamente durante uma queda pode multiplicar consultas até derrubar o serviço principal. Nunca armazene resultados falhos ou parcialmente autorizados.

Um Grant somente leitura, emitido enquanto um operador apenas visualiza a conta de outra pessoa, lê e nunca escreve: todo método que guarda, remove, esvazia ou trava uma entrada o recusa com `auth.ErrReadOnly`, e Remember e Flexible respondem o que calculam sem guardar. Pergunte a `Grant.ReadOnly` para deixar uma escrita fora de uma página assim, em vez de deixá-la falhar.

<a name="test-isolation-and-expiry-through-the-repository"></a>
## Teste isolamento e expiração pelo repositório

Escreva a mesma chave lógica com Grants de dois tenants e confirme que cada um lê apenas seu valor. Avance ou limite o tempo para comprovar expiração, verifique que Remember chama o loader só na ausência e execute o contrato compartilhado contra o driver de rede antes de escolhê-lo em produção.

<a name="repository-and-store-boundaries"></a>
## Fronteiras do repositório e armazenamento

Métodos de Repository cobrem valores, lotes, contadores, namespaces, tags quando suportadas e locks contra trabalho sobreposto. Store cobre persistência física. `WithoutOverlapping` recebe Grant e prefixa o lock por tenant como outras operações. `FlushLocks` é a exceção global: não recebe Grant porque limpa o namespace de locks do armazenamento. Isso não transforma stores brutos em API da aplicação.

<a name="a-protocol-compatible-server-is-a-deployment-choice"></a>
## Servidor compatível com o protocolo é escolha de implantação

O driver de rede fala o protocolo comum de servidores compatíveis com Redis sem incluir comandos específicos de produto no contrato de cache. Escolha conexão e store na configuração, mantenha chamadas da aplicação e dimensione TTL pela atualidade necessária dos dados de origem. Registros duráveis e sessões não viram cache só porque o backend retém bytes.
