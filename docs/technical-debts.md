---
title: Dívidas técnicas
aliases:
  - DT-NN
tags:
  - dividas-tecnicas
---

# Dívidas técnicas

Tudo o que ficou para depois ou precisa ser revisto no Janus.

## Abertas

- **DT-02** ⬜ — **O nome `janus` não foi conferido no AUR.** A consulta falhou em 2026-11-09. O servidor WebRTC Janus Gateway é a colisão mais provável. (registrada em 2026-11-09)
  - *Rever quando:* antes de qualquer publicação. Não afeta o uso pessoal.
- **DT-03** ⬜ — **O atalho do yazi não foi executado.** A sintaxe veio da documentação, e a tecla `c` seguida de `m` é sugestão. (registrada em 2026-11-09)
  - *Rever quando:* em [[roadmap|RM-01.02]].
- **DT-06** ⬜ — **O catálogo não declara a stack, e os padrões da casa não têm referência para Go.** O campo `stack` de `catalog-info.yaml` está vazio, e a lista de referências por stack dos padrões só cobre TypeScript, NestJS, React, React Native e TailwindCSS. As convenções deste projeto estão em [[decisions#DE-13 — Código Go segue as convenções do Go; o resto, as da casa|DE-13]]. (registrada em 2026-11-09)
  - *Rever quando:* quando a casa ganhar uma referência para Go, ou no primeiro projeto Go seguinte.
- **DT-08** ⬜ — **A extensão composta perde a primeira parte.** A extensão de referência é o trecho a partir do último ponto, então um modelo `backup.tar.gz` com o nome `foo` cria `foo.gz`. Nenhum modelo atual tem extensão composta. (registrada em 2026-11-16)
  - *Rever quando:* entrar na pasta de modelos um modelo com extensão composta. A saída muda a regra de [[bootstrap#Resolução do nome final]] e pede um DE.
  - *Ver também:* [[similar-tools]].
- **DT-09** ⬜ — **O manifesto não tem como escrever um `*` literal.** No padrão, `*` é sempre curinga, e todos os outros caracteres, inclusive `\`, são literais. Um modelo com `*` no nome só é alcançado por um curinga. Nenhum modelo atual tem `*` no nome. (registrada em 2026-11-16)
  - *Rever quando:* entrar na pasta de modelos um modelo com `*` no nome. Qualquer forma de escape muda o formato do manifesto, que é contrato, e pede um DE.
  - *Ver também:* [[similar-tools]].
- **DT-10** ⬜ — **O plugin do neovim precisa escapar o caminho antes do `:edit`.** O caminho devolvido por `janus new` pode conter `%`, `#` ou `|`, que o `:edit` interpreta. A saída é passar por `vim.fn.fnameescape` e testar com esses caracteres. (registrada em 2026-11-16)
  - *Rever quando:* em [[roadmap|RM-01.03]].
  - *Ver também:* [[similar-tools]].
- **DT-11** ⬜ — **A branch e o worktree de RM-01.01 fogem do padrão.** A branch `feat/rm-01.01_core-and-terminal` tem `.` e `_` no slug, e o worktree mora em `~/.herdr/worktrees/janus/feat-rm-01-01-core-and-terminal`, e não em `~/Projetos/dev/janus/feat_<slug>`, como pede o `AGENTS.md`. Nada foi renomeado. (registrada em 2026-11-16)
  - *Rever quando:* ao abrir a próxima branch. Ou ela segue o padrão, ou o `AGENTS.md` passa a aceitar os worktrees do herdr, com um DE.

## Finalizadas

- **DT-01** ✅ — **O manifesto oculto podia aparecer no menu do Nemo.** A hipótese era que o Nemo ignora itens ocultos no menu "Criar novo documento". O autor confirmou que a pasta oculta `.git`, que mora na pasta de modelos, não aparece no menu. O manifesto fica em `.janus.toml`, na raiz da pasta de modelos, como desenhado. Se o arquivo aparecer no menu quando for criado, abrir uma dívida nova: a saída prevista é movê-lo para `$XDG_CONFIG_HOME/janus/janus.toml`, como único local aceito. (registrada em 2026-11-09; resolvida em 2026-11-09)
- **DT-04** ✅ — **As ferramentas parecidas não foram estudadas.** Só os resultados de uma busca tinham sido lidos. O estudo está em [[similar-tools]]. As lições viraram [[decisions#DE-20 — O user-dirs.dirs é lido como a GLib lê|DE-20]], as dívidas DT-08, DT-09 e DT-10 e casos de teste dos requisitos. (registrada em 2026-11-09; resolvida em 2026-11-16)
- **DT-05** ✅ — **O remoto e o caminho do módulo Go não estavam definidos.** O remoto `origin` apareceu configurado como `git@gitlab.com:rg-dev/janus.git`, e o módulo segue o mesmo endereço: `gitlab.com/rg-dev/janus` ([[decisions#DE-16 — O módulo Go é gitlab.com/rg-dev/janus|DE-16]]). (registrada em 2026-11-09; resolvida em 2026-11-16)
- **DT-07** ✅ — **Os requisitos com ID ainda não tinham sido derivados.** Não existiam `FR` nem `NFR` em `docs/requirements/`. Derivados de [[bootstrap]] em [[core|núcleo]], [[cli|linha de comando]] e [[tui|modo interativo]]. Os módulos `NVIM` e `YAZI` ficam para [[roadmap|RM-01.03]] e [[roadmap|RM-01.02]]. (registrada em 2026-11-09; resolvida em 2026-11-16)
