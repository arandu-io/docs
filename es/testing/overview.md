# Resumen de pruebas

- [Las pruebas inician la misma aplicación sin servir un puerto](#tests-boot-the-same-application-graph-without-serving-a-port)
- [Genere las vistas y aísle el estado externo](#generate-views-and-isolate-external-state)
- [Ejecute las verificaciones del proyecto publicado](#run-the-released-project-gates)
- [Demuestre el comportamiento en la menor frontera completa](#prove-behavior-at-the-narrowest-complete-boundary)
- [Continúe con HTTP, consola, base o detalles de contratos](#choose-http-console-database-or-contract-detail-next)

> [!WARNING]
> **Limitaciones actuales.** Los helpers de prueba se publican en varios módulos, sin una fachada única. Componga pruebas Go ordinarias con el helper responsable de la frontera y deje explícita la preparación de cada paquete.

<a name="tests-boot-the-same-application-graph-without-serving-a-port"></a>
## Las pruebas inician la misma aplicación sin servir un puerto

Arandu usa pruebas Go ordinarias y helpers de aplicación. Una prueba puede construir un servicio, iniciar la aplicación compuesta y despachar su Handler o ejecutar un adapter de comando. La mejor frontera es la menor que aún demuestra el comportamiento, sin un framework de pruebas separado.

<a name="generate-views-and-isolate-external-state"></a>
## Genere las vistas y aísle el estado externo

Ejecute `aru view:build` en un checkout limpio, pues los paquetes de vistas generados pueden estar en gitignore. Use fixtures desechables de base y almacenamiento compartido, transportes fake explícitos, relojes controlados donde estén disponibles y restauración del entorno mediante `t.Setenv`.

- Mantenga reproducibles las verificaciones de release con GOWORK=off.
- Aísle base y tenant por prueba.
- Impida peticiones de red inesperadas.
- Use `t.Cleanup` para todo recurso perteneciente al proceso.

<a name="run-the-released-project-gates"></a>
## Ejecute las verificaciones del proyecto publicado

Ejecute formato, generación de vistas, vet, tests con race, build y el Doctor estricto desde la raíz del proyecto. Un parser o un test unitario en verde no basta cuando el comportamiento depende del bootstrap, de un driver seleccionado o de salida generada. `aru test` ejecuta `go test`, y `aru test --e2e` ejecuta los recorridos de tests/E2E, que están detrás de la build tag e2e, conducen un navegador real contra la aplicación en un puerto propio y se omiten, en la máquina de quien desarrolla, donde no hay navegador.

```sh
aru view:build
unformatted=$(find . -name '*.go' -not -path '*/testdata/*' -not -name '*.kyse.go' -exec gofmt -l {} +)
test -z "$unformatted" || { echo "$unformatted"; exit 1; }
GOWORK=off go vet ./...
GOWORK=off go test -race -count=1 ./...
GOWORK=off go build ./...
aru doctor --strict
```

<a name="prove-behavior-at-the-narrowest-complete-boundary"></a>
## Demuestre el comportamiento en la menor frontera completa

Use pruebas unitarias para policies puras y transformaciones, contratos para cada adapter, features mediante la composición real y recorridos similares a producción para estado entre réplicas. Verifique resultados y rechazos, no solo el parsing de valores o el retorno de constructores.

- Un cambio de configuración llega al comportamiento observable.
- Cada adapter de red tiene una prueba contra peticiones inesperadas.
- El aislamiento del tenant se ejercita en el almacenamiento.
- Las pruebas de race cubren estado compartido y cierre.

<a name="choose-http-console-database-or-contract-detail-next"></a>
## Continúe con HTTP, consola, base o detalles de contratos

Las pruebas HTTP despachan rutas, las de consola ejercitan comandos delegados, las de base demuestran transacciones y Grants, y los dobles capturan fronteras de componentes sin estado global.

- [Pruebas HTTP](/docs/{{version}}/testing/http)
- [Pruebas de consola](/docs/{{version}}/testing/console)
- [Pruebas de base de datos](/docs/{{version}}/testing/database)
- [Dobles de prueba](/docs/{{version}}/testing/test-doubles)
