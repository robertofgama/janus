---
title: Glossário
aliases:
  - glossário
tags:
  - glossario
---

# Glossário

Os termos do domínio do Janus. A coluna Inglês dá o nome do termo no código, no JSON e nos commits, e os identificadores seguem a convenção de cada linguagem (no Go, *templates dir* vira `TemplatesDir`). As escolhas que não são tradução literal, como *template* para modelo, estão justificadas em [[decisions#DE-15 — O glossário fixa os nomes em inglês; modelo é template|DE-15]].

| Termo | Inglês | Significado |
|---|---|---|
| **Modelo** | *template* | Um arquivo regular dentro da pasta de modelos, copiado tal como está. |
| **Pasta de modelos** | *templates dir* | A pasta XDG de modelos do usuário (`XDG_TEMPLATES_DIR`). Na máquina do autor, `~/Modelos`. |
| **Identificador** | *id* | O caminho do modelo relativo à raiz da pasta de modelos, com `/` como separador. |
| **Categoria** | *category* | O caminho da pasta que contém o modelo, relativo à raiz. Vazio para modelos da raiz. |
| **Manifesto** | *manifest* | O arquivo `.janus.toml` na raiz da pasta de modelos, que define nomes fixos. |
| **Regra** | *rule* | Um par padrão → nome fixo na tabela `[names]` do manifesto. |
| **Padrão** | *pattern* | A chave de uma regra. Casa identificadores, com `*` como único curinga. |
| **Nome fixo** | *fixed name* | O nome final predefinido de um modelo, dado por uma regra do manifesto. |
| **Modo de nome** | *name mode* | `fixed` quando o modelo tem nome fixo; `ask` quando o nome é perguntado. |
| **Nome final** | *final name* | O nome com que o arquivo nasce, depois das regras de [[bootstrap#Resolução do nome final]]. |
| **Extensão de referência** | *reference extension* | A extensão que o Janus acrescenta ao nome informado pelo usuário quando ela falta: a do nome fixo ou a do arquivo do modelo. |
| **Pasta de destino** | *destination dir* | A pasta onde o arquivo nasce: a de `--dir`, ou a atual. |
| **Destino** | *destination* | O caminho do arquivo a criar: a pasta de destino mais o nome final. |
| **Núcleo** | *core* | A parte do programa que descobre, lista, resolve nomes e copia. |
| **Contrato** | *contract* | Os subcomandos não interativos `list` e `new`, com suas saídas e códigos. |
| **Modo interativo** | *interactive mode* | O comando `janus` sem subcomando: seletor, campo de nome e criação. |
| **Seletor** | *selector* | Onde o usuário escolhe o modelo: o `fzf` no terminal, `vim.ui.select` no neovim. |
| **Campo de nome** | *name prompt* | Onde o usuário confirma ou edita o nome, a partir de um valor inicial. |
| **Prévia** | *preview* | O que aparece ao lado do seletor do `fzf`: o conteúdo do modelo, ou um aviso se ele for binário ou vazio. |
| **Adaptador** | *adapter* | O código que liga outro programa (yazi, neovim) ao Janus. |
| **Integração** | *integration* | A ligação de um programa ao Janus: o adaptador e as instruções de instalação, em `integrations/<programa>/`. |
