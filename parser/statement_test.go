package parser

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/VorobevPavel-dev/congenial-disco/tokenizer"
	token "github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

func TestSelectStatementParsing(t *testing.T) {
	t.Run("Test valid select parsing", func(t *testing.T) {
		inputs := []string{
			"Select a,b,c from test;",
			"select a1 from test;",
		}
		expectedOutputs := []*SelectStatement{
			{
				Item: []*token.Token{
					{Value: "a", Kind: token.IdentifierKind},
					{Value: "b", Kind: token.IdentifierKind},
					{Value: "c", Kind: token.IdentifierKind},
				},
				From: token.Token{
					Value: "test",
					Kind:  token.IdentifierKind,
				},
			},
			{
				Item: []*token.Token{
					{Value: "a1", Kind: token.IdentifierKind},
				},
				From: token.Token{
					Value: "test",
					Kind:  token.IdentifierKind,
				},
			},
		}
		for testCase := range inputs {
			tokenList := *token.ParseTokenSequence(inputs[testCase])
			actualResult, err := parseSelectStatement(tokenList)
			if err != nil {
				t.Errorf("Parsing failed on set #%d: %v",
					testCase, err)
			}
			if !actualResult.Equals(expectedOutputs[testCase]) {
				t.Errorf("Assertion failed. Expected: %s, got: %s",
					actualResult.String(), expectedOutputs[testCase].String())
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
	t.Run("Test valid select parsing", func(t *testing.T) {
		inputs := []string{
			"INsert into test values (1,2,3);",
		}
		expectedOutputs := []*InsertStatement{
			{
				Table: token.Token{
					Value: "test",
					Kind:  token.IdentifierKind,
				},
				Values: []*token.Token{
					{Value: "1", Kind: token.NumericKind},
					{Value: "2", Kind: token.NumericKind},
					{Value: "3", Kind: token.NumericKind},
				},
				ColumnNames: nil,
			},
		}
		for testCase := range inputs {
			tokenList := *token.ParseTokenSequence(inputs[testCase])
			actualResult, err := parseInsertIntoStatement(tokenList)
			if err != nil {
				t.Errorf("Parsing failed on set #%d: %v",
					testCase, err)
			}
			if !actualResult.Equals(expectedOutputs[testCase]) {
				t.Errorf("Assertion failed. Expected: %s, got: %s",
					actualResult.String(), expectedOutputs[testCase].String())
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
			tokens := tokenizer.ParseTokenSequence(inputSQL)
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
			tokens := tokenizer.ParseTokenSequence(inputSQL)
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
