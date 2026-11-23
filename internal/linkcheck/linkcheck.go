// Package linkcheck verifies references between the tracked Markdown files in
// the documentation repository.
package linkcheck

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	inlineLink = regexp.MustCompile(`!?\[([^\]]*)\]\(([^)\s]+)(?:\s+"[^"]*")?\)`)
	reference  = regexp.MustCompile(`^\s{0,3}\[([^\]]+)\]:\s*(\S+)`)
	htmlLink   = regexp.MustCompile(`<(?:img|a)\s[^>]*?(?:src|href)\s*=\s*"([^"]+)"`)
	heading    = regexp.MustCompile(`^#{1,6}\s+(.*?)\s*$`)
	codeInText = regexp.MustCompile("`([^`]*)`")
	codeSpan   = regexp.MustCompile("`[^`]*`")
	linkInText = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	fence      = regexp.MustCompile("^\\s*(?:```|~~~)")
	pointer    = regexp.MustCompile(`(?:[\p{L}\p{N}_.-]+/)*[\p{L}\p{N}_.-]+\.md`)
	shorthand  = regexp.MustCompile(`docs/(\d{2})`)
	anchorName = regexp.MustCompile(`<a\s+(?:name|id)\s*=\s*"([^"]+)"\s*>`)
)

// versionPrefix starts every link between usage pages. The site replaces
// {{version}} with the version it serves, so the checker resolves the rest as a
// page of the same language.
const versionPrefix = "/docs/{{version}}/"

// languages are the directories that hold the usage documentation, one per
// language with the same page names in each.
var languages = []string{"en", "pt", "es"}

// languageOf returns the usage-documentation language a file belongs to, or
// the empty string for every other file.
func languageOf(path string) string {
	first, _, found := strings.Cut(path, "/")
	if !found {
		return ""
	}
	for _, language := range languages {
		if first == language {
			return language
		}
	}
	return ""
}

type brokenLink struct {
	path   string
	line   int
	url    string
	reason string
}

type unresolvedPointer struct {
	path  string
	line  int
	name  string
	moved string
}

