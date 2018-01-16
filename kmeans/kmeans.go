package kmeans

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/DexterLB/search/indices"
)

func KMeans(index *indices.TotalIndex, k int) {
	index.Normalise()

	clusters := RealKMeans(index, k)
	for i, cluster := range clusters {
		fmt.Printf("Cluster %d has len %d\n", i, len(cluster))
	}

	fmt.Printf("%d\n", k)
}

func RealKMeans(index *indices.TotalIndex, k int) [][]int32 {
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
	for docId, _ := range centroidIndices {
		index.LoopOverDocumentPostings(docId, func(posting *indices.Posting) {
			centroids[i][posting.Index] = float32(posting.Count)
		})
		fmt.Printf("CentroidId %d\n", docId)
		i++
	}

	clusters := make([][]int32, k, k)

	for times := 0; times < 10; times++ {
		fmt.Printf("%d: Rss: %.3f\n", times, rss(clusters, centroids, index))
		for _, cl := range clusters {
			fmt.Printf("%d\n", len(cl))
		}
		fmt.Printf("=======\n")
		for i := 0; i < k; i++ {
			clearClusters(&clusters)
			for docId := int32(0); docId < int32(len(index.Forward.PostingLists)); docId++ {
				centroidIndex := closestCentroid(docId, &centroids, index)
				clusters[centroidIndex] = append(clusters[centroidIndex], docId)
			}
		}

		for i := 0; i < k; i++ {
			centroids[i] = newCentroid(clusters[i], index)
		}
	}
	return clusters
}

func sqr(x float32) float32 {
	return x * x
}

func rssK(cluster []int32, centroid []float32, index *indices.TotalIndex) float32 {
	var sum float32

	for _, docId := range cluster {
		index.LoopOverDocumentPostings(docId, func(posting *indices.Posting) {
			sum += sqr(float32(posting.Count) - centroid[posting.Index])
		})
	}

	return sum
}

func rss(clusters [][]int32, centroids [][]float32, index *indices.TotalIndex) float32 {
	var sum float32

	for k := 0; k < len(clusters); k++ {
		sum += rssK(clusters[k], centroids[k], index)
	}

	return sum
}

func clearClusters(clusters *[][]int32) {
	for i := 0; i < len(*clusters); i++ {
		(*clusters)[i] = nil
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

func newCentroid(documentIds []int32, index *indices.TotalIndex) []float32 {
	centroid := make([]float32, len(index.Inverse.PostingLists), len(index.Inverse.PostingLists))
	for _, docID := range documentIds {
		index.LoopOverDocumentPostings(docID, func(posting *indices.Posting) {
			centroid[posting.Index] += float32(posting.Count)
		})
	}

	for i := 0; i < len(centroid); i++ {
		centroid[i] /= float32(len(documentIds))
	}

	return centroid
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

// func distance(firstDocId, secondDocId int32, index *indices.TotalIndex) float32 {
// 	d1Posting := &index.Forward.Postings[index.Forward.PostingLists[firstDocId].FirstIndex]
// 	d2Posting := &index.Forward.Postings[index.Forward.PostingLists[secondDocId].FirstIndex]

// 	sum := float32(0)

// 	for {

// 		if d1Posting.Index == d2Posting.Index {
// 			fmt.Printf("They are the same on index: %d\n", d1Posting.Index)

// 			sum += d1Posting.NormalisedCount * d2Posting.NormalisedCount
// 			if d1Posting.NextPostingIndex != -1 {
// 				d1Posting = &index.Forward.Postings[d1Posting.NextPostingIndex]
// 			} else {
// 				break
// 			}

// 			if d2Posting.NextPostingIndex != -1 {
// 				d2Posting = &index.Forward.Postings[d2Posting.NextPostingIndex]
// 			} else {
// 				break
// 			}
// 		}

// 		if d1Posting.Index < d2Posting.Index {
// 			if d1Posting.NextPostingIndex != -1 {
// 				d1Posting = &index.Forward.Postings[d1Posting.NextPostingIndex]
// 			} else {
// 				break
// 			}
// 		}

// 		if d1Posting.Index > d2Posting.Index {
// 			if d2Posting.NextPostingIndex != -1 {
// 				d2Posting = &index.Forward.Postings[d2Posting.NextPostingIndex]
// 			} else {
// 				break
// 			}
// 		}
// 	}

// 	return sum
// }
