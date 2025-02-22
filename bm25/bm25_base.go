package bm25

func BM25BaseInit(corpus [][]string, b float64) BM25Base {
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

	return BM25Base{
		Corpus:        corpus,
		TermFreqInDoc: termFreqInDoc,
		DocFreq:       docFreq,
		AvgDocLen:     avgDocLen,
		TopN:          make([]int, 0, totalDocs),
		TopScores:     make([]float64, 0, totalDocs),
		B:             b,
	}
}

func BM25Init(corpus [][]string, b float64, k1 float64) BM25Model {
	bm25Base := BM25BaseInit(corpus, b)

	bm25Model := BM25Model{
		BM25Base: bm25Base,
		K1:       k1,
	}

	return bm25Model
}

func BM25LInit(corpus [][]string, b float64, k1 float64, delta float64) BM25LModel {
	bm25Base := BM25BaseInit(corpus, b)

	bm25LModel := BM25LModel{
		BM25Base: bm25Base,
		K1:       k1,
		Delta:    delta,
	}

	return bm25LModel
}

func BM25AdptInit(corpus [][]string, b float64, k1 float64) BM25AdptModel {
	bm25Base := BM25BaseInit(corpus, b)

	bm25AdptModel := BM25AdptModel{
		BM25Base: bm25Base,
		TermK1:   make(map[string]float64),
		G1Q:      make(map[string]float64),
	}

	bm25AdptModel = computeTermK1Adpt(bm25AdptModel)

	return bm25AdptModel
}

func BM25TInit(corpus [][]string, b float64, k1 float64) BM25TModel {
	bm25Base := BM25BaseInit(corpus, b)

	bm25TModel := BM25TModel{
		BM25Base: bm25Base,
		TermK1:   make(map[string]float64),
	}

	bm25TModel = computeTermK1(bm25TModel)

	return bm25TModel
}
