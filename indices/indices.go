package indices

type PostingData struct {
	Document uint32
	Term uint32
	Times uint32
}

type Posting struct {
	Data int32
	Next int32
}

type PostingList struct {
	First int32
	Last int32
}

type Index struct {
	PostingLists []PostingList
	Postings 	[]Posting
}

type TotalIndex struct {
	Forward Index
	Inverse Index
	Data []PostingData
	Documents []Document
}

type Document struct {
	Name string
	Class string
	Length uint32
}