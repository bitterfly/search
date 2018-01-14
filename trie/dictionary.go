package trie

type Dictionary struct {
	trie       Trie
	lastTermID int32
}

func NewDictionary() *Dictionary {
	return &Dictionary{
		trie:       *New(),
		lastTermID: 0,
	}
}

func (d *Dictionary) Get(word []byte) int32 {
	id := d.trie.GetOrPut(word, d.lastTermID)
	if id == d.lastTermID {
		d.lastTermID += 1
	}

	return id
}
