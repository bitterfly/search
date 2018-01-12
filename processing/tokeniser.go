package processing

type Tokeniser interface {
	Tokenise(text string) []string
	Normalise(token string) string
	IsStopWord(word string) bool
	GetTerms(text string, operation func(string))
}
