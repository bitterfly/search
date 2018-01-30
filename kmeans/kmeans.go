package kmeans

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/DexterLB/search/utils"
	"github.com/bitterfly/search/indices"
)

func ProcessArguments(index *indices.TotalIndex, k int) {
	index.Normalise()

	KMeans(index, k)

	// fmt.Printf("%d %.5f %d\n", k, rsss[len(rsss)-1], uniqueClasses(index, k))

	t := tableise(index, k)

	for i := 0; i < int(index.ClassNames.Size); i++ {
		for j := 0; j < k; j++ {
			if j < int(index.ClassNames.Size-1) {
				fmt.Printf("%d ", t[i][j])
			} else {
				fmt.Printf("%d", t[i][j])
			}
		}
		fmt.Printf("\n")
	}

	// fmt.Printf("\n")
	// for i := 1; i < len(rsss); i++ {
	// 	fmt.Printf("%d %7.5f\t%.5f-%.5f\n", i, rsss[i-1]-rsss[i], rsss[i-1], rsss[i])
	// }
	// PrintClusters(index, k)
	// fmt.Printf("%d\n", k)

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

func uniqueClasses(index *indices.TotalIndex, k int) int {
	uc := make([]map[int32]int, k, k)

	classes := make([]int, k, k)

	for i, _ := range uc {
		uc[i] = make(map[int32]int)
	}

	for _, doc := range index.Documents {
		for _, class := range doc.Classes {
			uc[doc.ClusterID][class] += 1
		}
		classes[doc.ClusterID] += 1
	}

	sum := 0
	for i := 0; i < k; i++ {
		sum += classes[i] - getMaxDict(uc[i])
	}

	return sum
}

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
	centroids := make([][]float64, k, k)
	for i := 0; i < k; i++ {
		centroids[i] = make([]float64, len(index.Inverse.PostingLists), len(index.Inverse.PostingLists))
	}

	// Make the k documents corresponding to the indices we've fetched in the previous step the new centroids
	i := 0
	for docID, _ := range centroidIndices {
		index.LoopOverDocumentPostings(docID, func(posting *indices.Posting) {
			centroids[i][posting.Index] = float64(posting.Count)
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

		rsss = append(rsss, rss(index, centroids))

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
				centroidIndex := closestCentroid(index, &centroids, docID)
				index.Documents[docID].ClusterID = centroidIndex
			}
		}, runtime.NumCPU())

		NewCentroids(index, k, &centroids)
	}
	fmt.Printf("\n")
	return rsss
}

// Empty old centroids
// Cycle through all documents and add to the corresponding index and count the number of documents in this
// cluster with the numberOfDocuments array in order to normalise later
func NewCentroids(index *indices.TotalIndex, k int, centroids *[][]float64) {
	for i := 0; i < k; i++ {
		for j := 0; j < len(index.Inverse.PostingLists); j++ {
			(*centroids)[i][j] = 0
		}
	}

	numberOfDocuments := make([]int32, k, k)
	for docID, doc := range index.Documents {
		index.LoopOverDocumentPostings(int32(docID), func(posting *indices.Posting) {
			(*centroids)[doc.ClusterID][posting.Index] += float64(posting.Count)
		})
		numberOfDocuments[doc.ClusterID] += 1
	}

	for i := 0; i < k; i++ {
		for j := 0; j < len(index.Inverse.PostingLists); j++ {
			(*centroids)[i][j] /= float64(numberOfDocuments[i])
		}
	}
}

func sqr(x float64) float64 {
	return x * x
}

// Returns the sum of the distance between a centroid and the documents in ints cluster for all the clusters
func rss(index *indices.TotalIndex, centroids [][]float64) float64 {
	var sum float64

	for docID, doc := range index.Documents {
		if doc.ClusterID == -1 {
			return math.MaxFloat64
		}

		sum += distance(index, centroids[doc.ClusterID], int32(docID))
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
func closestCentroid(index *indices.TotalIndex, centroids *[][]float64, documentId int32) int {
	min := float64(math.MaxFloat32)
	ind := -1

	for i := 0; i < len(*centroids); i++ {
		dist := distance(index, (*centroids)[i], documentId)
		if dist < min {
			min = dist
			ind = i

		}
	}

	return ind
}

// Finds the squared distance between the centroid (witch is an array with exact length of the total number of terms)
// and a document (which is a much sparser array)
func distance(index *indices.TotalIndex, centroid []float64, documentId int32) float64 {
	sum := float64(0)

	posting := &index.Forward.Postings[index.Forward.PostingLists[documentId].FirstIndex]
	ind := posting.Index

	for i := 0; i < len(centroid); i++ {
		if int32(i) == ind {
			sum += sqr(centroid[i] - float64(posting.Count))
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
