package textsplitter

import "github.com/wejick/gochain/document"

//go:generate moq -out textsplitter_moq.go . TextSplitter

// TextSplitter split text
type TextSplitter interface {
	SplitText(input string, maxChunkSize int, overlap int) []string
	SplitDocument(input document.Document, maxChunkSize int, overlap int) []document.Document
	Len(input string) int
}
