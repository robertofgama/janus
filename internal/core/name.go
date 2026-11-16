package core

import (
	"fmt"
	"strings"
	"unicode"
)

// NameProblem diz por que um nome é inválido.
type NameProblem int

const (
	NameEmpty   NameProblem = iota + 1 // vazio depois de aparar as pontas
	NameDots                           // "." ou ".."
	NameSlash                          // contém "/"
	NameControl                        // contém caractere de controle
)

// InvalidNameError é um nome inválido (código 7).
type InvalidNameError struct {
	Name    string // o nome como foi informado
	Problem NameProblem
}

func (e *InvalidNameError) Error() string {
	return fmt.Sprintf("invalid name %q (problem %d)", e.Name, e.Problem)
}

// NameRequiredError é um modelo ask sem nome informado, fora do modo interativo (código 6).
type NameRequiredError struct {
	ID string
}

func (e *NameRequiredError) Error() string {
	return fmt.Sprintf("template %q needs a name", e.ID)
}

// cleanName apara os espaços das pontas e valida o nome (FR-CORE-16). Só o espaço é
// aparado: uma tabulação na ponta fica e torna o nome inválido, como caractere de controle.
func cleanName(name string) (string, error) {
	n := strings.Trim(name, " ")
	var problem NameProblem
	switch {
	case n == "":
		problem = NameEmpty
	case n == "." || n == "..":
		problem = NameDots
	case strings.Contains(n, "/"):
		problem = NameSlash
	case strings.ContainsFunc(n, unicode.IsControl):
		problem = NameControl
	default:
		return n, nil
	}
	return "", &InvalidNameError{Name: name, Problem: problem}
}

// extension devolve a extensão de name, com o ponto, ou "" (FR-CORE-17): o trecho a
// partir do último ponto, se o ponto não é o primeiro caractere, há algo depois dele e o
// trecho não tem espaço. Extensão composta perde a primeira parte (DT-08).
func extension(name string) string {
	i := strings.LastIndexByte(name, '.')
	if i <= 0 || i == len(name)-1 || strings.Contains(name[i:], " ") {
		return ""
	}
	return name[i:]
}

// FinalName devolve o nome com que o arquivo de t vai nascer. input é o nome informado
// pelo usuário e só vale com informed. Sem ele, vale o nome fixo, e um modelo ask falha
// com NameRequiredError (FR-CLI-08). Com ele, o nome é validado e ganha a extensão de
// referência quando ela falta, mesmo num modelo fixed (FR-CORE-16, FR-CORE-17).
func FinalName(t Template, input string, informed bool) (string, error) {
	if !informed {
		if t.Mode == NameFixed {
			return t.FixedName, nil
		}
		return "", &NameRequiredError{ID: t.ID}
	}
	name, err := cleanName(input)
	if err != nil {
		return "", err
	}
	ext := t.ReferenceExtension()
	if ext == "" || hasSuffixFold(name, ext) {
		return name, nil
	}
	return name + ext, nil
}

// hasSuffixFold diz se s termina com suffix, sem diferenciar maiúsculas.
func hasSuffixFold(s, suffix string) bool {
	return len(s) >= len(suffix) && strings.EqualFold(s[len(s)-len(suffix):], suffix)
}
