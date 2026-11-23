# Console tests

- [Console tests execute the project-owned command boundary](#console-tests-execute-the-project-owned-command-boundary)
- [Dispatch one command and inspect its effect](#dispatch-one-command-and-inspect-its-effect)
- [Assert arguments, confirmation, effects, and exit](#assert-arguments-confirmation-effects-and-exit)
- [A registered name is not proof that its operation runs](#a-registered-name-is-not-proof-that-its-operation-runs)
- [Exercise the delegated command through Aru and the binary](#exercise-the-delegated-command-through-aru-and-the-binary)
- [Released console testing surfaces](#released-console-testing-surfaces)
- [Pair commands with their operational recovery guide](#pair-commands-with-their-operational-recovery-guide)

> [!WARNING]
> **Current limitations.** `Application.HandleCommand` with `WithConsole` is the released command harness: it takes the test's context and buffers and answers the exit code. Signals and the real binary's exit still need a subprocess.

<a name="console-tests-execute-the-project-owned-command-boundary"></a>
## Console tests execute the project-owned command boundary

Application-aware commands are implemented by the project binary and invoked through Aru. The released application seam is `Application.HandleCommand(ctx, args)`, which answers the exit code, and `WithConsole` gives it buffers in place of the process streams. A useful test asserts the output, the exit, and the durable effect; a command-name inventory proves only discoverability.

<a name="dispatch-one-command-and-inspect-its-effect"></a>
## Dispatch one command and inspect its effect

Set isolated application configuration, build the application with the same bootstrap the binary uses, hand it buffers with `WithConsole`, call `HandleCommand` with the exact command and arguments, and inspect the exit code beside the database, queue, scheduler, or filesystem outcome.

```go example
func runSchedule(ctx context.Context, t *testing.T, app *foundation.Application) {
	var out, errOut bytes.Buffer
	code := app.WithConsole(&out, &errOut, strings.NewReader("")).
		HandleCommand(ctx, []string{"schedule:run", "billing.close", "--tenant=acme"})
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut.String())
	}
}
```

<a name="assert-arguments-confirmation-effects-and-exit"></a>
## Assert arguments, confirmation, effects, and exit

Test invalid arguments, defaults versus flags, repeated invocation, and the selected adapter's actual effect. `HandleCommand` takes the context the test owns and the streams `WithConsole` set, boots the application only for a command that needs it, and closes what it opened before it returns. Use a subprocess only when the process boundary itself, a signal or the real binary's exit, is the contract.

- Assert exact delegated command and flags.
- Inspect durable side effects.
- Read output and errors from the `WithConsole` buffers.
- Run destructive commands only on disposable data.

<a name="a-registered-name-is-not-proof-that-its-operation-runs"></a>
## A registered name is not proof that its operation runs

A command can parse and still ignore a configuration value or call the wrong constructor default. Unknown names, missing project markers, refused production operations, held locks, unavailable dependencies, and partial writes all need their real exit and evidence.

> [!WARNING]
> Do not substitute a test that checks fourteen queue names for a journey that starts, pauses, retries, or drains a worker. Operational behavior needs the chosen store and process boundary.

<a name="exercise-the-delegated-command-through-aru-and-the-binary"></a>
## Exercise the delegated command through Aru and the binary

For each command family, run one successful effect, malformed invocation, dependency failure, cancellation, and repeated call. For delegated Aru commands, execute the pinned CLI against a temporary project binary and assert forwarded arguments and exit code.

- Defaults and overrides reach runtime options.
- Refused operations return an error.
- Repeated invocation preserves its idempotency contract.
- The selected persistent effect is observable.

<a name="released-console-testing-surfaces"></a>
## Released console testing surfaces

Aru owns top-level discovery and forwarding; `Application.HandleCommand(ctx, args)` int is the released application seam for the project-owned command, and `WithConsole(out, errOut, in)` sets the streams it writes to and reads from. help and list print the command list, and an unknown name prints it and answers 1.

| Surface | Released contract |
| --- | --- |
| `Application.HandleCommand` | context, argument slice, and exit code |
| `Application.WithConsole` | stdout, stderr, and stdin for one test |
| durable state | primary observable effect |
| subprocess | signals and the real binary's exit |

<a name="pair-commands-with-their-operational-recovery-guide"></a>
## Pair commands with their operational recovery guide

CLI reference lists commands, queue operations define worker effects, migrations define schema effects, and deployment defines production invocation.

- [CLI reference](/docs/{{version}}/cli/reference)
- [Queue operations](/docs/{{version}}/operations/queues)
- [Migrations](/docs/{{version}}/database/migrations)
- [Deployment](/docs/{{version}}/operations/deployment)
