# Almacenamiento de archivos

- [Disk controla el tenant antes de que el adapter reciba una ruta](#a-disk-owns-tenancy-before-an-adapter-sees-a-path)
- [Cree un disco y transmita un objeto](#create-a-disk-and-stream-an-object)
- [Enumere, mueva, descargue y firme acceso temporal](#list-move-download-and-sign-temporary-access)
- [Rechace una identidad insegura en lugar de reescribirla](#reject-unsafe-identity-instead-of-rewriting-it)
- [Use las aserciones de disco sobre el adapter configurado](#use-the-disk-assertions-against-the-configured-adapter)
- [Operaciones de Disk y capacidad del adapter](#disk-operations-and-adapter-capability)
- [Storage no es el filesystem común de la aplicación](#storage-is-not-the-application-s-plain-filesystem)

> [!WARNING]
> **Limitaciones actuales.** `foundation.NewDrivers` construye `drivers.Disks` a partir de `FILESYSTEM_DISK`, local, r2 o s3, así que toda aplicación publicada recibe el Disk con alcance de tenant; el recurso de ejemplo del esqueleto no guarda archivos, y ningún recorrido propio publicado sube y sirve uno. Verifique Put, Get y Serve en la aplicación consumidora antes de tratarlo como un servicio de la aplicación.

<a name="a-disk-owns-tenancy-before-an-adapter-sees-a-path"></a>
## Disk controla el tenant antes de que el adapter reciba una ruta

Disk es la frontera de aplicación sobre Adapter. Cada lectura y escritura recibe Grant, deriva el prefijo del tenant, valida la clave lógica y solo entonces entrega la ruta final al adapter. Ese orden impide que un driver olvide el aislamiento garantizado por otro.

<a name="create-a-disk-and-stream-an-object"></a>
## Cree un disco y transmita un objeto

`foundation.NewDrivers` construye `drivers.Disks` a partir de `FILESYSTEM_DISK`: local escribe bajo `FILESYSTEM_ROOT`, y r2 y s3 se construyen sobre el protocolo S3 a partir del bucket, endpoint, región y claves, con una variable ausente rechazada al arrancar. `Disk("")` responde el disco predeterminado. Put acepta un `io.Reader` para que una subida no tenga que quedar en memoria; Get devuelve la información y el contenido del archivo por el mismo límite de tenant. La clave de la aplicación sigue siendo lógica y nunca incluye un prefijo de tenant.

```go example
func StoreReport(ctx context.Context, drivers *foundation.Drivers, g security.Grant, body io.Reader) error {
	disk, err := drivers.Disks.Disk("")
	if err != nil {
		return err
	}
	if err := disk.Put(ctx, g, "reports/q1.pdf", body, "application/pdf"); err != nil {
		return err
	}
	file, err := disk.Get(ctx, g, "reports/q1.pdf")
	if err != nil {
		return err
	}
	defer file.Body.Close()
	return nil
}
```

<a name="list-move-download-and-sign-temporary-access"></a>
## Enumere, mueva, descargue y firme acceso temporal

Files y Directories enumeran el prefijo de un tenant; Copy y Move mantienen ambas claves dentro de él; Response y Download escriben respuestas HTTP; `TemporaryURL` delega acceso con expiración solo cuando el adapter lo proporciona. `PutFile` y `PutFileAs` reciben uploads preservando la misma política de rutas.

`filesystem.Serve` escribe un archivo guardado en una respuesta desde una ruta cuya Policy produjo el Grant. El tipo de contenido viene de los metadatos guardados, X-Content-Type-Options es nosniff y la Content-Security-Policy es default-src 'none'; sandbox, de modo que un archivo renderizado no ejecuta script. SVG, HTML y XML pueden llevar script y salen como adjunto, salvo que `ServeOptions.InlineActiveContent` pida mostrarlos en línea, aún en sandbox.

<a name="reject-unsafe-identity-instead-of-rewriting-it"></a>
## Rechace una identidad insegura en lugar de reescribirla

Un tenant de Grant vacío, un tenant con separador, una clave que escapa, un objeto ausente o una URL temporal no soportada devuelven un error. No normalice un tenant inseguro en otro identificador: dos identificadores originales podrían fundirse en un mismo directorio. Cierre el `File.Body` devuelto por Get y todo stream devuelto por `ReadStream`, y propague los fallos parciales de escritura. Una escritura bajo un Grant de solo lectura, emitido mientras un operador solo ve una cuenta, se rechaza con `auth.ErrReadOnly`.

<a name="use-the-disk-assertions-against-the-configured-adapter"></a>
## Use las aserciones de disco sobre el adapter configurado

`AssertExists`, `AssertMissing`, `AssertCount` y `AssertDirectoryEmpty` verifican por la frontera pública Disk. Añada pruebas entre tenants con claves lógicas iguales, traversal de rutas, cierre de streams y expiración de URLs temporales. Ejecute el contrato de cada driver de producción, sin confiar solo en la prueba local.

<a name="disk-operations-and-adapter-capability"></a>
## Operaciones de Disk y capacidad del adapter

Put, Get, `ReadStream`, `WriteStream`, Delete, Copy, Move, Stat, checksum, visibilidad y métodos de directorio componen el almacenamiento. `ProvidesTemporaryURLs` y `ProvidesTemporaryUploadURLs` revelan capacidades opcionales antes del uso. `GetAdapter` existe para integración; las operaciones normales permanecen en Disk para no eludir Grant.

<a name="storage-is-not-the-application-s-plain-filesystem"></a>
## Storage no es el filesystem común de la aplicación

Los objetos del tenant usan Disk. Configuración de proceso, assets incorporados, vistas generadas y otros archivos de aplicación usan el filesystem común, pues no son datos del cliente. Una fila de base puede poseer la clave del objeto, pero guardar la ruta absoluta del adapter expone detalles de despliegue y dificulta cambiar drivers.
