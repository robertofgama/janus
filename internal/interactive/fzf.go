package interactive

import (
	"errors"
	"io"
	"os/exec"
	"strings"

	"gitlab.com/rg-dev/janus/internal/failure"
	"gitlab.com/rg-dev/janus/internal/messages"
)

// FzfSelector é o seletor do terminal: o fzf, com a prévia do próprio janus (DE-06, DE-12).
type FzfSelector struct {
	Executable string    // o próprio janus, que a prévia chama
	Stderr     io.Writer // o terminal; recebe as mensagens de erro do fzf
}

// Select mostra ids no fzf e devolve o escolhido. Cancelar, ou confirmar sem nenhum
// resultado, devolve failure.ErrCanceled (FR-TUI-10).
func (s FzfSelector) Select(ids []string) (string, error) {
	out, code, err := runFzf(s.Stderr, strings.Join(ids, "\n")+"\n", selectorArgs(s.Executable))
	if err != nil {
		return "", err
	}
	id := firstLine(out)
	if code != 0 || id == "" {
		return "", failure.ErrCanceled
	}
	return id, nil
}

// selectorArgs monta as opções do seletor. Elas vêm depois das do FZF_DEFAULT_OPTS e
// prevalecem, então o tema do usuário fica e o contrato não muda.
func selectorArgs(executable string) []string {
	return []string{
		// Anulam opções do usuário que mudariam a saída (NFR-TUI-02).
		"--no-multi", "--no-print-query", "--no-expect",
		// A prévia roda com sh, e não com o $SHELL do usuário: o fish tem outras regras de
		// aspas, e o {} do fzf sai citado para sh (NFR-TUI-01). O identificador nunca é
		// juntado ao caminho aqui: quem o transforma em arquivo é o subcomando (DE-12).
		"--with-shell", "sh -c",
		"--preview", shellQuote(executable) + " " + PreviewCommand + " {}",
		"--prompt=" + messages.SelectorLabel,
	}
}

// FzfPrompt é o campo de nome: o fzf com a lista vazia (DE-17).
type FzfPrompt struct {
	Stderr io.Writer // o terminal; recebe as mensagens de erro do fzf
}

// Ask abre o campo com initial e, quando há, o aviso no cabeçalho. Devolve o texto
// confirmado; Esc, Ctrl-C e Ctrl-D devolvem failure.ErrCanceled (FR-TUI-11).
func (p FzfPrompt) Ask(initial, warning string) (string, error) {
	out, code, err := runFzf(p.Stderr, "", promptArgs(initial, warning))
	if err != nil {
		return "", err
	}
	if code == 130 {
		return "", failure.ErrCanceled
	}
	// Com a lista vazia, Enter sai com 1 (nenhum resultado); --print-query põe o texto na
	// primeira linha mesmo assim.
	return firstLine(out), nil
}

func promptArgs(initial, warning string) []string {
	args := []string{
		"--no-multi", "--no-expect", "--print-query",
		"--query=" + initial,
		"--prompt=" + messages.NameLabel,
		// Na configuração padrão, Ctrl-D só cancela com o campo vazio; a spec pede que
		// cancele sempre (DE-17).
		"--bind=ctrl-d:abort",
		"--info=hidden",
		"--height=~6",
	}
	if warning != "" {
		args = append(args, "--header="+warning)
	}
	return args
}

// runFzf roda o fzf com input na entrada e devolve a saída padrão e o código de saída. O
// fzf desenha no /dev/tty, então a saída de erro dele vai direto para o terminal, com as
// mensagens de erro que ele mesmo escreve. Os códigos 0, 1 e 130 são respostas; qualquer
// outro é falha (FR-TUI-10).
func runFzf(stderr io.Writer, input string, args []string) (string, int, error) {
	cmd := exec.Command("fzf", args...)
	cmd.Stdin = strings.NewReader(input)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = stderr
	err := cmd.Run()
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		code := exit.ExitCode()
		if code == 1 || code == 130 {
			return out.String(), code, nil
		}
		return "", 0, &failure.ProgramError{Program: "fzf", Code: code}
	}
	return out.String(), 0, err
}

func firstLine(s string) string {
	line, _, _ := strings.Cut(s, "\n")
	return line
}

// shellQuote cita s entre aspas simples, para sh.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
