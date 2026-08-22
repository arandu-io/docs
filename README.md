<p align="center">
  <img src=".github/logo.png" alt="Arandu" width="140" height="140">
</p>

<h1 align="center">arandu-io/docs</h1>

<p align="center">The decisions behind Arandu, and why each door was closed.</p>

<p align="center">
<a href="https://github.com/arandu-io/docs/tags"><img src="https://img.shields.io/github/v/tag/arandu-io/docs?label=version" alt="Latest Version"></a>
<a href="LICENSE.md"><img src="https://img.shields.io/github/license/arandu-io/docs" alt="License"></a>
</p>


## About this record

Not a usage guide — that is the next phase, and it will be an Arandu
application. This is the record: thirty-two numbered documents and the fifty-two
architecture decision records behind them.

An ADR captures the *why*, which is the part that cannot be reconstructed from
the code afterwards. Length follows the argument rather than a template:
fifteen lines where the decision was easy, four hundred and sixty-seven where
it was not, seventy-eight in the middle. One is written whenever a decision
closes a door: no dependency injection container, no runtime reflection, no
Active Record, no Node in the build, no second way to reach the database.

Written in Brazilian Portuguese, unlike everything else in the project. The code
is in English because the product is global and pkg.go.dev publishes the doc
comments; the decisions are in the language they were argued in.

## Learning Arandu

The API reference is generated from the doc comments and lives on
[pkg.go.dev](https://pkg.go.dev/github.com/arandu-io/framework). Every exported
symbol carries one, and that is deliberate: it is the documentation that cannot
drift from the code, because it sits in the same file.

The CLI documents itself. `aru help` lists every command, and each one explains
what it writes and what to do with it. `aru doctor` explains what it found and
what breaks, not which rule was violated.

A guide and a website do not exist yet, and that is a decision rather than a
gap: a guide written against an API that still moves is work done twice, and the
second time is worse — there is wrong documentation published. The site is the
next phase, and it will be an Arandu application.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Before opening a pull request, the check
at the top of that file has to pass. CI runs it, and one check the script does
not make on its own: that `_archive/` was not edited.

## Security Vulnerabilities

Please review [our security policy](SECURITY.md) on how to report a
vulnerability. Never open a public issue for one.

## License

Open-sourced software licensed under the [MIT license](LICENSE.md).
