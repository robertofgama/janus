package core

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"gitlab.com/rg-dev/janus/internal/testfixture"
)

// loadIDs lê a pasta de modelos dir e devolve os identificadores.
func loadIDs(t *testing.T, dir string) []string {
	t.Helper()
	cat, err := LoadCatalog(dir)
	if err != nil {
		t.Fatal(err)
	}
	return cat.IDs()
}

// find devolve o modelo id do catálogo.
func find(t *testing.T, cat Catalog, id string) Template {
	t.Helper()
	tpl, err := cat.Find(id)
	if err != nil {
		t.Fatal(err)
	}
	return tpl
}

// symlink cria o link simbólico link apontando para target.
func symlink(t *testing.T, target, link string) {
	t.Helper()
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
}

func TestLoadCatalog(t *testing.T) {
	t.Run("lists every template in order (FR-CORE-05, FR-CORE-09)", func(t *testing.T) {
		if got := loadIDs(t, testfixture.New(t)); !slices.Equal(got, testfixture.IDs) {
			t.Errorf("got  %q\nwant %q", got, testfixture.IDs)
		}
	})

	t.Run("fills category, file and path (FR-CORE-05)", func(t *testing.T) {
		dir := testfixture.New(t)
		cat, err := LoadCatalog(dir)
		if err != nil {
			t.Fatal(err)
		}
		react := find(t, cat, "info.ai/agents/react.AGENTS.md")
		if react.Category != "info.ai/agents" || react.File != "react.AGENTS.md" || react.Path != filepath.Join(dir, react.ID) {
			t.Errorf("got %+v", react)
		}
		if bash := find(t, cat, "Empty Bash"); bash.Category != "" {
			t.Errorf("Empty Bash has category %q", bash.Category)
		}
	})

	t.Run("ignores hidden items at any level (FR-CORE-06)", func(t *testing.T) {
		dir := testfixture.New(t)
		testfixture.Write(t, dir, "info.dev/.oculto", "x", 0o644)
		testfixture.Write(t, dir, ".pasta/dentro.txt", "x", 0o644)
		if got := loadIDs(t, dir); !slices.Equal(got, testfixture.IDs) {
			t.Errorf("got %q", got)
		}
	})

	t.Run("handles links inside the dir (FR-CORE-07)", func(t *testing.T) {
		dir := testfixture.New(t)
		symlink(t, "Empty Bash", filepath.Join(dir, "link"))
		symlink(t, "nao-existe", filepath.Join(dir, "quebrado"))
		symlink(t, "info.dev", filepath.Join(dir, "pasta-link"))
		want := append(slices.Clone(testfixture.IDs), "link")
		if got := loadIDs(t, dir); !slices.Equal(got, want) {
			t.Errorf("got %q", got)
		}
	})

	t.Run("ignores names with control characters (FR-CORE-08)", func(t *testing.T) {
		dir := testfixture.New(t)
		testfixture.Write(t, dir, "a\nb", "x", 0o644)
		testfixture.Write(t, dir, "c\td", "x", 0o644)
		if got := loadIDs(t, dir); !slices.Equal(got, testfixture.IDs) {
			t.Errorf("got %q", got)
		}
	})

	t.Run("breaks ties by byte order (FR-CORE-09)", func(t *testing.T) {
		dir := t.TempDir()
		for _, name := range []string{"b", "B", "a"} {
			testfixture.Write(t, dir, name, "", 0o644)
		}
		if got := loadIDs(t, dir); !slices.Equal(got, []string{"a", "B", "b"}) {
			t.Errorf("got %q", got)
		}
	})

	t.Run("follows the templates dir when it is a link (FR-CORE-04)", func(t *testing.T) {
		link := filepath.Join(t.TempDir(), "Modelos")
		symlink(t, testfixture.New(t), link)
		if got := loadIDs(t, link); !slices.Equal(got, testfixture.IDs) {
			t.Errorf("got %q", got)
		}
	})

	t.Run("resolves the fixed names of the manifest (FR-CORE-13)", func(t *testing.T) {
		cat, err := LoadCatalog(testfixture.New(t))
		if err != nil {
			t.Fatal(err)
		}
		fixed := map[string]string{
			"info.ai/agents/AGENTS.md":                 "AGENTS.md",
			"info.ai/agents/react.AGENTS.md":           "AGENTS.md",
			"info.ai/ai-memory/default.ai-memory.toml": ".ai-memory.toml",
			"info.dev/package.json.monorepo":           "package.json",
		}
		for _, tpl := range cat.Templates {
			name, isFixed := fixed[tpl.ID]
			if isFixed && (tpl.Mode != NameFixed || tpl.FixedName != name) {
				t.Errorf("%s: got %s %q, want fixed %q", tpl.ID, tpl.Mode, tpl.FixedName, name)
			}
			if !isFixed && tpl.Mode != NameAsk {
				t.Errorf("%s: got %s, want ask", tpl.ID, tpl.Mode)
			}
		}
	})

	t.Run("leaves every template in ask mode without a manifest (FR-CORE-10)", func(t *testing.T) {
		dir := testfixture.New(t)
		if err := os.Remove(filepath.Join(dir, ManifestFile)); err != nil {
			t.Fatal(err)
		}
		cat, err := LoadCatalog(dir)
		if err != nil {
			t.Fatal(err)
		}
		for _, tpl := range cat.Templates {
			if tpl.Mode != NameAsk {
				t.Errorf("%s is %s", tpl.ID, tpl.Mode)
			}
		}
	})

	t.Run("fails when two wildcard patterns match a template (FR-CORE-13)", func(t *testing.T) {
		dir := testfixture.New(t)
		testfixture.Write(t, dir, ManifestFile, testfixture.Ambiguous, 0o644)
		_, err := LoadCatalog(dir)
		var e *ManifestError
		if !errors.As(err, &e) || e.Problem != ManifestAmbiguous {
			t.Fatalf("got %v", err)
		}
		if e.ID != "info.ai/agents/react.AGENTS.md" || !slices.Equal(e.Patterns, []string{"info.ai/agents/*", "info.ai/*/react.*"}) {
			t.Errorf("got id %q, patterns %q", e.ID, e.Patterns)
		}
	})

	t.Run("lets an exact pattern win over wildcards (FR-CORE-13)", func(t *testing.T) {
		dir := testfixture.New(t)
		manifest := testfixture.Ambiguous + "\"info.ai/agents/react.AGENTS.md\" = \"X.md\"\n"
		testfixture.Write(t, dir, ManifestFile, manifest, 0o644)
		cat, err := LoadCatalog(dir)
		if err != nil {
			t.Fatal(err)
		}
		if react := find(t, cat, "info.ai/agents/react.AGENTS.md"); react.Mode != NameFixed || react.FixedName != "X.md" {
			t.Errorf("got %+v", react)
		}
	})

	t.Run("accepts a pattern that matches nothing (FR-CORE-12)", func(t *testing.T) {
		dir := testfixture.New(t)
		testfixture.Write(t, dir, ManifestFile, "[names]\n\"nada/*\" = \"x\"\n", 0o644)
		if _, err := LoadCatalog(dir); err != nil {
			t.Error(err)
		}
	})
}

func TestCatalogFind(t *testing.T) {
	cat, err := LoadCatalog(testfixture.New(t))
	if err != nil {
		t.Fatal(err)
	}

	t.Run("finds a listed template", func(t *testing.T) {
		if tpl := find(t, cat, "Empty Bash"); tpl.File != "Empty Bash" {
			t.Errorf("got %+v", tpl)
		}
	})

	t.Run("does not build paths from the id (FR-CORE-15)", func(t *testing.T) {
		for _, id := range []string{"nao/existe", "../x", "/etc/hostname", ".janus.toml", ".git/config", "./Empty Bash", "empty bash"} {
			_, err := cat.Find(id)
			var e *TemplateNotFoundError
			if !errors.As(err, &e) || e.ID != id {
				t.Errorf("Find(%q) = %v", id, err)
			}
		}
	})
}
