package processing

import (
	"strings"
	"unicode"

	"github.com/jdkato/prose/tokenize"
	"github.com/kljensen/snowball"
)

type EnglishTokeniser struct {
}

func NewEnglishTokeniser() *EnglishTokeniser {
	return &EnglishTokeniser{}
}

func (e *EnglishTokeniser) notPunctuation(word string) bool {
	if len(word) == 0 {
		return false
	}

	word = strings.TrimLeftFunc(word, func(r rune) bool { return r == '\'' })

	for _, symbol := range word {
		if symbol == '-' {
			continue
		}

		if !unicode.IsLetter(symbol) {
			return false
		}
	}
	return true
}

func (e *EnglishTokeniser) tokeniseSentence(sentence string) []string {
	tokeniser := tokenize.NewTreebankWordTokenizer()

	tokens := make([]string, 0)
	for _, word := range tokeniser.Tokenize(sentence) {
		if e.notPunctuation(word) {
			tokens = append(tokens, word)
		}
	}

	return tokens
}

func (e *EnglishTokeniser) Tokenise(text string) []string {
	sentenceSplitter, _ := tokenize.NewPragmaticSegmenter("en")
	sentences := sentenceSplitter.Tokenize(text)

	tokens := make([]string, 0)

	for _, sentence := range sentences {
		tokens = append(tokens, e.tokeniseSentence(sentence)...)
	}

	return tokens
}

func (e *EnglishTokeniser) Normalise(token string) string {
	lower := strings.ToLower(token)
	stemmed, err := snowball.Stem(lower, "english", true)
	if err == nil {
		return stemmed
	} else {
		return lower
	}
}

func (e *EnglishTokeniser) NormaliseMany(tokens []string) []string {
	normalised := make([]string, len(tokens))
	for i := range normalised {
		normalised[i] = e.Normalise(tokens[i])
	}

	return normalised
}
