package bm25

func BM25Init(corpus [][]string, b float64, k1 float64) BM25Model {
	totalDocs := len(corpus)
	docFreq := make(map[string]int)
	termFreqInDoc := make([]map[string]int, totalDocs)
	totalDocLen := 0

	for i, doc := range corpus {
		termFreqInDoc[i] = make(map[string]int)
		for _, term := range doc {
			termFreqInDoc[i][term]++
		}
		totalDocLen += len(doc)

		for term := range termFreqInDoc[i] {
			docFreq[term]++
		}
	}

	avgDocLen := float64(totalDocLen) / float64(totalDocs)

	return BM25Model{
		Corpus:        corpus,
		TermFreqInDoc: termFreqInDoc,
		DocFreq:       docFreq,
		AvgDocLen:     avgDocLen,
		TopN:          make([]int, 0, totalDocs),
		TopScores:     make([]float64, 0, totalDocs),
		B:             b,
		K1:            k1,
	}
}
