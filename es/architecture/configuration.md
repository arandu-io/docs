# Configuración

- [La configuración es una entrada tipada del inicio](#configuration-is-a-typed-boot-input)
- [Carga una vez y falla antes de atender tráfico](#load-once-and-fail-before-serving-traffic)
- [Usa un responsable para cada ajuste](#use-one-owner-for-each-setting)
- [Rechaza visiblemente valores inválidos y sin efecto en el runtime](#reject-invalid-and-unwired-values-visibly)
- [Rastrea un valor hasta el constructor que lo utiliza](#trace-a-value-into-the-constructor-that-uses-it)
- [Loaders y configuración tipada](#loaders-and-typed-configuration)
- [La configuración elige adapters; el código conserva el contrato](#configuration-chooses-adapters-code-keeps-the-contract)

> [!WARNING]
> **Limitaciones actuales.** El framework lee toda configuración de driver en `bootstrap.Settings` y construye los drivers a partir de ella, y el esqueleto publicado añade solo sus propias secciones. `LOG_SAMPLING` todavía se lee y valida sin llegar al logger, así que pruebe cada variable adicional en el consumidor que debería leerla.

<a name="configuration-is-a-typed-boot-input"></a>
## La configuración es una entrada tipada del inicio

La aplicación lee la configuración una vez, la valida antes de abrir el servidor y pasa valores tipados a los constructores. Los handlers no consultan variables de entorno y los componentes no buscan claves de texto al usarlas.

> [!WARNING]
> El framework es dueño de toda configuración que construye un driver. `bootstrap.LoadSettings` lee las variables de aplicación, HTTP, base de datos, caché, sesión, cola, correo, sistema de archivos y logging en `bootstrap.Settings`, y el `config.Config` del esqueleto incorpora esa configuración junto a sus propias secciones Auth, Geo y Services. No lea una variable del framework por segunda vez en `config/`: dos lectores de un mismo nombre divergen.

<a name="load-once-and-fail-before-serving-traffic"></a>
## Carga una vez y falla antes de atender tráfico

`bootstrap/app.go` carga la configuración, hace que `foundation.NewDrivers` construya cada driver que nombra la configuración, construye el provider de la aplicación sobre esos drivers y entrega los mismos drivers a `foundation.Configure`. Create devuelve la aplicación compuesta o la primera contradicción que encuentre.

```go
cfg, err := appconfig.Load()
if err != nil {
	return nil, err
}

drivers, err := foundation.NewDrivers(cfg.Settings)
if err != nil {
	return nil, err
}

app, err := providers.NewAppServiceProvider(cfg, drivers)
if err != nil {
	_ = drivers.Close()
	return nil, err
}

return foundation.Configure(drivers).WithProviders(app).Create()
```

> [!NOTE]
> Un driver cambia cambiando el entorno, nunca en el código. `CACHE_STORE`, `SESSION_DRIVER`, `QUEUE_CONNECTION`, `MAIL_URL` y `FILESYSTEM_DISK` nombran lo que construye `NewDrivers`, y un valor que nombra un driver que el binario no enlaza detiene el arranque con las dos líneas que lo añaden.

<a name="use-one-owner-for-each-setting"></a>
## Usa un responsable para cada ajuste

| Familia de configuración | Responsable en el esqueleto | Unidad o formato |
| --- | --- | --- |
| `APP_*`, `HTTP_*`, `TRUSTED_PROXIES` y observabilidad | `bootstrap.Settings` del framework | valores tipados |
| `DATABASE_URL` y `DB_*` | `bootstrap.Settings` del framework | una URL; cantidades del pool y segundos |
| `SESSION_DRIVER`, `SESSION_LIFETIME`, `SESSION_SECURE_COOKIE` y `CSRF_TTL` | `bootstrap.Settings` del framework | memory o redis; minutos; booleano; segundos positivos |
| `CACHE_STORE`, `CACHE_PREFIX`, `CACHE_TTL` y `REDIS_*` | `bootstrap.Settings` del framework | memory o redis; prefijo; segundos; una URL y archivos TLS |
| `QUEUE_CONNECTION`, `QUEUE_DEFAULT`, `QUEUE_WORKERS`, `QUEUE_RETRY_AFTER` y `QUEUE_MAX_ATTEMPTS` | `bootstrap.Settings` del framework | database o redis; nombre de la cola; cantidades positivas; retry en segundos |
| `MAIL_URL` y `MAIL_FROM_*` | `bootstrap.Settings` del framework | una URL de transporte; remitente |
| `FILESYSTEM_*` | `bootstrap.Settings` del framework | local, r2 o s3; raíz, bucket, endpoint, región y claves |
| `LOG_LEVEL`, `LOG_FORMAT` y `LOG_SAMPLING` | `bootstrap.Settings` del framework | nivel; json o text; uno de cada N |
| `ARANDU_TENANT_ID`, `AUTH_*`, `GEO_*` y credenciales de servicios | las secciones propias de la aplicación en `config/` | campos tipados |

El archivo `.env` completa valores que el entorno del proceso no haya proporcionado. Por tanto, las variables del entorno de despliegue tienen prioridad. Mantén los secretos fuera del control de versiones y usa el archivo de ejemplo solo como catálogo de entradas posibles, no como prueba de que todas lleguen a un consumidor.

### `APP_ENV` decide lo que producción rechaza

`APP_ENV` es dev, staging o prod. El `.env` copiado de `.env.example` dice dev; un proceso sin `APP_ENV` es prod, lo que desactiva la consola de depuración, la pila en la página de error y los comandos que vacían una base de datos. Una máquina de desarrollo declara `APP_ENV=dev`, y una imagen de contenedor define `APP_ENV=prod` en lugar de dejar la respuesta a una variable sin definir.

### Los nombres de sesión no son alias

La duración de la sesión tiene un único nombre. El framework lee `SESSION_LIFETIME`, en minutos, y `foundation.NewDrivers` construye el store de sesión a partir de él. `SESSION_TTL`, la cantidad de segundos que leía antes el esqueleto, detiene el arranque con el valor de `SESSION_LIFETIME` que debe escribirse en su lugar.

<a name="reject-invalid-and-unwired-values-visibly"></a>
## Rechaza visiblemente valores inválidos y sin efecto en el runtime

Booleanos malformados, recuentos no positivos, URL inválidas, claves ausentes y nombres de store no soportados deben detener el boot con la variable identificada en el error. No sustituyas un valor explícito inválido por el predeterminado: eso hace que el proceso ejecute una configuración distinta de la solicitada por el operador.

> [!WARNING]
> Producción rechaza dos valores en el arranque. `SESSION_SECURE_COOKIE=false` enviaría la sesión, la propia credencial, en texto claro; un despliegue servido por http es de staging y define `APP_ENV=staging`. `MAIL_URL=log://`, que también es el valor de un `MAIL_URL` sin definir, escribiría cada mensaje, incluidos los códigos de restablecimiento de contraseña, en el log; producción define `smtp://`, `smtps://`, `resend://` o `sendgrid://`.

> [!WARNING]
> En framework :frameworkVersion, el nombre y el alcance de la cookie de sesión no son configuraciones: es `arandu_session`, en la ruta /, para el host que respondió. `SESSION_PATH`, `SESSION_DOMAIN` o `SESSION_SAME_SITE` pidiendo otro alcance detienen el arranque, `SESSION_SECURE` lo detiene nombrando `SESSION_SECURE_COOKIE`, y `SESSION_COOKIE` no lo lee nada.

<a name="trace-a-value-into-the-constructor-that-uses-it"></a>
## Rastrea un valor hasta el constructor que lo utiliza

Una prueba del parser solo verifica la sintaxis. Para cada ajuste operativo, añade o encuentra una prueba que cambie el entorno, construya la aplicación y observe el componente consumidor. La prueba de réplicas de sesión y las pruebas de configuración de workers de cola son el modelo: fallan si la composición restaura silenciosamente el valor predeterminado del constructor.

```sh
aru doctor --strict
go test ./tests/Unit ./tests/Feature
```

- Los valores explícitos inválidos hacen fallar el boot e identifican el ajuste.
- Un valor ausente produce el predeterminado documentado.
- Un valor modificado llega al comportamiento observable del componente.
- Los secretos y las credenciales de conexión se ocultan en los errores.

<a name="loaders-and-typed-configuration"></a>
## Loaders y configuración tipada

`bootstrap.LoadSettings` lee el entorno una vez y responde `bootstrap.Settings`, una sección tipada por driver junto a la configuración propia del framework. En el esqueleto, `config.Load` añade las secciones de la aplicación y `config.From` hace lo mismo sobre la configuración que un test proporciona. `foundation.NewDrivers` construye la base de datos, los stores de caché, los locks, el store de sesión, el emisor de CSRF, el rate limiter, la cola, el mailer y los discos a partir de la configuración, y `foundation.Configure` compone la aplicación sobre esos drivers.

- [Configuración de bootstrap de Framework](https://pkg.go.dev/github.com/arandu-io/framework/foundation/bootstrap)
- [Drivers y builder del framework](https://pkg.go.dev/github.com/arandu-io/framework/foundation)
- [Utilidades de configuración de Hesape](https://pkg.go.dev/github.com/arandu-io/hesape/config)

<a name="configuration-chooses-adapters-code-keeps-the-contract"></a>
## La configuración elige adapters; el código conserva el contrato

La elección de store, driver o transporte pertenece a la configuración, y la raíz de composición la construye una vez mediante `foundation.NewDrivers`. Los servicios de la aplicación siguen dependiendo de la interfaz estable que reciben. Esto evita que la configuración se convierta en un service locator y permite que un test construya el mismo grafo con sustitutos explícitos.

- [Módulos y raíz de composición](/docs/{{version}}/architecture/modules-and-composition)
- [Sesión HTTP](/docs/{{version}}/http/session)
- [Contratos y adapters](/docs/{{version}}/architecture/contracts-adapters)
