package kmeans

import (
	"testing"

	"github.com/bitterfly/search/indices"
	"github.com/stretchr/testify/assert"
)

func makeIndex() *indices.TotalIndex {
	doc0 := indices.NewInfoAndTerms()
	doc0.TermsAndCounts.Put([]byte("foo"), 2)
	doc0.TermsAndCounts.Put([]byte("bar"), 1)

	doc1 := indices.NewInfoAndTerms()
	doc1.TermsAndCounts.Put([]byte("bar"), 1)
	doc1.TermsAndCounts.Put([]byte("qux"), 1)

	doc0.Name = "doc0"
	doc0.Classes = []string{"sports", "dodgeball"}
	doc0.Length = 3

	doc1.Name = "doc1"
	doc1.Classes = []string{"politics"}
	doc1.Length = 2

	ti := indices.NewTotalIndex()
	ti.Add(doc0)
	ti.Add(doc1)

	ti.Documents[0].UniqueLength = 2
	ti.Documents[1].UniqueLength = 2

	return ti
}

func TestNormalise(t *testing.T) {
	assert := assert.New(t)

	ti := makeIndex()

	ti.Normalise()

	assert.InDelta(float32(1), ti.Forward.Postings[ti.Forward.PostingLists[0].FirstIndex].NormalisedCount, float64(0.0001))
	assert.InDelta(float32(0.5), ti.Forward.Postings[ti.Forward.PostingLists[0].LastIndex].NormalisedCount, float64(0.0001))
	assert.InDelta(float32(0.5), ti.Forward.Postings[ti.Forward.PostingLists[1].FirstIndex].NormalisedCount, float64(0.0001))
	assert.InDelta(float32(0.5), ti.Forward.Postings[ti.Forward.PostingLists[1].LastIndex].NormalisedCount, float64(0.0001))
}

func TestDistance(t *testing.T) {
	assert := assert.New(t)

	ti := makeIndex()

	ti.Normalise()

	dist := ((0-2)*(0-2) + (0.5-1)*(0.5-1) + (1-0)*(1-0))

	assert.InDelta(dist, distance(0, []float32{0, 0.5, 1}, ti), float64(0.0001))
}

func TestNewCentroid(t *testing.T) {
	assert := assert.New(t)

	ti := makeIndex()

	ti.Normalise()

	assert.InDeltaSlice([]float32{1, 1, 0.5}, newCentroid([]int32{0, 1}, ti), 0.001)
}

func TestClosestCentroid(t *testing.T) {
	assert := assert.New(t)

	ti := makeIndex()

	ti.Normalise()

	// d0 = {1, 0.5}
	// d1 = {0.5, 0.5}

	centroids := make([][]float32, 2, 2)
	centroids[0] = []float32{1, 0, 0}
	centroids[1] = []float32{0, 1, 0}

	assert.Equal(0, closestCentroid(0, &centroids, ti))
	// assert.Equal(int32(1), closestCentroid(1, &centroids, ti))

}
