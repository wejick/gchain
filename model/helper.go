package model

import "strings"

func FlattenChatMessages(messages []ChatMessage) string {
	var result strings.Builder

	for _, message := range messages {
		result.WriteString(message.Role)
		result.WriteString(": ")
		result.WriteString(message.Content)
		result.WriteString("\n")
	}

	return result.String()
}
