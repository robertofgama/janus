# Janus

Ferramenta de linha de comando que leva o "Criar novo documento" do gerenciador de arquivos para o terminal, o yazi e o neovim, usando a mesma pasta de modelos (`~/Modelos`) que o Nemo já usa.

## Estado

Em especificação. Ainda não há código.

## Documentação

- [Especificação de bootstrap](docs/specs/bootstrap.md): o desenho completo e o ponto de partida do projeto.
- [Decisões](docs/decisions.md): o que foi decidido, por quê e o que foi rejeitado.
- [Roadmap](docs/roadmap.md): as entregas, na ordem.
- [Dívidas técnicas](docs/technical-debts.md): o que ficou para depois.

## Uso previsto

```
janus                          # escolhe um modelo e cria o arquivo na pasta atual
janus list --json              # lista os modelos, para outros programas
janus new <id> --name <nome>   # cria sem interação
```

## Padrões

Projeto da Darth Gama Software. Segue os padrões da casa, resumidos na especificação.
