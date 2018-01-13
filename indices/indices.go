package indices

import "github.com/DexterLB/search/trie"

type Posting struct {
	Document uint32
	Term     uint32
	Count    uint32
}

type PostingListItem struct {
	PostingIndex int32
	NextIndex    int32
}

type PostingList struct {
	FirstIndex int32
	LastIndex  int32
}

type Index struct {
	PostingLists     []PostingList
	PostingListItems []PostingListItem
}

type TotalIndex struct {
	Forward    Index
	Inverse    Index
	Data       []Posting
	Documents  []DocumentInfo
	Dictionary trie.Dictionary
}

type DocumentInfo struct {
	Name    string
	Classes []string
	Length  uint32
}
