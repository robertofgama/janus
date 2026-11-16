// Package messages concentra todo texto mostrado ao usuário, em português do Brasil
// (DE-11, DE-21). O texto não é contrato: os testes conferem o código de saída e o item
// citado, nunca a frase.
package messages

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gitlab.com/rg-dev/janus/internal/core"
	"gitlab.com/rg-dev/janus/internal/failure"
)

const (
	// SelectorLabel é o rótulo do seletor.
	SelectorLabel = "Modelo> "
	// NameLabel é o rótulo do campo de nome.
	NameLabel = "Nome do arquivo> "
	// HoldPrompt pede o Enter de --hold (FR-TUI-13).
	HoldPrompt = "Tecle Enter para sair. "
	// PreviewEmpty é a prévia de um modelo vazio (FR-TUI-12).
	PreviewEmpty = "(arquivo vazio)"
)

// PreviewBinary é a prévia de um modelo binário, com o tamanho (FR-TUI-12).
func PreviewBinary(size int64) string {
	return fmt.Sprintf("(arquivo binário de %s, sem prévia)", formatSize(size))
}

// ErrorLine é a linha escrita na saída de erro quando um comando falha.
func ErrorLine(err error) string {
	return "janus: " + Describe(err)
}

// Describe explica err em português, citando o item envolvido (FR-CLI-12). A ordem dos
// casos importa: o manifesto vem antes do nome inválido porque um nome fixo inválido chega
// embrulhado num ManifestError.
func Describe(err error) string {
	var (
		templatesDir *core.TemplatesDirError
		manifest     *core.ManifestError
		notFound     *core.TemplateNotFoundError
		destDir      *core.DestinationDirError
		required     *core.NameRequiredError
		invalid      *core.InvalidNameError
		exists       *core.DestinationExistsError
		usage        *failure.UsageError
		missing      *failure.MissingProgramError
		noTemplates  *failure.NoTemplatesError
		program      *failure.ProgramError
	)
	switch {
	case errors.As(err, &templatesDir):
		return describeTemplatesDir(templatesDir)
	case errors.As(err, &manifest):
		return describeManifest(manifest)
	case errors.As(err, &notFound):
		return fmt.Sprintf("modelo não encontrado: %q. Veja os modelos com janus list.", notFound.ID)
	case errors.As(err, &destDir):
		return fmt.Sprintf("pasta de destino não encontrada: %q não existe ou não é uma pasta.", destDir.Dir)
	case errors.As(err, &required):
		return fmt.Sprintf("o modelo %q pede um nome: informe --name.", required.ID)
	case errors.As(err, &invalid):
		return fmt.Sprintf("nome inválido %q: %s.", invalid.Name, nameProblem(invalid.Problem))
	case errors.As(err, &exists):
		return fmt.Sprintf("já existe: %s. Use outro nome.", exists.Path)
	case errors.As(err, &usage):
		return describeUsage(usage)
	case errors.As(err, &missing):
		return fmt.Sprintf("o modo interativo precisa do %s, que não está no PATH. Sem ele, use janus list e janus new.", missing.Program)
	case errors.As(err, &noTemplates):
		return fmt.Sprintf("a pasta de modelos %s não tem nenhum modelo.", noTemplates.Dir)
	case errors.As(err, &program):
		return fmt.Sprintf("o %s terminou com o código %d.", program.Program, program.Code)
	}
	return fmt.Sprintf("erro inesperado: %v", err)
}

func describeTemplatesDir(e *core.TemplatesDirError) string {
	where := "em XDG_TEMPLATES_DIR"
	if e.Source == core.SourceUserDirs {
		where = "em " + e.File
	}
	switch e.Reason {
	case core.DirFileMissing:
		return fmt.Sprintf("pasta de modelos não encontrada: XDG_TEMPLATES_DIR está vazia e %s não pode ser lido.", e.File)
	case core.DirLineMissing:
		return fmt.Sprintf("pasta de modelos não encontrada: %s não define XDG_TEMPLATES_DIR.", e.File)
	case core.DirRelative:
		return fmt.Sprintf("pasta de modelos inválida: %q, %s, não é um caminho absoluto nem começa por $HOME/.", e.Value, where)
	case core.DirNotDir:
		return fmt.Sprintf("pasta de modelos não encontrada: %s, definida %s, não existe ou não é uma pasta.", e.Value, where)
	case core.DirIsHome:
		return fmt.Sprintf("pasta de modelos desativada: o valor %s é a pasta pessoal (%s).", where, e.Value)
	case core.DirNoHome:
		return "pasta de modelos não encontrada: XDG_TEMPLATES_DIR e HOME estão vazias."
	}
	return fmt.Sprintf("pasta de modelos não encontrada: %s.", e.Value)
}

