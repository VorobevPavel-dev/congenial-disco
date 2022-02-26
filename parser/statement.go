package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"

	t "github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

var (
	ErrNoSemicolonAtTheEnd     error  = errors.New("provided request does not have ';' SymbolKind token at the end")
	ErrExpectedKeywordTemplate string = "expected %s keyword at %d"

	// ErrorBuilders are functions deicated to process data from parsing functions and return formatted error
	// There are only functions and errors required by all parsers. All specific errors declarated in parser
	// itself

	// ErrExpectedToken builds error according to desired token and position where that token must be
	ErrExpectedToken = func(e *t.Token, p int) error {
		return fmt.Errorf("expected %s %s at %d",
			e.Value,
			t.KindToString(e.Kind),
			p,
		)
	}

	// ErrInvalidTokenKind builds error according to current token and desired TokenKind
	ErrInvalidTokenKind = func(e *t.Token, ek t.TokenKind) error {
		return fmt.Errorf("expected %s but got %s at %d",
			t.KindToString(e.Kind),
			t.KindToString(ek),
			e.Position,
		)
	}
)

type queryKind int

const (
	CreateTableType queryKind = iota
	SelectType
	InsertType
	ShowCreateType
)

type Statement struct {
	SelectStatement      *SelectQuery
	CreateTableStatement *CreateTableQuery
	InsertStatement      *InsertIntoQuery
	ShowCreateStatement  *ShowCreateQuery
	// Experimental
	Type queryKind
}

// Experimental
type Query interface {
	// String() string
	Equals(*Query) bool
	Parse([]*t.Token) (*Query, bool, error)
	// CreateOriginal must return string containing original SQL request
	CreateOriginal() string
}

// Experimental
func QueryToString(q *Query) string {
	bytes, _ := json.Marshal(q)
	return string(bytes)
}

// Parse will try to parse statement with all parsers successively
// Returns a Statement struct with only one field not null
func Parse(request string) (resultStatement *Statement, err error) {

	// handling unexpected out-of-bounds
	defer func(err *error) {
		if panicError := recover(); panicError != nil {
			*err = fmt.Errorf("panic occured: %v", panicError)
		}
	}(&err)

	// Implement request string as a series of tokens
	tokens := *t.ParseTokenSequence(request)
	var (
		cursor int = 0
	)

	switch tokens[cursor].Value {
	case "create":
		cursor++
		query, err := parseCreateTableBranch(tokens, &cursor)
		if err != nil {
			return nil, err
		}
		resultStatement = &Statement{
			CreateTableStatement: query,
			Type:                 CreateTableType,
		}
		err = nil
	case "insert":
		cursor++
		query, err := parseInsertIntoBranch(tokens, &cursor)
		if err != nil {
			return nil, err
		}
		resultStatement = &Statement{
			InsertStatement: query,
			Type:            InsertType,
		}
		err = nil
	case "select":
		cursor++
		query, err := parseSelectBranch(tokens, &cursor)
		if err != nil {
			return nil, err
		}
		resultStatement = &Statement{
			SelectStatement: query,
			Type:            SelectType,
		}
		err = nil
	case "show":
		cursor++
		query, err := parseShowCreateBranch(tokens, &cursor)
		if err != nil {
			return nil, err
		}
		resultStatement = &Statement{
			ShowCreateStatement: query,
			Type:                ShowCreateType,
		}
		err = nil
	default:
		err = fmt.Errorf("current operation is not supported")
	}
	return resultStatement, err
}

