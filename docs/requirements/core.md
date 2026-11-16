---
title: Requisitos do núcleo
aliases:
  - FR-CORE-NN
  - NFR-CORE-NN
tags:
  - requisitos/funcionais
  - requisitos/nao-funcionais
---

# Requisitos do núcleo (`CORE`)

O que o núcleo (`internal/core`) deve atender: descobrir a pasta de modelos, listar, ler o manifesto, resolver o nome final e criar o arquivo. Derivados de [[bootstrap#Contratos e formatos]] e de [[bootstrap#Cenários de aceite]] em 2026-11-16 ([[technical-debts|DT-07]]). Os cenários citados como "cenário NN" são os da spec, e a "pasta de teste" é a pasta de modelos de teste descrita lá.

## Funcionais

### Descoberta da pasta de modelos

- **FR-CORE-01** ⬜ — **A variável de ambiente manda.** Se `XDG_TEMPLATES_DIR` está definida e não é vazia, ela é a pasta de modelos, sem segunda tentativa. (registrada em 2026-11-16)
  - *Critério:* com a variável apontando para a pasta de teste, a listagem traz os modelos dela; apontando para uma pasta inexistente, a descoberta falha com o erro de pasta de modelos, mesmo com um `user-dirs.dirs` válido (cenário 17). Com a variável vazia, vale o `user-dirs.dirs`.
- **FR-CORE-02** ⬜ — **Sem a variável, vale o `user-dirs.dirs`, lido como a GLib lê.** O arquivo é `$XDG_CONFIG_HOME/user-dirs.dirs`, ou `~/.config/user-dirs.dirs` sem `XDG_CONFIG_HOME`. A linha aceita espaços no início e em volta do `=`; vale a última linha que define `XDG_TEMPLATES_DIR`; o valor vai da primeira à última aspa, sem desfazer escapes; `$HOME` seguido de `/` ou da aspa final é expandido; barras no fim são removidas. (registrada em 2026-11-16)
  - *Critério:* cada regra tem um caso: `XDG_TEMPLATES_DIR="$HOME/Modelos"`; `  XDG_TEMPLATES_DIR = "/abs/x"`; duas linhas, vencendo a segunda; valor com `\"` no meio, mantido literal; `"$HOME/Modelos///"`; caminho absoluto; `XDG_CONFIG_HOME` ausente, caindo em `~/.config`. Os testes apontam `HOME` e `XDG_CONFIG_HOME` para pastas temporárias.
  - *Nota:* [[decisions#DE-20 — O user-dirs.dirs é lido como a GLib lê|DE-20]].
- **FR-CORE-03** ⬜ — **A descoberta falha dizendo o que tentou.** O erro de pasta de modelos traz o motivo e o item tentado (variável, arquivo, linha ou caminho) quando o arquivo não existe, a linha não existe, o valor é relativo (inclui `$HOME` seguido de algo que não é `/`), o caminho não existe ou não é pasta, ou o caminho é a pasta pessoal, que no padrão XDG significa pasta desativada. (registrada em 2026-11-16)
  - *Critério:* um caso por motivo, cada um com o motivo certo no erro; `"$HOME"`, `"$HOME/"` e o caminho absoluto da pasta pessoal dão o motivo "pasta pessoal" (cenário 17); `"$HOMEx"` e `"Modelos"` dão "relativo".
- **FR-CORE-04** ⬜ — **A pasta de modelos pode ser um link simbólico.** O link é seguido, como `~/Modelos` na máquina do autor. (registrada em 2026-11-16)
  - *Critério:* com a variável apontando para um link para a pasta de teste, a listagem é igual à da pasta original.

### Listagem

- **FR-CORE-05** ⬜ — **Todo arquivo regular é um modelo.** Em qualquer profundidade, inclusive os vazios. O identificador é o caminho relativo à raiz, com `/`; a categoria é o caminho da pasta que o contém, vazia na raiz; o arquivo é o nome dele. (registrada em 2026-11-16)
  - *Critério:* a pasta de teste dá nove modelos (cenário 01); `info.ai/agents/react.AGENTS.md` tem categoria `info.ai/agents` e arquivo `react.AGENTS.md`; `Empty Bash` tem categoria vazia (cenário 02); `Empty File`, com 0 bytes, está na lista.
- **FR-CORE-06** ⬜ — **Itens ocultos são ignorados em qualquer nível.** Arquivo ou pasta cujo nome começa com `.` não entra na listagem, nem nada dentro dela. (registrada em 2026-11-16)
  - *Critério:* `.git/config` e `.janus.toml` não aparecem (cenário 01); um `info.dev/.oculto` e um `.pasta/dentro.txt` acrescentados no teste também não.
- **FR-CORE-07** ⬜ — **Links simbólicos dentro da pasta seguem regras próprias.** Um link para arquivo é um modelo, e o que se copia é o arquivo apontado; um link quebrado é ignorado; um link para pasta não é seguido. (registrada em 2026-11-16)
  - *Critério:* com um link para arquivo, um link quebrado e um link para pasta acrescentados à pasta de teste, só o primeiro aparece na listagem, e criar a partir dele copia o conteúdo do arquivo apontado.
- **FR-CORE-08** ⬜ — **Nomes com caractere de controle são ignorados.** Arquivo ou pasta cujo nome contém quebra de linha ou outro caractere de controle não entra na listagem. (registrada em 2026-11-16)
  - *Critério:* um arquivo `a\nb` e um arquivo com tabulação no nome, acrescentados no teste, não aparecem.
- **FR-CORE-09** ⬜ — **A listagem é ordenada por identificador, sem diferenciar maiúsculas.** Empates são decididos pela ordem dos bytes. (registrada em 2026-11-16)
  - *Critério:* a pasta de teste sai na ordem `3d/…`, `direito/…`, `Empty Bash`, `Empty File`, `info.ai/…`, `info.dev/…` (cenário 01); `B` e `b` saem nessa ordem, por ser a ordem dos bytes.

### Manifesto

- **FR-CORE-10** ⬜ — **O manifesto é opcional.** Arquivo ausente, arquivo vazio ou tabela `[names]` ausente significam nenhuma regra: todos os modelos ficam em modo `ask`. (registrada em 2026-11-16)
  - *Critério:* um caso para cada uma das três situações, com todos os modelos em `ask` e a sugestão igual ao nome do arquivo.
- **FR-CORE-11** ⬜ — **O manifesto é validado com rigor.** A única tabela aceita é `[names]`, com chaves que são padrões e valores que são textos. Qualquer outra chave ou tabela, um valor que não é texto ou um erro de sintaxe tornam o manifesto inválido. (registrada em 2026-11-16)
  - *Critério:* chave solta fora de `[names]`, tabela `[outra]`, subtabela `[names.x]`, valor numérico e TOML malformado dão manifesto inválido, cada um citando a chave ou a linha (cenário 16).
  - *Nota:* [[decisions#DE-18 — O manifesto é lido com BurntSushi/toml|DE-18]].
- **FR-CORE-12** ⬜ — **O padrão casa o identificador inteiro, com `*` como único curinga.** `*` casa qualquer sequência, inclusive vazia, sem atravessar `/`; os demais caracteres, inclusive `?`, `[`, `]` e `\`, são literais. `**` em qualquer posição torna o manifesto inválido. Padrão que não casa nenhum modelo não é erro. (registrada em 2026-11-16)
  - *Critério:* tabela de casos: `info.ai/agents/*` casa `info.ai/agents/AGENTS.md` e não casa `info.ai/agents/x/y`; `*` não casa `a/b`; `a*` casa `a`; `info.dev/package.json.*` casa `info.dev/package.json.monorepo`; `a?c` não casa `abc`; `[ab]` só casa `[ab]`; `a**`, `**/x` e `a/**/b` dão manifesto inválido (cenário 16); um padrão sem casamento deixa o manifesto válido.
  - *Nota:* `path.Match` não serve, porque interpreta `?`, `[` e `\`. Um `*` literal não tem como ser escrito ([[technical-debts|DT-09]]).
- **FR-CORE-13** ⬜ — **A regra de cada modelo é resolvida sem ambiguidade.** Um padrão sem `*` igual ao identificador vence, não importa quantos padrões com `*` também casem. Senão, nenhum padrão com `*` casando deixa o modelo em `ask`, exatamente um define o nome fixo, e dois ou mais tornam o manifesto inválido, citando os padrões e o identificador. (registrada em 2026-11-16)
  - *Critério:* as três regras do exemplo da spec dão nome fixo aos modelos certos (cenário 02); `"info.ai/agents/*"` com `"info.ai/*/react.*"` dão manifesto inválido citando os dois padrões e `info.ai/agents/react.AGENTS.md` (cenário 14); somando a regra exata `"info.ai/agents/react.AGENTS.md" = "X.md"`, o manifesto fica válido e o nome fixo é `X.md` (cenário 15).
- **FR-CORE-14** ⬜ — **O nome fixo passa pela validação de nome.** Ele tem as pontas aparadas, e um nome fixo inválido (FR-CORE-16) torna o manifesto inválido, citando o padrão e o nome. (registrada em 2026-11-16)
  - *Critério:* `= "a/b"`, `= ".."` e `= "  "` dão manifesto inválido; `= "  AGENTS.md "` dá o nome fixo `AGENTS.md`.

### Busca

- **FR-CORE-15** ⬜ — **O identificador é procurado na listagem.** A busca nunca junta o identificador à pasta para formar um caminho: o que a listagem não traz não é modelo. (registrada em 2026-11-16)
  - *Critério:* `nao/existe` (cenário 13), `../x`, `/etc/hostname`, `.janus.toml`, `.git/config`, `./Empty Bash` e `empty bash` dão modelo não encontrado.
  - *Nota:* lição do [[similar-tools|estudo das ferramentas parecidas]].

### Nome final

- **FR-CORE-16** ⬜ — **Todo nome é aparado e validado.** Espaços nas pontas são removidos. O nome é inválido se ficar vazio, se for `.` ou `..`, ou se contiver `/` ou caractere de controle. Vale para o nome fixo, o digitado e o de `--name`. (registrada em 2026-11-16)
  - *Critério:* `a/b`, `..` e `  ` são inválidos (cenário 11); `.`, `a\tb` e `a\nb` também; `  x.md  ` vira `x.md`.
- **FR-CORE-17** ⬜ — **O nome informado ganha a extensão de referência quando ela falta.** A extensão de referência é a do nome fixo, se o modelo é `fixed`, ou a do arquivo do modelo, se é `ask`. Extensão é o trecho a partir do último ponto, se o ponto não é o primeiro caractere, há algo depois dele e o trecho não tem espaço. Se o nome não termina com ela, sem diferenciar maiúsculas, ela é acrescentada; sem extensão de referência, o nome fica como veio. O nome informado substitui o nome fixo. (registrada em 2026-11-16)
  - *Critério:* os sete exemplos da tabela de [[bootstrap#Resolução do nome final]]; cenários 06 a 10; `.ai-memory.toml` tem extensão `.toml`, e `.bashrc` não tem extensão.
  - *Nota:* extensão composta perde a primeira parte ([[technical-debts|DT-08]]).

### Criação

- **FR-CORE-18** ⬜ — **A pasta de destino precisa existir e ser pasta.** Um link simbólico para pasta é aceito. (registrada em 2026-11-16)
  - *Critério:* pasta inexistente (cenário 12) e arquivo comum dão erro de pasta de destino citando o caminho; um link para pasta é aceito, e o arquivo nasce na pasta apontada.
- **FR-CORE-19** ⬜ — **A criação é exclusiva.** Se já existe qualquer coisa com o nome final na pasta de destino (arquivo, pasta ou link, mesmo quebrado), a criação falha com o erro de destino existente, citando o caminho, e nada muda. Se o arquivo surgir entre a checagem e a escrita, o resultado é o mesmo. (registrada em 2026-11-16)
  - *Critério:* repetir a criação do cenário 03 falha e o arquivo existente fica igual (cenário 04); uma pasta e um link quebrado com o nome final também fazem falhar.
  - *Nota:* a exclusividade vem de abrir o destino com criação exclusiva (`O_EXCL`), não de uma checagem antes.
- **FR-CORE-20** ⬜ — **O conteúdo é copiado byte a byte.** Sem variáveis nem conversão ([[decisions#DE-05 — Conteúdo estático, sem variáveis|DE-05]]). (registrada em 2026-11-16)
  - *Critério:* o arquivo criado é idêntico ao modelo nos cenários 03 e 10, em `3d/fatiar.3mf`, com bytes nulos, e num modelo com quebras de linha CRLF; `Empty File` cria um arquivo de 0 bytes.
- **FR-CORE-21** ⬜ — **O arquivo novo tem exatamente as permissões do modelo.** Sem influência da `umask`. Datas não são preservadas. (registrada em 2026-11-16)
  - *Critério:* com `umask` 077, `Empty Bash` (755) cria um arquivo 755 (cenário 08), e um modelo 640 cria um arquivo 640.
- **FR-CORE-22** ⬜ — **Uma cópia que falha no meio não deixa arquivo.** O arquivo parcial é removido, e o erro devolvido não é de nenhum dos tipos com código próprio. (registrada em 2026-11-16)
  - *Critério:* com uma falha de leitura simulada no meio da cópia, a pasta de destino fica como estava.

### Prévia

- **FR-CORE-23** ⬜ — **O núcleo classifica o modelo para a prévia.** Arquivo de 0 bytes é vazio; arquivo com um byte nulo nos primeiros 8000 bytes é binário, com o tamanho; o resto é texto. (registrada em 2026-11-16)
  - *Critério:* `Empty File` é vazio, `3d/fatiar.3mf` é binário com o tamanho certo, `info.dev/react; script; simple.html` é texto; um arquivo com o primeiro byte nulo na posição 8001 é texto.

## Não funcionais

- **NFR-CORE-01** ⬜ — **O Janus nunca escreve na pasta de modelos.** Nem cria, nem altera, nem remove modelos ou o manifesto ([[decisions#DE-02 — A pasta XDG de modelos é a fonte única, compartilhada com o Nemo|DE-02]]). (registrada em 2026-11-16)
  - *Critério:* nos testes de criação, o conteúdo, as permissões e as datas de todos os itens da pasta de teste são iguais antes e depois.
- **NFR-CORE-02** ⬜ — **O núcleo não faz entrada e saída de terminal.** Ele não imprime, não lê do teclado e não conhece códigos de saída: devolve valores e erros tipados, um tipo por código de saída. (registrada em 2026-11-16)
  - *Critério:* um teste confere que nenhum arquivo de `internal/core`, fora os de teste, usa `os.Stdin`, `os.Stdout`, `os.Stderr`, `os.Exit` ou as funções `fmt.Print`, nem importa outro pacote de `internal/`.
