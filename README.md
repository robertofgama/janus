<p align="center">
  <img src="assets/icons/janus-icon-256.png" width="96" alt="Janus">
</p>

<h1 align="center">Janus</h1>

Ferramenta de linha de comando que leva o "Criar novo documento" do gerenciador de arquivos para o terminal, o yazi e o neovim, usando a mesma pasta de modelos (`~/Modelos`) que o Nemo já usa.

## Estado

Núcleo, terminal e yazi prontos ([RM-01.01 e RM-01.02](docs/roadmap.md)). A integração com o neovim é a próxima entrega.

## Documentação

- [Especificação de bootstrap](docs/specs/bootstrap.md): o desenho completo e o ponto de partida do projeto.
- [Glossário](docs/glossary.md): os termos do domínio, com o nome em inglês de cada um.
- [Decisões](docs/decisions.md): o que foi decidido, por quê e o que foi rejeitado.
- [Roadmap](docs/roadmap.md): as entregas, na ordem.
- [Dívidas técnicas](docs/technical-debts.md): o que ficou para depois.

## Instalação

```
go install ./cmd/janus
```

O binário vai para `~/go/bin`. Em tempo de execução, o modo interativo precisa do `fzf`; o `bat`, se existir, realça a prévia. `list` e `new` não dependem de nenhum dos dois.

## Uso

```
janus                          # escolhe um modelo com o fzf e cria o arquivo na pasta atual
janus list                     # lista os modelos, um por linha
janus list --json              # lista os modelos em JSON, para outros programas
janus new <id> --name <nome>   # cria sem perguntar nada
janus --help                   # ajuda, com os códigos de saída
```

O modo interativo imprime só o caminho criado na saída padrão, então `nvim (janus)` funciona no fish.

## Yazi

Tecle `c` e depois `m` no yazi para criar um arquivo, a partir de um modelo, na pasta em que ele está. A instalação do atalho está em [integrations/yazi](integrations/yazi/README.md).

## Nomes fixos

Sem configuração, todo modelo pergunta o nome, sugerindo o do próprio arquivo. Para dar um nome fixo, crie `.janus.toml` na raiz da pasta de modelos:

```toml
[names]
"info.ai/ai-memory/*"     = ".ai-memory.toml"
"info.ai/agents/*"        = "AGENTS.md"
"info.dev/package.json.*" = "package.json"
```

O padrão casa o identificador inteiro, e `*` casa qualquer trecho sem `/`. As regras completas estão na [especificação](docs/specs/bootstrap.md).

### Modelo de arquivo oculto

Itens ocultos da pasta de modelos são ignorados, como no Nemo. Para um modelo que deve nascer com um nome que começa com ponto, guarde-o sem o ponto e dê o nome no manifesto, como `.ai-memory.toml` acima.

## Padrões

Projeto da Darth Gama Software. Segue os padrões da casa, resumidos na especificação.
