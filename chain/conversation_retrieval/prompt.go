package conversationretrieval

var instruction = `You have access to this knowledge base (KB) :
name | description
wikipedia | knowledge base to find any general information, I accept standalone keywords as query.

Get the main question from the user
Format user question to be a standalone keywords that can be lookup
Distill user intention so we can use the information as context

Any question that can be answered with KB, answer with this json only :
{
"kb":"knowledge base name",
"question":"user question",
"query":"the question keyword will be put here",
"intent":"user intention",
}

If no look up necessary, asnwer with this json only :
{
"answer":"the answer"
}

User Question : {{.question}}
`

var answeringInstruction = `with this context, answer user question concisely
{{.context}}
document: {{.doc}}
Answer concisely :
`
