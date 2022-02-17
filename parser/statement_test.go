package parser

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	token "github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

//getRandomFromKinds returns you random element of given kinds
func getRandomFromKinds(kinds ...token.TokenKind) token.TokenKind {
	return kinds[rand.Intn(len(kinds))]
}

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
			tokens := token.ParseTokenSequence(inputSQL)
			if tokens == nil || len(*tokens) == 0 {
				t.Errorf("no tokens extracted from query %s", inputSQL)
			}
			result, err := parseSelectStatement(*tokens)
			if err != nil {
				t.Errorf("an error occured on request %s, tokens: %v, err: %v",
					inputSQL,
					tokens,
					err,
				)
			}
			if !inputQuery.Equals(result) {
				t.Errorf("requests are different on request %s, excpected: %v, got: %v",
					inputSQL,
					inputQuery,
					result,
				)
			}
		}
	})
	t.Run("Test invalid select parsing", func(t *testing.T) {
		inputs := []string{
			"Select 1,b,c from test;",
			"INsert into test values (1,2,3);",
			"Select from test;",
			"Select from test",
			"Select a,     b, c from",
		}
		for testCase := range inputs {
			tokenList := *token.ParseTokenSequence(inputs[testCase])
			actualResult, err := parseSelectStatement(tokenList)
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
			// Generate INSERT INTO query
			tableName := token.GenerateRandomToken(token.IdentifierKind)
			// Generate column names and values
			columnsCount := rand.Intn(100) + 1
			columns := make([]*token.Token, columnsCount)
			values := make([]*token.Token, columnsCount)
			for i := range columns {
				columns[i] = token.GenerateRandomToken(token.IdentifierKind)
				values[i] = token.GenerateRandomToken(getRandomFromKinds(
					token.NumericKind,
					token.IdentifierKind,
				))
			}
			inputQuery := &InsertIntoQuery{
				Table:       tableName,
				ColumnNames: columns,
				Values:      values,
			}
			inputSQL := (*inputQuery).CreateOriginal()
			// t.Logf("Generated SQL: %s", inputSQL)
			// Start reverse parsing
			tokens := token.ParseTokenSequence(inputSQL)
			if tokens == nil || len(*tokens) == 0 {
				t.Errorf("no tokens extracted from query %s", inputSQL)
			}
			result, err := parseInsertIntoStatement(*tokens)
			if err != nil {
				t.Errorf("an error occured on request %s, tokens: %v, err: %v",
					inputSQL,
					tokens,
					err,
				)
			}
			if !inputQuery.Equals(result) {
				t.Errorf("requests are different on request %s, excpected: %v, got: %v",
					inputSQL,
					inputQuery,
					result,
				)
			}
		}
	})
	t.Run("Test invalid insert parsing", func(t *testing.T) {
		nextPermutation := func(input []*token.Token) []*token.Token {
			return append(input[1:], input[0])
		}
		inputs := []string{}

		template := "INSERT INTO %s %s VALUES %s;"

		for _, tt := range []token.TokenKind{
			token.NumericKind,
			token.KeywordKind,
			token.SymbolKind,
			token.TypeKind,
		} {
			// Incorrect type of <table_name>
			incTableName := token.GenerateRandomToken(tt)
			inputs = append(inputs, fmt.Sprintf(template, incTableName.Value, "(test)", "1"))

			// Incorrect type of <column_name> in different postions
			incColumnName := token.GenerateRandomToken(tt)
			numberOfColumns := rand.Intn(10) + 1
			columns := make([]*token.Token, numberOfColumns)
			values := make([]*token.Token, numberOfColumns)
			//	Fill slices
			for i := range columns {
				columns[i] = token.GenerateRandomToken(token.IdentifierKind)
				values[i] = token.GenerateRandomToken(getRandomFromKinds(token.NumericKind, token.IdentifierKind))
			}
			columns[0] = incColumnName
			for i := 0; i < len(columns); i++ {
				inputs = append(inputs, (&InsertIntoQuery{
					Table:       token.GenerateRandomToken(token.IdentifierKind),
					ColumnNames: columns,
					Values:      values,
				}).CreateOriginal())
				columns = nextPermutation(columns)
			}
		}
		// Incorrect value type of <value> in different positions
		for _, tt := range []token.TokenKind{
			token.KeywordKind,
			token.SymbolKind,
			token.TypeKind,
		} {
			numberOfColumns := rand.Intn(10) + 1
			columns := make([]*token.Token, numberOfColumns)
			values := make([]*token.Token, numberOfColumns)
			//	Fill slices
			for i := range columns {
				columns[i] = token.GenerateRandomToken(token.IdentifierKind)
				values[i] = token.GenerateRandomToken(getRandomFromKinds(token.NumericKind, token.IdentifierKind))
			}
			values[0] = token.GenerateRandomToken(tt)
			for i := 0; i < len(columns); i++ {
				inputs = append(inputs, (&InsertIntoQuery{
					Table:       token.GenerateRandomToken(token.IdentifierKind),
					ColumnNames: columns,
					Values:      values,
				}).CreateOriginal())
				columns = nextPermutation(columns)
			}
		}

		// TODO:
		// 		Incorrect ";" position

		// Other cases
		inputs = append(inputs,
			"table_name",
			"ins into table_name values (1);",
			"insert in table_name values (1);",
			"insert into table_name val (1);",
			"insert table_name values (1);",
			"into table_name values (1);",
			"table_name values (1);",
		)

		t.Logf("Generated incorrect inputs: %d", len(inputs))
		for _, c := range inputs {
			tokenList := *token.ParseTokenSequence(c)
			actualResult, err := parseSelectStatement(tokenList)
			if err == nil {
				t.Errorf("Expected error on query %s. Values got: %v",
					c, actualResult)
			}
		}
	})
	// t.Run("Test invalid select parsing", func(t *testing.T) {
	// 	inputs := []string{
	// 		"Select 1,b,c from test;",
	// 		"INsert into test values (1,2,3);",
	// 		"Select from test;",
	// 		"Select from test",
	// 		"Select a,     b, c from",
	// 	}
	// 	for testCase := range inputs {
	// 		tokenList := *token.ParseTokenSequence(inputs[testCase])
	// 		actualResult, err := parseSelectStatement(tokenList)
	// 		if err == nil {
	// 			t.Errorf("Expected error on set #%d. Values got: %v",
	// 				testCase, actualResult)
	// 		}
	// 	}
	// })
}

