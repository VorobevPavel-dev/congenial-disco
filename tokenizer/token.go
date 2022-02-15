package tokenizer

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/VorobevPavel-dev/congenial-disco/utility"
)

var (
	KindNames               *map[TokenKind]string
	Reserved                map[TokenKind]map[string]*Token
	symbols                 *[]string
	keywords                *[]string
	types                   *[]string
	ErrUnsupportedTokenType error = errors.New("unsupported token type")
)

func init() {
	KindNames = KindMap()
	Reserved = Constants()
	symbols = Keys(Reserved[SymbolKind])
	keywords = Keys(Reserved[KeywordKind])
	types = Keys(Reserved[TypeKind])
}

// KindToString returns string representation of token kind
func KindToString(kind TokenKind) string {
	return (*KindNames)[kind]
}

type Token struct {
	Value    string    `json:"value"`
	Kind     TokenKind `json:"kind"`
	Position int       `json:"position"`
}

func (t *Token) Equals(other *Token) bool {
	return t.Value == other.Value && t.Kind == other.Kind
}

func (t *Token) String() string {
	data, _ := json.Marshal(t)
	return string(data)
}

// Tokenizer is a function that will parse strings to tokens and validate them
// inside function body (for example validate numbers, keywords, custom values)
type Tokenizer func(string) *Token

// TokenFromString parses given string (one word without spaces) into token (no matter what kind it will have).
func TokenFromString(value string, cursorPosition int) (*Token, error) {
	// Order matters.
	tokenizers := []Tokenizer{
		ParseNumericToken,
		ParseTypeToken,
		ParseKeywordToken,
		ParseSymbolToken,
		ParseIdentifierToken,
	}
	for _, function := range tokenizers {
		token := function(value)
		if token != nil {
			token.Position = cursorPosition + len(token.Value)
			return token, nil
		}
	}
	return nil, ErrUnsupportedTokenType
}

// ParseNumericToken parses given string as integer.
// TODO: Implement support for different types on numbers
func ParseNumericToken(value string) *Token {
	_, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return &Token{
		Value: value,
		Kind:  NumericKind,
	}
}

// ParseKeywordToken checks if given value is a reserved word.
// If it so returns &Token with KeywordKind
func ParseKeywordToken(value string) *Token {
	loweredValue := strings.ToLower(value)
	if value, ok := Reserved[KeywordKind][loweredValue]; ok {
		return value
	}
	return nil
}

// ParseSymbolToken checks if given value is a reserved symbol.
// If it so returns &Token with SymbolKind
func ParseSymbolToken(value string) *Token {
	if value, ok := Reserved[SymbolKind][value]; ok {
		return value
	}
	return nil
}

// ParseSymbolToken checks if given value is a
// reserved keyword describing types. If it so returns &Token with TypeKind
func ParseTypeToken(value string) *Token {
	loweredValue := strings.ToLower(value)
	if utility.StringIsIn(value, *types) {
		return &Token{
			Value: loweredValue,
			Kind:  TypeKind,
		}
	}
	return nil
}

// ParseIdentifierToken takes whole value and convert it into
// token with IdentifierKind
func ParseIdentifierToken(value string) *Token {
	return &Token{
		Value: value,
		Kind:  IdentifierKind,
	}
}

func TokenFromKeyword(value string) *Token {
	if !utility.StringIsIn(value, *keywords) {
		return nil
	}
	return &Token{
		Value: value,
		Kind:  KeywordKind,
	}
}

func TokenFromSymbol(value string) *Token {
	if !utility.StringIsIn(value, *symbols) {
		return nil
	}
	return &Token{
		Value: value,
		Kind:  SymbolKind,
	}
}

// ParseTokenSequence converts string into slice of tokens.
// It uses utility.DivideBySeparators function to divide string to
// set of individual parts and after that tries to parse every part
// as a token.
func ParseTokenSequence(expression string) *[]*Token {
	var (
		startPosition = 0
		resultTokens  []*Token
	)
	parts := utility.DivideBySeparators(expression, *symbols)
	for _, part := range parts {
		token, err := TokenFromString(part, startPosition)
		if err != nil {
			return nil
		}
		token.Position = startPosition
		startPosition += len(token.Value)
		// Removing spaces from token list
		if strings.TrimSpace(token.Value) != "" {
			resultTokens = append(resultTokens, token)
		}
	}
	return &resultTokens
}

// FindToken searches for provided token in a given list of tokens.
// Retuns -1 if token was not found. Otherwise returns its index.
func FindToken(tokens []*Token, expected *Token) int {
	for index := range tokens {
		if tokens[index].Equals(expected) {
			return index
		}
	}
	return -1
}

func GenerateRandomToken(kind TokenKind) *Token {
	switch int(kind) {
	case 0:
		min, max := -100000, 100000
		return &Token{
			Value: strconv.Itoa(rand.Intn(max-min+1) + min),
			Kind:  NumericKind,
		}
	case 1:
		index := rand.Intn(len(*keywords))
		return &Token{
			Value: (*keywords)[index],
			Kind:  KeywordKind,
		}
	case 2:
		// Excludes spaces
		value := " "
		for value == " " {
			index := rand.Intn(len(*symbols))
			value = (*symbols)[index]
		}
		return &Token{
			Value: value,
			Kind:  SymbolKind,
		}
	case 3:
		b := make([]byte, 20)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}
		return &Token{
			Value: string(b),
			Kind:  IdentifierKind,
		}
	case 4:
		index := rand.Intn(len(*types))
		return &Token{
			Value: (*types)[index],
			Kind:  TypeKind,
		}
	}
	return nil
}

// Bracketize will return set of token values inside brackets delimited by comma with space.
func Bracketize(input []*Token) string {
	if len(input) == 0 {
		return ""
	}
	result := fmt.Sprintf("(%s", input[0].Value)
	for _, el := range input[1:] {
		result += fmt.Sprintf(", %s", el.Value)
	}
	return result + ")"
}
