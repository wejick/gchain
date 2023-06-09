package textsplitter

import (
	"reflect"
	"testing"
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
