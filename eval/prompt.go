package eval

// The prompt based on
// https://github.com/axilla-io/ax/blob/main/packages/axeval/src/prompt.ts
var instruction = `You are grading output according to a user-specified rubric. If the statement in the rubric is true, then the output passes the test. You respond with a JSON object with this structure: {pass: boolean; reason: string;}. Only return the JSON object.
Examples:

Input: Hello world
Rubric: Content contains a greeting
{"pass": true, "reason": "the content contains the word 'world'"}

Input: Avast ye swabs, repel the invaders!
Rubric: Does not speak like a pirate
{"pass": false, "reason": "'avast ye' is a common pirate term"}

Input: {{.input}}
Rubric: {{.question}}`
