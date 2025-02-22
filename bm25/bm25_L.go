package bm25

import (
	"math"
	"sort"
)

func LQuery(query []string, bm25Model BM25Model, b float64, k1 float64, delta float64) BM25Model {
	corpus := bm25Model.Corpus
	var totalDocs = len(corpus)

	bm25Model.TopN = make([]int, 0, totalDocs)
	bm25Model.TopScores = make([]float64, 0, totalDocs)

	for i, doc := range corpus {
		var retrievalVal float64 = 0
		for _, term := range query {
			termFreqInDoc := bm25Model.TermFreqInDoc[i][term]

			logTerm := math.Log((float64(totalDocs) + 1) / (float64(bm25Model.DocFreq[term]) + 0.5))
			ctd := float64(termFreqInDoc) / (1 - b + b*(float64(len(doc))/bm25Model.AvgDocLen))
			numeratorTerm := (k1 + 1) * (ctd + delta)
			denominatorTerm := k1 + ctd + delta

			retrievalVal += logTerm * (numeratorTerm / denominatorTerm)

		}
		bm25Model.TopScores = append(bm25Model.TopScores, retrievalVal)
		bm25Model.TopN = append(bm25Model.TopN, i)
	}

	sort.SliceStable(bm25Model.TopN, func(i, j int) bool {
		return bm25Model.TopScores[bm25Model.TopN[i]] > bm25Model.TopScores[bm25Model.TopN[j]]
	})

	return bm25Model
}
