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

func parseInsertIntoStatement(tokens []*t.Token) (*InsertIntoQuery, error) {
	var (
		columnNames []*t.Token
		values      []*t.Token
		table       *t.Token
	)

	currentToken := 0

	//Process INSERT INTO sequense
	if err := AssertTokenSequence(tokens[:2], []*t.Token{
		{Value: "insert", Kind: t.KeywordKind},
		{Value: "into", Kind: t.KeywordKind},
	}); err != nil {
		return nil, fmt.Errorf("cannot process INSERT INTO sequence, err: %v", err)
	}
	currentToken += 2

	//Process table name
	if tokens[currentToken].Equals(t.Reserved[t.SymbolKind]["("]) {
		return nil, fmt.Errorf("expected table name at %d", tokens[currentToken].Position)
	}
	table = tokens[currentToken]

	currentToken++

	//Situation if column names specified
	if tokens[currentToken].Equals(t.Reserved[t.SymbolKind]["("]) {
		currentToken++
		for !tokens[currentToken].Equals(t.Reserved[t.SymbolKind][")"]) {
			if tokens[currentToken].Equals(t.Reserved[t.SymbolKind][","]) {
				currentToken++
				continue
			}
			if currentToken == len(tokens) {
				return nil, fmt.Errorf("expected \")\" symbol at %d", tokens[currentToken].Position)
			}
			tempToken := tokens[currentToken]
			if tempToken.Kind != t.IdentifierKind {
				return nil, fmt.Errorf("column names are only can be identifiers, got: %s", tempToken.String())
			}
			columnNames = append(columnNames, tempToken)
			currentToken++
		}
		currentToken++
	}

	// Process VALUES keyword
	if !tokens[currentToken].Equals(t.Reserved[t.KeywordKind]["values"]) {
		return nil, fmt.Errorf("expected VALUES keyword at %d", tokens[currentToken].Position)
	}
	currentToken++

	// Repeat but for values
	if tokens[currentToken].Equals(t.Reserved[t.SymbolKind]["("]) {
		currentToken++
		for !tokens[currentToken].Equals(t.Reserved[t.SymbolKind][")"]) {
			if tokens[currentToken].Equals(t.Reserved[t.SymbolKind][","]) {
				currentToken++
				continue
			}
			if currentToken == len(tokens) {
				return nil, fmt.Errorf("expected \")\" symbol at %d", tokens[currentToken].Position)
			}
			tempToken := tokens[currentToken]
			// if tempToken.Kind != t.IdentifierKind || tempToken.Kind != t.NumericKind {
			// 	return nil, fmt.Errorf("values can be only identifiers or numbers, got: %s", tempToken.String())
			// }
			values = append(values, tempToken)
			currentToken++
		}
		currentToken++
	} else {
		return nil, fmt.Errorf("expected \"(\" symbol at %d", tokens[currentToken].Position)
	}
	return &InsertIntoQuery{
		Table:       table,
		ColumnNames: columnNames,
		Values:      values,
	}, nil
}
