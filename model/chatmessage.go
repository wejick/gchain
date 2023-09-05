package model

type ChatMessage struct {
	Role          string
	Content       string
	Name          string `json:"name"`
	ParameterJson string // json string of function parameter
	PromptUsage   PromptUsage  `json:"usage"`
}

type PromptUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func (C *ChatMessage) String() string {
	return C.Role + ": " + C.Content
}

const (
	ChatMessageRoleSystem    = "system"
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
	ChatMessageRoleFunction  = "function"
)
