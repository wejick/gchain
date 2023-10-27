package wikipedia

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/wejick/gchain/model"
	"github.com/wejick/gchain/tools"

	gowiki "github.com/trietmn/go-wiki"
)

type WikipediaSearchTool struct {
	functionDefinition model.FunctionDefinition
}

func NewWikipediaSearchTool() *WikipediaSearchTool {
	return &WikipediaSearchTool{
		functionDefinition: model.FunctionDefinition{
			Name:        "wikipedia_search_tool",
			Description: "This tool is used to search to wikipedia.com for keyword, wikipedia is vast source of knowledge",
			Parameters: model.FunctionJsonSchema{
				Type: model.FunctionDataTypeObject,
				Properties: map[string]model.FunctionJsonSchema{
					"operation": {
						Type:        model.FunctionDataTypeString,
						Description: "specify whether to search or open wikipedia page. Possible value strickly OPEN, SEARCH",
					},
					"keyword": {
						Type:        model.FunctionDataTypeString,
						Description: "keyword to search or open",
					},
				},
				Required: []string{"operation", "keyword"},
			},
		},
	}
}

// Run give greeting to user, this is to demonstrate the simples form of tool
// Run expect as map with "user_name" key
func (WS *WikipediaSearchTool) Run(ctx context.Context, input map[string]string, options ...func(*model.Option)) (output map[string]string, err error) {
	if input == nil {
		return nil, errors.New("WikipediaSearchTool : Empty Input")
	}

	if input["operation"] == "SEARCH" {
		titles := WS.search(input["keyword"])
		output = map[string]string{"output": strings.Join(titles, ",")}
	} else if input["operation"] == "OPEN" {
		content := WS.open(input["keyword"])
		output = map[string]string{"output": content}
	} else {
		output = map[string]string{"output": "Supported operation are only SEARCH and OPEN"}
		err = errors.New("WikipediaSearchTool : Invalid operation")
	}

	return
}

// SimpleRun give greeting to user, this is to demonstrate the simples form of tool
// SimpleRun expect valid json string with "keyword" field
func (WS *WikipediaSearchTool) SimpleRun(ctx context.Context, prompt string, options ...func(*model.Option)) (output string, err error) {
	var parameter map[string]string
	err = json.Unmarshal([]byte(prompt), &parameter)
	if err != nil {
		return
	}
	runOutput, err := WS.Run(ctx, parameter, options...)
	if err != nil || runOutput == nil {
		return runOutput["output"], err
	}

	output = runOutput["output"]

	return
}

// search a keyword in wikipedia, return list of wikipedia pages
func (WS *WikipediaSearchTool) search(keyword string) (titles []string) {
	titles, _, err := gowiki.Search(keyword, 15, false)
	if err != nil {
		return
	}

	return titles
}

// open a page in wikipedia
func (WS *WikipediaSearchTool) open(title string) (pageContent string) {
	page, err := gowiki.GetPage(title, -1, true, true)
	if err != nil {
		return
	}
	pageContent, err = page.GetContent()
	if err != nil {
		return
	}

	return
}

// GetFunctionDefinition return function definition of the tool
func (WS *WikipediaSearchTool) GetFunctionDefinition() model.FunctionDefinition {
	return WS.functionDefinition
}

// GetDefinitionString tool definition in string format
func (WS *WikipediaSearchTool) GetDefinitionString() string {
	description := tools.GetDefinitionString(WS)

	return description
}
