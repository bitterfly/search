package kmeans

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/bitterfly/search/indices"
)

func KMeans(index *indices.TotalIndex, k int) {
	index.Normalise()

	RealKMeans(index, k)
	PrintClusters(index, k)

	fmt.Printf("%d\n", k)
}

func RealKMeans(index *indices.TotalIndex, k int) {
	centroidIndices := make(map[int32]struct{})

	rand.Seed(time.Now().UTC().UnixNano())
	for len(centroidIndices) < k {
		ind := rand.Int31n(int32(len(index.Forward.PostingLists)))
		if _, ok := centroidIndices[ind]; !ok {
			centroidIndices[ind] = struct{}{}
		}
	}

	centroids := make([][]float32, k, k)

	//initial centroids are k random documents
	for i := 0; i < k; i++ {
		centroids[i] = make([]float32, len(index.Inverse.PostingLists), len(index.Inverse.PostingLists))
	}

	i := 0
	for docID, _ := range centroidIndices {
		index.LoopOverDocumentPostings(docID, func(posting *indices.Posting) {
			centroids[i][posting.Index] = float32(posting.Count)
		})
		fmt.Printf("CentroidId %d\n", docID)
		i++
	}

	for times := 0; times < 10; times++ {
		fmt.Printf("%d: Rss: %.3f\n", times, rss(index, centroids))
		PrintClusters(index, k)

		fmt.Printf("=======\n")
		for i := 0; i < k; i++ {
			for docID := int32(0); docID < int32(len(index.Forward.PostingLists)); docID++ {
				centroidIndex := closestCentroid(docID, &centroids, index)
				index.Documents[docID].ClusterID = centroidIndex
			}
		}

		NewCentroids(index, k, &centroids)
	}
}

func NewCentroids(index *indices.TotalIndex, k int, centroids *[][]float32) {
	for i := 0; i < k; i++ {
		for j := 0; j < len(index.Inverse.PostingLists); j++ {
			(*centroids)[i][j] = 0
		}
	}

	numberOfDocuments := make([]int32, k, k)
	for docID, doc := range index.Documents {
		index.LoopOverDocumentPostings(int32(docID), func(posting *indices.Posting) {
			(*centroids)[doc.ClusterID][posting.Index] += float32(posting.Count)
		})
		numberOfDocuments[doc.ClusterID] += 1
	}

	for i := 0; i < k; i++ {
		for j := 0; j < len(index.Inverse.PostingLists); j++ {
			(*centroids)[i][j] /= float32(numberOfDocuments[i])
		}
	}
}

func sqr(x float32) float32 {
	return x * x
}

func rss(index *indices.TotalIndex, centroids [][]float32) float32 {
	var sum float32

	for docID, doc := range index.Documents {
		sum += distance(int32(docID), centroids[doc.ClusterID], index)
	}
	return sum
}

func PrintClusters(index *indices.TotalIndex, k int) {
	clusterNums := make([]int, k, k)
	for _, doc := range index.Documents {
		clusterNums[doc.ClusterID] += 1
	}

	for i := 0; i < len(clusterNums); i++ {
		fmt.Printf("%d: %d\n", i, clusterNums[i])
	}
}

func closestCentroid(documentId int32, centroids *[][]float32, index *indices.TotalIndex) int {
	min := float32(math.MaxFloat32)
	ind := -1

	for i := 0; i < len(*centroids); i++ {
		dist := distance(documentId, (*centroids)[i], index)
		if dist < min {
			min = dist
			ind = i

		}
	}

	return ind
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}

	return b
}

func distance(documentId int32, centroid []float32, index *indices.TotalIndex) float32 {
	sum := float32(0)

	posting := &index.Forward.Postings[index.Forward.PostingLists[documentId].FirstIndex]
	ind := posting.Index

	for i := 0; i < len(centroid); i++ {
		if int32(i) == ind {
			sum += sqr(centroid[i] - float32(posting.Count))
			if posting.NextPostingIndex != int32(-1) {
				posting = &index.Forward.Postings[posting.NextPostingIndex]
				ind = posting.Index
			} else {
				ind = -1
			}
		} else {
			sum += sqr(centroid[i])
		}
	}

	return sum
}

func similarity(documentId int32, centroid []float32, index *indices.TotalIndex) float32 {
	sum := float32(0)
	index.LoopOverDocumentPostings(documentId, func(posting *indices.Posting) {
		sum += posting.NormalisedCount * centroid[posting.Index]
	})

	return sum
}
