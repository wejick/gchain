package model

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
