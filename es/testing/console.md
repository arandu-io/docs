# Pruebas de consola

- [Las pruebas de consola ejecutan la frontera de comando del proyecto](#console-tests-execute-the-project-owned-command-boundary)
- [Despache un comando e inspeccione su efecto](#dispatch-one-command-and-inspect-its-effect)
- [Verifique argumentos, confirmación, efectos y salida](#assert-arguments-confirmation-effects-and-exit)
- [Un nombre registrado no demuestra que la operación se ejecute](#a-registered-name-is-not-proof-that-its-operation-runs)
- [Ejercite el comando delegado mediante Aru y el binario](#exercise-the-delegated-command-through-aru-and-the-binary)
- [Superficies publicadas de pruebas de consola](#released-console-testing-surfaces)
- [Asocie comandos con sus guías de recuperación operativa](#pair-commands-with-their-operational-recovery-guide)

> [!WARNING]
> **Limitaciones actuales.** `Application.HandleCommand` con `WithConsole` es el harness de comandos publicado: recibe el contexto y los buffers del test y responde el código de salida. Las señales y la salida real del binario siguen necesitando un subproceso.

<a name="console-tests-execute-the-project-owned-command-boundary"></a>
## Las pruebas de consola ejecutan la frontera de comando del proyecto

Los comandos que conocen la aplicación los implementa el binario del proyecto y se invocan a través de Aru. La frontera publicada de la aplicación es `Application.HandleCommand(ctx, args)`, que responde el código de salida, y `WithConsole` le da buffers en lugar de los streams del proceso. Un test útil verifica la salida, el código de salida y el efecto duradero; un inventario de nombres de comandos solo prueba que se pueden descubrir.

<a name="dispatch-one-command-and-inspect-its-effect"></a>
## Despache un comando e inspeccione su efecto

Defina una configuración de aplicación aislada, construya la aplicación con el mismo bootstrap que usa el binario, entréguele buffers con `WithConsole`, llame a `HandleCommand` con el comando y los argumentos exactos e inspeccione el código de salida junto al resultado en la base de datos, la cola, el programador o el sistema de archivos.

```go example
func runSchedule(ctx context.Context, t *testing.T, app *foundation.Application) {
	var out, errOut bytes.Buffer
	code := app.WithConsole(&out, &errOut, strings.NewReader("")).
		HandleCommand(ctx, []string{"schedule:run", "billing.close", "--tenant=acme"})
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut.String())
	}
}
```

<a name="assert-arguments-confirmation-effects-and-exit"></a>
## Verifique argumentos, confirmación, efectos y salida

Pruebe argumentos inválidos, valores predeterminados frente a flags, invocación repetida y el efecto real del adaptador seleccionado. `HandleCommand` recibe el contexto que posee el test y los streams que definió `WithConsole`, inicializa la aplicación solo para un comando que la necesita y cierra lo que abrió antes de retornar. Use un subproceso solo cuando la propia frontera del proceso, una señal o la salida real del binario, sea el contrato.

- Verifique el comando delegado y las flags exactas.
- Inspeccione los efectos persistidos.
- Lea la salida y los errores de los buffers de `WithConsole`.
- Ejecute comandos destructivos solo sobre datos desechables.

<a name="a-registered-name-is-not-proof-that-its-operation-runs"></a>
## Un nombre registrado no demuestra que la operación se ejecute

Un comando puede aceptar la sintaxis y aun ignorar un valor de configuración o usar el valor predeterminado incorrecto del constructor. Nombres desconocidos, marcadores de proyecto ausentes, operaciones de producción rechazadas, bloqueos retenidos, dependencias indisponibles y escrituras parciales necesitan su estado real de salida y evidencia.

> [!WARNING]
> Una prueba que comprueba catorce nombres de cola no sustituye un recorrido que inicia, pausa, reintenta o drena un worker. El comportamiento operativo exige el almacenamiento y la frontera de proceso elegidos.

<a name="exercise-the-delegated-command-through-aru-and-the-binary"></a>
## Ejercite el comando delegado mediante Aru y el binario

Para cada familia de comandos, ejecute un efecto correcto, llamada malformada, fallo de dependencia, cancelación y llamada repetida. Para comandos delegados de Aru, ejecute la CLI fijada contra un binario temporal del proyecto y verifique argumentos reenviados y código de salida.

- Los valores predeterminados y sus sustituciones llegan a las opciones del runtime.
- Las operaciones rechazadas devuelven error.
- Las llamadas repetidas preservan el contrato de idempotencia.
- El efecto persistente elegido es observable.

<a name="released-console-testing-surfaces"></a>
## Superficies publicadas de pruebas de consola

Aru es dueña del descubrimiento y el reenvío de nivel superior; `Application.HandleCommand(ctx, args)` int es la frontera publicada de la aplicación para el comando del proyecto, y `WithConsole(out, errOut, in)` define los streams en los que escribe y de los que lee. help y list imprimen la lista de comandos, y un nombre desconocido la imprime y responde 1.

| Superficie | Contrato publicado |
| --- | --- |
| `Application.HandleCommand` | contexto, slice de argumentos y código de salida |
| `Application.WithConsole` | stdout, stderr y stdin de un test |
| estado duradero | principal efecto observable |
| subproceso | señales y la salida real del binario |

<a name="pair-commands-with-their-operational-recovery-guide"></a>
## Asocie comandos con sus guías de recuperación operativa

La referencia de la CLI enumera comandos, las operaciones de cola definen los efectos de los workers, las migraciones definen efectos en el esquema y el despliegue define la ejecución en producción.

- [Referencia de la CLI](/docs/{{version}}/cli/reference)
- [Operaciones de cola](/docs/{{version}}/operations/queues)
- [Migraciones](/docs/{{version}}/database/migrations)
- [Despliegue](/docs/{{version}}/operations/deployment)
