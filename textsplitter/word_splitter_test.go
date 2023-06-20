package textsplitter

import (
	"reflect"
	"testing"

	"github.com/wejick/gochain/document"
)

func TestWordSplitter_SplitText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "This is a simple test case",
			input: "This is a simple test case",
			want:  []string{"This is a", "simple", "test case"},
		},
		{
			name:  "One two three four five six seven eight nine ten",
			input: "One two three four five six seven eight nine ten",
			want:  []string{"One two", "three", "four five", "six seven", "eight", "nine ten"},
		},
		{
			name:  "An extremely long word",
			input: "An extremely long word is pneumonoultramicroscopicsilicovolcanoconiosis",
			want:  []string{"An", "extremely", "long word", "is", "pneumonoultramicroscopicsilicovolcanoconiosis"},
		},
		{
			name:  "empty",
			input: "",
			want:  []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			W := &WordSplitter{}
			got := W.SplitText(tt.input, 10, 0)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WordSplitter.SplitText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWordSplitter_SplitDocument(t *testing.T) {
	tests := []struct {
		name  string
		W     *WordSplitter
		input document.Document
		want  []document.Document
	}{
		{
			name:  "This is a simple test case",
			input: document.Document{Text: "This is a simple test case"},
			want:  []document.Document{{Text: "This is a"}, {Text: "simple"}, {Text: "test case"}},
		},
		{
			name:  "One two three four five six seven eight nine ten",
			input: document.Document{Text: "One two three four five six seven eight nine ten"},
			want:  []document.Document{{Text: "One two"}, {Text: "three"}, {Text: "four five"}, {Text: "six seven"}, {Text: "eight"}, {Text: "nine ten"}},
		},
		{
			name:  "An extremely long word",
			input: document.Document{Text: "An extremely long word is pneumonoultramicroscopicsilicovolcanoconiosis"},
			want:  []document.Document{{Text: "An"}, {Text: "extremely"}, {Text: "long word"}, {Text: "is"}, {Text: "pneumonoultramicroscopicsilicovolcanoconiosis"}},
		},
		{
			name:  "empty",
			input: document.Document{Text: ""},
			want:  []document.Document{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			W := &WordSplitter{}
			if got := W.SplitDocument(tt.input, 10, 0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WordSplitter.SplitDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}
