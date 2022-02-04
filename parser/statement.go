package parser

import (
	"encoding/json"
	"fmt"

	"github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

// Definition of structures for Insert statements

type InsertStatement struct {
	Table       tokenizer.Token
	ColumnNames *[]tokenizer.Token
	Values      *[]string
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

	//Process SELECT keyword
	if !tokens[0].Equals(tokenizer.TokenFromKeyword("select")) {
		return nil, fmt.Errorf("expected SELECT keyword at %d", tokens[0].Position)
	}

	//Process columns
	//	Get FROM token position
	fromPosition := tokenizer.FindToken(tokens, tokenizer.TokenFromKeyword("from"))
	if fromPosition == -1 {
		return nil, fmt.Errorf("cannot find FROM keword in request")
	}
	//	Parse identifiers in loop
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
	if items == nil {
		return nil, fmt.Errorf("no identifiers provided for select")
	}

	//Process table name
	//	Check if ";" exists
	endPosition := tokenizer.FindToken(tokens, tokenizer.TokenFromSymbol(";"))
	if endPosition == -1 {
		return nil, fmt.Errorf("cannot find \";\"  in the end of request")
	}
	//	Check if there is a token between FROM and ; tokens
	if endPosition == fromPosition {
		return nil, fmt.Errorf("no table name provided in request")
	}
	tableToken := tokens[fromPosition+1]

	return &SelectStatement{
		Item: items,
		From: *tableToken,
	}, nil
}

// func parseInsertIntoStatement(tokens []*tokenizer.Token) (*InsertStatement, error) {
// 	// INSERT INTO table_name (column1, column2, column3, ...)
// 	// VALUES (value1, value2, value3, ...);
// 	// INSERT INTO table_name
// 	// VALUES (value1, value2, value3, ...);

// 	var (
// 		columnNames []tokenizer.Token
// 		values      []string
// 	)

// 	// if !tokens[0].Equals(tokenizer.TokenFromKeyword("insert")) {
// 	// 	return nil, fmt.Errorf("expected INSERT keyword at %d", tokens[0].Position)
// 	// }

// 	// if !tokens[1].Equals(tokenizer.TokenFromKeyword("into")) {
// 	// 	return nil, fmt.Errorf("expected INTO keyword at %d", tokens[0].Position)
// 	// }

// 	// tableName := tokens[2]

// 	// //Processs string from start to VALUES.pos-1 token
// 	// openBracketColumns := tokenizer.FindToken(tokens, tokenizer.TokenFromSymbol("("))
// 	// closeBracketColumns := tokenizer.FindToken(tokens, tokenizer.TokenFromSymbol(")"))
// 	// if !tokens[3].Equals(tokenizer.TokenFromKeyword("values")) {
// 	// 	for _, item := range tokens[openBracketColumns:closeBracketColumns] {
// 	// 		if item.Equals(tokenizer.TokenFromSymbol(",")) ||
// 	// 			item.Equals(tokenizer.TokenFromSymbol(" ")) {
// 	// 			continue
// 	// 		}
// 	// 		// if current token is a name
// 	// 		if item.Kind != tokenizer.IdentifierKind {
// 	// 			return nil, fmt.Errorf("only Identifiers allowed to be SELECTed")
// 	// 		}
// 	// 		columnNames = append(columnNames, *item)
// 	// 	}

// 	// }

// 	// openBracketColumns := tokenizer.FindToken(tokens, tokenizer.TokenFromSymbol("("))

// 	start, end := 0, len(tokens)
// 	// process "INSERT INTO"
// }
