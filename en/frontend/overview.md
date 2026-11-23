# Frontend overview

- [HTML is the application boundary](#html-is-the-application-boundary)
- [Edit, compile, and serve one screen](#edit-compile-and-serve-one-screen)
- [Give each frontend layer one job](#give-each-frontend-layer-one-job)
- [A fragment is a separate rendering contract](#a-fragment-is-a-separate-rendering-contract)
- [Test the no-JavaScript and enhanced paths](#test-the-no-javascript-and-enhanced-paths)
- [Released frontend surfaces](#released-frontend-surfaces)
- [Build the two response shapes](#build-the-two-response-shapes)

<a name="html-is-the-application-boundary"></a>
## HTML is the application boundary

Arandu renders HTML in Go, uses YPY to request server-rendered fragments, and embeds its scripts and stylesheet in the application binary. There is no client router, SPA data layer, Node runtime, CDN dependency, or second JSON API required to render a page.

<a name="edit-compile-and-serve-one-screen"></a>
## Edit, compile, and serve one screen

Edit a `.kyse.go` source under resources/views, compile the typed Go view and application stylesheet, then let the server return either a full page with `ctx.View` or a deliberately named fragment with `ctx.Fragment`.

```sh
aru view:build
aru dev
```

<a name="give-each-frontend-layer-one-job"></a>
## Give each frontend layer one job

Ordinary links and forms are the baseline. YPY enhances that path while the server continues to own status, CSRF, authentication, authorization, validation, and persistence.

- Controllers assemble typed page data and choose the response shape.
- Kyse views own markup but never service lookup or authorization.
- Typed components own reusable controls, not application records.
- YPY requests and swaps server fragments; it does not become domain state.

<a name="a-fragment-is-a-separate-rendering-contract"></a>
## A fragment is a separate rendering contract

`Context.Fragment` does not remove a layout from a full-page view. Use a separate layout-free view for a swap target and keep rejected 422 responses truthful: YPY renders a 422 inside the form that sent the request, so the status never has to be changed to success.

> [!WARNING]
> Every browser mutation still needs the session-bound CSRF token and server-side authorization. Do not load scripts from another origin while the released `SecurityHeaders` policy permits scripts only from self.

<a name="test-the-no-javascript-and-enhanced-paths"></a>
## Test the no-JavaScript and enhanced paths

Request the full page without enhanced headers, submit its ordinary form, then repeat with YPY-Request. Assert the document versus fragment roots, truthful statuses, CSRF refusal, visible labels, keyboard focus, and a useful result when scripts do not run.

- `aru view:build` succeeds from a clean checkout.
- Full pages contain one document root.
- Fragments contain no layout root.
- Ordinary and enhanced forms reach the same policy-backed service.

<a name="released-frontend-surfaces"></a>
## Released frontend surfaces

`Context.View` and `Context.Fragment` choose HTML responses. `framework/view.NewModule` supplies the renderer and content-addressed asset routes; `aru view:build` and `aru dev` compile and rebuild sources; Kyse components provide typed server-rendered controls.

| Surface | Released contract |
| --- | --- |
| `Context.View` / Fragment | full and partial HTML responses |
| `framework/view.NewModule` | renderer and embedded assets |
| `aru view:build` / `aru dev` | compile and development loop |
| `kyse/components` | typed server-rendered controls |

<a name="build-the-two-response-shapes"></a>
## Build the two response shapes

Views define the typed render contract, Kyse defines the source language, assets explain the embedded output, and HTTP responses preserve protocol status for both clients.

- [Views](/docs/{{version}}/frontend/views)
- [Kyse templates](/docs/{{version}}/frontend/kyse-templates)
- [Assets](/docs/{{version}}/frontend/assets)
- [HTTP responses](/docs/{{version}}/http/responses)
