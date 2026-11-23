package examples_test

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/arandu-io/docs/internal/examples"
)

func page(blocks ...string) *fstest.MapFS {
	return &fstest.MapFS{"en/http/page.md": {Data: []byte("# Page\n\n" + strings.Join(blocks, "\n\n") + "\n")}}
}

func TestExtractSeparatesProvedFromIllustrativeBlocks(t *testing.T) {
	fsys := page(
		"```go example\nfunc A() {}\n```",
		"```go\nx := 1\n```",
		"```sh\naru dev\n```",
	)

	blocks, err := examples.Extract(fsys, "en")
	if err != nil {
		t.Fatal(err)
	}

	if len(blocks) != 2 {
		t.Fatalf("got %d go blocks, want 2", len(blocks))
	}
	if !blocks[0].Proved || blocks[1].Proved {
		t.Errorf("proved = %v, %v; want true, false", blocks[0].Proved, blocks[1].Proved)
	}
	if got, want := blocks[0].Location(), "en/http/page.md:4"; got != want {
		t.Errorf("location = %q, want %q", got, want)
	}
}

func TestExtractRefusesAnExampleThatIsNotGo(t *testing.T) {
	_, err := examples.Extract(page("```sh example\naru dev\n```"), "en")
	if err == nil || !strings.Contains(err.Error(), "only a go block can be an example") {
		t.Fatalf("err = %v", err)
	}
}

func TestWrapAddsPackageAndImportsForDeclarations(t *testing.T) {
	file, err := examples.Wrap(examples.Block{Page: "en/a.md", Line: 4, Code: "func Name() string { return strings.ToUpper(fmt.Sprint(1)) }"})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"package example", `"fmt"`, `"strings"`, "//line en/a.md:4:1"} {
		if !strings.Contains(file.Source, want) {
			t.Errorf("source lacks %q:\n%s", want, file.Source)
		}
	}
	if file.Test {
		t.Error("a plain function is not a test")
	}
}

func TestWrapRefusesStatementsAndUnknownQualifiers(t *testing.T) {
	if _, err := examples.Wrap(examples.Block{Page: "en/a.md", Line: 1, Code: "x := 1"}); err == nil {
		t.Error("a statement outside a function was accepted")
	}
	_, err := examples.Wrap(examples.Block{Page: "en/a.md", Line: 1, Code: "func F() { missing.Call() }"})
	if err == nil || !strings.Contains(err.Error(), "missing") {
		t.Errorf("err = %v, want it to name the qualifier", err)
	}
}

func TestWrapKeepsAWholeFileAndSpotsTests(t *testing.T) {
	file, err := examples.Wrap(examples.Block{Page: "en/a.md", Line: 1, Code: "package demo\n\nimport \"testing\"\n\nfunc TestX(t *testing.T) {}"})
	if err != nil {
		t.Fatal(err)
	}
	if !file.Test || file.Name != "example_test.go" {
		t.Errorf("file = %s, test = %v", file.Name, file.Test)
	}
}

func goAvailable(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("the go command is not available")
	}
}

func TestProveAcceptsCompilingExamplesAndRunsTheirTests(t *testing.T) {
	goAvailable(t)
	blocks := []examples.Block{
		{Page: "en/a.md", Line: 3, Proved: true, Code: "func Double(n int) int { return n * 2 }"},
		{Page: "en/a.md", Line: 9, Proved: true, Code: "func TestDouble(t *testing.T) {\n\tif Double(2) != 4 {\n\t\tt.Fatal(\"wrong\")\n\t}\n}\n\nfunc Double(n int) int { return n * 2 }"},
		{Page: "en/a.md", Line: 20, Code: "this is not compiled"},
	}

	result, err := examples.Prove(context.Background(), blocks, examples.Options{Versions: examples.Versions{Go: "1.26"}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Proved != 2 || result.Tests != 1 || result.Illustrative != 1 {
		t.Errorf("result = %+v", result)
	}
}

func TestProveFailsAnExampleThatDoesNotCompileAndNamesTheMarkdownLine(t *testing.T) {
	goAvailable(t)
	var stderr bytes.Buffer
	blocks := []examples.Block{
		{Page: "en/a.md", Line: 14, Proved: true, Code: "func Broken() int { return \"text\" }"},
	}

	_, err := examples.Prove(context.Background(), blocks, examples.Options{Versions: examples.Versions{Go: "1.26"}, Stderr: &stderr})

	if err == nil {
		t.Fatal("a broken example was proved")
	}
	if !strings.Contains(stderr.String(), "en/a.md:14") {
		t.Errorf("the compiler output does not name the Markdown line:\n%s", stderr.String())
	}
}

func TestProveFailsAFailingTest(t *testing.T) {
	goAvailable(t)
	blocks := []examples.Block{
		{Page: "en/a.md", Line: 3, Proved: true, Code: "func TestFails(t *testing.T) { t.Fatal(\"no\") }"},
	}

	_, err := examples.Prove(context.Background(), blocks, examples.Options{Versions: examples.Versions{Go: "1.26"}, Stderr: &bytes.Buffer{}})

	if err == nil {
		t.Fatal("a failing test was proved")
	}
}

func TestParityComparesAnchorsAndFences(t *testing.T) {
	fsys := fstest.MapFS{
		"en/a.md":              {Data: []byte("# A\n\n<a name=\"x\"></a>\n## X\n\n```go example\nfunc A() {}\n```\n")},
		"pt/a.md":              {Data: []byte("# A\n\n<a name=\"x\"></a>\n## X\n\n```go\nfunc A() {}\n```\n")},
		"es/a.md":              {Data: []byte("# A\n\n<a name=\"x\"></a>\n## X\n\n```go example\nfunc A() {}\n```\n")},
		"en/from-laravel/b.md": {Data: []byte("# B\n")},
		"pt/stray.md":          {Data: []byte("# Stray\n")},
	}

	problems := examples.Parity(fsys, "en", []string{"pt", "es"})

	want := []string{
		"pt/a.md: the fenced blocks differ from en/a.md",
		"pt/stray.md: there is no en/stray.md",
	}
	if strings.Join(problems, "\n") != strings.Join(want, "\n") {
		t.Errorf("problems = %q, want %q", problems, want)
	}
}
