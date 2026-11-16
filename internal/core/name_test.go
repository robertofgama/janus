package core

import (
	"errors"
	"testing"
)

func TestFinalName(t *testing.T) {
	contestacao := Template{ID: "direito/01.2. Civ. Contestação.tex", File: "01.2. Civ. Contestação.tex", Mode: NameAsk}
	emptyBash := Template{ID: "Empty Bash", File: "Empty Bash", Mode: NameAsk}
	notas := Template{ID: "01.0. Notas", File: "01.0. Notas", Mode: NameAsk}
	html := Template{ID: "info.dev/react; script; simple.html", File: "react; script; simple.html", Mode: NameAsk}
	agents := Template{ID: "info.ai/agents/react.AGENTS.md", File: "react.AGENTS.md", Mode: NameFixed, FixedName: "AGENTS.md"}
	pkg := Template{ID: "info.dev/package.json.monorepo", File: "package.json.monorepo", Mode: NameFixed, FixedName: "package.json"}
	memory := Template{ID: "info.ai/ai-memory/default.ai-memory.toml", File: "default.ai-memory.toml", Mode: NameFixed, FixedName: ".ai-memory.toml"}
	bashrc := Template{ID: "bashrc", File: "bashrc", Mode: NameFixed, FixedName: ".bashrc"}

	t.Run("adds the reference extension when it is missing (FR-CORE-17)", func(t *testing.T) {
		cases := []struct {
			template Template
			input    string
			want     string
		}{
			{contestacao, "Contestação Fulano", "Contestação Fulano.tex"},
			{contestacao, "Fulano.TEX", "Fulano.TEX"},
			{contestacao, "  Fulano  ", "Fulano.tex"},
			{emptyBash, "deploy", "deploy"},
			{agents, "CLAUDE", "CLAUDE.md"},
			{pkg, "pacote", "pacote.json"},
			{memory, "outro", "outro.toml"},
			{notas, "resumo", "resumo"},
			{html, "pagina", "pagina.html"},
			{bashrc, "zshrc", "zshrc"},
		}
		for _, c := range cases {
			got, err := FinalName(c.template, c.input, true)
			if err != nil || got != c.want {
				t.Errorf("FinalName(%q, %q) = %q, %v; want %q", c.template.File, c.input, got, err, c.want)
			}
		}
	})

	t.Run("uses the fixed name when no name is informed (FR-CORE-13)", func(t *testing.T) {
		if got, err := FinalName(agents, "", false); err != nil || got != "AGENTS.md" {
			t.Errorf("got %q, %v", got, err)
		}
	})

	t.Run("requires a name for an ask template (FR-CLI-08)", func(t *testing.T) {
		_, err := FinalName(contestacao, "", false)
		var required *NameRequiredError
		if !errors.As(err, &required) || required.ID != contestacao.ID {
			t.Errorf("got %v, want NameRequiredError for %q", err, contestacao.ID)
		}
	})

	t.Run("rejects invalid names (FR-CORE-16)", func(t *testing.T) {
		cases := []struct {
			input   string
			problem NameProblem
		}{
			{"a/b", NameSlash},
			{"..", NameDots},
			{".", NameDots},
			{"  ", NameEmpty},
			{"", NameEmpty},
			{"a\tb", NameControl},
			{"a\nb", NameControl},
			{"\tx", NameControl},
		}
		for _, c := range cases {
			_, err := FinalName(contestacao, c.input, true)
			var invalid *InvalidNameError
			if !errors.As(err, &invalid) || invalid.Problem != c.problem || invalid.Name != c.input {
				t.Errorf("FinalName(%q) = %v; want problem %d", c.input, err, c.problem)
			}
		}
	})
}

func TestExtension(t *testing.T) {
	t.Run("follows the extension rule (FR-CORE-17)", func(t *testing.T) {
		cases := map[string]string{
			"01.2. Civ. Contestação.tex": ".tex",
			"Empty Bash":                 "",
			"AGENTS.md":                  ".md",
			"package.json":               ".json",
			".ai-memory.toml":            ".toml",
			".bashrc":                    "",
			"01.0. Notas":                "",
			"nome.":                      "",
			"a.b c":                      "",
			"x.tar.gz":                   ".gz",
		}
		for name, want := range cases {
			if got := extension(name); got != want {
				t.Errorf("extension(%q) = %q, want %q", name, got, want)
			}
		}
	})
}
