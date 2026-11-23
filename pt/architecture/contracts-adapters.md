# Contratos e adapters

- [O consumidor é responsável pela menor interface útil](#the-consumer-owns-the-smallest-useful-interface)
- [Inclua autorização no contrato do repositório](#put-authorization-in-the-repository-seam)
- [Selecione adapters na inicialização](#select-adapters-at-boot)
- [Evite interfaces que escondam autoridade ou responsabilidade](#avoid-interfaces-that-hide-authority-or-ownership)
- [Execute um contrato contra cada adapter](#run-one-contract-against-each-adapter)
- [Os contratos do Framework permanecem restritos](#framework-contracts-stay-narrow)
- [A composição torna a substituição visível](#composition-makes-substitution-visible)

<a name="the-consumer-owns-the-smallest-useful-interface"></a>
## O consumidor é responsável pela menor interface útil

Um contrato nomeia o comportamento de que um consumidor precisa. Um adapter o satisfaz estruturalmente e é selecionado na raiz de composição. O código de domínio permanece independente de SQL, RESP, SMTP ou de um serviço externo, sem acrescentar um registro ou pacote universal de abstrações.

<a name="put-authorization-in-the-repository-seam"></a>
## Inclua autorização no contrato do repositório

Inclua autorização no contrato do repositório para que trocar o armazenamento não remova o escopo do tenant. O serviço declara apenas a operação que chama, e a asserção de compilação mantém o adapter SQL concreto alinhado à interface pertencente a esse consumidor.

```go example
type Invoice struct{ ID string }

type InvoiceRepository interface {
	Find(context.Context, security.Grant, string) (Invoice, error)
}

type InvoiceService struct{ invoices InvoiceRepository }

func NewInvoiceService(invoices InvoiceRepository) *InvoiceService {
	return &InvoiceService{invoices: invoices}
}

type SQLInvoiceRepository struct{}

func (*SQLInvoiceRepository) Find(context.Context, security.Grant, string) (Invoice, error) {
	return Invoice{}, nil
}

var _ InvoiceRepository = (*SQLInvoiceRepository)(nil)
```

<a name="select-adapters-at-boot"></a>
## Selecione adapters na inicialização

`foundation.NewDrivers` constrói o conector de banco de dados, o store de cache, o transporte de e-mail e o disco do sistema de arquivos que a configuração nomeia, e um provider passa cada um, ou um fake de teste, uma vez. Os adaptadores traduzem erros de ausência, conflito e indisponibilidade do driver quando os chamadores precisam ramificar, preservando a causa para os operadores.

- Valide dependências obrigatórias nos construtores.
- Mantenha comportamentos opcionais em interfaces pequenas separadas.
- Encapsule erros de driver sem expor credenciais ou SQL.
- Coloque asserções de compilação ao lado dos adapters.

<a name="avoid-interfaces-that-hide-authority-or-ownership"></a>
## Evite interfaces que escondam autoridade ou responsabilidade

Um repositório que recebe apenas um id para dados de um tenant não tem onde comprovar autoridade. Interfaces grandes obrigam fakes a implementar métodos não usados e acoplam consumidores sem relação; interfaces pequenas duplicadas são mais seguras quando as semânticas realmente diferem.

> [!WARNING]
> Não compense autoridade ausente com um tenant global, singleton local da requisição ou leitura de ambiente dentro do adapter. Passe `security.Grant` na operação que acessa o armazenamento protegido.

<a name="run-one-contract-against-each-adapter"></a>
## Execute um contrato contra cada adapter

Execute um contrato compartilhado contra cada adapter. Teste sucesso, ausência, isolamento de tenant, cancelamento e mapeamento de erros pela interface do consumidor, sem acessar detalhes da implementação de um driver concreto.

- Grants zero e de outros tenants não retornam dados protegidos.
- Adapters de memória e de rede cumprem o mesmo comportamento.
- O cancelamento retorna prontamente.
- Asserções de compilação cobrem cada adapter concreto.

<a name="framework-contracts-stay-narrow"></a>
## Os contratos do Framework permanecem restritos

Escolha o contrato de comunicação pela interação, não por um ranking de popularidade. REST é requisição e resposta para recursos públicos; `gRPC` é uma fronteira de serviço tipada que também oferece streaming; uma fila persiste trabalho assíncrono; `WebSocket` é uma sessão bidirecional em tempo real; eventos enviados pelo servidor (SSE) transmitem atualizações do servidor em uma resposta HTTP; e um webhook entrega um evento assinado a um endpoint HTTP externo.

| Superfície | Contrato publicado |
| --- | --- |
| REST | recursos de requisição e resposta; use para APIs públicas ou voltadas ao navegador |
| `gRPC` | chamadas tipadas unárias ou por streaming; use Connect, `gRPC` ou gRPC-Web no servidor HTTP da aplicação |
| Fila | trabalho assíncrono durável com entrega pelo menos uma vez |
| `WebSocket` | sessões bidirecionais persistentes para interação de baixa latência |
| SSE | atualizações unidirecionais do servidor por uma resposta HTTP cancelável |
| Webhook | entrega durável e assinada de eventos para um serviço externo |

<a name="composition-makes-substitution-visible"></a>
## A composição torna a substituição visível

REST, `gRPC`, filas, `WebSocket`, SSE e webhooks resolvem problemas diferentes de propriedade e entrega. A configuração escolhe o adapter concreto, a raiz de composição o constrói e o serviço mantém apenas seu contrato. Context carrega a vida útil da requisição e Grant carrega o escopo de tenant autorizado, tornando a substituição visível em um único diff de composição.

Não substitua armazenamento durável por um transporte em tempo real: `WebSocket` e SSE não reproduzem sozinhos o trabalho perdido. Uma fila executa trabalho interno, enquanto um webhook cruza uma fronteira HTTP externa. O streaming `gRPC` compartilha o processo servidor e o ciclo de desligamento, em vez de criar um segundo listener.

- [Roteamento HTTP e `gRPC`](/docs/{{version}}/http/routing)
- [Respostas HTTP e SSE](/docs/{{version}}/http/responses)
- [Filas](/docs/{{version}}/async/queues)
- [`WebSocket` do Joaju](/docs/{{version}}/async/joaju)
- [Eventos e webhooks](/docs/{{version}}/async/events-outbox)
