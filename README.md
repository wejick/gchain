# GoChain
Langchain inspired framework to work with LLM in golang

## Example
```golang
llmModel = _openai.NewOpenAIModel(authToken, "", "text-davinci-003")
chain, err := llm_chain.NewLLMChain(llmModel, nil)
if err != nil {
    //handle error
}
outputMap, err := chain.Run(context.Background(), map[string]string{"input": "Indonesia Capital is Jakarta\nJakarta is the capital of "})
fmt.Println(outputMap["output"])
```
More example in the [example](./example/) folder

## Notice
1. Don't use it if you have better option
1. GoChain priority is golang idiomatic. So eventhough it happily use many langchain concept, don't expect exactly the same behavior as this is not reimplementation.
