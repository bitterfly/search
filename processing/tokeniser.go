package processing

type Tokeniser interface {
	Tokenise(text string) []string
	Normalise(token string) string
}
