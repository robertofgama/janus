package interactive

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"

	"gitlab.com/rg-dev/janus/internal/core"
	"gitlab.com/rg-dev/janus/internal/failure"
	"gitlab.com/rg-dev/janus/internal/testfixture"
)

// fakeSelector escolhe pick, ou devolve err, e guarda os identificadores recebidos.
type fakeSelector struct {
	pick string
	err  error
	ids  []string
}

func (s *fakeSelector) Select(ids []string) (string, error) {
	s.ids = ids
	return s.pick, s.err
}

// prompt é uma chamada ao campo de nome.
type prompt struct {
	initial string
	warning string
}

// fakePrompt responde com replies, na ordem, e cancela quando elas acabam.
type fakePrompt struct {
	replies []string
	calls   []prompt
}

func (p *fakePrompt) Ask(initial, warning string) (string, error) {
	p.calls = append(p.calls, prompt{initial, warning})
	if len(p.replies) == 0 {
		return "", failure.ErrCanceled
	}
	reply := p.replies[0]
	p.replies = p.replies[1:]
	return reply, nil
}

// setup aponta XDG_TEMPLATES_DIR para a pasta de teste e devolve as dependências com os
// substitutos e a pasta de destino, vazia. O PATH fica vazio: o laço não pode depender do
// fzf (NFR-TUI-03).
func setup(t *testing.T, sel *fakeSelector, p *fakePrompt) (Deps, string) {
	t.Helper()
	testfixture.Use(t)
	t.Setenv("PATH", t.TempDir())
	dest := t.TempDir()
	return Deps{
		Getenv:     os.Getenv,
		Getwd:      func() (string, error) { return dest, nil },
		LookPath:   func(string) (string, error) { return "/usr/bin/fzf", nil },
		IsTerminal: true,
		Selector:   sel,
		Prompt:     p,
	}, dest
}

// wantEmpty confere que dir está vazia.
func wantEmpty(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("%s has %d entries, want none", dir, len(entries))
	}
}

const (
	react       = "info.ai/agents/react.AGENTS.md"
	contestacao = "direito/01.2. Civ. Contestação.tex"
)

