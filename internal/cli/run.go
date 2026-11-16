package cli

import (
	"bufio"
	"fmt"
	"io"
	"runtime/debug"
	"strings"

	"gitlab.com/rg-dev/janus/internal/failure"
	"gitlab.com/rg-dev/janus/internal/interactive"
	"gitlab.com/rg-dev/janus/internal/messages"
)

// Run executa o janus com args, os argumentos sem o nome do programa, e devolve o código
// de saída.
func Run(args []string, sys System) int {
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		return finish(sys, runCommand(args[0], args[1:], sys), false)
	}
	p, err := parseArgs("", args, map[string]bool{"hold": false, "help": false, "version": false})
	switch {
	case err != nil:
		return finish(sys, err, false)
	case p.has("help"):
		return finish(sys, printAlone(sys, "", "--help", args, messages.Help), false)
	case p.has("version"):
		return finish(sys, printAlone(sys, "", "--version", args, version()+"\n"), false)
	case len(p.args) > 0:
		return finish(sys, &failure.UsageError{Problem: failure.UnexpectedArgument, Arg: p.args[0]}, false)
	}
	return finish(sys, runInteractive(sys), p.has("hold"))
}

// runCommand roda um subcomando.
func runCommand(name string, args []string, sys System) error {
	switch name {
	case "list":
		return runList(args, sys)
	case "new":
		return runNew(args, sys)
	case interactive.PreviewCommand:
		return runPreview(args, sys)
	}
	return &failure.UsageError{Problem: failure.UnknownCommand, Arg: name}
}

// printAlone imprime text quando option é o único argumento. Acompanhada de outros, ela é
// uso incorreto (FR-CLI-01).
func printAlone(sys System, command, option string, args []string, text string) error {
	for _, arg := range args {
		if arg != option {
			return &failure.UsageError{Problem: failure.UnexpectedArgument, Arg: arg, Command: command}
		}
	}
	_, err := io.WriteString(sys.Stdout, text)
	return err
}

// version devolve a versão que o Go grava no binário, a partir do git (DE-22).
func version() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return "(devel)"
}

// finish mostra o erro, quando há, e devolve o código de saída. Cancelar não mostra nada.
// Com hold, espera um Enter depois do erro, para o yazi não esconder a mensagem (FR-TUI-13).
func finish(sys System, err error, hold bool) int {
	code := exitCode(err)
	if code == codeOK || code == codeCanceled {
		return code
	}
	fmt.Fprintln(sys.Stderr, messages.ErrorLine(err))
	if hold {
		fmt.Fprint(sys.Stderr, messages.HoldPrompt)
		_, _ = bufio.NewReader(sys.Stdin).ReadString('\n')
	}
	return code
}
