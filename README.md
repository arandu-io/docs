<p align="center">
  <img src=".github/logo.svg" alt="Arandu" width="180">
</p>

<h1 align="center">arandu-io/docs</h1>

<p align="center">The Arandu documentation: how to use it, and the decisions behind it.</p>

<p align="center">
<a href="https://github.com/arandu-io/docs/tags"><img src="https://img.shields.io/github/v/tag/arandu-io/docs?label=version" alt="Latest Version"></a>
<a href="LICENSE.md"><img src="https://img.shields.io/github/license/arandu-io/docs" alt="License"></a>
</p>


## About this repository

This is the usage documentation of Arandu, written in Markdown the way
[laravel/docs](https://github.com/laravel/docs) is, and the archive of the
decisions that shaped the project. [arandu.io/docs](https://arandu.io/docs)
renders the usage documentation; this repository is where it is written.

| Path | What it holds |
|---|---|
| `en/`, `pt/`, `es/` | The usage guide in English, Brazilian Portuguese and Spanish: one file per page, under the same path in each language |
| `en/documentation.md` | The navigation: sections and pages, in reading order. `pt/` and `es/` carry their own |
| `versions.json` | The tags every Go example is proved against |
| `cmd/`, `internal/` | The standard-library Go commands that check the pages |
| `_archive/` | The numbered documents and the architecture decision records, frozen as they stood when they were archived |

The branches follow the versions, as `laravel/docs` does: `main` is the next
version and `0.x` is the published one. When 1.0 ships, `1.x` is created.

## How a page is written

Each page has one `# Title`, a list of its sections, and one
`<a name="anchor"></a>` and `## Heading` per section. Links between pages are
`/docs/{{version}}/<group>/<page>#anchor`, and a release is written as a
placeholder such as `:frameworkVersion`, which the site replaces with the tag it
publishes. [CONTRIBUTING.md](CONTRIBUTING.md) has the full format.

## Examples that compile

A Go block fenced as `go example` is a promise. `prove-examples` writes every
such block of `en/` into a temporary module that requires the tags in
`versions.json`, vets and builds it, and runs it when it declares a test. A block
that does not compile fails CI with the Markdown file and line, and the `pt/` and
`es/` pages must repeat the same blocks as the English page.

```
GOWORK=off go run ./cmd/prove-examples            # the pinned tags
GOWORK=off go run ./cmd/prove-examples -latest    # the newest tag of every module
```

A plain `go` block is illustration: it shows a shape that depends on generated
code or on an application's own services, and it is counted, never compiled.
`examples-latest` runs the second command every day, so a release that breaks a
page fails the day it is published.

## Learning Arandu

Start at the documentation index, `en/documentation.md`. The API reference is
generated from the doc comments and lives on
[pkg.go.dev](https://pkg.go.dev/github.com/arandu-io/framework): every exported
symbol carries one, which is the documentation that cannot drift from the code
because it sits in the same file.

The CLI documents itself. `aru help` lists every command, and each one explains
what it writes and what to do with it. `aru doctor` explains what it found and
what breaks, not which rule was violated.

## The record of decisions

`_archive/` holds the thirty-two numbered documents and the fifty-two
architecture decision records behind them, written in Brazilian Portuguese, the
language they were argued in. An ADR captures the *why*, which cannot be
reconstructed from the code afterwards, and one is written whenever a decision
closes a door: no dependency injection container, no runtime reflection, no
Active Record, no Node in the build, no second way to reach the database. The
archive is frozen; a decision that changed is corrected in the note that
supersedes it, never in the archived copy.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Before opening a pull request, the check
at the top of that file has to pass. CI runs it, and one check the script does
not make on its own: that `_archive/` was not edited.

## Security Vulnerabilities

Please review [our security policy](SECURITY.md) on how to report a
vulnerability. Never open a public issue for one.

## License

Open-sourced software licensed under the [MIT license](LICENSE.md).
