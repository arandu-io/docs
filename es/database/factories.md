# Factories de modelos

- [Las factories construyen filas tipadas para pruebas y seeds](#factories-build-typed-rows-for-tests-and-seeds)
- [Configure el model antes de definir su factory](#configure-the-model-before-defining-its-factory)
- [Genere datos deterministas o cree registros autorizados](#make-deterministic-data-or-create-authorized-rows)
- [Componga cantidad, estado, secuencia y relaciones](#compose-count-state-sequence-and-relationships)
- [Considere la persistencia parcial y los errores de callbacks](#account-for-partial-persistence-and-callback-errors)
- [La creación aún lleva el Grant del llamador](#creation-still-carries-the-caller-s-grant)
- [Demuestre repetibilidad, aislamiento y propiedad por tenant](#prove-repeatability-isolation-and-tenant-ownership)
- [Use factories en pruebas enfocadas y seeders idempotentes](#use-factories-from-focused-tests-and-idempotent-seeders)

<a name="factories-build-typed-rows-for-tests-and-seeds"></a>
## Las factories construyen filas tipadas para pruebas y seeds

`aru make:factory` escribe una factory tipada para una entidad, como `factories.Users(db)` en `database/factories/UserFactory.go`, y `aru model:build` mantiene al día todo lo que queda fuera del bloque custom. Guarda un valor predeterminado válido y compone estados deterministas sobre la `factories.Factory` no genérica. La definición devuelve User, los estados reciben `*User` y los resultados son las mismas filas que devuelven los terminales de consulta. Los campos con nombre incorrecto fallan al compilar en vez de desaparecer en un mapa.

<a name="configure-the-model-before-defining-its-factory"></a>
## Configure el model antes de definir su factory

Declare la entidad y su tabla en app/Models y ejecute `aru model:build`. La factory tipada se construye sobre la consulta generada de la entidad: la tabla aporta la configuración de clave y la convención de tenant, y el argumento db aporta la conexión; no hay valor de model que construir a mano.

- Declare la entidad y su tabla.
- Mantenga una definición predeterminada válida.
- Use valores Faker deterministas.
- Proporcione Grant solo a operaciones de persistencia.

<a name="make-deterministic-data-or-create-authorized-rows"></a>
## Genere datos deterministas o cree registros autorizados

`MakeOne` no hace I/O ni exige contexto o Grant; `CreateOne` escribe mediante el model con ambos. La semilla predeterminada es 1 y `Seed(n)` hace reproducibles los fallos.

```go
func defineUser(f faker.Faker) models.User {
	return models.User{Name: f.Name(), Email: f.Unique().Email(), Active: true}
}

made, err := factories.Users(db).MakeOne()
if err != nil { return err }
stored, err := factories.Users(db).Count(3).Create(ctx, grant)
if err != nil { return err }
```

<a name="compose-count-state-sequence-and-relationships"></a>
## Componga cantidad, estado, secuencia y relaciones

Count y Seed definen volumen y repetibilidad; State ejecuta en orden; Sequence se repite en ciclos; `AfterMaking` y `AfterCreating` añaden callbacks. La factory tipada reenvía esos métodos; Has, que crea hijos después de los padres, y `ForParent`, que crea primero un padre, están en la factory del núcleo y se alcanzan desde un método en el bloque custom. Cada modificador devuelve una nueva factory.

- `Count(n)` y `Seed(n)` controlan un lote reproducible.
- State y Sequence componen cambios tipados.
- `AfterCreating` recibe contexto y Grant.
- Has y `ForParent` requieren un callback explícito de asociación.

<a name="account-for-partial-persistence-and-callback-errors"></a>
## Considere la persistencia parcial y los errores de callbacks

Make puede fallar al asignar valores al modelo. Create guarda secuencialmente: un error posterior en Save o en un callback puede dejar registros anteriores persistidos, porque Factory no envuelve todo el lote en una transacción. Un Count negativo se convierte en cero, y un callback de relación válido pero incorrecto genera datos engañosos.

> [!WARNING]
> Envuelva la creación mediante la factory en una transacción controlada por el llamador cuando importe la atomicidad. Verifique las claves de las relaciones, sin suponer que un callback que compila conectó los registros previstos.

<a name="creation-still-carries-the-caller-s-grant"></a>
## La creación aún lleva el Grant del llamador

> [!WARNING]
> Una factory no elude autorización. Create y callbacks de relación usan el Grant proporcionado; nunca invente tenant ni use autoridad de sistema en peticiones. Los seeds de producción no contienen contraseñas fijas, tokens ni datos personales reales.

<a name="prove-repeatability-isolation-and-tenant-ownership"></a>
## Demuestre repetibilidad, aislamiento y propiedad por tenant

Demuestre que la misma semilla repite y otra difiere, las factories base no cambian, Make no llama a la base, Create escribe tenant y cantidad esperados, los Grants denegados rechazan, las claves de relación coinciden y el fallo parcial de callback se trata deliberadamente.

- Los valores con semilla se repiten.
- Make no realiza E/S.
- Create conserva la autoridad del tenant.
- Los modificadores no se filtran entre pruebas.

<a name="use-factories-from-focused-tests-and-idempotent-seeders"></a>
## Use factories en pruebas enfocadas y seeders idempotentes

Los modelos definen la persistencia, las relaciones definen los vínculos, los seeders controlan operaciones nombradas de inicialización y las pruebas de base de datos aíslan cada conjunto de datos creado.

- [Modelos](/docs/{{version}}/database/models)
- [Relaciones](/docs/{{version}}/database/relationships)
- [Seeders](/docs/{{version}}/database/seeders)
- [Pruebas de base de datos](/docs/{{version}}/testing/database)
