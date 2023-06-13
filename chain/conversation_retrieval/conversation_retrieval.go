package conversationretrieval

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/wejick/gochain/datastore"
	"github.com/wejick/gochain/model"
	"github.com/wejick/gochain/prompt"
	"github.com/wejick/gochain/textsplitter"
)

type AnswerOrLookupOutput struct {
	Question            string `json:"question"`
	Query               string `json:"query"`
	Intent              string `json:"intent"`
	Answer              string `json:"answer"`
	ConversationContext string `json:"conversation_context"`
}

// ConversationRetrievalChain conversation with ability to lookup data
type ConversationRetrievalChain struct {
	chatModel           model.ChatModel // only allow using ChatModel
	memory              []model.ChatMessage
	retriever           datastore.Retrieval
	textSplitter        textsplitter.TextSplitter
	instructionTemplate *prompt.PromptTemplate
	answerTemplate      *prompt.PromptTemplate
	maxToken            int
}

// FIXME : put some parameter as options
func NewConversationRetrievalChain(chatModel model.ChatModel, memory []model.ChatMessage, retriever datastore.Retrieval, textSplitter textsplitter.TextSplitter, firstSystemPrompt string, maxToken int) (chain *ConversationRetrievalChain) {
	instructionTemplate, _ := prompt.NewPromptTemplate("instruction", instruction)
	answerTemplate, _ := prompt.NewPromptTemplate("answer", answeringInstruction)
	memory = append(memory, model.ChatMessage{Role: model.ChatMessageRoleSystem, Content: firstSystemPrompt})
	if maxToken == 0 {
		maxToken = 1000
	}
	return &ConversationRetrievalChain{
		chatModel:           chatModel,
		memory:              memory,
		retriever:           retriever,
		instructionTemplate: instructionTemplate,
		answerTemplate:      answerTemplate,
		maxToken:            maxToken,
		textSplitter:        textSplitter,
	}
}

// Run expect chat["input"] as input, and put the result to output["output"]
func (C *ConversationRetrievalChain) Run(ctx context.Context, chat map[string]string, options ...func(*model.Option)) (output map[string]string, err error) {
	if _, ok := chat["input"]; !ok {
		return output, errors.New("input[\"input\"] is not specified")
	}
	output = make(map[string]string)

	inputChat := model.ChatMessage{Role: model.ChatMessageRoleUser, Content: chat["input"]}

	answerOrLookup, err := C.AnswerOrLookup(ctx, chat["input"], options...)
	if err != nil {
		return
	}

	// when we get direct answer
	if answerOrLookup.Answer != "" {
		C.AppendToMemory(inputChat)
		C.AppendToMemory(model.ChatMessage{Role: model.ChatMessageRoleAssistant, Content: answerOrLookup.Answer})
		output["output"] = answerOrLookup.Answer
		return
	}

	// when we need to look up
	retrieverOutput, err := C.retriever.Search(ctx, "Indonesia", answerOrLookup.Query)
	if err != nil {
		return
	}
	var retrieverResult string
	for _, resp := range retrieverOutput {
		if data, ok := resp.(map[string]interface{}); ok {
			retrieverResult += data["text"].(string)
		}
	}

	answer, err := C.answerFromDoc(ctx, answerOrLookup, retrieverResult, options...)

	// append answer to memory
	C.AppendToMemory(inputChat)
	C.AppendToMemory(model.ChatMessage{Role: model.ChatMessageRoleAssistant, Content: answer})
	if err != nil {
		return
	}

	output["output"] = answer

	return
}

// AnswerOrLookup will return answer if it can, or return lookup query
// This approach is faster because we will be able to get answer directly when possible
func (C *ConversationRetrievalChain) AnswerOrLookup(ctx context.Context, input string, options ...func(*model.Option)) (output AnswerOrLookupOutput, err error) {
	convoHistory := model.FlattenChatMessages(C.memory)
	instructionPrompt, err := C.instructionTemplate.FormatPrompt(map[string]string{"question": input, "history": convoHistory})
	if err != nil {
		return
	}

	response, err := C.chatModel.Chat(ctx, []model.ChatMessage{{Role: model.ChatMessageRoleUser, Content: instructionPrompt}}, options...)
	if err != nil {
		return
	}
	errUnmarshall := json.Unmarshal([]byte(response.Content), &output)
	if errUnmarshall != nil {
		return
	}

	return
}

// answerFromDoc based on the given context, will retrieve data and use it to answer question using llm
// this one off query is the key to make this more cost effective and save token usage
func (C *ConversationRetrievalChain) answerFromDoc(ctx context.Context, context AnswerOrLookupOutput, doc string, options ...func(*model.Option)) (output string, err error) {
	// cut to max token
	if len(doc) > C.maxToken {
		doc = C.textSplitter.SplitText(doc, C.maxToken, 0)[0]
	}

	b, err := json.Marshal(context)
	if err != nil {
		log.Println(err)
		return
	}

	instructionPrompt, err := C.answerTemplate.FormatPrompt(map[string]string{"doc": doc, "context": string(b)})
	if err != nil {
		log.Print(err)
		return
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
