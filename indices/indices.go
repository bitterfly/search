package indices

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"os"

	"github.com/DexterLB/search/trie"
)

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

func (t *TotalIndex) DeserialiseFromFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("unable to open file: %s", err)
	}
	return t.DeserialiseFrom(f)
}
