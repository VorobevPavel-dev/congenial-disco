package tokenizer

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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
		type generator func() *Token
		g := []generator{
			// Generates random NumericKind token
			func() *Token {
				min, max := -100000, 100000
				return &Token{
					Value: strconv.Itoa(rand.Intn(max-min+1) + min),
					Kind:  NumericKind,
				}
			},
			// Generates random IdentifierKind token
			func() *Token {
				b := make([]byte, 20)
				for i := range b {
					b[i] = charset[rand.Intn(len(charset))]
				}
				return &Token{
					Value: string(b),
					Kind:  IdentifierKind,
				}
			},
			// Generates random KeywordKind token
			func() *Token {
				index := rand.Intn(len(*Keywords()))
				return &Token{
					Value: (*Keywords())[index],
					Kind:  KeywordKind,
				}
			},
			// Generates random SymbolKind tokens
			func() *Token {
				// Excludes spaces
				value := " "
				for value == " " {
					index := rand.Intn(len(*Symbols()))
					value = (*Symbols())[index]
				}
				return &Token{
					Value: value,
					Kind:  SymbolKind,
				}
			},
			// Generates random TypeKind tokens
			func() *Token {
				index := rand.Intn(len(*Types()))
				return &Token{
					Value: (*Types())[index],
					Kind:  TypeKind,
				}
			},
		}

		numberOfTests := 10000
		sequenceLength := 20

		for j := 0; j < numberOfTests; j++ {
			generatedSequence := make([]string, sequenceLength)
			expectedOutputSequence := make([]*Token, sequenceLength)
			for i := range generatedSequence {
				genIndex := rand.Intn(len(g))
				tt := g[genIndex]()
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
