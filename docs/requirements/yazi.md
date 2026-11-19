---
title: Requisitos da integração com o yazi
aliases:
  - FR-YAZI-NN
  - NFR-YAZI-NN
tags:
  - requisitos/funcionais
  - requisitos/nao-funcionais
---

# Requisitos da integração com o yazi (`YAZI`)

O que a integração com o yazi (`integrations/yazi/`) deve atender: um atalho que roda o modo interativo na pasta em que o yazi está. É o nível trivial de integração ([[bootstrap#Dois níveis de integração]]): o adaptador não tem código, só o trecho do `keymap.toml`. Derivados de [[bootstrap#Yazi: atalho]] e de [[bootstrap#Roteiro de verificação manual: yazi]] em 2026-11-18, no desenho de [[roadmap|RM-01.02]]. A verificação é manual, pelo roteiro. O desenho previa rodá-lo todo num terminal simulado pelo `tmux`, com `YAZI_CONFIG_HOME` apontando para `integrations/yazi/`, para que o arquivo conferido fosse o entregue. Mas o `shell --block` do yazi não roda no `tmux` ([[technical-debts|DT-16]]). A conferência foi então dividida: o `janus --hold` de cada passo e o menu do `c` no `tmux`, e o yazi de verdade no terminal do autor, com uma cópia byte a byte do arquivo entregue ([[rm-01.02-manual-check]]).

## Funcionais

- **FR-YAZI-01** ✅ — **O atalho cria o arquivo na pasta em que o yazi está.** `c` seguido de `m` roda `janus --hold` nessa pasta, de forma bloqueante, e o arquivo criado aparece na listagem ao voltar ([[decisions#DE-24 — O atalho do yazi é c seguido de m|DE-24]]). (registrada em 2026-11-18; atendido em 2026-11-19)
  - *Critério:* passos 1, 2 e 4 do roteiro: o modelo `fixed` nasce com o nome fixo, e o `ask`, com o nome confirmado e a extensão do modelo.
- **FR-YAZI-02** ✅ — **Conflito e cancelamento voltam ao yazi sem criar nada.** (registrada em 2026-11-18; atendido em 2026-11-19)
  - *Critério:* passo 3 do roteiro: o campo de nome abre com o aviso de destino existente, e cancelar devolve ao yazi com a pasta igual.
- **FR-YAZI-03** ✅ — **Um erro fica na tela até o `Enter`.** O yazi esconde a tela do comando assim que ele termina, e o `--hold` segura a mensagem ([[tui|FR-TUI-13]]). (registrada em 2026-11-18; atendido em 2026-11-19)
  - *Critério:* passos 5 e 6 do roteiro: com o manifesto inválido, a mensagem do código 4; dentro da pasta de modelos, a do código 10, sem criar nada ([[core|FR-CORE-24]]).

## Não funcionais

- **NFR-YAZI-01** ✅ — **O adaptador é só o trecho do `keymap.toml`.** Ele não percorre a pasta de modelos nem lê o manifesto: chama o modo interativo ([[decisions#DE-03 — Núcleo em linha de comando com adaptadores finos|DE-03]]). (registrada em 2026-11-18; atendido em 2026-11-19)
  - *Critério:* `integrations/yazi/` tem só `keymap.toml` e `README.md`, e o atalho chama apenas `janus --hold`.
- **NFR-YAZI-02** ✅ — **O atalho tem precedência sem tomar o lugar das outras teclas do `c`.** Usa `[[mgr.prepend_keymap]]`, que vem antes do mapa padrão, em que `c` abre as ações de copiar caminhos. (registrada em 2026-11-18; atendido em 2026-11-19)
  - *Critério:* passo 7 do roteiro: o menu do `c` mostra a entrada do Janus junto das de copiar caminhos.
  - *Nota:* o roteiro não executa `c` seguido de `c`, que copiaria um caminho para a área de transferência do autor; o menu já mostra que o mapa padrão continua lá.
