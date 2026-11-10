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
- **DT-04** ⬜ — **As ferramentas parecidas não foram estudadas.** Só os resultados de uma busca foram lidos. (registrada em 2026-11-09)
  - *Rever quando:* ao escrever o plano de [[roadmap|RM-01.01]]. Não bloqueia.
  - *Ver também:* [[bootstrap#Referências]].
- **DT-05** ⬜ — **O remoto e o caminho do módulo Go não estão definidos.** O autor descartou a proposta `gitlab.com/rg-dev/janus` e vai informar o destino. Sem isso não há `go mod init` nem `git remote add`. (registrada em 2026-11-09)
  - *Rever quando:* antes de começar o código de [[roadmap|RM-01.01]].
- **DT-06** ⬜ — **O catálogo não declara a stack, e os padrões da casa não têm referência para Go.** O campo `stack` de `catalog-info.yaml` está vazio, e a lista de referências por stack dos padrões só cobre TypeScript, NestJS, React, React Native e TailwindCSS. As convenções deste projeto estão em [[decisions#DE-13 — Código Go segue as convenções do Go; o resto, as da casa|DE-13]]. (registrada em 2026-11-09)
  - *Rever quando:* quando a casa ganhar uma referência para Go, ou no primeiro projeto Go seguinte.
- **DT-07** ⬜ — **Os requisitos com ID ainda não foram derivados.** Não existem `FR` nem `NFR` em `docs/requirements/`. A fonte são as seções "Contratos e formatos" e "Cenários de aceite" de [[bootstrap]]. Módulos sugeridos: `CORE`, `CLI`, `TUI`, `NVIM` e `YAZI`. (registrada em 2026-11-09)
  - *Rever quando:* ao escrever o plano de [[roadmap|RM-01.01]], antes do primeiro teste.

## Finalizadas

- **DT-01** ✅ — **O manifesto oculto podia aparecer no menu do Nemo.** A hipótese era que o Nemo ignora itens ocultos no menu "Criar novo documento". O autor confirmou que a pasta oculta `.git`, que mora na pasta de modelos, não aparece no menu. O manifesto fica em `.janus.toml`, na raiz da pasta de modelos, como desenhado. Se o arquivo aparecer no menu quando for criado, abrir uma dívida nova: a saída prevista é movê-lo para `$XDG_CONFIG_HOME/janus/janus.toml`, como único local aceito. (registrada em 2026-11-09; resolvida em 2026-11-09)
