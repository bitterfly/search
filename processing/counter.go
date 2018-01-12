package processing

import (
	"github.com/DexterLB/search/documents"
	"github.com/DexterLB/search/indices"
)

func Count(doc *documents.Document, tokeniser Tokeniser) *indices.Document {
	idoc := indices.NewDocument()
	idoc.Name = doc.Title
	idoc.Classes = doc.Classes

	terms := tokeniser.Tokenise(doc.Body)

	for i := range terms {
		terms[i] = tokeniser.Normalise(terms[i])

		idoc.TermCounts.PutLambda(
			[]byte(terms[i]),
			func(x uint64) uint64 { return x + 1 },
			1,
		)
	}

	return idoc
}
