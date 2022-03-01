package table

import (
	"strconv"
	"strings"

	"github.com/VorobevPavel-dev/congenial-disco/parser"
	t "github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

var (
	typeCorreation = map[string]t.TokenKind{
		"int":  t.NumericKind,
		"text": t.IdentifierKind,
	}
	ErrIncorrectSetLenghtTemplate = "inserting set has incorrect number of values. expected: %s, got: %s"
	ErrIncorrectValueType         = "inserting value \"%s\" has incorrect type: %s, expected %s"
	ErrIncorrectSelectColumn      = "column %s is not provided in table %s"
	ErrIncorrectConditionColumn   = "column \"%s\" specified in condition is not provided in table %s"
	ErrConditionRuntime           = "error while processing condition: %v"
)

type Table interface {
	//IsInitialized will check if table has any columns created
	IsInitialized() bool
	//Create will handle CREATE TABLE requests
	Create(req *parser.CreateTableQuery) (Table, string, error)
	//Create will handle SELECT requests
	Select(req *parser.SelectQuery) ([][]string, error)
	//Create will handle INSERT requests
	Insert(req *parser.InsertIntoQuery) (Table, error)
	//Engine will return engine name as a string
	Engine() string
	// Name will return table name as a string
	Name() string
	//ShowCreate will return SQL-like string with original CREATE TABLE command
	ShowCreate() string
}

// element is a struct representing inner value.
type element struct {
	kind  t.TokenKind
	value string
}

func (e *element) asInt() (int, error) {
	val, err := strconv.Atoi(e.value)
	return val, err
}

func (e *element) equals(o *element) bool {
	return strings.Compare(e.value, o.value) == 0
}

func (e *element) less(o *element) bool {
	switch e.kind {
	case t.NumericKind:
		r, _ := e.asInt()
		l, _ := o.asInt()
		return r < l
	default:
		return strings.Compare(e.value, o.value) == -1
	}
}

func (e *element) greater(o *element) bool {
	switch e.kind {
	case t.NumericKind:
		r, _ := e.asInt()
		l, _ := o.asInt()
		return r > l
	default:
		return strings.Compare(e.value, o.value) == 1
	}
}

func tokenToElement(token *t.Token) element {
	return element{
		kind:  token.Kind,
		value: token.Value,
	}
}
