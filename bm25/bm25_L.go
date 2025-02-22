package bm25

import (
	"math"
	"sort"
)

func LQuery(query []string, bm25Model BM25LModel) BM25LModel {
	corpus := bm25Model.Corpus
	var totalDocs = len(corpus)

	bm25Model.TopN = make([]int, 0, totalDocs)
	bm25Model.TopScores = make([]float64, 0, totalDocs)

	for i, doc := range corpus {
		var retrievalVal float64 = 0
		for _, term := range query {
			if _, exists := bm25Model.TermFreqInDoc[i][term]; !exists {
				continue
			}
			termFreqInDoc := bm25Model.TermFreqInDoc[i][term]

			logTerm := math.Log((float64(totalDocs) + 1) / (float64(bm25Model.DocFreq[term]) + 0.5))
			ctd := float64(termFreqInDoc) / (1 - bm25Model.B + bm25Model.B*(float64(len(doc))/bm25Model.AvgDocLen))
			numeratorTerm := (bm25Model.K1 + 1) * (ctd + bm25Model.AvgDocLen)
			denominatorTerm := bm25Model.K1 + ctd + bm25Model.Delta

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
