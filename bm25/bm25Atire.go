package bm25

import (
	"math"

	"github.com/niranjanorkat/gokapiX/helper"
)

func AtireQuery(query []string, bm25Model BM25Model) BM25Model {
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

			logTerm := math.Log(float64(totalDocs) / float64(bm25Model.DocFreq[term]))
			numeratorTerm := (bm25Model.K1 + 1) * float64(termFreqInDoc)
			denominatorTerm := bm25Model.K1*(1-bm25Model.B+bm25Model.B*(float64(len(doc))/bm25Model.AvgDocLen)) + float64(termFreqInDoc)

			retrievalVal += logTerm * (numeratorTerm / denominatorTerm)

		}
		bm25Model.TopScores = append(bm25Model.TopScores, retrievalVal)
		bm25Model.TopN = append(bm25Model.TopN, i)
	}

	helper.SortTopResults(bm25Model.TopN, bm25Model.TopScores)

	return bm25Model
}
