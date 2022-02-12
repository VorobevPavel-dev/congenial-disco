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

// List of reserved words
const (
	SelectKeyword string = "select"
	FromKeyword   string = "from"
	AsKeyword     string = "as"
	TableKeyword  string = "table"
	CreateKeyword string = "create"
	InsertKeyword string = "insert"
	IntoKeyword   string = "into"
	ValuesKeyword string = "values"
	ShowKeyword   string = "show"
)

// List of constant symbols
const (
	SemicolonSymbol  string = ";"
	AsteriskSymbol   string = "*"
	CommaSymbol      string = ","
	LeftParenSymbol  string = "("
	RightParenSymbol string = ")"
	SpaceSymbol      string = " "
)

// List of reserver words describing type of something
const (
	IntType  string = "int"
	TextType string = "text"
	charset  string = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
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
	// TypeKind will correspond to every type in request
	TypeKind
)

// KindToString returns string representation of token kind
func KindToString(kind int) string {
	switch kind {
	case 0:
		return "number"
	case 1:
		return "Keyword"
	case 2:
		return "symbol"
	case 3:
		return "identifier"
	case 4:
		return "type"
	default:
		return "unknown"
	}
}

type TokenKind uint

func Keywords() *[]string {
	return &[]string{
		SelectKeyword,
		FromKeyword,
		AsKeyword,
		TableKeyword,
		CreateKeyword,
		InsertKeyword,
		IntoKeyword,
		ValuesKeyword,
		ShowKeyword,
	}
}

func Symbols() *[]string {
	return &[]string{
		CommaSymbol,
		SemicolonSymbol,
		SpaceSymbol,
		AsteriskSymbol,
		LeftParenSymbol,
		RightParenSymbol,
	}
}

func Types() *[]string {
	return &[]string{
		IntType,
		TextType,
	}
}

var (
	symbols                       = *Symbols()
	keywords                      = *Keywords()
	types                         = *Types()
	ErrUnsupportedTokenType error = errors.New("unsupported token type")
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
	if utility.StringIsIn(loweredValue, keywords) {
		return &Token{
			Value: loweredValue,
			Kind:  KeywordKind,
		}
	}
	return nil
}

// ParseSymbolToken checks if given value is a reserved symbol.
// If it so returns &Token with SymbolKind
func ParseSymbolToken(value string) *Token {
	if utility.StringIsIn(value, symbols) {
		return &Token{
			Value: value,
			Kind:  SymbolKind,
		}
	}
	return nil
}

// ParseSymbolToken checks if given value is a
// reserved keyword describing types. If it so returns &Token with TypeKind
func ParseTypeToken(value string) *Token {
	loweredValue := strings.ToLower(value)
	if utility.StringIsIn(value, types) {
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

// ParseTokenSequence converts string into slice of tokens.
// It uses utility.DivideBySeparators function to divide string to
// set of individual parts and after that tries to parse every part
// as a token.
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
		index := rand.Intn(len(*Keywords()))
		return &Token{
			Value: (*Keywords())[index],
			Kind:  KeywordKind,
		}
	case 2:
		// Excludes spaces
		value := " "
		for value == " " {
			index := rand.Intn(len(*Symbols()))
			value = (*Symbols())[index]
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
		index := rand.Intn(len(*Types()))
		return &Token{
			Value: (*Types())[index],
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
