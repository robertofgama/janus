package main

import (
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/rg-dev/janus/internal/testfixture"
)

// build compila o janus numa pasta temporária e devolve o caminho do binário.
func build(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "janus")
	if out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	return bin
}

// exitCode devolve o código de saída de err, ou 0.
func exitCode(err error) int {
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return exit.ExitCode()
	}
	return 0
}

func TestBinary(t *testing.T) {
	bin := build(t)
	testfixture.Use(t)

	t.Run("lists the templates (NFR-CLI-05)", func(t *testing.T) {
		out, err := exec.Command(bin, "list").Output()
		if want := strings.Join(testfixture.IDs, "\n") + "\n"; err != nil || string(out) != want {
			t.Errorf("got %q, %v", out, err)
		}
	})

	t.Run("creates a file (NFR-CLI-05)", func(t *testing.T) {
		dest := t.TempDir()
		out, err := exec.Command(bin, "new", "info.ai/agents/react.AGENTS.md", "--dir", dest).Output()
		if want := filepath.Join(dest, "AGENTS.md") + "\n"; err != nil || string(out) != want {
			t.Errorf("got %q, %v", out, err)
		}
	})

	t.Run("exits with the code of the error (NFR-CLI-05)", func(t *testing.T) {
		if code := exitCode(exec.Command(bin, "new", "nao/existe").Run()); code != 5 {
			t.Errorf("got code %d, want 5", code)
		}
	})

	t.Run("prints a version (FR-CLI-04)", func(t *testing.T) {
		out, err := exec.Command(bin, "--version").Output()
		if err != nil || strings.TrimSpace(string(out)) == "" || strings.Count(string(out), "\n") != 1 {
			t.Errorf("got %q, %v", out, err)
		}
	})

	t.Run("refuses the interactive mode without a terminal (FR-TUI-01)", func(t *testing.T) {
		var stderr strings.Builder
		cmd := exec.Command(bin)
		cmd.Stderr = &stderr
		if code := exitCode(cmd.Run()); code != 2 || !strings.Contains(stderr.String(), "janus new") {
			t.Errorf("got code %d, stderr %q", code, stderr.String())
		}
	})
}

func TestDependencies(t *testing.T) {
	t.Run("has BurntSushi/toml as the only direct dependency (NFR-CLI-04)", func(t *testing.T) {
		out, err := exec.Command("go", "list", "-m", "-f", "{{if not (or .Main .Indirect)}}{{.Path}}{{end}}", "all").Output()
		if err != nil {
			t.Fatal(err)
		}
		if got := strings.Fields(string(out)); len(got) != 1 || got[0] != "github.com/BurntSushi/toml" {
			t.Errorf("direct dependencies %q", got)
		}
	})
}
