package bm25

func getDocTermFreq(query []string, corpus [][]string) []int {
	docTermFrequency := make([]int, len(query))

	for i, term := range query {
		for _, doc := range corpus {
			if Contains(term, doc) {
				docTermFrequency[i] += 1
			}
		}
	}

	return docTermFrequency
}

func getAvgDocLength(corpus [][]string) int {
	if len(corpus) == 0 {
		return 0
	}

	totalDocs := len(corpus)
	totalLength := 0
	for _, doc := range corpus {
		totalLength += len(doc)
	}

	return totalLength / totalDocs
}

func CountOccurrences(word string, document []string) int {
	count := 0
	for _, w := range document {
		if w == word {
			count++
		}
	}
	return count
}

func Contains(word string, document []string) bool {
	for _, w := range document {
		if w == word {
			return true
		}
	}
	return false
}
