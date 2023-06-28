package llm_chain

import (
	"context"

	"github.com/wejick/gchain/callback"
	"github.com/wejick/gchain/chain"
	"github.com/wejick/gchain/model"
	"github.com/wejick/gchain/prompt"
)

var _ chain.BaseChain = &LLMChain{}

const defaultTemplate = `{{.input}}`

type LLMChain struct {
	llmModel        model.LLMModel
	callbackManager *callback.Manager
	promptTemplate  *prompt.PromptTemplate
}

// NewLLMChain create an LLMChain instance
// if nil promptTemplate provided, the default one will be used
// default promptTemplate expect prompt["input"] as template key
func NewLLMChain(llmModel model.LLMModel, callbackManager *callback.Manager, promptTemplate *prompt.PromptTemplate, verbose bool) (llmchain *LLMChain, err error) {
	if promptTemplate == nil {
		promptTemplate, err = prompt.NewPromptTemplate("default", defaultTemplate)
		if err != nil {
			return
		}
	}
	if verbose {
		callbackManager.RegisterCallback(chain.CallbackChainEnd, callback.VerboseCallback)
	}
	llmchain = &LLMChain{
		llmModel:        llmModel,
		promptTemplate:  promptTemplate,
		callbackManager: callbackManager,
	}

	return
}

// Run do completion
// the default template expect prompt["input"] as input, and put the result to output["output"]
func (L *LLMChain) Run(ctx context.Context, prompt map[string]string, options ...func(*model.Option)) (output map[string]string, err error) {
	output = make(map[string]string)

	//trigger callback chain start
	L.callbackManager.TriggerEvent(ctx, chain.CallbackChainStart, callback.CallbackData{
		EventName:    chain.CallbackChainStart,
		FunctionName: "LLMChain.Run",
		Input:        prompt,
		Output:       output,
	})

	promptStr, err := L.promptTemplate.FormatPrompt(prompt)
	if err != nil {
		return
	}

	output["output"], err = L.llmModel.Call(ctx, promptStr, options...)

	// trigger callback chain end
	L.callbackManager.TriggerEvent(ctx, chain.CallbackChainEnd, callback.CallbackData{
		EventName:    chain.CallbackChainEnd,
		FunctionName: "LLMChain.Run",
		Input:        prompt,
		Output:       output,
	})

	return
}

// SimpleRun will run the prompt string agains llmchain
func (L *LLMChain) SimpleRun(ctx context.Context, prompt string, options ...func(*model.Option)) (output string, err error) {
	//trigger callback chain start
	L.callbackManager.TriggerEvent(ctx, chain.CallbackChainStart, callback.CallbackData{
		EventName:    chain.CallbackChainStart,
		FunctionName: "LLMChain.SimpleRun",
		Input:        map[string]string{"input": prompt},
	})

	output, err = L.llmModel.Call(ctx, prompt, options...)

	//trigger callback chain end
	L.callbackManager.TriggerEvent(ctx, chain.CallbackChainEnd, callback.CallbackData{
		EventName:    chain.CallbackChainEnd,
		FunctionName: "LLMChain.SimpleRun",
		Input:        map[string]string{"input": prompt},
		Output:       map[string]string{"output": output},
	})

	return
}
