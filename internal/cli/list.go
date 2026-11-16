package cli

import (
	"encoding/json"
	"io"
	"strings"

	"gitlab.com/rg-dev/janus/internal/core"
	"gitlab.com/rg-dev/janus/internal/failure"
	"gitlab.com/rg-dev/janus/internal/messages"
)

// runList roda `janus list [--json]` (FR-CLI-05, FR-CLI-06). A ordem de verificação é 2,
// 3, 4, 1 (FR-CLI-11).
func runList(args []string, sys System) error {
	p, err := parseArgs("list", args, map[string]bool{"json": false, "help": false})
	if err != nil {
		return err
	}
	if p.has("help") {
		return printAlone(sys, "list", "--help", args, messages.ListHelp)
	}
	if len(p.args) > 0 {
		return &failure.UsageError{Problem: failure.UnexpectedArgument, Arg: p.args[0], Command: "list"}
	}
	cat, err := loadCatalog(sys)
	if err != nil {
		return err
	}
	if p.has("json") {
		return writeJSON(sys.Stdout, cat.Templates)
	}
	var b strings.Builder
	for _, t := range cat.Templates {
		b.WriteString(t.ID)
		b.WriteByte('\n')
	}
	_, err = io.WriteString(sys.Stdout, b.String())
	return err
}

// loadCatalog descobre a pasta de modelos e a lê (códigos 3 e 4).
func loadCatalog(sys System) (core.Catalog, error) {
	dir, err := core.FindTemplatesDir(sys.Getenv)
	if err != nil {
		return core.Catalog{}, err
	}
	return core.LoadCatalog(dir)
}

// jsonTemplate é um item de `janus list --json`. Os nomes das chaves são contrato.
type jsonTemplate struct {
	ID       string   `json:"id"`
	Category string   `json:"category"`
	File     string   `json:"file"`
	Name     jsonName `json:"name"`
}

type jsonName struct {
	Mode       string `json:"mode"`
	Value      string `json:"value,omitempty"`      // só em fixed
	Suggestion string `json:"suggestion,omitempty"` // só em ask
}

// writeJSON escreve a lista do contrato (FR-CLI-06).
func writeJSON(w io.Writer, templates []core.Template) error {
	items := make([]jsonTemplate, 0, len(templates)) // vazia sai como [], e não como null
	for _, t := range templates {
		name := jsonName{Mode: string(t.Mode)}
		if t.Mode == core.NameFixed {
			name.Value = t.FixedName
		} else {
			name.Suggestion = t.File
		}
		items = append(items, jsonTemplate{ID: t.ID, Category: t.Category, File: t.File, Name: name})
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false) // "<", ">" e "&" saem como são
	enc.SetIndent("", "  ")
	return enc.Encode(items)
}
