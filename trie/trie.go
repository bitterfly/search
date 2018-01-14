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
	children    map[int32][]Transition
	values      map[int32]int32
}

func New() *Trie {
	return &Trie{
		maxIndex:    0,
		transitions: make(map[Transition]int32),
		children:    make(map[int32][]Transition),
		values:      make(map[int32]int32),
	}
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

func (t *Trie) Put(word []byte, value int32) int32 {
	node, rest := t.traverseWith(word)
	if rest != nil {
		for _, letter := range rest {
			t.maxIndex += 1
			t.transitions[Transition{id: node, label: letter}] = t.maxIndex
			t.children[node] = append(t.children[node], Transition{id: t.maxIndex, label: letter})
			node = t.maxIndex
		}
	}

	t.values[node] = value
	return value
}

func (t *Trie) Get(word []byte) *int32 {
	destination, rest := t.traverseWith(word)
	if rest == nil {
		value, _ := t.values[destination]
		return &value
	}

	return nil
}

// Returns the value in the trie if the word is in
// else it returns the given value
func (t *Trie) GetOrPut(word []byte, value int32) int32 {
	inTrieValue := t.Get(word)
	if inTrieValue != nil {
		return *inTrieValue
	}

	return t.Put(word, value)
}

func (t *Trie) PutLambda(word []byte, lambda func(int32) int32, defaultValue int32) {
	inTrieValue := t.Get(word)
	if inTrieValue != nil {
		t.Put(word, lambda(*inTrieValue))
	} else {
		t.Put(word, defaultValue)
	}
}

func (t *Trie) walk(node int32, word *[]byte, operation func([]byte, int32)) {
	value, ok := t.values[node]

	//If final apply operation
	if ok {
		operation(*word, value)
	}

	for _, transition := range t.children[node] {
		*word = append(*word, transition.label)
		t.walk(transition.id, word, operation)
		*word = (*word)[:len(*word)-1]
	}

}

func (t *Trie) Walk(operation func([]byte, int32)) {
	var word []byte
	t.walk(0, &word, operation)
}
