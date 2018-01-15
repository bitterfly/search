package featureselection

import "github.com/DexterLB/search/indices"

type ClassInfo struct {
	DocumentsWhichHaveClass []int32
}

func ComputeClassInfo(ti *indices.TotalIndex) *ClassInfo {
	info := &ClassInfo{}

	numClasses := ti.ClassNames.Size

	info.DocumentsWhichHaveClass = make([]int32, numClasses)

	for termID := range ti.Inverse.PostingLists {
		ti.LoopOverTermPostings(termID, func(posting *indices.Posting) {
			for _, class := range ti.Documents[posting.Index].Classes {
				info.DocumentsWhichHaveClass[class] += 1
			}
		})
	}

	return info
}
