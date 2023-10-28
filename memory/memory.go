package memory

import "github.com/wejick/gchain/model"

// Memory retain all conversation messages
type Memory interface {
	// AddMessage adds a message to the memory
	AddMessage(message model.ChatMessage)
	// Get all messages
	GetAll() []model.ChatMessage
	// String flatten all messages to a single string
	String() string
}
