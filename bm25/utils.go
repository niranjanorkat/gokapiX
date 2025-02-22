package bm25

import "sort"

func SortTopResults(topN []int, topScores []float64) {
	sort.SliceStable(topN, func(i, j int) bool {
		return topScores[topN[i]] > topScores[topN[j]]
	})

	sortedScores := make([]float64, len(topScores))
	for i, id := range topN {
		sortedScores[i] = topScores[id]
	}
	copy(topScores, sortedScores)
}
