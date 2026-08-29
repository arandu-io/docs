package linkcheck_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/arandu-io/docs/internal/linkcheck"
)

func TestRunAcceptsALinkToAUnicodeHeading(t *testing.T) {
	repo := trackedRepository(t, map[string]string{
		"README.md": "# Index\n\n[Guide](guide.md#configuração-rápida)\n",
		"guide.md":  "# Configuração Rápida\n",
	})

	var stdout, stderr bytes.Buffer
	code := linkcheck.Run(repo, nil, &stdout, &stderr)

	if code != 0 {
		t.Errorf("exit code = %d, want 0; stderr = %q", code, stderr.String())
	}
	want := "2 markdown files, 1 internal links\n" +
		"0 pointers name a file that is not there: 0 inside _archive/, 0 outside it\n\n" +
		"_archive/ is frozen, so the pointers written inside it stay as they are.\n" +
		"They are counted here and they do not fail the run.\n\n" +
		"every markdown link resolves.\n"
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
}

func TestRunChecksEveryLinkFormAndStripsQueries(t *testing.T) {
	repo := trackedRepository(t, map[string]string{
		"README.md": "# Index\n\n" +
			"[Inline](guide.md?view=full#configuração-rápida \"Guide\")\n" +
			"![Image](logo.png)\n" +
			"[reference]: guide.md#configuração-rápida\n" +
			"<a href=\"guide.md#configuração-rápida\">Guide</a>\n" +
			"<img src=\"logo.png\">\n" +
			"[External](https://example.com/guide.md) [Local](#index)\n",
		"guide.md": "# Configuração Rápida\n",
		"logo.png": "not an image decoder's concern",
	})

	var stdout, stderr bytes.Buffer
	code := linkcheck.Run(repo, nil, &stdout, &stderr)

	if code != 0 {
		t.Errorf("exit code = %d, want 0; stderr = %q", code, stderr.String())
	}
	wantPrefix := "2 markdown files, 5 internal links\n"
	if !bytes.HasPrefix(stdout.Bytes(), []byte(wantPrefix)) {
		t.Errorf("stdout = %q, want prefix %q", stdout.String(), wantPrefix)
	}
}

func TestRunIgnoresLinksInFencesAndCodeSpans(t *testing.T) {
	repo := trackedRepository(t, map[string]string{
		"README.md": "# Index\n\n" +
			"`[generic](missing)`\n\n" +
			"```markdown\n[Broken](missing.md)\n```\n\n" +
			"[Guide](guide.md)\n",
		"guide.md": "# Guide\n",
	})

	var stdout, stderr bytes.Buffer
	code := linkcheck.Run(repo, nil, &stdout, &stderr)

	if code != 0 {
		t.Errorf("exit code = %d, want 0; stdout = %q; stderr = %q", code, stdout.String(), stderr.String())
	}
	wantPrefix := "2 markdown files, 1 internal links\n"
	if !bytes.HasPrefix(stdout.Bytes(), []byte(wantPrefix)) {
		t.Errorf("stdout = %q, want prefix %q", stdout.String(), wantPrefix)
	}
}

func TestRunCountsArchivePointersButFailsOnlyOutsideTheArchive(t *testing.T) {
	repo := trackedRepository(t, map[string]string{
		"README.md":            "The old guide is missing.md.\nThe numbered one is docs/07:42.\nhttps://example.com/ignored.md\n",
		"_archive/07-guide.md": "# Seven\n",
		"_archive/missing.md":  "# Moved\n",
		"_archive/source.md":   "The frozen pointer is ghost.md.\n",
	})

	var stdout, stderr bytes.Buffer
	code := linkcheck.Run(repo, nil, &stdout, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1; stderr = %q", code, stderr.String())
	}
	want := "4 markdown files, 0 internal links\n" +
		"3 pointers name a file that is not there: 1 inside _archive/, 2 outside it\n" +
		"       2  README.md\n" +
		"       1  _archive/source.md\n\n" +
		"_archive/ is frozen, so the pointers written inside it stay as they are.\n" +
		"They are counted here and they do not fail the run.\n\n" +
		"README.md:1: pointer names missing.md, which is not there, moved to _archive/missing.md\n" +
		"README.md:2: pointer names docs/07, which is not there, moved to _archive/07-guide.md\n\n" +
		"0 broken link(s), 2 unresolved pointer(s) outside _archive/.\n"
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
}

func TestRunQuietStillReportsBrokenFilesAndHeadings(t *testing.T) {
	repo := trackedRepository(t, map[string]string{
		"README.md": "[Missing](missing.png)\n[Bad heading](guide.md#absent)\n",
		"guide.md":  "# Present\n",
	})

	var stdout, stderr bytes.Buffer
	code := linkcheck.Run(repo, []string{"ignored", "--quiet"}, &stdout, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1; stderr = %q", code, stderr.String())
	}
	want := "README.md:1: no such file: missing.png\n" +
		"README.md:2: no such heading: guide.md#absent\n\n" +
		"2 broken link(s), 0 unresolved pointer(s) outside _archive/.\n"
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
}

func TestRunQuietAllowsPointersInsideTheArchive(t *testing.T) {
	repo := trackedRepository(t, map[string]string{
		"README.md":          "# Index\n",
		"_archive/source.md": "The frozen pointer is ghost.md.\n",
	})

	var stdout, stderr bytes.Buffer
	code := linkcheck.Run(repo, []string{"--quiet"}, &stdout, &stderr)

	if code != 0 {
		t.Errorf("exit code = %d, want 0; stderr = %q", code, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout = %q, want no output", stdout.String())
	}
}

func TestRunBuildsHeadingSlugsFromCodeAndLinkText(t *testing.T) {
	repo := trackedRepository(t, map[string]string{
		"README.md":    "[Guide](guide.md#configuração-rápida-agora)\n",
		"elsewhere.md": "# Elsewhere\n",
		"guide.md":     "# Configuração `Rápida` [Agora](elsewhere.md)!\n",
	})

	var stdout, stderr bytes.Buffer
	code := linkcheck.Run(repo, []string{"--quiet"}, &stdout, &stderr)

	if code != 0 {
		t.Errorf("exit code = %d, want 0; stdout = %q; stderr = %q", code, stdout.String(), stderr.String())
	}
}

func trackedRepository(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, contents := range files {
		path := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	git(t, dir, "init", "-q")
	git(t, dir, "add", ".")
	return dir
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
}
