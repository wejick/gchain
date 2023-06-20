//go:build integration
// +build integration

package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wejick/gochain/callback"
	"github.com/wejick/gochain/chain/conversation"
	"github.com/wejick/gochain/chain/conversational_retrieval"
	"github.com/wejick/gochain/chain/llm_chain"
	"github.com/wejick/gochain/chain/summarization"
	wikipedia "github.com/wejick/gochain/datastore/wikipedia_retriever"
	"github.com/wejick/gochain/model"
	_openai "github.com/wejick/gochain/model/openAI"
	"github.com/wejick/gochain/prompt"
	"github.com/wejick/gochain/textsplitter"
)

var llmModel model.LLMModel
var chatModel model.ChatModel

func TestMain(m *testing.M) {
	fmt.Println("Running integration tests...")
	// Perform any setup or initialization here

	var authToken = os.Getenv("OPENAI_API_KEY")
	llmModel = _openai.NewOpenAIModel(authToken, "", "text-davinci-003", callback.NewManager(), false)

	chatModel = _openai.NewOpenAIChatModel(authToken, "", _openai.GPT3Dot5Turbo0301, callback.NewManager(), false)

	exitCode := m.Run()

	// Perform any cleanup or teardown here

	// Exit with the appropriate exit code
	// (0 for success, non-zero for failure)
	os.Exit(exitCode)
}

func TestLlmChain(t *testing.T) {
	chain, err := llm_chain.NewLLMChain(llmModel, callback.NewManager(), nil, false)
	assert.NoError(t, err, "NewLLMChain")
	outputMap, err := chain.Run(context.Background(), map[string]string{"input": "Indonesia Capital is Jakarta\nJakarta is the capital of "})
	assert.NoError(t, err, "error Run")
	assert.Contains(t, outputMap["output"], "Indonesia", "unexpected result")

	customPrompt, err := prompt.NewPromptTemplate("customPrompt", "{{.text}}")
	customPromptChain, err := llm_chain.NewLLMChain(llmModel, callback.NewManager(), customPrompt, false)
	assert.NoError(t, err, "NewLLMChain")

	customOutputMap, err := customPromptChain.Run(context.Background(), map[string]string{"text": "Indonesia Capital is Jakarta\nJakarta is the capital of "})
	assert.NoError(t, err, "error Run")
	assert.Contains(t, customOutputMap["output"], "Indonesia", "unexpected result")
}

func TestStuffSummarizationChain(t *testing.T) {
	llmchain, err := llm_chain.NewLLMChain(llmModel, callback.NewManager(), nil, false)
	assert.NoError(t, err, "NewLLMChain")

	chain, err := summarization.NewStuffSummarizationChain(llmchain, "", "text")
	assert.NoError(t, err, "error NewStuffSummarizationChain")
	output, err := chain.SimpleRun(context.Background(), `Modular audio and video hardware for retro machines like the Commodore 64. Designed to use 74 series TTL through hole ICs available back in the 1980s, something you can solder at home from parts or order ready assembled.
	One of the most recent videos shows a "Shadow of the Beast" demonstration, to show parallax scrolling with precisely timed raster effects. Please do consider subscribing to the YouTube channel if you want to see more updates to this project: 
	This project started when old retro arcade hardware was being discussed. In the back of my mind was the often fabled "Mega games" by Imagine Software which were planned to use extra hardware on the Spectrum and Commodore 64 to augment the machine's capabilities. Since this hardware uses TTL logic available back from the same time period I was wondering exactly how much extra graphical grunt could have been engineered and interfaced with these old 8-bit computers.
	Truth be told, the Imagine hardware was pretty much just extra RAM https://www.gamesthatwerent.com/gtw64/mega-games/ but this was a fun project to see how far the arcade hardware was pushing the limits of board size and signal complexity.
	I was looking at Bomb Jack boards on ebay and pondering how they had enough fill-rate to draw 24 16x16 sprites and have the option for some to use 32x32 mode as well. A friend and I were discussing the clock speed and fill-rate while trying to deduce the operation of the hardware just by inspecting the hand drawn schematics, as you do.
	In the end to get some clarity on the sprite plotting specifically I started to transcribe what was thought to be the sprite logic portion of the schematic into Proteus, since it can simulate digital electronics really well.`)

	assert.NoError(t, err)
	t.Log(output)

}

