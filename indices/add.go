package indices

import "github.com/DexterLB/search/trie"

type Document struct {
	DocumentInfo

	TermCounts trie.Trie
}

func NewDocument() *Document {
	return &Document{
		DocumentInfo: DocumentInfo{},
		TermCounts:   *trie.New(),
	}
}

func (t *TotalIndex) Add(d *Document) {

}
