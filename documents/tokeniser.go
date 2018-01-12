package documents

import (
	"strings"
	"unicode"

	"github.com/jdkato/prose/tokenize"
)

func notPunctuation(word string) bool {
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

func tokeniseSentence(sentence string) []string {
	tokeniser := tokenize.NewTreebankWordTokenizer()

	tokens := make([]string, 0)
	for _, word := range tokeniser.Tokenize(sentence) {
		if notPunctuation(word) {
			tokens = append(tokens, word)
		}
	}

	return tokens
}

func Tokenise(text string) []string {
	sentenceSplitter, _ := tokenize.NewPragmaticSegmenter("en")
	sentences := sentenceSplitter.Tokenize(text)

	tokens := make([]string, 0)

	for _, sentence := range sentences {
		tokens = append(tokens, tokeniseSentence(sentence)...)
	}

	return tokens
}
