package cli

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"gitlab.com/rg-dev/janus/internal/testfixture"
)

// result é o que uma execução de Run produziu.
type result struct {
	code   int
	stdout string
	stderr string
}

// runWith roda o janus com args depois de change ajustar o System de teste: ambiente e
// PATH do processo (os testes os ajustam com t.Setenv), uma linha na entrada e um
// terminal fingido.
func runWith(t *testing.T, change func(*System), args ...string) result {
	t.Helper()
	var stdout, stderr bytes.Buffer
	sys := System{
		Stdin:            strings.NewReader("\n"),
		Stdout:           &stdout,
		Stderr:           &stderr,
		Getenv:           os.Getenv,
		Getwd:            os.Getwd,
		LookPath:         exec.LookPath,
		Executable:       "/usr/local/bin/janus",
		StdinIsTerminal:  true,
		StderrIsTerminal: true,
	}
	change(&sys)
	code := Run(args, sys)
	return result{code: code, stdout: stdout.String(), stderr: stderr.String()}
}

// run roda o janus com args e o System de teste padrão.
func run(t *testing.T, args ...string) result {
	t.Helper()
	return runWith(t, func(*System) {}, args...)
}

// setup monta a pasta de teste, aponta XDG_TEMPLATES_DIR para ela e entra numa pasta de
// destino vazia.
func setup(t *testing.T) (templates, dest string) {
	t.Helper()
	templates = testfixture.Use(t)
	dest = t.TempDir()
	t.Chdir(dest)
	return templates, dest
}

// wantFailure confere uma falha: o código, a saída padrão vazia e cada item citado na
// saída de erro (FR-CLI-12, NFR-CLI-01). A frase não é conferida.
func (r result) wantFailure(t *testing.T, code int, items ...string) {
	t.Helper()
	if r.code != code {
		t.Errorf("code %d, want %d (stderr %q)", r.code, code, r.stderr)
	}
	if r.stdout != "" {
		t.Errorf("stdout %q, want empty", r.stdout)
	}
	for _, item := range items {
		if !strings.Contains(r.stderr, item) {
			t.Errorf("stderr %q does not cite %q", r.stderr, item)
		}
	}
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
