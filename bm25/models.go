package bm25

type BM25Model struct {
	Corpus        [][]string
	TermFreqInDoc []map[string]int
	DocFreq       map[string]int
	AvgDocLen     float64
	TopN          []int
	TopScores     []float64
	B             float64
	K1            float64
}

type BM25LModel struct {
	BM25Model
	Delta float64
}

type BM25AdptModel struct {
	BM25Model
	TermK1 map[string]float64
}
