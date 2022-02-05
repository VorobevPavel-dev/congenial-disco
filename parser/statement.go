package parser

import "github.com/VorobevPavel-dev/congenial-disco/tokenizer"

type Statement struct {
	SelectStatement      *SelectStatement
	CreateTableStatement *CreateTableStatement
	InsertStatement      *InsertStatement
}

// Parse will try to parse statement with all parsers successively
// Returns a Statement struct with only one field not null
func Parse(request string) *Statement {
	// Implement request string as a series of tokens
	tokens := *tokenizer.ParseTokenSequence(request)

	createStatement, _ := parseCreateTableStatement(tokens)
	if createStatement != nil {
		return &Statement{
			CreateTableStatement: createStatement,
			InsertStatement:      nil,
			SelectStatement:      nil,
		}
	}
	insertStatement, _ := parseInsertIntoStatement(tokens)
	if insertStatement != nil {
		return &Statement{
			InsertStatement:      insertStatement,
			CreateTableStatement: nil,
			SelectStatement:      nil,
		}
	}
	selectStatement, _ := parseSelectStatement(tokens)
	if selectStatement != nil {
		return &Statement{
			SelectStatement:      selectStatement,
			CreateTableStatement: nil,
			InsertStatement:      nil,
		}
	}
	return nil
}
