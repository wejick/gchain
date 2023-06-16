/*
Gochain provide some abstraction and tools to easily create large language model based application.
It is heavily inspired by langchain, while also providing developer experience more idiomatic to go.

Currently there are 3 concept that we have in Gochain (which is not far different with langchain's) :
 1. Model
    Here is where LLM model get abstracted, currently there are 3 kind of model : llm model, chat model
    and embedding model. Currently those 3 doesn't share the same abstraction, probably at the future
    chat and llm can share the same abstraction.
 2. Datastore
    As the name implies it give access to where the data is stored. It can be a sql database, vector database
    or even an API. We split into 3 abstraction Datastore, vectorstore and retriever.
 3. Chain
    Chain is where model meet a usecase, this is where a problem get solved by using a model. Here where user input,
    datastore and model can meet. There are several ready made chain to use, or as inspiration to build your own chain.
*/
package gochain
