// Command janus cria arquivos a partir da pasta XDG de modelos, no terminal, no yazi e no
// neovim. A lógica está em internal/; aqui o processo só é repassado para cli.Run.
package main

import (
	"os"

	"gitlab.com/rg-dev/janus/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], cli.OSSystem()))
}
