package messages

// Help é a ajuda de `janus --help` (FR-CLI-03).
const Help = `janus cria arquivos a partir da pasta de modelos, a mesma do "Criar novo
documento" do Nemo.

Uso:
  janus [--hold]            escolhe um modelo com o fzf e cria o arquivo na pasta atual
  janus list [--json]       lista os modelos
  janus new <id> [--dir <pasta>] [--name <nome>]
                            cria um arquivo a partir de um modelo, sem perguntar nada
  janus <subcomando> --help mostra a ajuda do subcomando
  janus --version           mostra a versão

Opção do modo interativo:
  --hold   depois de um erro, espera um Enter antes de sair (para o yazi)

A pasta de modelos vem de XDG_TEMPLATES_DIR ou de ~/.config/user-dirs.dirs. Nomes
fixos vêm do manifesto .janus.toml, na raiz da pasta de modelos.

Códigos de saída:
  0  sucesso                       6  nome obrigatório não informado
  1  erro inesperado               7  nome inválido
  2  uso incorreto                 8  destino já existe
  3  pasta de modelos ausente      9  fzf ausente
  4  manifesto inválido           10  pasta de destino inválida
  5  modelo não encontrado       130  cancelado
`

// ListHelp é a ajuda de `janus list --help` (FR-CLI-03).
const ListHelp = `Uso: janus list [--json]

Lista os modelos, um identificador por linha, em ordem alfabética.

  --json   imprime uma lista JSON com id, category, file e name (mode, e value ou
           suggestion), para outros programas
`

// NewHelp é a ajuda de `janus new --help` (FR-CLI-03).
const NewHelp = `Uso: janus new <id> [--dir <pasta>] [--name <nome>]

Cria um arquivo a partir do modelo <id> e imprime o caminho dele. Nunca pergunta e nunca
sobrescreve.

  --dir <pasta>   pasta onde o arquivo nasce; sem ela, a pasta atual
  --name <nome>   nome do arquivo; obrigatório quando o modelo não tem nome fixo. A
                  extensão do modelo é acrescentada quando falta.

Use -- antes de um identificador que começa com "-".
`
