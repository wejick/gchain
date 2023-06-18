package callback

import (
	"context"
	"reflect"

	"github.com/k0kubun/pp"
)

type Event string

type Callback func(context context.Context, data CallbackData)
type CallbackIdentifier struct {
	EventName   Event
	FuncPointer uintptr
}

type CallbackData struct {
	RunID        string // to be populated with data from context
	EventName    string
	FunctionName string
	Input        map[string]string
	Output       map[string]string
	Data         interface{}
}

type Manager struct {
	callbacks           map[Event][]Callback
	callbackIdentifiers []CallbackIdentifier
}

func NewManager() *Manager {
	return &Manager{
		callbacks: make(map[Event][]Callback),
	}
}

func (m *Manager) RegisterCallback(event Event, callback Callback) {
	identifier := CallbackIdentifier{
		EventName:   event,
		FuncPointer: reflect.ValueOf(callback).Pointer(),
	}

	// Check if the identifier already exists in the map
	for _, id := range m.callbackIdentifiers {
		if id == identifier {
			// The callback already exists, so we skip adding it again
			return
		}
	}

	// Add the identifier and callback to the maps
	m.callbackIdentifiers = append(m.callbackIdentifiers, identifier)
	m.callbacks[event] = append(m.callbacks[event], callback)
}

func (m *Manager) TriggerEvent(ctx context.Context, event Event, data CallbackData) {
	if callbacks, ok := m.callbacks[event]; ok {
		for _, callback := range callbacks {
			callback(ctx, data)
		}
	}
}

// inspect return all the registered callback so we can test
func (m *Manager) inspect() map[Event][]Callback {
	return m.callbacks
}

func VerboseCallback(ctx context.Context, data CallbackData) {
	pp.Println(data)
}
