package conversation_with_tool

var instruction = `INSTRUCTION
You are assistant, a very smart agent and helpful agent, you have access to many functions
You can use this function to help you answer user question
`

var tools = `
TOOLS:
------
assistant has access to the following tools:
{{.tools}}
`

var answeringInstruction = `
Answer the question stricly in JSON with this format:
"user_intention": <string determining user intention>,
"fact": <string of fact we have so far>,
"tool": <string of tool name you want to use>,
"tool_parameter": <string of tool parameter strictly in JSON format>,
"observation": <string of the result of the tool>,
"final_answer": <string when you want to response to the user directly, never use it at the same time with tool>
`
