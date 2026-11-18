package core

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// DestinationDirReason diz por que a pasta de destino foi recusada.
type DestinationDirReason int

const (
	DestNotDir          DestinationDirReason = iota + 1 // inexistente ou que não é pasta
	DestInsideTemplates                                 // a pasta de modelos ou uma pasta dentro dela (DE-23)
)

// DestinationDirError é a pasta de destino recusada (código 10): inexistente, que não é
// pasta, ou a pasta de modelos e o que fica dentro dela.
type DestinationDirError struct {
	Dir          string // a pasta como foi informada
	Reason       DestinationDirReason
	TemplatesDir string // a pasta de modelos, só em DestInsideTemplates
	Err          error
}

func (e *DestinationDirError) Error() string {
	if e.Reason == DestInsideTemplates {
		return fmt.Sprintf("destination dir %q is inside the templates dir %q", e.Dir, e.TemplatesDir)
	}
	return fmt.Sprintf("destination dir %q not found", e.Dir)
}

func (e *DestinationDirError) Unwrap() error { return e.Err }

// DestinationExistsError é o destino que já existe (código 8).
type DestinationExistsError struct {
	Path string
}

func (e *DestinationExistsError) Error() string {
	return fmt.Sprintf("destination %s already exists", e.Path)
}

// DestinationDir confere a pasta de destino e devolve o caminho absoluto dela, sem resolver
// os links, para o caminho impresso ser o que o usuário vê. Um link para pasta é aceito
// (FR-CORE-18). O caminho vazio é recusado antes do filepath.Abs, que o trataria como a
// pasta atual. A pasta de modelos e o que fica dentro dela são recusados (FR-CORE-24): a
// verificação mora aqui porque `new` e o modo interativo já passam por esta função, cada
// um no ponto que a ordem do contrato pede (DE-26).
func DestinationDir(dir, templatesDir string) (string, error) {
	notDir := func(err error) error {
		return &DestinationDirError{Dir: dir, Reason: DestNotDir, Err: err}
	}
	if dir == "" {
		return "", notDir(nil)
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", notDir(err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", notDir(err)
	}
	if !info.IsDir() {
		return "", notDir(nil)
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", notDir(err)
	}
	inside, err := within(resolved, templatesDir)
	if err != nil {
		return "", err
	}
	if inside {
		return "", &DestinationDirError{Dir: dir, Reason: DestInsideTemplates, TemplatesDir: templatesDir}
	}
	return abs, nil
}

// within diz se realDir, um caminho já sem links, é root ou fica dentro dela. root é
// absoluta, como sai de FindTemplatesDir, e os links dela são resolvidos porque ~/Modelos
// costuma ser um link (DE-23). O filepath.Rel com o filepath.IsLocal evita o falso positivo
// da comparação por prefixo de texto, que acharia Modelos-antigo dentro de Modelos. Uma
// montagem com mount --bind escapa desta comparação; os.SameFile subindo até a raiz pegaria
// esse caso, mas custa um Stat por nível e foi rejeitado (DE-26). Uma falha ao resolver
// root, que já foi validada na descoberta, sai como erro inesperado.
func within(realDir, root string) (bool, error) {
	realRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return false, err
	}
	rel, err := filepath.Rel(realRoot, realDir)
	if err != nil {
		return false, err
	}
	return filepath.IsLocal(rel), nil
}

// copyContent copia o conteúdo do modelo. É variável só para o teste simular uma falha
// no meio da cópia (FR-CORE-22).
var copyContent = io.Copy

// Create cria name em destDir com o conteúdo e as permissões de t e devolve o caminho
// criado. destDir deve vir de DestinationDir.
//
// A criação é exclusiva (O_EXCL): falha se existir qualquer coisa com o nome, até um link
// quebrado, e fecha a janela entre checar e criar (FR-CORE-19). Checar antes com Stat foi
// rejeitado por deixar essa janela aberta. O destino é aberto antes do modelo para que
// "destino existe" (8) vença uma falha de leitura (1), como pede a ordem de verificação.
func Create(t Template, destDir, name string) (string, error) {
	dest := filepath.Join(destDir, name)
	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if errors.Is(err, fs.ErrExist) {
		return "", &DestinationExistsError{Path: dest}
	}
	if err != nil {
		return "", err
	}
	if err := fill(out, t.Path); err != nil {
		out.Close()
		os.Remove(dest) // FR-CORE-22
		return "", err
	}
	if err := out.Close(); err != nil {
		os.Remove(dest)
		return "", err
	}
	return dest, nil
}

// fill copia o modelo src para out, byte a byte (FR-CORE-20), e aplica as permissões dele.
// O Chmod no arquivo aberto aplica as permissões exatas; as do OpenFile passam pela umask
// (FR-CORE-21). Como o modelo é aberto com os.Open, um link é seguido até o alvo.
func fill(out *os.File, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	info, err := in.Stat()
	if err != nil {
		return err
	}
	if _, err := copyContent(out, in); err != nil {
		return err
	}
	return out.Chmod(info.Mode().Perm())
}
