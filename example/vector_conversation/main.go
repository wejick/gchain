package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
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
	embeddingModel := _openai.NewOpenAIEmbedModel(OAIauthToken, "", openai.AdaEmbeddingV2)
	chatModel = _openai.NewOpenAIChatModel(OAIauthToken, "", _openai.GPT3Dot5Turbo0301)

	wvClient, err = weaviateVS.NewWeaviateVectorStore(wvhost, wvscheme, wvApiKey, embeddingModel, nil)
	if err != nil {
		log.Println(err.Error() + "can't connect to weaviate")
	}

	return
}

func main() {
	indexFlag := flag.Bool("index", false, "Specify the --index flag to run the index function")

	flag.Parse()

	if Init() != nil {
		return
	}

	if *indexFlag {
		Indexing()
	} else {
		Chatting([]model.ChatMessage{})
	}
}

func Chatting(memory []model.ChatMessage) {
	chain := conversational_retrieval.NewConversationalRetrievalChain(chatModel, memory, wvClient, textplitter, "", 1000)
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

func Indexing() (err error) {
	sources := []source{
		{filename: "indonesia.txt", url: "https://en.wikipedia.org/wiki/Indonesia"},
		{filename: "history_of_indonesia.txt", url: "https://en.wikipedia.org/wiki/History_of_Indonesia"},
	}
	for idx, s := range sources {
		data, err := ioutil.ReadFile(s.filename)
		if err != nil {
			log.Println(err)
			continue
		}
		sources[idx].doc = string(data)
	}

	var docs []string
	docs = textplitter.SplitText(sources[0].doc, 1000, 0)
	docs = append(docs, textplitter.SplitText(sources[1].doc, 1000, 0)...)

	bErr, err := wvClient.AddDocuments(context.Background(), "Indonesia", docs)
	if err != nil {
		log.Println(err)
		log.Println(bErr)
		return
	}

	return
}
