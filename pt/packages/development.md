# Desenvolvimento de pacotes

- [Um pacote é um módulo Go com ponto de entrada explícito](#a-package-is-a-go-module-with-an-explicit-entry-point)
- [Implemente a interface obrigatória do módulo](#implement-the-required-module-surface)
- [Adicione contratos opcionais de ciclo de vida de forma deliberada](#add-optional-lifecycle-contracts-deliberately)
- [Faça a ausência de dependências falhar na construção](#make-missing-wiring-fail-at-construction)
- [Teste o pacote por seu contrato exportado](#test-the-package-through-its-exported-contract)
- [Interfaces obrigatórias e opcionais](#required-and-optional-interfaces)
- [Coloque primitivas em componentes e o estado junto de seu responsável](#put-primitives-in-components-and-state-with-its-owner)

> [!WARNING]
> **Limitações atuais.** O esqueleto de pacote está publicado e tem testes próprios, mas nenhuma aplicação independente com tag o importa. Trate suas interfaces de módulo, Policy, resource, migration e rota como um ponto de partida verificado e acrescente um fluxo de teste na aplicação consumidora antes de afirmar que há integração.

<a name="a-package-is-a-go-module-with-an-explicit-entry-point"></a>
## Um pacote é um módulo Go com ponto de entrada explícito

Um pacote instalável tem seu próprio caminho de módulo e expõe comportamento por APIs Go comuns. A integração com o Framework usa `foundation.Module`, cujo contrato obrigatório contém apenas um Name estável e o registro de rotas. As dependências ficam explícitas no construtor; nenhum container ou reflexão as resolve por trás da raiz de composição.

<a name="implement-the-required-module-surface"></a>
## Implemente a interface obrigatória do módulo

Name devolve um identificador estável em minúsculas e Routes recebe o router da aplicação. Capacidades opcionais são interfaces estruturais descobertas só quando implementadas; migrations e trabalho agendado não são métodos que todo módulo precisa fingir ter. A aplicação constrói o módulo e o lista em `WithProviders` em `bootstrap/app.go`.

```go example
type Handler struct{}

func (*Handler) Index(ctx *fhttp.Context) error { return ctx.Status(stdhttp.StatusOK) }

type Module struct {
	handler *Handler
}

func (*Module) Name() string { return "invoices" }

func (m *Module) Routes(r *fhttp.Router) {
	r.Action("GET", "/invoices", m.handler.Index).Name("invoices.index")
}
```

<a name="add-optional-lifecycle-contracts-deliberately"></a>
## Adicione contratos opcionais de ciclo de vida de forma deliberada

Implemente `foundation.Migratable` apenas quando o módulo possuir schema; seu método exato é `Migrations() []foundation.Migration`. Implemente `foundation.Schedulable` apenas quando declarar tarefas do kernel; seu método exato é `Schedule() []foundation.Task`. Mantenha views, traduções, policies, handlers e testes no módulo responsável por eles. Um pacote declara trabalho; não inicia goroutines sem gerenciamento durante o import.

<a name="make-missing-wiring-fail-at-construction"></a>
## Faça a ausência de dependências falhar na construção

Rejeite dependências obrigatórias nil no construtor do módulo e deixe nomes de rota duplicados, configuração inválida ou conflitos de migration aparecerem no bootstrap. Não substitua silenciosamente por adapters em memória quando a configuração de produção exigir um backend compartilhado. Uma interface opcional deve estar ausente, não implementada como uma operação vazia que sugere suporte.

<a name="test-the-package-through-its-exported-contract"></a>
## Teste o pacote por seu contrato exportado

Compile um consumidor pequeno que construa o módulo, registre rotas e execute uma requisição. Rode testes de migrations ou agendamento apenas para as interfaces implementadas pelo módulo. Os testes do pacote verificam seus repositórios e policies; a aplicação consumidora ainda testa a composição e a configuração que escolheu.

<a name="required-and-optional-interfaces"></a>
## Interfaces obrigatórias e opcionais

`foundation.Module` exige `Name()` string e `Routes(*http.Router)`. `foundation.Migratable` exige `Migrations() []foundation.Migration`, e `foundation.Schedulable` exige `Schedule() []foundation.Task`; ambas continuam opcionais e não ampliam Module. Os construtores continuam funções comuns, e os handlers de rota continuam compatíveis com `net/http` por meio do roteador do framework.

<a name="put-primitives-in-components-and-state-with-its-owner"></a>
## Coloque primitivas em componentes e o estado junto de seu responsável

Um cálculo ou adapter reutilizável pode viver na biblioteca de componentes sem virar um módulo instalável. Um recurso que possui schema, rotas e Policy pertence a um pacote. A raiz de composição da aplicação é onde os adapters do pacote encontram a configuração concreta, mantendo o pacote reutilizável sem esconder sua execução.
