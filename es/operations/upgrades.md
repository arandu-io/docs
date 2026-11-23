# Guía de actualización

- [Actualice el grafo de módulos y la composición perteneciente a la aplicación](#upgrade-the-module-graph-and-application-owned-wiring)
- [Registre el artefacto actual y el estado del esquema](#record-the-current-artifact-and-schema-state)
- [Actualice un conjunto compatible de versiones en una rama](#move-one-compatible-release-set-on-a-branch)
- [Revise API, configuración, esquema y salida generada](#review-api-configuration-schema-and-generated-output)
- [Separe fallos de compilación de cambios silenciosos de comportamiento](#separate-compile-failures-from-silent-behavior-changes)
- [Revise de nuevo cada frontera de identidad y datos](#recheck-every-identity-and-data-boundary)
- [Ejecute verificaciones de comportamiento en el conjunto de versiones deseado](#run-behavior-gates-against-the-target-release-set)
- [Restaure primero el código y revierta el esquema solo con evidencia](#restore-code-first-and-undo-schema-only-with-evidence)
- [Lea las evidencias del release antes de cambiar producción](#read-release-evidence-before-changing-production)

> [!WARNING]
> **Limitaciones actuales.** Las actualizaciones de módulos son explícitas y los cambios de esquema deben ser compatibles con rollback. La identidad de migración y el informe de rollback irreversible aún requieren cuidado; inspeccione el changelog exacto y ejercite rollout con versiones mixtas antes de producción.

<a name="upgrade-the-module-graph-and-application-owned-wiring"></a>
## Actualice el grafo de módulos y la composición perteneciente a la aplicación

Una aplicación Arandu reúne módulos Go ordinarios y código creado desde el esqueleto. Actualizar `go.mod` cambia APIs de bibliotecas; no reescribe bootstrap, configuración, Dockerfile, migraciones ni archivos UI pertenecientes al proyecto. Revise juntos los changelogs etiquetados y el diff del esqueleto.

<a name="record-the-current-artifact-and-schema-state"></a>
## Registre el artefacto actual y el estado del esquema

Comience por un commit verde con archivos de módulo restaurables, backup adecuado al driver, digest de imagen desplegada y lote de migraciones registrado en `arandu_migrations`. Lea todas las etiquetas entre la versión instalada y la deseada, nunca la rama de desarrollo móvil.

- Inventaríe `go list -m` all y `aru version`.
- Registre las versiones del esqueleto y de la imagen.
- Identifique salida Kyse generada y adaptaciones pertenecientes a la aplicación.
- Planifique la ventana de esquema con versiones distintas coexistiendo.

<a name="move-one-compatible-release-set-on-a-branch"></a>
## Actualice un conjunto compatible de versiones en una rama

Fije un conjunto revisado de releases, ordene el grafo y regenere las vistas. Después siga el `UPGRADE.md` del esqueleto, entrada por entrada, desde la release con la que se creó el proyecto: dice qué cambia un proyecto en `bootstrap/app.go`, app/Providers, la configuración, el entorno y los tests. El `UPGRADE.md` de cada módulo, en framework, hesape, kyse, ui, joaju y aru, dice qué deja de compilar o se comporta distinto y qué escribir en su lugar. Porte solo cambios intencionales, mediante revisión, porque `aru new` no es un comando de actualización.

```sh
go get github.com/arandu-io/framework@:frameworkVersion \
  github.com/arandu-io/hesape@:hesapeVersion \
  github.com/arandu-io/kyse@:kyseVersion \
  github.com/arandu-io/hesape/database/connectors/pgx@:pgxVersion \
  github.com/arandu-io/hesape/database/connectors/sqlite@:sqliteVersion \
  github.com/arandu-io/hesape/redis@:redisVersion
go mod tidy
aru view:build
```

<a name="review-api-configuration-schema-and-generated-output"></a>
## Revise API, configuración, esquema y salida generada

Compile tras actualizar pequeños grupos de módulos, busque símbolos eliminados, compare variables de entorno en parser y consumidor, regenere vistas tras cambios de Kyse, preserve bloques UI personalizados, prepare esquemas aditivos y actualice web y worker con el mismo artefacto.

- Revise símbolos de API y conexión de configuración por separado.
- Ejercite los drivers de caché, sesión y cola elegidos en producción.
- Genere de nuevo y compile las vistas.
- Pruebe el binario anterior con el esquema aditivo deseado.

<a name="separate-compile-failures-from-silent-behavior-changes"></a>
## Separe fallos de compilación de cambios silenciosos de comportamiento

Los símbolos eliminados fallan al compilar; una variable de entorno renombrada, un valor predeterminado cambiado, un campo leído pero inactivo o una reordenación de middleware pueden compilar y romper producción. Los imports generados fallan sin build de vistas, y un esquema destructivo puede romper el binario anterior.

> [!WARNING]
> No corrija una actualización añadiendo otro lector local de entorno sin rastrear el valor hasta el consumidor en ejecución. Los parsers de componentes del framework y los lectores de la aplicación pueden leer una variable que ningún constructor activo usa.

<a name="recheck-every-identity-and-data-boundary"></a>
## Revise de nuevo cada frontera de identidad y datos

> [!WARNING]
> Revise de nuevo backends y cookies de sesión, orden de CSRF, propagación de Grant, hashing de tokens, estado OAuth, exposición de diagnósticos, raíces TLS y almacenamiento del rate limit. Actualizar dependencias no demuestra que la composición de la aplicación adopte el camino más seguro.

<a name="run-behavior-gates-against-the-target-release-set"></a>
## Ejecute verificaciones de comportamiento en el conjunto de versiones deseado

Ejecute generación de vistas, formato, vet, pruebas con race, build, doctor, pruebas de contenedor similares a producción, sesión entre dos réplicas, migración dos veces y rollback una, compatibilidad del esquema con el binario anterior, conexión de opciones de cola y pruebas en las etiquetas exactas.

- GOWORK=off `go test -race ./`... pasa.
- `aru doctor --strict` pasa.
- Las vistas generadas corresponden al código fuente.
- Los drivers de producción pasan sus recorridos de feature.

<a name="restore-code-first-and-undo-schema-only-with-evidence"></a>
## Restaure primero el código y revierta el esquema solo con evidencia

Redespliegue la imagen anterior o restaure archivos de módulo y composición de aplicación. Mantenga el esquema aditivo durante el rollback del binario. Revierta una migración de release solo tras inspeccionar el lote y el Down probado; nunca use fresh, reset ni refresh como rollback.

<a name="read-release-evidence-before-changing-production"></a>
## Lea las evidencias del release antes de cambiar producción

Las notas de release definen cambios, las migraciones definen compatibilidad del esquema, el despliegue define rollout del artefacto y la configuración expone riesgos silenciosos de composición.

- [Notas de release](/docs/{{version}}/operations/release-notes)
- [Migraciones](/docs/{{version}}/database/migrations)
- [Despliegue](/docs/{{version}}/operations/deployment)
- [Configuración](/docs/{{version}}/architecture/configuration)
- [Guía de actualización del esqueleto](https://github.com/arandu-io/arandu/blob/:aranduVersion/UPGRADE.md)
