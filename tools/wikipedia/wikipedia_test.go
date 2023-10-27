//go:build integration
// +build integration

package wikipedia

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWikipediaTool(t *testing.T) {
	wikiSearchTool := NewWikipediaSearchTool()
	assert.NotNil(t, wikiSearchTool)

	assert.NotNil(t, wikiSearchTool.GetFunctionDefinition())

	// test simple run search
	output, err := wikiSearchTool.SimpleRun(context.Background(), `{"operation":"SEARCH","keyword":"Formula One"}`)
	assert.Nil(t, err)
	assert.Contains(t, output, "Formula One")

	// test simple run with invalid input
	output, err = wikiSearchTool.SimpleRun(context.Background(), `{"operation":"search","keyword":"Formula One"}`)
	assert.NotNil(t, err)
	assert.Equal(t, "Supported operation are only SEARCH and OPEN", output)

	// test simple run open
	output, err = wikiSearchTool.SimpleRun(context.Background(), `{"operation":"OPEN","keyword":"Formula One"}`)
	assert.Nil(t, err)
	assert.Contains(t, output, "Formula One")

	// get description
	description := `name = ` + wikiSearchTool.functionDefinition.Name + `
description = ` + wikiSearchTool.functionDefinition.Description + "\n"
	assert.Equal(t, description+wikiSearchTool.GetFunctionDefinition().Parameters.String(), wikiSearchTool.GetDefinitionString())
}
