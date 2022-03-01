package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/VorobevPavel-dev/congenial-disco/parser"
	"github.com/VorobevPavel-dev/congenial-disco/table"
	"github.com/VorobevPavel-dev/congenial-disco/table/column"
	"github.com/VorobevPavel-dev/congenial-disco/table/linear"
	"github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

// Session is a struct describing tables inside current run of engine
type Session struct {
	tables map[string]table.Table
}

// InitSession creates an initial map for tables and returns pointer for Session struct
func InitSession() Session {
	session := Session{
		tables: make(map[string]table.Table),
	}
	session.tables["system.tables"] = linear.LinearTable{
		Columns: []parser.ColumnDefinition{
			{
				Name: &tokenizer.Token{
					Value: "table_name",
					Kind:  tokenizer.IdentifierKind,
				},
				Datatype: &tokenizer.Token{
					Value: "text",
					Kind:  tokenizer.IdentifierKind,
				},
			},
			{
				Name: &tokenizer.Token{
					Value: "engine_type",
					Kind:  tokenizer.IdentifierKind,
				},
				Datatype: &tokenizer.Token{
					Value: "text",
					Kind:  tokenizer.IdentifierKind,
				},
			},
		},
		Name: &tokenizer.Token{
			Value: "system.tables",
			Kind:  tokenizer.IdentifierKind,
		},
		Elements: [][]*tokenizer.Token{{
			{Value: "system.tables", Kind: tokenizer.IdentifierKind},
			{Value: "\"linear\"", Kind: tokenizer.IdentifierKind},
		}},
	}
	return session
}

// CountTables returns number of table in current run of engine
func (s *Session) CountTables() int {
	return len(s.tables)
}

// ToString returns current state of session in JSON format where keys are table names and values are
// number of rows inside.
func (s *Session) ToString() string {
	mapping := make(map[string]int)
	for name, table := range s.tables {
		mapping[name] = table.Count()
	}
	data, _ := json.Marshal(mapping)
	return string(data)
}

// ExecuteCommand parses input string into struct implementing Query interface
// and executes query in engine. Return
//		string. If request returns string value it will be returned here
//		error
func (s *Session) ExecuteCommand(request string) (string, error) {
	statement, err := parser.Parse(strings.ToLower(request))
	if err != nil {
		return err.Error(), err
	}
	switch statement.Type {
	case parser.CreateTableType:
		err := s.executeCreate(statement)
		if err != nil {
			return fmt.Sprint(err), err
		} else {
			return "ok", nil
		}
	case parser.InsertType:
		err := s.executeInsert(statement)
		if err != nil {
			return fmt.Sprint(err), err
		}
		return "ok", nil
	case parser.SelectType:
		result, err := s.executeSelect(statement)
		if err != nil {
			return fmt.Sprint(err), err
		}
		return result, nil
	}
	return "", errors.New("current command is not supported. Only CREATE TABLE, INSERT INTO, SELECT")
}

func (s *Session) executeCreate(statement *parser.Statement) error {
	var (
		table_name string
		engine     string
		t          table.Table
		tn         string
		err        error
	)
	engine = statement.CreateTableStatement.Engine.Value
	switch engine {
	case "linear":
		t = table.Table(linear.LinearTable{})
	case "column":
		t = table.Table(column.ColumnTable{})
	default:
		return fmt.Errorf("current engine %s is not supported", engine)
	}
	t, tn, err = t.Create(statement.CreateTableStatement)
	if err != nil {
		return err
	}
	table_name = statement.CreateTableStatement.Name.Value
	s.tables[tn] = t

	// Add table to system.tables
	request, _ := parser.Parse(fmt.Sprintf("INSERT INTO system.tables VALUES (%s, \"%s\");", table_name, engine))
	s.executeInsert(request)
	return nil
}

func (s *Session) executeInsert(statement *parser.Statement) error {
	desiredTableName := statement.InsertStatement.Table.Value
	// Check if needed table actually exists
	if _, ok := s.tables[desiredTableName]; !ok {
		return fmt.Errorf("table %s does not exist", desiredTableName)
	}
	t, err := s.tables[desiredTableName].Insert(statement.InsertStatement)
	if err != nil {
		return err
	}
	s.tables[desiredTableName] = t
	return nil
}

func (s *Session) executeSelect(statement *parser.Statement) (string, error) {
	desiredTableName := statement.SelectStatement.From.Value
	if _, ok := s.tables[desiredTableName]; !ok {
		return "", fmt.Errorf("table %s does not exist", desiredTableName)
	}
	result, err := s.tables[desiredTableName].Select(statement.SelectStatement)
	if err != nil {
		return "", err
	}
	return formatCSV(result), nil
}

func formatCSV(input [][]*tokenizer.Token) string {
	var result string
	for _, line := range input {
		// extract values
		values := make([]string, len(line))
		for i := range values {
			values[i] = line[i].Value
		}
		result += strings.Join(values, ",") + "\n"
	}
	return result
}
