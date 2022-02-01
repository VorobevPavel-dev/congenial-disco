package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

// Definition of structures for Insert statements

type InsertStatement struct {
	Table  tokenizer.Token
	Values *[]string
}

// Definition of structures for Create statements

type ColumnDefinition struct {
	Name     tokenizer.Token
	Datatype tokenizer.Token
}

type CreateTableStatement struct {
	Name tokenizer.Token
	Cols *[]*ColumnDefinition
}

// Definition of structures for Select statements

type SelectStatement struct {
	Item []*tokenizer.Token `json:"item"`
	From tokenizer.Token    `json:"from"`
}

func (slct *SelectStatement) String() string {
	bytes, _ := json.Marshal(slct)
	return string(bytes)
}

func (slct *SelectStatement) Equals(other *SelectStatement) bool {
	if len(slct.Item) != len(other.Item) {
		return false
	}
	for index := range slct.Item {
		if !slct.Item[index].Equals(other.Item[index]) {
			return false
		}
	}
	return slct.From.Equals(&other.From)
}

type Statement struct {
	SelectStatement      *SelectStatement
	CreateTableStatement *CreateTableStatement
	InsertStatement      *InsertStatement
}

func Parse(request string) *Statement {
	return nil
}

func parseSelectStatement(tokens []*tokenizer.Token) (*SelectStatement, error) {
	// SELECT ... FROM table;

	var items []*tokenizer.Token

	if len(tokens) < 7 {
		return nil, errors.New("select statement contains less then 7 elements inside")
	}
	if !tokens[0].Equals(tokenizer.TokenFromKeyword("select")) {
		return nil, fmt.Errorf("expected SELECT keyword at %d", tokens[0].Position)
	}
	//tokens[1] is a space
	fromPosition := tokenizer.FindToken(tokens, tokenizer.TokenFromKeyword("from"))
	if fromPosition == -1 {
		return nil, fmt.Errorf("cannot find FROM keword in request")
	}

	//Parse values to select
	for _, item := range tokens[1:fromPosition] {
		if item.Equals(tokenizer.TokenFromSymbol(",")) ||
			item.Equals(tokenizer.TokenFromSymbol(" ")) {
			continue
		}
		// if current token is a name
		if item.Kind != tokenizer.IdentifierKind {
			return nil, fmt.Errorf("only Identifiers allowed to be SELECTed")
		}
		items = append(items, item)
	}

	//tokens[fromPosition+1] is a space
	//tokens[fromPosition+2] is a table
	//tokens[fromPosition+3] is a ;
	tableToken := tokens[fromPosition+2]

	if !tokens[fromPosition+3].Equals(tokenizer.TokenFromSymbol(";")) {
		return nil, fmt.Errorf("cannot find \";\"  in the end of request")
	}

	return &SelectStatement{
		Item: items,
		From: *tableToken,
	}, nil
}
