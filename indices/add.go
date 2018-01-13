package indices

import (
	"fmt"
	"strings"

	"github.com/DexterLB/search/trie"
)

type InfoAndTerms struct {
	DocumentInfo

	TermsAndCounts trie.Trie
}

func NewInfoAndTerms() *InfoAndTerms {
	return &InfoAndTerms{
		DocumentInfo:   DocumentInfo{},
		TermsAndCounts: *trie.New(),
	}
}

func (d *InfoAndTerms) Print() {
	fmt.Printf("************\n")
	fmt.Printf(
		"name: %s, classes: %s, length: %d\nterms:\n",
		d.Name,
		strings.Join(d.Classes, ", "),
		d.Length,
	)
	d.TermsAndCounts.Walk(func(term []byte, count uint64) {
		fmt.Printf("  %s: %d\n", string(term), count)
	})
}

func (t *TotalIndex) Add(d *InfoAndTerms) {

}
