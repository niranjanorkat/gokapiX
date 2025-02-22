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

	// Compute BM25 Methods

	fmt.Println("\n===== BM25 =====")
	bm25Model := bm25.BM25Init(tokenizedCorpus, 0.75, 1.2)
	bm25Result := bm25.AtireQuery(tokenizedQuery, bm25Model)
	printTopResult(query, corpus, bm25Result.TopN)

	fmt.Println("\n===== BM25L =====")
	bm25LModel := bm25.BM25LInit(tokenizedCorpus, 0.75, 1.2, 0.5)
	bm25LResult := bm25.LQuery(tokenizedQuery, bm25LModel)
	printTopResult(query, corpus, bm25LResult.TopN)

	fmt.Println("\n===== BM25+ =====")
	bm25PlusModel := bm25.BM25LInit(tokenizedCorpus, 0.75, 1.2, 0.5)
	bm25PlusResult := bm25.PlusQuery(tokenizedQuery, bm25PlusModel)
	printTopResult(query, corpus, bm25PlusResult.TopN)

	fmt.Println("\n===== BM25-adpt =====")
	bm25AdptModel := bm25.BM25AdptInit(tokenizedCorpus, 0.75, 1.2)
	bm25AdptResult := bm25.AdptQuery(tokenizedQuery, bm25AdptModel)
	printTopResult(query, corpus, bm25AdptResult.TopN)

	fmt.Println("\n===== BM25T =====")
	bm25TModel := bm25.BM25TInit(tokenizedCorpus, 0.75, 1.2)
	bm25TResult := bm25.TQuery(tokenizedQuery, bm25TModel)
	printTopResult(query, corpus, bm25TResult.TopN)
}

func printTopResult(query string, corpus []string, topN []int) {
	if len(topN) == 0 {
		fmt.Println("No results found.")
		return
	}
	topDocIndex := topN[0]
	topDoc := corpus[topDocIndex]
	fmt.Println("Query:", query)
	fmt.Println("Top-ranked document:", topDoc)
}
