package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"gitlab.com/rg-dev/janus/internal/testfixture"
)

func TestList(t *testing.T) {
	t.Run("prints one id per line in order (FR-CLI-05)", func(t *testing.T) {
		setup(t)
		r := run(t, "list")
		if want := strings.Join(testfixture.IDs, "\n") + "\n"; r.code != 0 || r.stdout != want || r.stderr != "" {
			t.Errorf("got %+v", r)
		}
	})

	t.Run("prints nothing for an empty dir (FR-CLI-05)", func(t *testing.T) {
		t.Setenv("XDG_TEMPLATES_DIR", t.TempDir())
		if r := run(t, "list"); r.code != 0 || r.stdout != "" {
			t.Errorf("got %+v", r)
		}
	})

	t.Run("prints the JSON of the contract (FR-CLI-06)", func(t *testing.T) {
		setup(t)
		r := run(t, "list", "--json")
		var items []map[string]any
		if err := json.Unmarshal([]byte(r.stdout), &items); err != nil || r.code != 0 {
			t.Fatalf("got %+v, %v", r, err)
		}
		byID := map[string]map[string]any{}
		for _, item := range items {
			keys := make([]string, 0, len(item))
			for k := range item {
				keys = append(keys, k)
			}
			slices.Sort(keys)
			if !slices.Equal(keys, []string{"category", "file", "id", "name"}) {
				t.Errorf("keys %q", keys)
			}
			byID[item["id"].(string)] = item
		}
		if len(byID) != len(testfixture.IDs) {
			t.Fatalf("got %d items", len(byID))
		}
		react := byID["info.ai/agents/react.AGENTS.md"]["name"].(map[string]any)
		if react["mode"] != "fixed" || react["value"] != "AGENTS.md" || react["suggestion"] != nil {
			t.Errorf("react.AGENTS.md name %v", react)
		}
		contestacao := byID["direito/01.2. Civ. Contestação.tex"]
		name := contestacao["name"].(map[string]any)
		if name["mode"] != "ask" || name["suggestion"] != contestacao["file"] || name["value"] != nil {
			t.Errorf("contestação %v", contestacao)
		}
		if byID["Empty Bash"]["category"] != "" {
			t.Errorf("Empty Bash %v", byID["Empty Bash"])
		}
		if !strings.Contains(r.stdout, "Contestação") {
			t.Error("accents are escaped instead of written as UTF-8")
		}
	})

	t.Run("prints [] for an empty dir (FR-CLI-06)", func(t *testing.T) {
		t.Setenv("XDG_TEMPLATES_DIR", t.TempDir())
		if r := run(t, "list", "--json"); r.code != 0 || strings.TrimSpace(r.stdout) != "[]" {
			t.Errorf("got %+v", r)
		}
	})

	t.Run("does not need fzf nor bat (NFR-CLI-02)", func(t *testing.T) {
		setup(t)
		t.Setenv("PATH", t.TempDir())
		if r := run(t, "list"); r.code != 0 {
			t.Errorf("got %+v", r)
		}
	})

	t.Run("fails with 3 when the templates dir is missing (FR-CLI-10)", func(t *testing.T) {
		t.Setenv("XDG_TEMPLATES_DIR", "/caminho/inexistente")
		run(t, "list").wantFailure(t, 3, "/caminho/inexistente")
	})

	t.Run("fails with 4 on an invalid manifest (FR-CLI-10)", func(t *testing.T) {
		templates, _ := setup(t)
		testfixture.Write(t, templates, ".janus.toml", testfixture.Ambiguous, 0o644)
		run(t, "list").wantFailure(t, 4, "info.ai/agents/*", "info.ai/*/react.*", "info.ai/agents/react.AGENTS.md")
		for content, item := range map[string]string{
			"[names]\n\"a/**\" = \"b\"\n": "a/**",
			"[outra]\n":                   "outra",
			"[names]\n\"a\" = 1\n":        "a",
		} {
			testfixture.Write(t, templates, ".janus.toml", content, 0o644)
			run(t, "list").wantFailure(t, 4, item)
		}
	})

	t.Run("fails with 1 when a folder cannot be read (FR-CLI-10)", func(t *testing.T) {
		templates, _ := setup(t)
		locked := filepath.Join(templates, "info.dev")
		if err := os.Chmod(locked, 0); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { os.Chmod(locked, 0o755) })
		run(t, "list").wantFailure(t, 1, "info.dev")
	})

	t.Run("checks in the order 2, 3, 4, 1 (FR-CLI-11)", func(t *testing.T) {
		t.Setenv("XDG_TEMPLATES_DIR", "/caminho/inexistente")
		run(t, "list", "--bogus").wantFailure(t, 2, "--bogus")

		home := t.TempDir()
		testfixture.Write(t, home, ".janus.toml", "[outra]\n", 0o644)
		t.Setenv("HOME", home)
		t.Setenv("XDG_TEMPLATES_DIR", home)
		run(t, "list").wantFailure(t, 3, home)

		templates, _ := setup(t)
		testfixture.Write(t, templates, ".janus.toml", "[outra]\n", 0o644)
		locked := filepath.Join(templates, "info.dev")
		if err := os.Chmod(locked, 0); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { os.Chmod(locked, 0o755) })
		run(t, "list").wantFailure(t, 4, "outra")
	})
}
