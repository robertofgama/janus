// Package failure reúne os erros do programa que não vêm do núcleo: uso incorreto,
// programa ausente, pasta sem modelos no modo interativo, falha do fzf e cancelamento.
// Ficam num pacote próprio para que internal/messages (o texto) e internal/cli (o código
// de saída) os conheçam sem que um importe o outro (DE-21).
package failure

import (
	"errors"
	"fmt"
)

// ErrCanceled é o cancelamento pelo usuário, no seletor ou no campo de nome (código 130).
var ErrCanceled = errors.New("canceled by the user")

// UsageProblem diz o que está errado na linha de comando.
type UsageProblem int

const (
	UnknownCommand     UsageProblem = iota + 1 // subcomando que não existe
	UnknownOption                              // opção que não existe ali, ou com um hífen só
	MissingValue                               // opção sem o valor
	RepeatedOption                             // opção repetida
	UnexpectedArgument                         // argumento a mais
	MissingArgument                            // falta um argumento
	NoTerminal                                 // modo interativo sem terminal
)

// UsageError é o uso incorreto (código 2).
type UsageError struct {
	Problem UsageProblem
	Arg     string // o argumento ou a opção envolvida, quando houver
	Command string // o subcomando em uso; "" no nível de cima
}

func (e *UsageError) Error() string {
	return fmt.Sprintf("usage error (problem %d) at %q", e.Problem, e.Arg)
}

// MissingProgramError é um programa externo obrigatório fora do PATH (código 9).
type MissingProgramError struct {
	Program string
}

func (e *MissingProgramError) Error() string {
	return fmt.Sprintf("program %q not found in PATH", e.Program)
}

// NoTemplatesError é a pasta de modelos sem nenhum modelo, no modo interativo (código 5).
type NoTemplatesError struct {
	Dir string
}

func (e *NoTemplatesError) Error() string {
	return fmt.Sprintf("no templates in %s", e.Dir)
}

// ProgramError é um programa externo que terminou com um código inesperado (código 1).
type ProgramError struct {
	Program string
	Code    int
}

func (e *ProgramError) Error() string {
	return fmt.Sprintf("%s exited with code %d", e.Program, e.Code)
}
