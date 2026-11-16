package core

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// ManifestFile é o nome do manifesto, na raiz da pasta de modelos.
const ManifestFile = ".janus.toml"

// Rule é uma regra do manifesto: um padrão e o nome fixo que ele dá.
type Rule struct {
	Pattern string
	Name    string
}

// Manifest guarda as regras na ordem em que aparecem no arquivo.
type Manifest struct {
	Path  string
	Rules []Rule
}

// ManifestProblem diz o que torna o manifesto inválido.
type ManifestProblem int

const (
	ManifestSyntax        ManifestProblem = iota + 1 // TOML malformado, inclusive chave repetida
	ManifestUnknownKey                               // chave ou tabela fora de [names]
	ManifestNamesNotTable                            // names existe, mas não é uma tabela
	ManifestNotString                                // valor de regra que não é texto
	ManifestDoubleStar                               // padrão com "**", reservado
	ManifestInvalidName                              // nome fixo inválido; Err traz o InvalidNameError
	ManifestAmbiguous                                // dois ou mais padrões com "*" casam o mesmo modelo
)

// ManifestError é o manifesto inválido (código 4).
type ManifestError struct {
	Path     string
	Problem  ManifestProblem
	Key      string   // a chave ou o padrão envolvido
	Line     int      // a linha, no erro de sintaxe
	ID       string   // o modelo, em ManifestAmbiguous
	Patterns []string // os padrões em conflito, na ordem do arquivo
	Err      error    // a causa, quando houver
}

func (e *ManifestError) Error() string {
	return fmt.Sprintf("invalid manifest %s (problem %d, key %q)", e.Path, e.Problem, e.Key)
}

func (e *ManifestError) Unwrap() error { return e.Err }

// LoadManifest lê o manifesto da pasta de modelos dir. Arquivo ausente, vazio ou sem a
// tabela [names] é um manifesto sem regras (FR-CORE-10).
//
// O TOML é lido como tabela genérica e conferido chave a chave, na ordem do arquivo
// (MetaData.Keys), para a mensagem citar sempre a primeira chave errada. Ler direto numa
// struct foi rejeitado: o erro de tipo da biblioteca não diz a chave de forma estruturada.
func LoadManifest(dir string) (Manifest, error) {
	path := filepath.Join(dir, ManifestFile)
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Manifest{Path: path}, nil
	}
	if err != nil {
		return Manifest{}, err
	}
	var raw map[string]any
	md, err := toml.Decode(string(data), &raw)
	if err != nil {
		e := &ManifestError{Path: path, Problem: ManifestSyntax, Err: err}
		var parseErr toml.ParseError
		if errors.As(err, &parseErr) {
			e.Line = parseErr.Position.Line
		}
		return Manifest{}, e
	}
	return parseRules(path, raw, md.Keys())
}

// parseRules confere a estrutura do manifesto e extrai as regras (FR-CORE-11, FR-CORE-12,
// FR-CORE-14).
func parseRules(path string, raw map[string]any, keys []toml.Key) (Manifest, error) {
	fail := func(problem ManifestProblem, key string, cause error) (Manifest, error) {
		return Manifest{}, &ManifestError{Path: path, Problem: problem, Key: key, Err: cause}
	}
	for _, k := range keys {
		if k[0] != "names" {
			return fail(ManifestUnknownKey, k[0], nil)
		}
	}
	m := Manifest{Path: path}
	value, ok := raw["names"]
	if !ok {
		return m, nil
	}
	names, ok := value.(map[string]any)
	if !ok {
		return fail(ManifestNamesNotTable, "names", nil)
	}
	seen := map[string]bool{}
	for _, k := range keys {
		if len(k) < 2 || seen[k[1]] {
			continue
		}
		pattern := k[1]
		seen[pattern] = true
		name, ok := names[pattern].(string)
		if !ok {
			return fail(ManifestNotString, pattern, nil)
		}
		if strings.Contains(pattern, "**") {
			return fail(ManifestDoubleStar, pattern, nil)
		}
		clean, err := cleanName(name)
		if err != nil {
			return fail(ManifestInvalidName, pattern, err)
		}
		m.Rules = append(m.Rules, Rule{Pattern: pattern, Name: clean})
	}
	return m, nil
}
