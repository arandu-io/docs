# Contexto

- [Contexto transporta cancelación, no dependencias](#context-carries-cancellation-not-dependencies)
- [Pasa el tiempo de vida de la solicitud al servicio](#pass-the-request-lifetime-into-the-service)
- [Deriva un contexto hijo con plazo para una dependencia](#derive-a-bounded-child-for-one-dependency)
- [La cancelación es un resultado, no una tormenta de logs](#cancellation-is-an-outcome-not-a-log-storm)
- [Cancela el contexto padre y observa la llamada final](#cancel-the-parent-and-observe-the-leaf)
- [API de contexto estándar y de solicitud](#standard-and-request-context-apis)
- [La cancelación atraviesa cada adapter explícito](#cancellation-crosses-every-explicit-adapter)

<a name="context-carries-cancellation-not-dependencies"></a>
## Contexto transporta cancelación, no dependencias

El contexto de solicitud indica cuándo deja de ser útil el trabajo y transporta un conjunto pequeño de hechos de esa solicitud entre paquetes. Pásalo primero a cada operación que pueda bloquear. Los constructores siguen transportando dependencias duraderas; configuración y servicios no pertenecen a los valores del contexto.

<a name="pass-the-request-lifetime-into-the-service"></a>
## Pasa el tiempo de vida de la solicitud al servicio

`http.Context.Ctx` devuelve el contexto de la solicitud subyacente. Pásalo mediante el controller a operaciones de base de datos, caché, HTTP externo, locks y concurrencia para que las desconexiones y los plazos lleguen a la llamada que realiza el trabajo. El Subject autenticado viaja junto a él hasta el service, porque un parámetro de ruta es un identificador no confiable y nunca prueba de acceso; una solicitud sin Subject se rechaza antes de que el service se ejecute.

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
## Deriva un contexto hijo con plazo para una dependencia

Deriva un contexto hijo con un plazo menor cuando una dependencia necesite un límite más estricto. El middleware de autenticación puede asociar un Subject con `auth.WithSubject`; los consumidores usan `auth.SubjectFrom` y tratan el booleano en vez de suponer que existe identidad.

- Usa `context.WithTimeout` para una llamada a una dependencia con plazo.
- Programa cancel con defer inmediatamente después de derivar el contexto.
- Lee Subject junto con el booleano que indica su presencia.
- Transfiere trabajo desacoplado solo a un responsable explícito.

<a name="cancellation-is-an-outcome-not-a-log-storm"></a>
## La cancelación es un resultado, no una tormenta de logs

Propaga `context.Canceled` y `context.DeadlineExceeded` para que la frontera decida si una respuesta aún es útil. No sustituyas el contexto de solicitud por `context.Background` para mantener en ejecución una llamada fallida y nunca retengas `*http.Context` después de que retorne el handler.

> [!WARNING]
> `context.WithoutCancel` elimina el plazo además de la cancelación. El trabajo desacoplado debe establecer un nuevo plazo, informar de sus fallos y hacer segura la ejecución repetida.

<a name="cancel-the-parent-and-observe-the-leaf"></a>
## Cancela el contexto padre y observa la llamada final

Cancela el contexto padre y comprueba que retorna una llamada bloqueada al repositorio o cliente. Haz expirar un plazo hijo mientras un contexto hermano conectado al padre sigue siendo utilizable y ejecuta el detector de carreras sobre los valores del ámbito de la solicitud.

- La cancelación llega a la llamada final.
- El timeout de un contexto hijo no cancela al padre.
- La ausencia de Subject significa falta de autenticación; solo `auth.Guest(tenant)` entra en la Policy como invitado.
- No se usa ningún objeto de solicitud después de que retorne el handler.

<a name="standard-and-request-context-apis"></a>
## API de contexto estándar y de solicitud

La biblioteca estándar proporciona Context, `WithCancel`, `WithTimeout`, `WithDeadline`, Cause y `WithoutCancel`. El `http.Context` de Framework expone Ctx para el tiempo de vida de la solicitud; los paquetes de autenticación proporcionan auxiliares específicos de Subject en el contexto.

<a name="cancellation-crosses-every-explicit-adapter"></a>
## La cancelación atraviesa cada adapter explícito

La concurrencia deriva contextos hermanos de este valor. Los controllers inician la cadena, los servicios y repositorios la conservan y los clientes externos construyen solicitudes con ella. Los módulos en segundo plano reciben un contexto perteneciente al proceso, sin tomar prestado el de una solicitud.

- [Contexto estándar](https://pkg.go.dev/context)
- [Controllers](/docs/{{version}}/http/controllers)
- [Concurrencia](/docs/{{version}}/architecture/concurrency)
