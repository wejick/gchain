package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/wejick/gchain/callback"
	"github.com/wejick/gchain/model"
	_openai "github.com/wejick/gchain/model/openAI"
)

func main() {
	var authToken = os.Getenv("OPENAI_API_KEY")
	var chatModel model.ChatModel
	chatModel = _openai.NewOpenAIChatModel(authToken, "", _openai.GPT3Dot5Turbo0301, callback.NewManager(), false)
	memory := []model.ChatMessage{}

	// prepare a function register
	functionList := map[string]func(map[string]string) string{
		"get_longname": func(parameter map[string]string) string {
			return getLongName(parameter["user_name"])
		},
	}

	// The first call to the model, to see whether function call is needed
	memory = append(memory, model.ChatMessage{Role: model.ChatMessageRoleUser, Content: "Hi I'm Gio, what's my long name?"})
	functionDef := model.FunctionDefinition{
		Name:        "get_longname",
		Description: "When user need to get the long name of user",
		Parameters: model.FunctionJsonSchema{
			Type: model.FunctionDataTypeObject,
			Properties: map[string]model.FunctionJsonSchema{
				"user_name": {
					Type:        model.FunctionDataTypeString,
					Description: "User name",
				},
			},
			Required: []string{"user_name"},
		},
	}
	response, err := chatModel.Chat(context.Background(), memory, model.WithFunctions([]model.FunctionDefinition{functionDef}))
	if err != nil {
		log.Println(err)
	}
	fmt.Println(response)

	// append the first response to memory
	memory = append(memory, model.ChatMessage{Role: model.ChatMessageRoleAssistant, Content: response.Content, FunctionName: response.FunctionName})

	// check if function call needed
	if response.FunctionName != "" {
		var parameter map[string]string
		err := json.Unmarshal([]byte(response.ParameterJson), &parameter)
		if err != nil {
			log.Println(err)
		}
		// call the function and get the result
		functionCallReturn := functionList[response.FunctionName](parameter)

		// The second call to the model, to give the function result to the model
		memory = append(memory, model.ChatMessage{Role: model.ChatMessageRoleFunction, FunctionName: response.FunctionName, Content: functionCallReturn})
		response, err = chatModel.Chat(context.Background(), memory, model.WithFunctions([]model.FunctionDefinition{functionDef}))
		if err != nil {
			log.Println(err)
		}
		fmt.Println(response.Content)
	}
}

func getLongName(username string) string {
	log.Println("Revealing the long name for", username)
	return fmt.Sprintf("Hi %s! your long name is %s SuperONE", username, username)
}
