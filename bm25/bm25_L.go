package bm25

import (
	"math"
	"sort"
)

func BM25LCompute(query []string, corpus [][]string, b float64, k1 float64, delta float64) BM25Result {
	var docTermFrequency = getDocTermFreq(query, corpus)
	var totalDocs = len(corpus)
	var avgDocLength = getAvgDocLength(corpus)

	bm25Okapi := BM25Result{
		TopN:      make([]int, 0, len(corpus)),
		TopScores: make([]float64, 0, len(corpus)),
	}
	for i, doc := range corpus {
		var retrievalVal float64 = 0
		for j, term := range query {
			termFreqInDoc := CountOccurrences(term, doc)

			logTerm := math.Log((float64(totalDocs) + 1) / (float64(docTermFrequency[j]) + 0.5))
			ctd := float64(termFreqInDoc) / (1 - b + b*float64((len(doc)/avgDocLength)))
			numeratorTerm := (k1 + 1) * (ctd + delta)
			denominatorTerm := k1 + ctd + delta

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
