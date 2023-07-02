package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionJsonSchema_String(t *testing.T) {
	parameter := FunctionJsonSchema{}
	parameter.Required = []string{"parameter1"}
	parameter.Properties = map[string]FunctionJsonSchema{
		"parameter1": {
			Type:        FunctionDataTypeString,
			Description: "parameter 1 description",
			Enum:        []string{"enum1", "enum2"},
		},
		"parameter2": {
			Type:        FunctionDataTypeString,
			Description: "parameter 2 description",
		},
	}

	paramString := parameter.String()
	expectedString := "required = parameter1\nparameter|description|type|enum" + "\n"
	expectedString += "parameter2|parameter 2 description|string|no" + "\n"
	expectedString += "parameter1|parameter 1 description|string|enum1,enum2"
	assert.Equal(t, expectedString, paramString)
}
