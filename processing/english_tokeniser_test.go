package processing

import (
	"strings"
	"testing"
)

func TestTokenise(t *testing.T) {
	e, err := NewEnglishTokeniser(strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	sentence := "In a hole in the ground there lived a hobbit. Not a nasty, dirty, wet hole, filled with the ends of worms and an oozy smell, nor yet a dry, bare, sandy hole with nothing in it to sit down on or to eat: it was a hobbit-hole, and that means comfort."

	correctTokens := []string{"In", "a", "hole", "in", "the", "ground", "there", "lived", "a", "hobbit", "Not", "a", "nasty", "dirty", "wet", "hole", "filled", "with", "the", "ends", "of", "worms", "and", "an", "oozy", "smell", "nor", "yet", "a", "dry", "bare", "sandy", "hole", "with", "nothing", "in", "it", "to", "sit", "down", "on", "or", "to", "eat", "it", "was", "a", "hobbit-hole", "and", "that", "means", "comfort"}

	tokens := e.Tokenise(sentence)

	for i, token := range tokens {
		if token != correctTokens[i] {
			t.Errorf("Token should be %s but is %s\n", token, correctTokens[i])
		}
	}

	sentence = "Could've. Would've! Should've?"
	tokens = e.Tokenise(sentence)

	correctTokens = []string{"Could", "Would", "Should"}

	for i, token := range tokens {
		if token != correctTokens[i] {
			t.Errorf("Token should be %s but is %s\n", token, correctTokens[i])
		}
	}
}

func TestNormalise(t *testing.T) {
	e, err := NewEnglishTokeniser(strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	tokens := e.Tokenise("I'm many cats.")
	normalisedTokens := e.NormaliseMany(tokens)

	correctNormalisedTokens := []string{"i", "mani", "cat"}

	for i, token := range normalisedTokens {
		if token != correctNormalisedTokens[i] {
			t.Errorf("Token should be %s but is %s\n", correctNormalisedTokens[i], token)
		}
	}

	tokens = e.Tokenise("My dogs are actually one dog!")
	normalisedTokens = e.NormaliseMany(tokens)

	correctNormalisedTokens = []string{"my", "dog", "are", "actual", "one", "dog"}

	for i, token := range normalisedTokens {
		if token != correctNormalisedTokens[i] {
			t.Errorf("Token should be /%s/ but is /%s/\n", correctNormalisedTokens[i], token)
		}
	}

	tokens = e.Tokenise("If I'm dying you all die with me")
	normalisedTokens = e.NormaliseMany(tokens)

	correctNormalisedTokens = []string{"if", "i", "die", "you", "all", "die", "with", "me"}

	for i, token := range normalisedTokens {
		if token != correctNormalisedTokens[i] {
			t.Errorf("Token should be /%s/ but is /%s/\n", correctNormalisedTokens[i], token)
		}
	}

}
