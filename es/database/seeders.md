# Seeders

- [Los seeders son operaciones explícitas de la aplicación, identificadas por nombre](#seeders-are-explicit-named-application-operations)
- [Defina dependencias limitadas y comportamiento idempotente](#define-narrow-dependencies-and-idempotent-intent)
- [Ejecute el seeder raíz o uno por su nombre](#run-the-root-or-one-named-seeder)
- [Cree cuentas iniciales mediante los servicios de la aplicación](#create-bootstrap-accounts-through-application-services)
- [Deténgase en la primera operación incompleta](#stop-at-the-first-incomplete-operation)
- [Mantenga secretos fuera de código, historial y lista de procesos](#keep-secrets-out-of-source-history-and-process-lists)
- [Ejecute cada seeder dos veces en una base no vacía](#run-every-seeder-twice-on-a-non-empty-database)
- [Use factories para los datos y servicios para las políticas](#use-factories-for-data-shape-and-services-for-policy)

<a name="seeders-are-explicit-named-application-operations"></a>
## Los seeders son operaciones explícitas de la aplicación, identificadas por nombre

Un seeder implementa `database.Seeder[Deps]` con Name y `Run(context.Context, Deps)`. El proyecto es dueño de Deps, del registro, de la raíz de fallback y del orden de ejecución, así que un seeder solo llega a los servicios que proporciona el provider. El Seed del provider, entregado al builder con `WithSeeders(app.Seed)` en `bootstrap/app.go`, es lo que ejecutan `db:seed` y migrate `--seed`. No hay reflexión ni descubrimiento por el sistema de archivos.

<a name="define-narrow-dependencies-and-idempotent-intent"></a>
## Defina dependencias limitadas y comportamiento idempotente

Declare dependencias limitadas a los servicios de la aplicación, al tenant configurado y a los argumentos restantes opcionales. Registre cada seeder y seleccione `DatabaseSeeder` como predeterminado cuando no haya nombre. La segunda ejecución debe converger, sin duplicar registros ni sustituir credenciales silenciosamente.

- Registre cada seeder nombrado.
- Elija una raíz de fallback.
- Pase explícitamente un tenant configurado.
- Haga observables e idempotentes las ejecuciones repetidas.

<a name="run-the-root-or-one-named-seeder"></a>
## Ejecute el seeder raíz o uno por su nombre

`aru db:seed` ejecuta el seeder predeterminado; un nombre posicional selecciona otro sin distinguir mayúsculas de minúsculas. Los nombres desconocidos fallan mostrando los disponibles, y la opción obsoleta `--class` se rechaza. Los argumentos restantes pertenecen al seeder, accesibles mediante `database.Flag` y `database.Switch`.

```sh
aru db:seed
aru db:seed UserSeeder -e you@example.com -p '<choose-a-password>'
aru db:seed UserSeeder -e you@example.com -p '<choose-a-password>' -r admin
```

<a name="create-bootstrap-accounts-through-application-services"></a>
## Cree cuentas iniciales mediante los servicios de la aplicación

La raíz del starter no crea ninguna cuenta: migrar o ejecutar `aru db:seed` nunca la crea, y en desarrollo siembra las notas de ejemplo, cuyos autores no tienen contraseña utilizable. Cuando una aplicación necesita un usuario de bootstrap o recuperación, ejecuta `UserSeeder` por nombre con correo y contraseña explícitos; añade `-r` admin sólo cuando ese rol sea intencional. La policy de la aplicación y el hashing permanecen en la ruta normal del servicio.

- Mantén las cuentas fuera de `DatabaseSeeder` hasta que la aplicación las añada deliberadamente.
- Ejecuta `UserSeeder` por nombre para crear una cuenta.
- Añade `-r` admin sólo como una decisión explícita de autorización.
- No modifique los roles existentes al recuperar solo la contraseña.

<a name="stop-at-the-first-incomplete-operation"></a>
## Deténgase en la primera operación incompleta

`database.Seed` envuelve un fallo de Run con el nombre del seeder. Si faltan el registro, una dependencia, el tenant o la credencial explícita de `UserSeeder`, falla antes de escribir, y una raíz se detiene en vez de ignorar un requisito incompleto. Producción pide confirmación, salvo que se fuerce explícitamente.

> [!WARNING]
> El framework no envuelve todos los seeders hijos en una única transacción. Use una transacción de aplicación cuando exija todo o nada, o haga cada paso idempotente y observable de forma independiente.

<a name="keep-secrets-out-of-source-history-and-process-lists"></a>
## Mantenga secretos fuera de código, historial y lista de procesos

> [!WARNING]
> Nunca crees una cuenta predeterminada, incluyas una contraseña en un commit, aceptes un tenant elegido por quien llama ni expongas el seeding mediante HTTP. Una contraseña explícita de `UserSeeder` aparece en el historial del shell y en la lista de procesos; borra esa entrada del historial, cambia la contraseña inmediatamente mediante la aplicación y limita la autoridad de bootstrap a un servicio auditado de la aplicación.

Una contraseña de recuperación pasada por línea de comandos debe cambiarse de inmediato mediante la aplicación; no es el recorrido ordinario de seed de despliegue.

<a name="run-every-seeder-twice-on-a-non-empty-database"></a>
## Ejecute cada seeder dos veces en una base no vacía

Migre una base de datos desechable con datos, ejecute dos veces el seeder raíz y cada seeder nombrado, verifique tenant y roles exactos sin duplicados, pruebe nombres desconocidos y opciones obsoletas, provoque el fallo de un hijo y compruebe que la confirmación de producción tiene no como valor predeterminado.

- Las ejecuciones repetidas del seeder raíz convergen.
- Los nombres desconocidos enumeran opciones válidas.
- El fallo de un hijo impide los siguientes.
- Los secretos ausentes no dejan una cuenta parcial.

<a name="use-factories-for-data-shape-and-services-for-policy"></a>
## Use factories para los datos y servicios para las políticas

Las factories construyen datos tipados, las migraciones poseen esquema, las policies llevan autoridad y el despliegue decide cuándo puede ejecutarse una operación inicial.

- [Factories](/docs/{{version}}/database/factories)
- [Migraciones](/docs/{{version}}/database/migrations)
- [Autorización](/docs/{{version}}/security/authorization)
- [Despliegue](/docs/{{version}}/operations/deployment)
