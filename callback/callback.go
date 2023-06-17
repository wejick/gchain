package callback

import (
	"context"

	"github.com/k0kubun/pp"
)

type Event string

type Callback func(context context.Context, data CallbackData)

type CallbackData struct {
	RunID        string // to be populated with data from context
	EventName    string
	FunctionName string
	Input        map[string]string
	Output       map[string]string
}

type Manager struct {
	callbacks map[Event][]Callback
}

func NewManager() *Manager {
	return &Manager{
		callbacks: make(map[Event][]Callback),
	}
}

func (m *Manager) RegisterCallback(event Event, callback Callback) {
	m.callbacks[event] = append(m.callbacks[event], callback)
}

func (m *Manager) TriggerEvent(ctx context.Context, event Event, data CallbackData) {
	if callbacks, ok := m.callbacks[event]; ok {
		for _, callback := range callbacks {
			callback(ctx, data)
		}
	}
}

func VerboseCallback(ctx context.Context, data CallbackData) {
	pp.Println(data)
}
