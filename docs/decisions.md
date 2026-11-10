---
title: Decisões
aliases:
  - DE-NN
tags:
  - decisoes
---

# Decisões

Decisões do Janus, cada uma com o motivo e as alternativas rejeitadas. As de DE-01 a DE-12 foram tomadas na sessão de desenho de 2026-11-09; DE-13 e DE-14, na migração para este repositório, e DE-15, na escrita do `AGENTS.md`, no mesmo dia. O desenho completo está em [[bootstrap]].

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
