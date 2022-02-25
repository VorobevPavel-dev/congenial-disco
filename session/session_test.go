package session

import (
	"math/rand"
	"testing"
	"time"

	"github.com/VorobevPavel-dev/congenial-disco/parser"
	"github.com/VorobevPavel-dev/congenial-disco/table/linear"
	token "github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

func TestCommandSequenceExecution(t *testing.T) {
	session := InitSession()
	t.Run("Check table creation sequence", func(t *testing.T) {
		// Create table and check if table was appended to session tables map
		// Also check if columns were appended
		createTableCommand := "CREATE TABLE test (id INT, name TEXT) engine linear;"
		expectedNumOfColumns := 2
		_, err := session.ExecuteCommand(createTableCommand)
		if err != nil {
			t.Error(err)
		}
		// Check if session.tables has "test" table
		if _, ok := session.tables["test"]; !ok {
			t.Error("No required table \"test\" was found in tables map")
		}
		columnNames := session.tables["test"].(linear.LinearTable).GetColumnsNames()
		if len(columnNames) != expectedNumOfColumns {
			t.Errorf("Count of columns differ. Expected: %d, got: %v", expectedNumOfColumns, columnNames)
		}
		t.Logf("Created table has columns: %s", session.tables["test"].GetColumns())
	})
	t.Run("Check table insertion", func(t *testing.T) {
		_, err := session.ExecuteCommand("INSERT INTO test VALUES (1, test);")
		if err != nil {
			t.Error(err)
		}
		if session.tables["test"].Count() != 1 {
			t.Errorf("expected only one row after insertion, got %d", session.tables["test"].Count())
		}
		_, err = session.ExecuteCommand("INSERT INTO test (id, name) VALUES (1, test_value);")
		if err != nil {
			t.Error(err)
		}
		if session.tables["test"].Count() != 2 {
			t.Errorf("expected only two rows after insertion, got %d", session.tables["test"].Count())
		}
	})
	t.Run("Check table selection", func(t *testing.T) {
		result, err := session.ExecuteCommand("SELECT (id, name) FROM test;")
		if err != nil {
			t.Error(err)
		}
		t.Logf("Result of select: \n%s", result)
	})
	t.Run("Check tables list", func(t *testing.T) {
		result, err := session.ExecuteCommand("SELECT (table_name, engine_type) FROM system.tables;")
		if err != nil {
			t.Error(err)
		}
		t.Logf("Result of select: \n%s", result)
	})
	t.Logf("State after all tests: %s", session.ToString())
}

func BenchmarkMassiveTableCreation(b *testing.B) {
	session := InitSession()
	inputs := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		// Generating random CREATE TABLE query
		tableName := token.GenerateRandomToken(token.IdentifierKind)
		//Generate column definition
		columns := make([]parser.ColumnDefinition, 10)
		for i := range columns {
			columns[i] = parser.ColumnDefinition{
				Name:     token.GenerateRandomToken(token.IdentifierKind),
				Datatype: token.GenerateRandomToken(token.TypeKind),
			}
		}
		inputQuery := &parser.CreateTableQuery{
			Name: tableName,
			Cols: &columns,
		}
		inputs[i] = (*inputQuery).CreateOriginal()
	}
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		_, err := session.ExecuteCommand(inputs[i])
		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
	}
}

// Will create table with 100 columns and b.N rows only with IdentifierKind tokens
func BenchmarkSimpleInsertion(b *testing.B) {
	rand.Seed(time.Now().Unix())
	session := InitSession()
	columnCount := 1
	// Generate table
	tableName := token.GenerateRandomToken(token.IdentifierKind)
	tableColumns := make([]parser.ColumnDefinition, columnCount)
	for i := range tableColumns {
		tableColumns[i] = parser.ColumnDefinition{
			Name: token.GenerateRandomToken(token.IdentifierKind),
			Datatype: &token.Token{
				Value: "text",
				Kind:  token.TypeKind,
			},
		}
	}
	inputQuery := &parser.CreateTableQuery{
		Name: tableName,
		Cols: &tableColumns,
	}
	session.ExecuteCommand((*inputQuery).CreateOriginal())
	// Generate payload
	for i := 0; i < b.N; i++ {
		columns := make([]*token.Token, columnCount)
		values := make([]*token.Token, columnCount)
		for i := range columns {
			columns[i] = tableColumns[i].Name
			values[i] = token.GenerateRandomToken(token.IdentifierKind)
		}
		tempQuery := &parser.InsertIntoQuery{
			Table:       tableName,
			ColumnNames: columns,
			Values:      values,
		}
		b.StartTimer()
		_, err := session.ExecuteCommand((*tempQuery).CreateOriginal())
		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
	}
}
