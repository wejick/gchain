package elasticsearchVS

import (
	"reflect"
	"testing"

	"github.com/wejick/gchain/document"
)

func Test_jsonStruct(t *testing.T) {
	type args struct {
		doc interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				doc: KNNSearchBody{
					KNN: KNNField{
						Field:         "vector",
						QueryVector:   []float32{},
						K:             1,
						NumCandidates: 1,
					},
					Fields: []string{"dense_vector"},
				},
			},
			want:    `{"knn":{"field":"vector","query_vector":[],"k":1,"num_candidates":1},"fields":["dense_vector"]}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonStruct(tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonStruct() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("jsonStruct() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_hitToDocs(t *testing.T) {
	type args struct {
		esRespBody       ElasticResponse
		additionalFields []string
	}
	tests := []struct {
		name     string
		args     args
		wantDocs []document.Document
		wantErr  bool
	}{
		{
			name: "success",
			args: args{
				esRespBody: ElasticResponse{
					Took:     100,
					TimedOut: false,
					Hits: struct {
						MaxScore float64      "json:\"max_score,omitempty\""
						Hits     []ElasticHit "json:\"hits,omitempty\""
					}{
						MaxScore: 1,
						Hits: []ElasticHit{
							{
								Index: "",
								ID:    "",
								Score: 1,
								Source_: []byte(`{
									"id": "1",
									"link": "https://seller.tokopedia.com/edu/seo-people-also-ask",
									"text": "lorem ipsum dolor sit amet"
								}`),
							},
						},
					},
				},
				additionalFields: []string{"id"},
			},
			wantDocs: []document.Document{
				{
					Text: "lorem ipsum dolor sit amet",
					Metadata: map[string]interface{}{
						"id": "1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error unmarshall invalid source_",
			args: args{
				esRespBody: ElasticResponse{
					Took:     100,
					TimedOut: false,
					Hits: struct {
						MaxScore float64      "json:\"max_score,omitempty\""
						Hits     []ElasticHit "json:\"hits,omitempty\""
					}{
						MaxScore: 1,
						Hits: []ElasticHit{
							{
								Index: "",
								ID:    "",
								Score: 1,
								Source_: []byte(`{
									"id": "1",
									"link": "https://seller.tokopedia.com/edu/seo-people-also-ask",
									"text": "lorem ipsum dolor sit amet"
								`),
							},
						},
					},
				},
				additionalFields: []string{"id"},
			},
			wantDocs: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDocs, err := hitToDocs(tt.args.esRespBody, tt.args.additionalFields)
			if (err != nil) != tt.wantErr {
				t.Errorf("hitToDocs() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDocs, tt.wantDocs) {
				t.Errorf("hitToDocs() = %+v, want %+v", gotDocs, tt.wantDocs)
			}
		})
	}
}

func Test_dataToESDoc(t *testing.T) {
	type args struct {
		documents []document.Document
		vector    [][]float32
	}
	tests := []struct {
		name       string
		args       args
		wantOutput []elasticDocument
	}{
		{
			name: "success",
			args: args{
				documents: []document.Document{
					{
						Text: "lorem ipsum dolor sit amet",
						Metadata: map[string]interface{}{
							"id": "1",
						},
					},
				},
				vector: [][]float32{
					{
						0.001,
					},
				},
			},
			wantOutput: []elasticDocument{
				map[string]interface{}{
					"id": "1",
					"text": "lorem ipsum dolor sit amet",
					"vector": []float32{0.001},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOutput := dataToESDoc(tt.args.documents, tt.args.vector); !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("dataToESDoc() = %+v, want %+v", gotOutput, tt.wantOutput)
			}
		})
	}
}
