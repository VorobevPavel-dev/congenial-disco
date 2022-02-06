package parser

import (
	"encoding/json"
	"errors"
	"fmt"

	t "github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

// ColumnDefinition contains information as Tokens about columns provided
// in SQL query
type ColumnDefinition struct {
	Name     *t.Token
	Datatype *t.Token
}

// Equals method is necessary for CreateTableQuery.Equals method.
// It compares tokens for name and type for provided column definitions.
func (cd *ColumnDefinition) Equals(other *ColumnDefinition) bool {
	return cd.Name.Equals(other.Name) && cd.Datatype.Equals(other.Datatype)
}

// CreateTableQuery give aceess to requests with syntax
//
// CREATE TABLE <table_name> (<column_name> <column_type> [, <column_name> <column_type>])
//
// where <table_name>, <column_name> must be tokens with IdentifierKind, but
// <column_type> with TypeKind
type CreateTableQuery struct {
	// Name of table requested to be created
	Name *t.Token
	// Slice of column definitions that table will contain
	Cols []*ColumnDefinition
}

// String method needs to be implemented in order to implement Query interface.
// Returns JSON object describing necessary information
func (ct CreateTableQuery) String() string {
	bytes, _ := json.Marshal(ct)
	return string(bytes)
}

// Equals method needs to be implemented in order to implement Query interface.
// Returns true if ColumnDefinitions of both queries are equal and names of reqested for
// creation tables has no differ
func (ct CreateTableQuery) Equals(other *CreateTableQuery) bool {
	if len(ct.Cols) != len(other.Cols) {
		return false
	}
	for index := range ct.Cols {
		if !ct.Cols[index].Equals(other.Cols[index]) {
			return false
		}
	}
	return ct.Name.Equals(other.Name)
}

// CreateOriginal method needs to be implemented in order to implement Query interface.
// Returns original SQL query representing data in current Query
func (ct CreateTableQuery) CreateOriginal() string {
	result := fmt.Sprintf("CREATE TABLE %s (%s %s",
		ct.Name.Value,
		ct.Cols[0].Name.Value,
		ct.Cols[0].Datatype.Value,
	)

	for _, col := range ct.Cols[1:] {
		result += fmt.Sprintf(", %s %s", col.Name.Value, col.Datatype.Value)
	}
	result += ");"
	return result
}

// parseCreateTableQuery will process set of tokens and try to parse it like CREATE TABLE
// query with syntax described in comments to CreateTableQuery struct.
// Input set must have SymbolToken with value ';' at the end in order to be parsed.
func parseCreateTableQuery(tokens []*t.Token) (*CreateTableQuery, error) {
	// Validate that set of tokens has ';' SymbolKind token at the end
	if !tokens[len(tokens)-1].Equals(t.TokenFromSymbol(";")) {
		return nil, ErrNoSemicolonAtTheEnd
	}

	var (
		tableName *t.Token
		columns   []*ColumnDefinition
		cursor    int = 0
	)

	// Process "CREATE TABLE " sequence
	if !tokens[cursor].Equals(t.TokenFromKeyword("create")) {
		return nil, ErrExpectedToken(t.TokenFromKeyword("create"), tokens[cursor].Position)
	}
	cursor++

	if !tokens[cursor].Equals(t.TokenFromKeyword("table")) {
		return nil, ErrExpectedToken(t.TokenFromKeyword("table"), tokens[cursor].Position)
	}
	cursor++

	// Process table name
	if tokens[cursor].Kind != t.IdentifierKind {
		return nil, ErrInvalidTokenKind(tokens[cursor], t.IdentifierKind)
	}
	tableName = tokens[cursor]
	cursor++

	if !tokens[cursor].Equals(t.TokenFromSymbol("(")) {
		return nil, ErrExpectedToken(t.TokenFromSymbol("("), tokens[cursor].Position)
	}
	colDefStartPos := cursor
	cursor++
	colDefEndPos := t.FindToken(tokens, t.TokenFromSymbol(")"))

	if colDefEndPos == colDefStartPos {
		return nil, errors.New("no columns specified")
	}

	for cursor = colDefStartPos + 1; cursor < colDefEndPos; cursor += 2 {
		if tokens[cursor].Equals(t.TokenFromSymbol(",")) {
			cursor++
		}

		tempColDef := &ColumnDefinition{}
		if tokens[cursor].Kind != t.IdentifierKind {
			return nil, ErrInvalidTokenKind(tokens[cursor], t.IdentifierKind)
		}
		tempColDef.Name = tokens[cursor]
		if tokens[cursor+1].Kind != t.TypeKind {
			return nil, ErrInvalidTokenKind(tokens[cursor], t.TypeKind)
		}
		tempColDef.Datatype = tokens[cursor+1]

		columns = append(columns, tempColDef)
		if tokens[cursor+2].Equals(t.TokenFromSymbol(")")) {
			break
		} else if !tokens[cursor+2].Equals(t.TokenFromSymbol(",")) {
			return nil, ErrExpectedToken(t.TokenFromSymbol(","), tokens[cursor+2].Position)
		}
	}
	return &CreateTableQuery{
		Name: tableName,
		Cols: columns,
	}, nil
}
