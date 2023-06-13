package model

import "testing"

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
