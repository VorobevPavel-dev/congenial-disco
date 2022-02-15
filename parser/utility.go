package parser

import (
	"fmt"

	t "github.com/VorobevPavel-dev/congenial-disco/tokenizer"
)

// AssertTokenSequence will compare all tokens from tokens slice with tokens from expected slice.
func AssertTokenSequence(tokens []*t.Token, expected []*t.Token) error {
	for index, token := range tokens {
		if !token.Equals(expected[index]) {
			return fmt.Errorf("expected: %s, got: %s", expected[index].Value, token.Value)
		}
	}
	return nil
}
