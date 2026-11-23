# Dobles de prueba y tiempo

- [Un doble de prueba implementa el contrato del consumidor](#a-test-double-implements-the-consumer-owned-contract)
- [Inyecte un fake local e inspeccione su llamada](#inject-one-local-fake-and-inspect-its-call)
- [Prefiera las APIs de captura de los componentes publicados](#prefer-capture-apis-for-released-components)
- [No sustituya por un fake el comportamiento que se está demostrando](#do-not-fake-away-the-behavior-being-proved)
- [Ejecute el mismo contrato en el fake y en el adapter real](#run-the-same-contract-against-fake-and-real-adapter)
- [Superficies publicadas de captura y fakes](#released-capture-and-fake-surfaces)
- [Combine dobles específicos con recorridos completos de feature](#combine-narrow-doubles-with-complete-feature-journeys)

> [!WARNING]
> **Limitaciones actuales.** Los fakes viven junto al componente sustituido y no forman un contenedor global de mocks. Dependa de interfaces pequeñas, pase relojes explícitamente cuando el tiempo controle el comportamiento y use fakes solo en la frontera que modelan.

<a name="a-test-double-implements-the-consumer-owned-contract"></a>
## Un doble de prueba implementa el contrato del consumidor

Un fake vive junto a la prueba de la menor interfaz exigida por el consumidor. Registra llamadas tipadas, devuelve resultados configurados y respeta el contexto y la semántica del Grant. Los registros globales de sustitución interfieren en pruebas paralelas y se evitan cuando existe captura por instancia.

<a name="inject-one-local-fake-and-inspect-its-call"></a>
## Inyecte un fake local e inspeccione su llamada

Defina un fake local para el colaborador del servicio, inyéctelo mediante el constructor normal, invoque el objeto probado y verifique el Subject o Grant exacto derivado del contexto y la entrada tipada, sin comparar un mapa sin estructura.

```go example
type invoiceSenderFake struct {
	grant security.Grant
	id    string
	err   error
}

func (f *invoiceSenderFake) Send(_ context.Context, g security.Grant, id string) error {
	f.grant, f.id = g, id
	return f.err
}

func TestFakeRecordsTheCall(t *testing.T) {
	fake := &invoiceSenderFake{}
	var grant security.Grant

	if err := fake.Send(context.Background(), grant, "inv_123"); err != nil {
		t.Fatal(err)
	}
	if fake.id != "inv_123" {
		t.Fatalf("id = %q", fake.id)
	}
}
```

<a name="prefer-capture-apis-for-released-components"></a>
## Prefiera las APIs de captura de los componentes publicados

Use fakes del cliente HTTP bloqueando peticiones inesperadas, Capture de notificaciones por prueba, transportes de correo capturados, adapters en memoria solo donde coincida el contrato y relojes o IDs explícitos cuando estén publicados. Cada captura pertenece a una sola prueba.

- Haga que fallen las llamadas inesperadas.
- Registre entradas tipadas y el orden de las llamadas.
- Respete las cancelaciones y los resultados de error.
- Evite fakes mutables globales en el paquete.

<a name="do-not-fake-away-the-behavior-being-proved"></a>
## No sustituya por un fake el comportamiento que se está demostrando

Un fake que siempre tiene éxito puede ocultar reintentos, rollback, cancelación y fallos parciales. Un adapter en memoria puede diferir en transacciones, orden, bloqueos, TTL y semántica entre procesos. Nunca use el doble para probar el comportamiento que omite intencionalmente.

> [!WARNING]
> Los fakes de parser o constructor no demuestran que un valor configurado llegue al componente en ejecución. Añada un recorrido de feature que cambie la configuración y observe el comportamiento resultante.

<a name="run-the-same-contract-against-fake-and-real-adapter"></a>
## Ejecute el mismo contrato en el fake y en el adapter real

Ejecute un contrato de comportamiento compartido en el fake y el adapter real cuando se prometa sustitución; luego añada una prueba de feature mediante el arranque. Ejecute pruebas en paralelo y con detección de race para revelar capturas globales o callbacks inseguros.

- Las llamadas inesperadas fallan.
- Se preservan el contexto y el Grant.
- Los fallos configurados llegan al objeto probado.
- Las pruebas paralelas no comparten capturas.

<a name="released-capture-and-fake-surfaces"></a>
## Superficies publicadas de captura y fakes

Factory del cliente HTTP proporciona Fake, `FakeSequence`, `StubUrl`, `PreventStrayRequests` y aserciones de envío; `notifications.Capture` devuelve canales y entregas por prueba; los demás dobles implementan directamente interfaces de consumidores.

| Superficie | Contrato publicado |
| --- | --- |
| fakes de Factory del cliente | captura de peticiones salientes |
| `notifications.Capture` | entregas por canal aisladas por prueba |
| fake de interfaz local | frontera del servicio o repositorio de la aplicación |
| contrato compartido | paridad de comportamiento entre adapters |

<a name="combine-narrow-doubles-with-complete-feature-journeys"></a>
## Combine dobles específicos con recorridos completos de feature

Las pruebas HTTP usan dobles en los bordes externos, las de base conservan almacenamiento concreto, los contratos definen sustitución y el resumen elige la menor frontera completa.

- [Pruebas HTTP](/docs/{{version}}/testing/http)
- [Pruebas de base de datos](/docs/{{version}}/testing/database)
- [Contratos](/docs/{{version}}/architecture/contracts-adapters)
- [Resumen de pruebas](/docs/{{version}}/testing/overview)
