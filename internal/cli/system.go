package cli

import (
	"io"
	"os"
	"os/exec"
)

// System reúne o que o programa recebe do sistema operacional. O main usa OSSystem; os
// testes montam um com buffers e um terminal fingido.
type System struct {
	Stdin            io.Reader
	Stdout           io.Writer
	Stderr           io.Writer
	Getenv           func(string) string
	Getwd            func() (string, error)
	LookPath         func(string) (string, error)
	Executable       string // o próprio binário, que a prévia do fzf chama
	StdinIsTerminal  bool
	StderrIsTerminal bool
}

// OSSystem monta o System do processo.
func OSSystem() System {
	executable, err := os.Executable()
	if err != nil {
		executable = "janus" // sem o caminho, a prévia recorre ao janus do PATH
	}
	return System{
		Stdin:            os.Stdin,
		Stdout:           os.Stdout,
		Stderr:           os.Stderr,
		Getenv:           os.Getenv,
		Getwd:            os.Getwd,
		LookPath:         exec.LookPath,
		Executable:       executable,
		StdinIsTerminal:  isTerminal(os.Stdin),
		StderrIsTerminal: isTerminal(os.Stderr),
	}
}
