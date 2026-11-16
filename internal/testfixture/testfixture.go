// Package testfixture monta, numa pasta temporária, a pasta de modelos de teste descrita
// em docs/specs/bootstrap.md ("Cenários de aceite"), tira fotos de pastas e põe programas
// falsos no PATH. Só os testes o importam.
package testfixture

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Manifest é o manifesto da pasta de teste: as três regras do exemplo da spec.
const Manifest = `[names]
"info.ai/ai-memory/*"     = ".ai-memory.toml"
"info.ai/agents/*"        = "AGENTS.md"
"info.dev/package.json.*" = "package.json"
`

// Ambiguous é o manifesto do cenário 14: dois padrões com "*" casam
// info.ai/agents/react.AGENTS.md.
const Ambiguous = `[names]
"info.ai/agents/*"  = "AGENTS.md"
"info.ai/*/react.*" = "REACT.md"
`

// Binary é o conteúdo de 3d/fatiar.3mf: tem bytes nulos, como um arquivo binário.
const Binary = "PK\x03\x04\x00\x00\x00\x00fatiar"

// IDs são os identificadores da pasta de teste, na ordem da listagem (FR-CORE-09).
var IDs = []string{
	"3d/fatiar.3mf",
	"direito/01.2. Civ. Contestação.tex",
	"Empty Bash",
	"Empty File",
	"info.ai/agents/AGENTS.md",
	"info.ai/agents/react.AGENTS.md",
	"info.ai/ai-memory/default.ai-memory.toml",
	"info.dev/package.json.monorepo",
	"info.dev/react; script; simple.html",
}

// files é o conteúdo da pasta de teste, fora o manifesto.
var files = []struct {
	rel     string
	content string
	perm    fs.FileMode
}{
	{".git/config", "[core]\n", 0o644},
	{"Empty Bash", "#!/usr/bin/env bash\n", 0o755},
	{"Empty File", "", 0o644},
	{"3d/fatiar.3mf", Binary, 0o644},
	{"direito/01.2. Civ. Contestação.tex", "\\documentclass{article}\n% Contestação\n", 0o644},
	{"info.ai/agents/AGENTS.md", "# AGENTS\n", 0o644},
	{"info.ai/agents/react.AGENTS.md", "# AGENTS para React\n", 0o644},
	{"info.ai/ai-memory/default.ai-memory.toml", "workspace = \"default\"\n", 0o644},
	{"info.dev/package.json.monorepo", "{\n  \"private\": true\n}\n", 0o644},
	{"info.dev/react; script; simple.html", "<!doctype html>\n<title>simple</title>\n", 0o644},
}

// New monta a pasta de teste numa pasta temporária, com o manifesto, e devolve o caminho.
func New(t testing.TB) string {
	t.Helper()
	dir := t.TempDir()
	for _, f := range files {
		Write(t, dir, f.rel, f.content, f.perm)
	}
	Write(t, dir, ".janus.toml", Manifest, 0o644)
	return dir
}

// Use monta a pasta de teste e aponta XDG_TEMPLATES_DIR para ela.
func Use(t testing.TB) string {
	t.Helper()
	dir := New(t)
	t.Setenv("XDG_TEMPLATES_DIR", dir)
	return dir
}

// Write cria rel dentro de dir, com as pastas que faltarem, aplica perm exatamente e
// devolve o caminho do arquivo.
func Write(t testing.TB, dir, rel, content string, perm fs.FileMode) string {
	t.Helper()
	path := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), perm); err != nil {
		t.Fatal(err)
	}
	// WriteFile passa pela umask; o Chmod garante a permissão pedida.
	if err := os.Chmod(path, perm); err != nil {
		t.Fatal(err)
	}
	return path
}

// Read devolve o conteúdo de path.
func Read(t testing.TB, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// Snapshot descreve tudo o que há em dir, inclusive os itens ocultos: caminho, tipo e
// permissões, data de modificação e resumo do conteúdo. Duas fotos iguais mostram que
// nada foi escrito na pasta (NFR-CORE-01).
func Snapshot(t testing.TB, dir string) string {
	t.Helper()
	var b strings.Builder
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		fmt.Fprintf(&b, "%s %v %d", rel, info.Mode(), info.ModTime().UnixNano())
		if info.Mode().IsRegular() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			fmt.Fprintf(&b, " %x", sha256.Sum256(data))
		}
		b.WriteByte('\n')
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return b.String()
}

// Reply é a resposta de um programa falso a uma chamada.
type Reply struct {
	Stdout string
	Stderr string
	Code   int
}

// Call é uma chamada recebida por um programa falso.
type Call struct {
	Args  []string
	Stdin string
}

// Fake é um programa falso posto no começo do PATH.
type Fake struct {
	dir string
}

// fakeScript registra os argumentos, separados por NUL, e a entrada de cada chamada, e
// responde com a resposta de mesmo número. A pasta vem de $0, que é o caminho do script.
const fakeScript = `#!/bin/sh
dir=${0%/*}
n=0
[ -f "$dir/count" ] && read -r n < "$dir/count"
n=$((n + 1))
echo "$n" > "$dir/count"
for arg in "$@"; do printf '%s\0' "$arg"; done > "$dir/call-$n.args"
cat > "$dir/call-$n.stdin"
[ -f "$dir/reply-$n.out" ] && cat "$dir/reply-$n.out"
[ -f "$dir/reply-$n.err" ] && cat "$dir/reply-$n.err" >&2
code=0
[ -f "$dir/reply-$n.code" ] && read -r code < "$dir/reply-$n.code"
exit "$code"
`

// FakeProgram põe no começo do PATH um programa falso chamado name. A chamada N recebe
// replies[N-1]; as chamadas além delas saem com 0 e sem saída.
func FakeProgram(t testing.TB, name string, replies ...Reply) *Fake {
	t.Helper()
	dir := t.TempDir()
	for i, r := range replies {
		n := i + 1
		Write(t, dir, fmt.Sprintf("reply-%d.out", n), r.Stdout, 0o644)
		Write(t, dir, fmt.Sprintf("reply-%d.err", n), r.Stderr, 0o644)
		Write(t, dir, fmt.Sprintf("reply-%d.code", n), strconv.Itoa(r.Code), 0o644)
	}
	Write(t, dir, name, fakeScript, 0o755)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return &Fake{dir: dir}
}

// Calls devolve as chamadas recebidas, na ordem.
func (f *Fake) Calls(t testing.TB) []Call {
	t.Helper()
	var calls []Call
	for n := 1; ; n++ {
		args, err := os.ReadFile(filepath.Join(f.dir, fmt.Sprintf("call-%d.args", n)))
		if errors.Is(err, fs.ErrNotExist) {
			return calls
		}
		if err != nil {
			t.Fatal(err)
		}
		list := strings.Split(string(args), "\x00")
		calls = append(calls, Call{
			Args:  list[:len(list)-1], // o último NUL deixa um elemento vazio no fim
			Stdin: Read(t, filepath.Join(f.dir, fmt.Sprintf("call-%d.stdin", n))),
		})
	}
}
