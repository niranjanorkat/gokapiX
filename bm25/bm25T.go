package bm25

import (
	"fmt"
	"math"

	"github.com/niranjanorkat/gokapiX/helper"
	"gonum.org/v1/gonum/optimize"
)

func TQuery(query []string, bm25Model BM25TModel) BM25TModel {
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
			numeratorTerm := (bm25Model.TermK1[term] + 1) * float64(termFreqInDoc)
			denominatorTerm := bm25Model.TermK1[term]*(1-bm25Model.B+bm25Model.B*(float64(len(doc))/bm25Model.AvgDocLen)) + float64(termFreqInDoc)

			retrievalVal += logTerm * (numeratorTerm / denominatorTerm)

		}
		bm25Model.TopScores = append(bm25Model.TopScores, retrievalVal)
		bm25Model.TopN = append(bm25Model.TopN, i)
	}

	helper.SortTopResults(bm25Model.TopN, bm25Model.TopScores)

	return bm25Model
}

func computeTermK1T(bm25TModel *BM25TModel) {
	var terms []string
	for term := range bm25TModel.DocFreq {
		terms = append(terms, term)
	}

	for _, term := range terms {
		ctdValues := []float64{}
		for i, doc := range bm25TModel.TermFreqInDoc {
			termFreq := bm25TModel.TermFreqInDoc[i][term]
			if termFreq == 0 {
				continue

			}
			ctd := float64(termFreq) / (1 - bm25TModel.B + bm25TModel.B*(float64(len(doc))/bm25TModel.AvgDocLen))
			ctdValues = append(ctdValues, ctd)
		}

		k1 := optimizeK1T(bm25TModel.DocFreq[term], ctdValues)
		bm25TModel.TermK1[term] = k1
	}
}

func optimizeK1T(docFreq int, ctdValues []float64, initGuess ...float64) float64 {
	k1 := 1.5
	if len(initGuess) > 0 {
		k1 = initGuess[0]
	}

	var sumLogCTD float64
	for _, ctd := range ctdValues {
		sumLogCTD += math.Log(ctd)
	}
	target := (sumLogCTD + 1) / float64(docFreq)

	objFunc := func(x []float64) float64 {
		k1 := x[0]
		gk1 := computeGk1(k1)
		diff := gk1 - target
		return diff * diff
	}

	problem := optimize.Problem{
		Func: objFunc,
	}

	result, err := optimize.Minimize(problem, []float64{k1}, nil, &optimize.NelderMead{})
	if err != nil {
		fmt.Println("Optimization failed:", err)
		return k1
	}

	return result.X[0]
}

func computeGk1(k1 float64) float64 {
	if k1 == 1 {
		return 1
	}
	return (k1 / (k1 - 1)) * math.Log(k1)
}
