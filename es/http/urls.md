# URLs y URLs firmadas

- [Los nombres de ruta son la fuente de direcciones](#route-names-are-the-address-source)
- [Genere y redirija por nombre](#build-and-redirect-by-name)
- [Separe rutas, hosts y firmas](#keep-paths-hosts-and-signatures-separate)
- [No oculte la brecha de integración del framework](#do-not-hide-the-framework-bridge-gap)
- [Altere cada byte firmado](#change-every-signed-byte)
- [Superficies de URL publicadas y su frontera](#released-url-surfaces-and-boundary)
- [Use enlaces según su nivel real de confianza](#use-links-at-their-real-trust-level)

<a name="route-names-are-the-address-source"></a>
## Los nombres de ruta son la fuente de direcciones

Las rutas con nombre son la fuente de direcciones de aplicación. Genere mediante la tabla en vez de concatenar segmentos; firme solo la ruta canónica y valores acotados de query cuando el enlace deba detectar alteraciones.

<a name="build-and-redirect-by-name"></a>
## Genere y redirija por nombre

Genere primero la ruta con nombre y falle ante nombres desconocidos o parámetros incompatibles. Un firmante puede añadir expiración y firma a un enlace de propósito único; verifique antes de que el handler use cualquier valor firmado de la query.

```go example
func SignInvoiceLink(router *fhttp.Router, signer *encryption.Signer, invoiceID string) (string, error) {
	path, err := router.Table().URL("invoices.show", invoiceID)
	if err != nil {
		return "", err
	}

	token := signer.Sign("invoice.show", path, 15*time.Minute)
	verifiedPath, err := signer.Verify("invoice.show", token)
	if err != nil {
		return "", err
	}
	return verifiedPath, nil
}
```

<a name="keep-paths-hosts-and-signatures-separate"></a>
## Separe rutas, hosts y firmas

Mantenga decisiones de host y protocolo en configuración fiable del despliegue. Genere URLs absolutas solo tras normalizar el proxy fiable, preserve el orden canónico de parámetros en la firma y use expiración corta en acciones por correo.

- Nombre las rutas durante el registro.
- Trate los errores de generación inversa de URL en el código de aplicación.
- Firme el destino canónico, no texto proporcionado por el llamador.
- Verifique expiración y firma antes de cualquier mutación.

<a name="do-not-hide-the-framework-bridge-gap"></a>
## No oculte la brecha de integración del framework

Nombre de ruta desconocido, cantidad incorrecta de parámetros, query alterada, firma expirada y host no fiable son rechazos distintos. Nunca use un destino sin firma como fallback si falla la firma. `hesape/routing.SignedRoute` recibe la tabla de rutas de `hesape/routing`, no la de `framework/http.Router`; con el router del framework, obtén la ruta con `Table().URL` y fírmala con `security.Signer`, como en el ejemplo anterior.

> [!WARNING]
> Una firma válida demuestra integridad y expiración, no identidad ni permiso. El handler aún carga el Subject actual, autoriza la acción y garantiza uso único en almacenamiento duradero.

<a name="change-every-signed-byte"></a>
## Altere cada byte firmado

Genere una URL con nombre, firme, verifique y altere un segmento de ruta, valor de query y expiración. Pruebe URLs absolutas en la topología real de proxies fiables e impida que un Host controlado por el atacante se convierta en origen canónico.

- Los parámetros de ruta se escapan y cuentan.
- Un valor alterado falla en la verificación.
- Un enlace expirado falla en la verificación.
- La autorización aún se ejecuta tras verificar la firma.

<a name="released-url-surfaces-and-boundary"></a>
## Superficies de URL publicadas y su frontera

`Route.Name` registra identidad estable y `Routes.URL` genera su ruta. El firmante publicado de security crea y verifica destinos firmados acotados. Estas APIs no realizan búsqueda de cuenta, autorización ni consumo de un solo uso.

<a name="use-links-at-their-real-trust-level"></a>
## Use enlaces según su nivel real de confianza

El enrutamiento posee la plantilla de ruta, los proxies fiables establecen protocolo y host externos, la autorización posee la acción y los códigos desechables poseen el uso único duradero.

- [Rutas](/docs/{{version}}/http/routing)
- [Autorización](/docs/{{version}}/security/authorization)
- [Verificación de correo](/docs/{{version}}/security/email-verification)
