package cli

import (
	"errors"

	"gitlab.com/rg-dev/janus/internal/core"
	"gitlab.com/rg-dev/janus/internal/failure"
)

// Códigos de saída do contrato (docs/specs/bootstrap.md, "Códigos de saída").
const (
	codeOK             = 0
	codeUnexpected     = 1
	codeUsage          = 2
	codeTemplatesDir   = 3
	codeManifest       = 4
	codeNotFound       = 5
	codeNameRequired   = 6
	codeInvalidName    = 7
	codeExists         = 8
	codeMissingProgram = 9
	codeDestinationDir = 10
	codeCanceled       = 130
)

// exitCode traduz err no código de saída (FR-CLI-10). O manifesto vem antes do nome
// inválido porque um nome fixo inválido chega embrulhado num ManifestError.
func exitCode(err error) int {
	switch {
	case err == nil:
		return codeOK
	case errors.Is(err, failure.ErrCanceled):
		return codeCanceled
	case errors.As(err, new(*failure.UsageError)):
		return codeUsage
	case errors.As(err, new(*core.TemplatesDirError)):
		return codeTemplatesDir
	case errors.As(err, new(*core.ManifestError)):
		return codeManifest
	case errors.As(err, new(*core.TemplateNotFoundError)), errors.As(err, new(*failure.NoTemplatesError)):
		return codeNotFound
	case errors.As(err, new(*core.NameRequiredError)):
		return codeNameRequired
	case errors.As(err, new(*core.InvalidNameError)):
		return codeInvalidName
	case errors.As(err, new(*core.DestinationExistsError)):
		return codeExists
	case errors.As(err, new(*failure.MissingProgramError)):
		return codeMissingProgram
	case errors.As(err, new(*core.DestinationDirError)):
		return codeDestinationDir
	}
	return codeUnexpected
}
