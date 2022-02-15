package parser

import (
	"encoding/json"
	"errors"
	"fmt"

	t "github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

// SelectQuery give access to requests with syntax
//
// SELECT (<column_name1>[...]) FROM <table_name>;
//
// where <column_nameX> and <table_name> must be tokens with IdentifierKind
type SelectQuery struct {
	Columns []*t.Token `json:"item"`
	From    *t.Token   `json:"from"`
}

// String method needs to be implemented in order to implement Query interface.
// Returns JSON object describing necessary information
func (slct SelectQuery) String() string {
	bytes, _ := json.Marshal(slct)
	return string(bytes)
}

// Equals method needs to be implemented in order to implement Query interface.
// Returns true if tokens for columns and table names are equal.
func (slct SelectQuery) Equals(other *SelectQuery) bool {
	if len(slct.Columns) != len(other.Columns) {
		return false
	}
	for index := range slct.Columns {
		if !slct.Columns[index].Equals(other.Columns[index]) {
			return false
		}
	}
	return slct.From.Equals(other.From)
}

// CreateOriginal method needs to be implemented in order to implement Query interface.
// Returns original SQL query representing data in current Query
func (slct SelectQuery) CreateOriginal() string {
	result := fmt.Sprintf("SELECT %s FROM %s;",
		t.Bracketize(slct.Columns),
		slct.From.Value,
	)
	return result
}

func parseSelectStatement(tokens []*t.Token) (*SelectQuery, error) {
	// Validate that set of tokens has ';' SymbolKind token at the end
	if !tokens[len(tokens)-1].Equals(t.Reserved[t.SymbolKind][";"]) {
		return nil, ErrNoSemicolonAtTheEnd
	}

	var (
		columns   []*t.Token
		tableName *t.Token
		cursor    int = 0
	)

	// Process SELECT keyword
	if !tokens[cursor].Equals(t.Reserved[t.KeywordKind]["select"]) {
		return nil, fmt.Errorf("expected SELECT keyword at %d", tokens[0].Position)
	}
	cursor++

	// Process set of columns if any were specified
	if !tokens[cursor].Equals(t.Reserved[t.SymbolKind]["("]) {
		return nil, ErrExpectedToken(t.Reserved[t.SymbolKind]["("], tokens[cursor].Position)
	}
	colDefStartPos := cursor
	cursor++
	colDefEndPos := t.FindToken(tokens, t.Reserved[t.SymbolKind][")"])
	if colDefEndPos == colDefStartPos {
		return nil, errors.New("no columns specified")
	}

	for cursor = colDefStartPos + 1; cursor < colDefEndPos; cursor++ {
		if tokens[cursor].Equals(t.Reserved[t.SymbolKind][","]) {
			cursor++
		}
		if tokens[cursor].Kind != t.IdentifierKind {
			return nil, ErrInvalidTokenKind(tokens[cursor], t.IdentifierKind)
		}
		columns = append(columns, tokens[cursor])
	}
	cursor++

	// Process FROM keyword
	if !tokens[cursor].Equals(t.Reserved[t.KeywordKind]["from"]) {
		return nil, fmt.Errorf("expected FROM keyword at %d", tokens[0].Position)
	}
	cursor++

	// Process table name
	if tokens[cursor].Kind != t.IdentifierKind {
		return nil, ErrInvalidTokenKind(tokens[cursor], t.IdentifierKind)
	}
	tableName = tokens[cursor]
	cursor++

	// Process ";" symbol at the end
	if !tokens[cursor].Equals(t.Reserved[t.SymbolKind][";"]) {
		return nil, ErrExpectedToken(t.Reserved[t.SymbolKind][";"], tokens[cursor].Position)
	}

	return &SelectQuery{
		Columns: columns,
		From:    tableName,
	}, nil
}
