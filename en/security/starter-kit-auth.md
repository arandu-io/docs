# Starter-kit authentication

- [The publisher writes an application-owned identity flow](#the-publisher-writes-an-application-owned-identity-flow)
- [Pin the UI release and protect the working tree](#pin-the-ui-release-and-protect-the-working-tree)
- [Preview, publish, wire, and compile](#preview-publish-wire-and-compile)
- [Preserve custom blocks during a reviewed republish](#preserve-custom-blocks-during-a-reviewed-republish)
- [Refuse a half-published authentication surface](#refuse-a-half-published-authentication-surface)
- [Generated screens do not become authorization](#generated-screens-do-not-become-authorization)
- [Test every published identity transition](#test-every-published-identity-transition)
- [Own the source and its future upgrades](#own-the-source-and-its-future-upgrades)

<a name="the-publisher-writes-an-application-owned-identity-flow"></a>
## The publisher writes an application-owned identity flow

The UI :uiVersion auth publisher writes views, controllers, page helpers, and mailables into an existing project. Those files become application source, not a runtime package that can be upgraded without reviewing bootstrap, policy, session, mail, and tenant decisions.

<a name="pin-the-ui-release-and-protect-the-working-tree"></a>
## Pin the UI release and protect the working tree

Start from a committed or backed-up project, pin the exact UI tag, and confirm `arandu.toml` declares the supported Aru floor. Use dry-run first and inspect the printed bootstrap imports and registration because the publisher does not edit `bootstrap/app.go`.

- Use `github.com/arandu-io/ui@:uiVersion`.
- Start with a clean diff.
- Run `--dry-run` before any write.
- Review the printed composition changes.

<a name="preview-publish-wire-and-compile"></a>
## Preview, publish, wire, and compile

Preview every path, publish the exact tag, add the printed `authui.New` call, with `drivers.Limiter`, to `WithProviders` before app, remove any former authentication module, build Kyse, and migrate. The publisher adds no module dependency and leaves all written files owned by the application.

```sh
go run github.com/arandu-io/ui@:uiVersion auth --dry-run
go run github.com/arandu-io/ui@:uiVersion auth
aru view:build
aru migrate
```

<a name="preserve-custom-blocks-during-a-reviewed-republish"></a>
## Preserve custom blocks during a reviewed republish

Use `--views` only to republish the screens: it leaves the published controllers alone, so a security fix in a controller arrives only through `--force` after diff review. Keep application changes inside `arandu:begin` custom markers. Five tightly coupled files are overwritten even without `--force` while preserving their marked blocks.

- Pin exact publisher versions in automation.
- Review every replaced path.
- Keep custom logic inside markers.
- Adjust registration, password, tenant, and account policy in owned code.

<a name="refuse-a-half-published-authentication-surface"></a>
## Refuse a half-published authentication surface

The publisher refuses a non-project root, invalid module path, unsupported positional arguments, and an old Aru floor before writing. Existing non-coupled files remain unless forced, but a successful write still needs the printed wiring, view generation, migration, and tests.

> [!WARNING]
> Custom blocks are position-matched and templates expect one per file. Code outside those markers can be replaced; dry-run and version-control diff are mandatory before a republish.

<a name="generated-screens-do-not-become-authorization"></a>
## Generated screens do not become authorization

The kit covers login, registration, verification, reset, confirmation, and two-factor screens, but the project owns whether registration is open, policy, tenant resolution, password rules, mail, session backend, and CSRF order. Preserve final-session creation only after required factors.

Published routes and forms are not permission. Keep Rotate, one-time purposes, recent-password middleware, and server-side Grant checks intact when editing.

<a name="test-every-published-identity-transition"></a>
## Test every published identity transition

Publish from a clean tree, inspect wiring, build views, format, vet, race-test and build; inventory auth routes; exercise ordinary and YPY login rejection, registration, verification, reset, two-factor recovery, logout, selected mail transport, custom-block preservation, and `--views` scope.

- Dry-run writes nothing.
- Published routes compile and render.
- Every identity journey reaches its released security boundary.
- Republish preserves custom blocks.

<a name="own-the-source-and-its-future-upgrades"></a>
## Own the source and its future upgrades

Authentication defines the identity contract, two-factor completes the flow, Kyse compiles the screens, and upgrades explain application-owned publisher diffs.

- [Authentication](/docs/{{version}}/security/authentication)
- [HTTP session](/docs/{{version}}/http/session)
- [Kyse templates](/docs/{{version}}/frontend/kyse-templates)
- [Upgrades](/docs/{{version}}/operations/upgrades)
