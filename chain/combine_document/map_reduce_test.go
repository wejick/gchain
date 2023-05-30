package combine_document

import (
	"reflect"
	"testing"
)

func Test_splitIntoBatches(t *testing.T) {
	type args struct {
		input    string
		maxToken int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitIntoBatches(tt.args.input, tt.args.maxToken); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitIntoBatches() = %v, want %v", got, tt.want)
			}
		})
	}
}
