package interactive

import (
	"bytes"
	"errors"
	"slices"
	"strings"
	"testing"

	"gitlab.com/rg-dev/janus/internal/failure"
	"gitlab.com/rg-dev/janus/internal/testfixture"
)

// optionValue devolve o argumento que vem logo depois de option, ou "".
func optionValue(args []string, option string) string {
	i := slices.Index(args, option)
	if i < 0 || i+1 == len(args) {
		return ""
	}
	return args[i+1]
}

// guardOptions são as opções que o seletor e o campo de nome precisam receber para anular
// as do usuário que mudariam a saída ou fariam o fzf sair sem perguntar (NFR-TUI-02).
var guardOptions = []string{"--no-multi", "--no-expect", "--no-exit-0", "--no-select-1", "--no-read0", "--no-print0"}

func TestFzfSelector(t *testing.T) {
	t.Run("sends one id per line and returns the choice (FR-TUI-10)", func(t *testing.T) {
		fake := testfixture.FakeProgram(t, "fzf", testfixture.Reply{Stdout: "b c\n"})
		id, err := FzfSelector{Executable: "/opt/janus", Stderr: &bytes.Buffer{}}.Select([]string{"a", "b c"})
		if err != nil || id != "b c" {
			t.Fatalf("got %q, %v", id, err)
		}
		if calls := fake.Calls(t); len(calls) != 1 || calls[0].Stdin != "a\nb c\n" {
			t.Errorf("got calls %+v", calls)
		}
	})

	t.Run("returns canceled on codes 1 and 130 (FR-TUI-10)", func(t *testing.T) {
		for _, code := range []int{1, 130} {
			testfixture.FakeProgram(t, "fzf", testfixture.Reply{Code: code})
			_, err := FzfSelector{Executable: "/opt/janus", Stderr: &bytes.Buffer{}}.Select([]string{"a"})
			if !errors.Is(err, failure.ErrCanceled) {
				t.Errorf("code %d: got %v", code, err)
			}
		}
	})

	t.Run("fails on other codes and passes the fzf message on (FR-TUI-10)", func(t *testing.T) {
		testfixture.FakeProgram(t, "fzf", testfixture.Reply{Stderr: "opção inválida\n", Code: 2})
		var stderr bytes.Buffer
		_, err := FzfSelector{Executable: "/opt/janus", Stderr: &stderr}.Select([]string{"a"})
		var e *failure.ProgramError
		if !errors.As(err, &e) || e.Program != "fzf" || e.Code != 2 {
			t.Errorf("got %v", err)
		}
		if stderr.String() != "opção inválida\n" {
			t.Errorf("stderr %q", stderr.String())
		}
	})

	t.Run("overrides the user options that change the output (NFR-TUI-02)", func(t *testing.T) {
		fake := testfixture.FakeProgram(t, "fzf", testfixture.Reply{Stdout: "a\n"})
		if _, err := (FzfSelector{Executable: "/opt/janus", Stderr: &bytes.Buffer{}}).Select([]string{"a"}); err != nil {
			t.Fatal(err)
		}
		args := fake.Calls(t)[0].Args
		for _, option := range append([]string{"--no-print-query"}, guardOptions...) {
			if !slices.Contains(args, option) {
				t.Errorf("missing %s in %q", option, args)
			}
		}
	})

	t.Run("runs the preview with sh and the quoted binary (NFR-TUI-01)", func(t *testing.T) {
		fake := testfixture.FakeProgram(t, "fzf", testfixture.Reply{Stdout: "a\n"})
		if _, err := (FzfSelector{Executable: "/opt/o'k dir/janus", Stderr: &bytes.Buffer{}}).Select([]string{"a"}); err != nil {
			t.Fatal(err)
		}
		args := fake.Calls(t)[0].Args
		if got := optionValue(args, "--with-shell"); got != "sh -c" {
			t.Errorf("--with-shell %q", got)
		}
		if got, want := optionValue(args, "--preview"), `'/opt/o'\''k dir/janus' __preview {}`; got != want {
			t.Errorf("--preview %q, want %q", got, want)
		}
	})
}

func TestFzfPrompt(t *testing.T) {
	t.Run("returns the typed name (FR-TUI-11)", func(t *testing.T) {
		for _, code := range []int{0, 1} {
			testfixture.FakeProgram(t, "fzf", testfixture.Reply{Stdout: "Fulano ção\n", Code: code})
			name, err := FzfPrompt{Stderr: &bytes.Buffer{}}.Ask("inicial", "")
			if err != nil || name != "Fulano ção" {
				t.Errorf("code %d: got %q, %v", code, name, err)
			}
		}
	})

	t.Run("returns canceled on code 130 (FR-TUI-11)", func(t *testing.T) {
		testfixture.FakeProgram(t, "fzf", testfixture.Reply{Code: 130})
		if _, err := (FzfPrompt{Stderr: &bytes.Buffer{}}).Ask("inicial", ""); !errors.Is(err, failure.ErrCanceled) {
			t.Errorf("got %v", err)
		}
	})

	t.Run("passes the initial value, the warning and the cancel key (FR-TUI-11)", func(t *testing.T) {
		fake := testfixture.FakeProgram(t, "fzf", testfixture.Reply{Stdout: "x\n", Code: 1})
		if _, err := (FzfPrompt{Stderr: &bytes.Buffer{}}).Ask("-inicial", "já existe"); err != nil {
			t.Fatal(err)
		}
		call := fake.Calls(t)[0]
		for _, option := range []string{"--print-query", "--query=-inicial", "--header=já existe", "--bind=ctrl-d:abort"} {
			if !slices.Contains(call.Args, option) {
				t.Errorf("missing %s in %q", option, call.Args)
			}
		}
		if call.Stdin != "" {
			t.Errorf("the list is %q, want empty", call.Stdin)
		}
	})

	// Regressão: com --exit-0 no FZF_DEFAULT_OPTS, o fzf com a lista vazia sai na hora com o
	// valor inicial. O arquivo nascia sem confirmação, e um destino existente reabria o
	// campo sem fim.
	t.Run("overrides the user options that change the output (NFR-TUI-02)", func(t *testing.T) {
		fake := testfixture.FakeProgram(t, "fzf", testfixture.Reply{Stdout: "x\n", Code: 1})
		if _, err := (FzfPrompt{Stderr: &bytes.Buffer{}}).Ask("inicial", ""); err != nil {
			t.Fatal(err)
		}
		args := fake.Calls(t)[0].Args
		for _, option := range guardOptions {
			if !slices.Contains(args, option) {
				t.Errorf("missing %s in %q", option, args)
			}
		}
	})

	t.Run("omits the header without a warning (FR-TUI-11)", func(t *testing.T) {
		fake := testfixture.FakeProgram(t, "fzf", testfixture.Reply{Stdout: "x\n", Code: 1})
		if _, err := (FzfPrompt{Stderr: &bytes.Buffer{}}).Ask("inicial", ""); err != nil {
			t.Fatal(err)
		}
		for _, arg := range fake.Calls(t)[0].Args {
			if strings.HasPrefix(arg, "--header") {
				t.Errorf("got %s", arg)
			}
		}
	})
}
