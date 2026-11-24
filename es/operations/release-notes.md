# Notas de release

- [Una nota de release describe una versión inmutable](#a-release-note-describes-an-immutable-version)
- [Comience por la etiqueta inmutable anterior](#start-from-the-previous-immutable-tag)
- [Describa el cambio observable](#describe-the-observable-change)
- [Agrupe adiciones, cambios, correcciones y eliminaciones](#group-additions-changes-fixes-and-removals)
- [Deténgase cuando la evidencia pertenezca a otra revisión](#stop-when-evidence-belongs-to-another-revision)
- [No publique secretos ni convierta una corrección en instrucciones de ataque](#do-not-publish-secrets-or-weaponize-a-fix)
- [Calcule la diferencia de API antes de crear la etiqueta](#compute-the-api-difference-before-tagging)
- [Corrija con una nueva versión](#correct-with-a-new-version)
- [Las guías de actualización describen transiciones más largas](#upgrade-guides-carry-the-longer-path)

<a name="a-release-note-describes-an-immutable-version"></a>
## Una nota de release describe una versión inmutable

Cada módulo publicado está por debajo de 1.0 y registra cambios incompatibles en su changelog. Una etiqueta es inmutable: corregir un artefacto publicado exige una nueva versión, nunca mover la etiqueta anterior. Así se recupera el mismo código en builds y revisiones de incidentes.

<a name="start-from-the-previous-immutable-tag"></a>
## Comience por la etiqueta inmutable anterior

Un candidato a release exige revisión limpia, verificaciones aprobadas, entrada de changelog y etiqueta anterior conocida. Compare esa etiqueta con el commit exacto previsto; comparar dos ramas móviles produce notas que nadie puede reproducir después.

<a name="describe-the-observable-change"></a>
## Describa el cambio observable

Describa qué debe cambiar el consumidor, el comportamiento diferente y la primera versión del módulo que lo incluye. Identifique con precisión símbolos eliminados o renombrados y su sustitución. Los comandos siguientes ilustran una inspección; no son el proceso de release ni prueba de verificaciones aprobadas. Los scripts y políticas del repositorio siguen siendo la autoridad.

```bash
PREVIOUS_TAG=$(git describe --tags --abbrev=0 HEAD^)
git diff --stat "$PREVIOUS_TAG"...HEAD
GOWORK=off go vet ./...
GOWORK=off go test -race -count=1 -p 1 ./...
```

<a name="group-additions-changes-fixes-and-removals"></a>
## Agrupe adiciones, cambios, correcciones y eliminaciones

Escriba para el consumidor y vincule cada nota de migración al módulo afectado. Destaque configuración, migraciones de datos, seguridad y ventanas de compatibilidad. Cuando varios módulos deban cambiar juntos, indique el conjunto compatible en lugar de dejar que el consumidor lo descubra mediante errores de build.

<a name="stop-when-evidence-belongs-to-another-revision"></a>
## Deténgase cuando la evidencia pertenezca a otra revisión

Una verificación fallida, archivo generado sucio, comparación de API ausente, changelog ausente o etiqueta fuera de la rama prevista bloquean el release. No copie resultados de ayer en las notas de hoy; los resultados demuestran un commit y pierden validez cuando cambia.

<a name="do-not-publish-secrets-or-weaponize-a-fix"></a>
## No publique secretos ni convierta una corrección en instrucciones de ataque

Las notas de release pueden describir el comportamiento afectado y la actualización segura sin credenciales, detalles privados de incidentes ni instrucciones de explotación antes de la divulgación coordinada. Los cambios de dependencias exigen repetir controles de seguridad y licencia en el grafo del release, no suponer resultados del workspace.

<a name="compute-the-api-difference-before-tagging"></a>
## Calcule la diferencia de API antes de crear la etiqueta

Los módulos públicos comparan la superficie exportada con la etiqueta anterior y exigen notas de actualización para cambios relevantes de firma. Formato, vet, race, vulnerabilidades, artefactos generados y ascendencia de etiqueta son verificaciones cuyos resultados actuales deben medirse. Enumerarlos aquí no demuestra aprobación.

<a name="correct-with-a-new-version"></a>
## Corrija con una nueva versión

Si un artefacto o nota publicados son incorrectos, preserve la etiqueta original, revierta o corrija en la rama de release y publique la siguiente versión con corrección explícita. Así los consumidores fijan, comparan y revierten sin que una etiqueta resuelva a códigos distintos en dos máquinas.

<a name="upgrade-guides-carry-the-longer-path"></a>
## Las guías de actualización describen transiciones más largas

Las notas de release resumen una versión. Las transiciones de datos o API en varios pasos pertenecen a la guía de actualización, enlazada desde la nota. Changelogs, política de versionado semántico, diferencias de API generada y evidencias de rollback son entradas independientes para la misma decisión; ninguna demuestra por sí sola que el release esté listo.
