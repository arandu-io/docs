package main

import (
	"os"

	"github.com/arandu-io/docs/internal/linkcheck"
)

func main() {
	os.Exit(linkcheck.Run(".", os.Args[1:], os.Stdout, os.Stderr))
}
