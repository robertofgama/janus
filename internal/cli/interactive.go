package cli

import (
	"fmt"

	"gitlab.com/rg-dev/janus/internal/failure"
	"gitlab.com/rg-dev/janus/internal/interactive"
)

// runInteractive roda o modo interativo com o fzf e imprime o caminho criado, e só ele, na
// saída padrão; o resto vai para o terminal (FR-TUI-09).
func runInteractive(sys System) error {
	path, err := interactive.Run(interactive.Deps{
		Getenv:     sys.Getenv,
		Getwd:      sys.Getwd,
		LookPath:   sys.LookPath,
		IsTerminal: sys.StdinIsTerminal && sys.StderrIsTerminal,
		Selector:   interactive.FzfSelector{Executable: sys.Executable, Stderr: sys.Stderr},
		Prompt:     interactive.FzfPrompt{Stderr: sys.Stderr},
	})
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(sys.Stdout, path)
	return err
}

// runPreview roda o subcomando interno de prévia, que o fzf chama com o identificador
// (DE-12). O identificador é tomado como veio, sem ler opções: ele pode começar com "-".
func runPreview(args []string, sys System) error {
	switch {
	case len(args) == 0:
		return &failure.UsageError{Problem: failure.MissingArgument, Arg: "<id>"}
	case len(args) > 1:
		return &failure.UsageError{Problem: failure.UnexpectedArgument, Arg: args[1]}
	}
	return interactive.Preview(args[0], interactive.PreviewDeps{
		Getenv:   sys.Getenv,
		LookPath: sys.LookPath,
		Stdout:   sys.Stdout,
	})
}
