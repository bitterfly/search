package indices

import (
	"fmt"
	"sort"
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

type TermAndCount struct {
	TermID uint32
	Count  uint32
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
	var sortedTermsAndCounts []TermAndCount

	d.TermsAndCounts.Walk(func(term []byte, count uint64) {
		sortedTermsAndCounts = append(
			sortedTermsAndCounts,
			TermAndCount{
				TermID: uint32(t.Dictionary.Get(term)),
				Count:  uint32(count),
			},
		)
	})

	sort.Slice(
		sortedTermsAndCounts,
		func(i, j int) bool {
			return sortedTermsAndCounts[i].TermID < sortedTermsAndCounts[j].TermID
		},
	)

}
