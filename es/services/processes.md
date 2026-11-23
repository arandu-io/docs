# Procesos

- [Ejecute un vector de argumentos, nunca una frase de shell](#execute-an-argument-vector-never-a-shell-sentence)
- [Ejecute un proceso hijo acotado](#run-one-bounded-child-process)
- [Controle explícitamente entorno, directorio y salida](#control-environment-directory-and-output-explicitly)
- [Separe fallos de inicio, salida, timeout y captura](#separate-start-exit-timeout-and-output-failures)
- [Use binarios de fixture, no comandos de la máquina](#use-fixture-binaries-instead-of-host-commands)
- [Frontera del componente de procesos publicado](#released-process-component-boundary)
- [Prefiera APIs Go y reserve procesos para herramientas reales](#prefer-go-apis-and-reserve-processes-for-real-tools)

> [!WARNING]
> **Limitaciones actuales.** La API de procesos está publicada y probada, pero no tiene consumidor de primera parte. Nunca interpole entrada no fiable en el shell; pase ejecutable y argumentos explícitamente, defina un plazo de contexto y limite la salida capturada.

<a name="execute-an-argument-vector-never-a-shell-sentence"></a>
## Ejecute un vector de argumentos, nunca una frase de shell

El componente de procesos publicado envuelve `os/exec` con contexto, vectores de argumentos, directorio de trabajo, entorno heredado con sustituciones, salida capturada e información de terminación. Ninguna aplicación actual etiquetada lo consume; esta guía limitada también explicita que stdout y stderr se retienen sin límite de bytes.

<a name="run-one-bounded-child-process"></a>
## Ejecute un proceso hijo acotado

Construya un vector de argumentos, defina directorio, plazos total y de inactividad y solo las sustituciones de entorno necesarias. Env hereda `os.Environ` y reemplaza claves coincidentes; no es una allowlist. Nunca concatene entrada de petición en un comando shell.

```go example
func RunTool(ctx context.Context, workDir, binary, inputPath string, environmentOverrides map[string]string) error {
	result, err := process.NewFactory().NewPendingProcess().
		Path(workDir).
		Env(environmentOverrides).
		Timeout(2 * time.Minute).
		IdleTimeout(30 * time.Second).
		Run(ctx, []string{binary, "--input", inputPath}, nil)
	if err != nil {
		return err
	}
	if _, err := result.Throw(nil); err != nil {
		return err
	}
	return nil
}
```

<a name="control-environment-directory-and-output-explicitly"></a>
## Controle explícitamente entorno, directorio y salida

Resuelva rutas de ejecutables al arrancar, pase cada argumento como una cadena, elija directorio validado, audite el entorno heredado y use Quietly u `OutputHandler` cuando retener salida ilimitada sea inaceptable. Registre código de salida y duración sin secretos de la salida.

- Pase argv sin sh `-c`.
- Defina plazos total y de inactividad.
- Audite las variables de entorno heredadas.
- Evite retener salida ilimitada de herramientas verbosas.

<a name="separate-start-exit-timeout-and-output-failures"></a>
## Separe fallos de inicio, salida, timeout y captura

Ejecutable no encontrado, permiso, fallo al iniciar, salida no cero, señal, cancelación, timeout total y de inactividad son distintos. Un hijo iniciado puede tener efectos externos parciales y la salida capturada puede crecer sin límite de biblioteca.

> [!WARNING]
> Cancelar el contexto termina el proceso hijo, pero no deshace sus efectos en filesystem ni en servicios externos. Un timeout no demuestra que no ocurrió nada.

<a name="use-fixture-binaries-instead-of-host-commands"></a>
## Use binarios de fixture, no comandos de la máquina

Construya pequeños binarios de prueba para éxito, stderr, código de salida elegido, espera, señal, salida excesiva, eco de entorno y limpieza de hijos. Evite aserciones sobre utilidades de shell específicas de plataforma o PATH del usuario.

- Los argumentos con espacios siguen siendo un solo argumento.
- El timeout termina el proceso auxiliar.
- Las sustituciones de entorno prevalecen y los demás valores del padre se heredan.
- La salida verbosa sigue la política de memoria de la aplicación.

<a name="released-process-component-boundary"></a>
## Frontera del componente de procesos publicado

`NewFactory`, `NewPendingProcess`, Path, Env, Timeout, `IdleTimeout`, Run, Quietly y `ProcessResult` forman la frontera publicada. Autorización, reducción de entorno, política de memoria de salida y disponibilidad de herramienta siguen siendo explícitas en la aplicación.

| Superficie | Contrato publicado |
| --- | --- |
| `NewPendingProcess` | constructor de procesos con vector de argumentos |
| Path / Env | sustituciones de directorio y entorno heredado |
| Timeout / `IdleTimeout` / Run | ejecución vinculada al contexto |
| `ProcessResult` | salida íntegra capturada y metadatos de terminación |

<a name="prefer-go-apis-and-reserve-processes-for-real-tools"></a>
## Prefiera APIs Go y reserve procesos para herramientas reales

El contexto posee cancelación, el filesystem posee archivos, el despliegue posee herramientas instaladas y los contratos mantienen procesos hijos detrás de un adapter específico de aplicación.

- [Contexto](/docs/{{version}}/architecture/context)
- [Sistema de archivos](/docs/{{version}}/services/filesystem)
- [Despliegue](/docs/{{version}}/operations/deployment)
- [Contratos](/docs/{{version}}/architecture/contracts-adapters)
