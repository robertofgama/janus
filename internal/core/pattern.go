// Package core é o núcleo do Janus: descobre a pasta de modelos, lista os modelos, lê o
// manifesto, resolve o nome final e cria o arquivo. Não imprime, não lê do teclado e não
// conhece códigos de saída: devolve valores e erros tipados (NFR-CORE-02).
package core

import "strings"

// matchPattern diz se pattern casa o identificador id inteiro (FR-CORE-12). O único
// curinga é "*", que casa qualquer sequência sem "/", inclusive vazia; todo o resto é
// literal. path.Match foi rejeitado porque interpreta "?", "[" e "\".
//
// Como "*" não atravessa "/", cada "/" do identificador casa com um "/" literal do padrão:
// basta comparar trecho a trecho.
func matchPattern(pattern, id string) bool {
	patternParts := strings.Split(pattern, "/")
	idParts := strings.Split(id, "/")
	if len(patternParts) != len(idParts) {
		return false
	}
	for i := range patternParts {
		if !matchPart(patternParts[i], idParts[i]) {
			return false
		}
	}
	return true
}

// matchPart casa um trecho sem "/". Ao falhar, volta ao último "*" visto e faz ele engolir
// mais um byte: é o algoritmo clássico de curinga, sem recursão.
func matchPart(pattern, s string) bool {
	p, i := 0, 0
	star, mark := -1, 0
	for i < len(s) {
		switch {
		case p < len(pattern) && pattern[p] == '*':
			star, mark = p, i
			p++
		case p < len(pattern) && pattern[p] == s[i]:
			p++
			i++
		case star >= 0:
			mark++
			i = mark
			p = star + 1
		default:
			return false
		}
	}
	for p < len(pattern) && pattern[p] == '*' {
		p++
	}
	return p == len(pattern)
}
