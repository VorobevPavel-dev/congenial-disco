package tokenizer

import (
	"testing"
)

func TestTokenParsing(t *testing.T) {
	t.Run("Parsing numeric tokens", func(t *testing.T) {
		input := []string{"123", "456", "-1234"}
		output := []Token{
			{
				Value:    "123",
				Kind:     NumericKind,
				Position: 0,
			},
			{
				Value:    "456",
				Kind:     NumericKind,
				Position: 0,
			},
			{
				Value:    "-1234",
				Kind:     NumericKind,
				Position: 0,
			},
		}
		for index, testCase := range input {
			actualValue, err := TokenFromString(testCase, 0)
			if err != nil {
				t.Errorf("Error in test case #%d: expected: %v, got: %v",
					index, output[index], actualValue)
			}
		}
	})
	//t.Run("Parse invalid numeric tokens", func(t *testing.T) {
	//
	//})
	t.Run("Parsing keyword tokens", func(t *testing.T) {
		input := []string{"select", "Into", "create"}
		output := []Token{
			{
				Value:    "select",
				Kind:     KeywordKind,
				Position: 0,
			},
			{
				Value:    "into",
				Kind:     KeywordKind,
				Position: 0,
			},
			{
				Value:    "create",
				Kind:     KeywordKind,
				Position: 0,
			},
		}
		for index, testCase := range input {
			actualValue, err := TokenFromString(testCase, 0)
			if err != nil {
				t.Errorf("Error in test case #%d: expected: %v, got: %v",
					index, output[index], actualValue)
			}
		}
	})
	t.Run("Parse invalid keyword tokens", func(t *testing.T) {
		input := []string{"selec1t", "Into2", "createz"}
		for _, testCase := range input {
			actualValue := ParseKeywordToken(testCase)
			if actualValue != nil {
				t.Errorf("Expected nil on parsing keyword: given: %v, got: %v",
					testCase, actualValue)
			}
		}
	})
}

func TestTokenSequenceParsing(t *testing.T) {
	t.Run("Parse token sequence", func(t *testing.T) {
		inputs := []string{
			"select from test(1234)",
		}
		expectedResults := [][]*Token{
			{
				{
					Value:    "select",
					Kind:     KeywordKind,
					Position: 0,
				},
				{
					Value:    " ",
					Kind:     SymbolKind,
					Position: 6,
				},
				{
					Value:    "from",
					Kind:     KeywordKind,
					Position: 7,
				},
				{
					Value:    " ",
					Kind:     SymbolKind,
					Position: 11,
				},
				{
					Value:    "test",
					Kind:     IdentifierKind,
					Position: 12,
				},
				{
					Value:    "(",
					Kind:     SymbolKind,
					Position: 16,
				},
				{
					Value:    "1234",
					Kind:     NumericKind,
					Position: 17,
				},
				{
					Value:    ")",
					Kind:     SymbolKind,
					Position: 21,
				},
			},
		}
		for testCase := range inputs {
			actualResult := *ParseTokenSequence(inputs[testCase])
			if len(actualResult) != len(expectedResults[testCase]) {
				t.Errorf("Function have returned unexpected number of tokens: %d (expected %d)",
					len(actualResult), len(expectedResults[testCase]))
			}
			for index := range actualResult {
				if !actualResult[index].equals(expectedResults[testCase][index]) {
					t.Errorf("Tokens on position %d are different. Expected: %s, got: %s",
						index+1,
						expectedResults[testCase][index],
						actualResult[index])
				}
			}
		}
	})
}
