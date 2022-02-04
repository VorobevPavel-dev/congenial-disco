package parser

import (
	"encoding/json"
	"fmt"

	"github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

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
