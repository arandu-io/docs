# Desarrollo de paquetes

- [Un paquete es un módulo Go con un punto de entrada explícito](#a-package-is-a-go-module-with-an-explicit-entry-point)
- [Implementa la interfaz obligatoria del módulo](#implement-the-required-module-surface)
- [Añade contratos opcionales de ciclo de vida de forma deliberada](#add-optional-lifecycle-contracts-deliberately)
- [Haz que la falta de dependencias falle durante la construcción](#make-missing-wiring-fail-at-construction)
- [Prueba el paquete mediante su contrato exportado](#test-the-package-through-its-exported-contract)
- [Interfaces obligatorias y opcionales](#required-and-optional-interfaces)
- [Coloca las primitivas en componentes y el estado junto a su responsable](#put-primitives-in-components-and-state-with-its-owner)

<a name="a-package-is-a-go-module-with-an-explicit-entry-point"></a>
## Un paquete es un módulo Go con un punto de entrada explícito

Un paquete instalable tiene su propia ruta de módulo y expone comportamiento mediante API Go comunes. La integración con Framework usa `foundation.Module`, cuyo contrato obligatorio contiene solo un Name estable y el registro de rutas. Las dependencias quedan explícitas en el constructor; ningún contenedor ni reflexión las resuelve detrás de la raíz de composición.

<a name="implement-the-required-module-surface"></a>
## Implementa la interfaz obligatoria del módulo

Name devuelve un identificador estable en minúsculas y Routes recibe el router de la aplicación. Las capacidades opcionales son interfaces estructurales que se descubren solo cuando se implementan; las migrations y el trabajo programado no son métodos que todo módulo deba fingir tener. La aplicación construye el módulo y lo lista en `WithProviders` en `bootstrap/app.go`.

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
## Añade contratos opcionales de ciclo de vida de forma deliberada

Implementa `foundation.Migratable` solo cuando el módulo posea schema; su método exacto es `Migrations() []foundation.Migration`. Implementa `foundation.Schedulable` solo cuando declare tareas del kernel; su método exacto es `Schedule() []foundation.Task`. Mantén vistas, traducciones, policies, handlers y pruebas dentro del módulo responsable. Un paquete declara trabajo; no inicia goroutines sin gestionar durante el import.

<a name="make-missing-wiring-fail-at-construction"></a>
## Haz que la falta de dependencias falle durante la construcción

Rechaza dependencias obligatorias nil en el constructor del módulo y deja que los nombres de ruta duplicados, la configuración inválida o los conflictos de migración aparezcan durante el bootstrap. No sustituyas silenciosamente por adapters en memoria cuando la configuración de producción exija un backend compartido. Una interfaz opcional debe estar ausente, no implementada como una operación vacía que sugiera soporte.

<a name="test-the-package-through-its-exported-contract"></a>
## Prueba el paquete mediante su contrato exportado

Compila un consumidor pequeño que construya el módulo, registre rutas y ejecute una solicitud. Ejecuta pruebas de migración o programación solo para las interfaces que implemente el módulo. Las pruebas del paquete verifican sus repositorios y policies; la aplicación consumidora aún prueba la composición y configuración elegidas.

<a name="required-and-optional-interfaces"></a>
## Interfaces obligatorias y opcionales

`foundation.Module` exige `Name()` string y `Routes(*http.Router)`. `foundation.Migratable` exige `Migrations() []foundation.Migration` y `foundation.Schedulable` exige `Schedule() []foundation.Task`; ambas siguen siendo opcionales y no amplían Module. Los constructores siguen siendo funciones comunes y los handlers de ruta siguen siendo compatibles con `net/http` mediante el router del framework.

<a name="put-primitives-in-components-and-state-with-its-owner"></a>
## Coloca las primitivas en componentes y el estado junto a su responsable

Un cálculo o adapter reutilizable puede vivir en la biblioteca de componentes sin convertirse en un módulo instalable. Una función que posee schema, rutas y Policy pertenece a un paquete. La raíz de composición de la aplicación es donde los adapters del paquete encuentran la configuración concreta, manteniéndolo reutilizable sin ocultar su ejecución.
