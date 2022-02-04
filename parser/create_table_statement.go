package parser

import (
	"encoding/json"
	"fmt"

	"github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

type ColumnDefinition struct {
	Name     tokenizer.Token
	Datatype tokenizer.Token
}

func (cd *ColumnDefinition) Equals(other *ColumnDefinition) bool {
	return cd.Name.Equals(&other.Name) && cd.Datatype.Equals(&other.Datatype)
}

func (ct *CreateTableStatement) String() string {
	bytes, _ := json.Marshal(ct)
	return string(bytes)
}

type CreateTableStatement struct {
	Name tokenizer.Token
	Cols []*ColumnDefinition
}

func (ct *CreateTableStatement) Equals(other *CreateTableStatement) bool {
	if len(ct.Cols) != len(other.Cols) {
		return false
	}
	for index := range ct.Cols {
		if !ct.Cols[index].Equals(other.Cols[index]) {
			return false
		}
	}
	return ct.Name.Equals(&other.Name)
}

func parseCreateTableStatement(tokens []*tokenizer.Token) (*CreateTableStatement, error) {
	// CREATE TABLE table_name (
	// 	column1 datatype,
	// 	column2 datatype,
	// 	column3 datatype,
	//    ....
	// );

	var (
		tableName *tokenizer.Token
		columns   []*ColumnDefinition
	)

	currentToken := 0

	// Process CREATE TABLE sequence
	if !tokens[currentToken].Equals(tokenizer.TokenFromKeyword("create")) {
		return nil, fmt.Errorf("expected CREATE keyword at %d", tokens[currentToken].Position)
	}
	currentToken++
	if !tokens[currentToken].Equals(tokenizer.TokenFromKeyword("table")) {
		return nil, fmt.Errorf("expected TABLE keyword at %d", tokens[currentToken].Position)
	}
	currentToken++

	// Process table name
	if tokens[currentToken].Equals(tokenizer.TokenFromSymbol("(")) {
		return nil, fmt.Errorf("expected \"(\" keyword at %d", tokens[currentToken].Position)
	}
	if tokens[currentToken].Kind != tokenizer.IdentifierKind {
		return nil, fmt.Errorf("expected table name identifier at %d, got: %s", tokens[currentToken].Position, tokens[currentToken].String())
	}
	tableName = tokens[currentToken]
	currentToken++

	// Process set of column definitions
	if !tokens[currentToken].Equals(tokenizer.TokenFromSymbol("(")) {
		return nil, fmt.Errorf("expected \"(\" symbol at %d", tokens[currentToken].Position)
	}
	currentToken++
	for !tokens[currentToken].Equals(tokenizer.TokenFromSymbol(")")) {
		if currentToken == len(tokens)-1 {
			return nil, fmt.Errorf("expected \")\" symbol at %d", tokens[currentToken].Position)
		}
		if tokens[currentToken].Equals(tokenizer.TokenFromSymbol(",")) {
			currentToken++
			continue
		}
		// Process column name
		if tokens[currentToken].Kind != tokenizer.IdentifierKind {
			return nil, fmt.Errorf("column names are only can be identifiers, got: %s", tokens[currentToken].String())
		}
		columnName := tokens[currentToken]
		currentToken++

		// Process column type
		if tokens[currentToken].Kind != tokenizer.TypeKind {
			return nil, fmt.Errorf("expected type of column, got: %s", tokens[currentToken].String())
		}
		columnType := tokens[currentToken]
		columns = append(columns, &ColumnDefinition{Name: *columnName, Datatype: *columnType})
		currentToken++
	}
	currentToken++
	if currentToken == len(tokens) {
		return nil, fmt.Errorf("expected \";\" symbol at the end of request")
	}

	return &CreateTableStatement{
		Name: *tableName,
		Cols: columns,
	}, nil
}