func describeManifest(e *core.ManifestError) string {
	var what string
	switch e.Problem {
	case core.ManifestSyntax:
		what = fmt.Sprintf("erro de sintaxe (%v)", e.Err)
		if e.Line > 0 {
			what = fmt.Sprintf("erro de sintaxe na linha %d (%v)", e.Line, e.Err)
		}
	case core.ManifestUnknownKey:
		what = fmt.Sprintf("a chave %q não é aceita; a única tabela é [names]", e.Key)
	case core.ManifestNamesNotTable:
		what = "names precisa ser uma tabela, escrita como [names]"
	case core.ManifestNotString:
		what = fmt.Sprintf("o valor do padrão %q precisa ser um texto entre aspas", e.Key)
	case core.ManifestDoubleStar:
		what = fmt.Sprintf("o padrão %q usa **, que é reservado", e.Key)
	case core.ManifestInvalidName:
		var invalid *core.InvalidNameError
		if errors.As(e.Err, &invalid) {
			what = fmt.Sprintf("o nome fixo %q, do padrão %q, %s", invalid.Name, e.Key, nameProblem(invalid.Problem))
		} else {
			what = fmt.Sprintf("o nome fixo do padrão %q é inválido", e.Key)
		}
	case core.ManifestAmbiguous:
		what = fmt.Sprintf("os padrões %s casam o mesmo modelo, %q; acrescente uma regra exata para ele", quoteAll(e.Patterns), e.ID)
	}
	return fmt.Sprintf("manifesto inválido (%s): %s.", e.Path, what)
}

func describeUsage(e *failure.UsageError) string {
	var what string
	switch e.Problem {
	case failure.NoTerminal:
		return "o modo interativo precisa de um terminal. Fora dele, use janus list e janus new."
	case failure.UnknownCommand:
		what = fmt.Sprintf("subcomando desconhecido: %s", e.Arg)
	case failure.UnknownOption:
		what = fmt.Sprintf("opção desconhecida: %s", e.Arg)
	case failure.MissingValue:
		what = fmt.Sprintf("a opção %s precisa de um valor", e.Arg)
	case failure.RepeatedOption:
		what = fmt.Sprintf("a opção %s aparece mais de uma vez", e.Arg)
	case failure.UnexpectedArgument:
		what = fmt.Sprintf("argumento inesperado: %s", e.Arg)
	case failure.MissingArgument:
		what = fmt.Sprintf("falta o argumento %s", e.Arg)
	default:
		what = "uso incorreto"
	}
	help := "janus --help"
	if e.Command != "" {
		help = "janus " + e.Command + " --help"
	}
	return fmt.Sprintf("%s. Veja %s.", what, help)
}

func nameProblem(p core.NameProblem) string {
	switch p {
	case core.NameEmpty:
		return "está vazio"
	case core.NameDots:
		return "não pode ser . nem .."
	case core.NameSlash:
		return "não pode ter /"
	case core.NameControl:
		return "não pode ter caractere de controle"
	}
	return "é inválido"
}

// quoteAll cita os itens entre aspas e os junta à moda brasileira: "a", "b" e "c".
func quoteAll(items []string) string {
	quoted := make([]string, len(items))
	for i, s := range items {
		quoted[i] = strconv.Quote(s)
	}
	if len(quoted) < 2 {
		return strings.Join(quoted, "")
	}
	return strings.Join(quoted[:len(quoted)-1], ", ") + " e " + quoted[len(quoted)-1]
}

// formatSize escreve um tamanho do jeito brasileiro: "24,0 kB".
func formatSize(n int64) string {
	switch {
	case n == 1:
		return "1 byte"
	case n < 1024:
		return fmt.Sprintf("%d bytes", n)
	case n < 1024*1024:
		return decimal(float64(n)/1024) + " kB"
	}
	return decimal(float64(n)/(1024*1024)) + " MB"
}

func decimal(f float64) string {
	return strings.Replace(strconv.FormatFloat(f, 'f', 1, 64), ".", ",", 1)
}
