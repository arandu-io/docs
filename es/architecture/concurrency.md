# Concurrencia

- [El trabajo concurrente sigue teniendo un responsable](#concurrent-work-still-has-one-owner)
- [Ejecuta lecturas independientes con un límite](#run-independent-reads-with-a-ceiling)
- [Difiere trabajo deliberadamente y recupera fallos en los límites bajo tu responsabilidad](#defer-work-deliberately-and-recover-owned-boundaries)
- [Un conjunto fallido no devuelve un slice parcial](#a-failed-set-returns-no-partial-slice)
- [Verifica el límite y la cancelación](#prove-the-ceiling-and-cancellation)
- [Operaciones con tareas finitas](#finite-task-operations)
- [Los bucles en segundo plano pertenecen a los módulos](#background-loops-belong-to-modules)

<a name="concurrent-work-still-has-one-owner"></a>
## El trabajo concurrente sigue teniendo un responsable

Los auxiliares de concurrencia ejecutan un conjunto finito de tareas tipadas bajo un contexto. Conservan el orden de los argumentos, cancelan las tareas hermanas tras el primer fallo, esperan a todas las goroutines iniciadas y convierten panics en errores, en vez de permitir que una tarea termine el proceso.

<a name="run-independent-reads-with-a-ceiling"></a>
## Ejecuta lecturas independientes con un límite

Usa `RunWith` cuando lecturas independientes puedan avanzar juntas, pero compartan un número finito de conexiones o una cuota externa. Los resultados mantienen el orden de los argumentos de las tareas, no el de finalización; un límite positivo impide que la distribución consuma todos los recursos disponibles.

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
## Difiere trabajo deliberadamente y recupera fallos en los límites bajo tu responsabilidad

Run no tiene límite y `RunWith` añade uno; cero o un valor negativo significan ausencia de límite. Recover protege un punto de llamada nombrado y devuelve `PanicError` con el valor y la stack. Defer devuelve un callback y no ejecuta nada hasta que se invoque.

- Usa un límite inferior al pool de base de datos o a la cuota del servicio externo.
- Asigna un plazo propio al trabajo diferido.
- Invoca el callback diferido devuelto exactamente una vez.
- Mantén los bucles de larga duración en los contratos de ciclo de vida del módulo.

<a name="a-failed-set-returns-no-partial-slice"></a>
## Un conjunto fallido no devuelve un slice parcial

El primer fallo de tarea cancela las demás y devuelve un slice nil de resultados. Una tarea que ignora su contexto derivado puede mantener toda la llamada esperando incluso después de descartar su resultado. Un contexto padre ya cancelado no inicia nada; la cancelación durante la ejecución devuelve el error del padre.

> [!WARNING]
> Defer desacopla la cancelación y los plazos de la solicitud mediante `context.WithoutCancel`. El nuevo responsable debe definir un plazo, informar de fallos y hacer segura la ejecución duplicada; invocar el callback dos veces ejecuta las tareas dos veces.

<a name="prove-the-ceiling-and-cancellation"></a>
## Verifica el límite y la cancelación

Cuenta las entradas simultáneas de tareas para demostrar el límite. Haz fallar una tarea para demostrar que su hermana observa la cancelación, invierte el orden de finalización para demostrar el orden de resultados y provoca un panic en una tarea para identificar `*concurrency.PanicError`.

- El máximo de tareas simultáneas coincide con el límite configurado.
- Una tarea hermana bloqueada termina después del primer error.
- Los resultados conservan el orden de los argumentos.
- El trabajo diferido no se ejecuta antes de su callback.

<a name="finite-task-operations"></a>
## Operaciones con tareas finitas

`Task[T]`, Run, `RunWith`, Defer, Recover y `PanicError` forman la API directa. `Manager[T]` resuelve drivers configurados, pero los nombres internos process, fork y sync ejecutan actualmente el mismo comportamiento con goroutines. Prefiere las funciones directas salvo cuando la propia configuración sea necesaria.

<a name="background-loops-belong-to-modules"></a>
## Los bucles en segundo plano pertenecen a los módulos

Un conjunto finito de tareas dentro de una solicitud difiere de un scheduler, worker de cola o bucle de relay. El trabajo de larga duración implementa el ciclo de vida de la aplicación, propaga un contexto de cierre perteneciente al proceso y expone salud. La exclusividad entre réplicas exige un lock compartido, no goroutines.

- [API de concurrencia](https://pkg.go.dev/github.com/arandu-io/hesape/concurrency)
- [Contexto](/docs/{{version}}/architecture/context)
- [Módulos](/docs/{{version}}/architecture/modules-and-composition)
