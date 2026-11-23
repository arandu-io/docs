package examples

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"sort"
	"strconv"
	"strings"
)

// Import is how a qualifier used in an example is imported.
type Import struct {
	// Alias is written before the path when the qualifier differs from the
	// last element of the path.
	Alias string
	Path  string
}

// Imports maps the qualifier an example uses to its package. An example
// carries no import block, so a qualifier is looked up here; two packages
// that share a name use different qualifiers in the pages (http is net/http,
// fhttp is the framework's and hhttp is Hesape's).
var Imports = map[string]Import{
	"database":      {Path: "github.com/arandu-io/hesape/database"},
	"test":          {Path: "github.com/arandu-io/framework/test"},
	"pagination":    {Path: "github.com/arandu-io/hesape/pagination"},
	"support":       {Path: "github.com/arandu-io/hesape/support"},
	"number":        {Path: "github.com/arandu-io/hesape/number"},
	"schema":        {Path: "github.com/arandu-io/hesape/database/schema"},
	"migrations":    {Path: "github.com/arandu-io/hesape/database/migrations"},
	"io":            {Path: "io"},
	"mcp":           {Path: "github.com/arandu-io/mcp"},
	"observability": {Path: "github.com/arandu-io/framework/observability"},
	"process":       {Path: "github.com/arandu-io/hesape/process"},
	"mail":          {Path: "github.com/arandu-io/framework/mail"},
	"translation":   {Path: "github.com/arandu-io/hesape/translation"},
	"image":         {Path: "github.com/arandu-io/hesape/image"},
	"client":        {Path: "github.com/arandu-io/hesape/http/client"},
	"providers":     {Path: "github.com/arandu-io/hesape/oauth/providers"},
	"bytes":         {Path: "bytes"},
	"context":       {Path: "context"},
	"errors":        {Path: "errors"},
	"fmt":           {Path: "fmt"},
	"http":          {Path: "net/http"},
	"httptest":      {Path: "net/http/httptest"},
	"maps":          {Path: "maps"},
	"os":            {Path: "os"},
	"filepath":      {Path: "path/filepath"},
	"slices":        {Path: "slices"},
	"strconv":       {Path: "strconv"},
	"strings":       {Path: "strings"},
	"testing":       {Path: "testing"},
	"time":          {Path: "time"},
	"url":           {Path: "net/url"},
	"stdhttp":       {Alias: "stdhttp", Path: "net/http"},

	"collections":   {Path: "github.com/arandu-io/hesape/collections"},
	"concurrency":   {Path: "github.com/arandu-io/hesape/concurrency"},
	"fhttp":         {Alias: "fhttp", Path: "github.com/arandu-io/framework/http"},
	"hhttp":         {Alias: "hhttp", Path: "github.com/arandu-io/hesape/http"},
	"security":      {Path: "github.com/arandu-io/framework/security"},
	"str":           {Path: "github.com/arandu-io/hesape/str"},
	"validation":    {Path: "github.com/arandu-io/hesape/validation"},
	"foundation":    {Path: "github.com/arandu-io/framework/foundation"},
	"exception":     {Path: "github.com/arandu-io/hesape/exception"},
	"cache":         {Path: "github.com/arandu-io/hesape/cache"},
	"middleware":    {Path: "github.com/arandu-io/framework/http/middleware"},
	"broadcasting":  {Path: "github.com/arandu-io/hesape/broadcasting"},
	"data":          {Path: "github.com/arandu-io/framework/data"},
	"notifications": {Path: "github.com/arandu-io/hesape/notifications"},
	"messages":      {Path: "github.com/arandu-io/hesape/notifications/messages"},
	"queue":         {Path: "github.com/arandu-io/hesape/queue"},
	"jobs":          {Path: "github.com/arandu-io/hesape/queue/jobs"},
	"auth":          {Path: "github.com/arandu-io/hesape/auth"},
	"session":       {Path: "github.com/arandu-io/hesape/session"},
	"model":         {Path: "github.com/arandu-io/hesape/database/model"},
	"hashing":       {Path: "github.com/arandu-io/hesape/hashing"},
	"encryption":    {Path: "github.com/arandu-io/hesape/encryption"},
	"events":        {Path: "github.com/arandu-io/framework/events"},
}