func TestRun(t *testing.T) {
	t.Run("creates a fixed template without asking (FR-TUI-04)", func(t *testing.T) {
		sel, p := &fakeSelector{pick: react}, &fakePrompt{}
		d, dest := setup(t, sel, p)
		path, err := Run(d)
		if want := filepath.Join(dest, "AGENTS.md"); err != nil || path != want {
			t.Fatalf("got %q, %v; want %q", path, err, want)
		}
		if len(p.calls) != 0 {
			t.Errorf("asked %d times", len(p.calls))
		}
		if !slices.Equal(sel.ids, testfixture.IDs) {
			t.Errorf("the selector got %q", sel.ids)
		}
	})

	t.Run("asks the name of an ask template with the suggestion (FR-TUI-05)", func(t *testing.T) {
		p := &fakePrompt{replies: []string{"Fulano"}}
		d, dest := setup(t, &fakeSelector{pick: contestacao}, p)
		path, err := Run(d)
		if want := filepath.Join(dest, "Fulano.tex"); err != nil || path != want {
			t.Fatalf("got %q, %v; want %q", path, err, want)
		}
		if want := []prompt{{"01.2. Civ. Contestação.tex", ""}}; !slices.Equal(p.calls, want) {
			t.Errorf("got calls %q", p.calls)
		}
	})

	t.Run("reopens the prompt with the name that collided (FR-TUI-06)", func(t *testing.T) {
		p := &fakePrompt{replies: []string{"AGENTS.md", "OUTRO"}}
		d, dest := setup(t, &fakeSelector{pick: react}, p)
		existing := testfixture.Write(t, dest, "AGENTS.md", "antigo", 0o644)
		path, err := Run(d)
		if want := filepath.Join(dest, "OUTRO.md"); err != nil || path != want {
			t.Fatalf("got %q, %v; want %q", path, err, want)
		}
		if len(p.calls) != 2 {
			t.Fatalf("got calls %q", p.calls)
		}
		for _, call := range p.calls {
			if call.initial != "AGENTS.md" || call.warning == "" {
				t.Errorf("got call %q", call)
			}
		}
		if got := testfixture.Read(t, existing); got != "antigo" {
			t.Errorf("the existing file changed to %q", got)
		}
	})

	t.Run("reopens the prompt with the invalid name (FR-TUI-07)", func(t *testing.T) {
		p := &fakePrompt{replies: []string{"a/b", "ok"}}
		d, dest := setup(t, &fakeSelector{pick: contestacao}, p)
		path, err := Run(d)
		if want := filepath.Join(dest, "ok.tex"); err != nil || path != want {
			t.Fatalf("got %q, %v; want %q", path, err, want)
		}
		if len(p.calls) != 2 || p.calls[1].initial != "a/b" || p.calls[1].warning == "" {
			t.Errorf("got calls %q", p.calls)
		}
	})

	t.Run("returns canceled when the selector is canceled (FR-TUI-08)", func(t *testing.T) {
		d, dest := setup(t, &fakeSelector{err: failure.ErrCanceled}, &fakePrompt{})
		if _, err := Run(d); !errors.Is(err, failure.ErrCanceled) {
			t.Errorf("got %v", err)
		}
		wantEmpty(t, dest)
	})

	t.Run("returns canceled when the prompt is canceled (FR-TUI-08)", func(t *testing.T) {
		d, dest := setup(t, &fakeSelector{pick: contestacao}, &fakePrompt{})
		if _, err := Run(d); !errors.Is(err, failure.ErrCanceled) {
			t.Errorf("got %v", err)
		}
		wantEmpty(t, dest)
	})

	t.Run("requires a terminal (FR-TUI-01)", func(t *testing.T) {
		d, _ := setup(t, &fakeSelector{pick: react}, &fakePrompt{})
		d.IsTerminal = false
		_, err := Run(d)
		var e *failure.UsageError
		if !errors.As(err, &e) || e.Problem != failure.NoTerminal {
			t.Errorf("got %v", err)
		}
	})

	t.Run("fails when there are no templates (FR-TUI-03)", func(t *testing.T) {
		sel := &fakeSelector{pick: react}
		d, _ := setup(t, sel, &fakePrompt{})
		empty := t.TempDir()
		testfixture.Write(t, empty, ".oculto", "", 0o644)
		t.Setenv("XDG_TEMPLATES_DIR", empty)
		_, err := Run(d)
		var e *failure.NoTemplatesError
		if !errors.As(err, &e) || e.Dir != empty {
			t.Errorf("got %v", err)
		}
		if sel.ids != nil {
			t.Error("the selector was opened")
		}
	})

	t.Run("fails when fzf is not in the PATH (FR-TUI-02)", func(t *testing.T) {
		d, _ := setup(t, &fakeSelector{pick: react}, &fakePrompt{})
		t.Setenv("PATH", t.TempDir())
		d.LookPath = exec.LookPath
		_, err := Run(d)
		var e *failure.MissingProgramError
		if !errors.As(err, &e) || e.Program != "fzf" {
			t.Errorf("got %v", err)
		}
	})

	t.Run("checks in the order 2, 3, 4, 5, 9 (FR-TUI-02)", func(t *testing.T) {
		d, _ := setup(t, &fakeSelector{pick: react}, &fakePrompt{})
		d.LookPath = func(string) (string, error) { return "", exec.ErrNotFound }

		d.IsTerminal = false
		t.Setenv("XDG_TEMPLATES_DIR", "/caminho/inexistente")
		if _, err := Run(d); !errors.As(err, new(*failure.UsageError)) {
			t.Errorf("no terminal and no templates dir: got %v", err)
		}

		d.IsTerminal = true
		empty := t.TempDir()
		testfixture.Write(t, empty, core.ManifestFile, "[outra]\n", 0o644)
		t.Setenv("XDG_TEMPLATES_DIR", empty)
		if _, err := Run(d); !errors.As(err, new(*core.ManifestError)) {
			t.Errorf("invalid manifest and no templates: got %v", err)
		}

		testfixture.Write(t, empty, core.ManifestFile, "", 0o644)
		if _, err := Run(d); !errors.As(err, new(*failure.NoTemplatesError)) {
			t.Errorf("no templates and no fzf: got %v", err)
		}
	})
}
