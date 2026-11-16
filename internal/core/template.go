package core

// NameMode diz como o nome final de um modelo é decidido.
type NameMode string

const (
	// NameAsk: o nome é perguntado, sugerindo o nome do arquivo do modelo.
	NameAsk NameMode = "ask"
	// NameFixed: o nome vem de uma regra do manifesto.
	NameFixed NameMode = "fixed"
)

// Template é um modelo da pasta de modelos.
type Template struct {
	ID        string // caminho relativo à raiz, com "/": o identificador
	Category  string // caminho da pasta que contém o modelo; "" na raiz
	File      string // nome do arquivo do modelo
	Path      string // caminho completo, de onde o conteúdo é lido
	Mode      NameMode
	FixedName string // o nome fixo; só em NameFixed
}

// ReferenceExtension é a extensão acrescentada ao nome informado quando ela falta: a do
// nome fixo, ou a do arquivo do modelo (FR-CORE-17).
func (t Template) ReferenceExtension() string {
	if t.Mode == NameFixed {
		return extension(t.FixedName)
	}
	return extension(t.File)
}
