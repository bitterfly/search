package kmeans

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"time"

	"github.com/bitterfly/search/indices"
	"github.com/bitterfly/search/utils"
)

func ProcessArguments(index *indices.TotalIndex, k int) {
	index.Normalise()

	KMeans(index, k)
}

func tableise(index *indices.TotalIndex, k int) [][]int {
	t := make([][]int, index.ClassNames.Size, index.ClassNames.Size)
	for i := 0; i < int(index.ClassNames.Size); i++ {
		t[i] = make([]int, k, k)
	}

	for _, doc := range index.Documents {
		for _, class := range doc.Classes {
			t[class][doc.ClusterID] += 1
		}
	}

	return t
}

func Purity(index *indices.TotalIndex, k int) (float64, []map[int32]int) {
	uc := make([]map[int32]int, k, k)
	docWithClasses := 0

	for i, _ := range uc {
		uc[i] = make(map[int32]int)
	}

	for _, doc := range index.Documents {
		if len(doc.Classes) != 0 {
			docWithClasses += 1
		}
		for _, class := range doc.Classes {
			uc[doc.ClusterID][class] += 1
		}
	}

	sum := 0
	var max int

	for i := 0; i < k; i++ {
		max = getMaxDict(uc[i])
		sum += max
	}

	return float64(sum) / float64(docWithClasses), uc
}

//returns how many documents in cluster have the most common class and the class id
func getMaxDict(d map[int32]int) int {
	max := 0
	for _, v := range d {
		if v > max {
			max = v
		}

	}

	return max
}

func KMeans(index *indices.TotalIndex, k int) []float64 {
	// Initialise this set in order to produce k random indices, because there isn't a way to get k
	// random numbers at once

	centroidIndices := make(map[int32]struct{})

	// Keep generating a new random number until there are k keys in centroidIndices
	rand.Seed(time.Now().UTC().UnixNano())
	for len(centroidIndices) < k {
		ind := rand.Int31n(int32(len(index.Forward.PostingLists)))
		if _, ok := centroidIndices[ind]; !ok {
			centroidIndices[ind] = struct{}{}
		}
	}

	// Create the the empty centroids which are k vectors each having length of the total number of terms
	index.Centroids = make([][]float64, k, k)
	for i := 0; i < k; i++ {
		index.Centroids[i] = make([]float64, len(index.Inverse.PostingLists), len(index.Inverse.PostingLists))
	}

	// Make the k documents corresponding to the indices we've fetched in the previous step the new centroids
	i := 0
	for docID, _ := range centroidIndices {
		index.LoopOverDocumentPostings(docID, func(posting *indices.Posting) {
			// index.Centroids[i][posting.Index] = tf(posting.Count, index.Documents[docID].Length) * idf(index, posting.Index)
			index.Centroids[i][posting.Index] = tf(posting.Count, index.Documents[docID].Length)
		})
		// fmt.Printf("CentroidId %d\n", docID)
		i++
	}

	iterations := 100
	rsss := make([]float64, 0, iterations)
	// Group the documents in clusters and recalculate the new centroid of the cluster
	// fmt.Printf("start")
	for times := 0; times < iterations; times++ {
		// fmt.Printf("\riteration %4d", times)

		rsss = append(rsss, rss(index))

		// break if there is no difference between new and old centroids
		if times > 1 && rsss[times-1]-rsss[times] < 0.00001 {
			break
		}

		docIdChannel := make(chan int32)

		go func() {
			for docID := int32(0); docID < int32(len(index.Forward.PostingLists)); docID++ {
				docIdChannel <- docID
			}
			close(docIdChannel)
		}()

		utils.Parallel(func() {
			for docID := range docIdChannel {
				centroidIndex := closestCentroid(index, docID)
				index.Documents[docID].ClusterID = centroidIndex
			}
		}, runtime.NumCPU())

		NewCentroids(index, k)
	}
	return rsss
}

func tf(count, len int32) float64 {
	if len == 0 {
		return float64(0)
	}
	return float64(count) / float64(len)
}

