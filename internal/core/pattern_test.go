package core

import (
	"fmt"
	"testing"
)

func TestMatchPattern(t *testing.T) {
	cases := []struct {
		pattern string
		id      string
		want    bool
	}{
		{"info.ai/agents/*", "info.ai/agents/AGENTS.md", true},
		{"info.ai/agents/*", "info.ai/agents/x/y", false},
		{"*", "Empty Bash", true},
		{"*", "a/b", false},
		{"a*", "a", true},
		{"info.dev/package.json.*", "info.dev/package.json.monorepo", true},
		{"info.ai/*/react.*", "info.ai/agents/react.AGENTS.md", true},
		{"direito/*.tex", "direito/01.2. Civ. Contestação.tex", true},
		{"*.tex", "direito/x.tex", false},
		{"a*b*c", "axxbyyc", true},
		{"a*b*c", "axxbyy", false},
		{"a?c", "abc", false},
		{"a?c", "a?c", true},
		{"[ab]", "a", false},
		{"[ab]", "[ab]", true},
		{`a\*`, `a\x`, true},
		{"exato", "exato", true},
		{"exato", "exatos", false},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%q against %q (FR-CORE-12)", c.pattern, c.id), func(t *testing.T) {
			if got := matchPattern(c.pattern, c.id); got != c.want {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}
