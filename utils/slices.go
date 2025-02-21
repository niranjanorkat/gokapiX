package utils

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
