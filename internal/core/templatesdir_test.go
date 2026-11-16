package core

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/rg-dev/janus/internal/testfixture"
)

// environment é um ambiente de teste; get faz o papel de os.Getenv sem tocar no processo.
type environment map[string]string

func (e environment) get(key string) string { return e[key] }

// homeWithUserDirs cria uma pasta pessoal com a pasta Modelos e um user-dirs.dirs com
// content em XDG_CONFIG_HOME. "<HOME>" em content vira o caminho da pasta pessoal.
func homeWithUserDirs(t *testing.T, content string) (environment, string) {
	t.Helper()
	home := t.TempDir()
	if err := os.Mkdir(filepath.Join(home, "Modelos"), 0o755); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(home, "config")
	testfixture.Write(t, config, "user-dirs.dirs", strings.ReplaceAll(content, "<HOME>", home), 0o644)
	return environment{"HOME": home, "XDG_CONFIG_HOME": config}, home
}

// wantDirError confere que err é um TemplatesDirError com a origem e o motivo esperados.
func wantDirError(t *testing.T, err error, source TemplatesDirSource, reason TemplatesDirReason) *TemplatesDirError {
	t.Helper()
	var e *TemplatesDirError
	if !errors.As(err, &e) {
		t.Fatalf("got %v, want *TemplatesDirError", err)
	}
	if e.Source != source || e.Reason != reason {
		t.Fatalf("got source %d reason %d, want source %d reason %d", e.Source, e.Reason, source, reason)
	}
	return e
}