// File is one generated Go file.
type File struct {
	// Name is the file name inside the example's directory.
	Name string
	// Source is the text of the file.
	Source string
	// Test reports whether the file declares a test.
	Test bool
}

// Wrap turns a block into a compilable file. A block that starts with a
// package clause is used as written. Any other block must hold only
// package-level declarations, and gets a package clause and the imports for the
// qualifiers it uses.
func Wrap(b Block) (File, error) {
	directive := fmt.Sprintf("//line %s:%d:1\n", b.Page, b.Line)
	code := b.Code
	whole := hasPackageClause(code)

	if whole {
		file, err := parser.ParseFile(token.NewFileSet(), "example.go", code, parser.ParseComments)
		if err != nil {
			return File{}, describe(b, err)
		}
		return File{Name: fileName(file), Source: directive + code + "\n", Test: declaresTest(file)}, nil
	}

	source := "package example\n\n" + code
	file, err := parser.ParseFile(token.NewFileSet(), "example.go", source, parser.ParseComments|parser.SkipObjectResolution)
	if err != nil {
		return File{}, describe(b, fmt.Errorf("an example holds package-level declarations or a whole file: %w", err))
	}
	// An import block of the example is kept as written, for the packages
	// imported for their effect; every other qualifier comes from the table.
	declared := map[string]bool{}
	for _, spec := range file.Imports {
		pathValue, _ := strconv.Unquote(spec.Path.Value)
		name := path.Base(pathValue)
		if spec.Name != nil {
			name = spec.Name.Name
		}
		declared[name] = true
	}
	used := qualifiers(code)
	var imports []string
	var unknown []string
	for _, name := range used {
		if declared[name] {
			continue
		}
		known, ok := Imports[name]
		if !ok {
			unknown = append(unknown, name)
			continue
		}
		if known.Alias != "" {
			imports = append(imports, known.Alias+" "+strconv.Quote(known.Path))
		} else {
			imports = append(imports, strconv.Quote(known.Path))
		}
	}
	if len(unknown) > 0 {
		return File{}, describe(b, fmt.Errorf("qualifier(s) %s are not defined by the example and are not in the imports table", strings.Join(unknown, ", ")))
	}
	sort.Strings(imports)
	var out bytes.Buffer
	out.WriteString("package example\n\n")
	if len(imports) > 0 {
		out.WriteString("import (\n")
		for _, line := range imports {
			out.WriteString("\t" + line + "\n")
		}
		out.WriteString(")\n\n")
	}
	out.WriteString(directive)
	out.WriteString(code + "\n")
	name := "example.go"
	test := declaresTest(file)
	if test {
		name = "example_test.go"
	}
	return File{Name: name, Source: out.String(), Test: test}, nil
}

func fileName(file *ast.File) string {
	if declaresTest(file) {
		return "example_test.go"
	}
	return "example.go"
}

func describe(b Block, err error) error { return fmt.Errorf("%s: %w", b.Location(), err) }

func hasPackageClause(code string) bool {
	for _, line := range strings.Split(code, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}
		return strings.HasPrefix(trimmed, "package ")
	}
	return false
}

// declaresTest reports whether the file holds a Test function over *testing.T.
func declaresTest(file *ast.File) bool {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || !strings.HasPrefix(fn.Name.Name, "Test") {
			continue
		}
		return true
	}
	return false
}

// qualifiers returns, sorted, the identifiers used as the left side of a
// selector and not declared anywhere in the code: the package names.
func qualifiers(code string) []string {
	file, err := parser.ParseFile(token.NewFileSet(), "example.go", "package example\n\n"+code, parser.ParseComments)
	if err != nil {
		return nil
	}
	unresolved := map[string]bool{}
	for _, ident := range file.Unresolved {
		unresolved[ident.Name] = true
	}
	found := map[string]bool{}
	ast.Inspect(file, func(n ast.Node) bool {
		selector, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if ident, ok := selector.X.(*ast.Ident); ok && ident.Obj == nil && unresolved[ident.Name] {
			found[ident.Name] = true
		}
		return true
	})
	out := make([]string, 0, len(found))
	for name := range found {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
