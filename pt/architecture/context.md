# Contexto

- [Contexto carrega cancelamento, não dependências](#context-carries-cancellation-not-dependencies)
- [Passe o tempo de vida da requisição ao serviço](#pass-the-request-lifetime-into-the-service)
- [Derive um contexto filho com prazo para uma dependência](#derive-a-bounded-child-for-one-dependency)
- [Cancelamento é um resultado, não uma tempestade de logs](#cancellation-is-an-outcome-not-a-log-storm)
- [Cancele o contexto pai e observe a chamada final](#cancel-the-parent-and-observe-the-leaf)
- [APIs de contexto padrão e de requisição](#standard-and-request-context-apis)
- [O cancelamento atravessa cada adapter explícito](#cancellation-crosses-every-explicit-adapter)

<a name="context-carries-cancellation-not-dependencies"></a>
## Contexto carrega cancelamento, não dependências

O contexto da requisição informa quando o trabalho deixa de ser útil e carrega um conjunto pequeno de fatos dessa requisição entre pacotes. Passe-o primeiro a cada operação que possa bloquear. Os construtores continuam carregando dependências duradouras; configuração e serviços não pertencem aos valores do contexto.

<a name="pass-the-request-lifetime-into-the-service"></a>
## Passe o tempo de vida da requisição ao serviço

`http.Context.Ctx` retorna o contexto da requisição subjacente. Passe-o pelo controller para operações de banco de dados, cache, HTTP externo, locks e concorrência, para que desconexões e prazos cheguem à chamada que executa o trabalho. O Subject autenticado segue ao lado dele até o service, porque um parâmetro de rota é um identificador não confiável e nunca prova de acesso; uma requisição sem Subject é recusada antes de o service rodar.

```go example
type Invoice struct{ ID, Reference string }

func (i Invoice) ToArray() map[string]any {
	return map[string]any{"id": i.ID, "reference": i.Reference}
}

func (Invoice) With() map[string]any { return nil }

type InvoiceFinder interface {
	Find(ctx context.Context, subject security.Subject, id string) (Invoice, error)
}

type InvoiceController struct{ invoices InvoiceFinder }

func (c *InvoiceController) Show(ctx *fhttp.Context) error {
	subject, ok := ctx.User()
	if !ok {
		return ctx.Status(stdhttp.StatusUnauthorized)
	}
	invoice, err := c.invoices.Find(ctx.Ctx(), subject, ctx.Param("invoice"))
	if err != nil {
		return err
	}
	return ctx.JSON(stdhttp.StatusOK, invoice)
}
```

<a name="derive-a-bounded-child-for-one-dependency"></a>
## Derive um contexto filho com prazo para uma dependência

Derive um contexto filho com prazo menor quando uma dependência precisar de um limite mais estrito. O middleware de autenticação pode associar um Subject com `auth.WithSubject`; consumidores usam `auth.SubjectFrom` e tratam o booleano em vez de presumir que há identidade.

- Use `context.WithTimeout` para uma chamada de dependência com prazo.
- Agende cancel com defer imediatamente após derivar o contexto.
- Leia Subject junto com o booleano que indica sua presença.
- Transfira trabalho desacoplado apenas para um responsável explícito.

<a name="cancellation-is-an-outcome-not-a-log-storm"></a>
## Cancelamento é um resultado, não uma tempestade de logs

Propague `context.Canceled` e `context.DeadlineExceeded` para que a fronteira decida se uma resposta ainda é útil. Não substitua o contexto da requisição por `context.Background` para manter uma chamada com falha em execução e nunca retenha `*http.Context` depois que o handler retorna.

> [!WARNING]
> `context.WithoutCancel` remove o prazo além do cancelamento. Trabalho desacoplado deve estabelecer um novo prazo, relatar suas falhas e tornar segura a execução repetida.

<a name="cancel-the-parent-and-observe-the-leaf"></a>
## Cancele o contexto pai e observe a chamada final

Cancele o contexto pai e verifique se uma chamada bloqueada ao repositório ou cliente retorna. Faça um prazo filho expirar enquanto um contexto irmão ligado ao pai continua utilizável e execute o detector de corrida sobre os valores do escopo da requisição.

- O cancelamento chega à chamada final.
- O timeout de um contexto filho não cancela o pai.
- Subject ausente significa ausência de autenticação; apenas `auth.Guest(tenant)` entra na Policy como convidado.
- Nenhum objeto da requisição é usado depois que o handler retorna.

<a name="standard-and-request-context-apis"></a>
## APIs de contexto padrão e de requisição

A biblioteca padrão fornece Context, `WithCancel`, `WithTimeout`, `WithDeadline`, Cause e `WithoutCancel`. O `http.Context` do Framework expõe Ctx para o tempo de vida da requisição; os pacotes de autenticação fornecem auxiliares específicos de Subject no contexto.

<a name="cancellation-crosses-every-explicit-adapter"></a>
## O cancelamento atravessa cada adapter explícito

A concorrência deriva contextos irmãos desse valor. Controllers iniciam a cadeia, serviços e repositórios a preservam, e clientes externos constroem requisições com ela. Módulos em segundo plano recebem um contexto pertencente ao processo, sem tomar emprestado o de uma requisição.

- [Contexto padrão](https://pkg.go.dev/context)
- [Controllers](/docs/{{version}}/http/controllers)
- [Concorrência](/docs/{{version}}/architecture/concurrency)
