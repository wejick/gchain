package textsplitter

import (
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

type TikTokenSplitter struct {
	tkm *tiktoken.Tiktoken
}

// NewTikTokenSplitter create new TikTokenSplitter instance
// if modelName empty, the default one is gpt-3.5-turbo-0301
func NewTikTokenSplitter(modelName string) (*TikTokenSplitter, error) {
	if modelName == "" {
		modelName = "gpt-3.5-turbo-0301"
	}

	tkm, err := tiktoken.EncodingForModel(modelName)
	return &TikTokenSplitter{
		tkm: tkm,
	}, err
}

// splitIntoBatches creates word batches where length's doesn't exceed maxChunkSize.
func (W *TikTokenSplitter) SplitText(input string, maxChunkSize int) []string {
	batches := []string{}

	words := strings.Fields(input)
	var batch []string
	var lenCounter int

	for _, word := range words {
		if lenCounter+W.len(word) > maxChunkSize {
			batches = append(batches, strings.Join(batch, " "))
			batch = []string{}
			lenCounter = 0
		}

		batch = append(batch, word)
		lenCounter += W.len(word)
	}

	if len(batch) > 0 {
		batches = append(batches, strings.Join(batch, " "))
	}

	return batches
}

func (W *TikTokenSplitter) len(input string) int {
	return len(W.tkm.Encode(input, nil, nil))
}
