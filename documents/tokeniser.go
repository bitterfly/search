package documents

import "github.com/jdkato/prose/tokenize"

func TokeniseSentence(sentence []byte) [][]byte {
	tokeniser := tokenize.NewTreebankWordTokenizer()

	tokens := make([][]byte, 1)

	for _, word := range tokeniser.Tokenize(sentence) {
		tokens = append(tokens, []byte(word))
	}

	return tokens
}

func Tokenise(text []byte) [][]byte {
	sentenceSplitter := tokenize.NewPragmaticSegmenter("en")
	sentences := sentenceSplitter.Tokenize(text)

	tokens := make([][]byte, len(sentences))

	for _, sentence := range sentences {
		tokens = append(tokes, TokeniseSentence(setence))
	}

	return tokens
}
