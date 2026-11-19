# Janus no yazi

Atalho do [yazi](https://yazi-rs.github.io) que abre o modo interativo do Janus na pasta em que o yazi está: escolha um modelo, confirme o nome, e o arquivo aparece na listagem ao voltar.

## Requisitos

- `janus` no `PATH` que o yazi enxerga.
- `fzf`, que o modo interativo usa para escolher o modelo e pedir o nome.
- `bat`, opcional, para realçar a prévia.

## Instalação

1. Na raiz do repositório do Janus, instale o binário:

   ```
   go install ./cmd/janus
   ```

   Ele vai para a pasta de `$GOBIN`, se ela estiver definida, ou para `$GOPATH/bin`, que por padrão é `~/go/bin`. Essa pasta precisa estar no `PATH`.

2. Leve o bloco de [`keymap.toml`](keymap.toml) para `~/.config/yazi/keymap.toml`. Se o arquivo ainda não existe, copie-o inteiro:

   ```
   mkdir -p ~/.config/yazi
   cp integrations/yazi/keymap.toml ~/.config/yazi/keymap.toml
   ```

   Se ele já existe, acrescente o bloco `[[mgr.prepend_keymap]]` ao fim, sem apagar os outros. A exceção é o arquivo que já tem, dentro de `[mgr]`, um `prepend_keymap = [ … ]` em forma de lista: ali o bloco deixaria o TOML inválido, e a entrada vai como mais um item da lista:

   ```toml
   { on = ["c", "m"], run = 'shell "janus --hold" --block', desc = "Novo arquivo a partir de modelo" },
   ```

3. Feche e abra o yazi de novo, para ele ler o `keymap.toml`.

## Uso

Na pasta onde o arquivo deve nascer, tecle `c` e depois `m`.

- Um modelo com nome fixo é criado direto; os outros pedem o nome, sugerindo o do modelo.
- Se o arquivo já existe, o campo de nome abre com um aviso, para você escolher outro nome.
- `Esc` cancela e volta ao yazi sem criar nada.
- Um erro, como um manifesto inválido, fica na tela até o `Enter`.
- Dentro da pasta de modelos o atalho recusa criar o arquivo: o Janus nunca escreve nela.

As outras teclas do `c`, que copiam caminhos, continuam funcionando: o menu do `c` mostra todas.

## Problemas comuns

- **O menu do `c` não mostra "Novo arquivo a partir de modelo":** o yazi não leu o bloco. Confira o caminho do arquivo, ou a pasta de `YAZI_CONFIG_HOME`, se ela estiver definida, e reinicie o yazi.
- **A tela pisca e volta ao yazi sem fazer nada:** o `janus` provavelmente não está no `PATH` que o yazi enxerga, e a mensagem de erro do shell pode sumir junto com a tela. Acontece quando o yazi é aberto por um lançador gráfico que não leva o `~/go/bin` no `PATH`. Abra o yazi pelo terminal, ou ponha o `~/go/bin` no `PATH` da sessão gráfica.
