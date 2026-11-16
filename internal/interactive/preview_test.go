package interactive

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/rg-dev/janus/internal/core"
	"gitlab.com/rg-dev/janus/internal/testfixture"
)

// previewOf roda a prévia de id e devolve a saída.
func previewOf(t *testing.T, id string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	err := Preview(id, PreviewDeps{Getenv: os.Getenv, LookPath: exec.LookPath, Stdout: &out})
	return out.String(), err
}

func TestPreview(t *testing.T) {
	const html = "info.dev/react; script; simple.html"

	t.Run("prints the raw text without bat (FR-TUI-12)", func(t *testing.T) {
		dir := testfixture.Use(t)
		t.Setenv("PATH", t.TempDir())
		got, err := previewOf(t, html)
		if want := testfixture.Read(t, filepath.Join(dir, html)); err != nil || got != want {
			t.Errorf("got %q, %v; want %q", got, err, want)
		}
	})

	t.Run("hands the path to bat as a single argument (FR-TUI-12)", func(t *testing.T) {
		dir := testfixture.Use(t)
		fake := testfixture.FakeProgram(t, "bat", testfixture.Reply{Stdout: "realçado"})
		got, err := previewOf(t, html)
		if err != nil || got != "realçado" {
			t.Fatalf("got %q, %v", got, err)
		}
		args := fake.Calls(t)[0].Args
		if last := args[len(args)-1]; last != filepath.Join(dir, html) {
			t.Errorf("bat got %q", args)
		}
	})

	t.Run("prints one warning line for a binary and for an empty file (FR-TUI-12)", func(t *testing.T) {
		testfixture.Use(t)
		t.Setenv("PATH", t.TempDir())
		for _, id := range []string{"3d/fatiar.3mf", "Empty File"} {
			got, err := previewOf(t, id)
			if err != nil || strings.Count(got, "\n") != 1 || strings.Contains(got, "fatiar") {
				t.Errorf("%s: got %q, %v", id, got, err)
			}
		}
	})

	t.Run("fails on an id that is not listed (FR-CORE-15)", func(t *testing.T) {
		testfixture.Use(t)
		if _, err := previewOf(t, "../x"); !errors.As(err, new(*core.TemplateNotFoundError)) {
			t.Errorf("got %v", err)
		}
	})
}
