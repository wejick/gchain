# GoChain
[![Go Reference](https://pkg.go.dev/badge/github.com/wejick/gochain.svg)](https://pkg.go.dev/github.com/wejick/gochain)
![Build workflow](https://github.com/wejick/gochain/actions/workflows/go.yml/badge.svg)
[![Integration test](https://github.com/wejick/gochain/actions/workflows/integration.yml/badge.svg)](https://github.com/wejick/gochain/actions/workflows/integration.yml)

Langchain-inspired framework to work with LLM in Golang

## Example
```golang
llmModel = _openai.NewOpenAIModel(authToken, "", "text-davinci-003",callback.NewManager(), true)
chain, err := llm_chain.NewLLMChain(llmModel, nil)
if err != nil {
    //handle error
}
outputMap, err := chain.Run(context.Background(), map[string]string{"input": "Indonesia Capital is Jakarta\nJakarta is the capital of "})
fmt.Println(outputMap["output"])
```
More example in the [example](./example/) folder

As our documentation is not yet complete, please refer to examples and integration test for reference.

## Notice
1. Don't use it if you have better option
1. GoChain priority is golang idiomatic. So eventhough it happily use many langchain concept, don't expect exactly the same behavior as this is not reimplementation.
