package callback

import (
	"context"
	"reflect"

	"github.com/google/uuid"
	"github.com/k0kubun/pp"
)

type contextKey int

const (
	contextKeySessionID contextKey = iota
	contextKeyID
	contextKeyParentID
)

type Event string

type Callback func(context context.Context, data CallbackData)
type CallbackIdentifier struct {
	EventName   Event
	FuncPointer uintptr
}

type CallbackData struct {
	SessionID    string // to be populated with data from context
	ID           string
	ParentID     string
	EventName    string
	FunctionName string
	Input        map[string]string
	Output       map[string]string
	Data         interface{}
	UnixTime     int64
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
	// get session id from context
	if sessionID, ok := ctx.Value(contextKeyParentID).(string); ok {
		data.SessionID = sessionID
	}
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

// NewContext create new context with new generated ID and update the parentID
func NewContext(ctx context.Context) context.Context {
	// generate new UUID
	id := uuid.New()

	// get parentID from context
	// previousID will be parentID
	var parentID string
	if parentIDFromContext, ok := ctx.Value("ID").(string); ok {
		parentID = parentIDFromContext
	}

	// create new context with new ID and parentID
	newContext := context.WithValue(ctx, contextKeyID, id.String())
	newContext = context.WithValue(newContext, contextKeyParentID, parentID)

	return newContext
}

func VerboseCallback(ctx context.Context, data CallbackData) {
	pp.Println(data)
}
