package trie

import "testing"

func TestTrie_Put_and_Get(t *testing.T) {
	trie := New()

	trie.Put([]byte("foo"), 42)
	trie.Put([]byte("fob"), 43)
	trie.Put([]byte("bar"), 44)

	if *trie.Get([]byte("foo")) != 42 {
		t.Errorf("cannot get foo")
	}

	if *trie.Get([]byte("fob")) != 43 {
		t.Errorf("cannot get fob")
	}

	if *trie.Get([]byte("bar")) != 44 {
		t.Errorf("cannot get bar")
	}

	if trie.GetOrPut([]byte("qux"), 5) != 5 {
		t.Errorf("cannot put or get qux when it's not present")
	}

	if *trie.Get([]byte("qux")) != 5 {
		t.Errorf("cannot get qux")
	}

	if trie.GetOrPut([]byte("qux"), 1) != 5 {
		t.Errorf("cannot put or get qux when it is present")
	}
}