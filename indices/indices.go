package indices

import "github.com/DexterLB/search/trie"

type Posting struct {
	Index int32
	Count int32

	NextPostingIndex int32
}

type PostingList struct {
	FirstIndex int32
	LastIndex  int32
}

type Index struct {
	PostingLists []PostingList
	Postings     []Posting
}

type TotalIndex struct {
	Forward    Index
	Inverse    Index
	Documents  []DocumentInfo
	Dictionary trie.Dictionary
}

type DocumentInfo struct {
	Name    string
	Classes []string
	Length  int32
}

func NewTotalIndex() *TotalIndex {
	return &TotalIndex{
		Dictionary: *trie.NewDictionary(),
	}
}
