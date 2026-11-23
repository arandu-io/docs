package examples

import (
	"fmt"
	"io/fs"
	"regexp"
	"sort"
	"strings"
)

var anchorLine = regexp.MustCompile(`^<a name="([^"]+)"></a>$`)

// Parity checks that every translation is the same page as the English one:
// the same files, the same anchors in the same order, and the same fenced
// blocks, with the same fence, so that a translation can never promise an
// example the source does not prove or hide one it does.
func Parity(fsys fs.FS, source string, translations []string) []string {
	pages := func(dir string) map[string]string {
		out := map[string]string{}
		_ = fs.WalkDir(fsys, dir, func(name string, entry fs.DirEntry, err error) error {
			if err == nil && !entry.IsDir() && strings.HasSuffix(name, ".md") {
				out[strings.TrimPrefix(name, dir+"/")] = name
			}
			return nil
		})
		return out
	}
	english := pages(source)
	var problems []string
	for _, language := range translations {
		translated := pages(language)
		for name := range translated {
			if _, ok := english[name]; !ok {
				problems = append(problems, fmt.Sprintf("%s: there is no %s/%s", translated[name], source, name))
			}
		}
		names := make([]string, 0, len(english))
		for name := range english {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			other, ok := translated[name]
			if !ok {
				if englishOnly(fsys, english[name]) {
					continue
				}
				problems = append(problems, fmt.Sprintf("%s: missing, %s exists", language+"/"+name, english[name]))
				continue
			}
			a, _ := fs.ReadFile(fsys, english[name])
			b, _ := fs.ReadFile(fsys, other)
			if !equal(anchors(string(a)), anchors(string(b))) {
				problems = append(problems, fmt.Sprintf("%s: the anchors differ from %s", other, english[name]))
			}
			if !equal(fences(string(a)), fences(string(b))) {
				problems = append(problems, fmt.Sprintf("%s: the fenced blocks differ from %s", other, english[name]))
			}
		}
	}
	sort.Strings(problems)
	return problems
}

// englishOnly reports the pages that carry no translation by decision: the
// pages that tell a reader coming from Laravel why an expectation has no
// answer here.
func englishOnly(fsys fs.FS, name string) bool {
	return strings.HasPrefix(name, "en/from-laravel/")
}

func anchors(text string) []string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		if m := anchorLine.FindStringSubmatch(line); m != nil {
			out = append(out, m[1])
		}
	}
	return out
}

// fences returns every fenced block, with its info string, as one string each.
func fences(text string) []string {
	var out []string
	var current []string
	open := false
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			if open {
				current = append(current, line)
				out = append(out, strings.Join(current, "\n"))
				current = nil
			} else {
				current = []string{line}
			}
			open = !open
			continue
		}
		if open {
			current = append(current, line)
		}
	}
	return out
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