func TestFindTemplatesDir(t *testing.T) {
	t.Run("uses XDG_TEMPLATES_DIR first (FR-CORE-01)", func(t *testing.T) {
		env, _ := homeWithUserDirs(t, `XDG_TEMPLATES_DIR="/nao/existe"`+"\n")
		templates := t.TempDir()
		env["XDG_TEMPLATES_DIR"] = templates
		if got, err := FindTemplatesDir(env.get); err != nil || got != templates {
			t.Fatalf("got %q, %v; want %q", got, err, templates)
		}
	})

	t.Run("does not fall back when XDG_TEMPLATES_DIR is bad (FR-CORE-01)", func(t *testing.T) {
		env, _ := homeWithUserDirs(t, `XDG_TEMPLATES_DIR="$HOME/Modelos"`+"\n")
		env["XDG_TEMPLATES_DIR"] = "/caminho/inexistente"
		_, err := FindTemplatesDir(env.get)
		if e := wantDirError(t, err, SourceEnv, DirNotDir); e.Value != "/caminho/inexistente" {
			t.Errorf("value %q", e.Value)
		}
	})

	t.Run("reads user-dirs.dirs like GLib (FR-CORE-02)", func(t *testing.T) {
		cases := []struct {
			name    string
			content string
			want    string
		}{
			{"with the home prefix", `XDG_TEMPLATES_DIR="$HOME/Modelos"` + "\n", "Modelos"},
			{"with spaces around the equal sign", "  \tXDG_TEMPLATES_DIR \t= \"$HOME/Modelos\"\n", "Modelos"},
			{"with the last line winning", "XDG_TEMPLATES_DIR=\"$HOME/Outra\"\nXDG_TEMPLATES_DIR=\"$HOME/Modelos\"\n", "Modelos"},
			{"with trailing slashes", `XDG_TEMPLATES_DIR="$HOME/Modelos///"` + "\n", "Modelos"},
			{"with other lines and comments", "# comentário\nXDG_DESKTOP_DIR=\"$HOME/Desktop\"\nXDG_TEMPLATES_DIR=\"$HOME/Modelos\"\n", "Modelos"},
			{"with escapes kept as they are", `XDG_TEMPLATES_DIR="$HOME/Mo\"delos"` + "\n", `Mo\"delos`},
			{"with an absolute path", `XDG_TEMPLATES_DIR="<HOME>/Modelos"` + "\n", "Modelos"},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				env, home := homeWithUserDirs(t, c.content)
				want := filepath.Join(home, c.want)
				if err := os.MkdirAll(want, 0o755); err != nil {
					t.Fatal(err)
				}
				if got, err := FindTemplatesDir(env.get); err != nil || got != want {
					t.Errorf("got %q, %v; want %q", got, err, want)
				}
			})
		}
	})

	t.Run("falls back to ~/.config without XDG_CONFIG_HOME (FR-CORE-02)", func(t *testing.T) {
		home := t.TempDir()
		want := filepath.Join(home, "Modelos")
		if err := os.Mkdir(want, 0o755); err != nil {
			t.Fatal(err)
		}
		testfixture.Write(t, filepath.Join(home, ".config"), "user-dirs.dirs", `XDG_TEMPLATES_DIR="$HOME/Modelos"`+"\n", 0o644)
		if got, err := FindTemplatesDir(environment{"HOME": home}.get); err != nil || got != want {
			t.Errorf("got %q, %v; want %q", got, err, want)
		}
	})

	t.Run("explains each failure of user-dirs.dirs (FR-CORE-03)", func(t *testing.T) {
		cases := []struct {
			name    string
			content string
			reason  TemplatesDirReason
		}{
			{"without the templates line", `XDG_DESKTOP_DIR="$HOME/Desktop"` + "\n", DirLineMissing},
			{"with a relative value", `XDG_TEMPLATES_DIR="Modelos"` + "\n", DirRelative},
			{"with home followed by a name", `XDG_TEMPLATES_DIR="$HOMEx"` + "\n", DirRelative},
			{"with a missing dir", `XDG_TEMPLATES_DIR="$HOME/nao-existe"` + "\n", DirNotDir},
			{"with home alone", `XDG_TEMPLATES_DIR="$HOME"` + "\n", DirIsHome},
			{"with home and a slash", `XDG_TEMPLATES_DIR="$HOME/"` + "\n", DirIsHome},
			{"with home as an absolute path", `XDG_TEMPLATES_DIR="<HOME>"` + "\n", DirIsHome},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				env, _ := homeWithUserDirs(t, c.content)
				_, err := FindTemplatesDir(env.get)
				e := wantDirError(t, err, SourceUserDirs, c.reason)
				if want := filepath.Join(env["XDG_CONFIG_HOME"], "user-dirs.dirs"); e.File != want {
					t.Errorf("file %q, want %q", e.File, want)
				}
			})
		}
	})

	t.Run("fails when user-dirs.dirs is missing (FR-CORE-03)", func(t *testing.T) {
		home := t.TempDir()
		_, err := FindTemplatesDir(environment{"HOME": home}.get)
		e := wantDirError(t, err, SourceUserDirs, DirFileMissing)
		if want := filepath.Join(home, ".config", "user-dirs.dirs"); e.File != want {
			t.Errorf("file %q, want %q", e.File, want)
		}
	})

	t.Run("fails when XDG_TEMPLATES_DIR is a file or relative (FR-CORE-03)", func(t *testing.T) {
		file := testfixture.Write(t, t.TempDir(), "arquivo", "", 0o644)
		_, err := FindTemplatesDir(environment{"XDG_TEMPLATES_DIR": file}.get)
		wantDirError(t, err, SourceEnv, DirNotDir)
		_, err = FindTemplatesDir(environment{"XDG_TEMPLATES_DIR": "Modelos"}.get)
		wantDirError(t, err, SourceEnv, DirRelative)
	})

	t.Run("treats the home dir as disabled (FR-CORE-03)", func(t *testing.T) {
		home := t.TempDir()
		for _, value := range []string{home, home + "/"} {
			_, err := FindTemplatesDir(environment{"HOME": home, "XDG_TEMPLATES_DIR": value}.get)
			wantDirError(t, err, SourceEnv, DirIsHome)
		}
	})

	t.Run("fails without HOME (FR-CORE-03)", func(t *testing.T) {
		_, err := FindTemplatesDir(environment{}.get)
		wantDirError(t, err, SourceUserDirs, DirNoHome)
	})

	t.Run("follows a link to the templates dir (FR-CORE-04)", func(t *testing.T) {
		target := t.TempDir()
		link := filepath.Join(t.TempDir(), "Modelos")
		if err := os.Symlink(target, link); err != nil {
			t.Fatal(err)
		}
		if got, err := FindTemplatesDir(environment{"XDG_TEMPLATES_DIR": link}.get); err != nil || got != link {
			t.Errorf("got %q, %v; want %q", got, err, link)
		}
	})
}