func TestMapReduceSummarizationChain(t *testing.T) {
	llmchain, err := llm_chain.NewLLMChain(llmModel, callback.NewManager(), nil, false)
	assert.NoError(t, err, "NewLLMChain")

	splitter, err := textsplitter.NewTikTokenSplitter("")
	assert.NoError(t, err, "NewTikTokenSplitter")
	chain, err := summarization.NewMapReduceSummarizationChain(llmchain, "", "", "text", splitter, 1000)
	assert.NoError(t, err, "error NewMapReduceSummarizationChain")

	testDoc := make(map[string]string)
	testDoc["input"] = `Modular audio and video hardware for retro machines like the Commodore 64. Designed to use 74 series TTL through hole ICs available back in the 1980s, something you can solder at home from parts or order ready assembled.
	One of the most recent videos shows a "Shadow of the Beast" demonstration, to show parallax scrolling with precisely timed raster effects. Please do consider subscribing to the YouTube channel if you want to see more updates to this project: 
	This project started when old retro arcade hardware was being discussed. In the back of my mind was the often fabled "Mega games" by Imagine Software which were planned to use extra hardware on the Spectrum and Commodore 64 to augment the machine's capabilities. Since this hardware uses TTL logic available back from the same time period I was wondering exactly how much extra graphical grunt could have been engineered and interfaced with these old 8-bit computers.
	Truth be told, the Imagine hardware was pretty much just extra RAM https://www.gamesthatwerent.com/gtw64/mega-games/ but this was a fun project to see how far the arcade hardware was pushing the limits of board size and signal complexity.
	I was looking at Bomb Jack boards on ebay and pondering how they had enough fill-rate to draw 24 16x16 sprites and have the option for some to use 32x32 mode as well. A friend and I were discussing the clock speed and fill-rate while trying to deduce the operation of the hardware just by inspecting the hand drawn schematics, as you do.
	In the end to get some clarity on the sprite plotting specifically I started to transcribe what was thought to be the sprite logic portion of the schematic into Proteus, since it can simulate digital electronics really well.`

	output, err := chain.Run(context.Background(), testDoc, model.WithMaxToken(200))
	assert.NoError(t, err, "error Run(context.Background(), testDoc, model.MaxToken(200))")

	t.Log(output)

}

func TestStuffSummarizationChainChat(t *testing.T) {
	llmchain, err := llm_chain.NewLLMChain(llmModel, callback.NewManager(), nil, false)
	assert.NoError(t, err, "NewLLMChain")

	chain, err := summarization.NewStuffSummarizationChain(llmchain, "", "text")
	assert.NoError(t, err, "error NewStuffSummarizationChain")
	output, err := chain.SimpleRun(context.Background(), `Modular audio and video hardware for retro machines like the Commodore 64. Designed to use 74 series TTL through hole ICs available back in the 1980s, something you can solder at home from parts or order ready assembled.
	One of the most recent videos shows a "Shadow of the Beast" demonstration, to show parallax scrolling with precisely timed raster effects. Please do consider subscribing to the YouTube channel if you want to see more updates to this project: 
	This project started when old retro arcade hardware was being discussed. In the back of my mind was the often fabled "Mega games" by Imagine Software which were planned to use extra hardware on the Spectrum and Commodore 64 to augment the machine's capabilities. Since this hardware uses TTL logic available back from the same time period I was wondering exactly how much extra graphical grunt could have been engineered and interfaced with these old 8-bit computers.
	Truth be told, the Imagine hardware was pretty much just extra RAM https://www.gamesthatwerent.com/gtw64/mega-games/ but this was a fun project to see how far the arcade hardware was pushing the limits of board size and signal complexity.
	I was looking at Bomb Jack boards on ebay and pondering how they had enough fill-rate to draw 24 16x16 sprites and have the option for some to use 32x32 mode as well. A friend and I were discussing the clock speed and fill-rate while trying to deduce the operation of the hardware just by inspecting the hand drawn schematics, as you do.
	In the end to get some clarity on the sprite plotting specifically I started to transcribe what was thought to be the sprite logic portion of the schematic into Proteus, since it can simulate digital electronics really well.`)

	assert.NoError(t, err)
	t.Log(output)
}

func TestConversationChainChat(t *testing.T) {
	memory := []model.ChatMessage{}
	convoChain := conversation.NewConversationChain(chatModel, memory, callback.NewManager(), "You're helpful chatbot that answer very concisely", false)

	convoChain.AppendToMemory(model.ChatMessage{Role: model.ChatMessageRoleAssistant, Content: "Hi, My name is GioAI"})
	output, err := convoChain.Run(context.Background(), map[string]string{"input": "what's your name?"}, model.WithTemperature(0), model.WithMaxToken(100))
	assert.NoError(t, err)

	outputString, err := convoChain.SimpleRun(context.Background(), "so your name is gioAI", model.WithTemperature(0), model.WithMaxToken(100))

	t.Log(output["output"])
	t.Log(outputString)
}

func TestConversationalRetrievalChainChat(t *testing.T) {
	memory := []model.ChatMessage{}
	splitter, err := textsplitter.NewTikTokenSplitter(_openai.GPT3Dot5Turbo0301)
	assert.NoError(t, err)
	convoChain := conversational_retrieval.NewConversationalRetrievalChain(chatModel, memory, &wikipedia.Wikipedia{}, "", splitter, callback.NewManager(), "You're helpful chatbot that answer very concisely", 1000, false)

	convoChain.AppendToMemory(model.ChatMessage{Role: model.ChatMessageRoleAssistant, Content: "Hi, My name is GioAI"})
	convoChain.AppendToMemory(model.ChatMessage{Role: model.ChatMessageRoleUser, Content: "Who is the first president of Indonesia?"})
	convoChain.AppendToMemory(model.ChatMessage{Role: model.ChatMessageRoleAssistant, Content: "The first president of indonesia was Soekarno"})

	response, err := convoChain.Run(context.Background(), map[string]string{"input": "tell me little bit more about soekarno?"}, model.WithTemperature(0.3), model.WithMaxToken(1000))
	t.Log(response)
	assert.NoError(t, err)
}
