package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/wejick/gchain/callback"
	"github.com/wejick/gchain/chain/conversation"
	"github.com/wejick/gchain/model"
	_openai "github.com/wejick/gchain/model/openAI"

	"github.com/gorilla/websocket"
)

type message struct {
	Text     string `json:"text"`
	Finished bool   `json:"finished"`
}

var authToken = os.Getenv("OPENAI_API_KEY")
var chatModel *_openai.OpenAIChatModel

func main() {
	chatModel = _openai.NewOpenAIChatModel(authToken, _openai.GPT3Dot5Turbo0301, callback.NewManager())

	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	// websocket route
	http.HandleFunc("/chat", wshandler)

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	fmt.Println("Program exited.")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	log.Println("serving ws connection")
	memory := []model.ChatMessage{}
	streamingChannel := make(chan model.ChatMessage, 100)

	// setup new conversation
	convoChain := conversation.NewConversationChain(chatModel, memory, callback.NewManager(), "You're helpful chatbot that answer human question very concisely, response with formatted html.", false)
	convoChain.AppendToMemory(model.ChatMessage{Role: model.ChatMessageRoleAssistant, Content: "Hi, My name is GioAI"})

	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// send greetings
	m, err := json.Marshal(message{Text: "Hi, My name is GioAI", Finished: true})
	if err != nil {
		log.Println(err)
		return
	}
	ws.WriteMessage(websocket.TextMessage, m)

	for {
		// Read in a new requestMessage as JSON and map it to a Message object
		_, requestMessage, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		var output string

		// send request to model
		go func() {
			var err error
			output, err = convoChain.SimpleRun(context.Background(), string(requestMessage[:]), model.WithStreaming(true), model.WithStreamingChannel(streamingChannel))
			if err != nil {
				fmt.Println("error " + err.Error())
				return
			}
		}()

		// handle the response streaming
		for {
			value, ok := <-streamingChannel
			if ok && !model.IsStreamFinished(value) {
				m, err := json.Marshal(message{Text: value.Content, Finished: false})
				if err != nil {
					log.Println(err)
					continue
				}
				ws.WriteMessage(websocket.TextMessage, []byte(m))
			} else {
				m, err := json.Marshal(message{Finished: true})
				if err != nil {
					log.Println(err)
					continue
				}
				ws.WriteMessage(websocket.TextMessage, m)
				break
			}
		}

		// put use message and model response to conversation memory
		convoChain.AppendToMemory(model.ChatMessage{Role: model.ChatMessageRoleUser, Content: string(requestMessage[:])})
		convoChain.AppendToMemory(model.ChatMessage{Role: model.ChatMessageRoleAssistant, Content: output})
	}
}
