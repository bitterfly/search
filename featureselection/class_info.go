package featureselection

import "github.com/DexterLB/search/indices"

type ClassInfo struct {
	DocumentsWhichHaveClass   []int32
	DocumentsWhichContainTerm []int32
	NumClasses                int32
}

func ComputeClassInfo(ti *indices.TotalIndex) *ClassInfo {
	info := &ClassInfo{}

	numClasses := ti.ClassNames.Size
	numTerms := len(ti.Inverse.PostingLists)

	info.NumClasses = numClasses
	info.DocumentsWhichHaveClass = make([]int32, numClasses)
	info.DocumentsWhichContainTerm = make([]int32, numTerms)

	for termID := range ti.Inverse.PostingLists {
		ti.LoopOverTermPostings(termID, func(posting *indices.Posting) {
			for _, class := range ti.Documents[posting.Index].Classes {
				info.DocumentsWhichHaveClass[class] += 1
			}
			info.DocumentsWhichContainTerm[termID] += 1
		})
	}

	return info
}
