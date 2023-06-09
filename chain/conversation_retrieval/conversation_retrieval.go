package conversationretrieval

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/wejick/gochain/datastore"
	wikipedia "github.com/wejick/gochain/datastore/wikipedia_retrieval"
	"github.com/wejick/gochain/model"
	"github.com/wejick/gochain/prompt"
)

type AnswerOrLookupOutput struct {
	Question string
	Query    string
	Intent   string
	Answer   string
}

// ConversationRetrievalChain conversation with ability to lookup data
type ConversationRetrievalChain struct {
	chatModel           model.ChatModel // only allow using ChatModel
	memory              []model.ChatMessage
	retriever           datastore.Retrieval
	instructionTemplate *prompt.PromptTemplate
	answerTemplate      *prompt.PromptTemplate
}

func NewConversationRetrievalChain(chatModel model.ChatModel, memory []model.ChatMessage, firstSystemPrompt string) (chain *ConversationRetrievalChain) {
	instructionTemplate, _ := prompt.NewPromptTemplate("instruction", instruction)
	answerTemplate, _ := prompt.NewPromptTemplate("answer", answeringInstruction)
	memory = append(memory, model.ChatMessage{Role: model.ChatMessageRoleSystem, Content: firstSystemPrompt})
	return &ConversationRetrievalChain{
		chatModel:           chatModel,
		memory:              memory,
		retriever:           &wikipedia.Wikipedia{},
		instructionTemplate: instructionTemplate,
		answerTemplate:      answerTemplate,
	}
}

// Run expect chat["input"] as input, and put the result to output["output"]
func (C *ConversationRetrievalChain) Run(ctx context.Context, chat map[string]string, options ...func(*model.Option)) (output map[string]string, err error) {
	if _, ok := chat["input"]; !ok {
		return output, errors.New("input[\"input\"] is not specified")
	}
	output = make(map[string]string)

	answerOrLookup, err := C.AnswerOrLookup(ctx, chat["input"], options...)
	if err != nil {
		return
	}

	// when we get direct answer
	if answerOrLookup.Answer != "" {
		C.AppendToMemory(model.ChatMessage{Role: model.ChatMessageRoleUser, Content: chat["input"]})
		C.AppendToMemory(model.ChatMessage{Role: model.ChatMessageRoleAssistant, Content: answerOrLookup.Answer})
		output["output"] = answerOrLookup.Answer
		return
	}

	// when we need to look up
	retrieverOutput, err := C.retriever.Search(ctx, "", answerOrLookup.Query)
	if err != nil {
		return
	}
	var retrieverResult string
	for _, resp := range retrieverOutput {
		if data, ok := resp.(string); ok {
			retrieverResult = data
		}
	}

	answer, err := C.AnswerFromDoc(ctx, answerOrLookup, retrieverResult, options...)

	// append answer to memory
	C.AppendToMemory(model.ChatMessage{Role: model.ChatMessageRoleAssistant, Content: answer})
	if err != nil {
		return
	}

	output["output"] = answer

	return
}

func (C *ConversationRetrievalChain) AnswerOrLookup(ctx context.Context, input string, options ...func(*model.Option)) (output AnswerOrLookupOutput, err error) {
	instructionPrompt, err := C.instructionTemplate.FormatPrompt(map[string]string{"question": input})
	if err != nil {
		return
	}

	// copy memory
	contextMem := C.memory[:]
	// add instruction and user query
	contextMem = append(contextMem, model.ChatMessage{Role: model.ChatMessageRoleUser, Content: instructionPrompt})
	// get the answer
	response, err := C.chatModel.Chat(ctx, contextMem, options...)
	if err != nil {
		return
	}
	errUnmarshall := json.Unmarshal([]byte(response.Content), &output)
	if errUnmarshall == nil {
		return
	}

	output.Answer = response.Content

	return
}

func (C *ConversationRetrievalChain) AnswerFromDoc(ctx context.Context, context AnswerOrLookupOutput, doc string, options ...func(*model.Option)) (output string, err error) {
	// FIXME cut the doc to 300
	if len(doc) > 300 {
		doc = doc[0:300]
	}

	b, err := json.Marshal(context)
	if err != nil {
		log.Fatal(err)
	}

	instructionPrompt, err := C.answerTemplate.FormatPrompt(map[string]string{"doc": doc, "context": string(b)})
	if err != nil {
		log.Fatal(err)
	}
	message := model.ChatMessage{
		Role:    model.ChatMessageRoleUser,
		Content: instructionPrompt,
	}
	AIResponse, err := C.chatModel.Chat(ctx, []model.ChatMessage{message}, options...)
	if err != nil {
		return
	}

	output = AIResponse.Content

	return
}

// AppendMemory to add conversation to the memory
func (C *ConversationRetrievalChain) AppendToMemory(message model.ChatMessage) {
	C.memory = append(C.memory, message)
}
