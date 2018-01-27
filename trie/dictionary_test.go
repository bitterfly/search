package trie

import (
	"fmt"
	"testing"
)

func TestDictionary(t *testing.T) {
	words := []string{
		"foobar",
		"provid",
		"pro",
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
