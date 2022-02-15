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

// parseShowCreateQuery will process set of tokens and try to parse it like SHOW CREATE
// query with syntax described in comments to ShowCreateQuery struct.
// Input set must have SymbolToken with value ';' at the end in order to be parsed.
func parseShowCreateQuery(tokens []*t.Token) (*ShowCreateQuery, error) {
	// Validate that set of tokens has ';' SymbolKind token at the end
	if !tokens[len(tokens)-1].Equals(t.TokenFromSymbol(";")) {
		return nil, ErrNoSemicolonAtTheEnd
	}

	var tableName *t.Token

	// Process SHOW CREATE ( sequence
	cursor := 0

	if err := AssertTokenSequence(tokens[:3], []*t.Token{
		t.Reserved[t.KeywordKind]["show"],
		t.Reserved[t.KeywordKind]["create"],
		t.Reserved[t.SymbolKind]["("],
	}); err != nil {
		return nil, fmt.Errorf("cannot process \"SHOW CREATE (\" sequence, err: %v", err)
	}
	cursor += 3

	if tokens[cursor].Kind != t.IdentifierKind {
		return nil, fmt.Errorf("table names are only can be identifiers, got: %s", tokens[cursor].String())
	}
	tableName = tokens[cursor]
	cursor++

	//process ); sequence
	if err := AssertTokenSequence(tokens[cursor:], []*t.Token{
		t.Reserved[t.SymbolKind][")"],
		t.Reserved[t.SymbolKind][";"],
	}); err != nil {
		return nil, fmt.Errorf("cannot process \");\" sequence, err: %v", err)
	}
	return &ShowCreateQuery{
		TableName: tableName,
	}, nil
}
