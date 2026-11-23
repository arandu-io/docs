# Serialização de models

- [O JSON padrão é conveniente, mas não define a API](#default-json-is-convenient-and-intentionally-not-the-api)
- [Mapeie um model para uma resposta explícita](#map-a-model-into-an-explicit-response)
- [Represente coleções e relacionamentos opcionais](#represent-collections-and-optional-relations)
- [Pare em valores não suportados ou mapeamento incorreto](#stop-on-unsupported-values-or-broken-mapping)
- [Fixe chaves públicas e ausência de campos sensíveis](#lock-the-public-keys-and-sensitive-absences)
- [DTOs e `JsonResource` têm contratos diferentes](#dtos-and-jsonresource-have-different-contracts)
- [Logs precisam de representação segura própria](#logging-needs-its-own-safe-representation)

> [!WARNING]
> **Limitações atuais.** Representação de model publicada e DTO oficial com campos permitidos existem, mas nenhum fluxo em tag comprova `ToArray` ou `ToJSON` do model como contrato público. Prefira DTO ou `JsonResource` explícito e teste suas chaves exatas.

<a name="default-json-is-convenient-and-intentionally-not-the-api"></a>
## O JSON padrão é conveniente, mas não define a API

`encoding/json` serializa campos Go exportados salvo tag ou marshaler customizado. Esse padrão serve a fronteiras confiáveis, mas é amplo demais para endpoint público: adicionar coluna exportada pode mudar a saída sem editar handler ou contrato de resposta.

<a name="map-a-model-into-an-explicit-response"></a>
## Mapeie um model para uma resposta explícita

Declare os campos prometidos pelo endpoint e construa esse tipo com model já autorizado. A listagem é DTO Go comum serializado por `encoding/json`, não contrato `hhttp.JsonResource` do framework. Reutilize IDs estáveis e valores de domínio omitindo tenant, hashes, caminhos de storage, timestamps internos e relações desnecessárias ao cliente.

```go example
type User struct{ ID, Name string }

type UserResource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewUserResource(user User) UserResource {
	return UserResource{ID: user.ID, Name: user.Name}
}
```

<a name="represent-collections-and-optional-relations"></a>
## Represente coleções e relacionamentos opcionais

Mapeie cada model da coleção e envolva itens com fatos de paginação quando o endpoint for paginado. Inclua recurso relacionado só se a relação foi carregada deliberadamente e o contrato a exige. Use omitempty apenas quando ausência tiver significado definido; ele não substitui a escolha de campos públicos.

<a name="stop-on-unsupported-values-or-broken-mapping"></a>
## Pare em valores não suportados ou mapeamento incorreto

A codificação JSON pode falhar com valores não suportados, e implementações de `MarshalJSON` podem retornar erros. A ausência de um valor obrigatório do domínio deve interromper o mapeamento, não gerar uma resposta enganosa. Se o construtor do recurso falhar, não serialize o model como alternativa: isso reabre o vazamento de campos que o recurso deveria impedir.

<a name="lock-the-public-keys-and-sensitive-absences"></a>
## Fixe chaves públicas e ausência de campos sensíveis

Verifique o formato exato da resposta ou o conjunto de chaves decodificadas e confirme explicitamente a ausência de campos sensíveis do model. Adicione ao model ou à entrada do mapeamento um campo exportado exclusivo da fixture e prove que ele não aparece automaticamente. Cubra relacionamentos opcionais carregados e não carregados, além dos caminhos de erro de codificação.

<a name="dtos-and-jsonresource-have-different-contracts"></a>
## DTOs e `JsonResource` têm contratos diferentes

`ToArray` de um registro ou de uma coleção, controles de visibilidade, campos ocultos e atributos acrescentados definem a representação em array do pacote model. Um DTO comum segue as tags de `encoding/json` e `MarshalJSON`. Já `Context.JSON` aceita `hhttp.JsonResource`, cujos métodos `ToArray` e With retornam mapas em runtime; suas chaves exigem testes exatos de resposta. Escolha explicitamente o contrato de transporte, sem tratar as duas formas como o mesmo recurso.

<a name="logging-needs-its-own-safe-representation"></a>
## Logs precisam de representação segura própria

Lista de campos públicos não remove automaticamente dados sensíveis de logs, eventos ou dumps. Tipos com credenciais ou chaves precisam de comportamento seguro de logging e payloads deliberados. Criptografia protege texto armazenado, mas não torna texto cifrado ou metadados apropriados para toda resposta pública.