func TestCreateTableParsing(t *testing.T) {
	t.Run("Test valid create table parsing", func(t *testing.T) {
		numberOfTests := 100
		for i := 0; i < numberOfTests; i++ {
			// Generating random CREATE TABLE query
			tableName := token.GenerateRandomToken(token.IdentifierKind)
			//Generate column definition
			columns := make([]*ColumnDefinition, rand.Intn(100)+1)
			for i := range columns {
				columns[i] = &ColumnDefinition{
					Name:     token.GenerateRandomToken(token.IdentifierKind),
					Datatype: token.GenerateRandomToken(token.TypeKind),
				}
			}
			inputQuery := &CreateTableQuery{
				Name: tableName,
				Cols: columns,
			}
			inputSQL := (*inputQuery).CreateOriginal()
			// t.Logf("Generated SQL: %s", inputSQL)
			// Test SQL parsing
			tokens := token.ParseTokenSequence(inputSQL)
			if tokens == nil || len(*tokens) == 0 {
				t.Errorf("no tokens extracted from query %s", inputSQL)
			}
			// Try to parse token set as CREATE TABLE query
			result, err := parseCreateTableQuery(*tokens)
			if err != nil {
				t.Errorf("an error occured on request %s, tokens: %v, err: %v",
					inputSQL,
					tokens,
					err,
				)
			}
			// Assert input data equals data in CreateTableQuery
			if !inputQuery.Equals(result) {
				t.Errorf("requests are different on request %s, excpected: %v, got: %v",
					inputSQL,
					inputQuery,
					result,
				)
			}
		}
	})
	t.Run("Test invalid table creation statement parsing", func(t *testing.T) {
		inputs := []string{
			"create table test ;(id int, name text)",
			"create table test id int, name text;",
		}
		for testCase := range inputs {
			tokenList := *token.ParseTokenSequence(inputs[testCase])
			actualResult, err := parseCreateTableQuery(tokenList)
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
			result, err := parseShowCreateQuery(*tokens)
			if err != nil {
				t.Errorf("an error occured on request %s, tokens: %v, err: %v",
					inputSQL,
					tokens,
					err,
				)
			}
			// Assert input data equals data in ShowCreateQuery
			if !inputQuery.Equals(result) {
				t.Errorf("requests are different on request %s, excpected: %v, got: %v",
					inputSQL,
					inputQuery,
					result,
				)
			}
		}
	})
	t.Run("Test invalid show create statement parsing", func(t *testing.T) {
		// Generate inputs
		inputs := []string{}

		//		Not identifier kind of <table_name>
		// TODO: Dynamic search for token kinds
		for i := 0; i < 10; i++ {
			tokenKind := rand.Intn(5)
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

		for testCase := range inputs {
			tokenList := *token.ParseTokenSequence(inputs[testCase])
			actualResult, err := parseShowCreateQuery(tokenList)
			if err == nil {
				t.Errorf("Expected error on set #%d. Values got: %v",
					testCase, actualResult)
			}
		}
	})
}
