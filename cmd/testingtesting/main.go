package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/DexterLB/search/documents"
	"github.com/DexterLB/search/indices"
	"github.com/DexterLB/search/processing"
	"github.com/DexterLB/search/utils"
)

func GetXMLs(folder string, into chan<- string) {
	files, err := filepath.Glob(filepath.Join(folder, "*.xml"))
	if err != nil {
		log.Fatal("unable to get files in folder %s: %s", os.Args[1], err)
	}

	for i := range files {
		into <- files[i]
	}
}

func main() {
	files := make(chan string, 200)
	docs := make(chan *documents.Document, 2000)
	countDocs := make(chan *indices.Document, 2000)

	tokeniser, err := processing.NewEnglishTokeniserFromFile(
		filepath.Join(os.Args[1], "stopwords"),
	)
	if err != nil {
		log.Fatal("unable to get stopwords: %s", err)
	}

	go func() {
		GetXMLs(os.Args[1], files)
		close(files)
	}()

	go func() {
		utils.Parallel(func() {
			documents.NewReutersParser().ParseFiles(files, docs)
		}, runtime.NumCPU())
		close(docs)
	}()

	go func() {
		utils.Parallel(func() {
			processing.CountInDocuments(docs, tokeniser, countDocs)
		}, runtime.NumCPU())
		close(countDocs)
	}()

	for countDoc := range countDocs {
		countDoc.Print()
	}

}
