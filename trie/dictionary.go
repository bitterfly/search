package trie

type Dictionary struct {
	trie       Trie
	lastTermID uint64
}

func NewDictionary() *Dictionary {
	return &Dictionary{
		trie:       *New(),
		lastTermID: 0,
	}
}

func (d *Dictionary) Get(word []byte) uint64 {
	id := d.trie.GetOrPut(word, d.lastTermID+1)
	if id > d.lastTermID {
		d.lastTermID = id
	}

	return id
}
