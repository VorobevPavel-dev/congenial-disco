package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/VorobevPavel-dev/congenial-disco/parser"
	"github.com/VorobevPavel-dev/congenial-disco/table"
	"github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

// Session is a struct describing tables inside current run of engine
type Session struct {
	tables map[string]table.Table
}

// InitSession creates an initial map for tables and returns pointer for Session struct
func InitSession() Session {
	return Session{
		tables: make(map[string]table.Table),
	}
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
	statement := parser.Parse(strings.ToLower(request))
	if statement == nil {
		return "", errors.New("current command is not supported. Only CREATE TABLE, SHOW CREATE(), INSERT INTO, SELECT")
	}
	switch statement.Type {
	case parser.ShowCreateType:
		if val, ok := s.tables[statement.ShowCreateStatement.TableName.Value]; ok {
			return val.ShowCreate(), nil
		}
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
	return "", errors.New("current command is not supported. Only CREATE TABLE, SHOW CREATE(), INSERT INTO, SELECT")
}

func (s *Session) executeCreate(statement *parser.Statement) error {
	t := table.Table(table.LinearTable{})
	t, tn, err := t.Create(statement.CreateTableStatement)
	if err != nil {
		return err
	}
	s.tables[tn] = t
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
