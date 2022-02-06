package tokenizer

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestTokenParsing(t *testing.T) {
	// Generates numberOfTests numbers, converts them to strings and
	// tries to represent it as list of tokens
	t.Run("Parsing numeric tokens", func(t *testing.T) {
		// Generate input and expected output data
		numberOfTests := 10000
		rand.Seed(time.Now().UnixNano())
		min, max := -100000, 100000
		for i := 0; i < numberOfTests; i++ {
			ti := strconv.Itoa(rand.Intn(max-min+1) + min)
			eo := &Token{Value: ti, Kind: NumericKind}
			ar, err := TokenFromString(ti, 0)
			if err != nil {
				t.Errorf("Cannot convert number to token with NumericKind. Input: %s, expected: %v, got: %v, err: %v",
					ti, eo, ar, err)
			}
		}
	})

	// Check if all keyword-strings defined as constants
	// can be parsed as tokens with KeywordKind
	t.Run("Parsing keyword tokens", func(t *testing.T) {
		inputs := []string{"SELECT", "FROM", "AS", "TABLE", "CREATE", "INSERT", "INTO", "VALUES", "SHOW"}
		for _, ti := range inputs {
			eo := &Token{Value: ti, Kind: KeywordKind}
			ar, err := TokenFromString(ti, 0)
			if err != nil {
				t.Errorf("Cannot convert number to token with NumericKind. Input: %s, expected: %v, got: %v, err: %v",
					ti, eo, ar, err)
			}
		}
	})
}

func TestTokenSequenceParsing(t *testing.T) {
	// Checks if sequense of tokens can be parsed into tokens.
	// Uses generators for creating sequences
	t.Run("Parse token sequence", func(t *testing.T) {
		rand.Seed(time.Now().UnixNano())
		numberOfTests := 10000
		sequenceLength := 20

		for j := 0; j < numberOfTests; j++ {
			kind := rand.Intn(5)
			generatedSequence := make([]string, sequenceLength)
			expectedOutputSequence := make([]*Token, sequenceLength)
			for i := range generatedSequence {
				tt := GenerateRandomToken(TokenKind(kind))
				expectedOutputSequence[i] = tt
				generatedSequence[i] = tt.Value
			}
			request := strings.Join(generatedSequence, " ")
			// t.Logf("Generated request: %s", request)
			actalResult := *ParseTokenSequence(request)
			for i := range actalResult {
				if !expectedOutputSequence[i].Equals(actalResult[i]) {
					t.Errorf("Tokens are different on position %d. Sequence: %s, expected: %v, got: %v",
						i, request, expectedOutputSequence, actalResult)
				}
			}
		}
	})
}
