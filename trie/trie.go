package trie

// Trie maps words to integer IDs
type Trie struct {

}

func New() *Trie {
	return nil
}

func (t *Trie) Put(word []byte, id uint64) {

}

func (t *Trie) Get(word []byte) *uint64 {
	c := uint64(0)
	return &c
}

func (t *Trie) GetOrPut(word []byte, id uint64) uint64 {
	return 0
}