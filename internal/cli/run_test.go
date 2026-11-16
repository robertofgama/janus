package cli

import (
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/rg-dev/janus/internal/interactive"
	"gitlab.com/rg-dev/janus/internal/testfixture"
)

func TestRun(t *testing.T) {
	t.Run("rejects what the contract does not accept (FR-CLI-01)", func(t *testing.T) {
		setup(t)
		for _, args := range [][]string{
			{"foo"}, {"list", "x"}, {"list", "--dir", "x"}, {"list", "--hold"}, {"list", "--json=1"},
			{"list", "--help", "--json"}, {"--hold", "--hold"}, {"--hold", "list"}, {"-h"},
			{"--version", "x"}, {"--help", "--hold"}, {"--bogus"},
		} {
			r := run(t, args...)
			if r.code != 2 || r.stdout != "" {
				t.Errorf("%q: got %+v", args, r)
			}
		}
	})

	t.Run("cites the unknown subcommand (FR-CLI-12)", func(t *testing.T) {
		run(t, "foo").wantFailure(t, 2, "foo")
	})

	t.Run("prints the help on stdout (FR-CLI-03)", func(t *testing.T) {
		for _, args := range [][]string{{"--help"}, {"list", "--help"}} {
			r := run(t, args...)
			if r.code != 0 || r.stdout == "" || r.stderr != "" || strings.Contains(r.stdout, interactive.PreviewCommand) {
				t.Errorf("%q: got %+v", args, r)
			}
		}
	})

	t.Run("prints the version (FR-CLI-04)", func(t *testing.T) {
		if r := run(t, "--version"); r.code != 0 || r.stdout != "(devel)\n" || r.stderr != "" {
			t.Errorf("got %+v", r)
		}
	})
}

func TestInteractiveMode(t *testing.T) {
	const react = "info.ai/agents/react.AGENTS.md"

	t.Run("prints only the created path (FR-TUI-09)", func(t *testing.T) {
		_, dest := setup(t)
		testfixture.FakeProgram(t, "fzf", testfixture.Reply{Stdout: react + "\n"})
		r := run(t)
		if want := filepath.Join(dest, "AGENTS.md") + "\n"; r.code != 0 || r.stdout != want {
			t.Errorf("got %+v, want stdout %q", r, want)
		}
	})

	t.Run("exits with 130 and no message when canceled (FR-TUI-08)", func(t *testing.T) {
		_, dest := setup(t)
		testfixture.FakeProgram(t, "fzf", testfixture.Reply{Code: 130})
		if r := run(t); r.code != 130 || r.stdout != "" || r.stderr != "" {
			t.Errorf("got %+v", r)
		}
		wantEmpty(t, dest)
	})

	t.Run("requires a terminal on stdin and stderr (FR-TUI-01)", func(t *testing.T) {
		setup(t)
		runWith(t, func(s *System) { s.StdinIsTerminal = false }).wantFailure(t, 2, "janus new")
		runWith(t, func(s *System) { s.StderrIsTerminal = false }).wantFailure(t, 2, "janus new")
	})

	t.Run("fails with 9 without fzf (FR-TUI-02)", func(t *testing.T) {
		setup(t)
		t.Setenv("PATH", t.TempDir())
		run(t).wantFailure(t, 9, "fzf")
	})

	t.Run("waits for Enter after an error with --hold (FR-TUI-13)", func(t *testing.T) {
		templates, _ := setup(t)
		testfixture.Write(t, templates, ".janus.toml", "[outra]\n", 0o644)
		stdin := strings.NewReader("\n")
		runWith(t, func(s *System) { s.Stdin = stdin }, "--hold").wantFailure(t, 4, "outra")
		if stdin.Len() != 0 {
			t.Error("did not wait for the Enter")
		}
	})

	t.Run("does not wait when canceled or without --hold (FR-TUI-13)", func(t *testing.T) {
		templates, _ := setup(t)
		testfixture.FakeProgram(t, "fzf", testfixture.Reply{Code: 130})
		stdin := strings.NewReader("\n")
		if r := runWith(t, func(s *System) { s.Stdin = stdin }, "--hold"); r.code != 130 || stdin.Len() != 1 {
			t.Errorf("canceled: got %+v, %d bytes left", r, stdin.Len())
		}
		testfixture.Write(t, templates, ".janus.toml", "[outra]\n", 0o644)
		if r := runWith(t, func(s *System) { s.Stdin = stdin }); r.code != 4 || stdin.Len() != 1 {
			t.Errorf("without --hold: got %+v, %d bytes left", r, stdin.Len())
		}
	})
}

func TestPreviewCommand(t *testing.T) {
	t.Run("prints the preview of an id with spaces and semicolons (FR-TUI-12)", func(t *testing.T) {
		templates, _ := setup(t)
		t.Setenv("PATH", t.TempDir())
		const id = "info.dev/react; script; simple.html"
		r := run(t, interactive.PreviewCommand, id)
		if want := testfixture.Read(t, filepath.Join(templates, id)); r.code != 0 || r.stdout != want {
			t.Errorf("got %+v, want stdout %q", r, want)
		}
	})

	t.Run("takes an id that starts with a dash as it is (FR-TUI-12)", func(t *testing.T) {
		templates, _ := setup(t)
		t.Setenv("PATH", t.TempDir())
		testfixture.Write(t, templates, "-x.md", "traço\n", 0o644)
		if r := run(t, interactive.PreviewCommand, "-x.md"); r.code != 0 || r.stdout != "traço\n" {
			t.Errorf("got %+v", r)
		}
	})
}
