package main

import (
	"testing"

	"github.com/VorobevPavel-dev/congenial-disco/table"
)

func TestCommandSequenceExecution(t *testing.T) {
	// Init session
	session := InitSession()
	t.Run("Check table creation sequence", func(t *testing.T) {
		// Create table and check if table was appended to session tables map
		// Also check if columns were appended
		createTableCommand := "CREATE TABLE test (id INT, name TEXT);"
		expectedNumOfColumns := 2
		_, _, _, err := session.ExecuteCommand(createTableCommand)
		if err != nil {
			t.Error(err)
		}
		// Check if session.tables has "test" table
		if _, ok := session.tables["test"]; !ok {
			t.Error("No required table \"test\" was found in tables map")
		}
		columnNames := session.tables["test"].(table.LinearTable).GetColumnsNames()
		if len(columnNames) != expectedNumOfColumns {
			t.Errorf("Count of columns differ. Expected: %d, got: %v", expectedNumOfColumns, columnNames)
		}
		t.Logf("Created table has columns: %s", session.tables["test"].GetColumns())
	})
	t.Run("Check table insertion", func(t *testing.T) {
		_, _, _, err := session.ExecuteCommand("INSERT INTO test VALUES (1, test);")
		if err != nil {
			t.Error(err)
		}
		if session.tables["test"].Count() != 1 {
			t.Errorf("expected only one row after insertion, got %d", session.tables["test"].Count())
		}
		_, _, _, err = session.ExecuteCommand("INSERT INTO test (id, name) VALUES (1, test_value);")
		if err != nil {
			t.Error(err)
		}
		if session.tables["test"].Count() != 2 {
			t.Errorf("expected only two rows after insertion, got %d", session.tables["test"].Count())
		}

	})
}
