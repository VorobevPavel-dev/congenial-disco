package parser

import (
	"testing"

	"github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

func TestSelectStatementParsing(t *testing.T) {
	t.Run("Test valid select parsing", func(t *testing.T) {
		inputs := []string{
			"Select a,b,c from test;",
			"select a1 from test;",
		}
		expectedOutputs := []*SelectStatement{
			{
				Item: []*tokenizer.Token{
					{Value: "a", Kind: tokenizer.IdentifierKind},
					{Value: "b", Kind: tokenizer.IdentifierKind},
					{Value: "c", Kind: tokenizer.IdentifierKind},
				},
				From: tokenizer.Token{
					Value: "test",
					Kind:  tokenizer.IdentifierKind,
				},
			},
			{
				Item: []*tokenizer.Token{
					{Value: "a1", Kind: tokenizer.IdentifierKind},
				},
				From: tokenizer.Token{
					Value: "test",
					Kind:  tokenizer.IdentifierKind,
				},
			},
		}
		for testCase := range inputs {
			tokenList := *tokenizer.ParseTokenSequence(inputs[testCase])
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
			tokenList := *tokenizer.ParseTokenSequence(inputs[testCase])
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
				Table: tokenizer.Token{
					Value: "test",
					Kind:  tokenizer.IdentifierKind,
				},
				Values: []*tokenizer.Token{
					{Value: "1", Kind: tokenizer.NumericKind},
					{Value: "2", Kind: tokenizer.NumericKind},
					{Value: "3", Kind: tokenizer.NumericKind},
				},
				ColumnNames: nil,
			},
		}
		for testCase := range inputs {
			tokenList := *tokenizer.ParseTokenSequence(inputs[testCase])
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
	// 		tokenList := *tokenizer.ParseTokenSequence(inputs[testCase])
	// 		actualResult, err := parseSelectStatement(tokenList)
	// 		if err == nil {
	// 			t.Errorf("Expected error on set #%d. Values got: %v",
	// 				testCase, actualResult)
	// 		}
	// 	}
	// })
}

func TestCreateTableParsing(t *testing.T) {
	t.Run("Test valid select parsing", func(t *testing.T) {
		inputs := []string{
			"create table test (id int, name text);",
		}
		expectedOutputs := []*CreateTableStatement{
			{
				Name: tokenizer.Token{Value: "test", Kind: tokenizer.IdentifierKind},
				Cols: []*ColumnDefinition{
					{
						Name:     tokenizer.Token{Value: "id", Kind: tokenizer.IdentifierKind},
						Datatype: tokenizer.Token{Value: "int", Kind: tokenizer.TypeKind},
					},
					{
						Name:     tokenizer.Token{Value: "name", Kind: tokenizer.IdentifierKind},
						Datatype: tokenizer.Token{Value: "text", Kind: tokenizer.TypeKind},
					},
				},
			},
		}
		for testCase := range inputs {
			tokenList := *tokenizer.ParseTokenSequence(inputs[testCase])
			actualResult, err := parseCreateTableStatement(tokenList)
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
	t.Run("Test invalid table creation statement parsing", func(t *testing.T) {
		inputs := []string{
			"create table test (id int, name text)",
			"create table test id int, name text;",
		}
		for testCase := range inputs {
			tokenList := *tokenizer.ParseTokenSequence(inputs[testCase])
			actualResult, err := parseCreateTableStatement(tokenList)
			if err == nil {
				t.Errorf("Expected error on set #%d. Values got: %v",
					testCase, actualResult)
			}
		}
	})
}
