package parser

import (
	"encoding/json"
	"fmt"

	"github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

type InsertStatement struct {
	Table       tokenizer.Token
	ColumnNames []*tokenizer.Token
	Values      []*tokenizer.Token
}

func (ins *InsertStatement) String() string {
	bytes, _ := json.Marshal(ins)
	return string(bytes)
}

func (ins *InsertStatement) Equals(other *InsertStatement) bool {
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
	return ins.Table.Equals(&other.Table)
}

func parseInsertIntoStatement(tokens []*tokenizer.Token) (*InsertStatement, error) {
	var (
		columnNames []*tokenizer.Token
		values      []*tokenizer.Token
		table       tokenizer.Token
	)

	currentToken := 0

	//Process INSERT INTO sequense
	if !tokens[currentToken].Equals(tokenizer.TokenFromKeyword("insert")) {
		return nil, fmt.Errorf("expected INSERT keyword at %d", tokens[currentToken].Position)
	}
	currentToken++
	if !tokens[currentToken].Equals(tokenizer.TokenFromKeyword("into")) {
		return nil, fmt.Errorf("expected INTO keyword at %d", tokens[currentToken].Position)
	}
	currentToken++

	//Process table name
	if tokens[currentToken].Equals(tokenizer.TokenFromSymbol("(")) {
		return nil, fmt.Errorf("expected table name at %d", tokens[currentToken].Position)
	} else {
		table = *tokens[currentToken]
	}

	currentToken++

	//Situation if column names specified
	if tokens[currentToken].Equals(tokenizer.TokenFromSymbol("(")) {
		currentToken++
		for !tokens[currentToken].Equals(tokenizer.TokenFromSymbol(")")) {
			if tokens[currentToken].Equals(tokenizer.TokenFromSymbol(",")) {
				continue
			}
			if currentToken == len(tokens) {
				return nil, fmt.Errorf("expected \")\" symbol at %d", tokens[currentToken].Position)
			}
			tempToken := tokens[currentToken]
			if tempToken.Kind != tokenizer.IdentifierKind {
				return nil, fmt.Errorf("column names are only can be identifiers, got: %s", tempToken.String())
			}
			columnNames = append(columnNames, tempToken)
			currentToken++
		}
		currentToken++
	}

	// Process VALUES keyword
	if !tokens[currentToken].Equals(tokenizer.TokenFromKeyword("values")) {
		return nil, fmt.Errorf("expected VALUES keyword at %d", tokens[currentToken].Position)
	}
	currentToken++

	// Repeat but for values
	if tokens[currentToken].Equals(tokenizer.TokenFromSymbol("(")) {
		currentToken++
		for !tokens[currentToken].Equals(tokenizer.TokenFromSymbol(")")) {
			if tokens[currentToken].Equals(tokenizer.TokenFromSymbol(",")) {
				currentToken++
				continue
			}
			if currentToken == len(tokens) {
				return nil, fmt.Errorf("expected \")\" symbol at %d", tokens[currentToken].Position)
			}
			tempToken := tokens[currentToken]
			// if tempToken.Kind != tokenizer.IdentifierKind || tempToken.Kind != tokenizer.NumericKind {
			// 	return nil, fmt.Errorf("values can be only identifiers or numbers, got: %s", tempToken.String())
			// }
			values = append(values, tempToken)
			currentToken++
		}
		currentToken++
	} else {
		return nil, fmt.Errorf("expected \"(\" symbol at %d", tokens[currentToken].Position)
	}
	return &InsertStatement{
		Table:       table,
		ColumnNames: columnNames,
		Values:      values,
	}, nil
}
