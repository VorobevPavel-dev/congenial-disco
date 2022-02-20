package parser

import (
	"encoding/json"
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
