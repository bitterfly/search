package processing

import (
	"bufio"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/DexterLB/prose/tokenize"
	"github.com/DexterLB/search/trie"
	"github.com/kljensen/snowball"
)

type EnglishTokeniser struct {
	stopWords trie.Trie
}

func NewEnglishTokeniserFromFile(stopWordFile string) (*EnglishTokeniser, error) {
	f, err := os.Open(stopWordFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return NewEnglishTokeniser(f)
}

func (e *EnglishTokeniser) stem(word string) string {
	stemmed, err := snowball.Stem(word, "english", true)
	if err == nil {
		return stemmed
	} else {
		return lower
	}
}

func NewEnglishTokeniser(stopWordList io.Reader) (*EnglishTokeniser, error) {
	tok := &EnglishTokeniser{
		stopWords: *trie.New(),
	}

	scanner := bufio.NewScanner(stopWordList)
	for scanner.Scan() {
		tok.stopWords.Put([]byte(scanner.Text()), 1)
		tok.stopWords.Put(tok.stem(scanner.Text()), 1)
	}

	return tok, scanner.Err()
}

func (e *EnglishTokeniser) notPunctuation(word string) bool {
	if len(word) == 0 {
		return false
	}

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

func (e *EnglishTokeniser) Tokenise(text string) []string {
	sentenceSplitter, _ := tokenize.NewThreadSafePragmaticSegmenter("en")
	sentences := sentenceSplitter.Tokenize(text)

	tokens := make([]string, 0)
	tokeniser := tokenize.NewTreebankWordTokenizer()

	for _, sentence := range sentences {
		for _, word := range tokeniser.Tokenize(sentence) {
			if e.notPunctuation(word) {
				tokens = append(tokens, word)
			}
		}
	}

	return tokens
}

func (e *EnglishTokeniser) Normalise(token string) string {
	return e.stem(strings.ToLower(token))
}

func (e *EnglishTokeniser) NormaliseMany(tokens []string) []string {
	normalised := make([]string, len(tokens))
	for i := range normalised {
		normalised[i] = e.Normalise(tokens[i])
	}

	return normalised
}

func (e *EnglishTokeniser) IsStopWord(word string) bool {
	return e.stopWords.Get([]byte(word)) != nil
}

func (e *EnglishTokeniser) GetTerms(text string, operation func(string)) {
	terms := e.Tokenise(text)
	for i := range terms {
		terms[i] = e.Normalise(terms[i])
		if e.IsStopWord(terms[i]) {
			continue // maybe don't need the second check?
		}
		operation(terms[i])
	}
}
