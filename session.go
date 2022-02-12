package main

import (
	"errors"
	"strings"

	"github.com/VorobevPavel-dev/congenial-disco/parser"
	"github.com/VorobevPavel-dev/congenial-disco/table"
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

// ExecuteCommand parses input string into struct implementing Query interface
// and executes query in engine
func (s *Session) ExecuteCommand(request string) (string, *[][]table.Element, error) {
	statement := parser.Parse(strings.ToLower(request))
	switch statement.Type {
	case parser.ShowCreateType:
		if val, ok := s.tables[statement.ShowCreateStatement.TableName.Value]; ok {
			return val.ShowCreate(), nil, nil
		}
	case parser.CreateTableType:
		t := table.Table(table.LinearTable{})
		t, tn, err := t.Create(statement.CreateTableStatement)
		if err != nil {
			return "", nil, err
		}
		s.tables[tn] = t
		return tn, nil, nil
	default:
		return "", nil, errors.New("current command is not supported. Only CREATE TABLE, SHOW CREATE ()")
	}
	return "", nil, nil
}