func idf(index *indices.TotalIndex, termID int32) float64 {
	return math.Log(float64(len(index.Documents)) / float64(1+index.Inverse.PostingLists[termID].Len))
}

// Empty old centroids
// Cycle through all documents and add to the corresponding index and count the number of documents in this
// cluster with the numberOfDocuments array in order to normalise later
func NewCentroids(index *indices.TotalIndex, k int) {
	for i := 0; i < k; i++ {
		for j := 0; j < len(index.Inverse.PostingLists); j++ {
			(*index).Centroids[i][j] = 0
		}
	}

	numberOfDocuments := make([]int32, k, k)
	for docID, doc := range index.Documents {
		index.LoopOverDocumentPostings(int32(docID), func(posting *indices.Posting) {
			// (*index).Centroids[doc.ClusterID][posting.Index] += tf(posting.Count, doc.Length) * idf(index, posting.Index)
			(*index).Centroids[doc.ClusterID][posting.Index] += tf(posting.Count, doc.Length)
		})
		numberOfDocuments[doc.ClusterID] += 1
	}

	for i := 0; i < k; i++ {
		for j := 0; j < len(index.Inverse.PostingLists); j++ {
			(*index).Centroids[i][j] /= float64(numberOfDocuments[i])
		}
	}
}

func sqr(x float64) float64 {
	return x * x
}

// Returns the sum of the distance between a centroid and the documents in ints cluster for all the clusters
func rss(index *indices.TotalIndex) float64 {
	var sum float64

	for docID, doc := range index.Documents {
		if doc.ClusterID == -1 {
			return -1
		}

		sum += distance(index, doc.ClusterID, int32(docID))
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

// Finds the centroid with minimal distance to the document
func closestCentroid(index *indices.TotalIndex, documentId int32) int {
	min := float64(math.MaxFloat32)
	ind := -1

	for i := 0; i < len((*index).Centroids); i++ {
		dist := distance(index, i, documentId)
		if dist < min {
			min = dist
			ind = i

		}
	}

	return ind
}

// Finds the squared distance between the centroid (witch is an array with exact length of the total number of terms)
// and a document (which is a much sparser array)
func distance(index *indices.TotalIndex, centroidIndex int, documentId int32) float64 {
	sum := float64(0)
	doclen := index.Documents[documentId].Length
	centroid := (*index).Centroids[centroidIndex]

	posting := &index.Forward.Postings[index.Forward.PostingLists[documentId].FirstIndex]
	ind := posting.Index

	for i := 0; i < len(centroid); i++ {
		if int32(i) == ind {
			// sum += sqr(centroid[i] - tf(posting.Count, doclen)*idf(index, posting.Index))
			sum += sqr(centroid[i] - tf(posting.Count, doclen))
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

type IndexedTerm struct {
	index int32
	count int32
}

func distanceToInfo(index *indices.TotalIndex, centroidIndex int, info *indices.InfoAndTerms) float64 {
	sum := float64(0)
	centroid := (*index).Centroids[centroidIndex]

	termIndices := make([]IndexedTerm, 0, info.TermsAndCounts.Size())
	var ind int32
	info.TermsAndCounts.Walk(func(word []byte, value int32) {
		ind = index.Dictionary.Get(word)
		if ind != -1 {
			termIndices = append(termIndices, IndexedTerm{index: ind, count: value})
		}
	})

	sort.Slice(termIndices, func(i, j int) bool { return termIndices[i].index < termIndices[j].index })

	j := 0
	for i := 0; i < len(centroid); i++ {
		if i == int(termIndices[j].index) {
			sum += sqr(centroid[i] - tf(termIndices[j].count, info.Length)*idf(index, int32(i)))
			if j < len(termIndices)-1 {
				j++
			}
		}

		sum += sqr(centroid[i])
	}
	return sum
}

func ClosestCentroidToInfo(index *indices.TotalIndex, info *indices.InfoAndTerms) int {
	min := float64(math.MaxFloat32)
	ind := -1

	for i := 0; i < len((*index).Centroids); i++ {
		dist := distanceToInfo(index, i, info)
		if dist < min {
			min = dist
			ind = i

		}
	}

	return ind
}
