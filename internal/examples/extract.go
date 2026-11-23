// Package examples proves the Go examples of the usage documentation.
//
// A fenced block whose info string is "go example" is a promise: it compiles,
// and when it declares a test, the test passes, against the published tags of
// the framework, Hesape and Kyse. A block fenced as plain "go" is illustration
// and is counted, never compiled.
package examples

import (
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

// Block is one fenced Go block of a documentation page.
type Block struct {
	// Page is the slash-separated path of the Markdown file.
	Page string
	// Line is the line of the first line of code, counted from one.
	Line int
	// Code is the text between the fences, without the final newline.
	Code string
	// Proved reports whether the fence carries the word "example".
	Proved bool
}

// Location is where the block starts, in the form compilers print.
func (b Block) Location() string { return fmt.Sprintf("%s:%d", b.Page, b.Line) }

// Extract returns every Go block of every Markdown file below dir of fsys, in
// path order, and an error for a fence that marks as an example something that
// is not Go.
func Extract(fsys fs.FS, dir string) ([]Block, error) {
	var pages []string
	err := fs.WalkDir(fsys, dir, func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && strings.HasSuffix(name, ".md") {
			pages = append(pages, name)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(pages)

	var blocks []Block
	for _, page := range pages {
		data, err := fs.ReadFile(fsys, page)
		if err != nil {
			return nil, err
		}
		found, err := parse(page, string(data))
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, found...)
	}
	return blocks, nil
}

// parse reads the fenced blocks of one page.
func parse(page, text string) ([]Block, error) {
	var blocks []Block
	var current *Block
	var body []string
	for number, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case current == nil && strings.HasPrefix(trimmed, "```"):
			fields := strings.Fields(strings.TrimPrefix(trimmed, "```"))
			current = nil
			if len(fields) == 0 {
				current = &Block{Page: "", Line: -1}
				body = nil
				continue
			}
			language, rest := fields[0], fields[1:]
			proved := false
			for _, word := range rest {
				if word == "example" {
					proved = true
				}
			}
			if proved && language != "go" {
				return nil, fmt.Errorf("%s:%d: only a go block can be an example, not %q", page, number+1, language)
			}
			if language == "go" {
				current = &Block{Page: page, Line: number + 2, Proved: proved}
			} else {
				current = &Block{Page: "", Line: -1}
			}
			body = nil
		case current != nil && strings.HasPrefix(trimmed, "```"):
			if current.Page != "" {
				current.Code = strings.Join(body, "\n")
				blocks = append(blocks, *current)
			}
			current = nil
		case current != nil:
			body = append(body, line)
		}
	}
	if current != nil {
		return nil, fmt.Errorf("%s: a fenced block is never closed", page)
	}
	return blocks, nil
}

// Language returns the first path segment of a page, such as "en".
func Language(page string) string {
	first, _, _ := strings.Cut(path.Clean(page), "/")
	return first
}
