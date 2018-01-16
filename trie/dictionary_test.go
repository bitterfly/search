package trie

import (
	"fmt"
	"testing"
)

func TestDictionary(t *testing.T) {
	words := []string{
		"midcon",
		"corp",
		"chief",
		"chairman",
		"subsidiari",
		"said",
		"occident",
		"oxi",
		"oper",
		"offic",
		"petroleum",
		"presid",
		"william",
		"terpstra",
		"resign",
		"repons",
		"reason",
		"reuter",
		"assum",
		"davi",
		"given",
		"bank",
		"band",
		"bring",
		"bought",
		"england",
		"said",
		"stg",
		"session",
		"shortag",
		"provid",
		"purchas",
		"pct",
		"money",
		"market",
		"mln",
		"assist",
		"afternoon",
		"addit",
		"total",
		"today",
		"treasuri",
		"help",
		"far",
		"forecast",
		"compar",
		"compris",
		"central",
		"revis",
		"reuter",
		"outright",
		"money-fx",
		"interest",
		"shr",
		"sharehold",
		"loss",
		"cts",
		"compani",
		"vs",
		"pro",
		"profit",
		"purchas",
		"public",
		"net",
		"note",
		"rev",
	}

	dic := NewDictionary()

	seen := make(map[int32]string)
	for i := range words {
		id := dic.Get([]byte(words[i]))
		if oldTerm, ok := seen[id]; ok && oldTerm != words[i] {
			panic(fmt.Sprintf("oh no, id %d for terms %s and %s", id, oldTerm, words[i]))
		}
		seen[id] = words[i]
	}
}