// Run checks the Git repository containing dir and writes the report to
// stdout. It returns the process exit code the command should use.
func Run(dir string, args []string, stdout, stderr io.Writer) int {
	quiet := false
	for _, arg := range args {
		if arg == "--quiet" {
			quiet = true
		}
	}

	root, err := gitOutput(dir, "rev-parse", "--show-toplevel")
	if err != nil {
		fmt.Fprintf(stderr, "check-links: %v\n", err)
		return 1
	}
	listed, err := gitOutput(root, "ls-files")
	if err != nil {
		fmt.Fprintf(stderr, "check-links: %v\n", err)
		return 1
	}
	tracked := make(map[string]bool)
	for _, path := range strings.Fields(listed) {
		tracked[path] = true
	}
	files := make([]string, 0)
	for path := range tracked {
		if strings.HasSuffix(path, ".md") {
			files = append(files, path)
		}
	}
	sort.Strings(files)

	checked := 0
	broken := make([]brokenLink, 0)
	pointers := make([]unresolvedPointer, 0)
	anchors := make(map[string]map[string]bool)
	for _, path := range files {
		lines, err := readLines(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			fmt.Fprintf(stderr, "check-links: %v\n", err)
			return 1
		}
		directory := filepath.Dir(path)
		fenced := false
		for number, raw := range lines {
			if fence.MatchString(raw) {
				fenced = !fenced
				continue
			}
			if fenced {
				continue
			}
			line := codeSpan.ReplaceAllStringFunc(raw, func(match string) string {
				return strings.Repeat(" ", len(match))
			})
			for _, url := range linkURLs(line) {
				if isExternal(url) {
					continue
				}
				language := languageOf(path)
				if language != "" && strings.HasPrefix(url, "/") && !strings.HasPrefix(url, versionPrefix) {
					checked++
					broken = append(broken, brokenLink{path, number + 1, url, "site link without " + versionPrefix})
					continue
				}
				if strings.HasPrefix(url, "#") && language == "" {
					continue
				}
				target, fragment, _ := strings.Cut(url, "#")
				target, _, _ = strings.Cut(target, "?")
				if target == "" && language == "" {
					continue
				}
				checked++
				resolved := filepath.Clean(filepath.Join(directory, filepath.FromSlash(target)))
				if target == "" {
					resolved = path
				}
				if strings.HasPrefix(url, versionPrefix) {
					resolved = language + "/" + strings.TrimPrefix(target, versionPrefix) + ".md"
				}
				resolved = filepath.ToSlash(resolved)
				if !tracked[resolved] {
					broken = append(broken, brokenLink{path, number + 1, url, "no such file"})
					continue
				}
				if fragment != "" && strings.HasSuffix(resolved, ".md") {
					if anchors[resolved] == nil {
						anchors[resolved], err = headings(filepath.Join(root, filepath.FromSlash(resolved)))
						if err != nil {
							fmt.Fprintf(stderr, "check-links: %v\n", err)
							return 1
						}
					}
					if !anchors[resolved][strings.ToLower(fragment)] {
						broken = append(broken, brokenLink{path, number + 1, url, "no such heading"})
					}
				}
			}
			for _, match := range pointerMatches(path, raw) {
				if !pointerBoundary(raw, match[0], match[1]) {
					continue
				}
				named := raw[match[0]:match[1]]
				candidates := []string{
					filepath.ToSlash(filepath.Clean(filepath.Join(directory, filepath.FromSlash(named)))),
					filepath.ToSlash(filepath.Clean(filepath.FromSlash(named))),
				}
				if tracked[candidates[0]] || tracked[candidates[1]] {
					continue
				}
				pointers = append(pointers, unresolvedPointer{
					path: path, line: number + 1, name: named, moved: guess(named, tracked),
				})
			}
			for _, match := range shorthandMatches(path, raw) {
				if !shorthandBoundary(raw, match[0], match[1]) {
					continue
				}
				pointers = append(pointers, unresolvedPointer{
					path: path, line: number + 1, name: raw[match[0]:match[1]],
					moved: guess(raw[match[2]:match[3]]+"-", tracked),
				})
			}
		}
	}
	inside := make([]unresolvedPointer, 0)
	outside := make([]unresolvedPointer, 0)
	for _, item := range pointers {
		if strings.HasPrefix(item.path, "_archive/") {
			inside = append(inside, item)
		} else {
			outside = append(outside, item)
		}
	}

	if !quiet {
		fmt.Fprintf(stdout, "%d markdown files, %d internal links\n", len(files), checked)
		fmt.Fprintf(stdout, "%d pointers name a file that is not there: %d inside _archive/, %d outside it\n", len(pointers), len(inside), len(outside))
		writePointerCounts(stdout, pointers)
		fmt.Fprintln(stdout)
		fmt.Fprintln(stdout, "_archive/ is frozen, so the pointers written inside it stay as they are.")
		fmt.Fprintln(stdout, "They are counted here and they do not fail the run.")
		fmt.Fprintln(stdout)
	}
	for _, item := range outside {
		where := ""
		if item.moved != "" {
			where = ", moved to " + item.moved
		}
		fmt.Fprintf(stdout, "%s:%d: pointer names %s, which is not there%s\n", item.path, item.line, item.name, where)
	}
	for _, item := range broken {
		fmt.Fprintf(stdout, "%s:%d: %s: %s\n", item.path, item.line, item.reason, item.url)
	}
	if len(broken) != 0 || len(outside) != 0 {
		fmt.Fprintf(stdout, "\n%d broken link(s), %d unresolved pointer(s) outside _archive/.\n", len(broken), len(outside))
		return 1
	}
	if !quiet {
		fmt.Fprintln(stdout, "every markdown link resolves.")
	}
	return 0
}

// pointerMatches finds the file names a line mentions. The usage pages name
// files of the projects they describe, such as UPGRADE.md, so only the
// repository's own records are checked for pointers.
func pointerMatches(path, line string) [][]int {
	if languageOf(path) != "" {
		return nil
	}
	return pointer.FindAllStringIndex(line, -1)
}

func shorthandMatches(path, line string) [][]int {
	if languageOf(path) != "" {
		return nil
	}
	return shorthand.FindAllStringSubmatchIndex(line, -1)
}

