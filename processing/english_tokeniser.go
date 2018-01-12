package processing

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"
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

func NewEnglishTokeniser(stopWordList io.Reader) (*EnglishTokeniser, error) {
	tok := &EnglishTokeniser{
		stopWords: *trie.New(),
	}

	scanner := bufio.NewScanner(stopWordList)
	for scanner.Scan() {
		tok.stopWords.Put([]byte(scanner.Text()), 1)
	}

	return tok, scanner.Err()
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

var tokeniserLock sync.Mutex

func (e *EnglishTokeniser) Tokenise(text string) []string {
	sentenceSplitter, _ := tokenize.NewThreadSafePragmaticSegmenter("en")
	sentences := sentenceSplitter.Tokenize(text)

	tokens := make([]string, 0)

	tokeniserLock.Lock() // the PragmaticSegmenter is unsafe :(
	for _, sentence := range sentences {
		tokens = append(tokens, e.tokeniseSentence(sentence)...)
	}
	tokeniserLock.Unlock()

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

func (e *EnglishTokeniser) IsStopWord(word string) bool {
	return e.stopWords.Get([]byte(word)) != nil
}

func (e *EnglishTokeniser) GetTerms(text string, operation func(string)) {
	terms := e.Tokenise(text)
	for i := range terms {
		if e.IsStopWord(terms[i]) {
			continue
		}
		terms[i] = e.Normalise(terms[i])
		if e.IsStopWord(terms[i]) {
			continue // maybe don't need the second check?
		}
		operation(terms[i])
	}
}
