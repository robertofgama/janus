package interactive

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"gitlab.com/rg-dev/janus/internal/core"
	"gitlab.com/rg-dev/janus/internal/messages"
)

// PreviewCommand é o subcomando interno que o fzf chama para a prévia, com o
// identificador. Não aparece na ajuda e não é contrato (DE-12).
const PreviewCommand = "__preview"

// PreviewDeps é o que a prévia usa de fora.
type PreviewDeps struct {
	Getenv   func(string) string
	LookPath func(string) (string, error)
	Stdout   io.Writer
}

// Preview escreve a prévia do modelo id (FR-TUI-12): texto realçado pelo bat quando ele
// está no PATH, ou cru; uma linha de aviso para binário e para vazio. O identificador
// passa pela listagem, como em new; nunca vira caminho direto (FR-CORE-15).
func Preview(id string, d PreviewDeps) error {
	dir, err := core.FindTemplatesDir(d.Getenv)
	if err != nil {
		return err
	}
	cat, err := core.LoadCatalog(dir)
	if err != nil {
		return err
	}
	t, err := cat.Find(id)
	if err != nil {
		return err
	}
	kind, size, err := core.Inspect(t.Path)
	if err != nil {
		return err
	}
	switch kind {
	case core.PreviewEmpty:
		_, err = fmt.Fprintln(d.Stdout, messages.PreviewEmpty)
		return err
	case core.PreviewBinary:
		_, err = fmt.Fprintln(d.Stdout, messages.PreviewBinary(size))
		return err
	}
	if bat, err := d.LookPath("bat"); err == nil {
		// O caminho vai como argumento, sem shell: nomes com espaço e ";" chegam intactos.
		cmd := exec.Command(bat, "--color=always", "--style=plain", "--paging=never", t.Path)
		cmd.Stdout = d.Stdout
		cmd.Stderr = d.Stdout
		return cmd.Run()
	}
	f, err := os.Open(t.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(d.Stdout, f)
	return err
}
