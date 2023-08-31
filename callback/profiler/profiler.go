package profiler

import (
	"context"
	"sync"

	"github.com/wejick/gchain/callback"
)

type Profiler struct {
	Events map[string][]*callback.CallbackData // map of callbackData by sessionID
	mutext sync.RWMutex
}

func NewProfiler() *Profiler {
	return &Profiler{
		Events: make(map[string][]*callback.CallbackData),
	}
}

func (P *Profiler) Callback(context context.Context, data callback.CallbackData) {
	P.mutext.Lock()
	defer P.mutext.Unlock()

	P.Events[data.SessionID] = append(P.Events[data.SessionID], &data)
}

func (P *Profiler) GetEvents(sessionID string) []*callback.CallbackData {
	P.mutext.RLock()
	defer P.mutext.RUnlock()

	return P.Events[sessionID]
}
