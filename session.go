package main

import (
	"errors"

	"github.com/VorobevPavel-dev/congenial-disco/parser"
	"github.com/VorobevPavel-dev/congenial-disco/table"
)

type Session struct {
	tables map[string]interface{}
}

func InitSession() *Session {
	return &Session{
		tables: make(map[string]interface{}),
	}
}

func (s *Session) CountTables() int {
	return len(s.tables)
}

func (s *Session) ExecuteCommand(request string) error {
	statement := parser.Parse(request)
	// statement can have only one non-null field
	if statement.CreateTableStatement != nil {
		t := table.LinearTable{}
		tn, _ := t.Create(statement.CreateTableStatement)
		s.tables[tn] = t
		return nil
	}
	return errors.New("current command is not supported. Only CREATE TABLE")
}
