package main

import (
	"errors"
	"fmt"
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
// and executes query in engine. Return
//		table.Table. If data inside table was modified returned table must replace old version in s.tables map
//		string. If request returns string value it will be returned here
//		*[][]table.Element
//		error
func (s *Session) ExecuteCommand(request string) (table.Table, string, *[][]table.Element, error) {
	statement := parser.Parse(strings.ToLower(request))
	switch statement.Type {
	case parser.ShowCreateType:
		if val, ok := s.tables[statement.ShowCreateStatement.TableName.Value]; ok {
			return nil, val.ShowCreate(), nil, nil
		}
	case parser.CreateTableType:
		t := table.Table(table.LinearTable{})
		t, tn, err := t.Create(statement.CreateTableStatement)
		if err != nil {
			return nil, "", nil, err
		}
		s.tables[tn] = t
		return nil, tn, nil, nil
	case parser.InsertType:
		desiredTableName := statement.InsertStatement.Table.Value
		// Check if needed table actually exists
		if _, ok := s.tables[desiredTableName]; !ok {
			return nil, "", nil, fmt.Errorf("table %s does not exist", desiredTableName)
		}
		t, err := s.tables[desiredTableName].Insert(statement.InsertStatement)
		if err != nil {
			return nil, "", nil, err
		}
		s.tables[desiredTableName] = t
	default:
		return nil, "", nil, errors.New("current command is not supported. Only CREATE TABLE, SHOW CREATE ()")
	}
	return nil, "", nil, nil
}
