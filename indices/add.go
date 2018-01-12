package indices

import (
	"fmt"
	"strings"

	"github.com/DexterLB/search/trie"
)

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

func (d *Document) Print() {
	fmt.Printf("************\n")
	fmt.Printf(
		"name: %s, classes: %s, length: %d\nterms:\n",
		d.Name,
		strings.Join(d.Classes, ", "),
		d.Length,
	)
	d.TermCounts.Walk(func(term []byte, count uint64) {
		fmt.Printf("  %s: %d\n", string(term), count)
	})
}

func (t *TotalIndex) Add(d *Document) {

}
