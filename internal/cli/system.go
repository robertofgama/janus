package cli

import "io"

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
