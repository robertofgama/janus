package core

import (
	"bytes"
	"errors"
	"io"
	"os"
)

// PreviewKind diz como o modelo aparece na prévia.
type PreviewKind int

const (
	PreviewText   PreviewKind = iota + 1 // conteúdo mostrado
	PreviewEmpty                         // aviso de arquivo vazio
	PreviewBinary                        // aviso de arquivo binário, com o tamanho
)

// binarySniffLen é quantos bytes do começo são examinados à procura de um byte nulo, o
// mesmo critério do git.
const binarySniffLen = 8000

// Inspect classifica o arquivo path para a prévia e devolve o tamanho dele (FR-CORE-23).
func Inspect(path string) (PreviewKind, int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return 0, 0, err
	}
	if info.Size() == 0 {
		return PreviewEmpty, 0, nil
	}
	head := make([]byte, binarySniffLen)
	n, err := io.ReadFull(f, head)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return 0, 0, err
	}
	if bytes.IndexByte(head[:n], 0) >= 0 {
		return PreviewBinary, info.Size(), nil
	}
	return PreviewText, info.Size(), nil
}
