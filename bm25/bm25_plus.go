package bm25

import (
	"math"
	"sort"
)

func BM25PlusCompute(query []string, corpus [][]string, b float64, k1 float64, delta float64) BM25Result {
	var docTermFrequency = getDocTermFreq(query, corpus)
	var totalDocs = len(corpus)
	var avgDocLength = getAvgDocLength(corpus)

	bm25Result := BM25Result{
		TopN:      make([]int, 0, len(corpus)),
		TopScores: make([]float64, 0, len(corpus)),
	}
	for i, doc := range corpus {
		var retrievalVal float64 = 0
		for j, term := range query {
			termFreqInDoc := CountOccurrences(term, doc)

			logTerm := math.Log((float64(totalDocs) + 1) / float64(docTermFrequency[j]))
			numeratorTerm := (k1 + 1) * float64(termFreqInDoc)
			denominatorTerm := k1*(1-b+b*float64((len(doc)/avgDocLength))) + float64(termFreqInDoc)

			retrievalVal += logTerm * ((numeratorTerm / denominatorTerm) + delta)

		}
		bm25Result.TopScores = append(bm25Result.TopScores, retrievalVal)
		bm25Result.TopN = append(bm25Result.TopN, i)
	}

	sort.SliceStable(bm25Result.TopN, func(i, j int) bool {
		return bm25Result.TopScores[bm25Result.TopN[i]] > bm25Result.TopScores[bm25Result.TopN[j]]
	})

	return bm25Result
}
