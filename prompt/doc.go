/*
	A “prompt” refers to the input to the model.

This input is rarely hard coded, but rather is often constructed from multiple components.
A PromptTemplate is responsible for the construction of this input.

PromptTemplate is powered by Go text/template, however we provided simplistic interface to interact with.
Everything is a string and data passed along to the template as map[string]string.

Example :

	template, _ := NewPromptTemplate("template_name", "{{.string}} {{.stringfloat}} {{.stringinteger}}")
	Data := map[string]string{"string": "string", "stringfloat": "0.1", "stringinteger": "1"}
	outputPrompt, _ := P.FormatPrompt(Data)

outputPrompt will containt "string 0.1 1"
*/
package prompt