func parseCreateTableBranch(tokens []*t.Token, cursor *int) (*CreateTableQuery, error) {
	var (
		parsingInProgress       bool               = true
		step                    parsingStep        = stepTableKeyword
		tableName               *t.Token           = nil
		columnDefinitions       []ColumnDefinition = []ColumnDefinition{}
		currentColumnDefinition ColumnDefinition   = ColumnDefinition{}
		engine                  *t.Token           = nil
	)
	for *cursor < len(tokens) && parsingInProgress {
		switch step {
		case stepTableKeyword:
			if !tokens[*cursor].Equals(t.Reserved[t.KeywordKind]["table"]) {
				return nil, fmt.Errorf("expected table keyword at %d", tokens[*cursor].Position)
			}
			*cursor++
			step = stepTableName
			continue
		case stepTableName:
			if tokens[*cursor].Kind != t.IdentifierKind {
				return nil, fmt.Errorf("expected table name at %d", tokens[*cursor].Position)
			}
			tableName = tokens[*cursor]
			*cursor++
			if !tokens[*cursor].Equals(t.Reserved[t.SymbolKind]["("]) {
				return nil, fmt.Errorf("expected \"(\" symbol at %d", tokens[*cursor].Position)
			} else {
				*cursor++
				step = stepColumnName
				continue
			}
		case stepColumnName:
			if tokens[*cursor].Kind != t.IdentifierKind {
				return nil, fmt.Errorf("expected column name at %d", tokens[*cursor].Position)
			}
			currentColumnDefinition.Name = tokens[*cursor]
			*cursor++
			step = stepColumnType
		case stepColumnType:
			if tokens[*cursor].Kind != t.TypeKind {
				return nil, fmt.Errorf("expected column type at %d", tokens[*cursor].Position)
			}
			currentColumnDefinition.Datatype = tokens[*cursor]
			columnDefinitions = append(columnDefinitions, currentColumnDefinition)
			currentColumnDefinition = ColumnDefinition{}
			*cursor++
			if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][")"]) {
				step = stepEngineKeyword
				*cursor++
			} else if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][","]) {
				step = stepColumnName
				*cursor++
			} else {
				return nil, fmt.Errorf("expected comma or closed bracket at %d", tokens[*cursor].Position)
			}
			continue
		case stepEngineKeyword:
			if !tokens[*cursor].Equals(t.Reserved[t.KeywordKind]["engine"]) {
				return nil, fmt.Errorf("expected ENGINE keyword at %d", tokens[*cursor].Position)
			}
			*cursor++
			step = stepEngineName
			continue
		case stepEngineName:
			if tokens[*cursor].Kind != t.EngineKind {
				return nil, fmt.Errorf("expected engine name at %d", tokens[*cursor].Position)
			}
			engine = tokens[*cursor]
			*cursor++
			if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][";"]) {
				parsingInProgress = false
			} else {
				return nil, fmt.Errorf("expected SETTINGS or \";\" symbol at %d", tokens[*cursor].Position)
			}
		}
	}
	return &CreateTableQuery{
		Name:   tableName,
		Cols:   &columnDefinitions,
		Engine: engine,
	}, nil
}

func parseInsertIntoBranch(tokens []*t.Token, cursor *int) (*InsertIntoQuery, error) {
	var (
		parsingInProgress bool         = true
		step              parsingStep  = stepInsIntoKeyword
		tableName         *t.Token     = nil
		values            [][]*t.Token = [][]*t.Token{}
		currentValueSet   []*t.Token
	)
	for *cursor < len(tokens) && parsingInProgress {
		switch step {
		case stepInsIntoKeyword:
			if !tokens[*cursor].Equals(t.Reserved[t.KeywordKind]["into"]) {
				return nil, fmt.Errorf("expected into keyword at %d", tokens[*cursor].Position)
			}
			*cursor++
			step = stepInsTableName
			continue
		case stepInsTableName:
			if tokens[*cursor].Kind != t.IdentifierKind {
				return nil, fmt.Errorf("expected table name at %d", tokens[*cursor].Position)
			}
			tableName = tokens[*cursor]
			*cursor++
			step = stepInsValuesKeyword
		case stepInsValuesKeyword:
			if !tokens[*cursor].Equals(t.Reserved[t.KeywordKind]["values"]) {
				return nil, fmt.Errorf("expected values keyword at %d", tokens[*cursor].Position)
			}
			*cursor++
			step = stepInsValuesetOpenBracket
			continue
		case stepInsValuesetOpenBracket:
			if !tokens[*cursor].Equals(t.Reserved[t.SymbolKind]["("]) {
				return nil, fmt.Errorf("expected opening bracket at %d", tokens[*cursor].Position)
			}
			*cursor++
			step = stepInsValuesetValue
			currentValueSet = []*t.Token{}
		case stepInsValuesetValue:
			if !((tokens[*cursor].Kind == t.IdentifierKind) || (tokens[*cursor].Kind == t.NumericKind)) {
				return nil, fmt.Errorf("values can be only identifiers or numbers but got %s at %d", tokens[*cursor].Value, tokens[*cursor].Position)
			}
			currentValueSet = append(currentValueSet, tokens[*cursor])
			*cursor++
			if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][")"]) {
				step = stepInsValuesetCloseBracket
			} else if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][","]) {
				step = stepInsValuesetValue
				*cursor++
			} else {
				return nil, fmt.Errorf("expected comma ar closing bracket at %d", tokens[*cursor].Position)
			}
		case stepInsValuesetCloseBracket:
			if !(tokens[*cursor].Equals(t.Reserved[t.SymbolKind][")"])) {
				return nil, fmt.Errorf("expected closing bracket at %d", tokens[*cursor].Position)
			}
			values = append(values, currentValueSet)
			*cursor++
			if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][","]) {
				step = stepInsValuesetOpenBracket
				*cursor++
			} else {
				step = stepEnd
			}
		case stepEnd:
			if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][";"]) {
				return &InsertIntoQuery{
					Table:  tableName,
					Values: values,
				}, nil
			} else {
				return nil, fmt.Errorf("expected \";\" at %d", tokens[*cursor].Position)
			}
		}
	}
	return nil, fmt.Errorf("cannot parse query as INSERT INTO query")
}

