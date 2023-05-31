package textsplitter

import "strings"

type WordSplitter struct {
	maxToken int
}

// splitIntoBatches creates word batches where length's doesn't exceed maxToken.
func (W *WordSplitter) SplitText(input string) []string {
	batches := []string{}

	words := strings.Fields(input)
	var batch []string
	var lenCounter int

	for _, word := range words {
		// +1 is for a possible space character
		if lenCounter+len(word)+1 > W.maxToken {
			batches = append(batches, strings.Join(batch, " "))
			batch = []string{}
			lenCounter = 0
		}

		batch = append(batch, word)
		lenCounter += len(word) + 1
	}

	if len(batch) > 0 {
		batches = append(batches, strings.Join(batch, " "))
	}

	return batches
}
