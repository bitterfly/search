package kmeans

import (
	"github.com/DexterLB/search/indices"
)

func NormaliseIndex(index indices.TotalIndex) {
	index.Normalise()
}