func parseSelectBranch(tokens []*t.Token, cursor *int) (*SelectQuery, error) {
	var (
		columns   []*t.Token  = []*t.Token{}
		tableName *t.Token    = nil
		step      parsingStep = stepSelColName
	)
	// Initial step assertion
	if !tokens[*cursor].Equals(t.Reserved[t.SymbolKind]["("]) {
		return nil, fmt.Errorf("expected \"(\" symbol at %d", tokens[*cursor].Position)
	}
	*cursor++
	for *cursor < len(tokens) {
		switch step {
		case stepSelColName:
			if tokens[*cursor].Kind != t.IdentifierKind {
				return nil, fmt.Errorf("expected column name at %d", tokens[*cursor].Position)
			}
			columns = append(columns, tokens[*cursor])
			*cursor++
			if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][","]) {
				*cursor++
				continue
			} else if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][")"]) {
				*cursor++
				step = stepSelFromKeyword
				continue
			} else {
				return nil, fmt.Errorf("expected \")\" or comma at %d", tokens[*cursor].Position)
			}
		case stepSelFromKeyword:
			if !tokens[*cursor].Equals(t.Reserved[t.KeywordKind]["from"]) {
				return nil, fmt.Errorf("expected from keyword at %d", tokens[*cursor].Position)
			}
			*cursor++
			step = stepSelTableName
			continue
		case stepSelTableName:
			if tokens[*cursor].Kind != t.IdentifierKind {
				return nil, fmt.Errorf("expected table name at %d", tokens[*cursor].Position)
			}
			tableName = tokens[*cursor]
			*cursor++
			step = stepEnd
		case stepEnd:
			if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][";"]) {
				return &SelectQuery{
					Columns: columns,
					From:    tableName,
				}, nil
			} else {
				return nil, fmt.Errorf("expected \";\" at %d", tokens[*cursor].Position)
			}
		}
	}
	return nil, fmt.Errorf("cannot parse query as SELECT FROM query")
}

func parseShowCreateBranch(tokens []*t.Token, cursor *int) (*ShowCreateQuery, error) {
	var (
		tableName *t.Token    = nil
		step      parsingStep = stepShCreateKeyword
	)
	for *cursor < len(tokens) {
		switch step {
		case stepShCreateKeyword:
			if !tokens[*cursor].Equals(t.Reserved[t.KeywordKind]["create"]) {
				return nil, fmt.Errorf("expected create keyword at %d", tokens[*cursor].Position)
			}
			*cursor++
			if tokens[*cursor].Equals(t.Reserved[t.SymbolKind]["("]) {
				*cursor++
				step = stepShTableName
			} else {
				return nil, fmt.Errorf("expected \"(\" symbol at %d", tokens[*cursor].Position)
			}
		case stepShTableName:
			if tokens[*cursor].Kind != t.IdentifierKind {
				return nil, fmt.Errorf("expected table name at %d", tokens[*cursor].Position)
			}
			tableName = tokens[*cursor]
			*cursor++
			if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][")"]) {
				*cursor++
				step = stepEnd
			} else {
				return nil, fmt.Errorf("expected \")\" symbol at %d", tokens[*cursor].Position)
			}
		case stepEnd:
			if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][";"]) {
				return &ShowCreateQuery{
					TableName: tableName,
				}, nil
			} else {
				return nil, fmt.Errorf("expected \";\" at %d", tokens[*cursor].Position)
			}
		}
	}
	return nil, fmt.Errorf("cannot parse query as SHOW CREATE query")
}

//getRandomFromKinds returns you random element of given kinds
func getRandomFromKinds(kinds ...t.TokenKind) t.TokenKind {
	return kinds[rand.Intn(len(kinds))]
}

func GenerateStatement(kind queryKind) (*Statement, string) {
	switch kind {
	case InsertType:
		var (
			values [][]*t.Token
		)
		tableName := t.GenerateRandomToken(t.IdentifierKind)
		columnCount := rand.Intn(10) + 1
		// number of inserting sets
		for i := 0; i < rand.Intn(5)+1; i++ {
			tempValues := []*t.Token{}
			// number of values in each set
			for j := 0; j < columnCount; j++ {
				randomKind := getRandomFromKinds(
					t.NumericKind,
					t.IdentifierKind,
				)
				tempValues = append(tempValues, t.GenerateRandomToken(randomKind))
			}
			values = append(values, tempValues)
		}
		inputQuery := InsertIntoQuery{
			Table:  tableName,
			Values: values,
		}
		inputSQL := inputQuery.CreateOriginal()
		return &Statement{
			InsertStatement: &inputQuery,
			Type:            InsertType,
		}, inputSQL
	}
	return nil, ""
}
