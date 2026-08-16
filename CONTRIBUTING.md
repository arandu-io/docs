# Contributing

## Sign your commits

Every commit needs a `Signed-off-by` line:

```
git commit -s -m "what changed and why"
```

That line is the [Developer Certificate of Origin](https://developercertificate.org/):
you are stating that you wrote the change, or that you have the right to submit
it under this project's license. It is not a copyright assignment — you keep
your copyright, and this project can never be relicensed behind your back.

We use DCO rather than a CLA on purpose. A CLA would let the project relicense
later, and the price is that every contributor has to sign a legal document
before their first patch.

## Before you open a pull request

```
gofmt -l $(find . -name '*.go' -not -path '*/testdata/*' -not -name '*.kyse.go')   # no output
go vet ./...
go test -race ./...
```

Both filters earn their place, and CI runs this exact `find`:

  - a `.kyse.go` file is a view template that a build tag keeps out of the
    compiler, and `gofmt` is the one tool in the chain that ignores build tags —
    so a plain `gofmt -l .` reports a parse error on every view in the skeleton
    repositories and looks like a broken checkout
  - `testdata/` holds fixtures that are invalid on purpose. `aru` keeps a policy
    file with a missing brace so that `aru doctor` has broken code to find, and
    `gofmt` reports it every time

CI runs exactly this, plus a check that no new dependency entered the core: the
framework depends on the standard library and `golang.org/x/crypto`, and nothing
else. A pull request that adds a dependency there needs to argue for it first,
in an issue.

## Where a test goes

Beside the code it tests, named `*_test.go`, in the same directory. There is no
`tests/` directory, and that is not style: `go test` attributes coverage per
directory, so a test filed elsewhere leaves the package under test reporting
0% -- and it can only reach what the package exports.

Which package the test declares is a real choice, and it answers one question:

| declare | when |
|---|---|
| `package X_test` | this is the **contract**. The test sees what a caller sees, which is the point |
| `package X` | this is the **implementation**, and the test genuinely needs something the package does not export |

Prefer the first. Take the second only when you use it -- `plans/testpackages.go`
in the arandu-io working tree checks exactly that, by intersecting the
identifiers a test names with what its package declares unexported, and the
checklist runs it across every repository.

A `package main` has no external form: it cannot be imported, so its tests are
internal and that is the end of it.

## What the commit message says

What changed and why. The why is the part that is not in the diff, and it is the
part someone will need in two years.

Write the subject as an imperative sentence, with no prefix:

    Add the community health files GitHub reads from fixed paths
    Correct two comments that name a removed command

This project does not use conventional commits. No `feat:`, `fix:` or `chore:`
tag — the subject is a sentence, and it says what the commit does to the tree.
One commit per logical unit: formatting does not travel with a change in
behaviour.

No AI attribution of any kind: no `Co-Authored-By` for an assistant, no
"generated with" footer. Commits are authored by the people who submit them.

## Architecture decisions

The decisions this project has already made live in `arandu-io/docs`, and every
one that closed a door has an ADR. If your change contradicts one, say so in the
pull request and argue for the change of decision — that is a normal thing to
do, and it is better than a patch that quietly works around it.
