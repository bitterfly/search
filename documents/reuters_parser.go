package documents

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/DexterLB/htmlparsing"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

func ParseFile(filename string) ([]Document, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file: %s", err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("Unable to read file: %s", err)
	}

	return Parse(data)
}

func Parse(data []byte) ([]Document, error) {
	xml, err := gokogiri.ParseXml(data)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse file: %s", err)
	}

	docnodes, err := xml.Search("//REUTERS")
	if err != nil {
		return nil, fmt.Errorf("Unable to find documents in file: %s", err)
	}

	documents := make([]Document, len(docnodes))
	for i := range docnodes {
		err = parseDocument(docnodes[i], &documents[i])
		if err != nil {
			return nil, fmt.Errorf("Unable to parse document: %s", err)
		}
	}

	return documents, nil
}

func parseDocument(node xml.Node, document *Document) error {
	titleNode, err := htmlparsing.First(node, ".//TITLE")
	if err != nil {
		return fmt.Errorf("Unable to parse document title: %s", err)
	}

	document.Title = titleNode.Content()

	bodyNode, err := htmlparsing.First(node, ".//BODY")
	if err != nil {
		return fmt.Errorf("Unable to parse document body: %s", err)
	}

	document.Body = bodyNode.Content()

	dateNode, err := htmlparsing.First(node, ".//DATE")
	if err != nil {
		return fmt.Errorf("Unable to parse document date: %s", err)
	}

	document.Date = dateNode.Content()

	topicNodes, err := node.Search(".//TOPICS/d")
	if err != nil {
		return fmt.Errorf("Unable to parse document topics: %s", err)
	}

	document.Classes = make([]string, len(topicNodes))
	for i := range topicNodes {
		document.Classes[i] = topicNodes[i].Content()
	}

	return nil
}
