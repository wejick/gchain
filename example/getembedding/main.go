package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
	_openai "github.com/wejick/gchain/model/openAI"
)

type vector struct {
	Embedding []float32 `json:"embedding"`
}

func main() {
	// get the string from command line
	input := os.Args[1]
	// get from env variable openai ke
	OAIauthToken := os.Getenv("OPENAI_API_KEY")

	v := vector{}
	var err error

	// create new openai embedding model
	embeddingModel := _openai.NewOpenAIEmbedModel(OAIauthToken, "", "", openai.AdaEmbeddingV2)
	v.Embedding, err = embeddingModel.EmbedQuery(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	output, _ := json.Marshal(v)
	fmt.Println("input :" + input)
	fmt.Print(string(output[:]))
}
