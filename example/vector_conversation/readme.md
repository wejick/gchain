# Example of using conversation retrieval
There are 2 part of this example, the indexer and chat interface.

1. On the indexer, we index text of Indonesia and History of Indonesia wikipedia page. Put it to weaviate vector db.
2. Using data in weaviate, we accept user queries.

Run indexing
```sh
$go build
$vector_conversation --index
```

Run chat
```sh
$go build
$vector_conversation
AI : How can I help you, I know many things about indonesia
User : who is the first president of indonesia?
AI : The first president of Indonesia was Sukarno.
```