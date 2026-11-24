# Localización

- [Un traductor, una decisión de idioma por petición](#one-translator-one-locale-decision-per-request)
- [Registre catálogos y resuelva una clave](#register-catalogues-and-resolve-a-key)
- [Sustituya valores, elija plurales y añada namespaces](#replace-values-choose-plurals-and-add-namespaces)
- [Un fallback sin resolución devuelve la clave](#an-unresolved-fallback-returns-the-key)
- [Pruebe el origen del idioma y el catálogo](#test-the-locale-source-and-the-catalogue)
- [Operaciones del traductor](#translator-operations)
- [Las rutas y el texto tienen contratos de idioma distintos](#routes-and-prose-have-separate-language-contracts)

<a name="one-translator-one-locale-decision-per-request"></a>
## Un traductor, una decisión de idioma por petición

Translator carga frases por clave desde rutas registradas y mantiene grupos en caché de forma segura para peticiones concurrentes. El idioma pertenece al contexto de petición, elegido una vez por middleware o ámbito explícito. El código posterior lee esa decisión sin reinterpretar URL ni header de idioma.

<a name="register-catalogues-and-resolve-a-key"></a>
## Registre catálogos y resuelva una clave

Construya el traductor una vez a partir del directorio de catálogos con `translation.NewFromDirectory`, nombrando el locale de respuesta y el de fallback, y entréguelo al builder con `WithTranslator`. Cada petición lleva entonces el catálogo junto al nombre de la aplicación, de modo que `view.Page` responde Translator y Locale y un componente dibuja sus propias palabras en el idioma de la petición. Get comprueba la cadena de fallback configurada; cuando ni el catálogo pedido ni el de fallback resuelven el elemento, se ejecuta el handler de clave ausente y se devuelve la propia clave sin resolver. El esqueleto ya lo hace: su catálogo es `lang/en.json`, y `bootstrap/app.go` entrega el traductor al builder con `WithTranslator`.

```go example
func WithTranslations(drivers *foundation.Drivers) (*foundation.Application, error) {
	translator, err := translation.NewFromDirectory(os.DirFS("."), "lang", "en", "en")
	if err != nil {
		return nil, err
	}

	return foundation.Configure(drivers).WithTranslator(translator).Create()
}
```

<a name="replace-values-choose-plurals-and-add-namespaces"></a>
## Sustituya valores, elija plurales y añada namespaces

Get sustituye placeholders con nombre, Choice selecciona el plural por número, Has verifica cobertura y `AddNamespace` asocia traducciones de paquete a su propio hint. `InLocale` sirve para trabajo fuera de la petición; cuando middleware ya proporciona idioma, use esa fuente única en lugar de anidar otra elección.

<a name="an-unresolved-fallback-returns-the-key"></a>
## Un fallback sin resolución devuelve la clave

El fallback mantiene utilizable una página cuando el idioma está incompleto, pero no inventa texto. Si el elemento sigue ausente y el handler de claves no proporciona sustitución, Get devuelve la clave, exponiendo el defecto del catálogo. Las pruebas de paridad deben comparar conjuntos de claves y rechazar valores copiados del idioma fuente, salvo excepciones breves y justificadas. Un idioma no compatible nunca debe convertirse en una ruta de lectura arbitraria de archivos.

<a name="test-the-locale-source-and-the-catalogue"></a>
## Pruebe el origen del idioma y el catálogo

Envíe peticiones por todos los prefijos de idioma compatibles y compruebe interfaz traducida y rutas estables. Ejercite sustituciones y límites de plural, luego busque claves ausentes, adicionales y valores sospechosamente idénticos. La prueba también debe verificar una cantidad plausible de claves para que un loader roto no pase sin comprobar nada.

<a name="translator-operations"></a>
## Operaciones del traductor

`NewFromDirectory` construye el traductor de la aplicación y Locales lista los locales que trae su directorio, que es la lista contra la que negocia el middleware de locale. New, Get, Choice, Has y `HasForLocale` forman el camino normal de lectura. `AddPath` y `AddJSONPath` registran catálogos de la aplicación; `AddNamespace` registra catálogos de paquetes; `DetermineLocalesUsing` personaliza el orden de fallback; `HandleMissingKeysUsing` observa las claves sin resolver. `SetLocale` cambia el valor predeterminado y no sustituye el contexto de locale por petición.

<a name="routes-and-prose-have-separate-language-contracts"></a>
## Las rutas y el texto tienen contratos de idioma distintos

La interfaz, fechas, números y frases del producto pertenecen a la traducción. Identificadores de API, ejemplos Go y cuerpo del manual pueden permanecer deliberadamente en el idioma fuente cuando ese sea el contrato publicado. Los helpers de URL deben recibir el idioma elegido por el contexto de rutas, evitando discrepancias entre enlaces, middleware y catálogos.
