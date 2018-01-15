package kmeans

import (
	"fmt"

	"github.com/DexterLB/search/indices"
)

func KMeans(index *indices.TotalIndex, k int) {
	index.Normalise()

	index.LoopOverDocumentPostings(0, func(posting *indices.Posting) { fmt.Printf("%d(%.3f) ", posting.Index, posting.NormalisedCount) })
	fmt.Printf("\n\n")
	index.LoopOverDocumentPostings(1, func(posting *indices.Posting) { fmt.Printf("%d(%.3f) ", posting.Index, posting.NormalisedCount) })
	fmt.Printf("\n\n")

	fmt.Printf("sum: %.3f\n", distance(0, 1, index))

	fmt.Printf("%d\n", k)
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}

	return b
}

func distance(firstDocId, secondDocId int32, index *indices.TotalIndex) float32 {
	d1Posting := &index.Forward.Postings[index.Forward.PostingLists[firstDocId].FirstIndex]
	d2Posting := &index.Forward.Postings[index.Forward.PostingLists[secondDocId].FirstIndex]

	sum := float32(0)

	for {

		if d1Posting.Index == d2Posting.Index {
			fmt.Printf("They are the same on index: %d\n", d1Posting.Index)

			sum += d1Posting.NormalisedCount * d2Posting.NormalisedCount
			if d1Posting.NextPostingIndex != -1 {
				d1Posting = &index.Forward.Postings[d1Posting.NextPostingIndex]
			} else {
				break
			}

			if d2Posting.NextPostingIndex != -1 {
				d2Posting = &index.Forward.Postings[d2Posting.NextPostingIndex]
			} else {
				break
			}
		}

		if d1Posting.Index < d2Posting.Index {
			if d1Posting.NextPostingIndex != -1 {
				d1Posting = &index.Forward.Postings[d1Posting.NextPostingIndex]
			} else {
				break
			}
		}

		if d1Posting.Index > d2Posting.Index {
			if d2Posting.NextPostingIndex != -1 {
				d2Posting = &index.Forward.Postings[d2Posting.NextPostingIndex]
			} else {
				break
			}
		}
	}

	return sum

	// for posting := &t.Forward.Postings[postingList.FirstIndex]; posting.NextPostingIndex != -1; posting = &t.Forward.Postings[posting.NextPostingIndex] {
}
