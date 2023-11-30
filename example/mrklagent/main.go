package main

import (
	"context"
	"fmt"
	"os"

	"github.com/wejick/gchain/agent"
	"github.com/wejick/gchain/agent/mrkl"
	"github.com/wejick/gchain/callback"
	"github.com/wejick/gchain/chain/llm_chain"
	_openai "github.com/wejick/gchain/model/openAI"
	"github.com/wejick/gchain/tools/greeting"
	"github.com/wejick/gchain/tools/wikipedia"
)

func main() {
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Println("Type .quit to exit")

	var authToken = os.Getenv("OPENAI_API_KEY")
	chatModel := _openai.NewOpenAIChatModel(authToken, _openai.GPT3Dot5Turbo0301, callback.NewManager())

	llmChain, err := llm_chain.NewLLMChain(chatModel, callback.NewManager(), nil, true)
	if err != nil {
		panic(err)
	}

	mrklAgent, err := mrkl.NewMRKLAgent(llmChain)
	agent := agent.NewExecutor(mrklAgent, 10)
	agent.RegisterTool(wikipedia.NewWikipediaSearchTool())
	agent.RegisterTool(greeting.NewGreetingTool())

	output, err := agent.Run(context.Background(), map[string]string{"input": "Hi, how are you?"})

	if err != nil {
		panic(err)
	}

	fmt.Println(output)

	fmt.Println("Program exited.")
}
