# Release notes

- [A release note describes an immutable version](#a-release-note-describes-an-immutable-version)
- [Start from the previous immutable tag](#start-from-the-previous-immutable-tag)
- [Describe the observable change](#describe-the-observable-change)
- [Group additions, changes, fixes, and removals](#group-additions-changes-fixes-and-removals)
- [Stop when evidence belongs to another revision](#stop-when-evidence-belongs-to-another-revision)
- [Do not publish secrets or weaponize a fix](#do-not-publish-secrets-or-weaponize-a-fix)
- [Compute the API difference before tagging](#compute-the-api-difference-before-tagging)
- [Correct with a new version](#correct-with-a-new-version)
- [Upgrade guides carry the longer path](#upgrade-guides-carry-the-longer-path)

<a name="a-release-note-describes-an-immutable-version"></a>
## A release note describes an immutable version

Each published module is below 1.0 and records incompatible surface changes in its changelog. A tag is immutable: correcting a released artifact means publishing a new version, never moving the old tag. That property lets builds and incident reviews recover the same source later.

<a name="start-from-the-previous-immutable-tag"></a>
## Start from the previous immutable tag

A release candidate needs a clean revision, passing repository gates, a changelog entry, and a known previous tag. Compare from that tag to the exact commit intended for release; comparing a moving branch against another moving branch produces notes nobody can reproduce later.

<a name="describe-the-observable-change"></a>
## Describe the observable change

Write what a consumer must change, which behavior differs, and which module version first carries it. Name removed or renamed symbols precisely and include the replacement path. The commands below are an illustrative inspection sequence, not the repository's release procedure or evidence that any gate passed; the repository-owned release scripts and policy remain authoritative.

```bash
PREVIOUS_TAG=$(git describe --tags --abbrev=0 HEAD^)
git diff --stat "$PREVIOUS_TAG"...HEAD
GOWORK=off go vet ./...
GOWORK=off go test -race -count=1 -p 1 ./...
```

<a name="group-additions-changes-fixes-and-removals"></a>
## Group additions, changes, fixes, and removals

Keep entries in consumer language and link each migration note to the affected module. Call out configuration changes, data migrations, security behavior, and compatibility windows. When several modules must move together, state the compatible version set rather than leaving the consumer to discover it through build errors.

<a name="stop-when-evidence-belongs-to-another-revision"></a>
## Stop when evidence belongs to another revision

A failed gate, dirty generated file, missing API comparison, absent changelog entry, or tag that does not descend from the intended branch blocks the release. Do not copy yesterday's successful command output into today's notes; results are evidence about one commit and expire when the commit changes.

<a name="do-not-publish-secrets-or-weaponize-a-fix"></a>
## Do not publish secrets or weaponize a fix

Release notes may name affected behavior and the safe upgrade without including credentials, private incident details, or exploit-ready steps before coordinated disclosure. Dependency changes need their security and license checks repeated on the release graph, not assumed from a development workspace.

<a name="compute-the-api-difference-before-tagging"></a>
## Compute the API difference before tagging

Public modules compare their exported surface with the previous tag and require upgrade notes for meaningful signature changes. Formatting, vet, race tests, vulnerability checks, generated artifacts, and tag ancestry are gates whose current results must be measured. Merely listing a gate here is not proof that it passed.

<a name="correct-with-a-new-version"></a>
## Correct with a new version

If a published artifact or note is wrong, leave the original tag immutable, revert or correct on the release branch, and publish the next version with an explicit correction. Consumers can then pin, compare, and roll back without a tag resolving to different source on two machines.

<a name="upgrade-guides-carry-the-longer-path"></a>
## Upgrade guides carry the longer path

Release notes summarize one version. A multi-step data or API transition belongs in the upgrade guide, with the note linking to it. Module changelogs, semantic version policy, generated API differences, and deployment rollback evidence are separate inputs to the same decision; none alone proves a release is ready.
