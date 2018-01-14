package trie

type Dictionary struct {
	Trie       Trie
	LastTermID int32
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
