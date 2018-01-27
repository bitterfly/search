package documents

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/bitterfly/htmlparsing"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type ReutersParser struct{}

func NewReutersParser() *ReutersParser {
	return &ReutersParser{}
}

func (r *ReutersParser) ParseFiles(filenames <-chan string, documents chan<- *Document) {
	for f := range filenames {
		log.Printf("start parsing %s", f)
		documentsInFile, err := r.ParseFile(f)
		if err != nil {
			log.Printf("Unable to parse file %s: %s", f, err)
			// too lazy for proper error handling
		} else {
			log.Printf("finish parsing %s", f)
		}

		for _, doc := range documentsInFile {
			documents <- doc
		}
	}
}

func (r *ReutersParser) ParseFile(filename string) ([]*Document, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file: %s", err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("Unable to read file: %s", err)
	}

	return r.Parse(data)
}

func (r *ReutersParser) Parse(data []byte) ([]*Document, error) {
	xml, err := gokogiri.ParseXml(data)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse file: %s", err)
	}

	docnodes, err := xml.Search("//REUTERS")
	if err != nil {
		return nil, fmt.Errorf("Unable to find documents in file: %s", err)
	}

	var documents []*Document

	for i := range docnodes {
		doc := Document{}
		err = r.parseDocument(docnodes[i], &doc)
		if err != nil {
			// return nil, fmt.Errorf("Unable to parse document: %s", err)
			if strings.Contains(err.Error(), "Unable to parse document body") {
				continue // this message is too irritating
			}
			log.Printf("Unable to parse document: %s", err)
			continue
		}
		documents = append(documents, &doc)
	}

	return documents, nil
}

func (r *ReutersParser) parseDocument(node xml.Node, document *Document) error {
	titleNode, err := htmlparsing.First(node, ".//TITLE")
	if err == nil {
		// don't care if there's no title
		document.Title = titleNode.Content()
	}

	bodyNode, err := htmlparsing.First(node, ".//BODY")
	if err != nil {
		return fmt.Errorf("Unable to parse document body: %s", err)
	}

	document.Body = bodyNode.Content()

	dateNode, err := htmlparsing.First(node, ".//DATE")
	if err == nil {
		document.Date = dateNode.Content()
	}

	topicNodes, err := node.Search(".//TOPICS/D")
	if err != nil {
		return fmt.Errorf("Unable to parse document topics: %s", err)
	}

	document.Classes = make([]string, len(topicNodes))
	for i := range topicNodes {
		document.Classes[i] = topicNodes[i].Content()
	}

	return nil
}
