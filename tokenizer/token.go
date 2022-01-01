package tokenizer

import (
	"errors"
	"github.com/VorobevPavel-dev/congenial-disco/utility"
	"strconv"
	"strings"
)

type keyword string

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
	ErrUnsupportedTokenType = errors.New("unsupported token type")
)

const (
	NumericKind TokenKind = iota
	KeywordKind TokenKind = iota
)

type Token struct {
	Value    string
	Kind     TokenKind
	Position int
}

func (t *Token) equals(other *Token) bool {
	return t.Value == other.Value && t.Kind == other.Kind
}

type Tokenizer func(string) *Token

// TokenFromString is a function which will parse given string to token (no matter what kind
// it will have)
func TokenFromString(value string, cursorPosition int) (*Token, error) {
	tokenizers := []Tokenizer{ParseNumericToken, ParseKeywordToken}
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
