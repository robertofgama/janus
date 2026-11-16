package core

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"unicode"
)

// Catalog é a pasta de modelos lida: os modelos, em ordem, com o modo de nome resolvido.
type Catalog struct {
	Dir       string
	Templates []Template
}

// TemplateNotFoundError é um identificador que a listagem não traz (código 5).
type TemplateNotFoundError struct {
	ID string
}

func (e *TemplateNotFoundError) Error() string {
	return fmt.Sprintf("template %q not found", e.ID)
}

// LoadCatalog lê o manifesto e a pasta de modelos dir e resolve a regra de cada modelo.
// O manifesto é lido antes da pasta para que um manifesto malformado (código 4) vença uma
// falha de leitura da pasta (código 1), como pede a ordem de verificação (FR-CLI-11).
func LoadCatalog(dir string) (Catalog, error) {
	m, err := LoadManifest(dir)
	if err != nil {
		return Catalog{}, err
	}
	templates, err := listTemplates(dir)
	if err != nil {
		return Catalog{}, err
	}
	for i := range templates {
		if err := m.resolve(&templates[i]); err != nil {
			return Catalog{}, err
		}
	}
	return Catalog{Dir: dir, Templates: templates}, nil
}

// IDs devolve os identificadores, na ordem da listagem.
func (c Catalog) IDs() []string {
	ids := make([]string, len(c.Templates))
	for i, t := range c.Templates {
		ids[i] = t.ID
	}
	return ids
}

// Find procura o modelo pelo identificador, só na listagem. O identificador nunca é
// juntado à pasta para formar um caminho, então "../x" ou ".janus.toml" não são modelos
// (FR-CORE-15).
func (c Catalog) Find(id string) (Template, error) {
	for _, t := range c.Templates {
		if t.ID == id {
			return t, nil
		}
	}
	return Template{}, &TemplateNotFoundError{ID: id}
}

// listTemplates percorre a pasta de modelos e devolve os modelos em ordem, todos em modo
// ask (FR-CORE-05 a FR-CORE-09).
func listTemplates(root string) ([]Template, error) {
	var templates []Template
	var walk func(category string) error
	walk = func(category string) error {
		// os.ReadDir segue a própria pasta quando ela é um link (FR-CORE-04); as entradas vêm
		// com o tipo do link, e não do alvo.
		entries, err := os.ReadDir(filepath.Join(root, category))
		if err != nil {
			return err
		}
		for _, e := range entries {
			name := e.Name()
			if strings.HasPrefix(name, ".") || strings.ContainsFunc(name, unicode.IsControl) {
				continue // FR-CORE-06, FR-CORE-08
			}
			id := path.Join(category, name)
			switch {
			case e.IsDir():
				if err := walk(id); err != nil {
					return err
				}
			case e.Type()&fs.ModeSymlink != 0:
				// Link para arquivo é modelo; quebrado é ignorado; link para pasta não é
				// seguido, o que também evita laços (FR-CORE-07).
				if info, err := os.Stat(filepath.Join(root, id)); err == nil && info.Mode().IsRegular() {
					templates = append(templates, newTemplate(root, category, name))
				}
			case e.Type().IsRegular():
				templates = append(templates, newTemplate(root, category, name))
			}
		}
		return nil
	}
	if err := walk(""); err != nil {
		return nil, err
	}
	slices.SortFunc(templates, func(a, b Template) int {
		if c := strings.Compare(strings.ToLower(a.ID), strings.ToLower(b.ID)); c != 0 {
			return c
		}
		return strings.Compare(a.ID, b.ID)
	})
	return templates, nil
}

func newTemplate(root, category, file string) Template {
	id := path.Join(category, file)
	return Template{ID: id, Category: category, File: file, Path: filepath.Join(root, id), Mode: NameAsk}
}

// resolve aplica as regras do manifesto a t (FR-CORE-13). Um padrão sem "*" igual ao
// identificador vence; senão, conta quantos padrões com "*" casam.
func (m Manifest) resolve(t *Template) error {
	var matches []Rule
	for _, r := range m.Rules {
		if !strings.Contains(r.Pattern, "*") {
			if r.Pattern == t.ID {
				t.Mode, t.FixedName = NameFixed, r.Name
				return nil
			}
			continue
		}
		if matchPattern(r.Pattern, t.ID) {
			matches = append(matches, r)
		}
	}
	switch len(matches) {
	case 0:
		return nil
	case 1:
		t.Mode, t.FixedName = NameFixed, matches[0].Name
		return nil
	}
	patterns := make([]string, len(matches))
	for i, r := range matches {
		patterns[i] = r.Pattern
	}
	return &ManifestError{Path: m.Path, Problem: ManifestAmbiguous, ID: t.ID, Patterns: patterns}
}
