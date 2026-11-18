---
title: Decisões
aliases:
  - DE-NN
tags:
  - decisoes
---

# Decisões

Decisões do Janus, cada uma com o motivo e as alternativas rejeitadas. As de DE-01 a DE-12 foram tomadas na sessão de desenho de 2026-11-09; DE-13 e DE-14, na migração para este repositório, e DE-15, na escrita do `AGENTS.md`, no mesmo dia. As de DE-16 a DE-22 saíram do planejamento de [[roadmap|RM-01.01]], em 2026-11-16, e DE-23 e DE-24, da revisão que fechou essa entrega, no mesmo dia. DE-25 e DE-26 saíram do desenho de [[roadmap|RM-01.02]], em 2026-11-18. O desenho completo está em [[bootstrap]].

## Registro

### DE-01 — O escopo é só modelos de arquivo

- **Decisão:** o Janus cria arquivos inteiros a partir de modelos. Não gerencia snippets.
- **Por quê:** modelo é um arquivo criado numa pasta; snippet é um trecho inserido dentro de um arquivo ou da linha de comando. São problemas diferentes, com integrações diferentes, principalmente no neovim.
- **Rejeitadas:** uma ferramenta única para modelos e snippets. Dobraria o escopo antes de resolver a dor original. Snippets seguem como ideia separada na incubadora.

### DE-02 — A pasta XDG de modelos é a fonte única, compartilhada com o Nemo

- **Decisão:** o Janus lê a mesma pasta que o Nemo usa. Não renomeia, não reorganiza e não altera modelos, e nunca escreve nessa pasta. O único arquivo novo ali é o manifesto oculto, criado e mantido pelo autor.
- **Por quê:** um só lugar para manter, e o menu do botão direito continua funcionando igual.
- **Rejeitadas:** reorganizar os modelos para seguir uma convenção de nomes, o que mudaria o que aparece no menu do Nemo; migrar para um formato próprio, o que deixaria o Nemo para trás.

### DE-03 — Núcleo em linha de comando com adaptadores finos

- **Decisão:** um único binário concentra tudo o que se sabe sobre modelos e expõe um contrato pequeno. Cada programa se liga a ele por um adaptador de poucas linhas.
- **Por quê:** as regras ficam em um lugar só, e um programa novo custa apenas um adaptador.
- **Rejeitadas:** colar ferramentas prontas (um utilitário de terminal, um plugin de neovim, um atalho no yazi), o que daria três comportamentos diferentes, sem regra de nome fixo e sem peça para o yazi; uma convenção compartilhada com implementação nativa em cada programa, o que escreveria a lógica de nomes e de conflito três vezes.

### DE-04 — Nome perguntado por padrão, nome fixo por manifesto oculto

