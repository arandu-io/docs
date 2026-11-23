# Conexiones de base de datos

- [Un recorrido de conexión con instrumentación](#one-connection-path-with-instrumentation)
- [Vincule un controlador y abra la URL configurada](#link-one-driver-and-open-the-configured-url)
- [Configure pools y transacciones en la frontera](#configure-pools-and-transactions-at-the-boundary)
- [Interrumpa el inicio si falta el controlador o la conexión no es válida](#fail-boot-on-an-absent-driver-or-invalid-connection)
- [Pruebe el contrato del conector en motores reales](#run-the-connector-contract-against-real-engines)
- [Configuración, DB y Connection](#configuration-db-and-connection)
- [Los repositorios son la puerta de aplicación](#repositories-are-the-application-door)

<a name="one-connection-path-with-instrumentation"></a>
## Un recorrido de conexión con instrumentación

`database.Open` construye el DB de la aplicación a partir de configuración tipada. El handle conoce el dialecto, adapta los marcadores, aplica la política del pool e informa las sentencias al colector configurado. Mantener un único handle encapsulado evita rutas alternativas que ejecuten consultas sin los diagnósticos utilizados por el resto de la aplicación.

<a name="link-one-driver-and-open-the-configured-url"></a>
## Vincule un controlador y abra la URL configurada

Un conector es un módulo Go separado y se registra durante la inicialización del paquete. Haga import en blanco exactamente de los motores que este binario puede hablar. En una aplicación, `foundation.NewDrivers` hace esta llamada con la configuración que produjeron `DATABASE_URL` y los ajustes del pool, responde el handle como `drivers.DB` y cierra el pool en el apagado y al final de cada comando; los repositorios y handlers no deben llamar a otro lector de entorno.

```go example
import (
	"github.com/arandu-io/hesape/database"
	_ "github.com/arandu-io/hesape/database/connectors/sqlite"
)

func OpenDatabase(cfg database.Config) (*database.DB, func(), error) {
	db, closeDB, err := database.Open(cfg)
	if err != nil {
		return nil, nil, err
	}
	return db, closeDB, nil
}
```

<a name="configure-pools-and-transactions-at-the-boundary"></a>
## Configure pools y transacciones en la frontera

La URL selecciona el dialecto registrado y contiene los parámetros de conexión del controlador. Los límites y la duración del pool pertenecen a la configuración de la base de datos, no a repositorios individuales. `database.Transaction` mantiene el trabajo en la transacción asociada al contexto; pase ese contexto por todas las llamadas para que las operaciones anidadas no escapen al pool base.

<a name="fail-boot-on-an-absent-driver-or-invalid-connection"></a>
## Interrumpa el inicio si falta el controlador o la conexión no es válida

Una URL que nombra un driver no enlazado al binario no puede corregirse en runtime; detenga el arranque con el error. Distinga servidor inaccesible, credenciales inválidas, archivo SQLite ausente y configuración malformada. No use fallback hacia otra base con distinta propiedad de datos.

<a name="run-the-connector-contract-against-real-engines"></a>
## Pruebe el contrato del conector en motores reales

La suite de conformidad de conectores demuestra el contrato compartido entre `PostgreSQL`, `MySQL` y SQLite. En la aplicación, compruebe que el controlador seleccionado está disponible, que una consulta parametrizada llega al colector, que el rollback no deja registros y que el cierre libera el pool. Un mock por sí solo no demuestra la vinculación en el dialecto.

<a name="configuration-db-and-connection"></a>
## Configuración, DB y Connection

La configuración y `ParseURL` producen Config; Open crea DB; Transaction delimita trabajo; `AvailableDrivers` enumera conectores registrados. Elija un dueño de configuración, sin llamar otra vez a `database.Load` en el constructor. DB y Connection son infraestructura debajo del repositorio. `ForMigrations` y `ForSchema` exponen interfaces restringidas sin convertir repositorios en herramientas de esquema.

<a name="repositories-are-the-application-door"></a>
## Los repositorios son la puerta de aplicación

La conexión no aplica policy de tenant a SQL arbitrario; los módulos acceden a filas mediante repositorio o terminal de model que exige Grant. Las migraciones operan debajo de esa frontera al actuar sobre esquema antes de peticiones. Elija conectores en bootstrap; ningún handler importa un driver ni abre otro pool.
