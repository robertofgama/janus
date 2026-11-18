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
- **DT-12** ⬜ — **Nomes de arquivo que não são UTF-8 válido não têm regra.** A listagem aceita o nome como vem, mas `list --json` troca os bytes inválidos por U+FFFD, e o identificador impresso deixa de ser encontrado por `janus new`. Nenhum modelo atual tem um nome assim. (registrada em 2026-11-16)
  - *Rever quando:* aparecer na pasta de modelos um nome que não é UTF-8. Saídas possíveis: ignorar esses nomes na listagem, como os que têm caractere de controle, ou recusá-los com um aviso. Pede um DE.
- **DT-13** ⬜ — **A pasta de modelos pode ser o destino.** A [[decisions#DE-02 — A pasta XDG de modelos é a fonte única, compartilhada com o Nemo|DE-02]] diz que o Janus nunca escreve na pasta de modelos, mas `janus new --dir` apontando para ela, ou `janus` rodado dentro dela, cria o arquivo lá: o contrato não tem erro para esse caso. Visto na revisão de RM-01.01. A saída foi decidida em [[decisions#DE-23 — A pasta de destino nunca fica dentro da pasta de modelos|DE-23]]: recusar com o código 10, já na spec e em FR-CORE-24; falta o código. (registrada em 2026-11-16)
  - *Rever quando:* em [[roadmap|RM-01.02]], antes do atalho do yazi ir para o dotfiles do autor. Junto com o código, o critério de FR-TUI-02 ganha o caso do 10.
- **DT-14** ⬜ — **O `--accept-nth` do usuário muda a saída do seletor.** O seletor anula as opções do `FZF_DEFAULT_OPTS` que mudam a saída ([[tui|NFR-TUI-02]]), mas o `fzf` 0.74 não tem `--no-accept-nth`. Com `--accept-nth` nas opções do usuário, o seletor devolve só parte do identificador, e o `janus` falha com o código 5. (registrada em 2026-11-16)
  - *Rever quando:* o `fzf` ganhar uma forma de anular a opção, ou alguém usar `--accept-nth` no `FZF_DEFAULT_OPTS`. Saída possível: passar um `--accept-nth` que devolva a linha inteira, conferido com identificadores que têm espaço e `;`.
- **DT-15** ⬜ — **O aviso de destino existente pode ser cortado no campo de nome.** O aviso cita o caminho completo, e com uma pasta longa o `fzf` o corta na largura do terminal. O nome que colidiu continua visível no campo. Visto na conferência manual. (registrada em 2026-11-16)
  - *Rever quando:* incomodar no uso. Saída: citar só o nome no aviso do campo de nome e manter o caminho completo na mensagem de `janus new`.
  - *Ver também:* [[rm-01.01-manual-check]].

## Finalizadas

- **DT-01** ✅ — **O manifesto oculto podia aparecer no menu do Nemo.** A hipótese era que o Nemo ignora itens ocultos no menu "Criar novo documento". O autor confirmou que a pasta oculta `.git`, que mora na pasta de modelos, não aparece no menu. O manifesto fica em `.janus.toml`, na raiz da pasta de modelos, como desenhado. Se o arquivo aparecer no menu quando for criado, abrir uma dívida nova: a saída prevista é movê-lo para `$XDG_CONFIG_HOME/janus/janus.toml`, como único local aceito. (registrada em 2026-11-09; resolvida em 2026-11-09)
- **DT-03** ✅ — **O atalho do yazi não foi executado.** A sintaxe veio da documentação, e a tecla `c` seguida de `m` era sugestão. O autor executou os cinco passos do roteiro manual da spec na yazi 26.9.1, com uma configuração temporária (`YAZI_CONFIG_HOME`) e uma cópia da pasta de modelos, e tudo saiu conforme, inclusive o erro com `--hold` esperando o `Enter`. A tecla ficou fixada em [[decisions#DE-24 — O atalho do yazi é c seguido de m|DE-24]]. (registrada em 2026-11-09; resolvida em 2026-11-16)
- **DT-04** ✅ — **As ferramentas parecidas não foram estudadas.** Só os resultados de uma busca tinham sido lidos. O estudo está em [[similar-tools]]. As lições viraram [[decisions#DE-20 — O user-dirs.dirs é lido como a GLib lê|DE-20]], as dívidas DT-08, DT-09 e DT-10 e casos de teste dos requisitos. (registrada em 2026-11-09; resolvida em 2026-11-16)
- **DT-05** ✅ — **O remoto e o caminho do módulo Go não estavam definidos.** O remoto `origin` apareceu configurado como `git@gitlab.com:rg-dev/janus.git`, e o módulo segue o mesmo endereço: `gitlab.com/rg-dev/janus` ([[decisions#DE-16 — O módulo Go é gitlab.com/rg-dev/janus|DE-16]]). (registrada em 2026-11-09; resolvida em 2026-11-16)
- **DT-07** ✅ — **Os requisitos com ID ainda não tinham sido derivados.** Não existiam `FR` nem `NFR` em `docs/requirements/`. Derivados de [[bootstrap]] em [[core|núcleo]], [[cli|linha de comando]] e [[tui|modo interativo]]. Os módulos `NVIM` e `YAZI` ficam para [[roadmap|RM-01.03]] e [[roadmap|RM-01.02]]. (registrada em 2026-11-09; resolvida em 2026-11-16)
- **DT-11** ✅ — **A branch e o worktree de RM-01.01 fogem do padrão.** A branch `feat/rm-01.01_core-and-terminal` tem `.` e `_` no slug, e o worktree mora em `~/.herdr/worktrees/janus/feat-rm-01-01-core-and-terminal`, e não em `~/Projetos/dev/janus/feat_<slug>`, como pedia o `AGENTS.md`. A branch de RM-01.02 repetiu o formato. O autor adotou o formato para as entregas do roadmap em [[decisions#DE-25 — A branch de entrega do roadmap leva o ID e mora no worktree do herdr|DE-25]], e o `AGENTS.md` passou a descrevê-lo. (registrada em 2026-11-16; resolvida em 2026-11-18)
