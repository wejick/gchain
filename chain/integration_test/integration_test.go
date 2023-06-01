//go:build integration
// +build integration

package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wejick/gochain/chain/llm_chain"
	"github.com/wejick/gochain/chain/summarization"
	"github.com/wejick/gochain/model"
	_openai "github.com/wejick/gochain/model/openAI"
)

var llmModel model.LLMModel
var chatModel model.LLMModel

func TestMain(m *testing.M) {
	fmt.Println("Running integration tests...")
	// Perform any setup or initialization here

	var authToken = os.Getenv("OPENAI_API_KEY")
	llmModel = _openai.NewOpenAIModel(authToken, "", "text-davinci-003")

	chatModel = _openai.NewOpenAIChatModel(authToken, "", _openai.GPT3Dot5Turbo0301)

	exitCode := m.Run()

	// Perform any cleanup or teardown here

	// Exit with the appropriate exit code
	// (0 for success, non-zero for failure)
	os.Exit(exitCode)
}

func TestLlmChain(t *testing.T) {
	chain := llm_chain.NewLLMChain(llmModel)
	outputMap, err := chain.Run(context.Background(), map[string]string{"input": "Indonesia Capital is Jakarta\nJakarta is the capital of "})
	assert.NoError(t, err, "error Run")
	assert.Contains(t, outputMap["output"], "Indonesia", "unexpected result")
}

func TestStuffSummarizationChain(t *testing.T) {
	llmchain := llm_chain.NewLLMChain(llmModel)
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
	llmchain := llm_chain.NewLLMChain(llmModel)
	chain, err := summarization.NewMapReduceSummarizationChain(llmchain, "", "", "text", 1000)
	assert.NoError(t, err, "error NewMapReduceSummarizationChain")

	testDoc := make(map[string]string)
	testDoc["input"] = `Modular audio and video hardware for retro machines like the Commodore 64. Designed to use 74 series TTL through hole ICs available back in the 1980s, something you can solder at home from parts or order ready assembled.
	One of the most recent videos shows a "Shadow of the Beast" demonstration, to show parallax scrolling with precisely timed raster effects. Please do consider subscribing to the YouTube channel if you want to see more updates to this project: 
	This project started when old retro arcade hardware was being discussed. In the back of my mind was the often fabled "Mega games" by Imagine Software which were planned to use extra hardware on the Spectrum and Commodore 64 to augment the machine's capabilities. Since this hardware uses TTL logic available back from the same time period I was wondering exactly how much extra graphical grunt could have been engineered and interfaced with these old 8-bit computers.
	Truth be told, the Imagine hardware was pretty much just extra RAM https://www.gamesthatwerent.com/gtw64/mega-games/ but this was a fun project to see how far the arcade hardware was pushing the limits of board size and signal complexity.
	I was looking at Bomb Jack boards on ebay and pondering how they had enough fill-rate to draw 24 16x16 sprites and have the option for some to use 32x32 mode as well. A friend and I were discussing the clock speed and fill-rate while trying to deduce the operation of the hardware just by inspecting the hand drawn schematics, as you do.
	In the end to get some clarity on the sprite plotting specifically I started to transcribe what was thought to be the sprite logic portion of the schematic into Proteus, since it can simulate digital electronics really well.`

	output, err := chain.Run(context.Background(), testDoc, model.MaxToken(200))
	assert.NoError(t, err, "error Run(context.Background(), testDoc, model.MaxToken(200))")

	t.Log(output)

}

func TestStuffSummarizationChainChat(t *testing.T) {
	llmchain := llm_chain.NewLLMChain(chatModel)
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
