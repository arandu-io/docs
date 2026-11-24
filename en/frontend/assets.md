# Assets without Node

- [The binary serves content-addressed assets](#the-binary-serves-content-addressed-assets)
- [Pin the CLI floor and downloaded tool](#pin-the-cli-floor-and-downloaded-tool)
- [Compile one stylesheet into the application](#compile-one-stylesheet-into-the-application)
- [Keep source, output, and serving visible](#keep-source-output-and-serving-visible)
- [Recognize stale and unconsumed output](#recognize-stale-and-unconsumed-output)
- [Keep scripts local and bytes attributable](#keep-scripts-local-and-bytes-attributable)
- [Prove bytes, hashes, headers, and classes](#prove-bytes-hashes-headers-and-classes)
- [Rebuild from sources and verify the embedded result](#rebuild-from-sources-and-verify-the-embedded-result)
- [Trace an asset from source to response](#trace-an-asset-from-source-to-response)

<a name="the-binary-serves-content-addressed-assets"></a>
## The binary serves content-addressed assets

`aru view:build` compiles project CSS with the pinned standalone Tailwind executable and writes `assets/app.css`. The application embeds those bytes and the framework view module serves CSS, YPY, and interaction scripts under a SHA-256-derived path that changes with content.

<a name="pin-the-cli-floor-and-downloaded-tool"></a>
## Pin the CLI floor and downloaded tool

Set the released Aru floor and Tailwind version in `arandu.toml`. A Tailwind version Aru ships no checksums for also needs tailwindcss-sha256sums, the SHA-256 of that release's checksum file, under `[tools]`; without it the download is refused before it starts. The Kyse compiler is inside Aru; only the standalone stylesheet compiler is fetched into the verified tool cache, and the running server needs neither Node nor the build tool.

- Pin `[arandu].aru` to a compatible floor.
- Pin `[tools].tailwindcss` to the reviewed binary version.
- Keep `resources/css/app.css` and Kyse sources under review.
- Import the application assets package in bootstrap.

<a name="compile-one-stylesheet-into-the-application"></a>
## Compile one stylesheet into the application

Run the normal view build before Go compilation. It generates application views and writes the minified embedded stylesheet; importing the assets package registers those bytes and `framework/view.NewModule` mounts the hashed route.

```sh
aru view:build
go build ./...
```

<a name="keep-source-output-and-serving-visible"></a>
## Keep source, output, and serving visible

Source CSS and Kyse remain application-owned. Generated views are normally gitignored, while the official consumer commits the compiled `assets/app.css`. framework scripts travel with the released module and `view.URL` emits a local path for the bytes actually registered.

- Edit `resources/css/app.css` and Kyse sources.
- Regenerate `storage/framework/views` instead of editing it.
- Commit the official consumer's `assets/app.css` output.
- Keep normal `aru build` view generation enabled.

<a name="recognize-stale-and-unconsumed-output"></a>
## Recognize stale and unconsumed output

Missing asset imports silently leave the framework default stylesheet, absent generated packages break clean builds, stale deleted sources can duplicate names, and an unknown asset name panics during rendering rather than inventing a 404 URL.

> [!WARNING]
> `aru build --skip-views` can package stale generated views and CSS. Use it only in a pipeline that has already proved current outputs from the same source commit.

<a name="keep-scripts-local-and-bytes-attributable"></a>
## Keep scripts local and bytes attributable

> [!WARNING]
> Keep scripts on the same origin, preserve license notices, use `view.URL` instead of hand-written hashes, register accurate content types, and never serve user uploads from the trusted application asset registry or execute an unpinned checkout-supplied tool.

<a name="prove-bytes-hashes-headers-and-classes"></a>
## Prove bytes, hashes, headers, and classes

Build twice and compare bytes, request the content-addressed stylesheet, assert Content-Type, nosniff, immutable caching, representative application classes, and a changed URL after one byte changes.

- Unchanged sources produce identical CSS.
- The served body matches embedded bytes.
- A byte change changes the URL.
- Application-only classes survive compilation.

<a name="rebuild-from-sources-and-verify-the-embedded-result"></a>
## Rebuild from sources and verify the embedded result

Confirm the Aru floor and Tailwind pin, regenerate sources, rebuild and test the asset response, then inspect the running page's hashed URL. Clearing a browser cache cannot repair a binary built with the wrong bytes.

If download is unavailable, restore the verified cache or build in a networked preparation stage; the production server has no asset publishing step.

<a name="trace-an-asset-from-source-to-response"></a>
## Trace an asset from source to response

The overview explains the server-rendered contract, views own typed HTML, Kyse owns source compilation, and middleware constrains where scripts can load from.

- [Frontend overview](/docs/{{version}}/frontend/overview)
- [Views](/docs/{{version}}/frontend/views)
- [Kyse templates](/docs/{{version}}/frontend/kyse-templates)
- [Middleware](/docs/{{version}}/http/middleware)
