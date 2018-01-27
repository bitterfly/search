package trie

type Transition struct {
	Id    int32
	Label byte
}

// Trie maps words to integer IDs
// root always starts from zero
type Trie struct {
	MaxIndex    int32
	Transitions map[Transition]int32
	Children    map[int32][]Transition
	Values      map[int32]int32
}

func (t *Trie) Empty() bool {
	return t.MaxIndex == 0
}

func New() *Trie {
	return &Trie{
		MaxIndex:    0,
		Transitions: make(map[Transition]int32),
		Children:    make(map[int32][]Transition),
		Values:      make(map[int32]int32),
	}
}

func (t *Trie) traverseWith(word []byte) (int32, []byte) {
	node := int32(0)
	var destination int32
	var ok bool
	for i, letter := range word {
		destination, ok = t.Transitions[Transition{Id: node, Label: letter}]
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
			t.MaxIndex += 1
			t.Transitions[Transition{Id: node, Label: letter}] = t.MaxIndex
			t.Children[node] = append(t.Children[node], Transition{Id: t.MaxIndex, Label: letter})
			node = t.MaxIndex
		}
	}

	t.Values[node] = value
	return value
}

func (t *Trie) Get(word []byte) *int32 {
	destination, rest := t.traverseWith(word)
	if rest == nil {
		value, ok := t.Values[destination]
		if ok {
			return &value
		} else {
			return nil
		}
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
	value, ok := t.Values[node]

	//If final apply operation
	if ok {
		operation(*word, value)
	}

	for _, transition := range t.Children[node] {
		*word = append(*word, transition.Label)
		t.walk(transition.Id, word, operation)
		*word = (*word)[:len(*word)-1]
	}

}

func (t *Trie) Walk(operation func([]byte, int32)) {
	var word []byte
	t.walk(0, &word, operation)
}
