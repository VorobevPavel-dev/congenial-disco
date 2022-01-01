package tokenizer

import "testing"

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
