# Módulos y raíz de composición

- [La raíz de composición es Go legible](#the-composition-root-is-readable-go)
- [Registra un módulo pequeño](#register-one-small-module)
- [Añade comportamiento de ciclo de vida mediante interfaces opcionales](#add-lifecycle-behavior-through-optional-interfaces)
- [Rechaza una composición ambigua](#refuse-ambiguous-composition)
- [Inicializa el mismo grafo usado por el servidor](#boot-the-same-graph-a-server-uses)
- [Contratos obligatorios y opcionales de módulo](#required-and-optional-module-contracts)
- [Las interfaces se encuentran en la raíz](#interfaces-meet-at-the-root)

<a name="the-composition-root-is-readable-go"></a>
## La raíz de composición es Go legible

Un módulo aporta un nombre estable y rutas. Los constructores reciben cada dependencia explícitamente, `foundation.NewDrivers` construye los adaptadores concretos a partir de la configuración y `bootstrap/app.go` lista los módulos en `WithProviders`. No hay contenedor de servicios, escaneo por reflexión ni resolutor oculto entre un handler y sus colaboradores.

<a name="register-one-small-module"></a>
## Registra un módulo pequeño

Defina el módulo más pequeño que posee el grupo de rutas, construya su controller a partir de un servicio de la aplicación en un provider y liste ese valor en `WithProviders`. El builder acepta `foundation.Module`, así que el compilador prueba Name y Routes mientras el bootstrap prueba la implementación elegida.

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
> `WithProviders` inicializa los módulos en el orden en que los lista. Un módulo escrito contra `hesape/foundation` cuyo Routes todavía recibe el router solo de handlers se lista como `foundation.Component(module)`, lo que mantiene visible para el kernel cada interfaz opcional que el módulo implementa.

<a name="add-lifecycle-behavior-through-optional-interfaces"></a>
## Añade comportamiento de ciclo de vida mediante interfaces opcionales

El comportamiento del ciclo de vida se expresa mediante interfaces opcionales detectadas en el valor registrado. Implementa solo los contratos de los que sea responsable el módulo y nunca inicies una goroutine en el constructor ni en Routes, porque la aplicación controla el inicio y cierre del proceso.

- Bootable prepara recursos y valida la composición.
- Background inicia un bucle solo cuando se ejecuta la aplicación.
- Closable libera recursos durante el cierre.
- Migratable, Schedulable y Health declaran schema, tareas y comprobaciones de salud.
- `RendererProvider` proporciona el único renderizador usado por las rutas.
- Commandable, CSRFExempt y `JobHandlers` declaran comandos de consola, exenciones de CSRF para los receptores del propio módulo y handlers de jobs; Create los recoge de todos los providers.

<a name="refuse-ambiguous-composition"></a>
## Rechaza una composición ambigua

Los nombres vacíos o duplicados, dos proveedores de renderer, un fallo de arranque y los conflictos de rutas son errores de cableado. Create y Boot también rechazan un nombre de comando declarado dos veces, un nombre de job con dos handlers, una ruta de la aplicación que responde a una ruta que el framework monta bajo `/_arandu/` y una exención de CSRF de módulo que cubre la ruta de escritura de otro dueño. Fallan antes del tráfico. El orden de registro es observable en el arranque, las rutas, las migrations y el apagado, así que manténgalo intencional en lugar de depender de la iteración de un map.

<a name="boot-the-same-graph-a-server-uses"></a>
## Inicializa el mismo grafo usado por el servidor

Inicializa el mismo grafo que usará el servidor e inspecciona sus rutas y ciclo de vida en vez de probar un controller aislado. Comprueba el orden del boot, el inicio del procesamiento en segundo plano, el resultado de salud y la liberación en orden inverso seguro de cada recurso del módulo.

- Boot no devuelve error para el grafo completo.
- El inventario de rutas identifica el módulo y los handlers esperados.
- El trabajo en segundo plano comienza en Run, no en la construcción.
- El cierre libera los recursos propios en un orden seguro.

<a name="required-and-optional-module-contracts"></a>
## Contratos obligatorios y opcionales de módulo

`foundation.Module` exige `Name()` string y `Routes(*http.Router)`. `foundation.Configure`, `WithProviders` y Create componen el grafo; `Application.HandleCommand`, Boot, Run, Handler y Shutdown son sus dueños. Los contratos opcionales incluyen Bootable, Background, `RendererProvider`, Closable, Diagnostic, Schedulable, Migratable, Health, Ready, Commandable, CSRFExempt y `JobHandlers`, y `foundation.Component` adapta un módulo de `hesape/foundation`.

<a name="interfaces-meet-at-the-root"></a>
## Las interfaces se encuentran en la raíz

Los contratos mantienen los servicios independientes de adapters; la configuración selecciona esos adapters; las rutas exponen controllers ya construidos. La estructura de directorios identifica a cada responsable, mientras el ciclo de la solicitud muestra cuándo se ejecuta el grafo resultante.

- [Configuración](/docs/{{version}}/architecture/configuration)
- [Contratos y adapters](/docs/{{version}}/architecture/contracts-adapters)
- [Module de Foundation](https://pkg.go.dev/github.com/arandu-io/framework/foundation#Module)
