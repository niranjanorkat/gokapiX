package bm25

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/optimize"
)

func AdptQuery(query []string, bm25Model BM25Model) BM25Model {
	corpus := bm25Model.Corpus
	var totalDocs = len(corpus)

	bm25Model.TopN = make([]int, 0, totalDocs)
	bm25Model.TopScores = make([]float64, 0, totalDocs)

	for i, doc := range corpus {
		var retrievalVal float64 = 0
		for _, term := range query {
			termFreqInDoc := bm25Model.TermFreqInDoc[i][term]

			logTerm := math.Log(float64(totalDocs) / float64(bm25Model.DocFreq[term]))
			numeratorTerm := (bm25Model.K1 + 1) * float64(termFreqInDoc)
			denominatorTerm := bm25Model.K1*(1-bm25Model.B+bm25Model.B*(float64(len(doc))/bm25Model.AvgDocLen)) + float64(termFreqInDoc)

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

func BM25AdptInit(corpus [][]string, b float64, k1 float64) BM25AdptModel {
	bm25Model := BM25Init(corpus, b, k1)

	bm25AdptModel := BM25AdptModel{
		BM25Model: bm25Model,
		TermK1:    make(map[string]float64),
	}

	return bm25AdptModel
}

func ComputeTermK1(bm25AdptModel BM25AdptModel) BM25AdptModel {
	var terms []string
	for term := range bm25AdptModel.DocFreq {
		terms = append(terms, term)
	}

	corpusSize := len(bm25AdptModel.Corpus)

	for _, term := range terms {
		uniqueTermFreq := map[int]struct{}{}
		for _, doc := range bm25AdptModel.TermFreqInDoc {
			if freq, exists := doc[term]; exists {
				uniqueTermFreq[freq] = struct{}{}
			}
		}

		var uniqueTermFreqSlice []int
		for freq := range uniqueTermFreq {
			uniqueTermFreqSlice = append(uniqueTermFreqSlice, freq)
		}

		grq := make([]float64, 0)
		for _, numOccurrences := range uniqueTermFreqSlice {
			var dfr int

			if numOccurrences == 0 {
				dfr = corpusSize
			} else if numOccurrences == 1 {
				dfr = bm25AdptModel.DocFreq[term]
			} else if numOccurrences > 1 {
				nDT := 0
				for _, doc := range bm25AdptModel.TermFreqInDoc {
					if termFreq, exists := doc[term]; exists && termFreq > 0 {
						ctd := float64(termFreq) / (1 - bm25AdptModel.B + bm25AdptModel.B*(float64(len(doc))/bm25AdptModel.AvgDocLen))
						if ctd >= float64(numOccurrences) {
							nDT++
						}
					}
				}
				dfr = nDT
			}

			grq = append(grq, float64(dfr))
		}
	}
}

func OptimizeK1(GqR []float64, Gq float64, rValues []float64, initGuess ...float64) float64 {
	k1 := 1.5
	if len(initGuess) > 0 {
		k1 = initGuess[0]
	}

	objFunc := func(x []float64) float64 {
		k1 := x[0]
		var sum float64
		for i, r := range rValues {
			bm25Term := ((k1 + 1) * r) / (k1 + r)
			infoGainRatio := GqR[i] / Gq
			diff := infoGainRatio - bm25Term
			sum += diff * diff
		}
		return sum
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
