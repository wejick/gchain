# Eval Package

Eval is a set of functionality to do integration test againts LLM application.

There are several llm based evaluator gchain provides :
1. Correctness evaluation
Correctness evaluation provide a way to test whether the input is valid based on the set of the expectation.

Example :
```go
NewCorrectnessEval(model, "The color of the sky is blue").Evaluate("When I go outside, I see the sky is blue")
>> True
```

2. QA Relevance evalation: 
This evaluation is grading whether the the answer relevant to the question according to the given fact. This can be used to evaluate a conversational application, to see wether the answer for the given answer is relevant.

Example :
```go
fact := "The color of the sky is blue"
question := "What is the color of the sky?"
answer := "When I go outside, I see the sky is blue"

NewQARelevanceEval(model,fact,question).Evaluate(answer)

>> True
```

Other than that, gchain also provide utilitary evaluator :
1. Json validation
