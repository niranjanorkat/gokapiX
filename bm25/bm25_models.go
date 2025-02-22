package bm25

type BM25Base struct {
	Corpus        [][]string
	TermFreqInDoc []map[string]int
	DocFreq       map[string]int
	AvgDocLen     float64
	TopN          []int
	TopScores     []float64
	B             float64
}

type BM25Model struct {
	BM25Base
	K1 float64
}

type BM25LModel struct {
	BM25Base
	K1    float64
	Delta float64
}

type BM25AdptModel struct {
	BM25Base
	G1Q    map[string]float64
	TermK1 map[string]float64
}

type BM25TModel struct {
	BM25Base
	TermK1 map[string]float64
}
