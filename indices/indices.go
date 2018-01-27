package indices

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"os"

	"github.com/bitterfly/search/trie"
)

type Posting struct {
	Index           int32
	Count           int32
	NormalisedCount float32

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
	Dictionary trie.BiDictionary // bidictionary is better for debugging
	ClassNames trie.BiDictionary
}

type DocumentInfo struct {
	Name         string
	Classes      []int32
	Length       int32
	UniqueLength int32
	ClusterID    int
}

func NewTotalIndex() *TotalIndex {
	return &TotalIndex{
		Dictionary: *trie.NewBiDictionary(),
		ClassNames: *trie.NewBiDictionary(),
	}
}

func (t *TotalIndex) LoopOverTermPostings(termID int32, operation func(posting *Posting)) {
	postingList := &t.Inverse.PostingLists[termID]

	for posting := &t.Inverse.Postings[postingList.FirstIndex]; ; posting = &t.Inverse.Postings[posting.NextPostingIndex] {
		operation(posting)

		if posting.NextPostingIndex == -1 {
			break
		}
	}
}

func (t *TotalIndex) LoopOverDocumentPostings(docID int32, operation func(posting *Posting)) {
	postingList := &t.Forward.PostingLists[docID]

	if postingList.FirstIndex == -1 {
		fmt.Printf("DocId has first index -1: %d, Postinglist: %v\n", docID, postingList)
		return
	}

	for posting := &t.Forward.Postings[postingList.FirstIndex]; ; posting = &t.Forward.Postings[posting.NextPostingIndex] {
		operation(posting)
		if posting.NextPostingIndex == -1 {
			break
		}
	}
}

func (t *TotalIndex) SerialiseTo(w io.Writer) error {
	gzWriter := gzip.NewWriter(w)
	encoder := gob.NewEncoder(gzWriter)
	err := encoder.Encode(t)
	if err != nil {
		return err
	}
	return gzWriter.Close()
}

func (t *TotalIndex) SerialiseToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("unable to open file: %s", err)
	}
	return t.SerialiseTo(f)
}

func (t *TotalIndex) DeserialiseFrom(r io.Reader) error {
	gzReader, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(gzReader)
	err = decoder.Decode(t)
	if err != nil {
		return err
	}

	return gzReader.Close()
}

func (t *TotalIndex) Normalise() {
	for docId := int32(0); docId < int32(len(t.Forward.PostingLists)); docId++ {
		normalise := func(posting *Posting) {
			if t.Documents[docId].UniqueLength != 0 {
				posting.NormalisedCount = float32(posting.Count) / float32(t.Documents[docId].UniqueLength)
			}
		}

		t.LoopOverDocumentPostings(docId, normalise)
	}
}

func (t *TotalIndex) DeserialiseFromFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("unable to open file: %s", err)
	}
	return t.DeserialiseFrom(f)
}

func (t *TotalIndex) Verify() {
	for docID := range t.Forward.PostingLists {
		var lastPosting *Posting
		t.LoopOverDocumentPostings(int32(docID), func(posting *Posting) {
			if lastPosting != nil {
				if posting.Index <= lastPosting.Index {
					panic(fmt.Sprintf(
						"consecutive postings of document %d have out of order term indices: %d, %d",
						docID, lastPosting.Index, posting.Index,
					))
				}
			}
			lastPosting = posting
		})
		if lastPosting == nil {
			panic(fmt.Sprintf("document %d has no terms", docID))
		}
	}

	for termID := range t.Inverse.PostingLists {
		var lastPosting *Posting
		t.LoopOverTermPostings(int32(termID), func(posting *Posting) {
			if lastPosting != nil {
				if posting.Index <= lastPosting.Index {
					panic(fmt.Sprintf(
						"consecutive postings of term %d have out of order document indices: %d, %d",
						termID, lastPosting.Index, posting.Index,
					))
				}
			}
			lastPosting = posting
		})
		if lastPosting == nil {
			panic(fmt.Sprintf("term %d has no documents", termID))
		}
	}
}
