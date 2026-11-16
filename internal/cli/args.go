// Package cli é a linha de comando: lê os argumentos, roda list, new e o modo interativo e
// traduz os erros em códigos de saída. É o contrato de que os adaptadores dependem.
package cli

import (
	"strings"

	"gitlab.com/rg-dev/janus/internal/failure"
)

// parsed é o resultado da leitura dos argumentos de um nível da linha de comando.
type parsed struct {
	options map[string]string // opção, sem os hífens, e o valor; "" nas que não recebem valor
	args    []string          // os argumentos posicionais, na ordem
}

func (p parsed) has(option string) bool {
	_, ok := p.options[option]
	return ok
}

// parseArgs lê args conforme options, que diz quais opções existem e se recebem valor
// (DE-19). Aceita "--opcao valor" e "--opcao=valor" em qualquer posição; "--" encerra as
// opções. Opção desconhecida, com um hífen só, repetida ou sem valor é uso incorreto
// (FR-CLI-01, FR-CLI-02). O pacote flag foi rejeitado porque para no primeiro argumento
// posicional, e o contrato põe as opções depois do identificador.
func parseArgs(command string, args []string, options map[string]bool) (parsed, error) {
	usage := func(problem failure.UsageProblem, arg string) (parsed, error) {
		return parsed{}, &failure.UsageError{Problem: problem, Arg: arg, Command: command}
	}
	p := parsed{options: map[string]string{}}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--":
			p.args = append(p.args, args[i+1:]...)
			return p, nil
		case strings.HasPrefix(arg, "--"):
			name, value, hasValue := strings.Cut(arg[2:], "=")
			takesValue, known := options[name]
			if !known || (hasValue && !takesValue) {
				return usage(failure.UnknownOption, arg)
			}
			if p.has(name) {
				return usage(failure.RepeatedOption, "--"+name)
			}
			if takesValue && !hasValue {
				if i+1 == len(args) {
					return usage(failure.MissingValue, arg)
				}
				i++
				value = args[i]
			}
			p.options[name] = value
		case strings.HasPrefix(arg, "-"):
			return usage(failure.UnknownOption, arg)
		default:
			p.args = append(p.args, arg)
		}
	}
	return p, nil
}
