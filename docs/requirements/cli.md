---
title: Requisitos da linha de comando
aliases:
  - FR-CLI-NN
  - NFR-CLI-NN
tags:
  - requisitos/funcionais
  - requisitos/nao-funcionais
---

# Requisitos da linha de comando (`CLI`)

O que a linha de comando (`internal/cli`) deve atender: ler os argumentos, rodar `list` e `new`, imprimir ajuda e versão e traduzir erros em códigos de saída. É o contrato de que os adaptadores dependem. Derivados de [[bootstrap#Contratos e formatos]] e de [[bootstrap#Cenários de aceite]] em 2026-11-16 ([[technical-debts|DT-07]]). Os testes rodam `cli.Run` no mesmo processo, com a pasta de teste apontada por `XDG_TEMPLATES_DIR`, e conferem o código de saída e o item que a mensagem cita, nunca a frase.

## Funcionais

### Argumentos

- **FR-CLI-01** ✅ — **Só as formas do contrato são aceitas.** `janus [--hold]`, `janus list [--json]`, `janus new <id> [--dir <pasta>] [--name <nome>]`, `janus --help`, `janus <subcomando> --help` e `janus --version`. Qualquer outro argumento, subcomando ou opção é uso incorreto, com o código 2. (registrada em 2026-11-16; atendido em 2026-11-16)
  - *Critério:* dão código 2: `janus foo`, `janus list x`, `janus list --dir x`, `janus new` sem identificador, `janus new a b`, `janus new a --name`, sem valor, `janus new a --json`, `janus list --hold`, `janus --hold --hold`, `janus new a --dir x --dir y`, `janus new a -name x`, `janus -h`, `janus new a --help` e `janus --version x`.
  - *Nota:* as regras que a spec não escreveu saem de "qualquer outro argumento é uso incorreto": opção repetida, opção com um hífen só e `--help` ou `--version` acompanhados de outros argumentos são uso incorreto; `--hold` só existe sem subcomando.
- **FR-CLI-02** ✅ — **Opções aceitam as duas formas e podem vir em qualquer posição.** `--opcao valor` e `--opcao=valor`, antes ou depois do identificador. `--` encerra as opções, e o que vem depois é posicional. (registrada em 2026-11-16; atendido em 2026-11-16)
  - *Critério:* `new <id> --dir=X`, `new --dir X <id>` e `new <id> --name=a=b` funcionam, o último com o nome `a=b`; com um modelo `-x.md` na pasta, `new --name y -- -x.md` cria `y.md`, e `new -x.md` dá código 2.
  - *Nota:* [[decisions#DE-19 — Os argumentos são lidos por código próprio, sem o pacote flag|DE-19]]. Um valor vazio (`--name=`, `--dir ""`) não é uso incorreto: segue para a validação do nome (7) ou da pasta (10).
- **FR-CLI-03** ✅ — **A ajuda vai para a saída padrão, com código 0.** `janus --help` descreve o comando inteiro, e `janus list --help` e `janus new --help`, cada subcomando, em português. O subcomando interno de prévia não aparece. (registrada em 2026-11-16; atendido em 2026-11-16)
  - *Critério:* as três formas saem com 0, com a saída de erro vazia e a ajuda na saída padrão; nenhuma cita `__preview`.
- **FR-CLI-04** ✅ — **`janus --version` imprime a versão.** Na saída padrão, com código 0. A versão é a que o Go grava no binário, ou `(devel)` sem ela. (registrada em 2026-11-16; atendido em 2026-11-16)
  - *Critério:* no teste, que não tem versão gravada, a saída é `(devel)` com quebra de linha; o teste de ponta a ponta confere que o binário compilado imprime uma linha não vazia.
  - *Nota:* [[decisions#DE-22 — A versão vem do que o Go grava no binário|DE-22]].

### `list`

- **FR-CLI-05** ✅ — **`list` imprime um identificador por linha.** Pasta sem modelos não imprime nada. O código é 0 nos dois casos. (registrada em 2026-11-16; atendido em 2026-11-16)
  - *Critério:* a pasta de teste dá as nove linhas de FR-CORE-09, na mesma ordem (cenário 01); uma pasta vazia dá saída vazia e código 0.
- **FR-CLI-06** ✅ — **`list --json` imprime a lista JSON do contrato.** Uma lista em UTF-8, na ordem da listagem, com `id`, `category`, `file` e `name.mode`, mais `name.value` só em `fixed` e `name.suggestion`, igual a `file`, só em `ask`. Pasta sem modelos imprime `[]`. (registrada em 2026-11-16; atendido em 2026-11-16)
  - *Critério:* `react.AGENTS.md` tem `mode` `fixed` e `value` `AGENTS.md`, sem `suggestion`; a contestação tem `mode` `ask` e `suggestion` igual ao nome do arquivo, sem `value`; `Empty Bash` tem `category` `""` (cenário 02); acentos saem como UTF-8, e não como `\u`; pasta vazia dá `[]`.

### `new`

- **FR-CLI-07** ✅ — **`new` cria o arquivo e imprime só o caminho absoluto.** Na pasta de `--dir`, ou na atual sem a opção. A saída padrão recebe o caminho absoluto do arquivo criado, uma quebra de linha e nada mais. (registrada em 2026-11-16; atendido em 2026-11-16)
  - *Critério:* o cenário 03 cria `AGENTS.md` idêntico ao modelo e imprime só o caminho absoluto, com código 0; com `--dir` relativo, o caminho impresso é absoluto.
- **FR-CLI-08** ✅ — **`new` nunca pergunta.** Modelo `ask` sem `--name` falha com o código 6, citando o identificador, e nada é criado. (registrada em 2026-11-16; atendido em 2026-11-16)
  - *Critério:* cenário 05.
- **FR-CLI-09** ✅ — **`--name` passa pelas regras de nome.** O nome é validado e ganha a extensão de referência quando ela falta (FR-CORE-16 e FR-CORE-17); num modelo `fixed`, substitui o nome fixo. (registrada em 2026-11-16; atendido em 2026-11-16)
  - *Critério:* cenários 06 a 11; `new "info.ai/agents/react.AGENTS.md" --name CLAUDE` cria `CLAUDE.md`.

### Erros

- **FR-CLI-10** ✅ — **Cada erro tem o seu código de saída.** 1 inesperado, 2 uso incorreto, 3 pasta de modelos, 4 manifesto, 5 modelo não encontrado, 6 nome obrigatório, 7 nome inválido, 8 destino existente e 10 pasta de destino. Os códigos 9 e 130 são do modo interativo. (registrada em 2026-11-16; atendido em 2026-11-16)
  - *Critério:* um caso por código: os cenários 04 (8), 05 (6), 11 (7), 12 (10), 13 (5), 14 e 16 (4, em `list` e em `new`) e 17 (3), mais uma falha de escrita na pasta de destino (1).
- **FR-CLI-11** ✅ — **Quando mais de um erro se aplica, vale a ordem de verificação.** Em `list`: 2, 3, 4, 1. Em `new`: 2, 3, 4, 5, 10, 6, 7, 8, 1. (registrada em 2026-11-16; atendido em 2026-11-16)
  - *Critério:* cada par vizinho tem um caso: uso incorreto com pasta de modelos ruim dá 2; pasta de modelos ruim com manifesto inválido dá 3; manifesto inválido com identificador inexistente dá 4; identificador inexistente com `--dir` inexistente, ou com `--dir` na pasta de modelos, dá 5; `--dir` inexistente, ou na pasta de modelos, com modelo `ask` sem `--name` dá 10; e, com `--name "a/b"` e um `a` já existente na pasta de destino, dá 7.
  - *Nota:* os casos com `--dir` na pasta de modelos entraram em 2026-11-18, com [[core|FR-CORE-24]], em [[roadmap|RM-01.02]].
- **FR-CLI-12** ✅ — **Toda falha explica, na saída de erro, o que aconteceu.** A mensagem é em português e cita o item envolvido (caminho, identificador, padrão ou nome); a saída padrão fica vazia. (registrada em 2026-11-16; atendido em 2026-11-16)
  - *Critério:* em cada caso de FR-CLI-10, a saída padrão está vazia e a saída de erro contém o item citado: o caminho em 3, 8 e 10, o padrão ou a chave em 4, o identificador em 5 e 6 e o nome em 7.

## Não funcionais

- **NFR-CLI-01** ✅ — **A saída padrão leva só o resultado.** Lista, JSON, caminho criado, ajuda ou versão. Avisos e erros vão para a saída de erro. É o que permite `nvim (janus …)` no fish. (registrada em 2026-11-16; atendido em 2026-11-16)
  - *Critério:* em todo teste da linha de comando, a saída padrão é exatamente o resultado esperado, ou vazia na falha.
- **NFR-CLI-02** ✅ — **`list` e `new` não dependem de `fzf` nem de `bat`.** (registrada em 2026-11-16; atendido em 2026-11-16)
  - *Critério:* com `PATH` apontando para uma pasta vazia, os cenários 01 a 03 passam.
- **NFR-CLI-03** ✅ — **As mensagens ficam num lugar só.** Todo texto ao usuário, inclusive a ajuda, é em português do Brasil e mora em `internal/messages`. (registrada em 2026-11-16; atendido em 2026-11-16)
  - *Critério:* revisão do código: nenhum texto ao usuário fora de `internal/messages`.
  - *Nota:* [[decisions#DE-11 — Contrato em inglês, mensagens em português|DE-11]] e [[decisions#DE-21 — As mensagens moram em internal/messages|DE-21]].
- **NFR-CLI-04** ✅ — **A única dependência externa é a do TOML.** (registrada em 2026-11-16; atendido em 2026-11-16)
  - *Critério:* o `go.mod` tem um só `require` direto, `github.com/BurntSushi/toml`.
  - *Nota:* [[decisions#DE-18 — O manifesto é lido com BurntSushi/toml|DE-18]].
- **NFR-CLI-05** ✅ — **O binário funciona de ponta a ponta.** O `main` só repassa os argumentos, a entrada e as saídas do processo para `cli.Run` e encerra com o código devolvido. (registrada em 2026-11-16; atendido em 2026-11-16)
  - *Critério:* um teste compila `cmd/janus` e confere `list`, `new` e um código de erro rodando o binário de verdade.
