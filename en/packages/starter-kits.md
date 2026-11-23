# Starter kits

- [Start from owned source and keep owning it](#start-from-owned-source-and-keep-owning-it)
- [Create a project and pin the publisher](#create-a-project-and-pin-the-publisher)
- [Preview and publish the authentication kit](#preview-and-publish-the-authentication-kit)
- [Republish without silently losing application code](#republish-without-silently-losing-application-code)
- [Refuse a half-published or unbuildable flow](#refuse-a-half-published-or-unbuildable-flow)
- [Recheck the published identity flow as application code](#recheck-the-published-identity-flow-as-application-code)
- [Test the generated journey, not only the file list](#test-the-generated-journey-not-only-the-file-list)
- [Continue with authentication, Kyse, and upgrades](#continue-with-authentication-kyse-and-upgrades)

<a name="start-from-owned-source-and-keep-owning-it"></a>
## Start from owned source and keep owning it

Arandu has two source-delivery paths. `aru new` creates the base application from the skeleton; the UI module publishes an authentication flow into an existing project. Neither is a runtime abstraction: generated and published files become application source to review, edit, test, and upgrade deliberately.

<a name="create-a-project-and-pin-the-publisher"></a>
## Create a project and pin the publisher

Install Aru, Git, and Go. `aru new` creates a project, removes skeleton history, writes a private environment with fresh secrets, and attempts view compilation. Commit or back up the project before the pinned UI publisher and ensure `arandu.toml` declares its supported Aru floor.

- Record or verify the starter commit.
- Use an exact UI publisher tag.
- Start from a clean worktree.
- Run publisher dry-run first.

<a name="preview-and-publish-the-authentication-kit"></a>
## Preview and publish the authentication kit

Preview then publish UI :uiVersion. It writes application-owned Kyse and Go files and prints the wiring without editing `bootstrap/app.go`: the import of the published controllers and an `authui.New` call to list in `WithProviders`, before app, built from the provider's account rules and the drivers, `drivers.Limiter` included. Apply that wiring, remove any former authentication module from `WithProviders`, build views, and migrate.

```sh
go run github.com/arandu-io/ui@:uiVersion auth --dry-run
go run github.com/arandu-io/ui@:uiVersion auth
aru view:build
aru migrate
```

<a name="republish-without-silently-losing-application-code"></a>
## Republish without silently losing application code

Pin the publisher in automation and preserve application changes inside arandu custom markers. `--views` republishes only the screens and leaves published controllers alone, so a security fix in a controller needs `--force` after diff review; the code-form limits and timed responses of ui :uiVersion live in the controllers and arrive only with auth `--force`. Five coupled layout and page files replace outside-marker content even without force.

- Review all publisher paths.
- Follow the printed `WithProviders` wiring.
- Keep custom code inside markers.
- Compile and route-test after each publish.

<a name="refuse-a-half-published-or-unbuildable-flow"></a>
## Refuse a half-published or unbuildable flow

The publisher refuses a non-project root, invalid module path, unexpected arguments, and an old Aru floor before writing. A sequential filesystem error can still leave a partial set; version control and rerunning the same pinned tag provide recovery.

> [!WARNING]
> `aru new` :aruVersion selects the immutable skeleton :aranduVersion and exposes no override flag. Verify that tag's `go.mod` and the generated project before treating the baseline as reproducible.

<a name="recheck-the-published-identity-flow-as-application-code"></a>
## Recheck the published identity flow as application code

> [!WARNING]
> The kit supplies transport flows but the project owns policy, tenant resolution, password rules, mail, session, CSRF, and registration availability. Preserve session rotation, purpose-bound codes, factor completion, and recent-password gates when editing.

Generated screens and routes do not authorize application data. Every operation still reaches the application service and Policy.

<a name="test-the-generated-journey-not-only-the-file-list"></a>
## Test the generated journey, not only the file list

Dry-run from clean state, publish the exact tag, inspect changes and wiring, build views, vet, race-test and build, inventory auth routes, run ordinary and YPY identity journeys, and republish to prove marker preservation and `--views` scope.

- Dry-run writes nothing.
- Published source compiles.
- Authentication journeys reach released services.
- Custom blocks survive republish.

<a name="continue-with-authentication-kyse-and-upgrades"></a>
## Continue with authentication, Kyse, and upgrades

Authentication defines identity, Kyse compiles the published views, the editor extension understands the project, and upgrades carry application-owned diffs forward.

- [Authentication](/docs/{{version}}/security/authentication)
- [Kyse templates](/docs/{{version}}/frontend/kyse-templates)
- [Editor extension](/docs/{{version}}/ai/assisted-development)
- [Upgrades](/docs/{{version}}/operations/upgrades)
