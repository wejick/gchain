package prompt

import (
	"bytes"
	"text/template"
)

type PromptTemplate struct {
	tmplt *template.Template
}

// NewPromptTemplate create a new prompt template
func NewPromptTemplate(name string, templateString string) (output_template *PromptTemplate, err error) {
	tplt, err := template.New(name).Parse(templateString)
	output_template = &PromptTemplate{
		tmplt: tplt,
	}
	return
}

// FormatPrompt to generate prompt from template
func (P *PromptTemplate) FormatPrompt(Data map[string]string) (output_prompt string, err error) {
	var buf bytes.Buffer
	err = P.tmplt.Execute(&buf, Data)
	output_prompt = buf.String()

	return
}
