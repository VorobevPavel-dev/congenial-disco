package parser

import (
	"encoding/json"
	"fmt"

	t "github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

// InsertIntoQuery give access to requests with syntax
//
// INSERT INTO <table_name> [(<column_name1>,...,<column_nameN>)] VALUES (<value1,...,<valueN>);
//
// where <table_name>, <column_nameX> must be token with IdentifierKind,
// <valueX> can be only NumberKind or IdentifierKind.
// Number of columns must be exact same as number of provided values.
// If no <column_name> specified then number of <valueX> must be exact as number of columns in target table.
type InsertIntoQuery struct {
	Table       *t.Token
	ColumnNames []*t.Token
	Values      []*t.Token
}

// String method needs to be implemented in order to implement Query interface.
// Returns JSON object describing necessary information
func (ins InsertIntoQuery) String() string {
	bytes, _ := json.Marshal(ins)
	return string(bytes)
}

// Equals method needs to be implemented in order to implement Query interface.
// Returns true if tokens for values, column names and table are same as in other
// InsertIntoQuery.
func (ins InsertIntoQuery) Equals(other *InsertIntoQuery) bool {
	if len(ins.Values) != len(other.Values) {
		return false
	}

	for index := range ins.Values {
		if !ins.Values[index].Equals(other.Values[index]) {
			return false
		}
	}

	if len(ins.ColumnNames) != len(other.ColumnNames) {
		return false
	}

	for index := range ins.ColumnNames {
		if !ins.ColumnNames[index].Equals(other.ColumnNames[index]) {
			return false
		}
	}
	return ins.Table.Equals(other.Table)
}

// CreateOriginal method needs to be implemented in order to implement Query interface.
// Returns original SQL query representing data in current Query.
func (ins InsertIntoQuery) CreateOriginal() string {
	result := fmt.Sprintf("INSERT INTO %s %s VALUES %s;",
		ins.Table.Value,
		t.Bracketize(ins.ColumnNames),
		t.Bracketize(ins.Values),
	)
	return result
}
