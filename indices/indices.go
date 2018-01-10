package indices

type PostingData struct {
	Document uint64
	Term uint64
	Times uint64
}

type Posting struct {
	Data int64
	Next int64
}

type PostingList struct {
	First int64
	Last int64
}

type Index struct {
	PostingLists []PostingList
	Postings 	[]Posting
}

type TotalIndex struct {
	Forward Index
	Inverse Index
	Data []PostingData
}