package indices

import "github.com/DexterLB/search/trie"

type Document struct {
	DocumentInfo

	TermCounts trie.Trie
}

func (t *TotalIndex) Add(d *Document) {

}
