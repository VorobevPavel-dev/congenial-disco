package parser

import (
	"encoding/json"
	"fmt"

	t "github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

// ShowCreateQuery is an utility function
// describing with what query table was created
// Syntax:
//
// SHOW CREATE (<table_name>);
//
//where <table_name> must be an IdentifierKind Token
type ShowCreateQuery struct {
	TableName *t.Token
}

// Equals method needs to be implemented in order to implement Query interface.
// Returns true if table names of both queries are equal
func (scs ShowCreateQuery) Equals(other *ShowCreateQuery) bool {
	return scs.TableName.Equals(other.TableName)
}

// String method needs to be implemented in order to implement Query interface.
// Returns JSON object describing necessary information
func (scs ShowCreateQuery) String() string {
	bytes, _ := json.Marshal(scs)
	return string(bytes)
}

// CreateOriginal method needs to be implemented in order to implement Query interface.
// Returns original SQL query representing data in current Query
func (scs ShowCreateQuery) CreateOriginal() string {
	return fmt.Sprintf("SHOW CREATE (%s);", scs.TableName.Value)
}
