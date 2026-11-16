package core

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// DestinationDirError é a pasta de destino inexistente ou que não é pasta (código 10).
type DestinationDirError struct {
	Dir string // a pasta como foi informada
	Err error
}

func (e *DestinationDirError) Error() string {
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

// DestinationDir confere a pasta de destino e devolve o caminho absoluto dela. Um link
// para pasta é aceito (FR-CORE-18). O caminho vazio é recusado antes do filepath.Abs, que
// o trataria como a pasta atual.
func DestinationDir(dir string) (string, error) {
	if dir == "" {
		return "", &DestinationDirError{Dir: dir}
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", &DestinationDirError{Dir: dir, Err: err}
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", &DestinationDirError{Dir: dir, Err: err}
	}
	if !info.IsDir() {
		return "", &DestinationDirError{Dir: dir}
	}
	return abs, nil
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
