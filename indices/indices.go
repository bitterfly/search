package indices

import "github.com/DexterLB/search/trie"

type Posting struct {
	Index uint32
	Count uint32

	NextPostingIndex uint32
}

type PostingList struct {
	FirstIndex uint32
	LastIndex  uint32
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
	Length  uint32
}
