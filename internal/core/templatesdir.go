package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TemplatesDirSource diz de onde veio a pasta de modelos.
type TemplatesDirSource int

const (
	SourceEnv      TemplatesDirSource = iota + 1 // a variável XDG_TEMPLATES_DIR
	SourceUserDirs                               // o arquivo user-dirs.dirs
)

// TemplatesDirReason diz por que a pasta de modelos não foi encontrada.
type TemplatesDirReason int

const (
	DirFileMissing TemplatesDirReason = iota + 1 // user-dirs.dirs ausente ou ilegível
	DirLineMissing                               // user-dirs.dirs sem XDG_TEMPLATES_DIR
	DirRelative                                  // valor que não é absoluto nem começa por $HOME
	DirNotDir                                    // caminho inexistente ou que não é pasta
	DirIsHome                                    // a pasta pessoal: pasta de modelos desativada
	DirNoHome                                    // HOME vazia, sem XDG_TEMPLATES_DIR
)

// TemplatesDirError é a pasta de modelos não encontrada (código 3).
type TemplatesDirError struct {
	Source TemplatesDirSource
	Reason TemplatesDirReason
	File   string // o user-dirs.dirs, quando Source é SourceUserDirs
	Value  string // o valor ou o caminho tentado
	Err    error  // a causa, quando houver
}

func (e *TemplatesDirError) Error() string {
	return fmt.Sprintf("templates dir not found (source %d, reason %d): %q", e.Source, e.Reason, e.Value)
}

func (e *TemplatesDirError) Unwrap() error { return e.Err }

// FindTemplatesDir descobre a pasta de modelos (FR-CORE-01 a FR-CORE-04). getenv lê o
// ambiente; os testes passam um mapa, sem tocar no ambiente do processo.
func FindTemplatesDir(getenv func(string) string) (string, error) {
	home := getenv("HOME")
	if value := getenv("XDG_TEMPLATES_DIR"); value != "" {
		// Sem segunda tentativa: um valor ruim falha, mesmo com um user-dirs.dirs válido.
		return checkTemplatesDir(value, home, &TemplatesDirError{Source: SourceEnv})
	}
	fail := &TemplatesDirError{Source: SourceUserDirs}
	if home == "" {
		fail.Reason = DirNoHome
		return "", fail
	}
	config := getenv("XDG_CONFIG_HOME")
	if config == "" {
		config = filepath.Join(home, ".config")
	}
	fail.File = filepath.Join(config, "user-dirs.dirs")
	data, err := os.ReadFile(fail.File)
	if err != nil {
		fail.Reason, fail.Err = DirFileMissing, err
		return "", fail
	}
	value, ok := userDirsValue(string(data))
	if !ok {
		fail.Reason = DirLineMissing
		return "", fail
	}
	dir, ok := expandHome(value, home)
	if !ok {
		fail.Reason, fail.Value = DirRelative, value
		return "", fail
	}
	return checkTemplatesDir(dir, home, fail)
}

// userDirsValue devolve o valor de XDG_TEMPLATES_DIR em user-dirs.dirs, lido como a GLib
// lê, porque é ela que dá a pasta ao Nemo (DE-20): espaços no início e em volta do "=",
// a última linha vence, e o valor vai da primeira à última aspa, sem desfazer escapes.
func userDirsValue(content string) (string, bool) {
	value, found := "", false
	for line := range strings.Lines(content) {
		rest, ok := strings.CutPrefix(strings.TrimLeft(line, " \t"), "XDG_TEMPLATES_DIR")
		if !ok {
			continue
		}
		rest, ok = strings.CutPrefix(strings.TrimLeft(rest, " \t"), "=")
		if !ok {
			continue
		}
		rest, ok = strings.CutPrefix(strings.TrimLeft(rest, " \t"), `"`)
		end := strings.LastIndexByte(rest, '"')
		if !ok || end < 0 {
			continue
		}
		value, found = rest[:end], true
	}
	return value, found
}

// expandHome expande o $HOME do começo do valor. Como na GLib, só "$HOME" sozinho ou
// seguido de "/" é aceito; um valor que não começa por "$HOME" precisa ser absoluto.
func expandHome(value, home string) (string, bool) {
	rest, ok := strings.CutPrefix(value, "$HOME")
	if !ok {
		return value, strings.HasPrefix(value, "/")
	}
	if rest != "" && rest[0] != '/' {
		return "", false
	}
	return filepath.Join(home, rest), true
}

// checkTemplatesDir confere que dir é absoluto, não é a pasta pessoal e é uma pasta. As
// barras do fim saem antes da comparação.
func checkTemplatesDir(dir, home string, fail *TemplatesDirError) (string, error) {
	fail.Value = dir
	if !filepath.IsAbs(dir) {
		fail.Reason = DirRelative
		return "", fail
	}
	dir = filepath.Clean(dir)
	fail.Value = dir
	if home != "" && dir == filepath.Clean(home) {
		fail.Reason = DirIsHome
		return "", fail
	}
	// Stat segue o link: ~/Modelos é um link na máquina do autor (FR-CORE-04).
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		fail.Reason, fail.Err = DirNotDir, err
		return "", fail
	}
	return dir, nil
}
