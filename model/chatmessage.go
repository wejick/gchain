package model

type ChatMessage struct {
	Role          string
	Content       string
	Name          string `json:"name"`
	ParameterJson string // json string of function parameter
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
