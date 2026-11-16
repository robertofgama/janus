---
title: Roadmap
aliases:
  - RM-NN
tags:
  - roadmap
---

# Roadmap

Para onde o Janus vai, em entregas. O que cada entrega precisa fazer está especificado em [[bootstrap]]; as tarefas de cada uma nascem no plano correspondente, em `docs/plans/`.

## Abertas

### RM-01 — MVP

- **RM-01.01** 🔄 — **Núcleo e terminal.** Descoberta da pasta, listagem, manifesto, resolução de nome, `list`, `new`, modo interativo com `fzf`, prévia, campo de nome, `--hold`, códigos de saída e os testes automatizados. (registrada em 2026-11-09)
  - *Por que aqui:* resolve sozinha a dor original e é a base de todas as outras entregas.
- **RM-01.02** ⬜ — **Integração com o yazi.** Trecho de `keymap.toml`, instruções de instalação e o roteiro manual executado. (registrada em 2026-11-09)
  - *Por que aqui:* custa poucas linhas e prova o nível trivial de integração.
  - *Depende de:* RM-01.01.
- **RM-01.03** ⬜ — **Integração com o neovim.** Plugin com `:Janus [pasta]`, instruções de instalação e o roteiro manual executado. (registrada em 2026-11-09)
  - *Por que aqui:* prova o nível nativo de integração, que é o que torna a ferramenta extensível.
  - *Depende de:* RM-01.01.

### RM-02 — Depois do MVP

- **RM-02.01** ⬜ — **Preencher o buffer atual no neovim.** Subcomando público `janus show <id>`, que imprime o conteúdo de um modelo, e um comando do plugin que despeja esse conteúdo no buffer. (registrada em 2026-11-09)
  - *Por que aqui:* é a primeira evolução pedida pelo autor ([[decisions#DE-07 — No neovim, criar o arquivo e abrir|DE-07]]).
  - *Depende de:* RM-01.03.
- **RM-02.02** ⬜ — **Autocompletar no fish** para `janus new`. (registrada em 2026-11-09)
  - *Depende de:* RM-01.01.
- **RM-02.03** ⬜ — **Cursor no arquivo criado**, no yazi. (registrada em 2026-11-09)
  - *Depende de:* RM-01.02.
- **RM-02.04** ⬜ — **Abrir no editor** logo após criar pelo terminal. (registrada em 2026-11-09)
  - *Depende de:* RM-01.01.
- **RM-02.05** ⬜ — **Adaptadores para outros programas**, conforme a necessidade aparecer. (registrada em 2026-11-09)
  - *Depende de:* RM-01.01.

## Finalizadas

_Nenhuma entrega concluída ainda._
