package bm25

import (
	"math"
	"sort"

	"github.com/niranjankrishna-acad/gokapiX/utils"
)

type BM25Okapi struct {
	TopN      []int
	TopScores []float64
}

func BM25OkapiCompute(query []string, corpus [][]string, b float64, k1 float64) BM25Okapi {
	// Assuming query and corpus is received tokenized.

	docTermFrequency := make([]int, len(query))
	for i, term := range query {
		for _, doc := range corpus {
			termInDoc := utils.Contains(term, doc)
			if termInDoc {
				docTermFrequency[i] += 1
			}
		}
	}

	var totalDocs = len(corpus)
	var avgDocLength int
	for _, doc := range corpus {
		avgDocLength += len(doc)
	}
	avgDocLength /= totalDocs

	bm25Okapi := BM25Okapi{
		TopN:      make([]int, 0, len(corpus)),
		TopScores: make([]float64, 0, len(corpus)),
	}
	for i, doc := range corpus {
		var retrievalVal float64 = 0
		for j, term := range query {
			termFreqInDoc := utils.CountOccurrences(term, doc)

			logTerm := math.Log(float64(totalDocs) / float64(docTermFrequency[j]))
			numeratorTerm := (k1 + 1) * float64(termFreqInDoc)
			denominatorTerm := k1*(1-b+b*float64((len(doc)/avgDocLength))) + float64(termFreqInDoc)

			retrievalVal += logTerm * (numeratorTerm / denominatorTerm)

		}
		bm25Okapi.TopScores = append(bm25Okapi.TopScores, retrievalVal)
		bm25Okapi.TopN = append(bm25Okapi.TopN, i)
	}

	sort.SliceStable(bm25Okapi.TopN, func(i, j int) bool {
		return bm25Okapi.TopScores[bm25Okapi.TopN[i]] > bm25Okapi.TopScores[bm25Okapi.TopN[j]]
	})

	return bm25Okapi
}
