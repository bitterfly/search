package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/bitterfly/search/documents"
	"github.com/bitterfly/search/indices"
	"github.com/bitterfly/search/kmeans"
	"github.com/bitterfly/search/processing"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "k_means"
	app.Usage = "Takes a ready index and a file, containing document, and finds its closest centroid"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "index, i",
			Usage: "File with index",
			Value: "/tmp/kmeans.gob.gz",
		},
		cli.StringFlag{
			Name:  "d",
			Usage: "File containing document",
			Value: "/tmp/document.gob.gz",
		},

		cli.StringFlag{
			Name:  "stopwords, s",
			Usage: "Stopwords file. If not specified, defaults to ${xmldir}/stopwords",
			Value: "",
		},
	}

	app.Action = mainCommand

	app.Run(os.Args)
}

func mainCommand(c *cli.Context) {
	ti := indices.NewTotalIndex()
	err := ti.DeserialiseFromFile(c.String("index"))
	if err != nil {
		log.Fatal(err)
	}

	parser := documents.NewReutersParser()
	documents, error := parser.ParseFile(c.String("d"))
	if error != nil {
		log.Fatal(err)
	}

	tokeniser, err := processing.NewEnglishTokeniserFromFile(c.String("s"))
	if err != nil {
		log.Fatal("unable to get stopwords: %s", err)
	}

	docInfo := processing.Count(documents[0], tokeniser)
	i := kmeans.ClosestCentroidToInfo(ti, docInfo)

	purity, allClasses := kmeans.Purity(ti, len((*ti).Centroids))

	fmt.Printf("Purity %.3f\n\nDocument: %s\nis in cluster %d with most common classes:\n", purity, docInfo, i)
	clusterClasses := sortDict(allClasses[i])

	for i := 0; i < min(8, len(clusterClasses)); i++ {
		fmt.Printf("%s\n", ti.ClassNames.GetInverse(clusterClasses[i]))
	}
}

type IndexCount struct {
	index int32
	count int
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func sortDict(dict map[int32]int) []int32 {
	classesCounts := make([]IndexCount, 0, len(dict))
	for k, v := range dict {
		classesCounts = append(classesCounts, IndexCount{index: k, count: v})
	}

	sort.Slice(classesCounts, func(i, j int) bool { return classesCounts[i].count > classesCounts[j].count })

	classes := make([]int32, len(dict), len(dict))

	for i := 0; i < len(classesCounts); i++ {
		classes[i] = classesCounts[i].index
	}

	return classes
}
