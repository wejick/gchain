package conversation

import (
	"context"
	"errors"

	"github.com/wejick/gchain/callback"
	basechain "github.com/wejick/gchain/chain"
	"github.com/wejick/gchain/model"
)

// ConversationChain that carries on a conversation from a prompt plus history.
type ConversationChain struct {
	chatModel       model.ChatModel // only allow using ChatModel
	memory          []model.ChatMessage
	callbackManager *callback.Manager
}

// NewConversationChain create new conversation chain
// firstSystemPrompt will be the first chat in chat memory, with role as System
func NewConversationChain(chatModel model.ChatModel,
	memory []model.ChatMessage,
	callbackManager *callback.Manager,
	firstSystemPrompt string,
	verbose bool) (chain *ConversationChain) {

	if verbose {
		callbackManager.RegisterCallback(basechain.CallbackChainEnd, callback.VerboseCallback)
	}
	memory = append(memory, model.ChatMessage{Role: model.ChatMessageRoleSystem, Content: firstSystemPrompt})
	return &ConversationChain{
		chatModel:       chatModel,
		callbackManager: callbackManager,
		memory:          memory,
	}
}

// AppendMemory to add conversation to the memory
func (C *ConversationChain) AppendToMemory(message model.ChatMessage) {
	C.memory = append(C.memory, message)
}

// Run expect chat["input"] as input, and put the result to output["output"]
func (C *ConversationChain) Run(ctx context.Context, chat map[string]string, options ...func(*model.Option)) (output map[string]string, err error) {
	if _, ok := chat["input"]; !ok {
		return output, errors.New("input[\"input\"] is not specified")
	}
	output = make(map[string]string)

	newContext := callback.NewContext(ctx)

	//trigger callback
	C.callbackManager.TriggerEvent(newContext, basechain.CallbackChainStart, callback.CallbackData{
		EventName:    basechain.CallbackChainStart,
		FunctionName: "ConversationChain.Run",
		Input:        chat,
		Output:       output,
	})

	C.AppendToMemory(model.ChatMessage{Role: model.ChatMessageRoleUser, Content: chat["input"]})
	message, err := C.chatModel.Chat(newContext, C.memory, options...)

	// add response message to memory
	C.AppendToMemory(message)

	output["output"] = message.Content

	//trigger callback
	C.callbackManager.TriggerEvent(newContext, basechain.CallbackChainEnd, callback.CallbackData{
		EventName:    basechain.CallbackChainEnd,
		FunctionName: "ConversationChain.Run",
		Input:        chat,
		Output:       output,
	})

	return
}

// SimpleRun will run the prompt string agains llmchain
func (C *ConversationChain) SimpleRun(ctx context.Context, chat string, options ...func(*model.Option)) (output string, err error) {
	response, err := C.Run(ctx, map[string]string{"input": chat}, options...)
	output = response["output"]

	return
}
