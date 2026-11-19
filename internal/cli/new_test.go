package cli

import (
	"os"
	"path/filepath"
	"testing"

	"gitlab.com/rg-dev/janus/internal/interactive"
	"gitlab.com/rg-dev/janus/internal/testfixture"
)

func TestNew(t *testing.T) {
	const (
		react       = "info.ai/agents/react.AGENTS.md"
		contestacao = "direito/01.2. Civ. Contestação.tex"
	)

	t.Run("creates the file and prints only its absolute path (FR-CLI-07)", func(t *testing.T) {
		templates, dest := setup(t)
		r := run(t, "new", react)
		want := filepath.Join(dest, "AGENTS.md")
		if r.code != 0 || r.stdout != want+"\n" || r.stderr != "" {
			t.Fatalf("got %+v", r)
		}
		if testfixture.Read(t, want) != testfixture.Read(t, filepath.Join(templates, react)) {
			t.Error("the content differs from the template")
		}
	})

	t.Run("fails with 8 and keeps the existing file (FR-CLI-10)", func(t *testing.T) {
		_, dest := setup(t)
		existing := testfixture.Write(t, dest, "AGENTS.md", "antigo", 0o644)
		run(t, "new", react).wantFailure(t, 8, existing)
		if testfixture.Read(t, existing) != "antigo" {
			t.Error("the existing file changed")
		}
	})

	t.Run("fails with 6 when an ask template has no name (FR-CLI-08)", func(t *testing.T) {
		_, dest := setup(t)
		run(t, "new", contestacao).wantFailure(t, 6, contestacao)
		wantEmpty(t, dest)
	})

	t.Run("applies the name rules to --name (FR-CLI-09)", func(t *testing.T) {
		cases := []struct {
			id   string
			name string
			file string
			perm os.FileMode
		}{
			{contestacao, "Contestação Fulano", "Contestação Fulano.tex", 0o644},
			{contestacao, "Fulano.TEX", "Fulano.TEX", 0o644},
			{"Empty Bash", "deploy", "deploy", 0o755},
			{"info.dev/package.json.monorepo", "pacote", "pacote.json", 0o644},
			{"info.dev/react; script; simple.html", "pagina", "pagina.html", 0o644},
			{react, "CLAUDE", "CLAUDE.md", 0o644},
		}
		for _, c := range cases {
			_, dest := setup(t)
			r := run(t, "new", c.id, "--name", c.name)
			want := filepath.Join(dest, c.file)
			if r.code != 0 || r.stdout != want+"\n" {
				t.Errorf("%s --name %q: got %+v", c.id, c.name, r)
				continue
			}
			if info, err := os.Stat(want); err != nil || info.Mode().Perm() != c.perm {
				t.Errorf("%s: got %v, %v; want permission %o", want, info, err, c.perm)
			}
		}
	})

	t.Run("fails with 7 on invalid names (FR-CLI-09)", func(t *testing.T) {
		_, dest := setup(t)
		for _, name := range []string{"a/b", "..", "  "} {
			run(t, "new", contestacao, "--name", name).wantFailure(t, 7, name)
		}
		wantEmpty(t, dest)
	})

	t.Run("creates in --dir and prints an absolute path (FR-CLI-07)", func(t *testing.T) {
		_, dest := setup(t)
		if err := os.Mkdir(filepath.Join(dest, "sub"), 0o755); err != nil {
			t.Fatal(err)
		}
		r := run(t, "new", react, "--dir", "sub")
		if want := filepath.Join(dest, "sub", "AGENTS.md") + "\n"; r.code != 0 || r.stdout != want {
			t.Errorf("got %+v, want stdout %q", r, want)
		}
	})

	t.Run("accepts options before the id and after -- (FR-CLI-02)", func(t *testing.T) {
		templates, dest := setup(t)
		testfixture.Write(t, templates, "-x.md", "x", 0o644)
		if r := run(t, "new", "--name", "y", "--", "-x.md"); r.code != 0 || r.stdout != filepath.Join(dest, "y.md")+"\n" {
			t.Errorf("got %+v", r)
		}
		run(t, "new", "-x.md").wantFailure(t, 2, "-x.md")
		if r := run(t, "new", "--dir=.", "--name=z", contestacao); r.code != 0 || r.stdout != filepath.Join(dest, "z.tex")+"\n" {
			t.Errorf("got %+v", r)
		}
	})

	t.Run("fails with 10 on a missing or empty --dir (FR-CLI-10)", func(t *testing.T) {
		setup(t)
		run(t, "new", react, "--dir", "/caminho/inexistente").wantFailure(t, 10, "/caminho/inexistente")
		run(t, "new", react, "--dir=").wantFailure(t, 10)
	})

	t.Run("fails with 10 when the destination dir is inside the templates dir (FR-CORE-24)", func(t *testing.T) {
		templates, _ := setup(t)
		link := filepath.Join(t.TempDir(), "Modelos")
		if err := os.Symlink(templates, link); err != nil {
			t.Fatal(err)
		}
		before := testfixture.Snapshot(t, templates)
		for _, dir := range []string{templates, filepath.Join(templates, "info.ai"), link} {
			run(t, "new", react, "--dir", dir).wantFailure(t, 10, dir, templates)
		}
		t.Chdir(filepath.Join(templates, "info.dev"))
		run(t, "new", react).wantFailure(t, 10, templates)
		for _, dir := range []string{".", ".."} {
			run(t, "new", react, "--dir", dir).wantFailure(t, 10, templates)
		}
		if after := testfixture.Snapshot(t, templates); after != before {
			t.Error("the templates dir changed")
		}
	})

	t.Run("fails with 5 on an unknown id (FR-CLI-10)", func(t *testing.T) {
		setup(t)
		run(t, "new", "nao/existe").wantFailure(t, 5, "nao/existe")
	})

	t.Run("fails with 4 on an invalid manifest (FR-CLI-10)", func(t *testing.T) {
		templates, _ := setup(t)
		testfixture.Write(t, templates, ".janus.toml", testfixture.Ambiguous, 0o644)
		run(t, "new", react).wantFailure(t, 4, "info.ai/agents/*", "info.ai/*/react.*", react)
	})

	t.Run("fails with 3 when the templates dir is the home dir (FR-CLI-10)", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		t.Setenv("XDG_TEMPLATES_DIR", home)
		run(t, "new", react).wantFailure(t, 3, home)
	})

	t.Run("fails with 1 when it cannot write (FR-CLI-10)", func(t *testing.T) {
		_, dest := setup(t)
		readonly := filepath.Join(dest, "ro")
		if err := os.Mkdir(readonly, 0o555); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { os.Chmod(readonly, 0o755) })
		run(t, "new", react, "--dir", readonly).wantFailure(t, 1)
	})

	t.Run("rejects what the contract does not accept (FR-CLI-01)", func(t *testing.T) {
		setup(t)
		for _, args := range [][]string{
			{"new"}, {"new", "a", "b"}, {"new", "a", "--name"}, {"new", "a", "--json"},
			{"new", "a", "--dir", "x", "--dir", "y"}, {"new", "a", "-name", "x"}, {"new", "a", "--help"},
		} {
			if r := run(t, args...); r.code != 2 || r.stdout != "" {
				t.Errorf("%q: got %+v", args, r)
			}
		}
	})

	t.Run("prints its help (FR-CLI-03)", func(t *testing.T) {
		r := run(t, "new", "--help")
		if r.code != 0 || r.stdout == "" || r.stderr != "" {
			t.Errorf("got %+v", r)
		}
	})

	t.Run("checks in the order 2, 3, 4, 5, 10, 6, 7, 8 (FR-CLI-11)", func(t *testing.T) {
		t.Setenv("XDG_TEMPLATES_DIR", "/caminho/inexistente")
		run(t, "new", react, "--bogus").wantFailure(t, 2, "--bogus")

		home := t.TempDir()
		testfixture.Write(t, home, ".janus.toml", "[outra]\n", 0o644)
		t.Setenv("HOME", home)
		t.Setenv("XDG_TEMPLATES_DIR", home)
		run(t, "new", "nao/existe").wantFailure(t, 3, home)

		templates, dest := setup(t)
		testfixture.Write(t, templates, ".janus.toml", "[outra]\n", 0o644)
		run(t, "new", "nao/existe").wantFailure(t, 4, "outra")

		testfixture.Write(t, templates, ".janus.toml", testfixture.Manifest, 0o644)
		run(t, "new", "nao/existe", "--dir", "/caminho/inexistente").wantFailure(t, 5, "nao/existe")
		run(t, "new", "nao/existe", "--dir", templates).wantFailure(t, 5, "nao/existe")
		run(t, "new", contestacao, "--dir", "/caminho/inexistente").wantFailure(t, 10, "/caminho/inexistente")
		run(t, "new", contestacao, "--dir", templates).wantFailure(t, 10, templates)

		testfixture.Write(t, dest, "a", "", 0o644)
		run(t, "new", contestacao, "--name", "a/b").wantFailure(t, 7, "a/b")
	})

	t.Run("does not need fzf nor bat (NFR-CLI-02)", func(t *testing.T) {
		_, dest := setup(t)
		t.Setenv("PATH", t.TempDir())
		if r := run(t, "new", react); r.code != 0 || r.stdout != filepath.Join(dest, "AGENTS.md")+"\n" {
			t.Errorf("got %+v", r)
		}
	})

	t.Run("never writes to the templates dir (NFR-CORE-01)", func(t *testing.T) {
		templates, _ := setup(t)
		before := testfixture.Snapshot(t, templates)
		run(t, "new", react)
		run(t, "new", contestacao, "--name", "x")
		run(t, interactive.PreviewCommand, react)
		if after := testfixture.Snapshot(t, templates); after != before {
			t.Error("the templates dir changed")
		}
	})
}
