/*
Language model has token limit, managing text/prompt to fit the limit is easier with text splitter.
Text splitter will split text into chunk of text with max size defined maxChunkSize.

There are 2 types of text splitter available :
1. Word splitter
Split the text word by word and make sure the chunk size is not exceed the maxChunkSize.
The maxChunkSize is in character.
2. Tiktoken splitter (I think it should be called tiktoken word splitter)
Split the text word by word and make sure the chunk size is not exceed the maxChunkSize.
The maxChunkSize is according to tiktoken definition of token.
*/
package textsplitter

import "github.com/wejick/gochain/document"

//go:generate moq -out textsplitter_moq.go . TextSplitter

// TextSplitter split text into chunk of text
type TextSplitter interface {
	SplitText(input string, maxChunkSize int, overlap int) []string
	SplitDocument(input document.Document, maxChunkSize int, overlap int) []document.Document
	Len(input string) int
}
