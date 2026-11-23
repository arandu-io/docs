# Strings

- [Primero la biblioteca estándar](#standard-library-first)
- [Convenciones de nombres compartidas](#shared-naming-conventions)
- [Límites y slugs vacíos](#boundaries-and-empty-slugs)

<a name="standard-library-first"></a>
## Primero la biblioteca estándar

Usa la biblioteca estándar para operaciones comunes con strings. strings proporciona búsqueda, división, sustitución, recorte, comparación y conversión de mayúsculas/minúsculas sin añadir convenciones de la aplicación. strconv interpreta y formatea valores primitivos; unicode y utf8 gestionan caracteres y codificación. Mantener visibles esas llamadas facilita reconocer el comportamiento familiar de Go.

Usa str cuando varios paquetes deban coincidir en una convención de nombres que la biblioteca estándar no define: convertir un campo en etiqueta, derivar un identificador de estilo Go, flexionar un sustantivo inglés, convertir texto latino soportado en ASCII o producir un slug de dirección. Una transformación usada por una sola función permanece con ella en vez de convertirse en un auxiliar global.

<a name="shared-naming-conventions"></a>
## Convenciones de nombres compartidas

Studly y Camel construyen identificadores desde palabras separadas por guion, subrayado o espacio. Snake y Kebab hacen explícita la conversión opuesta. Headline transforma un valor parecido a un campo en una etiqueta; framework :frameworkVersion usa exactamente esa función cuando la validación hace legible el nombre de un campo. Así, Password Confirmation es la salida compartida, no un texto elegido por separado por cada formulario.

Slug recibe su separador explícitamente porque Go no tiene argumentos predeterminados. Las funciones de paquete son la forma más clara para una operación. `str.Of` está disponible cuando una secuencia de transformaciones describe un valor, y `ToString` devuelve el string final; ambos caminos ejecutan el mismo comportamiento del paquete.

```go example
func NamingExamples() (string, string, string) {
	return str.Studly("invoice_item"), // InvoiceItem
		str.Headline("password_confirmation"), // Password Confirmation
		str.Slug("Quarterly Report 2026", "-") // quarterly-report-2026
}
```

> [!NOTE]
> Mantén el valor original como fuente de verdad. Un identificador derivado puede volver a producirse; almacenar varias variantes de capitalización crea oportunidades independientes de divergencia.

<a name="boundaries-and-empty-slugs"></a>
## Límites y slugs vacíos

ASCII usa una tabla de conversión documentada, no un transliterador general. Cubre Latin-1 Supplement y Latin Extended-A; los sistemas de escritura no soportados se descartan en vez de adivinarse. Slug comienza por esa conversión, por lo que un título escrito solo con caracteres no soportados puede producir un slug vacío. Comprueba el resultado y usa una clave estable emitida por separado cuando pueda haber una salida vacía.

Un slug es presentación, no identidad, autorización ni lista de permitidos de la base de datos. Los identificadores de tenant, claves de almacenamiento, parámetros de ruta y elecciones de columnas aún deben validarse en el subsistema que los consume. Las frases en lenguaje humano pertenecen a catálogos de traducción: cambiar límites de palabras o capitalización no traduce significado.
