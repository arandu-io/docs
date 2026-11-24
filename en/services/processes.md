# Processes

- [Execute an argument vector, never a shell sentence](#execute-an-argument-vector-never-a-shell-sentence)
- [Run one bounded child process](#run-one-bounded-child-process)
- [Control environment, directory, and output explicitly](#control-environment-directory-and-output-explicitly)
- [Separate start, exit, timeout, and output failures](#separate-start-exit-timeout-and-output-failures)
- [Use fixture binaries instead of host commands](#use-fixture-binaries-instead-of-host-commands)
- [Released process component boundary](#released-process-component-boundary)
- [Prefer Go APIs and reserve processes for real tools](#prefer-go-apis-and-reserve-processes-for-real-tools)

<a name="execute-an-argument-vector-never-a-shell-sentence"></a>
## Execute an argument vector, never a shell sentence

The released process component wraps `os/exec` with context, argument vectors, working directory, inherited-plus-overridden environment, captured output, and exit information. Stdout and stderr are retained without a byte ceiling.

<a name="run-one-bounded-child-process"></a>
## Run one bounded child process

Construct an argument vector, set working directory, total and idle deadlines, and only the environment overrides the child needs. Env inherits `os.Environ` and replaces matching keys; it is not an allowlist. Never concatenate request input into a shell command.

```go example
func RunTool(ctx context.Context, workDir, binary, inputPath string, environmentOverrides map[string]string) error {
	result, err := process.NewFactory().NewPendingProcess().
		Path(workDir).
		Env(environmentOverrides).
		Timeout(2 * time.Minute).
		IdleTimeout(30 * time.Second).
		Run(ctx, []string{binary, "--input", inputPath}, nil)
	if err != nil {
		return err
	}
	if _, err := result.Throw(nil); err != nil {
		return err
	}
	return nil
}
```

<a name="control-environment-directory-and-output-explicitly"></a>
## Control environment, directory, and output explicitly

Resolve executable paths at boot, pass each argument as one string, choose a validated working directory, audit the inherited environment, and use Quietly or an `OutputHandler` when retaining unbounded output is unacceptable. Record exit code and duration without output secrets.

- Pass argv without sh `-c`.
- Set total and idle deadlines.
- Audit inherited environment variables.
- Avoid unbounded retained output for chatty tools.

<a name="separate-start-exit-timeout-and-output-failures"></a>
## Separate start, exit, timeout, and output failures

Executable not found, permission, start failure, non-zero exit, signal, cancellation, total timeout, and idle timeout are distinct. A child that starts may have partial external side effects, and captured output can grow without a library ceiling.

> [!WARNING]
> Context cancellation terminates the child boundary but cannot undo its filesystem or remote effects. Do not treat a timeout as proof that nothing happened.

<a name="use-fixture-binaries-instead-of-host-commands"></a>
## Use fixture binaries instead of host commands

Build tiny test helper binaries for success, stderr, chosen exit, sleep, signal, oversized output, environment echo, and child cleanup. Avoid assertions against platform-specific shell utilities or user PATH.

- Arguments with spaces remain one argument.
- Timeout terminates the helper.
- Environment overrides win while parent values remain inherited.
- Chatty output follows the application's memory policy.

<a name="released-process-component-boundary"></a>
## Released process component boundary

`NewFactory`, `NewPendingProcess`, Path, Env, Timeout, `IdleTimeout`, Run, Quietly, and `ProcessResult` form the released boundary. Application authorization, environment minimization, output-memory policy, and tool availability remain explicit.

| Surface | Released contract |
| --- | --- |
| `NewPendingProcess` | argument-vector process builder |
| Path / Env | directory and inherited environment overrides |
| Timeout / `IdleTimeout` / Run | context-bound execution |
| `ProcessResult` | fully captured output and exit metadata |

<a name="prefer-go-apis-and-reserve-processes-for-real-tools"></a>
## Prefer Go APIs and reserve processes for real tools

Context owns cancellation, filesystem owns files, deployment owns installed tools, and contracts keep child processes behind an application-specific adapter.

- [Context](/docs/{{version}}/architecture/context)
- [Filesystem](/docs/{{version}}/services/filesystem)
- [Deployment](/docs/{{version}}/operations/deployment)
- [Contracts](/docs/{{version}}/architecture/contracts-adapters)
