package parser

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	token "github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

func TestSelectStatementParsing(t *testing.T) {
	t.Run("Test valid select parsing", func(t *testing.T) {
		// Generate valid select parsing
		numberOfTests := 100
		for i := 0; i < numberOfTests; i++ {
			tempTableName := token.GenerateRandomToken(token.IdentifierKind)
			columnsCount := rand.Intn(100) + 1
			columns := make([]*token.Token, columnsCount)
			for i := range columns {
				columns[i] = token.GenerateRandomToken(token.IdentifierKind)
			}
			inputQuery := &SelectQuery{
				Columns: columns,
				From:    tempTableName,
			}
			inputSQL := inputQuery.CreateOriginal()
			// t.Logf("Generated SQL: %s", inputSQL)
			result, err := Parse(inputSQL)
			if err != nil {
				t.Errorf("an error occured on request %s, err: %v",
					inputSQL,
					err,
				)
			}
			if result == nil {
				t.Errorf("no result from parsing %s, error: %s", inputSQL, err.Error())
			} else {
				if !inputQuery.Equals(result.SelectStatement) {
					t.Errorf("requests are different on request %s, excpected: %v, got: %v",
						inputSQL,
						inputQuery,
						result,
					)
				}
			}
		}
	})
	t.Run("Test invalid select parsing", func(t *testing.T) {
		inputs := []string{
			"Select 1,b,c from test;",
			"Select from test;",
			"Select from test",
			"Select a,     b, c from",
		}
		for testCase := range inputs {
			actualResult, err := Parse(inputs[testCase])
			if err == nil {
				t.Errorf("Expected error on set #%d. Values got: %v",
					testCase, actualResult)
			}
		}
	})
}

func TestInsertStatementParsing(t *testing.T) {
	t.Run("Test valid insert parsing", func(t *testing.T) {
		numberOfTests := 100
		for i := 0; i < numberOfTests; i++ {
			inputStatement, inputSQL := GenerateStatement(InsertType)
			inputQuery := inputStatement.InsertStatement
			// t.Logf("Generated SQL: %s", inputSQL)
			// Start reverse parsing
			result, err := Parse(inputSQL)
			if err != nil {
				t.Errorf("an error occured on request %s, err: %v",
					inputSQL,
					err,
				)
			}
			if result == nil {
				t.Errorf("no result from parsing %s, error: %s", inputSQL, err.Error())
			} else {
				if !inputQuery.Equals(result.InsertStatement) {
					t.Errorf("requests are different on request %s, excpected: %v, got: %v",
						inputSQL,
						inputQuery,
						result,
					)
				}
			}
		}
	})
	t.Run("Test invalid insert parsing", func(t *testing.T) {

		// Other cases
		inputs := []string{
			"table_name",
			"ins into table_name values (1);",
			"insert in table_name values (1);",
			"insert into table_name val (1);",
			"insert table_name values (1);",
			"into table_name values (1);",
			"table_name values (1);",
		}

		t.Logf("Generated incorrect inputs: %d", len(inputs))
		for _, c := range inputs {
			actualResult, err := Parse(c)
			if err == nil {
				t.Errorf("Expected error on query %s. Values got: %v",
					c, actualResult)
			}
		}
	})
}

