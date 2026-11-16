---
title: Ferramentas parecidas
aliases:
  - similar-tools
  - DT-04
tags:
  - relatorios/estudo
---

# Ferramentas parecidas

Estudo das cinco ferramentas citadas em [[bootstrap#Referências]], feito em 2026-11-16 no planejamento de [[roadmap|RM-01.01]] para resolver [[technical-debts|DT-04]]. O código-fonte das cinco foi lido; o touchtemplate, pelo repositório no GitHub (versão 0.8.0), sem abrir a página do lib.rs. O comportamento do `xdg-user-dir` 0.20 foi testado em pastas temporárias.

O achado principal: **as cinco sobrescrevem o destino sem avisar**. A criação exclusiva do Janus não tem equivalente em nenhuma delas.

## touchtemplate (Rust)

- **Pasta**: lê só o `user-dirs.dirs` e ignora a variável de ambiente, embora o README cite `$XDG_TEMPLATE_DIR`. Expande só `$HOME/`, desfaz escapes com `\` e trata `"$HOME/"` como pasta desativada. Em linha repetida, vale a primeira.
- **Nome**: obrigatório, sem extensão automática. Aceita caminho absoluto como modelo.
- **Conflito**: `fs::copy` sobrescreve. O erro sai na saída padrão com código 0.
- **Prévia**: não tem.

## NeFT (Python)

- **Pasta**: uma lista num arquivo YAML. Com `-xdg`, usa o `xdg-user-dir`, que devolve a pasta pessoal quando a entrada falta ou está desativada.
- **Listagem**: só o primeiro nível, incluindo os itens ocultos.
- **Nome**: pergunta, sugerindo o nome do modelo, sem extensão automática. Aceita `/` no nome.
- **Conflito**: `shutil.copy` sobrescreve, ou copia para dentro de uma pasta com o nome do destino.
- **Prévia**: não tem. Tem ícones e um modo que repete a criação em laço.

## tmpl.sh (Bash)

- **Pasta**: `$TMPL_DIR`, depois `$XDG_TEMPLATES_DIR`, depois `~/Templates`. Não lê o `user-dirs.dirs`.
- **Extras**: `fetch -` manda o conteúdo para a saída padrão, e `-e` abre o arquivo criado no editor.
- **Conflito**: `cp -r` sobrescreve e aceita pasta como modelo. Sai sempre com o código 1.

## ttouch (Python)

- **Pasta**: `~/.templates`, ou uma variável de ambiente cujo nome no README difere do código.
- **Nome**: o modelo é escolhido pelo nome do destino.
- **Conflito**: sem modelo, `open("w")` trunca um arquivo existente. Com modelo, a cópia perde o bit de execução.
- **Nomes estranhos**: o autocompletar quebra em espaço. Usa arquivos ocultos como modelos.

## template.nvim (Lua)

- **Pasta**: `temp_dir` configurado no plugin, sem XDG.
- **Listagem**: descarta o modelo cujo tipo de arquivo o neovim não detecta, como os sem extensão.
- **Nome**: corta o nome no último ponto num lugar e no primeiro em outro. Busca com padrão do Lua.
- **Conflito**: trunca o arquivo existente, força a permissão 0644 e abre com `:e` sem escapar o caminho.
- **Conteúdo**: converte CRLF. A prévia vem do Telescope.

## Lições para o Janus

O que virou decisão, dívida ou requisito:

1. **O identificador é procurado na listagem**, nunca juntado a um caminho. `../x`, `/etc/hostname`, `.janus.toml`, `.git/config` e `./Empty Bash` dão o código 5. Virou critério de requisito.
2. **Lacunas da leitura do `user-dirs.dirs`**: escapes, espaço em volta do `=`, linha repetida e barra no fim. Foram resolvidas em [[decisions#DE-20 — O user-dirs.dirs é lido como a GLib lê|DE-20]], seguindo a GLib, que é quem dá a pasta ao Nemo, e não o `xdg-user-dir`. A sugestão de desfazer `\"` e `\$` foi rejeitada, porque a GLib não desfaz.
3. **O `xdg-user-dir` ignora a variável de ambiente**, e o Janus a lê primeiro, como a spec já define. A divergência é intencional.
4. **Extensão composta**: `backup.tar.gz` com o nome `foo` vira `foo.gz`. Virou [[technical-debts|DT-08]].
5. **`*` literal no manifesto**: não há como escrever. Virou [[technical-debts|DT-09]].
6. **Caminho no `:edit` do neovim**: precisa de `fnameescape`, por causa de `%`, `#` e `|`. Virou [[technical-debts|DT-10]].
7. **Modelo de arquivo oculto**: como os itens ocultos são ignorados, o modelo fica guardado sem o ponto, e o nome com ponto vem do manifesto. Entra no README.
8. **CRLF**: a cópia byte a byte precisa preservar as quebras de linha do Windows. Virou critério de requisito.
9. **Ideias fora do escopo**: mandar o conteúdo para a saída padrão já está previsto em [[roadmap|RM-02.01]]; criar vários arquivos seguidos não tem entrega e pediria um RM.
