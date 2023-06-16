package model

import "strings"

func FlattenChatMessages(messages []ChatMessage) string {
	var result strings.Builder

	for _, message := range messages {
		result.WriteString(message.String())
		result.WriteString("\n")
	}

	return result.String()
}

const signalRole = "signal"
const signalContentFinished = "finished"

// IsStreamFinished check if it's the end of message
func IsStreamFinished(message ChatMessage) bool {
	if message.Role == signalRole && message.Content == signalContentFinished {
		return true
	}
	return false
}
