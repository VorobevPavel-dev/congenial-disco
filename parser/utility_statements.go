package parser

import (
	"encoding/json"
	"fmt"

	"github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

type ShowCreateStatement struct {
	TableName tokenizer.Token
}

func (scs *ShowCreateStatement) Equals(other *ShowCreateStatement) bool {
	return scs.TableName.Equals(&other.TableName)
}

func (scs *ShowCreateStatement) String() string {
	bytes, _ := json.Marshal(scs)
	return string(bytes)
}

func parseShowCreateStatement(tokens []*tokenizer.Token) (*ShowCreateStatement, error) {
	// SHOW CREATE (...);

	var tableName *tokenizer.Token

	// Process SHOW CREATE sequence
	currentToken := 0

	if !tokens[currentToken].Equals(tokenizer.TokenFromKeyword("show")) {
		return nil, fmt.Errorf("expected SHOW keyword at %d", tokens[currentToken].Position)
	}
	currentToken++
	if !tokens[currentToken].Equals(tokenizer.TokenFromKeyword("create")) {
		return nil, fmt.Errorf("expected CREATE keyword at %d", tokens[currentToken].Position)
	}
	currentToken++

	// Process set of table name
	if !tokens[currentToken].Equals(tokenizer.TokenFromSymbol("(")) {
		return nil, fmt.Errorf("expected \"(\" symbol at %d", tokens[currentToken].Position)
	}
	currentToken++
	if tokens[currentToken].Kind != tokenizer.IdentifierKind {
		return nil, fmt.Errorf("table names are only can be identifiers, got: %s", tokens[currentToken].String())
	} else {
		tableName = tokens[currentToken]
	}
	currentToken++

	//process ); sequence
	if !tokens[currentToken].Equals(tokenizer.TokenFromSymbol(")")) {
		return nil, fmt.Errorf("expected \")\" symbol at %d", tokens[currentToken].Position)
	}
	currentToken++
	if !tokens[currentToken].Equals(tokenizer.TokenFromSymbol(";")) {
		return nil, fmt.Errorf("expected \";\" symbol at %d", tokens[currentToken].Position)
	}
	currentToken++
	return &ShowCreateStatement{
		TableName: *tableName,
	}, nil
}
