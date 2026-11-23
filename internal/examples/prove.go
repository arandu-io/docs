package examples

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Versions pins what the examples are proved against.
type Versions struct {
	// Go is the minimum Go version, written as in a go.mod file.
	Go string `json:"go"`
	// Modules are the published tags the examples require, by module path.
	Modules map[string]string `json:"modules"`
}

// LoadVersions reads a versions file.
func LoadVersions(name string) (Versions, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return Versions{}, err
	}
	var v Versions
	if err := json.Unmarshal(data, &v); err != nil {
		return Versions{}, fmt.Errorf("%s: %w", name, err)
	}
	if v.Go == "" || len(v.Modules) == 0 {
		return Versions{}, fmt.Errorf("%s: needs a go version and at least one module", name)
	}
	return v, nil
}

// Options says how to prove a set of examples.
type Options struct {
	Versions Versions
	// Latest resolves every module to its newest tag instead of the pinned one.
	Latest bool
	// Dir is the working directory the module is written to. It is a new
	// temporary directory when empty, and removed afterwards unless Keep is set.
	Dir  string
	Keep bool
	// Stdout and Stderr receive the output of the go commands.
	Stdout, Stderr io.Writer
}

// Result counts what was proved.
type Result struct {
	Proved       int
	Tests        int
	Illustrative int
}

// Prove compiles every proved block, and runs the tests among them. It returns
// an error naming the Markdown file and line of whatever does not hold.
func Prove(ctx context.Context, blocks []Block, opts Options) (Result, error) {
	var result Result
	var proved []Block
	for _, b := range blocks {
		if b.Proved {
			proved = append(proved, b)
		} else {
			result.Illustrative++
		}
	}
	files := make([]File, len(proved))
	var problems []string
	for i, b := range proved {
		file, err := Wrap(b)
		if err != nil {
			problems = append(problems, err.Error())
			continue
		}
		files[i] = file
		result.Proved++
		if file.Test {
			result.Tests++
		}
	}
	if len(problems) > 0 {
		return result, fmt.Errorf("%s", strings.Join(problems, "\n"))
	}
	if len(proved) == 0 {
		return result, nil
	}

	dir := opts.Dir
	if dir == "" {
		temp, err := os.MkdirTemp("", "docs-examples-")
		if err != nil {
			return result, err
		}
		dir = temp
		if !opts.Keep {
			defer os.RemoveAll(dir)
		}
	}
	if err := writeModule(dir, proved, files, opts.Versions); err != nil {
		return result, err
	}

	run := func(args ...string) error {
		cmd := exec.CommandContext(ctx, "go", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=-mod=mod")
		var out bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &out
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(orDiscard(opts.Stderr), "$ go %s\n%s", strings.Join(args, " "), packageDirectory.ReplaceAllString(out.String(), ""))
			return fmt.Errorf("go %s failed: %w", strings.Join(args, " "), err)
		}
		return nil
	}

	if opts.Latest {
		for _, module := range sortedKeys(opts.Versions.Modules) {
			if err := run("get", module+"@latest"); err != nil {
				return result, err
			}
		}
	}
	if err := run("mod", "tidy"); err != nil {
		return result, err
	}
	if err := run("vet", "./..."); err != nil {
		return result, err
	}
	if err := run("build", "./..."); err != nil {
		return result, err
	}
	if result.Tests > 0 {
		if err := run("test", "-count=1", "./..."); err != nil {
			return result, err
		}
	}
	return result, nil
}

func writeModule(dir string, blocks []Block, files []File, v Versions) error {
	var mod bytes.Buffer
	fmt.Fprintf(&mod, "module example.test\n\ngo %s\n\nrequire (\n", v.Go)
	for _, module := range sortedKeys(v.Modules) {
		fmt.Fprintf(&mod, "\t%s %s\n", module, v.Modules[module])
	}
	mod.WriteString(")\n")
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), mod.Bytes(), 0o644); err != nil {
		return err
	}
	for i, file := range files {
		pkg := filepath.Join(dir, fmt.Sprintf("ex%03d", i+1))
		if err := os.MkdirAll(pkg, 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(pkg, file.Name), []byte(file.Source), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// packageDirectory is the directory of a generated package, which a line
// directive turns into a prefix of the Markdown path in compiler messages.
var packageDirectory = regexp.MustCompile(`\bex\d{3}/`)

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func orDiscard(w io.Writer) io.Writer {
	if w == nil {
		return io.Discard
	}
	return w
}
