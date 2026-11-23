# Contribución

- [Las comprobaciones del repositorio definen un cambio revisable](#repository-gates-define-a-reviewable-change)
- [Lee el contrato del repositorio antes de editar](#read-the-repository-contract-before-editing)
- [Haz un cambio concreto con su prueba de regresión](#make-one-focused-change-and-its-regression)
- [Demuestra que cada nueva protección puede fallar](#prove-every-new-guard-can-fail)
- [Asocia cada cambio con un motivo](#keep-one-change-attributable-to-one-reason)

> [!WARNING]
> **Limitaciones actuales.** El contrato de contribución se verifica mediante la CI y la protección de estructura de la versión etiquetada del propio repositorio. No tiene un consumidor independiente con tag porque rige cambios en ese repositorio, no en un paquete de runtime.

<a name="repository-gates-define-a-reviewable-change"></a>
## Las comprobaciones del repositorio definen un cambio revisable

Formato, vet y pruebas con detector de carreras son la base común, complementada por las comprobaciones de generación, dependencias, vulnerabilidades y conformidad de cada repositorio. Los comandos se ejecutan localmente y en CI; una contribución informa de sus resultados reales, sin tratar los nombres de los comandos como aprobación.

```bash
unformatted=$(find . -name '*.go' -not -path '*/testdata/*' -not -name '*.kyse.go' -exec gofmt -l {} +)
test -z "$unformatted" || { echo "$unformatted"; exit 1; }
GOWORK=off go vet ./...
GOWORK=off go test -race -count=1 -p 1 ./...
```

<a name="read-the-repository-contract-before-editing"></a>
## Lee el contrato del repositorio antes de editar

Usa la versión de Go y la política de workspace que declara el repositorio, comprueba la responsabilidad sobre archivos generados e identifica las pruebas específicas del comportamiento modificado. Los informes de seguridad siguen la política de comunicación privada; una issue pública no es el primer lugar para divulgar una vulnerabilidad.

<a name="make-one-focused-change-and-its-regression"></a>
## Haz un cambio concreto con su prueba de regresión

Reproduce el fallo, añade una prueba que observe el comportamiento incorrecto, implementa la corrección coherente más pequeña y vuelve a ejecutar la prueba. Actualiza la documentación pública cuando cambie la API observable. Deja fuera del patch el formato no relacionado y los cambios innecesarios de archivos generados para mantener clara la revisión.

<a name="prove-every-new-guard-can-fail"></a>
## Demuestra que cada nueva protección puede fallar

Introduce la infracción, ejecuta la protección específica, observa el fallo, restaura el código y observa el éxito. Una prueba de escaneo también comprueba haber inspeccionado un número plausible de objetivos, para que un patrón obsoleto no pase sin encontrar nada. Excluye las plantillas de vistas y fixtures deliberadamente inválidas de las comprobaciones genéricas de formato por sus motivos documentados.

<a name="keep-one-change-attributable-to-one-reason"></a>
## Asocia cada cambio con un motivo

Separa los cambios mecánicos de los de comportamiento cuando facilite revisar las evidencias. El mensaje de commit explica el motivo que no puede reconstruirse desde el diff. Las instrucciones locales, la responsabilidad sobre archivos generados, las reglas de release y los procedimientos para reportar problemas de seguridad siguen vigentes además de los comandos comunes.
