package trie

type Dictionary struct {
	Trie       Trie
	LastTermID int32
}

type BiDictionary struct {
	Dictionary

	Inverse map[int32][]byte
}

func NewDictionary() *Dictionary {
	return &Dictionary{
		Trie:       *New(),
		LastTermID: 0,
	}
}

func (d *Dictionary) Get(word []byte) int32 {
	id := d.Trie.GetOrPut(word, d.LastTermID)
	if id == d.LastTermID {
		d.LastTermID += 1
	}

	return id
}

func NewBiDictionary() *BiDictionary {
	return &BiDictionary{
		Dictionary: *NewDictionary(),
		Inverse:    make(map[int32][]byte),
	}
}

func (b *BiDictionary) Get(word []byte) int32 {
	id := b.Dictionary.Get(word)
	if _, ok := b.Inverse[id]; !ok {
		b.Inverse[id] = append([]byte(nil), word...) // make a copy of word
	}
	return id
}

func (b *BiDictionary) GetInverse(id int32) []byte {
	word, ok := b.Inverse[id]
	if ok {
		return word
	} else {
		return nil
	}
}
