package parser

import (
	"encoding/json"
	"errors"
	"fmt"

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
func Parse(request string) (*Statement, error) {
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
		return &Statement{
			CreateTableStatement: query,
			Type:                 CreateTableType,
		}, nil
	case "insert":
		cursor++
		query, err := parseInsertIntoBranch(tokens, &cursor)
		if err != nil {
			return nil, err
		}
		return &Statement{
			InsertStatement: query,
			Type:            InsertType,
		}, nil
	}
	return nil, nil
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
			step = stepColOpenBracket
		case stepColOpenBracket:
			if !tokens[*cursor].Equals(t.Reserved[t.SymbolKind]["("]) {
				return nil, fmt.Errorf("expected \"(\" at %d", tokens[*cursor].Position)
			}
			*cursor++
			step = stepColumnName
			continue
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
			if tokens[*cursor].Equals(t.Reserved[t.KeywordKind]["settings"]) {
				step = stepSettingsName
				*cursor++
				continue
			} else if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][";"]) {
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
		parsingInProgress bool        = true
		step              parsingStep = stepInsIntoKeyword
		tableName         *t.Token    = nil
		columnNames       []*t.Token  = []*t.Token{}
		values            []*t.Token  = []*t.Token{}
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
			if tokens[*cursor].Equals(t.Reserved[t.SymbolKind]["("]) {
				step = stepInsColsetName
				*cursor++
			} else {
				step = stepInsValuesKeyword
				*cursor++
			}
			continue
		case stepInsColsetName:
			if tokens[*cursor].Kind != t.IdentifierKind {
				return nil, fmt.Errorf("expected column name at %d", tokens[*cursor].Position)
			}
			columnNames = append(columnNames, tokens[*cursor])
			*cursor++
			if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][")"]) {
				step = stepInsValuesKeyword
				*cursor++
			} else if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][","]) {
				step = stepInsColsetName
				*cursor++
			} else {
				return nil, fmt.Errorf("exected comma or closing bracket at %d", tokens[*cursor].Position)
			}
			continue
		case stepInsValuesKeyword:
			if !tokens[*cursor].Equals(t.Reserved[t.KeywordKind]["values"]) {
				return nil, fmt.Errorf("expected values keyword at %d", tokens[*cursor].Position)
			}
			*cursor++
			if !tokens[*cursor].Equals(t.Reserved[t.SymbolKind]["("]) {
				return nil, fmt.Errorf("expected opening bracket for values set at %d", tokens[*cursor].Position)
			}
			*cursor++
			step = stepInsValueValue
			continue
		case stepInsValueValue:
			if !((tokens[*cursor].Kind == t.IdentifierKind) || (tokens[*cursor].Kind == t.NumericKind)) {
				return nil, fmt.Errorf("values can be only identifiers or numbers but got %s at %d", tokens[*cursor].Value, tokens[*cursor].Position)
			}
			values = append(values, tokens[*cursor])
			*cursor++
			if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][")"]) {
				step = stepEnd
				*cursor++
			} else if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][","]) {
				step = stepInsValueValue
				*cursor++
			} else {
				return nil, fmt.Errorf("expected comma ar closing bracket at %d", tokens[*cursor].Position)
			}
		case stepEnd:
			if tokens[*cursor].Equals(t.Reserved[t.SymbolKind][";"]) {
				return &InsertIntoQuery{
					Table:       tableName,
					ColumnNames: columnNames,
					Values:      values,
				}, nil
			} else {
				return nil, fmt.Errorf("expected \";\" at %d", tokens[*cursor].Position)
			}
		}
	}
	return nil, fmt.Errorf("cannot parse query as INSERT INTO query")
}
