package textsplitter

import (
	"reflect"
	"testing"

	"github.com/wejick/gochain/document"
)

func TestTikTokenSplitter_SplitText(t *testing.T) {
	tkm, _ := NewTikTokenSplitter("")
	type args struct {
		input        string
		maxChunkSize int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "This is a simple test case",
			args: args{
				input:        "This is a simple test case",
				maxChunkSize: 10,
			},
			want: []string{"This is a simple test case"},
		},
		{
			name: "One two three four five six seven eight nine ten",
			args: args{
				input:        "One two three four five six seven eight nine ten",
				maxChunkSize: 10,
			},
			want: []string{"One two three four five six seven eight nine ten"},
		},
		{
			name: "An extremely long word",
			args: args{
				input:        "An extremely long word is pneumonoultramicroscopicsilicovolcanoconiosis",
				maxChunkSize: 10,
			},
			want: []string{"An extremely long word is", "pneumonoultramicroscopicsilicovolcanoconiosis"},
		},
		{
			name: "An extremely long word",
			args: args{
				input:        "An extremely long word is pneumonoultramicroscopicsilicovolcanoconiosis",
				maxChunkSize: 5,
			},
			want: []string{"An extremely long word", "is", "pneumonoultramicroscopicsilicovolcanoconiosis"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			W := tkm
			if got := W.SplitText(tt.args.input, tt.args.maxChunkSize, 0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TikTokenSplitter.SplitText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTikTokenSplitter_SplitDocument(t *testing.T) {
	tkm, _ := NewTikTokenSplitter("")

	type args struct {
		input        document.Document
		maxChunkSize int
		overlap      int
	}
	tests := []struct {
		name string
		args args
		want []document.Document
	}{
		{
			name: "This is a simple test case",
			args: args{
				input:        document.Document{Text: "This is a simple test case"},
				maxChunkSize: 10,
			},
			want: []document.Document{{Text: "This is a simple test case"}},
		},
		{
			name: "One two three four five six seven eight nine ten",
			args: args{
				input:        document.Document{Text: "One two three four five six seven eight nine ten"},
				maxChunkSize: 10,
			},
			want: []document.Document{{Text: "One two three four five six seven eight nine ten"}},
		},
		{
			name: "An extremely long word",
			args: args{
				input:        document.Document{Text: "An extremely long word is pneumonoultramicroscopicsilicovolcanoconiosis"},
				maxChunkSize: 10,
			},
			want: []document.Document{{Text: "An extremely long word is"}, {Text: "pneumonoultramicroscopicsilicovolcanoconiosis"}},
		},
		{
			name: "An extremely long word",
			args: args{
				input:        document.Document{Text: "An extremely long word is pneumonoultramicroscopicsilicovolcanoconiosis"},
				maxChunkSize: 5,
			},
			want: []document.Document{{Text: "An extremely long word"}, {Text: "is"}, {Text: "pneumonoultramicroscopicsilicovolcanoconiosis"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			T := tkm
			if got := T.SplitDocument(tt.args.input, tt.args.maxChunkSize, tt.args.overlap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TikTokenSplitter.SplitDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}
