# Autenticación del starter kit

- [El publicador escribe un flujo de identidad perteneciente a la aplicación](#the-publisher-writes-an-application-owned-identity-flow)
- [Fije la versión UI y proteja el árbol de trabajo](#pin-the-ui-release-and-protect-the-working-tree)
- [Previsualice, publique, conecte y compile](#preview-publish-wire-and-compile)
- [Preserve bloques personalizados en una republicación revisada](#preserve-custom-blocks-during-a-reviewed-republish)
- [Rechace una superficie de autenticación publicada a medias](#refuse-a-half-published-authentication-surface)
- [Las pantallas generadas no se vuelven autorización](#generated-screens-do-not-become-authorization)
- [Pruebe cada transición de identidad publicada](#test-every-published-identity-transition)
- [Asuma el código fuente y sus futuras actualizaciones](#own-the-source-and-its-future-upgrades)

<a name="the-publisher-writes-an-application-owned-identity-flow"></a>
## El publicador escribe un flujo de identidad perteneciente a la aplicación

El publicador auth de UI :uiVersion escribe vistas, controllers, helpers y mailables en un proyecto existente. Se vuelven código de aplicación, no paquete de runtime actualizable sin revisar bootstrap, policy, sesión, correo y decisiones de tenant.

<a name="pin-the-ui-release-and-protect-the-working-tree"></a>
## Fije la versión UI y proteja el árbol de trabajo

Comience con proyecto commiteado o con backup, fije la etiqueta exacta de UI y confirme la versión mínima de Aru en `arandu.toml`. Haga dry-run e inspeccione imports y registro de bootstrap impresos: el publicador no edita `bootstrap/app.go`.

- Use `github.com/arandu-io/ui@:uiVersion`.
- Comience con un diff limpio.
- Ejecute `--dry-run` antes de escribir.
- Revise los cambios de composición impresos.

<a name="preview-publish-wire-and-compile"></a>
## Previsualice, publique, conecte y compile

Previsualice cada ruta, publique la tag exacta, añada a `WithProviders`, antes de app, la llamada impresa a `authui.New` con `drivers.Limiter`, quite cualquier módulo de autenticación anterior, construya el Kyse y migre. El publicador no añade dependencia de módulo y deja todos los archivos escritos como propiedad de la aplicación.

```sh
go run github.com/arandu-io/ui@:uiVersion auth --dry-run
go run github.com/arandu-io/ui@:uiVersion auth
aru view:build
aru migrate
```

<a name="preserve-custom-blocks-during-a-reviewed-republish"></a>
## Preserve bloques personalizados en una republicación revisada

Use `--views` solo para republicar las pantallas: deja intactos los controllers publicados, así que una corrección de seguridad en un controller solo llega mediante `--force` tras revisar el diff. Mantenga cambios de aplicación en marcadores `arandu:begin` custom. Cinco archivos fuertemente acoplados se sobrescriben incluso sin `--force`, preservando bloques marcados.

- Fije versiones exactas del publicador en automatización.
- Revise cada ruta sustituida.
- Mantenga la lógica personalizada dentro de los marcadores.
- Ajuste políticas de registro, contraseña, tenant y cuenta en código propio.

<a name="refuse-a-half-published-authentication-surface"></a>
## Rechace una superficie de autenticación publicada a medias

El publicador rechaza raíz sin proyecto, ruta de módulo inválida, argumentos posicionales no compatibles y versión mínima antigua de Aru antes de escribir. Los archivos existentes no acoplados permanecen salvo force; una escritura correcta aún exige composición impresa, generación de vistas, migración y pruebas.

> [!WARNING]
> Los bloques personalizados se corresponden por posición y los templates esperan actualmente uno por archivo. El código fuera de los marcadores puede sustituirse; dry-run y diff versionado son obligatorios antes de republicar.

<a name="generated-screens-do-not-become-authorization"></a>
## Las pantallas generadas no se vuelven autorización

> [!WARNING]
> El kit incluye pantallas de login, registro, verificación, recuperación, confirmación y dos factores. El proyecto decide apertura de registro, policy, resolución de tenant, reglas de contraseña, correo, backend de sesión y orden CSRF. Cree la sesión final solo tras los factores exigidos.

Las rutas y formularios publicados no son permiso. Preserve Rotate, propósitos de un solo uso, middleware de contraseña reciente y verificaciones de Grant en servidor al editar.

<a name="test-every-published-identity-transition"></a>
## Pruebe cada transición de identidad publicada

Publique desde un árbol limpio, inspeccione composición, genere vistas, formatee, ejecute vet, race y build; inventaríe rutas de auth; ejercite rechazo de login normal y YPY, registro, verificación, recuperación, recuperación de dos factores, logout, transporte de correo, preservación de bloques personalizados y ámbito `--views`.

- Dry-run no escribe nada.
- Las rutas publicadas compilan y renderizan.
- Cada flujo de identidad alcanza su frontera de seguridad publicada.
- Republicar preserva bloques personalizados.

<a name="own-the-source-and-its-future-upgrades"></a>
## Asuma el código fuente y sus futuras actualizaciones

La autenticación define identidad, los dos factores completan el flujo, Kyse compila pantallas y las actualizaciones explican los diffs del publicador pertenecientes a la aplicación.

- [Autenticación](/docs/{{version}}/security/authentication)
- [Sesión HTTP](/docs/{{version}}/http/session)
- [Templates Kyse](/docs/{{version}}/frontend/kyse-templates)
- [Actualizaciones](/docs/{{version}}/operations/upgrades)
