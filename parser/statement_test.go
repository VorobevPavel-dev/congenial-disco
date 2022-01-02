package parser

import (
	"github.com/VorobevPavel-dev/congenial-disco/tokenizer"
	"testing"
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
