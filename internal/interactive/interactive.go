// Package interactive é o modo interativo: o seletor, a prévia e o campo de nome. Ele chama
// o núcleo no mesmo processo, nunca o próprio binário. O seletor e o campo ficam atrás de
// interfaces para os testes os substituírem (NFR-TUI-03).
package interactive

import (
	"errors"

	"gitlab.com/rg-dev/janus/internal/core"
	"gitlab.com/rg-dev/janus/internal/failure"
	"gitlab.com/rg-dev/janus/internal/messages"
)

// Selector escolhe um modelo entre os identificadores. Cancelar devolve failure.ErrCanceled.
type Selector interface {
	Select(ids []string) (string, error)
}

// NamePrompt pede o nome do arquivo a partir de um valor inicial editável, mostrando o
// aviso quando houver. Cancelar devolve failure.ErrCanceled.
type NamePrompt interface {
	Ask(initial, warning string) (string, error)
}

// Deps é o que o modo interativo usa de fora.
type Deps struct {
	Getenv     func(string) string
	Getwd      func() (string, error)
	LookPath   func(string) (string, error)
	IsTerminal bool // a entrada padrão e a saída de erro são um terminal
	Selector   Selector
	Prompt     NamePrompt
}

// Run escolhe um modelo, cria o arquivo na pasta atual e devolve o caminho criado. As
// verificações seguem a ordem do contrato: 2, 3, 4, 5, 9 e 10 (FR-TUI-02). O 10, que
// inclui a pasta de modelos (FR-CORE-24), vem antes de abrir o seletor.
func Run(d Deps) (string, error) {
	if !d.IsTerminal {
		return "", &failure.UsageError{Problem: failure.NoTerminal}
	}
	dir, err := core.FindTemplatesDir(d.Getenv)
	if err != nil {
		return "", err
	}
	cat, err := core.LoadCatalog(dir)
	if err != nil {
		return "", err
	}
	if len(cat.Templates) == 0 {
		return "", &failure.NoTemplatesError{Dir: dir}
	}
	if _, err := d.LookPath("fzf"); err != nil {
		return "", &failure.MissingProgramError{Program: "fzf"}
	}
	cwd, err := d.Getwd()
	if err != nil {
		return "", err
	}
	destDir, err := core.DestinationDir(cwd, cat.Dir)
	if err != nil {
		return "", err
	}
	id, err := d.Selector.Select(cat.IDs())
	if err != nil {
		return "", err
	}
	t, err := cat.Find(id)
	if err != nil {
		return "", err
	}
	return create(t, destDir, d.Prompt)
}

// create tenta criar o arquivo e, quando o nome digitado é inválido ou o destino já existe,
// avisa e reabre o campo com o nome que falhou, até dar certo ou o usuário cancelar
// (FR-TUI-04 a FR-TUI-07). Um modelo fixed não pergunta na primeira tentativa.
func create(t core.Template, destDir string, prompt NamePrompt) (string, error) {
	input, warning := t.File, ""
	asking := t.Mode == core.NameAsk
	for {
		if asking {
			typed, err := prompt.Ask(input, warning)
			if err != nil {
				return "", err
			}
			input = typed
		}
		name, err := core.FinalName(t, input, asking)
		// Só o nome digitado pode ser inválido: o nome fixo já foi validado com o manifesto.
		if asking && errors.As(err, new(*core.InvalidNameError)) {
			warning = messages.Describe(err)
			continue
		}
		if err != nil {
			return "", err
		}
		path, err := core.Create(t, destDir, name)
		if errors.As(err, new(*core.DestinationExistsError)) {
			input, warning, asking = name, messages.Describe(err), true
			continue
		}
		return path, err
	}
}
