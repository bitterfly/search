package trie

type Transition struct {
	id    int32
	label byte
}

// Trie maps words to integer IDs
// root always starts from zero
type Trie struct {
	maxIndex    int32
	transitions map[Transition]int32
	values      map[int32]uint64
}

func (t *Trie) traverseWith(word []byte) (int32, []byte) {
	node := int32(0)
	var destination int32
	var ok bool
	for i, letter := range word {
		destination, ok = t.transitions[Transition{id: node, label: letter}]
		if ok {
			node = destination
		} else {
			return node, word[i:]
		}
	}
	return destination, nil
}

func New() *Trie {
	return &Trie{
		maxIndex:    0,
		transitions: make(map[Transition]int32),
		values:      make(map[int32]uint64),
	}
}

func (t *Trie) Put(word []byte, value uint64) uint64 {
	node, rest := t.traverseWith(word)
	if rest != nil {
		for _, letter := range rest {
			t.maxIndex += 1
			t.transitions[Transition{id: node, label: letter}] = t.maxIndex
			node = t.maxIndex
		}
	}

	t.values[node] = value
	return value
}

func (t *Trie) Get(word []byte) *uint64 {
	destination, rest := t.traverseWith(word)
	if rest == nil {
		value, _ := t.values[destination]
		return &value
	}

	return nil
}

// Returns the value in the trie if the word is in
// else it returns the given value
func (t *Trie) GetOrPut(word []byte, value uint64) uint64 {
	inTrieValue := t.Get(word)
	if inTrieValue != nil {
		return *inTrieValue
	}

	return t.Put(word, value)
}

func (t *Trie) PutLambda(word []byte, lambda func(uint64) uint64, defaultValue uint64) {
	inTrieValue := t.Get(word)
	if inTrieValue != nil {
		t.Put(word, lambda(*inTrieValue))
	} else {
		t.Put(word, defaultValue)
	}
}

func (t *Trie) Walk(operation func([]byte, uint64)) {

}