- **Decisão:** sem configuração, todo modelo pede o nome, sugerindo o nome do próprio modelo. Nomes fixos vêm de `.janus.toml`, na raiz da pasta de modelos.
- **Por quê:** o padrão funciona em qualquer pasta de modelos, e a regra de nome fixo descreve o modelo, então viaja com ele no mesmo repositório git.
- **Rejeitadas:** deduzir o nome fixo do nome do arquivo, porque os modelos atuais não seguem um padrão único (a variante é prefixo em `react.AGENTS.md` e sufixo em `package.json.monorepo`) e [[#DE-02 — A pasta XDG de modelos é a fonte única, compartilhada com o Nemo|DE-02]] proíbe renomeá-los; guardar as regras em `~/.config`, porque elas se separariam dos modelos.

### DE-05 — Conteúdo estático, sem variáveis

- **Decisão:** o modelo é copiado byte a byte. Não há marcadores de substituição.
- **Por quê:** o Nemo copiaria os marcadores crus, o que quebraria [[#DE-02 — A pasta XDG de modelos é a fonte única, compartilhada com o Nemo|DE-02]].
- **Rejeitadas:** variáveis com sintaxe própria, que exigiriam modelos separados para o Nemo e para o Janus.

### DE-06 — O seletor do terminal é o fzf

- **Decisão:** o modo interativo delega a escolha do modelo ao `fzf`. Só o modo interativo depende dele.
- **Por quê:** já está instalado e em uso pelo autor; busca aproximada, prévia e atalhos de teclado vêm prontos.
- **Rejeitadas:** um seletor embutido no binário, que traria uma biblioteca de interface de terminal para refazer o que o `fzf` já faz.

### DE-07 — No neovim, criar o arquivo e abrir

- **Decisão:** no MVP, o comando do neovim cria um arquivo novo a partir do modelo e o abre. Preencher o buffer atual com o conteúdo de um modelo é a primeira entrega depois do MVP ([[roadmap|RM-02.01]]).
- **Por quê:** é o mesmo comportamento do terminal, então o adaptador fica mínimo.
- **Rejeitadas:** entregar as duas formas no MVP.

### DE-08 — Projeto próprio, de uso pessoal primeiro

- **Decisão:** repositório próprio, comando instalável e testado. Feito para o ambiente do autor, mas sem depender dele: a pasta é descoberta pelo padrão XDG e nada é fixo no código.
- **Por quê:** publicar depois vira uma decisão barata.
- **Rejeitadas:** scripts em fish ou bash dentro dos dotfiles, frágeis com os nomes de arquivo do inventário e difíceis de testar; projeto público desde o início, com documentação para terceiros e pacote no AUR antes de existir algo funcionando.

### DE-09 — Go

- **Decisão:** o Janus é escrito em Go.
- **Por quê:** binário único e estático, partida instantânea (o comando roda a cada uso, inclusive chamado pelo neovim), biblioteca padrão que cobre quase tudo e testes embutidos.
- **Rejeitadas:** TypeScript com Bun, a linguagem dos outros projetos do autor, por gerar binário grande e partida mais lenta ou depender do Bun instalado; Python, por empacotamento e partida mais lenta; Rust, por mais esforço de escrita sem ganho para um programa deste tamanho.

### DE-10 — O nome é Janus

- **Decisão:** o comando é `janus`, o manifesto é `.janus.toml` e o comando do neovim é `:Janus`.
- **Por quê:** escolha do autor. O comando estava livre no `PATH` em 2026-11-09, e não há pacote `janus` nos repositórios oficiais do Manjaro.
- **Rejeitadas:** `molde`, `modelo`, `novo` e `tpl`.

### DE-11 — Contrato em inglês, mensagens em português

- **Decisão:** subcomandos, opções, chaves do JSON e chaves do manifesto são em inglês. Mensagens ao usuário (erros, perguntas, ajuda) são em português do Brasil, concentradas em um único lugar do código.
- **Por quê:** é o padrão de idioma da casa: a máquina e o schema leem inglês, o operador lê português. Concentrar as mensagens barateia uma tradução futura.
- **Rejeitadas:** a tabela do manifesto em português (`[nomes]`), que apareceu no primeiro rascunho e misturava idiomas no mesmo contrato.

### DE-12 — A prévia é produzida pelo próprio binário

- **Decisão:** a prévia mostrada pelo `fzf` vem de um subcomando interno do Janus, que recebe o identificador e imprime o que deve aparecer.
- **Por quê:** mantém a regra de que só o núcleo transforma identificador em arquivo, evita montar comandos de shell com nomes que têm espaço e `;`, e permite tratar binários e arquivos vazios de forma testável.
- **Rejeitadas:** passar o caminho do arquivo direto para o `bat` na linha de comando do `fzf`; antecipar o subcomando público `show` para o MVP só por causa da prévia.

### DE-13 — Código Go segue as convenções do Go; o resto, as da casa

- **Decisão:** arquivos `.go` e pastas de pacote usam a convenção do Go: minúsculas, sem hífen, testes em `_test.go` ao lado do código. Todo o resto do repositório (documentação, integrações, scripts) segue `kebab-case`, e as regras de idioma da casa valem para tudo.
- **Por quê:** os padrões da casa foram escritos para monorepos TypeScript e não têm referência para Go. O ferramental do Go exige o sufixo `_test.go` e não aceita hífen em nome de pacote.
- **Rejeitadas:** `kebab-case` também nas pastas de pacote, o que obrigaria o nome do pacote a diferir do nome da pasta.

### DE-14 — Escopos de commit dentro do vocabulário da casa

- **Decisão:** `core` para o código Go, `specs` e `plans` para a documentação correspondente, `deps` para dependências. Mudanças nas integrações (neovim, yazi) e mudanças gerais vão sem escopo.
- **Por quê:** o vocabulário de escopos da casa é fechado, e escopo inventado torna o histórico impossível de filtrar.
- **Rejeitadas:** criar escopos locais como `cli`, `nvim` e `yazi` sem antes acrescentá-los ao padrão da casa.

### DE-15 — O glossário fixa os nomes em inglês; modelo é template

- **Decisão:** o [[glossary|glossário]] do projeto, num arquivo próprio, tem uma coluna Inglês com o nome de cada termo do domínio, e o código, o JSON e os commits usam esse nome. Modelo é *template*. Campo de nome é *name prompt*. Integração (*integration*) e adaptador (*adapter*) são termos distintos: a integração é o conjunto em `integrations/<programa>/`, com as instruções de instalação, e o adaptador é o código que ela contém.
- **Por quê:** o padrão da casa pede um glossário PT → EN em `docs/`, e sem ele o mesmo conceito ganha nomes diferentes: a descrição do `catalog-info.yaml` chamava os modelos de *models*. O glossário tem arquivo próprio porque vale para o projeto inteiro, e não só para [[bootstrap]]. *Template* é o nome que o padrão XDG dá à pasta (`XDG_TEMPLATES_DIR`). *Prompt* é o termo comum em programas de terminal para pedir um texto ao usuário. A spec e o roadmap já usavam "integração" para o conjunto e "adaptador" para o código.
- **Rejeitadas:** *model*, que em inglês sugere modelo de dados ou de IA; um termo só, *adapter*, com a pasta renomeada para `adapters/`, o que obrigaria a reescrever a spec e o roadmap, onde "integração" já nomeia o conjunto; *name field*, a tradução literal, e *name input*, que nomeiam um elemento de formulário, e não o ato de pedir o nome; o glossário como seção de [[bootstrap]], onde nasceu, que prenderia à spec de partida um documento que os requisitos, os planos e o código também consultam.

### DE-16 — O módulo Go é gitlab.com/rg-dev/janus

- **Decisão:** o caminho do módulo é `gitlab.com/rg-dev/janus`, igual ao remoto `origin` (`git@gitlab.com:rg-dev/janus.git`). Trocar o remoto exige trocar o módulo no mesmo commit.
- **Por quê:** a proposta tinha sido descartada na migração ([[technical-debts|DT-05]]), mas em 2026-11-16 o remoto já estava configurado com esse endereço e com a `main` enviada. Com o caminho igual ao remoto, `go install gitlab.com/rg-dev/janus/cmd/janus@latest` funciona quando o repositório for publicado; enquanto for privado, isso pede `GOPRIVATE=gitlab.com/rg-dev/*` e acesso ao git por SSH, porque o proxy público do Go não enxerga repositório privado. Trocar depois custa pouco: `go mod edit -module` e a troca dos imports internos, num commit mecânico. Todo o código fica em `internal/`, que nenhum outro módulo pode importar, e os adaptadores chamam o binário, não o módulo.
- **Rejeitadas:** só `janus`, sem domínio, que impede instalar a partir do remoto e teria de mudar numa publicação; um caminho no domínio da casa, servido por uma página que redireciona ao repositório, que desliga o módulo do remoto mas exige domínio e hospedagem.

### DE-17 — O campo de nome é o próprio fzf

- **Decisão:** o campo de nome do modo interativo é o `fzf` com a lista vazia: `--query` traz o valor inicial, `--print-query` devolve o texto digitado, o cabeçalho mostra o aviso (destino existente, nome inválido) e `Ctrl-D` cancela sempre (`--bind ctrl-d:abort`). Fica atrás da interface do campo de nome, como o seletor fica atrás da dele.
- **Por quê:** o `fzf` já é obrigatório no modo interativo, então nenhuma dependência entra, e o campo tem a mesma aparência do seletor. Na configuração padrão do `fzf`, `Ctrl-D` só cancela com o campo vazio; a spec pede que cancele sempre. O custo aceito é usar o `fzf` para o que ele não foi feito: sobra uma área de lista vazia, reduzida com `--height`.
- **Rejeitadas:** uma biblioteca Go de edição de linha com texto inicial, que a spec autorizava, mas traria uma dependência e um visual diferente do seletor; um editor de linha próprio com o terminal em modo cru, que pediria `golang.org/x/term` e código nosso para cursor, acentos e largura de caracteres.

### DE-18 — O manifesto é lido com BurntSushi/toml

- **Decisão:** o `.janus.toml` é lido com `github.com/BurntSushi/toml`, versão 1.6.0. É a única dependência externa do projeto.
- **Por quê:** `MetaData.Keys()` devolve as chaves na ordem do arquivo, o que permite conferir cada uma e citar sempre a primeira errada, e o erro de sintaxe (`ParseError`), inclusive o de chave repetida, traz a linha para a mensagem. O manifesto é lido como tabela genérica e conferido chave a chave, porque o erro de tipo da biblioteca não diz a chave de forma estruturada; a primeira ideia, `MetaData.Undecoded()` com uma struct, esbarrava nisso. Conferido no código da versão 1.6.0 em 2026-11-16.
- **Rejeitadas:** `github.com/pelletier/go-toml/v2`, que também tem modo estrito; a vantagem dela é a velocidade, irrelevante para um arquivo de poucas linhas.

### DE-19 — Os argumentos são lidos por código próprio, sem o pacote flag

- **Decisão:** a linha de comando é lida por um leitor próprio e pequeno, em `internal/cli`. Ele aceita `--opcao valor`, `--opcao=valor` e `--`, com as opções antes ou depois do identificador.
- **Por quê:** o pacote `flag` para de ler opções no primeiro argumento que não é opção, e o contrato põe as opções depois do identificador (`janus new <id> --dir x`). Ele também aceita `-dir` com um hífen só e trata `-h`, formas que o contrato considera uso incorreto, e escreve as próprias mensagens em inglês.
- **Rejeitadas:** o `flag` com várias leituras seguidas, uma a cada argumento que não é opção, que perde o efeito do `--`; bibliotecas como a `cobra`, uma dependência grande demais para dois subcomandos, que a stack só aceitaria com um DE.

### DE-20 — O user-dirs.dirs é lido como a GLib lê

- **Decisão:** a linha de `XDG_TEMPLATES_DIR` no `user-dirs.dirs` é interpretada como na função `load_user_special_dirs_from_string` da GLib. Espaços e tabulações são aceitos no início da linha e em volta do `=`. Vale a última linha que define a variável. O valor vai da primeira aspa até a última aspa da linha, sem desfazer escapes. Ele começa por `/` ou por `$HOME` seguido de `/` ou da aspa final. Barras no fim são removidas. `"$HOME"` sozinho é a pasta pessoal e, como a spec manda, falha com o código 3.
- **Por quê:** o Nemo obtém a pasta pela GLib (`g_get_user_special_dir`), e o Janus deve achar a mesma pasta que ele. A spec não definia esses casos. O código foi lido em `glib/gutilsprivate.c` em 2026-11-16. Uma diferença é assumida: a GLib ignora uma linha com valor relativo e fica com a anterior; o Janus falha com o código 3, como a spec já pedia, para que o erro apareça.
- **Rejeitadas:** ler o arquivo como o shell lê (o `xdg-user-dir` faz `source` dele), desfazendo `\"` e `\$`, como sugeriu o [[similar-tools|estudo das ferramentas parecidas]], o que divergiria do Nemo justamente nesses casos; aceitar só a forma exata que a spec descrevia, o que recusaria um arquivo que o Nemo aceita.

### DE-21 — As mensagens moram em internal/messages

- **Decisão:** todo texto mostrado ao usuário fica no pacote `internal/messages`, que também transforma cada erro tipado em mensagem. O `cli` transforma o erro em código de saída, e o `interactive` usa o mesmo pacote para os avisos do laço de escolha e nome. Os erros que não vêm do núcleo (uso incorreto, programa ausente, pasta sem modelos no modo interativo, falha do `fzf` e cancelamento) moram em `internal/failure`, que `messages` e `cli` importam sem importar um ao outro.
- **Por quê:** a estrutura proposta em [[bootstrap]] punha as mensagens em `internal/cli`, mas o modo interativo também mostra avisos, e o `cli` importa o `interactive`. Com as mensagens no `cli`, os dois pacotes se importariam, o que o Go proíbe. Um pacote próprio mantém o lugar único de [[#DE-11 — Contrato em inglês, mensagens em português|DE-11]].
- **Rejeitadas:** o `cli` passar ao `interactive` uma função que formata mensagens, o que espalharia a montagem dos avisos por dois pacotes; mensagens em cada pacote de código, o que quebraria DE-11.

### DE-22 — A versão vem do que o Go grava no binário

- **Decisão:** `janus --version` imprime a versão do módulo principal, lida com `debug.ReadBuildInfo`. Desde o Go 1.24, o `go build` e o `go install` gravam essa versão a partir do git: a tag (`v0.1.0`) ou uma pseudoversão com data e commit, com `+dirty` quando há mudança não registrada. Sem essa informação, imprime `(devel)`.
- **Por quê:** não há constante de versão para esquecer de atualizar, e a versão aponta o commit exato do binário instalado.
- **Rejeitadas:** uma constante no código, que envelhece; injetar a versão com `-ldflags`, que exigiria um script de build no lugar do `go install ./cmd/janus` da spec.

### DE-23 — A pasta de destino nunca fica dentro da pasta de modelos

- **Decisão:** a pasta de destino não pode ser a pasta de modelos nem ficar dentro dela. A comparação usa os caminhos reais, com os links resolvidos, porque `~/Modelos` costuma ser um link. `janus new` falha com o código 10, depois do 5, e o modo interativo também, antes de abrir o seletor: a ordem dele passa a ser 2, 3, 4, 5, 9, 10. A mensagem explica que o Janus não escreve na pasta de modelos.
- **Por quê:** fecha a lacuna de [[technical-debts|DT-13]]: a [[#DE-02 — A pasta XDG de modelos é a fonte única, compartilhada com o Nemo|DE-02]] diz que o Janus nunca escreve na pasta de modelos, mas nada no contrato impedia criar o arquivo lá. No yazi, é fácil acionar o atalho com essa pasta aberta, e o arquivo novo apareceria sem querer como modelo, no menu do Nemo e em `janus list`.
- **Rejeitadas:** permitir e restringir a DE-02 a "não alterar os modelos existentes", o que deixaria o modelo acidental aparecer; um código de saída novo, porque o 10 já é o da pasta de destino e nenhum adaptador precisa distinguir os casos; perguntar antes de criar, porque `janus new` nunca pergunta.

### DE-24 — O atalho do yazi é c seguido de m

- **Decisão:** o trecho do `keymap.toml` do yazi usa `on = ["c", "m"]` com `run = 'shell "janus --hold" --block'`, em `[[mgr.prepend_keymap]]`.
- **Por quê:** o autor executou os cinco passos do roteiro manual na yazi 26.9.1 em 2026-11-16, com essa tecla, e aprovou ([[technical-debts|DT-03]]). A leitura é "criar a partir de modelo", e o `prepend_keymap` dá precedência sobre o mapa padrão do yazi, em que `c` abre as ações de copiar caminhos.
- **Rejeitadas:** manter a tecla como sugestão até a [[roadmap|RM-01.02]], o que deixaria em aberto uma escolha já testada; outra combinação, que ninguém pediu.

### DE-25 — A branch de entrega do roadmap leva o ID e mora no worktree do herdr

- **Decisão:** a branch de uma entrega do roadmap se chama `feat/rm-NN.NN_<slug>`, com o ID da entrega antes do `_` e o slug em `kebab-case` depois, e o worktree dela é o que o herdr cria em `~/.herdr/worktrees/janus/`. O resto do trabalho segue o padrão da casa: `feat/<slug>` e `fix/<slug>`, com o worktree irmão em `~/Projetos/dev/janus/`. Vale só para o Janus; o padrão da casa não muda.
- **Por quê:** as duas primeiras entregas nasceram assim, pelo herdr, que o autor usa para abrir as sessões de trabalho ([[technical-debts|DT-11]]). O ID na branch liga o histórico à entrega sem precisar abrir o roadmap. Mover o worktree obrigaria a trocar de pasta no meio da sessão e quebraria o vínculo com o herdr.
- **Rejeitadas:** seguir o padrão da casa também nas entregas, renomeando a branch e movendo o worktree para `~/Projetos/dev/janus/feat_<slug>`, o que pediria a troca de pasta no meio da sessão; adiar de novo a [[technical-debts|DT-11]], o que deixaria cada branch nova repetindo a divergência sem regra. O erro de digitação de `feat/rm-01.02_yazy-integration` fica, porque a branch some depois do merge e renomeá-la pode confundir o herdr.

### DE-26 — A recusa da pasta de modelos mora em DestinationDir

- **Decisão:** `core.DestinationDir` recebe a pasta de modelos e, além de conferir que a pasta de destino existe e é pasta, recusa a pasta de modelos e o que fica dentro dela ([[#DE-23 — A pasta de destino nunca fica dentro da pasta de modelos|DE-23]]). A comparação resolve os links dos dois caminhos com `filepath.EvalSymlinks` e usa `filepath.Rel`: a pasta é recusada quando o resultado é `.` ou não sobe com `..`. O erro continua o `DestinationDirError`, com um campo `Reason` que separa "inexistente ou não é pasta" de "dentro da pasta de modelos" para a mensagem; o código de saída continua 10. O caminho devolvido continua o absoluto, sem os links resolvidos, para o caminho impresso ser o que o usuário vê.
- **Por quê:** `new` e o modo interativo já passam por essa função, cada um no ponto que a ordem do contrato pede (no `new`, depois do 5 e antes do 6; no modo interativo, depois do 9 e antes do seletor). Assim a regra vale nos dois sem código novo nos chamadores, e a assinatura nova faz o compilador acusar qualquer chamador esquecido. O `filepath.Rel` evita o falso positivo da comparação por prefixo de texto, que recusaria `Modelos-antigo` com a pasta de modelos em `Modelos`.
- **Rejeitadas:** uma função separada, chamada por cada frontend depois de `DestinationDir`, em que esquecer uma chamada quebraria [[core|NFR-CORE-01]] sem aviso do compilador; a verificação dentro de `core.Create`, que chega tarde: o modo interativo já teria aberto o seletor, e no `new` o 10 viria depois do 6 e do 7; comparar com `os.SameFile` subindo da pasta de destino até a raiz, que pegaria até uma montagem com `mount --bind`, mas foge da comparação de caminhos decidida na DE-23 e custa um `Stat` por nível, por um caso que ninguém usa.
