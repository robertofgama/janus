---
title: Requisitos do modo interativo
aliases:
  - FR-TUI-NN
  - NFR-TUI-NN
tags:
  - requisitos/funcionais
  - requisitos/nao-funcionais
---

# Requisitos do modo interativo (`TUI`)

O que o modo interativo (`internal/interactive`) deve atender: o comando `janus` sem subcomando, com o seletor, a prévia, o campo de nome e `--hold`. Derivados de [[bootstrap#Contratos e formatos]] e de [[bootstrap#Cenários de aceite]] em 2026-11-16 ([[technical-debts|DT-07]]). Os testes do laço usam substitutos do seletor e do campo de nome; os do adaptador do `fzf` usam um `fzf` falso no `PATH`, que registra as opções recebidas e devolve saídas e códigos combinados.

## Funcionais

### Entrada

- **FR-TUI-01** ⬜ — **O modo interativo exige um terminal.** Se a entrada padrão ou a saída de erro não forem um terminal, ele falha com o código 2 e sugere `janus new`. A saída padrão pode estar redirecionada. (registrada em 2026-11-16)
  - *Critério:* entrada padrão vinda de um arquivo dá código 2 e uma mensagem que cita `janus new`; a mesma execução com a saída padrão num arquivo e os outros dois num terminal segue adiante.
- **FR-TUI-02** ⬜ — **A ordem de verificação é 2, 3, 4, 5 e 9.** Depois vem o laço de escolha e nome. (registrada em 2026-11-16)
  - *Critério:* sem terminal e com pasta de modelos ruim dá 2; pasta de modelos ruim com manifesto inválido dá 3; manifesto inválido com pasta vazia dá 4; pasta sem modelos e sem `fzf` no `PATH` dá 5; sem `fzf` no `PATH` dá 9 (cenário 22).
- **FR-TUI-03** ⬜ — **Pasta sem nenhum modelo encerra com o código 5.** Com um aviso que cita a pasta. (registrada em 2026-11-16)
  - *Critério:* pasta de modelos vazia, ou só com itens ocultos, dá código 5 sem abrir o seletor.

### Laço de escolha e nome

- **FR-TUI-04** ⬜ — **Modelo `fixed` é criado sem pergunta na primeira tentativa.** (registrada em 2026-11-16)
  - *Critério:* escolhendo `info.ai/agents/react.AGENTS.md`, `AGENTS.md` é criado e o campo de nome não é chamado (cenário 18).
- **FR-TUI-05** ⬜ — **Modelo `ask` abre o campo com a sugestão.** O nome confirmado passa pelas regras de nome (FR-CORE-16 e FR-CORE-17). (registrada em 2026-11-16)
  - *Critério:* escolhendo a contestação, o campo recebe `01.2. Civ. Contestação.tex` como valor inicial, e confirmar `Fulano` cria `Fulano.tex` (cenário 19).
- **FR-TUI-06** ⬜ — **Destino existente reabre o campo com o nome que colidiu.** Com um aviso, inclusive para modelo `fixed`, repetindo até dar certo ou o usuário cancelar. (registrada em 2026-11-16)
  - *Critério:* com `AGENTS.md` já na pasta, escolher `react.AGENTS.md` abre o campo com `AGENTS.md` e um aviso de destino existente; confirmar `AGENTS.md` de novo reabre o campo; confirmar `OUTRO` cria `OUTRO.md`, e o `AGENTS.md` existente não muda (cenário 20).
- **FR-TUI-07** ⬜ — **Nome inválido reabre o campo com o que foi digitado.** Com um aviso. (registrada em 2026-11-16)
  - *Critério:* confirmar `a/b` reabre o campo com `a/b` e um aviso de nome inválido; confirmar `ok` em seguida cria o arquivo.
- **FR-TUI-08** ⬜ — **Cancelar encerra com o código 130, em silêncio.** No seletor ou no campo de nome, nada é criado e não aparece mensagem de erro. (registrada em 2026-11-16)
  - *Critério:* cancelar no seletor e cancelar no campo, cada um, dão código 130, pasta de destino vazia e saída de erro vazia (cenário 21).
- **FR-TUI-09** ⬜ — **A saída padrão recebe só o caminho criado.** O caminho absoluto, uma quebra de linha e nada mais; seletor, campo, avisos e erros ficam no terminal. É o que permite `nvim (janus)` no fish. (registrada em 2026-11-16)
  - *Critério:* nos cenários 18 a 20, a saída padrão é exatamente o caminho absoluto do arquivo criado; no 21, vazia.

### `fzf`

- **FR-TUI-10** ⬜ — **O seletor é o `fzf`.** Ele recebe um identificador por linha na entrada. Escolher segue adiante; cancelar (`Esc`, `Ctrl-C`) ou confirmar sem resultado devolve cancelamento (130); qualquer outra falha dá código 1 com a mensagem do `fzf`. (registrada em 2026-11-16)
  - *Critério:* com o `fzf` falso: a entrada recebida é a listagem, um id por linha; saída 0 com um id devolve o id; saídas 1 e 130 devolvem cancelamento; saída 2 com texto na saída de erro dá 1 e repassa o texto.
  - *Nota:* [[decisions#DE-06 — O seletor do terminal é o fzf|DE-06]].
- **FR-TUI-11** ⬜ — **O campo de nome é o `fzf` com lista vazia.** Ele abre com o valor inicial editável e o aviso no cabeçalho, quando houver. `Enter` confirma; `Esc`, `Ctrl-C` e `Ctrl-D` cancelam. (registrada em 2026-11-16)
  - *Critério:* com o `fzf` falso: as opções recebidas trazem `--print-query`, `--query` com o valor inicial, o aviso no cabeçalho e o atalho que faz `Ctrl-D` cancelar; a primeira linha da saída é devolvida como o nome, com a saída 0 ou 1; a saída 130 devolve cancelamento. Na conferência manual com o `fzf` real, as três teclas cancelam e um nome com acento e espaço volta intacto.
  - *Nota:* [[decisions#DE-17 — O campo de nome é o próprio fzf|DE-17]].
- **FR-TUI-12** ⬜ — **A prévia vem de um subcomando interno do próprio binário.** O `fzf` chama `janus __preview <id>`. Texto aparece com realce pelo `bat` quando ele está no `PATH`, e cru sem ele; binário mostra uma linha de aviso com o tamanho; vazio mostra uma linha de aviso. O subcomando não aparece na ajuda e não é contrato. (registrada em 2026-11-16)
  - *Critério:* sem `bat` no `PATH`, a prévia de um modelo de texto é idêntica ao arquivo; com um `bat` falso, ele recebe o caminho do modelo como um argumento só; `3d/fatiar.3mf` e `Empty File` dão uma linha de aviso cada, sem conteúdo cru (cenário 23); `info.dev/react; script; simple.html` chega intacto ao subcomando.
  - *Nota:* [[decisions#DE-12 — A prévia é produzida pelo próprio binário|DE-12]].

### `--hold`

- **FR-TUI-13** ⬜ — **`--hold` espera um `Enter` depois de um erro.** Só quando o código não é 0 nem 130, depois de mostrar a mensagem. Existe para programas que escondem a tela assim que o comando termina, como o yazi. (registrada em 2026-11-16)
  - *Critério:* com `--hold` e manifesto inválido, o código 4 só é devolvido depois de ler uma linha da entrada; com `--hold` e cancelamento, nada é lido; sem `--hold`, nada é lido em nenhum caso.

## Não funcionais

- **NFR-TUI-01** ⬜ — **A prévia funciona com qualquer shell do usuário.** O `fzf` roda a prévia com `sh -c`, e não com o `$SHELL`, que na máquina do autor é o fish, com outras regras de aspas. O caminho do binário vai entre aspas simples. (registrada em 2026-11-16)
  - *Critério:* com o `fzf` falso, as opções trazem `--with-shell` com `sh -c` e a prévia com o caminho do binário entre aspas; na conferência manual, a prévia funciona com `SHELL=fish` e com o binário numa pasta com espaço no nome.
- **NFR-TUI-02** ⬜ — **O `FZF_DEFAULT_OPTS` do usuário é respeitado, sem quebrar o contrato.** O tema e o layout do usuário valem, mas as opções de que o Janus depende vêm depois e prevalecem. (registrada em 2026-11-16)
  - *Critério:* com o `fzf` falso, as opções do seletor incluem `--no-multi`, `--no-print-query` e `--no-expect`; na conferência manual, `FZF_DEFAULT_OPTS='--multi --print-query --expect=ctrl-x'` não muda o resultado do seletor.
- **NFR-TUI-03** ⬜ — **O laço não depende do `fzf`.** O seletor e o campo de nome ficam atrás de interfaces, e o modo interativo chama o núcleo no mesmo processo, nunca o próprio binário. (registrada em 2026-11-16)
  - *Critério:* os testes dos cenários 18 a 21 rodam sem `fzf` no `PATH`.