func TestCreateTableParsing(t *testing.T) {
	t.Run("Test valid create table parsing", func(t *testing.T) {
		numberOfTests := 100
		for i := 0; i < numberOfTests; i++ {
			// Generating random CREATE TABLE query
			tableName := token.GenerateRandomToken(token.IdentifierKind)
			//Generate column definition
			columns := make([]ColumnDefinition, rand.Intn(100)+1)
			for i := range columns {
				columns[i] = ColumnDefinition{
					Name:     token.GenerateRandomToken(token.IdentifierKind),
					Datatype: token.GenerateRandomToken(token.TypeKind),
				}
			}
			// Choose engine
			engine := token.GenerateRandomToken(token.EngineKind)
			inputQuery := &CreateTableQuery{
				Name:   tableName,
				Cols:   &columns,
				Engine: engine,
			}
			inputSQL := (*inputQuery).CreateOriginal()
			// Try to parse token set as CREATE TABLE query
			result, err := Parse(inputSQL)
			if err != nil {
				t.Errorf("an error occured on request %s, err: %v",
					inputSQL,
					err,
				)
			}
			if result == nil {
				t.Errorf("no result from parsing %s, error: %s", inputSQL, err.Error())
			} else {
				if !inputQuery.Equals(result.CreateTableStatement) {
					t.Errorf("requests are different on request %s, excpected: %v, got: %v",
						inputSQL,
						inputQuery,
						result,
					)
				}
			}
		}
	})
	t.Run("Test invalid table creation statement parsing", func(t *testing.T) {
		inputs := []string{
			"create table test ;(id int, name text)",
			"create table test id int, name text;",
		}
		for testCase := range inputs {
			actualResult, err := Parse(inputs[testCase])
			if err == nil {
				t.Errorf("Expected error on set #%d. Values got: %v",
					testCase, actualResult)
			}
		}
	})
}

func TestShowCreateParsing(t *testing.T) {
	t.Run("Test valid show create parsing", func(t *testing.T) {
		numberOfTests := 10
		for i := 0; i < numberOfTests; i++ {
			// Generate SHOW CREATE query
			tableName := token.GenerateRandomToken(token.IdentifierKind)
			inputQuery := &ShowCreateQuery{
				TableName: tableName,
			}
			inputSQL := inputQuery.CreateOriginal()
			// t.Logf("Generated SQL: %s", inputSQL)
			tokens := token.ParseTokenSequence(inputSQL)
			if tokens == nil || len(*tokens) == 0 {
				t.Errorf("no tokens extracted from query %s", inputSQL)
			}
			// Try to parse token set as SHOW CREATE query
			result, err := Parse(inputSQL)
			if err != nil {
				t.Errorf("an error occured on request %s, err: %v",
					inputSQL,
					err,
				)
			}
			if result == nil {
				t.Errorf("no result from parsing %s, error: %s", inputSQL, err.Error())
			} else {
				if !inputQuery.Equals(result.ShowCreateStatement) {
					t.Errorf("requests are different on request %s, excpected: %v, got: %v",
						inputSQL,
						inputQuery,
						result,
					)
				}
			}
		}
	})
	t.Run("Test invalid show create statement parsing", func(t *testing.T) {
		// Generate inputs
		inputs := []string{}

		//		Not identifier kind of <table_name>
		// TODO: Dynamic search for token kinds
		for i := 0; i < 10; i++ {
			tokenKind := rand.Intn(len(token.Reserved))
			tableName := token.GenerateRandomToken(token.TokenKind(tokenKind))
			if tableName.Kind == token.IdentifierKind {
				i--
			} else {
				inputs = append(inputs, fmt.Sprintf("show create (%s);", tableName.Value))
			}
		}

		//		Incorrect ";" placement
		tk := token.GenerateRandomToken(token.IdentifierKind)
		tokenList := token.ParseTokenSequence(fmt.Sprintf("show create (%s)", tk.Value))
		//		Generate list of strings
		stringParts := make([]string, len(*tokenList))
		for i := range stringParts {
			stringParts[i] = (*tokenList)[i].Value
		}
		for i := 1; i < len(*tokenList); i++ {
			tempSQL := strings.Join(stringParts[:i], " ") + ";" + strings.Join(stringParts[i:], " ")
			inputs = append(inputs, tempSQL)
			// t.Log(tempSQL)
		}

		//		Other cases
		inputs = append(inputs, []string{
			"(table_name);",
			"sh create (table_name);",
			"show cr (table_name);",
			"show create (table_name;",
			"show create (table_name",
			"show create table_name)",
			"show create table_name);",
			"show create (table_name)",
		}...)
		for _, c := range inputs {
			actualResult, err := Parse(c)
			if err == nil {
				t.Errorf("Expected error on query %s. Values got: %v",
					c, actualResult)
			} else {
				t.Logf("got error: %v", err)
			}
		}
	})
}
