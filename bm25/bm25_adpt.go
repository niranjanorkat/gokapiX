package bm25

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/optimize"
)

func AdptQuery(query []string, bm25Model BM25AdptModel) BM25AdptModel {
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

			numeratorTerm := (bm25Model.TermK1[term] + 1) * float64(termFreqInDoc)
			denominatorTerm := bm25Model.TermK1[term]*(1-bm25Model.B+bm25Model.B*(float64(len(doc))/bm25Model.AvgDocLen)) + float64(termFreqInDoc)

			retrievalVal += bm25Model.G1Q[term] * (numeratorTerm / denominatorTerm)

		}
		bm25Model.TopScores = append(bm25Model.TopScores, retrievalVal)
		bm25Model.TopN = append(bm25Model.TopN, i)
	}

	sort.SliceStable(bm25Model.TopN, func(i, j int) bool {
		return bm25Model.TopScores[bm25Model.TopN[i]] > bm25Model.TopScores[bm25Model.TopN[j]]
	})

	return bm25Model
}

func computeTermK1Adpt(bm25AdptModel BM25AdptModel) BM25AdptModel {
	var terms []string
	for term := range bm25AdptModel.DocFreq {
		terms = append(terms, term)
	}

	corpusSize := len(bm25AdptModel.Corpus)

	for _, term := range terms {
		uniqueTermFreq := map[int]struct{}{0: {}}
		for _, doc := range bm25AdptModel.TermFreqInDoc {
			if freq, exists := doc[term]; exists {
				uniqueTermFreq[freq] = struct{}{}
			}
		}

		var uniqueTermFreqSlice []int
		for freq := range uniqueTermFreq {
			uniqueTermFreqSlice = append(uniqueTermFreqSlice, freq)
		}

		gqrs := make([]float64, 0)

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

			df_r1 := float64(dfr) + 0.5
			df_r := float64(dfr) + 1
			gqr := math.Log2(df_r1/df_r) - math.Log2((df_r1-0.5)/(float64(corpusSize)+1))

			gqrs = append(gqrs, gqr)

		}
		if len(gqrs) > 1 {
			bm25AdptModel.G1Q[term] = gqrs[1]
		} else {
			bm25AdptModel.G1Q[term] = 0
		}

		bm25AdptModel.G1Q[term] = gqrs[1]
		k1 := optimizeK1Adpt(gqrs, bm25AdptModel.G1Q[term], uniqueTermFreqSlice)
		bm25AdptModel.TermK1[term] = k1
	}
	return bm25AdptModel
}

func optimizeK1Adpt(GqR []float64, Gq float64, rValues []int, initGuess ...float64) float64 {
	k1 := 1.5
	if len(initGuess) > 0 {
		k1 = initGuess[0]
	}

	objFunc := func(x []float64) float64 {
		k1 := x[0]
		var sum float64
		for i, r := range rValues {
			bm25Term := ((k1 + 1) * float64(r)) / (k1 + float64(r))
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
