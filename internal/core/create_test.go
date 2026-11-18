package core

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"gitlab.com/rg-dev/janus/internal/testfixture"
)

// templateOf lê a pasta de modelos dir e devolve o modelo id.
func templateOf(t *testing.T, dir, id string) Template {
	t.Helper()
	cat, err := LoadCatalog(dir)
	if err != nil {
		t.Fatal(err)
	}
	return find(t, cat, id)
}

func TestDestinationDir(t *testing.T) {
	t.Run("returns the absolute path (FR-CORE-18)", func(t *testing.T) {
		parent := t.TempDir()
		if err := os.Mkdir(filepath.Join(parent, "sub"), 0o755); err != nil {
			t.Fatal(err)
		}
		t.Chdir(parent)
		if got, err := DestinationDir("sub", t.TempDir()); err != nil || got != filepath.Join(parent, "sub") {
			t.Errorf("got %q, %v", got, err)
		}
	})

	t.Run("accepts a link to a dir and returns it unresolved (FR-CORE-18, DE-26)", func(t *testing.T) {
		link := filepath.Join(t.TempDir(), "link")
		symlink(t, t.TempDir(), link)
		if got, err := DestinationDir(link, t.TempDir()); err != nil || got != link {
			t.Errorf("got %q, %v", got, err)
		}
	})

	t.Run("rejects what is not a dir (FR-CORE-18)", func(t *testing.T) {
		file := testfixture.Write(t, t.TempDir(), "arquivo", "", 0o644)
		for _, dir := range []string{"/caminho/inexistente", file, ""} {
			_, err := DestinationDir(dir, t.TempDir())
			var e *DestinationDirError
			if !errors.As(err, &e) || e.Dir != dir || e.Reason != DestNotDir {
				t.Errorf("DestinationDir(%q) = %v", dir, err)
			}
		}
	})

	t.Run("rejects the templates dir and what is inside it (FR-CORE-24)", func(t *testing.T) {
		base := t.TempDir()
		templates := filepath.Join(base, "modelos")
		testfixture.Write(t, templates, "sub/x", "", 0o644)
		testfixture.Write(t, templates, ".git/config", "", 0o644)
		testfixture.Write(t, templates, "..foo/x", "", 0o644)
		// Como o ~/Modelos do autor: a pasta de modelos é um link para a pasta real.
		link := filepath.Join(base, "Modelos")
		symlink(t, templates, link)
		toSub := filepath.Join(base, "atalho")
		symlink(t, filepath.Join(templates, "sub"), toSub)
		cases := []struct{ dir, templatesDir string }{
			{templates, templates},
			{filepath.Join(templates, "sub"), templates},
			{filepath.Join(templates, ".git"), templates},
			{filepath.Join(templates, "..foo"), templates},
			{toSub, templates},
			{link, templates},
			{filepath.Join(templates, "sub"), link},
			{filepath.Join(link, "sub"), link},
		}
		for _, c := range cases {
			_, err := DestinationDir(c.dir, c.templatesDir)
			var e *DestinationDirError
			if !errors.As(err, &e) || e.Reason != DestInsideTemplates || e.Dir != c.dir || e.TemplatesDir != c.templatesDir {
				t.Errorf("DestinationDir(%q, %q) = %v", c.dir, c.templatesDir, err)
			}
		}
	})

	t.Run("accepts the parent and a sibling with the same prefix (FR-CORE-24)", func(t *testing.T) {
		base := t.TempDir()
		templates := filepath.Join(base, "Modelos")
		testfixture.Write(t, templates, "x", "", 0o644)
		sibling := filepath.Join(base, "Modelos-antigo")
		if err := os.Mkdir(sibling, 0o755); err != nil {
			t.Fatal(err)
		}
		for _, dir := range []string{base, sibling} {
			if got, err := DestinationDir(dir, templates); err != nil || got != dir {
				t.Errorf("DestinationDir(%q) = %q, %v", dir, got, err)
			}
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("copies the content byte by byte (FR-CORE-20)", func(t *testing.T) {
		dir := testfixture.New(t)
		testfixture.Write(t, dir, "crlf.txt", "linha 1\r\nlinha 2\r\n", 0o644)
		ids := []string{"info.ai/agents/react.AGENTS.md", "info.dev/react; script; simple.html", "3d/fatiar.3mf", "Empty File", "crlf.txt"}
		for _, id := range ids {
			tpl := templateOf(t, dir, id)
			dest := t.TempDir()
			path, err := Create(tpl, dest, "novo")
			if err != nil {
				t.Fatalf("%s: %v", id, err)
			}
			if path != filepath.Join(dest, "novo") {
				t.Errorf("%s: path %q", id, path)
			}
			if got, want := testfixture.Read(t, path), testfixture.Read(t, tpl.Path); got != want {
				t.Errorf("%s: got %q, want %q", id, got, want)
			}
		}
	})

	t.Run("copies the target of a link (FR-CORE-07)", func(t *testing.T) {
		dir := testfixture.New(t)
		symlink(t, "Empty Bash", filepath.Join(dir, "link"))
		path, err := Create(templateOf(t, dir, "link"), t.TempDir(), "novo")
		if err != nil {
			t.Fatal(err)
		}
		if got, want := testfixture.Read(t, path), testfixture.Read(t, filepath.Join(dir, "Empty Bash")); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("applies the template permissions regardless of the umask (FR-CORE-21)", func(t *testing.T) {
		old := syscall.Umask(0o077)
		t.Cleanup(func() { syscall.Umask(old) })
		dir := testfixture.New(t)
		testfixture.Write(t, dir, "restrito", "x", 0o640)
		for id, want := range map[string]os.FileMode{"Empty Bash": 0o755, "restrito": 0o640} {
			path, err := Create(templateOf(t, dir, id), t.TempDir(), "novo")
			if err != nil {
				t.Fatal(err)
			}
			info, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			if got := info.Mode().Perm(); got != want {
				t.Errorf("%s: got %o, want %o", id, got, want)
			}
		}
	})

	t.Run("never overwrites anything (FR-CORE-19)", func(t *testing.T) {
		tpl := templateOf(t, testfixture.New(t), "info.ai/agents/react.AGENTS.md")
		dest := t.TempDir()
		testfixture.Write(t, dest, "arquivo", "antigo", 0o644)
		if err := os.Mkdir(filepath.Join(dest, "pasta"), 0o755); err != nil {
			t.Fatal(err)
		}
		symlink(t, "nao-existe", filepath.Join(dest, "quebrado"))
		for _, name := range []string{"arquivo", "pasta", "quebrado"} {
			_, err := Create(tpl, dest, name)
			var e *DestinationExistsError
			if !errors.As(err, &e) || e.Path != filepath.Join(dest, name) {
				t.Errorf("%s: got %v", name, err)
			}
		}
		if got := testfixture.Read(t, filepath.Join(dest, "arquivo")); got != "antigo" {
			t.Errorf("the existing file changed to %q", got)
		}
	})

	t.Run("removes the partial file when the copy fails (FR-CORE-22)", func(t *testing.T) {
		copyContent = func(dst io.Writer, src io.Reader) (int64, error) {
			n, _ := io.CopyN(dst, src, 3)
			return n, errors.New("falha simulada")
		}
		t.Cleanup(func() { copyContent = io.Copy })
		dest := t.TempDir()
		_, err := Create(templateOf(t, testfixture.New(t), "info.ai/agents/react.AGENTS.md"), dest, "novo")
		if err == nil || errors.As(err, new(*DestinationExistsError)) {
			t.Fatalf("got %v, want the copy error", err)
		}
		if entries, _ := os.ReadDir(dest); len(entries) != 0 {
			t.Errorf("left %d entries in the destination dir", len(entries))
		}
	})

	t.Run("never writes to the templates dir (NFR-CORE-01)", func(t *testing.T) {
		dir := testfixture.New(t)
		before := testfixture.Snapshot(t, dir)
		dest := t.TempDir()
		for _, id := range testfixture.IDs {
			if _, err := Create(templateOf(t, dir, id), dest, filepath.Base(id)); err != nil {
				t.Fatal(err)
			}
		}
		if after := testfixture.Snapshot(t, dir); after != before {
			t.Errorf("the templates dir changed:\n%s\n%s", before, after)
		}
	})
}
