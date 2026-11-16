package core

import (
	"errors"
	"slices"
	"testing"

	"gitlab.com/rg-dev/janus/internal/testfixture"
)

// loadManifest grava content como manifesto numa pasta nova e o lê.
func loadManifest(t *testing.T, content string) (Manifest, error) {
	t.Helper()
	dir := t.TempDir()
	testfixture.Write(t, dir, ManifestFile, content, 0o644)
	return LoadManifest(dir)
}

func TestLoadManifest(t *testing.T) {
	t.Run("has no rules when the file is missing (FR-CORE-10)", func(t *testing.T) {
		m, err := LoadManifest(t.TempDir())
		if err != nil || len(m.Rules) != 0 {
			t.Errorf("got %+v, %v", m, err)
		}
	})

	t.Run("has no rules without the names table (FR-CORE-10)", func(t *testing.T) {
		for _, content := range []string{"", "# só um comentário\n", "[names]\n"} {
			m, err := loadManifest(t, content)
			if err != nil || len(m.Rules) != 0 {
				t.Errorf("%q: got %+v, %v", content, m, err)
			}
		}
	})

	t.Run("reads the rules in file order", func(t *testing.T) {
		m, err := loadManifest(t, testfixture.Manifest)
		want := []Rule{
			{Pattern: "info.ai/ai-memory/*", Name: ".ai-memory.toml"},
			{Pattern: "info.ai/agents/*", Name: "AGENTS.md"},
			{Pattern: "info.dev/package.json.*", Name: "package.json"},
		}
		if err != nil || !slices.Equal(m.Rules, want) {
			t.Errorf("got %+v, %v", m.Rules, err)
		}
	})

	t.Run("trims the fixed name (FR-CORE-14)", func(t *testing.T) {
		m, err := loadManifest(t, "[names]\n\"a\" = \"  AGENTS.md \"\n")
		if err != nil || len(m.Rules) != 1 || m.Rules[0].Name != "AGENTS.md" {
			t.Errorf("got %+v, %v", m.Rules, err)
		}
	})

	cases := []struct {
		name    string
		content string
		problem ManifestProblem
		key     string
	}{
		{"fails on a key outside names (FR-CORE-11)", "versao = 1\n[names]\n\"a\" = \"b\"\n", ManifestUnknownKey, "versao"},
		{"fails on an unknown table (FR-CORE-11)", "[outra]\nx = \"y\"\n", ManifestUnknownKey, "outra"},
		{"fails when names is not a table (FR-CORE-11)", "names = \"x\"\n", ManifestNamesNotTable, "names"},
		{"fails when names is an array of tables (FR-CORE-11)", "[[names]]\na = \"b\"\n", ManifestNamesNotTable, "names"},
		{"fails on a subtable of names (FR-CORE-11)", "[names.x]\ny = \"z\"\n", ManifestNotString, "x"},
		{"fails on a numeric value (FR-CORE-11)", "[names]\n\"a\" = 1\n", ManifestNotString, "a"},
		{"fails on a double star at the end (FR-CORE-12)", "[names]\n\"a**\" = \"b\"\n", ManifestDoubleStar, "a**"},
		{"fails on a double star at the start (FR-CORE-12)", "[names]\n\"**/x\" = \"b\"\n", ManifestDoubleStar, "**/x"},
		{"fails on a double star in the middle (FR-CORE-12)", "[names]\n\"a/**/b\" = \"b\"\n", ManifestDoubleStar, "a/**/b"},
		{"fails on a fixed name with a slash (FR-CORE-14)", "[names]\n\"a\" = \"a/b\"\n", ManifestInvalidName, "a"},
		{"fails on a fixed name that is .. (FR-CORE-14)", "[names]\n\"a\" = \"..\"\n", ManifestInvalidName, "a"},
		{"fails on a blank fixed name (FR-CORE-14)", "[names]\n\"a\" = \"  \"\n", ManifestInvalidName, "a"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := loadManifest(t, c.content)
			var e *ManifestError
			if !errors.As(err, &e) || e.Problem != c.problem || e.Key != c.key {
				t.Errorf("got %v (%+v), want problem %d on %q", err, e, c.problem, c.key)
			}
		})
	}

	t.Run("reports the line of a syntax error (FR-CORE-11)", func(t *testing.T) {
		for content, line := range map[string]int{
			"[names]\n\"a\" = \n":                     2,
			"[names]\n\"a\" = \"x\"\n\"a\" = \"y\"\n": 3,
		} {
			_, err := loadManifest(t, content)
			var e *ManifestError
			if !errors.As(err, &e) || e.Problem != ManifestSyntax || e.Line != line {
				t.Errorf("%q: got %v (%+v), want syntax error on line %d", content, err, e, line)
			}
		}
	})
}
