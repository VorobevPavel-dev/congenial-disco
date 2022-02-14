package tokenizer

import (
	"fmt"
	"math/rand"
	"strconv"
)

// FindToken searches for provided token in a given list of tokens.
// Retuns -1 if token was not found. Otherwise returns its index.
func FindToken(tokens []Token, expected Token) int {
	for index := range tokens {
		if tokens[index].Equals(expected) {
			return index
		}
	}
	return -1
}

// GenerateRandomToken will return token of given kind with
// randomly generated value inside.
func GenerateRandomToken(kind TokenKind) Token {
	switch kind {
	case NumericKind:
		min, max := -100000, 100000
		return Token{
			Value: strconv.Itoa(rand.Intn(max-min+1) + min),
			Kind:  NumericKind,
		}
	case KeywordKind:
		return Token{
			Value: Keywords[rand.Intn(len(Keywords))],
			Kind:  KeywordKind,
		}
	case SymbolKind:
		// Excludes spaces
		value := " "
		for value == " " {
			value = Symbols[rand.Intn(len(Symbols))]
		}
		return Token{
			Value: value,
			Kind:  SymbolKind,
		}
	case IdentifierKind:
		b := make([]byte, 20)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}
		return Token{
			Value: string(b),
			Kind:  IdentifierKind,
		}
	case TypeKind:
		return Token{
			Value: Types[rand.Intn(len(Types))],
			Kind:  TypeKind,
		}
	}
	// Never returns empty token
	return Token{}
}

// Bracketize will return set of token values inside brackets delimited by comma with space.
func Bracketize(input []Token) string {
	if len(input) == 0 {
		return ""
	}
	result := fmt.Sprintf("(%s", input[0].Value)
	for _, el := range input[1:] {
		result += fmt.Sprintf(", %s", el.Value)
	}
	return result + ")"
}

// getKeys will return only keys from map[string]Token
func getKeys(inputMap map[string]Token) []string {
	result := make([]string, len(inputMap))
	i := 0
	for key := range inputMap {
		result[i] = key
		i++
	}
	return result
}
