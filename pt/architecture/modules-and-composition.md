# Módulos e raiz de composição

- [A raiz de composição é Go legível](#the-composition-root-is-readable-go)
- [Registre um módulo pequeno](#register-one-small-module)
- [Adicione comportamento de ciclo de vida por interfaces opcionais](#add-lifecycle-behavior-through-optional-interfaces)
- [Rejeite uma composição ambígua](#refuse-ambiguous-composition)
- [Inicialize o mesmo grafo usado pelo servidor](#boot-the-same-graph-a-server-uses)
- [Contratos obrigatórios e opcionais de módulo](#required-and-optional-module-contracts)
- [As interfaces se encontram na raiz](#interfaces-meet-at-the-root)

<a name="the-composition-root-is-readable-go"></a>
## A raiz de composição é Go legível

Um módulo contribui com um nome estável e rotas. Os construtores recebem cada dependência explicitamente, `foundation.NewDrivers` constrói os adaptadores concretos a partir da configuração e `bootstrap/app.go` lista os módulos em `WithProviders`. Não há container de serviços, varredura por reflexão nem resolvedor oculto entre um handler e seus colaboradores.

<a name="register-one-small-module"></a>
## Registre um módulo pequeno

Defina o menor módulo que possui o grupo de rotas, construa seu controller a partir de um serviço da aplicação num provider e liste esse valor em `WithProviders`. O builder aceita `foundation.Module`, então o compilador prova Name e Routes enquanto o bootstrap prova a implementação escolhida.

```go example
type InvoiceService struct{}

type InvoiceController struct{ service *InvoiceService }

func NewInvoiceController(service *InvoiceService) *InvoiceController {
	return &InvoiceController{service: service}
}

func (c *InvoiceController) Index(ctx *fhttp.Context) error { return ctx.Status(stdhttp.StatusOK) }

type BillingModule struct{ invoices *InvoiceController }

func (m *BillingModule) Name() string { return "billing" }
func (m *BillingModule) Routes(r *fhttp.Router) {
	r.Action("GET", "/invoices", m.invoices.Index).Name("invoices.index")
}

func Compose(drivers *foundation.Drivers, app foundation.Module, service *InvoiceService) (*foundation.Application, error) {
	billing := &BillingModule{invoices: NewInvoiceController(service)}
	return foundation.Configure(drivers).WithProviders(billing, app).Create()
}
```

> [!NOTE]
> `WithProviders` inicializa os módulos na ordem em que os lista. Um módulo escrito contra `hesape/foundation` cujo Routes ainda recebe o router só de handlers é listado como `foundation.Component(module)`, o que mantém visível ao kernel cada interface opcional que o módulo implementa.

<a name="add-lifecycle-behavior-through-optional-interfaces"></a>
## Adicione comportamento de ciclo de vida por interfaces opcionais

O comportamento do ciclo de vida é expresso por interfaces opcionais detectadas no valor registrado. Implemente apenas os contratos pelos quais o módulo é responsável e nunca inicie uma goroutine no construtor ou em Routes, pois a aplicação controla o início e o encerramento do processo.

- Bootable prepara recursos e valida a composição.
- Background inicia um loop apenas quando a aplicação executa.
- Closable libera recursos durante o encerramento.
- Migratable, Schedulable e Health declaram schema, tarefas e verificações de saúde.
- `RendererProvider` fornece o renderizador único usado pelas rotas.
- Commandable, CSRFExempt e `JobHandlers` declaram comandos de console, isenções de CSRF para os receptores do próprio módulo e handlers de jobs; Create os coleta de todos os providers.

<a name="refuse-ambiguous-composition"></a>
## Rejeite uma composição ambígua

Nomes vazios ou duplicados, dois provedores de renderer, uma falha de boot e conflitos de rota são erros de ligação. Create e Boot também recusam um nome de comando declarado duas vezes, um nome de job com dois handlers, uma rota da aplicação que responde a um caminho que o framework monta sob `/_arandu/` e uma isenção de CSRF de módulo que cobre a rota de escrita de outro dono. Eles falham antes do tráfego. A ordem de registro é observável no boot, nas rotas, nas migrations e no desligamento, então mantenha-a intencional em vez de depender da iteração de um map.

<a name="boot-the-same-graph-a-server-uses"></a>
## Inicialize o mesmo grafo usado pelo servidor

Inicialize o mesmo grafo que o servidor usará e inspecione suas rotas e ciclo de vida em vez de testar um controller isolado. Verifique a ordem do boot, o início do processamento em segundo plano, o resultado de saúde e a liberação em ordem inversa segura de cada recurso do módulo.

- Boot não retorna erro para o grafo completo.
- O inventário de rotas identifica o módulo e os handlers esperados.
- O trabalho em segundo plano começa em Run, não na construção.
- O encerramento libera os recursos próprios em uma ordem segura.

<a name="required-and-optional-module-contracts"></a>
## Contratos obrigatórios e opcionais de módulo

`foundation.Module` exige `Name()` string e `Routes(*http.Router)`. `foundation.Configure`, `WithProviders` e Create compõem o grafo; `Application.HandleCommand`, Boot, Run, Handler e Shutdown são donos dele. Os contratos opcionais incluem Bootable, Background, `RendererProvider`, Closable, Diagnostic, Schedulable, Migratable, Health, Ready, Commandable, CSRFExempt e `JobHandlers`, e `foundation.Component` adapta um módulo de `hesape/foundation`.

<a name="interfaces-meet-at-the-root"></a>
## As interfaces se encontram na raiz

Contratos mantêm os serviços independentes de adapters; a configuração seleciona esses adapters; as rotas expõem controllers já construídos. A estrutura de diretórios identifica cada responsável, enquanto o ciclo da requisição mostra quando o grafo resultante executa.

- [Configuração](/docs/{{version}}/architecture/configuration)
- [Contratos e adapters](/docs/{{version}}/architecture/contracts-adapters)
- [Module de Foundation](https://pkg.go.dev/github.com/arandu-io/framework/foundation#Module)
