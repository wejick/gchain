package model

import (
	"testing"

	goopenai "github.com/sashabaranov/go-openai"
)

func TestFlattenChatMessages(t *testing.T) {
	messages := []ChatMessage{
		{Role: "User", Content: "Hello"},
		{Role: "Bot", Content: "Hi there!"},
		{Role: "User", Content: "How are you?"},
	}

	expected := "User: Hello\nBot: Hi there!\nUser: How are you?\n"
	result := FlattenChatMessages(messages)

	if result != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result)
	}
}

func TestIsStreamFinished(t *testing.T) {
	const signalRole = "signal"
	const signalContentFinished = "finished"

	// A message that should signal the end of the stream
	endMessage := ChatMessage{
		Role:    signalRole,
		Content: signalContentFinished,
	}
	if !IsStreamFinished(endMessage) {
		t.Error("Expected stream to be finished, but it was not")
	}

	// A normal message that should not signal the end of the stream
	normalMessage := ChatMessage{
		Role:    goopenai.ChatMessageRoleAssistant,
		Content: "hello",
	}
	if IsStreamFinished(normalMessage) {
		t.Error("Expected stream to not be finished, but it was")
	}
}
