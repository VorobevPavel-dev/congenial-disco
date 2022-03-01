package parser

import (
	"encoding/json"
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
func (cd ColumnDefinition) Equals(other ColumnDefinition) bool {
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
	Cols    *[]ColumnDefinition
	Engine  *t.Token
	OrderBy *t.Token
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
	if len(*ct.Cols) != len(*other.Cols) {
		return false
	}

	if ct.OrderBy != nil && !ct.OrderBy.Equals(other.OrderBy) {
		return false
	}

	for index := range *ct.Cols {
		if !(*ct.Cols)[index].Equals((*other.Cols)[index]) {
			return false
		}
	}
	return ct.Name.Equals(other.Name) && ct.Engine.Equals(other.Engine)
}

// CreateOriginal method needs to be implemented in order to implement Query interface.
// Returns original SQL query representing data in current Query
func (ct CreateTableQuery) CreateOriginal() string {
	result := fmt.Sprintf("CREATE TABLE %s (%s %s",
		ct.Name.Value,
		(*ct.Cols)[0].Name.Value,
		(*ct.Cols)[0].Datatype.Value,
	)

	for _, col := range (*ct.Cols)[1:] {
		result += fmt.Sprintf(", %s %s", col.Name.Value, col.Datatype.Value)
	}
	result += fmt.Sprintf(") ENGINE %s;", ct.Engine.Value)
	// result += ");"
	return result
}
