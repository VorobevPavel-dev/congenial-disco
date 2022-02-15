package parser

import (
	"encoding/json"
	"errors"
	"fmt"

	t "github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

var (
	ErrNoSemicolonAtTheEnd     error  = errors.New("provided request does not have ';' SymbolKind token at the end")
	ErrExpectedKeywordTemplate string = "expected %s keyword at %d"

	// ErrorBuilders are functions deicated to process data from parsing functions and return formatted error
	// There are only functions and errors required by all parsers. All specific errors declarated in parser
	// itself

	// ErrExpectedToken builds error according to desired token and position where that token must be
	ErrExpectedToken = func(e *t.Token, p int) error {
		return fmt.Errorf("expected %s %s at %d",
			e.Value,
			t.KindToString(e.Kind),
			p,
		)
	}

	// ErrInvalidTokenKind builds error according to current token and desired TokenKind
	ErrInvalidTokenKind = func(e *t.Token, ek t.TokenKind) error {
		return fmt.Errorf("expected %s but got %s at %d",
			t.KindToString(e.Kind),
			t.KindToString(ek),
			e.Position,
		)
	}
)

const (
	CreateTableType int = iota
	SelectType
	InsertType
	ShowCreateType
)

type Statement struct {
	SelectStatement      *SelectQuery
	CreateTableStatement *CreateTableQuery
	InsertStatement      *InsertIntoQuery
	ShowCreateStatement  *ShowCreateQuery
	// Experimental
	Type int
}

// Experimental
type Query interface {
	// String() string
	Equals(*Query) bool
	Parse([]*t.Token) (*Query, bool, error)
	// CreateOriginal must return string containing original SQL request
	CreateOriginal() string
}

// Experimental
func QueryToString(q *Query) string {
	bytes, _ := json.Marshal(q)
	return string(bytes)
}

// Parse will try to parse statement with all parsers successively
// Returns a Statement struct with only one field not null
func Parse(request string) *Statement {
	// Implement request string as a series of tokens
	tokens := *t.ParseTokenSequence(request)

	createStatement, _ := parseCreateTableQuery(tokens)
	if createStatement != nil {
		return &Statement{
			CreateTableStatement: createStatement,
			Type:                 CreateTableType,
		}
	}
	insertStatement, _ := parseInsertIntoStatement(tokens)
	if insertStatement != nil {
		return &Statement{
			InsertStatement: insertStatement,
			Type:            InsertType,
		}
	}
	selectStatement, _ := parseSelectStatement(tokens)
	if selectStatement != nil {
		return &Statement{
			SelectStatement: selectStatement,
			Type:            SelectType,
		}
	}
	showCreateStatement, _ := parseShowCreateQuery(tokens)
	if showCreateStatement != nil {
		return &Statement{
			ShowCreateStatement: showCreateStatement,
			Type:                ShowCreateType,
		}
	}
	return nil
}
