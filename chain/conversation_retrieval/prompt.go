package conversationretrieval

var instruction = `INSTRUCTION
You have access to this knowledge base (KB) :
name | description
wikipedia | knowledge base to find any general information, I accept standalone keywords as query.

CONVERSATION
{{.history}}
User: {{.question}}

When responding to me, please output a response in one of two formats:

**Option 1:**
use this if you can answer directly without KB lookup
Answer in this following json schema:

{
"conversation_context":"longer additional context to understand user question",
"intent":"user intention",
"answer":"the answer here",
"lookup":false
}

**Option 2:**
Use this if you want to lookup
Answer in this following json schema:

{
"kb":"knowledge base name",
"question":"user question",
"query":"the question keyword will be put here",
"conversation_context":"longer additional context to understand user question",
"intent":"user intention",
"lookup":true
}
`

var answeringInstruction = `with this context, answer user question in kind and concise manner
{{.context}}
document: {{.doc}}
Answer concisely :
`
