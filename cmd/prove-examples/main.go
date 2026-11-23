// Command prove-examples compiles the Go examples of the usage documentation.
//
// Every fenced block whose info string is "go example" in en/ is written into a
// temporary module that requires the tags named in versions.json, then vetted
// and built, and tested when it declares a test. The pages in pt/ and es/ must
// carry the same anchors and fenced blocks as the English ones. A block that
// does not compile fails the command with the Markdown file and line.
//
// Run it from the repository root:
//
//	GOWORK=off go run ./cmd/prove-examples
//	GOWORK=off go run ./cmd/prove-examples -latest
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/arandu-io/docs/internal/examples"
)

func main() {
	os.Exit(run())
}

func run() int {
	versionsFile := flag.String("versions", "versions.json", "the tags the examples are proved against")
	latest := flag.Bool("latest", false, "prove against the newest tag of every module instead of the pinned ones")
	keep := flag.Bool("keep", false, "keep the temporary module and print where it is")
	dir := flag.String("dir", "", "write the temporary module here instead of a new temporary directory")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	root := os.DirFS(".")
	failed := false
	for _, problem := range examples.Parity(root, "en", []string{"pt", "es"}) {
		fmt.Fprintln(os.Stderr, problem)
		failed = true
	}

	blocks, err := examples.Extract(root, "en")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	versions, err := examples.LoadVersions(*versionsFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	result, err := examples.Prove(ctx, blocks, examples.Options{
		Versions: versions, Latest: *latest, Dir: *dir, Keep: *keep,
		Stdout: os.Stdout, Stderr: os.Stderr,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		failed = true
	}
	fmt.Printf("%d examples proved (%d of them tests), %d illustrative blocks left uncompiled\n",
		result.Proved, result.Tests, result.Illustrative)
	if failed {
		return 1
	}
	return 0
}