func pointerBoundary(line string, start, end int) bool {
	if start > 0 {
		char, _ := utf8.DecodeLastRuneInString(line[:start])
		if isWord(char) || char == '/' || char == '-' {
			return false
		}
	}
	if end < len(line) {
		char, _ := utf8.DecodeRuneInString(line[end:])
		if isWord(char) {
			return false
		}
	}
	return true
}

func shorthandBoundary(line string, start, end int) bool {
	if !pointerBoundary(line, start, end) {
		return false
	}
	if end < len(line) {
		char, _ := utf8.DecodeRuneInString(line[end:])
		if char == '.' || char == '-' {
			return false
		}
	}
	return true
}

func isWord(char rune) bool {
	// The former checker's Unicode-aware word class meant letters, numbers, and
	// underscore for the repository names and headings this checker accepts.
	return unicode.IsLetter(char) || unicode.IsNumber(char) || char == '_'
}

func guess(named string, tracked map[string]bool) string {
	base := filepath.Base(filepath.FromSlash(named))
	exact := make([]string, 0)
	prefix := make([]string, 0)
	for path := range tracked {
		candidate := filepath.Base(filepath.FromSlash(path))
		if candidate == base {
			exact = append(exact, path)
		}
		if strings.HasPrefix(candidate, base) {
			prefix = append(prefix, path)
		}
	}
	if len(exact) == 1 {
		return exact[0]
	}
	if len(exact) == 0 && len(prefix) == 1 {
		return prefix[0]
	}
	return ""
}

func writePointerCounts(output io.Writer, pointers []unresolvedPointer) {
	type count struct {
		path  string
		value int
		order int
	}
	byPath := make(map[string]*count)
	ordered := make([]*count, 0)
	for _, item := range pointers {
		entry := byPath[item.path]
		if entry == nil {
			entry = &count{path: item.path, order: len(ordered)}
			byPath[item.path] = entry
			ordered = append(ordered, entry)
		}
		entry.value++
	}
	sort.SliceStable(ordered, func(left, right int) bool {
		if ordered[left].value != ordered[right].value {
			return ordered[left].value > ordered[right].value
		}
		return ordered[left].order < ordered[right].order
	})
	limit := min(10, len(ordered))
	for _, item := range ordered[:limit] {
		fmt.Fprintf(output, "%8d  %s\n", item.value, item.path)
	}
	if len(ordered) > limit {
		fmt.Fprintf(output, "    and %d more files\n", len(ordered)-limit)
	}
}

func linkURLs(line string) []string {
	var urls []string
	for _, match := range inlineLink.FindAllStringSubmatch(line, -1) {
		urls = append(urls, match[2])
	}
	for _, match := range htmlLink.FindAllStringSubmatch(line, -1) {
		urls = append(urls, match[1])
	}
	if match := reference.FindStringSubmatch(line); match != nil {
		urls = append(urls, match[2])
	}
	return urls
}

func gitOutput(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(output)), nil
}

func readLines(path string) ([]string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(contents), "\n"), nil
}

func headings(path string) (map[string]bool, error) {
	lines, err := readLines(path)
	if err != nil {
		return nil, err
	}
	result := make(map[string]bool)
	for _, line := range lines {
		for _, named := range anchorName.FindAllStringSubmatch(line, -1) {
			result[strings.ToLower(named[1])] = true
		}
		match := heading.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		text := codeInText.ReplaceAllString(match[1], "$1")
		text = linkInText.ReplaceAllString(text, "$1")
		result[slug(text)] = true
	}
	return result, nil
}

func slug(text string) string {
	var builder strings.Builder
	for _, char := range strings.ToLower(text) {
		if isWord(char) || char == '-' || char == ' ' {
			builder.WriteRune(char)
		}
	}
	return strings.ReplaceAll(strings.TrimSpace(builder.String()), " ", "-")
}

func isExternal(url string) bool {
	return strings.HasPrefix(url, "http:") || strings.HasPrefix(url, "https:") ||
		strings.HasPrefix(url, "mailto:") || strings.HasPrefix(url, "tel:") ||
		strings.HasPrefix(url, "ftp:")
}
