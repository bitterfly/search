package documents

import "testing"

func TestTokeniseSentence(t *testing.T) {
	tokens := TokeniseSentence([]byte("I am a cat."))

	if len(tokens) != 4 {
		t.Errorf("Tokeniser found %d tokens in the sentence 'I am a cat'\n", len(tokens))
	}
}
