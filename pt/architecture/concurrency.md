# Concorrência

- [O trabalho concorrente ainda tem um responsável](#concurrent-work-still-has-one-owner)
- [Execute leituras independentes com um limite](#run-independent-reads-with-a-ceiling)
- [Adie trabalho deliberadamente e recupere falhas nos limites sob sua responsabilidade](#defer-work-deliberately-and-recover-owned-boundaries)
- [Um conjunto com falha não retorna slice parcial](#a-failed-set-returns-no-partial-slice)
- [Verifique o limite e o cancelamento](#prove-the-ceiling-and-cancellation)
- [Operações com tarefas finitas](#finite-task-operations)
- [Loops em segundo plano pertencem aos módulos](#background-loops-belong-to-modules)

<a name="concurrent-work-still-has-one-owner"></a>
## O trabalho concorrente ainda tem um responsável

Os auxiliares de concorrência executam um conjunto finito de tarefas tipadas sob um contexto. Preservam a ordem dos argumentos, cancelam as tarefas irmãs após a primeira falha, aguardam todas as goroutines iniciadas e convertem panics em erros, em vez de deixar uma tarefa encerrar o processo.

<a name="run-independent-reads-with-a-ceiling"></a>
## Execute leituras independentes com um limite

Use `RunWith` quando leituras independentes puderem avançar juntas, mas compartilharem um número finito de conexões ou uma cota externa. Os resultados seguem a ordem dos argumentos das tarefas, não a ordem de conclusão; um limite positivo impede que a distribuição consuma todos os recursos disponíveis.

```go example
type Invoice struct{ ID string }

type InvoiceFinder interface {
	Find(ctx context.Context, id string) (Invoice, error)
}

func LoadPair(ctx context.Context, invoices InvoiceFinder, firstID, secondID string) ([]Invoice, error) {
	results, err := concurrency.RunWith(ctx, 4,
		func(ctx context.Context) (Invoice, error) { return invoices.Find(ctx, firstID) },
		func(ctx context.Context) (Invoice, error) { return invoices.Find(ctx, secondID) },
	)
	if err != nil {
		return nil, err
	}
	return results, nil
}
```

<a name="defer-work-deliberately-and-recover-owned-boundaries"></a>
## Adie trabalho deliberadamente e recupere falhas nos limites sob sua responsabilidade

Run não tem limite, e `RunWith` acrescenta um; zero ou um valor negativo significam ausência de limite. Recover protege um ponto de chamada nomeado e retorna `PanicError` com o valor e a stack. Defer retorna um callback e não executa nada até ele ser chamado.

- Use um limite abaixo do pool de banco de dados ou da cota do serviço externo.
- Dê um prazo próprio ao trabalho adiado.
- Chame o callback adiado retornado exatamente uma vez.
- Mantenha loops de longa duração nos contratos de ciclo de vida do módulo.

<a name="a-failed-set-returns-no-partial-slice"></a>
## Um conjunto com falha não retorna slice parcial

A primeira falha de tarefa cancela as demais e retorna um slice nil de resultados. Uma tarefa que ignora seu contexto derivado pode manter a chamada inteira esperando mesmo após seu resultado ser descartado. Um contexto pai já cancelado não inicia nada; o cancelamento durante a execução retorna o erro do pai.

> [!WARNING]
> Defer desacopla o cancelamento e os prazos da requisição usando `context.WithoutCancel`. O novo responsável deve definir um prazo, relatar falhas e tornar segura a execução duplicada; chamar o callback duas vezes executa as tarefas duas vezes.

<a name="prove-the-ceiling-and-cancellation"></a>
## Verifique o limite e o cancelamento

Conte entradas simultâneas de tarefas para comprovar o limite. Faça uma tarefa falhar para comprovar que a irmã observa o cancelamento, inverta a ordem de conclusão para comprovar a ordem dos resultados e provoque panic em uma tarefa para identificar `*concurrency.PanicError`.

- O máximo de tarefas simultâneas corresponde ao limite configurado.
- Uma tarefa irmã bloqueada termina após o primeiro erro.
- Os resultados preservam a ordem dos argumentos.
- O trabalho adiado não executa antes de seu callback.

<a name="finite-task-operations"></a>
## Operações com tarefas finitas

`Task[T]`, Run, `RunWith`, Defer, Recover e `PanicError` formam a API direta. `Manager[T]` resolve drivers configurados, mas os nomes internos process, fork e sync atualmente executam o mesmo comportamento com goroutines. Prefira as funções diretas, salvo quando a própria configuração for necessária.

<a name="background-loops-belong-to-modules"></a>
## Loops em segundo plano pertencem aos módulos

Um conjunto finito de tarefas dentro de uma requisição difere de um scheduler, worker de fila ou loop de relay. Trabalho de longa duração implementa o ciclo de vida da aplicação, propaga um contexto de encerramento pertencente ao processo e expõe saúde. Exclusividade entre réplicas exige um lock compartilhado, não goroutines.

- [API de concorrência](https://pkg.go.dev/github.com/arandu-io/hesape/concurrency)
- [Contexto](/docs/{{version}}/architecture/context)
- [Módulos](/docs/{{version}}/architecture/modules-and-composition)
