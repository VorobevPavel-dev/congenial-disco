package main

import (
	"testing"

	"github.com/VorobevPavel-dev/congenial-disco/table"
)

func TestCommandSequenceExecution(t *testing.T) {
	t.Run("Check table creation sequence", func(t *testing.T) {
		// Init session
		session := InitSession()
		// Create table and check if table was appended to session tables map
		// Also check if columns were appended
		createTableCommand := "CREATE TABLE test (id INT, name TEXT);"
		expectedNumOfColumns := 2
		_, _, err := session.ExecuteCommand(createTableCommand)
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
	})
}
