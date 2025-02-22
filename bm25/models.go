package bm25

type BM25Model struct {
	Corpus        [][]string
	TermFreqInDoc []map[string]int
	DocFreq       map[string]int
	AvgDocLen     float64
	TopN          []int
	TopScores     []float64
}
