package main

import (
	"errors"

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
func (s *Session) ExecuteCommand(request string) error {
	statement := parser.Parse(request)
	// statement can have only one non-null field
	if (*statement).CreateTableStatement != nil {
		t := table.LinearTable{}
		tn, _ := t.Create(statement.CreateTableStatement)
		s.tables[tn] = t
		return nil
	} else if (*statement).ShowCreateStatement != nil {
		// Check if table exists
		if val, ok := s.tables[statement.ShowCreateStatement.TableName.Value]; ok {
			val.ShowCreate()
		}
	}
	return errors.New("current command is not supported. Only CREATE TABLE, SHOW CREATE ()")
}
