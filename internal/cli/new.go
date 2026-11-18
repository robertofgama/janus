package cli

import (
	"fmt"

	"gitlab.com/rg-dev/janus/internal/core"
	"gitlab.com/rg-dev/janus/internal/failure"
	"gitlab.com/rg-dev/janus/internal/messages"
)

// runNew roda `janus new <id> [--dir <pasta>] [--name <nome>]`. Nunca pergunta e nunca
// sobrescreve. As verificações seguem a ordem do contrato: 2, 3, 4, 5, 10, 6, 7, 8 e 1
// (FR-CLI-11).
func runNew(args []string, sys System) error {
	p, err := parseArgs("new", args, map[string]bool{"dir": true, "name": true, "help": false})
	if err != nil {
		return err
	}
	if p.has("help") {
		return printAlone(sys, "new", "--help", args, messages.NewHelp)
	}
	switch {
	case len(p.args) == 0:
		return &failure.UsageError{Problem: failure.MissingArgument, Arg: "<id>", Command: "new"}
	case len(p.args) > 1:
		return &failure.UsageError{Problem: failure.UnexpectedArgument, Arg: p.args[1], Command: "new"}
	}
	cat, err := loadCatalog(sys)
	if err != nil {
		return err
	}
	tpl, err := cat.Find(p.args[0])
	if err != nil {
		return err
	}
	dir, hasDir := p.options["dir"]
	if !hasDir {
		if dir, err = sys.Getwd(); err != nil {
			return err
		}
	}
	destDir, err := core.DestinationDir(dir, cat.Dir)
	if err != nil {
		return err
	}
	name, hasName := p.options["name"]
	finalName, err := core.FinalName(tpl, name, hasName)
	if err != nil {
		return err
	}
	path, err := core.Create(tpl, destDir, finalName)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(sys.Stdout, path)
	return err
}
