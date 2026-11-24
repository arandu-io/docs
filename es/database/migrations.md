# Migraciones y esquema

- [Una migración es código Go registrado con identidad inmutable](#a-migration-is-registered-go-code-with-immutable-identity)
- [Planifique el despliegue con versiones coexistentes](#design-for-the-mixed-version-rollout-window)
- [Registre un cambio reversible de esquema](#register-one-reversible-schema-change)
- [Inspeccione, simule, aplique y revierta de forma deliberada](#inspect-preview-apply-and-roll-back-deliberately)
- [Entienda resultados de omisión, transacción y rollback](#understand-skip-transaction-and-rollback-outcomes)
- [Separe la autoridad de despliegue de la autoridad de la solicitud](#separate-deployment-authority-from-request-authority)
- [Pruebe ambas direcciones del esquema y el binario anterior](#test-both-schema-directions-and-the-old-binary)
- [Restaure el servicio antes de decidir revertir el esquema](#restore-service-before-choosing-a-schema-down-action)
- [Coordine esquema, despliegue y datos iniciales](#coordinate-schema-deployment-and-seed-data)

<a name="a-migration-is-registered-go-code-with-immutable-identity"></a>
## Una migración es código Go registrado con identidad inmutable

Cada migración implementa `migrations.Migration`, se registra en init y se ordena por `GetName`. El repositorio registra ese nombre exacto y el lote solo después de que Up termine correctamente. El starter usa `arandu_migrations` e importa `database/migrations` en el bootstrap para que todos los comandos vean el registro.

<a name="design-for-the-mixed-version-rollout-window"></a>
## Planifique el despliegue con versiones coexistentes

Documente cómo se comporta el binario anterior tras la llegada del nuevo esquema. Añada columnas anulables o con valor predeterminado y tablas ignoradas por el código antiguo; distribuya los cambios de nombre y las eliminaciones entre versiones. Los comandos de producción requieren una caché RESP compartida para el bloqueo de migraciones.

- Elija un nombre inmutable para la migración.
- Planifique la compatibilidad con el binario anterior.
- Pruebe la gramática de la base elegida.
- Configure bloqueo entre réplicas en producción.

<a name="register-one-reversible-schema-change"></a>
## Registre un cambio reversible de esquema

Implemente Up y Down sobre la conexión de migración, registre el valor y compruebe `ReversibleMigration` al compilar. Incorpore `BaseMigration` para conexión predeterminada, `ShouldRun` normal e intención transaccional publicada.

```go example
type AddLocaleToUsers struct{ migrations.BaseMigration }

func (AddLocaleToUsers) GetName() string { return "20260830_0004_add_locale_to_users" }
func (AddLocaleToUsers) Up(ctx context.Context, conn migrations.Connection) error {
	return conn.Schema().Table(ctx, "users", func(t *schema.Blueprint) { t.String("locale").Nullable() })
}
func (AddLocaleToUsers) Down(ctx context.Context, conn migrations.Connection) error {
	return conn.Schema().Table(ctx, "users", func(t *schema.Blueprint) { t.DropColumn("locale") })
}
func init() { migrations.Register(AddLocaleToUsers{}) }

var _ migrations.ReversibleMigration = AddLocaleToUsers{}
```

<a name="inspect-preview-apply-and-roll-back-deliberately"></a>
## Inspeccione, simule, aplique y revierta de forma deliberada

Genere una migration, inspeccione el estado y la salida de pretend, aplique normalmente o por pasos y revierta solo un último lote revisado o un paso explícito. Use el schema builder para DDL portable y Statement o Select explícitos para un backfill revisado. Las columnas de tiempo usan precisión de microsegundos por defecto, para que dos filas escritas en el mismo segundo conserven su orden; `schema.DefaultTimePrecision` cambia ese valor, y una migration `migrations.WidenTimePrecision` eleva las columnas de una tabla creada con segundos enteros, que es como se ponen al día las tablas de los propios paquetes.

- Use `aru make:migration` con `--create` o `--table`.
- Use `aru migrate:status --pending` y `aru migrate --pretend`.
- Use `aru migrate` o `aru migrate --step`.
- Use `aru migrate:rollback` con una decisión explícita de lote.

<a name="understand-skip-transaction-and-rollback-outcomes"></a>
## Entienda resultados de omisión, transacción y rollback

ShouldRun=false se omite y no se registra, por lo que se evaluará de nuevo. El DDL transaccional varía según la gramática, y `OutsideTransaction` puede dejar el esquema parcialmente modificado. El rollback valida el lote antes de ejecutar: la ausencia de Down sin declaración explícita impide el trabajo; una migración explícitamente Irreversible permanece aplicada y registrada.

> [!WARNING]
> Un fallo de Down en tiempo de ejecución puede detener el rollback después de que entradas reversibles anteriores ya se hayan completado.

<a name="separate-deployment-authority-from-request-authority"></a>
## Separe la autoridad de despliegue de la autoridad de la solicitud

> [!WARNING]
> Las migraciones no reciben Grant: son DDL de despliegue, no peticiones. Dé al pipeline una credencial restringida, vincule valores, revise SQL sin procesar, no exponga secretos en la salida y nunca ejecute migraciones en handlers HTTP ni al iniciar réplicas.

<a name="test-both-schema-directions-and-the-old-binary"></a>
## Pruebe ambas direcciones del esquema y el binario anterior

En cada base de datos desechable compatible, migre dos veces, inspeccione los lotes, ejecute Down y luego Up, revise el SQL simulado, inyecte fallos en Up y Down, pruebe la validación de irreversibilidad, la disputa y la indisponibilidad del bloqueo, e inicie el binario anterior con el esquema aditivo.

- Repetir migrate es idempotente.
- Down y Up restauran el catálogo esperado.
- La disputa por el bloqueo permite una sola aplicación.
- El binario anterior sigue leyendo y escribiendo.

<a name="restore-service-before-choosing-a-schema-down-action"></a>
## Restaure el servicio antes de decidir revertir el esquema

Vuelva a desplegar el binario anterior cuando el esquema aditivo sea compatible. Para un lote reversible con errores, inspeccione el estado y ejecute el rollback específico ya probado. Ante DDL parcial fuera de transacción o cambios irreversibles, inspeccione el catálogo real y la copia de seguridad y escriba una corrección hacia adelante, en lugar de editar los registros del repositorio de migraciones.

`migrate:fresh`, `migrate:reset` y `migrate:refresh` solo se ejecutan con `APP_ENV=dev`, y un `APP_ENV` ausente es producción, porque vacían el schema; el rollback sigue siendo una decisión explícita de lote.

<a name="coordinate-schema-deployment-and-seed-data"></a>
## Coordine esquema, despliegue y datos iniciales

El despliegue posee la ventana de rollout, las actualizaciones eligen versiones, los seeders insertan registros deliberados y las pruebas de base ejercitan ambas direcciones.

- [Despliegue](/docs/{{version}}/operations/deployment)
- [Actualizaciones](/docs/{{version}}/operations/upgrades)
- [Seeders](/docs/{{version}}/database/seeders)
- [Pruebas de base de datos](/docs/{{version}}/testing/database)
