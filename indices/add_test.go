package indices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	assert := assert.New(t)

	doc0 := NewInfoAndTerms()
	doc0.TermsAndCounts.Put([]byte("foo"), 2)
	doc0.TermsAndCounts.Put([]byte("bar"), 1)
	doc0.TermsAndCounts.Put([]byte("foo"), 1)

	doc1 := NewInfoAndTerms()
	doc1.TermsAndCounts.Put([]byte("bar"), 1)
	doc1.TermsAndCounts.Put([]byte("qux"), 1)

	doc0.DocumentInfo.Name = "doc0"
	doc0.Classes = []string{"sports", "dodgeball"}
	doc0.Length = 3

	doc1.DocumentInfo.Name = "doc1"
	doc1.Classes = []string{"politics"}
	doc1.Length = 2

	ti := NewTotalIndex()
	ti.Add(doc0)
	ti.Add(doc1)

	// Ensure order of terms in dictionary for easier testing
	ti.Dictionary.Get([]byte("foo"))
	ti.Dictionary.Get([]byte("bar"))
	ti.Dictionary.Get([]byte("qux"))

	// IDs of terms in the dictionary are in the same order as
	// they've been put
	assert.Equal(0, ti.Dictionary.Get([]byte("foo")))
	assert.Equal(1, ti.Dictionary.Get([]byte("bar")))
	assert.Equal(2, ti.Dictionary.Get([]byte("qux")))

	// Verify document information
	assert.Equal("doc0", ti.Documents[0].Name)
	assert.Equal([]string{"sports", "dodgeball"}, ti.Documents[0].Classes)
	assert.Equal(3, ti.Documents[0].Length)

	assert.Equal("doc1", ti.Documents[1].Name)
	assert.Equal([]string{"politics"}, ti.Documents[1].Classes)
	assert.Equal(2, ti.Documents[1].Length)

	// Forward index for doc0
	assert.Equal(
		ti.Dictionary.Get([]byte("foo")),

		ti.Forward.Postings[ti.Forward.PostingLists[0].FirstIndex].Index,
	)
	assert.Equal(
		2, // "foo" occured twice in doc0

		ti.Forward.Postings[ti.Forward.PostingLists[0].FirstIndex].Count,
	)

	assert.Equal(
		ti.Dictionary.Get([]byte("bar")),

		ti.Forward.Postings[ti.Forward.PostingLists[0].LastIndex].Index,
	)
	assert.Equal(
		1, // "bar" occured once in doc0

		ti.Forward.Postings[ti.Forward.PostingLists[0].LastIndex].Count,
	)
	assert.Equal(
		ti.Forward.PostingLists[0].LastIndex,
		ti.Forward.Postings[ti.Forward.PostingLists[0].FirstIndex].NextPostingIndex,
	)
	assert.Equal(
		-1,
		ti.Forward.Postings[ti.Forward.PostingLists[0].LastIndex].NextPostingIndex,
	)

	// Forward index for doc1
	assert.Equal(
		ti.Dictionary.Get([]byte("bar")),

		ti.Forward.Postings[ti.Forward.PostingLists[1].FirstIndex].Index,
	)
	assert.Equal(
		1, // "bar" occured once in doc1

		ti.Forward.Postings[ti.Forward.PostingLists[1].FirstIndex].Count,
	)

	assert.Equal(
		ti.Dictionary.Get([]byte("qux")),

		ti.Forward.Postings[ti.Forward.PostingLists[1].LastIndex].Index,
	)
	assert.Equal(
		1, // "qux" occured once in doc1

		ti.Forward.Postings[ti.Forward.PostingLists[1].LastIndex].Count,
	)
	assert.Equal(
		ti.Forward.PostingLists[1].LastIndex,
		ti.Forward.Postings[ti.Forward.PostingLists[1].FirstIndex].NextPostingIndex,
	)
	assert.Equal(
		-1,
		ti.Forward.Postings[ti.Forward.PostingLists[1].LastIndex].NextPostingIndex,
	)

	// Inverse index for "foo"
	assert.Equal(
		0,	// "foo" occured in doc0

		ti.Inverse.Postings[ti.Inverse.PostingLists[ti.Dictionary.Get([]byte("foo"))].FirstIndex].Index
	)
	assert.Equal(
		2,	// "foo" occured twice in doc0

		ti.Inverse.Postings[ti.Inverse.PostingLists[ti.dictionary.Get([]byte("foo"))].FirstIndex].Count
	)
	assert.Equal(
		// "foo" occurs only in one document, so first and last are the same
		ti.Inverse.PostingLists[ti.Dictionary.Get([]byte("foo"))].LastIndex,
		ti.Inverse.PostingLists[ti.Dictionary.Get([]byte("foo"))].FirstIndex,
	)
	assert.Equal(
		// No next
		-1,
		ti.Inverse.Postings[ti.Inverse.PostingLists[ti.Dictionary.Get([]byte("foo"))].FirstIndex].NextPosting,
	)

	// Inverse index for "qux"
	assert.Equal(
		1,	// "qux" occured in doc1

		ti.Inverse.Postings[ti.Inverse.PostingLists[ti.Dictionary.Get([]byte("qux"))].FirstIndex].Index
	)
	assert.Equal(
		1,	// "qux" occured once in doc1

		ti.Inverse.Postings[ti.Inverse.PostingLists[ti.dictionary.Get([]byte("qux"))].FirstIndex].Count
	)
	assert.Equal(
		// "qux" occurs only in one document, so first and last are the same
		ti.Inverse.PostingLists[ti.Dictionary.Get([]byte("qux"))].LastIndex,
		ti.Inverse.PostingLists[ti.Dictionary.Get([]byte("qux"))].FirstIndex,
	)
	assert.Equal(
		// No next
		-1,
		ti.Inverse.Postings[ti.Inverse.PostingLists[ti.Dictionary.Get([]byte("qux"))].FirstIndex].NextPosting,
	)

	// Inverse index for "bar"
	assert.Equal(
		0,	// "bar" occured in doc0

		ti.Inverse.Postings[ti.Inverse.PostingLists[ti.Dictionary.Get([]byte("bar"))].FirstIndex].Index
	)
	assert.Equal(
		1,	// "bar" occured once in doc0

		ti.Inverse.Postings[ti.Inverse.PostingLists[ti.dictionary.Get([]byte("bar"))].FirstIndex].Count
	)
	assert.Equal(
		1,	// "bar" occured in doc1

		ti.Inverse.Postings[ti.Inverse.PostingLists[ti.Dictionary.Get([]byte("bar"))].LastIndex].Index
	)
	assert.Equal(
		1,	// "bar" occured once in doc1

		ti.Inverse.Postings[ti.Inverse.PostingLists[ti.dictionary.Get([]byte("bar"))].LastIndex].Count
	)
	assert.Equal(
		ti.Inverse.PostingLists[ti.Inverse.PostingLists[ti.dictionary.Get([]byte("bar"))].LastIndex,
		ti.Inverse.Postings[ti.Forward.PostingLists[ti.Inverse.PostingLists[ti.dictionary.Get([]byte("bar"))].FirstIndex].NextPostingIndex,
	)
	assert.Equal(
		-1,
		ti.Inverse.Postings[ti.Inverse.PostingLists[ti.Dictionary.Get([]byte("bar"))].LastIndex].NextPosting,
	)
}
