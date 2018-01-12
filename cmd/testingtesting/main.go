package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/DexterLB/search/documents"
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
	files := make(chan string)
	docs := make(chan *documents.Document)

	go func() {
		GetXMLs(os.Args[1], files)
		close(files)
	}()

	utils.Parallel(func() {
		documents.ParseFiles(files, docs)
	}, runtime.NumCPU())
}
