package core

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCoreHasNoTerminalIO(t *testing.T) {
	t.Run("does not touch the terminal nor other internal packages (NFR-CORE-02)", func(t *testing.T) {
		files, err := filepath.Glob("*.go")
		if err != nil {
			t.Fatal(err)
		}
		forbidden := []string{"os.Stdin", "os.Stdout", "os.Stderr", "os.Exit", "fmt.Print"}
		for _, file := range files {
			if strings.HasSuffix(file, "_test.go") {
				continue
			}
			src, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			for _, word := range forbidden {
				if strings.Contains(string(src), word) {
					t.Errorf("%s uses %s", file, word)
				}
			}
			parsed, err := parser.ParseFile(token.NewFileSet(), file, src, parser.ImportsOnly)
			if err != nil {
				t.Fatal(err)
			}
			for _, imp := range parsed.Imports {
				if strings.Contains(imp.Path.Value, "/internal/") {
					t.Errorf("%s imports %s", file, imp.Path.Value)
				}
			}
		}
	})
}
