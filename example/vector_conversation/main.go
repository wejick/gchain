package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/wejick/gochain/callback"
	"github.com/wejick/gochain/chain/conversational_retrieval"
	weaviateVS "github.com/wejick/gochain/datastore/weaviate_vector"
	"github.com/wejick/gochain/model"
	_openai "github.com/wejick/gochain/model/openAI"
	"github.com/wejick/gochain/textsplitter"
)

var OAIauthToken = os.Getenv("OPENAI_API_KEY")
var chatModel *_openai.OpenAIChatModel
var embeddingModel *_openai.OpenAIEmbedModel
var wvClient *weaviateVS.WeaviateVectorStore
var textplitter *textsplitter.TikTokenSplitter

const (
	wvhost   = "question-testing-twjfnqyp.weaviate.network"
	wvscheme = "https"
	wvApiKey = ""
)

type source struct {
	filename string
	url      string
	doc      string
}

func Init() (err error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	textplitter, err = textsplitter.NewTikTokenSplitter(_openai.GPT3Dot5Turbo0301)
	if err != nil {
		log.Println(err)
		return
	}
	embeddingModel = _openai.NewOpenAIEmbedModel(OAIauthToken, "", openai.AdaEmbeddingV2)
	chatModel = _openai.NewOpenAIChatModel(OAIauthToken, "", _openai.GPT3Dot5Turbo0301, callback.NewManager(), false)

	wvClient, err = weaviateVS.NewWeaviateVectorStore(wvhost, wvscheme, wvApiKey, embeddingModel, nil)
	if err != nil {
		log.Println(err.Error() + "can't connect to weaviate")
	}

	return
}

func main() {
	indexFlag := flag.Bool("index", false, "Specify the --index flag to run the index function")
	deleteFlag := flag.Bool("delete", false, "Specify the --delete flag to run the delete function")

	flag.Parse()

	if Init() != nil {
		return
	}

	if *indexFlag {
		Indexing()
	} else if *deleteFlag {
		DeleteIndex()
	} else {
		Chatting([]model.ChatMessage{})
	}
}

func Chatting(memory []model.ChatMessage) {
	chain := conversational_retrieval.NewConversationalRetrievalChain(chatModel, memory, wvClient, "Indonesia", textplitter, callback.NewManager(), "", 2000, false)
	fmt.Println("AI : How can I help you, I know many things about indonesia")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("User : ")
		scanner.Scan()
		input := scanner.Text()
		output, err := chain.Run(context.Background(), map[string]string{"input": input})
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("AI : " + output["output"])
	}
}

func DeleteIndex() (err error) {
	err = wvClient.DeleteIndex(context.Background(), "Indonesia")
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func Indexing() (err error) {
	indexingplitter, err := textsplitter.NewTikTokenSplitter(openai.AdaEmbeddingV2.String())
	if err != nil {
		log.Println(err)
		return
	}
	sources := []source{
		{filename: "indonesia.txt", url: "https://en.wikipedia.org/wiki/Indonesia"},
		{filename: "history_of_indonesia.txt", url: "https://en.wikipedia.org/wiki/History_of_Indonesia"},
	}
	for idx, s := range sources {
		data, err := os.ReadFile(s.filename)
		if err != nil {
			log.Println(err)
			continue
		}
		sources[idx].doc = string(data)
	}

	var docs []string
	docs = indexingplitter.SplitText(sources[0].doc, 500, 0)
	docs = append(docs, indexingplitter.SplitText(sources[1].doc, 500, 0)...)

	bErr, err := wvClient.AddDocuments(context.Background(), "Indonesia", docs)
	if err != nil {
		log.Println(err)
		log.Println(bErr)
		return
	}

	return
}
