package main

import (
	"fmt"
	"strings"

	"github.com/niranjankrishna-acad/gokapiX/bm25"
)

func main() {
	// Original corpus
	corpus := []string{
		"Hello there good man!",
		"It is quite windy in London",
		"How is the weather today?",
	}

	// Tokenize corpus
	tokenizedCorpus := make([][]string, len(corpus))
	for i, doc := range corpus {
		tokenizedCorpus[i] = strings.Fields(doc)
	}

	// Query
	query := "windy London"
	tokenizedQuery := strings.Fields(query)

	// Compute BM25
	bm25Result := bm25.BM25OkapiCompute(tokenizedQuery, tokenizedCorpus, 0.75, 1.2)

	// Retrieve the highest-ranked document
	topDocIndex := bm25Result.TopN[0]
	topDoc := corpus[topDocIndex]

	fmt.Println("Query:", query)
	fmt.Println("Top-ranked document:", topDoc)
}
