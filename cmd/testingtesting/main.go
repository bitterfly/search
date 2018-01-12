package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DexterLB/search/documents"
)

func main() {
	docs, err := documents.ParseFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	for i := range docs {
		fmt.Printf("********\n%s\n", &docs[i])
	}
}
