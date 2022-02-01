package tokenizer

import (
	"encoding/json"
	"errors"
	"github.com/VorobevPavel-dev/congenial-disco/utility"
	"strconv"
	"strings"
)

//SQL-reserved words
const (
	SelectKeyword string = "select"
	FromKeyword   string = "from"
	AsKeyword     string = "as"
	TableKeyword  string = "table"
	CreateKeyword string = "create"
	InsertKeyword string = "insert"
	IntoKeyword   string = "into"
	ValuesKeyword string = "values"
)

// Symbol constants
const (
	SemicolonSymbol  string = ";"
	AsteriskSymbol   string = "*"
	CommaSymbol      string = ","
	LeftParenSymbol  string = "("
	RightParenSymbol string = ")"
	SpaceSymbol      string = " "
)

const (
	// NumericKind will correspond to all numeric values
	NumericKind TokenKind = iota
	// KeywordKind will correspond to all string equal to one of keywords
	KeywordKind
	// SymbolKind will correspond to every specified utility symbol
	SymbolKind
	// IdentifierKind will correspond to every custom value (table name, column name, values etc...)
	IdentifierKind
)

type TokenKind uint

var (
	keywords = []string{
		SelectKeyword,
		FromKeyword,
		AsKeyword,
		TableKeyword,
		CreateKeyword,
		InsertKeyword,
		IntoKeyword,
		ValuesKeyword,
	}
	symbols = []string{
		CommaSymbol,
		SemicolonSymbol,
		SpaceSymbol,
		AsteriskSymbol,
		LeftParenSymbol,
		RightParenSymbol,
	}
	ErrUnsupportedTokenType = errors.New("unsupported token type")
)

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

// TokenFromString is a function which will parse given string to token (no matter what kind
// it will have)
func TokenFromString(value string, cursorPosition int) (*Token, error) {
	tokenizers := []Tokenizer{ParseNumericToken, ParseKeywordToken, ParseSymbolToken, ParseIdentifierToken}
	for _, function := range tokenizers {
		token := function(value)
		if token != nil {
			token.Position = cursorPosition + len(token.Value)
			return token, nil
		}
	}
	return nil, ErrUnsupportedTokenType
}

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

func ParseKeywordToken(value string) *Token {
	loweredValue := strings.ToLower(value)
	if utility.StringIsIn(loweredValue, keywords) {
		return &Token{
			Value: loweredValue,
			Kind:  KeywordKind,
		}
	}
	return nil
}

func ParseSymbolToken(value string) *Token {
	if utility.StringIsIn(value, symbols) {
		return &Token{
			Value: value,
			Kind:  SymbolKind,
		}
	}
	return nil
}

func ParseIdentifierToken(value string) *Token {
	return &Token{
		Value: value,
		Kind:  IdentifierKind,
	}
}

func TokenFromKeyword(value string) *Token {
	if !utility.StringIsIn(value, keywords) {
		return nil
	}
	return &Token{
		Value: value,
		Kind:  KeywordKind,
	}
}

func TokenFromSymbol(value string) *Token {
	if !utility.StringIsIn(value, symbols) {
		return nil
	}
	return &Token{
		Value: value,
		Kind:  SymbolKind,
	}
}

func ParseTokenSequence(expression string) *[]*Token {
	var (
		startPosition = 0
		resultTokens  []*Token
	)
	parts := utility.DivideBySeparators(expression, symbols)
	for _, part := range parts {
		token, err := TokenFromString(part, startPosition)
		if err != nil {
			return nil
		}
		token.Position = startPosition
		// FIXME: replace it with actual length (for different languages)
		startPosition += len(token.Value)
		resultTokens = append(resultTokens, token)
	}
	return &resultTokens
}

func FindToken(tokens []*Token, expected *Token) int {
	for index := range tokens {
		if tokens[index].Equals(expected) {
			return index
		}
	}
	return -1
}
