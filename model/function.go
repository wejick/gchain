package model

import "strings"

type DataType string

const (
	FunctionDataTypeString DataType = "string"
	FunctionDataTypeObject DataType = "object"
)

type FunctionJsonSchema struct {
	Type        DataType                      `json:"type,omitempty"`
	Properties  map[string]FunctionJsonSchema `json:"properties,omitempty"`
	Required    []string                      `json:"required,omitempty"`
	Description string                        `json:"description,omitempty"`
	Enum        []string                      `json:"enum,omitempty"`
}

// FunctionDefinition is to describe function to model
// currently only being supported by few openAI's Model
type FunctionDefinition struct {
	Name        string             `json:"name,omitempty"`
	Type        string             `json:"type,omitempty"`
	Description string             `json:"description,omitempty"`
	Parameters  FunctionJsonSchema `json:"parameters,omitempty"`
}

func (F FunctionJsonSchema) String() string {
	var paramString string
	if len(F.Required) > 0 {
		paramString += "required = " + strings.Join(F.Required, ",")
	}
	paramString += "\nparameter|description|type|enum"
	for key, value := range F.Properties {
		var enumString string
		if len(value.Enum) > 0 {
			enumString = strings.Join(value.Enum, ",")
		} else {
			enumString = "no"
		}
		paramString += "\n" + key + "|" + value.Description + "|" + string(value.Type) + "|" + enumString
	}
	return paramString
}
