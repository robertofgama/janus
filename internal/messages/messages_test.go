package messages

import (
	"errors"
	"strings"
	"testing"

	"gitlab.com/rg-dev/janus/internal/core"
	"gitlab.com/rg-dev/janus/internal/failure"
)

func TestDescribe(t *testing.T) {
	t.Run("cites the item of each error (FR-CLI-12)", func(t *testing.T) {
		cases := []struct {
			name  string
			err   error
			items []string
		}{
			{"templates dir from the variable", &core.TemplatesDirError{Source: core.SourceEnv, Reason: core.DirNotDir, Value: "/x/Modelos"}, []string{"/x/Modelos", "XDG_TEMPLATES_DIR"}},
			{"missing user-dirs.dirs", &core.TemplatesDirError{Source: core.SourceUserDirs, Reason: core.DirFileMissing, File: "/h/.config/user-dirs.dirs"}, []string{"/h/.config/user-dirs.dirs"}},
			{"missing line", &core.TemplatesDirError{Source: core.SourceUserDirs, Reason: core.DirLineMissing, File: "/h/u"}, []string{"/h/u", "XDG_TEMPLATES_DIR"}},
			{"relative value", &core.TemplatesDirError{Source: core.SourceUserDirs, Reason: core.DirRelative, File: "/h/u", Value: "Modelos"}, []string{"Modelos", "/h/u"}},
			{"home dir", &core.TemplatesDirError{Source: core.SourceEnv, Reason: core.DirIsHome, Value: "/home/x"}, []string{"/home/x"}},
			{"no home", &core.TemplatesDirError{Source: core.SourceUserDirs, Reason: core.DirNoHome}, []string{"HOME"}},
			{"manifest syntax", &core.ManifestError{Path: "/m/.janus.toml", Problem: core.ManifestSyntax, Line: 2, Err: errors.New("toml: x")}, []string{"/m/.janus.toml", "2"}},
			{"unknown key", &core.ManifestError{Path: "/m", Problem: core.ManifestUnknownKey, Key: "versao"}, []string{"versao"}},
			{"names not a table", &core.ManifestError{Path: "/m", Problem: core.ManifestNamesNotTable, Key: "names"}, []string{"names"}},
			{"value not a string", &core.ManifestError{Path: "/m", Problem: core.ManifestNotString, Key: "a/*"}, []string{"a/*"}},
			{"double star", &core.ManifestError{Path: "/m", Problem: core.ManifestDoubleStar, Key: "a/**"}, []string{"a/**"}},
			{"invalid fixed name", &core.ManifestError{Path: "/m", Problem: core.ManifestInvalidName, Key: "a/*", Err: &core.InvalidNameError{Name: "a/b", Problem: core.NameSlash}}, []string{"a/*", "a/b"}},
			{"ambiguous patterns", &core.ManifestError{Path: "/m", Problem: core.ManifestAmbiguous, ID: "a/b", Patterns: []string{"a/*", "*/b"}}, []string{"a/*", "*/b", "a/b"}},
			{"template not found", &core.TemplateNotFoundError{ID: "x/y"}, []string{"x/y"}},
			{"destination dir", &core.DestinationDirError{Dir: "/nao", Reason: core.DestNotDir}, []string{"/nao"}},
			{"destination dir inside the templates dir", &core.DestinationDirError{Dir: "sub", Reason: core.DestInsideTemplates, TemplatesDir: "/x/Modelos"}, []string{"sub", "/x/Modelos"}},
			{"name required", &core.NameRequiredError{ID: "x/y"}, []string{"x/y", "--name"}},
			{"invalid name", &core.InvalidNameError{Name: "a/b", Problem: core.NameSlash}, []string{"a/b"}},
			{"destination exists", &core.DestinationExistsError{Path: "/d/AGENTS.md"}, []string{"/d/AGENTS.md"}},
			{"unknown command", &failure.UsageError{Problem: failure.UnknownCommand, Arg: "foo"}, []string{"foo", "janus --help"}},
			{"unknown option of new", &failure.UsageError{Problem: failure.UnknownOption, Arg: "--x", Command: "new"}, []string{"--x", "janus new --help"}},
			{"no terminal", &failure.UsageError{Problem: failure.NoTerminal}, []string{"janus new"}},
			{"missing program", &failure.MissingProgramError{Program: "fzf"}, []string{"fzf"}},
			{"no templates", &failure.NoTemplatesError{Dir: "/x/Modelos"}, []string{"/x/Modelos"}},
			{"program failure", &failure.ProgramError{Program: "fzf", Code: 2}, []string{"fzf", "2"}},
			{"unexpected error", errors.New("boom"), []string{"boom"}},
		}
		for _, c := range cases {
			got := Describe(c.err)
			for _, item := range c.items {
				if !strings.Contains(got, item) {
					t.Errorf("%s: %q does not cite %q", c.name, got, item)
				}
			}
		}
	})

	t.Run("prefixes the error line with the program name", func(t *testing.T) {
		if got := ErrorLine(errors.New("x")); !strings.HasPrefix(got, "janus: ") {
			t.Errorf("got %q", got)
		}
	})
}

func TestPreviewBinary(t *testing.T) {
	t.Run("writes the size in Brazilian style (FR-TUI-12)", func(t *testing.T) {
		for size, want := range map[int64]string{1: "1 byte", 500: "500 bytes", 24576: "24,0 kB", 3 * 1024 * 1024: "3,0 MB"} {
			if got := PreviewBinary(size); !strings.Contains(got, want) {
				t.Errorf("PreviewBinary(%d) = %q, want %q in it", size, got, want)
			}
		}
	})
}
